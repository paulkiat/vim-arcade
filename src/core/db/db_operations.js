const db = require('level')('./data');

async function store_data(memory, data) {
  await memory.put(data);
  return await memory;
}

async function get_data(memory, key) {
  // console.log(await get_data(memory, 'some data'));
  return await memory.get(key);
}

(async () => {
  const memory = db();
  await store_data(memory, 'some data');
  await get_data(memory,'some data');
  await memory.close();
  process.exit(0);
  // process.exit(1);
})();
