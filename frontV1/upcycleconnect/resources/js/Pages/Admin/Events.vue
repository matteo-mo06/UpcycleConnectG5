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

    <div v-if="pendingEvents.length > 0" class="bg-amber-50 border border-amber-200 rounded-xl mb-6 overflow-hidden">
      <div class="px-5 py-3.5 border-b border-amber-200 flex items-center gap-2">
        <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" class="w-4 h-4 text-amber-600">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <h2 class="text-sm font-semibold text-amber-800">
          {{pendingEvents.length}} événement(s) en attente de validation
        </h2>
      </div>
      <div class="divide-y divide-amber-100">
        <div
          v-for="event in pendingEvents"
          :key="event.id"
          class="px-5 py-4 flex items-start gap-4">
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <p class="font-medium text-gray-800 text-sm">{{event.title}}</p>
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="categoryBadge(event.category)">
                {{event.category}}
              </span>
            </div>
            <p class="text-xs text-gray-500 mt-0.5">
              Par <span class="font-medium text-gray-700">{{event.organizer}}</span>
              · {{event.location}} · {{event.time}}
            </p>
            <p class="text-xs text-gray-600 mt-1.5 italic line-clamp-2">"{{event.description}}"</p>
          </div>

          <div class="flex-shrink-0 flex gap-2">
            <button
              @click="approveEvent(event.id)"
              class="px-3 py-1.5 text-xs font-medium bg-secondary text-white rounded-lg hover:bg-secondary-dark transition-colors">
              Valider
            </button>
            <button
              @click="rejectEvent(event.id)"
              class="px-3 py-1.5 text-xs font-medium bg-white text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Refuser
            </button>
          </div>
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
          v-model="filterCategory"
          class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
          <option value="">Toutes les catégories</option>
          <option>Atelier</option>
          <option>Conférence</option>
          <option>Animation</option>
          <option>Formation</option>
        </select>

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
              <th class="text-left text-white font-medium px-5 py-3">Catégorie</th>
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
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="categoryBadge(event.category)">
                  {{event.category}}
                </span>
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
                    v-if="event.status === 'À venir' || event.status === 'En cours'"
                    @click="cancelEvent(event)"
                    title="Annuler l'événement"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
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
              <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
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
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="categoryBadge(detailEvent.category)">
                {{detailEvent.category}}
              </span>
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
            v-if="detailEvent.status === 'À venir' || detailEvent.status === 'En cours'"
            @click="cancelEvent(detailEvent); detailEvent = null"
            class="px-3 py-1.5 text-sm font-medium text-amber-600 border border-amber-200 rounded-lg hover:bg-amber-50 transition-colors">
            Annuler l'événement
          </button>
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
import {ref, computed} from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'

const pendingEvents = ref([
  {id: 101, title: 'Événement test en attente', category: 'Catégorie', organizer: 'Organisateur', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, description: 'Description de l\'événement en attente de validation.'},
  {id: 102, title: 'Événement test en attente', category: 'Catégorie', organizer: 'Organisateur', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, description: 'Description de l\'événement en attente de validation.'},
])

function approveEvent(id) {
  const idx = pendingEvents.value.findIndex(e => e.id === id)
  const event = pendingEvents.value[idx]
  events.value.unshift({...event, status: 'À venir', registered: 0})
  pendingEvents.value.splice(idx, 1)
}

function rejectEvent(id) {
  pendingEvents.value = pendingEvents.value.filter(e => e.id !== id)
}

const events = ref([
  {id: 1, title: 'Événement test', organizer: 'Organisateur', category: 'Catégorie', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, registered: 10, status: 'À venir', description: 'Description de l\'événement test.'},
  {id: 2, title: 'Événement test', organizer: 'Organisateur', category: 'Catégorie', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, registered: 10, status: 'À venir', description: 'Description de l\'événement test.'},
  {id: 3, title: 'Événement test', organizer: 'Organisateur', category: 'Catégorie', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, registered: 10, status: 'À venir', description: 'Description de l\'événement test.'},
  {id: 4, title: 'Événement test', organizer: 'Organisateur', category: 'Catégorie', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, registered: 10, status: 'À venir', description: 'Description de l\'événement test.'},
  {id: 5, title: 'Événement test', organizer: 'Organisateur', category: 'Catégorie', location: 'Salle', date: '01/01/2025', time: '10h00', capacity: 20, registered: 10, status: 'À venir', description: 'Description de l\'événement test.'},
])

const stats = [
  {
    label: 'Total événements',
    value: 5,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
 },
  {
    label: 'À venir',
    value: 5,
    bgClass: 'bg-green-100',
    iconClass: 'text-green-600',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
 },
  {
    label: 'En attente de validation',
    value: 2,
    bgClass: 'bg-amber-100',
    iconClass: 'text-amber-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
 },
  {
    label: 'Annulés',
    value: 0,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/></svg>`,
 },
]

const search = ref('')
const filterCategory = ref('')
const filterStatus = ref('')

const filteredEvents = computed(() => {
  return events.value.filter(e => {
    const q = search.value.toLowerCase()
    if (q && !e.title.toLowerCase().includes(q) && !e.organizer.toLowerCase().includes(q)) return false
    if (filterCategory.value && e.category !== filterCategory.value) return false
    if (filterStatus.value && e.status !== filterStatus.value) return false
    return true
 })
})

const detailEvent = ref(null)

function openDetail(event) {
  detailEvent.value = event
}

function cancelEvent(event) {
  if (!confirm(`Annuler l'événement "${event.title}" ?`)) return
  events.value.find(e => e.id === event.id).status = 'Annulé'
}

function deleteEvent(event) {
  if (!confirm(`Supprimer définitivement "${event.title}" ?`)) return
  events.value = events.value.filter(e => e.id !== event.id)
  if (detailEvent.value?.id === event.id) detailEvent.value = null
}

function categoryBadge(category) {
  const map = {
    'Atelier': 'bg-blue-100 text-blue-700',
    'Conférence': 'bg-purple-100 text-purple-700',
    'Animation': 'bg-green-100 text-green-700',
    'Formation': 'bg-amber-100 text-amber-700',
 }
  return map[category] ?? 'bg-gray-100 text-gray-600'
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
