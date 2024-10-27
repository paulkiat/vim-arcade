// Server-side code (app.js)
const express = require('express');
const app = express();
const http = require('http').createServer(app);
const io = require('socket.io')(http);
const path = require('path');

app.use(express.static('public'));

app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

const players = {};

io.on('connection', (socket) => {
  console.log('A player connected');
  
  players[socket.id] = {
    x: Math.floor(Math.random() * 10),
    y: 0,
    z: Math.floor(Math.random() * 10)
  };

  socket.emit('initGame', players[socket.id]);
  socket.broadcast.emit('newPlayer', { id: socket.id, ...players[socket.id] });
  socket.emit('players', players);

  socket.on('move', (data) => {
    players[socket.id] = { ...players[socket.id], ...data };
    socket.broadcast.emit('playerMoved', { id: socket.id, ...players[socket.id] });
  });

  socket.on('disconnect', () => {
    console.log('A player disconnected');
    delete players[socket.id];
    io.emit('playerDisconnected', socket.id);
  });
});

const PORT = process.env.PORT || 3000;
http.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});