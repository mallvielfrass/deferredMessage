<template>
  <v-app-bar color="#6A76AB" dark shrink-on-scroll prominent class="my-app-bar">
    <div v-if="isLogin" class="ml-auto padding-right">
      <v-btn variant="outlined" @click="logOut"> Logout</v-btn>
    </div>
    <div v-else class="ml-auto padding-right">
      <v-btn variant="outlined" @click="AuthFrame"> Login</v-btn>
    </div>
    <v-dialog v-model="dialog" max-width="500px">
      <v-card>
        <div v-if="mode === 'login'">
          <LoginForm ref="loginFormRef" />
        </div>
        <div v-else>
          <RegisterForm ref="registerFormRef" />
        </div>
        <v-card-actions>
          <div v-if="mode === 'login'" class="change-mode-container">
            Don't have an account?
            <v-btn color="secondary" @click="toggleMode">Register</v-btn>
          </div>
          <div v-else class="change-mode-container">
            Already have an account?
            <v-btn color="secondary" @click="toggleMode">Login</v-btn>
          </div>
          <v-spacer></v-spacer>
          <v-btn
            color="secondary"
            @click="closeDialog"
            justify="right"
            align="right"
            >Close</v-btn
          >
          <v-btn
            color="secondary"
            @click="sendForm"
            justify="right"
            align="right"
            >Send</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app-bar>
</template>
<script>
import LoginForm from "@/components/LoginForm.vue";
import RegisterForm from "@/components/RegisterForm.vue";
export default {
  data() {
    return {
      dialog: true,

      errors: [],
      mode: "login",
      username: "",
      mail: "",
      password: "",
      passwordRepeat: "",
      parentData: "initial value",
    };
  },
  components: {
    LoginForm,
    RegisterForm,
  },
  methods: {
    checkLoginForm() {
      return this.$refs.loginFormRef.getValidateState();
    },
    logOut() {
      this.isLogin = false;
      console.log("logout");
    },

    AuthFrame() {
      this.dialog = true;
    },
    sendForm() {
      console.log("sendForm:", this.checkLoginForm());
    },
    closeDialog() {
      this.dialog = false;
    },
    toggleMode() {
      this.mode = this.mode === "login" ? "register" : "login";
    },
  },
};
</script>
<style>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.login-form {
  border: 1px solid #ccc;
  padding: 20px;
  border-radius: 5px;
  background-color: #fff;
}
</style>
