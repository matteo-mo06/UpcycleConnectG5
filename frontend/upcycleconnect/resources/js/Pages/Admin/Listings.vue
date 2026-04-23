<template>
    <AdminLayout title="Annonces">

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-4 gap-5 mb-8">
                <div
                    v-for="stat in stats"
                    :key="stat.label"
                    class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                    <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                        <div :class="stat.iconClass" v-html="stat.icon" />
                    </div>
                    <div class="min-w-0">
                        <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                        <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
                <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h2 class="text-base font-semibold text-gray-800">Dernières annonces ajoutées</h2>
                    <span class="text-xs text-gray-400">Vérifiez la conformité des nouvelles publications</span>
                </div>
                <div class="divide-y divide-gray-50">
                    <div
                        v-for="listing in recentListings"
                        :key="listing.id"
                        class="px-5 py-3.5 flex items-center justify-between hover:bg-gray-50/60 transition-colors">
                        <div class="flex-1 mr-4">
                            <div class="flex items-center gap-2">
                                <p class="font-medium text-gray-800 text-sm truncate">{{ listing.title }}</p>
                                <span
                                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                                    :class="typeBadge(listing.type)">
                                    {{ listing.type }}
                                </span>
                            </div>
                            <p class="text-xs text-gray-500 mt-0.5">
                                Par <span class="font-medium text-gray-700">{{ listing.author }}</span>
                                · {{ listing.category }} · {{ listing.date }}
                            </p>
                        </div>
                        <div class="flex-shrink-0 flex items-center gap-2">
                            <button
                                @click="openDetail(listing)"
                                class="px-3 py-1.5 text-xs font-medium text-primary border border-primary/30 rounded-lg hover:bg-primary/5 transition-colors">
                                Voir
                            </button>
                            <button
                                @click="confirmDelete(listing)"
                                class="px-3 py-1.5 text-xs font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                                Retirer
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
                            placeholder="Rechercher par titre ou auteur…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>

                    <select
                        v-model="filterType"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les types</option>
                        <option>Don</option>
                        <option>Vente</option>
                    </select>

                    <select
                        v-model="filterCategory"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Toutes les catégories</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
                    </select>

                    <select
                        v-model="filterStatus"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option>Active</option>
                        <option>Expirée</option>
                        <option>Supprimée</option>
                    </select>

                    <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">
                        {{ filteredListings.length }} résultat(s)
                    </span>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Annonce</th>
                                <th class="text-left text-white font-medium px-5 py-3">Type</th>
                                <th class="text-left text-white font-medium px-5 py-3">Catégorie</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-left text-white font-medium px-5 py-3">Date</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(listing, i) in filteredListings"
                                :key="listing.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                                <td class="px-5 py-3">
                                    <div class="flex items-center gap-2">
                                        <svg
                                            v-if="listing.featured"
                                            class="w-3.5 h-3.5 text-amber-400 flex-shrink-0"
                                            viewBox="0 0 24 24"
                                            fill="currentColor">
                                            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                                        </svg>
                                        <div class="min-w-0">
                                            <p class="font-medium text-gray-800 truncate max-w-48">{{ listing.title }}</p>
                                            <p class="text-xs text-gray-500 truncate">{{ listing.author }}</p>
                                        </div>
                                    </div>
                                </td>

                                <td class="px-5 py-3">
                                    <span
                                        class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                                        :class="typeBadge(listing.type)">
                                        {{ listing.type }}
                                    </span>
                                </td>

                                <td class="px-5 py-3 text-gray-500 text-xs">{{ listing.category }}</td>

                                <td class="px-5 py-3">
                                    <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                        <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(listing.status)" />
                                        <span :class="statusText(listing.status)">{{ listing.status }}</span>
                                    </span>
                                </td>

                                <td class="px-5 py-3 text-gray-500 text-xs">{{ listing.date }}</td>

                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button
                                            @click="openDetail(listing)"
                                            title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>

                                        <button
                                            @click="toggleFeatured(listing)"
                                            :title="listing.featured ? 'Retirer la mise en avant' : 'Mettre en avant'"
                                            :class="[
                                                'p-1.5 rounded-lg transition-colors',
                                                listing.featured
                                                    ? 'text-amber-400 hover:text-amber-600 hover:bg-amber-50'
                                                    : 'text-gray-400 hover:text-amber-400 hover:bg-amber-50'
                                            ]">
                                            <svg class="w-4 h-4" viewBox="0 0 24 24" :fill="listing.featured ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                                            </svg>
                                        </button>

                                        <button
                                            @click="confirmDelete(listing)"
                                            title="Supprimer l'annonce"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>

                            <tr v-if="filteredListings.length === 0">
                                <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucune annonce ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ filteredListings.length }} annonce(s)
                </div>
            </div>

        </template>

        <div
            v-if="detailListing"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="detailListing = null">
            <div class="absolute inset-0 bg-black/40" @click="detailListing = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
                    <div class="min-w-0">
                        <div class="flex items-center gap-2">
                            <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                                {{ detailListing.title }}
                            </h3>
                            <span
                                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                                :class="typeBadge(detailListing.type)">
                                {{ detailListing.type }}
                            </span>
                        </div>
                        <p class="text-xs text-gray-500 mt-0.5">{{ detailListing.category }}</p>
                    </div>
                    <button
                        @click="detailListing = null"
                        class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Auteur</p>
                            <p class="font-medium text-gray-800">{{ detailListing.author }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailListing.status)" />
                                <span :class="statusText(detailListing.status)">{{ detailListing.status }}</span>
                            </span>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Date de publication</p>
                            <p class="text-gray-700">{{ detailListing.date }}</p>
                        </div>
                    </div>

                    <div>
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailListing.description }}</p>
                    </div>

                    <div v-if="detailListing.tags && detailListing.tags.length">
                        <p class="text-xs text-gray-400 mb-1.5">Tags</p>
                        <div class="flex gap-1.5">
                            <span
                                v-for="tag in detailListing.tags"
                                :key="tag"
                                class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">
                                {{ tag }}
                            </span>
                        </div>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-between items-center">
                    <button
                        @click="toggleFeatured(detailListing); detailListing = null"
                        :class="[
                            'flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg border transition-colors',
                            detailListing.featured
                                ? 'border-amber-300 text-amber-600 hover:bg-amber-50'
                                : 'border-gray-200 text-gray-600 hover:bg-gray-50'
                        ]">
                        <svg class="w-4 h-4" viewBox="0 0 24 24" :fill="detailListing.featured ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                        {{ detailListing.featured ? 'Retirer la mise en avant' : 'Mettre en avant' }}
                    </button>
                    <button
                        @click="confirmDelete(detailListing); detailListing = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Supprimer
                    </button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer l'annonce ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimée.</p>
                <div class="flex gap-3">
                    <button
                        @click="toDelete = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="deleteListing"
                        :disabled="deleting"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const listings = ref([])
