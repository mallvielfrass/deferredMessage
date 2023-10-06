<template>
  <v-container>
    <v-row>
      <v-col>
        <v-btn
          size="x-large"
          variant="outlined"
          color="green"
          @click="dialog = true"
        >
          Create New Post</v-btn
        >
      </v-col>
    </v-row>
  </v-container>

  <v-dialog v-model="dialog" max-width="500px" max-height="70%">
    <v-card>
      <div>
        <v-form v-model="valid" @submit.prevent="">
          <v-container>
            <v-row justify="center" align="center">
              <v-card-title>New Post</v-card-title>
            </v-row>

            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-text-field
                  v-model="title"
                  label="Title"
                  required
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-textarea label="Text of the post"></v-textarea>
              </v-col>
            </v-row>
            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-select
                  label="Select bot"
                  v-model="socialBot"
                  @update:modelValue="selectSocialBot"
                  :items="socialBotList"
                  :item-props="itemProps"
                >
                </v-select>
              </v-col>
            </v-row>
          </v-container>
        </v-form>
      </div>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="secondary"
          @click="closeDialog"
          justify="right"
          align="right"
          >Close</v-btn
        >
        <v-btn color="secondary" @click="sendForm" justify="right" align="right"
          >Send</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { getBots, getChats } from "@/api/bots.js";
export default {
  data() {
    return {
      dialog: false,
      socialBot: {},
      socialBotList: [],
      showChats: false,
      chats: [],
    };
  },
  mounted() {
    this.getSocialBotList();
  },
  methods: {
    async getChatsList() {
      const response = await getChats();
      console.log(response);
      this.chats = response;
    },
    async getSocialBotList() {
      this.socialBotList = await getBots();
    },
    itemProps(item) {
      return {
        title: item.name,
        subtitle: item.identifier,
      };
    },
    async selectSocialBot(bot) {
      console.log("Bot selected:", bot);
      await this.getChatsList();
    },
    closeDialog() {
      this.dialog = false;
    },
  },
};
</script>
<style>
.box {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 400px;
}
</style>
