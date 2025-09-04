<template>
  <div class="login-page">
    <!-- Fond animé isolé pour le form -->
    <div id="org-form-particles"></div>

    <div class="login-wrapper">
      <div class="login-container">
        <div class="login-card">
          <h2 class="login-title">NetSecure-IQ</h2>
          <h3 class="login-subtitle">Register Your Organization</h3>

          <form @submit.prevent="submitForm" class="login-form">
            <!-- Organization Info -->
            <div class="form-section">
              <h4><i class="fas fa-building"></i> Organization Information</h4>

              <div class="form-group">
                <label>Organization Name</label>
                <input v-model="form.name" placeholder="Organization Name" required />
              </div>

              <div class="form-group">
                <label>VAT Number</label>
                <input v-model="form.vat_number" placeholder="VAT Number" required />
              </div>

              <div class="form-group">
                <label>Address</label>
                <input v-model="form.address" placeholder="Full Address" required />
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label>City</label>
                  <input v-model="form.city" placeholder="City" required />
                </div>
                <div class="form-group">
                  <label>State</label>
                  <input v-model="form.state" placeholder="State/Region" required />
                </div>
                <div class="form-group">
                  <label>ZIP Code</label>
                  <input v-model="form.zip_code" placeholder="ZIP/Postal Code" required />
                </div>
              </div>

              <div class="form-group">
                <label>Contact Email</label>
                <input v-model="form.contact_email" type="email" placeholder="Email" required />
              </div>

              <div class="form-group">
                <label>PEC Email</label>
                <input v-model="form.pec_email" type="email" placeholder="PEC Email" required />
              </div>

              <div class="form-group">
                <label>SDI Code</label>
                <input v-model="form.sdi_code" placeholder="SDI Code" required />
              </div>

              <div class="form-group">
                <label>Contact Phone</label>
                <input v-model="form.contact_phone" type="tel" placeholder="Phone Number" required />
              </div>

              <div class="form-group">
                <label>Personnel Information</label>
                <textarea v-model="form.personnel_info" placeholder="Personnel information (roles, contacts, etc.)" rows="4" required></textarea>
              </div>
            </div>

            <div class="form-actions">
              <button type="submit" :disabled="loading" class="btn-primary">
                {{ loading ? 'Saving...' : 'Save Organization' }}
              </button>
              <button type="button" class="btn-secondary" @click="goToDashboard">
                Cancel
              </button>
            </div>

            <p v-if="message" class="login-message" :class="messageType">{{ message }}</p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { API } from '@/appCore.js'
import {
  ensurePJSDom, loadParticlesScript, defaultConfig,
  safeRender, observeTheme, destroyForId, themeIsDark
} from '@/appCore.js'

const router = useRouter()
const message = ref('')
const loading = ref(false)
const successMessage = ref(false)
const messageType = computed(() => (successMessage.value ? 'success' : 'error'))

const token  = ref(localStorage.getItem('token')   || '')
const userId = ref(localStorage.getItem('user_id') || '')

/* ---- helpers: decode JWT -> user_id si absent ---- */
function b64UrlDecode (str) {
  try {
    str = str.replace(/-/g, '+').replace(/_/g, '/')
    const pad = str.length % 4
    if (pad) str += '='.repeat(4 - pad)
    return atob(str)
  } catch (e) {
    console.debug('[jwt] base64 decode failed', e)
    return ''
  }
}
function parseJwt (t) {
  try { return JSON.parse(b64UrlDecode((t || '').split('.')[1] || '')) }
  catch (e) { console.debug('[jwt] parse failed', e); return null }
}

/* ---- Particles ---- */
const CONTAINER_ID = 'org-form-particles'
let stopObs = () => {}
function renderParticles() { return safeRender(CONTAINER_ID, defaultConfig(themeIsDark())) }

/* ---- Etat du formulaire ---- */
const form = reactive({
  name: '',
  vat_number: '',
  address: '',
  state: '',
  city: '',
  zip_code: '',
  contact_email: '',
  contact_phone: '',
  pec_email: '',
  sdi_code: '',
  personnel_info: ''
})

