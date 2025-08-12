import { createRouter, createWebHashHistory } from 'vue-router'

// Components/views (chemins corrigés)
import HomePage from '@/components/home/HomePage.vue'
import LoginForm from '@/views/auth/LoginForm.vue'
import RegisterForm from '@/views/auth/RegisterForm.vue'
import DashboardPage from '@/views/dashboard/DashboardPage.vue'
import OrganizationForm from '@/components/organization/OrganizationForm.vue'
import OrganizationProfile from '@/components/organization/OrganizationProfile.vue'
import RouterTable from '@/components/organization/RouterTable.vue'
import AddOperatorForm from '@/components/forms/AddOperatorForm.vue'
import AddUserForm from '@/components/forms/AddUserForm.vue'
// Lazy views
const AgentDashboard = () => import('@/views/agents/AgentDashboard.vue');
const RegisterAgent = () => import('@/views/agents/RegisterAgent.vue');

const routes = [
  { path: '/', component: HomePage },
  { path: '/login', component: LoginForm },
  { path: '/register', component: RegisterForm },

  // Commun à tous les utilisateurs authentifiés
  { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
  { path: '/routertable', component: RouterTable, meta: { requiresAuth: true } },

  // Rôles : admin + operator
  { path: '/organization', component: OrganizationProfile, meta: { requiresAuth: true, roles: ['administrator', 'operator'] } },
  { path: '/organization/edit', component: OrganizationForm, meta: { requiresAuth: true, roles: ['administrator', 'operator'] } },

  // Add User : admin OU operator (operator doit avoir une org)
  { path: '/adduser', component: AddUserForm, meta: { requiresAuth: true, roles: ['administrator', 'operator'], requireOrgForOperator: true } },

  // Add Operator : admin uniquement
  { path: '/addoperator', component: AddOperatorForm, meta: { requiresAuth: true, roles: ['administrator'] } },

  // Agents : admin OU operator (operator doit avoir une org)
  { path: '/agents', component: AgentDashboard, meta: { requiresAuth: true, roles: ['administrator', 'operator'], requireOrgForOperator: true } },
  { path: '/agents/register', component: RegisterAgent, meta: { requiresAuth: true, roles: ['administrator', 'operator'], requireOrgForOperator: true } },

  // Legacy alias conservé si besoin
  { path: '/organizationForm', component: OrganizationForm, meta: { requiresAuth: true, roles: ['administrator', 'operator'] } },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// Guard d'auth + rôles ultra-simple basé sur localStorage
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const isAuthed = !!token

  const rawRole = (localStorage.getItem('role') || '').toLowerCase()
  // normalisation : certains endroits peuvent mettre "admin"
  const role = rawRole === 'admin' ? 'administrator' : rawRole
  const orgId = localStorage.getItem('organization_id') || ''

  // Bloque les pages protégées
  if (to.meta.requiresAuth && !isAuthed) {
    return next('/login')
  }

  // Évite d'aller sur login/register quand déjà connecté
  if (isAuthed && (to.path === '/login' || to.path === '/register')) {
    return next('/dashboard')
  }

  // Contrôle de rôle si précisé
  const allowedRoles = to.meta.roles
  if (Array.isArray(allowedRoles) && allowedRoles.length > 0) {
    if (!allowedRoles.includes(role)) {
      // rôle non autorisé -> on renvoie sur dashboard (UX soft)
      return next('/dashboard')
    }
  }

  // Pour operator, certaines routes exigent d'avoir une organisation
  if (to.meta.requireOrgForOperator && role === 'operator' && !orgId) {
    return next('/organization/edit')
  }

  return next()
})

export default router
