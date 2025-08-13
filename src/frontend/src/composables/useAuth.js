import { ref } from 'vue';

const isAuthenticated = ref(false);

function login() {
  isAuthenticated.value = true;
}

function logout() {
  isAuthenticated.value = false;
  localStorage.removeItem('token');
  localStorage.removeItem('role');
  localStorage.removeItem('organization_id');
  localStorage.removeItem('user_id');
  window.dispatchEvent(new Event('auth-changed'));
  // redirection propre
  window.location.replace('/#/login');
}

export function useAuth() {
  return { isAuthenticated, login, logout };
}
