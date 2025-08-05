<template>
  <div class="agent-dashboard">
    <h1>Gestion des agents MikroTik</h1>
    <div class="actions">
      <button :disabled="!selectedAgent" @click="onAssociate">Associer</button>
      <button :disabled="!selectedAgent" @click="onDeactivate">Désactiver</button>
      <button :disabled="!selectedAgent" @click="onDelete">Supprimer</button>
      <router-link to="/agents/register"><button>Pré-enregistrer un agent</button></router-link>
    </div>
    <table class="agent-table">
      <thead>
        <tr>
          <th>Sélection</th>
          <th>MAC Adresse</th>
          <th>Statut</th>
          <th>Organisation</th>
          <th>Site</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="agent in agents" :key="agent.mac" :class="[statusClass(agent.status), {selected: agent === selectedAgent}]" @click="selectAgent(agent)">
          <td><input type="radio" :checked="agent === selectedAgent" /></td>
          <td>{{ agent.mac }}</td>
          <td><span :class="statusBadgeClass(agent.status)">{{ agent.statusLabel }}</span></td>
          <td>{{ agent.organization || '-' }}</td>
          <td>{{ agent.site || '-' }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const agents = ref([
  { mac: 'AA:BB:CC:DD:EE:01', status: 'associated', organization: 'Org1', site: 'SiteA' },
  { mac: 'AA:BB:CC:DD:EE:02', status: 'unassociated' },
  { mac: 'AA:BB:CC:DD:EE:03', status: 'deactivated', organization: 'Org2' },
]);
const selectedAgent = ref(null);

function selectAgent(agent) {
  selectedAgent.value = agent;
}
function statusClass(status) {
  return {
    'row-associated': status === 'associated',
    'row-unassociated': status === 'unassociated',
    'row-deactivated': status === 'deactivated',
  };
}
function statusBadgeClass(status) {
  return {
    badge: true,
    'badge-green': status === 'associated',
    'badge-orange': status === 'unassociated',
    'badge-red': status === 'deactivated',
  };
}
function onAssociate() {
  // TODO: implémenter
}
function onDeactivate() {
  // TODO: implémenter
}
function onDelete() {
  // TODO: implémenter
}
</script>

<style scoped>
.agent-dashboard { padding: 2rem; }
.actions { margin-bottom: 1rem; display: flex; gap: 1rem; }
.agent-table { width: 100%; border-collapse: collapse; }
.agent-table th, .agent-table td { padding: 0.5rem 1rem; border-bottom: 1px solid #222; }
.row-associated { background: #e8fbe8; }
.row-unassociated { background: #fff8e1; }
.row-deactivated { background: #fdeaea; }
.selected { outline: 2px solid #1976d2; }
.badge { padding: 0.2em 0.7em; border-radius: 1em; font-weight: bold; }
.badge-green { background: #4caf50; color: #fff; }
.badge-orange { background: #ff9800; color: #fff; }
.badge-red { background: #f44336; color: #fff; }
</style>
