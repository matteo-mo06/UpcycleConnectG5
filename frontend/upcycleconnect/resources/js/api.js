import axios from "axios";
import router from "./router.js";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        "Content-Type": "application/json",
    },
});

export function setAuthToken(token) {
    api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
}

export function clearAuthToken() {
    delete api.defaults.headers.common["Authorization"];
}

let loggingOut = false;

api.interceptors.response.use(
    (response) => {
        if (response.config?.url?.includes("/auth/login")) {
            loggingOut = false;
        }
        return response;
    },
    (error) => {
        const status = error.response?.status;
        const url = error.config?.url ?? "";

        const isAuthEndpoint = url.includes("/auth/login") || url.includes("/auth/register");

        if (!error.response) {
            router.push("/error/500");
        } else if (status === 401 && !isAuthEndpoint) {
            if (!loggingOut) {
                loggingOut = true;
                import("./stores/auth.js").then(({ useAuthStore }) => {
                    useAuthStore().logout();
                    router.push("/login");
                });
            }
        } else if (status >= 500) {
            router.push("/error/500");
        }

        return Promise.reject(error);
    }
);

export default api;
