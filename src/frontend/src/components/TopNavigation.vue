<template>
  <nav class="navbar">
    <!-- Left side: Brand and navigation links -->
    <div class="nav-left">
      <router-link to="/" class="brand">NetSecure-IQ</router-link>
      <router-link 
        to="/" 
        class="nav-link"
        :class="{ 'active': $route.path === '/' }"
      >
        Home
      </router-link>
      <router-link 
        to="/login" 
        class="nav-link"
        :class="{ 'active': $route.path === '/login' }"
      >
        Login
      </router-link>
      <router-link 
        to="/register" 
        class="nav-link"
        :class="{ 'active': $route.path === '/register' }"
      >
        Register
      </router-link>
    </div>
    
    <!-- Right side: Theme toggle -->
    <div class="nav-right">
      <button 
        class="theme-toggle" 
        @click="toggleTheme" 
        :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
      >
        <span v-if="isDark">☀️</span>
        <span v-else>🌙</span>
      </button>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'

const theme = ref(localStorage.getItem('theme') || 'dark')
const isDark = computed(() => theme.value === 'dark')

// Fonction pour basculer le thème
const toggleTheme = () => {
  theme.value = isDark.value ? 'light' : 'dark'
  localStorage.setItem('theme', theme.value)
  document.documentElement.setAttribute('data-theme', theme.value)
}

// Mettre à jour l'attribut data-theme au chargement et lors des changements
onMounted(() => {
  document.documentElement.setAttribute('data-theme', theme.value)
  
  // Observer les changements de thème
  const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.attributeName === 'data-theme') {
        theme.value = document.documentElement.getAttribute('data-theme') || 'dark'
      }
    })
  })
  
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })
  
  return () => observer.disconnect()
})


/**
 * Initialize particles background for navigation
 */
function initNavParticles() {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const container = document.getElementById('nav-particles-js')
  if (!container) return

  // Set canvas size to match container
  const updateCanvasSize = () => {
    canvas.width = container.clientWidth
    canvas.height = container.clientHeight
  }
  
  updateCanvasSize()
  
  // Style canvas
  canvas.style.position = 'absolute'
  canvas.style.top = '0'
  canvas.style.left = '0'
  canvas.style.width = '100%'
  canvas.style.height = '100%'
  canvas.style.pointerEvents = 'none'
  canvas.style.zIndex = '0'

  container.appendChild(canvas)

  // Create particles
  const particles = []
  const particleCount = 30 // Fewer particles for nav bar
  
  // Theme-aware colors
  const getThemeColors = () => {
    const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
    
    if (isDark) {
      return {
        particle: 'rgba(255, 255, 255, 0.5)',
        line: 'rgba(255, 255, 255, 0.2)'
      };
    } else {
      return {
        particle: 'rgba(85, 85, 85, 0.5)',
        line: 'rgba(136, 136, 136, 0.2)'
      };
    }
  }

  // Create particles
  for (let i = 0; i < particleCount; i++) {
    particles.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      vx: (Math.random() - 0.5) * 0.3,
      vy: (Math.random() - 0.5) * 0.3,
      size: Math.random() * 1.5 + 0.5
    })
  }

  let animationId
  // Watch for theme changes
  const observer = new MutationObserver(() => {
    // Re-render particles when theme changes
    const newColors = getThemeColors();
    Object.assign(colors, newColors);
  });
  
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })
  
  // Initial colors
  const colors = getThemeColors();

  const animate = () => {
    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // Update and draw particles
    particles.forEach((particle, i) => {
      // Update position
      particle.x += particle.vx
      particle.y += particle.vy

      // Bounce off edges
      if (particle.x < 0 || particle.x > canvas.width) particle.vx *= -1
      if (particle.y < 0 || particle.y > canvas.height) particle.vy *= -1

      // Draw particle
      ctx.beginPath()
      ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2)
      ctx.fillStyle = colors.particle
      ctx.fill()

      // Draw connections
      particles.slice(i + 1).forEach(otherParticle => {
        const dx = particle.x - otherParticle.x
        const dy = particle.y - otherParticle.y
        const distance = Math.sqrt(dx * dx + dy * dy)
        const maxDistance = 100

        if (distance < maxDistance) {
          ctx.beginPath()
          ctx.moveTo(particle.x, particle.y)
          ctx.lineTo(otherParticle.x, otherParticle.y)
          const opacity = (0.2 * (1 - distance / maxDistance)).toFixed(2);
          ctx.strokeStyle = colors.line.replace(/0\.2/, opacity);
          ctx.lineWidth = 0.5
          ctx.stroke()
        }
      })
    })

    animationId = requestAnimationFrame(animate)
  }

  // Start animation
  animate()

  // Handle resize
  const handleResize = () => {
    updateCanvasSize()
  }

  window.addEventListener('resize', handleResize)

  // Cleanup function
  return () => {
    window.removeEventListener('resize', handleResize)
    cancelAnimationFrame(animationId)
    observer.disconnect()
    if (container && container.contains(canvas)) {
      container.removeChild(canvas)
    }
  }
}

