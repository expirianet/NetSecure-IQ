import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { API } from '@/appCore.js';

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter();
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'));
  const token = ref(localStorage.getItem('token') || null);
  const returnUrl = ref(null);

  const isAuthenticated = computed(() => !!token.value);

  async function login(email, password) {
    try {
      const response = await fetch(`${API}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Login failed');
      }

      // Mise à jour de l'état
      user.value = data.user;
      token.value = data.token;

      // Stockage dans le localStorage
      localStorage.setItem('user', JSON.stringify(data.user));
      localStorage.setItem('token', data.token);

      // Redirection vers la page précédente ou la page d'accueil
      router.push(returnUrl.value || '/dashboard');
      return data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  function logout() {
    // Suppression des données d'authentification
    user.value = null;
    token.value = null;
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    
    // Redirection vers la page de connexion
    router.push('/login');
  }

  function setReturnUrl(url) {
    returnUrl.value = url;
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    logout,
    setReturnUrl,
  };
});
