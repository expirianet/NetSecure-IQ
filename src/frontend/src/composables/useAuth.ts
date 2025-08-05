import { ref, computed } from 'vue';

const user = ref<{ email: string } | null>(null);

export function useAuth() {
  function login(email: string) {
    user.value = { email };
  }
  function logout() {
    user.value = null;
  }
  const isAuthenticated = computed(() => !!user.value);
  return { user, isAuthenticated, login, logout };
}
