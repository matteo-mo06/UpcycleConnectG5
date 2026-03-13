<template>
  <div class="flex h-screen overflow-hidden bg-bg">

    <aside class="flex flex-col flex-shrink-0 w-64 bg-sidebar">
      <div class="flex items-center h-16 px-4 border-b border-white/10">
        <div class="flex items-center gap-3">
          <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
            <span class="text-white font-bold text-sm">U</span>
          </div>
          <span class="text-white font-semibold text-base whitespace-nowrap" >UpcycleConnect</span>
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
              : 'text-white/70 hover:text-white hover:bg-white/10'
          ]">
          <div class="flex-shrink-0 w-5 h-5" v-html="item.icon" />
          <span class="text-sm font-medium whitespace-nowrap">
            {{item.label}}
          </span>
        </RouterLink>
      </nav>

      <div class="p-4 border-t border-white/10 flex items-center gap-3">
        <div class="w-9 h-9 rounded-full bg-gray-600 flex items-center justify-center flex-shrink-0">
          <svg class="w-5 h-5 text-gray-300" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
          </svg>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-white text-sm font-medium truncate">Administrateur</p>
          <p class="text-white/50 text-xs truncate">admin@upcycle.fr</p>
        </div>
        <RouterLink
          to="/"
          title="Retour au site"
          class="flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white/50 hover:text-white hover:bg-white/10 transition-colors duration-150">
          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" class="w-4 h-4">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
          </svg>
        </RouterLink>
      </div>
    </aside>

    <div class="flex flex-col flex-1 overflow-hidden">

      <header class="flex items-center h-16 px-6 bg-white border-b border-gray-200 flex-shrink-0 shadow-sm">
        <h1
          class="text-xl font-semibold text-gray-800">
          {{title}}
        </h1>
      </header>

      <main class="flex-1 overflow-y-auto p-6">
        <slot />
      </main>

    </div>
  </div>
</template>

<script setup>
import { RouterLink, useRoute } from 'vue-router'

defineProps({
  title: {
    type: String,
    default: 'Administration'
  }
})

const route = useRoute()

const navItems = [
  {
    label: 'Dashboard',
    href: '/admin/dashboard',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/></svg>`,
  },
  {
    label: 'Utilisateurs',
    href: '/admin/users',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
  },
  {
    label: 'Rôles',
    href: '/admin/roles',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>`,
  },
  {
    label: 'Annonces',
    href: '/admin/listings',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
  },
  {
    label: 'Événements',
    href: '/admin/events',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
  },
  {
    label: 'Signalements',
    href: '/admin/reports',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
  },
]

function isActive(href) {
  return route.path.startsWith(href)
}
</script>
