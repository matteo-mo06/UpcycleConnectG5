import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api, { setAuthToken, clearAuthToken } from "../api";

export const useAuthStore = defineStore("auth", () => {
    const token = ref(sessionStorage.getItem("token") ?? null);
    const user = ref(JSON.parse(sessionStorage.getItem("user") ?? "null"));

    const isLoggedIn = computed(() => !!token.value);
    const isAdmin = computed(
        () => user.value?.roles?.includes("admin") ?? false,
    );
    const isArtisan = computed(
        () => user.value?.roles?.includes("artisan") ?? false,
    );

    async function login(email, password) {
        const { data } = await api.post("/auth/login", {
            email,
            password_user: password,
        });

        token.value = data.token;
        user.value = data.user;

        sessionStorage.setItem("token", data.token);
        sessionStorage.setItem("user", JSON.stringify(data.user));

        setAuthToken(data.token);
    }

    function logout() {
        token.value = null;
        user.value = null;
        sessionStorage.removeItem("token");
        sessionStorage.removeItem("user");
        clearAuthToken();
    }

    function init() {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        if (token.value) {
            setAuthToken(token.value);
        }
    }

    function hasPermission(permission) {
        return user.value?.permissions?.includes(permission) ?? false;
    }

    function setTutorialDone() {
        if (user.value) {
            user.value.tutorial_done = true;
            sessionStorage.setItem("user", JSON.stringify(user.value));
        }
    }

    return {
        token,
        user,
        isLoggedIn,
        isAdmin,
        isArtisan,
        hasPermission,
        login,
        logout,
        init,
        setTutorialDone,
    };
});
