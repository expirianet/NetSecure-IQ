<!-- src/frontend/src/views/agents/AgentDashboard.vue -->
<template>
  <div class="agents-page">
    <!-- Particules de fond -->
    <BackgroundParticles />

    <div class="agents-wrapper">
      <div class="agents-container">
        <div class="agents-card">
          <!-- Header -->
          <div class="header-row">
            <div class="titles">
              <h2>Gestion des agents MikroTik</h2>
              <p class="subtitle">Surveillez, associez et gérez vos routeurs déployés.</p>
            </div>

            <RouterLink to="/agents/register" class="btn primary">
              <i class="fas fa-microchip"></i>
              <span>Pré-enregistrer un agent</span>
            </RouterLink>
          </div>

          <!-- Barre outils -->
          <div class="toolbar">
            <div class="search">
              <i class="fas fa-search"></i>
              <input
                v-model.trim="query"
                type="text"
                placeholder="Rechercher par MAC, organisation ou site"
                aria-label="Rechercher"
              />
            </div>

            <div class="actions">
              <button class="btn ghost" :disabled="!selectedAgent" @click="onAssociate">
                <i class="fas fa-link"></i><span>Associer</span>
              </button>
              <button class="btn ghost" :disabled="!selectedAgent" @click="onDeactivate">
                <i class="fas fa-ban"></i><span>Désactiver</span>
              </button>
              <button class="btn ghost danger" :disabled="!selectedAgent" @click="onDelete">
                <i class="fas fa-trash"></i><span>Supprimer</span>
              </button>
            </div>

            <div class="legend">
              <span class="chip green"><span class="dot"></span>Associé</span>
              <span class="chip orange"><span class="dot"></span>Non associé</span>
              <span class="chip red"><span class="dot"></span>Désactivé</span>
            </div>
          </div>

          <!-- Tableau -->
          <div class="table-wrapper">
            <table class="agents-table" role="table" aria-label="Liste des agents">
              <thead>
                <tr>
                  <th style="width:72px">Sélection</th>
                  <th>MAC Adresse</th>
                  <th style="width:170px">Statut</th>
                  <th>Organisation</th>
                  <th>Site</th>
                </tr>
              </thead>

              <tbody>
                <tr
                  v-for="(a, idx) in filteredAgents"
                  :key="a.mac"
                  :class="['row', rowTint(a.status), { selected: idx === selectedIndex } ]"
                  @click="selectAgent(a, idx)"
                >
                  <td>
                    <span class="radio" :class="{ checked: idx === selectedIndex }" aria-hidden="true"></span>
                    <span class="sr-only">Sélectionner {{ a.mac }}</span>
                  </td>
                  <td class="mono">{{ a.mac }}</td>
                  <td>
                    <span :class="['state', stateClass(a.status)]">{{ label(a.status) }}</span>
                  </td>
                  <td>{{ a.organization || '—' }}</td>
                  <td>{{ a.site || '—' }}</td>
                </tr>

                <tr v-if="!filteredAgents.length">
                  <td colspan="5" class="empty">Aucun agent ne correspond à la recherche.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Toast -->
          <p v-if="toast" class="toast" :class="toastType">{{ toast }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/BackgroundParticles.vue'
import { ref, computed } from 'vue'

/** Données de démonstration (branche sur ton API quand dispo) */
const agents = ref([
  { mac: 'AA:BB:CC:DD:EE:01', status: 'associated',   organization: 'Org1', site: 'SiteA' },
  { mac: 'AA:BB:CC:DD:EE:02', status: 'unassociated', organization: '',     site: ''     },
  { mac: 'AA:BB:CC:DD:EE:03', status: 'deactivated',  organization: 'Org2', site: ''     }
])

/** Recherche & sélection */
const query = ref('')
const selectedAgent = ref(null)
const selectedIndex = ref(-1)

const filteredAgents = computed(() => {
  const q = query.value.toLowerCase()
  if (!q) return agents.value
  return agents.value.filter(a =>
    a.mac.toLowerCase().includes(q) ||
    a.organization?.toLowerCase().includes(q) ||
    a.site?.toLowerCase().includes(q)
  )
})

function selectAgent(a, idx) {
  selectedAgent.value = a
  selectedIndex.value = idx
}

/** Actions (stubs) */
const toast = ref('')
const toastType = ref('success')
const showToast = (msg, type = 'success') => {
  toast.value = msg
  toastType.value = type
  setTimeout(() => (toast.value = ''), 2200)
}

function onAssociate() {
  if (!selectedAgent.value) return
  selectedAgent.value.status = 'associated'
  showToast('Agent associé.')
}
function onDeactivate() {
  if (!selectedAgent.value) return
  selectedAgent.value.status = 'deactivated'
  showToast('Agent désactivé.', 'error')
}
function onDelete() {
  if (!selectedAgent.value) return
  agents.value = agents.value.filter((_, i) => i !== selectedIndex.value)
  selectedAgent.value = null
  selectedIndex.value = -1
  showToast('Agent supprimé.', 'error')
}

/** Helpers d’affichage */
function label(status) {
  return status === 'associated'
    ? 'Associé'
    : status === 'unassociated'
    ? 'Non associé'
    : 'Désactivé'
}
function stateClass(status) {
  return status === 'associated'
    ? 'green'
    : status === 'unassociated'
    ? 'orange'
    : 'red'
}
function rowTint(status) {
  return status === 'associated'
    ? 'tint-green'
    : status === 'unassociated'
    ? 'tint-orange'
    : 'tint-red'
}
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

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

/* Layout de page */
.agents-page { position: relative; min-height: 100vh; overflow: hidden; }
.agents-wrapper { position: relative; z-index: 10; display: flex; justify-content: center; padding: 32px; }
.agents-container { width: 100%; max-width: 1120px; }
.agents-card {
  background: var(--panel-grey);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.05);
  box-shadow: 0 0 40px rgba(0,194,194,.05);
}

