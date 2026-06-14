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

            <div class="grid grid-cols-2 gap-5">

                <!-- Annonces -->
                <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                        <p class="text-sm font-semibold text-gray-800">Dernières annonces</p>
                        <RouterLink to="/admin/annonces" class="text-xs text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div v-for="a in lastItems.announcements" :key="a.id" class="px-4 py-2 flex items-center justify-between gap-2">
                            <p class="text-xs text-gray-700 truncate flex-1">{{ a.title }}</p>
                            <p class="text-xs text-gray-400 flex-shrink-0">{{ a.date }}</p>
                        </div>
                        <p v-if="lastItems.announcements.length === 0" class="px-4 py-4 text-xs text-gray-400">Aucune annonce</p>
                    </div>
                </div>

                <!-- Événements -->
                <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                        <p class="text-sm font-semibold text-gray-800">Derniers événements</p>
                        <RouterLink to="/admin/events" class="text-xs text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div v-for="e in lastItems.events" :key="e.id" class="px-4 py-2 flex items-center justify-between gap-2">
                            <p class="text-xs text-gray-700 truncate flex-1">{{ e.title }}</p>
                            <p class="text-xs text-gray-400 flex-shrink-0">{{ e.date }}</p>
                        </div>
                        <p v-if="lastItems.events.length === 0" class="px-4 py-4 text-xs text-gray-400">Aucun événement</p>
                    </div>
                </div>

                <!-- Formations -->
                <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                        <p class="text-sm font-semibold text-gray-800">Dernières formations</p>
                        <RouterLink to="/admin/formations" class="text-xs text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div v-for="f in lastItems.formations" :key="f.id" class="px-4 py-2 flex items-center justify-between gap-2">
                            <p class="text-xs text-gray-700 truncate flex-1">{{ f.title }}</p>
                            <p class="text-xs text-gray-400 flex-shrink-0">{{ f.date }}</p>
                        </div>
                        <p v-if="lastItems.formations.length === 0" class="px-4 py-4 text-xs text-gray-400">Aucune formation</p>
                    </div>
                </div>

                <!-- Projets -->
                <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                        <p class="text-sm font-semibold text-gray-800">Derniers projets</p>
                        <RouterLink to="/admin/projets" class="text-xs text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div v-for="p in lastItems.projects" :key="p.id" class="px-4 py-2 flex items-center justify-between gap-2">
                            <p class="text-xs text-gray-700 truncate flex-1">{{ p.title }}</p>
                            <p class="text-xs text-gray-400 flex-shrink-0">{{ p.date }}</p>
                        </div>
                        <p v-if="lastItems.projects.length === 0" class="px-4 py-4 text-xs text-gray-400">Aucun projet</p>
                    </div>
                </div>

                <!-- Derniers utilisateurs -->
                <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                        <p class="text-sm font-semibold text-gray-800">Derniers utilisateurs inscrits</p>
                        <RouterLink to="/admin/utilisateurs" class="text-xs text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div v-for="u in lastItems.users" :key="u.id" class="px-4 py-2 flex items-center justify-between gap-2">
                            <p class="text-xs text-gray-700 truncate flex-1">{{ u.name }}</p>
                            <p class="text-xs text-gray-400 flex-shrink-0">{{ u.date }}</p>
                        </div>
                        <p v-if="lastItems.users.length === 0" class="px-4 py-4 text-xs text-gray-400">Aucun utilisateur</p>
                    </div>
                </div>

                <!-- Derniers signalements -->
                <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                        <p class="text-sm font-semibold text-gray-800">Derniers signalements</p>
                        <RouterLink to="/admin/signalements" class="text-xs text-primary hover:underline">Voir tout</RouterLink>
                    </div>
                    <div class="divide-y divide-gray-50">
                        <div v-for="r in lastItems.reports" :key="r.id" class="px-4 py-2 flex items-center justify-between gap-2">
                            <p class="text-xs text-gray-700 truncate flex-1">{{ r.reason }}</p>
                            <p class="text-xs text-gray-400 flex-shrink-0">{{ r.date }}</p>
                        </div>
                        <p v-if="lastItems.reports.length === 0" class="px-4 py-4 text-xs text-gray-400">Aucun signalement</p>
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

const auth = useAuthStore()
const loading = ref(true)
const error = ref('')
const statsData = ref({})
const lastItems = ref({ announcements: [], events: [], formations: [], projects: [], users: [], reports: [] })

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
        const [{ data: apiStats }, { data: annData }, { data: evtData }, { data: fmtData }, { data: prjData }, { data: usrData }, { data: rptData }] = await Promise.all([
            api.get('/admin/stats'),
            api.get('/admin/announcements', { params: { limit: 10 } }),
            api.get('/admin/events', { params: { limit: 10 } }),
            api.get('/admin/formations', { params: { limit: 10 } }),
            api.get('/admin/projects', { params: { limit: 10 } }),
            api.get('/admin/users', { params: { limit: 10 } }),
            api.get('/admin/reports', { params: { limit: 10 } }),
        ])

        statsData.value = apiStats
        lastItems.value.announcements = (annData.data ?? annData).map(a => ({ id: a.id_announcement ?? a.id, title: a.title, date: a.created_at?.slice(0, 10) ?? '-' }))
        lastItems.value.events = (evtData.data ?? evtData).map(e => ({ id: e.id_event ?? e.id, title: e.title, date: e.date?.slice(0, 10) ?? e.created_at?.slice(0, 10) ?? '-' }))
        lastItems.value.formations = (fmtData.data ?? fmtData).map(f => ({ id: f.id_formation ?? f.id, title: f.title, date: f.date?.slice(0, 10) ?? f.created_at?.slice(0, 10) ?? '-' }))
        lastItems.value.projects = (prjData.data ?? prjData).map(p => ({ id: p.id_project ?? p.id, title: p.title, date: p.created_at?.slice(0, 10) ?? '-' }))
        lastItems.value.users = (usrData.data ?? usrData).map(u => ({ id: u.id, name: `${u.first_name ?? ''} ${u.last_name ?? ''}`.trim() || u.email, date: u.created_at?.slice(0, 10) ?? '-' }))
        lastItems.value.reports = (rptData.data ?? rptData).map(r => ({ id: r.id_report, reason: r.reason || r.content_type || 'Signalement', date: r.created_at?.slice(0, 10) ?? '-' }))
    } catch {
        error.value = 'Impossible de charger les données.'
    } finally {
        loading.value = false
    }
})
</script>
