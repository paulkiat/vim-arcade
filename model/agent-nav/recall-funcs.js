// src/recall-funcs.js
import recallHome from './Recall_Home';

// Initialize Recall Home at the start
recallHome.initializeState();

// Function to save the last output
export const saveLastOutput = (output) => {
  recallHome.updateState('lastOutput', output);
};

// Function to retrieve the last output
export const retrieveLastOutput = () => {
  return recallHome.getState().lastOutput;
};

// Function to save the last action
export const saveLastAction = (action) => {
  recallHome.updateState('lastAction', action);
};

// Function to retrieve the last action
export const retrieveLastAction = () => {
  return recallHome.getState().lastAction;
};