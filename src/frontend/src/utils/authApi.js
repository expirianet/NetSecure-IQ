// src/frontend/src/utils/authApi.js
import { API } from './api'

export const authApi = {
  async login(email, password) {
    console.log('Making login request to:', `${API}/login`);
    const response = await fetch(`${API}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password })
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Login failed')
    }
    
    return response.json()
  },
  
  async logout() {
    // Add logout logic if needed
  },
  
  async getCurrentUser() {
    const token = localStorage.getItem('token')
    if (!token) return null
    
    const response = await fetch(`${API}/me`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch user')
    }
    
    return response.json()
  }
}
