<template>
    <UserLayout>

        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                    Événements
                </h1>
                <p class="text-sm text-gray-500 mt-1">Découvrez et inscrivez-vous aux événements disponibles</p>
            </div>
            <button
                v-if="canCreate"
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un événement
            </button>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6 mb-6">
            <h2 class="font-semibold text-gray-800 mb-4" style="font-family: var(--font-family-title)">
                Mes événements
            </h2>
            <FullCalendar :options="calendarOptions" />
        </div>

        <div v-if="canCreate && myEvents.length > 0" class="bg-white rounded-2xl shadow-sm overflow-hidden mb-6">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">Mes événements créés</h2>
                <span class="text-xs text-gray-400">{{ myEvents.length }} événement(s)</span>
            </div>
            <div class="divide-y divide-gray-50">
                <div v-for="event in myEvents" :key="event.id" class="flex items-center gap-4 px-6 py-4">
                    <div class="flex-shrink-0 w-14 text-center">
                        <div class="bg-gray-100 rounded-xl px-2 py-2">
                            <p class="text-xs font-semibold text-gray-500 uppercase leading-none">{{ event.monthShort }}</p>
                            <p class="text-2xl font-bold text-gray-700 leading-none mt-0.5">{{ event.dayNum }}</p>
                        </div>
                    </div>
                    <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-0.5">
                            <p class="font-semibold text-gray-800 truncate">{{ event.title }}</p>
                            <span v-if="event.dbStatus === 'pending'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-600">En attente</span>
                            <span v-else-if="event.dbStatus === 'approved'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700">Validé</span>
                            <span v-else-if="event.dbStatus === 'rejected'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-600">Refusé</span>
                        </div>
                        <div v-if="event.dbStatus === 'rejected' && event.rejectionReason" class="text-xs text-red-500 mt-0.5">{{ event.rejectionReason }}</div>
                        <div class="flex items-center gap-3 text-xs text-gray-500">
                            <span v-if="event.location !== '-'">{{ event.location }}</span>
                            <span v-if="event.time !== '-'">{{ event.time }}</span>
                            <span v-if="event.capacity">{{ event.registered }}/{{ event.capacity }} places</span>
                        </div>
                    </div>
                    <button
                        @click="confirmDeleteOwn(event)"
                        class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        title="Supprimer">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        type="text"
                        placeholder="Rechercher un événement…"
                        class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                    />
                </div>
                <select
                    v-model="filterStatus"
                    class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                    <option value="">Tous les événements</option>
                    <option value="upcoming">À venir</option>
                    <option value="registered">Mes inscriptions</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">
                    {{ filteredEvents.length }} résultat(s)
                </span>
            </div>

            <div class="divide-y divide-gray-50">
                <div
                    v-for="event in filteredEvents"
                    :key="event.id"
                    class="flex items-center gap-4 px-6 py-4 hover:bg-gray-50/50 transition-colors">

                    <div class="flex-shrink-0 w-14 text-center">
                        <div class="bg-primary/10 rounded-xl px-2 py-2">
                            <p class="text-xs font-semibold text-primary uppercase leading-none">{{ event.monthShort }}</p>
                            <p class="text-2xl font-bold text-primary leading-none mt-0.5">{{ event.dayNum }}</p>
                        </div>
                    </div>

                    <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-0.5">
                            <p class="font-semibold text-gray-800 truncate">{{ event.title }}</p>
                            <span v-if="event.isRegistered" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-secondary/15 text-secondary">
                                Inscrit
                            </span>
                            <span v-else-if="event.isPast" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500">
                                Passé
                            </span>
                        </div>
                        <div class="flex items-center gap-3 text-xs text-gray-500">
                            <span v-if="event.location !== '-'">{{ event.location }}</span>
                            <span v-if="event.time !== '-'">{{ event.time }}</span>
                            <span v-if="event.capacity">{{ event.registered }}/{{ event.capacity }} places</span>
                        </div>
                    </div>

                    <div class="flex-shrink-0 flex items-center gap-2">
                        <button
                            v-if="canDelete(event)"
                            @click.stop="confirmDelete(event)"
                            class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                            title="Supprimer"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </button>
                        <button
                            @click="openDetail(event)"
                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                            </svg>
                        </button>
                        <button
                            v-if="!event.isPast && event.isRegistered"
                            @click="toggleRegistration(event)"
                            :disabled="event.loading"
                            class="px-3 py-1.5 text-xs font-medium rounded-lg border border-gray-200 text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50">
                            Se désinscrire
                        </button>
                        <button
                            v-else-if="!event.isPast && !event.isRegistered"
                            @click="toggleRegistration(event)"
                            :disabled="event.loading || isFull(event)"
                            class="px-3 py-1.5 text-xs font-semibold rounded-lg bg-secondary text-white hover:bg-secondary-dark transition-colors disabled:opacity-50">
                            {{ isFull(event) ? 'Complet' : "S'inscrire" }}
                        </button>
                    </div>
                </div>

                <div v-if="filteredEvents.length === 0" class="px-6 py-16 text-center text-gray-400 text-sm">
                    Aucun événement disponible pour le moment.
                </div>
            </div>
            <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />
        </div>

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
                            <p class="text-xs text-gray-400 mb-0.5">Tarif</p>
                            <p class="font-medium text-gray-800">{{ detailEvent.price === 0 ? 'Gratuit' : (detailEvent.price / 100).toFixed(2) + ' €' }}</p>
                        </div>
                    </div>
                    <div>
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailEvent.description }}</p>
                    </div>
                    <div v-if="detailEvent.isRegistered" class="flex items-center gap-2 p-3 rounded-xl bg-secondary/10">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        <p class="text-sm text-secondary font-medium">Vous êtes inscrit à cet événement</p>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button
                        v-if="!detailEvent.isPast && detailEvent.isRegistered"
                        @click="toggleRegistration(detailEvent); detailEvent = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Se désinscrire
                    </button>
                    <button
                        v-else-if="!detailEvent.isPast && !detailEvent.isRegistered && !isFull(detailEvent)"
                        @click="toggleRegistration(detailEvent); detailEvent = null"
                        class="px-3 py-1.5 text-sm font-semibold bg-secondary text-white rounded-lg hover:bg-secondary-dark transition-colors">
                        S'inscrire
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
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteEvent" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="createModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="createModal = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        Créer un événement
                    </h3>
                    <button @click="createModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input v-model="eventForm.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre de l'événement" />
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea v-model="eventForm.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de l'événement…" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date & heure <span class="text-red-400">*</span></label>
                            <input v-model="eventForm.date" type="datetime-local" :min="minDateTime"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Lieu</label>
                            <input v-model="eventForm.location" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse ou ville" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Capacité</label>
                            <input v-model.number="eventForm.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Prix (€)</label>
                            <input v-model.number="eventForm.priceEuros" type="number" min="0" step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = gratuit" />
                        </div>
                    </div>
                    <p v-if="createError" class="text-xs text-red-500">{{ createError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button @click="createModal = false"
                        class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitCreate" :disabled="creating"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ creating ? 'Création…' : 'Créer' }}
                    </button>
                </div>
            </div>
        </div>

    </UserLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'
import frLocale from '@fullcalendar/core/locales/fr'
import UserLayout from '@/Layouts/UserLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'
import { useAuthStore } from '@/stores/auth.js'

const auth = useAuthStore()
const events = ref([])
const myEvents = ref([])
const calendarEvents = ref([])
const page = ref(1)
const total = ref(0)
const search = ref('')
const filterStatus = ref('')
const detailEvent = ref(null)
const toDelete = ref(null)
const deleting = ref(false)
const createModal = ref(false)
const creating = ref(false)
const createError = ref('')
const eventForm = ref({ title: '', description: '', date: '', location: '', capacity: null, priceEuros: 0 })

const canCreate = computed(() => auth.hasPermission('create_event'))

const minDateTime = computed(() => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 1)
    return now.toISOString().slice(0, 16)
})

