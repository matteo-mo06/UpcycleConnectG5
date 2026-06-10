<template>
    <UserLayout>
        <div class="mb-6">
            <h1
                class="text-3xl font-bold text-gray-800"
                style="font-family: var(--font-family-title)"
            >
                Paramètres
            </h1>
            <p class="text-sm text-gray-400 mt-1">
                Gérez votre profil et vos préférences
            </p>
        </div>

        <div class="grid grid-cols-3 gap-6">
            <div class="col-span-2 space-y-5">
                <div class="bg-white rounded-2xl shadow-sm p-5">
                    <h2
                        class="font-semibold text-gray-800 mb-4"
                        style="font-family: var(--font-family-title)"
                    >
                        Informations personnelles
                    </h2>
                    <form @submit.prevent="saveProfile" class="space-y-4">
                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <label
                                    class="block text-xs font-medium text-gray-600 mb-1.5"
                                    >Prénom</label
                                >
                                <input
                                    v-model="profile.first_name"
                                    type="text"
                                    class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                                />
                            </div>
                            <div>
                                <label
                                    class="block text-xs font-medium text-gray-600 mb-1.5"
                                    >Nom</label
                                >
                                <input
                                    v-model="profile.last_name"
                                    type="text"
                                    class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                                />
                            </div>
                        </div>
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Adresse e-mail</label
                            >
                            <input
                                v-model="profile.email"
                                type="email"
                                class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                            />
                        </div>
                        <p v-if="profileError" class="text-xs text-red-500">
                            {{ profileError }}
                        </p>
                        <p v-if="profileSuccess" class="text-xs text-secondary">
                            {{ profileSuccess }}
                        </p>
                        <div class="flex justify-end">
                            <button
                                type="submit"
                                :disabled="profileLoading"
                                class="px-5 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                            >
                                {{
                                    profileLoading
                                        ? "Enregistrement…"
                                        : "Enregistrer"
                                }}
                            </button>
                        </div>
                    </form>
                </div>

                <div class="bg-white rounded-2xl shadow-sm p-5">
                    <h2
                        class="font-semibold text-gray-800 mb-4"
                        style="font-family: var(--font-family-title)"
                    >
                        Mot de passe
                    </h2>
                    <form @submit.prevent="savePassword" class="space-y-4">
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Mot de passe actuel</label
                            >
                            <div class="relative">
                                <input
                                    v-model="passwords.current"
                                    :type="showCurrent ? 'text' : 'password'"
                                    class="w-full px-3 py-2 pr-10 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                                />
                                <button
                                    type="button"
                                    @click="showCurrent = !showCurrent"
                                    class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                                >
                                    <svg v-if="!showCurrent" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                    </svg>
                                    <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                                    </svg>
                                </button>
                            </div>
                        </div>
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Nouveau mot de passe</label
                            >
                            <div class="relative">
                                <input
                                    v-model="passwords.new"
                                    :type="showNew ? 'text' : 'password'"
                                    class="w-full px-3 py-2 pr-10 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                                />
                                <button
                                    type="button"
                                    @click="showNew = !showNew"
                                    class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                                >
                                    <svg v-if="!showNew" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                    </svg>
                                    <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                                    </svg>
                                </button>
                            </div>
                        </div>
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Confirmer le nouveau mot de passe</label
                            >
                            <div class="relative">
                                <input
                                    v-model="passwords.confirm"
                                    :type="showConfirm ? 'text' : 'password'"
                                    class="w-full px-3 py-2 pr-10 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                                />
                                <button
                                    type="button"
                                    @click="showConfirm = !showConfirm"
                                    class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                                >
                                    <svg v-if="!showConfirm" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                    </svg>
                                    <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                                    </svg>
                                </button>
                            </div>
                        </div>
                        <p v-if="passwordError" class="text-xs text-red-500">
                            {{ passwordError }}
                        </p>
                        <p
                            v-if="passwordSuccess"
                            class="text-xs text-secondary"
                        >
                            {{ passwordSuccess }}
                        </p>
                        <div class="flex justify-end">
                            <button
                                type="submit"
                                :disabled="passwordLoading"
                                class="px-5 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                            >
                                {{
                                    passwordLoading
                                        ? "Mise à jour…"
                                        : "Mettre à jour"
                                }}
                            </button>
                        </div>
                    </form>
                </div>

                <div v-if="!isArtisan" class="bg-white rounded-2xl shadow-sm p-5">
                    <h2
                        class="font-semibold text-gray-800 mb-1"
                        style="font-family: var(--font-family-title)"
                    >
                        Compte professionnel
                    </h2>
                    <p class="text-xs text-gray-400 mb-4">
                        Devenez artisan pour publier des annonces et proposer vos services.
                    </p>

                    <template v-if="!proRequest">
                        <button
                            @click="submitProRequest"
                            :disabled="proSubmitting"
                            class="px-4 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                        >
                            {{ proSubmitting ? "Envoi en cours…" : "Devenir professionnel" }}
                        </button>
                    </template>

                    <template v-else-if="proRequest.status === 'pending'">
                        <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-amber-100 text-amber-700">
                            <span class="w-1.5 h-1.5 rounded-full bg-amber-500" />
                            En cours de traitement
                        </span>
                        <p class="text-xs text-gray-400 mt-2">
                            Demande envoyée le {{ proRequest.created_at?.slice(0, 10) }}
                        </p>
                    </template>

                    <template v-else-if="proRequest.status === 'rejected'">
                        <p class="text-sm text-red-600 mb-1">Votre demande a été refusée.</p>
                        <p class="text-xs text-gray-400 mb-3">Vous pouvez soumettre une nouvelle demande.</p>
                        <button
                            @click="submitProRequest"
                            :disabled="proSubmitting"
                            class="px-4 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                        >
                            {{ proSubmitting ? "Envoi en cours…" : "Renvoyer une demande" }}
                        </button>
                    </template>
                </div>

                <div class="bg-white rounded-2xl shadow-sm p-5">
                    <h2
                        class="font-semibold text-gray-800 mb-1"
                        style="font-family: var(--font-family-title)"
                    >
                        Notifications
                    </h2>
                    <p class="text-xs text-gray-400 mb-4">
                        Gérez vos notifications push sur cet appareil.
                    </p>
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-700 font-medium">Notifications push</p>
                            <p class="text-xs text-gray-400 mt-0.5">
                                {{ notifOptedIn ? "Activées - vous recevez les alertes en temps réel." : "Désactivées." }}
                            </p>
                        </div>
                        <button
                            @click="toggleNotifications"
                            :disabled="notifLoading"
                            :class="[
                                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none disabled:opacity-50',
                                notifOptedIn ? 'bg-primary' : 'bg-gray-200'
                            ]"
                        >
                            <span
                                :class="[
                                    'inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform',
                                    notifOptedIn ? 'translate-x-6' : 'translate-x-1'
                                ]"
                            />
                        </button>
                    </div>
                    <p v-if="notifBlocked" class="text-xs text-amber-600 mt-3">
                        Les notifications sont bloquées dans votre navigateur. Autorisez-les dans les paramètres du site pour les réactiver.
                    </p>
                </div>

                <div
                    class="bg-white rounded-2xl shadow-sm p-5 border border-red-100"
                >
                    <h2
                        class="font-semibold text-red-600 mb-1"
                        style="font-family: var(--font-family-title)"
                    >
                        Zone de danger
                    </h2>
                    <p class="text-xs text-gray-400 mb-4">
                        Ces actions sont irréversibles, soyez prudent.
                    </p>
                    <button
                        @click="showDeleteConfirm = true"
                        class="px-4 py-2 rounded-xl border border-red-300 text-red-600 text-sm font-medium hover:bg-red-50 transition-colors"
                    >
                        Supprimer mon compte
                    </button>
                </div>
            </div>

            <div class="space-y-5">
                <div
                    class="bg-white rounded-2xl shadow-sm p-5 flex flex-col items-center gap-4"
                >
                    <h2
                        class="font-semibold text-gray-800 self-start"
                        style="font-family: var(--font-family-title)"
                    >
                        Photo de profil
                    </h2>
                    <div
                        class="w-24 h-24 rounded-full bg-primary/10 flex items-center justify-center overflow-hidden"
                    >
                        <img
                            v-if="profile.avatar_url"
                            :src="profile.avatar_url"
                            alt="Avatar"
                            class="w-full h-full object-cover"
                        />
                        <span v-else class="text-3xl font-bold text-primary">
                            {{
                                (profile.first_name?.[0] ?? "") +
                                (profile.last_name?.[0] ?? "")
                            }}
                        </span>
                    </div>
                    <label
                        class="px-4 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors cursor-pointer"
                    >
                        Changer la photo
                        <input
                            type="file"
                            accept="image/*"
                            class="hidden"
                            @change="uploadAvatar"
                        />
                    </label>
                </div>

            </div>
        </div>

        <div
            v-if="showDeleteConfirm"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="showDeleteConfirm = false"
        >
            <div
                class="bg-white rounded-2xl shadow-xl w-full max-w-sm overflow-hidden"
            >
                <div class="px-6 py-5">
                    <h3
                        class="font-semibold text-gray-800 mb-2"
                        style="font-family: var(--font-family-title)"
                    >
                        Supprimer mon compte
                    </h3>
                    <p class="text-sm text-gray-500">
                        Cette action est irréversible. Toutes vos données seront
                        supprimées définitivement.
                    </p>
                </div>
                <div class="px-6 pb-5 flex gap-3">
                    <button
                        @click="showDeleteConfirm = false"
                        class="flex-1 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Annuler
                    </button>
                    <button
                        @click="deleteAccount"
                        :disabled="deleteLoading"
                        class="flex-1 py-2 rounded-xl bg-red-500 text-white text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-50"
                    >
                        {{ deleteLoading ? "Suppression…" : "Confirmer" }}
                    </button>
                </div>
            </div>
        </div>
    </UserLayout>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import UserLayout from "@/Layouts/UserLayout.vue";
