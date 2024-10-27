// server/middleware/crawlerMiddleware.js

const path = require('path');

  

const crawlerMiddleware = (req, res, next) => {

    const userAgent = req.headers['user-agent'] || '';

    const isCrawler = /Googlebot|Bingbot|Slurp|DuckDuckBot|Baiduspider|YandexBot/.test(userAgent);

  

    if (isCrawler) {

        // Serve the crawler-optimized page

        res.sendFile(path.join(__dirname, '../views/crawler.html'));

    } else {

        // Proceed to serve the regular SPA

        next();

    }


};

  

module.exports = crawlerMiddleware;
