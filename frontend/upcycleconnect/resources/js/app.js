import { createApp } from "vue";
import { createPinia } from "pinia";
import OneSignalVuePlugin from "@onesignal/onesignal-vue3";
import App from "./App.vue";
import router from "./router.js";
import i18n from "./i18n/index.js";
import "../css/app.css";
import { useAuthStore } from "./stores/auth";

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);
app.use(i18n);
app.use(OneSignalVuePlugin, {
    appId: import.meta.env.VITE_ONESIGNAL_APP_ID,
    allowLocalhostAsSecureOrigin: import.meta.env.DEV,
});

useAuthStore().init();

app.mount("#app");
