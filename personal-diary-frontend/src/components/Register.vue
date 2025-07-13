<template>
  <div class="form-page">
    <div class="form-container card">
      <div class="form-header">
        <div class="form-icon">✨</div>
        <h2>Create an Account</h2>
        <p>Join us to start your personal diary journey</p>
      </div>
      
      <form @submit.prevent="register">
        <div class="input-group">
          <label for="name">Name</label>
          <input 
            id="name"
            v-model="name" 
            placeholder="Your name" 
            required 
          />
        </div>
        <div class="input-group">
          <label for="email">Email</label>
          <input 
            id="email"
            v-model="email" 
            placeholder="your@email.com" 
            type="email" 
            required 
          />
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <div class="password-input-wrapper">
            <input 
              id="password"
              v-model="password" 
              placeholder="Choose a password" 
              :type="showPassword ? 'text' : 'password'" 
              required 
            />
            <button 
              type="button" 
              class="password-toggle" 
              @click="showPassword = !showPassword"
              tabindex="-1"
              aria-label="Toggle password visibility"
            >
              <span v-if="showPassword">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
              </span>
              <span v-else>
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                  <line x1="1" y1="1" x2="23" y2="23"></line>
                </svg>
              </span>
            </button>
          </div>
        </div>

        <button type="submit" class="btn">
          Register
        </button>
      </form>
      
      <div class="form-footer">
        <p>Already have an account? <router-link to="/login">Login</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'
import { alertService } from '../utils/alert'
import { useLoaderStore } from '../stores/loader'
export default {
  data() {
    return { 
      name: '', 
      email: '', 
      password: '',
      showPassword: false,
      loaderStore: null  // placeholder
    }
  },
  created() {
    // now Pinia is active, so we can grab the store
    this.loaderStore = useLoaderStore()
  },
  methods: {
    async register() {
      this.loaderStore.show()
      try {
        const response = await api.post('/register', {
          name: this.name,
          email: this.email,
          password: this.password,
        })

        if (response.status == "success") {
          alertService.success(response.message)
          this.$router.push('/login')
        } else if (response.status == "error") {
          alertService.error(response.message)
        }
      } catch (err) {
        alertService.error('Registration failed.')
      } finally {
        this.loaderStore.hide()
      }
    },
  },
}
</script>

<style scoped>
.form-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 90vh;
  padding: 20px;
}

.form-container {
  width: 100%;
  max-width: 450px;
  padding: 40px;
}

.form-header {
  text-align: center;
  margin-bottom: 30px;
}

.form-icon {
  font-size: 40px;
  margin-bottom: 16px;
  display: inline-block;
}

h2 {
  margin-bottom: 10px;
  font-size: 28px;
  color: var(--text-dark);
}

.form-header p {
  color: var(--text-light);
  margin-bottom: 0;
}

.input-group {
  margin-bottom: 16px;
  margin-right: 36px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: var(--text-dark);
}

button {
  margin-top: 10px;
  width: 100%;
}

.password-input-wrapper > button {
  margin: 0px 0px 0px 0px;
  padding: 0px 0px 0px 0px;
  width: 0px;
}


.form-footer {
  text-align: center;
  margin-top: 30px;
  color: var(--text-light);
}

.form-footer a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.form-footer a:hover {
  text-decoration: underline;
}

/* Password input styling */
.password-input-wrapper {
  position: relative;
  /* display: flex; */
  align-items: center;
}

.password-input-wrapper input {
  width: 100%;
  /* padding-right: 45px; Make space for the icon */
}

.password-toggle {
  position: absolute;
  right: -1px;
  top: 43%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
  opacity: 0.7;
  transition: opacity 0.2s;
  z-index: 2;
}

.password-toggle:hover {
  opacity: 1;
}

.password-toggle span {
  font-size: 20px;
  color: #6b7280;
}
</style>