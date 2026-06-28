<template>
    <AdminLayout>
        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Revenus</h1>
            <div class="flex gap-2">
                <button
                    :class="['px-4 py-2 rounded-lg text-sm font-medium transition-colors', activeTab === 'transactions' ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50']"
                    @click="activeTab = 'transactions'"
                >Transactions</button>
                <button
                    :class="['px-4 py-2 rounded-lg text-sm font-medium transition-colors', activeTab === 'abonnements' ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50']"
                    @click="activeTab = 'abonnements'"
                >Abonnements</button>
                <button
                    :class="['px-4 py-2 rounded-lg text-sm font-medium transition-colors', activeTab === 'parametres' ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50']"
                    @click="activeTab = 'parametres'"
                >Paramètres</button>
            </div>
        </div>

        <div class="grid grid-cols-4 gap-5 mb-8">
            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-green-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                </div>
                <div class="min-w-0">
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ formatEuros(summary.total_amount_cents) }}</p>
                    <p class="text-sm text-gray-500 mt-1">Volume total</p>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-blue-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 14l6-6m-5.5.5h.01m4.99 5h.01M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16l3.5-2 3.5 2 3.5-2 3.5 2z"/>
                    </svg>
                </div>
                <div class="min-w-0">
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ formatEuros(summary.total_commission_cents) }}</p>
                    <p class="text-sm text-gray-500 mt-1">Commissions perçues</p>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-purple-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                    </svg>
                </div>
                <div class="min-w-0">
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ summary.total_transactions ?? 0 }}</p>
                    <p class="text-sm text-gray-500 mt-1">Transactions</p>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-orange-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-orange-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"/>
                    </svg>
                </div>
                <div class="min-w-0">
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ summary.commission_rate ?? 0 }}%</p>
                    <p class="text-sm text-gray-500 mt-1">Taux de commission</p>
                </div>
            </div>
        </div>

        <div v-if="activeTab === 'transactions'">
            <div v-if="loadingTx" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
            <div v-else-if="transactions.length === 0" class="py-12 text-center text-sm text-gray-400">Aucune transaction pour le moment.</div>
            <div v-else class="bg-white rounded-xl shadow-sm overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Annonce</th>
                                <th class="text-left text-white font-medium px-5 py-3">Acheteur</th>
                                <th class="text-right text-white font-medium px-5 py-3">Montant</th>
                                <th class="text-right text-white font-medium px-5 py-3">Commission</th>
                                <th class="text-left text-white font-medium px-5 py-3">Date</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(tx, i) in transactions"
                                :key="tx.id_payement"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']"
                            >
                                <td class="px-5 py-3 text-gray-700 max-w-xs truncate">
                                    {{ tx.announcement_title || '-' }}
                                </td>
                                <td class="px-5 py-3 text-gray-600">{{ tx.buyer_name || '-' }}</td>
                                <td class="px-5 py-3 text-right font-medium text-gray-800">{{ formatEuros(tx.amount_cents) }}</td>
                                <td class="px-5 py-3 text-right">
                                    <span class="bg-green-50 text-green-700 px-2 py-0.5 rounded-full text-xs font-semibold">
                                        +{{ formatEuros(tx.commission_amount_cents) }}
                                    </span>
                                </td>
                                <td class="px-5 py-3 text-gray-500 text-xs">{{ formatDate(tx.created_at) }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="px-5 py-3 border-t border-gray-100 flex items-center justify-between text-xs text-gray-400">
                    <span>{{ total }} transaction(s)</span>
                    <div class="flex gap-2">
                        <button
                            :disabled="page <= 1"
                            @click="changePage(page - 1)"
                            class="px-3 py-1 rounded border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
                        >Précédent</button>
                        <span class="px-3 py-1">{{ page }} / {{ totalPages }}</span>
                        <button
                            :disabled="page >= totalPages"
                            @click="changePage(page + 1)"
                            class="px-3 py-1 rounded border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
                        >Suivant</button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="activeTab === 'abonnements'">
            <div v-if="loadingSub" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
            <div v-else-if="subscriptions.length === 0" class="py-12 text-center text-sm text-gray-400">Aucun abonnement pour le moment.</div>
            <div v-else class="bg-white rounded-xl shadow-sm overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Utilisateur</th>
                                <th class="text-left text-white font-medium px-5 py-3">Email</th>
                                <th class="text-left text-white font-medium px-5 py-3">Plan</th>
                                <th class="text-right text-white font-medium px-5 py-3">Prix</th>
                                <th class="text-left text-white font-medium px-5 py-3">Début</th>
                                <th class="text-left text-white font-medium px-5 py-3">Fin</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(sub, i) in subscriptions"
                                :key="sub.user_id + sub.start_date"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']"
                            >
                                <td class="px-5 py-3 text-gray-700 font-medium">{{ sub.first_name }} {{ sub.last_name }}</td>
                                <td class="px-5 py-3 text-gray-500 text-xs">{{ sub.email }}</td>
                                <td class="px-5 py-3 text-gray-700">{{ sub.plan_name }}</td>
                                <td class="px-5 py-3 text-right font-medium text-gray-800">{{ formatEuros(sub.price_cents) }}</td>
                                <td class="px-5 py-3 text-gray-500 text-xs">{{ formatDate(sub.start_date) }}</td>
                                <td class="px-5 py-3 text-gray-500 text-xs">{{ sub.end_date ? formatDate(sub.end_date) : '-' }}</td>
                                <td class="px-5 py-3">
                                    <span v-if="sub.cancelled" class="bg-red-50 text-red-600 px-2 py-0.5 rounded-full text-xs font-semibold">Annulé</span>
                                    <span v-else-if="sub.end_date && new Date(sub.end_date) < new Date()" class="bg-orange-50 text-orange-600 px-2 py-0.5 rounded-full text-xs font-semibold">Expiré</span>
                                    <span v-else class="bg-green-50 text-green-700 px-2 py-0.5 rounded-full text-xs font-semibold">Actif</span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="px-5 py-3 border-t border-gray-100 flex items-center justify-between text-xs text-gray-400">
                    <span>{{ totalSub }} abonnement(s)</span>
                    <div class="flex gap-2">
                        <button
                            :disabled="pageSub <= 1"
                            @click="changeSubPage(pageSub - 1)"
                            class="px-3 py-1 rounded border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
                        >Précédent</button>
                        <span class="px-3 py-1">{{ pageSub }} / {{ totalSubPages }}</span>
                        <button
                            :disabled="pageSub >= totalSubPages"
                            @click="changeSubPage(pageSub + 1)"
                            class="px-3 py-1 rounded border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
                        >Suivant</button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="activeTab === 'parametres'">
            <div class="bg-white rounded-xl shadow-sm p-6 max-w-md">
                <h2 class="font-semibold text-gray-800 mb-1">Taux de commission</h2>
                <p class="text-sm text-gray-500 mb-4">
                    Pourcentage prélevé par UpcycleConnect sur chaque vente d'annonce. S'applique aux prochaines transactions.
                </p>
                <div class="flex items-center gap-3">
                    <div class="relative flex-1">
                        <input
                            v-model.number="editRate"
                            type="number"
                            min="0"
                            max="100"
                            step="0.5"
                            class="w-full px-3 py-2 pr-8 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"
                        />
                        <span class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm">%</span>
                    </div>
                    <button
                        @click="saveRate"
                        :disabled="savingRate"
                        class="px-4 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-60"
                    >
                        {{ savingRate ? 'Enregistrement…' : 'Enregistrer' }}
                    </button>
                </div>
                <p v-if="rateSuccess" class="text-xs text-green-600 mt-2">Taux mis à jour avec succès.</p>
                <p v-if="rateError" class="text-xs text-red-500 mt-2">{{ rateError }}</p>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const activeTab = ref('transactions')

const summary = ref({ total_transactions: 0, total_amount_cents: 0, total_commission_cents: 0, commission_rate: 0 })
const transactions = ref([])
const total = ref(0)
const page = ref(1)
const limit = 20
const loadingTx = ref(true)

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)))

