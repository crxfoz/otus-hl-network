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
    }
  },

  mounted() {
    axios.get("/api/v1/profile", {
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
  },

  methods: {
    profileEdit(event) {
      this.$router.push('/update');
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


  <button @click="profileEdit">Edit Profile</button>

</template>