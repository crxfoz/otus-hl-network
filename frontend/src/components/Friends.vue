<script>
import axios from 'axios'

export default {
  data() {
    return {
      friends: [],
    }
  },

  mounted() {
    axios.get("/api/v1/friends", {
      headers: {
        'Authorization': 'Bearer '+ this.$store.state.token
      }
    }).then((response) => {
      for(let user of response.data) {
        this.friends.push({
          id: user.user_id,
          first_name: user.first_name,
          last_name: user.last_name,
        })
      }
    })
  },

  methods: {
    removeFriend(id) {
      axios.delete("/api/v1/friends/" + id, {
        headers: {
          'Authorization': 'Bearer '+ this.$store.state.token
        }
      }).then((_) => {
        console.log("friend removed")
      })
    },

    goToUser(id) {
      this.$router.push('/profile/' + id);
    }
  }
}

</script>

<template>
  List: <br />
  <ul>
    <li v-for="friend in friends" @click="goToUser(friend.id)"> {{ friend.first_name }} {{ friend.last_name }}</li>
  </ul>

</template>

