<template>
    <AdminLayout title="Formations">

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-4 gap-5 mb-8">
                <div v-for="stat in stats" :key="stat.label"
                    class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                    <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                        <div :class="stat.iconClass" v-html="stat.icon"/>
                    </div>
                    <div class="min-w-0">
                        <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                        <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm overflow-hidden">

                <div class="px-5 py-4 border-b border-gray-100 flex items-center gap-3">
                    <div class="relative flex-1">
                        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                        </svg>
                        <input v-model="search" type="text" placeholder="Rechercher une formation…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                    </div>
                    <select v-model="filterStatus"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option value="pending">En attente</option>
                        <option value="approved">Approuvées</option>
                        <option value="rejected">Rejetées</option>
                    </select>
                    <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} résultat(s)</span>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Formation</th>
                                <th class="text-left text-white font-medium px-5 py-3">Date & lieu</th>
                                <th class="text-left text-white font-medium px-5 py-3">Niveau</th>
                                <th class="text-left text-white font-medium px-5 py-3">Inscrits</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(f, i) in formations" :key="f.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                                <td class="px-5 py-3">
                                    <p class="font-medium text-gray-800 truncate max-w-52">{{ f.title }}</p>
                                    <p class="text-xs text-gray-500 truncate">{{ f.creator_name }}</p>
                                </td>
                                <td class="px-5 py-3">
                                    <p class="text-gray-800 text-xs font-medium">{{ f.date ? f.date.slice(0,10) : '—' }}</p>
                                    <p class="text-gray-500 text-xs truncate">{{ f.location ?? '—' }}</p>
                                </td>
                                <td class="px-5 py-3">
                                    <span :class="levelClass(f.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                        {{ levelLabel(f.level) }}
                                    </span>
                                </td>
                                <td class="px-5 py-3 text-xs text-gray-700 font-medium">
                                    {{ f.inscription_count }}<span v-if="f.capacity">/{{ f.capacity }}</span>
                                </td>
                                <td class="px-5 py-3">
                                    <span :class="statusClass(f.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                        {{ statusLabel(f.status) }}
                                    </span>
                                </td>
                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button v-if="f.status === 'pending'" @click="approveFormation(f)"
                                            title="Approuver"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                            </svg>
                                        </button>
                                        <button v-if="f.status === 'pending'" @click="openRejectModal(f)"
                                            title="Rejeter"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                            </svg>
                                        </button>
                                        <button @click="openDetail(f)" title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>
                                        <button @click="confirmDelete(f)" title="Supprimer"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                            <tr v-if="formations.length === 0">
                                <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucune formation ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ total }} formation(s) au total
                </div>
            </div>

            <Pagination v-if="total > 20" :page="page" :total="total" :limit="20" @update:page="changePage"/>

        </template>

        <!-- Modal détail -->
        <div v-if="detailFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="detailFormation = null">
            <div class="absolute inset-0 bg-black/40" @click="detailFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ detailFormation.title }}</h3>
                        <p class="text-xs text-gray-500 mt-0.5">Par {{ detailFormation.creator_name }}</p>
                    </div>
                    <button @click="detailFormation = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span :class="statusClass(detailFormation.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ statusLabel(detailFormation.status) }}
                            </span>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Niveau</p>
                            <span :class="levelClass(detailFormation.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ levelLabel(detailFormation.level) }}
                            </span>
                        </div>
                        <div v-if="detailFormation.date">
                            <p class="text-xs text-gray-400 mb-0.5">Date</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.date.slice(0,10) }}</p>
                        </div>
                        <div v-if="detailFormation.location">
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.location }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Inscriptions</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.inscription_count }}<span v-if="detailFormation.capacity">/{{ detailFormation.capacity }}</span></p>
                        </div>
                        <div v-if="detailFormation.duration_hours">
                            <p class="text-xs text-gray-400 mb-0.5">Durée</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.duration_hours }}h</p>
                        </div>
                    </div>
                    <div v-if="detailFormation.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailFormation.description }}</p>
                    </div>
                    <div v-if="detailFormation.status === 'rejected' && detailFormation.rejection_reason">
                        <p class="text-xs text-gray-400 mb-1">Raison du rejet</p>
                        <p class="text-sm text-red-600">{{ detailFormation.rejection_reason }}</p>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-between gap-2">
                    <div class="flex gap-2">
                        <button v-if="detailFormation.status === 'pending'" @click="approveFormation(detailFormation); detailFormation = null"
                            class="px-3 py-1.5 text-sm font-medium text-secondary border border-secondary/30 rounded-lg hover:bg-secondary/5 transition-colors">
                            Approuver
                        </button>
                        <button v-if="detailFormation.status === 'pending'" @click="openRejectModal(detailFormation); detailFormation = null"
                            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Rejeter
                        </button>
                    </div>
                    <button @click="confirmDelete(detailFormation); detailFormation = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Supprimer
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal rejet -->
        <div v-if="rejectModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="rejectModal = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Rejeter la formation</h3>
                <p class="text-sm text-gray-500 mb-4">« {{ rejectModal.title }} »</p>
                <div class="mb-4">
                    <label class="block text-xs text-gray-400 mb-1">Raison (optionnelle)</label>
                    <textarea v-model="rejectReason" rows="3"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                        placeholder="Expliquez pourquoi la formation est rejetée…"/>
                </div>
                <div class="flex gap-3">
                    <button @click="rejectModal = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50">Annuler</button>
                    <button @click="submitReject" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600">Rejeter</button>
                </div>
            </div>
        </div>

        <!-- Modal suppression -->
        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer la formation ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimée avec toutes ses inscriptions.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50">Annuler</button>
                    <button @click="deleteFormation" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const formations = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const error = ref('')
