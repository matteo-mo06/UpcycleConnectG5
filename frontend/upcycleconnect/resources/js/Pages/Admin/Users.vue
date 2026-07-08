<template>
    <AdminLayout>
        <div class="mb-6">
            <h1
                class="text-3xl font-bold text-gray-800"
                style="font-family: var(--font-family-title)"
            >
                Gestion des utilisateurs
            </h1>
        </div>

        <button
            @click="notifyRecent"
            class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
        >
            <svg
                class="w-4 h-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 4v16m8-8H4"
                />
            </svg>
            Envoie de mail au 10 derniers users valides
        </button>

        <div
            class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 items-center"
        >
            <input
                v-model="search"
                type="text"
                placeholder="Rechercher par nom, email..."
                class="flex-1 px-4 py-2 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
            <select
                v-model="filterType"
                class="px-4 py-2 rounded-lg border border-gray-200 text-sm text-gray-600 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            >
                <option value="">Tous les types</option>
                <option
                    v-for="role in roles"
                    :key="role.value"
                    :value="role.value"
                >
                    {{ role.label }}
                </option>
            </select>
            <select
                v-model="filterStatus"
                class="px-4 py-2 rounded-lg border border-gray-200 text-sm text-gray-600 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            >
                <option value="">Tous les statuts</option>
                <option value="active">Actif</option>
                <option value="suspended">Suspendu</option>
                <option value="blacklisted">Blacklisté</option>
            </select>
            <span class="text-xs text-gray-400 ml-auto"
                >{{ total }} utilisateur(s)</span
            >
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">
            Chargement...
        </div>

        <div v-if="error" class="py-12 text-center text-sm text-red-500">
            {{ error }}
        </div>

        <div
            v-if="!loading && !error"
            class="bg-white rounded-xl shadow-sm overflow-hidden"
        >
            <div class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead>
                        <tr class="bg-primary">
                            <th
                                class="text-left text-white font-medium px-5 py-3"
                            >
                                Utilisateur
                            </th>
                            <th
                                class="text-left text-white font-medium px-5 py-3"
                            >
                                Email
                            </th>
                            <th
                                class="text-left text-white font-medium px-5 py-3"
                            >
                                Type
                            </th>
                            <th
                                class="text-left text-white font-medium px-5 py-3"
                            >
                                Statut
                            </th>
                            <th
                                class="text-left text-white font-medium px-5 py-3"
                            >
                                Inscription
                            </th>
                            <th
                                class="text-right text-white font-medium px-5 py-3"
                            >
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="(user, i) in users"
                            :key="user.id"
                            :class="[
                                'border-b border-gray-50',
                                i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40',
                            ]"
                        >
                            <td class="px-5 py-3">
                                <div class="flex items-center gap-2.5">
                                    <div
                                        class="w-7 h-7 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0"
                                    >
                                        <svg
                                            class="w-4 h-4 text-gray-400"
                                            fill="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"
                                            />
                                        </svg>
                                    </div>
                                    <div class="min-w-0">
                                        <p
                                            class="font-medium text-gray-800 truncate"
                                        >
                                            {{ user.name }}
                                        </p>
                                        <p class="text-xs text-gray-400">
                                            #{{ user.id }}
                                        </p>
                                    </div>
                                </div>
                            </td>

                            <td class="px-5 py-3 text-gray-500">
                                {{ user.email }}
                            </td>

                            <td class="px-5 py-3">
                                <span
                                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600"
                                >
                                    {{ user.type }}
                                </span>
                            </td>

                            <td class="px-5 py-3">
                                <span
                                    class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium"
                                    :class="statusBadge(user.status)"
                                >
                                    <span
                                        class="w-1.5 h-1.5 rounded-full"
                                        :class="statusDot(user.status)"
                                    />
                                    {{ user.status }}
                                </span>
                            </td>

                            <td class="px-5 py-3 text-gray-500">
                                {{ user.date }}
                            </td>

                            <td class="px-5 py-3">
                                <div
                                    class="flex items-center justify-end gap-1"
                                >
                                    <button
                                        @click="openRights(user)"
                                        class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors duration-150"
                                        title="Gérer les droits"
                                    >
                                        <svg
                                            class="w-4 h-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                                            />
                                        </svg>
                                    </button>

                                    <button
                                        @click="openHistory(user)"
                                        class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors duration-150"
                                        title="Voir l'historique"
                                    >
                                        <svg
                                            class="w-4 h-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                    </button>

                                    <button
                                        @click="deleteUser(user)"
                                        class="p-1.5 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors duration-150"
                                        title="Supprimer le compte"
                                    >
                                        <svg
                                            class="w-4 h-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                            />
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>

                        <tr v-if="users.length === 0">
                            <td
                                colspan="6"
                                class="px-5 py-12 text-center text-gray-400 text-sm"
                            >
                                Aucun utilisateur ne correspond à votre
                                recherche.
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div
                class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400"
            >
                {{ total }} utilisateur(s) au total
            </div>
        </div>
        <Pagination
            v-if="total > 50"
            :page="page"
            :total="total"
            :limit="50"
            @update:page="changePage"
        />

        <div
            v-if="toDelete"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
        >
            <div
                class="absolute inset-0 bg-black/40"
                @click="toDelete = null"
            />
            <div
                class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6"
            >
                <h3 class="font-semibold text-gray-800 mb-2">
                    Supprimer le compte ?
                </h3>
                <p class="text-sm text-gray-500 mb-5">
                    Le compte de
                    <span class="font-medium text-gray-700">{{
                        toDelete.name
                    }}</span>
                    sera définitivement supprimé.
                </p>
                <div class="flex gap-3">
                    <button
                        @click="toDelete = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors duration-150"
                    >
                        Annuler
                    </button>
                    <button
                        @click="confirmDeleteUser"
                        :disabled="deleting"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors duration-150 disabled:opacity-60"
                    >
                        {{ deleting ? "Suppression…" : "Supprimer" }}
                    </button>
                </div>
            </div>
        </div>

        <div
            v-if="rightsModal.open"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="rightsModal.open = false"
        >
            <div
                class="absolute inset-0 bg-black/40"
                @click="rightsModal.open = false"
            />
            <div
                class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6 z-10"
            >
                <h3 class="text-base font-semibold text-gray-800 mb-5">
                    Rôle :
                </h3>
                <div class="space-y-2">
                    <label
                        v-for="role in roles"
                        :key="role.value"
                        class="flex items-center gap-3 p-3 rounded-lg border transition-colors duration-150"
                        :class="[
                            role.value === 'admin' &&
                            rightsModal.user?.id === auth.user?.id
                                ? 'opacity-50 cursor-not-allowed'
                                : 'cursor-pointer',
                            role.value === 'admin'
                                ? rightsModal.selectedRole === 'admin'
                                    ? 'border-red-300 bg-red-50'
                                    : 'border-red-100 hover:border-red-200'
                                : rightsModal.selectedRole === role.value
                                  ? 'border-primary bg-primary/5'
                                  : 'border-gray-100 hover:border-gray-200',
                        ]"
                    >
                        <input
                            type="radio"
                            :value="role.value"
                            v-model="rightsModal.selectedRole"
                            :disabled="
                                role.value === 'admin' &&
                                rightsModal.user?.id === auth.user?.id
                            "
                            class="accent-primary"
                        />
                        <div class="flex items-center gap-2 flex-1">
                            <p class="text-sm font-medium text-gray-800">
                                {{ role.label }}
                            </p>
                            <span
                                v-if="role.value === 'admin'"
                                class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs font-medium bg-red-100 text-red-700"
                            >
                                <svg
                                    class="w-3 h-3"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                                    />
                                </svg>
                                Admin
                            </span>
                        </div>
                    </label>
                </div>
                <div class="flex gap-2 mt-5">
                    <button
                        @click="rightsModal.open = false"
                        class="flex-1 py-2 rounded-lg text-sm text-gray-500 border border-gray-200 hover:bg-gray-50 transition-colors duration-150"
                    >
                        Annuler
                    </button>
                    <button
                        @click="saveRights"
                        class="flex-1 py-2 rounded-lg text-sm font-medium bg-primary text-white hover:bg-primary-dark transition-colors duration-150"
                    >
                        Enregistrer
                    </button>
                </div>
            </div>
        </div>
        <div
            v-if="adminConfirm.open"
            class="fixed inset-0 z-[60] flex items-center justify-center p-4"
        >
            <div
                class="absolute inset-0 bg-black/40"
                @click="adminConfirm.open = false"
            />
            <div
                class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6 z-10"
            >
                <h3 class="font-semibold text-gray-800 mb-2">
                    Attribuer le rôle administrateur ?
                </h3>
                <p class="text-sm text-gray-500 mb-5">
                    Vous allez attribuer le rôle administrateur à cet
                    utilisateur. Cette action lui donnera un accès complet à la
                    plateforme.
                </p>
                <div class="flex gap-3">
                    <button
                        @click="adminConfirm.open = false"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors duration-150"
                    >
                        Annuler
                    </button>
                    <button
                        @click="doSaveRights"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors duration-150"
                    >
                        Confirmer
                    </button>
                </div>
            </div>
        </div>

        <div
            v-if="historyModal.open"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="historyModal.open = false"
        >
            <div
                class="absolute inset-0 bg-black/40"
                @click="historyModal.open = false"
            />
            <div
                class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden flex flex-col max-h-[90vh]"
            >
                <div
                    class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0"
                >
                    <h3
                        class="font-semibold text-gray-800"
                        style="font-family: var(--font-family-title)"
                    >
                        Historique : {{ historyModal.user?.name }}
                    </h3>
                    <button
                        @click="historyModal.open = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
                    >
                        <svg
                            class="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 overflow-y-auto space-y-5 flex-1">
                    <div
                        v-if="historyModal.loading"
                        class="text-sm text-gray-400 text-center py-6"
                    >
                        Chargement…
                    </div>

                    <template v-else>
                        <div>
                            <p
                                class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2"
                            >
                                Sanctions reçues
                            </p>
                            <p
                                v-if="historyModal.data.sanctions.length === 0"
                                class="text-xs text-gray-400"
                            >
                                Aucune sanction.
                            </p>
                            <div v-else class="space-y-2">
                                <div
                                    v-for="s in historyModal.data.sanctions"
                                    :key="s.id_sanction"
                                    class="p-3 rounded-lg bg-gray-50 border border-gray-100"
                                >
                                    <div
                                        class="flex items-center justify-between"
                                    >
                                        <span
                                            class="text-sm font-medium text-gray-800 capitalize"
                                            >{{ s.type }}</span
                                        >
                                        <span class="text-xs text-gray-400">{{
                                            s.created_at?.slice(0, 10)
                                        }}</span>
                                    </div>
                                    <p
                                        v-if="s.reason"
                                        class="text-xs text-gray-500 mt-1"
                                    >
                                        {{ s.reason }}
                                    </p>
                                    <p class="text-xs text-gray-400 mt-1">
                                        Par {{ s.admin_name || "-" }}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div>
                            <p
                                class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2"
                            >
                                Demandes de compte pro
                            </p>
                            <p
                                v-if="
                                    historyModal.data.professional_requests
                                        .length === 0
                                "
                                class="text-xs text-gray-400"
                            >
                                Aucune demande.
                            </p>
                            <div v-else class="space-y-2">
                                <div
                                    v-for="req in historyModal.data
                                        .professional_requests"
                                    :key="req.id"
                                    class="p-3 rounded-lg bg-gray-50 border border-gray-100 flex items-center justify-between"
                                >
                                    <span
                                        class="text-sm text-gray-700 capitalize"
                                        >{{ req.status }}</span
                                    >
                                    <span class="text-xs text-gray-400">{{
                                        req.created_at?.slice(0, 10)
                                    }}</span>
                                </div>
                            </div>
                        </div>

                        <div>
                            <p
                                class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2"
                            >
                                Signalements reçus
                            </p>
                            <p
                                v-if="
                                    historyModal.data.reports_received
                                        .length === 0
                                "
                                class="text-xs text-gray-400"
                            >
                                Aucun signalement.
                            </p>
                            <div v-else class="space-y-2">
                                <div
                                    v-for="r in historyModal.data
                                        .reports_received"
                                    :key="r.id_report"
                                    class="p-3 rounded-lg bg-gray-50 border border-gray-100"
                                >
                                    <div
                                        class="flex items-center justify-between"
                                    >
                                        <span
                                            class="text-sm font-medium text-gray-800"
                                            >{{
                                                r.content_title ||
                                                r.content_type
                                            }}</span
                                        >
                                        <span class="text-xs text-gray-400">{{
                                            r.created_at?.slice(0, 10)
                                        }}</span>
                                    </div>
                                    <p class="text-xs text-gray-500 mt-1">
                                        {{ r.reason }}
                                    </p>
                                    <p class="text-xs text-gray-400 mt-1">
                                        Statut : {{ r.status }}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </template>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import api from "@/api.js";
