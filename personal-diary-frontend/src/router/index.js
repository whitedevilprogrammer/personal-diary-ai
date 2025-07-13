import { createRouter, createWebHistory } from 'vue-router'
import Home from '../Views/Home.vue'
import Register from '../components/Register.vue'
import Login from '../components/Login.vue'
import Dashboard from '../Views/Dashboard.vue'
import ForgotPassword from '../components/ForgotPassword.vue'
import { useLoaderStore } from '../stores/loader.js'

const routes = [
  { path: '/', component: Home },
  { path: '/register', component: Register },
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard },
  { path: '/forgot-password', component: ForgotPassword },


]

const router = createRouter({
  history: createWebHistory(),
  routes,
})


router.beforeEach((to, from, next) => {
  const loader = useLoaderStore()
  loader.show()
  next()
})

router.afterEach(() => {
  const loader = useLoaderStore()
  // Add a short delay for better UX
  setTimeout(() => loader.hide(), 300)
})




export default router
