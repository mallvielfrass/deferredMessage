export const registerUser = async (username, mail, password) => {
  const errors = [];
  let isLogin = false;
  console.log("register");
  await fetch("/api/nauth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      mail: mail,
      name: username,
      password: password,
    }),
  })
    .then(async (response) => {
      const body = await response.json();
      if (response.status == 400) {
        switch (body.error) {
          case "user already exist": {
            errors.push("User already exist");
            break;
          }

          case "no body": {
            errors.push("Empty field for mail or name");
          }
          default: {
            errors.push("Unknown error");
            console.log("Unknown error", body);
          }
        }
        return;
      }
      if (!response.status == 200) {
        errors.push("Unknown error");
        console.log("unknown error", body, response.status);
        return;
      }
      if (body.status !== "success") {
        errors.push("Unknown register error");
        console.log("unknown register error", body, response.status);
        return;
      }

      const session = body.session;
      localStorage.setItem("token", session.id);
      isLogin = true;
    })

    .catch((error) => {
      this.errors.push("Unknown register error");
      console.log("unknown catch error", error);
    });
  return { isLogin, errors };
};
export const loginUser = async (mail, password) => {
  const errors = [];
  let isLogin = false;

  console.log("login");
  await fetch("/api/nauth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      mail: mail,
      password: password,
    }),
  }).then(async (response) => {
    const body = await response.json();
    if (response.status == 400) {
      switch (body.error) {
        case "no body": {
          errors.push("Empty field for mail or password");
          break;
        }
        case "user or password incorrect": {
          errors.push("User or password incorrect");
          break;
        }

        default:
          errors.push("Unknown error");
          console.log("Unknown error", body);
          break;
      }
      return;
    }
    if (!response.status == 200) {
      errors.push("Unknown error");
      console.log("unknown error", body, response.status);
      return;
    }
    if (body.status !== "success") {
      errors.push("Unknown login error");
      console.log("unknown login error", body, response.status);
      return;
    }
    const session = body.session;
    localStorage.setItem("token", session.id);
    isLogin = true;
  });
  return { isLogin, errors };
};
export const getToken = () => {
  return localStorage.getItem("token");
};
