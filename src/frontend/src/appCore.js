// src/frontend/src/appCore.js
// Fusion de tous les utilitaires, stores, composables, helpers, API

// ---- API base & HTTP helper ----
export const API = (import.meta?.env?.VITE_API_BASE || '').replace(/\/+$/, '') || 'http://localhost:8000';
const defaultHeaders = { 'Content-Type': 'application/json' };
export const authHeaders = () => {
  const t = localStorage.getItem('token');
  return t ? { Authorization: `Bearer ${t}` } : {};
};
async function http(method, path, body, extraHeaders = {}) {
  const res = await fetch(`${API}${path}`, {
    method,
    headers: { ...defaultHeaders, ...authHeaders(), ...extraHeaders },
    body: body != null ? JSON.stringify(body) : undefined,
  });
  let data = null;
  try { data = await res.json(); } catch { /* no json body */ }
  if (!res.ok) {
    const msg = data?.error || data?.message || `HTTP ${res.status}`;
    throw new Error(msg);
  }
  return data;
}
export const api = {
  get: (p, h) => http('GET', p, undefined, h),
  post: (p, b, h) => http('POST', p, b, h),
  del: (p, b, h) => http('DELETE', p, b, h),
};
// ---- Domain APIs ----
export const authApi = {
  login: (email, password) => api.post('/api/login', { email, password }),
  register: (email) => api.post('/api/register', { email }),
  createUser: (payload) => api.post('/api/users', payload),
  completeOrganization: (payload) => api.post('/api/complete-organization', payload),
  async getCurrentUser() {
    const token = localStorage.getItem('token');
    if (!token) return null;
    const response = await fetch(`${API}/me`, { headers: { 'Authorization': `Bearer ${token}` } });
    if (!response.ok) throw new Error('Failed to fetch user');
    return response.json();
  }
};
export const agentsApi = {
  list: () => api.get('/api/mikrotik/list'),
  preregister: (mac, siteId = null) => api.post('/api/mikrotik/preregister', siteId ? { mac, site_id: siteId } : { mac }),
  enable: (mac) => api.post('/api/mikrotik/enable', { mac }),
  disable: (mac) => api.post('/api/mikrotik/disable', { mac }),
  remove: (mac) => api.del('/api/mikrotik', { mac }),
  test: (mac) => api.post('/api/mikrotik/test', { mac }),
  associate: (mac, siteId) => api.post('/api/mikrotik/associate', { mac, site_id: siteId }),
  delete: (mac) => api.post('/api/mikrotik/delete', { mac }),
};

// ---- Auth composable ----
import { ref, computed } from 'vue';
export const isAuthenticated = ref(!!localStorage.getItem('token'));
export function login() { isAuthenticated.value = true; }
export function logout() {
  isAuthenticated.value = false;
  localStorage.removeItem('token');
  localStorage.removeItem('role');
  localStorage.removeItem('organization_id');
  localStorage.removeItem('user_id');
  window.dispatchEvent(new Event('auth-changed'));
  window.location.replace('/#/login');
}
export function useAuth() {
  return { isAuthenticated, login, logout };
}

// ---- Pinia Auth Store ----
import { defineStore } from 'pinia';
export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '');
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'));
  const isAuthenticatedStore = computed(() => !!token.value);
  const role = computed(() => user.value?.role || null);
  function setAuth(newToken, newUser) {
    token.value = newToken;
    user.value = newUser;
    localStorage.setItem('token', newToken);
    localStorage.setItem('user', JSON.stringify(newUser));
  }
  function logoutStore() {
    token.value = '';
    user.value = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }
  return { token, user, isAuthenticated: isAuthenticatedStore, role, setAuth, logout: logoutStore };
});

// ---- Pinia Site Store ----
export const useSiteStore = defineStore('site', {
  state: () => ({ sites: [] }),
  actions: {
    async fetchSites() {
      const auth = useAuthStore();
      let url = '/api/sites';
      this.sites = await fetch(url, { headers: { Authorization: `Bearer ${auth.token}` } }).then(r => r.json());
    }
  },
  getters: {
    filteredSites(state) {
      const auth = useAuthStore();
      if (auth.role === 'Administrator') return state.sites;
      if (auth.role === 'Operator') return state.sites.filter(s => s.organization_id === auth.user.organization_id);
      if (auth.role === 'User') return state.sites.filter(s => auth.user.site_ids?.includes(s.id));
      return [];
    }
  }
});

