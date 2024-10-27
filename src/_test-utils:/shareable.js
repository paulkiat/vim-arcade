// shareable.js

const EventEmitter = require('events');

/**
 * @file shareable.js
 * @description Provides a base class with event handling capabilities for sharing resources and communication.
 */

/**
 * Base class representing a shareable entity with event capabilities.
 */
class Shareable extends EventEmitter {
  constructor() {
    super();
  }

  /**
   * Broadcast an event to all listeners.
   * @param {string} event 
   * @param  {...any} args 
   */
  broadcast(event, ...args) {
    this.emit(event, ...args);
  }
}

module.exports = Shareable;