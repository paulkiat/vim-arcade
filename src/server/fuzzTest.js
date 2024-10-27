// server/fuzzTest.js
const serverConfigSchema = require('./configSchema');
const Joi = require('joi');
const faker = require('faker');
const { compressDataZlib, decompressDataZlib } = require('./compression/zlibCompression');
const { compressDataZstd, decompressDataZstd } = require('./compression/zstdCompression');

// Function to generate random server configurations
function generateRandomConfig() {
  const numServers = faker.datatype.number({ min: 1, max: 10 });
  const servers = [];

  for (let i = 0; i < numServers; i++) {
    servers.push({
      id: i + 1,
      status: faker.random.arrayElement(['active', 'inactive', 'maintenance']),
      load: faker.datatype.float({ min: 0, max: 1, precision: 0.01 })
    });
  }

  return { server: servers };
}

// Fuzz-Testing Function with Compression
async function fuzzTest(iterations = 1000, compression = 'zlib') {
  let passed = 0;
  let failed = 0;

  for (let i = 0; i < iterations; i++) {
    const config = generateRandomConfig();
    const { error } = serverConfigSchema.validate(config);

    if (error) {
      console.error(`Validation failed on iteration ${i + 1}:`, error.details);
      failed++;
      continue;
    }

    // Compress the config
    try {
      let compressed;
      if (compression === 'zlib') {
        compressed = await compressDataZlib(config);
      } else if (compression === 'zstd') {
        compressed = await compressDataZstd(config);
      }

      // Decompress the config
      let decompressed;
      if (compression === 'zlib') {
        decompressed = await decompressDataZlib(compressed);
      } else if (compression === 'zstd') {
        decompressed = await decompressDataZstd(compressed);
      }

      // Validate decompressed data matches original
      const { error: decompressionError } = serverConfigSchema.validate(decompressed);
      if (decompressionError) {
        console.error(`Decompression validation failed on iteration ${i + 1}:`, decompressionError.details);
        failed++;
      } else {
        passed++;
      }
    } catch (err) {
      console.error(`Compression/Decompression error on iteration ${i + 1}:`, err);
      failed++;
    }
  }

  console.log(`Fuzz Testing Completed with ${compression}: ${passed} passed, ${failed} failed.`);
}

// Run Fuzz Testing for both zlib and zstd
(async () => {
  console.log('Starting Fuzz Testing with zlib...');
  await fuzzTest(1000, 'zlib');
  
  console.log('Starting Fuzz Testing with zstd...');
  await fuzzTest(1000, 'zstd');
})();