Vue.component("my-component", {
  template: `#my-component-template`,
  data() {
    return {
      message: "Привет, мир!",
    };
  },
  methods: {
    changeMessage() {
      this.message = "Новое сообщение!";
    },
  },
});