const MONTHS_SHORT = ['jan','fév','mar','avr','mai','jun','jul','aoû','sep','oct','nov','déc']

const today = new Date()

function mapEvent(e) {
    return {
        id: e.id,
        title: e.title,
        creatorId: e.id_creator ?? null,
        organizer: e.creator_name ?? '-',
        location: e.location ?? '-',
        date: e.date?.slice(0, 10) ?? '-',
        time: e.date?.slice(11, 16) ?? '-',
        dateStr: e.date?.slice(0, 10) ?? null,
        dayNum: e.date?.slice(8, 10) ?? '-',
        monthShort: e.date ? MONTHS_SHORT[parseInt(e.date.slice(5, 7)) - 1] : '-',
        capacity: e.capacity ?? 0,
        registered: e.inscription_count ?? 0,
        price: e.price_cents ?? 0,
        isRegistered: e.is_registered ?? false,
        isPast: e.date ? e.date.slice(0, 10) < today.toISOString().slice(0, 10) : false,
        description: e.description ?? '-',
        dbStatus: e.status ?? 'approved',
        rejectionReason: e.rejection_reason ?? null,
        loading: false,
    }
}

const calendarOptions = computed(() => ({
    plugins: [dayGridPlugin, interactionPlugin],
    initialView: 'dayGridMonth',
    locale: frLocale,
    headerToolbar: {
        left: 'prev,next today',
        center: 'title',
        right: '',
    },
    events: calendarEvents.value
        .filter(e => e.dateStr)
        .map(e => ({
            id: e.id,
            title: e.title,
            date: e.dateStr,
            backgroundColor: '#8dc734',
            borderColor: '#6fa028',
        })),
    eventClick: (info) => {
        const event = calendarEvents.value.find(e => e.id === info.event.id)
            || events.value.find(e => e.id === info.event.id)
        if (event) openDetail(event)
    },
    height: 'auto',
}))

