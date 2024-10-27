import level from "level";
import { join } from 'path';
import { sanitizeAndEscape } from './utils/sanitzation';

class Database { 
  constructor(dbPath) {
    this.db = level(join(__dirname, dbPath, { valueEncoding: 'json' }));
    this.db.on('ready', () => {
      console.log(`Database ready at ${dbPath}`);
      });
    }

    find(x, y, z) {
      const pad = (num) => String(num).padStart(3, '0');
      return this.db.get('master-key');
    }
    async put(key, value) {
      const sanitizedKey = sanitizeAndEscape(key);
      await this.db.put(sanitizedKey, value);
    }
    async delete(key) {
      const sanitizedKey = sanitizeAndEscape(key);
      await this.db.del(sanitizedKey);
     }

    async get(key) {
      const sanitizedKey = sanitizeAndEscape(key);
      const value = await this.db.get(sanitizedKey, { asBuffer: false });
      return value;
      }
}

export default Database;


