<template>
    <div class="flex h-screen overflow-hidden bg-bg">
        <aside class="flex flex-col flex-shrink-0 w-64 bg-sidebar">
            <div class="flex items-center h-16 px-4 border-b border-white/10">
                <div class="flex items-center gap-3">
                    <div
                        class="flex-shrink-0 w-8 h-8 rounded-lg bg-primary flex items-center justify-center"
                    >
                        <span class="text-white font-bold text-sm">U</span>
                    </div>
                    <span
                        class="text-white font-semibold text-base whitespace-nowrap"
                        style="font-family: var(--font-family-title)"
                        >Upcycle Connect</span
                    >
                </div>
            </div>

            <nav class="flex-1 py-4 space-y-1 overflow-y-auto">
                <RouterLink
                    v-for="item in navItems"
                    :key="item.href"
                    :to="item.href"
                    :class="[
                        'flex items-center gap-3 px-4 py-2.5 mx-2 rounded-lg transition-colors duration-150',
                        isActive(item.href)
                            ? 'bg-primary text-white'
                            : 'text-white/70 hover:text-white hover:bg-white/10',
                    ]"
                >
                    <div class="flex-shrink-0 w-5 h-5" v-html="item.icon" />
                    <span class="text-sm font-medium whitespace-nowrap">
                        {{ item.label }}
                    </span>
                </RouterLink>
            </nav>

            <div class="p-4 border-t border-white/10 flex items-center gap-3">
                <div class="w-9 h-9 rounded-full flex items-center justify-center flex-shrink-0 overflow-hidden" style="background-color: #f9f0ef">
                    <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" alt="" class="w-full h-full object-cover" />
                    <span v-else class="text-sm font-semibold" style="color: #2D2D2D">{{ initials }}</span>
                </div>
                <div class="min-w-0 flex-1">
                    <p class="text-white text-sm font-medium truncate">
                        {{ auth.user?.first_name }} {{ auth.user?.last_name }}
                    </p>
                    <p class="text-white/50 text-xs truncate capitalize">{{ primaryRole }}</p>
                </div>
                <RouterLink
                    to="/accueil"
                    title="Espace utilisateur"
                    class="flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white/50 hover:text-white hover:bg-white/10 transition-colors duration-150"
                >
                    <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
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
    </div>
</template>

<script setup>
import { computed } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth.js";
import { initials as getInitials, primaryRole as getPrimaryRole } from "@/utils.js";

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();

const initials = computed(() => getInitials(auth.user));
const primaryRole = computed(() => getPrimaryRole(auth.user?.roles ?? []));

function handleLogout() {
    auth.logout();
    router.push("/login");
}

const navItems = [
    {
        label: "Accueil",
        href: "/salarie/dashboard",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/></svg>`,
    },
    {
        label: "Formations",
        href: "/salarie/formations",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>`,
    },
    {
        label: "Événements",
        href: "/salarie/evenements",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
    },
    {
        label: "Conseils",
        href: "/salarie/conseils",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/></svg>`,
    },
    {
        label: "Forum",
        href: "/salarie/forum",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>`,
    },
    {
        label: "Calendrier",
        href: "/salarie/calendrier",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/></svg>`,
    },
    {
        label: "Abonnement",
        href: "/abonnement",
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M15 9a2 2 0 10-4 0v5a2 2 0 01-2 2h6m-6-4h4m8 0a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
];

function isActive(href) {
    return route.path.startsWith(href);
}
</script>
