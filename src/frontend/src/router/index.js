import { createRouter, createWebHashHistory } from 'vue-router'
import HomePage from '@/components/HomePage.vue'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import DashboardPage from '@/components/DashboardPage.vue'
import OrganizationForm from '@/components/OrganizationForm.vue'
import OrganizationProfile from '@/components/organization/OrganizationProfile.vue'
import RouterTable from '@/components/RouterTable.vue'
import AddUserForm from '@/components/AddUserForm.vue'
import AddOperatorForm from '@/components/AddOperatorForm.vue'
import DashboardOperator from '@/components/DashboardOperator.vue'

import { useAuth } from '../appCore.js'

const AgentDashboard = () => import('@/views/agents/AgentDashboard.vue')
const RegisterAgent = () => import('@/views/agents/RegisterAgent.vue')

const routes = [
  { path: '/', component: HomePage },
  { path: '/login', component: LoginForm },
  { path: '/register', component: RegisterForm },

  // Dashboards
  { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
  { path: '/dashboard-operator', component: DashboardOperator, meta: { requiresAuth: true } },

  // Organization
  { path: '/organization', component: OrganizationProfile, meta: { requiresAuth: true } },
  { path: '/organization/edit', component: OrganizationForm, meta: { requiresAuth: true } },
  { path: '/organizationForm', component: OrganizationForm, meta: { requiresAuth: true } },

  // Features
  { path: '/routertable', component: RouterTable, meta: { requiresAuth: true } },
  { path: '/adduser', component: AddUserForm, meta: { requiresAuth: true } },
  { path: '/addoperator', component: AddOperatorForm, meta: { requiresAuth: true } },

  // Agents
  { path: '/agents', component: AgentDashboard, meta: { requiresAuth: true } },
  { path: '/agents/register', component: RegisterAgent, meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// Auth guard + redirection opérateur -> dashboard-operator
router.beforeEach((to, from, next) => {
  const { isAuthenticated } = useAuth()
  const role = (localStorage.getItem('role') || '').toLowerCase()

  // Besoin d'être connecté
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/login')
    return
  }

  // Si déjà connecté et on va vers login/register -> dashboard adapté au rôle
  if ((to.path === '/login' || to.path === '/register') && isAuthenticated.value) {
    next(role === 'operator' ? '/dashboard-operator' : '/dashboard')
    return
  }

  // Un opérateur qui vise /dashboard est redirigé vers /dashboard-operator
  if (to.path === '/dashboard' && role === 'operator') {
    next('/dashboard-operator')
    return
  }

  next()
})

export default router
