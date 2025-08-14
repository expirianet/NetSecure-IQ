// Utils robustes pour particles.js (sans crash dans toutes les situations)

export function ensurePJSDom() {
    if (!Array.isArray(window.pJSDom)) window.pJSDom = []
  }
  
  export function loadParticlesScript() {
    return new Promise((resolve, reject) => {
      if (window.particlesJS) return resolve()
      // Déjà injecté ?
      if ([...document.scripts].some(s => (s.src || '').includes('/particles/particles.min.js'))) {
        // Attendre que particlesJS apparaisse
        const iv = setInterval(() => {
          if (window.particlesJS) { clearInterval(iv); resolve() }
        }, 20)
        setTimeout(() => { clearInterval(iv); resolve() }, 1000) // fail-soft
        return
      }
      const s = document.createElement('script')
      s.src = '/particles/particles.min.js'
      s.async = true
      s.onload = () => resolve()
      s.onerror = async () => {
        // Fallback CDN
        const cdn = document.createElement('script')
        cdn.src = 'https://cdn.jsdelivr.net/npm/particles.js@2.0.0/particles.min.js'
        cdn.async = true
        cdn.onload = () => resolve()
        cdn.onerror = () => reject(new Error('particles.js load failed'))
        document.head.appendChild(cdn)
      }
      document.head.appendChild(s)
    })
  }
  
  export function destroyForId(id) {
    const el = document.getElementById(id)
    if (!el) return
    const listBefore = Array.isArray(window.pJSDom) ? window.pJSDom : []
    const entry = listBefore.find(d => d?.pJS?.canvas?.el?.parentElement === el)
    try { entry?.pJS?.fn?.vendors?.destroypJS?.() } catch {}
    const current = Array.isArray(window.pJSDom) ? window.pJSDom : listBefore
    try {
      window.pJSDom = current.filter(d => d?.pJS?.canvas?.el?.parentElement !== el)
    } catch {
      window.pJSDom = []
    }
    try { el.querySelectorAll('canvas').forEach(c => c.remove()) } catch {}
  }
  
  export function hasContainer(id) {
    return !!document.getElementById(id)
  }
  
  export function themeIsDark() {
    const de = document.documentElement
    return de.getAttribute('data-theme') === 'dark' || de.classList.contains('dark')
  }
  
  export function safeRender(id, config) {
    const el = document.getElementById(id)
    if (!el || !window.particlesJS) return false
    ensurePJSDom()
    destroyForId(id)
    try {
      window.particlesJS(id, config)
      return true
    } catch {
      return false
    }
  }
  
  /**
   * Observe les changements de thème. Retourne une fonction "stop()" pour débrancher.
   * Si le conteneur n'existe pas (page démontée), l'observer continue silencieusement,
   * mais le render n'est appelé que si le conteneur est présent.
   */
  export function observeTheme(id, onThemeChange) {
    const obs = new MutationObserver(muts => {
      if (!hasContainer(id)) return
      if (muts.some(m => m.attributeName === 'data-theme' || m.attributeName === 'class')) {
        onThemeChange()
      }
    })
    obs.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme', 'class']
    })
    return () => obs.disconnect()
  }
  
  /** Config par défaut, sensible au thème */
  export function defaultConfig(dark = themeIsDark()) {
    return {
      particles: {
        number: { value: 80, density: { enable: true, value_area: 800 } },
        color: { value: dark ? '#ffffff' : '#334155' },      // ← plus foncé en light
        shape: { type: 'circle' },
        opacity: { value: dark ? 0.5 : 0.55 },               // ← un peu plus visible en light
        size: { value: 3, random: true },
        line_linked: {
          enable: true,
          distance: 150,
          color: dark ? '#ffffff' : '#475569',               // ← liens plus foncés en light
          opacity: dark ? 0.4 : 0.45,                        // ← un peu plus opaque
          width: 1
        },
        move: { enable: true, speed: 3.6, direction: 'none', out_mode: 'bounce' }
      },
      interactivity: {
        detect_on: 'canvas',
        events: { onhover: { enable: true, mode: 'repulse' }, onclick: { enable: true, mode: 'push' }, resize: true },
        modes: { repulse: { distance: 200 }, push: { particles_nb: 4 } }
      },
      retina_detect: true
    }
  }
  
  