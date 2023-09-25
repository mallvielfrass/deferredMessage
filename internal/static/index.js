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
      this.errors = [];
      console.log("register");
      fetch("/api/nauth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          mail: this.mail,
          name: this.username,
          password: this.password,
        }),
      })
        .then(async (response) => {
          const body = await response.json();
          if (response.status == 400) {
            switch (body.error) {
              case "user already exist": {
                this.errors.push("User already exist");
                break;
              }

              case "no body": {
                this.errors.push("Empty field for mail or name");
              }
              default: {
                this.errors.push("Unknown error");
                console.log("Unknown error", body);
              }
            }
            return;
          }
          if (!response.status == 200) {
            this.errors.push("Unknown error");
            console.log("unknown error", body, response.status);
            return;
          }
          if (body.status !== "success") {
            this.errors.push("Unknown register error");
            console.log("unknown register error", body, response.status);
            return;
          }

          const session = body.session;
          localStorage.setItem("token", session.id);
          this.isLogin = true;
        })

        .catch((error) => {
          this.errors.push("Unknown register error");
          console.log("unknown catch error", error);
        });
    },
    login: async function () {
      console.log("login");
      this.errors = [];
      fetch("/api/nauth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          mail: this.mail,
          password: this.password,
        }),
      }).then(async (response) => {
        const body = await response.json();
        if (response.status == 400) {
          switch (body.error) {
            case "no body": {
              this.errors.push("Empty field for mail or password");
              break;
            }
            case "user or password incorrect": {
              this.errors.push("User or password incorrect");
              break;
            }

            default:
              this.errors.push("Unknown error");
              console.log("Unknown error", body);
              break;
          }
          return;
        }
        if (!response.status == 200) {
          this.errors.push("Unknown error");
          console.log("unknown error", body, response.status);
          return;
        }
        if (body.status !== "success") {
          this.errors.push("Unknown login error");
          console.log("unknown login error", body, response.status);
          return;
        }
        const session = body.session;
        localStorage.setItem("token", session.id);
        this.isLogin = true;
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
          .then(async (responseStream) => {
            const response = await responseStream.json();
            if (responseStream.status === 200 && response.message === "pong") {
              this.isLogin = true;
              return console.log("login success:", this.isLogin);
            }
            localStorage.removeItem("token");
            return console.log("login fail:", this.isLogin, response);
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
