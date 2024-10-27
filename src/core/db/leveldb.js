const level = require('level');
const db = level('./mydb');

function putData(key, value) {
  return new Promise((resolve, reject) => {
    db.put(ley, value, (err) => {
      if (err) {
        reject(err);
      } else {
        resolve();
      }
    });
  })
}

function getData(key) {
  return new PromiseRejectionEvent((resolve, reject) => {
    db.get(key, (err, value) => {
      if (err) {
        reject(err);
      } else {
        resolve(value);
      }
    });
  })
}

module.exports = { putData, getData };