<!-- src/frontend/src/views/organization/OrganizationForm.vue -->
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

              <!-- Personnel Information (amélioré) -->
              <div class="form-group">
                <div class="pi-head">
                  <label><i class="fas fa-users"></i> Personnel Information</label>

                  <div class="mode-toggle" role="tablist" aria-label="Personnel editor mode">
                    <button
                      type="button"
                      class="pill"
                      :class="{ active: piMode === 'structured' }"
                      @click="switchMode('structured')"
                    >
                      Structured
                    </button>
                    <button
                      type="button"
                      class="pill"
                      :class="{ active: piMode === 'text' }"
                      @click="switchMode('text')"
                    >
                      Free text
                    </button>
                  </div>
                </div>

                <!-- Structured editor -->
                <div v-if="piMode === 'structured'" class="pi-structured">
                  <!-- Inline add row -->
                  <div class="person-inline">
                    <input v-model.trim="draft.name" placeholder="Full name (e.g., Jane Doe)" />
                    <input
                      v-model.trim="draft.role"
                      placeholder="Role"
                      list="role-suggestions"
                    />
                    <input v-model.trim="draft.email" type="email" placeholder="Email" />
                    <input v-model.trim="draft.phone" type="tel" placeholder="Phone" />
                    <button
                      class="icon-btn"
                      type="button"
                      title="Add person"
                      @click="addPerson"
                    >
                      <i class="fas fa-plus"></i>
                    </button>
                  </div>

                  <!-- Suggestions for roles -->
                  <datalist id="role-suggestions">
                    <option>CEO</option>
                    <option>CTO</option>
                    <option>IT Manager</option>
                    <option>Security Officer</option>
                    <option>Network Engineer</option>
                    <option>Helpdesk</option>
                  </datalist>

                  <!-- People list -->
                  <div v-if="people.length" class="people-list">
                    <div class="person-card" v-for="(p, i) in people" :key="i">
                      <div class="pc-main">
                        <span class="name">{{ p.name || '—' }}</span>
                        <span class="role">{{ p.role || '—' }}</span>
                      </div>
                      <div class="pc-meta">
                        <span class="mono">{{ p.email || '—' }}</span>
                        <span class="mono">{{ p.phone || '—' }}</span>
                      </div>
                      <div class="pc-actions">
                        <button
                          type="button"
                          class="icon-btn"
                          title="Edit"
                          @click="startEdit(i)"
                        >
                          <i class="fas fa-pen"></i>
                        </button>
                        <button
                          type="button"
                          class="icon-btn danger"
                          title="Remove"
                          @click="removePerson(i)"
                        >
                          <i class="fas fa-trash"></i>
                        </button>
                      </div>
                    </div>
                  </div>

                  <!-- Edit dialog (inline row) -->
                  <div v-if="editIndex !== -1" class="person-inline edit">
                    <input v-model.trim="editDraft.name" placeholder="Full name" />
                    <input v-model.trim="editDraft.role" placeholder="Role" list="role-suggestions" />
                    <input v-model.trim="editDraft.email" type="email" placeholder="Email" />
                    <input v-model.trim="editDraft.phone" type="tel" placeholder="Phone" />

                    <button class="icon-btn" type="button" title="Save" @click="applyEdit">
                      <i class="fas fa-check"></i>
                    </button>
                    <button class="icon-btn" type="button" title="Cancel" @click="cancelEdit">
                      <i class="fas fa-xmark"></i>
                    </button>
                  </div>

                  <!-- Generated text preview (this is what will be saved) -->
                  <div class="preview">
                    <div class="preview-head">
                      <span class="hint">Preview (stored as text)</span>
                      <button class="icon-btn ghost" type="button" title="Copy" @click="copyPreview">
                        <i class="fas fa-copy"></i>
                      </button>
                    </div>
                    <pre class="pre">{{ generatedPersonnelText }}</pre>
                  </div>
                </div>

                <!-- Free text mode -->
                <div v-else class="pi-text">
                  <textarea
                    v-model="form.personnel_info"
                    placeholder="Personnel information (roles, contacts, on-call, escalation, etc.)"
                    rows="6"
                    @input="autoGrow($event)"
                    ref="textAreaRef"
                  ></textarea>
                  <small class="hint">
                    Tip: One person per line — “Full Name | Role | email | phone”. You can paste here
                    and switch to <b>Structured</b> to parse it automatically.
                  </small>

                  <div class="pi-actions">
                    <button type="button" class="btn-mini" @click="tryParseFromText">
                      <i class="fas fa-wand-magic-sparkles"></i> Parse to structured
                    </button>
                  </div>
                </div>
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
import { ref, computed, onMounted, onBeforeUnmount, reactive, watch } from 'vue'
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
  } catch (e) { console.debug('[jwt] base64 decode failed', e); return '' }
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

