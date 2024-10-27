// config/studentConfig.js

const studentConfig = {
  modelName: 'GPT-2-Student',
  modelType: 'GPT-2',
  parameters: {
    layers: 12,
    hiddenSize: 768,
    attentionHeads: 12,
    maxSequenceLength: 1024,
    contextWindow: 2048,
  },
  training: {
    dataset: 'Distilled Knowledge from Oracle',
    epochs: 5,
    batchSize: 32,
    learningRate: 5e-5,
    optimizer: 'AdamW',
  },
  deployment: {
    endpoint: 'https://api.zettelkasten-quest.com/student',
    gpuResources: {
      type: 'NVIDIA V100',
      memory: '32GB',
    },
    scaling: {
      minInstances: 1,
      maxInstances: 5,
      autoscale: true,
    },
  },
  knowledgeTransfer: {
    method: 'Knowledge Distillation',
    temperature: 1.0,
    alpha: 0.7,
    lossFunction: 'Mean Squared Error',
  },
  monitoring: {
    logging: true,
    metrics: ['perplexity', 'surprisal', 'entropy'],
    alerts: {
      cpuUsage: 85,
      memoryUsage: 75,
      latency: 150, // in milliseconds
    },
  },
};

export default studentConfig;