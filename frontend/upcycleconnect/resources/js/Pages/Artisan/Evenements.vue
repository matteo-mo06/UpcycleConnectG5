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

        <div v-if="!loading" class="flex items-center gap-3 mb-4">
            <select
                v-model="filterStatus"
                class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600 bg-white"
            >
                <option value="">Tous les statuts</option>
                <option value="pending">En attente</option>
                <option value="approved">Validé</option>
                <option value="rejected">Refusé</option>
            </select>
            <span class="text-xs text-gray-400">{{ filteredEvents.length }} événement(s)</span>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">Chargement…</div>

        <div v-else-if="filteredEvents.length === 0" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">
            {{ events.length === 0 ? 'Vous n\'avez créé aucun événement.' : 'Aucun événement pour ce filtre.' }}
        </div>

        <div v-else class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="divide-y divide-gray-50">
                <div
                    v-for="event in filteredEvents"
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
                            <span v-else-if="event.status === 'approved'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700">Validé</span>
                            <span v-else-if="event.status === 'rejected'" class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-600">Refusé</span>
                        </div>
                        <div v-if="event.status === 'rejected' && event.rejection_reason" class="text-xs text-red-500 mt-0.5">{{ event.rejection_reason }}</div>
                        <div class="flex items-center gap-3 text-xs text-gray-500">
                            <span v-if="event.location">{{ event.location }}</span>
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
        </div>

        <div v-if="selectedEvent" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="selectedEvent = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selectedEvent.title }}</h3>
                    <button @click="selectedEvent = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto flex-1">
                    <div class="flex gap-2 flex-wrap">
                        <span v-if="selectedEvent.status === 'pending'" class="px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-600">En attente de validation</span>
                        <span v-else-if="selectedEvent.status === 'approved'" class="px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700">Validé</span>
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
                        <div v-if="selectedEvent.location"><span class="font-medium text-gray-700">Lieu :</span> {{ selectedEvent.location }}</div>
                        <div v-if="selectedEvent.capacity"><span class="font-medium text-gray-700">Places :</span> {{ selectedEvent.registered }}/{{ selectedEvent.capacity }}</div>
                        <div>
                            <span class="font-medium text-gray-700">Tarif :</span>
                            {{ selectedEvent._raw?.price_cents ? (selectedEvent._raw.price_cents / 100).toFixed(2) + ' €' : 'Gratuit' }}
                        </div>
                    </div>
                    <p v-if="selectedEvent._raw?.description" class="leading-relaxed text-gray-600">{{ selectedEvent._raw.description }}</p>
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
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ editTarget ? 'Modifier l\'événement' : 'Créer un événement' }}
                    </h3>
                    <button @click="formModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input v-model="form.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre de l'événement" />
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea v-model="form.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de l'événement…" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date & heure <span class="text-red-400">*</span></label>
                            <input v-model="form.date" type="datetime-local" :min="minDateTime"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Lieu</label>
                            <input v-model="form.location" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse ou ville" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Capacité</label>
                            <input v-model.number="form.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Prix (€)</label>
                            <input v-model.number="form.priceEuros" type="number" min="0" step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = gratuit" />
                        </div>
                    </div>
                    <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button @click="formModal = false" class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitForm" :disabled="submitting"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ submitting ? 'Enregistrement…' : (editTarget ? 'Mettre à jour' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

    </ArtisanLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import api from '@/api.js'

const route = useRoute()

const events = ref([])
const loading = ref(true)
const toDelete = ref(null)
const deleting = ref(false)
const selectedEvent = ref(null)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', date: '', location: '', capacity: null, priceEuros: 0 })
const filterStatus = ref('')
const filteredEvents = computed(() =>
    filterStatus.value ? events.value.filter(e => e.status === filterStatus.value) : events.value
)

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
        location: e.location ?? null,
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

async function fetchEvents() {
    loading.value = true
    try {
        const { data } = await api.get('/user/my-events')
        events.value = (data ?? []).map(mapEvent)
    } catch {
        events.value = []
    } finally {
        loading.value = false
    }
}

async function deleteEvent() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/event/${toDelete.value.id}`)
        events.value = events.value.filter(e => e.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', date: '', location: '', capacity: null, priceEuros: 0 }
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
        location: event.location ?? '',
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
            location: form.value.location || null,
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
