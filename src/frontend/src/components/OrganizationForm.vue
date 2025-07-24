<template>
  <div class="login-page">
    <!-- Canvas animé -->
    <div id="particles-js"></div>

    <!-- Formulaire d'organisation -->
    <div class="login-wrapper">
      <div class="login-container">
        <div class="login-card">
          <h2 class="login-title">NetSecure-IQ</h2>
          <h3 class="login-subtitle">Register Your Organization</h3>

          <form @submit.prevent="submitForm" class="login-form">
            <!-- Organization Info -->
            <div class="form-section">
              <h4><i class="fas fa-building"></i> Organization Information</h4>
              <input v-model="form.name" placeholder="Organization Name" required />
              <input v-model="form.vat_number" placeholder="VAT Number or Fiscal Code" required />
              <input v-model="form.address" placeholder="Address" required />
              <div class="form-row">
                <input v-model="form.city" placeholder="City" required />
                <input v-model="form.state" placeholder="State" required />
                <input v-model="form.zip_code" placeholder="ZIP Code" required />
              </div>
              <input v-model="form.email" type="email" placeholder="Email" required />
              <input v-model="form.phone" type="tel" placeholder="Phone Number" required />
              <input v-model="form.pec_email" type="email" placeholder="PEC Email (Optional)" />
              <input v-model="form.sdi" placeholder="SDI Code (Optional)" />
            </div>

            <!-- Company Manager -->
            <div class="form-section">
              <h4><i class="fas fa-user-tie"></i> Company Manager</h4>
              <input v-model="form.manager_name" placeholder="Name and Surname" required />
              <input v-model="form.manager_email" type="email" placeholder="Email" required />
              <input v-model="form.manager_phone" placeholder="Phone Number" required />
            </div>

            <!-- Technical Manager -->
            <div class="form-section">
              <h4><i class="fas fa-user-cog"></i> Technical Manager</h4>
              <input v-model="form.technical_name" placeholder="Name and Surname" required />
              <input v-model="form.technical_email" type="email" placeholder="Email" required />
              <input v-model="form.technical_phone" placeholder="Phone Number" required />
            </div>

            <!-- Data Controller -->
            <div class="form-section">
              <h4><i class="fas fa-shield-alt"></i> Data Controller</h4>
              <input v-model="form.controller_name" placeholder="Name and Surname" required />
              <input v-model="form.controller_email" type="email" placeholder="Email" required />
              <input v-model="form.controller_phone" placeholder="Phone Number" required />
            </div>

            <!-- Data Processor -->
            <div class="form-section">
              <h4><i class="fas fa-database"></i> Data Processor</h4>
              <input v-model="form.processor_name" placeholder="Name and Surname" required />
              <input v-model="form.processor_email" type="email" placeholder="Email" required />
              <input v-model="form.processor_phone" placeholder="Phone Number" required />
            </div>

            <div class="form-actions">
              <button type="submit" :disabled="loading">
                {{ loading ? 'Submitting...' : 'Submit' }}
              </button>
              <button type="button" class="btn-secondary" @click="goToDashboard">
                Go to Dashboard
              </button>
            </div>

            <p v-if="message" class="login-message" :class="messageType">
              {{ message }}
            </p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const message = ref('')
const loading = ref(false)
const successMessage = ref(false)

const messageType = computed(() => {
  return successMessage.value ? 'success' : 'error'
})

