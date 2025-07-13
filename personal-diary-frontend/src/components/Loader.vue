<template>
  <div class="loader-overlay" v-if="loading">
    <svg class="svg-loader" viewBox="0 0 50 50">
      <circle class="path" cx="25" cy="25" r="20" fill="none" stroke-width="5"></circle>
    </svg>
  </div>
</template>


<script setup>
import { computed } from 'vue'
import { useLoaderStore } from '../stores/loader'


const loaderStore = useLoaderStore()
const loading = computed(() => loaderStore.loading)
</script>

<style scoped>

.loader-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(255, 255, 255, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.loader {
  display: inline-block;
  width: 80px;
  height: 80px;
}

.loader:after {
  content: " ";
  display: block;
  width: 64px;
  height: 64px;
  margin: 8px;
  border-radius: 50%;
  border: 6px solid #3498db;
  border-color: #3498db transparent #3498db transparent;
  animation: loader-dual-ring 1.2s linear infinite;
}

@keyframes loader-dual-ring {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.svg-loader {
  animation: rotate 2s linear infinite;
  width: 60px;
  height: 60px;
}

.path {
  stroke: #3498db;
  stroke-linecap: round;
  animation: dash 1.5s ease-in-out infinite;
}

@keyframes rotate {
  100% {
    transform: rotate(360deg);
  }
}

@keyframes dash {
  0% {
    stroke-dasharray: 1, 150;
    stroke-dashoffset: 0;
  }
  50% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -35;
  }
  100% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -124;
  }
}
</style>
