//src/client.js

"use strict";

var argv = require('yargs')
    .usage('Usage: -c [-a] [-b] [-d] [-h] [-i] [-o] [-p] [-r [number of seconds]] [-s [number of bytes per second]] [-u [number of kbits/sec]]  [-v [level]] [-w [name of weapon]] [-z [max number of players]]')
    .boolean('a') // -a: Ace Attack
    .boolean('b') // -b: Blackjack
    .boolean('c') // -c: Cannon
    .boolean('d') // -d: Drill
    .boolean('h') // -h: Hammers
    .boolean('i') // -i: Ice Age
    .boolean('l') // -l: Long Swords
    .boolean('m') // -m: Marathon
    .boolean('p') // -p: Punch
    .number(['r', 'rotation'])  // -r [seconds]: rotation period. Default 120 seconds
    .number(['s', 'speed'])    // -s [bytes per second]: Data transfer speed in bytes/sec. Default 384 bytes/sec
    .boolean('x')               // -x: Exit without closing window
    .number(['u', 'userlimit']) // -u [max number of players]: Maximum number of players. Default 200 players
    .number(['v', 'verbose'])   // -v [level]: Verbosity level. Defaults to 3. Higher levels include more debugging information
    .string(['w', 'weapon'])    // -w [weapon]: Weapon to be used in the combat rounds
    .number(['z', 'maxplayers']) // -z [max number of players]: Maximum number of players. Default is 200 players
    .argv;

/** @type {number} */
const rotationPeriod = argv.rotation || 120;
let bytesPerSecond = argv.speed || 384,
    useVerbose = !!argv['verbose'],
    weapon = argv['weapon'] || 'Punch',
    maxPlayers = argv['maxplayers'] || 200;


function main() {
  console.log('Client started');

  require('./client-lib').run(rotationPeriod / 10, bytesPerSecond, maxPlayers, useVerbose, weapon);
}
main();

