<template>
  <div class="login">
    <div class="login-box">
      <h3>Login to your account</h3>
      <form @submit.prevent="login">
        <input
          v-model="email"
          type="email"
          placeholder="Email"
          required
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          required
        />
        <button :disabled="loading" type="submit">
          {{ loading ? "Loading..." : "Login" }}
        </button>
      </form>

      <div class="link">
        <span>Don't have an account?</span>
        <router-link to="/register">Register</router-link>
      </div>

      <p v-if="message" :class="{ success: successMessage, error: !successMessage }">
        {{ message }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"

const email = ref("")
const password = ref("")
const message = ref("")
const successMessage = ref(false)
const loading = ref(false)
const router = useRouter()

const login = async () => {
  loading.value = true;
  try {
    const res = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: email.value, password: password.value }),
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || data.message || "Login failed");
    }

    message.value = "Login successful! Redirecting...";
    successMessage.value = true;

    const role = data.role?.toLowerCase();
    localStorage.setItem("token", data.token);
    localStorage.setItem("user_id", data.user_id);
    localStorage.setItem("role", role);
    console.log("User role:", role);

    // 🔁 Decide final redirect path
    let redirectTo = "/dashboard";

    if (role === "user") {
      redirectTo = "/routertable";
    } else if (role === "operator" && !data.organization_id) {
      redirectTo = "/organizationForm";
    }

    // ✅ Single final redirect
    setTimeout(() => {
      router.push(redirectTo);
    }, 500);

  } catch (err) {
    message.value = "Error : " + err.message;
    successMessage.value = false;
  } finally {
    loading.value = false;
  }
};
</script>

