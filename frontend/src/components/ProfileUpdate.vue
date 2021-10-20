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
      interests: "",

      error: "",
    }
  },

  methods: {
    updateBtn(event) {
      let interests = this.interests.split(';')

      if(this.first_name.length === 0) {
        this.error = "FirstName should be filled"
        return
      }

      if(this.last_name.length === 0) {
        this.error = "LastName should be filled"
        return
      }

      axios.post("/api/v1/profile", {
        first_name: this.first_name,
        last_name: this.last_name,
        age: parseInt(this.age),
        city: this.city,
        gender: this.gender,
        interests: interests,
      }, {
        headers: {
          'Authorization': 'Bearer '+ this.$store.state.token
        }
      }).then((_) => {
        this.$router.push("/me")
      }).catch((_) => {
        this.error = "something went wrong"
      })
    }
  }
}

</script>

<template>
  <form>
    FirstName: <input v-model="first_name" type="text" /> <br />
    LastName: <input v-model="last_name" type="text" /> <br />
    Age: <input v-model="age" type="text" /> <br />
    Gender: <select v-model="gender">
    <option value="female">Female</option>
    <option value="male">Male</option>
  </select> <br />
    City: <input v-model="city" type="text" /> <br />
    Interests: <input v-model="interests" type="text" /> <br />
  </form> <br />

  <button @click="updateBtn">Update</button> <br />

  {{ error }}


</template>