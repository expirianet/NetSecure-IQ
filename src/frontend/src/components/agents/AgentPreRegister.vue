<template>
  <div class="prereg-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <div class="card">
          <h2 class="title">Agent Pre-Registration</h2>
          <p class="muted">Generate WireGuard keys & MikroTik script (server side).</p>

          <form class="form" @submit.prevent="submit">
            <div class="row">
              <label>MAC Address</label>
              <div class="mac-row">
                <input
                  class="input"
                  v-model="mac"
                  placeholder="AA:BB:CC:DD:EE:FF"
                  @input="formatMac"
                  required
                />
                <button class="btn" type="submit" :disabled="loading || !macValid">
                  <i class="fas" :class="loading ? 'fa-spinner fa-spin' : 'fa-wand-magic-sparkles'"></i>
                  <span>{{ loading ? 'Generating…' : 'Generate' }}</span>
                </button>
              </div>
              <small class="hint">MAC will be normalized to upper-case with separators.</small>
            </div>

            <div class="row">
              <label>Site ID (optional)</label>
              <input class="input" v-model="siteId" placeholder="site uuid…" />
            </div>
          </form>

          <div v-if="result" class="result">
            <div class="grid">
              <div class="box">
                <div class="box-head">
                  <h3>Assigned IP</h3>
                  <button class="mini" @click="copy(result.internal_ip)">Copy</button>
                </div>
                <pre class="mono">{{ result.internal_ip }}</pre>
              </div>

              <div class="box">
                <div class="box-head">
                  <h3>Private Key</h3>
                  <button class="mini" @click="copy(result.private_key)">Copy</button>
                </div>
                <pre class="mono">{{ result.private_key }}</pre>
              </div>

              <div class="box">
                <div class="box-head">
                  <h3>Public Key</h3>
                  <button class="mini" @click="copy(result.public_key)">Copy</button>
                </div>
                <pre class="mono">{{ result.public_key }}</pre>
              </div>
            </div>

            <div class="box wide">
              <div class="box-head">
                <h3>MikroTik Configuration (.rsc)</h3>
                <div class="actions">
                  <button class="mini" @click="copy(result.mikrotik_config)">Copy</button>
                  <button class="mini" @click="downloadRsc()">Download .rsc</button>
                </div>
              </div>
              <pre class="mono">{{ result.mikrotik_config }}</pre>
            </div>

            <div class="box wide">
              <div class="box-head">
                <h3>Server Peer (wg0)</h3>
                <button class="mini" @click="copy(result.server_peer)">Copy</button>
              </div>
              <pre class="mono">{{ result.server_peer }}</pre>
            </div>

            <div class="cta-row">
              <RouterLink class="btn primary" to="/agents">Back to Agents</RouterLink>
            </div>
          </div>

          <p v-if="toast.text" :class="['toast', toast.ok ? 'success' : 'error']">{{ toast.text }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/BackgroundParticles.vue'
import { API } from '@/appCore.js'
import { ref, reactive, computed } from 'vue'
import { RouterLink } from 'vue-router'

const mac = ref('')
const siteId = ref('')
const loading = ref(false)
const result = ref(null)
const toast = reactive({ text: '', ok: true })

function showToast(text, ok = true) {
  toast.text = text
  toast.ok = ok
  setTimeout(() => (toast.text = ''), 2400)
}

function formatMac() {
  let v = mac.value.replace(/[^0-9a-f]/gi, '').toUpperCase().slice(0, 12)
  mac.value = v.match(/.{1,2}/g)?.join(':') ?? ''
}
const macValid = computed(() => /^[0-9A-F]{2}(:[0-9A-F]{2}){5}$/.test(mac.value))

import { preregister } from '@/appCore.js'

async function submit() {
  if (!macValid.value) return showToast('Invalid MAC address', false)
  loading.value = true
  result.value = null
  try {
    const data = await preregister(mac.value, siteId.value.trim() || undefined)
    result.value = data
    showToast('✅ Pre-registration complete')
  } catch (e) {
    showToast(e.message || 'Pre-registration failed', false)
  } finally {
    loading.value = false
  }
}


function copy(text) {
  navigator.clipboard.writeText(String(text || '')).then(
    () => showToast('Copied'),
    () => showToast('Copy failed', false)
  )
}

function downloadRsc() {
  if (!result.value?.mikrotik_config) return
  const blob = new Blob([result.value.mikrotik_config], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${mac.value.replace(/:/g, '')}_netsecure_wg.rsc`
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

.prereg-page { position:relative; min-height:100vh; overflow:hidden; }
.wrapper { position:relative; z-index:10; display:flex; justify-content:center; padding:32px; }
.container { width:100%; max-width:1000px; }
.card {
  background: var(--panel-grey);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.05);
  box-shadow: 0 0 40px rgba(0,194,194,.05);
}

.title { margin:0 0 6px; font-size:20px; }
.muted { color: var(--text-secondary); margin-bottom: 16px; }

.form { display:flex; flex-direction:column; gap:14px; background: rgba(255,255,255,.04); padding:14px; border-radius:12px; border:1px solid rgba(255,255,255,.06); }
.row { display:flex; flex-direction:column; gap:6px; }
.input {
  border: 1px solid var(--divider-grey); background: #0b0e16; color: var(--text-primary);
  border-radius: 8px; padding: 10px 12px; font-size: 14px;
}
.hint { color: var(--text-secondary); font-size: 12px; }
.mac-row { display: grid; grid-template-columns: 1fr auto; gap: 8px; }

.btn {
  display:inline-flex; align-items:center; gap:8px;
  border-radius:10px; padding:10px 12px; font-weight:600;
  border:1px solid rgba(255,255,255,.08); background:rgba(255,255,255,.06);
  color:var(--text-primary); cursor:pointer; transition:.15s;
}
.btn:hover { background:rgba(255,255,255,.10); }
.btn.primary { background: var(--primary-accent); color: #0e111a; border-color: transparent; }
.btn.primary:hover { background: var(--primary-hover); color: #fff; }
.mini { padding:6px 8px; font-size:12px; }

.result { margin-top: 16px; display:grid; gap:12px; }
.grid { display:grid; gap:12px; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); }
.box { border:1px solid rgba(255,255,255,.06); background: rgba(12,16,24,.35); padding:12px; border-radius:12px; }
.box.wide { grid-column: 1 / -1; }
.box-head { display:flex; align-items:center; justify-content:space-between; gap:8px; margin-bottom:6px; }
.box-head h3 { margin:0; font-size:14px; color:var(--text-secondary); }
.box-head .actions { display:flex; gap:8px; }
.mono { white-space: pre-wrap; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 13px; }

.cta-row { display:flex; justify-content:flex-end; margin-top:6px; }

.toast { margin-top:12px; text-align:center; padding:10px 12px; border-radius:8px; font-size:14px; }
.toast.success { background:rgba(34,197,94,.12); color:#22c55e; }
.toast.error { background:rgba(239,68,68,.12); color:#f87171; }
</style>
