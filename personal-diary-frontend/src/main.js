// import { createApp } from 'vue'
// import './style.css'
// import App from './App.vue'

// createApp(App).mount('#app')

// import { createApp } from 'vue'
// import App from './App.vue'
// import router from './router'
// import AlertService from './utils/alert.js'
// import { createPinia } from 'pinia'

// const app = createApp(App)   // create app once
// const pinia = createPinia()

// app.use(AlertService)
// app.use(router)
// app.use(pinia)

// app.mount('#app')            // mount the same app instance


import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import AlertService from './utils/alert.js'


// 🌟 Step 1: Read token from URL
const urlParams = new URLSearchParams(window.location.search)
const token = urlParams.get('token')

if (token) {
  // 🌟 Step 2: Store in localStorage (or cookie if preferred)
  localStorage.setItem('token', token)

  // 🌟 Step 3: Remove token from URL (optional clean-up)
  const cleanUrl = window.location.origin + window.location.pathname
  window.history.replaceState({}, document.title, cleanUrl)
}

const app = createApp(App)
app.use(AlertService)
app.use(router)
app.use(createPinia())
app.mount('#app')


