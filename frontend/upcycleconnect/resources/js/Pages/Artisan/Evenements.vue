<template>
    <ArtisanLayout>

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Événements</h1>
            <button
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un événement
            </button>
        </div>

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex flex-wrap items-center gap-3">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        @input="onSearchInput"
                        type="text"
                        placeholder="Rechercher un événement…"
                        class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                </div>
                <select
                    v-model="filterStatus"
                    @change="resetAndFetch"
                    class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
                >
                    <option value="">Tous les statuts</option>
                    <option value="pending">En attente</option>
                    <option value="approved">À venir</option>
                    <option value="rejected">Refusé</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} événement(s)</span>
            </div>

            <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">Chargement…</div>

            <div v-else-if="events.length === 0" class="py-12 text-center text-gray-400 text-sm">
                {{ filterStatus || search ? 'Aucun événement pour ce filtre.' : 'Vous n\'avez créé aucun événement.' }}
            </div>

            <div v-else class="divide-y divide-gray-50">
                <div
                    v-for="event in events"
                    :key="event.id"
                    @click="selectedEvent = event"
                    class="flex items-center gap-4 px-6 py-4 cursor-pointer hover:bg-gray-50/60 transition-colors"
                >
                    <div class="flex-shrink-0 w-14 text-center">
                        <div class="bg-primary/10 rounded-xl px-2 py-2">
                            <p class="text-xs font-semibold text-primary uppercase leading-none">{{ event.monthShort }}</p>
                            <p class="text-2xl font-bold text-primary leading-none mt-0.5">{{ event.dayNum }}</p>
                        </div>
                    </div>
                    <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-0.5">
                            <p class="font-semibold text-gray-800 truncate">{{ event.title }}</p>
                            <span v-if="event.status === 'pending'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-600">En attente</span>
                            <span v-else-if="event.status === 'approved'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700">À venir</span>
                            <span v-else-if="event.status === 'rejected'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-600">Refusé</span>
                        </div>
                        <div v-if="event.status === 'rejected' && event.rejection_reason" class="text-xs text-red-500 mt-0.5">{{ event.rejection_reason }}</div>
                        <div class="flex items-center gap-3 text-xs text-gray-500">
                            <span v-if="event.city || event.address">{{ event.city || event.address }}</span>
                            <span v-if="event.time">{{ event.time }}</span>
                            <span v-if="event.capacity">{{ event.registered }}/{{ event.capacity }} places</span>
                        </div>
                    </div>
                    <button
                        @click.stop="toDelete = event"
                        class="flex-shrink-0 p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        title="Supprimer"
                    >
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                </div>
            </div>

            <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />
        </div>

        <div v-if="selectedEvent" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="selectedEvent = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selectedEvent.title }}</h3>
                    <button @click="selectedEvent = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto flex-1">
                    <div class="flex gap-2 flex-wrap">
                        <span v-if="selectedEvent.status === 'pending'" class="px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-600">En attente de validation</span>
                        <span v-else-if="selectedEvent.status === 'approved'" class="px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700">À venir</span>
                        <span v-else-if="selectedEvent.status === 'rejected'" class="px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-600">Refusé</span>
                    </div>
                    <p v-if="selectedEvent.status === 'rejected' && selectedEvent.rejection_reason" class="text-xs text-red-500">
                        Raison : {{ selectedEvent.rejection_reason }}
                    </p>
                    <div class="grid grid-cols-2 gap-3 text-xs text-gray-500">
                        <div v-if="selectedEvent.date !== '-'">
                            <span class="font-medium text-gray-700">Date :</span>
                            {{ selectedEvent.dayNum }} {{ selectedEvent.monthShort }}
                            <span v-if="selectedEvent.time"> · {{ selectedEvent.time }}</span>
                        </div>
                        <div v-if="selectedEvent.address"><span class="font-medium text-gray-700">Adresse :</span> {{ selectedEvent.address }}</div>
                        <div v-if="selectedEvent.city"><span class="font-medium text-gray-700">Ville :</span> {{ selectedEvent.postal ? selectedEvent.postal + ' ' : '' }}{{ selectedEvent.city }}</div>
                        <div v-if="selectedEvent.capacity"><span class="font-medium text-gray-700">Places :</span> {{ selectedEvent.registered }}/{{ selectedEvent.capacity }}</div>
                        <div>
                            <span class="font-medium text-gray-700">Tarif :</span>
                            {{ selectedEvent._raw?.price_cents ? (selectedEvent._raw.price_cents / 100).toFixed(2) + ' €' : 'Gratuit' }}
                        </div>
                    </div>
                    <p v-if="selectedEvent._raw?.description" class="leading-relaxed text-gray-600">{{ selectedEvent._raw.description }}</p>

                    <div class="border-t border-gray-100 pt-4">
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
                <div class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button
                        v-if="selectedEvent.status === 'pending'"
                        @click="openEditFromDetail(selectedEvent)"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
                    >Modifier</button>
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

        <div v-if="formModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ editTarget ? 'Modifier l\'événement' : 'Créer un événement' }}
                    </h3>
                    <button @click="formModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input v-model="form.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre de l'événement" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea v-model="form.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de l'événement…" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-gray-700 mb-1">Adresse <span class="text-red-400">*</span></label>
                            <input v-model="form.address" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse de l'événement" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Ville</label>
                            <input v-model="form.city" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Paris" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Code postal</label>
                            <input v-model="form.postal" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="75001" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Date & heure <span class="text-red-400">*</span></label>
                            <input v-model="form.date" type="datetime-local" :min="minDateTime"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Capacité</label>
                            <input v-model.number="form.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée" />
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-gray-700 mb-1">Prix (€)</label>
                            <input v-model.number="form.priceEuros" type="number" min="0" step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = gratuit" />
                        </div>
                    </div>
                    <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2 flex-shrink-0">
                    <button @click="formModal = false" class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitForm" :disabled="submitting"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ submitting ? 'Enregistrement…' : (editTarget ? 'Enregistrer' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

    </ArtisanLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const route = useRoute()

