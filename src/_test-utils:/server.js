// server.js
const express = require('express');
const http = require('http');
const socketIo = require('socket.io');
const Redis = require('ioredis');
const zstd = require('@gfx/zstd');

const app = express();
const server = http.createServer(app);
const io = socketIo(server);

const redis = new Redis({
    host: 'localhost',
    port: 6379,
});

const PLAYERS_PER_MATCH = 150; // Assuming 150 players per match like in Warzone
const SKILL_RANGE = 100; // Skill range for matchmaking

// Compress data using Zstandard
function compress(data) {
    return zstd.compress(Buffer.from(JSON.stringify(data)));
}

// Decompress data using Zstandard
function decompress(data) {
    return JSON.parse(zstd.decompress(data).toString());
}

// Matchmaking queue
class MatchmakingQueue {
    constructor() {
        this.queue = {};
    }

    addPlayer(player) {
        const skillRange = Math.floor(player.skill / SKILL_RANGE);
        if (!this.queue[skillRange]) {
            this.queue[skillRange] = [];
        }
        this.queue[skillRange].push(player);

        if (this.queue[skillRange].length >= PLAYERS_PER_MATCH) {
            return this.createMatch(skillRange);
        }
        return null;
    }

    createMatch(skillRange) {
        const players = this.queue[skillRange].splice(0, PLAYERS_PER_MATCH);
        return {
            id: Date.now(),
            players: players
        };
    }
}

const matchmakingQueue = new MatchmakingQueue();

io.on('connection', (socket) => {
    console.log('A user connected');

    socket.on('find_match', async (playerData) => {
        const player = {
            id: socket.id,
            skill: playerData.skill
        };

        const compressedPlayer = compress(player);
        await redis.set(`player:${socket.id}`, compressedPlayer);

        const match = matchmakingQueue.addPlayer(player);
        if (match) {
            match.players.forEach(async (p) => {
                const playerSocket = io.sockets.sockets.get(p.id);
                if (playerSocket) {
                    playerSocket.emit('match_found', { matchId: match.id });
                }
                await redis.del(`player:${p.id}`);
            });
        } else {
            socket.emit('waiting_for_match');
        }
    });

    socket.on('disconnect', async () => {
        console.log('User disconnected');
        await redis.del(`player:${socket.id}`);
    });
});

const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});

// Load testing script
const { performance } = require('perf_hooks');
const { io: ioClient } = require("socket.io-client");

async function simulatePlayer() {
    const socket = ioClient(`http://localhost:${PORT}`);
    
    socket.on('connect', () => {
        socket.emit('find_match', { skill: Math.floor(Math.random() * 1000) });
    });

    socket.on('match_found', (data) => {
        console.log(`Match found: ${data.matchId}`);
        socket.disconnect();
    });

    socket.on('waiting_for_match', () => {
        console.log('Waiting for match...');
    });
}

async function runLoadTest(numPlayers) {
    const start = performance.now();
    const players = Array(numPlayers).fill().map(() => simulatePlayer());
    await Promise.all(players);
    const end = performance.now();
    console.log(`Load test completed for ${numPlayers} players in ${(end - start) / 1000} seconds`);
}

// Run load test with 1 million players
runLoadTest(1000000);