/* Prefill depuis localStorage (si on revient sur le formulaire) */
function loadLocalProfileIntoForm() {
  try {
    const saved = JSON.parse(localStorage.getItem('organization_profile') || 'null')
    if (saved && typeof saved === 'object') {
      Object.assign(form, {
        name: saved.name || '',
        vat_number: saved.vat_number || '',
        address: saved.address || '',
        state: saved.state || '',
        city: saved.city || '',
        zip_code: saved.zip_code || '',
        contact_email: saved.contact_email || '',
        contact_phone: saved.contact_phone || '',
        pec_email: saved.pec_email || '',
        sdi_code: saved.sdi_code || '',
        personnel_info: saved.personnel_info || ''
      })
    }
  } catch (e) {
    console.debug('[org] failed to prefill from localStorage', e)
  }
}

onMounted(async () => {
  if (!userId.value && token.value) {
    const claims = parseJwt(token.value)
    if (claims?.user_id) {
      userId.value = claims.user_id
      localStorage.setItem('user_id', userId.value)
    }
  }

  loadLocalProfileIntoForm()

  try { await loadParticlesScript() }
  catch (e) { console.debug('[particles] load failed (non-blocking)', e) }
  ensurePJSDom()
  renderParticles()
  stopObs = observeTheme(CONTAINER_ID, renderParticles)
})

onBeforeUnmount(() => {
  stopObs?.()
  ensurePJSDom()
  destroyForId(CONTAINER_ID)
})

const goToDashboard = () => router.push('/dashboard')

/* ---- Soumission ---- */
async function submitForm() {
  if (!token.value) {
    message.value = 'You must be logged in to save organization data.'
    successMessage.value = false
    return
  }
  if (!userId.value) {
    message.value = 'Unable to determine the current user. Please sign in again.'
    successMessage.value = false
    return
  }

  loading.value = true
  message.value = ''

  try {
    // Clés exactement comme attendues par le backend (+ user_id)
    const payload = {
      name: form.name.trim(),
      vat_number: form.vat_number.trim(),
      address: form.address.trim(),
      state: form.state.trim(),
      city: form.city.trim(),
      zip_code: form.zip_code.trim(),
      contact_email: form.contact_email.trim(),
      contact_phone: form.contact_phone.trim(),
      pec_email: form.pec_email.trim(),
      sdi_code: form.sdi_code.trim(),
      personnel_info: form.personnel_info.trim(),
      user_id: userId.value
    }

    const resp = await fetch(`${API}/api/complete-organization`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify(payload)
    })

    const raw = await resp.text()
    let data = {}
    try { data = raw ? JSON.parse(raw) : {} } catch (e) { data = {} }

    if (!resp.ok) {
      const errorMsg = (data && (data.error || data.message)) || raw || 'Failed to save organization'
      throw new Error(errorMsg)
    }

    // Sauvegarde locale pour OrganizationProfile
    const profileToStore = {
      organization_id: data.organization_id || null,
      name: payload.name,
      vat_number: payload.vat_number,
      address: payload.address,
      state: payload.state,
      city: payload.city,
      zip_code: payload.zip_code,
      contact_email: payload.contact_email,
      contact_phone: payload.contact_phone,
      pec_email: payload.pec_email,
      sdi_code: payload.sdi_code,
      personnel_info: payload.personnel_info
    }
    localStorage.setItem('organization_profile', JSON.stringify(profileToStore))

    successMessage.value = true
    message.value = (data.message || 'Organization saved successfully!') +
      (data.organization_id ? ` (ID: ${data.organization_id})` : '')

    // Redirection douce (la page profil lira le localStorage)
    setTimeout(() => { router.push('/dashboard') }, 1200)
  } catch (e) {
    console.error('Error saving organization:', e)
    message.value = e?.message || 'An error occurred while saving the organization.'
    successMessage.value = false
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css');

:root {
  --bg-dark: #0e111a; --panel-grey: #1a1d26; --divider-grey: #2a2d36;
  --text-primary: #f5f7fa; --text-secondary: #9ca3af;
  --primary-accent: #00c2c2; --primary-hover: #00a7a7;
  --danger: #ef4444; --success: #22c55e; --border-radius: 12px; --transition: all .2s;
}

.login-page { position: relative; min-height: 100vh; overflow: hidden; background-color: var(--bg-dark); }

#org-form-particles {
  position: fixed; inset: 0; width: 100vw; height: 100vh; z-index: 0;
  background-color: var(--bg-dark); pointer-events: none; transition: background-color .3s;
}
[data-theme='light'] #org-form-particles { background-color: #f6f8fb; }

