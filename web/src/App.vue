<template>
  <v-app>
    <Bar @setloginstate="setloginstate" />
    <div class="padding-top"></div>
    <div v-if="isLogin">
      <router-view />
    </div>
  </v-app>
</template>

<script>
import Bar from "./components/Bar.vue";
export default {
  name: "App",
  data: () => ({
    isLogin: false,
  }),
  components: { Bar },
  methods: {
    setloginstate(state) {
      console.log("setloginstate", state);
      this.isLogin = state;
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
