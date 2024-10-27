// document storage, retrieval, and query library based on level-db (lib/store)
// storage also chunks and creates embeddings for each document with some level
// of scoping, which is usually by application id

const net = require('net');
const fsp = require('fs/promises');
const util = require('../lib/util');
const state = require('./service').init();
const fetchp = import('node-fetch');
const cli_server = require('../cli/store');
const log = util.logpre('doc');
const { args, env } = util;
const { debug } = args;
const { PromiseQ } = require('../lib/pqueue');
const drop_host = args["drop-host"];
const drop_port = args["drop-port"] || 0;
const cli_store = args["cli-store"];

const debug_embed = debug === 'embed';

async function setup_node() {
  const { node } = state;
  // re-connect the doc app when the proxy connection bounces
  node.on_reconnect(register_service);
}

// respond to rewquests looking for docs servicing an app
async function register_service() {
  const { app_id, net_addrs, node } = state;
  // announce presence
  node.publish("service-up", {
    app_id,
    type: "doc-server",
    subtype: "rawh-level-v0",
    net_addrs
  });
  // bind api service endpoints for document operations
  node.handle([ "doc-load", app_id ], doc_load);
  node.handle([ "doc-list", app_id ], doc_list);
  node.handle([ "doc-delete", app_id ], doc_delete);
  node.handle([ "docs-query", app_id ], docs_query);
}

// request to a tcp socket for bulk loading a document
// which will then be stored, chunked, and vector embedded
async function doc_load(msg = {}) {
  log({ doc_load: msg });
  const { url, name, type } = msg;
  const { node, app_id } = state;
  // create the file drop target with time+random file uid
  const fuid = util.uid();
  const fdir = `${state.data_dir}/docs`;
  const path = `${fdir}/${fuid}`;
  await fsp.mkdir(fdir, { recursive: true }).catch(e => e);
  const frec = { 
    id: fuid, 
    url,  // optional if url was the source
    name, // file name without path (for display)
    type, // pdf, csv, html, text, etc.
    path, // location of artifact on disk
    length: 0, // artifact length in bytes,
    state: "loading",
    added: Date.now() 
  };
  const fdel = async function () {
    // delete partial file data
    await fsp.rm(path).catch(error => log({ bulk_delete_error: error }));
  }
  node.publish([ 'doc-loading', app_id ], frec);
  if (url) {
    return doc_load_url(frec, fdel);
  } else {
    return doc_load_socket(frec, fdel);
  }
}

// load a document stream from url
async function doc_load_url(frec, fdel) {
  const fetch = (await fetchp).default;
  const data = await fetch(frec.url, { redirect: 'follow', follow: 20 });
  const bufr = Buffer.from(await data.arrayBuffer());
  await fsp.writeFile(frec.path, bufr);
  // do file analysis
  await doc_embed(frec).catch(error => {
        log({ embed_error: error, frec });
        fdel();
  });
  return { loaded: frec.url };
}

// accept an inbound document stream over tcp
async function doc_load_socket(frec, fdel) {
  const { net_addrs } = state;
  const file = await fsp.open(frec.path, 'w');
  const fstr = await file.createWriteStream();
  // listen on random tcp port for incoming file dump
  const srv = net.createServer(conn => {
    // log({ bulk_conn: conn });
    conn.on('error', (error) => {
         log({ bulk_error: error });
         fdel();
         srv.close();
    });
    conn.on('data', (data) => {
      // log({ bulk_data: data });
      fstr.write(data);
      frec.length += (data.length || data.byteLength);
    })
    conn.on('end', async () => {
      // log({ bulk_end: name, type });
      await file.datasync();
      await fstr.end();
      await file.close();
      // do file analysis
      await doc_embed(frec).catch(error => {
            log({ embed_error: error, frec });
      });
      srv.close();
    });
  })
  .on('error', error => {
      log({ bulk_listen_error: error });
      fdel();
      srv.close();
  });
  // once tcp server is up, send port back to caller
  return await new Promise(reply => {
    srv.listen({ port: drop_port }, () => {
      // log({ bulk_load: srv.address() });
      // send addr to requestor to complete bulk load
      reply({ 
        host: drop_host ? [ drop_host ] : net_addrs,
        port: srv.address().port
      });
    });
  });
}

