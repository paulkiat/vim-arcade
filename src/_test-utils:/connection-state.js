// connectionStates.js

const Shareable = require('./shareable');
const DescriptorData = require('./descriptorData');

/**
 * @file connectionStates.js
 * @description Defines the ConnectionState classes and the ConnectionStateMachine for managing connection states.
 */

/**
 * Base class for all connection states.
 */
class ConnectionState extends Shareable {
  /**
   * Creates an instance of ConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super();
    this.m_pDesc = descriptor;
    this.stateMachine = null; // Will be set when pushed to the state machine
  }

  /**
   * Idle Update
   */
  idle() {
    // To be overridden by subclasses if needed
  }

  /**
   * Entering State
   */
  enter() {
    // To be overridden by subclasses
  }

  /**
   * Exiting State
   */
  exit() {
    // To be overridden by subclasses
  }

  /**
   * Pause entering State
   */
  pause() {
    // To be overridden by subclasses if needed
  }

  /**
   * Resume entering State
   */
  resume() {
    // To be overridden by subclasses if needed
  }

  /**
   * Handle user input
   * @param {string} arg 
   */
  parse(arg) {
    throw new Error("parse() must be implemented by subclass");
  }

  /**
   * Get the name of the state
   * @returns {string}
   */
  getName() {
    throw new Error("getName() must be implemented by subclass");
  }

  /**
   * Check if the state is playing
   * @returns {boolean}
   */
  isPlaying() {
    return false;
  }

  /**
   * Check if the state is in game
   * @returns {boolean}
   */
  isInGame() {
    return false;
  }

  /**
   * Called when the editor saves
   */
  onEditorSave() {
    // To be overridden by subclasses if needed
  }

  /**
   * Called when the editor aborts
   */
  onEditorAbort() {
    // To be overridden by subclasses if needed
  }

  /**
   * Called when the editor finishes
   */
  onEditorFinished() {
    // To be overridden by subclasses if needed
  }

  /**
   * Push a new state onto the state machine
   * @param {ConnectionState} state 
   */
  push(state) {
    if (this.stateMachine) {
      this.stateMachine.push(state);
    }
  }

  /**
   * Pop the current state from the state machine
   */
  pop() {
    if (this.stateMachine) {
      this.stateMachine.pop();
    }
  }

  /**
   * Swap the current state with a new state
   * @param {ConnectionState} state 
   */
  swap(state) {
    if (this.stateMachine) {
      this.stateMachine.swap(state);
    }
  }

  /**
   * Set descriptor data
   * @param {DescriptorData} d 
   */
  setDesc(d) {
    this.m_pDesc = d;
  }

  /**
   * Get descriptor data
   * @returns {DescriptorData}
   */
  getDesc() {
    return this.m_pDesc;
  }

  /**
   * Get the state machine managing this state
   * @returns {ConnectionStateMachine}
   */
  getStateMachine() {
    return this.stateMachine;
  }

  /**
   * Destructor equivalent in JavaScript (cleanup if necessary)
   */
  destroy() {
    // Cleanup resources if needed
  }
}

/**
 * Manages the state transitions for connections.
 */
class ConnectionStateMachine extends Shareable {
  /**
   * Creates an instance of ConnectionStateMachine.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super();
    this.m_pDesc = descriptor;
    this.m_States = [];
  }

  /**
   * Get the current state.
   * @returns {ConnectionState | null}
   */
  getState() {
    return this.m_States.length > 0 ? this.m_States[this.m_States.length - 1] : null;
  }

  /**
   * Idle Update
   */
  idle() {
    const currentState = this.getState();
    if (currentState) {
      currentState.idle();
    }
  }

  /**
   * Parse user input
   * @param {string} arg 
   */
  parse(arg) {
    const currentState = this.getState();
    if (currentState) {
      currentState.parse(arg);
    }
  }

