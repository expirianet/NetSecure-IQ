<template>
  <!-- Un seul conteneur pour cette instance -->
  <div id="particles-js" class="particles-layer"></div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'

/** Charge un script une seule fois (par URL) */
function loadScriptOnce(src) {
  return new Promise((resolve, reject) => {
    if ([...document.scripts].some(s => (s.src || '').includes(src))) {
      resolve(); return
    }
    const s = document.createElement('script')
    s.src = src
    s.async = true
    s.onload = resolve
    s.onerror = reject
    document.head.appendChild(s)
  })
}

function destroyInstance(el) {
  const list = (window.pJSDom || [])
  const entry = list.find(d => d && d.pJS && d.pJS.canvas.el && d.pJS.canvas.el.parentElement === el)
  if (entry?.pJS?.fn?.vendors?.destroypJS) {
    entry.pJS.fn.vendors.destroypJS()
  }
}

function initParticles() {
  const el = document.getElementById('particles-js')
  if (!el || !window.particlesJS) return

  // détruire l’instance précédente sur ce conteneur (si re-montage/re-theme)
  destroyInstance(el)

  const dark = document.documentElement.getAttribute('data-theme') === 'dark'

  window.particlesJS('particles-js', {
    particles: {
      number: { value: 80, density: { enable: true, value_area: 800 } },
      color:  { value: dark ? '#ffffff' : '#888888' },
      shape:  { type: 'circle' },
      opacity:{ value: dark ? 0.5 : 0.35 },
      size:   { value: 3, random: true },
      line_linked: {
        enable: true, distance: 150,
        color: dark ? '#ffffff' : '#E0E0E0', opacity: dark ? 0.4 : 0.3, width: 1
      },
      move: { enable: true, speed: 2, direction: 'none', out_mode: 'out' }
    },
    interactivity: {
      detect_on: 'canvas',
      events: { onhover: { enable: true, mode: 'grab' }, resize: true },
      modes: { grab: { distance: 140, line_linked: { opacity: 1 } } }
    },
    retina_detect: true
  })
}

let themeObserver

onMounted(async () => {
  // Charge local si dispo, sinon fallback CDN
  try {
    await loadScriptOnce('/particles/particles.min.js')
  } catch {
    await loadScriptOnce('https://cdn.jsdelivr.net/npm/particles.js@2.0.0/particles.min.js')
  }

  initParticles()

  // Re-render au changement de thème
  themeObserver = new MutationObserver(muts => {
    if (muts.some(m => m.attributeName === 'data-theme')) initParticles()
  })
  themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
})

onUnmounted(() => {
  if (themeObserver?.disconnect) themeObserver.disconnect()
  const el = document.getElementById('particles-js')
  if (el) destroyInstance(el)
})
</script>

<style scoped>
.particles-layer {
  position: fixed;
  inset: 0;
  width: 100vw;
  height: 100vh;
  z-index: 0;            /* visible derrière le contenu, plus de -1 */
  pointer-events: none;  /* clics passent au travers */
  background: transparent;
  opacity: .9;
}
</style>
