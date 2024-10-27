// stats.js

/**
 * Calculate the mean of a vector
 * @param {number[]} vector
 * @returns {number}
 */
function calculateMean(vector) {
  const sum = vector.reduce((acc, val) => acc + val, 0);
  return sum / vector.length;
}

/**
* Mean-center the vector
* @param {number[]} vector
* @returns {number[]}
*/
function meanCenter(vector) {
  const mean = calculateMean(vector);
  return vector.map(element => element - mean);
}

/**
* Calculate Root Mean Square (RMS) of a vector
* @param {number[]} vector
* @returns {number}
*/
function calculateRMS(vector) {
  const sumOfSquares = vector.reduce((acc, val) => acc + val * val, 0);
  const meanOfSquares = sumOfSquares / vector.length;
  return Math.sqrt(meanOfSquares);
}

/**
* Calculate Population Standard Deviation
* @param {number[]} vector
* @returns {number}
*/
function calculatePopulationSD(vector) {
  const meanCentered = meanCenter(vector);
  return calculateRMS(meanCentered);
}

/**
* Calculate Sample Standard Deviation
* @param {number[]} vector
* @returns {number}
*/
function calculateSampleSD(vector) {
  const mean = calculateMean(vector);
  const sumOfSquares = vector.reduce((acc, val) => acc + Math.pow(val - mean, 2), 0);
  const variance = sumOfSquares / (vector.length - 1);
  return Math.sqrt(variance);
}

/**
* Verify SD properties
* @param {number[]} vector
* @param {number} b
* @param {number} m
*/
function verifySDProperties(vector, b, m) {
  // SD(x + b) = SD(x)
  const transformedVector1 = vector.map(x => x + b);
  const sdOriginal = calculatePopulationSD(vector);
  const sdTransformed1 = calculatePopulationSD(transformedVector1);
  console.log(`SD(x) = ${sdOriginal}`);
  console.log(`SD(x + ${b}) = ${sdTransformed1}`);
  console.log(`Property SD(x + b) = SD(x) holds: ${sdOriginal === sdTransformed1}`);

  // SD(m * x) = |m| * SD(x)
  const transformedVector2 = vector.map(x => m * x);
  const sdTransformed2 = calculatePopulationSD(transformedVector2);
  const expectedSD = Math.abs(m) * sdOriginal;
  console.log(`SD(m * x) = ${sdTransformed2}`);
  console.log(`Expected SD = ${expectedSD}`);
  console.log(`Property SD(m * x) = |m| * SD(x) holds: ${sdTransformed2 === expectedSD}`);
}

// Export functions for use in other modules
export default {
  calculateMean,
  meanCenter,
  calculateRMS,
  calculatePopulationSD,
  calculateSampleSD,
  verifySDProperties
};