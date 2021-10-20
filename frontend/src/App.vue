<script setup>
</script>

<script>

export default {
  mounted() {
  }
}
</script>

<template>
  <router-link :to="{ name: 'ProfileMy'}" v-if="this.$store.state.authorized">Profile</router-link> <br />
  <router-link :to="{ name: 'Users'}" v-if="this.$store.state.authorized">Users</router-link> <br />
  <router-link :to="{ name: 'Friends'}" v-if="this.$store.state.authorized">Friends</router-link> <br />

  <router-link :to="{ name: 'Registration'}" v-if="!this.$store.state.authorized">Registration</router-link> <br />
  <router-link :to="{ name: 'Login'}" v-if="!this.$store.state.authorized">Login</router-link> <br />

  <router-view></router-view>
</template>

<style>
</style>
