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
  --panel-grey: #1A1D26;
  --divider-grey: #2A2D36;
  --text-primary: #F5F7FA;
  --text-secondary: #9CA3AF;
  --primary-accent: #00C2C2;
  --primary-hover: #00A7A7;
  --danger: #EF4444;
  --success: #22C55E;
}

/* **Cette règle s’appliquera enfin** */
body {
  margin: 0;
  background-color: var(--bg-dark);
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
}

/* Le nav reste au-dessus */
nav {
  position: relative;
  z-index: 20;
  padding: 16px;
  background-color: var(--panel-grey);
  border-bottom: 1px solid var(--divider-grey);
}

nav a {
  color: var(--primary-accent);
}

nav a:hover {
  color: var(--primary-hover);
}

/* Et par défaut, tous tes titres, paragraphes, liens... */
h1,h2,h3,h4,h5,h6,p,a,span,label {
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

#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: -1;
  background-color: #0E111A;
}

