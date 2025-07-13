// // src/utils/alert.js
// import { toast } from 'vue3-toastify';
// import 'vue3-toastify/dist/index.css';

// let confirmDialogInstance = null; // to register custom confirm dialog

// export const alert = {
//   success(message) {
//     toast.success(message, {
//       position: toast.POSITION.TOP_RIGHT,
//     });
//   },

//   error(message) {
//     toast.error(message, {
//       position: toast.POSITION.TOP_RIGHT,
//     });
//   },

//   warning(message) {
//     toast.warning(message, {
//       position: toast.POSITION.TOP_RIGHT,
//     });
//   },

//   info(message) {
//     toast.info(message, {
//       position: toast.POSITION.TOP_RIGHT,
//     });
//   },

//   custom(message, options = {}) {
//     toast(message, options);
//   },

//   async confirm(message) {
//     if (!confirmDialogInstance) {
//       throw new Error("Confirm dialog not registered!");
//     }
//     return await confirmDialogInstance(message);
//   },

//   // Call this once during app init to register the confirm function
//   registerConfirmHandler(handlerFn) {
//     confirmDialogInstance = handlerFn;
//   }
// };




// // alert.js
// import VueSimpleAlert from 'vue-simple-alert'


// // Create your custom alert functions
// const AlertService = {
//   // Simple alert message
//   showAlert(message, title = '') {
//     return Vue.prototype.$alert(message, title)
//   },
  
//   // Success message with predefined title
//   showSuccess(message) {
//     return Vue.prototype.$alert(message, 'Success')
//   },
  
//   // Error message with predefined title
//   showError(message) {
//     return Vue.prototype.$alert(message, 'Error')
//   },
  
//   // Confirmation dialog with customizable options
//   showConfirm(message, title = 'Confirm') {
//     return Vue.prototype.$confirm(message, title)
//   },
  
//   // Prompt dialog
//   showPrompt(message, title = 'Input Required') {
//     return Vue.prototype.$prompt(message, title)
//   }
// }

// export default AlertService



// alert.js
// import { createApp, h } from 'vue'

// alert.js
class AlertService {
  constructor() {
    this.alertQueue = []
    this.isProcessingQueue = false
    
    // Inject styles once
    this.injectStyles()
  }

