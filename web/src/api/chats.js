import { getToken } from "@/api/auth";
export const createChat = async (name, botIdentifier) => {
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
      botIdentifier: botIdentifier,
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
  resp.chat = chat;
  return resp;
};
export const updateChat = async (id, updatedObj) => {
  const token = getToken();
  let resp = {
    chat: null,
    error: null,
  };
  const chat = await fetch(`/api/auth/user/chats/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: `${token}`,
    },
    body: JSON.stringify(updatedObj),
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
  resp.chat = chat;
  return resp;
};