import { useAuthStore } from "@/stores/auth.js";
import { useOneSignal } from "@onesignal/onesignal-vue3";
import api from "@/api.js";

const oneSignal = useOneSignal();

const auth = useAuthStore();
const router = useRouter();

const profile = ref({
    first_name: auth.user?.first_name ?? "",
    last_name: auth.user?.last_name ?? "",
    email: auth.user?.email ?? "",
    avatar_url: auth.user?.avatar_url ?? null,
});

const passwords = ref({ current: "", new: "", confirm: "" });

const profileLoading = ref(false);
const profileError = ref("");
const profileSuccess = ref("");

const passwordLoading = ref(false);
const passwordError = ref("");
const passwordSuccess = ref("");

const deleteLoading = ref(false);
const showDeleteConfirm = ref(false);

const showCurrent = ref(false);
const showNew = ref(false);
const showConfirm = ref(false);

const proRequest = ref(null);
const proSubmitting = ref(false);
const isArtisan = computed(() => auth.user?.roles?.includes("artisan") ?? false);

const notifOptedIn = ref(false);
const notifLoading = ref(false);
const notifBlocked = ref(false);

function refreshNotifState() {
    notifOptedIn.value = oneSignal.User.PushSubscription.optedIn ?? false;
    notifBlocked.value = Notification.permission === "denied";
}

