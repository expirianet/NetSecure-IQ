
<!-- src/frontend/src/components/AddUserForm.vue -->
<template>
  <div class="adduser-page">
    <!-- Particules -->
    <div id="adduser-particles"></div>

    <!-- Carte centrée -->
    <div class="adduser-wrapper">
      <div class="adduser-container">
        <div class="adduser-card">
          <h2 class="adduser-title">NetSecure-IQ</h2>
          <h3 class="adduser-subtitle">Create User</h3>

          <form @submit.prevent="submitForm" class="adduser-form" novalidate>
            <!-- Identité -->
            <div class="form-section">
              <h4><i class="fas fa-id-card"></i> Identity</h4>
              <div class="form-row">
                <div class="field">
                  <label for="firstName">First name</label>
                  <input
                    id="firstName"
                    v-model.trim="firstName"
                    type="text"
                    :class="{ invalid: showErrors && !firstNameValid }"
                    placeholder="Jane"
                    required
                  />
                  <small v-if="showErrors && !firstNameValid">First name is required (min 2 chars).</small>
                </div>

                <div class="field">
                  <label for="lastName">Last name</label>
                  <input
                    id="lastName"
                    v-model.trim="lastName"
                    type="text"
                    :class="{ invalid: showErrors && !lastNameValid }"
                    placeholder="Doe"
                    required
                  />
                  <small v-if="showErrors && !lastNameValid">Last name is required (min 2 chars).</small>
                </div>
              </div>

              <div class="field">
                <label for="email">Email</label>
                <input
                  id="email"
                  v-model.trim="email"
                  type="email"
                  :class="{ invalid: showErrors && !emailValid }"
                  placeholder="user@example.com"
                  required
                />
                <small v-if="showErrors && !emailValid">Please enter a valid email address.</small>
              </div>

              <div class="field">
                <label for="phone">Phone (optional)</label>
                <input id="phone" v-model.trim="phone" type="tel" placeholder="+33 6 12 34 56 78" />
              </div>
            </div>

            <!-- Contexte / Organisation -->
            <div class="form-section">
              <h4><i class="fas fa-building"></i> Organization</h4>

              <!-- Operator : org imposée depuis la session -->
              <div v-if="isOperator" class="readonly-pill" title="Taken from your session">
                <span class="dot"></span>
                Using your organization: <code>{{ userOrgId || '—' }}</code>
              </div>

              <!-- Admin : pas d’endpoint de liste => champ optionnel -->
              <div v-if="isAdmin" class="field">
                <label for="orgIdManual">Organization ID (optional)</label>
                <input
                  id="orgIdManual"
                  v-model.trim="orgIdManual"
                  type="text"
                  placeholder="UUID… (leave empty to create user without org)"
                />
                <small class="hint">No organization list endpoint available. You can paste a UUID or leave empty.</small>
              </div>

              <div class="readonly-pill" title="Role is fixed for this screen">
                <span class="dot cyan"></span>
                Role: <code>User</code>
              </div>
            </div>

            <!-- Sécurité / Options -->
            <div class="form-section">
              <h4><i class="fas fa-shield-alt"></i> Security & Options</h4>

              <div class="options">
                <label class="checkbox">
                  <input type="checkbox" v-model="sendInvite" />
                  <span>Send invitation email</span>
                </label>

                <label class="checkbox">
                  <input type="checkbox" v-model="requireReset" />
                  <span>Require password reset on first login</span>
                </label>

                <label class="checkbox">
                  <input type="checkbox" v-model="active" />
                  <span>Active account</span>
                </label>
              </div>

              <div class="field">
                <label for="tempPass">Temporary password (optional)</label>
                <div class="pass-row">
                  <input
                    id="tempPass"
                    :type="showPassword ? 'text' : 'password'"
                    v-model="tempPassword"
                    placeholder="Leave empty to auto-generate server side"
                  />
                  <button type="button" class="btn-ghost" @click="toggleShowPassword">
                    <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
                  </button>
                  <button type="button" class="btn-ghost" @click="generatePassword">
                    <i class="fas fa-magic"></i> Generate
                  </button>
                </div>

                <div class="strength" v-if="tempPassword">
                  <div class="bar" :style="{ width: strengthPct + '%'}"></div>
                </div>
                <small v-if="tempPassword">Strength: {{ strengthLabel }}</small>
              </div>
            </div>

            <!-- Actions -->
            <button class="submit" type="submit" :disabled="!formValid || loading">
              <span v-if="loading" class="spinner"></span>
              {{ loading ? 'Creating…' : 'Create user' }}
            </button>

            <p v-if="message" class="adduser-message" :class="successMessage ? 'success' : 'error'">
              {{ message }}
            </p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { API } from '@/utils/api.js'