// Initialize particles.js
const initializeParticles = () => {
  // Only run on client-side
  if (typeof window !== 'undefined' && window.particlesJS) {
    // Destroy existing instance if it exists
    if (window.pJSDom && window.pJSDom.length > 0) {
      window.pJSDom[0].pJS.fn.vendors.destroypJS();
      window.pJSDom = [];
    }

    // Initialize particles with theme-appropriate colors
    const isDark = document.documentElement.classList.contains('dark')
    
    window.particlesJS('particles-js', {
      particles: {
        number: { value: 80, density: { enable: true, value_area: 800 } },
        color: { value: isDark ? '#ffffff' : '#000000' },
        shape: { type: 'circle' },
        opacity: {
          value: 0.5,
          random: true,
          anim: { enable: true, speed: 1, opacity_min: 0.1, sync: false }
        },
        size: {
          value: 3,
          random: true,
          anim: { enable: true, speed: 2, size_min: 0.1, sync: false }
        },
        line_linked: {
          enable: true,
          distance: 150,
          color: isDark ? '#ffffff' : '#000000',
          opacity: 0.2,
          width: 1
        },
        move: {
          enable: true,
          speed: 1,
          direction: 'none',
          random: false,
          straight: false,
          out_mode: 'out',
          bounce: false,
          attract: { enable: false, rotateX: 600, rotateY: 1200 }
        }
      },
      interactivity: {
        detect_on: 'canvas',
        events: {
          onhover: { enable: true, mode: 'grab' },
          onclick: { enable: true, mode: 'push' },
          resize: true
        },
        modes: {
          grab: { distance: 140, line_linked: { opacity: 0.5 } },
          push: { particles_nb: 4 }
        }
      },
      retina_detect: true
    });
  }
}

// Handle theme changes
const handleThemeChange = (mutations) => {
  for (const mutation of mutations) {
    if (mutation.attributeName === 'class') {
      initializeParticles()
    }
  }
}

// Initialize particles when component is mounted
onMounted(() => {
  // Wait for particles.js to be loaded
  if (typeof window !== 'undefined' && !window.particlesJS) {
    const script = document.createElement('script')
    script.src = 'https://cdn.jsdelivr.net/particles.js/2.0.0/particles.min.js'
    script.onload = () => {
      initializeParticles()
      // Observe theme changes
      const observer = new MutationObserver(handleThemeChange)
      observer.observe(document.documentElement, { attributes: true })
      
      // Cleanup observer when component is unmounted
      onUnmounted(() => {
        observer.disconnect()
      })
    }
    document.head.appendChild(script)
  } else {
    initializeParticles()
  }
})

const form = reactive({
  name: '', vat_number: '', address: '', state: '', city: '', zip_code: '',
  email: '', pec_email: '', sdi: '', phone: '',
  manager_name: '', manager_email: '', manager_phone: '',
  technical_name: '', technical_email: '', technical_phone: '',
  controller_name: '', controller_email: '', controller_phone: '',
  processor_name: '', processor_email: '', processor_phone: ''
})

const goToDashboard = () => {
  router.push('/dashboard')
}

