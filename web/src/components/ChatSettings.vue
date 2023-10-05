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
                  v-model="chat.name"
                  label="Title"
                  required
                  :rules="[(v) => !!v || 'Title is required']"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-text-field
                  v-model="chat.linkOrIdInNetwork"
                  label="Link or ID in the network"
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
import { getNetworks, getChats } from "@/api/networks.js";
import { createChat } from "@/api/chats.js";
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
      linkOrIdInNetwork: this.chat.linkOrIdInNetwork,
    };
  },
  methods: {
    chatIsChanged() {
      return (
        this.chat.name !== this.chatCopy.name ||
        this.chat.linkOrIdInNetwork !== this.chatCopy.linkOrIdInNetwork
      );
    },
    async sendForm() {
      const checkForm = await this.$refs.chatSettingsForm.validate();
      if (!checkForm.valid) return;
      if (!this.chatIsChanged()) {
        this.closeDialog();
        return;
      }

      //   const response = await createChat({
      // 	name: this.chatCopy.name,
      // 	linkOrIdInNetwork: this.chatCopy.linkOrIdInNetwork,
      // 	networkIdentifier: this.chat.networkIdentifier,
      //   });
      //   console.log(response);
      //   this.$emit("eventchatcreated");
      //   this.closeDialog();
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
