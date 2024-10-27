// server/matchmaking.js
const ClientFactory = require('../client/clientFactory');

class MatchmakingServer {
  constructor(compression = 'zlib', preSeededConfigs = null) {
    this.clientFactory = new ClientFactory(compression);
    this.preSeededConfigs = preSeededConfigs || { server: []}
  }

  initialize_configs() {
    // Pre-seed the server with clients that are already connected
    this.preSeededConfigs.server.forEach(config => {
      // TODO: Logic to initialize server forms based on config
      // e.g., setting load, status, etc.

    // for (const config of this.preSeededConfigs.server) {
    //   config.client = this.addClient(config.id, config.compression);
    //   config.server = this.clientFactory.getServer(config.id, config.compression);
    // }
    });
  }



  addClient(id, compression) {
    const client = this.clientFactory.createClient(id, compression);
    // Additional matchmaking logic can be added here
    return client;
  }

  connectAllClients() {
    this.clientFactory.batchConnect(this.clientFactory.getClients());
  }

  assertClientEnvironments(expectedStats) {
    return this.clientFactory.assertClients(expectedStats);
  }

  getClients() {
    return this.clientFactory.getClients();
  }
}

module.exports = MatchmakingServer;