// ---- ACL Helper ----
export function can(user, needed) {
  if (!user) return false;
  const perms = new Set((user.permissions || []).map(String));
  const list = Array.isArray(needed) ? needed : (needed ? [needed] : []);
  return list.every(p => perms.has(String(p)));
}
export default can;

// ---- Particles Utils ----
let _scriptPromise = null;
export function loadParticlesScript () {
  if (typeof window !== 'undefined' && window.particlesJS) {
    return Promise.resolve();
  }
  if (_scriptPromise) return _scriptPromise;
  _scriptPromise = new Promise((resolve, reject) => {
    try {
      const s = document.createElement('script');
      s.src = 'https://cdn.jsdelivr.net/npm/particles.js@2.0.0/particles.min.js';
      s.async = true;
      s.onload = () => resolve();
      s.onerror = () => reject(new Error('Failed to load particles.js'));
      document.head.appendChild(s);
    } catch (e) {
      console.debug('[particles] creating script tag failed', e);
      reject(e);
    }
  }).catch((e) => {
    console.debug('[particles] script load failed', e);
    throw e;
  });
  return _scriptPromise;
}
export function ensurePJSDom () {
  try {
    if (typeof window !== 'undefined' && !window.pJSDom) {
      window.pJSDom = [];
    }
  } catch (e) { console.debug('[particles] ensurePJSDom error', e); }
}
export function destroyForId (id) {
  try {
    const list = (typeof window !== 'undefined' && window.pJSDom) ? window.pJSDom : [];
    for (let i = list.length - 1; i >= 0; i--) {
      try {
        const inst = list[i];
        const el = inst?.pJS?.canvas?.el;
        const hostId = (el && el.parentElement && el.parentElement.id) || inst?.tag?.id || null;
        if (hostId === id) {
          const vendors = inst?.pJS?.fn?.vendors;
          if (vendors?.destroypJS) vendors.destroypJS();
          if (vendors?.destroy) vendors.destroy();
          list.splice(i, 1);
        }
      } catch (e) { console.debug('[particles] single destroy failed', e); }
    }
  } catch (e) { console.debug('[particles] destroyForId failed', e); }
}
export function safeRender (id, config) {
  try {
    const el = typeof document !== 'undefined' ? document.getElementById(id) : null;
    if (!el) return false;
    if (typeof window === 'undefined' || !window.particlesJS) return false;
    destroyForId(id);
    window.particlesJS(id, config);
    return true;
  } catch (e) { console.debug('[particles] render failed', e); return false; }
}
export function observeTheme (id, renderFn) {
  const target = typeof document !== 'undefined' ? document.documentElement : null;
  if (!target) return () => {};
  const observer = new MutationObserver(() => {
    try {
      destroyForId(id);
      renderFn?.();
    } catch (e) { console.debug('[particles] theme change error', e); }
  });
  try { observer.observe(target, { attributes: true, attributeFilter: ['data-theme'] }); } catch (e) { console.debug('[particles] observer start failed', e); }
  return () => { try { observer.disconnect(); } catch (e) { console.debug('[particles] observer cleanup failed', e); } };
}
export function themeIsDark () {
  try {
    const t = typeof document !== 'undefined' ? document.documentElement.getAttribute('data-theme') : null;
    return t !== 'light';
  } catch (e) { console.debug('[particles] themeIsDark error', e); return true; }
}
export function defaultConfig (isDark = true) {
  const baseColor = isDark ? '#00c2c2' : '#0ea5a5';
  const linkColor = isDark ? '#5eead4' : '#14b8a6';
  const bgColor = isDark ? '#0e111a' : '#f6f8fb';
  return {
    particles: {
      number: { value: 60, density: { enable: true, value_area: 800 } },
      color: { value: baseColor },
      shape: { type: 'circle' },
      opacity: { value: 0.4, random: false },
      size: { value: 3, random: true },
      line_linked: { enable: true, distance: 130, color: linkColor, opacity: 0.25, width: 1 },
      move: { enable: true, speed: 1.2, direction: 'none', out_mode: 'out' }
    },
    interactivity: {
      detect_on: 'canvas',
      events: { onhover: { enable: false }, onclick: { enable: false }, resize: true }
    },
    retina_detect: true,
    background: { color: bgColor }
  };
}
