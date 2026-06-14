<template>
    <ArtisanLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Espace professionnel
            </h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="relative mb-8">
                <div class="grid grid-cols-4 gap-5" :class="{ 'blur-sm pointer-events-none select-none': !isPremium }">
                    <div
                        v-for="stat in stats"
                        :key="stat.label"
                        class="bg-white rounded-2xl shadow-sm p-5 flex items-center gap-4"
                    >
                        <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                            <div :class="stat.iconClass" v-html="stat.icon" />
                        </div>
                        <div class="min-w-0">
                            <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                            <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                        </div>
                    </div>
                </div>
                <div v-if="!isPremium" class="absolute inset-0 flex items-center justify-center">
                    <div class="bg-white rounded-2xl shadow-lg px-6 py-4 text-center border border-primary/20">
                        <p class="font-semibold text-gray-800 text-sm mb-1">Statistiques avancées</p>
                        <p class="text-xs text-gray-500 mb-3">Disponibles avec l'abonnement Premium</p>
                        <RouterLink to="/abonnement"
                            class="inline-block px-4 py-1.5 bg-primary text-white text-xs font-semibold rounded-lg hover:bg-primary-dark transition-colors">
                            Passer Premium
                        </RouterLink>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-5 gap-6">

                <div class="col-span-3 bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
                        <h2 class="text-base font-semibold text-gray-800">Dernières annonces</h2>
                        <RouterLink to="/artisan/annonces" class="text-xs font-medium text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div
                            v-for="ann in recentAnnouncements"
                            :key="ann.id"
                            class="px-5 py-3 flex items-center gap-3"
                        >
                            <div class="flex-1 min-w-0">
                                <p class="font-medium text-gray-800 text-sm truncate">{{ ann.title }}</p>
                                <p class="text-xs text-gray-400 mt-0.5 capitalize">{{ ann.type }}</p>
                            </div>
                            <span :class="['inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium', stateBadge(ann.state)]">{{ ann.state }}</span>
                            <span class="text-xs text-gray-400 flex-shrink-0">{{ ann.date }}</span>
                        </div>
                        <div v-if="recentAnnouncements.length === 0" class="px-5 py-8 text-center text-gray-400 text-sm">
                            Aucune annonce publiée
                        </div>
                    </div>
                </div>

                <div class="col-span-2 bg-white rounded-2xl shadow-sm p-5">
                    <h2 class="text-base font-semibold text-gray-800 mb-1">Actions rapides</h2>
                    <p class="text-xs text-gray-400 mb-4">Accès direct à vos actions</p>
                    <div class="grid grid-cols-2 gap-3">
                        <button
                            @click="showCreateAnnonce = true"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-primary text-white hover:bg-primary-dark transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Publier une annonce</span>
                        </button>
                        <RouterLink
                            to="/artisan/evenements"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-primary text-primary hover:bg-primary/5 transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Créer un événement</span>
                        </RouterLink>
                        <RouterLink
                            to="/artisan/formations"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-secondary text-white hover:bg-secondary-dark transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Créer une formation</span>
                        </RouterLink>
                        <RouterLink
                            to="/artisan/projets"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-secondary text-secondary hover:bg-secondary/5 transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Créer un projet</span>
                        </RouterLink>
                    </div>
                </div>

            </div>

        </template>

        <CreateAnnouncementModal v-model="showCreateAnnonce" />

    </ArtisanLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const auth = useAuthStore()
const loading = ref(true)
const error = ref('')
const showCreateAnnonce = ref(false)
const announcements = ref([])
const statsData = ref({})
const isPremium = ref(false)

const stats = computed(() => [
    {
        label: 'Annonces publiées',
        value: announcements.value.length,
        bgClass: 'bg-primary/10',
        iconClass: 'text-primary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
    },
    {
        label: 'Annonces actives',
        value: statsData.value.active_announcements ?? '-',
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Dépôts en attente',
        value: statsData.value.pending_deposits ?? '-',
        bgClass: 'bg-amber-100',
        iconClass: 'text-amber-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>`,
    },
    {
        label: 'Annonces vendues / données',
        value: announcements.value.filter(a => a.state === 'Vendu' || a.state === 'Donné').length,
        bgClass: 'bg-secondary/20',
        iconClass: 'text-secondary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg>`,
    },
])

const recentAnnouncements = computed(() =>
    announcements.value.slice(0, 5).map(a => ({
        id: a.id,
        title: a.title,
        type: a.type,
        state: a.state,
        date: a.availability_date?.slice(0, 10) ?? '-',
    }))
)

function stateBadge(state) {
    if (state === 'Active') return 'bg-green-100 text-green-700'
    if (state === 'En attente') return 'bg-amber-100 text-amber-700'
    if (state === 'Refusée') return 'bg-red-100 text-red-600'
    return 'bg-gray-100 text-gray-600'
}

onMounted(async () => {
    try {
        const [{ data: annData }, { data: userStats }, { data: limData }] = await Promise.all([
            api.get('/user/announcements'),
            api.get('/user/stats'),
            api.get('/user/limits'),
        ])
        announcements.value = Array.isArray(annData) ? annData : (annData.data ?? [])
        statsData.value = userStats
        isPremium.value = limData?.is_premium ?? false
    } catch {
        error.value = 'Impossible de charger les données.'
    } finally {
        loading.value = false
    }
})
</script>
