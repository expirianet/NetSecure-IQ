<template>
  <div>
    <h1>Welcome to the Dashboard</h1>

    <div v-if="isAdminOrOperator">
      <button @click="goToRouterInfo">Router Info</button>
      <button @click="addUser">Add User</button>
      <button v-if="needsOrganization" @click="goToOrganizationInfo">Organization info</button>
      <button v-if="isAdmin" @click="addOperator">Add Operator</button>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    role() {
      return localStorage.getItem("role");
    },
    isAdmin() {
      return this.role === "administrator";
    },
    isOperator() {
      return this.role === "operator";
    },
    isAdminOrOperator() {
      return this.isAdmin || this.isOperator;
    },
    hasOrganization() {
      return localStorage.getItem("organization_id") !== "";
    },
    needsOrganization() {
      return this.isOperator && !this.hasOrganization;
    } 
  },
  methods: {
    addUser() {
      this.$router.push("/adduser");
    },
    addOperator() {
      this.$router.push("/addoperator");
    },
    goToRouterInfo() {
      this.$router.push("/routertable");
    },
    goToOrganizationInfo() {
      this.$router.push("/organizationForm");
    }
  }
};
</script>
