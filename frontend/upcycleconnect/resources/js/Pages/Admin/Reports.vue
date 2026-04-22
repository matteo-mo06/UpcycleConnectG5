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
          <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
          <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
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
          <option>Topic</option>
          <option>Post</option>
        </select>

        <select
          v-model="filterStatus"
          class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
          <option value="">Tous les statuts</option>
          <option>À traiter</option>
          <option>Résolu</option>
          <option>Ignoré</option>
        </select>

        <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">{{ filteredReports.length }} signalement(s)</span>
      </div>

      <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
      <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

      <div v-else class="overflow-x-auto">
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
                <p class="font-medium text-gray-800 truncate max-w-48">{{ report.contentTitle }}</p>
                <p class="text-xs text-gray-500 truncate">{{ report.contentAuthor }}</p>
              </td>

              <td class="px-5 py-3">
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="typeBadge(report.type)">
                  {{ report.type }}
                </span>
              </td>

              <td class="px-5 py-3">
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600 max-w-32 truncate">
                  {{ report.reason }}
                </span>
              </td>

              <td class="px-5 py-3 text-gray-500 text-xs">{{ report.reporter }}</td>

              <td class="px-5 py-3 text-gray-500 text-xs">{{ report.date }}</td>

              <td class="px-5 py-3">
                <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                  <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(report.status)" />
                  <span :class="statusText(report.status)">{{ report.status }}</span>
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
                    title="Ignorer"
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
        {{ filteredReports.length }} signalement(s) affiché(s)
      </div>
    </div>

    <div class="mt-8 bg-white rounded-xl shadow-sm overflow-hidden">
      <div class="px-5 py-4 border-b border-gray-100 flex items-center gap-3">
        <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">Historique utilisateurs</h2>
        <div class="relative flex-1 max-w-xs">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
          </svg>
          <input
            v-model="histoSearch"
            @input="debouncedHistoFetch"
            type="text"
            placeholder="Rechercher par nom ou email…"
            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
          />
        </div>
      </div>

      <div v-if="histoLoading" class="py-8 text-center text-sm text-gray-400">Chargement…</div>

      <div v-else-if="histoUsers.length === 0" class="py-8 text-center text-sm text-gray-400">
        Aucun utilisateur trouvé.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-100">
              <th class="text-left text-gray-500 font-medium px-5 py-3">Utilisateur</th>
              <th class="text-left text-gray-500 font-medium px-5 py-3">Email</th>
              <th class="text-left text-gray-500 font-medium px-5 py-3">Statut</th>
              <th class="text-center text-gray-500 font-medium px-5 py-3">Signalements reçus</th>
              <th class="text-center text-gray-500 font-medium px-5 py-3">Sanctions</th>
              <th class="px-5 py-3" />
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(u, i) in histoUsers"
              :key="u.id_user"
              :class="['border-b border-gray-50 cursor-pointer hover:bg-gray-50/70 transition-colors', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/30']"
              @click="openHistoModal(u)">
              <td class="px-5 py-3 font-medium text-gray-800">{{ u.name }}</td>
              <td class="px-5 py-3 text-gray-500 text-xs">{{ u.email }}</td>
              <td class="px-5 py-3">
                <span :class="['inline-flex px-2 py-0.5 rounded-full text-xs font-medium', userStatusBadge(u.status)]">
                  {{ userStatusLabel(u.status) }}
                </span>
              </td>
              <td class="px-5 py-3 text-center">
                <span :class="['inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-semibold', u.report_count > 0 ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-400']">
                  {{ u.report_count }}
                </span>
              </td>
              <td class="px-5 py-3 text-center">
                <span :class="['inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-semibold', u.sanction_count > 0 ? 'bg-amber-100 text-amber-700' : 'bg-gray-100 text-gray-400']">
                  {{ u.sanction_count }}
                </span>
              </td>
              <td class="px-5 py-3 text-right">
                <span class="text-xs text-primary hover:underline">Voir →</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div
      v-if="detailReport"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="closeDetailModal">
      <div class="absolute inset-0 bg-black/40" @click="closeDetailModal" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">

        <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
          <div class="min-w-0">
            <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
              Signalement #{{ detailReport.id.slice(0, 8) }}
            </h3>
            <p class="text-xs text-gray-500 mt-0.5">Reçu le {{ detailReport.date }}</p>
          </div>
          <button
            @click="closeDetailModal"
            class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="px-6 py-5 space-y-4 max-h-[70vh] overflow-y-auto">

          <div class="bg-gray-50 rounded-xl p-4">
            <div class="flex items-start gap-3">
              <span class="flex-shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium mt-0.5" :class="typeBadge(detailReport.type)">
                {{ detailReport.type }}
              </span>
              <div class="min-w-0">
                <p class="font-medium text-gray-800 text-sm">{{ detailReport.contentTitle }}</p>
                <p class="text-xs text-gray-500 mt-0.5">Par {{ detailReport.contentAuthor }}</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Signalé par</p>
              <p class="font-medium text-gray-800">{{ detailReport.reporter }}</p>
              <p class="text-xs text-gray-500">{{ detailReport.reporterEmail }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400 mb-0.5">Statut actuel</p>
              <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailReport.status)" />
                <span :class="statusText(detailReport.status)">{{ detailReport.status }}</span>
              </span>
            </div>
          </div>

          <div>
            <p class="text-xs text-gray-400 mb-1">Motif du signalement</p>
            <p class="text-sm text-gray-700 leading-relaxed italic">"{{ detailReport.reason }}"</p>
          </div>

          <div v-if="detailReport.status === 'À traiter'">
            <p class="text-xs text-gray-400 mb-2">Actions</p>
            <div class="flex flex-wrap gap-2">
              <button
                @click="removeContent(detailReport)"
                class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
                Supprimer le contenu
              </button>
              <template v-if="detailReport.contentAuthorId">
                <button
                  @click="openSanction('warning', detailReport)"
                  class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-amber-700 bg-amber-50 border border-amber-200 rounded-lg hover:bg-amber-100 transition-colors">
                  Avertir l'auteur
                </button>
                <button
                  @click="openSanction('suspension', detailReport)"
                  class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-orange-700 bg-orange-50 border border-orange-200 rounded-lg hover:bg-orange-100 transition-colors">
                  Suspendre
                </button>
                <button
                  @click="openSanction('ban', detailReport)"
                  class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors">
                  Bannir
                </button>
              </template>
            </div>

            <div v-if="sanctionType && sanctionTarget && sanctionTarget.source === 'report'" class="mt-3 pt-3 border-t border-gray-100 space-y-2">
              <p class="text-xs font-medium text-gray-600">
                {{ sanctionLabels[sanctionType] }} — {{ sanctionTarget.name }}
              </p>
              <div v-if="sanctionType === 'suspension'" class="flex items-center gap-2">
                <label class="text-xs text-gray-500">Durée :</label>
                <input v-model.number="sanctionDays" type="number" min="1" max="365"
                  class="w-20 text-sm border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary/30" />
                <span class="text-xs text-gray-500">jours</span>
              </div>
              <textarea
                v-model="sanctionReason"
                rows="2"
                placeholder="Raison de la sanction…"
                class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none" />
              <div class="flex gap-2">
                <button @click="sanctionType = ''; sanctionTarget = null"
                  class="px-3 py-1.5 text-xs border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors">
                  Annuler
                </button>
                <button @click="submitSanction" :disabled="sanctionLoading"
                  class="px-3 py-1.5 text-xs font-medium bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors disabled:opacity-60">
                  {{ sanctionLoading ? '…' : 'Confirmer' }}
                </button>
              </div>
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

    <div
      v-if="histoUser"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="histoUser = null">
      <div class="absolute inset-0 bg-black/40" @click="histoUser = null" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden">

        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <div>
            <div class="flex items-center gap-2">
              <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ histoUser.name }}</h3>
              <span :class="['inline-flex px-2 py-0.5 rounded-full text-xs font-medium', userStatusBadge(histoUser.status)]">
                {{ userStatusLabel(histoUser.status) }}
              </span>
            </div>
            <p class="text-xs text-gray-500 mt-0.5">{{ histoUser.email }}</p>
          </div>
          <button @click="histoUser = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="px-6 py-5 space-y-5 max-h-[65vh] overflow-y-auto">

          <div>
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Sanctions ({{ histoUser.sanctions.length }})</p>
            <div v-if="histoUser.sanctions.length === 0" class="text-xs text-gray-400">Aucune sanction.</div>
            <div v-else class="space-y-1.5">
              <div v-for="s in histoUser.sanctions" :key="s.id_sanction"
                class="flex items-start gap-3 text-xs text-gray-600 bg-gray-50 rounded-lg px-3 py-2">
                <span :class="['flex-shrink-0 font-medium px-2 py-0.5 rounded-full text-xs', sanctionBadge(s.type)]">
                  {{ sanctionLabels[s.type] ?? s.type }}
                </span>
                <div class="min-w-0 flex-1">
                  <p class="text-gray-700">{{ s.reason || '—' }}</p>
                  <p class="text-gray-400 mt-0.5">{{ s.created_at?.slice(0, 10) }} · par {{ s.admin_name }}</p>
                </div>
              </div>
            </div>
          </div>

          <div>
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Signalements reçus ({{ histoUser.reportsReceived.length }})</p>
            <div v-if="histoUser.reportsReceived.length === 0" class="text-xs text-gray-400">Aucun signalement reçu.</div>
            <div v-else class="space-y-1.5">
              <div v-for="r in histoUser.reportsReceived" :key="r.id_report"
                class="flex items-start gap-3 text-xs text-gray-600 bg-gray-50 rounded-lg px-3 py-2">
                <span class="flex-shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="typeBadge(contentTypeLabel(r.content_type))">
                  {{ contentTypeLabel(r.content_type) }}
                </span>
                <div class="min-w-0">
                  <p class="font-medium text-gray-700 truncate">{{ r.content_title || '—' }}</p>
                  <p class="text-gray-500 truncate">{{ r.reason }}</p>
                  <p class="text-gray-400 mt-0.5">{{ r.created_at?.slice(0, 10) }} ·
                    <span :class="statusText(statusLabel(r.status))">{{ statusLabel(r.status) }}</span>
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div class="border-t border-gray-100 pt-4">
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Sanctionner cet utilisateur</p>
            <div class="flex flex-wrap gap-2">
              <button @click="openSanctionForUser('warning')"
                class="px-3 py-1.5 text-xs font-medium text-amber-700 bg-amber-50 border border-amber-200 rounded-lg hover:bg-amber-100 transition-colors">
                Avertir
              </button>
              <button @click="openSanctionForUser('suspension')"
                class="px-3 py-1.5 text-xs font-medium text-orange-700 bg-orange-50 border border-orange-200 rounded-lg hover:bg-orange-100 transition-colors">
                Suspendre
              </button>
              <button @click="openSanctionForUser('ban')"
                class="px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors">
                Bannir
              </button>
              <button
                v-if="histoUser.status === 'suspended' || histoUser.status === 'blacklisted'"
                @click="liftSanction"
                class="px-3 py-1.5 text-xs font-medium text-green-700 bg-green-50 border border-green-200 rounded-lg hover:bg-green-100 transition-colors">
                Lever la sanction
              </button>
            </div>

            <div v-if="sanctionType && sanctionTarget && sanctionTarget.source === 'histo'" class="mt-3 space-y-2">
              <p class="text-xs font-medium text-gray-600">{{ sanctionLabels[sanctionType] }}</p>
              <div v-if="sanctionType === 'suspension'" class="flex items-center gap-2">
                <label class="text-xs text-gray-500">Durée :</label>
                <input v-model.number="sanctionDays" type="number" min="1" max="365"
                  class="w-20 text-sm border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary/30" />
                <span class="text-xs text-gray-500">jours</span>
              </div>
              <textarea
                v-model="sanctionReason"
                rows="2"
                placeholder="Raison de la sanction…"
                class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none" />
              <div class="flex gap-2">
                <button @click="sanctionType = ''; sanctionTarget = null"
                  class="px-3 py-1.5 text-xs border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors">
                  Annuler
                </button>
                <button @click="submitSanction" :disabled="sanctionLoading"
                  class="px-3 py-1.5 text-xs font-medium bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors disabled:opacity-60">
                  {{ sanctionLoading ? '…' : 'Confirmer' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const reports = ref([])
const stats = ref([
  {
    label: 'Total signalements', key: 'total', value: 0, bgClass: 'bg-red-100', iconClass: 'text-red-500',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"/></svg>`,
  },
  {
    label: 'À traiter', key: 'pending', value: 0, bgClass: 'bg-red-100', iconClass: 'text-red-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
  },
  {
    label: 'Résolus', key: 'resolved', value: 0, bgClass: 'bg-green-100', iconClass: 'text-green-600',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
  },
  {
    label: 'Ignorés', key: 'ignored', value: 0, bgClass: 'bg-gray-100', iconClass: 'text-gray-400',
    icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>`,
  },
])

const loading = ref(false)
const error = ref('')
const search = ref('')
const filterType = ref('')
const filterStatus = ref('')
const detailReport = ref(null)

const sanctionType = ref('')
const sanctionReason = ref('')
const sanctionDays = ref(3)
const sanctionTarget = ref(null)
const sanctionLoading = ref(false)

const sanctionLabels = { warning: 'Avertissement', suspension: 'Suspension', ban: 'Bannissement' }

const histoSearch = ref('')
const histoUsers = ref([])
const histoLoading = ref(false)
const histoUser = ref(null)

const filteredReports = computed(() => {
  return reports.value.filter(r => {
    const q = search.value.toLowerCase()
    if (q && !r.contentTitle.toLowerCase().includes(q) && !r.reporter.toLowerCase().includes(q)) return false
    if (filterType.value && r.type !== filterType.value) return false
    if (filterStatus.value && r.status !== filterStatus.value) return false
    return true
  })
})

function mapReport(r) {
  return {
    id: r.id_report,
    contentAuthorId: r.id_reported_user,
    type: { announcement: 'Annonce', topic: 'Topic', post: 'Post' }[r.content_type] ?? r.content_type,
    contentTitle: r.content_title || r.id_announcement || r.id_topic || r.id_post || '—',
    contentAuthor: r.reported_user_name || '—',
    reporter: r.reporter_name || '—',
    reporterEmail: r.reporter_email || '',
    reason: r.reason,
    date: r.created_at?.slice(0, 10) ?? '—',
    status: statusLabel(r.status),
  }
}

function statusLabel(s) {
  return { pending: 'À traiter', resolved: 'Résolu', ignored: 'Ignoré' }[s] ?? s
}

function contentTypeLabel(t) {
  return { announcement: 'Annonce', topic: 'Topic', post: 'Post' }[t] ?? t
}

async function fetchStats() {
  try {
    const { data } = await api.get('/admin/reports/stats')
    stats.value.forEach(s => { if (data[s.key] !== undefined) s.value = data[s.key] })
  } catch {}
}

async function fetchReports() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await api.get('/admin/reports')
    reports.value = data.map(mapReport)
  } catch {
    error.value = 'Impossible de charger les signalements.'
  } finally {
    loading.value = false
  }
}

let histoTimer = null
function debouncedHistoFetch() {
  clearTimeout(histoTimer)
  histoTimer = setTimeout(fetchHistoUsers, 300)
}

async function fetchHistoUsers() {
  histoLoading.value = true
  try {
    const params = {}
    if (histoSearch.value) params.search = histoSearch.value
    const { data } = await api.get('/admin/reports/users', { params })
    histoUsers.value = data
  } catch {
    histoUsers.value = []
  } finally {
    histoLoading.value = false
  }
}

async function openHistoModal(u) {
  try {
    const { data } = await api.get(`/admin/user/${u.id_user}/history`)
    histoUser.value = {
      id: u.id_user,
      name: u.name,
      email: u.email,
      status: data.status ?? 'active',
      sanctions: data.sanctions ?? [],
      reportsReceived: data.reports_received ?? [],
    }
  } catch {
    alert('Erreur lors du chargement de l\'historique.')
  }
}

async function liftSanction() {
  try {
    await api.patch(`/admin/user/${histoUser.value.id}/status`, { status: 'active' })
    histoUser.value.status = 'active'
    const idx = histoUsers.value.findIndex(u => u.id_user === histoUser.value.id)
    if (idx !== -1) histoUsers.value[idx].status = 'active'
  } catch (e) {
    alert(e.response?.data?.error ?? 'Erreur lors de la levée de sanction.')
  }
}

function openDetail(report) {
  sanctionType.value = ''
  sanctionTarget.value = null
  detailReport.value = reports.value.find(r => r.id === report.id) ?? report
}

function closeDetailModal() {
  detailReport.value = null
  sanctionType.value = ''
  sanctionTarget.value = null
}

async function ignoreReport(report) {
  try {
    await api.patch(`/admin/report/${report.id}`, { status: 'ignored' })
    report.status = 'Ignoré'
    await fetchStats()
  } catch (e) {
    alert(e.response?.data?.error ?? 'Erreur lors de la mise à jour.')
  }
}

async function resolveReport(report) {
  try {
    await api.patch(`/admin/report/${report.id}`, { status: 'resolved' })
    report.status = 'Résolu'
    await fetchStats()
  } catch (e) {
    alert(e.response?.data?.error ?? 'Erreur lors de la mise à jour.')
  }
}

async function removeContent(report) {
  try {
    await api.delete(`/admin/report/${report.id}/content`)
    report.status = 'Résolu'
    detailReport.value = null
    await fetchStats()
  } catch (e) {
    alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
  }
}

function openSanction(type, report) {
  sanctionType.value = type
  sanctionReason.value = ''
  sanctionDays.value = 3
  sanctionTarget.value = { id: report.contentAuthorId, name: report.contentAuthor, idReport: report.id, source: 'report' }
}

function openSanctionForUser(type) {
  sanctionType.value = type
  sanctionReason.value = ''
  sanctionDays.value = 3
  sanctionTarget.value = { id: histoUser.value.id, name: histoUser.value.name, source: 'histo' }
}

async function submitSanction() {
  if (!sanctionTarget.value || !sanctionType.value) return
  sanctionLoading.value = true
  const savedType = sanctionType.value
  const payload = {
    type: savedType,
    reason: sanctionReason.value,
    duration_days: savedType === 'suspension' ? sanctionDays.value : 0,
  }
  if (sanctionTarget.value.idReport) {
    payload.id_report = sanctionTarget.value.idReport
  }
  try {
    await api.post(`/admin/user/${sanctionTarget.value.id}/sanctions`, payload)
    if (sanctionTarget.value.idReport) {
      await api.patch(`/admin/report/${sanctionTarget.value.idReport}`, { status: 'resolved', action_taken: savedType })
      const r = reports.value.find(x => x.id === sanctionTarget.value.idReport)
      if (r) r.status = 'Résolu'
      if (detailReport.value?.id === sanctionTarget.value.idReport) {
        detailReport.value.status = 'Résolu'
      }
      await fetchStats()
    }
    if (histoUser.value) {
      await openHistoModal({ id_user: histoUser.value.id, name: histoUser.value.name, email: histoUser.value.email })
    }
    sanctionType.value = ''
    sanctionReason.value = ''
    sanctionTarget.value = null
  } catch (e) {
    alert(e.response?.data?.error ?? 'Erreur lors de la sanction.')
  } finally {
    sanctionLoading.value = false
  }
}

function userStatusBadge(status) {
  const map = { active: 'bg-green-100 text-green-700', suspended: 'bg-orange-100 text-orange-700', blacklisted: 'bg-red-100 text-red-700' }
  return map[status] ?? 'bg-gray-100 text-gray-500'
}

function userStatusLabel(status) {
  return { active: 'Actif', suspended: 'Suspendu', blacklisted: 'Banni' }[status] ?? status
}

function typeBadge(type) {
  const map = { 'Annonce': 'bg-blue-100 text-blue-700', 'Topic': 'bg-purple-100 text-purple-700', 'Post': 'bg-indigo-100 text-indigo-700' }
  return map[type] ?? 'bg-gray-100 text-gray-600'
}

function sanctionBadge(type) {
  const map = { warning: 'bg-amber-100 text-amber-700', suspension: 'bg-orange-100 text-orange-700', ban: 'bg-red-100 text-red-700' }
  return map[type] ?? 'bg-gray-100 text-gray-600'
}

function statusDot(status) {
  const map = { 'À traiter': 'bg-red-500', 'Résolu': 'bg-green-500', 'Ignoré': 'bg-gray-400' }
  return map[status] ?? 'bg-gray-300'
}

function statusText(status) {
  const map = { 'À traiter': 'text-red-600', 'Résolu': 'text-green-700', 'Ignoré': 'text-gray-500' }
  return map[status] ?? 'text-gray-500'
}

onMounted(async () => {
  await Promise.all([fetchReports(), fetchStats(), fetchHistoUsers()])
})
</script>
