<template>
    <SalarieLayout>
        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Abonnement Premium</h1>
            <p class="text-sm text-gray-400 mt-1">Accédez aux fonctionnalités avancées d'UpcycleConnect</p>
        </div>

        <div v-if="subscription" class="bg-white rounded-2xl shadow-sm p-6 mb-6 border border-secondary/20">
            <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-full bg-secondary/10 flex items-center justify-center flex-shrink-0">
                    <svg class="w-6 h-6 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                </div>
                <div class="flex-1">
                    <p class="font-semibold text-gray-800">{{ subscription.plan?.name }} actif</p>
                    <p class="text-sm text-gray-500 mt-0.5">
                        <span v-if="subscription.end_timestamp">
                            Valide jusqu'au {{ formatDate(subscription.end_timestamp) }}
                        </span>
                        <span v-else>Renouvellement automatique</span>
                    </p>
                </div>
                <button @click="goToPortal" :disabled="portalLoading"
                    class="px-4 py-2 text-sm font-medium text-gray-600 border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors disabled:opacity-50">
                    {{ portalLoading ? 'Chargement...' : 'Gérer mon abonnement' }}
                </button>
            </div>
        </div>

        <div v-if="successMsg" class="mb-6 p-4 rounded-xl bg-green-50 text-green-700 text-sm font-medium">
            {{ successMsg }}
        </div>
        <div v-if="cancelledMsg" class="mb-6 p-4 rounded-xl bg-amber-50 text-amber-700 text-sm font-medium">
            {{ cancelledMsg }}
        </div>

        <div v-if="!subscription" class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-6">
            <div v-for="plan in plans" :key="plan.id"
                class="bg-white rounded-2xl shadow-sm p-6 flex flex-col gap-4 border border-gray-100 hover:border-primary/30 transition-colors">
                <div>
                    <h2 class="text-xl font-bold text-gray-800">{{ plan.name }}</h2>
                    <p class="text-3xl font-bold text-primary mt-2">
                        {{ formatPrice(plan.price_cents) }}
                        <span class="text-sm font-normal text-gray-400">/ {{ plan.interval_unit === 'month' ? 'mois' : 'an' }}</span>
                    </p>
                </div>
                <ul class="space-y-2 flex-1">
                    <li v-for="feature in features" :key="feature" class="flex items-center gap-2 text-sm text-gray-600">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                        </svg>
                        {{ feature }}
                    </li>
                </ul>
                <button @click="subscribe(plan)" :disabled="!plan.stripe_price_id || checkoutLoading === plan.id"
                    class="w-full py-2.5 text-sm font-semibold bg-primary text-white rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    <span v-if="checkoutLoading === plan.id">Redirection...</span>
                    <span v-else-if="!plan.stripe_price_id">Bientôt disponible</span>
                    <span v-else>S'abonner</span>
                </button>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6">
            <h2 class="font-semibold text-gray-800 mb-4">Ce que comprend le Premium</h2>
            <div class="grid grid-cols-2 gap-4">
                <div v-for="feature in featureDetails" :key="feature.title" class="flex items-start gap-3">
                    <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0 mt-0.5" v-html="feature.icon"/>
                    <div>
                        <p class="text-sm font-medium text-gray-800">{{ feature.title }}</p>
                        <p class="text-xs text-gray-500 mt-0.5">{{ feature.desc }}</p>
                    </div>
                </div>
            </div>
        </div>
    </SalarieLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SalarieLayout from '@/Layouts/SalarieLayout.vue'
import api from '@/api.js'

const route = useRoute()

const plans = ref([])
const subscription = ref(null)
const checkoutLoading = ref(null)
const portalLoading = ref(false)

const successMsg = computed(() => route.query.success ? 'Votre abonnement a été activé avec succès !' : '')
const cancelledMsg = computed(() => route.query.cancelled ? 'Paiement annulé. Votre abonnement n\'a pas été modifié.' : '')

const features = [
    'Achats d\'annonces illimités',
    'Annonces illimitées',
    'Projets illimités',
    'Tableau de bord avancé',
    'Upcycling Score détaillé',
]

const featureDetails = [
    {
        title: 'Achats illimités',
        desc: 'Achetez autant d\'annonces que vous souhaitez sans restriction.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/></svg>`,
    },
    {
        title: 'Annonces illimitées',
        desc: 'Publiez autant d\'annonces que nécessaire.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>`,
    },
    {
        title: 'Tableau de bord avancé',
        desc: 'Statistiques détaillées sur vos ventes et projets.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>`,
    },
    {
        title: 'Upcycling Score',
        desc: 'Mesurez votre impact écologique en détail.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
]

function formatPrice(cents) {
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

function formatDate(dateStr) {
    return new Date(dateStr).toLocaleDateString('fr-FR')
}

async function subscribe(plan) {
    checkoutLoading.value = plan.id
    try {
        const { data } = await api.post('/user/subscription/checkout', { plan_id: plan.id })
        window.location.href = data.url
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la redirection vers le paiement.')
        checkoutLoading.value = null
    }
}

async function goToPortal() {
    portalLoading.value = true
    try {
        const { data } = await api.post('/user/subscription/portal')
        window.location.href = data.url
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de l\'ouverture du portail.')
        portalLoading.value = false
    }
}

onMounted(async () => {
    const [plansRes, subRes] = await Promise.all([
        api.get('/subscription/plans').catch(() => ({ data: [] })),
        api.get('/user/subscription').catch(() => ({ data: {} })),
    ])
    plans.value = (plansRes.data ?? []).filter(p => p.is_active)
    subscription.value = subRes.data?.active ? subRes.data.subscription : null
})
</script>
