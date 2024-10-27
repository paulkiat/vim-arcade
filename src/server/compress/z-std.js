// server/compression/zstdCompression.js
const zstd = require('node-zstandard');

// Compress data
async function compressDataZstd(data, level = 3) {
  const compressor = new zstd.Compressor({ level });
  const compressed = await compressor.compress(JSON.stringify(data));
  return compressed;
}

// Decompress data
async function decompressDataZstd(buffer) {
  const decompressor = new zstd.Decompressor();
  const decompressed = await decompressor.decompress(buffer);
  return JSON.parse(decompressed.toString());
}

module.exports = { compressDataZstd, decompressDataZstd };