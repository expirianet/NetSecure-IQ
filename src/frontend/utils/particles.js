// Utilities to safely use particles.js (Vincent Garreau) with Vue & theme switching

let _scriptPromise = null

/**
 * Load particles.js from CDN once.
 * Resolves immediately if already present.
 */
export function loadParticlesScript () {
  if (typeof window !== 'undefined' && window.particlesJS) {
    return Promise.resolve()
  }
  if (_scriptPromise) return _scriptPromise

  _scriptPromise = new Promise((resolve, reject) => {
    try {
      const s = document.createElement('script')
      s.src = 'https://cdn.jsdelivr.net/npm/particles.js@2.0.0/particles.min.js'
      s.async = true
      s.onload = () => resolve()
      s.onerror = (e) => reject(new Error('Failed to load particles.js'))
      document.head.appendChild(s)
    } catch (e) {
      reject(e)
    }
  }).catch((e) => {
    console.debug('[particles] script load failed', e)
    throw e
  })

  return _scriptPromise
}

/**
 * Ensure the global pJSDom array exists (created by the lib).
 */
export function ensurePJSDom () {
  try {
    if (typeof window !== 'undefined' && !window.pJSDom) {
      // create an empty array so we can manage instances before the lib populates it
      window.pJSDom = []
    }
  } catch (e) {
    console.debug('[particles] ensurePJSDom error', e)
  }
}

/**
 * Destroy particles instance bound to a given container id.
 */
export function destroyForId (id) {
  try {
    const list = (typeof window !== 'undefined' && window.pJSDom) ? window.pJSDom : []
    for (let i = list.length - 1; i >= 0; i--) {
      try {
        const inst = list[i]
        const el = inst?.pJS?.canvas?.el
        const hostId =
          (el && el.parentElement && el.parentElement.id) ||
          inst?.tag?.id ||
          null
        if (hostId === id) {
          const vendors = inst?.pJS?.fn?.vendors
          // particles.js exposes destroypJS(); some forks expose destroy()
          if (vendors?.destroypJS) vendors.destroypJS()
          if (vendors?.destroy) vendors.destroy()
          list.splice(i, 1)
        }
      } catch (e) {
        console.debug('[particles] single destroy failed', e)
      }
    }
  } catch (e) {
    console.debug('[particles] destroyForId failed', e)
  }
}

/**
 * Render particles into an element by id, defensively.
 * Returns true on success, false if skipped (no elem or no lib).
 */
export function safeRender (id, config) {
  try {
    const el = typeof document !== 'undefined' ? document.getElementById(id) : null
    if (!el) return false
    if (typeof window === 'undefined' || !window.particlesJS) return false

    // kill previous instance (if any) then render
    destroyForId(id)
    window.particlesJS(id, config)
    return true
  } catch (e) {
    console.debug('[particles] render failed', e)
    return false
  }
}

/**
 * Observe theme changes (data-theme attr on <html>) and re-render particles.
 * Returns a function to stop observing.
 */
export function observeTheme (id, renderFn) {
  const target = typeof document !== 'undefined' ? document.documentElement : null
  if (!target) return () => {}

  const observer = new MutationObserver(() => {
    try {
      destroyForId(id)
      renderFn?.()
    } catch (e) {
      console.debug('[particles] theme change error', e)
    }
  })

  try {
    observer.observe(target, { attributes: true, attributeFilter: ['data-theme'] })
  } catch (e) {
    console.debug('[particles] observer start failed', e)
  }

  return () => {
    try { observer.disconnect() } catch (e) {
      console.debug('[particles] observer cleanup failed', e)
    }
  }
}

/**
 * Detect if current theme is dark based on data-theme attribute.
 */
export function themeIsDark () {
  try {
    const t = typeof document !== 'undefined'
      ? document.documentElement.getAttribute('data-theme')
      : null
    return t !== 'light'
  } catch (e) {
    console.debug('[particles] themeIsDark error', e)
    return true
  }
}

/**
 * Minimal default config that adapts to theme.
 */
export function defaultConfig (isDark = true) {
  const baseColor = isDark ? '#00c2c2' : '#0ea5a5'
  const linkColor = isDark ? '#5eead4' : '#14b8a6'
  const bgColor = isDark ? '#0e111a' : '#f6f8fb'

  return {
    particles: {
      number: { value: 60, density: { enable: true, value_area: 800 } },
      color: { value: baseColor },
      shape: { type: 'circle' },
      opacity: { value: 0.4, random: false },
      size: { value: 3, random: true },
      line_linked: { enable: true, distance: 130, color: linkColor, opacity: 0.25, width: 1 },
      move: { enable: true, speed: 1.2, direction: 'none', out_mode: 'out' }
    },
    interactivity: {
      detect_on: 'canvas',
      events: { onhover: { enable: false }, onclick: { enable: false }, resize: true }
    },
    retina_detect: true,
    // for some forks that support background
    background: { color: bgColor }
  }
}
