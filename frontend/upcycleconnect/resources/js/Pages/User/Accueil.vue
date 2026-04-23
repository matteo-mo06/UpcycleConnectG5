<template>
    <UserLayout @openDepot="showDepot = true">

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Bonjour {{ auth.user?.first_name ?? '' }}
            </h1>
        </div>

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
                    <p class="text-sm text-gray-500 mt-1">{{ stat.label }}</p>
                    <span v-if="stat.badge" class="inline-block mt-1 px-1.5 py-0.5 rounded text-xs font-medium bg-primary/10 text-primary">
                        {{ stat.badge }}
                    </span>
                </div>
            </div>
        </div>

        <div class="grid grid-cols-5 gap-6">

            <div class="col-span-3 bg-white rounded-2xl shadow-sm p-5">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Activité récente</h2>
                <p class="text-xs text-gray-400 mb-4">Dernières actualités de la plateforme</p>
                <div class="space-y-4">
                    <p v-if="activity.length === 0" class="text-sm text-gray-400 text-center py-6">Aucune activité récente pour le moment.</p>
                    <div v-for="item in activity" :key="item.id" class="flex items-start gap-3">
                        <div :class="['w-9 h-9 rounded-full flex items-center justify-center flex-shrink-0 text-white text-sm font-semibold', item.avatarColor]">
                            {{ item.initials }}
                        </div>
                        <div class="min-w-0 flex-1">
                            <p class="text-sm text-gray-700">
                                <span class="font-medium">{{ item.author }}</span>
                                {{ item.action }}
                                <span class="font-medium">{{ item.subject }}</span>
                            </p>
                            <div class="flex items-center gap-2 mt-1">
                                <span class="text-xs text-gray-400">{{ item.time }}</span>
                                <span :class="['px-2 py-0.5 rounded-full text-xs font-medium', item.tagClass]">{{ item.tag }}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="col-span-2 bg-white rounded-2xl shadow-sm p-5">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Actions rapides</h2>
                <p class="text-xs text-gray-400 mb-4">Accès direct à vos actions</p>
                <div class="grid grid-cols-2 gap-3">
                    <button
                        @click="showDepot = true"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-primary text-white hover:bg-primary-dark transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                        </svg>
                        <span class="text-sm font-medium text-center">Déposer un objet</span>
                    </button>

                    <button
                        @click="showDepot = true"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-primary text-primary hover:bg-primary/5 transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                        </svg>
                        <span class="text-sm font-medium text-center">+ Publier une annonce</span>
                    </button>

                    <RouterLink
                        to="/formations"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-secondary text-white hover:bg-secondary-dark transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                        </svg>
                        <span class="text-sm font-medium text-center">+ Voir les formations</span>
                    </RouterLink>

                    <RouterLink
                        to="/projets"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-secondary text-secondary hover:bg-secondary/5 transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                        </svg>
                        <span class="text-sm font-medium text-center">+ Rejoindre un projet</span>
                    </RouterLink>
                </div>
            </div>

        </div>

        <CreateAnnouncementModal v-model="showDepot" />

    </UserLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import UserLayout from '@/Layouts/UserLayout.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const auth = useAuthStore()
const showDepot = ref(false)

const stats = ref([
    {
        label: 'Upcycling Score',
        value: '—',
        key: 'upcycling_score',
        badge: null,
        bgClass: 'bg-secondary/20',
        iconClass: 'text-secondary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/></svg>`,
    },
    {
        label: 'Annonces actives',
        value: '—',
        key: 'active_announcements',
        badge: null,
        bgClass: 'bg-primary/10',
        iconClass: 'text-primary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"/></svg>`,
    },
    {
        label: 'Dépôts en attente',
        value: '—',
        key: 'pending_deposits',
        badge: 'En attente',
        bgClass: 'bg-primary/10',
        iconClass: 'text-primary',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>`,
    },
    {
        label: 'Événements à venir',
        value: '—',
        key: 'upcoming_events',
        badge: null,
        bgClass: 'bg-blue-100',
        iconClass: 'text-blue-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
    },
])

const activity = []

onMounted(async () => {
    try {
        const { data } = await api.get('/user/stats')
        stats.value.forEach(s => {
            if (data[s.key] !== undefined) s.value = data[s.key]
        })
    } catch {}
})
</script>
