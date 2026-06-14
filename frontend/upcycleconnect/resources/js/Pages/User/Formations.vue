<template>
    <UserLayout>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                    Formations
                </h1>
                <p class="text-sm text-gray-400 mt-1">Apprenez les techniques d'upcycling avec nos formations</p>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        @input="onSearchInput"
                        type="text"
                        placeholder="Rechercher une formation…"
                        class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                </div>
                <select v-model="filterLevel" @change="resetAndFetch"
                    class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Tous les niveaux</option>
                    <option value="beginner">Débutant</option>
                    <option value="intermediate">Intermédiaire</option>
                    <option value="advanced">Avancé</option>
                </select>
                <select v-if="categories.length" v-model="filterCategory" @change="resetAndFetch"
                    class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Toutes les thématiques</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} résultat(s)</span>
            </div>

            <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>

            <div v-else-if="formations.length === 0" class="py-12 text-center">
                <svg class="w-12 h-12 text-gray-200 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                <p class="text-gray-400 text-sm">Aucune formation disponible pour le moment.</p>
            </div>

            <div v-else class="grid grid-cols-3 gap-5 p-6">
                <div
                    v-for="f in formations"
                    :key="f.id"
                    class="bg-gray-50 rounded-2xl overflow-hidden hover:shadow-md transition-shadow cursor-pointer flex flex-col border border-gray-100"
                    @click="openDetail(f)">
                    <div class="h-32 bg-secondary/10 flex items-center justify-center">
                        <svg class="w-10 h-10 text-secondary/40" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                        </svg>
                    </div>
                    <div class="p-4 flex flex-col gap-2 flex-1">
                        <div class="flex items-center gap-2">
                            <span :class="levelClass(f.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ levelLabel(f.level) }}
                            </span>
                            <span v-if="f.duration_hours" class="text-xs text-gray-400">{{ f.duration_hours }}h</span>
                        </div>
                        <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ f.title }}</h3>
                        <p class="text-xs text-gray-500 line-clamp-2">{{ f.description }}</p>
                        <div class="mt-auto pt-2 flex items-center justify-between">
                            <span class="text-xs text-gray-400">Par {{ f.creator_name }}</span>
                            <span v-if="f.is_registered" class="text-xs text-secondary font-medium">Inscrit</span>
                            <span v-else-if="f.capacity && f.inscription_count >= f.capacity" class="text-xs text-gray-400">Complet</span>
                        </div>
                    </div>
                </div>
            </div>

            <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />
        </div>

        <div v-if="detailFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="detailFormation = null">
            <div class="absolute inset-0 bg-black/40" @click="detailFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ detailFormation.title }}</h3>
                        <p class="text-xs text-gray-500 mt-0.5">Animée par {{ detailFormation.formateur_name ?? detailFormation.creator_name }}</p>
                    </div>
                    <button @click="detailFormation = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Niveau</p>
                            <span :class="levelClass(detailFormation.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ levelLabel(detailFormation.level) }}
                            </span>
                        </div>
                        <div v-if="detailFormation.duration_hours">
                            <p class="text-xs text-gray-400 mb-0.5">Durée</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.duration_hours }}h</p>
                        </div>
                        <div v-if="detailFormation.date">
                            <p class="text-xs text-gray-400 mb-0.5">Date</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.date.slice(0,10) }}</p>
                            <p class="text-xs text-gray-500">{{ detailFormation.date.slice(11,16) }}</p>
                        </div>
                        <div v-if="detailFormation.location">
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.location }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Inscriptions</p>
                            <p class="font-medium text-gray-800">
                                {{ detailFormation.inscription_count }}
                                <span v-if="detailFormation.capacity"> / {{ detailFormation.capacity }}</span>
                                participants
                            </p>
                        </div>
                        <div v-if="detailFormation.category_name">
                            <p class="text-xs text-gray-400 mb-0.5">Thématique</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.category_name }}</p>
                        </div>
                    </div>
                    <div v-if="detailFormation.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailFormation.description }}</p>
                    </div>
                    <div v-if="detailFormation.is_registered" class="flex items-center gap-2 p-3 rounded-xl bg-secondary/10">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        <p class="text-sm text-secondary font-medium">Vous êtes inscrit à cette formation</p>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button
                        v-if="detailFormation.is_registered"
                        @click="toggleRegistration(detailFormation); detailFormation = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Se désinscrire
                    </button>
                    <button
                        v-else-if="!isFull(detailFormation)"
                        @click="toggleRegistration(detailFormation); detailFormation = null"
                        class="px-3 py-1.5 text-sm font-semibold bg-secondary text-white rounded-lg hover:bg-secondary-dark transition-colors">
                        S'inscrire
                    </button>
                    <span v-else class="px-3 py-1.5 text-sm text-gray-400">Complet</span>
                </div>
            </div>
        </div>
    </UserLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
import UserLayout from '@/Layouts/UserLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'
import { useAuthStore } from '@/stores/auth.js'

const auth = useAuthStore()

const formations = ref([])
const categories = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const search = ref('')
const filterLevel = ref('')
const filterCategory = ref('')
const detailFormation = ref(null)

const canRegister = computed(() => auth.hasPermission('register_formation'))

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(resetAndFetch, 300)
}

function resetAndFetch() {
    page.value = 1
    fetchFormations()
}

function levelLabel(level) {
    return { beginner: 'Débutant', intermediate: 'Intermédiaire', advanced: 'Avancé' }[level] ?? level
}

function levelClass(level) {
    return {
        beginner: 'bg-secondary/10 text-secondary',
        intermediate: 'bg-primary/10 text-primary',
        advanced: 'bg-red-100 text-red-500',
    }[level] ?? 'bg-gray-100 text-gray-500'
}

function isFull(f) {
    return f.capacity > 0 && f.inscription_count >= f.capacity
}

function openDetail(f) {
    detailFormation.value = f
}

async function toggleRegistration(f) {
    if (!canRegister.value) return
    f.loading = true
    try {
        if (f.is_registered) {
            await api.delete(`/user/formation/${f.id}/unregister`)
            f.is_registered = false
            f.inscription_count = Math.max(0, f.inscription_count - 1)
        } else {
            await api.post(`/user/formation/${f.id}/register`)
            f.is_registered = true
            f.inscription_count++
        }
    } catch (err) {
        alert(err.response?.data?.error ?? 'Erreur lors de l\'inscription.')
    } finally {
        f.loading = false
    }
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function fetchFormations(silent = false) {
    if (!silent) loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        if (filterLevel.value) params.level = filterLevel.value
        if (filterCategory.value) params.id_category = filterCategory.value
        const { data } = await api.get('/formations', { params })
        formations.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (e) {
        console.error('fetchFormations error:', e)
    } finally {
        if (!silent) loading.value = false
    }
}

async function fetchCategories() {
    try {
        const { data } = await api.get('/categories')
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {}
}

async function fetchAll(silent = false) {
    await Promise.all([fetchFormations(silent), fetchCategories()])
}

usePolling(fetchAll)
</script>
