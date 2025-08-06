import { defineStore } from 'pinia'
import { useAuthStore } from './authStore'

export const useSiteStore = defineStore('site', {
  state: () => ({
    sites: [],
  }),
  actions: {
    async fetchSites() {
      const auth = useAuthStore()
      let url = '/api/sites'
      // Optionally add org/user filter params based on role
      this.sites = await fetch(url, {
        headers: { Authorization: `Bearer ${auth.token}` }
      }).then(r => r.json())
    }
  },
  getters: {
    filteredSites(state) {
      const auth = useAuthStore()
      if (auth.role === 'Administrator') return state.sites
      if (auth.role === 'Operator') return state.sites.filter(s => s.organization_id === auth.user.organization_id)
      if (auth.role === 'User') return state.sites.filter(s => auth.user.site_ids?.includes(s.id))
      return []
    }
  }
})