  /**
   * Push a new state onto the stack
   * @param {ConnectionState} state 
   */
  push(state) {
    if (state) {
      state.stateMachine = this;
      this.m_States.push(state);
      state.enter();
      this.emit('stateChanged', this.getState());
    }
  }

  /**
   * Pop the current state from the stack
   */
  pop() {
    if (this.m_States.length > 0) {
      const state = this.m_States.pop();
      state.exit();
      this.emit('stateChanged', this.getState());
    }
  }

  /**
   * Swap the current state with a new state
   * @param {ConnectionState} state 
   */
  swap(state) {
    if (this.m_States.length > 0) {
      const currentState = this.m_States.pop();
      currentState.exit();
    }
    this.push(state);
  }

  /**
   * Get descriptor data
   * @returns {DescriptorData}
   */
  getDesc() {
    return this.m_pDesc;
  }

  /**
   * Destructor equivalent in JavaScript (cleanup if necessary)
   */
  destroy() {
    // Cleanup resources if needed
  }
}

/**
 * Represents the Playing state.
 */
class PlayingConnectionState extends ConnectionState {
  /**
   * Creates an instance of PlayingConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Handle user input in Playing state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Playing state
    console.log(`PlayingState parsing input: ${arg}`);
    // Example: Emit an event to update the UI or handle game actions
    this.emit('gameAction', arg);
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Playing";
  }

  /**
   * Indicates that this state is playing.
   * @returns {boolean}
   */
  isPlaying() {
    return true;
  }

  /**
   * Indicates that this state is in game.
   * @returns {boolean}
   */
  isInGame() {
    return true;
  }

  /**
   * Called when the editor aborts in Playing state.
   */
  onEditorAbort() {
    // Implement abort logic for Playing state
    console.log("Editor aborted in PlayingState.");
  }

  /**
   * Enter the Playing state.
   */
  enter() {
    super.enter();
    console.log("Entering PlayingState.");
    // Example: Update connection status
    this.m_pDesc.updateConnectionStatus('playing');
    this.emit('stateEntered', this.getName());
  }

  /**
   * Exit the Playing state.
   */
  exit() {
    super.exit();
    console.log("Exiting PlayingState.");
    this.emit('stateExited', this.getName());
  }
}

/**
 * Represents the Disconnecting state.
 */
class DisconnectConnectionState extends ConnectionState {
  /**
   * Creates an instance of DisconnectConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Idle Update for Disconnecting state.
   */
  idle() {
    // Implement idle logic for Disconnecting state if needed
    console.log("DisconnectingState is idling.");
  }

  /**
   * Handle user input in Disconnecting state.
   * @param {string} arg 
   */
  parse(arg) {
    // Nothing to do, we're disconnecting
    console.log(`DisconnectingState received input: ${arg}, no action taken.`);
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Disconnecting";
  }

  /**
   * Enter the Disconnecting state.
   */
  enter() {
    super.enter();
    console.log("Entering DisconnectingState.");
    this.m_pDesc.updateConnectionStatus('disconnecting');
    // Example: Perform cleanup or notify server
    this.emit('stateEntered', this.getName());
  }

  /**
   * Exit the Disconnecting state.
   */
  exit() {
    super.exit();
    console.log("Exiting DisconnectingState.");
    this.emit('stateExited', this.getName());
  }
}

/**
 * Represents the Lookup Host (Ident) state.
 */
class IdentConnectionState extends ConnectionState {
  /**
   * Creates an instance of IdentConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Enter the Ident state.
   */
  enter() {
    super.enter();
    console.log("Entering IdentConnectionState.");
    this.m_pDesc.updateConnectionStatus('identifying');
    // Example: Initiate host lookup
    this.emit('stateEntered', this.getName());
  }

  /**
   * Exit the Ident state.
   */
  exit() {
    super.exit();
    console.log("Exiting IdentConnectionState.");
    this.emit('stateExited', this.getName());
  }

