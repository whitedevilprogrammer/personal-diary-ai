<template>
  <div class="dashboard">
    <div class="logout-container">
      <button class="logout-btn" @click="logout">Logout</button>
    </div>

    <div class="dashboard-header">
      <h1>📘 My Diary Entries</h1>
      <p>Record your thoughts, reflect on your day</p>
    </div>

    <div class="dashboard-content">
      <div class="entry-form card">
        <h2>Create New Entry</h2>
        <form @submit.prevent="createEntry">
          <div class="input-group" style="position:relative;">
            <label for="title">Title</label>
            <input id="title" v-model="title" placeholder="Entry title" required />
          </div>
          <div class="input-group" style="position:relative;">
            <label for="content">Content</label>
            <textarea
              id="content"
              ref="entryTextarea"
              v-model="content"
              placeholder="Write your thoughts..."
              rows="5"
              required
              @mouseup="showRefineButton"
              @keyup="showRefineButton"
              @input="autoResizeTextarea($event)"
            ></textarea>
            <RefineButton
              :visible="refineVisible"
              :position="refinePosition"
              @refine="openRefineModal"
            />
            <RefineModal
              :visible="refineModalVisible"
              :refinedText="refinedText"
              :originalText="originalSelectedText"
              :loading="loading"
              @replace="replaceText"
              @close="refineModalVisible = false"
            />
          </div>
          
          <!-- Background Image Generation Section -->
          <div class="background-section">
            <div class="background-header">
              <label>Background Image</label>
              <button 
                type="button" 
                class="btn-generate" 
                @click="generateBackgroundImage"
                :disabled="generatingImage || !content.trim()"
              >
                <span v-if="generatingImage" class="generating">
                  <span class="spinner"></span>
                  Generating...
                </span>
                <span v-else>✨ Generate Background</span>
              </button>
            </div>
            
            <div v-if="generatedImageUrl" class="image-preview">
              <img :src="generatedImageUrl" alt="Generated background" />
              <div class="image-actions">
                <button type="button" class="btn-icon success" @click="useGeneratedImage" title="Use this image">
                  ✓
                </button>
                <button type="button" class="btn-icon danger" @click="removeGeneratedImage" title="Remove image">
                  ✕
                </button>
              </div>
            </div>
            
            <div v-if="selectedBackgroundUrl" class="selected-background">
              <p class="selected-label">Selected background:</p>
              <div class="selected-preview">
                <img :src="selectedBackgroundUrl" alt="Selected background" />
                <button type="button" class="btn-icon danger" @click="removeSelectedBackground" title="Remove background">
                  ✕
                </button>
              </div>
            </div>
          </div>
          
          <button class="btn">Save Entry</button>
        </form>
      </div>

      <div class="entries-container">
        <div class="entries-header">
          <h2>Your Entries</h2>
          <div v-if="entries.length === 0" class="no-entries">
            <p>No entries yet. Create your first one!</p>
          </div>
        </div>

        <transition-group name="entry-list" tag="div" class="entries-list">
          <div
            v-for="entry in entries"
            :key="entry._id"
            class="entry-card card"
            :style="getEntryCardStyle(entry)"
          >
            <div v-if="editingId === entry._id" class="edit-mode">
              <div class="edit-form-header">
                <h2>Edit Entry</h2>
                <div class="edit-actions">
                  <button class="btn-icon success" @click="saveEdit(entry._id)" title="Save">
                    ✓
                  </button>
                  <button class="btn-icon cancel" @click="cancelEdit" title="Cancel">
                    ✕
                  </button>
                </div>
              </div>
              <form @submit.prevent="saveEdit(entry._id)">
                <div class="input-group">
                  <label for="edit-title">Title</label>
                  <input
                    id="edit-title"
                    v-model="editForm.title"
                    placeholder="Entry title"
                    required
                    @keyup.esc="cancelEdit"
                  />
                </div>
                <div class="input-group" style="position:relative;">
                  <label for="edit-content">Content</label>
                  <textarea
                    id="edit-content"
                    :ref="`editEntryTextarea_${entry._id}`"
                    v-model="editForm.content"
                    placeholder="Write your thoughts..."
                    rows="5"
                    required
                    @keyup.ctrl.enter="saveEdit(entry._id)"
                    @keyup.esc="cancelEdit"
                    @mouseup="showEditRefineButton(entry._id)"
                    @keyup="showEditRefineButton(entry._id)"
                    @input="autoResizeEditTextarea($event, entry._id)"
                  ></textarea>
                  <RefineButton
                    :visible="editRefineVisible && editingId === entry._id"
                    :position="editRefinePosition"
                    @refine="openEditRefineModal(entry._id)"
                  />
                  <RefineModal
                    :visible="editRefineModalVisible && editingId === entry._id"
                    :refinedText="editRefinedText"
                    :originalText="originalEditSelectedText"
                    :loading="editLoading"
                    @replace="replaceEditText"
                    @close="editRefineModalVisible = false"
                  />
                </div>
                
                <!-- Background Image Edit Section -->
                <div class="background-section">
                  <div class="background-header">
                    <label>Background Image</label>
                    <button 
                      type="button" 
                      class="btn-generate" 
                      @click="generateEditBackgroundImage(entry._id)"
                      :disabled="generatingEditImage || !editForm.content.trim()"
                    >
                      <span v-if="generatingEditImage" class="generating">
                        <span class="spinner"></span>
                        Generating...
                      </span>
                      <span v-else>✨ Generate New Background</span>
                    </button>
                  </div>
                  
                  <div v-if="editGeneratedImageUrl" class="image-preview">
                    <img :src="editGeneratedImageUrl" alt="Generated background" />
                    <div class="image-actions">
                      <button type="button" class="btn-icon success" @click="useEditGeneratedImage" title="Use this image">
                        ✓
                      </button>
                      <button type="button" class="btn-icon danger" @click="removeEditGeneratedImage" title="Remove image">
                        ✕
                      </button>
                    </div>
                  </div>
                  
                  <div v-if="editForm.backgroundImageUrl" class="selected-background">
                    <p class="selected-label">Current background:</p>
                    <div class="selected-preview">
                      <img :src="editForm.backgroundImageUrl" alt="Current background" />
                      <button type="button" class="btn-icon danger" @click="removeEditSelectedBackground" title="Remove background">
                        ✕
                      </button>
                    </div>
                  </div>
                </div>
                
                <div class="edit-form-actions">
                  <button type="submit" class="btn">Save Changes</button>
                  <button type="button" class="btn btn-outline" @click="cancelEdit">Cancel</button>
                </div>
              </form>
              <div class="edit-footer">
                <span class="edit-hint">Press Ctrl+Enter to save, Esc to cancel</span>
              </div>
            </div>

            <div v-else class="view-mode">
              <div class="entry-header">
                <h3>{{ entry.title }}</h3>
                <div class="entry-actions">
                  <button class="btn-icon edit" @click="startEdit(entry)" title="Edit">✏️</button>
                  <button class="btn-icon danger" @click="deleteEntry(entry._id)" title="Delete">🗑</button>
                </div>
              </div>
              <div class="entry-content">{{ entry.content }}</div>
              <div class="entry-footer">
                <span class="entry-date">{{ formatDate(entry.createdAt) }}</span>
              </div>
              
              <!-- Background overlay for text readability -->
              <div v-if="entry.backgroundImageUrl" class="entry-overlay"></div>
            </div>
          </div>
        </transition-group>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'
