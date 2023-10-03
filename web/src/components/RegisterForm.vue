<template>
  <v-form v-model="valid" @submit.prevent="">
    <v-container>
      <v-row justify="center" align="center">
        <v-card-title>Register</v-card-title></v-row
      >
      <v-row justify="center" align="center">
        <div v-if="errors.length">
          <b>Please correct the following error(s):</b>
          <ul>
            <li v-for="error in errors">{{ error }}</li>
          </ul>
        </div>
      </v-row>
      <v-row justify="center" align="center">
        <v-col cols="1" md="6">
          <v-text-field
            v-model="username"
            :rules="nameRules"
            label="Username"
            required
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center" align="center">
        <v-col cols="1" md="6">
          <v-text-field
            v-model="email"
            :rules="emailRules"
            label="E-mail"
            required
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center" align="center">
        <v-col ccols="1" md="6">
          <v-text-field
            v-model="password"
            :rules="passRules"
            label="password"
            type="password"
            required
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center" align="center">
        <v-col cols="1" md="6">
          <v-text-field
            v-model="passwordRepeat"
            :rules="passRepeatRules"
            label="repeat password"
            type="password"
            required
          ></v-text-field>
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>
<script>
export default {
  data: () => ({
    valid: false,

    errors: [],
    username: "",
    nameRules: [
      (value) => {
        if (value) return true;

        return "Name is required.";
      },
      (value) => {
        if (4 <= value?.length <= 10) return true;
        return "Name must be greater then 4 and less than 10 characters.";
      },
    ],
    email: "",
    emailRules: [
      (value) => {
        if (value) return true;

        return "E-mail is requred.";
      },
      (value) => {
        if (/.+@.+\..+/.test(value)) return true;

        return "E-mail must be valid.";
      },
    ],
    password: "",
    passRules: [
      (value) => {
        if (value) return true;
        return "Password is requred.";
      },
      (value) => {
        if (8 <= value.length) return true;
        return "Password must be greater than 8 characters.";
      },
    ],
    passwordRepeat: "",
  }),
  computed: {
    passRepeatRules() {
      return [
        (value) => {
          if (value) return true;
          return "Repeat password is requred.";
        },
        (value) => {
          console.log("repeat", value, this.password);
          if (this.password === value) return true;
          return "Passwords must be the same.";
        },
      ];
    },
  },
  methods: {
    getValidateState() {
      if (this.password !== this.passwordRepeat) {
        this.valid = false;
      }
      return this.valid ? this.valid : false;
    },
    getData() {
      return {
        username: this.username,
        email: this.email,
        password: this.password,
      };
    },
  },
};
</script>
