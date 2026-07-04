<template>
    <UserLayout>

        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                    Événements
                </h1>
                <p class="text-sm text-gray-500 mt-1">Découvrez et inscrivez-vous aux événements disponibles</p>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6 mb-6">
            <h2 class="font-semibold text-gray-800 mb-4" style="font-family: var(--font-family-title)">
                Mes événements
            </h2>
            <FullCalendar :options="calendarOptions" />
        </div>

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        @input="onSearchInput"
                        type="text"
                        placeholder="Rechercher un événement…"
                        class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                    />
                </div>
                <select
                    v-model="filterStatus"
                    @change="resetAndFetch"
                    class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                    <option value="">Tous les événements</option>
                    <option value="upcoming">À venir</option>
                    <option value="registered">Mes inscriptions</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">
                    {{ total }} résultat(s)
                </span>
            </div>

            <div class="divide-y divide-gray-50">
                <div
                    v-for="event in events"
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

                <div v-if="events.length === 0" class="px-6 py-16 text-center text-gray-400 text-sm">
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

                    <div v-if="auth.user?.id === detailEvent?.id_creator || auth.isAdmin" class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Participants</p>
                        <div class="space-y-1">
                            <p v-if="participantsLoading" class="text-xs text-gray-400">Chargement…</p>
                            <template v-else>
                                <p v-if="participants.length === 0" class="text-xs text-gray-400">Aucun participant</p>
                                <div v-for="p in participants" :key="p.id" class="flex items-center gap-2 py-1">
                                    <img v-if="p.avatar_url" :src="p.avatar_url" class="w-6 h-6 rounded-full object-cover" />
                                    <span v-else class="w-6 h-6 rounded-full bg-gray-200 text-xs flex items-center justify-center font-medium text-gray-500">
                                        {{ p.first_name[0] }}{{ p.last_name[0] }}
                                    </span>
                                    <span class="text-sm text-gray-700">{{ p.first_name }} {{ p.last_name }}</span>
                                </div>
                            </template>
                        </div>
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

    </UserLayout>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
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
const calendarEvents = ref([])
const page = ref(1)
const total = ref(0)
const search = ref('')
const filterStatus = ref('')
const detailEvent = ref(null)
const participants = ref([])
const participantsLoading = ref(false)

const MONTHS_SHORT = ['jan','fév','mar','avr','mai','jun','jul','aoû','sep','oct','nov','déc']

const today = new Date()

function mapEvent(e) {
    return {
        id: e.id,
        title: e.title,
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
        loading: false,
        id_creator: e.id_creator ?? null,
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

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(resetAndFetch, 300)
}

function resetAndFetch() {
    page.value = 1
    fetchEvents()
}

function isFull(event) {
    return event.capacity > 0 && event.registered >= event.capacity
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
        alert(err?.response?.data?.error ?? 'Erreur lors de l\'inscription / désinscription.')
    } finally {
        event.loading = false
    }
}

async function fetchEvents() {
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        if (filterStatus.value) params.status = filterStatus.value
        const { data } = await api.get('/events', { params })
        events.value = (data.data ?? []).map(mapEvent)
        total.value = data.total ?? 0
    } catch (e) {
        console.error('fetchEvents error:', e)
    }
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(page, fetchEvents)

watch(detailEvent, (val) => {
    if (!val) {
        participants.value = []
        participantsLoading.value = false
    } else {
        loadParticipants(`/events/${val.id}/participants`)
    }
})

async function loadParticipants(url) {
    participantsLoading.value = true
    try {
        const { data } = await api.get(url)
        participants.value = data ?? []
    } catch {
        participants.value = []
    } finally {
        participantsLoading.value = false
    }
}

async function fetchAll() {
    await Promise.all([
        fetchEvents(),
        api.get('/user/events').then(({ data }) => {
            calendarEvents.value = (data ?? []).map(mapEvent)
        }).catch(() => {}),
    ])
}

usePolling(fetchAll)
</script>
