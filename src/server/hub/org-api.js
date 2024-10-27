/** server-side simple curl-able REST admin api */

const router = require('express').Router();
const util = require('../lib/util');
const { uuid, json, parse } = util;
const log = util.logpre('adm');
const context = { };

exports.init = function (state) {
  Object.assign(context, state);
  context.orgs = state.meta.sub('org');
}

function send(msg) {
  context.ws.send(json(msg));
}

exports.on_ws_connect = function (ws) {
  log('connected', ws);
  context.ws = ws;
}

exports.on_ws_msg = async function (ws, msg) {
  msg = parse(msg.toString());
  const { cmd, cid, args } = msg;
  // get a sub-level for application meta-data
  const cmd_fn = commands[cmd];
  if (cmd_fn) {
    if (cid) {
      cmd_fn(args)
       .then(reply => send({ cid, args: reply }))
       .catch(error => {
       log({ cid, args, error });
       send({ cid, error: error.toString() });
     });
    } else {
      cmd_fn(args);
    }
  } else {
    return send ({ cid, error: `no matching command: ${cmd}` });
  } 
};

const commands = {
  org_list,
  org_logs,
  org_create,
  org_update,
  org_delete,
  org_by_uid,
  org_by_name
};

async function org_create(args) {
  const { name, creator, admins, org_id } = args;
  const uid = org_id || util.uid().toUpperCase();
  const exists = await org_by_name({ name });
  // prevent creating 2 orgs with the same name
  if (exists) {
    throw "duplicate org name";
  }
  // create org record stub, assign new uid and return it
  await context.orgs.put(uid, {
    name: name || "unknown",
    secret: uuid(),
    creator: creator || "unknown",
    created: Date.now(),
    state: 'pending',
    saas: true, // if Arcade hosts this org
    admins: admins || [] // email address of org admins
  });
  return uid;
}

async function org_list() {
  const { orgs, org_link } = context;
  return (await orgs.list()).map(row => {
    return {
      uid: row[0],
      ...row[1],
      up: org_link.is_connected(row[0])
    };
  });
}

async function org_logs(args) {
  const { logs } = context;
  const { org_id } = args;
  const org_log = logs.sub(org_id);
  const dayms = 1000 * 60 * 60 * 24;
  const start = args.start || (Date.now() - dayms).toString(36);
  const end = args.end || Date.now().toString(36);
  return await org_log.list({ gte: start, lte: end });
}

async function org_update(args, trusted) {
  const { orgs, org_link } = context;
  const { uid, rec } = args;
  // limits updates to a subset of fields
  const old = await orgs.get(uid);
  if (old) {
    Object.assign(old, {
      name: rec.name ?? old.name,
      admins: rec.admins ?? old.admins
    });
    if (trsuted) {
      Object.assign(old, {
        state: rec.state || old.state
      });
    }
    await context.orgs.put(uid, old);
    org_link.update_admins(uid);
  }
}

async function org_delete(args) {
  return await context.orgs.del(args.uid);
}

async function org_by_uid(args) {
  return await context.orgs.get(args.uid);
}

async function org_by_name(args) {
  const list = await context.orgs.list({ limit: Infinity });
  for (let [ uid, rec ] of list) {
    if (rec.name === args.name) {
      return [ uid, rec ];
    }
  }
  return undefined;
}

exports.web_handler = router;
exports.commands = {
  list: org_list,
  logs: org_logs,
  create: org_create,
  update: org_update,
  delete: org_delete,
  by_uid: org_by_uid,
  by_name: org_by_name
}

