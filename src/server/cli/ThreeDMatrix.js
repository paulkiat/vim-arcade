// ThreeDMatrix.js

const   { uuid } = require('uuid');
const { levelup } = require('levelup');
const { leveldown } = require('leveldown');
const level = require('level');
const path = require('path');
const db = levelup(path.join(__dirname, 'db'), { db: leveldown });
const _ = require('lodash');
const { sanitizeAndEscape } = require('./utils/sanitization');

class ThreeDMatrix {
  constructor(dbPath) {
    this.db = db;
    this.dbPath = dbPath;
    this.db = level(dbPath, { valueEncoding: 'json' });
  }

  /**
   * Generates a zero-padded key based on coordinates.
   * @param {number} x 
   * @param {number} y 
   * @param {number} z 
   * @returns {string} Key in format 'xxx:yyy:zzz'
   */
  generateKey(x, y, z) {
    const pad = (num) => String(num).padStart(3, '0');
    return `${pad(x)}:${pad(y)}:${pad(z)}`;
  }

  /**
   * Inserts or updates a blob at the specified coordinates.
   * @param {number} x 
   * @param {number} y 
   * @param {number} z 
   * @param {Object} blobData 
   */
  async setBlob(x, y, z, blobData) {
    const key = this.generateKey(x, y, z);
    const sanitizedData = {
      ...blobData,
      content: sanitizeAndEscape(JSON.stringify(blobData.content)),
      style: sanitizeAndEscape(JSON.stringify(blobData.style)),
      interactions: sanitizeAndEscape(JSON.stringify(blobData.interactions)),
      updatedAt: new Date().toISOString(),
    };
    await this.db.put(key, sanitizedData);
    console.log(`Blob set at (${x}, ${y}, ${z}) with key ${key}.`);
  }

  /**
   * Retrieves a blob from the specified coordinates.
   * @param {number} x 
   * @param {number} y 
   * @param {number} z 
   * @returns {Object|null} The blob data or null if not found.
   */
  async getBlob(x, y, z) {
    const key = this.generateKey(x, y, z);
    try {
      const blob = await this.db.get(key);
      return blob;
    } catch (error) {
      if (error.notFound) {
        console.warn(`Blob not found at (${x}, ${y}, ${z}) with key ${key}.`);
        return null;
      } else {
        throw error;
      }
    }
  }

  /**
   * Deletes a blob from the specified coordinates.
   * @param {number} x 
   * @param {number} y 
   * @param {number} z 
   */
  async deleteBlob(x, y, z) {
    const key = this.generateKey(x, y, z);
    try {
      await this.db.del(key);
      console.log(`Blob deleted from (${x}, ${y}, ${z}) with key ${key}.`);
    } catch (error) {
      if (error.notFound) {
        console.warn(`No blob to delete at (${x}, ${y}, ${z}) with key ${key}.`);
      } else {
        throw error;
      }
    }
  }

  /**
   * Traverses the 3D grid based on a direction vector.
   * @param {number} x 
   * @param {number} y 
   * @param {number} z 
   * @param {Object} direction - { dx, dy, dz }
   * @returns {Object} New coordinates.
   */
  traverse(x, y, z, direction) {
    const { dx, dy, dz } = direction;
    return {
      x: x + dx,
      y: y + dy,
      z: z + dz,
    };
  }

  /**
   * Retrieves all blobs within a specified range.
   * @param {Object} range - { xMin, xMax, yMin, yMax, zMin, zMax }
   * @returns {Array<Object>} List of blobs.
   */
  async getBlobsInRange(range) {
    const { xMin, xMax, yMin, yMax, zMin, zMax } = range;
    const pad = (num) => String(num).padStart(3, '0');
    const startKey = `${pad(xMin)}:${pad(yMin)}:${pad(zMin)}`;
    const endKey = `${pad(xMax)}:${pad(yMax)}:${pad(zMax)}`;

    const blobs = [];

    return new Promise((resolve, reject) => {
      this.db.createReadStream({
        gte: startKey,
        lte: endKey,
      })
        .on('data', (data) => {
          blobs.push(data.value);
        })
        .on('error', (err) => {
          reject(err);
        })
        .on('end', () => {
          resolve(blobs);
        });
    });
  }
}

module.exports = ThreeDMatrix;