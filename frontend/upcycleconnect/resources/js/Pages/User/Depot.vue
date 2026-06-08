<template>
    <UserLayout>

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Dépôt d'objet</h1>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6 mb-6">
            <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Mes demandes de dépôt</h2>
            <p class="text-xs text-gray-400 mb-4">Statut de vos demandes d'assignation de casier</p>

            <div v-if="loadingMine" class="text-center py-8 text-gray-400 text-sm">Chargement…</div>

            <p v-else-if="myRequests.length === 0" class="text-sm text-gray-400 text-center py-6">
                Vous n'avez aucune demande de dépôt. Faites une demande depuis vos annonces.
            </p>

            <div v-else class="space-y-3">
                <div
                    v-for="req in myRequests"
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
                            <p class="text-xs text-gray-400 capitalize mt-0.5">{{ req.announcement_type === 'vente' ? 'Vente' : 'Don' }}</p>
                        </div>
                    </div>

                    <div class="flex items-center gap-6 flex-shrink-0">
                        <div v-if="req.locker_number" class="text-right">
                            <p class="text-xs text-gray-400">Casier assigné</p>
                            <p class="text-sm font-semibold text-gray-800">N° {{ req.locker_number }}</p>
                        </div>

                        <span class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap bg-yellow-100 text-yellow-700">
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

        <div v-if="auth.hasPermission('validate_deposit')" class="bg-white rounded-2xl shadow-sm p-6">
            <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Demandes en attente</h2>
            <p class="text-xs text-gray-400 mb-4">Demandes d'assignation de casier à valider</p>

            <div v-if="loadingPending" class="text-center py-8 text-gray-400 text-sm">Chargement…</div>

            <p v-else-if="pendingRequests.length === 0" class="text-sm text-gray-400 text-center py-6">
                Aucune demande en attente de validation.
            </p>

            <div v-else class="space-y-3">
                <div
                    v-for="ann in pendingRequests"
                    :key="ann.id"
                    class="flex items-center justify-between p-4 rounded-xl border border-gray-100 hover:border-gray-200 transition-colors"
                >
                    <div class="flex items-center gap-4 min-w-0">
                        <div class="w-10 h-10 rounded-xl bg-secondary/10 flex items-center justify-center flex-shrink-0">
                            <svg class="w-5 h-5 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                            </svg>
                        </div>
                        <div class="min-w-0">
                            <p class="font-medium text-gray-800 truncate">{{ ann.title || 'Annonce sans titre' }}</p>
                            <div class="flex items-center gap-2 mt-0.5">
                                <span :class="['px-1.5 py-0.5 rounded text-xs font-medium text-white', ann.type === 'vente' ? 'bg-primary' : 'bg-secondary']">
                                    {{ ann.type === 'vente' ? 'Vente' : 'Don' }}
                                </span>
                                <span class="text-xs text-gray-400">{{ ann.author_name }}</span>
                                <span class="text-xs text-gray-400">· Dispo. {{ formatDate(ann.availability_date) }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="flex items-center gap-2 flex-shrink-0">
                        <button
                            @click="rejectRequest(ann.id)"
                            :disabled="validating === ann.id || rejecting === ann.id"
                            class="flex items-center gap-2 px-4 py-2 bg-red-100 text-red-700 text-sm font-medium rounded-lg hover:bg-red-200 transition-colors disabled:opacity-60"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                            </svg>
                            {{ rejecting === ann.id ? 'Refus…' : 'Refuser' }}
                        </button>
                        <button
                            @click="validateRequest(ann.id)"
                            :disabled="validating === ann.id || rejecting === ann.id"
                            class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors disabled:opacity-60"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                            </svg>
                            {{ validating === ann.id ? 'Validation…' : 'Valider' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>

    </UserLayout>
</template>

<script setup>
import { ref } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
import UserLayout from '@/Layouts/UserLayout.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'
import { formatDate } from '@/utils.js'

const auth = useAuthStore()

const myRequests = ref([])
const pendingRequests = ref([])
const loadingMine = ref(true)
const loadingPending = ref(true)
const validating = ref(null)
const rejecting = ref(null)

async function loadMyRequests(silent = false) {
    if (!silent) loadingMine.value = true
    try {
        const { data } = await api.get('/user/deposit-requests')
        myRequests.value = data ?? []
    } catch {
        myRequests.value = []
    } finally {
        loadingMine.value = false
    }
}

async function loadPendingRequests(silent = false) {
    if (!auth.hasPermission('validate_deposit')) {
        if (!silent) loadingPending.value = false
        return
    }
    if (!silent) loadingPending.value = true
    try {
        const { data } = await api.get('/deposit-requests/pending')
        pendingRequests.value = data ?? []
    } catch {
        pendingRequests.value = []
    } finally {
        loadingPending.value = false
    }
}


async function cancelRequest(announcementId) {
    try {
        await api.delete(`/announcements/${announcementId}/deposit-request`)
        await loadMyRequests()
    } catch (e) {
        alert(e.response?.data ?? 'Erreur lors de l\'annulation.')
    }
}

async function validateRequest(announcementId) {
    validating.value = announcementId
    try {
        await api.post(`/deposit-requests/${announcementId}/validate`)
        await loadPendingRequests()
    } catch (e) {
        if (e.response?.status === 409) {
            alert('Aucun casier disponible. La demande ne peut pas être validée.')
        } else {
            alert(e.response?.data ?? 'Erreur lors de la validation.')
        }
    } finally {
        validating.value = null
    }
}

async function rejectRequest(announcementId) {
    rejecting.value = announcementId
    try {
        await api.post(`/deposit-requests/${announcementId}/reject`)
        await loadPendingRequests()
    } catch (e) {
        alert(e.response?.data ?? 'Erreur lors du refus.')
    } finally {
        rejecting.value = null
    }
}

async function fetchAll(silent = false) {
    await Promise.all([loadMyRequests(silent), loadPendingRequests(silent)])
}

usePolling(fetchAll)
</script>
