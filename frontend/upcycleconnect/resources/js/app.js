import { createApp } from "vue";
import { createPinia } from "pinia";
import OneSignalVuePlugin from "@onesignal/onesignal-vue3";
import App from "./App.vue";
import router from "./router.js";
import "../css/app.css";
import { useAuthStore } from "./stores/auth";

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);
app.use(OneSignalVuePlugin, {
    appId: import.meta.env.VITE_ONESIGNAL_APP_ID,
    allowLocalhostAsSecureOrigin: true,
});

useAuthStore().init();

app.mount("#app");
