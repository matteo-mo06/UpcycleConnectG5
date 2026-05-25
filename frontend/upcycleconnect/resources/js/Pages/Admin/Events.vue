<template>
    <AdminLayout title="Événements">

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-4 gap-5 mb-8">
                <div
                    v-for="stat in stats"
                    :key="stat.label"
                    class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                    <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                        <div :class="stat.iconClass" v-html="stat.icon" />
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
                        <input
                            v-model="search"
                            type="text"
                            placeholder="Rechercher par titre ou organisateur…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>

                    <select
                        v-model="filterStatus"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option value="pending">En attente</option>
                        <option value="approved-upcoming">À venir</option>
                        <option value="approved-past">Passé</option>
                        <option value="rejected">Refusé</option>
                    </select>

                    <button
                        @click="openCreate"
                        class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors whitespace-nowrap">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                        </svg>
                        Créer
                    </button>

                    <span class="text-xs text-gray-400 whitespace-nowrap">
                        {{ total }} résultat(s)
                    </span>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Événement</th>
                                <th class="text-left text-white font-medium px-5 py-3">Date & lieu</th>
                                <th class="text-left text-white font-medium px-5 py-3">Inscrits</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(event, i) in events"
                                :key="event.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                                <td class="px-5 py-3">
                                    <p class="font-medium text-gray-800 truncate max-w-52">{{ event.title }}</p>
                                    <p class="text-xs text-gray-500 truncate">{{ event.organizer }}</p>
                                </td>

                                <td class="px-5 py-3">
                                    <p class="text-gray-800 text-xs font-medium">{{ event.date }} · {{ event.time }}</p>
                                    <p class="text-gray-500 text-xs truncate">{{ event.location }}</p>
                                </td>

                                <td class="px-5 py-3 text-xs text-gray-700 font-medium">
                                    {{ event.registered }}/{{ event.capacity }}
                                </td>

                                <td class="px-5 py-3">
                                    <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                        <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(event.status)" />
                                        <span :class="statusText(event.status)">{{ event.status }}</span>
                                    </span>
                                </td>

                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button
                                            @click="openDetail(event)"
                                            title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>

                                        <template v-if="event.dbStatus === 'pending'">
                                            <button
                                                @click="approveEvent(event)"
                                                title="Approuver"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                </svg>
                                            </button>
                                            <button
                                                @click="openRejectEvent(event)"
                                                title="Refuser"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                </svg>
                                            </button>
                                        </template>
                                        <template v-else>
                                            <button
                                                @click="openEdit(event)"
                                                title="Modifier l'événement"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                                                </svg>
                                            </button>
                                        </template>

                                        <button
                                            @click="confirmDelete(event)"
                                            title="Supprimer l'événement"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>

                            <tr v-if="events.length === 0">
                                <td colspan="5" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucun événement ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ total }} événement(s) au total
                </div>
            </div>
            <Pagination v-if="total > 20" :page="page" :total="total" :limit="20" @update:page="changePage" />

        </template>

        <div
            v-if="detailEvent"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="detailEvent = null">
            <div class="absolute inset-0 bg-black/40" @click="detailEvent = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                            {{ detailEvent.title }}
                        </h3>
                        <p class="text-xs text-gray-500 mt-0.5">Organisé par {{ detailEvent.organizer }}</p>
                    </div>
                    <button
                        @click="detailEvent = null"
                        class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Date & heure</p>
                            <p class="font-medium text-gray-800">{{ detailEvent.date }}</p>
                            <p class="text-xs text-gray-500">{{ detailEvent.time }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ detailEvent.location }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Inscriptions</p>
                            <p class="font-medium text-gray-800">{{ detailEvent.registered }} / {{ detailEvent.capacity }} participants</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailEvent.status)" />
                                <span :class="statusText(detailEvent.status)">{{ detailEvent.status }}</span>
                            </span>
                        </div>
                    </div>

                    <div>
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailEvent.description }}</p>
                    </div>
                    <div v-if="detailEvent.rejectionReason" class="p-3 rounded-xl bg-red-50 border border-red-100">
                        <p class="text-xs font-medium text-red-600 mb-0.5">Motif de refus</p>
                        <p class="text-sm text-red-700">{{ detailEvent.rejectionReason }}</p>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-between gap-2">
                    <div class="flex gap-2">
                        <button
                            v-if="detailEvent.dbStatus === 'pending'"
                            @click="approveEvent(detailEvent); detailEvent = null"
                            class="px-3 py-1.5 text-sm font-medium text-secondary border border-secondary/30 rounded-lg hover:bg-secondary/5 transition-colors">
                            Approuver
                        </button>
                        <button
                            v-if="detailEvent.dbStatus === 'pending'"
                            @click="openRejectEvent(detailEvent); detailEvent = null"
                            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Refuser
                        </button>
                        <button
                            v-if="detailEvent.dbStatus !== 'pending'"
                            @click="openEdit(detailEvent)"
                            class="px-3 py-1.5 text-sm font-medium text-primary border border-primary/30 rounded-lg hover:bg-primary/5 transition-colors">
                            Modifier
                        </button>
                    </div>
                    <button
                        @click="confirmDelete(detailEvent); detailEvent = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Supprimer
                    </button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer l'événement ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimé.</p>
                <div class="flex gap-3">
                    <button
                        @click="toDelete = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="deleteEvent"
                        :disabled="deleting"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="rejectModal.open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="rejectModal.open = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-1">Refuser l'événement</h3>
                <p class="text-sm text-gray-500 mb-4">« {{ rejectModal.event?.title }} »</p>
                <div class="mb-4">
                    <label class="block text-xs text-gray-400 mb-1">Motif de refus (optionnel)</label>
                    <textarea
                        v-model="rejectModal.reason"
                        rows="3"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-300 focus:border-red-400 resize-none"
                        placeholder="Expliquez pourquoi cet événement est refusé…" />
                </div>
                <div class="flex gap-3">
                    <button @click="rejectModal.open = false" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="confirmRejectEvent" :disabled="rejectModal.loading" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ rejectModal.loading ? 'Refus…' : 'Refuser' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="formModal.open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal.open = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ formModal.mode === 'create' ? 'Créer un événement' : 'Modifier l\'événement' }}
                    </h3>
                    <button
                        @click="formModal.open = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4">
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input
                            v-model="eventForm.title"
                            type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre de l'événement" />
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea
                            v-model="eventForm.description"
                            rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de l'événement…" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date & heure</label>
                            <input
                                v-model="eventForm.date"
                                type="datetime-local"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Lieu</label>
                            <input
                                v-model="eventForm.location"
                                type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse ou ville" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Capacité</label>
                            <input
                                v-model.number="eventForm.capacity"
                                type="number"
                                min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Prix (€)</label>
                            <input
                                v-model.number="eventForm.priceEuros"
                                type="number"
                                min="0"
                                step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = gratuit" />
                        </div>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button
                        @click="formModal.open = false"
                        class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="saveEvent"
                        :disabled="saving"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ saving ? 'Enregistrement…' : (formModal.mode === 'create' ? 'Créer' : 'Enregistrer') }}
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

