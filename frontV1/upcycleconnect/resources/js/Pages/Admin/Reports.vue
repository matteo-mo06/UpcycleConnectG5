<template>
  <AdminLayout title="Signalements">

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

    <div class="bg-white rounded-xl shadow-sm overflow-hidden">

      <div class="px-5 py-4 border-b border-gray-100 flex items-center gap-3">
        <div class="relative flex-1">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
          </svg>
          <input
            v-model="search"
            type="text"
            placeholder="Rechercher par contenu ou signalant…"
            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
          />
        </div>

        <select
          v-model="filterType"
          class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
          <option value="">Tous les types</option>
          <option>Annonce</option>
          <option>Utilisateur</option>
        </select>

        <select
          v-model="filterStatus"
          class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
          <option value="">Tous les statuts</option>
          <option>À traiter</option>
          <option>Résolu</option>
          <option>Ignoré</option>
        </select>

        <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">{{filteredReports.length}} signalement(s)</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-primary">
              <th class="text-left text-white font-medium px-5 py-3">Contenu signalé</th>
              <th class="text-left text-white font-medium px-5 py-3">Type</th>
              <th class="text-left text-white font-medium px-5 py-3">Motif</th>
              <th class="text-left text-white font-medium px-5 py-3">Signalé par</th>
              <th class="text-left text-white font-medium px-5 py-3">Date</th>
              <th class="text-left text-white font-medium px-5 py-3">Statut</th>
              <th class="text-right text-white font-medium px-5 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(report, i) in filteredReports"
              :key="report.id"
              :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
              <td class="px-5 py-3">
                <p class="font-medium text-gray-800 truncate max-w-48">{{report.contentTitle}}</p>
                <p class="text-xs text-gray-500 truncate">{{report.contentAuthor}}</p>
              </td>

              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="typeBadge(report.type)">
                  {{report.type}}
                </span>
              </td>

              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="reasonBadge(report.reason)">
                  {{report.reason}}
                </span>
              </td>

              <td class="px-5 py-3 text-gray-500 text-xs">{{report.reporter}}</td>

              <td class="px-5 py-3 text-gray-500 text-xs">{{report.date}}</td>

              <td class="px-5 py-3">
                <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                  <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(report.status)" />
                  <span :class="statusText(report.status)">{{report.status}}</span>
                </span>
              </td>

              <td class="px-5 py-3">
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click="openDetail(report)"
                    title="Voir le détail"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>

                  <button
                    v-if="report.status === 'À traiter'"
                    @click="ignoreReport(report)"
                    title="Ignorer le signalement"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                    </svg>
                  </button>

                  <button
                    v-if="report.status === 'À traiter'"
                    @click="resolveReport(report)"
                    title="Marquer comme résolu"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredReports.length === 0">
              <td colspan="7" class="px-5 py-12 text-center text-gray-400 text-sm">
                Aucun signalement ne correspond à vos filtres.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
        {{filteredReports.length}} signalement(s) affiché(s)
      </div>
    </div>

    <div
      v-if="detailReport"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="detailReport = null">
      <div class="absolute inset-0 bg-black/40" @click="detailReport = null" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">

        <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                Signalement #{{detailReport.id}}
              </h3>
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="reasonBadge(detailReport.reason)">
                {{detailReport.reason}}
              </span>
            </div>
            <p class="text-xs text-gray-500 mt-0.5">Reçu le {{detailReport.date}}</p>
          </div>
          <button
            @click="detailReport = null"
            class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="px-6 py-5 space-y-5">

          <div class="bg-gray-50 rounded-xl p-4">
            <div class="flex items-start gap-3">
              <span
                class="flex-shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium mt-0.5"
                :class="typeBadge(detailReport.type)">
                {{detailReport.type}}
              </span>
              <div class="min-w-0">
                <p class="font-medium text-gray-800 text-sm">{{detailReport.contentTitle}}</p>
                <p class="text-xs text-gray-500 mt-0.5">Par {{detailReport.contentAuthor}}</p>
                <p class="text-xs text-gray-600 mt-2 leading-relaxed">{{detailReport.contentPreview}}</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Signalé par</p>
              <p class="font-medium text-gray-800">{{detailReport.reporter}}</p>
              <p class="text-xs text-gray-500">{{detailReport.reporterEmail}}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Statut actuel</p>
              <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailReport.status)" />
                <span :class="statusText(detailReport.status)">{{detailReport.status}}</span>
              </span>
            </div>
          </div>

          <div>
            <p class="text-xs text-gray-400 mb-1">Message du signalement</p>
            <p class="text-sm text-gray-700 leading-relaxed italic">"{{detailReport.message}}"</p>
          </div>

          <div v-if="detailReport.status === 'À traiter'">
            <p class="text-xs text-gray-400 mb-2">Action sur le contenu</p>
            <div class="flex gap-2">
              <button
                @click="warnAuthor()"
                class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-amber-700 bg-amber-50 border border-amber-200 rounded-lg hover:bg-amber-100 transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>
                Avertir l'auteur
              </button>
              <button
                @click="removeContent(detailReport)"
                class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
                Supprimer le contenu
              </button>
            </div>
          </div>
        </div>

        <div class="px-6 py-4 border-t border-gray-100 flex justify-between items-center">
          <button
            v-if="detailReport.status === 'À traiter'"
            @click="ignoreReport(detailReport); detailReport = null"
            class="px-3 py-1.5 text-sm font-medium text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
            Ignorer
          </button>
          <div v-else />
          <button
            v-if="detailReport.status === 'À traiter'"
            @click="resolveReport(detailReport); detailReport = null"
            class="px-3 py-1.5 text-sm font-medium bg-secondary text-white rounded-lg hover:bg-secondary-dark transition-colors">
            Marquer comme résolu
          </button>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import {ref, computed} from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'

