// Petit client fetch qui imite un peu axios ({ data, status })
const RAW_BASE =
  (typeof process !== 'undefined' && process.env && process.env.VUE_APP_BACKEND_URL) ||
  (typeof import !== 'undefined' && typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.VITE_API_URL) ||
  'http://localhost:8081'; // ton backend écoute sur :8081

const BASE = String(RAW_BASE).replace(/\/+$/, '') + '/api';

async function request(method, url, body) {
  const headers = { 'Content-Type': 'application/json' };
  const t = localStorage.getItem('token');
  if (t) headers.Authorization = `Bearer ${t}`;

  const res = await fetch(BASE + url, {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  });

  const text = await res.text();
  let data;
  try { data = text ? JSON.parse(text) : null; } catch { data = text; }

  if (!res.ok) {
    const err = new Error((data && (data.error || data.message)) || res.statusText);
    err.status = res.status; err.data = data;
    throw err;
  }
  return { data, status: res.status, ok: true };
}

const api = {
  baseURL: BASE,
  get: (u) => request('GET', u),
  post: (u, b) => request('POST', u, b),
  put: (u, b) => request('PUT', u, b),
  patch: (u, b) => request('PATCH', u, b),
  delete: (u) => request('DELETE', u),
};

export default api;
