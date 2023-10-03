<template>
  <v-form v-model="valid" @submit.prevent="">
    <v-container>
      <v-row justify="center" align="center">
        <v-card-title>Login</v-card-title></v-row
      >

      <v-row justify="center" align="center">
        <v-col cols="12" md="4">
          <v-text-field
            v-model="email"
            :rules="emailRules"
            label="E-mail"
            required
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center" align="center">
        <v-col cols="12" md="4">
          <v-text-field
            v-model="password"
            :rules="passRules"
            label="password"
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
      },
      (value) => {
        if (8 <= value.length) return true;
        return "Password must be greater than 8 characters.";
      },
    ],
  }),
  methods: {
    getValidateState() {
      return this.valid ? this.valid : false;
    },
    getData() {
      return {
        email: this.email,
        password: this.password,
      };
    },
  },
};
</script>
