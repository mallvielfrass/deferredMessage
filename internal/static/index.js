var app = new Vue({
  el: "#app",
  data: {
    msg: "Hello Vue!",
    isLogin: false,
    username: "",
    password: "",
    mail: "",
    passwordRepeat: "",
    mode: "login",
    errors: [],
  },
  mounted: async function () {
    console.log("mounted");
    this.checkLoginFromLocalStorage();
  },
  methods: {
    toggleMode: function () {
      this.mode = this.mode === "login" ? "register" : "login";
      console.log("toggleMode set:", this.mode);
    },
    checkFormAndAuth: async function (e) {
      this.errors = [];
      console.log(" checkFormAndAuth");
      if (this.mail == "") {
        this.errors.push("Mail required.");
      }

      if (this.password.length < 8) {
        this.errors.push("Password must be at least 8 characters.");
      }
      if (this.mode == "register") {
        if (this.name == "") {
          this.errors.push("Name required.");
        }
        if (this.password != this.passwordRepeat) {
          this.errors.push("Password repeat do not match.");
        }
      }

      if (0 < this.errors.length) {
        return;
      }
      switch (this.mode) {
        case "login":
          return this.login();
        case "register":
          return this.register();

        default:
          break;
      }
    },
    register: async function () {
      console.log("register");
      fetch("/api/nauth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          mail: this.username,
          password: this.password,
        }),
      });
    },
    login: async function () {
      console.log("login");
      fetch("/api/nauth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          mail: this.username,
          password: this.password,
        }),
      });
    },
    checkLoginFromLocalStorage: async function () {
      const token = localStorage.getItem("token");
      //const dateExpiration = localStorage.getItem("expiration");
      //const logTime= localStorage.getItem("logtime");
      if (token) {
        this.isLogin = false;
        fetch("/api/auth/user/ping", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `${token}`,
          },
        })
          .then((response) => {
            if (response.status === 200 && response.body?.message === "pong") {
              this.isLogin = true;
              return console.log("login success:", this.isLogin);
            }
            localStorage.removeItem("token");
            return console.log("login fail:", this.isLogin, response.body);
          })
          .catch((error) => {
            console.log("login fail:", this.isLogin, error);
            this.isLogin = false;
            localStorage.removeItem("token");
          });
      } else {
        console.log("not login");
        this.isLogin = false;
      }
    },
  },
});
