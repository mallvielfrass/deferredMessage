import { getToken } from "@/api/auth";
export const getNetworks = async () => {
  const networks = [];
  const token = getToken();
  await fetch("/api/auth/user/networks", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `${token}`,
    },
  }).then(async (response) => {
    const body = await response.json();
    if (response.status != 200) {
      return;
    }

    networks.push(...body.networks);
  });
  return networks;
};
export const getChats = async () => {
  const chats = [];
  const token = getToken();
  await fetch("/api/auth/user/chats", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `${token}`,
    },
  }).then(async (response) => {
    const body = await response.json();
    if (response.status != 200) {
      return;
    }
    if (!body.chats) {
      return;
    }
    if (body.chats.length == 0) {
      return;
    }

    chats.push(...body.chats);
  });
  return chats;
};
