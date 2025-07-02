import { createRouter, createWebHashHistory } from 'vue-router'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import DashboardPage from '@/components/DashboardPage.vue'
import OrganizationForm from '@/components/OrganizationForm.vue'


const routes = [
  { path: '/login', component: LoginForm },
  { path: '/register', component: RegisterForm },
  { path: '/dashboard', component: DashboardPage },
  { path: '/organizationForm', component: OrganizationForm }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
