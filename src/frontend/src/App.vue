<template>
  <div id="app">
    <nav>
      <router-link to="/login">Login</router-link> |
      <router-link to="/register">Register</router-link>
    </nav>
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
:root {
  --bg-dark: #0E111A;
  --text-primary: #F5F7FA;
  --text-secondary: #9CA3AF;
  --primary-accent: #00C2C2;
  --primary-hover: #00A7A7;
  --panel-grey: #1A1D26;
  --divider-grey: #2A2D36;
  --danger: #EF4444;
  --success: #22C55E;
}

body {
  margin: 0;
  background-color: var(--bg-dark);
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
}

h1, h2, h3, h4, h5, h6, p, a, span, label {
  color: var(--text-primary);
}
</style>

nav {
  padding: 16px;
  background-color: var(--panel-grey);
  border-bottom: 1px solid var(--divider-grey);
}

nav a {
  margin: 0 8px;
  color: var(--primary-accent);
  text-decoration: none;
}

nav a:hover {
  color: var(--primary-hover);
}