// split document into bite sized chunks
async function doc_chunk(frec) {
  const { node, token, app_id, doc_meta, doc_chnk, llama_token } = state;

  // store and publish meta-data about doc
  frec.state = "tokenizing";
  await doc_meta.put(frec.uid, frec);
  node.publish([ 'doc-loading', app_id ], frec);

  // tokenize and embed
  const chunks = await token.load_path(frec.path, {
    type: frec.type,
    chunkSize: 350,
    chunksSizeMax: 700,
    paraChunks: true, // new paragraph splitter
  });

  // extract usable data from langchain nested structure
  const sub = doc_chnk.sub(frec.uid);
  const batch = await sub.batch();
  const map = chunks.map((chunk, idx) => {
      const { pageContent, metadata } = chunk;
      const { loc } = metadata;
      const { pageNumber, lines } = loc;
      const { from, to } = lines;
      if (pageContent.length === 0) {
          // eliminate empty chunks
          return undefined;
      }
      const tokens = llama_token.encode(pageContent).length;
      const rec = [ pageNumber, from, to, tokens, pageContent ];
      batch.put(idx.toString().padStart(4,0), rec);
      return rec;
  }).filter(r => r); // filter empties

  // persist to level storage
  await batch.write();

  return map;
}

// create embeddings from one or more document chunks
async function doc_embed(frec) {
  const { node, embed, app_id, doc_meta, doc_mbed, lama_token } = state;
  const start = Date.now();

  log({ doc_embed: frec });

  // create text chunks from raw or paged document
  const chunks = await doc_chunk(frec);

  // store and publish meta-data about doc
  frec.state = "embedding";
  await doc_meta.put(frec.uid, frec);
  node.publish([ 'doc-loading', app_id ], frec);

  // annotate chunks with their vector and db index (also used for cosine sim)
  const low = 200;
  const high = 600;
  const acc = { cnks: [], toks: [], text: [], tokens: 0, count: 0, new: 0 };
  const batch = await doc_mbed.batch();
  const pq = new PromiseQ({ concurrent: 1 });

  async function (emit) {
      if (acc.new === 0) {
          return emit(acc, 0); // not sure if this is right
      }

      const cnks = acc.cnks.slice();
      const from = cnks[0];
      const to = cnks(cnks.length -1);
      const tokens = acc.toks.slice();
      const tokcnt = acc.tokens;
      const txtlen = acc.text.length;
      const pro = node.promise.call('', "embed/org", { text: acc.text.join("\n") }).then(vector => {
        const index = embed.vector_to_index(vector);
        const key = `${index.toString().padEnd(18, 0)}:${frec.uid}`;
      if (debug_embed) {
          console.log(`---[ ${tokcnt} ]---[ ${from}:${to} ]---[ ${cnks.length}:${tokens.length}:${txtlen} ]---`);
      }
      return batch.put(key, {
          from: acc.cnks[0],
          to: acc.cnks[acc.cnks.length -1],
          tokens: acc.toks,
          index,
          vector,
      });
    });
    acc.count++;
    acc.new = 0;
    return await pq.add(pro);
  }

  // loop over chunks and create rolling embeds with overlapping regions
  for (let i = 0; i < chunks.length; i++) {
       const tokens = chunks[i][3];
       // roll in
       if (acc.tokens + tokens > high) {
           await emit();
           // roll out
           while (acc.cnks.length && acc.tokens > low) {
              acc.cnks.shift();
              acc.tokens -= acc.toks.shift();
              acc.text.shift();
           }
       }
       acc.tokens += tokens;
       acc.cnks.push(i);
       acc.cnks.push(tokens);
       acc.text.path(chunks[i][4]);
       acc.new++;
  }
  await emit();
  await pq.drain();
  await batch.write();

  // store and publish meta-data about doc
  frec.elapsed = Date.now() - start;
  frec.chunks = chunks.length;
  frec.embeds = acc.count;
  frec.state = "ready";
  await doc_meta.put(frec.uid, frec);
  node.publish([ 'doc-loading', app_id ], frec);
  log({ doc_loaded: frec });
}

