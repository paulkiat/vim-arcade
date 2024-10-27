//src/client-lib.js
/** @type {string[]} */
const weaponList = [
  'Punch',
  'Hammers',
  'Cannon',
  'Drill',
  'Long Swords'
];

let currentWeapon = 0;


function run(rotationPeriod, bytesPerSecond, maxPlayers, useVerbose, weapon) {
  if (weapon !== undefined && !weaponList.includes(weapon)) {
    console.log(`Invalid weapon ${weapon} specified - falling back to default`);
    weapon = 'Punch';
  }

  currentWeapon = weapon === undefined ? 0 : weaponList.indexOf(weapon);

  require('child_process').spawn(__dirname + '/client-socket', [maxPlayers, bytesPerSecond]).stdout
      .on('data', data => {
        process.stdout.write(data.toString());
    }).stderr.pipe(process.stderr).on('close', () => {
    console.log("Connection closed");
  });

  if (useVerbose) {
    console.log(`Using ${rotationPeriod} seconds as rotation period`);
    console.log(`Bytes per second is ${bytesPerSecond}`);

    let weaponIndex = 0;
    setInterval(()=> {
      process.stdout.write(
          `Setting the weapon to ${weaponList[currentWeapon]} (${++weaponIndex % weaponList.length})` +
          '\n');
    }, rotationPeriod * 1000);
  }
}

exports.run = run;
