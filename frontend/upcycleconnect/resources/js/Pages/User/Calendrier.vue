<template>
    <UserLayout>

        <div class="mb-4">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Calendrier
            </h1>
            <p class="text-sm text-gray-500 mt-1">Vos événements, formations et projets en un coup d'œil</p>
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
            <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full flex-shrink-0" style="background-color: #2D2D2D"></span>
                <span class="text-xs text-gray-500">Projets</span>
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
                <button
                    @click="selectedItem = null"
                    class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                    <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                </button>
            </div>
        </div>

    </UserLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'
import frLocale from '@fullcalendar/core/locales/fr'
import UserLayout from '@/Layouts/UserLayout.vue'
import api from '@/api.js'

const events = ref([])
const formations = ref([])
const projects = ref([])
const selectedItem = ref(null)

const TYPE_COLORS = {
    event:     '#c46d68',
    formation: '#8dc734',
    project:   '#2D2D2D',
}

const TYPE_LABELS = {
    event:     'Événement',
    formation: 'Formation',
    project:   'Projet',
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
    } else if (s.type === 'project') {
        const start = formatDate(s.start_date)
        const end = formatDate(s.end_date)
        if (start && end) parts.push(`${start} → ${end}`)
        else if (start) parts.push(`Dès le ${start}`)
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

    const projectEntries = projects.value
        .filter(p => p.start_date)
        .map(p => ({
            id: `project-${p.id}`,
            title: p.title,
            start: p.start_date.slice(0, 10),
            end: p.end_date ? p.end_date.slice(0, 10) : undefined,
            backgroundColor: TYPE_COLORS.project,
            borderColor: '#1a1a1a',
            textColor: '#ffffff',
            extendedProps: { type: 'project', start_date: p.start_date, end_date: p.end_date ?? null, location: p.location ?? null, description: p.description ?? null },
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
        events: [...eventEntries, ...formationEntries, ...projectEntries],
        eventContent: (arg) => {
            const ep = arg.event.extendedProps
            const t = arg.event.title.replace(/</g, '&lt;').replace(/>/g, '&gt;')
            const raw = ep.description?.replace(/</g, '&lt;').replace(/>/g, '&gt;')
            const m = raw ? `<div class="fc-custom-meta">${raw.slice(0, 55)}</div>` : ''
            return { html: `<div class="fc-custom-event"><div class="fc-custom-title">${t}</div>${m}</div>` }
        },
        eventClick: (info) => {
            const ep = info.event.extendedProps
            selectedItem.value = {
                title:       info.event.title,
                type:        ep.type,
                date:        ep.date        ?? null,
                location:    ep.location    ?? null,
                level:       ep.level       ?? null,
                start_date:  ep.start_date  ?? null,
                end_date:    ep.end_date    ?? null,
                description: ep.description ?? null,
            }
        },
        height: 'auto',
    }
})

onMounted(async () => {
    const [evRes, foRes, prRes] = await Promise.all([
        api.get('/user/events').catch(() => ({ data: [] })),
        api.get('/user/formations').catch(() => ({ data: [] })),
        api.get('/user/projects').catch(() => ({ data: [] })),
    ])
    events.value     = evRes.data ?? []
    formations.value = foRes.data ?? []
    projects.value   = prRes.data ?? []
})
</script>
