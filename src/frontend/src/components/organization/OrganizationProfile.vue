===== BEGIN: src/frontend/src/components/organization/OrganizationProfile.vue =====
<template>
  <div class="org-page">
    <!-- Fond animé isolé pour la page profil -->
    <div id="org-particles" class="particles-layer"></div>

    <div class="wrapper">
      <div class="container">
        <div class="card">
          <h2 class="title">NetSecure-IQ</h2>
          <h3 class="subtitle">Organization Profile</h3>

          <div v-if="!org" class="empty">
            <p>No organization information saved yet.</p>
            <RouterLink class="btn primary" to="/organization/edit">Fill the form</RouterLink>
          </div>

          <div v-else class="content">
            <section class="section">
              <h4 class="section-title"><i class="fas fa-building"></i> Organization Information</h4>
              <div class="grid">
                <div class="row"><label>Name</label><span>{{ safe(org.name) }}</span></div>
                <div class="row"><label>VAT / Fiscal Code</label><span>{{ safe(org.vat_number) }}</span></div>
                <div class="row wide"><label>Address</label><span>{{ safe(org.address) }}</span></div>
                <div class="row"><label>City</label><span>{{ safe(org.city) }}</span></div>
                <div class="row"><label>State</label><span>{{ safe(org.state) }}</span></div>
                <div class="row"><label>ZIP</label><span>{{ safe(org.zip_code) }}</span></div>
                <div class="row"><label>Email</label><span>{{ safe(org.contact_email) }}</span></div>
                <div class="row"><label>Phone</label><span>{{ safe(org.contact_phone) }}</span></div>
                <div class="row"><label>PEC Email</label><span>{{ safe(org.pec_email) }}</span></div>
                <div class="row"><label>SDI Code</label><span>{{ safe(org.sdi_code) }}</span></div>
              </div>
            </section>

            <section class="section">
              <h4 class="section-title"><i class="fas fa-user-tie"></i> Company Manager</h4>
              <div class="grid">
                <div class="row"><label>Name</label><span>{{ person(manager).name }}</span></div>
                <div class="row"><label>Email</label><span>{{ person(manager).email }}</span></div>
                <div class="row"><label>Phone</label><span>{{ person(manager).phone }}</span></div>
              </div>
            </section>

            <section v-if="hasTechnical" class="section">
              <h4 class="section-title"><i class="fas fa-user-cog"></i> Technical Manager</h4>
              <div class="grid">
                <div class="row"><label>Name</label><span>{{ person(technical).name }}</span></div>
                <div class="row"><label>Email</label><span>{{ person(technical).email }}</span></div>
                <div class="row"><label>Phone</label><span>{{ person(technical).phone }}</span></div>
              </div>
            </section>

            <section class="section">
              <h4 class="section-title"><i class="fas fa-shield-alt"></i> Data Controller</h4>
              <div class="grid">
                <div class="row"><label>Name</label><span>{{ person(controller).name }}</span></div>
                <div class="row"><label>Email</label><span>{{ person(controller).email }}</span></div>
                <div class="row"><label>Phone</label><span>{{ person(controller).phone }}</span></div>
              </div>
            </section>

            <section class="section">
              <h4 class="section-title"><i class="fas fa-database"></i> Data Processor</h4>
              <div class="grid">
                <div class="row"><label>Name</label><span>{{ person(processor).name }}</span></div>
                <div class="row"><label>Email</label><span>{{ person(processor).email }}</span></div>
                <div class="row"><label>Phone</label><span>{{ person(processor).phone }}</span></div>
              </div>
            </section>

            <section v-if="org.personnel_info" class="section">
              <h4 class="section-title"><i class="fas fa-file-alt"></i> Personnel Info</h4>
              <pre class="pre">{{ org.personnel_info }}</pre>
            </section>

            <div class="actions">
              <RouterLink class="btn primary" to="/organization/edit">Edit</RouterLink>
              <button class="btn ghost" @click="clearLocal">Reset local data</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * Important :
 * - Pas de <BackgroundParticles/> ici.
 * - On utilise un ID dédié + garde-fou pJSDom.
 */
import { computed, ref, onMounted, onBeforeUnmount } from 'vue'

const CONTAINER_ID = 'org-particles'
let themeObserver

function ensurePJSDom() {
  // 👉 corrige le crash "reading 'push' of null" dans particles.js
  if (!Array.isArray(window.pJSDom)) window.pJSDom = []
}

function loadScriptOnce(src) {
  return new Promise((resolve, reject) => {
    if ([...document.scripts].some(s => (s.src || '').includes(src))) return resolve()
    const s = document.createElement('script')
    s.src = src; s.async = true
    s.onload = resolve; s.onerror = reject
    document.head.appendChild(s)
  })
}

function destroyFor(el) {
  try {
    if (!el) return
    if (Array.isArray(window.pJSDom)) {
      window.pJSDom = window.pJSDom.filter(entry => {
        const same = entry?.pJS?.canvas?.el?.parentElement === el
        if (same) { try { entry.pJS.fn.vendors.destroypJS() } catch {} }
        return !same
      })
    }
    el.querySelectorAll('canvas').forEach(c => c.remove())
  } catch {}
}

