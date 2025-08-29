// src/frontend/src/utils/agentsApi.js
// Client minimal pour les endpoints MikroTik (agents) du backend.

import { API } from '@/utils/api.js'

function authHeaders() {
  const t = localStorage.getItem('token') || ''
  return {
    'Content-Type': 'application/json',
    ...(t ? { Authorization: 'Bearer ' + t } : {})
  }
}

export async function listAgents() {
  const res = await fetch(`${API}/mikrotik/list`, { headers: authHeaders() })
  if (!res.ok) throw new Error('List failed')
  return res.json() // [{ mac, ip, status, site_id }]
}

export async function preregister(mac, siteId = null) {
  const body = siteId ? { mac, site_id: siteId } : { mac }
  const res = await fetch(`${API}/mikrotik/preregister`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify(body)
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || data.message || 'Pre-register failed')
  return data // { internal_ip, private_key, public_key, mikrotik_config, server_peer }
}

export async function testAgent(mac) {
  const res = await fetch(`${API}/mikrotik/test`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ mac })
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || data.message || 'Test failed')
  return data // { success: true, message: 'Ping successful' }
}

export async function disableAgent(mac) {
  const res = await fetch(`${API}/mikrotik/disable`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ mac })
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || data.message || 'Disable failed')
  return data
}

export async function enableAgent(mac) {
  const res = await fetch(`${API}/mikrotik/enable`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ mac })
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || data.message || 'Enable failed')
  return data
}

export async function associateAgent(mac, siteId) {
  const res = await fetch(`${API}/mikrotik/associate`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ mac, site_id: siteId })
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || data.message || 'Association failed')
  return data
}

export async function deleteAgent(mac) {
  const res = await fetch(`${API}/mikrotik/delete`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ mac })
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || data.message || 'Deletion failed')
  return data
}