  // Inject required CSS styles
  injectStyles() {
    if (document.getElementById('vue-alert-styles')) return

    const styleEl = document.createElement('style')
    styleEl.id = 'vue-alert-styles'
    styleEl.textContent = `
      @keyframes vue-alert-fade-in {
        from { opacity: 0; transform: translateY(-10px); }
        to { opacity: 1; transform: translateY(0); }
      }
      
      @keyframes vue-alert-fade-out {
        from { opacity: 1; transform: translateY(0); }
        to { opacity: 0; transform: translateY(-10px); }
      }
      
      .vue-alert {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
        position: relative;
        border-radius: 8px;
        box-shadow: 0 5px 15px rgba(0, 0, 0, 0.08), 0 1px 3px rgba(0, 0, 0, 0.1);
        overflow: hidden;
        max-width: 400px;
        width: 100%;
        margin-bottom: 10px;
        animation: vue-alert-fade-in 0.2s ease forwards;
        pointer-events: auto;
      }
      
      .vue-alert-modal {
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15), 0 5px 10px rgba(0, 0, 0, 0.12);
        margin-bottom: 0;
        width: 95%;
        max-width: 420px;
      }
      
      .vue-alert-content {
        display: flex;
        padding: 16px;
        background-color: #ffffff;
      }
      
      .vue-alert-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 12px;
      }
      
      .vue-alert-icon svg {
        flex-shrink: 0;
      }
      
      .vue-alert-body {
        flex: 1;
      }
      
      .vue-alert-title {
        font-weight: 600;
        font-size: 15px;
        margin-bottom: 4px;
        color: #111827;
      }
      
      .vue-alert-message {
        font-size: 14px;
        line-height: 1.5;
        color: #4B5563;
      }
      
      .vue-alert-progress {
        height: 3px;
        width: 100%;
        position: absolute;
        bottom: 0;
        left: 0;
        transform-origin: left;
      }
      
      .vue-alert-close {
        position: absolute;
        top: 12px;
        right: 12px;
        background: transparent;
        border: none;
        color: #9CA3AF;
        cursor: pointer;
        font-size: 18px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 24px;
        height: 24px;
        border-radius: 4px;
        transition: all 0.2s;
      }
      
      .vue-alert-close:hover {
        background-color: rgba(0, 0, 0, 0.05);
        color: #4B5563;
      }
      
      .vue-alert-input {
        border: 1px solid #E5E7EB;
        border-radius: 6px;
        padding: 10px 12px;
        font-size: 14px;
        width: 100%;
        margin-top: 12px;
        outline: none;
        transition: border-color 0.2s ease;
        box-sizing: border-box;
      }
      
      .vue-alert-input:focus {
        border-color: #3B82F6;
        box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
      }
      
      .vue-alert-buttons {
        display: flex;
        justify-content: flex-end;
        margin-top: 16px;
        gap: 8px;
      }
      
      .vue-alert-button {
        padding: 8px 16px;
        border-radius: 6px;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.15s ease;
        border: 1px solid transparent;
      }
      
      .vue-alert-button-cancel {
        background-color: #F9FAFB;
        color: #4B5563;
        border-color: #E5E7EB;
      }
      
      .vue-alert-button-cancel:hover {
        background-color: #F3F4F6;
      }
      
      .vue-alert-button-primary {
        background-color: #3B82F6;
        color: white;
      }
      
      .vue-alert-button-primary:hover {
        background-color: #2563EB;
      }
      
      .vue-alert-info { border-top: 3px solid #3B82F6; }
      .vue-alert-success { border-top: 3px solid #10B981; }
      .vue-alert-error { border-top: 3px solid #EF4444; }
      .vue-alert-confirm { border-top: 3px solid #6366F1; }
      .vue-alert-prompt { border-top: 3px solid #8B5CF6; }
      
      .vue-alert-info .vue-alert-progress { background-color: #3B82F6; }
      .vue-alert-success .vue-alert-progress { background-color: #10B981; }
      .vue-alert-error .vue-alert-progress { background-color: #EF4444; }
    `
    document.head.appendChild(styleEl)
  }

  // Create modal container for confirm/prompt dialogs
  ensureModalContainer() {
    let container = document.getElementById('vue-alert-modal-container')
    if (!container) {
      container = document.createElement('div')
      container.id = 'vue-alert-modal-container'
      container.style.position = 'fixed'
      container.style.top = '0'
      container.style.left = '0'
      container.style.right = '0'
      container.style.bottom = '0'
      container.style.zIndex = '9999'
      container.style.display = 'flex'
      container.style.justifyContent = 'center'
      container.style.alignItems = 'center'
      container.style.backgroundColor = 'rgba(0, 0, 0, 0.4)'
      container.style.opacity = '0'
      container.style.transition = 'opacity 0.2s ease'
      container.style.pointerEvents = 'auto'
      document.body.appendChild(container)
      
      // Fade in the backdrop
      setTimeout(() => {
        container.style.opacity = '1'
      }, 10)
    }
    return container
  }

  // Create alert container if it doesn't exist
  ensureContainer() {
    let container = document.getElementById('vue-alert-container')
    if (!container) {
      container = document.createElement('div')
      container.id = 'vue-alert-container'
      container.style.position = 'fixed'
      container.style.top = '20px'
      container.style.right = '20px'
      container.style.zIndex = '9999'
      container.style.display = 'flex'
      container.style.flexDirection = 'column'
      container.style.alignItems = 'flex-end'
      container.style.maxWidth = '400px'
      container.style.width = '100%'
      container.style.pointerEvents = 'none'
      document.body.appendChild(container)
    }
    return container
  }

