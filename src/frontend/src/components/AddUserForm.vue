<template>
  <div>
    <h2>Add New User</h2>

    <form @submit.prevent="submitForm">
      <div>
        <label>First Name:</label>
        <input v-model="firstName" type="text" required />
      </div>

      <div>
        <label>Last Name:</label>
        <input v-model="lastName" type="text" required />
      </div>

      <div>
        <label>Email:</label>
        <input v-model="email" type="email" required />
      </div>

      <div v-if="isAdmin">
        <label>Select Organization:</label>
        <select v-model="selectedOrg" required>
          <option v-for="org in organizations" :key="org.id" :value="org.id">
            {{ org.name }}
          </option>
        </select>
      </div>

      <button type="submit">Add User</button>
      <p v-if="message">{{ message }}</p>
    </form>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";

export default {
  setup() {
    const email = ref("");
    const firstName = ref("");
    const lastName = ref("");
    const organizations = ref([]);
    const selectedOrg = ref(null);
    const message = ref("");
    const role = localStorage.getItem("role");
    const userOrgId = localStorage.getItem("organization_id");

    const isAdmin = role === "administrator";

    onMounted(async () => {
      if (isAdmin) {
        const res = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/organizations`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token"),
          },
        });
        const data = await res.json();
        organizations.value = data.organizations;
      }
    });

    const submitForm = async () => {
      const payload = {
        email: email.value,
        first_name: firstName.value,
        last_name: lastName.value,
        organization_id: isAdmin ? selectedOrg.value : userOrgId,
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
        message.value = "✅ User created and password sent via email!";
      } else {
        message.value = "❌ Failed: " + (data.error || data.message);
      }
    };

    return {
      email,
      firstName,
      lastName,
      organizations,
      selectedOrg,
      submitForm,
      message,
      isAdmin,
    };
  },
};
</script>
