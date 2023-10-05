<template>
  <div>
    <NewChat @eventchatcreated="eventChatCreated" class="ms-10 pa-1" />
    <h1 class="ms-10 pa-2">Chats</h1>
    <div v-for="chat in chats" :key="chat._id" class="ms-10 pa-2">
      <v-card max-width="500">
        <v-card-text>
          <div>{{ chat.networkIdentifier }}</div>
          <p class="text-h4 text--primary">{{ chat.name }}</p>
        </v-card-text>
        <v-card-actions>
          <v-btn variant="text" color="teal-accent-4" @click="reveal = true">
            Settings
          </v-btn>
          <v-btn
            variant="text"
            color="teal-accent-4"
            :to="'/posts?chatid=' + chat._id"
          >
            Posts
          </v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>
<script>
import NewChat from "@/components/NewChat.vue";
import { getChats } from "@/api/networks";
export default {
  components: {
    NewChat,
  },
  data() {
    return {
      chats: [],
      page: 1,
      limit: 10,
    };
  },
  created() {
    this.getChats();
  },

  methods: {
    eventChatCreated() {
      console.log("eventChatCreated");
      this.getChats();
    },
    async getChats() {
      this.chats = await getChats();
    },
  },
};
</script>
<style>
.card-with-padding {
  padding-left: 20px; /* Adjust the padding value as needed */
}
</style>
