<template>
    <ArtisanLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Mon Upcycling Score
            </h1>
            <p class="text-sm text-gray-400 mt-1">Suivez votre progression et découvrez comment gagner des points.</p>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-sm text-gray-400">
            Chargement…
        </div>

        <div v-else-if="premiumRequired" class="bg-white rounded-2xl shadow-sm p-12 flex flex-col items-center gap-4 text-center">
            <div class="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center">
                <svg class="w-8 h-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                </svg>
            </div>
            <div>
                <p class="font-semibold text-gray-800 mb-1">Upcycling Score — fonctionnalité Premium</p>
                <p class="text-sm text-gray-500 mb-4">Accédez au détail de votre score écologique avec l'abonnement Premium.</p>
                <RouterLink to="/artisan/abonnement"
                    class="inline-block px-5 py-2 bg-primary text-white text-sm font-semibold rounded-xl hover:bg-primary-dark transition-colors">
                    Découvrir Premium
                </RouterLink>
            </div>
        </div>

        <template v-else-if="!premiumRequired">

        <div class="bg-white rounded-2xl shadow-sm p-6 flex items-center gap-6 mb-6">
            <div class="flex-shrink-0 w-16 h-16 rounded-2xl bg-secondary/20 flex items-center justify-center">
                <svg class="w-8 h-8 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                </svg>
            </div>
            <div>
                <p class="text-sm text-gray-500">Score total</p>
                <p class="text-5xl font-bold text-gray-800 leading-none mt-1">{{ score !== null ? score : '–' }}</p>
            </div>
        </div>

        <div class="grid grid-cols-5 gap-6">

            <div class="col-span-3 bg-white rounded-2xl shadow-sm p-5">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Ce que j'ai gagné</h2>
                <p class="text-xs text-gray-400 mb-4">Actions complétées et points obtenus</p>

                <p v-if="earned.length === 0" class="text-sm text-gray-400 text-center py-6">
                    Vous n'avez pas encore gagné de points.
                </p>

                <div v-else class="space-y-3">
                    <div
                        v-for="item in earned"
                        :key="item.action_type"
                        class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0"
                    >
                        <div class="min-w-0">
                            <p class="text-sm text-gray-700">{{ item.description }}</p>
                            <p class="text-xs text-gray-400 mt-0.5">{{ item.count }} fois</p>
                        </div>
                        <span class="flex-shrink-0 ml-4 text-sm font-semibold text-secondary">+{{ item.subtotal }} pts</span>
                    </div>
                </div>

                <div v-if="sanctions.length > 0" class="mt-4 pt-4 border-t border-gray-100 space-y-3">
                    <div
                        v-for="item in sanctions"
                        :key="item.action_type"
                        class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0"
                    >
                        <div class="min-w-0">
                            <p class="text-sm text-gray-700">{{ item.description }}</p>
                            <p class="text-xs text-gray-400 mt-0.5">{{ item.count }} fois</p>
                        </div>
                        <span class="flex-shrink-0 ml-4 text-sm font-semibold text-red-500">{{ item.subtotal }} pts</span>
                    </div>
                </div>
            </div>

            <div class="col-span-2 bg-white rounded-2xl shadow-sm p-5">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Comment ça marche</h2>
                <p class="text-xs text-gray-400 mb-4">Points attribués par type d'action</p>

                <div class="space-y-1">
                    <div
                        v-for="item in allPositive"
                        :key="item.action_type"
                        class="flex items-center justify-between py-1.5"
                    >
                        <span class="text-sm text-gray-600">{{ item.description }}</span>
                        <span class="flex-shrink-0 ml-3 px-2 py-0.5 rounded-full text-xs font-semibold bg-green-50 text-green-700">+{{ item.points }} pts</span>
                    </div>
                </div>

                <div v-if="allNegative.length > 0" class="mt-3 pt-3 border-t border-gray-100 space-y-1">
                    <div
                        v-for="item in allNegative"
                        :key="item.action_type"
                        class="flex items-center justify-between py-1.5"
                    >
                        <span class="text-sm text-gray-600">{{ item.description }}</span>
                        <span class="flex-shrink-0 ml-3 px-2 py-0.5 rounded-full text-xs font-semibold bg-red-50 text-red-600">{{ item.points }} pts</span>
                    </div>
                </div>
            </div>

        </div>

        </template>

    </ArtisanLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import api from '@/api.js'

const loading = ref(false)
const premiumRequired = ref(false)
const score = ref(null)
const breakdown = ref([])

const earned = computed(() => breakdown.value.filter(a => a.count > 0 && a.points > 0))
const sanctions = computed(() => breakdown.value.filter(a => a.count > 0 && a.points < 0))
const allPositive = computed(() => breakdown.value.filter(a => a.points > 0))
const allNegative = computed(() => breakdown.value.filter(a => a.points < 0))

onMounted(async () => {
    loading.value = true
    try {
        const { data } = await api.get('/user/score-breakdown')
        score.value = data.score
        breakdown.value = data.breakdown
    } catch (e) {
        if (e.response?.status === 403) premiumRequired.value = true
    } finally {
        loading.value = false
    }
})
</script>
