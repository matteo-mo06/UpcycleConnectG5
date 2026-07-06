<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Publicités</h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="loadError" class="py-12 text-center text-sm text-red-500">{{ loadError }}</div>

        <template v-else>

            <div v-if="actionError" class="mb-4 px-4 py-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-600 flex items-center justify-between">
                <span>{{ actionError }}</span>
                <button @click="actionError = ''" class="ml-3 text-red-400 hover:text-red-600">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                </button>
            </div>

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

            <!-- Plans tarifaires -->
            <div class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
                <div class="px-5 py-4 border-b border-gray-100">
                    <h2 class="font-semibold text-gray-800">Plans tarifaires</h2>
                </div>
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Durée (semaines)</th>
                                <th class="text-left text-white font-medium px-5 py-3">Prix (€)</th>
                                <th class="text-right text-white font-medium px-5 py-3">Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(plan, i) in adPlans" :key="plan.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']">
                                <td class="px-5 py-3">
                                    <input v-model.number="plan.weeks" type="number" min="1"
                                        class="w-24 px-2 py-1 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"/>
                                </td>
                                <td class="px-5 py-3">
                                    <input
                                        type="number" min="0.01" step="0.01"
                                        :value="(plan.price_cents / 100).toFixed(2)"
                                        @input="plan.price_cents = Math.round($event.target.value * 100)"
                                        class="w-32 px-2 py-1 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"/>
                                </td>
                                <td class="px-5 py-3 text-right">
                                    <button @click="savePlan(plan)"
                                        class="px-3 py-1 text-xs font-medium text-white bg-primary rounded-lg hover:bg-primary-dark transition-colors">
                                        Enregistrer
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
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
                            @input="onSearchInput"
                            type="text"
                            placeholder="Rechercher par titre ou artisan…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>

                    <select
                        v-model="filterState"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option value="pending">En attente</option>
                        <option value="approved">Approuvée</option>
                        <option value="active">Active</option>
                        <option value="rejected">Refusée</option>
                        <option value="expired">Expirée</option>
                    </select>

                    <span class="text-xs text-gray-400 whitespace-nowrap">
                        {{ total }} résultat(s)
                    </span>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Aperçu</th>
                                <th class="text-left text-white font-medium px-5 py-3">Artisan</th>
                                <th class="text-left text-white font-medium px-5 py-3">Titre</th>
                                <th class="text-right text-white font-medium px-5 py-3">Prix</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-left text-white font-medium px-5 py-3">Date</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(ad, i) in advertisements"
                                :key="ad.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">

                                <td class="px-5 py-3">
                                    <img
                                        :src="ad.image_url"
                                        :alt="ad.title"
                                        class="w-16 h-10 object-cover rounded-lg cursor-pointer hover:opacity-80 transition-opacity"
                                        @click="openDetail(ad)" />
                                </td>

                                <td class="px-5 py-3">
                                    <p class="font-medium text-gray-800 truncate max-w-36">{{ ad.user_first_name }} {{ ad.user_last_name }}</p>
                                    <p class="text-xs text-gray-400 truncate max-w-36">{{ ad.user_email }}</p>
                                </td>

                                <td class="px-5 py-3">
                                    <p class="text-gray-800 font-medium truncate max-w-48">{{ ad.title }}</p>
                                    <a v-if="ad.link_url" :href="ad.link_url" target="_blank" rel="noopener noreferrer"
                                        class="text-xs text-primary hover:underline truncate block max-w-48">{{ ad.link_url }}</a>
                                </td>

                                <td class="px-5 py-3 text-right font-medium text-gray-800">{{ formatEuros(ad.price_cents) }}</td>

                                <td class="px-5 py-3">
                                    <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                        <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(ad.state)" />
                                        <span :class="statusText(ad.state)">{{ statusLabel(ad.state) }}</span>
                                    </span>
                                </td>

                                <td class="px-5 py-3 text-xs text-gray-500">
                                    {{ formatDate(ad.created_at) }}
                                    <span v-if="ad.expires_at" class="block text-amber-500 mt-0.5">Expire {{ formatDate(ad.expires_at) }}</span>
                                </td>

                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button
                                            @click="openDetail(ad)"
                                            title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>

                                        <template v-if="ad.state === 'pending'">
                                            <button
                                                @click="approveAd(ad)"
                                                title="Approuver"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                </svg>
                                            </button>
                                            <button
                                                @click="openReject(ad)"
                                                title="Refuser"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                </svg>
                                            </button>
                                        </template>

                                        <button
                                            v-if="ad.state === 'active'"
                                            @click="confirmDeactivate(ad)"
                                            title="Désactiver"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
                                            </svg>
                                        </button>
                                        <button
                                            v-if="ad.state === 'expired'"
                                            @click="doReactivate(ad)"
                                            title="Réactiver"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-green-600 hover:bg-green-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>

                            <tr v-if="advertisements.length === 0">
                                <td colspan="7" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucune publicité ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ total }} publicité(s) au total
                </div>
            </div>

            <Pagination v-if="total > 20" :page="page" :total="total" :limit="20" @update:page="changePage" />

        </template>

        <!-- Detail modal -->
        <div
            v-if="detailAd"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="detailAd = null">
            <div class="absolute inset-0 bg-black/40" @click="detailAd = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4 flex-shrink-0">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                            {{ detailAd.title }}
                        </h3>
                        <p class="text-xs text-gray-500 mt-0.5">{{ detailAd.user_first_name }} {{ detailAd.user_last_name }} · {{ detailAd.user_email }}</p>
                    </div>
                    <button
                        @click="detailAd = null"
                        class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
                    <img :src="detailAd.image_url" :alt="detailAd.title" class="w-full h-48 object-cover rounded-xl" />

                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Prix</p>
                            <p class="font-semibold text-gray-800">{{ formatEuros(detailAd.price_cents) }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailAd.state)" />
                                <span :class="statusText(detailAd.state)">{{ statusLabel(detailAd.state) }}</span>
                            </span>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Soumis le</p>
                            <p class="font-medium text-gray-800">{{ formatDate(detailAd.created_at) }}</p>
                        </div>
                        <div v-if="detailAd.approved_at">
                            <p class="text-xs text-gray-400 mb-0.5">Approuvé le</p>
                            <p class="font-medium text-gray-800">{{ formatDate(detailAd.approved_at) }}</p>
                        </div>
                        <div v-if="detailAd.paid_at">
                            <p class="text-xs text-gray-400 mb-0.5">Payé le</p>
                            <p class="font-medium text-gray-800">{{ formatDate(detailAd.paid_at) }}</p>
                        </div>
                        <div v-if="detailAd.expires_at">
                            <p class="text-xs text-gray-400 mb-0.5">Expire le</p>
                            <p class="font-medium text-gray-800">{{ formatDate(detailAd.expires_at) }}</p>
                        </div>
                        <div v-if="detailAd.link_url" class="col-span-2">
                            <p class="text-xs text-gray-400 mb-0.5">Lien</p>
                            <a :href="detailAd.link_url" target="_blank" rel="noopener noreferrer"
                                class="text-sm text-primary hover:underline break-all">{{ detailAd.link_url }}</a>
                        </div>
                    </div>

                    <div v-if="detailAd.rejection_reason" class="p-3 rounded-xl bg-red-50 border border-red-100">
                        <p class="text-xs font-medium text-red-600 mb-0.5">Motif de refus</p>
                        <p class="text-sm text-red-700">{{ detailAd.rejection_reason }}</p>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-between gap-2 flex-shrink-0">
                    <div class="flex gap-2">
                        <button
                            v-if="detailAd.state === 'pending'"
                            @click="approveAd(detailAd); detailAd = null"
                            class="px-3 py-1.5 text-sm font-medium text-secondary border border-secondary/30 rounded-lg hover:bg-secondary/5 transition-colors">
                            Approuver
                        </button>
                        <button
                            v-if="detailAd.state === 'pending'"
                            @click="openReject(detailAd); detailAd = null"
                            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Refuser
                        </button>
                        <button
                            v-if="detailAd.state === 'active'"
                            @click="confirmDeactivate(detailAd); detailAd = null"
                            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Désactiver
                        </button>
                    </div>
                    <button
                        @click="detailAd = null"
                        class="px-3 py-1.5 text-sm text-gray-500 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Fermer
                    </button>
                </div>
            </div>
        </div>

        <!-- Reject modal -->
        <div
            v-if="rejectTarget"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="rejectTarget = null">
            <div class="absolute inset-0 bg-black/40" @click="rejectTarget = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800">Refuser la publicité</h3>
                    <button @click="rejectTarget = null" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-3">
                    <p class="text-sm text-gray-600">Publicité : <span class="font-medium text-gray-800">{{ rejectTarget.title }}</span></p>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Motif de refus <span class="text-gray-400 font-normal">(optionnel)</span></label>
                        <textarea
                            v-model="rejectReason"
                            rows="3"
                            placeholder="Expliquez pourquoi la publicité est refusée…"
                            class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-red-300 focus:border-red-400 resize-none"
                        />
                    </div>
                    <p v-if="rejectError" class="text-sm text-red-500">{{ rejectError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button
                        @click="rejectTarget = null; rejectError = ''"
                        class="px-4 py-2 text-sm text-gray-500 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="doReject"
                        :disabled="rejectSubmitting"
                        class="px-4 py-2 text-sm font-medium text-white bg-red-500 rounded-lg hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ rejectSubmitting ? 'Envoi…' : 'Confirmer le refus' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Deactivate confirm modal -->
        <div
            v-if="deactivateTarget"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="deactivateTarget = null">
            <div class="absolute inset-0 bg-black/40" @click="deactivateTarget = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm">
                <div class="px-6 py-5 space-y-2">
                    <h3 class="font-semibold text-gray-800">Désactiver la publicité</h3>
                    <p class="text-sm text-gray-500">
                        Voulez-vous désactiver <span class="font-medium text-gray-700">« {{ deactivateTarget.title }} »</span> ?
                        Elle ne sera plus visible sur l'accueil.
                    </p>
                    <p v-if="deactivateError" class="text-sm text-red-500">{{ deactivateError }}</p>
                </div>
                <div class="px-6 pb-5 flex justify-end gap-2">
                    <button
                        @click="deactivateTarget = null; deactivateError = ''"
                        class="px-4 py-2 text-sm text-gray-500 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="doDeactivate"
                        :disabled="deactivateSubmitting"
                        class="px-4 py-2 text-sm font-medium text-white bg-red-500 rounded-lg hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deactivateSubmitting ? 'Désactivation…' : 'Désactiver' }}
                    </button>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const loading = ref(true)
const loadError = ref('')
const actionError = ref('')
const adPlans = ref([])
const advertisements = ref([])
const total = ref(0)
const page = ref(1)
const search = ref('')
const filterState = ref('')

const detailAd = ref(null)
const rejectTarget = ref(null)
const rejectReason = ref('')
const rejectError = ref('')
const rejectSubmitting = ref(false)
const deactivateTarget = ref(null)
const deactivateError = ref('')
const deactivateSubmitting = ref(false)

const statsData = ref({ total: 0, pending: 0, approved: 0, active: 0, rejected: 0, expired: 0 })

const stats = computed(() => [
    {
        label: 'Total publicités',
        value: statsData.value.total,
        bgClass: 'bg-blue-100',
        iconClass: 'w-6 h-6 text-blue-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"/></svg>`,
    },
    {
        label: 'En attente',
        value: statsData.value.pending,
        bgClass: 'bg-yellow-100',
        iconClass: 'w-6 h-6 text-yellow-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Actives',
        value: statsData.value.active,
        bgClass: 'bg-green-100',
        iconClass: 'w-6 h-6 text-green-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Refusées',
        value: statsData.value.rejected,
        bgClass: 'bg-red-100',
        iconClass: 'w-6 h-6 text-red-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
])

function formatEuros(cents) {
    if (!cents) return '0,00 €'
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

function formatDate(str) {
    if (!str) return '-'
    return new Date(str).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function statusDot(state) {
    return {
        pending:  'bg-yellow-400',
        approved: 'bg-blue-400',
        active:   'bg-green-500',
        rejected: 'bg-red-500',
        expired:  'bg-gray-400',
    }[state] ?? 'bg-gray-400'
}

function statusText(state) {
    return {
        pending:  'text-yellow-700',
        approved: 'text-blue-700',
        active:   'text-green-700',
        rejected: 'text-red-600',
        expired:  'text-gray-500',
    }[state] ?? 'text-gray-500'
}

function statusLabel(state) {
    return {
        pending:  'En attente',
        approved: 'Approuvée',
        active:   'Active',
        rejected: 'Refusée',
        expired:  'Expirée',
    }[state] ?? state
}

function openDetail(ad) {
    detailAd.value = ad
}

function openReject(ad) {
    rejectTarget.value = ad
    rejectReason.value = ''
    rejectError.value = ''
}

function confirmDeactivate(ad) {
    deactivateTarget.value = ad
    deactivateError.value = ''
}

async function fetchAds() {
    const params = { page: page.value, limit: 20 }
    if (search.value.trim()) params.search = search.value.trim()
    if (filterState.value) params.state = filterState.value
    const { data } = await api.get('/admin/advertisements', { params })
    advertisements.value = data.advertisements ?? []
    total.value = data.total ?? 0
}

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(() => { page.value = 1; fetchAds() }, 300)
}

async function fetchStats() {
    const { data } = await api.get('/admin/advertisements/stats')
    statsData.value = data
}

async function reload() {
    await Promise.all([fetchAds(), fetchStats()])
}

async function approveAd(ad) {
    try {
        await api.patch(`/admin/advertisement/${ad.id}/approve`)
        await reload()
    } catch (err) {
        actionError.value = err?.response?.data?.error ?? 'Erreur lors de l\'approbation.'
    }
}

async function doReject() {
    if (!rejectTarget.value) return
    rejectError.value = ''
    rejectSubmitting.value = true
    try {
        await api.patch(`/admin/advertisement/${rejectTarget.value.id}/reject`, { reason: rejectReason.value })
        rejectTarget.value = null
        rejectReason.value = ''
        await reload()
    } catch (err) {
        rejectError.value = err?.response?.data?.error ?? 'Erreur lors du refus.'
    } finally {
        rejectSubmitting.value = false
    }
}

async function fetchAdPlans() {
    try {
        const { data } = await api.get('/admin/advertisement/plans')
        adPlans.value = data ?? []
    } catch (err) {
        console.error('fetchAdPlans error:', err)
    }
}

async function savePlan(plan) {
    try {
        await api.put(`/admin/advertisement/plan/${plan.id}`, { weeks: plan.weeks, price_cents: plan.price_cents })
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de la sauvegarde du plan.')
    }
}

async function doReactivate(ad) {
    try {
        await api.patch(`/admin/advertisement/${ad.id}/reactivate`)
        await reload()
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de la réactivation.')
    }
}

async function doDeactivate() {
    if (!deactivateTarget.value) return
    deactivateError.value = ''
    deactivateSubmitting.value = true
    try {
        await api.patch(`/admin/advertisement/${deactivateTarget.value.id}/deactivate`)
        deactivateTarget.value = null
        await reload()
    } catch (err) {
        deactivateError.value = err?.response?.data?.error ?? 'Erreur lors de la désactivation.'
    } finally {
        deactivateSubmitting.value = false
    }
}

async function changePage(p) {
    page.value = p
    try {
        await fetchAds()
    } catch (err) {
        actionError.value = err?.response?.data?.error ?? 'Erreur lors du chargement de la page.'
    }
}

watch(filterState, () => { page.value = 1; fetchAds() })

onMounted(async () => {
    try {
        await Promise.all([reload(), fetchAdPlans()])
    } catch {
        loadError.value = 'Impossible de charger les publicités. Vérifiez votre connexion et rechargez la page.'
    } finally {
        loading.value = false
    }
})
</script>
