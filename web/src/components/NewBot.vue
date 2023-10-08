<template>
  <v-container>
    <v-row>
      <v-col>
        <v-btn
          size="x-large"
          variant="outlined"
          color="green"
          @click="dialog = true"
        >
          Create New Bot</v-btn
        >
      </v-col>
    </v-row>
  </v-container>

  <v-dialog v-model="dialog" max-width="500px" max-height="70%">
    <v-card>
      <div>
        <v-form v-model="valid" @submit.prevent="" ref="newChatForm">
          <v-container>
            <v-row justify="center" align="center">
              <v-card-title>New Bot</v-card-title>
            </v-row>
            <v-row justify="center" align="center">
              <div v-if="errors.length">
                <b>Please correct the following error(s):</b>
                <ul>
                  <li v-for="error in errors">{{ error }}</li>
                </ul>
              </div>
            </v-row>
          </v-container>
        </v-form>
      </div></v-card
    ></v-dialog
  >
</template>

<script>
export default {
  data() {
    return {
      dialog: false,
      valid: false,
      name: "",
      nameRules: [
        (v) => !!v || "Name is required",
        (v) => (v && v.length <= 10) || "Name must be less than 10 characters",
      ],
      errors: [],
    };
  },
  methods: {
    closeDialog() {
      this.dialog = false;
      this.errors = [];
    },
    sendForm() {
      if (this.$refs.newChatForm.validate()) {
        this.dialog = false;
        this.errors = [];
      }
    },
  },
};
</script>