import AdminLayout from "@/Layouts/AdminLayout.vue";
import Pagination from "@/Components/Pagination.vue";
import { roleLabel, fullName } from "@/utils.js";
import { useAuthStore } from "@/stores/auth.js";

const search = ref("");
const filterType = ref("");
const filterStatus = ref("");

const users = ref([]);
const total = ref(0);
const page = ref(1);
const loading = ref(false);
const error = ref("");
const roleDisplayMap = ref({});
const toDelete = ref(null);
const deleting = ref(false);

const rightsModal = ref({ open: false, user: null, selectedRole: "" });
const roles = ref([]);

async function notifyRecent() {
    const subject = prompt("Sujet du mail ?");
    if (!subject) return;
    const message = prompt("Message ?");
    if (!message) return;
    try {
        const { data } = await api.post("/admin/users/notify-recent", {
            subject,
            message,
        });
        alert(data.message);
    } catch (e) {
        alert(e.response?.data?.error ?? "Erreur lors de l'envoi.");
    }
}
const auth = useAuthStore();
const adminConfirm = ref({ open: false });
const historyModal = ref({
    open: false,
    user: null,
    loading: false,
    data: { sanctions: [], reports_received: [], professional_requests: [] },
});

function changePage(p) {
    page.value = p;
}

async function fetchUsers() {
    loading.value = true;
    try {
        const params = { page: page.value, limit: 50 };
        if (search.value) params.search = search.value;
        if (filterStatus.value) params.status = filterStatus.value;
        if (filterType.value) params.role = filterType.value;
        const { data } = await api.get("/admin/users", { params });
        users.value = data.data.map((u) => ({
            id: u.id,
            email: u.email,
            name: fullName(u),
            type: roleToType(u.roles),
            status: statusToLabel(u.status),
            date: u.created_at?.slice(0, 10) ?? "-",
        }));
        total.value = data.total;
    } catch (e) {
        error.value = "Impossible de charger les utilisateurs";
        console.error(e);
    } finally {
        loading.value = false;
    }
}

