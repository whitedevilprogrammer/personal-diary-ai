<template>
  <div class="form-page">
    <div class="form-container card">
      <div class="form-header">
        <div class="form-icon">🔑</div>
        <h2>Reset Your Password</h2>
        <p>Enter your email address and we'll send you a link to reset your password</p>
      </div>

      <!-- Step 1: Email Input -->
      <form v-if="step === 1" @submit.prevent="sendResetEmail">
        <div class="input-group">
          <label for="email">Email Address</label>
          <input 
            id="email" 
            v-model="email" 
            placeholder="your@email.com" 
            type="email" 
            required 
            :disabled="isLoading"
          />
        </div>

        <button type="submit" class="btn secondary" :disabled="isLoading || !email">
          <span v-if="isLoading">Sending...</span>
          <span v-else>Send Reset Link</span>
        </button>
      </form>

      <!-- Step 2: Success Message -->
      <div v-if="step === 2" class="success-message">
        <div class="success-icon">✅</div>
        <h3>Check Your Email</h3>
        <p>We've sent a password reset link to <strong>{{ email }}</strong></p>
        <p class="small-text">
          Didn't receive the email? Check your spam folder or 
          <button type="button" class="link-button" @click="resendEmail" :disabled="resendCooldown > 0">
            <span v-if="resendCooldown > 0">resend in {{ resendCooldown }}s</span>
            <span v-else>click here to resend</span>
          </button>
        </p>
      </div>

      <!-- Step 3: Reset Password Form (when accessed via reset token) -->
      <form v-if="step === 3" @submit.prevent="resetPassword">
        <div class="input-group">
          <label for="newPassword">New Password</label>
          <div class="password-input-wrapper">
            <input 
              id="newPassword" 
              v-model="newPassword" 
              placeholder="Enter new password"
              :type="showNewPassword ? 'text' : 'password'" 
              required 
              minlength="6"
              :disabled="isLoading"
            />
            <button 
              type="button" 
              class="password-toggle" 
              @click="showNewPassword = !showNewPassword" 
              tabindex="-1"
              aria-label="Toggle password visibility"
            >
              <span v-if="showNewPassword">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
              </span>
              <span v-else>
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                  <line x1="1" y1="1" x2="23" y2="23"></line>
                </svg>
              </span>
            </button>
          </div>
        </div>

        <div class="input-group">
          <label for="confirmPassword">Confirm New Password</label>
          <div class="password-input-wrapper">
            <input 
              id="confirmPassword" 
              v-model="confirmPassword" 
              placeholder="Confirm new password"
              :type="showConfirmPassword ? 'text' : 'password'" 
              required 
              minlength="6"
              :disabled="isLoading"
            />
            <button 
              type="button" 
              class="password-toggle" 
              @click="showConfirmPassword = !showConfirmPassword" 
              tabindex="-1"
              aria-label="Toggle password visibility"
            >
              <span v-if="showConfirmPassword">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
              </span>
              <span v-else>
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                  <line x1="1" y1="1" x2="23" y2="23"></line>
                </svg>
              </span>
            </button>
          </div>
        </div>

        <div v-if="passwordError" class="error-message">
          {{ passwordError }}
        </div>

        <button type="submit" class="btn secondary" :disabled="isLoading || !isPasswordValid">
          <span v-if="isLoading">Resetting...</span>
          <span v-else>Reset Password</span>
        </button>
      </form>

      <div class="form-footer">
        <p>
          <router-link to="/login">← Back to Login</router-link>
        </p>
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
      step: 1, // 1: email input, 2: success message, 3: reset form
      email: '',
      newPassword: '',
      confirmPassword: '',
      showNewPassword: false,
      showConfirmPassword: false,
      isLoading: false,
      resetToken: null,
      resendCooldown: 0,
      loaderStore: null,
    }
  },
  computed: {
    isPasswordValid() {
      return this.newPassword.length >= 6 && 
             this.confirmPassword.length >= 6 && 
             this.newPassword === this.confirmPassword
    },
    passwordError() {
      if (!this.newPassword && !this.confirmPassword) return null
      if (this.newPassword.length < 6) return 'Password must be at least 6 characters long'
      if (this.newPassword !== this.confirmPassword) return 'Passwords do not match'
      return null
    }
  },
  created() {
    this.loaderStore = useLoaderStore()
    this.checkResetToken()
  },
  methods: {
    checkResetToken() {
      // Check if accessing via reset token URL
      debugger
      const urlParams = new URLSearchParams(window.location.search)
      const token = urlParams.get('token')
      
      if (token) {
        this.resetToken = token
        this.step = 3
        this.validateResetToken()
      }
    },

    async validateResetToken() {
      try {
        this.isLoading = true
        const res = await api.post('/validate-reset-token', {
          token: this.resetToken
        })
        
        if (res.status !== 'success') {
          alertService.error('Invalid or expired reset link. Please request a new one.', 'Error')
          this.step = 1
          this.$router.push('/forgot-password')
        }
      } catch (error) {
        console.error('Token validation error:', error)
        alertService.error('Invalid reset link. Please request a new one.', 'Error')
        this.step = 1
        this.$router.push('/forgot-password')
      } finally {
        this.isLoading = false
      }
    },

    async sendResetEmail() {
      if (this.isLoading) return
      
      this.isLoading = true
      this.loaderStore.show()
      
      try {
        const res = await api.post('/forgot-password', {
          email: this.email
        })
        
        if (res.status === 'success') {
          this.step = 2
          alertService.success(res.message || 'Reset link sent to your email', 'Success')
        } else {
          alertService.error(res.message || 'Failed to send reset email', 'Error')
        }
      } catch (error) {
        console.error('Forgot password error:', error)
        alertService.error('Failed to send reset email. Please try again.', 'Error')
      } finally {
        this.isLoading = false
        this.loaderStore.hide()
      }
    },

    async resendEmail() {
      if (this.resendCooldown > 0) return
      
      await this.sendResetEmail()
      
      // Start cooldown
      this.resendCooldown = 60
      const interval = setInterval(() => {
        this.resendCooldown--
        if (this.resendCooldown <= 0) {
          clearInterval(interval)
        }
      }, 1000)
    },

    async resetPassword() {
      if (this.isLoading || !this.isPasswordValid) return
      
      this.isLoading = true
      this.loaderStore.show()
      
      try {
        const res = await api.post('/reset-password', {
          token: this.resetToken,
          password: this.newPassword,
          confirmPassword: this.confirmPassword
        })
        
        if (res.status === 'success') {
          alertService.success(res.message || 'Password reset successfully!', 'Success')
          // Clean URL and redirect to login
          window.history.replaceState({}, document.title, '/login')
          this.$router.push('/login')
        } else {
          alertService.error(res.message || 'Failed to reset password', 'Error')
        }
      } catch (error) {
        console.error('Reset password error:', error)
        alertService.error('Failed to reset password. Please try again.', 'Error')
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

input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s ease;
}


input:focus {
  outline: none;
  border-color: var(--primary);
}

input:disabled {
  background-color: #f9fafb;
  cursor: not-allowed;
  opacity: 0.6;
}

.btn {
  width: 100%;
  padding: 12px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 10px;
}

.btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}


.password-input-wrapper {
  position: relative;
  /* display: flex; */
  align-items: center;
}

/* .password-input-wrapper input {
  padding-right: 45px;
} */

.password-toggle {
  position: absolute;
  right: -16px;
  top: 14%;
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
  transition: opacity 0.2s;
  width: auto;
}

.password-toggle:hover {
  opacity: 0.7;
}

.success-message {
  text-align: center;
  padding: 20px 0;
}

.success-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.success-message h3 {
  color: var(--text-dark);
  margin-bottom: 12px;
  font-size: 24px;
}

.success-message p {
  color: var(--text-light);
  margin-bottom: 12px;
  line-height: 1.5;
}

.small-text {
  font-size: 14px;
  margin-top: 20px;
}

.link-button {
  background: none;
  border: none;
  color: var(--primary);
  cursor: pointer;
  padding: 0;
  font-size: inherit;
  text-decoration: underline;
  width: auto;
  margin: 0;
}

.link-button:hover:not(:disabled) {
  opacity: 0.8;
}

.link-button:disabled {
  color: var(--text-light);
  cursor: not-allowed;
  text-decoration: none;
}

.error-message {
  color: #dc2626;
  font-size: 14px;
  margin-top: 8px;
  margin-bottom: 16px;
  text-align: center;
  padding: 8px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
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

/* Animations */
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

/* Responsive Design */
@media (max-width: 480px) {
  .form-container {
    padding: 24px;
    margin: 10px;
  }
  
  h2 {
    font-size: 24px;
  }
  
  .success-icon {
    font-size: 36px;
  }
  
  .success-message h3 {
    font-size: 20px;
  }
}
</style>