/* Header */
.header-row { display: grid; grid-template-columns: 1fr auto; align-items: center; gap: 16px; }
.titles h2 { margin: 0; font-size: 20px; color: var(--text-primary); }
.subtitle { margin: 4px 0 0; color: var(--text-secondary); }

/* Buttons */
.btn {
  display: inline-flex; align-items: center; gap: 8px;
  border-radius: 10px; padding: 10px 14px; font-weight: 600;
  border: 1px solid rgba(255,255,255,.08); background: rgba(255,255,255,.06);
  color: var(--text-primary); cursor: pointer; transition: .15s;
}
.btn:hover { background: rgba(255,255,255,.10); }
.btn.primary {
  background: var(--primary-accent); color: #0e111a; border-color: transparent;
}
.btn.primary:hover { background: var(--primary-hover); color: #fff; }
.btn.ghost.danger { border-color: rgba(239,68,68,.35); color: #fca5a5; }
.btn:disabled { opacity: .5; cursor: not-allowed; }

/* Toolbar */
.toolbar {
  margin-top: 16px; display: grid; gap: 12px;
  grid-template-columns: minmax(240px, 1fr) auto auto;
  align-items: center;
}
.search {
  display: flex; align-items: center; gap: 10px;
  background: rgba(255,255,255,.06); border: 1px solid rgba(255,255,255,.08);
  padding: 10px 12px; border-radius: 10px;
}
.search input {
  border: none; outline: none; background: transparent; color: var(--text-primary); width: 100%;
}
.actions { display: flex; gap: 8px; }
.legend { display: flex; gap: 8px; justify-self: end; }

/* Chips */
.chip {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 10px; border-radius: 999px; font-size: 12px;
  border: 1px solid rgba(255,255,255,.08); background: rgba(255,255,255,.06);
}
.chip .dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.chip.green  .dot { background: #16a34a; }
.chip.orange .dot { background: #f59e0b; }
.chip.red    .dot { background: #ef4444; }

/* Table */
.table-wrapper { margin-top: 14px; overflow: auto; border-radius: 12px; border: 1px solid rgba(255,255,255,.06); }
.agents-table { width: 100%; border-collapse: separate; border-spacing: 0; font-size: 14px; }
.agents-table thead th {
  text-align: left; padding: 12px 14px; background: rgba(255,255,255,.04);
  border-bottom: 1px solid rgba(255,255,255,.08); color: var(--text-secondary);
}
.agents-table tbody td { padding: 14px; border-bottom: 1px solid rgba(255,255,255,.05); }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; letter-spacing: .2px; }

.row { transition: background .15s, outline .15s; cursor: pointer; }
.row:hover { outline: 1px solid rgba(255,255,255,.08); }
.selected { outline: 2px solid rgba(0,194,194,.35); }

/* Teinte de ligne selon statut (subtile) */
.tint-green  { background: rgba(22,163,74,.06); }
.tint-orange { background: rgba(245,158,11,.06); }
.tint-red    { background: rgba(239,68,68,.06); }

/* Radio custom */
.radio { display:inline-block; width:12px; height:12px; border-radius:50%;
  border:2px solid rgba(255,255,255,.35); vertical-align: middle; margin-right: 6px; }
.radio.checked { background: var(--primary-accent); border-color: var(--primary-accent); }

/* Badges état */
.state {
  display: inline-block; padding: 6px 10px; border-radius: 999px; font-weight: 600; font-size: 12px;
}
.state.green  { background: rgba(22,163,74,.15); color:#22c55e; }
.state.orange { background: rgba(245,158,11,.15); color:#fbbf24; }
.state.red    { background: rgba(239,68,68,.15); color:#f87171; }

/* Vide */
.empty { text-align: center; color: var(--text-secondary); padding: 20px; }

/* Toast */
.toast {
  margin-top: 12px; text-align:center; padding: 10px 12px; border-radius: 8px; font-size: 14px;
}
.toast.success { background: rgba(34,197,94,.12); color: var(--success); }
.toast.error   { background: rgba(239,68,68,.12); color: var(--danger); }

/* Responsive */
@media (max-width: 980px) {
  .toolbar { grid-template-columns: 1fr; }
  .legend { justify-self: start; }
}
</style>
