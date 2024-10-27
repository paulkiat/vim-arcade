// server/server.js
const zmq = require('zeromq');
const WebSocket = require('ws');
const http = require('http');
const bodyParser = require('body-parser');
const WebSocketServer = require('ws').Server;
const port = process.env.PORT || 3001;
const server = http.createServer(app);

server.listen(port);
console.log(`Server listening on port ${port}`);

// server/app.js
app.use(bodyParser.json());

// server/auth.js
const passport = require('passport')

// server/matchmaking.js
const express = require('express');
const session = require('express-session');
const passport = require('./auth');
const MatchmakingServer = require('./matchmaking');
const { compressDataZlib, decompressDataZlib } = require('./compression/zlibCompression');
const { compressDataZstd, decompressDataZstd } = require('./compression/zstdCompression');

const app = express();
app.use(express.json());

// Initialize Matchmaking Server with zstd compression
const matchmakingServer = new MatchmakingServer('zstd', {
  server: [
    { id: 1, status: 'active', load: 0.75 },
    { id: 2, status: 'active', load: 0.60 },
    // Add more pre-seeded server configurations
  ]
});
matchmakingServer.initializeConfigs();

// Session and Passport Initialization
app.use(session({
  secret: 'your_secret_key',
  resave: false,
  saveUninitialized: true
}));

app.use(passport.initialize());
app.use(passport.session());

// Authentication Routes
app.get('/auth/provider', passport.authenticate('oauth2'));

app.get('/auth/callback',
  passport.authenticate('oauth2', { failureRedirect: '/login' }),
  function(req, res) {
    res.redirect('/matchmaking');
  });

app.get('/logout', (req, res) => {
  req.logout();
  res.redirect('/');
});

// Middleware to check authentication
function ensureAuthenticated(req, res, next) {
  if (req.isAuthenticated()) { return next(); }
  res.redirect('/login');
}

// Matchmaking Routes
app.get('/matchmaking', ensureAuthenticated, (req, res) => {
  res.send('Welcome to Matchmaking!');
});

// API Endpoints with Compression
app.post('/api/matchmaking/addClient', ensureAuthenticated, async (req, res) => {
  const clientId = req.body.id;
  const compression = req.body.compression || 'zstd'; // 'zlib' or 'zstd'
  const client = matchmakingServer.addClient(clientId, compression);
  res.status(200).json({ message: `Client ${clientId} added.` });
});

app.post('/api/matchmaking/connectAll', ensureAuthenticated, async (req, res) => {
  matchmakingServer.connectAllClients();
  res.status(200).json({ message: 'All clients connected.' });
});

app.post('/api/matchmaking/assertEnvironments', ensureAuthenticated, async (req, res) => {
  const expectedStats = req.body.expectedStats;
  const result = matchmakingServer.assertClientEnvironments(expectedStats);
  res.status(200).json({ allMatch: result });
});

// Example Endpoint with Compression
app.post('/api/data', ensureAuthenticated, async (req, res) => {
  const data = req.body.data;
  const compressionMethod = req.body.compression || 'zlib';

  try {
    let compressed;
    if (compressionMethod === 'zlib') {
      compressed = await compressDataZlib(data);
      res.setHeader('Content-Encoding', 'gzip');
      res.setHeader('Content-Type', 'application/json');
    } else if (compressionMethod === 'zstd') {
      compressed = await compressDataZstd(data, 5);
      res.setHeader('Content-Encoding', 'zstd');
      res.setHeader('Content-Type', 'application/octet-stream');
    } else {
      compressed = Buffer.from(JSON.stringify(data));
    }

    res.send(compressed);
  } catch (error) {
    res.status(500).send('Compression Error');
  }
});

// Start Server
app.listen(4000, () => {
  console.log('Matchmaking server running on port 4000');
});