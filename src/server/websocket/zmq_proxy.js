// zmq proxy
const { Router, Dealer } = require('zeromq');

// Define the endpoints
const frontEnd = 'tcp://127.0.0.1:5556';
const backEnd = 'tcp://127.0.0.1:5558';
async function startProxy() {
  const frontend = new Router();
  const backend = new Dealer();

  await frontend.bind(frontendEndpoint);
  console.log(`Proxy frontend bound to ${frontendEndpoint}`);
  
  await backend.bind(backendEndpoint);
  console.log(`Proxy backend bound to ${backendEndpoint}`);
  
  // Start the proxy
  proxy(frontend, backend);

  console.log('ZMQ Proxy is running...');
}

startProxy().catch(err => {
  console.error('Failed to start ZMQ Proxy:', err);
});
