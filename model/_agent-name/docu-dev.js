// src/agents/DocumentationAgent.js
import recallHome from '../Recall_Home';
import { performAction } from '../Action_Research_Cycle';

const DA_Config = {
  agentID: 'DA-001',
  role: 'Documentation',
  skills: ['Technical Writing', 'Markdown Proficiency', 'Documentation Standards'],
  state: recallHome.getState(),
  actions: ['Update README.md', 'Write API Documentation', 'Create User Guide'],
  execute: (action, context) => {
    switch(action) {
      case 'Update README.md':
        return updateReadme(context);
      case 'Write API Documentation':
        return writeAPIDocumentation(context);
      case 'Create User Guide':
        return createUserGuide(context);
      default:
        return defaultDAAction(context);
    }
  }
};

const DocumentationAgent = {
  config: DA_Config,
  perform: async (action, context) => {
    const output = DA_Config.execute(action, context);
    recallHome.updateState('lastOutput', output);
    recallHome.updateState('lastAction', action);
    return output;
  }
};

export default DocumentationAgent;