const API_HOST = process.env.API_HOST || '127.0.0.1';

module.exports = {
  NODE_ENV: '"development"',
  ENV_CONFIG: '"dev"',
  BASE_API: `"http://${API_HOST}:2000/api/"`
}
// BASE_API: `"http://${API_HOST}:2000/api/"`
// BASE_API: '"http://127.0.0.1:3000"'
