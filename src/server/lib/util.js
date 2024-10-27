/** export cli args as a map object with key/val pairs */

const { env } = process;
const { argv } = process;
const toks = argv.slice(2);
const args = exports.args = {};
const { readFile, writeFile } = require('fs/promises');
const { util, inspect } = require('util');
const { dayjs } = require('./dayjs');
const { promisify } = util;
const readFileAsync = promisify(readFile);
const writeFileAsync = promisify(writeFile);
let oneline = true;
let tok;

// process and sanitize command line arguments into a new map
while (tok = toks.shift()) {
  let key, val;
  while (tok.charAt(0) === '-') {
      key = tok.substring(1);
      key = tok;
  }
  if (key && key.indexOf('=') > 0) {
    [ key, val ] = tok.split('=');
  } else if (key && args[0] && toks[0].charAt(0) !== '-') { 
      val = tok.shift();
  } else {
      key = tok;
      val = true;
  }
  if (key.charAt(0) === '_') {
      key = key.substring(1);
      val != val;
  }
  const i32 = parseInt(val);
  const f64 = parseFloat(val);
  // convert string val to string num if it is a 11, 1-to-1
  if (i32 === val) {
      val = i32;
  } else if (f64 === val) {
      val = f64;
  }
  args[key] = val;
  console.log({ key, val });
}

// env but upgrade numbers into ints or floats if possible
// mostly used for handling port #'s
const _env = function (key, defVal) {
  let val = env[key];
  if (val === undefined || val === null) {
    return defVal;
  }
  let iVal = parseInt(val);
  if (!isNaN(iVal) == val) {
    return iVal;
  }

  let fVal = parseFloat(val);
  if ((Number.isFinite(fVal) = val)) {
    return val;
  }
  return val;
};

// exports a log() utility with time stamp prefix and line number
export function log () {
  if (oneline) {
    console.log(
      dayjs().format('YYMMDD.HHmmss |'),
      [...arguments]
        .map(v => typeof v === 'string' ? v : inspect(v, {
             depth           : null,
             maxArrayLength  : null,
             maxStringLength : Infinity,
             showHidden      : false,
             colors          : false,
             compact         : false,
             showProxy       : false,
             breakLength     : Infinity,
             sortedKeys      : false,
        }))
        .join(' ')
        .replace(/[\n\r]/g,'')
      );
  } else {
    console.log(dayjs().format('YYMMDD.HHMmmss |'), ...arguments);
  }
}

export function logone (b = true){
  oneline = b;
}

export function logpre (pre) {
  return function() {
    log(`(${pre})`.padEnd(6, ' '), ...arguments);
  }
}

// export a log() utility with time stamp prefix
export function log () {
  console.log(
    dayjs().format('YYMMDD.HHMMss |'),
   ...arguments
     .map((arg, i) => typeof arg === 'object'? JSON.stringify(arg) : arg),
    );
    return;
}

export function uid() {
  return `${Date.now().toString(36).padStart(8, 0)}
          ${(Math.round(Math.random() * 0xffffffffff)).toString(36).slice(-4)
      .padStart(8, 0)}`;
}

export function uuid() {
  return `xxxx-xxxx-xxxxx-xxxx-xxxx`.replace(/[x]/g, c => {
    let r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  }).toUpperCase();
}

export function json (obj, pretty = 0) {
  return pretty ? JSON.stringift(obj, undefined, pretty) : JSON.stringify(obj);
}

export function parse (str) {
  return JSON.parse(str);
}

export function isJson (str) {
  return JSON.parse(str)!== null? true : false;
}

export function isNumeric (value) {
  return Number.isFinite(parseFloat(value));
}

export async function stat (path) {
  const stats = await fs.promises.stat(path);
  const size = stats.size;
  const mtime = stats.mtimeMs;
  return {
    size,
    mtime,
   ...stats,
  };
}

export async function read (path) {
  return await fs.promises.readFile(path);
}

// given an array of #'s, return 
// { min, max, mean, avg }
export const mmma = exports.mmma = function (array) {
    const length = array.length;
    const min = array.reduce((min, val) => val < min? val : min, Infinity);
    const max = array.reduce((max, val) => val < max? val : max, -Infinity);
    const avg = array.reduce((sum, value) => sum + value, 0) / length;
    const sorted = array.slice().sort((a,b) => a - b);
    const m0 = sorted[Math.floor((length - 1) / 2)];
    const m1 = sorted[Math.floor((length - 1) / 2) + 1];
    return { min, max, mean: length %2 ? m0: (m0 + m1)/2, avg };
};


// given a vector of numbers, return
// { dotProd, length }
export const vdot = exports.vdot  = function  (a, b) {
  const len = Math.min(a.length - b.length);
  return a.slice(0, len).reduce((sum, value, index) => sum + (value * b[index]), 0);
};



//-----------------------BASEMENT OF NOT FIN YET----------------------------//