function initParticles() {
  const el = document.getElementById(CONTAINER_ID)
  if (!el || !window.particlesJS) return
  ensurePJSDom()
  destroyFor(el)

  const dark = document.documentElement.getAttribute('data-theme') === 'dark'
  window.particlesJS(CONTAINER_ID, {
    particles: {
      number: { value: 80, density: { enable: true, value_area: 800 } },
      color:  { value: dark ? '#ffffff' : '#888888' },
      shape:  { type: 'circle' },
      opacity:{ value: dark ? 0.5 : 0.35 },
      size:   { value: 3, random: true },
      line_linked: { enable: true, distance: 150, color: dark ? '#ffffff' : '#E0E0E0', opacity: 0.35, width: 1 },
      move: { enable: true, speed: 2, direction: 'none', out_mode: 'out' }
    },
    interactivity: { detect_on: 'canvas', events: { onhover: { enable: true, mode: 'grab' }, resize: true },
      modes: { grab: { distance: 140, line_linked: { opacity: 1 } } } },
    retina_detect: true
  })
}

/* ----- Données locales (inchangé) ----- */
function readLocal() { try { return JSON.parse(localStorage.getItem('organization_profile') || 'null') } catch { return null } }
const org = ref(readLocal())
const manager    = computed(() => org.value?.manager    || {})
const technical  = computed(() => org.value?.technical  || null)
const controller = computed(() => org.value?.controller || {})
const processor  = computed(() => org.value?.processor  || {})
const hasTechnical = computed(() => !!(technical.value && (technical.value.name || technical.value.email || technical.value.phone)))
function safe(v) { return (v ?? '').toString().trim() || '—' }
function person(p) { return { name: safe(p?.name), email: safe(p?.email), phone: safe(p?.phone) } }
function clearLocal() { localStorage.removeItem('organization_profile'); org.value = null }

onMounted(async () => {
  document.title = 'NetSecure-IQ - Organization Profile'
  org.value = readLocal()

  try { await loadScriptOnce('/particles/particles.min.js') }
  catch { await loadScriptOnce('https://cdn.jsdelivr.net/npm/particles.js@2.0.0/particles.min.js') }

  ensurePJSDom()
  initParticles()

  themeObserver = new MutationObserver(m => {
    if (m.some(x => x.attributeName === 'data-theme')) initParticles()
  })
  themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
})

onBeforeUnmount(() => {
  themeObserver?.disconnect?.()
  destroyFor(document.getElementById(CONTAINER_ID))
})
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root{ --bg-dark:#0e111a; --panel-grey:#1a1d26; --divider-grey:#2a2d36;
--text-primary:#f5f7fa; --text-secondary:#9ca3af; --primary-accent:#00c2c2; --primary-hover:#00a7a7;
--radius:16px; --danger:#ef4444; --success:#22c55e; }

.org-page{ position:relative; min-height:100vh; overflow:hidden; }
.particles-layer{ position:fixed; inset:0; width:100vw; height:100vh; z-index:0; pointer-events:none; }

.wrapper{ position:relative; z-index:10; display:flex; align-items:flex-start; justify-content:center; padding:32px; min-height:100vh; }
.container{ width:100%; max-width:1000px; }
.card{ background:var(--panel-grey); border-radius:var(--radius); padding:26px; border:1px solid rgba(255,255,255,.06); box-shadow:0 0 40px rgba(0,194,194,.05); }

.title{ text-align:center; font-size:20px; font-weight:600; color:var(--primary-accent); margin-bottom:6px; }
.subtitle{ text-align:center; font-size:16px; margin-bottom:20px; }

.section{ margin-top:16px; padding:14px; border:1px solid rgba(255,255,255,.05); border-radius:12px; background:rgba(31,41,55,.30); }
.section-title{ margin:0 0 10px; color:var(--primary-accent); font-size:15px; font-weight:600; display:flex; align-items:center; gap:8px; }

.grid{ display:grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap:10px 14px; }
.row{ display:flex; flex-direction:column; gap:6px; }
.row.wide{ grid-column: 1 / -1; }
label{ font-size:12px; color:var(--text-secondary); }
span{ font-size:14px; color:var(--text-primary); }

.pre{ white-space:pre-wrap; background:#0b0e16; border:1px solid var(--divider-grey); border-radius:8px; padding:10px 12px; color:#e5e7eb; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }

.actions{ display:flex; gap:10px; justify-content:flex-end; margin-top:14px; flex-wrap:wrap; }
.btn{ display:inline-flex; align-items:center; justify-content:center; gap:8px; border:none; border-radius:8px; font-weight:700; cursor:pointer; padding:10px 14px; transition:.2s; }
.btn.primary{ background:var(--primary-accent); color:#0e111a; }
.btn.primary:hover{ background:var(--primary-hover); color:#fff; }
.btn.ghost{ background:rgba(255,255,255,.06); border:1px solid rgba(255,255,255,.08); color:var(--text-primary); }
.btn.ghost:hover{ background:rgba(255,255,255,.10); }

.empty{ text-align:center; color:var(--text-secondary); padding:10px 0 4px; display:grid; gap:12px; }
@media (max-width: 720px){ .wrapper{ padding:16px; } .card{ padding:20px; } }
</style>
===== END: src/frontend/src/components/organization/OrganizationProfile.vue =====