  /**
   * Handle user input in Ident state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Ident state
    console.log(`IdentConnectionState parsing input: ${arg}`);
    // Example: Process identification response
    this.emit('identResponse', arg);
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Lookup Host";
  }
}

/**
 * Represents the Old Address state.
 */
class OldAddressConnectionState extends ConnectionState {
  /**
   * Creates an instance of OldAddressConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Enter the Old Address state.
   */
  enter() {
    super.enter();
    console.log("Entering OldAddressConnectionState.");
    this.m_pDesc.updateConnectionStatus('using_old_address');
    // Example: Handle old address logic
    this.emit('stateEntered', this.getName());
  }

  /**
   * Exit the Old Address state.
   */
  exit() {
    super.exit();
    console.log("Exiting OldAddressConnectionState.");
    this.emit('stateExited', this.getName());
  }

  /**
   * Handle user input in Old Address state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Old Address state
    console.log(`OldAddressConnectionState parsing input: ${arg}`);
    // Example: Validate old address or prompt for new address
    this.emit('oldAddressResponse', arg);
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Old Address";
  }
}

/**
 * Represents the Idle Timeout state.
 */
class IdleTimeoutConnectionState extends ConnectionState {
  /**
   * Creates an instance of IdleTimeoutConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
    this.m_StartTime = Date.now();
  }

  /**
   * Idle Update for Idle Timeout state.
   */
  idle() {
    const currentTime = Date.now();
    if (currentTime - this.m_StartTime > IdleTimeoutConnectionState.ms_IdleTimeout) {
      console.log("Idle timeout reached.");
      this.emit('timeout', this.getName());
      // Example: Transition to timeout state or disconnect
      this.pop(); // Pop current state
    }
  }

  /**
   * Reset the idle timeout timer.
   */
  resetTimeout() {
    this.m_StartTime = Date.now();
    console.log("Idle timeout reset.");
  }
}

// Static properties for IdleTimeoutConnectionState.
IdleTimeoutConnectionState.ms_IdleTimeout = 2 * 60 * 1000; // 2 minutes in milliseconds

/**
 * Represents the Get Name state.
 */
class GetNameConnectionState extends IdleTimeoutConnectionState {
  /**
   * Creates an instance of GetNameConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Enter the Get Name state.
   */
  enter() {
    super.enter();
    console.log("Entering GetNameConnectionState.");
    this.m_pDesc.updateConnectionStatus('getting_name');
    // Example: Prompt user for name
    this.emit('stateEntered', this.getName());
  }

  /**
   * Handle user input in Get Name state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Get Name state
    console.log(`GetNameConnectionState parsing input: ${arg}`);
    // Example: Validate and set user name
    if (arg && arg.trim() !== "") {
      this.m_pDesc.updateUserInfo({ userName: arg.trim() });
      this.emit('nameSet', arg.trim());
      this.pop(); // Pop current state
    } else {
      console.log("Invalid name input.");
      this.emit('invalidName', arg);
    }
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Get Name";
  }
}

/**
 * Represents the Confirm Name state.
 */
class ConfirmNameConnectionState extends IdleTimeoutConnectionState {
  /**
   * Creates an instance of ConfirmNameConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Enter the Confirm Name state.
   */
  enter() {
    super.enter();
    console.log("Entering ConfirmNameConnectionState.");
    this.m_pDesc.updateConnectionStatus('confirming_name');
    // Example: Prompt user to confirm name
    this.emit('stateEntered', this.getName());
  }

