const { start } = require('repl');
const { Worker } = require('worker_threads');

function startCrawler() {
  const worker = new Worker('./crawler.js', {
    workerData: {
      url: 'https://www.google.com/',
    },
  });
}

(async () => {
  startCrawler();
})();
module.exports = { startCrawler };