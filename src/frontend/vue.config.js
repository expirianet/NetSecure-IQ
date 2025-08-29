// src/frontend/vue.config.js
module.exports = {
  devServer: {
    port: 8080,
    proxy: {
      '^/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        ws: false,
        logLevel: 'debug',
      },
    },
  },
};
