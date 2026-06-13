<template>
    <ArtisanLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Dépôt d'objet</h1>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6">
            <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Demandes de casier</h2>
            <p class="text-xs text-gray-400 mb-4">Statut de vos demandes d'assignation de casier</p>

            <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Chargement…</div>

            <p v-else-if="requests.length === 0" class="text-sm text-gray-400 text-center py-6">
                Aucune demande de dépôt en cours. Faites une demande depuis vos annonces vendues.
            </p>

            <div v-else class="space-y-3">
                <div
                    v-for="req in requests"
                    :key="req.announcement_id"
                    class="flex items-center justify-between p-4 rounded-xl border border-gray-100 hover:border-gray-200 transition-colors"
                >
                    <div class="flex items-center gap-4 min-w-0">
                        <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                            <svg class="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                            </svg>
                        </div>
                        <div class="min-w-0">
                            <p class="font-medium text-gray-800 truncate">{{ req.announcement_title || 'Annonce sans titre' }}</p>
                            <p class="text-xs text-gray-400 mt-0.5">{{ req.announcement_type === 'vente' ? 'Vente' : 'Don' }}</p>
                        </div>
                    </div>

                    <div class="flex items-center gap-6 flex-shrink-0">
                        <div v-if="req.locker_number" class="text-right">
                            <p class="text-xs text-gray-400">Casier assigné</p>
                            <p class="text-sm font-semibold text-gray-800">N° {{ req.locker_number }}</p>
                        </div>
                        <span class="px-3 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-700">
                            En attente de validation
                        </span>
                        <button
                            @click="cancelRequest(req.announcement_id)"
                            class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                            title="Annuler la demande"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>

    </ArtisanLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import api from '@/api.js'

const requests = ref([])
const loading = ref(true)

async function fetchRequests() {
    loading.value = true
    try {
        const { data } = await api.get('/user/deposit-requests')
        requests.value = data ?? []
    } catch {
        requests.value = []
    } finally {
        loading.value = false
    }
}

async function cancelRequest(announcementId) {
    try {
        await api.delete(`/announcements/${announcementId}/deposit-request`)
        await fetchRequests()
    } catch (e) {
        alert(e.response?.data ?? 'Erreur lors de l\'annulation.')
    }
}

onMounted(fetchRequests)
</script>