import {
  ensurePJSDom, loadParticlesScript, defaultConfig,
  safeRender, observeTheme, destroyForId, themeIsDark
} from '@/utils/particles.js'

/* ---------- Champs ---------- */
const firstName = ref('')
const lastName  = ref('')
const email     = ref('')
const phone     = ref('')

/* ---------- Options ---------- */
const sendInvite   = ref(true)
const requireReset = ref(true)
const active       = ref(true)
const tempPassword = ref('')
const showPassword = ref(false)

/* ---------- UI ---------- */
const loading        = ref(false)
const message        = ref('')
const successMessage = ref(false)
const showErrors     = ref(false)

/* ---------- Contexte auth ---------- */
const token   = localStorage.getItem('token') || ''
const roleStr = (localStorage.getItem('role') || '').toLowerCase()
const isAdmin    = computed(() => roleStr === 'administrator')
const isOperator = computed(() => roleStr === 'operator')
const userOrgId  = (localStorage.getItem('organization_id') || '').toString()

/* ---------- Admin: org manuelle optionnelle ---------- */
const orgIdManual = ref('')

/* ---------- Validation ---------- */
const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/i
const firstNameValid = computed(() => firstName.value.trim().length >= 2)
const lastNameValid  = computed(() => lastName.value.trim().length >= 2)
const emailValid     = computed(() => emailRe.test(email.value))

// orgValid :
// - operator -> doit avoir une org en session
// - admin    -> champ optionnel (toujours true)
const orgValid = computed(() => (isOperator.value ? !!userOrgId : true))

const formValid = computed(() =>
  firstNameValid.value && lastNameValid.value && emailValid.value && orgValid.value
)

/* ---------- Password strength ---------- */
function scorePassword(pw) {
  let score = 0
  if (!pw) return 0
  if (pw.length >= 8) score++
  if (/[A-Z]/.test(pw) && /[a-z]/.test(pw)) score++
  if (/\d/.test(pw)) score++
  if (/[^A-Za-z0-9]/.test(pw)) score++
  if (pw.length >= 12) score++
  return Math.min(score, 4)
}
const strengthScore = computed(() => scorePassword(tempPassword.value))
const strengthPct   = computed(() => [0, 25, 50, 75, 100][strengthScore.value])
const strengthLabel = computed(() => ['Very weak', 'Weak', 'Medium', 'Strong', 'Very strong'][strengthScore.value])

function generatePassword() {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789!@#$%*?'
  let out = ''
  for (let i = 0; i < 14; i++) out += chars[Math.floor(Math.random() * chars.length)]
  tempPassword.value = out
}
function toggleShowPassword() { showPassword.value = !showPassword.value }

