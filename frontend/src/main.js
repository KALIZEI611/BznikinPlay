import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "axios";

// Настройка базового URL для API
const API_URL =
  import.meta.env.VITE_API_URL || "https://your-backend-url.railway.app";
axios.defaults.baseURL = API_URL;

const app = createApp(App);
app.use(router);
app.mount("#app");
