<template>
  <Transition name="refine-button" appear>
    <div
      v-if="visible"
      class="refine-button"
      :style="{ top: `${position.y}px`, left: `${position.x}px` }"
      @click="refine"
      @mousedown.prevent
    >
      <div class="refine-icon">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 2L13.09 8.26L20 9L13.09 9.74L12 16L10.91 9.74L4 9L10.91 8.26L12 2Z" fill="currentColor"/>
          <path d="M19 15L20.09 18.26L23 19L20.09 19.74L19 23L17.91 19.74L15 19L17.91 18.26L19 15Z" fill="currentColor"/>
          <path d="M5 6L5.5 7.5L7 8L5.5 8.5L5 10L4.5 8.5L3 8L4.5 7.5L5 6Z" fill="currentColor"/>
        </svg>
      </div>
      <span class="refine-text">Refine</span>
      <div class="refine-shine"></div>
    </div>
  </Transition>
</template>

<script>
export default {
  name: 'RefineButton',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    position: {
      type: Object,
      default: () => ({ x: 0, y: 0 })
    }
  },
  emits: ['refine'],
  methods: {
    refine() {
      this.$emit('refine')
    }
  }
}
</script>

<style scoped>
.refine-button {
  position: absolute;
  z-index: 999;
  display: flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  box-shadow: 
    0 4px 12px rgba(102, 126, 234, 0.3),
    0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.refine-button:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 8px 20px rgba(102, 126, 234, 0.4),
    0 4px 8px rgba(0, 0, 0, 0.15);
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
}

.refine-button:active {
  transform: translateY(-1px);
  transition: all 0.1s ease;
}

.refine-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  animation: sparkle 2s ease-in-out infinite;
}

.refine-text {
  letter-spacing: 0.5px;
  font-weight: 600;
}

.refine-shine {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  transition: left 0.5s ease;
}

.refine-button:hover .refine-shine {
  left: 100%;
}

@keyframes sparkle {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}

/* Transition animations */
.refine-button-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.refine-button-leave-active {
  transition: all 0.2s ease-in;
}

.refine-button-enter-from {
  opacity: 0;
  transform: translateY(10px) scale(0.8);
}

.refine-button-leave-to {
  opacity: 0;
  transform: translateY(-5px) scale(0.9);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .refine-button {
    padding: 10px 14px;
    font-size: 14px;
    border-radius: 10px;
  }
  
  .refine-icon {
    width: 18px;
    height: 18px;
    margin-right: 2px;
  }
}

/* Dark mode support */
@media (prefers-color-scheme: dark) {
  .refine-button {
    background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
    border-color: rgba(255, 255, 255, 0.1);
  }
  
  .refine-button:hover {
    background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
  }
}
</style>