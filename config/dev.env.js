const API_HOST = process.env.API_HOST || '127.0.0.1';
const ENV_CONFIG = process.env.ENV_CONFIG || 'dev';

module.exports = {
  NODE_ENV: '"development"',
  ENV_CONFIG: `"${ENV_CONFIG}"`,
  BASE_API: `"http://${API_HOST}:2000/api/"`
}
// BASE_API: `"http://${API_HOST}:2000/api/"`
// BASE_API: '"http://127.0.0.1:3000"'
