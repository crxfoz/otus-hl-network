<script>
import axios from 'axios'


export default {
  data() {
    return {
      error: "",
      login: "",
      password: "",
    }
  },

  methods: {
    loginBtn(event) {
      let bodyFormData = new FormData();
      bodyFormData.append("username", this.login);
      bodyFormData.append("password", this.password);

      axios.post("/api/v1/auth", bodyFormData).then((response) => {
        this.$store.commit('auth', {
          authorized: true,
          user_id: response.data.id,
          username: response.data.username,
          token: response.data.token,
        });

        console.log("auth successful")

        this.$router.push('/me')
      }).catch((_) => {
        this.error = "wrong creds"
      })
    }
  }
}
</script>

<template>
  <form>
    Username: <input v-model="login" type="text" />
    Password: <input v-model="password" type="password" />
  </form>

  <button @click="loginBtn">Login</button> <br />

  {{ error }}

</template>