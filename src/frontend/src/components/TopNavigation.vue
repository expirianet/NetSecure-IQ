<template>
  <nav class="top-navigation">
    <!-- Background particles canvas -->
    <div id="nav-particles-js"></div>
    
    <!-- Navigation content -->
    <div class="nav-content">
      <div class="nav-container">
        <!-- Logo/Brand -->
        <div class="nav-brand">
          <router-link to="/" class="brand-link">NetSecure-IQ</router-link>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="nav-links desktop-nav">
          <router-link 
            to="/" 
            class="nav-link"
            :class="{ 'active': $route.path === '/' }"
          >
            Home
          </router-link>
          <span class="nav-divider">|</span>
          <router-link 
            to="/login" 
            class="nav-link"
            :class="{ 'active': $route.path === '/login' }"
          >
            Login
          </router-link>
          <span class="nav-divider">|</span>
          <router-link 
            to="/register" 
            class="nav-link"
            :class="{ 'active': $route.path === '/register' }"
          >
            Register
          </router-link>
        </div>

        <!-- Theme Toggle -->
        <button 
          class="theme-toggle" 
          @click="toggleTheme" 
          :aria-label="isDark ? 'Passer en mode clair' : 'Passer en mode sombre'"
        >
          <span v-if="isDark">☀️</span>
          <span v-else>🌙</span>
        </button>

        <!-- Mobile Menu Toggle -->
        <button 
          class="mobile-menu-toggle md:hidden"
          @click="toggleMobileMenu"
          :class="{ 'active': mobileMenuOpen }"
        >
          <span class="hamburger-line"></span>
          <span class="hamburger-line"></span>
          <span class="hamburger-line"></span>
        </button>
      </div>

      <!-- Mobile Navigation Menu -->
      <div 
        class="mobile-nav"
        :class="{ 'open': mobileMenuOpen }"
      >
        <router-link 
          to="/" 
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/' }"
          @click="closeMobileMenu"
        >
          Home
        </router-link>
        <router-link 
          to="/login" 
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/login' }"
          @click="closeMobileMenu"
        >
          Login
        </router-link>
        <router-link 
          to="/register" 
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/register' }"
          @click="closeMobileMenu"
        >
          Register
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'

const mobileMenuOpen = ref(false)
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

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}

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
  
  // Set up click outside handler
  const handleClickOutside = (event) => {
    if (mobileMenuOpen.value && !event.target.closest('.top-navigation')) {
      closeMobileMenu()
    }
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

.top-navigation {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 80px;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  margin: 0;
  box-sizing: border-box;
}

#nav-particles-js {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}

.nav-content {
  position: relative;
  z-index: 10;
  height: 100%;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: relative;
}

.nav-brand {
  flex-shrink: 0;
  z-index: 20;
}

.brand-link {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--primary-accent);
  text-decoration: none;
  transition: color 0.2s ease;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

.brand-link:hover {
  color: var(--primary-hover);
}

.desktop-nav {
  display: flex;
  align-items: center;
  gap: 16px;
  z-index: 20;
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  font-weight: 500;
  padding: 8px 12px;
  border-radius: 6px;
  transition: all 0.2s ease;
  position: relative;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

.nav-link:hover {
  color: var(--primary-accent);
  background-color: rgba(0, 180, 170, 0.05);
}

.nav-link.active {
  color: var(--primary-accent);
  font-weight: 600;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background-color: var(--primary-accent);
  border-radius: 50%;
}

.nav-divider {
  color: var(--divider-grey);
  font-size: 1.25rem;
  user-select: none;
  opacity: 0.3;
}

/* Theme Toggle Button */
.theme-toggle {
  background: none;
  border: 2px solid var(--primary-accent);
  border-radius: 8px;
  color: var(--primary-accent);
  font-size: 1.1em;
  cursor: pointer;
  padding: 4px 10px;
  margin-left: 15px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
}

.theme-toggle:hover {
  background-color: var(--primary-accent);
  color: var(--bg-dark);
}

/* Mobile Menu Toggle */
.mobile-menu-toggle {
  display: none;
  flex-direction: column;
  justify-content: space-between;
  width: 30px;
  height: 24px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  gap: 6px;
  z-index: 30;
  transition: all 0.3s ease;
  border-radius: 8px;
  margin-left: 15px;
}

.mobile-menu-toggle:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

[data-theme='dark'] .mobile-menu-toggle:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.hamburger-line {
  display: block;
  width: 100%;
  height: 2px;
  background-color: var(--text-primary);
  transition: all 0.3s ease;
}

.mobile-menu-toggle.active .hamburger-line:nth-child(1) {
  transform: translateY(8px) rotate(45deg);
}

.mobile-menu-toggle.active .hamburger-line:nth-child(2) {
  opacity: 0;
}

.mobile-menu-toggle.active .hamburger-line:nth-child(3) {
  transform: translateY(-8px) rotate(-45deg);
}

/* Mobile Navigation */
.mobile-nav {
  position: fixed;
  top: 80px;
  left: 0;
  right: 0;
  background-color: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transform: translateY(-100%);
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 10;
}

[data-theme='dark'] .mobile-nav {
  background-color: rgba(26, 29, 38, 0.98);
  border-bottom-color: rgba(255, 255, 255, 0.05);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.mobile-nav.open {
  transform: translateY(0);
  opacity: 1;
  visibility: visible;
}

.mobile-nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 1rem;
  font-weight: 500;
  padding: 12px 24px;
  margin: 0 16px;
  border-radius: 8px;
  transition: all 0.2s ease;
  text-align: left;
  display: flex;
  align-items: center;
}

.mobile-nav-link:hover {
  color: var(--primary-accent);
  background-color: rgba(0, 180, 170, 0.05);
}

.mobile-nav-link.active {
  color: var(--primary-accent);
  font-weight: 600;
  background-color: rgba(0, 180, 170, 0.1);
}

/* Responsive Design */
@media (max-width: 768px) {
  .desktop-nav {
    display: none;
  }
  
  .mobile-menu-toggle {
    display: flex;
  }
  
  .nav-container {
    padding: 0 16px;
  }
  
  .brand-link {
    font-size: 1.125rem;
  }
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