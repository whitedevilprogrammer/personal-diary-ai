<template>
  <div ref="threeContainer" class="three-scene-container"></div>
</template>

<script>
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
// Optional: If you want to use a more robust noise library, install it first:
// npm install simplex-noise
// import SimplexNoise from 'simplex-noise';

export default {
  name: 'ThreeScene',
  props: {
    baseColor: {
      type: String,
      default: '#7b68ee', // Your primary color
    },
    accentColor: {
      type: String,
      default: '#a48ffc', // A lighter, complementary shade for highlights
    },
    animationSpeed: {
      type: Number,
      default: 0.0005, // Slower for subtle, organic movement
    },
    blobScale: {
      type: Number,
      default: 0.8, // Initial size of the blob
    },
    displacementIntensity: {
      type: Number,
      default: 0.2, // How much the blob "breathes"
    }
  },
  
  // Three.js properties (non-reactive)
  scene: null,
  camera: null,
  renderer: null,
  blobMesh: null, // Renamed from cube
  controls: null,
  animationFrameId: null,
  noise: null, // For noise generation
  startTime: null, // For animation timing

  mounted() {
    this.startTime = Date.now(); // Initialize start time for animation
    // If using a library: this.noise = new SimplexNoise();
    this.initThreeScene();
    this.animate();
    window.addEventListener('resize', this.onWindowResize);
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.onWindowResize);
    if (this.animationFrameId) {
      cancelAnimationFrame(this.animationFrameId);
    }
    if (this.renderer) {
      this.renderer.dispose();
      if (this.$refs.threeContainer && this.renderer.domElement) {
        this.$refs.threeContainer.removeChild(this.renderer.domElement);
      }
    }
    if (this.scene) {
      this.scene.traverse((object) => {
        if (object.isMesh) {
          if (object.geometry) object.geometry.dispose();
          if (object.material) {
            for (const key of Object.keys(object.material)) {
              if (object.material[key] && typeof object.material[key].dispose === 'function') {
                object.material[key].dispose();
              }
            }
            object.material.dispose();
          }
        }
      });
      this.scene.clear();
    }
    if (this.controls) {
      this.controls.dispose();
    }
  },
  methods: {
    initThreeScene() {
      // 1. Scene
      this.scene = new THREE.Scene();
      this.scene.background = new THREE.Color(0xf8f9fa); // Matches dashboard background

      // 2. Camera
      const width = this.$refs.threeContainer.clientWidth;
      const height = this.$refs.threeContainer.clientHeight;
      this.camera = new THREE.PerspectiveCamera(
        50, // Slightly narrower FOV for a less distorted view
        width / height,
        0.1,
        1000
      );
      this.camera.position.z = 3; // Move camera back to view the larger blob

      // 3. Renderer
      this.renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true }); // alpha: true for transparent background
      this.renderer.setSize(width, height);
      this.renderer.setPixelRatio(window.devicePixelRatio);
      this.$refs.threeContainer.appendChild(this.renderer.domElement);

      // 4. Geometry and Material (The Blob)
      // Icosahedron: A good base for organic shapes
      const geometry = new THREE.IcosahedronGeometry(this.blobScale, 5); // Radius, subdivisions (higher for smoother blob)

      // Phong Material: Responds well to light, good for glossy/smooth surfaces
      const material = new THREE.MeshPhongMaterial({
        color: new THREE.Color(this.baseColor),
        shininess: 50, // Adjust for more or less specularity
        specular: new THREE.Color(this.accentColor), // Specular highlight color
        flatShading: false // Smoother shading
      });

      this.blobMesh = new THREE.Mesh(geometry, material);
      this.scene.add(this.blobMesh);

      // Store original positions for displacement
      this.blobMesh.originalPositions = new Float32Array(geometry.attributes.position.array);

      // 5. Lighting
      // Soft Ambient Light
      const ambientLight = new THREE.AmbientLight(0xffffff, 0.4); // Softer ambient light
      this.scene.add(ambientLight);

      // Key Light (main light source)
      const keyLight = new THREE.SpotLight(0xffffff, 1.2, 0, Math.PI / 4, 0.5, 2);
      keyLight.position.set(5, 5, 5);
      keyLight.castShadow = true; // Enable shadows (requires renderer.shadowMap.enabled = true)
      this.scene.add(keyLight);

      // Fill Light (softens shadows from key light)
      const fillLight = new THREE.PointLight(0xffffff, 0.6);
      fillLight.position.set(-5, -5, 5);
      this.scene.add(fillLight);

      // Back Light (for rim lighting and depth)
      const backLight = new THREE.DirectionalLight(0xffffff, 0.3);
      backLight.position.set(0, 5, -5).normalize();
      this.scene.add(backLight);

      // Renderer for shadows
      this.renderer.shadowMap.enabled = true;
      this.renderer.shadowMap.type = THREE.PCFSoftShadowMap; // Softer shadows

      // 6. Controls (Optional, for exploration)
      this.controls = new OrbitControls(this.camera, this.renderer.domElement);
      this.controls.enableDamping = true;
      this.controls.dampingFactor = 0.05;
      this.controls.enableZoom = false; // Disable zoom for a more fixed view
      this.controls.enablePan = false; // Disable pan
    },

    animate() {
      this.animationFrameId = requestAnimationFrame(this.animate);

      // Update blob displacement
      if (this.blobMesh && this.blobMesh.geometry && this.blobMesh.originalPositions) {
        const time = (Date.now() - this.startTime) * this.animationSpeed;
        const positions = this.blobMesh.geometry.attributes.position.array;
        const originalPositions = this.blobMesh.originalPositions;

        for (let i = 0; i < originalPositions.length; i += 3) {
          const x = originalPositions[i];
          const y = originalPositions[i + 1];
          const z = originalPositions[i + 2];

          // Use simple noise to displace vertices
          // You can replace this with a more advanced noise function if needed
          const noiseValue = this.simplexNoise3D(
            x * 1 + time,
            y * 1 + time,
            z * 1 + time
          );

          // Displace vertices along their normals (or just radially from center)
          // For simplicity, displacing radially from center (normalized original position)
          const vec = new THREE.Vector3(x, y, z).normalize();
          positions[i] = originalPositions[i] + vec.x * noiseValue * this.displacementIntensity;
          positions[i + 1] = originalPositions[i + 1] + vec.y * noiseValue * this.displacementIntensity;
          positions[i + 2] = originalPositions[i + 2] + vec.z * noiseValue * this.displacementIntensity;
        }

        this.blobMesh.geometry.attributes.position.needsUpdate = true; // Tell Three.js to update vertex positions
        this.blobMesh.geometry.computeVertexNormals(); // Recompute normals for correct lighting
      }

      // Smooth rotation for the entire blob
      if (this.blobMesh) {
          this.blobMesh.rotation.y += 0.001; // Slow, continuous rotation
          this.blobMesh.rotation.x += 0.0005;
      }


      if (this.controls) {
        this.controls.update(); // Only needed if damping is enabled
      }

      if (this.renderer && this.scene && this.camera) {
        this.renderer.render(this.scene, this.camera);
      }
    },

    onWindowResize() {
      if (this.camera && this.renderer && this.$refs.threeContainer) {
        const width = this.$refs.threeContainer.clientWidth;
        const height = this.$refs.threeContainer.clientHeight;
        this.camera.aspect = width / height;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(width, height);
        this.renderer.setPixelRatio(window.devicePixelRatio);
      }
    },

    // --- Simplex Noise Implementation (A basic 3D version) ---
    // This is a simplified version. For production, consider using a dedicated library like 'simplex-noise'.
    // Source: https://gist.github.com/banksean/3045220
    // Simplified for a quick example. More robust implementations exist.
    simplexNoise3D(x, y, z) {
        let n0, n1, n2, n3, n4, n5, n6;
        let p = new Array(512); // Permutation table for Simplex Noise
        let perm = new Array(512);
        for (let i = 0; i < 256; i++) {
            p[i] = Math.floor(Math.random() * 256);
            perm[i] = p[i];
            perm[i + 256] = p[i];
        }

        function grad3d(i, x, y, z) {
            let h = i & 15;
            let u = h < 8 ? x : y;
            let v = h < 4 ? y : (h === 12 || h === 14 ? x : z);
            return ((h & 1) === 0 ? u : -u) + ((h & 2) === 0 ? v : -v);
        }

        return (function() {
            let F3 = 1.0 / 3.0;
            let G3 = 1.0 / 6.0;

            return function(xin, yin, zin) {
                let s = (xin + yin + zin) * F3;
                let i = Math.floor(xin + s);
                let j = Math.floor(yin + s);
                let k = Math.floor(zin + s);

                let t = (i + j + k) * G3;
                let X0 = i - t;
                let Y0 = j - t;
                let Z0 = k - t;

                let x0 = xin - X0;
                let y0 = yin - Y0;
                let z0 = zin - Z0;

                let i1, j1, k1;
                let i2, j2, k2;

                if (x0 >= y0) {
                    if (y0 >= z0) { i1 = 1; j1 = 0; k1 = 0; i2 = 1; j2 = 1; k2 = 0; }
                    else if (x0 >= z0) { i1 = 1; j1 = 0; k1 = 0; i2 = 1; j2 = 0; k2 = 1; }
                    else { i1 = 0; j1 = 0; k1 = 1; i2 = 1; j2 = 0; k2 = 1; }
                } else {
                    if (y0 < z0) { i1 = 0; j1 = 0; k1 = 1; i2 = 0; j2 = 1; k2 = 1; }
                    else if (x0 < z0) { i1 = 0; j1 = 1; k1 = 0; i2 = 0; j2 = 1; k2 = 1; }
                    else { i1 = 0; j1 = 1; k1 = 0; i2 = 1; j2 = 1; k2 = 0; }
                }

                let x1 = x0 - i1 + G3;
                let y1 = y0 - j1 + G3;
                let z1 = z0 - k1 + G3;
                let x2 = x0 - i2 + 2 * G3;
                let y2 = y0 - j2 + 2 * G3;
                let z2 = z0 - k2 + 2 * G3;
                let x3 = x0 - 1 + 3 * G3;
                let y3 = y0 - 1 + 3 * G3;
                let z3 = z0 - 1 + 3 * G3;

                let ii = i & 255;
                let jj = j & 255;
                let kk = k & 255;

                let t0 = 0.6 - x0 * x0 - y0 * y0 - z0 * z0;
                if (t0 < 0) n0 = 0.0;
                else {
                    t0 *= t0;
                    n0 = t0 * t0 * grad3d(perm[ii + perm[jj + perm[kk]]], x0, y0, z0);
                }

                let t1 = 0.6 - x1 * x1 - y1 * y1 - z1 * z1;
                if (t1 < 0) n1 = 0.0;
                else {
                    t1 *= t1;
                    n1 = t1 * t1 * grad3d(perm[ii + i1 + perm[jj + j1 + perm[kk + k1]]], x1, y1, z1);
                }

                let t2 = 0.6 - x2 * x2 - y2 * y2 - z2 * z2;
                if (t2 < 0) n2 = 0.0;
                else {
                    t2 *= t2;
                    n2 = t2 * t2 * grad3d(perm[ii + i2 + perm[jj + j2 + perm[kk + k2]]], x2, y2, z2);
                }

                let t3 = 0.6 - x3 * x3 - y3 * y3 - z3 * z3;
                if (t3 < 0) n3 = 0.0;
                else {
                    t3 *= t3;
                    n3 = t3 * t3 * grad3d(perm[ii + 1 + perm[jj + 1 + perm[kk + 1]]], x3, y3, z3);
                }

                return 32.0 * (n0 + n1 + n2 + n3);
            };
        })();
    }
    // --- End Simplex Noise Implementation ---
  },
};
</script>

<style scoped>
.three-scene-container {
  width: 100%;
  height: 300px; /* Or adjust as needed for your layout */
  background-color: transparent; /* Allows the dashboard background to show through */
  border-radius: 12px;
  overflow: hidden; /* Ensures the blob stays within the rounded corners */
  position: relative;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05); /* Match card shadow */
}

.three-scene-container canvas {
  display: block; /* Remove extra space below canvas */
  width: 100% !important;
  height: 100% !important;
}
</style>