const events = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(true)
const search = ref('')
const filterStatus = ref('')
const toDelete = ref(null)
const deleting = ref(false)
const selectedEvent = ref(null)
const participants = ref([])
const participantsLoading = ref(false)

watch(selectedEvent, (val) => {
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

const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', date: '', address: '', city: '', postal: '', capacity: null, priceEuros: 0 })

const MONTHS_SHORT = ['jan','fév','mar','avr','mai','jun','jul','aoû','sep','oct','nov','déc']

const minDateTime = computed(() => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 1)
    return now.toISOString().slice(0, 16)
})

function mapEvent(e) {
    return {
        id: e.id,
        title: e.title,
        address: e.address ?? null,
        city: e.city ?? null,
        postal: e.postal ?? null,
        date: e.date?.slice(0, 10) ?? '-',
        time: e.date?.slice(11, 16) ?? null,
        dayNum: e.date?.slice(8, 10) ?? '-',
        monthShort: e.date ? MONTHS_SHORT[parseInt(e.date.slice(5, 7)) - 1] : '-',
        capacity: e.capacity ?? 0,
        registered: e.inscription_count ?? 0,
        status: e.status ?? 'approved',
        rejection_reason: e.rejection_reason ?? null,
        _raw: e,
    }
}

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(resetAndFetch, 300)
}

function resetAndFetch() {
    page.value = 1
    fetchEvents()
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
    fetchEvents()
}

async function fetchEvents() {
    loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        if (filterStatus.value) params.status = filterStatus.value
        const { data } = await api.get('/user/my-events', { params })
        events.value = (data.data ?? []).map(mapEvent)
        total.value = data.total ?? 0
    } catch {
        events.value = []
        total.value = 0
    } finally {
        loading.value = false
    }
}

async function deleteEvent() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/event/${toDelete.value.id}`)
        toDelete.value = null
        await fetchEvents()
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', date: '', address: '', city: '', postal: '', capacity: null, priceEuros: 0 }
    formError.value = ''
    formModal.value = true
}

function openEditFromDetail(event) {
    selectedEvent.value = null
    editTarget.value = event
    form.value = {
        title: event.title,
        description: event._raw?.description ?? '',
        date: event._raw?.date ? event._raw.date.slice(0, 16) : '',
        address: event.address ?? '',
        city: event.city ?? '',
        postal: event.postal ?? '',
        capacity: event.capacity || null,
        priceEuros: event._raw?.price_cents ? event._raw.price_cents / 100 : 0,
    }
    formError.value = ''
    formModal.value = true
}

async function submitForm() {
    if (!form.value.title.trim()) { formError.value = 'Le titre est requis.'; return }
    if (!form.value.date) { formError.value = 'La date est requise.'; return }
    if (!editTarget.value && new Date(form.value.date) <= new Date()) { formError.value = 'La date doit être dans le futur.'; return }
    submitting.value = true
    formError.value = ''
    try {
        const payload = {
            title: form.value.title,
            description: form.value.description || null,
            date: form.value.date || null,
            address: form.value.address || null,
            city: form.value.city || null,
            postal: form.value.postal || null,
            capacity: form.value.capacity || null,
            price_cents: Math.round((form.value.priceEuros || 0) * 100),
        }
        if (editTarget.value) {
            await api.patch(`/user/event/${editTarget.value.id}`, payload)
        } else {
            await api.post('/events', payload)
        }
        formModal.value = false
        await fetchEvents()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await fetchEvents()
    const openId = route.query.open
    if (openId) {
        const found = events.value.find(e => String(e.id) === String(openId))
        if (found) selectedEvent.value = found
    }
})
</script>
