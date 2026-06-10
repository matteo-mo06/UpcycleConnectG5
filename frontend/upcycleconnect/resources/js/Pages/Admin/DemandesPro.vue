<template>
    <AdminLayout>
        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Demandes professionnelles</h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">
            Chargement...
        </div>

        <div v-if="error" class="py-12 text-center text-sm text-red-500">
            {{ error }}
        </div>

        <template v-if="!loading && !error">
            <div class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
                <div class="px-5 py-4 border-b border-gray-100 flex items-center gap-2">
                    <h2 class="text-base font-semibold text-gray-800">En attente</h2>
                    <span
                        v-if="pending.length"
                        class="inline-flex items-center justify-center w-5 h-5 rounded-full bg-primary text-white text-xs font-bold"
                    >{{ pending.length }}</span>
                </div>
                <div class="divide-y divide-gray-50">
                    <div v-if="pending.length === 0" class="px-5 py-8 text-center text-sm text-gray-400">
                        Aucune demande en attente
                    </div>
                    <div
                        v-for="req in pending"
                        :key="req.id"
                        class="px-5 py-3.5 flex items-center justify-between gap-4"
                    >
                        <div class="flex items-center gap-3">
                            <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
                                <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z" />
                                </svg>
                            </div>
                            <div class="min-w-0">
                                <p class="text-sm font-medium text-gray-800 truncate">{{ req.name }}</p>
                                <p class="text-xs text-gray-500 truncate">
                                    {{ req.email }} · Demande le {{ req.requestDate }}
                                </p>
                            </div>
                        </div>
                        <div class="flex items-center gap-2 flex-shrink-0">
                            <button
                                @click="approveRequest(req.id)"
                                class="px-3 py-1.5 rounded-lg text-xs font-medium bg-secondary text-white hover:bg-secondary-dark transition-colors duration-150"
                            >Valider</button>
                            <button
                                @click="rejectRequest(req.id)"
                                class="px-3 py-1.5 rounded-lg text-xs font-medium bg-gray-100 text-gray-600 hover:bg-red-100 hover:text-red-600 transition-colors duration-150"
                            >Refuser</button>
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm overflow-hidden">
                <div class="px-5 py-4 border-b border-gray-100">
                    <h2 class="text-base font-semibold text-gray-800">Historique</h2>
                </div>
                <div class="divide-y divide-gray-50">
                    <div v-if="history.length === 0" class="px-5 py-8 text-center text-sm text-gray-400">
                        Aucune demande traitée
                    </div>
                    <div
                        v-for="req in history"
                        :key="req.id"
                        class="px-5 py-3.5 flex items-center justify-between gap-4"
                    >
                        <div class="flex items-center gap-3">
                            <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
                                <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z" />
                                </svg>
                            </div>
                            <div class="min-w-0">
                                <p class="text-sm font-medium text-gray-800 truncate">{{ req.name }}</p>
                                <p class="text-xs text-gray-500 truncate">
                                    {{ req.email }} · Demande le {{ req.requestDate }}
                                </p>
                            </div>
                        </div>
                        <span
                            class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0"
                            :class="statusBadge(req.status)"
                        >
                            <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(req.status)" />
                            {{ statusLabel(req.status) }}
                        </span>
                    </div>
                </div>
            </div>
        </template>
    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import api from "@/api.js";
import AdminLayout from "@/Layouts/AdminLayout.vue";
import { fullName } from "@/utils.js";

const requests = ref([]);
const loading = ref(false);
const error = ref("");

const pending = computed(() => requests.value.filter((r) => r.status === "pending"));
const history = computed(() => requests.value.filter((r) => r.status !== "pending"));

onMounted(fetchRequests);

async function fetchRequests() {
    loading.value = true;
    try {
        const { data } = await api.get("/admin/professional-requests");
        requests.value = (data ?? []).map((r) => ({
            id: r.id,
            name: fullName(r),
            email: r.email,
            requestDate: r.created_at?.slice(0, 10) ?? "-",
            status: r.status,
        }));
    } catch (e) {
        error.value = "Impossible de charger les demandes";
    } finally {
        loading.value = false;
    }
}

async function approveRequest(id) {
    try {
        await api.put(`/admin/professional-requests/${id}/validate`);
        const req = requests.value.find((r) => r.id === id);
        if (req) req.status = "approved";
    } catch (e) {
        alert("Erreur lors de la validation");
    }
}

async function rejectRequest(id) {
    try {
        await api.put(`/admin/professional-requests/${id}/reject`);
        const req = requests.value.find((r) => r.id === id);
        if (req) req.status = "rejected";
    } catch (e) {
        alert("Erreur lors du refus");
    }
}

function statusBadge(status) {
    return status === "approved" ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700";
}

function statusDot(status) {
    return status === "approved" ? "bg-green-500" : "bg-red-500";
}

function statusLabel(status) {
    return status === "approved" ? "Approuvée" : "Refusée";
}
</script>