.login-wrapper { position: relative; z-index: 10; display: flex; align-items: center; justify-content: center; padding: 32px; min-height: 100vh; }
.login-container { width: 100%; max-width: 800px; }
.login-card { background-color: var(--panel-grey); border-radius: 16px; padding: 32px; box-shadow: 0 0 40px rgba(0, 194, 194, 0.05); border: 1px solid rgba(255,255,255,.05); }

.login-title { text-align: center; font-size: 24px; font-weight: 600; color: var(--primary-accent); margin-bottom: 8px; }
.login-subtitle { text-align: center; font-size: 18px; margin-bottom: 32px; font-weight: 500; color: var(--text-primary); }

.login-form { display: flex; flex-direction: column; gap: 16px; }
.form-section { background-color: rgba(31,41,55,.30); border-radius: 8px; padding: 16px; margin-bottom: 16px; transition: background-color .3s; }
:root:not([data-theme='dark']) .form-section { background-color: rgba(243,244,246,.5); border: 1px solid rgba(209,213,219,.5); }
.form-section h4 { color: var(--primary-accent); margin: 0 0 16px; font-size: 16px; font-weight: 500; display: flex; gap: 8px; align-items: center; }
.login-form input { width: 100%; background: var(--panel-grey); border: 1px solid var(--divider-grey); border-radius: 6px; padding: 12px 14px; font-size: 14px; color: var(--text-primary); transition: var(--transition); margin-bottom: 8px; }
.login-form input:focus { outline: none; border-color: var(--primary-accent); box-shadow: 0 0 0 2px rgba(0,194,194,.2); }
.login-form input::placeholder { color: var(--text-secondary); opacity: .7; }

.form-row { display: flex; gap: 16px; margin-bottom: 8px; }
.form-row input { margin-bottom: 0; }

.form-actions { display: flex; justify-content: space-between; margin-top: 24px; gap: 16px; }
button { flex: 1; padding: 12px 20px; border: none; border-radius: 6px; font-weight: 600; font-size: 14px; cursor: pointer; transition: var(--transition); display:inline-flex; align-items:center; justify-content:center; gap:8px; }
button:disabled { opacity: .6; cursor: not-allowed; }
button:not(:disabled):hover { transform: translateY(-1px); }
button[type='submit'] { background-color: var(--primary-accent); color: #0e111a; }
button[type='submit']:not(:disabled):hover { background-color: var(--primary-hover); }
.btn-secondary { background: transparent; color: var(--primary-accent); border: 1px solid var(--primary-accent)!important; }
.btn-secondary:not(:disabled):hover { background: rgba(0,194,194,.1); }

.login-message { margin-top: 16px; padding: 12px 16px; border-radius: 6px; font-size: 14px; text-align: center; }
.login-message.success { background-color: rgba(34,197,94,.1); color: var(--success); border: 1px solid rgba(34,197,94,.2); }
.login-message.error   { background-color: rgba(239,68,68,.1); color: var(--danger); border: 1px solid rgba(239,68,68,.2); }

@media (max-width: 768px) {
  .login-wrapper { padding: 16px; }
  .login-card { padding: 24px; }
  .form-row { flex-direction: column; gap: 8px; }
  .form-actions { flex-direction: column; }
  button { width: 100%; }
}
</style>
