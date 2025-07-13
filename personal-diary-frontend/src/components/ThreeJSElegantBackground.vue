<template>
  <div ref="threeContainer" class="three-js-background"></div>
</template>

<script>
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'; // Optional, for development/testing

export default {
  name: 'ThreeJSElegantBackground',
  data() {
    return {
      scene: null,
      camera: null,
      renderer: null,
      animateId: null, // To store the animation frame ID for cleanup
      objects: [], // Array to hold our 3D objects
      controls: null, // For OrbitControls if used
    };
  },
  mounted() {
    this.initThreeJs();
    this.addElegantObjects();
    this.animate();
    window.addEventListener('resize', this.onWindowResize);
  },
  beforeUnmount() {
    // Cleanup Three.js resources
    window.removeEventListener('resize', this.onWindowResize);
    cancelAnimationFrame(this.animateId);
    this.renderer.dispose();
    this.scene.clear();
    // You might need to dispose of geometries and materials if you create many
    // objects or load complex models. For simple objects, renderer.dispose() helps.
  },
  methods: {
    initThreeJs() {
      const container = this.$refs.threeContainer;
      const width = container.clientWidth;
      const height = container.clientHeight;

      // Scene
      this.scene = new THREE.Scene();

      // Camera
      this.camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
      this.camera.position.z = 5;

      // Renderer
      this.renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true }); // alpha: true for transparent background
      this.renderer.setSize(width, height);
      this.renderer.setPixelRatio(window.devicePixelRatio); // For sharper rendering on high-DPI screens
      container.appendChild(this.renderer.domElement);

      // Optional: OrbitControls for easy camera manipulation during development
      // this.controls = new OrbitControls(this.camera, this.renderer.domElement);
      // this.controls.enableDamping = true; // For a smoother feel
      // this.controls.dampingFactor = 0.05;

      // Lighting - subtle ambient light
      const ambientLight = new THREE.AmbientLight(0xffffff, 0.6);
      this.scene.add(ambientLight);

      // Directional light for subtle shadows/highlights
      const directionalLight = new THREE.DirectionalLight(0xffffff, 0.4);
      directionalLight.position.set(0, 10, 5);
      this.scene.add(directionalLight);
    },

    addElegantObjects() {
      // Define a material for glowing, transparent spheres
      const material = new THREE.MeshPhongMaterial({
        color: 0x7b68ee, // Using your primary color
        emissive: 0x7b68ee, // Emit light of the same color
        emissiveIntensity: 0.3, // Subtle glow
        specular: 0xaaaaaa,
        shininess: 50,
        transparent: true,
        opacity: 0.7,
        side: THREE.DoubleSide,
      });

      // Create a few spheres
      const sphereGeometry = new THREE.SphereGeometry(0.5, 32, 32);
      const numSpheres = 5;
      for (let i = 0; i < numSpheres; i++) {
        const sphere = new THREE.Mesh(sphereGeometry, material);
        sphere.position.x = (Math.random() - 0.5) * 10;
        sphere.position.y = (Math.random() - 0.5) * 10;
        sphere.position.z = (Math.random() - 0.5) * 5 - 2; // Slightly behind content
        sphere.userData.speedX = (Math.random() - 0.5) * 0.005;
        sphere.userData.speedY = (Math.random() - 0.5) * 0.005;
        this.scene.add(sphere);
        this.objects.push(sphere);
      }
    },

    animate() {
      this.animateId = requestAnimationFrame(this.animate);

      // this.controls.update(); // Only if OrbitControls are enabled

      // Animate our elegant objects
      this.objects.forEach(obj => {
        obj.position.x += obj.userData.speedX;
        obj.position.y += obj.userData.speedY;

        // Wrap around effect
        if (obj.position.x > 5 || obj.position.x < -5) obj.userData.speedX *= -1;
        if (obj.position.y > 5 || obj.position.y < -5) obj.userData.speedY *= -1;

        // Subtle rotation
        obj.rotation.x += 0.001;
        obj.rotation.y += 0.002;
      });

      this.renderer.render(this.scene, this.camera);
    },

    onWindowResize() {
      const container = this.$refs.threeContainer;
      if (!container) return; // Ensure container exists

      const width = container.clientWidth;
      const height = container.clientHeight;

      this.camera.aspect = width / height;
      this.camera.updateProjectionMatrix();
      this.renderer.setSize(width, height);
      this.renderer.setPixelRatio(window.devicePixelRatio);
    },
  },
};
</script>

<style scoped>
.three-js-background {
  position: absolute; /* Position absolutely to act as a background */
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0; /* Place it behind other content */
  pointer-events: none; /* Allow clicks to pass through to elements below */
  overflow: hidden; /* Prevent scrollbars if elements slightly outside */
}
</style>