const submitForm = async () => {
  console.log("📤 submitForm() triggered")

  const personnelInfo = `
Company Manager:
  Name: ${form.manager_name}
  Email: ${form.manager_email}
  Phone: ${form.manager_phone}

Technical Manager:
  Name: ${form.technical_name}
  Email: ${form.technical_email}
  Phone: ${form.technical_phone}

Data Controller:
  Name: ${form.controller_name}
  Email: ${form.controller_email}
  Phone: ${form.controller_phone}

Data Processor:
  Name: ${form.processor_name}
  Email: ${form.processor_email}
  Phone: ${form.processor_phone}
`.trim()

  const payload = {
    name: form.name,
    vat_number: form.vat_number,
    address: form.address,
    state: form.state,
    city: form.city,
    zip_code: form.zip_code,
    contact_email: form.email,
    pec_email: form.pec_email,
    sdi_code: form.sdi,
    contact_phone: form.phone,
    personnel_info: personnelInfo,
    user_id: localStorage.getItem("user_id"),
  }

  console.log("📤 Sending payload:", JSON.stringify(payload, null, 2))

  try {
    const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/complete-organization`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    })

    const text = await response.text()
    let data
    try {
      data = JSON.parse(text)
    } catch {
      throw new Error(text)
    }

    if (!response.ok) throw new Error(data.error || data.message)

    message.value = "Organization info submitted! Redirecting..."
    setTimeout(() => router.push('/dashboard'), 1000)

  } catch (err) {
    message.value = "Submission failed: " + err.message
  }
}
</script>


<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css');

/* ===== Base Variables ===== */
:root {
  --bg-dark: #0e111a;
  --panel-grey: #1a1d26;
  --divider-grey: #2a2d36;
  --text-primary: #f5f7fa;
  --text-secondary: #9ca3af;
  --primary-accent: #00c2c2;
  --primary-hover: #00a7a7;
  --danger: #ef4444;
  --success: #22c55e;
  --border-radius: 12px;
  --transition: all 0.2s ease;
}

/* ===== Base Styles ===== */
.login-page {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background-color: var(--bg-dark);
}

/* ===== Particles ===== */
#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 0;
  background-color: var(--bg-dark);
  transition: background-color 0.3s ease;
  pointer-events: none;
}

[data-theme='light'] #particles-js {
  background-color: #E0E0E0;
}

/* ===== Login Wrapper ===== */
.login-wrapper {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  min-height: 100vh;
}

/* ===== Container ===== */
.login-container {
  width: 100%;
  max-width: 800px;
}

/* ===== Login Card ===== */
.login-card {
  background-color: var(--panel-grey);
  border-radius: var(--border-radius);
  padding: 32px;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  box-sizing: border-box;
  margin: 20px 0;
}

/* ===== Headers ===== */
.login-title {
  text-align: center;
  font-size: 24px;
  font-weight: 600;
  color: var(--primary-accent);
  margin-bottom: 8px;
}

.login-subtitle {
  text-align: center;
  font-size: 18px;
  color: var(--text-primary);
  margin-bottom: 32px;
  font-weight: 500;
}

/* ===== Form ===== */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ===== Form Sections ===== */
.form-section {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
}

.form-section h4 {
  color: var(--primary-accent);
  margin: 0 0 16px;
  font-size: 16px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-section h4 i {
  font-size: 14px;
}

/* ===== Form Inputs ===== */
.login-form input {
  width: 100%;
  background-color: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 6px;
  padding: 12px 14px;
  font-size: 14px;
  color: var(--text-primary);
  transition: var(--transition);
  margin-bottom: 8px;
}

.login-form input:focus {
  outline: none;
  border-color: var(--primary-accent);
  box-shadow: 0 0 0 2px rgba(0, 194, 194, 0.2);
}

.login-form input::placeholder {
  color: var(--text-secondary);
  opacity: 0.7;
}

/* ===== Form Rows ===== */
.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}

.form-row input {
  margin-bottom: 0;
}

/* ===== Form Actions ===== */
.form-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 24px;
  gap: 16px;
}

/* ===== Buttons ===== */
button {
  flex: 1;
  padding: 12px 20px;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

button:not(:disabled):hover {
  transform: translateY(-1px);
}

/* Primary Button */
button[type='submit'] {
  background-color: var(--primary-accent);
  color: #0e111a;
}

button[type='submit']:not(:disabled):hover {
  background-color: var(--primary-hover);
}

/* Secondary Button */
.btn-secondary {
  background-color: transparent;
  color: var(--primary-accent);
  border: 1px solid var(--primary-accent) !important;
}

.btn-secondary:not(:disabled):hover {
  background-color: rgba(0, 194, 194, 0.1);
}

/* ===== Messages ===== */
.login-message {
  margin-top: 16px;
  padding: 12px 16px;
  border-radius: 6px;
  font-size: 14px;
  text-align: center;
  transition: var(--transition);
}

.login-message.success {
  background-color: rgba(34, 197, 94, 0.1);
  color: var(--success);
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.login-message.error {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

/* ===== Responsive Design ===== */
@media (max-width: 768px) {
  .login-wrapper {
    padding: 16px;
  }
  
  .login-card {
    padding: 24px;
  }
  
  .form-row {
    flex-direction: column;
    gap: 8px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  button {
    width: 100%;
  }
}
</style>
