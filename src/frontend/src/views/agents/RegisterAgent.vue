<!-- src/frontend/src/views/agents/RegisterAgent.vue -->
<template>
  <div class="register-agent-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <div class="card">
          <!-- Header -->
          <div class="header-row">
            <div class="titles">
              <h2>Pré-enregistrement d’un agent MikroTik</h2>
              <p class="subtitle">
                Saisis l’adresse MAC (ether1), puis génère le script <code>.rsc</code> à exécuter sur le routeur.
              </p>
            </div>

            <button v-if="script" class="btn ghost" @click="resetAll">
              <i class="fas fa-eraser"></i><span>Réinitialiser</span>
            </button>
          </div>

          <!-- Formulaire compact -->
          <form class="form" @submit.prevent="onGenerate">
            <label for="mac" class="label">Adresse MAC principale (ether1)</label>

            <div class="row">
              <input
                id="mac"
                class="input mono"
                v-model.trim="mac"
                placeholder="AA:BB:CC:DD:EE:FF"
                maxlength="17"
                @input="formatMac"
                autocomplete="off"
              />
              <button class="btn primary" type="submit" :disabled="loading || !macValid">
                <i class="fas fa-microchip"></i>
                <span>{{ loading ? 'Génération…' : 'Générer le script' }}</span>
              </button>
            </div>

            <small class="hint">Format attendu : 6 octets hexadécimaux séparés par « : ».</small>
          </form>

          <!-- Résultat -->
          <div v-if="script" class="result">
            <div class="toolbar">
              <span class="tag">
                <i class="fas fa-file-code"></i> script.rsc
              </span>
              <div class="actions">
                <button class="btn ghost" @click="copyScript"><i class="fas fa-copy"></i><span>Copier</span></button>
                <button class="btn ghost" @click="downloadScript"><i class="fas fa-download"></i><span>Télécharger</span></button>
              </div>
            </div>

            <textarea class="code" readonly rows="12" :value="script"></textarea>
            <p class="note">
              Applique le script via <em>Files → Run Script</em> ou en CLI sur le routeur.
            </p>
          </div>

          <p v-if="toast" class="toast" :class="toastType">{{ toast }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/BackgroundParticles.vue'
import { ref, computed } from 'vue'

const mac = ref('')
const script = ref('')
const loading = ref(false)
const toast = ref('')
const toastType = ref('success')

const macRe = /^[0-9A-F]{2}(:[0-9A-F]{2}){5}$/i
const macValid = computed(() => macRe.test(mac.value))

function formatMac() {
  // Uppercase + auto-insertion des ":" (AA:BB:…)
  let v = mac.value.replace(/[^0-9a-f]/gi, '').toUpperCase().slice(0, 12)
  mac.value = v.match(/.{1,2}/g)?.join(':') ?? ''
}

function makeDummyScript(macAddr) {
  // Génère un script Mikrotik plausible (démo locale)
  const last = macAddr.split(':').pop() || 'FF'
  const ipOctet = parseInt(last, 16) || Math.floor(Math.random() * 200) + 10
  return `
# MikroTik WireGuard onboarding script (démo)
# MAC: ${macAddr}
# Exécuter : /import file=script.rsc

/interface wireguard add name=wg0 listen-port=51820 private-key="CHANGEME_PRIVATE"
/ip address add address=10.10.10.${ipOctet}/32 interface=wg0

# Pair serveur (à ajuster)
/interface wireguard peers add \
  interface=wg0 \
  public-key="SERVER_PUBLIC_KEY" \
  endpoint-address=vpn.example.com \
  endpoint-port=51820 \
  allowed-address=0.0.0.0/0 \
  persistent-keepalive=25s

# Tag local
:put "Agent pré-enregistré (${macAddr})"
`.trim()
}

async function onGenerate() {
  if (!macValid.value) return
  loading.value = true
  script.value = ''
  try {
    // Ici on pourrait appeler un vrai endpoint backend si dispo.
    // Pour le moment, on génère un script local « démo » cohérent.
    await new Promise(r => setTimeout(r, 500))
    script.value = makeDummyScript(mac.value)
  } catch (e) {
    showToast('Erreur pendant la génération.', 'error')
  } finally {
    loading.value = false
  }
}

function copyScript() {
  navigator.clipboard.writeText(script.value)
  showToast('Script copié dans le presse-papiers.')
}
function downloadScript() {
  const blob = new Blob([script.value], { type: 'text/plain' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `mikrotik-${mac.value.replace(/:/g, '')}.rsc`
  a.click()
  URL.revokeObjectURL(a.href)
}
function resetAll() {
  mac.value = ''
  script.value = ''
}
function showToast(msg, type = 'success') {
  toast.value = msg
  toastType.value = type
  setTimeout(() => (toast.value = ''), 2200)
}
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root{
  --bg-dark:#0e111a; --panel-grey:#1a1d26; --divider-grey:#2a2d36;
  --text-primary:#f5f7fa; --text-secondary:#9ca3af;
  --primary-accent:#00c2c2; --primary-hover:#00a7a7;
  --danger:#ef4444; --success:#22c55e; --radius:12px;
}

/* Layout carte */
.register-agent-page{ position:relative; min-height:100vh; overflow:hidden; }
.wrapper{ position:relative; z-index:10; display:flex; justify-content:center; padding:32px; }
.container{ width:100%; max-width:900px; }
.card{
  background:var(--panel-grey);
  border-radius:16px;
  padding:20px;
  border:1px solid rgba(255,255,255,.05);
  box-shadow:0 0 40px rgba(0,194,194,.05);
}

/* Header */
.header-row{ display:flex; align-items:flex-start; justify-content:space-between; gap:12px; }
.titles h2{ margin:0; font-size:20px; color:var(--text-primary); }
.subtitle{ margin:.25rem 0 0; color:var(--text-secondary); }

/* Form */
.form{ margin-top:14px; display:flex; flex-direction:column; gap:8px; }
.label{ font-size:13px; color:var(--text-secondary); }
.row{ display:grid; grid-template-columns: 1fr auto; gap:10px; }
.input{
  width:100%; padding:12px 14px; border-radius:8px;
  border:1px solid var(--divider-grey); background:#121521; color:var(--text-primary);
}
.input:focus{ outline:none; border-color:var(--primary-accent); background:#0b0e16; }
.mono{ font-family: ui-monospace,SFMono-Regular,Menlo,monospace; letter-spacing:.2px; }
.hint{ color:var(--text-secondary); }

/* Buttons */
.btn{
  display:inline-flex; align-items:center; gap:8px;
  border-radius:10px; padding:10px 14px; font-weight:600;
  border:1px solid rgba(255,255,255,.08); background:rgba(255,255,255,.06);
  color:var(--text-primary); cursor:pointer; transition:.15s;
}
.btn:hover{ background:rgba(255,255,255,.10); }
.btn:disabled{ opacity:.5; cursor:not-allowed; }
.btn.primary{ background:var(--primary-accent); color:#0e111a; border-color:transparent; }
.btn.primary:hover{ background:var(--primary-hover); color:#fff; }
.btn.ghost{ background:rgba(255,255,255,.06); }

/* Résultat */
.result{ margin-top:16px; }
.toolbar{
  display:flex; align-items:center; justify-content:space-between;
  gap:8px; margin-bottom:8px;
}
.tag{
  display:inline-flex; align-items:center; gap:8px;
  padding:6px 10px; border-radius:999px;
  background:rgba(255,255,255,.06); border:1px solid rgba(255,255,255,.08);
  color:var(--text-primary); font-size:12px;
}
.actions{ display:flex; gap:8px; }
.code{
  width:100%; border-radius:10px; padding:12px 14px;
  border:1px solid var(--divider-grey); background:#0b0e16; color:#e5e7eb;
  font-family: ui-monospace,SFMono-Regular,Menlo,monospace;
}
.note{ color:var(--text-secondary); margin-top:8px; }

/* Toast */
.toast{
  margin-top:12px; text-align:center; padding:10px 12px; border-radius:8px; font-size:14px;
}
.toast.success{ background:rgba(34,197,94,.12); color:var(--success); }
.toast.error{ background:rgba(239,68,68,.12); color:var(--danger); }

/* Responsive */
@media (max-width: 820px){
  .row{ grid-template-columns: 1fr; }
}
</style>