let debounce = null;
watch(search, () => {
    clearTimeout(debounce);
    debounce = setTimeout(() => {
        page.value = 1;
        fetchUsers();
    }, 300);
});
watch([filterType, filterStatus], () => {
    page.value = 1;
    fetchUsers();
});
watch(page, fetchUsers);

onMounted(async () => {
    try {
        const { data: rolesData } = await api.get("/admin/roles");
        roleDisplayMap.value = Object.fromEntries(
            rolesData.map((r) => [r.name, formatRoleName(r.name)]),
        );
        roles.value = rolesData.map((r) => ({
            id: r.id,
            value: r.name,
            label: formatRoleName(r.name),
        }));

        await fetchUsers();
    } catch (e) {
        error.value = "Impossible de charger les données";
        console.error(e);
    }
});

function formatRoleName(name) {
    return roleLabel(name);
}

function roleToType(userRoles) {
    if (!Array.isArray(userRoles) || userRoles.length === 0)
        return "Particulier";
    for (const role of userRoles) {
        if (roleDisplayMap.value[role]) return roleDisplayMap.value[role];
    }
    return "Particulier";
}

function statusToLabel(status) {
    const map = {
        active: "Actif",
        suspended: "Suspendu",
        blacklisted: "Blacklisté",
    };
    return map[status] ?? "Actif";
}