/* ==== Personnel Information helpers ==== */
const piMode = ref('structured') // 'structured' | 'text'
const people = ref([])           // [{name, role, email, phone}]
const draft = reactive({ name: '', role: '', email: '', phone: '' })
const editIndex = ref(-1)
const editDraft = reactive({ name: '', role: '', email: '', phone: '' })
const textAreaRef = ref(null)

function switchMode(mode) {
  piMode.value = mode
  if (mode === 'structured') {
    // si on vient du mode texte, essaie de parser
    tryParseFromText()
  }
}

/* Generate the text saved into form.personnel_info from people[] */
const generatedPersonnelText = computed(() => {
  if (!people.value.length) return (form.personnel_info || '').trim()
  const lines = ['Contacts:']
  for (const p of people.value) {
    const parts = []
    if (p.name) parts.push(p.name)
    if (p.role) parts.push(`(${p.role})`)
    const head = parts.join(' ')
    const tail = [p.email, p.phone].filter(Boolean).join(', ')
    lines.push(`- ${head}${tail ? ` — ${tail}` : ''}`)
  }
  return lines.join('\n')
})

/* Keep form.personnel_info in sync while in structured mode */
watch([people, piMode], () => {
  if (piMode.value === 'structured') {
    form.personnel_info = generatedPersonnelText.value
  }
}, { deep: true })

function addPerson () {
  if (!draft.name && !draft.role && !draft.email && !draft.phone) return
  people.value.push({ name: draft.name, role: draft.role, email: draft.email, phone: draft.phone })
  draft.name = draft.role = draft.email = draft.phone = ''
}
function removePerson (i) { people.value.splice(i, 1) }
function startEdit (i) { editIndex.value = i; Object.assign(editDraft, people.value[i]) }
function cancelEdit () { editIndex.value = -1; Object.assign(editDraft, { name:'', role:'', email:'', phone:'' }) }
function applyEdit () {
  if (editIndex.value < 0) return
  people.value[editIndex.value] = { ...editDraft }
  cancelEdit()
}

/* Parse free text (best effort): "name | role | email | phone" OR "- Name (Role) — email, phone" */
function tryParseFromText () {
  const txt = (form.personnel_info || '').trim()
  if (!txt) return
  const out = []
  const lines = txt.split(/\r?\n/).map(l => l.trim()).filter(Boolean)
  for (const line of lines) {
    if (/^contacts?:/i.test(line)) continue
    // pipe-separated
    if (line.includes('|')) {
      const [name, role, email, phone] = line.split('|').map(s => s.trim())
      out.push({ name, role, email, phone })
      continue
    }
    // bullet with (role) — email, phone
    const m = line.match(/^-?\s*(.+?)\s*(?:\((.*?)\))?(?:\s*[—-]\s*(.*))?$/)
    if (m) {
      const name = (m[1] || '').trim()
      const role = (m[2] || '').trim()
      const rest = (m[3] || '').trim()
      let email = '', phone = ''
      if (rest) {
        const parts = rest.split(/\s*,\s*/g)
        for (const p of parts) {
          if (!email && /\S+@\S+\.\S+/.test(p)) email = p
          else if (!phone && /[\d+\s().-]{6,}/.test(p)) phone = p
        }
      }
      if (name || role || email || phone) out.push({ name, role, email, phone })
    }
  }
  if (out.length) people.value = out
}

