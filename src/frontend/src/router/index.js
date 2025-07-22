import { createRouter, createWebHashHistory } from 'vue-router'
import HomePage from '@/components/HomePage.vue'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import DashboardPage from '@/components/DashboardPage.vue'
import OrganizationForm from '@/components/OrganizationForm.vue'
import RouterTable from '@/components/RouterTable.vue'
import AddUserForm from '@/components/AddUserForm.vue'
import AddOperatorForm from '@/components/AddOperatorForm.vue'

const routes = [
  { path: '/', component: HomePage },
  { path: '/login', component: LoginForm },
  { path: '/register', component: RegisterForm },
  { path: '/dashboard', component: DashboardPage },
  { path: '/organizationForm', component: OrganizationForm },
  { path: '/routertable', component: RouterTable },
  { path: '/adduser', component: AddUserForm },
  { path: '/addoperator', component: AddOperatorForm }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