  // Process alerts in queue
  async processQueue() {
    if (this.isProcessingQueue || this.alertQueue.length === 0) return
    
    this.isProcessingQueue = true
    const alert = this.alertQueue.shift()
    await this.showAlertElement(alert)
    
    this.isProcessingQueue = false
    this.processQueue()
  }

  // Get icon for each alert type
  getIconSVG(type) {
    const colors = {
      info: '#3B82F6',
      success: '#10B981',
      error: '#EF4444',
      confirm: '#6366F1',
      prompt: '#8B5CF6'
    }
    
    switch(type) {
      case 'error':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${colors.error}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>`;
      case 'success':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${colors.success}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`;
      case 'info':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${colors.info}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`;
      case 'confirm':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${colors.confirm}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`;
      case 'prompt':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${colors.prompt}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>`;
      default:
        return `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${colors.info}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`;
    }
  }

  // Show the actual alert element
  showAlertElement(config) {
    return new Promise(resolve => {
      // For confirm and prompt, use a centered modal container
      const isModal = config.type === 'confirm' || config.type === 'prompt';
      const container = isModal ? this.ensureModalContainer() : this.ensureContainer();
      
      // Create alert element
      const alertEl = document.createElement('div')
      alertEl.className = `vue-alert vue-alert-${config.type}${isModal ? ' vue-alert-modal' : ''}`
      
      // Create alert content
      const contentEl = document.createElement('div')
      contentEl.className = 'vue-alert-content'
      alertEl.appendChild(contentEl)
      
      // Icon
      const iconEl = document.createElement('div')
      iconEl.className = 'vue-alert-icon'
      iconEl.innerHTML = this.getIconSVG(config.type)
      contentEl.appendChild(iconEl)
      
      // Body
      const bodyEl = document.createElement('div')
      bodyEl.className = 'vue-alert-body'
      contentEl.appendChild(bodyEl)
      
      // Title
      if (config.title) {
        const titleEl = document.createElement('div')
        titleEl.className = 'vue-alert-title'
        titleEl.textContent = config.title
        bodyEl.appendChild(titleEl)
      }
      
      // Message
      const messageEl = document.createElement('div')
      messageEl.className = 'vue-alert-message'
      messageEl.textContent = config.message
      bodyEl.appendChild(messageEl)
      
      // Close button
      const closeBtn = document.createElement('button')
      closeBtn.className = 'vue-alert-close'
      closeBtn.innerHTML = '&times;'
      closeBtn.onclick = () => {
        this.removeAlert(alertEl, container, resolve, config.reject)
      }
      contentEl.appendChild(closeBtn)
      
      // For confirm and prompt, add the necessary elements
      if (config.type === 'confirm' || config.type === 'prompt') {
        // Input field for prompt
        let inputEl
        if (config.type === 'prompt') {
          inputEl = document.createElement('input')
          inputEl.type = 'text'
          inputEl.className = 'vue-alert-input'
          inputEl.placeholder = 'Enter your response'
          bodyEl.appendChild(inputEl)
          
          // Focus the input field
          setTimeout(() => {
            inputEl.focus()
          }, 100)
        }
        
        // Buttons container
        const btnContainer = document.createElement('div')
        btnContainer.className = 'vue-alert-buttons'
        bodyEl.appendChild(btnContainer)
        
        // Cancel button
        const cancelBtn = document.createElement('button')
        cancelBtn.className = 'vue-alert-button vue-alert-button-cancel'
        cancelBtn.textContent = 'Cancel'
        cancelBtn.onclick = () => {
          this.removeAlert(alertEl, container, resolve, config.reject)
        }
        btnContainer.appendChild(cancelBtn)
        
        // OK button
        const okBtn = document.createElement('button')
        okBtn.className = 'vue-alert-button vue-alert-button-primary'
        okBtn.textContent = 'OK'
        okBtn.onclick = () => {
          if (config.type === 'prompt' && inputEl) {
            const value = inputEl.value
            this.removeAlert(alertEl, container, resolve, null, () => config.resolve(value))
          } else {
            this.removeAlert(alertEl, container, resolve, null, config.resolve)
          }
        }
        btnContainer.appendChild(okBtn)
        
        // Add keyboard event listeners for confirm/prompt dialogs
        const handleKeydown = (e) => {
          if (e.key === 'Enter') {
            okBtn.click()
            document.removeEventListener('keydown', handleKeydown)
          } else if (e.key === 'Escape') {
            cancelBtn.click()
            document.removeEventListener('keydown', handleKeydown)
          }
        }
        document.addEventListener('keydown', handleKeydown)
      }
      
      // Progress bar for timed alerts
      if (config.duration) {
        const progressEl = document.createElement('div')
        progressEl.className = 'vue-alert-progress'
        alertEl.appendChild(progressEl)
        
        progressEl.style.transform = 'scaleX(1)'
        progressEl.style.transition = `transform ${config.duration}ms linear`
        
        // Start progress animation
        setTimeout(() => {
          progressEl.style.transform = 'scaleX(0)'
        }, 10)
        
        // Auto close after duration
        setTimeout(() => {
          if (container.contains(alertEl)) {
            this.removeAlert(alertEl, container, resolve)
          }
        }, config.duration)
      }
      
      // Add to container
      container.appendChild(alertEl)
    })
  }
  
  // Remove alert with animation
  removeAlert(alertEl, container, resolve, reject = null, callback = null) {
    alertEl.style.animation = 'vue-alert-fade-out 0.2s ease forwards'
    
    // Check if this is a modal container
    const isModal = container.id === 'vue-alert-modal-container'
    
    if (isModal) {
      // Fade out the backdrop
      container.style.opacity = '0'
    }
    
    // Remove after animation completes
    setTimeout(() => {
      if (container.contains(alertEl)) {
        container.removeChild(alertEl)
        
        // If this was the last alert in a modal container, remove the container
        if (isModal && container.children.length === 0) {
          document.body.removeChild(container)
        }
        
        if (callback) callback()
        if (reject) reject()
        resolve()
      }
    }, 200)
  }

  // Basic alert
  alert(message, title = '', duration = 3000) {
    return new Promise(resolve => {
      this.alertQueue.push({
        type: 'info',
        message,
        title,
        duration,
        resolve
      })
      this.processQueue()
    })
  }

  // Success alert
  success(message, title = 'Success', duration = 3000) {
    return new Promise(resolve => {
      this.alertQueue.push({
        type: 'success',
        message,
        title,
        duration,
        resolve
      })
      this.processQueue()
    })
  }

  // Error alert
  error(message, title = 'Error', duration = 3000) {
    return new Promise(resolve => {
      this.alertQueue.push({
        type: 'error',
        message,
        title,
        duration,
        resolve
      })
      this.processQueue()
    })
  }

  // Confirmation dialog
  confirm(message, title = 'Confirm') {
    return new Promise((resolve, reject) => {
      this.alertQueue.push({
        type: 'confirm',
        message,
        title,
        duration: 0,
        resolve,
        reject
      })
      this.processQueue()
    })
  }

  // Prompt dialog
  prompt(message, title = 'Input Required') {
    return new Promise((resolve, reject) => {
      this.alertQueue.push({
        type: 'prompt',
        message,
        title,
        duration: 0,
        resolve,
        reject
      })
      this.processQueue()
    })
  }
}

// Create a singleton instance
const alertService = new AlertService()

// Export a Vue plugin
export default {
  install: (app) => {
    app.config.globalProperties.$alert = alertService.alert.bind(alertService)
    app.config.globalProperties.$alertSuccess = alertService.success.bind(alertService)
    app.config.globalProperties.$alertError = alertService.error.bind(alertService)
    app.config.globalProperties.$confirm = alertService.confirm.bind(alertService)
    app.config.globalProperties.$prompt = alertService.prompt.bind(alertService)
    
    // Also make it available via injection
    app.provide('alertService', alertService)
  }
}

// Export the direct service for use outside of components
export { alertService }