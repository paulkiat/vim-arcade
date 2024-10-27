// src/cycleController.js
import agents from './agents/InitializeAgents';
import { actionResearchCycle, startCycle, stopCycle } from './Action_Research_Cycle';

/**
 * CycleController
 * Manages the execution of iterative development cycles.
 */
const CycleController = {
  currentIteration: 1,
  totalIterations: 3,
  
  /**
   * Executes a single iteration cycle following the Function → Feature → Function pattern.
   */
  executeIteration: async () => {
    if (CycleController.currentIteration > CycleController.totalIterations) {
      console.log('All iterations completed.');
      return;
    }
    
    console.log(`\n=== Starting Iteration ${CycleController.currentIteration} ===`);
    
    // Phase 1: Function - Optimize existing functionalities
    await CycleController.executePhase('Function', 'Optimize existing functionalities');
    
    // Phase 2: Feature - Implement new features
    await CycleController.executePhase('Feature', 'Implement new features');
    
    // Phase 3: Function - Refine and test functionalities
    await CycleController.executePhase('Function', 'Refine and test functionalities');
    
    console.log(`=== Iteration ${CycleController.currentIteration} completed ===\n`);
    CycleController.currentIteration += 1;
  },
  
  /**
   * Executes a specific phase within an iteration.
   * 
   * @param {string} phaseType - The type of phase ('Function' or 'Feature').
   * @param {string} phaseDescription - A brief description of the phase.
   */
  executePhase: async (phaseType, phaseDescription) => {
    console.log(`\n--- Phase: ${phaseType} ---`);
    console.log(phaseDescription);
    
    // Define actions based on phase type and workflow
    let actions = [];
    if (phaseType === 'Function') {
      // Function phase: Delete, Simplify, Automate
      actions = ['Delete Redundant Code', 'Simplify Complex Functions', 'Automate Repetitive Tasks'];
    } else if (phaseType === 'Feature') {
      // Feature phase: Feature-specific actions
      actions = ['Implement New Feature', 'Enhance Existing Feature', 'Automate Feature Deployment'];
    }
    
    // Assign actions to appropriate agents
    for (const action of actions) {
      await CycleController.assignActionToAgents(action);
    }
  },
  
  /**
   * Assigns a specific action to relevant agents based on the action type.
   * 
   * @param {string} action - The action to be performed.
   */
  assignActionToAgents: async (action) => {
    // Define mapping between actions and agents
    const actionAgentMap = {
      'Delete Redundant Code': ['Code Review', 'Feature Development'],
      'Simplify Complex Functions': ['Feature Development', 'Performance Optimization'],
      'Automate Repetitive Tasks': ['Performance Optimization', 'Deployment'],
      'Implement New Feature': ['Feature Development'],
      'Enhance Existing Feature': ['Feature Development', 'Code Review'],
      'Automate Feature Deployment': ['Deployment']
    };
    
    const agentRoles = actionAgentMap[action];
    
    if (!agentRoles) {
      console.warn(`No agents assigned for action: ${action}`);
      return;
    }
    
    // Assign action to each relevant agent
    for (const role of agentRoles) {
      const agent = agents.find(a => a.config.role === role);
      if (agent) {
        console.log(`Assigning action "${action}" to ${agent.config.role} (${agent.config.agentID})`);
        await agent.perform(action, agent.config.state)
          .then(output => {
            console.log(`Agent ${agent.config.agentID} output: ${output}`);
          })
          .catch(err => {
            console.error(`Agent ${agent.config.agentID} encountered an error:`, err);
          });
      } else {
        console.warn(`Agent with role "${role}" not found.`);
      }
    }
  },
  
  /**
   * Initiates the iterative development cycles.
   */
  start: async () => {
    for (let i = 0; i < CycleController.totalIterations; i++) {
      await CycleController.executeIteration();
    }
    
    console.log('All iterative cycles have been executed.');
  }
};

export default CycleController;