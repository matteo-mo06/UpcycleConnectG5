<template>
    <SalarieLayout>

        <div class="mb-4">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Calendrier
            </h1>
            <p class="text-sm text-gray-500 mt-1">Vos événements et formations créés</p>
        </div>

        <div class="flex items-center gap-6 mb-4">
            <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full flex-shrink-0" style="background-color: #c46d68"></span>
                <span class="text-xs text-gray-500">Événements</span>
            </div>
            <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full flex-shrink-0" style="background-color: #8dc734"></span>
                <span class="text-xs text-gray-500">Formations</span>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6">
            <FullCalendar :options="calendarOptions" />

            <div v-if="selectedItem" class="mt-4 pt-4 border-t border-gray-100 flex items-start gap-3">
                <span
                    class="mt-1 px-2.5 py-1 rounded-full text-xs font-medium text-white flex-shrink-0"
                    :style="{ backgroundColor: TYPE_COLORS[selectedItem.type] }">
                    {{ TYPE_LABELS[selectedItem.type] }}
                </span>
                <div class="flex-1 min-w-0">
                    <p class="text-base font-semibold text-gray-800 truncate">{{ selectedItem.title }}</p>
                    <p v-if="selectedMeta.length" class="text-sm text-gray-500 mt-1">{{ selectedMeta.join(' · ') }}</p>
                    <p v-if="selectedItem.description" class="text-sm text-gray-400 mt-1 leading-snug">{{ selectedItem.description }}</p>
                </div>
                <div class="flex-shrink-0 flex items-center gap-2">
                    <button
                        @click="viewItem"
                        class="px-3 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors">
                        Voir
                    </button>
                    <button
                        @click="selectedItem = null"
                        class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
            </div>
        </div>

    </SalarieLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'
import frLocale from '@fullcalendar/core/locales/fr'
import SalarieLayout from '@/Layouts/SalarieLayout.vue'
import api from '@/api.js'

const router = useRouter()

const events = ref([])
const formations = ref([])
const selectedItem = ref(null)

const VIEW_ROUTES = {
    event:     '/salarie/evenements',
    formation: '/salarie/formations',
}

function viewItem() {
    if (!selectedItem.value) return
    router.push(`${VIEW_ROUTES[selectedItem.value.type]}?open=${selectedItem.value.id}`)
}

const TYPE_COLORS = {
    event:     '#c46d68',
    formation: '#8dc734',
}

const TYPE_LABELS = {
    event:     'Événement',
    formation: 'Formation',
}

const LEVEL_LABELS = {
    beginner:     'Débutant',
    intermediate: 'Intermédiaire',
    advanced:     'Avancé',
}

const MONTHS_SHORT = ['jan.','fév.','mar.','avr.','mai','jun.','jul.','aoû.','sep.','oct.','nov.','déc.']

function formatDate(dateStr) {
    if (!dateStr) return null
    const [year, month, day] = dateStr.slice(0, 10).split('-')
    return `${parseInt(day)} ${MONTHS_SHORT[parseInt(month) - 1]} ${year}`
}

function formatTime(dateStr) {
    if (!dateStr || dateStr.length < 16) return null
    return dateStr.slice(11, 16)
}

const selectedMeta = computed(() => {
    if (!selectedItem.value) return []
    const s = selectedItem.value
    const parts = []

    if (s.type === 'event') {
        const date = formatDate(s.date)
        const time = formatTime(s.date)
        if (date && time) parts.push(`${date} · ${time}`)
        else if (date) parts.push(date)
        if (s.location) parts.push(s.location)
    } else if (s.type === 'formation') {
        const date = formatDate(s.date)
        if (date) parts.push(date)
        if (s.level) parts.push(LEVEL_LABELS[s.level] ?? s.level)
        if (s.location) parts.push(s.location)
    }

    return parts
})

function composeAddr(o) {
    const cityLine = [o.postal, o.city].filter(Boolean).join(' ')
    return [o.address, cityLine].filter(Boolean).join(', ') || null
}

const calendarOptions = computed(() => {
    const eventEntries = events.value
        .filter(e => e.date)
        .map(e => ({
            id: `event-${e.id}`,
            title: e.title,
            date: e.date.slice(0, 10),
            backgroundColor: TYPE_COLORS.event,
            borderColor: '#a85a55',
            textColor: '#ffffff',
            extendedProps: { type: 'event', date: e.date, location: composeAddr(e), description: e.description ?? null },
        }))

    const formationEntries = formations.value
        .filter(f => f.date)
        .map(f => ({
            id: `formation-${f.id}`,
            title: f.title,
            date: f.date.slice(0, 10),
            backgroundColor: TYPE_COLORS.formation,
            borderColor: '#6fa028',
            textColor: '#ffffff',
            extendedProps: { type: 'formation', date: f.date, level: f.level ?? null, location: composeAddr(f), description: f.description ?? null },
        }))

    return {
        plugins: [dayGridPlugin, interactionPlugin],
        initialView: 'dayGridMonth',
        locale: frLocale,
        headerToolbar: {
            left: 'prev,next today',
            center: 'title',
            right: '',
        },
        events: [...eventEntries, ...formationEntries],
        eventContent: (arg) => {
            const ep = arg.event.extendedProps
            const t = arg.event.title.replace(/</g, '&lt;').replace(/>/g, '&gt;')
            const raw = ep.description?.replace(/</g, '&lt;').replace(/>/g, '&gt;')
            const m = raw ? `<div class="fc-custom-meta">${raw.slice(0, 55)}</div>` : ''
            return { html: `<div class="fc-custom-event"><div class="fc-custom-title">${t}</div>${m}</div>` }
        },
        eventClick: (info) => {
            const ep = info.event.extendedProps
            const rawId = info.event.id
            const id = rawId.replace(/^(event|formation)-/, '')
            selectedItem.value = {
                id,
                title:       info.event.title,
                type:        ep.type,
                date:        ep.date     ?? null,
                location:    ep.location ?? null,
                level:       ep.level    ?? null,
                description: ep.description ?? null,
            }
        },
        height: 'auto',
    }
})

onMounted(async () => {
    const [evRes, foRes] = await Promise.all([
        api.get('/user/my-events').catch(() => ({ data: [] })),
        api.get('/user/my-formations').catch(() => ({ data: [] })),
    ])
    events.value     = evRes.data ?? []
    formations.value = foRes.data ?? []
})
</script>