function openRights(user) {
    rightsModal.value = { open: true, user, selectedRole: user.type };
}

async function saveRights() {
    if (rightsModal.value.selectedRole === "admin") {
        adminConfirm.value.open = true;
        return;
    }
    await doSaveRights();
}

async function doSaveRights() {
    const user = rightsModal.value.user;
    const newRole = rightsModal.value.selectedRole;

    const roleObj = roles.value.find((r) => r.value === newRole);
    if (!roleObj) return;

    try {
        const { data: currentRoles } = await api.get(
            `/admin/user/${user.id}/roles`,
        );
        await Promise.all(
            (currentRoles ?? []).map((r) =>
                api.delete(`/admin/user/${user.id}/roles/${r.id}`),
            ),
        );
        await api.post(`/admin/user/${user.id}/roles`, { role_id: roleObj.id });
        user.type = roleToType([newRole]);
        rightsModal.value.open = false;
        adminConfirm.value.open = false;
    } catch (e) {
        alert("Erreur lors de la mise à jour du rôle");
    }
}

async function openHistory(user) {
    historyModal.value = {
        open: true,
        user,
        loading: true,
        data: {
            sanctions: [],
            reports_received: [],
            professional_requests: [],
        },
    };
    try {
        const { data } = await api.get(`/admin/user/${user.id}/history`);
        historyModal.value.data = data;
    } catch (e) {
        console.error("openHistory error:", e);
    } finally {
        historyModal.value.loading = false;
    }
}

function deleteUser(user) {
    toDelete.value = user;
}

async function confirmDeleteUser() {
    if (!toDelete.value) return;
    deleting.value = true;
    try {
        await api.delete(`/admin/user/${toDelete.value.id}`);
        toDelete.value = null;
        await fetchUsers();
    } catch (e) {
        alert("Erreur lors de la suppression");
    } finally {
        deleting.value = false;
    }
}

function statusBadge(status) {
    const map = {
        Actif: "bg-green-100 text-green-700",
        Suspendu: "bg-amber-100 text-amber-700",
        Blacklisté: "bg-red-100 text-red-700",
    };
    return map[status] ?? "bg-gray-100 text-gray-600";
}

function statusDot(status) {
    const map = {
        Actif: "bg-green-500",
        Suspendu: "bg-amber-500",
        Blacklisté: "bg-red-500",
    };
    return map[status] ?? "bg-gray-400";
}
</script>
