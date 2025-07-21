const { defineConfig } = require('@vue/cli-service');
const path = require('path');
const fs = require('fs');
const dotenv = require('dotenv');

// Manually load environment variables from custom location
const envPath = path.resolve(__dirname, '../../config/frontend_env/frontend.env');
if (fs.existsSync(envPath)) {
  dotenv.config({ path: envPath });
  console.log('[INFO] Loaded frontend.env from config/frontend_env/');
} else {
  console.warn('[WARNING] frontend.env file not found at config/frontend_env/');
}

module.exports = defineConfig({
  transpileDependencies: true
});
