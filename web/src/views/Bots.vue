<template>
  <div>
    <NewBot class="ms-10 pa-1" />

    <h1 class="ms-10 pa-2">Chats</h1>
    <div v-for="chat in chats" :key="chat._id" class="ms-10 pa-2">
      <v-card max-width="500">
        <v-card-text>
          <v-row>
            <v-col cols="1">
              <v-tooltip>
                <span v-if="chat.verified">chat verified</span>
                <span v-if="!chat.verified"
                  >chat not verified in bot. Check in settings</span
                >
                <template v-slot:activator="{ props }">
                  <v-icon
                    v-if="chat.verified"
                    v-bind="props"
                    color="success"
                    icon="mdi-check-circle"
                  ></v-icon>

                  <v-icon
                    v-else
                    color="error"
                    v-bind="props"
                    icon="mdi-close-circle"
                  ></v-icon>
                </template>
              </v-tooltip>
            </v-col>
            <v-col cols="4">
              <div>{{ chat.botIdentifier }}</div>
            </v-col>
          </v-row>

          <p class="text-h4 text--primary">{{ chat.name }}</p>
        </v-card-text>
        <v-card-actions>
          <ChatSettings :chat="chat" @eventchatupdated="eventChatUpdated" />

          <v-btn
            variant="text"
            color="teal-accent-4"
            :to="'/posts?chatid=' + chat._id"
            :disabled="!chat.verified"
          >
            Posts
          </v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>
<script>
import NewBot from "@/components/NewBot.vue";
export default {
  components: {
    NewBot,
  },
  data() {
    return {
      posts: [],
      page: 1,
      limit: 10,
    };
  },
  created() {
    this.getPosts();
  },
  methods: {
    async getPosts() {
      const response = await fetch(
        `https://jsonplaceholder.typicode.com/posts?_page=${this.page}&_limit=${this.limit}`
      );
      const data = await response.json();
      this.posts = data;
    },
  },
};
</script>
