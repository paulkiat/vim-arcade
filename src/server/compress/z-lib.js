// server/compression/z-lib.js
const zlib = require('zlib');

// Compress data
function compressDataZlib(data) {
  return new Promise((resolve, reject) => {
    zlib.gzip(JSON.stringify(data), (err, buffer) => {
      if (err) reject(err);
      else resolve(buffer);
    });
  });
}

// Decompress data
function decompressDataZlib(buffer) {
  return new Promise((resolve, reject) => {
    zlib.gunzip(buffer, (err, buffer) => {
      if (err) reject(err);
      else resolve(JSON.parse(buffer.toString()));
    });
  });
}

module.exports = { compressDataZlib, decompressDataZlib };