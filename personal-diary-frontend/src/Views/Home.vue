<template>
  <div class="home-container">
    <div class="home-card card">
      <div class="logo">📝</div>
      <h1>Your Personal Diary</h1>
      <p class="tagline">
        Write your thoughts, track your memories, and keep everything safe in one elegant space.
      </p>
      <div class="buttons">
        <router-link to="/register" class="btn">Get Started</router-link>
        <router-link to="/login" class="btn secondary">Login</router-link>
      </div>
      <div class="decoration">
        <div class="circle c1"></div>
        <div class="circle c2"></div>
        <div class="circle c3"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 90vh;
  padding: 20px;
}

.home-card {
  max-width: 700px;
  width: 100%;
  text-align: center;
  padding: 60px 40px;
  position: relative;
  overflow: hidden;
}

.logo {
  font-size: 48px;
  margin-bottom: 20px;
  display: inline-block;
  background: linear-gradient(135deg, var(--primary), var(--secondary));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  filter: drop-shadow(0 2px 5px rgba(0,0,0,0.1));
}

h1 {
  font-size: 36px;
  margin-bottom: 16px;
  background: linear-gradient(135deg, var(--text-dark) 30%, var(--primary) 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.tagline {
  font-size: 18px;
  color: var(--text-light);
  margin-bottom: 40px;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.buttons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
  z-index: 10;
  position: relative;
}

.btn {
  min-width: 140px;
}

/* Decorative elements */
.decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
}

.circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
}

.c1 {
  background: linear-gradient(to right, var(--primary), var(--secondary));
  width: 300px;
  height: 300px;
  top: -100px;
  right: -100px;
}

.c2 {
  background: var(--primary);
  width: 200px;
  height: 200px;
  bottom: -70px;
  left: -70px;
}

.c3 {
  background: var(--secondary);
  width: 150px;
  height: 150px;
  bottom: 50px;
  right: 50px;
}
</style>