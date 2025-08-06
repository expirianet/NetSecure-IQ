import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', {
  state: () => ({
  isAuthenticated: false,
  user: null, // { email, role_id, organization_id }
  token: ''
    isAuthenticated: false,
    user: null, // { email, role_id, organization_id }
    token: ''
  }),
  getters: {
  isAdmin: (state) => state.user?.role_id === 1,
  isOperator: (state) => state.user?.role_id === 2,
  isUser: (state) => state.user?.role_id === 3,
    isAdmin: (state) => state.user?.role_id === 1,
    isOperator: (state) => state.user?.role_id === 2,
    isUser: (state) => state.user?.role_id === 3,
  },
  actions: {
  logout() {
    this.isAuthenticated = false
    this.user = null
    this.token = ''
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  },
    logout() {
      this.isAuthenticated = false
      this.user = null
      this.token = ''
      localStorage.removeItem('token')
      // Optionally: use router.push('/login') if router is available here
    }
  }
})
