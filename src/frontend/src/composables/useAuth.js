import { ref } from 'vue';

const isAuthenticated = ref(false);

function login() {
  isAuthenticated.value = true;
}

function logout() {
  isAuthenticated.value = false;
}

export function useAuth() {
  return {
    isAuthenticated,
    login,
    logout,
  };
}
