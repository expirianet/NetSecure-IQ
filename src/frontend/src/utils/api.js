// Single source of truth for the backend base URL.
export const API =
  (typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.VITE_API_URL) ||
  (typeof process !== 'undefined' && process.env && process.env.VUE_APP_BACKEND_URL) ||
  'http://localhost:8081';