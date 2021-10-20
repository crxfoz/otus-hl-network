<script>
import axios from 'axios'

export default {
  data() {
    return {
      users: [],
    }
  },

  mounted() {
    axios.get("/api/v1/users", {
      headers: {
        'Authorization': 'Bearer '+ this.$store.state.token
      }
    }).then((response) => {
      for(let user of response.data) {
        this.users.push({
          id: user.user_id,
          first_name: user.first_name,
          last_name: user.last_name,
        })
      }
    })
  },

  methods: {
    goToUser(id) {
      this.$router.push('/profile/' + id);
    }
  }
}

</script>

<template>
  Users: <br />
  <ul>
    <li v-for="user in users" @click="goToUser(user.id)"> {{ user.first_name }} {{ user.last_name }}</li>
  </ul>

</template>