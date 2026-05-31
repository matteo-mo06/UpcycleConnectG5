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
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Biographie</label
                            >
                            <textarea
                                v-model="profile.bio"
                                rows="3"
                                placeholder="Présentez-vous en quelques mots…"
                                class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
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
                            <input
                                v-model="passwords.current"
                                type="password"
                                class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                            />
                        </div>
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Nouveau mot de passe</label
                            >
                            <input
                                v-model="passwords.new"
                                type="password"
                                class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                            />
                        </div>
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-600 mb-1.5"
                                >Confirmer le nouveau mot de passe</label
                            >
                            <input
                                v-model="passwords.confirm"
                                type="password"
                                class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                            />
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

                <div class="bg-white rounded-2xl shadow-sm p-5">
                    <h2
                        class="font-semibold text-gray-800 mb-4"
                        style="font-family: var(--font-family-title)"
                    >
                        Notifications
                    </h2>
                    <div class="space-y-3">
                        <label
                            v-for="pref in notifPrefs"
                            :key="pref.key"
                            class="flex items-center justify-between cursor-pointer"
                        >
                            <span class="text-sm text-gray-700">{{
                                pref.label
                            }}</span>
                            <button
                                @click="pref.enabled = !pref.enabled"
                                :class="[
                                    'relative inline-flex h-5 w-9 items-center rounded-full transition-colors',
                                    pref.enabled
                                        ? 'bg-secondary'
                                        : 'bg-gray-200',
                                ]"
                            >
                                <span
                                    :class="[
                                        'inline-block h-3.5 w-3.5 rounded-full bg-white shadow transition-transform',
                                        pref.enabled
                                            ? 'translate-x-4.5'
                                            : 'translate-x-0.5',
                                    ]"
                                />
                            </button>
                        </label>
                    </div>
                    <button
                        @click="saveNotifPrefs"
                        class="mt-4 w-full py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Enregistrer les préférences
                    </button>
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
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import UserLayout from "@/Layouts/UserLayout.vue";
import { useAuthStore } from "@/stores/auth.js";
import api from "@/api.js";

const auth = useAuthStore();
const router = useRouter();

const profile = ref({
    first_name: auth.user?.first_name ?? "",
    last_name: auth.user?.last_name ?? "",
    email: auth.user?.email ?? "",
    bio: "",
    avatar_url: null,
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

// Les préférences de notification seront chargées depuis l'API
const notifPrefs = ref([
    { key: "deposit_update", label: "Mises à jour de dépôt", enabled: true },
    { key: "event_reminder", label: "Rappels d'événements", enabled: false },
    { key: "forum_reply", label: "Réponses au forum", enabled: true },
]);

async function saveProfile() {
    profileLoading.value = true;
    profileError.value = "";
    profileSuccess.value = "";
    try {
        // await api.patch('/user/profile', profile.value)
        profileSuccess.value = "Profil mis à jour avec succès.";
    } catch {
        profileError.value = "Une erreur est survenue.";
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
        // await api.patch('/user/password', { current: passwords.value.current, new: passwords.value.new })
        passwordSuccess.value = "Mot de passe mis à jour avec succès.";
        passwords.value = { current: "", new: "", confirm: "" };
    } catch {
        passwordError.value = "Mot de passe actuel incorrect.";
    }
    passwordLoading.value = false;
}

async function uploadAvatar(e) {
    const file = e.target.files[0];
    if (!file) return;
    // L'upload de l'avatar sera soumis via l'API
    // const form = new FormData()
    // form.append('file', file)
    // const { data } = await api.post('/upload', form)
    // await api.patch('/user/profile', { avatar_url: data.url })
    // profile.value.avatar_url = data.url
}

async function saveNotifPrefs() {
    // Les préférences de notification seront soumises via l'API
    // await api.patch('/user/notifications', Object.fromEntries(notifPrefs.value.map(p => [p.key, p.enabled])))
}

async function deleteAccount() {
    deleteLoading.value = true;
    try {
        // await api.delete('/user/account')
        auth.logout();
        router.push("/login");
    } catch {}
    deleteLoading.value = false;
}

// Les données du profil seront chargées depuis l'API
// onMounted(async () => {
//     const { data } = await api.get('/user/profile')
//     Object.assign(profile.value, data)
// })
</script>
