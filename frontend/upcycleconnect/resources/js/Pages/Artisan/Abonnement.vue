<template>
    <ArtisanLayout>
        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Abonnement</h1>
            <p class="text-sm text-gray-400 mt-1">Choisissez le palier qui correspond à votre activité</p>
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

        <div v-if="!polling" class="grid grid-cols-1 md:grid-cols-3 gap-5 mb-6">
            <div v-for="plan in plans" :key="plan.id"
                :class="['bg-white rounded-2xl shadow-sm p-6 flex flex-col gap-4 border transition-colors',
                    isCurrentPlan(plan) ? 'border-secondary ring-1 ring-secondary/30' : 'border-gray-100 hover:border-primary/30']">
                <div>
                    <div class="flex items-center gap-2 mb-1">
                        <h2 class="text-xl font-bold text-gray-800">{{ plan.name }}</h2>
                        <span v-if="isCurrentPlan(plan)" class="px-2 py-0.5 rounded-full text-xs font-semibold bg-secondary/10 text-secondary">
                            Abonnement actuel
                        </span>
                    </div>
                    <p class="text-3xl font-bold text-primary mt-2">
                        {{ formatPrice(plan.price_cents) }}
                        <span class="text-sm font-normal text-gray-400">/ {{ plan.interval_unit === 'month' ? 'mois' : 'an' }}</span>
                    </p>
                </div>

                <ul class="space-y-2 flex-1">
                    <li class="flex items-center gap-2 text-sm text-gray-600">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                        </svg>
                        <span v-if="isCurrentPlan(plan) && limits.announcements">{{ limits.announcements.used }}/{{ limits.announcements.max }} annonces ce mois-ci</span>
                        <span v-else-if="plan.max_announcements_per_month != null">{{ plan.max_announcements_per_month }} annonces / mois</span>
                        <span v-else>Annonces illimitées</span>
                    </li>
                    <li class="flex items-center gap-2 text-sm text-gray-600">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                        </svg>
                        <span v-if="isCurrentPlan(plan) && limits.projects">{{ limits.projects.used }}/{{ limits.projects.max }} projets ce mois-ci</span>
                        <span v-else-if="plan.max_projects_per_month != null">{{ plan.max_projects_per_month }} projets / mois</span>
                        <span v-else>Projets illimités</span>
                    </li>
                    <li class="flex items-center gap-2 text-sm text-gray-600">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                        </svg>
                        <span v-if="isCurrentPlan(plan) && limits.features">{{ limits.features.used }}/{{ limits.features.max }} mise(s) en avant ce mois-ci</span>
                        <span v-else-if="plan.max_features_per_month != null">{{ plan.max_features_per_month }} mise(s) en avant / mois</span>
                        <span v-else>Mises en avant illimitées</span>
                    </li>
                    <li class="flex items-center gap-2 text-sm text-gray-600">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                        </svg>
                        Achats d'annonces illimités
                    </li>
                </ul>

                <p v-if="isCurrentPlan(plan)" class="text-center text-sm font-medium text-secondary py-2.5">
                    Votre abonnement en cours
                </p>
                <button v-else-if="!subscription" @click="subscribe(plan)" :disabled="!plan.stripe_price_id || checkoutLoading === plan.id"
                    class="w-full py-2.5 text-sm font-semibold bg-primary text-white rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    <span v-if="checkoutLoading === plan.id">Redirection...</span>
                    <span v-else-if="!plan.stripe_price_id">Bientôt disponible</span>
                    <span v-else>S'abonner</span>
                </button>
                <button v-else @click="changePlan(plan)" :disabled="!plan.stripe_price_id || changeLoading === plan.id"
                    class="w-full py-2.5 text-sm font-semibold bg-primary text-white rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    <span v-if="changeLoading === plan.id">Changement en cours...</span>
                    <span v-else-if="!plan.stripe_price_id">Bientôt disponible</span>
                    <span v-else>Changer pour ce plan</span>
                </button>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6">
            <h2 class="font-semibold text-gray-800 mb-4">Pourquoi passer à un abonnement ?</h2>
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
    </ArtisanLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import api from '@/api.js'
import { usePolling } from '@/composables/usePolling'

const route = useRoute()

const plans = ref([])
const subscription = ref(null)
const limits = ref({})
const checkoutLoading = ref(null)
const changeLoading = ref(null)
const portalLoading = ref(false)
const polling = ref(false)
const pollTries = ref(0)

const successMsg = computed(() => route.query.success ? 'Votre abonnement a été activé avec succès !' : '')
const cancelledMsg = computed(() => route.query.cancelled ? 'Paiement annulé. Votre abonnement n\'a pas été modifié.' : '')

const featureDetails = [
    {
        title: 'Publiez davantage',
        desc: 'Chaque palier augmente votre quota mensuel d\'annonces et de projets par rapport à l\'offre gratuite.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>`,
    },
    {
        title: 'Achats d\'annonces illimités',
        desc: 'Achetez autant d\'annonces que vous souhaitez, sans restriction.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/></svg>`,
    },
    {
        title: 'Mise en avant selon votre palier',
        desc: 'Boostez la visibilité de vos annonces sur la page d\'accueil, avec un nombre de mises en avant mensuelles qui dépend de votre palier.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>`,
    },
    {
        title: 'Suivi avancé',
        desc: 'Statistiques détaillées sur vos ventes, vos projets et votre Upcycling Score.',
        icon: `<svg class="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
]

function isCurrentPlan(plan) {
    return subscription.value?.plan?.id === plan.id
}

function formatPrice(cents) {
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

function formatDate(dateStr) {
    return new Date(dateStr).toLocaleDateString('fr-FR')
}

async function fetchState() {
    const [plansRes, subRes, limitsRes] = await Promise.all([
        api.get('/subscription/plans').catch(() => ({ data: [] })),
        api.get('/user/subscription').catch(() => ({ data: {} })),
        api.get('/user/limits').catch(() => ({ data: {} })),
    ])
    plans.value = (plansRes.data ?? []).filter(p => p.is_active)
    subscription.value = subRes.data?.active ? subRes.data.subscription : null
    limits.value = limitsRes.data ?? {}
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

async function changePlan(plan) {
    changeLoading.value = plan.id
    try {
        const { data } = await api.post('/user/subscription/change', { plan_id: plan.id })
        sessionStorage.setItem('pendingPlanId', plan.id)
        window.location.href = data.url
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors du changement de plan.')
        changeLoading.value = null
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

function stillWaitingForStripe() {
    if (!route.query.success || pollTries.value >= 5) return false
    const pendingPlanId = sessionStorage.getItem('pendingPlanId')
    if (pendingPlanId) return subscription.value?.plan?.id !== pendingPlanId
    return !subscription.value
}

onMounted(fetchState)

usePolling(async () => {
    if (!stillWaitingForStripe()) {
        polling.value = false
        sessionStorage.removeItem('pendingPlanId')
        return
    }
    polling.value = true
    pollTries.value++
    await fetchState()
    if (!stillWaitingForStripe()) {
        sessionStorage.removeItem('pendingPlanId')
    }
}, 2000, stillWaitingForStripe)
</script>