  /**
   * Handle user input in Confirm Name state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Confirm Name state
    console.log(`ConfirmNameConnectionState parsing input: ${arg}`);
    // Example: Confirm or reject the name
    if (arg.toLowerCase() === 'confirm') {
      this.emit('nameConfirmed', this.m_pDesc.userName);
      this.pop(); // Pop current state
    } else if (arg.toLowerCase() === 'reject') {
      this.emit('nameRejected');
      this.pop(); // Pop current state
    } else {
      console.log("Invalid confirmation input.");
      this.emit('invalidConfirmation', arg);
    }
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Confirm Name";
  }
}

/**
 * Represents the Main Menu state.
 */
class MenuConnectionState extends ConnectionState {
  /**
   * Creates an instance of MenuConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Resume the Main Menu state.
   */
  resume() {
    // Implement resume logic for Main Menu state
    console.log("Resuming MenuConnectionState.");
    this.emit('stateResumed', this.getName());
  }

  /**
   * Handle user input in Main Menu state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Main Menu state
    console.log(`MenuConnectionState parsing input: ${arg}`);
    // Example: Navigate to different menu options based on input
    this.emit('menuAction', arg);
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Main Menu";
  }

  /**
   * Enter the Main Menu state.
   */
  enter() {
    super.enter();
    console.log("Entering MenuConnectionState.");
    this.m_pDesc.updateConnectionStatus('in_menu');
    // Example: Render main menu UI
    this.emit('stateEntered', this.getName());
  }

  /**
   * Exit the Main Menu state.
   */
  exit() {
    super.exit();
    console.log("Exiting MenuConnectionState.");
    this.emit('stateExited', this.getName());
  }
}

/**
 * Represents the Online Creation (OLC) state.
 */
class OLCConnectionState extends ConnectionState {
  /**
   * Creates an instance of OLCConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Enter the OLC state.
   */
  enter() {
    super.enter();
    console.log("Entering OLCConnectionState.");
    this.m_pDesc.updateConnectionStatus('online_creation');
    // Example: Initialize OLC tools
    this.emit('stateEntered', this.getName());
  }

  /**
   * Exit the OLC state.
   */
  exit() {
    super.exit();
    console.log("Exiting OLCConnectionState.");
    this.emit('stateExited', this.getName());
  }

  /**
   * Indicates that this state is in game.
   * @returns {boolean}
   */
  isInGame() {
    return true;
  }
}

/**
 * Represents the Object Edit OLC state.
 */
class ObjectEditOLCConnectionState extends OLCConnectionState {
  /**
   * Creates an instance of ObjectEditOLCConnectionState.
   * @param {DescriptorData} descriptor 
   */
  constructor(descriptor) {
    super(descriptor);
  }

  /**
   * Handle user input in Object Edit OLC state.
   * @param {string} arg 
   */
  parse(arg) {
    // Implement parsing logic for Object Edit OLC state
    console.log(`ObjectEditOLCConnectionState parsing input: ${arg}`);
    // Example: Edit object properties
    this.emit('objectEdited', arg);
  }

  /**
   * Get the name of the state.
   * @returns {string}
   */
  getName() {
    return "Object edit";
  }

  /**
   * Called when the editor finishes in Object Edit OLC state.
   */
  onEditorFinished() {
    // Implement logic when editor finishes
    console.log("Editor finished in ObjectEditOLCConnectionState.");
    this.emit('editorFinished', this.getName());
  }

  /**
   * Enter the Object Edit OLC state.
   */
  enter() {
    super.enter();
    // Additional enter logic if necessary
  }

  /**
   * Exit the Object Edit OLC state.
   */
  exit() {
    super.exit();
    // Additional exit logic if necessary
  }
}

// Similarly, define other OLC states (RoomEdit, ZoneEdit, etc.) following the same pattern.
// For brevity, they are not all listed here, but the pattern remains consistent.

module.exports = {
  ConnectionState,
  ConnectionStateMachine,
  PlayingConnectionState,
  DisconnectConnectionState,
  IdentConnectionState,
  OldAddressConnectionState,
  IdleTimeoutConnectionState,
  GetNameConnectionState,
  ConfirmNameConnectionState,
  MenuConnectionState,
  OLCConnectionState,
  ObjectEditOLCConnectionState,
  // ... Add other states as needed
};