<template>
  <div id="particles-js"></div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue';
import 'particles.js';

let particlesInstance = null;

const initializeParticles = () => {
  if (window.particlesJS) {
    particlesInstance = window.particlesJS('particles-js', {
      particles: {
        number: {
          value: 80,
          density: {
            enable: true,
            value_area: 800
          }
        },
        color: {
          value: '#3B82F6' // blue-500
        },
        shape: {
          type: 'circle'
        },
        opacity: {
          value: 0.5,
          random: true
        },
        size: {
          value: 3,
          random: true
        },
        line_linked: {
          enable: true,
          distance: 150,
          color: '#93C5FD', // blue-300
          opacity: 0.4,
          width: 1
        },
        move: {
          enable: true,
          speed: 2,
          direction: 'none',
          random: true,
          straight: false,
          out_mode: 'out',
          bounce: false
        }
      },
      interactivity: {
        detect_on: 'canvas',
        events: {
          onhover: {
            enable: true,
            mode: 'grab'
          },
          onclick: {
            enable: false
          },
          resize: true
        },
        modes: {
          grab: {
            distance: 140,
            line_linked: {
              opacity: 1
            }
          }
        }
      },
      retina_detect: true
    });
  }
};

const handleResize = () => {
  if (particlesInstance && particlesInstance.pJSDom && particlesInstance.pJSDom.length > 0) {
    particlesInstance.pJSDom[0].pJS.particles.move.speed = window.innerWidth < 768 ? 1 : 2;
    particlesInstance.pJSDom[0].pJS.fn.particlesRefresh();
  }
};

onMounted(() => {
  // Load particles.js dynamically if not already loaded
  if (!window.particlesJS) {
    const script = document.createElement('script');
    script.src = 'https://cdn.jsdelivr.net/particles.js/2.0.0/particles.min.js';
    script.onload = () => {
      initializeParticles();
      window.addEventListener('resize', handleResize);
    };
    document.head.appendChild(script);
  } else {
    initializeParticles();
    window.addEventListener('resize', handleResize);
  }
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  if (particlesInstance && particlesInstance.pJSDom && particlesInstance.pJSDom[0]) {
    particlesInstance.pJSDom[0].pJS.fn.vendors.destroypJS();
    particlesInstance = null;
  }
});
</script>

<style scoped>
#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
  background-color: var(--bg-page);
  opacity: 0.6;
}

/* Dark mode adjustments */
:root[data-theme="dark"] #particles-js {
  background-color: var(--bg-page);
}
</style>