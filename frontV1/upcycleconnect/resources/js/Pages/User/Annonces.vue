<template>
  <UserLayout @openDepot="showDepot = true">

    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Annonces</h1>
      <button
        @click="showDepot = true"
        class="flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-lg hover:bg-primary-dark transition-colors"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
        </svg>
        Publier
      </button>
    </div>

    <div class="flex items-center gap-4 mb-6">
      <div class="flex gap-1 bg-gray-100 rounded-lg p-1">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          @click="switchTab(tab.value)"
          :class="[
            'px-4 py-1.5 rounded-md text-sm font-medium transition-colors',
            activeTab === tab.value ? 'bg-white text-gray-800 shadow-sm' : 'text-gray-500 hover:text-gray-700'
          ]"
        >{{ tab.label }}</button>
      </div>

      <div v-if="activeTab === 'all'" class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
        </svg>
        <input
          v-model="search"
          @input="debouncedFetch"
          type="text"
          placeholder="Rechercher une annonce…"
          class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
        />
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-400 text-sm">Chargement…</div>

    <div v-else-if="announcements.length === 0" class="text-center py-12 text-gray-400 text-sm">
      {{ activeTab === 'mine' ? 'Vous n\'avez aucune annonce.' : 'Aucune annonce disponible.' }}
    </div>

    <div v-else class="grid grid-cols-3 gap-5">
      <div
        v-for="a in announcements"
        :key="a.id"
        class="bg-white rounded-2xl shadow-sm overflow-hidden hover:shadow-md transition-shadow"
      >
        <div class="relative h-40 bg-gray-100 flex items-center justify-center cursor-pointer" @click="openDetail(a)">
          <img v-if="a.first_photo" :src="a.first_photo" class="w-full h-full object-cover" :alt="a.title" />
          <svg v-else class="w-10 h-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"/>
          </svg>
          <span :class="['absolute top-2 right-2 px-2 py-0.5 rounded-full text-xs font-semibold text-white', a.type === 'vente' ? 'bg-primary' : 'bg-secondary']">
            {{ a.type === 'vente' ? 'VENTE' : 'DON' }}
          </span>
        </div>

        <div class="p-4">
          <h3 class="font-semibold text-gray-800 truncate">{{ a.title }}</h3>
          <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ a.description }}</p>

          <div class="mt-3 flex items-center gap-2">
            <span v-if="a.condition" class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">{{ conditionLabel(a.condition) }}</span>
          </div>

          <p class="text-xs text-gray-400 mt-2">Disponible le {{ formatDate(a.availability_date) }}</p>

          <div class="mt-3 flex items-center justify-between">
            <div class="flex items-center gap-1.5">
              <div class="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-white text-xs font-semibold">
                {{ a.author_name?.[0] ?? '?' }}
              </div>
              <span class="text-xs text-gray-500">{{ a.author_name || 'Inconnu' }}</span>
            </div>
            <div class="flex items-center gap-2">
              <button
                v-if="canDelete()"
                @click.stop="confirmDelete(a)"
                class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                title="Supprimer"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
              </button>
              <button @click="openDetail(a)" class="px-3 py-1 bg-primary text-white text-xs font-medium rounded-lg hover:bg-primary-dark transition-colors">
                Voir
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <CreateAnnouncementModal v-model="showDepot" @created="onCreated" />

    <!-- Detail modal -->
    <div v-if="selected" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="closeDetail">
      <div class="absolute inset-0 bg-black/40" @click="closeDetail" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selected.title }}</h3>
          <button @click="closeDetail" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <!-- Photos -->
        <div v-if="selectedPhotos.length" class="relative bg-gray-100">
          <img :src="selectedPhotos[photoIndex]" class="w-full h-56 object-cover" :alt="selected.title" />
          <template v-if="selectedPhotos.length > 1">
            <button @click="photoIndex = (photoIndex - 1 + selectedPhotos.length) % selectedPhotos.length"
              class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/40 hover:bg-black/60 rounded-full flex items-center justify-center text-white transition-colors">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>
            </button>
            <button @click="photoIndex = (photoIndex + 1) % selectedPhotos.length"
              class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/40 hover:bg-black/60 rounded-full flex items-center justify-center text-white transition-colors">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/></svg>
            </button>
            <div class="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
              <button v-for="(_, i) in selectedPhotos" :key="i" @click="photoIndex = i"
                :class="['w-2 h-2 rounded-full transition-colors', i === photoIndex ? 'bg-white' : 'bg-white/50']" />
            </div>
          </template>
        </div>

        <div class="px-6 py-5 space-y-3 text-sm text-gray-700">
          <div class="flex gap-2">
            <span :class="['px-2 py-0.5 rounded-full text-xs font-semibold text-white', selected.type === 'vente' ? 'bg-primary' : 'bg-secondary']">{{ selected.type === 'vente' ? 'Vente' : 'Don' }}</span>
            <span v-if="selected.condition" class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">{{ conditionLabel(selected.condition) }}</span>
          </div>
          <p class="leading-relaxed">{{ selected.description }}</p>
          <div class="grid grid-cols-2 gap-4 text-xs text-gray-500">
            <div><span class="font-medium text-gray-700">Ville :</span> {{ selected.city }}</div>
            <div><span class="font-medium text-gray-700">Prix :</span> {{ selected.price ? selected.price + ' €' : 'Gratuit' }}</div>
            <div><span class="font-medium text-gray-700">Disponible :</span> {{ formatDate(selected.availability_date) }}</div>
            <div><span class="font-medium text-gray-700">Par :</span> {{ selected.author_name || 'Inconnu' }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete confirm modal -->
    <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/40" @click="toDelete = null" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
        <h3 class="font-semibold text-gray-800 mb-2">Supprimer l'annonce ?</h3>
        <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimée.</p>
        <div class="flex gap-3">
          <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
          <button @click="deleteAnnouncement" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
            {{ deleting ? 'Suppression…' : 'Supprimer' }}
          </button>
        </div>
      </div>
    </div>

  </UserLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import UserLayout from '@/Layouts/UserLayout.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const auth = useAuthStore()
const announcements = ref([])
const loading = ref(true)
const search = ref('')
const showDepot = ref(false)
const selected = ref(null)
const selectedPhotos = ref([])
const photoIndex = ref(0)
const toDelete = ref(null)
const deleting = ref(false)
const activeTab = ref('all')

const tabs = [
  { label: 'Toutes', value: 'all' },
  { label: 'Mes annonces', value: 'mine' },
]

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchAnnouncements, 300)
}

