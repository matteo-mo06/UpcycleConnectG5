import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api, { setAuthToken, clearAuthToken } from "../api";
import { useOneSignal } from "@onesignal/onesignal-vue3";

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
    const isSalarie = computed(
        () => user.value?.roles?.includes("salarie") ?? false,
    );

    // Accès aux espaces : piloté par les permissions access_admin / access_artisan
    // / access_salarie (voir Admin > Rôles), pas par le rôle brut.
    const canAccessAdmin = computed(
        () => user.value?.permissions?.includes("access_admin") ?? false,
    );
    const canAccessArtisan = computed(
        () => user.value?.permissions?.includes("access_artisan") ?? false,
    );
    const canAccessSalarie = computed(
        () => user.value?.permissions?.includes("access_salarie") ?? false,
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

        registerOnesignalPlayer();
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

    async function registerOnesignalPlayer() {
        try {
            const oneSignal = useOneSignal();
            oneSignal.User.PushSubscription.addEventListener("change", async (event) => {
                const id = event?.current?.id;
                if (id) {
                    await api.post("/user/onesignal-player-id", { player_id: id });
                }
            });
            await oneSignal.User.PushSubscription.optIn();
            const playerId = oneSignal.User.PushSubscription.id;
            if (playerId) {
                await api.post("/user/onesignal-player-id", { player_id: playerId });
            }
        } catch {
        }
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
        isSalarie,
        canAccessAdmin,
        canAccessArtisan,
        canAccessSalarie,
        hasPermission,
        login,
        logout,
        init,
        setTutorialDone,
    };
});
