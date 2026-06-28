<template>
    <AdminLayout>
        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Abonnements</h1>
        </div>

        <div class="grid grid-cols-3 gap-5 mb-8">
            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-green-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                    </svg>
                </div>
                <div>
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ activeSubscribersCount }}</p>
                    <p class="text-sm text-gray-500 mt-1">Abonnés actifs</p>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-purple-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"/>
                    </svg>
                </div>
                <div>
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ plans.length }}</p>
                    <p class="text-sm text-gray-500 mt-1">Plans disponibles</p>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-blue-100 flex items-center justify-center">
                    <svg class="w-6 h-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                </div>
                <div>
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ activePlansCount }}</p>
                    <p class="text-sm text-gray-500 mt-1">Plans actifs</p>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
            <div class="px-5 py-4 border-b border-gray-100">
                <h2 class="font-semibold text-gray-800">Plans d'abonnement</h2>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead>
                        <tr class="bg-primary">
                            <th class="text-left text-white font-medium px-5 py-3">Nom</th>
                            <th class="text-left text-white font-medium px-5 py-3">Prix</th>
                            <th class="text-left text-white font-medium px-5 py-3">Intervalle</th>
                            <th class="text-left text-white font-medium px-5 py-3">Stripe</th>
                            <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                            <th class="text-left text-white font-medium px-5 py-3">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(plan, i) in plans" :key="plan.id"
                            :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']">
                            <td class="px-5 py-3 font-medium text-gray-800">{{ plan.name }}</td>
                            <td class="px-5 py-3 text-gray-700">{{ formatPrice(plan.price_cents) }}</td>
                            <td class="px-5 py-3 text-gray-500">{{ plan.interval_count }} {{ plan.interval_unit === 'month' ? 'mois' : 'an' }}</td>
                            <td class="px-5 py-3">
                                <span v-if="plan.stripe_price_id" class="text-xs text-green-600 font-medium">Synchronisé</span>
                                <span v-else class="text-xs text-amber-500">Non synchronisé</span>
                            </td>
                            <td class="px-5 py-3">
                                <span :class="plan.is_active ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'"
                                    class="px-2 py-0.5 rounded-full text-xs font-semibold">
                                    {{ plan.is_active ? 'Actif' : 'Inactif' }}
                                </span>
                            </td>
                            <td class="px-5 py-3">
                                <button @click="openEdit(plan)"
                                    class="text-primary hover:underline text-xs font-medium">
                                    Modifier
                                </button>
                            </td>
                        </tr>
                        <tr v-if="plans.length === 0">
                            <td colspan="6" class="px-5 py-10 text-center text-gray-400 text-sm">Aucun plan</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm overflow-hidden">
            <div class="px-5 py-4 border-b border-gray-100">
                <h2 class="font-semibold text-gray-800">Abonnés</h2>
            </div>
            <div v-if="subscriptions.length === 0" class="py-10 text-center text-sm text-gray-400">Aucun abonnement pour le moment.</div>
            <div v-else class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead>
                        <tr class="bg-primary">
                            <th class="text-left text-white font-medium px-5 py-3">Utilisateur</th>
                            <th class="text-left text-white font-medium px-5 py-3">Email</th>
                            <th class="text-left text-white font-medium px-5 py-3">Plan</th>
                            <th class="text-left text-white font-medium px-5 py-3">Début</th>
                            <th class="text-left text-white font-medium px-5 py-3">Fin</th>
                            <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(sub, i) in subscriptions" :key="sub.user_id + sub.plan_name"
                            :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']">
                            <td class="px-5 py-3 font-medium text-gray-800">{{ sub.first_name }} {{ sub.last_name }}</td>
                            <td class="px-5 py-3 text-gray-500">{{ sub.email }}</td>
                            <td class="px-5 py-3 text-gray-700">{{ sub.plan_name }}</td>
                            <td class="px-5 py-3 text-gray-500 text-xs">{{ sub.start_date ? sub.start_date.slice(0, 10) : '-' }}</td>
                            <td class="px-5 py-3 text-gray-500 text-xs">{{ sub.end_date ? sub.end_date.slice(0, 10) : '-' }}</td>
                            <td class="px-5 py-3">
                                <span :class="sub.cancelled ? 'bg-red-50 text-red-600' : 'bg-green-50 text-green-700'"
                                    class="px-2 py-0.5 rounded-full text-xs font-semibold">
                                    {{ sub.cancelled ? 'Annulé' : 'Actif' }}
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div v-if="modal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="modal = false"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md p-6 flex flex-col gap-4">
                <h3 class="font-semibold text-gray-800">{{ editTarget ? 'Modifier le plan' : 'Nouveau plan' }}</h3>

                <div class="grid grid-cols-2 gap-4">
                    <div class="col-span-2">
                        <label class="text-xs text-gray-500 mb-1 block">Nom</label>
                        <input v-model="form.name" type="text"
                            class="w-full px-3 py-2 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"/>
                    </div>
                    <div>
                        <label class="text-xs text-gray-500 mb-1 block">Prix (centimes)</label>
                        <input v-model.number="form.price_cents" type="number" min="0"
                            class="w-full px-3 py-2 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"/>
                        <p class="text-xs text-gray-400 mt-0.5">= {{ formatPrice(form.price_cents) }}</p>
                    </div>
                    <div>
                        <label class="text-xs text-gray-500 mb-1 block">Intervalle</label>
                        <select v-model="form.interval_unit"
                            class="w-full px-3 py-2 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary">
                            <option value="month">Mensuel</option>
                            <option value="year">Annuel</option>
                        </select>
                    </div>
                    <div class="col-span-2 flex items-center gap-2">
                        <input v-model="form.is_active" type="checkbox" id="is_active" class="rounded"/>
                        <label for="is_active" class="text-sm text-gray-700">Plan actif</label>
                    </div>
                </div>

                <p v-if="error" class="text-sm text-red-500">{{ error }}</p>

                <div class="flex gap-2 justify-end">
                    <button @click="modal = false"
                        class="px-4 py-2 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitPlan" :disabled="submitting"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50">
                        {{ submitting ? 'Synchronisation Stripe...' : 'Enregistrer' }}
                    </button>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const plans = ref([])
