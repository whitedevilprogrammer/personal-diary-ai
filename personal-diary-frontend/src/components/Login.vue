<template>
  <div class="form-page">
    <div class="form-container card">
      <div class="form-header">
        <div class="form-icon">🔐</div>
        <h2>Login to Your Diary</h2>
        <p>Welcome back! Please enter your credentials</p>
      </div>

      <!-- Google OAuth Button -->
      <div class="oauth-section"> 
        <button type="button" class="google-btn" @click="loginWithGoogle" :disabled="isLoading">
          <div class="google-icon">
            <svg width="18" height="18" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
          </div>
          Continue with Google
        </button>
      </div>

      <!-- Divider -->
      <div class="divider">
        <span>or</span>
      </div>

      <!-- Email/Password Form -->
      <form @submit.prevent="login">
        <div class="input-group">
          <label for="email">Email</label>
          <input id="email" v-model="email" placeholder="your@email.com" type="email" required />
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <div class="password-input-wrapper">
            <input id="password" v-model="password" placeholder="Your password"
              :type="showPassword ? 'text' : 'password'" required />
            <button type="button" class="password-toggle" @click="showPassword = !showPassword" tabindex="-1"
              aria-label="Toggle password visibility">
              <span v-if="showPassword">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
              </span>
              <span v-else>
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path
                    d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24">
                  </path>
                  <line x1="1" y1="1" x2="23" y2="23"></line>
                </svg>
              </span>
            </button>
          </div>
        </div>

        <div class="forgot-password">
          <router-link to="/forgot-password">Forgot your password?</router-link>
        </div>

        <button type="submit" class="btn secondary" :disabled="isLoading">
          <span v-if="isLoading">Logging in...</span>
          <span v-else>Login</span>
        </button>
      </form>

      <div class="form-footer">
        <p>Don't have an account? <router-link to="/register">Register</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'
import { alertService } from '../utils/alert'
import { useLoaderStore } from '../stores/loader'

const googleAuthURL = "http://localhost:8080/auth/google/login"

export default {
  data() {
    return {
      email: '',
      password: '',
      showPassword: false,
      isLoading: false,
      loaderStore: null,
      googleAuthURL,
    }
  },
  created() {
    this.loaderStore = useLoaderStore()
    // Check if user is returning from OAuth with token
    this.checkForOAuthToken()
  },
  methods: {
    checkForOAuthToken() {
      debugger
      const urlParams = new URLSearchParams(window.location.search)
      const token = urlParams.get('token')
      if (token) {
        localStorage.setItem('token', token)
        alertService.success('Login successful!', 'Welcome')
        // Clean URL and redirect to dashboard
        window.history.replaceState({}, document.title, window.location.pathname)
        this.$router.push('/dashboard')
      }
    },

    async loginWithGoogle() {
      try {
        this.isLoading = true
        // Redirect to the OAuth endpoint
        window.location.href = this.googleAuthURL
      } catch (error) {
        console.error('Google OAuth error:', error)
        alertService.error('Google login failed. Please try again.', 'Error')
      } finally {
        this.isLoading = false
      }
    },

    async login() {
      if (this.isLoading) return
      
      this.isLoading = true
      this.loaderStore.show()
      
      try {
        const res = await api.post('/login', {
          email: this.email,
          password: this.password
        })
        
        if (res.status === 'success') {
          alertService.success(res.message, 'Login')
          localStorage.setItem('token', res.data.token)
          this.$router.push('/dashboard')
        } else if (res.status === 'error') {
          alertService.error(res.message, 'Error')
        }
      } catch (err) {
        console.error('Login error:', err)
        alertService.error('Login failed. Please check your credentials.', 'Error')
      } finally {
        this.isLoading = false
        this.loaderStore.hide()
      }
    }
  }
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

/* OAuth Section */
.oauth-section {
  margin-bottom: 24px;
}

.google-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px 24px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #374151;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-top: 0;
}

.google-btn:hover:not(:disabled) {
  border-color: #d1d5db;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
  background: #fafafa;
}

.google-btn:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.google-btn:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(66, 133, 244, 0.3);
}

.google-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.google-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* Divider */
.divider {
  position: relative;
  text-align: center;
  margin: 24px 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(to right, transparent, #e5e7eb 20%, #e5e7eb 80%, transparent);
}

.divider span {
  background: white;
  padding: 0 16px;
  color: var(--text-light);
  font-size: 14px;
  font-weight: 500;
  position: relative;
  z-index: 1;
}

/* Form Inputs */
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

.forgot-password {
  text-align: right;
  margin-bottom: 20px;
}

.forgot-password a {
  color: var(--primary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s ease;
}

.forgot-password a:hover {
  text-decoration: underline;
  opacity: 0.8;
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

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
  transition: color 0.2s ease;
}

.form-footer a:hover {
  text-decoration: underline;
  opacity: 0.8;
}

/* Password input styling */
.password-input-wrapper {
  position: relative;
  align-items: center;
}

.password-input-wrapper input {
  width: 100%;
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

/* Enhanced animations and effects */
.btn.secondary {
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
}

.btn.secondary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* Responsive Design */
@media (max-width: 480px) {
  .form-container {
    padding: 24px;
    margin: 10px;
  }
  
  .google-btn {
    font-size: 15px;
    padding: 10px 20px;
  }
  
  h2 {
    font-size: 24px;
  }
  
  .input-group {
    margin-right: 0;
  }
}

/* Subtle entrance animation */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.form-container {
  animation: fadeInUp 0.6s ease-out;
}

/* Loading state for Google button */
.google-btn.loading {
  pointer-events: none;
  opacity: 0.7;
}

.google-btn.loading::after {
  content: '';
  position: absolute;
  width: 16px;
  height: 16px;
  margin: auto;
  border: 2px solid transparent;
  border-top-color: #4285F4;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>