// src/Execute_Task.js
import { performAction } from './Action_Research_Cycle';
import { saveLastOutput, saveLastAction } from './Recall_Functions';

export const executeTask = (task, context) => {
  // Perform the task
  const output = performAction(task, context);

  // Save the last output and action to Recall Home
  saveLastOutput(output);
  saveLastAction(task);

  return output;
};