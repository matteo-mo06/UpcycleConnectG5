<template>
    <SalarieLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Espace salarié
            </h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-4 gap-5 mb-8">
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

            <div class="grid grid-cols-5 gap-6">

                <div class="col-span-3 bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
                        <h2 class="text-base font-semibold text-gray-800">Mes derniers conseils</h2>
                        <RouterLink to="/salarie/conseils" class="text-xs font-medium text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div
                            v-for="conseil in recentConseils"
                            :key="conseil.id"
                            class="px-5 py-3 flex items-center gap-3"
                        >
                            <div class="flex-1 min-w-0">
                                <p class="font-medium text-gray-800 text-sm truncate">{{ conseil.title }}</p>
                                <p v-if="conseil.category_name" class="text-xs text-primary mt-0.5">{{ conseil.category_name }}</p>
                            </div>
                            <span class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0" :class="priorityBadge(conseil.priority).class">{{ priorityBadge(conseil.priority).label }}</span>
                            <span v-if="conseil.is_expired" class="px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500 flex-shrink-0">Expiré</span>
                            <span class="text-xs text-gray-400 flex-shrink-0">{{ conseil.created_at?.slice(0, 10) ?? '-' }}</span>
                        </div>
                        <div v-if="recentConseils.length === 0" class="px-5 py-8 text-center text-gray-400 text-sm">
                            Aucun conseil publié
                        </div>
                    </div>
                </div>

                <div class="col-span-2 bg-white rounded-2xl shadow-sm p-5">
                    <h2 class="text-base font-semibold text-gray-800 mb-1">Actions rapides</h2>
                    <p class="text-xs text-gray-400 mb-4">Accès direct à vos actions</p>
                    <div class="grid grid-cols-2 gap-3">
                        <RouterLink
                            to="/salarie/conseils"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-primary text-white hover:bg-primary-dark transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Créer un conseil</span>
                        </RouterLink>
                        <RouterLink
                            to="/salarie/formations"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-primary text-primary hover:bg-primary/5 transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Créer une formation</span>
                        </RouterLink>
                        <RouterLink
                            to="/salarie/evenements"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-secondary text-white hover:bg-secondary-dark transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Créer un événement</span>
                        </RouterLink>
                        <RouterLink
                            to="/salarie/forum"
                            class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-secondary text-secondary hover:bg-secondary/5 transition-colors min-h-24"
                        >
                            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                            </svg>
                            <span class="text-sm font-medium text-center">Aller au forum</span>
                        </RouterLink>
                    </div>
                </div>

            </div>

        </template>

    </SalarieLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import SalarieLayout from '@/Layouts/SalarieLayout.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const auth = useAuthStore()
const loading = ref(true)
const error = ref('')
const formations = ref([])
const events = ref([])
const conseils = ref([])
const topics = ref([])

const recentConseils = computed(() => conseils.value.slice(0, 5))

const priorityBadges = {
    1: { label: 'Haute', class: 'bg-primary/10 text-primary' },
    2: { label: 'Moyenne', class: 'bg-amber-100 text-amber-700' },
    3: { label: 'Basse', class: 'bg-gray-100 text-gray-500' },
}

function priorityBadge(p) {
    return priorityBadges[p] ?? priorityBadges[2]
}

const stats = computed(() => [
    {
        label: 'Formations créées',
        value: formations.value.length,
        bgClass: 'bg-primary/10',
        iconClass: 'text-primary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>`,
    },
    {
        label: 'Événements créés',
        value: events.value.length,
        bgClass: 'bg-secondary/20',
        iconClass: 'text-secondary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
    },
    {
        label: 'Conseils publiés',
        value: conseils.value.length,
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>`,
    },
    {
        label: 'Sujets forum créés',
        value: topics.value.length,
        bgClass: 'bg-amber-100',
        iconClass: 'text-amber-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>`,
    },
])

onMounted(async () => {
    try {
        const userId = auth.user?.id
        const [foRes, evRes, coRes, toRes] = await Promise.allSettled([
            api.get('/user/my-formations'),
            api.get('/user/my-events'),
            api.get('/conseils'),
            api.get('/forum/topics', { params: { page: 1, limit: 100 } }),
        ])
        formations.value = foRes.status === 'fulfilled' ? (foRes.value.data?.data ?? foRes.value.data ?? []) : []
        events.value     = evRes.status === 'fulfilled' ? (evRes.value.data ?? []) : []

        const allConseils = coRes.status === 'fulfilled' ? (coRes.value.data?.data ?? coRes.value.data ?? []) : []
        conseils.value = userId ? allConseils.filter(c => c.id_creator === userId) : []

        const allTopics = toRes.status === 'fulfilled' ? (toRes.value.data?.data ?? []) : []
        topics.value = userId ? allTopics.filter(t => t.id_author === userId) : []
    } catch {
        error.value = 'Impossible de charger les données.'
    } finally {
        loading.value = false
    }
})
</script>
