/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Import Bootstrap and BootstrapVue CSS files (order is important)

// Components
import App from "./App.vue";
import { createApp } from "vue";
// Plugins
import { registerPlugins } from "@/plugins";

const app = createApp(App);

registerPlugins(app);

app.mount("#app");
