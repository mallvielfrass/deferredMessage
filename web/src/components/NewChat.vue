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
          Create New Chat</v-btn
        >
      </v-col>
    </v-row>
  </v-container>

  <v-dialog v-model="dialog" max-width="500px" max-height="70%">
    <v-card>
      <div>
        <v-form v-model="valid" @submit.prevent="" ref="newChatForm">
          <v-container>
            <v-row justify="center" align="center">
              <v-card-title>New Chat</v-card-title>
            </v-row>
            <v-row justify="center" align="center">
              <div v-if="errors.length">
                <b>Please correct the following error(s):</b>
                <ul>
                  <li v-for="error in errors">{{ error }}</li>
                </ul>
              </div>
            </v-row>
            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-text-field
                  v-model="title"
                  label="chat name"
                  required
                  :rules="nameRules"
                ></v-text-field>
              </v-col>
            </v-row>

            <v-row justify="center" align="center">
              <v-col cols="12" md="7">
                <v-select
                  ref="socialNetworkSelect"
                  label="Select network"
                  v-model="socialNetwork"
                  :required="!selected"
                  :rules="socialNetworkRules"
                  @update:modelValue="selectSocialNetwork"
                  :items="socialNetworkList"
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
import { getNetworks, getChats } from "@/api/networks.js";
import { createChat } from "@/api/chats.js";
export default {
  data() {
    return {
      dialog: false,
      socialNetwork: {},
      socialNetworkList: [],
      showChats: false,
      selected: null,

      errors: [],
      title: "",

      nameRules: [
        (value) => {
          if (value) return true;
          return "Name is required.";
        },
        (value) => {
          if (4 <= value?.length <= 10) return true;
          return "Name must be greater then 4 and less than 10 characters.";
        },
      ],
      socialNetworkRules: [
        (v) => {
          console.log("V:", v);
          if (!v.identifier) {
            console.log("v req");
            return "social network is required";
          }
          return;
        },
      ],
    };
  },
  mounted() {
    this.getSocialNetworkList();
  },
  methods: {
    async sendForm() {
      this.errors = [];
      const checkForm = await this.$refs.newChatForm.validate();
      // console.log("sendForm:", checkForm);
      if (!checkForm.valid) return;
      const resp = await createChat(this.title, this.socialNetwork.identifier);
      if (resp.error) {
        console.log("error", resp.error);
        this.errors.push(resp.error);
        return;
      }
      this.dialog = false;
      this.$emit("eventchatcreated");
      //  this.$refs.socialNetworkSelect.validate();
    },

    async getSocialNetworkList() {
      this.socialNetworkList = await getNetworks();
    },
    itemProps(item) {
      return {
        title: item.name,
        subtitle: item.identifier,
      };
    },
    async selectSocialNetwork(network) {
      console.log("Network selected:", network);
      //  await this.getChatsList();
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