// list all docs along with the status (loading, embedding, ready)
async function doc_list(msg) {
  return state.doc_meta.list();
}

// delete a document and all of its associated embeddings
async function doc_delete(msg) {
  const { node, app_id, doc_meta, doc_chnk, doc_mbed } = state;
  const { uid } = msg;
  // delete matching meta-data and embeds
  const rec = await doc_chnk.get(uid);
  if (!rec) {
      log({ doc_delete_missing: uid });
      return `delete: invalid document ${uid}`;
  }
  const batch = await doc_mbed.batch();
  let recs = 0, match = 0;
  for await (const [key] of doc_mbed.iter({ values: false })) {
    const [ index, doc_uid ] = key.split(":");
    if (doc_uid === uid) {
      batch.del(key);
      match++;
    }
    recs++;
  }
  await batch.write();
  await doc_meta.del(uid);
  await doc_chunk.sub(rec.uid).clear();
  // delete file artifact
  const fdir = `${state.data_dir}/docs`;
  const path = `${fdir}/${uid}`;
  await fsp.unlink(path);
  // log({ del_doc_info: rec, chunks: match });
  node.publish([ "doc-delete", app_id ], rec);
  return `analyzed ${recs} recs, ${match} matched`;
}

// given a query, get matching embed chunks from loaded docs
async function docs_query(msg) {
  const { node, embed, doc_mbed, doc_chnk } = state;
  const { query, max_tokens, min_match, llm, topic } = msg;

  const vector = await (await node.promise.call('', "embed/org", { text: query }));
  const index = embed.vector_to_index(vector);
  const key = `${index.toString().padEnd(18, 0)}`;
  log({ docs_query: msg });

  const iter_lt = {
    iter: doc_mbed.iter({ lt: key }),
    coss: 1,
    pos: -1,
    add: -1
  };
  const iter_gt = {
    iter: doc_mbed.iter({ gte: key }),
    coss: 1,
    pos: 1,
    add: 1
  };
  iter_lt.next = iter_gt;
  iter_gt.next = iter_lt;

  const dmap = {}; // document chunk index to token count map
  const found = []; // found chunks ranked by vector index
  const search = [ iter_gt, iter_lt ];
  const max_t = max_tokens || 2048;
  let tokens = 0;
  let which = 0;

    for (let i=0; i<500; i++) {
      const iter = search[which];
      const next = await iter.iter.next();

      if (!next) {
        // log({ iter_ended: iter.add });
        // current iterator ehausted
        iter.dead = true;
        which = 1 - which;
        // both iterators exhausted
        if (search[which].dead) {
          log("iterators exhaausted", i);
          break;
        }
        continue;
      }
      const [ key, rec ] = next;
      const coss = embed.cosine_similarity({ vector, index }, {
        vector: rec.vector,
        index: rec.index
      });

      // find how many new tokens this chunk range provides
      const [ fidx, fuid ] = key.split(":");
      const dtoks = dmap[fuid] = dmap[fuid] || {};
      const toks = rec.tokens;
      // log({ process : { rec.from, to: rec.to, tok: rec.tokens } });
      for (let x = 0; x < toks.length; x++) {
          let cidx = rec.from + x;
          if (dtoks[ cidx ] !== toks[x]) {
              // console.log({ set: cidx, val: toks[x], ln: toks.length });
              tokens += toks[x];
          }
      }

      found.push({
          i: iter.pos,
          coss,
          fuid,
          from: rec.from,
          to: rec.to,
      });

      iter.pos += iter.add;
      iter.coss += coss;

      if (coss < iter.next.coss && !iter.next.dead) {
        which = 1 - which;
      }
      if (tokens >= 8192) { //max_t * 2
        break;
      }
    }

    // close iterators to prevent mem leak
    search.forEach(iter => iter.iter.close());

    // sort by relevance, allows limiting to top matches
    found.sort((a, b) => { return b.coss - a.coss });

    // console.log({ found });
    // console.log({ dmap });

    // reduce to max chunks that will fit in embed window and
    // limit to chunks within 90% of max cosine_similarity value
    let min_coss = found[0].coss * (min_match || 0.90);
    let cnk_used = 0;
    let rec_used = 0;
    let tleft = max_t;
    let rank = {};

    outer: for (let rec of found) {
       if (rec.coss < min_coss) {
         if (debug) {
             console.log({ similarity_limit: { cnk_used, rec_used } });
         }
         tleft = 0;
         break;
       }
       const dtoks = dmap[rec.fuid];
       // attempt to use the whole chunk range unless
       // we hit a token limit while iterating
       if (debug) {
           console.log({ using: rec });
       }
       cnk_used++;
       for (let x = rec.from; x <= rec.to; x++) {
          tleft -= dtoks[x];
          if (tleft <= 0) {
              console.log({ tokens_limit: { cnk_used, rec_used } });
              break outer;
          }
          dtoks[x] = 0;
          rank[x] = cnk_used;
          rec_used++;
       }
    }

    // console.llog({ dmap });

    const subs = [];
    const segs = [];
    const plist = {};
    let totxt = 0;
    let totkn = 0;
    for (let [ fuid, dtoks ] of Object.entries(dmap)) {
        const sub = subs[fuid] = subs[fuid] || doc_chnk.sub(fuid);
        const dcnk = Object.entries(dtoks).map(r => [ parseInt(r[0]), r[1] ]);
        dcnk.sort((a,b) => a[0] - b[0]);
        for (let [ chunk, tokens ] of dcnk) {
           if (tokens === 0 ) {
             const rec = await sub.get(parseInt(chunk).toString().padStart(4,0));
             const [ page, from, to, tokens, text ] = rec;
             if (debug) {
                 console.log({ page, from, to, rank: rank[chunk] });
             }
             if (rank[chunk] === 1) {
                 plist[page] = page;
             }
             segs.push({ fuid, chunk, text, tln: text.length, tokens });
             totkn += tokens;
             totxt += text.length;
            }
          }
        }
 
    // time to consult the llm
    if (llm) {
      const embed = [
        "Based on the following context, succinctly answer the question at the end.\n",
        "Prepend the Paragraph ID(s) used for the answer in [brackers]\n",
        // "using only the text available in the fragments, Do not improvise.\n",
        "If the answer is not found in the provided context, reply that you do not ",
        "have a document related to the question,\n",
        "-----\n",
        ...segs.map(r => `\n<P>${r.text}"</P>\n`),
        `\nQuestion: ${query}\n\nAnswer: `
      ].join('');
      log({
        records: cnk_used,
        embeds: segs.length,
        tokens: totkn,
        text: embed.length,
      });
      if (debug_embed) {
          console.log("------------( embed )------------");
          console.log(embed);
          console.log("---------------------------------")
      }
      const start = Date.now();
      const pages = Object.keys(plist).map(r => parseInt(r));
      const answer = await node.call(llm, { query: embed, topic })
          .catch(error => {
              log({ llm_error: error });
              return { error: "llm service not responding" };
        });
      if (debug) {
          log({ time: Date.now() - start, pages, ...answer });
      }
      return { ...answer, pages };
    }

    return segs;
}

(async () => {
  const { embed, token } = await require('../llm/api').init();

  const store = await require('../lib/store').open(`${state.data_dir}/embed`);
  const doc_meta = store.sub('docs'); // document meta-data (size, type, source)
  const doc_chnk = store.sub('chunks'); // document chunks (page, line, text)
  const doc_mbed = store.sub('embed'); // embeds composed of chunks (vector, index)
  const llama_token = (await import('llama-tokenizer-js')).default;

  if (cli_store) {
      cli_server.server(store, 0, log);
  }

  Object.assign(state, {
    store,
    embed,
    token,
    doc_meta,
    doc_chnk,
    doc_mbed,
    llama_token,
  });

  await setup_node();
  await register_service();
})();