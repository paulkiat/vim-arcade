// client/clientFactory.js
const Client = require('./client');

class ClientFactory {
  constructor(compression = 'zlib') {
    this.clients = [];
    this.compression = compression; // Default compression
  }

  createClient(id, compression) {
    const clientCompression = compression || this.compression;
    const connection = null; // Initialize connection as null
    const client = new Client(id, connection, clientCompression);
    this.clients.push(client);
    return client;
  }

  batchConnect(clients) {
    clients.forEach(client => client.connect());
  }

  assertClients(expectedStats) {
    return this.clients.every(client => client.assertEnvironment(expectedStats));
  }

  getClients() {
    return this.clients;
  }
}

module.exports = ClientFactory;