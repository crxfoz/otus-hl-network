<script>

import axios from 'axios'

export default {
  data() {
    return {
      first_name: "",
      last_name: "",
      age: 0,
      gender: "",
      city: "",
      interests: [],

      is_friend: false,
      self_page: false,
    }
  },

  mounted() {
    if (this.$route.params.id == this.$store.state.user_id) {
      this.self_page = true;
    }

    axios.get("/api/v1/profile/" + this.$route.params.id, {
      headers: {
        'Authorization': 'Bearer '+ this.$store.state.token
      }
    }).then((response) => {
      this.first_name = response.data.first_name;
      this.last_name = response.data.last_name;
      this.age = response.data.age;
      this.gender = response.data.gender;
      this.city = response.data.city;
      this.interests = response.data.interests;
    })

    // TODO: use state for friends checking
    axios.get("/api/v1/friends", {
      headers: {
        'Authorization': 'Bearer '+ this.$store.state.token
      }
    }).then((response) => {
      for(let friend of response.data) {
        if (friend.user_id == this.$route.params.id) {
          this.is_friend = true;
        }
      }
    })
  },

  methods: {
    addFriend(id) {
      axios.post("/api/v1/friends/" + id, {}, {
        headers: {
          'Authorization': 'Bearer '+ this.$store.state.token
        }
      }).then((_) => {
        console.log("friend added")
        this.is_friend = true;
      })
    },
    removeFriend(id) {
      axios.delete("/api/v1/friends/" + id, {
        headers: {
          'Authorization': 'Bearer '+ this.$store.state.token
        }
      }).then((_) => {
        console.log("friend removed")
        this.is_friend = false;
      })
    }
  }
}

</script>

<template>
  Name: {{ first_name }} {{ last_name }} <br />
  Age: {{ age }} <br />
  Gender: {{ gender }} <br />
  City: {{ city }} <br />
  Interests:
  <ul>
    <li v-for="interest in interests">{{ interest }}</li>
  </ul>

  <button @click="removeFriend(this.$route.params.id)" v-if="this.is_friend && !this.self_page">Remove Friend</button>
  <button @click="addFriend(this.$route.params.id)" v-if="!this.is_friend && !this.self_page">Add Friend</button>

</template>