// client/client.js
const zlib = require('zlib');
const zstd = require('node-zstandard');

class Client {
  constructor(id, connection, compression = 'zlib') {
    this.id = id;
    this.connection = connection;
    this.environmentStats = {};
    this.compression = compression; // 'zlib' or 'zstd'
  }

  connect() {
    // Implement connection logic, e.g., WebSocket connection
    console.log(`Client ${this.id} connected.`);
    this.connection = `Connection for client ${this.id}`;
  }

  disconnect() {
    // Implement disconnection logic
    console.log(`Client ${this.id} disconnected.`);
    this.connection = null;
  }

  async sendData(data) {
    if (this.compression === 'zlib') {
      const compressed = await new Promise((resolve, reject) => {
        zlib.gzip(JSON.stringify(data), (err, buffer) => {
          if (err) reject(err);
          else resolve(buffer);
        });
      });
      this.connection.send(compressed);
    } else if (this.compression === 'zstd') {
      const compressor = new zstd.Compressor();
      const compressed = await compressor.compress(JSON.stringify(data));
      this.connection.send(compressed);
    }
  }

  async receiveData(buffer) {
    if (this.compression === 'zlib') {
      const decompressed = await new Promise((resolve, reject) => {
        zlib.gunzip(buffer, (err, buffer) => {
          if (err) reject(err);
          else resolve(JSON.parse(buffer.toString()));
        });
      });
      return decompressed;
    } else if (this.compression === 'zstd') {
      const decompressor = new zstd.Decompressor();
      const decompressed = await decompressor.decompress(buffer);
      return JSON.parse(decompressed.toString());
    }
  }

  updateStats(stats) {
    this.environmentStats = stats;
  }

  assertEnvironment(expectedStats) {
    return JSON.stringify(this.environmentStats) === JSON.stringify(expectedStats);
  }
}

module.exports = Client;