/* Copy preview text */
function copyPreview () {
  navigator.clipboard.writeText(generatedPersonnelText.value).catch(() => {})
}

/* Auto-grow textarea height */
function autoGrow (e) {
  const el = e?.target || textAreaRef.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = Math.min(el.scrollHeight, 520) + 'px'
}

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
      // try to parse personnel to people for structured mode
      tryParseFromText()
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

  // Keep payload stable: always send 'personnel_info' text (from structured or text)
  if (piMode.value === 'structured') {
    form.personnel_info = generatedPersonnelText.value
  }

  loading.value = true
  message.value = ''

  try {
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
  --danger: #ef4444; --success: #22c55e;
  --border-radius: 12px; --transition: all .2s;
}

.login-page { position: relative; min-height: 100vh; overflow: hidden; background-color: var(--bg-dark); }

#org-form-particles {
  position: fixed; inset: 0; width: 100vw; height: 100vh; z-index: 0;
  background-color: var(--bg-dark); pointer-events: none; transition: background-color .3s;
}
[data-theme='light'] #org-form-particles { background-color: #f6f8fb; }

.login-wrapper { position: relative; z-index: 10; display: flex; align-items: center; justify-content: center; padding: 32px; min-height: 100vh; }
.login-container { width: 100%; max-width: 840px; }
.login-card { background-color: var(--panel-grey); border-radius: 16px; padding: 32px; box-shadow: 0 0 40px rgba(0, 194, 194, 0.05); border: 1px solid rgba(255,255,255,.05); }

.login-title { text-align: center; font-size: 24px; font-weight: 600; color: var(--primary-accent); margin-bottom: 8px; }
.login-subtitle { text-align: center; font-size: 18px; margin-bottom: 28px; font-weight: 500; color: var(--text-primary); }

.login-form { display: flex; flex-direction: column; gap: 16px; }
.form-section { background-color: rgba(31,41,55,.30); border-radius: 8px; padding: 16px; margin-bottom: 16px; transition: background-color .3s; }
.form-section h4 { color: var(--primary-accent); margin: 0 0 16px; font-size: 16px; font-weight: 500; display: flex; gap: 8px; align-items: center; }

.login-form input,
.login-form textarea {
  width: 100%;
  background: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 6px; padding: 12px 14px; font-size: 14px; color: var(--text-primary);
  transition: var(--transition);
}
.login-form textarea { resize: none; }
.login-form input:focus, .login-form textarea:focus {
  outline: none; border-color: var(--primary-accent); box-shadow: 0 0 0 2px rgba(0,194,194,.18);
}
.login-form input::placeholder, .login-form textarea::placeholder { color: var(--text-secondary); opacity: .7; }

.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; }

