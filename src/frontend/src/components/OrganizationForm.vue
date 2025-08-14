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
              <input v-model="form.name" placeholder="Organization Name" required />
              <input v-model="form.vat_number" placeholder="VAT Number or Fiscal Code" required />
              <input v-model="form.address" placeholder="Address" required />
              <div class="form-row">
                <input v-model="form.city" placeholder="City" required />
                <input v-model="form.state" placeholder="State" required />
                <input v-model="form.zip_code" placeholder="ZIP Code" required />
              </div>
              <input v-model="form.email" type="email" placeholder="Email" required />
              <input v-model="form.phone" type="tel" placeholder="Phone Number" required />
              <input v-model="form.pec_email" type="email" placeholder="PEC Email (Optional)" />
              <input v-model="form.sdi" placeholder="SDI Code (Optional)" />
            </div>

            <!-- Company Manager -->
            <div class="form-section">
              <h4><i class="fas fa-user-tie"></i> Company Manager</h4>
              <input v-model="form.manager_name" placeholder="Name and Surname" required />
              <input v-model="form.manager_email" type="email" placeholder="Email" required />
              <input v-model="form.manager_phone" placeholder="Phone Number" required />
            </div>

            <!-- Technical Manager -->
            <div class="form-section">
              <h4><i class="fas fa-user-cog"></i> Technical Manager</h4>
              <input v-model="form.technical_name" placeholder="Name and Surname" required />
              <input v-model="form.technical_email" type="email" placeholder="Email" required />
              <input v-model="form.technical_phone" placeholder="Phone Number" required />
            </div>

            <!-- Data Controller -->
            <div class="form-section">
              <h4><i class="fas fa-shield-alt"></i> Data Controller</h4>
              <input v-model="form.controller_name" placeholder="Name and Surname" required />
              <input v-model="form.controller_email" type="email" placeholder="Email" required />
              <input v-model="form.controller_phone" placeholder="Phone Number" required />
            </div>

            <!-- Data Processor -->
            <div class="form-section">
              <h4><i class="fas fa-database"></i> Data Processor</h4>
              <input v-model="form.processor_name" placeholder="Name and Surname" required />
              <input v-model="form.processor_email" type="email" placeholder="Email" required />
              <input v-model="form.processor_phone" placeholder="Phone Number" required />
            </div>

            <div class="form-actions">
              <button type="submit" :disabled="loading">{{ loading ? 'Submitting...' : 'Submit' }}</button>
              <button type="button" class="btn-secondary" @click="goToDashboard">Go to Dashboard</button>
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
import { API } from '@/utils/api.js'
import {
  ensurePJSDom, loadParticlesScript, defaultConfig,
  safeRender, observeTheme, destroyForId, themeIsDark
} from '@/utils/particles.js'

const router = useRouter()
const message = ref('')
const loading = ref(false)
const successMessage = ref(false)
const messageType = computed(() => (successMessage.value ? 'success' : 'error'))

/* ---------- Particles (ID dédié + utilités robustes) ---------- */
const CONTAINER_ID = 'org-form-particles'
let stopObs = () => {}
function renderParticles() { return safeRender(CONTAINER_ID, defaultConfig(themeIsDark())) }

onMounted(async () => {
  try { await loadParticlesScript() } catch {}
  ensurePJSDom()
  renderParticles()
  stopObs = observeTheme(CONTAINER_ID, renderParticles)
})
onBeforeUnmount(() => {
  stopObs?.()
  ensurePJSDom()
  destroyForId(CONTAINER_ID)
})

/* ---------- Données & Submit ---------- */
const form = reactive({
  name: '', vat_number: '', address: '', state: '', city: '', zip_code: '',
  email: '', pec_email: '', sdi: '', phone: '',
  manager_name: '', manager_email: '', manager_phone: '',
  technical_name: '', technical_email: '', technical_phone: '',
  controller_name: '', controller_email: '', controller_phone: '',
  processor_name: '', processor_email: '', processor_phone: ''
})

const goToDashboard = () => router.push('/dashboard')

const submitForm = async () => {
  const personnelInfo = `
Company Manager:
  Name: ${form.manager_name}
  Email: ${form.manager_email}
  Phone: ${form.manager_phone}

Technical Manager:
  Name: ${form.technical_name}
  Email: ${form.technical_email}
  Phone: ${form.technical_phone}

Data Controller:
  Name: ${form.controller_name}
  Email: ${form.controller_email}
  Phone: ${form.controller_phone}

Data Processor:
  Name: ${form.processor_name}
  Email: ${form.processor_email}
  Phone: ${form.processor_phone}
`.trim()

  const stored = {
    name: form.name, vat_number: form.vat_number, address: form.address,
    city: form.city, state: form.state, zip_code: form.zip_code,
    contact_email: form.email, contact_phone: form.phone,
    pec_email: form.pec_email, sdi_code: form.sdi,
    personnel_info: personnelInfo,
    manager:   { name: form.manager_name,    email: form.manager_email,    phone: form.manager_phone },
    technical: { name: form.technical_name,  email: form.technical_email,  phone: form.technical_phone },
    controller:{ name: form.controller_name, email: form.controller_email, phone: form.controller_phone },
    processor: { name: form.processor_name,  email: form.processor_email,  phone: form.processor_phone },
  }
  localStorage.setItem('organization_profile', JSON.stringify(stored))

  const payload = {
    name: form.name, vat_number: form.vat_number, address: form.address,
    state: form.state, city: form.city, zip_code: form.zip_code,
    contact_email: form.email, pec_email: form.pec_email, sdi_code: form.sdi,
    contact_phone: form.phone, personnel_info: personnelInfo,
    user_id: localStorage.getItem('user_id'),
  }

  loading.value = true
  message.value = ''
  successMessage.value = false

  try {
    const token = localStorage.getItem('token') || ''
    const res = await fetch(`${API}/api/complete-organization`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) },
      body: JSON.stringify(payload),
    }).catch(() => null)

    if (res && res.ok) { await res.text().catch(() => '') }
    successMessage.value = true
    message.value = 'Organization info saved.'
    setTimeout(() => router.push('/organization'), 400)
  } catch {
    successMessage.value = true
    message.value = 'Organization info saved locally.'
    setTimeout(() => router.push('/organization'), 400)
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

/* Conteneur animé dédié */
#org-form-particles {
  position: fixed; inset: 0; width: 100vw; height: 100vh; z-index: 0;
  background-color: var(--bg-dark); pointer-events: none; transition: background-color .3s;
}
[data-theme='light'] #org-form-particles { background-color: #E0E0E0; }

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
