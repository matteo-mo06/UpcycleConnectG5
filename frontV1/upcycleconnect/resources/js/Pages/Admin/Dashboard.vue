<template>
  <AdminLayout title="Dashboard">

    <div class="grid grid-cols-4 gap-5 mb-8">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
        <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
          <div :class="stat.iconClass" v-html="stat.icon" />
        </div>
        <div class="min-w-0">
          <p class="text-2xl font-bold text-gray-800 leading-none">{{stat.value}}</p>
          <p class="text-sm text-gray-500 mt-1 truncate">{{stat.label}}</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-5 gap-6">

      <div class="col-span-3 bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2 class="text-base font-semibold text-gray-800">Derniers utilisateurs</h2>
          <a href="/admin/users" class="text-xs font-medium text-primary hover:underline">Voir tout</a>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="bg-primary">
                <th class="text-left text-white font-medium px-5 py-3">Nom</th>
                <th class="text-left text-white font-medium px-5 py-3">Email</th>
                <th class="text-left text-white font-medium px-5 py-3">Type</th>
                <th class="text-left text-white font-medium px-5 py-3">Inscription</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(user, i) in recentUsers"
                :key="user.email"
                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                <td class="px-5 py-3 font-medium text-gray-800">
                  <div class="flex items-center gap-2">
                    <div class="w-7 h-7 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
                      <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
                      </svg>
                    </div>
                    {{user.name}}
                  </div>
                </td>
                <td class="px-5 py-3 text-gray-500">{{user.email}}</td>
                <td class="px-5 py-3">
                  <span
                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="typeBadge(user.type)">
                    {{user.type}}
                  </span>
                </td>
                <td class="px-5 py-3 text-gray-500">{{user.date}}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="col-span-2 bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2 class="text-base font-semibold text-gray-800">Prochains événements</h2>
          <a href="/admin/events" class="text-xs font-medium text-primary hover:underline">Voir tout</a>
        </div>
        <div class="divide-y divide-gray-50">
          <div
            v-for="event in upcomingEvents"
            :key="event.id"
            class="px-5 py-3.5 flex items-center justify-between hover:bg-gray-50/60 transition-colors duration-100"
          >
            <div class="min-w-0 mr-3">
              <p class="font-medium text-gray-800 text-sm truncate">{{ event.title }}</p>
              <p class="text-xs text-gray-500 mt-0.5">{{ event.date }} · {{ event.location }}</p>
            </div>
            <span
              class="flex-shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
              :class="categoryBadge(event.category)"
            >
              {{ event.category }}
            </span>
          </div>
        </div>
      </div>

    </div>

  </AdminLayout>
</template>

<script setup>
import AdminLayout from '@/Layouts/AdminLayout.vue'

const stats = [
  {
    label: 'Utilisateurs total',
    value: 1248,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
  },
  {
    label: 'Annonces actives',
    value: 342,
    bgClass: 'bg-green-100',
    iconClass: 'text-green-600',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
  },
  {
    label: 'Événements à venir',
    value: 9,
    bgClass: 'bg-amber-100',
    iconClass: 'text-amber-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
  },
  {
    label: 'Signalements ouverts',
    value: 5,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
  },
]

const recentUsers = [
  { name: 'Particulier', email: 'particulier@test.fr', type: 'Particulier', date: '01/01/2025' },
  { name: 'Artisan', email: 'artisan@test.fr', type: 'Artisan', date: '01/01/2025' },
  { name: 'Salarié', email: 'salarie@test.fr', type: 'Salarié', date: '01/01/2025' },
  { name: 'Administrateur', email: 'admin@test.fr', type: 'Administrateur', date: '01/01/2025' },
  { name: 'Utilisateur', email: 'user@test.fr', type: 'Particulier', date: '01/01/2025' },
]

const upcomingEvents = [
  { id: 1, title: 'Événement test', date: '01/01/2025', location: 'Salle', category: 'Catégorie' },
  { id: 2, title: 'Événement test', date: '01/01/2025', location: 'Salle', category: 'Catégorie' },
  { id: 3, title: 'Événement test', date: '01/01/2025', location: 'Salle', category: 'Catégorie' },
  { id: 4, title: 'Événement test', date: '01/01/2025', location: 'Salle', category: 'Catégorie' },
  { id: 5, title: 'Événement test', date: '01/01/2025', location: 'Salle', category: 'Catégorie' },
]

function categoryBadge(category) {
  const map = {
    'Catégorie': 'bg-blue-100 text-blue-700',
  }
  return map[category] ?? 'bg-gray-100 text-gray-600'
}

function typeBadge(type) {
  const map = {
    'Particulier': 'bg-blue-100 text-blue-700',
    'Artisan': 'bg-secondary/20 text-secondary-dark',
    'Salarié': 'bg-purple-100 text-purple-700',
    'Administrateur': 'bg-primary/15 text-primary-dark',
  }
  return map[type] ?? 'bg-gray-100 text-gray-600'
}
</script>
