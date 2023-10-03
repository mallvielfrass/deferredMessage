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
        <NotifyAlert ref="notifyAlert" />
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
import NotifyAlert from "@/components/NotifyAlert.vue";
import { registerUser, loginUser } from "@/api/auth.js";

export default {
  data() {
    return {
      dialog: true,
      errors: [],
      mode: "login",
    };
  },
  components: {
    LoginForm,
    RegisterForm,
    NotifyAlert,
  },
  methods: {
    logOut() {
      this.isLogin = false;
      console.log("logout");
    },

    AuthFrame() {
      this.dialog = true;
    },
    async sendFormRegister() {
      const checkForm = this.$refs.registerFormRef.getValidateState();
      console.log("sendForm:", checkForm);
      if (!checkForm) {
        return this.$refs.notifyAlert.showNotification(
          "Please fill in all fields correctly"
        );
      }
      const { username, email, password } =
        this.$refs.registerFormRef.getData();
      console.log(username, email, password);
      const response = await registerUser(username, email, password);
      console.log(response);
      if (response.errors.length != 0) {
        response.errors.map((error) => {
          this.$refs.notifyAlert.showNotification(error);
        });
        return;
      }
      this.dialog = false;
      this.isLogin = response.isLogin;
    },
    async sendFormLogin() {
      const checkForm = this.$refs.loginFormRef.getValidateState();
      console.log("sendForm:", checkForm);
      if (!checkForm) {
        return this.$refs.notifyAlert.showNotification(
          "Please fill in all fields correctly"
        );
      }
      const { email, password } = this.$refs.loginFormRef.getData();
      console.log(email, password);
      const response = await loginUser(email, password);
      console.log(response);
      if (response.errors.length != 0) {
        response.errors.map((error) => {
          this.$refs.notifyAlert.showNotification(error);
        });
        return;
      }
      this.dialog = false;
      this.isLogin = response.isLogin;
    },
    sendForm() {
      if (this.mode === "login") {
        return this.sendFormLogin();
      }

      if (this.mode === "register") {
        return this.sendFormRegister();
      }
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