.form-actions { display: flex; justify-content: space-between; margin-top: 22px; gap: 16px; }
button { flex: 1; padding: 12px 20px; border: none; border-radius: 8px; font-weight: 700; font-size: 14px; cursor: pointer; transition: var(--transition); display:inline-flex; align-items:center; justify-content:center; gap:8px; }
button:disabled { opacity: .6; cursor: not-allowed; }
button:not(:disabled):hover { transform: translateY(-1px); }
button[type='submit'] { background-color: var(--primary-accent); color: #0e111a; }
button[type='submit']:not(:disabled):hover { background-color: var(--primary-hover); }
.btn-secondary { background: transparent; color: var(--primary-accent); border: 1px solid var(--primary-accent)!important; }
.btn-secondary:not(:disabled):hover { background: rgba(0,194,194,.1); }

/* Personnel information styles */
.pi-head { display:flex; align-items:center; justify-content:space-between; gap:12px; }
.mode-toggle { display:flex; gap:8px; }
.pill {
  background: rgba(255,255,255,.06); color: var(--text-primary);
  border: 1px solid rgba(255,255,255,.1); padding: 6px 10px; border-radius: 999px;
  font-weight: 700; font-size: 12px; cursor: pointer;
}
.pill.active { background: var(--primary-accent); color: #0e111a; border-color: transparent; }
.pill:not(.active):hover { background: rgba(255,255,255,.10); }

.pi-structured { display: grid; gap: 12px; }
.person-inline{
  display: grid; grid-template-columns: 1.3fr 1fr 1.2fr 1fr auto;
  gap: 10px; align-items: center;
}
.person-inline input{ margin: 0; }
.person-inline.edit{ background: rgba(255,255,255,.04); border:1px solid rgba(255,255,255,.06); padding:10px; border-radius:8px; }

.icon-btn{
  padding: 10px; border-radius: 8px;
  background: rgba(255,255,255,.06); border: 1px solid rgba(255,255,255,.08);
  color: var(--text-primary); cursor:pointer; font-weight: 700;
}
.icon-btn:hover{ background: rgba(255,255,255,.10); }
.icon-btn.ghost{ background: rgba(255,255,255,.04); }
.icon-btn.danger{ border-color: rgba(239,68,68,.35); color: #f87171; }
.icon-btn.danger:hover{ background: rgba(239,68,68,.15); }

.people-list{ display:grid; gap:10px; }
.person-card{
  display:grid; grid-template-columns: 1fr 1fr auto;
  gap:10px; align-items:center;
  background: rgba(12,16,24,.35); border:1px solid rgba(255,255,255,.06);
  border-radius: 10px; padding: 10px 12px;
}
.pc-main{ display:flex; gap:10px; align-items:center; }
.pc-main .name{ font-weight:700; color: var(--text-primary); }
.pc-main .role{ color: var(--text-secondary); }
.pc-meta{ display:flex; gap:14px; color: var(--text-secondary); }
.pc-actions{ display:flex; gap:8px; }

.preview{ display:grid; gap:6px; }
.preview-head{ display:flex; align-items:center; justify-content:space-between; }
.pre{
  white-space: pre-wrap; background:#0b0e16; color:#e5e7eb;
  border:1px solid var(--divider-grey); border-radius:8px; padding:10px 12px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 13px;
}

.pi-text textarea{ min-height: 140px; }
.pi-actions{ margin-top:8px; display:flex; gap:8px; }
.btn-mini{
  background: rgba(255,255,255,.06); border:1px solid rgba(255,255,255,.08);
  color: var(--text-primary); padding: 8px 10px; border-radius: 8px; font-weight: 700; cursor: pointer;
}
.btn-mini:hover{ background: rgba(255,255,255,.10); }

.mono{ font-family: ui-monospace, SFMono-Regular, Menlo, monospace; letter-spacing:.2px; }

.login-message { margin-top: 16px; padding: 12px 16px; border-radius: 6px; font-size: 14px; text-align: center; }
.login-message.success { background-color: rgba(34,197,94,.1); color: var(--success); border: 1px solid rgba(34,197,94,.2); }
.login-message.error   { background-color: rgba(239,68,68,.1); color: var(--danger); border: 1px solid rgba(239,68,68,.2); }

@media (max-width: 980px) {
  .person-inline{ grid-template-columns: 1fr 1fr 1fr 1fr auto; }
}
@media (max-width: 760px) {
  .login-wrapper { padding: 16px; }
  .login-card { padding: 24px; }
  .form-row { grid-template-columns: 1fr; }
  .form-actions { flex-direction: column; }
  button { width: 100%; }
  .person-inline{ grid-template-columns: 1fr 1fr; }
  .pc-meta{ flex-direction: column; gap:4px; }
}
</style>
