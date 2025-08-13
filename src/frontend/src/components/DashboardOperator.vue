<!-- src/frontend/src/views/dashboard/DashboardOperator.vue -->
<template>
    <section class="dashboard-operator">
      <header class="header">
        <h1>Dashboard — Operator</h1>
        <p class="subtitle">
          Vue opérateur (portée limitée à votre organisation).
        </p>
      </header>
  
      <!-- État : pas encore rattaché à une organisation -->
      <div v-if="!hasOrg" class="card warning">
        <h2>Onboarding requis</h2>
        <p>
          Vous devez être rattaché(e) à une organisation pour accéder à toutes les fonctionnalités.
        </p>
        <div class="actions">
          <RouterLink class="btn primary" to="/organization/edit">
            Créer mon organisation
          </RouterLink>
          <RouterLink class="btn" to="/onboarding-operator">
            Demander un rattachement
          </RouterLink>
        </div>
        <p class="hint">
          Une fois l’organisation créée ou validée par un administrateur, vous aurez accès aux actions ci-dessous.
        </p>
      </div>
  
      <!-- Contenu principal quand l'opérateur est rattaché -->
      <div v-else class="grid">
        <!-- Raccourcis principaux -->
        <div class="card">
          <h3>Raccourcis</h3>
          <div class="quick-actions">
            <RouterLink class="shortcut" to="/routertable" aria-label="Accéder au RouterTable">
              <span class="title">RouterTable</span>
              <span class="desc">Liste des routeurs et état (scopé à l’organisation)</span>
            </RouterLink>
  
            <RouterLink class="shortcut" to="/organization" aria-label="Profil organisation">
              <span class="title">Organisation</span>
              <span class="desc">Profil / informations de l’organisation</span>
            </RouterLink>
  
            <!-- IMPORTANT : l'opérateur PEUT ajouter des utilisateurs de son org -->
            <RouterLink class="shortcut" to="/adduser" aria-label="Ajouter un utilisateur">
              <span class="title">Ajouter un utilisateur</span>
              <span class="desc">Créer des comptes “User” pour vos sites</span>
            </RouterLink>
  
            <RouterLink class="shortcut" to="/agents" aria-label="Gérer les agents">
              <span class="title">Agents</span>
              <span class="desc">Voir et gérer les agents (portée org)</span>
            </RouterLink>
  
            <RouterLink class="shortcut" to="/agents/register" aria-label="Pré-enregistrer un agent">
              <span class="title">Pré-enregistrement Agent</span>
              <span class="desc">Associer un routeur MikroTik (MAC) à un site</span>
            </RouterLink>
          </div>
        </div>
  
        <!-- Widget léger d’information -->
        <div class="card">
          <h3>Conseils</h3>
          <ul class="tips">
            <li>Utilisez <em>Pré-enregistrement Agent</em> pour l’auto-provisioning via l’adresse MAC.</li>
            <li>Créez des <em>Users</em> pour donner un accès ciblé à certains sites.</li>
            <li>Le <em>RouterTable</em> affiche uniquement les équipements de votre organisation.</li>
          </ul>
        </div>
      </div>
  
      <!-- Note de conformité au périmètre de l'opérateur -->
      <footer class="footnote">
        <small>
          Cette vue respecte les droits de l’Operator&nbsp;: gestion de son organisation, création d’utilisateurs de son org,
          gestion/lecture des sites et agents de son org — sans création d’opérateurs. 
        </small>
      </footer>
    </section>
  </template>
  
  <script setup lang="ts">
  import { computed } from 'vue'
  import { useAuthStore } from '@/stores/authStore'
  
  const auth = useAuthStore()
  
  // On considère que le store expose un user avec organization_id.
  // Adapter si votre store a une autre forme.
  const hasOrg = computed(() => {
    const orgId = (auth.user as any)?.organization_id
    return !!orgId
  })
  </script>
  
  <style scoped>
  .dashboard-operator {
    display: grid;
    gap: 1rem;
    padding: 1rem;
  }
  .header .subtitle {
    opacity: 0.8;
    margin-top: .25rem;
  }
  .grid {
    display: grid;
    gap: 1rem;
  }
  .card {
    border: 1px solid #e5e7eb;
    border-radius: .75rem;
    padding: 1rem;
    background: #fff;
  }
  .card.warning {
    border-color: #f59e0b;
    background: #fffbeb;
  }
  .actions {
    display: flex;
    gap: .5rem;
    margin-top: .75rem;
    flex-wrap: wrap;
  }
  .btn {
    display: inline-block;
    padding: .5rem .75rem;
    border-radius: .5rem;
    border: 1px solid #d1d5db;
    text-decoration: none;
  }
  .btn.primary {
    background: #111827;
    color: #fff;
    border-color: #111827;
  }
  .quick-actions {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: .75rem;
    margin-top: .5rem;
  }
  .shortcut {
    display: block;
    padding: .75rem;
    border: 1px solid #e5e7eb;
    border-radius: .75rem;
    text-decoration: none;
    color: inherit;
  }
  .shortcut:hover {
    background: #f9fafb;
  }
  .shortcut .title {
    font-weight: 600;
  }
  .shortcut .desc {
    display: block;
    opacity: 0.7;
    font-size: .9rem;
  }
  .tips {
    margin: 0;
    padding-left: 1rem;
  }
  .footnote {
    opacity: 0.7;
  }
  </style>
  