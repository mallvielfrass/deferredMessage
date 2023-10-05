import { getToken } from "@/api/auth";
export const createChat = async (name, networkIdentifier) => {
  const token = getToken();
  let resp = {
    chat: null,
    error: null,
  };
  const chat = await fetch("/api/auth/user/chats", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `${token}`,
    },
    body: JSON.stringify({
      name: name,
      networkIdentifier: networkIdentifier,
    }),
  })
    .then(async (response) => {
      const body = await response.json();
      if (response.status != 200) {
        resp.error = body.error;
        return;
      }
      return body;
    })
    .catch((error) => {
      console.log(error);
      resp.error = error;
    });

  console.log("chat", chat);
  return resp;
};
