<template>
  <Teleport to="body">
    <Transition name="modal" appear>
      <div v-if="visible" class="modal-overlay">
        <div class="modal-backdrop" @click="$emit('close')"></div>
        <div class="modal-container">
          <div class="modal-content">
            <!-- Header -->
            <div class="modal-header">
              <div class="modal-title">
                <div class="title-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L13.09 8.26L20 9L13.09 9.74L12 16L10.91 9.74L4 9L10.91 8.26L12 2Z" fill="currentColor"/>
                    <path d="M19 15L20.09 18.26L23 19L20.09 19.74L19 23L17.91 19.74L15 19L17.91 18.26L19 15Z" fill="currentColor"/>
                  </svg>
                </div>
                <h3>AI Text Refinement</h3>
              </div>
              <button class="close-button" @click="$emit('close')" aria-label="Close">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M18 6L6 18M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </button>
            </div>

            <!-- Content -->
            <div class="modal-body">
              <div v-if="loading" class="loading-state">
                <div class="loading-spinner">
                  <div class="spinner-ring"></div>
                  <div class="spinner-ring"></div>
                  <div class="spinner-ring"></div>
                </div>
                <p class="loading-text">Refining your text with AI...</p>
                <div class="loading-progress">
                  <div class="progress-bar"></div>
                </div>
              </div>

              <div v-else class="refined-content">
                <div class="content-comparison">
                  <div class="comparison-section">
                    <label class="section-label">Original Text</label>
                    <div class="text-preview original-text">{{ originalText || 'No original text provided' }}</div>
                  </div>
                  
                  <div class="comparison-arrow">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M5 12H19M19 12L12 5M19 12L12 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </div>
                  
                  <div class="comparison-section">
                    <label class="section-label">Refined Text</label>
                    <div class="text-preview refined-text">{{ refinedText || 'Refinement will appear here' }}</div>
                  </div>
                </div>

                <div class="improvement-tags" v-if="!loading && refinedText">
                  <span class="tag">Grammar Enhanced</span>
                  <span class="tag">Clarity Improved</span>
                  <span class="tag">Tone Adjusted</span>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="modal-footer">
              <div class="action-buttons">
                <button 
                  @click="$emit('close')" 
                  class="btn btn-secondary"
                  :disabled="loading"
                >
                  Cancel
                </button>
                <button 
                  @click="$emit('replace', refinedText)" 
                  class="btn btn-primary"
                  :disabled="loading || !refinedText"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M20 6L9 17L4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Replace Text
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script>
export default {
  name: 'RefineModal',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    refinedText: {
      type: String,
      default: ''
    },
    originalText: {
      type: String,
      default: ''
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['close', 'replace'],
  mounted() {
    // Prevent body scroll when modal is open
    if (this.visible) {
      document.body.style.overflow = 'hidden'
    }
  },
  beforeUnmount() {
    // Restore body scroll
    document.body.style.overflow = ''
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = ''
      }
    }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
  box-sizing: border-box;
}

.modal-backdrop {
  position: absolute;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
}

.modal-container {
  position: relative;
  z-index: 1001;
  width: 100%;
  max-width: 700px;
  max-height: 80vh;
  overflow: hidden;
}

.modal-content {
  background: white;
  border-radius: 16px;
  box-shadow: 
    0 25px 50px -12px rgba(0, 0, 0, 0.25),
    0 0 0 1px rgba(255, 255, 255, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Header */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px 20px;
  border-bottom: 1px solid #f1f5f9;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
}

.modal-title h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  font-family: 'Inter', sans-serif;
}

.close-button {
  width: 36px;
  height: 36px;
  border: none;
  background: #f8fafc;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #64748b;
  transition: all 0.2s ease;
}

.close-button:hover {
  background: #e2e8f0;
  color: #475569;
}

/* Body */
.modal-body {
  padding: 28px;
  flex: 1;
  overflow-y: auto;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.loading-spinner {
  position: relative;
  width: 60px;
  height: 60px;
  margin-bottom: 24px;
}

.spinner-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 3px solid transparent;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1.5s linear infinite;
}

.spinner-ring:nth-child(2) {
  animation-delay: -0.5s;
  border-top-color: #764ba2;
}

.spinner-ring:nth-child(3) {
  animation-delay: -1s;
  border-top-color: #f093fb;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 16px;
  color: #64748b;
  margin: 0 0 20px 0;
  font-weight: 500;
}

.loading-progress {
  width: 200px;
  height: 4px;
  background: #e2e8f0;
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
  animation: progress 2s ease-in-out infinite;
}

@keyframes progress {
  0% {
    width: 0%;
    margin-left: 0%;
  }
  50% {
    width: 75%;
    margin-left: 25%;
  }
  100% {
    width: 0%;
    margin-left: 100%;
  }
}

.refined-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.content-comparison {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 20px;
  align-items: start;
}

.comparison-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-label {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #64748b;
}

.text-preview {
  padding: 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  min-height: 80px;
  max-height: 200px;
  overflow-y: auto;
  word-wrap: break-word;
}

.original-text {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #7f1d1d;
}

.refined-text {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #14532d;
}

.comparison-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  margin-top: 20px;
}

.improvement-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  background: linear-gradient(135deg, #667eea20, #764ba220);
  color: #4c1d95;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid #667eea30;
}

/* Footer */
.modal-footer {
  padding: 20px 28px 28px;
  border-top: 1px solid #f1f5f9;
  background: #fafbfc;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  font-family: 'Inter', sans-serif;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 120px;
  justify-content: center;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: white;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

/* Transitions */
.modal-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.modal-leave-active {
  transition: all 0.2s ease-in;
}

.modal-enter-from {
  opacity: 0;
  transform: scale(0.8) translateY(20px);
}

.modal-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-10px);
}

/* Responsive */
@media (max-width: 768px) {
  .modal-overlay {
    padding: 16px;
  }
  
  .modal-header {
    padding: 20px 20px 16px;
  }
  
  .modal-body {
    padding: 20px;
  }
  
  .modal-footer {
    padding: 16px 20px 20px;
  }
  
  .content-comparison {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .comparison-arrow {
    transform: rotate(90deg);
    margin: 0;
  }
  
  .action-buttons {
    flex-direction: column-reverse;
  }
  
  .btn {
    width: 100%;
  }
}

/* Dark mode support */
@media (prefers-color-scheme: dark) {
  .modal-content {
    background: #1e293b;
    color: #f1f5f9;
  }
  
  .modal-header {
    background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
    border-bottom-color: #334155;
  }
  
  .modal-title h3 {
    color: #f1f5f9;
  }
  
  .close-button {
    background: #334155;
    color: #94a3b8;
  }
  
  .close-button:hover {
    background: #475569;
    color: #cbd5e1;
  }
  
  .modal-footer {
    background: #0f172a;
    border-top-color: #334155;
  }
  
  .section-label {
    color: #94a3b8;
  }
  
  .original-text {
    background: #4c1d1d;
    border-color: #7f1d1d;
    color: #fca5a5;
  }
  
  .refined-text {
    background: #14532d;
    border-color: #166534;
    color: #86efac;
  }
  
  .btn-secondary {
    background: #334155;
    color: #cbd5e1;
    border-color: #475569;
  }
  
  .btn-secondary:hover:not(:disabled) {
    background: #475569;
    border-color: #64748b;
  }
}
</style>