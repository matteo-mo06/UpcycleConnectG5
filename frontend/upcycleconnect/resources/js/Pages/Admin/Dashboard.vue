<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Espace admin
            </h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-4 gap-5 mb-8">
                <div
                    v-for="stat in stats"
                    :key="stat.label"
                    class="bg-white rounded-2xl shadow-sm p-5 flex items-center gap-4">
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
                        <h2 class="text-base font-semibold text-gray-800">Derniers utilisateurs</h2>
                        <RouterLink to="/admin/users" class="text-xs font-medium text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div
                            v-for="user in recentUsers"
                            :key="user.email"
                            class="px-5 py-3 flex items-center gap-3">
                            <div class="w-7 h-7 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
                                <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
                                </svg>
                            </div>
                            <div class="flex-1 min-w-0">
                                <p class="font-medium text-gray-800 text-sm truncate">{{ user.name }}</p>
                                <p class="text-xs text-gray-500 truncate">{{ user.email }}</p>
                            </div>
                            <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600">{{ user.type }}</span>
                            <span class="text-xs text-gray-400 flex-shrink-0">{{ user.date }}</span>
                        </div>
                        <div v-if="recentUsers.length === 0" class="px-5 py-8 text-center text-gray-400 text-sm">
                            Aucun utilisateur
                        </div>
                    </div>
                </div>

                <div class="col-span-2 bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
                        <h2 class="text-base font-semibold text-gray-800">Prochains événements</h2>
                        <RouterLink to="/admin/events" class="text-xs font-medium text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div
                            v-for="event in upcomingEvents"
                            :key="event.id"
                            class="px-5 py-3.5 flex items-center justify-between hover:bg-gray-50/60 transition-colors duration-100">
                            <div class="min-w-0 mr-3">
                                <p class="font-medium text-gray-800 text-sm truncate">{{ event.title }}</p>
                                <p class="text-xs text-gray-500 mt-0.5">{{ event.date }} · {{ event.location }}</p>
                            </div>
                        </div>
                        <div v-if="upcomingEvents.length === 0" class="px-5 py-8 text-center text-gray-400 text-sm">
                            Aucun événement
                        </div>
                    </div>
                </div>

            </div>

        </template>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'
import { roleLabel, fullName } from '@/utils.js'

const auth = useAuthStore()
const loading = ref(true)
const error = ref('')
const statsData = ref({})
const recentUsers = ref([])
const upcomingEvents = ref([])

const stats = computed(() => [
    {
        label: 'Utilisateurs total',
        value: statsData.value.user_count ?? '-',
        bgClass: 'bg-red-100',
        iconClass: 'text-red-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
    },
    {
        label: 'Annonces actives',
        value: statsData.value.announcement_count ?? '-',
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
    },
    {
        label: 'Événements à venir',
        value: statsData.value.event_count ?? '-',
        bgClass: 'bg-amber-100',
        iconClass: 'text-amber-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
    },
    {
        label: 'Signalements ouverts',
        value: statsData.value.report_count ?? '-',
        bgClass: 'bg-red-100',
        iconClass: 'text-red-400',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
    },
])

onMounted(async () => {
    try {
        const [{ data: usersData }, { data: eventsData }, { data: apiStats }] = await Promise.all([
            api.get('/admin/users'),
            api.get('/admin/events'),
            api.get('/admin/stats'),
        ])

        statsData.value = apiStats

        recentUsers.value = (usersData.data ?? usersData).slice(0, 5).map(u => ({
            name: fullName(u),
            email: u.email,
            type: roleLabel(u.roles?.[0] ?? ''),
            date: u.created_at?.slice(0, 10) ?? '-',
        }))

        upcomingEvents.value = (eventsData.data ?? eventsData).slice(0, 5).map(e => ({
            id: e.id,
            title: e.title,
            date: e.date ?? '-',
            location: e.location ?? '-',
        }))
    } catch {
        error.value = 'Impossible de charger les données.'
    } finally {
        loading.value = false
    }
})
</script>