const events = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(true)
const error = ref('')
const search = ref('')
const filterStatus = ref('')
const detailEvent = ref(null)
const toDelete = ref(null)
const deleting = ref(false)
const formModal = ref({ open: false, mode: 'create', id: null })
const rejectModal = ref({ open: false, event: null, reason: '', loading: false })
const eventForm = ref({ title: '', description: '', date: '', location: '', capacity: null, priceEuros: 0 })
const saving = ref(false)

const stats = computed(() => [
    {
        label: 'Total événements',
        value: total.value,
        bgClass: 'bg-red-100',
        iconClass: 'text-red-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
    },
    {
        label: 'En attente',
        value: events.value.filter(e => e.dbStatus === 'pending').length,
        bgClass: 'bg-orange-100',
        iconClass: 'text-orange-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'À venir',
        value: events.value.filter(e => e.status === 'À venir').length,
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Refusés',
        value: events.value.filter(e => e.dbStatus === 'rejected').length,
        bgClass: 'bg-red-100',
        iconClass: 'text-red-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
])

async function fetchEvents() {
    loading.value = true
    error.value = ''
    try {
        const params = { page: page.value, limit: 20 }
        if (search.value) params.search = search.value
        const statusMap = { 'pending': 'pending', 'approved-upcoming': 'upcoming', 'approved-past': 'past', 'rejected': 'rejected' }
        if (filterStatus.value) params.status = statusMap[filterStatus.value] ?? ''
        const { data } = await api.get('/admin/events', { params })
        events.value = data.data.map(mapEvent)
        total.value = data.total
    } catch {
        error.value = 'Impossible de charger les événements.'
    } finally {
        loading.value = false
    }
}

let debounce = null
watch(search, () => {
    clearTimeout(debounce)
    debounce = setTimeout(() => { page.value = 1; fetchEvents() }, 300)
})
watch(filterStatus, () => { page.value = 1; fetchEvents() })
watch(page, fetchEvents)

function computeStatus(e) {
    if (e.status === 'pending') return 'En attente'
    if (e.status === 'rejected') return 'Refusé'
    if (!e.date) return '-'
    return new Date(e.date) > new Date() ? 'À venir' : 'Passé'
}

function mapEvent(e) {
    return {
        id: e.id,
        title: e.title,
        organizer: e.creator_name ?? '-',
        location: e.location ?? '-',
        date: e.date?.slice(0, 10) ?? '-',
        time: e.date?.slice(11, 16) ?? '-',
        capacity: e.capacity ?? 0,
        registered: e.inscription_count ?? 0,
        dbStatus: e.status ?? 'approved',
        status: computeStatus(e),
        rejectionReason: e.rejection_reason ?? null,
        description: e.description ?? '-',
        priceCents: e.price_cents ?? 0,
        rawDate: e.date ?? null,
        idCreator: e.id_creator ?? null,
    }
}

onMounted(fetchEvents)

function openDetail(event) {
    detailEvent.value = event
}

function confirmDelete(event) {
    toDelete.value = event
}

async function deleteEvent() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/event/${toDelete.value.id}`)
        if (detailEvent.value?.id === toDelete.value.id) detailEvent.value = null
        toDelete.value = null
        await fetchEvents()
    } catch {
        alert('Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function changePage(p) { page.value = p }

function openCreate() {
    eventForm.value = { title: '', description: '', date: '', location: '', capacity: null, priceEuros: 0 }
    formModal.value = { open: true, mode: 'create', id: null }
}

function openEdit(event) {
    eventForm.value = {
        title: event.title,
        description: event.description === '-' ? '' : event.description,
        date: event.rawDate ? event.rawDate.slice(0, 16) : '',
        location: event.location === '-' ? '' : event.location,
        capacity: event.capacity || null,
        priceEuros: parseFloat((event.priceCents / 100).toFixed(2)),
        idCreator: event.idCreator ?? null,
    }
    formModal.value = { open: true, mode: 'edit', id: event.id }
    detailEvent.value = null
}

async function saveEvent() {
    if (!eventForm.value.title.trim()) return
    saving.value = true
    try {
        const payload = {
            title: eventForm.value.title,
            description: eventForm.value.description || null,
            date: eventForm.value.date || null,
            location: eventForm.value.location || null,
            capacity: eventForm.value.capacity || null,
            price_cents: Math.round((eventForm.value.priceEuros || 0) * 100),
            id_creator: eventForm.value.idCreator ?? null,
        }
        if (formModal.value.mode === 'create') {
            await api.post('/admin/events', payload)
            page.value = 1
        } else {
            await api.put(`/admin/event/${formModal.value.id}`, payload)
        }
        formModal.value.open = false
        await fetchEvents()
    } catch {
        alert('Erreur lors de l\'enregistrement.')
    } finally {
        saving.value = false
    }
}

function statusDot(status) {
    const map = { 'À venir': 'bg-green-500', 'Passé': 'bg-gray-400', 'En attente': 'bg-orange-400', 'Refusé': 'bg-red-400' }
    return map[status] ?? 'bg-gray-300'
}

function statusText(status) {
    const map = { 'À venir': 'text-green-700', 'Passé': 'text-gray-500', 'En attente': 'text-orange-600', 'Refusé': 'text-red-600' }
    return map[status] ?? 'text-gray-500'
}

async function approveEvent(event) {
    try {
        await api.patch(`/admin/event/${event.id}/approve`)
        await fetchEvents()
    } catch {
        alert('Erreur lors de l\'approbation.')
    }
}

function openRejectEvent(event) {
    rejectModal.value = { open: true, event, reason: '', loading: false }
}

async function confirmRejectEvent() {
    if (!rejectModal.value.event) return
    rejectModal.value.loading = true
    try {
        await api.patch(`/admin/event/${rejectModal.value.event.id}/reject`, { reason: rejectModal.value.reason })
        rejectModal.value.open = false
        await fetchEvents()
    } catch {
        alert('Erreur lors du refus.')
    } finally {
        rejectModal.value.loading = false
    }
}
</script>
