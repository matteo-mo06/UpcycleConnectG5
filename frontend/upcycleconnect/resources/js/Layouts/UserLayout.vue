<template>
    <div class="flex h-screen overflow-hidden bg-bg">

        <aside class="flex flex-col flex-shrink-0 w-64 bg-sidebar">
            <div class="flex items-center h-16 px-4 border-b border-white/10">
                <div class="flex items-center gap-3">
                    <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
                        <span class="text-white font-bold text-sm">U</span>
                    </div>
                    <span class="text-white font-semibold text-base whitespace-nowrap" style="font-family: var(--font-family-title)">
                        Upcycle Connect
                    </span>
                </div>
            </div>

            <nav class="flex-1 py-4 space-y-1 overflow-y-auto">
                <template v-for="item in navItems" :key="item.label">
                    <button
                        v-if="item.modal"
                        @click="$emit('openDepot')"
                        class="w-full flex items-center gap-3 px-4 py-2.5 mx-2 rounded-lg transition-colors duration-150 text-white/70 hover:text-white hover:bg-white/10"
                        style="width: calc(100% - 1rem)"
                    >
                        <div class="flex-shrink-0 w-5 h-5" v-html="item.icon" />
                        <span class="text-sm font-medium whitespace-nowrap">{{ item.label }}</span>
                    </button>
                    <RouterLink
                        v-else
                        :to="item.href"
                        :class="[
                            'flex items-center gap-3 px-4 py-2.5 mx-2 rounded-lg transition-colors duration-150',
                            isActive(item.href) ? 'bg-primary text-white' : 'text-white/70 hover:text-white hover:bg-white/10',
                        ]"
                    >
                        <div class="flex-shrink-0 w-5 h-5" v-html="item.icon" />
                        <span class="text-sm font-medium whitespace-nowrap">{{ item.label }}</span>
                    </RouterLink>
                </template>
            </nav>

            <div class="p-4 border-t border-white/10 flex items-center gap-3">
                <div class="w-9 h-9 rounded-full bg-primary flex items-center justify-center flex-shrink-0">
                    <span class="text-white text-sm font-semibold">{{ initials }}</span>
                </div>
                <div class="min-w-0 flex-1">
                    <p class="text-white text-sm font-medium truncate">{{ auth.user?.first_name }} {{ auth.user?.last_name }}</p>
                    <p class="text-white/50 text-xs truncate capitalize">{{ primaryRole }}</p>
                </div>
                <RouterLink
                    v-if="auth.isAdmin"
                    to="/admin/dashboard"
                    title="Interface admin"
                    class="flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white/50 hover:text-white hover:bg-white/10 transition-colors duration-150"
                >
                    <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    </svg>
                </RouterLink>
                <button
                    @click="handleLogout"
                    title="Déconnexion"
                    class="flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white/50 hover:text-white hover:bg-white/10 transition-colors duration-150"
                >
                    <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
                    </svg>
                </button>
            </div>
        </aside>

        <div class="flex flex-col flex-1 overflow-hidden">
            <main class="flex-1 overflow-y-auto p-6">
                <slot />
            </main>
        </div>

        <div v-if="toast" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/50" />
            <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden">
                <div class="h-1.5 w-full" :class="toast.type === 'warning' ? 'bg-amber-500' : 'bg-red-600'" />
                <div class="px-6 py-5 flex flex-col items-center text-center gap-4">
                    <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0"
                        :class="toast.type === 'warning' ? 'bg-amber-100' : 'bg-red-100'">
                        <svg class="w-6 h-6" :class="toast.type === 'warning' ? 'text-amber-500' : 'text-red-600'" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                        </svg>
                    </div>
                    <div>
                        <p class="font-semibold text-gray-800 text-base" style="font-family: var(--font-family-title)">
                            {{ toastTitle }}
                        </p>
                        <p class="mt-2 text-sm text-gray-600 leading-relaxed">{{ toast.message }}</p>
                    </div>
                    <button
                        @click="dismissToast"
                        class="w-full py-2.5 rounded-xl text-sm font-medium text-white transition-colors"
                        :class="toast.type === 'warning' ? 'bg-amber-500 hover:bg-amber-600' : 'bg-red-600 hover:bg-red-700'">
                        J'ai compris
                    </button>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import { initials as getInitials, primaryRole as getPrimaryRole } from '@/utils.js'

defineEmits(['openDepot'])

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const initials = computed(() => getInitials(auth.user))
const primaryRole = computed(() => getPrimaryRole(auth.user?.roles ?? []))

function isActive(href) {
    return route.path.startsWith(href)
}

function handleLogout() {
    auth.logout()
    router.push('/login')
}

const toast = ref(null)
let sseAbort = null

const toastTitles = {
    warning: 'Avertissement des modérateurs',
    suspension: 'Votre compte a été suspendu',
    ban: 'Votre compte a été banni',
}
const toastTitle = computed(() => toastTitles[toast.value?.sanctionType] ?? 'Action sur votre compte')

function showToast(message, type, sanctionType) {
    toast.value = { message, type, sanctionType }
}

function dismissToast() {
    const sanctionType = toast.value?.sanctionType
    toast.value = null
    if (sanctionType === 'suspension' || sanctionType === 'ban') {
        auth.logout()
        router.push('/login')
    }
}

async function connectSSE() {
    const token = sessionStorage.getItem('token')
    if (!token) return

    sseAbort = new AbortController()

    try {
        const response = await fetch('http://localhost:8084/sse', {
            headers: { Authorization: `Bearer ${token}` },
            signal: sseAbort.signal,
        })

        if (!response.ok || !response.body) {
            setTimeout(connectSSE, 5000)
            return
        }

        const reader = response.body.getReader()
        const decoder = new TextDecoder()

        while (true) {
            const { done, value } = await reader.read()
            if (done) break

            const lines = decoder.decode(value).split('\n')
            for (const line of lines) {
                if (!line.startsWith('data: ')) continue
                try {
                    const data = JSON.parse(line.slice(6))
                    if (data.type === 'sanction' && !auth.isAdmin) {
                        showToast(data.message, data.sanction_type === 'warning' ? 'warning' : 'error', data.sanction_type)
                    }
                } catch {}
            }
        }
    } catch (e) {
        if (e?.name === 'AbortError') return
    }

    setTimeout(connectSSE, 5000)
}

onMounted(connectSSE)
onUnmounted(() => {
    sseAbort?.abort()
})

const navItems = [
    {
        label: 'Accueil',
        href: '/accueil',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/></svg>`,
    },
    {
        label: 'Annonces',
        href: '/annonces',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/></svg>`,
    },
    {
        label: "Dépôt d'objet",
        href: '/depot',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>`,
    },
    {
        label: 'Projets',
        href: '/projets',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
    },
    {
        label: 'Forum',
        href: '/forum',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>`,
    },
    {
        label: 'Formations',
        href: '/formations',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>`,
    },
    {
        label: 'Événements',
        href: '/evenements',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
    },
    {
        label: 'Conseils',
        href: '/conseils',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/></svg>`,
    },
    {
        label: 'Paramètres',
        href: '/parametres',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
    },
]
</script>
