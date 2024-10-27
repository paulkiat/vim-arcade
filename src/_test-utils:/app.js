// app.js

const { args, env } = require('../lib/util');
const { file_drop } = require('../app/doc-client');
const log = require('../lib/util').logpre('app');
const web = require('../lib/web');
const express = require('express');
const ConnectionHandler = require('./connectionHandler.js');
const serveStatic = require('serve-static');

const app_handler = express();
const state = require('./service.js').init();
const connectionHandler = new ConnectionHandler();

// WebSocket proxy handler
const ws_handler = web.ws_proxy_path(state.node);

// Assign additional properties to the state object
Object.assign(state, {
  direct: args['direct'] || false,
  app_dir: env('APP_DIR') || args['app-dir'] || 'app',
  app_port: env('APP_PORT') ?? args['app-port'] ?? (args.prod ? 80 : 7000),
  app_handler,
  ws_handler
});

/**
 * Sets up the node with event listeners for connection states.
 */
async function setup_node() {
  // Re-announce the web app when the proxy connection bounces
  state.node.on_reconnect(announce_service);

  // Listen for events from the ConnectionStateMachine
  connectionHandler.stateMachine.on('stateChanged', (newState) => {
    if (newState) {
      // Example: Handle state change (e.g., update UI, log, etc.)
      log(`New state: ${newState.getName()}`);
    } else {
      log('No active state.');
    }
  });
}

/**
 * Announces the web service to the network or discovery service.
 */
async function announce_service() {
  const { node, app_id, net_addrs } = state;
  log({ register_app_web: app_id, static_dir: state.app_dir });
  node.publish("service-up", {
    app_id,
    type: "web-server",
    // Put app into direct access dev mode vs production proxy
    direct: state.direct,
    web_port: state.app_port,
    web_addr: net_addrs,
  });
}

/**
 * Middleware to inject app-id context into requests.
 * For app services like "file_drop" to have app-id context.
 */
function injectXAppId(req, res, next) {
  req.headers['x-app-id'] = state.app_id;
  next();
}

/**
 * Serves local app web assets and sets up middleware.
 */
async function setup_app_handlers() {
  app_handler
    .use(injectXAppId)
    .use(web.parse_query)
    .use(file_drop(state))
    .use((req, res, next) => {
      const url = req.url;
      const appurl = `/app/${state.app_id}`;
      // Limit requests to contents of /app (e.g., ignore hub and org)
      if (!(url === appurl || url.startsWith(`${appurl}`))) {
        res.writeHead(404, { 'Content-Type': 'text/plain' });
        res.end('404 Invalid URL');
      } else {
        // Rewrite app url to remove /app/<app-id> prefix
        req.url = req.url.substring(appurl.length) || '/';
        next();
      }
    })
    .use(serveStatic(`web/${state.app_dir}`, { index: [ "index.html=" ]}))
    .use((req, res) => {
      res.writeHead(404, { 'Content-Type': 'text/plain' });
      res.end('404 Not Found');
    });
}

/**
 * Sets up the ConnectionStateMachine integration with the web server.
 */
async function setup_connection_handlers() {
  // Example: Listen for specific events from the ConnectionStateMachine
  connectionHandler.stateMachine.on('stateChanged', (newState) => {
    if (newState) {
      // For instance, log the state change or perform actions based on state
      log(`Connection state changed to: ${newState.getName()}`);
    }
  });

  // Handle specific events from states if needed
  connectionHandler.stateMachine.on('stateChanged', (newState) => {
    if (newState instanceof require('./connectionStates.js').PlayingConnectionState) {
      // Perform actions specific to Playing state
      log('Entered Playing state.');
    }
    // Handle other states similarly
  });
}

(async () => {
  await setup_node();
  await setup_app_handlers();
  await setup_connection_handlers();
  await web.start_web_listeners(state);
  await announce_service();

  // Periodically perform idle updates
  setInterval(() => {
    connectionHandler.performIdleUpdate();
  }, 1000); // Every second
})();