function switchTab(tab) {
  activeTab.value = tab
  search.value = ''
  fetchAnnouncements()
}

async function fetchAnnouncements() {
  loading.value = true
  try {
    if (activeTab.value === 'mine') {
      const { data } = await api.get('/user/announcements')
      announcements.value = data
    } else {
      const params = {}
      if (search.value) params.search = search.value
      const { data } = await api.get('/announcements', { params })
      announcements.value = data
    }
  } catch {
    announcements.value = []
  } finally {
    loading.value = false
  }
}

async function openDetail(a) {
  selected.value = a
  photoIndex.value = 0
  try {
    const { data } = await api.get(`/announcements/${a.id}/documents`)
    selectedPhotos.value = data.map(d => d.link).filter(Boolean)
  } catch {
    selectedPhotos.value = a.first_photo ? [a.first_photo] : []
  }
}

function closeDetail() {
  selected.value = null
  selectedPhotos.value = []
  photoIndex.value = 0
}

function canDelete() {
  if (activeTab.value === 'mine') return true
  return auth.hasPermission('manage_announcements')
}

function confirmDelete(a) {
  toDelete.value = a
}

async function deleteAnnouncement() {
  if (!toDelete.value) return
  deleting.value = true
  try {
    if (activeTab.value === 'mine') {
      await api.delete(`/user/announcement/${toDelete.value.id}`)
    } else {
      await api.delete(`/announcements/${toDelete.value.id}`)
    }
    announcements.value = announcements.value.filter(a => a.id !== toDelete.value.id)
    toDelete.value = null
  } catch (e) {
    alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
  } finally {
    deleting.value = false
  }
}

function onCreated() {
  if (activeTab.value === 'mine') fetchAnnouncements()
  else fetchAnnouncements()
}

function conditionLabel(v) {
  return { neuf: 'Neuf', bon_etat: 'Bon état', use: 'Usé' }[v] ?? v
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(fetchAnnouncements)
</script>
