<template>
  <div>
    <h2>Add New Operator</h2>

    <form @submit.prevent="submitForm">
      <div>
        <label>Email:</label>
        <input v-model="email" type="email" required />
      </div>

      <div>
        <label>First Name:</label>
        <input v-model="firstName" type="text" required />
      </div>

      <div>
        <label>Last Name:</label>
        <input v-model="lastName" type="text" required />
      </div>

      <button type="submit">Add Operator</button>
      <p v-if="message">{{ message }}</p>
    </form>
  </div>
</template>

<script>
import { ref } from "vue";
import { jwtDecode } from 'jwt-decode';

export default {
  setup() {
    const email = ref("");
    const firstName = ref("");
    const lastName = ref("");
    const message = ref("");

    const submitForm = async () => {
  try {
    const token = localStorage.getItem("token");
    const decoded = jwtDecode(token);
    const organization_id = decoded.organization_id;

    const payload = {
      email: email.value,
      first_name: firstName.value,   // snake_case ici
      last_name: lastName.value,     // snake_case ici
      role: "operator",
    };

    if (organization_id) {
      payload.organization_id = String(organization_id);  // Convertir en string
    }

    console.log("Payload to send:", JSON.stringify(payload, null, 2));

    const res = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/users`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
      body: JSON.stringify(payload),
    });

    let data;
    try {
      data = await res.json();
    } catch (err) {
      data = { error: "Invalid response from server" };
    }

    if (res.ok) {
      message.value = "✅ Operator created successfully!";
    } else {
      message.value = "❌ Failed: " + (data.error || data.message);
    }
  } catch (err) {
    message.value = "❌ Internal error: " + err.message;
  }
};


    return {
      email,
      firstName,
      lastName,
      submitForm,
      message,
    };
  },
};
</script>