const search = ref('')
const filterStatus = ref('')
const detailFormation = ref(null)
const rejectModal = ref(null)
const rejectReason = ref('')
const toDelete = ref(null)
const deleting = ref(false)

const stats = computed(() => {
    const all = formations.value
    return [
        {
            label: 'Total formations',
            value: total.value,
            bgClass: 'bg-primary/10',
            iconClass: 'w-6 h-6 text-primary',
            icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>`,
        },
        {
            label: 'En attente',
            value: all.filter(f => f.status === 'pending').length,
            bgClass: 'bg-yellow-100',
            iconClass: 'w-6 h-6 text-yellow-600',
            icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
        },
        {
            label: 'Approuvées',
            value: all.filter(f => f.status === 'approved').length,
            bgClass: 'bg-green-100',
            iconClass: 'w-6 h-6 text-green-600',
            icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
        },
        {
            label: 'Rejetées',
            value: all.filter(f => f.status === 'rejected').length,
            bgClass: 'bg-red-100',
            iconClass: 'w-6 h-6 text-red-500',
            icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
        },
    ]
})

function levelLabel(level) {
    return { beginner: 'Débutant', intermediate: 'Intermédiaire', advanced: 'Avancé' }[level] ?? level
}

function levelClass(level) {
    return {
        beginner: 'bg-secondary/10 text-secondary',
        intermediate: 'bg-primary/10 text-primary',
        advanced: 'bg-red-100 text-red-500',
    }[level] ?? 'bg-gray-100 text-gray-500'
}

function statusLabel(status) {
    return { pending: 'En attente', approved: 'Approuvée', rejected: 'Rejetée' }[status] ?? status
}

function statusClass(status) {
    return {
        pending: 'bg-yellow-100 text-yellow-700',
        approved: 'bg-green-100 text-green-700',
        rejected: 'bg-red-100 text-red-600',
    }[status] ?? 'bg-gray-100 text-gray-500'
}

function openDetail(f) {
    detailFormation.value = f
}

function confirmDelete(f) {
    toDelete.value = f
}

function openRejectModal(f) {
    rejectModal.value = f
    rejectReason.value = ''
}

async function approveFormation(f) {
    try {
        await api.patch(`/formations/${f.id}/approve`)
        f.status = 'approved'
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de l\'approbation.')
    }
}

async function submitReject() {
    if (!rejectModal.value) return
    try {
        await api.patch(`/formations/${rejectModal.value.id}/reject`, { reason: rejectReason.value || null })
        const f = formations.value.find(x => x.id === rejectModal.value.id)
        if (f) {
            f.status = 'rejected'
            f.rejection_reason = rejectReason.value || null
        }
        rejectModal.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors du rejet.')
    }
}

async function deleteFormation() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/formation/${toDelete.value.id}`)
        formations.value = formations.value.filter(f => f.id !== toDelete.value.id)
        total.value = Math.max(0, total.value - 1)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function changePage(p) {
    page.value = p
}

async function fetchFormations() {
    loading.value = true
    error.value = ''
    try {
        const { data } = await api.get('/admin/formations', {
            params: { page: page.value, limit: 20, search: search.value, status: filterStatus.value }
        })
        formations.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (e) {
        error.value = 'Erreur lors du chargement des formations.'
    } finally {
        loading.value = false
    }
}

watch([page, filterStatus], fetchFormations)
watch(search, () => { page.value = 1; fetchFormations() })

onMounted(fetchFormations)
</script>
