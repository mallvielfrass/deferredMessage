<template>
  <div>
    <div class="notification-container">
      <v-alert
        variant="tonal"
        closable
        border-color="error"
        color="error"
        close-label="Close Alert"
        v-for="(notification, index) in notifications"
        :key="index"
      >
        {{ notification.message }}
      </v-alert>
    </div>

    <!-- <button @click="showNotification">Показать уведомление</button> -->
  </div>
</template>

<script>
export default {
  data() {
    return {
      dialog: true,
      notifications: [],
    };
  },
  methods: {
    showNotification(text) {
      this.notifications.push({
        message: text, // текст уведомления
      });
    },
  },
};
</script>

<style>
.notification-container {
  top: 40px;
  right: 10px;
}
</style>