const subscriptions = ref([])
const modal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const error = ref('')

const form = ref({ name: '', price_cents: 1500, interval_unit: 'month', interval_count: 1, is_active: true })

const activeSubscribersCount = computed(() => subscriptions.value.filter(s => !s.cancelled).length)
const activePlansCount = computed(() => plans.value.filter(p => p.is_active).length)

function formatPrice(cents) {
    if (!cents) return '0,00 €'
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

function openEdit(plan) {
    editTarget.value = plan
    form.value = {
        name: plan.name,
        price_cents: plan.price_cents,
        interval_unit: plan.interval_unit,
        interval_count: plan.interval_count,
        is_active: plan.is_active,
    }
    error.value = ''
    modal.value = true
}

async function submitPlan() {
    if (!form.value.name || !form.value.price_cents) {
        error.value = 'Nom et prix sont requis.'
        return
    }
    submitting.value = true
    error.value = ''
    const payload = { ...form.value, interval_count: 1 }
    try {
        if (editTarget.value) {
            await api.put(`/admin/subscription/plan/${editTarget.value.id}`, payload)
        } else {
            await api.post('/admin/subscription/plans', payload)
        }
        modal.value = false
        await fetchPlans()
    } catch (e) {
        error.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

async function fetchPlans() {
    try {
        const { data } = await api.get('/admin/subscription/plans')
        plans.value = data ?? []
    } catch (err) {
        console.error('fetchPlans error:', err)
    }
}

async function fetchSubscriptions() {
    try {
        const { data } = await api.get('/admin/subscriptions')
        subscriptions.value = data.subscriptions ?? []
    } catch (err) {
        console.error('fetchSubscriptions error:', err)
    }
}

onMounted(() => {
    fetchPlans()
    fetchSubscriptions()
})
</script>
