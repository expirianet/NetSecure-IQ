<template>
  <div class="organization-profile">
    <BackgroundParticles />
    
    <div class="form-container">
      <div class="form-header">
        <h1 class="form-title">Organization Profile</h1>
        <p class="form-subtitle">View your organization's information</p>
      </div>

      <div class="form-content">
        <!-- Organization Information -->
        <div class="form-section">
          <h2 class="section-title">
            <i class="fas fa-building"></i> Organization Information
          </h2>
          <div class="form-grid">
            <div class="form-group">
              <label>Organization Name</label>
              <input type="text" :value="organization.name" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>VAT Number</label>
              <input type="text" :value="organization.vatNumber" readonly class="form-input" />
            </div>
            <div class="form-group full-width">
              <label>Address</label>
              <input type="text" :value="organization.address" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>City</label>
              <input type="text" :value="organization.city" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>ZIP/Postal Code</label>
              <input type="text" :value="organization.zipCode" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>Province</label>
              <input type="text" :value="organization.state" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>Email</label>
              <input type="email" :value="organization.email" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>Phone</label>
              <input type="tel" :value="organization.phone" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>SDI Code</label>
              <input type="text" :value="organization.sdiCode || 'Not provided'" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>PEC Email</label>
              <input type="email" :value="organization.pecEmail || 'Not provided'" readonly class="form-input" />
            </div>
          </div>
        </div>

        <!-- Manager Information -->
        <div class="form-section">
          <h2 class="section-title">
            <i class="fas fa-user-tie"></i> Company Manager
          </h2>
          <div class="form-grid">
            <div class="form-group">
              <label>Full Name</label>
              <input type="text" :value="organization.manager.name" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>Email</label>
              <input type="email" :value="organization.manager.email" readonly class="form-input" />
            </div>
            <div class="form-group">
              <label>Phone</label>
              <input type="tel" :value="organization.manager.phone" readonly class="form-input" />
            </div>
          </div>
        </div>

        <div class="form-actions">
          <router-link to="/organization/edit" class="btn btn-primary">
            <i class="fas fa-edit"></i> Edit Profile
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import BackgroundParticles from '../BackgroundParticles.vue';

// Sample data - in a real app, this would come from an API
const organization = ref({
  name: 'Acme Corporation',
  vatNumber: 'IT12345678901',
  address: '123 Business Street',
  city: 'Milan',
  state: 'MI',
  zipCode: '20100',
  email: 'info@acmecorp.com',
  phone: '+39 02 1234567',
  pecEmail: 'acmecorp@pec.it',
  sdiCode: 'ABC12345',
  manager: {
    name: 'Mario Rossi',
    email: 'mario.rossi@acmecorp.com',
    phone: '+39 333 1234567'
  }
});

// Set page title on mount
onMounted(() => {
  document.title = 'NetSecure-IQ - Organization Profile';
});
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css');

.organization-profile {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  position: relative;
  background-color: var(--bg-page);
  color: var(--text-primary);
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
}

.form-container {
  width: 100%;
  max-width: 900px;
  background: var(--bg-container);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 2.5rem;
  position: relative;
  z-index: 1;
  border: 1px solid var(--border-color);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.form-header {
  text-align: center;
  margin-bottom: 2rem;
}

.form-title {
  color: var(--text-primary);
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 0.5rem;
  line-height: 1.2;
}

.form-subtitle {
  color: var(--text-secondary);
  font-size: 1.125rem;
  margin: 0;
  opacity: 0.9;
}

.form-section {
  margin-bottom: 2.5rem;
  padding: 1.5rem;
  background: var(--bg-section);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.section-title {
  color: var(--text-primary);
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.section-title i {
  color: var(--primary);
  font-size: 1.1em;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.25rem;
  margin-bottom: 0.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: var(--text-primary);
  background-color: var(--bg-input);
  background-clip: padding-box;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
  appearance: none;
  cursor: not-allowed;
}

.form-input:focus {
  outline: 0;
  border-color: var(--primary);
  box-shadow: 0 0 0 0.2rem rgba(var(--primary-rgb), 0.25);
}

/* Dark mode overrides */
:root[data-theme="dark"] {
  --bg-page: #121212;
  --bg-container: rgba(30, 30, 30, 0.9);
  --bg-section: rgba(40, 40, 40, 0.7);
  --bg-input: #2d2d2d;
  --text-primary: #f5f5f5;
  --text-secondary: #a0a0a0;
  --border-color: #444;
  --primary: #4f46e5;
  --primary-rgb: 79, 70, 229;
}

/* Light mode overrides */
:root[data-theme="light"] {
  --bg-page: #f5f7fa;
  --bg-container: rgba(255, 255, 255, 0.95);
  --bg-section: rgba(245, 247, 250, 0.8);
  --bg-input: #ffffff;
  --text-primary: #1a1a1a;
  --text-secondary: #4b5563;
  --border-color: #e5e7eb;
  --primary: #4f46e5;
  --primary-rgb: 79, 70, 229;
}

/* Responsive styles */
@media (max-width: 768px) {
  .organization-profile {
    padding: 1rem;
  }
  
  .form-container {
    padding: 1.5rem;
    margin: 1rem;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .form-section {
    padding: 1rem;
  }
}

/* Button styles */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  text-align: center;
  white-space: nowrap;
  vertical-align: middle;
  user-select: none;
  border: 1px solid transparent;
  padding: 0.625rem 1.25rem;
  font-size: 0.9375rem;
  line-height: 1.5;
  border-radius: 6px;
  transition: all 0.2s ease-in-out;
  text-decoration: none;
  cursor: pointer;
}

.btn-primary {
  color: white;
  background-color: var(--primary);
  border-color: var(--primary);
}

.btn-primary:hover {
  background-color: #4338ca;
  border-color: #4338ca;
}

.btn i {
  margin-right: 0.5rem;
  font-size: 0.875em;
}

/* Form actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 2rem;
}

/* Focus styles */
:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

/* Disabled state for inputs */
.form-input:disabled,
.form-input[readonly] {
  background-color: var(--bg-input);
  opacity: 0.8;
  cursor: not-allowed;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background-color: var(--bg-input);
  color: var(--text-primary);
  font-size: 0.9375rem;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-input:read-only {
  background-color: var(--bg-secondary);
  cursor: not-allowed;
  opacity: 0.8;
}

.form-actions {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.5rem;
  border-radius: var(--radius-md);
  font-weight: 500;
  font-size: 1rem;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.btn-primary:hover {
  background-color: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.btn i {
  margin-right: 0.5rem;
}

/* Dark mode adjustments */
:root[data-theme="dark"] .form-container {
  background: rgba(30, 41, 59, 0.7);
  border-color: rgba(255, 255, 255, 0.1);
}

:root[data-theme="dark"] .form-input:read-only {
  background-color: rgba(30, 41, 59, 0.5);
  color: var(--text-secondary);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .organization-profile {
    padding: 1rem;
  }
  
  .form-container {
    padding: 1.5rem;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
