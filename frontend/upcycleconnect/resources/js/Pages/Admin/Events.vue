<template>
  <AdminLayout title="Événements">

    <div class="grid grid-cols-4 gap-5 mb-8">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
        <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
          <div :class="stat.iconClass" v-html="stat.icon" />
        </div>
        <div class="min-w-0">
          <p class="text-2xl font-bold text-gray-800 leading-none">
            {{stat.value}}
          </p>
          <p class="text-sm text-gray-500 mt-1 truncate">{{stat.label}}</p>
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
          <option>À venir</option>
          <option>En cours</option>
          <option>Passé</option>
          <option>Annulé</option>
        </select>

        <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">
          {{filteredEvents.length}} résultat(s)
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
              v-for="(event, i) in filteredEvents"
              :key="event.id"
              :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
              <td class="px-5 py-3">
                <p class="font-medium text-gray-800 truncate max-w-52">{{event.title}}</p>
                <p class="text-xs text-gray-500 truncate">{{event.organizer}}</p>
              </td>

              <td class="px-5 py-3">
                <p class="text-gray-800 text-xs font-medium">{{event.date}} · {{event.time}}</p>
                <p class="text-gray-500 text-xs truncate">{{event.location}}</p>
              </td>

              <td class="px-5 py-3 text-xs text-gray-700 font-medium">
                {{event.registered}}/{{event.capacity}}
              </td>

              <td class="px-5 py-3">
                <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                  <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(event.status)" />
                  <span :class="statusText(event.status)">{{event.status}}</span>
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

                  <button
                    @click="deleteEvent(event)"
                    title="Supprimer l'événement"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredEvents.length === 0">
              <td colspan="5" class="px-5 py-12 text-center text-gray-400 text-sm">
                Aucun événement ne correspond à vos filtres.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
        {{filteredEvents.length}} événement(s) au total
      </div>
    </div>

    <div
      v-if="detailEvent"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="detailEvent = null">
      <div class="absolute inset-0 bg-black/40" @click="detailEvent = null" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                {{detailEvent.title}}
              </h3>
            </div>
            <p class="text-xs text-gray-500 mt-0.5">Organisé par {{detailEvent.organizer}}</p>
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
              <p class="font-medium text-gray-800">{{detailEvent.date}}</p>
              <p class="text-xs text-gray-500">{{detailEvent.time}}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
              <p class="font-medium text-gray-800">{{detailEvent.location}}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Inscriptions</p>
              <p class="font-medium text-gray-800">{{detailEvent.registered}} / {{detailEvent.capacity}} participants</p>
            </div>
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Statut</p>
              <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailEvent.status)" />
                <span :class="statusText(detailEvent.status)">{{detailEvent.status}}</span>
              </span>
            </div>
          </div>

          <div>
            <p class="text-xs text-gray-400 mb-1">Description</p>
            <p class="text-sm text-gray-700 leading-relaxed">{{detailEvent.description}}</p>
          </div>
        </div>

        <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
          <button
            @click="deleteEvent(detailEvent); detailEvent = null"
            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
            Supprimer
          </button>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const events = ref([])

function computeStatus(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr) > new Date() ? 'À venir' : 'Passé'
}

function mapEvent(e) {
  return {
    id: e.id,
    title: e.title,
    organizer: e.creator_name ?? '—',
    location: e.location ?? '—',
    date: e.date?.slice(0, 10) ?? '—',
    time: e.date?.slice(11, 16) ?? '—',
    capacity: e.capacity ?? 0,
    registered: e.inscription_count ?? 0,
    status: computeStatus(e.date),
    description: e.description ?? '—',
  }
}

onMounted(async () => {
  try {
    const { data } = await api.get('/admin/events')
    events.value = data.map(mapEvent)
  } catch (e) {
    console.error('Events fetch error:', e)
  }
})

const stats = computed(() => [
  {
    label: 'Total événements',
    value: events.value.length,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
  },
  {
    label: 'À venir',
    value: events.value.filter(e => e.status === 'À venir').length,
    bgClass: 'bg-green-100',
    iconClass: 'text-green-600',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
  },
  {
    label: 'Passés',
    value: events.value.filter(e => e.status === 'Passé').length,
    bgClass: 'bg-amber-100',
    iconClass: 'text-amber-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
  },
  {
    label: 'Sans date',
    value: events.value.filter(e => e.status === '—').length,
    bgClass: 'bg-gray-100',
    iconClass: 'text-gray-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/></svg>`,
  },
])

const search = ref('')
const filterStatus = ref('')

const filteredEvents = computed(() => {
  return events.value.filter(e => {
    const q = search.value.toLowerCase()
    if (q && !e.title.toLowerCase().includes(q) && !e.organizer.toLowerCase().includes(q)) return false
    if (filterStatus.value && e.status !== filterStatus.value) return false
    return true
  })
})

const detailEvent = ref(null)

function openDetail(event) {
  detailEvent.value = event
}

async function deleteEvent(event) {
  if (!confirm(`Supprimer définitivement "${event.title}" ?`)) return
  try {
    await api.delete(`/admin/event/${event.id}`)
    events.value = events.value.filter(e => e.id !== event.id)
    if (detailEvent.value?.id === event.id) detailEvent.value = null
  } catch (e) {
    console.error('deleteEvent error:', e)
  }
}

function statusDot(status) {
  const map = {
    'À venir': 'bg-green-500',
    'En cours': 'bg-blue-500',
    'Passé': 'bg-gray-400',
    'Annulé': 'bg-red-400',
 }
  return map[status] ?? 'bg-gray-300'
}

function statusText(status) {
  const map = {
    'À venir': 'text-green-700',
    'En cours': 'text-blue-700',
    'Passé': 'text-gray-500',
    'Annulé': 'text-red-600',
 }
  return map[status] ?? 'text-gray-500'
}
</script>