const categories = ref([])
const loading = ref(true)
const error = ref('')
const statsData = ref({})
const search = ref('')
const filterType = ref('')
const filterCategory = ref('')
const filterStatus = ref('')
const detailListing = ref(null)
const toDelete = ref(null)
const deleting = ref(false)

const stats = computed(() => [
    {
        label: 'Total annonces',
        value: statsData.value.total ?? '—',
        bgClass: 'bg-red-100',
        iconClass: 'text-red-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
    },
    {
        label: 'Actives',
        value: statsData.value.active ?? '—',
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Mises en avant',
        value: statsData.value.featured ?? '—',
        bgClass: 'bg-amber-100',
        iconClass: 'text-amber-500',
        icon: `<svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>`,
    },
    {
        label: 'Signalées',
        value: statsData.value.reported ?? '—',
        bgClass: 'bg-red-100',
        iconClass: 'text-red-400',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"/></svg>`,
    },
])

const recentListings = computed(() => listings.value.slice(0, 3))

const filteredListings = computed(() => {
    const filtered = listings.value.filter(l => {
        const q = search.value.toLowerCase()
        if (q && !l.title.toLowerCase().includes(q) && !l.author.toLowerCase().includes(q)) return false
        if (filterType.value && l.type !== filterType.value) return false
        if (filterCategory.value && l.category !== filterCategory.value) return false
        if (filterStatus.value && l.status !== filterStatus.value) return false
        return true
    })
    return [...filtered].sort((a, b) => (b.featured ? 1 : 0) - (a.featured ? 1 : 0))
})

onMounted(async () => {
    try {
        const [{ data: announcementsData }, { data: catsData }, { data: announcementStats }] = await Promise.all([
            api.get('/admin/announcements'),
            api.get('/admin/categories'),
            api.get('/admin/announcements/stats'),
        ])

        categories.value = catsData ?? []
        statsData.value = announcementStats

        const catMap = Object.fromEntries(categories.value.map(c => [c.id, c.name]))
        listings.value = announcementsData.map(a => ({
            id: a.id,
            title: a.title,
            author: a.author_name || '—',
            type: a.type === 'vente' ? 'Vente' : 'Don',
            category: catMap[a.id_category] ?? '—',
            status: a.state ?? 'Active',
            featured: false,
            date: a.availability_date?.slice(0, 10) ?? '—',
            description: a.description ?? '—',
            tags: [],
        }))
    } catch {
        error.value = 'Impossible de charger les annonces.'
    } finally {
        loading.value = false
    }
})

function openDetail(listing) {
    detailListing.value = listing
}

function confirmDelete(listing) {
    toDelete.value = listing
}

function toggleFeatured(listing) {
    listings.value.find(l => l.id === listing.id).featured ^= true
}

async function deleteListing() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/announcement/${toDelete.value.id}`)
        listings.value = listings.value.filter(l => l.id !== toDelete.value.id)
        if (detailListing.value?.id === toDelete.value.id) detailListing.value = null
        toDelete.value = null
    } catch {
        alert('Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function typeBadge(type) {
    const map = { 'Don': 'bg-green-100 text-green-700', 'Vente': 'bg-blue-100 text-blue-700' }
    return map[type] ?? 'bg-gray-100 text-gray-600'
}

function statusDot(status) {
    const map = { 'Active': 'bg-green-500', 'Expirée': 'bg-gray-400', 'Supprimée': 'bg-red-400' }
    return map[status] ?? 'bg-gray-300'
}

function statusText(status) {
    const map = { 'Active': 'text-green-700', 'Expirée': 'text-gray-500', 'Supprimée': 'text-red-600' }
    return map[status] ?? 'text-gray-500'
}
</script>
