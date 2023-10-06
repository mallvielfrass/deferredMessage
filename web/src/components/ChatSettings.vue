<template>
  <v-btn variant="text" color="teal-accent-4" @click="dialog = true">
    Settings
  </v-btn>

  <v-dialog v-model="dialog" max-width="500px" max-height="70%">
    <v-card>
      <div>
        <v-form v-model="valid" @submit.prevent="" ref="chatSettingsForm">
          <v-container>
            <v-row justify="center" align="center">
              <v-card-title>Chat Settings</v-card-title>
            </v-row>

            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-text-field
                  v-model="chatCopy.name"
                  label="Title"
                  required
                  :rules="[(v) => !!v || 'Title is required']"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-text-field
                  v-model="chatCopy.linkOrIdInBot"
                  label="Link or ID in the bot"
                ></v-text-field>
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
import { updateChat } from "@/api/chats.js";
export default {
  props: {
    chat: Object,
  },
  data() {
    return {
      dialog: false,
      chatCopy: {},
    };
  },
  mounted() {
    console.log("setting chat:", this.chat._id);
    this.chatCopy = {
      name: this.chat.name,
      linkOrIdInBot: this.chat.linkOrIdInBot,
    };
  },
  methods: {
    chatChanges() {
      let changed = {};
      if (this.chat.name !== this.chatCopy.name) {
        changed.name = this.chatCopy.name;
      }
      if (this.chat.linkOrIdInBot !== this.chatCopy.linkOrIdInBot) {
        changed.linkOrIdInBot = this.chatCopy.linkOrIdInBot;
      }
      return changed;
    },
    async sendForm() {
      const checkForm = await this.$refs.chatSettingsForm.validate();
      if (!checkForm.valid) return;
      const changes = this.chatChanges();
      if (Object.keys(changes).length === 0) {
        this.closeDialog();
        return;
      }

      const response = await updateChat(this.chat._id, changes);
      this.$emit("eventchatupdated");
      this.closeDialog();
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
