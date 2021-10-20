<script>
import axios from 'axios'


export default {
  data() {
    return {
      username: "",
      password: "",
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
    register(event) {
      this.error = ""

      if(this.username.length < 5) {
        this.error = "Username should be at least 5 characters"
        return
      }

      if(this.password.length < 5) {
        this.error = "Password should be at least 5 characters"
        return
      }

      if(this.first_name.length === 0) {
        this.error = "FirstName should be filled"
        return
      }

      if(this.last_name.length === 0) {
        this.error = "LastName should be filled"
        return
      }

      let interests = this.interests.split(';')

      axios.post("/api/v1/register", {
        username: this.username,
        password: this.password,
        first_name: this.first_name,
        last_name: this.last_name,
        age: parseInt(this.age),
        city: this.city,
        gender: this.gender,
        interests: interests,
      }).then((_) => {
        this.$router.push("/login")
      }).catch((_) => {
        this.error = "something went wrong"
      })

    }
  }
}
</script>

<template>
  <form>
    Username: <input v-model="username" type="text" /> <br />
    Password: <input v-model="password" type="password"/> <br />
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

  <button @click="register">Register</button> <br />

  {{ error }}


</template>