// Initialize particles and set up event listeners
onMounted(() => {
  // Initialize particles
  const cleanupParticles = initNavParticles()
  
  // Handle click outside (placeholder for future use)
  const handleClickOutside = () => {
    // No-op for now
  }
  
  document.addEventListener('click', handleClickOutside)
  
  // Cleanup function
  onUnmounted(() => {
    if (cleanupParticles) cleanupParticles()
    document.removeEventListener('click', handleClickOutside)
  })
})
</script>

<style scoped>
:root {
  --bg-dark: #0e111a;
  --panel-grey: #1a1d26;
  --divider-grey: #2a2d36;
  --text-primary: #f5f7fa;
  --text-secondary: #9ca3af;
  --primary-accent: #00b4aa;
  --primary-hover: #008a8a;
}

/* Navbar styles */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 32px;
  height: 64px;
  background-color: var(--bg-dark);
  border-bottom: 1px solid var(--divider-grey);
  transition: all 0.3s ease;
  width: 100%;
}

[data-theme='light'] .navbar {
  background-color: #ffffff;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

/* Left side - Brand and navigation links */
.nav-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.brand {
  font-weight: bold;
  font-size: 18px;
  color: var(--primary-accent);
  text-decoration: none;
  transition: color 0.2s ease;
}

.brand:hover {
  color: var(--primary-hover);
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 14px;
  padding: 8px 12px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.nav-link:hover,
.nav-link.active {
  color: var(--primary-accent);
  background-color: rgba(0, 194, 194, 0.1);
}

.nav-link.active {
  font-weight: 600;
}

/* Right side - Theme toggle */
.nav-right {
  margin-left: auto;
}

.theme-toggle {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 1.25rem;
  padding: 8px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.theme-toggle:hover {
  background-color: var(--divider-grey);
  color: var(--text-primary);
}

[data-theme='light'] .theme-toggle:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.nav-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.nav-container {
  display: flex;
  align-items: center;
  height: 4.5rem;
  position: relative;
  width: 100%;
}

.nav-brand {
  font-size: 1.25rem;
  font-weight: 600;
  margin-right: 2rem;
  white-space: nowrap;
}

.nav-links-container {
  display: flex;
  flex: 1;
  justify-content: center;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
  align-items: center;
  margin: 0;
  padding: 0;
}

.theme-toggle-container {
  margin-left: auto;
  display: flex;
  align-items: center;
}

/* Ensure the theme toggle button is properly aligned */
.theme-toggle {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s ease;
}

.theme-toggle:hover {
  background-color: var(--divider-grey);
}

[data-theme='light'] .theme-toggle:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.nav-brand {
  flex-shrink: 0;
  z-index: 20;
}

.brand-link {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--primary-accent);
  text-decoration: none;
  transition: color 0.2s ease;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

.brand-link:hover {
  color: var(--primary-hover);
}

/* Navigation Links */
.nav-links {
  display: none;
}

@media (min-width: 768px) {
  .nav-links {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    margin-left: 2rem;
  }
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  font-weight: 500;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  transition: all 0.2s ease;
  position: relative;
  white-space: nowrap;
}

.nav-link:hover,
.nav-link.active {
  color: var(--text-primary);
  background-color: rgba(0, 194, 194, 0.1);
}

.nav-link.active {
  color: var(--primary-accent);
  font-weight: 600;
}

/* Remove divider styles since we're using gap instead */
.nav-divider {
  display: none;
}

/* Responsive styles */
@media (max-width: 768px) {
  .nav-link {
    font-size: 0.9rem;
    padding: 0.5rem 0.75rem;
  }
  
  .brand {
    font-size: 1rem;
  }
  
  .nav-right {
    margin-left: auto;
  }
  
  .nav-container {
    padding: 0 16px;
  }
}

.brand-link {
  font-size: 1.125rem;
}

@media (min-width: 769px) {
  .mobile-nav {
    display: none !important;
  }
}

/* Dark theme support */
[data-theme='dark'] .top-navigation {
  background-color: rgba(26, 29, 38, 0.95);
  border-bottom-color: rgba(42, 45, 54, 0.5);
}

[data-theme='dark'] .mobile-nav {
  background-color: rgba(26, 29, 38, 0.98);
  border-bottom-color: rgba(42, 45, 54, 0.5);
}
</style>