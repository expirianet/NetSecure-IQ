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

import { useAuth } from '../composables/useAuth';

const AgentDashboard = () => import('@/views/agents/AgentDashboard.vue');
const RegisterAgent = () => import('@/views/agents/RegisterAgent.vue');

const routes = [
  { path: '/', component: HomePage },
  { path: '/login', component: LoginForm },
  { path: '/register', component: RegisterForm },
<<<<<<< HEAD
  { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
  { path: '/organizationForm', component: OrganizationForm, meta: { requiresAuth: true } },
  { path: '/routertable', component: RouterTable, meta: { requiresAuth: true } },
  { path: '/adduser', component: AddUserForm, meta: { requiresAuth: true } },
  { path: '/addoperator', component: AddOperatorForm, meta: { requiresAuth: true } },
  { path: '/agents', component: AgentDashboard, meta: { requiresAuth: true } },
  { path: '/agents/register', component: RegisterAgent, meta: { requiresAuth: true } }
];
=======
  { path: '/dashboard', component: DashboardPage },
  { path: '/organization', component: OrganizationProfile },
  { path: '/organization/edit', component: OrganizationForm }, // ✅ c’est ici qu’on redirige
  { path: '/routertable', component: RouterTable },
  { path: '/adduser', component: AddUserForm },
  { path: '/addoperator', component: AddOperatorForm }
]
>>>>>>> 4e564af9b7e58dea5d6ffe261121df93dc7b54ed



const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const { isAuthenticated } = useAuth();
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/login');
  } else if ((to.path === '/login' || to.path === '/register') && isAuthenticated.value) {
    next('/dashboard');
  } else {
    next();
  }
});

export default router
