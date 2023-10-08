<template>
  <v-container>
    <v-row no-gutters justify="end">
      <v-col v-for="page in pages" :key="page.url" cols="auto" class="pa-2">
        <v-btn variant="outlined" :to="page.url">{{ page.name }}</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
export default {
  data() {
    return {
      pages: [
        {
          name: "posts",
          url: "/posts",
        },

        { name: "chats", url: "/chats" },
        { name: "bots", url: "/bots" },
        { name: "account", url: "/account" },
      ],
    };
  },
};
</script>