const filteredEvents = computed(() => {
    return events.value.filter(e => {
        const q = search.value.toLowerCase()
        if (q && !e.title.toLowerCase().includes(q) && !e.location.toLowerCase().includes(q)) return false
        if (filterStatus.value === 'upcoming' && e.isPast) return false
        if (filterStatus.value === 'registered' && !e.isRegistered) return false
        return true
    })
})

function isFull(event) {
    return event.capacity > 0 && event.registered >= event.capacity
}

function canDelete(event) {
    if (event.creatorId === auth.user?.id) return true
    return auth.isAdmin
}

function confirmDelete(event) {
    toDelete.value = event
}

function confirmDeleteOwn(event) {
    toDelete.value = event
}

async function deleteEvent() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/event/${toDelete.value.id}`)
        events.value = events.value.filter(e => e.id !== toDelete.value.id)
        myEvents.value = myEvents.value.filter(e => e.id !== toDelete.value.id)
        calendarEvents.value = calendarEvents.value.filter(e => e.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function openDetail(event) {
    detailEvent.value = event
}

async function toggleRegistration(event) {
    event.loading = true
    try {
        if (event.isRegistered) {
            await api.delete(`/user/event/${event.id}/unregister`)
            event.isRegistered = false
            event.registered = Math.max(0, event.registered - 1)
            calendarEvents.value = calendarEvents.value.filter(e => e.id !== event.id)
        } else {
            await api.post(`/user/event/${event.id}/register`)
            event.isRegistered = true
            event.registered++
            if (!calendarEvents.value.find(e => e.id === event.id)) {
                calendarEvents.value.push({ ...event, isRegistered: true })
            }
        }
    } catch (err) {
        console.error('toggleRegistration error:', err)
    } finally {
        event.loading = false
    }
}

function openCreate() {
    eventForm.value = { title: '', description: '', date: '', location: '', capacity: null, priceEuros: 0 }
    createError.value = ''
    createModal.value = true
}

async function submitCreate() {
    if (!eventForm.value.title.trim()) {
        createError.value = 'Le titre est requis.'
        return
    }
    if (!eventForm.value.date) {
        createError.value = 'La date est requise.'
        return
    }
    if (new Date(eventForm.value.date) <= new Date()) {
        createError.value = 'La date doit être dans le futur.'
        return
    }
    creating.value = true
    createError.value = ''
    try {
        const payload = {
            title: eventForm.value.title,
            description: eventForm.value.description || null,
            date: eventForm.value.date || null,
            location: eventForm.value.location || null,
            capacity: eventForm.value.capacity || null,
            price_cents: Math.round((eventForm.value.priceEuros || 0) * 100),
        }
        await api.post('/events', payload)
        page.value = 1
        await fetchEvents()
        if (auth.hasPermission('create_event')) {
            api.get('/user/my-events').then(({ data }) => {
                myEvents.value = (data ?? []).map(mapEvent)
            }).catch(() => {})
        }
        createModal.value = false
    } catch (e) {
        createError.value = e.response?.data?.error ?? 'Erreur lors de la création.'
    } finally {
        creating.value = false
    }
}

async function fetchEvents() {
    try {
        const { data } = await api.get('/events', { params: { page: page.value, limit: 15 } })
        events.value = (data.data ?? []).map(mapEvent)
        total.value = data.total ?? 0
    } catch (e) {
        console.error('Events fetch error:', e)
    }
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(page, fetchEvents)

onMounted(async () => {
    const promises = [
        fetchEvents(),
        api.get('/user/events').then(({ data }) => {
            calendarEvents.value = (data ?? []).map(mapEvent)
        }).catch(() => {}),
    ]
    if (auth.hasPermission('create_event')) {
        promises.push(
            api.get('/user/my-events').then(({ data }) => {
                myEvents.value = (data ?? []).map(mapEvent)
            }).catch(() => {})
        )
    }
    await Promise.all(promises)
})
</script>
