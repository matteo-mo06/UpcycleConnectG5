<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Casiers</h1>
        </div>

        <div class="grid grid-cols-3 gap-5 mb-8">
            <div v-for="stat in stats" :key="stat.label" class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                    <div :class="stat.iconClass" v-html="stat.icon" />
                </div>
                <div class="min-w-0">
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                    <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 items-center">
            <select
                v-model="filterStatus"
                @change="fetchLockers"
                class="px-4 py-2 rounded-lg border border-gray-200 text-sm text-gray-600 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            >
                <option value="">Tous les statuts</option>
                <option value="free">Libre</option>
                <option value="occupied">Occupé</option>
            </select>
            <span class="text-xs text-gray-400 ml-auto">{{ lockers.length }} casier(s)</span>
            <button
                @click="addLocker"
                :disabled="adding"
                class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60 whitespace-nowrap"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                {{ adding ? 'Ajout…' : 'Ajouter un casier' }}
            </button>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <div v-if="!loading && !error" class="bg-white rounded-xl shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead>
                        <tr class="bg-primary">
                            <th class="text-left text-white font-medium px-5 py-3">N° casier</th>
                            <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                            <th class="text-left text-white font-medium px-5 py-3">Annonce associée</th>
                            <th class="text-left text-white font-medium px-5 py-3">Vendeur</th>
                            <th class="text-left text-white font-medium px-5 py-3">Acheteur</th>
                            <th class="text-left text-white font-medium px-5 py-3">Code d'accès</th>
                            <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="(locker, i) in lockers"
                            :key="locker.id"
                            :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']"
                        >
                            <td class="px-5 py-3 font-semibold text-gray-800">#{{ locker.number }}</td>

                            <td class="px-5 py-3">
                                <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                    <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="locker.occupied ? 'bg-amber-400' : 'bg-green-500'" />
                                    <span :class="locker.occupied ? 'text-amber-600' : 'text-green-700'">
                                        {{ locker.occupied ? 'Occupé' : 'Libre' }}
                                    </span>
                                </span>
                            </td>

                            <td class="px-5 py-3 text-gray-500 text-xs">{{ locker.announcement_title || '-' }}</td>

                            <td class="px-5 py-3 text-gray-500 text-xs">{{ locker.seller_name || '-' }}</td>

                            <td class="px-5 py-3 text-gray-500 text-xs">{{ locker.buyer_name || '-' }}</td>

                            <td class="px-5 py-3 font-mono text-xs text-gray-600">{{ locker.access_code || '-' }}</td>

                            <td class="px-5 py-3">
                                <div class="flex items-center justify-end gap-1">
                                    <button
                                        v-if="locker.occupied"
                                        @click="openEditCode(locker)"
                                        title="Modifier le code d'accès"
                                        class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors"
                                    >
                                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                                        </svg>
                                    </button>
                                    <button
                                        v-if="!locker.occupied"
                                        @click="confirmDelete(locker)"
                                        title="Supprimer le casier"
                                        class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors"
                                    >
                                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>

                        <tr v-if="lockers.length === 0">
                            <td colspan="7" class="px-5 py-12 text-center text-gray-400 text-sm">
                                Aucun casier ne correspond à vos filtres.
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                {{ lockers.length }} casier(s) au total
            </div>
        </div>

        <div v-if="toEditCode" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toEditCode = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-1">Modifier le code d'accès</h3>
                <p class="text-sm text-gray-500 mb-3">Casier #{{ toEditCode.number }}</p>
                <label class="block text-sm font-medium text-gray-700 mb-1">Code d'accès <span class="text-red-400">*</span></label>
                <input
                    v-model="newAccessCode"
                    type="text"
                    maxlength="8"
                    placeholder="12345678"
                    class="w-full px-3 py-2 border rounded-lg text-sm font-mono focus:outline-none focus:ring-1 mb-1"
                    :class="newAccessCode && !/^\d{8}$/.test(newAccessCode) ? 'border-red-300 focus:border-red-400 focus:ring-red-300' : 'border-gray-200 focus:border-primary focus:ring-primary'"
                    @keydown.enter="saveAccessCode"
                />
                <p v-if="newAccessCode && !/^\d{8}$/.test(newAccessCode)" class="text-xs text-red-500 mb-3">
                    8 chiffres exactement (0–9)
                </p>
                <div v-else class="mb-3" />
                <div class="flex gap-3">
                    <button
                        @click="toEditCode = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Annuler
                    </button>
                    <button
                        @click="saveAccessCode"
                        :disabled="savingCode || !/^\d{8}$/.test(newAccessCode)"
                        class="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-60"
                    >
                        {{ savingCode ? 'Enregistrement…' : 'Enregistrer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le casier #{{ toDelete.number }} ?</h3>
                <p class="text-sm text-gray-500 mb-5">Cette action est irréversible.</p>
                <div class="flex gap-3">
                    <button
                        @click="toDelete = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Annuler
                    </button>
                    <button
                        @click="deleteLocker"
                        :disabled="deleting"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60"
                    >
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const lockers = ref([])
const loading = ref(true)
const error = ref('')
const adding = ref(false)
const toDelete = ref(null)
const deleting = ref(false)
const filterStatus = ref('')
const toEditCode = ref(null)
const newAccessCode = ref('')
const savingCode = ref(false)

const stats = computed(() => [
    {
        label: 'Total casiers',
        value: lockers.value.length,
        bgClass: 'bg-blue-100',
        iconClass: 'text-blue-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>`,
    },
    {
        label: 'Libres',
        value: lockers.value.filter(l => !l.occupied).length,
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Occupés',
        value: lockers.value.filter(l => l.occupied).length,
        bgClass: 'bg-amber-100',
        iconClass: 'text-amber-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>`,
    },
])

async function fetchLockers() {
    try {
        const params = {}
        if (filterStatus.value) params.status = filterStatus.value
        const { data } = await api.get('/admin/lockers', { params })
        lockers.value = data ?? []
    } catch {
        error.value = 'Impossible de charger les casiers.'
    }
}

onMounted(async () => {
    await fetchLockers()
    loading.value = false
})

async function addLocker() {
    adding.value = true
    try {
        await api.post('/admin/lockers')
        await fetchLockers()
    } catch {
        alert('Erreur lors de la création du casier.')
    } finally {
        adding.value = false
    }
}

function openEditCode(locker) {
    toEditCode.value = locker
    newAccessCode.value = locker.access_code || ''
}

async function saveAccessCode() {
    if (!toEditCode.value || !newAccessCode.value.trim()) return
    savingCode.value = true
    try {
        await api.patch(`/admin/locker/${toEditCode.value.id}/access-code`, { access_code: newAccessCode.value.trim() })
        toEditCode.value = null
        await fetchLockers()
    } catch {
        alert('Erreur lors de la mise à jour du code.')
    } finally {
        savingCode.value = false
    }
}

function confirmDelete(locker) {
    toDelete.value = locker
}

async function deleteLocker() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/locker/${toDelete.value.id}`)
        toDelete.value = null
        await fetchLockers()
    } catch (e) {
        if (e?.response?.status === 409) {
            alert('Ce casier est actuellement occupé et ne peut pas être supprimé.')
        } else {
            alert('Erreur lors de la suppression.')
        }
    } finally {
        deleting.value = false
    }
}
</script>
