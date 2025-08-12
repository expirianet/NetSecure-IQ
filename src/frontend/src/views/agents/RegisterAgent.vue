<template>
  <div class="register-agent-page">
    <BackgroundParticles />
    <div class="register-agent">
      <h1>Pré-enregistrement d’un agent MikroTik</h1>
      <form @submit.prevent="onSubmit">
        <label for="mac">Adresse MAC principale (ether1):</label>
        <input id="mac" v-model="mac" required placeholder="AA:BB:CC:DD:EE:FF" />
        <button type="submit" :disabled="loading">Générer le script</button>
      </form>
      <div v-if="loading" class="loading">Génération du script...</div>
      <div v-if="script">
        <label>Script généré :</label>
        <textarea readonly rows="8" :value="script"></textarea>
        <div class="actions">
          <button @click="copyScript">Copier</button>
          <button @click="downloadScript">Télécharger</button>
          <button @click="testAgent" :disabled="testing">Test</button>
          <button @click="cancelRegistration" :disabled="testing">Annuler l’inscription</button>
        </div>
        <div v-if="testing" class="loading">Test de connexion...</div>
        <div v-if="testResult === true" class="success">Agent en ligne, peer désactivé et retour à la liste.</div>
        <div v-if="testResult === false" class="error">Échec du test, vérifiez la connexion de l’agent.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/common/BackgroundParticles.vue';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const mac = ref('');
const script = ref('');
const loading = ref(false);
const testing = ref(false);
const testResult = ref(null);
const router = useRouter();

function onSubmit() {
  loading.value = true;
  setTimeout(() => {
    // Simuler la génération de script
    script.value = `/interface/wireguard/add listen-port=13231 private-key=...\n/interface/wireguard/peers/add public-key=... allowed-address=...`; 
    loading.value = false;
  }, 1500);
}
function copyScript() {
  navigator.clipboard.writeText(script.value);
}
function downloadScript() {
  const blob = new Blob([script.value], { type: 'text/plain' });
  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = `mikrotik-script-${mac.value}.rsc`;
  link.click();
}
function testAgent() {
  testing.value = true;
  testResult.value = null;
  setTimeout(() => {
    // Simuler le test (ping)
    const success = Math.random() > 0.5;
    testResult.value = success;
    testing.value = false;
    if (success) {
      // Simuler désactivation peer côté backend et retour à la liste
      setTimeout(() => router.push('/agents'), 1000);
    }
  }, 1500);
}
function cancelRegistration() {
  script.value = '';
  mac.value = '';
  testResult.value = null;
}
</script>

<style scoped>
.register-agent { padding: 2rem; max-width: 600px; margin: auto; }
form { margin-bottom: 1rem; }
input { margin-right: 1rem; }
textarea { width: 100%; margin-bottom: 1rem; }
.actions { display: flex; gap: 1rem; }
.loading { color: #1976d2; font-weight: bold; }
.success { color: #4caf50; font-weight: bold; }
.error { color: #f44336; font-weight: bold; }
</style>
