// config/oracleConfig.js

const oracleConfig = {
  modelName: 'GPT-4-Oracle',
  modelType: 'GPT-4',
  parameters: {
    layers: 48,
    hiddenSize: 12288,
    attentionHeads: 96,
    maxSequenceLength: 2048,
    contextWindow: 4096,
  },
  training: {
    dataset: 'Extensive Knowledge Corpus',
    epochs: 10,
    batchSize: 16,
    learningRate: 2e-5,
    optimizer: 'AdamW',
  },
  deployment: {
    endpoint: 'https://api.zettelkasten-quest.com/oracle',
    gpuResources: {
      type: 'NVIDIA A100',
      memory: '80GB',
    },
    scaling: {
      minInstances: 2,
      maxInstances: 10,
      autoscale: true,
    },
  },
  knowledgeTransfer: {
    method: 'Knowledge Distillation',
    temperature: 2.0,
    alpha: 0.5,
    lossFunction: 'KL-Divergence',
  },
  monitoring: {
    logging: true,
    metrics: ['perplexity', 'surprisal', 'entropy'],
    alerts: {
      cpuUsage: 90,
      memoryUsage: 80,
      latency: 200, // in milliseconds
    },
  },
};

export default oracleConfig;