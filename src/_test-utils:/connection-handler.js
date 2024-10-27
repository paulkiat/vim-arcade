// connectionHandler.js

const {
  ConnectionStateMachine,
  MenuConnectionState,
  PlayingConnectionState,
  DisconnectConnectionState,
  // ... import other states as needed
} = require('./connectionStates');
const DescriptorData = require('./descriptorData');

/**
 * @file connectionHandler.js
 * @description Initializes and manages the ConnectionStateMachine.
 */

/**
 * ConnectionHandler class to manage connection states.
 */
class ConnectionHandler {
  /**
   * Creates an instance of ConnectionHandler.
   */
  constructor() {
    // Initialize DescriptorData
    this.descriptor = new DescriptorData({
      connectionId: 'conn-001',
      userId: 'user-123',
      userName: 'Alice',
      connectionStatus: 'disconnected',
      metadata: { /* Additional data */ }
    });

    // Initialize ConnectionStateMachine
    this.stateMachine = new ConnectionStateMachine(this.descriptor);

    // Register event listeners
    this.stateMachine.on('stateChanged', (newState) => {
      if (newState) {
        console.log(`State changed to: ${newState.getName()}`);
        // Here, you can integrate with your web application's state management or logging
      } else {
        console.log('No active state.');
      }
    });

    // Initialize with Menu state
    this.initMenuState();
  }

  /**
   * Initializes the Menu state.
   */
  initMenuState() {
    const menuState = new MenuConnectionState(this.descriptor);
    menuState.on('menuAction', (action) => {
      if (action === 'startGame') {
        this.startGame();
      } else if (action === 'disconnect') {
        this.disconnect();
      }
      // Handle other menu actions
    });
    this.stateMachine.push(menuState);
  }

  /**
   * Starts the game by transitioning to the Playing state.
   */
  startGame() {
    const playingState = new PlayingConnectionState(this.descriptor);
    playingState.on('gameAction', (action) => {
      console.log(`Game action received: ${action}`);
      // Handle game-specific actions
    });
    this.stateMachine.push(playingState);
  }

  /**
   * Disconnects by transitioning to the Disconnecting state.
   */
  disconnect() {
    const disconnectState = new DisconnectConnectionState(this.descriptor);
    disconnectState.on('stateExited', (stateName) => {
      console.log(`${stateName} has exited. Connection is now disconnected.`);
      // Update web application state to reflect disconnection
    });
    this.stateMachine.push(disconnectState);
  }

  /**
   * Handles user input by passing it to the state machine.
   * @param {string} input 
   */
  handleInput(input) {
    this.stateMachine.parse(input);
  }

  /**
   * Performs idle updates, should be called periodically (e.g., via setInterval).
   */
  performIdleUpdate() {
    this.stateMachine.idle();
  }
}

module.exports = ConnectionHandler;