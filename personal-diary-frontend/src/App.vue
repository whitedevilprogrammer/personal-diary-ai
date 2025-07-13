<template>
  <div id="app">
    <!-- <nav>
      <h1>📝 Personal Diary</h1>
      <div class="nav-links">
        <router-link to="/">Home</router-link>
        <router-link to="/">Logout</router-link>
        <router-link to="/register">Register</router-link>
        <router-link to="/login">Login</router-link>
        <router-link to="/dashboard">Dashboard</router-link>
      </div>
    </nav> -->
     <!-- loader will show automatically during API calls -->
    <Loader />
    <!-- <div>
    <h1>Login with Google</h1>
    <a :href="googleAuthURL">Login</a>
  </div> -->
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
import Loader from './components/Loader.vue'
const googleAuthURL = "http://localhost:8080/auth/google/login"

export default {
  name: 'App',
  components: {
    Loader
  },
  data() {
    return {
      googleAuthURL
    }
  }
}
</script>

<style scoped>

/* Global styles */
body {
  margin: 0;
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#app {
  width: 100%;
  min-height: 100vh;
}

:root {
  --primary: #3a4bc9;
  --primary-dark: #2c3aa3;
  --secondary: #19b787;
  --secondary-dark: #159770;
  --danger: #e74c3c;
  --danger-dark: #c0392b;
  --text-dark: #2d3436;
  --text-light: #636e72;
  --bg-light: #f7fafc;
  --card-bg: #ffffff;
  --border-radius: 10px;
  --box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  --box-shadow-hover: 0 15px 35px rgba(0, 0, 0, 0.1);
  --transition: all 0.3s ease;
}

body {
  font-family: 'Inter', 'Segoe UI', -apple-system, sans-serif;
  background-color: var(--bg-light);
  color: var(--text-dark);
  line-height: 1.6;
}

h1, h2, h3, h4, h5, h6 {
  font-weight: 600;
  line-height: 1.3;
}

.btn {
  display: inline-block;
  padding: 12px 24px;
  background-color: var(--primary);
  color: white;
  text-decoration: none;
  font-weight: 600;
  border-radius: var(--border-radius);
  border: none;
  cursor: pointer;
  transition: var(--transition);
  text-align: center;
}

.btn:hover {
  background-color: var(--primary-dark);
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(58, 75, 201, 0.2);
}

.btn.secondary {
  background-color: var(--secondary);
}

.btn.secondary:hover {
  background-color: var(--secondary-dark);
  box-shadow: 0 5px 15px rgba(25, 183, 135, 0.2);
}

.btn.danger {
  background-color: var(--danger);
}

.btn.danger:hover {
  background-color: var(--danger-dark);
  box-shadow: 0 5px 15px rgba(231, 76, 60, 0.2);
}

.card {
  background: var(--card-bg);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 25px;
  transition: var(--transition);
}

.card:hover {
  box-shadow: var(--box-shadow-hover);
  transform: translateY(-3px);
}

input, textarea {
  width: 100%;
  padding: 12px 16px;
  margin-bottom: 16px;
  border-radius: var(--border-radius);
  border: 1px solid #e0e0e0;
  font-family: inherit;
  font-size: 16px;
  transition: var(--transition);
}

input:focus, textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(58, 75, 201, 0.1);
} 
</style>
