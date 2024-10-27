// server/simulate-client.js
const MatchmakingServer = require('../server/lib/matchmaking');

const matchmakingServer = new MatchmakingServer();

// Simulate adding clients
for (let i = 1; i <= 5; i++) {
  matchmakingServer.addClient(i);
}

// Connect all clients
matchmakingServer.connectAllClients();

// Update environment stats
matchmakingServer.getClients().forEach(client => {
  client.updateStats({
    cpu: `${Math.floor(Math.random() * 100)}%`,
    memory: `${Math.floor(Math.random() * 16)}GB`
  });
});

// Assert environments
const expectedStats = {
  cpu: '75%',
  memory: '8GB'
};
const allMatch = matchmakingServer.assertClientEnvironments(expectedStats);
console.log('All clients match expected stats:', allMatch);