<template>
  <div id="bg-particles" class="particles-layer"></div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import {
  ensurePJSDom, loadParticlesScript, defaultConfig,
  safeRender, observeTheme, destroyForId, themeIsDark
} from '@/utils/particles.js'

const ID = 'bg-particles'
let stopObs = () => {}

function render() {
  const ok = safeRender(ID, defaultConfig(themeIsDark()))
  return ok
}

onMounted(async () => {
  try { await loadParticlesScript() } catch {}
  ensurePJSDom()
  render()
  stopObs = observeTheme(ID, render)
})

onUnmounted(() => {
  stopObs?.()
  ensurePJSDom()
  destroyForId(ID)
})
</script>

<style scoped>
.particles-layer {
  position: fixed;
  inset: 0;
  width: 100vw;
  height: 100vh;
  z-index: 0;
  pointer-events: none;
  background: transparent;
  opacity: .9;
}
</style>
