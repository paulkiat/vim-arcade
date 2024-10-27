// src/agents/TestingAgent.js
import recallHome from '../Recall_Home';
import { performAction } from '../Action_Research_Cycle';

const TA_Config = {
  agentID: 'TA-001',
  role: 'Testing',
  skills: ['Automated Testing', 'Manual Testing', 'Bug Reporting'],
  state: recallHome.getState(),
  actions: ['Run Unit Tests', 'Perform Integration Testing', 'Report Bug #456'],
  execute: (action, context) => {
    switch(action) {
      case 'Run Unit Tests':
        return runUnitTests(context);
      case 'Perform Integration Testing':
        return performIntegrationTesting(context);
      case 'Report Bug #456':
        return reportBug(456, context);
      default:
        return defaultTAAction(context);
    }
  }
};

const TestingAgent = {
  config: TA_Config,
  perform: async (action, context) => {
    const output = TA_Config.execute(action, context);
    recallHome.updateState('lastOutput', output);
    recallHome.updateState('lastAction', action);
    return output;
  }
};

export default TestingAgent;