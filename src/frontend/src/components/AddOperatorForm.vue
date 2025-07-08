<template>
  <div>
    <h2>Add New Operator</h2>

    <form @submit.prevent="submitForm">
      <div>
        <label>Email:</label>
        <input v-model="email" type="email" required />
      </div>

      <div>
        <label>Password:</label>
        <input v-model="password" type="password" required />
      </div>

      <div>
        <label>Assign to Organization (optional):</label>
        <select v-model="selectedOrg">
          <option value="">-- None --</option>
          <option v-for="org in organizations" :key="org.id" :value="org.id">
            {{ org.name }}
          </option>
        </select>
      </div>

      <button type="submit">Add Operator</button>
      <p v-if="message">{{ message }}</p>
    </form>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";

export default {
  setup() {
    const email = ref("");
    const password = ref("");
    const organizations = ref([]);
    const selectedOrg = ref("");
    const message = ref("");

    onMounted(async () => {
      const res = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/organizations`, {
        headers: {
          Authorization: "Bearer " + localStorage.getItem("token"),
        },
      });
      const data = await res.json();
      organizations.value = data.organizations;
    });

    const submitForm = async () => {
      const payload = {
        email: email.value,
        password: password.value,
        role: "Operator", // fixed role
        organization_id: selectedOrg.value || null, // optional
      };

      const res = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/users`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + localStorage.getItem("token"),
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      if (res.ok) {
        message.value = "✅ Operator created successfully!";
      } else {
        message.value = "❌ Failed: " + (data.error || data.message);
      }
    };

    return {
      email,
      password,
      organizations,
      selectedOrg,
      submitForm,
      message,
    };
  },
};
</script>
