// ../model/init-agent.js
import recallHome from '../Recall_Home';
import FeatureDevelopmentAgent from './FeatureDevelopmentAgent';
import CodeReviewAgent from './CodeReviewAgent';
import TestingAgent from './TestingAgent';
import DocumentationAgent from './DocumentationAgent';
import DeploymentAgent from './DeploymentAgent';
import PerformanceOptimizationAgent from './PerformanceOptimizationAgent';

// Initialize Recall Home
recallHome.initializeState();

// List of agents
const agents = [
  FeatureDevelopmentAgent,
  CodeReviewAgent,
  TestingAgent,
  DocumentationAgent,
  DeploymentAgent,
  PerformanceOptimizationAgent
];

// Initialize each agent's state
agents.forEach(agent => {
  agent.config.state = recallHome.getState();
});

export default agents;