const subscriptions = ref([])
const totalSub = ref(0)
const pageSub = ref(1)
const limitSub = 20
const loadingSub = ref(false)

const totalSubPages = computed(() => Math.max(1, Math.ceil(totalSub.value / limitSub)))

const editRate = ref(0)
const savingRate = ref(false)
const rateSuccess = ref(false)
const rateError = ref('')

function formatEuros(cents) {
    if (!cents) return '0,00 €'
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

function formatDate(str) {
    if (!str) return '-'
    return new Date(str).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

async function fetchSummary() {
    try {
        const { data } = await api.get('/admin/revenue/summary')
        summary.value = data
        editRate.value = data.commission_rate
    } catch { /* silencieux */ }
}

async function fetchTransactions() {
    loadingTx.value = true
    try {
        const { data } = await api.get('/admin/revenue/transactions', { params: { page: page.value, limit } })
        transactions.value = data.transactions ?? []
        total.value = data.total ?? 0
    } catch { /* silencieux */ } finally {
        loadingTx.value = false
    }
}

async function changePage(p) {
    page.value = p
    await fetchTransactions()
}

async function fetchSubscriptions() {
    loadingSub.value = true
    try {
        const { data } = await api.get('/admin/subscriptions', { params: { page: pageSub.value, limit: limitSub } })
        subscriptions.value = data.subscriptions ?? []
        totalSub.value = data.total ?? 0
    } catch { /* silencieux */ } finally {
        loadingSub.value = false
    }
}

async function changeSubPage(p) {
    pageSub.value = p
    await fetchSubscriptions()
}

async function saveRate() {
    rateSuccess.value = false
    rateError.value = ''
    if (editRate.value < 0 || editRate.value > 100) {
        rateError.value = 'Le taux doit être entre 0 et 100.'
        return
    }
    savingRate.value = true
    try {
        await api.put('/admin/revenue/commission-rate', { commission_rate: editRate.value })
        summary.value.commission_rate = editRate.value
        rateSuccess.value = true
    } catch {
        rateError.value = 'Erreur lors de la mise à jour.'
    } finally {
        savingRate.value = false
    }
}

onMounted(async () => {
    await Promise.all([fetchSummary(), fetchTransactions(), fetchSubscriptions()])
})
</script>