const reports = ref([
  {
    id: 1,
    type: 'Annonce',
    contentTitle: 'Annonce test',
    contentAuthor: 'Utilisateur 1',
    contentPreview: 'Description de l\'annonce signalée.',
    reporter: 'Utilisateur 2',
    reporterEmail: 'user2@test.fr',
    reason: 'Contenu inapproprié',
    date: '01/01/2025',
    status: 'À traiter',
    message: 'Message du signalement test.',
 },
  {
    id: 2,
    type: 'Utilisateur',
    contentTitle: 'Profil Utilisateur',
    contentAuthor: 'Utilisateur',
    contentPreview: 'Description du comportement signalé.',
    reporter: 'Utilisateur',
    reporterEmail: 'user@test.fr',
    reason: 'Spam',
    date: '01/01/2025',
    status: 'À traiter',
    message: 'Message du signalement test.',
 },
])

const stats = [
  {
    label: 'Total signalements',
    value: 2,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"/></svg>`,
 },
  {
    label: 'À traiter',
    value: 2,
    bgClass: 'bg-red-100',
    iconClass: 'text-red-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
 },
  {
    label: 'Résolus',
    value: 0,
    bgClass: 'bg-green-100',
    iconClass: 'text-green-600',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
 },
  {
    label: 'Ignorés',
    value: 0,
    bgClass: 'bg-gray-100',
    iconClass: 'text-gray-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>`,
 },
]

const search = ref('')
const filterType = ref('')
const filterStatus = ref('')

const filteredReports = computed(() => {
  return reports.value.filter(r => {
    const q = search.value.toLowerCase()
    if (q && !r.contentTitle.toLowerCase().includes(q) && !r.reporter.toLowerCase().includes(q)) return false
    if (filterType.value && r.type !== filterType.value) return false
    if (filterStatus.value && r.status !== filterStatus.value) return false
    return true
 })
})

const detailReport = ref(null)

function openDetail(report) {
  detailReport.value = reports.value.find(r => r.id === report.id)
}

function ignoreReport(report) {
  reports.value.find(r => r.id === report.id).status = 'Ignoré'
}

function resolveReport(report) {
  reports.value.find(r => r.id === report.id).status = 'Résolu'
}

function warnAuthor() {
  detailReport.value = null
}


function removeContent(report) {
  reports.value.find(r => r.id === report.id).status = 'Résolu'
  detailReport.value = null
}

function typeBadge(type) {
  const map = {
    'Annonce': 'bg-blue-100 text-blue-700',
    'Utilisateur': 'bg-purple-100 text-purple-700',
 }
  return map[type] ?? 'bg-gray-100 text-gray-600'
}

function reasonBadge(reason) {
  const map = {
    'Contenu inapproprié': 'bg-red-100 text-red-700',
    'Spam': 'bg-orange-100 text-orange-700',
    'Hors charte': 'bg-amber-100 text-amber-700',
    'Doublon': 'bg-gray-100 text-gray-600',
    'Comportement abusif': 'bg-rose-100 text-rose-700',
 }
  return map[reason] ?? 'bg-gray-100 text-gray-600'
}

function statusDot(status) {
  const map = {
    'À traiter': 'bg-red-500',
    'Résolu': 'bg-green-500',
    'Ignoré': 'bg-gray-400',
 }
  return map[status] ?? 'bg-gray-300'
}

function statusText(status) {
  const map = {
    'À traiter': 'text-red-600',
    'Résolu': 'text-green-700',
    'Ignoré': 'text-gray-500',
 }
  return map[status] ?? 'text-gray-500'
}
</script>