oneSignal.User.PushSubscription.addEventListener("change", refreshNotifState);

async function toggleNotifications() {
    if (notifBlocked.value) return;
    notifLoading.value = true;
    try {
        if (notifOptedIn.value) {
            await oneSignal.User.PushSubscription.optOut();
            await api.delete("/user/onesignal-player-id");
        } else {
            await oneSignal.User.PushSubscription.optIn();
            const playerId = oneSignal.User.PushSubscription.id;
            if (playerId) await api.post("/user/onesignal-player-id", { player_id: playerId });
        }
        refreshNotifState();
    } catch (e) {
        console.error("Erreur notifications:", e);
    }
    notifLoading.value = false;
}

async function saveProfile() {
    profileLoading.value = true;
    profileError.value = "";
    profileSuccess.value = "";
    try {
        await api.patch("/user/profile", {
            first_name: profile.value.first_name,
            last_name: profile.value.last_name,
            email: profile.value.email,
        });
        const updated = { ...auth.user, first_name: profile.value.first_name, last_name: profile.value.last_name, email: profile.value.email };
        auth.user = updated;
        sessionStorage.setItem("user", JSON.stringify(updated));
        profileSuccess.value = "Profil mis à jour avec succès.";
    } catch (e) {
        profileError.value = e.response?.data?.error ?? "Une erreur est survenue.";
    }
    profileLoading.value = false;
}

async function savePassword() {
    passwordError.value = "";
    passwordSuccess.value = "";
    if (passwords.value.new !== passwords.value.confirm) {
        passwordError.value = "Les mots de passe ne correspondent pas.";
        return;
    }
    passwordLoading.value = true;
    try {
        await api.patch("/user/password", {
            current: passwords.value.current,
            new: passwords.value.new,
        });
        passwordSuccess.value = "Mot de passe mis à jour avec succès.";
        passwords.value = { current: "", new: "", confirm: "" };
    } catch (e) {
        passwordError.value = e.response?.data?.error ?? "Mot de passe actuel incorrect.";
    }
    passwordLoading.value = false;
}

async function uploadAvatar(e) {
    const file = e.target.files[0];
    if (!file) return;
    const form = new FormData();
    form.append("file", file);
    try {
        const { data: upload } = await api.post("/upload", form, { headers: { "Content-Type": undefined } });
        await api.patch("/user/avatar", { avatar_url: upload.url });
        profile.value.avatar_url = upload.url;
        const updated = { ...auth.user, avatar_url: upload.url };
        auth.user = updated;
        sessionStorage.setItem("user", JSON.stringify(updated));
    } catch {}
}

async function submitProRequest() {
    proSubmitting.value = true;
    try {
        await api.post("/professional-request");
        const { data } = await api.get("/user/professional-request");
        proRequest.value = data;
    } catch (e) {
        alert(e.response?.data?.error ?? "Erreur lors de l'envoi de la demande");
    } finally {
        proSubmitting.value = false;
    }
}

async function deleteAccount() {
    deleteLoading.value = true;
    try {
        await api.delete("/user/account");
        auth.logout();
        router.push("/login");
    } catch {}
    deleteLoading.value = false;
}

onMounted(async () => {
    setTimeout(refreshNotifState, 500);

    try {
        const { data } = await api.get("/user/me");
        profile.value.first_name = data.first_name;
        profile.value.last_name = data.last_name;
        profile.value.email = data.email;
        profile.value.avatar_url = data.avatar_url ?? null;
    } catch {}
    try {
        const { data } = await api.get("/user/professional-request");
        proRequest.value = data;
    } catch {}
});
</script>
