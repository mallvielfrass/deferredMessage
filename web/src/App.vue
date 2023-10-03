<template>
  <v-app>
    <Bar />
    <div class="padding-top"></div>
  </v-app>
</template>

<script>
import Bar from "./components/Bar.vue";
export default {
  name: "App",
  data: () => ({
    isLogin: false,
    dialog: false,
  }),
  components: { Bar },
  methods: {
    logout() {
      console.log("logout");
      this.isLogin = false;
    },
    login() {
      console.log("login");
      this.isLogin = true;
    },
    openDialog() {
      this.dialog = true;
    },
    closeDialog() {
      this.dialog = false;
    },
  },
};
</script>
<style>
.padding-top {
  padding-top: 80px;
}
.my-app-bar {
  z-index: 0; /* Установите подходящее значение Z-индекса */
}
.padding-right {
  padding-right: 20px;
}
.bar-container {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  border: 1px solid #ccc;

  border-radius: 5px;
  background-color: #fff;
}
.bar {
  border: 1px solid #ccc;
  padding: 20px;
  border-radius: 5px;
  background-color: #fff;
}
</style>