import { alertService } from '../utils/alert'
import { useLoaderStore } from '../stores/loader'
import RefineButton from '../components/RefineButton.vue'
import RefineModal from '../components/RefineModal.vue'

export default {
  components: {
    RefineButton,
    RefineModal
  },
  data() {
    return {
      title: '',
      content: '',
      entries: [],
      loaderStore: null,
      editingId: null,
      editForm: {
        title: '',
        content: '',
        backgroundImageUrl: ''
      },
      // New Entry Refine
      refineVisible: false,
      refinePosition: { x: 0, y: 0 },
      selectedText: '',
      refinedText: '',
      refineModalVisible: false,
      loading: false,
      // Edit Entry Refine
      editRefineVisible: false,
      editRefinePosition: { x: 0, y: 0 },
      editSelectedText: '',
      editRefinedText: '',
      editRefineModalVisible: false,
      editLoading: false,
      originalSelectedText: '',
      originalEditSelectedText: '',
      // Background Image Generation
      generatingImage: false,
      generatedImageUrl: '',
      selectedBackgroundUrl: '',
      // Edit Background Image Generation
      generatingEditImage: false,
      editGeneratedImageUrl: '',
    }
  },
  created() {
    const token = new URLSearchParams(window.location.search).get('token')
    if (!localStorage.getItem('token')) {
      if (token) {
        localStorage.setItem('token', token)
        const cleanUrl = window.location.origin + '/dashboard'
        window.history.replaceState({}, document.title, cleanUrl)
        this.$router.push('/dashboard')
      } else {
        this.$router.push('/')
      }
    }
    this.loaderStore = useLoaderStore()
  },
  async mounted() {
    await this.fetchEntries()

    document.addEventListener('selectionchange', this.handleSelectionChange)
    document.addEventListener('click', this.handleDocumentClick)
  },
  beforeUnmount() {
    document.removeEventListener('selectionchange', this.handleSelectionChange)
    document.removeEventListener('click', this.handleDocumentClick)
  },
  methods: {
    async fetchEntries() {
      this.loaderStore.show()
      try {
        const res = await api.get('/diary')
        this.entries = res.data
      } catch (err) {
        this.logout()
        alertService.error('Failed to fetch entries. Please login again.', 'Session Expired')
      } finally {
        this.loaderStore.hide()
      }
    },
    async createEntry() {
      if (!this.title || !this.content) {
        alertService.error('Please fill in all fields', 'Error')
        return
      }
      this.loaderStore.show()
      try {
        const payload = {
          title: this.title,
          content: this.content,
        }

        if (this.selectedBackgroundUrl) {
          payload.backgroundImageUrl = this.selectedBackgroundUrl
        }

        const res = await api.post('/diary', payload)
        if (res.status === 'success') {
          alertService.success(res.message, 'Create Entry')
        } else {
          alertService.error(res.message, 'Error')
        }
      } catch (err) {
        alertService.error('Failed to create entry', 'Error')
      } finally {
        this.loaderStore.hide()
        this.title = ''
        this.content = ''
        this.selectedBackgroundUrl = ''
        this.generatedImageUrl = ''
        // Ensure entry textarea is reset if it had dynamic height
        if (this.$refs.entryTextarea) {
          this.$refs.entryTextarea.style.height = 'auto'; // Or your default height
          this.$nextTick(() => { // Ensure DOM is updated before resizing
            this.resizeTextarea(this.$refs.entryTextarea);
          });
        }
        await this.fetchEntries()
      }
    },
    async deleteEntry(id) {
      this.loaderStore.show()
      try {
        await alertService.confirm('Are you sure you want to delete this item?', 'Confirm Delete')
        await api.delete(`/diary/${id}`)
        alertService.success('Entry deleted successfully', 'Delete Entry') // Changed to success type
        await this.fetchEntries()
      } catch (e) {
        // Cancelled or error
        if (e && e.message !== 'Dialog cancelled by user') { // Example check if alertService throws specific error
          alertService.error('Failed to delete entry', 'Error') // Add error alert for actual errors
        }
      } finally {
        this.loaderStore.hide()
      }
    },
    startEdit(entry) {
      this.editingId = entry._id
      this.editForm = {
        title: entry.title,
        content: entry.content,
        backgroundImageUrl: entry.backgroundImageUrl || ''
      }
      this.editRefineVisible = false
      this.editRefineModalVisible = false
      this.editGeneratedImageUrl = ''

      this.$nextTick(() => {
        const textareaRef = this.$refs[`editEntryTextarea_${entry._id}`];
        const textarea = Array.isArray(textareaRef) ? textareaRef[0] : textareaRef;
        if (textarea) {
          this.resizeTextarea(textarea);
        }
      });
    },

    resizeTextarea(textarea) {
      if (!textarea) return;
      textarea.style.height = 'auto';
      const minHeight = 120;
      const maxHeight = 400;
      const newHeight = Math.min(Math.max(textarea.scrollHeight, minHeight), maxHeight);
      textarea.style.height = newHeight + 'px';
      textarea.style.overflowY = textarea.scrollHeight > maxHeight ? 'auto' : 'hidden';
    },

    autoResizeTextarea(event) {
      this.resizeTextarea(event.target);
    },

    autoResizeEditTextarea(event) { // Removed entryId as it's not directly used for resize logic
      this.resizeTextarea(event.target);
    },

    async saveEdit(id) {
      if (!this.editForm.title.trim() || !this.editForm.content.trim()) {
        alertService.error('Please fill in all fields', 'Error')
        return
      }
      this.loaderStore.show()
      try {
        const payload = {
          title: this.editForm.title,
          content: this.editForm.content,
        }

        // Only include backgroundImageUrl if it's present, otherwise omit it to not change existing
        // If you want to explicitly *remove* a background, you'd set it to null or an empty string based on your backend API.
        // The current logic means if editForm.backgroundImageUrl is empty, it won't be sent in the payload.
        // If the backend has a default for missing fields, this is fine. If it explicitly needs null/empty string to clear, adjust here.
        if (this.editForm.backgroundImageUrl) {
          payload.backgroundImageUrl = this.editForm.backgroundImageUrl
        } else {
            // Explicitly set to null if it's empty in the form and we want to clear it on the backend
            payload.backgroundImageUrl = null;
        }

        const res = await api.put(`/diary/${id}`, payload)
        if (res.status === 'success') {
          alertService.success(res.message || 'Entry updated successfully', 'Update Entry')
          this.cancelEdit()
          await this.fetchEntries()
        } else {
          alertService.error(res.message, 'Error')
        }
      } catch (err) {
        alertService.error('Failed to update entry', 'Error')
      } finally {
        this.loaderStore.hide()
      }
    },
    cancelEdit() {
      this.editingId = null
      this.editForm = { title: '', content: '', backgroundImageUrl: '' }
      this.editRefineVisible = false
      this.editRefineModalVisible = false
      this.editGeneratedImageUrl = ''
    },
    formatDate(dateString) {
      const options = { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }
      return new Date(dateString).toLocaleDateString(undefined, options)
    },
    logout() {
      localStorage.removeItem('token')
      this.$router.push('/')
      alertService.success('Logged out successfully', 'Logout') // Changed to success type
    },

    async generateBackgroundImage() {
      if (!this.content.trim()) {
        alertService.error('Please write some content first to generate a background image.', 'Error');
        return;
      }

      this.generatingImage = true;
      this.generatedImageUrl = ''; // Reset previous generated image

      try {
        const response = await api.post('/gemini/generate-background', {
          content: this.content,
          title: this.title || 'Diary Entry'
        });

        if (response && response.status === 'success' && response.data && typeof response.data.imageUrl === 'string' && response.data.imageUrl.trim() !== '') {
          this.generatedImageUrl = response.data.imageUrl;
          alertService.success('Background image generated successfully!', 'Image Generated');
        } else {
          console.warn('Background generation API call succeeded but response format was unexpected or imageUrl missing:', response);
          let errorMessage = 'Failed to get a valid image URL from the server.';
          if (response && response.message) {
            errorMessage = response.message;
          } else if (response && response.data && response.data.message) {
            errorMessage = response.data.message;
          }
          throw new Error(errorMessage);
        }
      } catch (error) {
        console.error('Background generation error:', error);
        let displayMessage = 'Failed to generate background image. Please try again.';
        if (error.response && error.response.data && error.response.data.message) {
          displayMessage = error.response.data.message;
        } else if (error.message) {
          displayMessage = error.message;
        }
        alertService.error(displayMessage, 'Generation Error');
        this.generatedImageUrl = ''; // Ensure no broken preview if generation fails
      } finally {
        this.generatingImage = false;
      }
    },

    useGeneratedImage() {
      this.selectedBackgroundUrl = this.generatedImageUrl
      this.generatedImageUrl = ''
      alertService.success('Background image selected!', 'Image Selected')
    },

    removeGeneratedImage() {
      this.generatedImageUrl = ''
    },

    removeSelectedBackground() {
      this.selectedBackgroundUrl = ''
    },

    async generateEditBackgroundImage() { // Removed entryId as it's not directly used for generation logic
      if (!this.editForm.content.trim()) {
        alertService.error('Please write some content first to generate a background image.', 'Error');
        return;
      }

      this.generatingEditImage = true;
      this.editGeneratedImageUrl = ''; // Reset previous generated image

      try {
        const response = await api.post('/gemini/generate-background', {
          content: this.editForm.content,
          title: this.editForm.title || 'Diary Entry'
        });

        if (response && response.status === 'success' && response.data && typeof response.data.imageUrl === 'string' && response.data.imageUrl.trim() !== '') {
          this.editGeneratedImageUrl = response.data.imageUrl;
          alertService.success('Background image generated successfully!', 'Image Generated');
        } else {
          console.warn('Edit background generation API call succeeded but response format was unexpected or imageUrl missing:', response);
          let errorMessage = 'Failed to get a valid image URL from the server for editing.';
          if (response && response.message) {
            errorMessage = response.message;
          } else if (response && response.data && response.data.message) {
            errorMessage = response.data.message;
          }
          throw new Error(errorMessage);
        }
      } catch (error) {
        console.error('Edit background generation error:', error);
        let displayMessage = 'Failed to generate background image for editing. Please try again.';
        if (error.response && error.response.data && error.response.data.message) {
          displayMessage = error.response.data.message;
        } else if (error.message) {
          displayMessage = error.message;
        }
        alertService.error(displayMessage, 'Generation Error');
        this.editGeneratedImageUrl = ''; // Ensure no broken preview if generation fails
      } finally {
        this.generatingEditImage = false;
      }
    },

    useEditGeneratedImage() {
      this.editForm.backgroundImageUrl = this.editGeneratedImageUrl
      this.editGeneratedImageUrl = ''
      alertService.success('Background image selected!', 'Image Selected')
    },

    removeEditGeneratedImage() {
      this.editGeneratedImageUrl = ''
    },

    removeEditSelectedBackground() {
      this.editForm.backgroundImageUrl = ''
    },

    getEntryCardStyle(entry) {
      if (entry.backgroundImageUrl) {
        return {
          backgroundImage: `linear-gradient(rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.85)), url(${entry.backgroundImageUrl})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundRepeat: 'no-repeat',
          position: 'relative'
        }
      }
      return {}
    },

    handleSelectionChange() {
      // Delay to ensure the selection is fully updated by the browser
      setTimeout(() => {
        const selection = window.getSelection();
        if (selection.rangeCount === 0 || selection.isCollapsed) {
          this.refineVisible = false;
          this.editRefineVisible = false;
          return;
        }

        const activeElement = document.activeElement;

        // Check if the selection is within the new entry textarea
        if (this.$refs.entryTextarea && this.$refs.entryTextarea === activeElement) {
          const start = this.$refs.entryTextarea.selectionStart;
          const end = this.$refs.entryTextarea.selectionEnd;
          if (start !== end) {
            this.showRefineButton();
          } else {
            this.refineVisible = false;
          }
          this.editRefineVisible = false; // Hide edit refine button if new entry textarea is active
        }
        // Check if the selection is within an edit entry textarea
        else {
          let foundEditArea = false;
          for (let key in this.$refs) {
            if (key.startsWith('editEntryTextarea_')) {
              const textareaRef = this.$refs[key];
              const textarea = Array.isArray(textareaRef) ? textareaRef[0] : textareaRef;
              if (textarea && textarea === activeElement) {
                const start = textarea.selectionStart;
                const end = textarea.selectionEnd;
                if (start !== end) {
                  // Ensure we pass the entry ID to showEditRefineButton
                  const entryId = key.split('_')[1];
                  this.showEditRefineButton(entryId);
                } else {
                  this.editRefineVisible = false;
                }
                foundEditArea = true;
                this.refineVisible = false; // Hide new entry refine button if edit textarea is active
                break;
              }
            }
          }
          if (!foundEditArea) {
            // If neither textarea is active, hide both buttons
            this.refineVisible = false;
            this.editRefineVisible = false;
          }
        }
      }, 10);
    },

    handleDocumentClick(event) {
      const entryTextarea = this.$refs.entryTextarea
      const editTextareas = Object.keys(this.$refs).filter(key => key.startsWith('editEntryTextarea_'))

      let clickedInTextarea = false

      if (entryTextarea && entryTextarea.contains(event.target)) {
        clickedInTextarea = true
      }

      if (!clickedInTextarea) { // Only check edit textareas if not already in new entry textarea
        for (let refKey of editTextareas) {
          const textareaRef = this.$refs[refKey]
          const textarea = Array.isArray(textareaRef) ? textareaRef[0] : textareaRef
          if (textarea && textarea.contains(event.target)) {
            clickedInTextarea = true
            break
          }
        }
      }

      // If clicked on a refine button or inside a refine modal, don't hide
      if (event.target.closest('.refine-button-component-class') || event.target.closest('.refine-modal-component-class')) { // Use actual classes of your components
        return;
      }

      if (!clickedInTextarea) {
        this.refineVisible = false
        this.editRefineVisible = false
      }
    },

    getTextareaCoordinates(textarea, selectionStart, selectionEnd) {
      const text = textarea.value;
      const textBeforeSelection = text.substring(0, selectionStart);
      // Create a temporary span to measure the text and get its coordinates
      const span = document.createElement('span');
      span.style.whiteSpace = 'pre-wrap'; // Important for multi-line text
      span.style.visibility = 'hidden';
      span.style.position = 'absolute';
      span.style.display = 'inline-block'; // To allow getBoundingClientRect to work
      // Copy computed styles for accurate measurement
      const computedStyle = window.getComputedStyle(textarea);
      for (let prop of ['font', 'fontSize', 'fontFamily', 'lineHeight', 'paddingTop', 'paddingLeft', 'borderTopWidth', 'borderLeftWidth', 'width']) {
        span.style[prop] = computedStyle[prop];
      }
      span.style.boxSizing = 'border-box'; // Match textarea's box model

      span.textContent = textBeforeSelection;
      document.body.appendChild(span);

      // Get the bounding rectangle of the span to find the cursor's position
      const spanRect = span.getBoundingClientRect();
      const textareaRect = textarea.getBoundingClientRect();

      document.body.removeChild(span);

      // Calculate the X and Y coordinates relative to the textarea's top-left corner
      // Adjust for textarea's padding and border, and scroll position
      const x = spanRect.left - textareaRect.left + textarea.scrollLeft;
      const y = spanRect.bottom - textareaRect.top + textarea.scrollTop;

      // Position the button above the selected text, with an offset
      const buttonOffsetX = 0; // Adjust as needed
      const buttonOffsetY = 40; // Height of the button + some margin

      return {
        x: x + buttonOffsetX,
        y: y - buttonOffsetY
      };
    },

    showRefineButton() {
      const textarea = this.$refs.entryTextarea;
      if (!textarea) return;

      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;

      if (start === end) {
        this.refineVisible = false;
        return;
      }

      const currentSelectedText = textarea.value.substring(start, end);
      if (currentSelectedText.trim().length === 0) {
        this.refineVisible = false;
        return;
      }

      this.selectedText = currentSelectedText; // This is used by openRefineModal
      this.originalSelectedText = currentSelectedText; // For the modal display

      const coords = this.getTextareaCoordinates(textarea, start, end);

      this.refinePosition = {
        x: coords.x,
        y: coords.y
      };

      this.refineVisible = true;
    },

    async openRefineModal() {
      this.refineModalVisible = true;
      this.refineVisible = false; // Hide the button once modal is open
      this.loading = true;
      this.refinedText = '';

      try {
        const response = await api.post('/diary/refine', {
          text: this.selectedText, // Use the stored selected text
          context: 'diary_entry',
          tone: 'professional'
        });

        if (response.status === 'success' && response.data && response.data.refinedText) {
          this.refinedText = response.data.refinedText;
        } else {
          throw new Error(response.message || response.data?.message || 'Failed to refine text');
        }
      } catch (error) {
        console.error('Refine API error:', error);
        alertService.error(
          error.response?.data?.message || error.message || 'Failed to refine text. Please try again.',
          'Refinement Error'
        );
        this.refineModalVisible = false;
      } finally {
        this.loading = false;
      }
    },

    replaceText(newText) {
      const textarea = this.$refs.entryTextarea;
      if (!textarea) return;

      // We need to use the selection start and end that were active when the refine button was clicked,
      // or re-evaluate the selection if the content might have changed.
      // For a robust solution, store `selectionStart` and `selectionEnd` when `showRefineButton` is called.
      // For now, let's re-find the `originalSelectedText` in the `content`.
      const currentContent = this.content;
      const indexOfSelected = currentContent.indexOf(this.originalSelectedText);

      if (indexOfSelected !== -1) {
        const start = indexOfSelected;
        const end = indexOfSelected + this.originalSelectedText.length;

        this.content =
          currentContent.slice(0, start) + newText + currentContent.slice(end);
        this.refineModalVisible = false;

        this.$nextTick(() => {
          textarea.focus();
          textarea.setSelectionRange(start, start + newText.length);
          this.autoResizeTextarea({ target: textarea }); // Resize after content change
        });
      } else {
        alertService.error('Could not find original selected text to replace. Please re-select.', 'Replacement Error');
        this.refineModalVisible = false; // Close modal as replacement failed
      }
    },

    showEditRefineButton(entryId) {
      // Ensure we are targeting the correct textarea for the currently editing entry
      const textareaRef = this.$refs[`editEntryTextarea_${this.editingId}`];
      const textarea = Array.isArray(textareaRef) ? textareaRef[0] : textareaRef;

      if (!textarea) return;

      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;

      if (start === end) {
        this.editRefineVisible = false;
        return;
      }

      const currentSelectedText = textarea.value.substring(start, end);
      if (currentSelectedText.trim().length === 0) {
        this.editRefineVisible = false;
        return;
      }

      this.editSelectedText = currentSelectedText; // For openEditRefineModal
      this.originalEditSelectedText = currentSelectedText; // For modal display

      const coords = this.getTextareaCoordinates(textarea, start, end);

      this.editRefinePosition = {
        x: coords.x,
        y: coords.y
      };

      this.editRefineVisible = true;
    },

    async openEditRefineModal() { // entryId is not needed here as this.editingId is already set
      this.editRefineModalVisible = true;
      this.editRefineVisible = false; // Hide the button once modal is open
      this.editLoading = true;
      this.editRefinedText = '';

      try {
        const response = await api.post('/diary/refine', {
          text: this.editSelectedText, // Use the stored selected text
          context: 'diary_entry',
          tone: 'professional'
        });

        if (response.status === 'success' && response.data && response.data.refinedText) {
          this.editRefinedText = response.data.refinedText;
        } else {
          throw new Error(response.message || response.data?.message || 'Failed to refine text');
        }
      } catch (error) {
        console.error('Edit refine API error:', error);
        alertService.error(
          error.response?.data?.message || error.message || 'Failed to refine text. Please try again.',
          'Refinement Error'
        );
        this.editRefineModalVisible = false;
      } finally {
        this.editLoading = false;
      }
    },

    replaceEditText(newText) {
      const textareaRef = this.$refs[`editEntryTextarea_${this.editingId}`];
      const textarea = Array.isArray(textareaRef) ? textareaRef[0] : textareaRef;

      if (!textarea) return;

      // Similar to `replaceText`, we need to find the specific range to replace.
      const currentContent = this.editForm.content;
      const indexOfSelected = currentContent.indexOf(this.originalEditSelectedText);

      if (indexOfSelected !== -1) {
        const start = indexOfSelected;
        const end = indexOfSelected + this.originalEditSelectedText.length;

        this.editForm.content =
          currentContent.slice(0, start) + newText + currentContent.slice(end);
        this.editRefineModalVisible = false;

        this.$nextTick(() => {
          textarea.focus();
          textarea.setSelectionRange(start, start + newText.length);
          this.autoResizeEditTextarea({ target: textarea }); // Resize after content change
        });
      } else {
        alertService.error('Could not find original selected text to replace in edit mode. Please re-select.', 'Replacement Error');
        this.editRefineModalVisible = false;
      }
    }
  }
}
</script>

<style scoped>
/* Your existing styles */
.dashboard {
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 20px;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 40px;
}

.dashboard-header h1 {
  font-size: 32px;
  margin-bottom: 8px;
  background: linear-gradient(135deg, var(--text-dark) 30%, var(--primary) 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.dashboard-header p {
  color: var(--text-light);
  font-size: 18px;
}

.dashboard-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 30px;
}

@media (max-width: 768px) {
  .dashboard-content {
    grid-template-columns: 1fr;
  }
}

.entry-form {
  padding: 30px;
}

.entry-form h2 {
  margin-bottom: 20px;
  font-size: 24px;
}

.input-group {
  margin-bottom: 16px;
  position: relative; /* Ensure this is set for RefineButton positioning */
}

input, textarea {
  width: -webkit-fill-available;
  transition: height 0.2s ease;
  resize: none; /* Disable manual resize since we're auto-resizing */
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.entries-container {
  display: flex;
  flex-direction: column;
}

.entries-header {
  margin-bottom: 20px;
}

.entries-header h2 {
  font-size: 24px;
  margin-bottom: 10px;
}

.no-entries {
  text-align: center;
  padding: 40px 0;
  color: var(--text-light);
}

.entries-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.entry-card {
  padding: 20px;
}

.entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.entry-actions {
  display: flex;
  gap: 8px;
}

.card {
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.06);
  padding: 20px;
  transition: box-shadow 0.3s ease;
}

.card:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
}

.entry-content {
  white-space: pre-wrap;
  line-height: 1.6;
  color: var(--text-dark);
  margin-bottom: 12px;
  overflow-wrap: anywhere;
}

.entry-footer {
  display: flex;
  justify-content: flex-end;
  font-size: 14px;
  color: var(--text-light);
}

/* Edit Mode Styles */
.edit-mode {
  animation: editModeEnter 0.3s ease;
}

@keyframes editModeEnter {
  from {
    background-color: rgba(108, 92, 231, 0.02);
    transform: scale(0.998);
  }
  to {
    background-color: transparent;
    transform: scale(1);
  }
}

.edit-form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.edit-form-header h2 {
  margin: 0;
  font-size: 24px;
}

.edit-actions {
  display: flex;
  gap: 8px;
}

.edit-form-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.edit-footer {
  margin-top: 12px;
  text-align: right;
}

.edit-hint {
  font-size: 12px;
  color: var(--text-light);
  font-style: italic;
}

.entry-list-enter-active,
.entry-list-leave-active {
  transition: all 0.5s ease;
}

.entry-list-enter-from,
.entry-list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

:root {
  --primary: #6c5ce7;
  --danger: #e74c3c;
  --success: #00b894;
  --text-dark: #2d3436;
  --text-light: #636e72;
}

.entry-header h3 {
  margin: 0;
  font-size: 20px;
}

/* Elegant Button Styles */
.btn {
  font-family: 'Inter', sans-serif;
  font-weight: 500;
  font-size: 16px;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, var(--primary) 0%, #8e44ad 100%);
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(142, 68, 173, 0.2);
  position: relative;
  overflow: hidden;
  letter-spacing: 0.5px;
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(142, 68, 173, 0.3);
}

.btn:active {
  transform: translateY(1px);
  box-shadow: 0 2px 8px rgba(142, 68, 173, 0.2);
}

.btn:after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.btn:hover:after {
  opacity: 1;
}

/* Secondary/outline button variant */
.btn-outline {
  background: transparent;
  border: 2px solid var(--primary);
  color: var(--primary);
  box-shadow: 0 4px 12px rgba(142, 68, 173, 0.1);
}

.btn-outline:hover {
  background: rgba(142, 68, 173, 0.05);
}

/* Elegant Logout Button Styling */
.logout-container {
  display: flex;
  justify-content: flex-end;
  width: 100%;
  padding: 16px 0;
  position: relative;
  z-index: 10;
}

.logout-btn {
  font-family: 'Inter', 'Helvetica', sans-serif;
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.6px;
  text-transform: uppercase;
  padding: 10px 20px;
  color: #5d4777;
  background: transparent;
  border: 1px solid rgba(93, 71, 119, 0.15);
  border-radius: 4px;
  transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.02);
}

.logout-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(93, 71, 119, 0.05) 0%, rgba(93, 71, 119, 0.01) 100%);
  opacity: 0;
  transition: opacity 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
}

.logout-btn:hover {
  border-color: rgba(93, 71, 119, 0.3);
  color: #4a2f70;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(93, 71, 119, 0.06);
}

.logout-btn:hover::before {
  opacity: 1;
}

.logout-btn:active {
  transform: translateY(0);
  box-shadow: 0 1px 3px rgba(93, 71, 119, 0.03);
}

/* Icon button styling refinements */
.btn-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid #eee;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  cursor: pointer;
}

.btn-icon:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.btn-icon.danger {
  color: var(--danger);
}

.btn-icon.danger:hover {
  background-color: rgba(231, 76, 60, 0.08);
}

.btn-icon.edit {
  color: var(--primary);
}

.btn-icon.edit:hover {
  background-color: rgba(108, 92, 231, 0.08);
}

.btn-icon.success {
  color: var(--success);
}

.btn-icon.success:hover {
  background-color: rgba(0, 184, 148, 0.08);
}

.btn-icon.cancel {
  color: var(--text-light);
}

.btn-icon.cancel:hover {
  background-color: rgba(99, 110, 114, 0.08);
}

.refine-button {
  position: absolute;
  /* Add these styles for the refine button */
  z-index: 10;
  background-color: var(--primary);
  color: white;
  padding: 6px 12px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.refine-button:hover {
  transform: translateY(-2px);
}

.refine-button:active {
  transform: translateY(0);
}

.refine-button .refine-icon {
  display: flex;
  align-items: center;
}

.refine-button svg {
  fill: currentColor;
}

/* Transition for the refine button */
.refine-button-enter-active,
.refine-button-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.refine-button-enter-from,
.refine-button-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

/* Background Image Generation Section Styles */
.background-section {
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, rgba(108, 92, 231, 0.02) 0%, rgba(142, 68, 173, 0.01) 100%);
  border: 1px solid rgba(108, 92, 231, 0.08);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.background-section:hover {
  border-color: rgba(108, 92, 231, 0.12);
  background: linear-gradient(135deg, rgba(108, 92, 231, 0.03) 0%, rgba(142, 68, 173, 0.02) 100%);
}

.background-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.background-header label {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.background-header label::before {
  content: '🎨';
  font-size: 18px;
}

/* Generate Button Styling */
.btn-generate {
  font-family: 'Inter', sans-serif;
  font-weight: 500;
  font-size: 14px;
  padding: 10px 18px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, var(--primary) 0%, #8e44ad 100%);
  color: white;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.165, 0.84, 0.44, 1);
  box-shadow: 0 3px 10px rgba(142, 68, 173, 0.25);
  position: relative;
  overflow: hidden;
  letter-spacing: 0.3px;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 160px;
  justify-content: center;
}

.btn-generate:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(142, 68, 173, 0.35);
}

.btn-generate:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(142, 68, 173, 0.25);
}

.btn-generate:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: 0 3px 10px rgba(142, 68, 173, 0.15);
}

/* Generating State */
.generating {
  display: flex;
  align-items: center;
  gap: 8px;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Image Preview Styling */
.image-preview {
  margin-top: 16px;
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease;
  animation: imageAppear 0.5s ease-out;
}

@keyframes imageAppear {
  0% {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.image-preview:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.image-preview img {
  width: 100%;
  height: 200px;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}

.image-preview:hover img {
  transform: scale(1.02);
}

.image-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.image-preview:hover .image-actions {
  opacity: 1;
}

/* Selected Background Styling */
.selected-background {
  margin-top: 16px;
  padding: 16px;
  background: rgba(0, 184, 148, 0.05);
  border: 1px solid rgba(0, 184, 148, 0.15);
  border-radius: 12px;
  animation: selectedAppear 0.4s ease-out;
}

@keyframes selectedAppear {
  0% {
    opacity: 0;
    transform: scale(0.98) translateY(10px);
    background: rgba(0, 184, 148, 0.1);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
    background: rgba(0, 184, 148, 0.05);
  }
}

.selected-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--success);
  margin: 0 0 12px 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.selected-label::before {
  content: '✓';
  background: var(--success);
  color: white;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}

.selected-preview {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 184, 148, 0.15);
}

.selected-preview img {
  width: 100%;
  height: 120px;
  object-fit: cover;
  display: block;
}

.selected-preview .btn-icon {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.selected-preview:hover .btn-icon {
  opacity: 1;
}

/* Entry Card Background Styling */
.entry-card {
  position: relative;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  overflow: hidden;
}

.entry-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.95) 0%,
    rgba(255, 255, 255, 0.85) 50%,
    rgba(255, 255, 255, 0.9) 100%
  );
  pointer-events: none;
  border-radius: 12px;
}

/* Ensure content is above overlay */
.entry-card .view-mode {
  position: relative;
  z-index: 2;
}

.entry-card .edit-mode {
  position: relative;
  z-index: 2;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 20px;
  backdrop-filter: blur(10px);
}

/* Enhanced button icon styles for image actions */
.btn-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s cubic-bezier(0.165, 0.84, 0.44, 1);
  cursor: pointer;
  backdrop-filter: blur(10px);
  font-size: 14px;
  font-weight: 600;
}

.btn-icon:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.btn-icon:active {
  transform: translateY(0) scale(1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* Responsive Design */
@media (max-width: 768px) {
  .background-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .btn-generate {
    width: 100%;
    min-width: unset;
  }
  
  .image-preview img,
  .selected-preview img {
    height: 150px;
  }
  
  .background-section {
    padding: 16px;
  }
}

/* Focus states for accessibility */
.btn-generate:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(108, 92, 231, 0.3), 0 6px 20px rgba(142, 68, 173, 0.35);
}

.btn-icon:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(108, 92, 231, 0.3), 0 6px 16px rgba(0, 0, 0, 0.2);
}

/* Loading state enhancement */
.btn-generate:disabled .spinner {
  animation: spin 1s linear infinite, pulse 2s ease-in-out infinite alternate;
}

@keyframes pulse {
  0% { opacity: 0.8; }
  100% { opacity: 1; }
}
</style>