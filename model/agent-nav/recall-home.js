// src/Recall_Home.js
const recallHome = {
  state: {},

  initializeState: () => {
    recallHome.state = {
      lastAction: null,
      lastOutput: null,
      currentTasks: [],
      insights: null,
      // Add more state properties as needed
    };
    console.log('Recall Home initialized.');
  },

  updateState: (key, value) => {
    recallHome.state[key] = value;
    console.log(`State updated: ${key} = ${value}`);
  },

  getState: () => {
    return recallHome.state;
  },

  resetState: () => {
    recallHome.initializeState();
    console.log('Recall Home state has been reset.');
  }
};

export default recallHome;