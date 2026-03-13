<template>
  <AdminLayout title="Dashboard">

    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-5 mb-8">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4"
      >
        <div
          class="flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center"
          :style="{ backgroundColor: stat.color + '1a' }"
        >
          <div :style="{ color: stat.color }" v-html="stat.icon" />
        </div>
        <div class="min-w-0">
          <p
            class="text-2xl font-bold text-gray-800 leading-none"
            style="font-family: var(--font-family-title)"
          >
            {{ stat.value.toLocaleString('fr-FR') }}
          </p>
          <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-5 gap-6">

      <div class="xl:col-span-3 bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2
            class="text-base font-semibold text-gray-800"
            style="font-family: var(--font-family-title)"
          >
            Derniers utilisateurs
          </h2>
          <a href="/admin/users" class="text-xs font-medium text-primary hover:underline">
            Voir tout
          </a>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="bg-primary">
                <th class="text-left text-white font-medium px-5 py-3">Nom</th>
                <th class="text-left text-white font-medium px-5 py-3 hidden md:table-cell">Email</th>
                <th class="text-left text-white font-medium px-5 py-3">Type</th>
                <th class="text-left text-white font-medium px-5 py-3 hidden sm:table-cell">Inscription</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(user, i) in recentUsers"
                :key="user.email"
                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']"
              >
                <td class="px-5 py-3 font-medium text-gray-800">
                  <div class="flex items-center gap-2">
                    <div
                      class="w-7 h-7 rounded-full flex items-center justify-center text-white text-xs font-semibold flex-shrink-0"
                      :style="{ backgroundColor: user.avatarColor }"
                    >
                      {{ user.name.charAt(0).toUpperCase() }}
                    </div>
                    {{ user.name }}
                  </div>
                </td>
                <td class="px-5 py-3 text-gray-500 hidden md:table-cell">{{ user.email }}</td>
                <td class="px-5 py-3">
                  <span
                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="typeBadge(user.type)"
                  >
                    {{ user.type }}
                  </span>
                </td>
                <td class="px-5 py-3 text-gray-500 hidden sm:table-cell">{{ user.date }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="xl:col-span-2 bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2
            class="text-base font-semibold text-gray-800"
            style="font-family: var(--font-family-title)"
          >
            Prochains événements
          </h2>
          <a href="/admin/events" class="text-xs font-medium text-primary hover:underline">
            Voir tout
          </a>
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
              :class="eventBadge(event.category)"
            >
              {{ event.category }}
            </span>
          </div>
        </div>
        <div class="px-5 py-3 border-t border-gray-100">
          <a
            href="/admin/events"
            class="block w-full text-center py-2 rounded-lg text-sm font-medium text-primary border border-primary hover:bg-primary hover:text-white transition-colors duration-150"
          >
            Gérer tous les événements
          </a>
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
    color: '#c46d68',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:24px;height:24px"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
  },
  {
    label: 'Annonces actives',
    value: 342,
    color: '#8dc734',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:24px;height:24px"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
  },
  {
    label: 'Événements à venir',
    value: 9,
    color: '#f59e0b',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:24px;height:24px"><path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>`,
  },
  {
    label: 'Signalements ouverts',
    value: 5,
    color: '#ef4444',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:24px;height:24px"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
  },
]

const recentUsers = [
  { name: 'Camille Dupont', email: 'camille.dupont@mail.com', type: 'Particulier', date: '08/03/2026', avatarColor: '#c46d68' },
  { name: 'Marc Artisan', email: 'marc.artisan@pro.fr', type: 'Artisan', date: '07/03/2026', avatarColor: '#8dc734' },
  { name: 'Lucie Bernard', email: 'lucie.b@gmail.com', type: 'Particulier', date: '06/03/2026', avatarColor: '#6366f1' },
  { name: 'Thomas Salarié', email: 'thomas.s@upcycle.fr', type: 'Salarié', date: '05/03/2026', avatarColor: '#f59e0b' },
  { name: 'Sophie Martin', email: 'sophie.m@mail.com', type: 'Particulier', date: '04/03/2026', avatarColor: '#c46d68' },
]

const upcomingEvents = [
  { id: 1, title: 'Atelier couture zéro déchet', date: '12/03/2026', location: 'Salle A', category: 'Atelier' },
  { id: 2, title: 'Conférence upcycling textile', date: '15/03/2026', location: 'Amphithéâtre', category: 'Conférence' },
  { id: 3, title: 'Repair café mensuel', date: '20/03/2026', location: 'Espace commun', category: 'Animation' },
  { id: 4, title: 'Formation bois & récup', date: '25/03/2026', location: 'Atelier bois', category: 'Formation' },
  { id: 5, title: 'Marché de la seconde main', date: '29/03/2026', location: 'Parvis extérieur', category: 'Animation' },
]

function eventBadge(category) {
  const map = {
    'Atelier': 'bg-blue-100 text-blue-700',
    'Conférence': 'bg-purple-100 text-purple-700',
    'Animation': 'bg-green-100 text-green-700',
    'Formation': 'bg-amber-100 text-amber-700',
  }
  return map[category] ?? 'bg-gray-100 text-gray-600'
}

function typeBadge(type) {
  const map = {
    'Particulier': 'bg-blue-100 text-blue-700',
    'Artisan': 'bg-green-100 text-green-700',
    'Salarié': 'bg-purple-100 text-purple-700',
    'Administrateur': 'bg-red-100 text-red-700',
  }
  return map[type] ?? 'bg-gray-100 text-gray-600'
}
</script>
