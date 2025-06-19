import { createRouter, createWebHistory } from 'vue-router'
import RegisterForm from '../components/RegisterForm.vue'
import LoginForm from '../components/LoginForm.vue'

const routes = [
  { path: '/register', component: RegisterForm },
  { path: '/login', component: LoginForm }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