/* ---------- Submit ---------- */
async function submitForm() {
  showErrors.value = true
  if (!formValid.value) return

  // Détermine l’orgId selon le rôle
  const resolvedOrgId = isOperator.value
    ? String(userOrgId || '')
    : String(orgIdManual.value || '')

  const payload = {
    email: email.value,
    first_name: firstName.value,
    last_name: lastName.value,
    phone: phone.value || undefined, // ignoré par le back actuel (sans impact)
    role: 'user',
    ...(resolvedOrgId ? { organization_id: resolvedOrgId } : {}),
    // options purement front pour l’instant (le back les ignore sans erreur)
    send_invite: sendInvite.value,
    require_password_reset: requireReset.value,
    active: active.value,
    temp_password: tempPassword.value || undefined
  }

  loading.value = true
  message.value = ''
  successMessage.value = false

  try {
    const res = await fetch(`${API}/api/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify(payload)
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok) throw new Error(data.error || data.message || 'Request failed')

    successMessage.value = true
    message.value = '✅ User created successfully.'
    // reset soft
    email.value = ''
    firstName.value = ''
    lastName.value = ''
    phone.value = ''
    tempPassword.value = ''
    if (isAdmin.value) orgIdManual.value = ''
    showErrors.value = false
  } catch (err) {
    successMessage.value = false
    message.value = '❌ Failed: ' + (err.message || 'Unknown error')
  } finally {
    loading.value = false
  }
}

/* ---------- Particles ---------- */
const ID = 'adduser-particles'
let stopObs = () => {}
function renderParticles() { return safeRender(ID, defaultConfig(themeIsDark())) }

onMounted(async () => {
  try {
    await loadParticlesScript()
  } catch (e) {
    console.debug('[particles] load failed (non-blocking)', e)
  }
  ensurePJSDom()
  renderParticles()
  stopObs = observeTheme(ID, renderParticles)
})

onBeforeUnmount(() => {
  stopObs?.()
  ensurePJSDom()
  destroyForId(ID)
})
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css');

:root {
  --bg-dark: #0e111a;
  --panel-grey: #1a1d26;
  --divider-grey: #2a2d36;
  --text-primary: #f5f7fa;
  --text-secondary: #9ca3af;
  --primary-accent: #00c2c2;
  --primary-hover: #00a7a7;
  --danger: #ef4444;
  --success: #22c55e;
  --radius: 12px;
}

.adduser-page { position: relative; min-height: 100vh; overflow: hidden; }
#adduser-particles {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  z-index: 0; background-color: var(--bg-dark);
  transition: background-color .3s ease; pointer-events: none;
}
[data-theme='light'] #adduser-particles { background-color: #f6f8fb; }

.adduser-wrapper { position: relative; z-index: 10; display: flex; align-items: center; justify-content: center; padding: 32px; min-height: 100vh; }
.adduser-container { width: 100%; max-width: 720px; }
.adduser-card {
  background-color: var(--panel-grey);
  border-radius: 16px; padding: 28px;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  border: 1px solid rgba(255,255,255,.05);
  box-sizing: border-box;
}

.adduser-title { text-align: center; font-size: 20px; font-weight: 600; color: var(--primary-accent); margin-bottom: 6px; }
.adduser-subtitle { text-align: center; font-size: 16px; margin-bottom: 20px; }

.adduser-form { display: flex; flex-direction: column; gap: 16px; }
.form-section {
  background-color: rgba(31,41,55,.35);
  border-radius: 10px; padding: 16px;
  border: 1px solid rgba(255,255,255,.05);
}
.form-section h4 {
  margin: 0 0 12px; color: var(--primary-accent);
  font-size: 15px; font-weight: 600; display: flex; gap: 8px; align-items: center;
}

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.field { display: flex; flex-direction: column; gap: 6px; }
label { font-size: 13px; color: var(--text-secondary); }
input, select {
  width: 100%; background-color: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 6px; padding: 12px 14px; font-size: 14px; color: var(--text-primary);
  transition: border-color .2s;
}
input::placeholder { color: var(--text-secondary); }
input:focus, select:focus { outline: none; border-color: var(--primary-accent); background-color: var(--bg-dark); }
.invalid { border-color: var(--danger)!important; }
small { color: var(--danger); }
small.hint { color: var(--text-secondary); }

.readonly-pill {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 8px 10px; border-radius: 999px;
  background: rgba(255,255,255,.04); border: 1px solid rgba(255,255,255,.06);
  font-size: 13px; color: var(--text-primary);
}
.readonly-pill .dot { width: 8px; height: 8px; border-radius: 50%; background: #aaa; }
.readonly-pill .dot.cyan { background: var(--primary-accent); }
.readonly-pill code { opacity: .9; }

.options { display: grid; grid-template-columns: repeat(auto-fit,minmax(220px,1fr)); gap: 10px; }
.checkbox { display: inline-flex; align-items: center; gap: 8px; color: var(--text-primary); font-size: 14px; }
.checkbox input { width: 16px; height: 16px; }

.pass-row { display: flex; gap: 8px; }
.btn-ghost {
  background: rgba(255,255,255,.06); border: 1px solid rgba(255,255,255,.08);
  color: var(--text-primary); padding: 10px 12px; border-radius: 6px;
  cursor: pointer; font-weight: 600;
}
.btn-ghost:hover { background: rgba(255,255,255,.10); }
.strength { margin-top: 8px; height: 6px; background: rgba(255,255,255,.08); border-radius: 6px; overflow: hidden; }
.strength .bar { height: 100%; background: linear-gradient(90deg,#ff4d4d,#ffc107,#22c55e); }

.submit {
  background-color: var(--primary-accent); color: var(--bg-dark);
  border: none; border-radius: 8px; font-weight: 700; padding: 12px 20px;
  cursor: pointer; transition: background-color .2s; width: 100%;
}
.submit:hover { background-color: var(--primary-hover); color: #fff; }
.submit:disabled { background-color: #2f333d; color: #666; cursor: not-allowed; }

.spinner {
  display: inline-block; width: 14px; height: 14px; border: 2px solid rgba(0,0,0,.2);
  border-top-color: rgba(0,0,0,.6); border-radius: 50%; margin-right: 8px;
  animation: spin .8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.adduser-message { margin-top: 10px; padding: 10px 12px; border-radius: 6px; text-align: center; font-size: 14px; }
.adduser-message.success { background-color: rgba(34,197,94,.1); color: var(--success); }
.adduser-message.error   { background-color: rgba(239,68,68,.1); color: var(--danger); }

@media (max-width: 760px) { .form-row { grid-template-columns: 1fr; } }
</style>

