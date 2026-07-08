<template>
    <ArtisanLayout>
        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Mes publicités</h1>
            <button
                v-if="isPremium && !showForm"
                @click="showForm = true"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Nouvelle publicité
            </button>
        </div>

        <div v-if="successMessage" class="mb-4 px-4 py-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700">
            {{ successMessage }}
        </div>
        <div v-if="fetchError" class="mb-4 px-4 py-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-600">
            {{ fetchError }}
        </div>

        <div v-if="!isPremium && !loading" class="bg-white rounded-2xl shadow-sm p-10 flex flex-col items-center justify-center text-center gap-4">
            <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
                <svg class="w-8 h-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                </svg>
            </div>
            <div>
                <p class="font-semibold text-gray-700">Abonnement premium requis</p>
                <p class="text-sm text-gray-400 mt-1">Souscrivez à un abonnement premium pour publier des publicités sur l'accueil.</p>
            </div>
            <RouterLink to="/artisan/abonnement" class="px-4 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors">
                Voir les abonnements
            </RouterLink>
        </div>

        <div v-else-if="isPremium">

            <div v-if="showForm" class="bg-white rounded-2xl shadow-sm p-6 mb-6">
                <h2 class="font-semibold text-gray-800 mb-4">Nouvelle publicité</h2>
                <div class="space-y-4 max-w-lg">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input
                            v-model="form.title"
                            type="text"
                            placeholder="Nom de votre publicité"
                            class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Image <span class="text-red-400">*</span></label>
                        <div v-if="form.image_url" class="mb-2">
                            <img :src="form.image_url" class="h-32 rounded-lg object-cover" />
                        </div>
                        <label class="flex items-center gap-2 px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-500 cursor-pointer hover:border-primary hover:text-primary transition-colors">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                            </svg>
                            <span>{{ uploadingImage ? 'Envoi en cours…' : (form.image_url ? 'Changer l\'image' : 'Choisir une image') }}</span>
                            <input type="file" accept="image/*" class="hidden" @change="uploadImage" :disabled="uploadingImage" />
                        </label>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Lien de redirection <span class="text-gray-400 font-normal">(optionnel)</span></label>
                        <input
                            v-model="form.link_url"
                            type="url"
                            placeholder="https://votre-site.com"
                            class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary"
                        />
                    </div>
                    <p v-if="formError" class="text-sm text-red-500">{{ formError }}</p>
                    <div class="flex gap-3 pt-1">
                        <button
                            @click="submitAd"
                            :disabled="submitting || uploadingImage"
                            class="px-4 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-60"
                        >
                            {{ submitting ? 'Envoi…' : 'Soumettre pour validation' }}
                        </button>
                        <button @click="cancelForm" class="px-4 py-2 border border-gray-200 text-gray-600 rounded-lg text-sm hover:bg-gray-50 transition-colors">
                            Annuler
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
            <div v-else-if="ads.length === 0 && !showForm" class="bg-white rounded-2xl shadow-sm p-10 text-center text-sm text-gray-400">
                Vous n'avez pas encore de publicité. Cliquez sur "Nouvelle publicité" pour commencer.
            </div>
            <div v-else class="space-y-4">
                <div
                    v-for="ad in ads"
                    :key="ad.id"
                    class="bg-white rounded-2xl shadow-sm p-5 flex gap-5 items-start"
                >
                    <img :src="ad.image_url" :alt="ad.title" class="w-32 h-20 object-cover rounded-xl flex-shrink-0" />
                    <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                            <p class="font-semibold text-gray-800 truncate">{{ ad.title }}</p>
                            <span :class="stateBadge(ad.state).class" class="px-2 py-0.5 rounded-full text-xs font-semibold flex-shrink-0">
                                {{ stateBadge(ad.state).label }}
                            </span>
                        </div>
                        <p v-if="ad.link_url" class="text-xs text-gray-400 truncate mb-1">{{ ad.link_url }}</p>
                        <p class="text-xs text-gray-400">Soumis le {{ formatDate(ad.created_at) }}</p>
                        <p v-if="ad.expires_at" class="text-xs mt-0.5" :class="ad.state === 'expired' ? 'text-red-400' : 'text-gray-400'">
                            {{ ad.state === 'expired' ? 'Expirée le' : 'Expire le' }} {{ formatDate(ad.expires_at) }}
                        </p>
                        <p v-if="ad.rejection_reason" class="mt-1 text-xs text-red-500">Raison : {{ ad.rejection_reason }}</p>
                        <div v-if="ad.state === 'approved'" class="mt-3">
                            <p class="text-xs text-gray-500 mb-2">Votre publicité a été approuvée. Choisissez une durée de campagne :</p>
                            <div class="flex flex-wrap gap-2 mb-3">
                                <button
                                    v-for="plan in plans"
                                    :key="plan.id"
                                    @click="selectedPlan[ad.id] = plan.id"
                                    :class="selectedPlan[ad.id] === plan.id
                                        ? 'border-primary bg-primary/5 text-primary'
                                        : 'border-gray-200 text-gray-600 hover:border-primary/50'"
                                    class="flex flex-col items-center px-4 py-2 border rounded-lg text-sm transition-colors cursor-pointer"
                                >
                                    <span class="font-semibold">{{ plan.weeks }} semaines</span>
                                    <span class="text-xs">{{ formatEuros(plan.price_cents) }}</span>
                                </button>
                            </div>
                            <button
                                @click="pay(ad.id)"
                                :disabled="paying === ad.id || !selectedPlan[ad.id]"
                                class="px-4 py-2 bg-secondary text-white rounded-lg text-sm font-medium hover:bg-secondary-dark transition-colors disabled:opacity-60"
                            >
                                {{ paying === ad.id ? 'Redirection…' : selectedPlan[ad.id] ? 'Payer ' + formatEuros(planPrice(selectedPlan[ad.id])) : 'Sélectionnez une durée' }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </ArtisanLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import api from '@/api.js'

const route = useRoute()

const isPremium = ref(false)
const loading = ref(true)
const ads = ref([])
const plans = ref([])
const selectedPlan = ref({})
const showForm = ref(false)
const submitting = ref(false)
const uploadingImage = ref(false)
const formError = ref('')
const paying = ref(null)
const successMessage = ref('')
const fetchError = ref('')

const form = ref({ title: '', image_url: '', link_url: '' })

function formatEuros(cents) {
    if (!cents) return '0,00 €'
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

function formatDate(str) {
    if (!str) return '-'
    return new Date(str).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function stateBadge(state) {
    const map = {
        pending:  { label: 'En attente',  class: 'bg-yellow-50 text-yellow-700' },
        approved: { label: 'Approuvée',   class: 'bg-blue-50 text-blue-700' },
        rejected: { label: 'Refusée',     class: 'bg-red-50 text-red-600' },
        active:   { label: 'Active',      class: 'bg-green-50 text-green-700' },
        expired:  { label: 'Expirée',     class: 'bg-gray-100 text-gray-500' },
    }
    return map[state] ?? { label: state, class: 'bg-gray-100 text-gray-500' }
}

async function uploadImage(e) {
    const file = e.target.files?.[0]
    if (!file) return
    uploadingImage.value = true
    try {
        const fd = new FormData()
        fd.append('file', file)
        const { data } = await api.post('/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
        form.value.image_url = data.url
    } catch {
        formError.value = 'Erreur lors de l\'envoi de l\'image.'
    } finally {
        uploadingImage.value = false
    }
}

async function submitAd() {
    formError.value = ''
    if (!form.value.title.trim()) { formError.value = 'Le titre est requis.'; return }
    if (!form.value.image_url) { formError.value = 'Veuillez téléverser une image.'; return }
    submitting.value = true
    try {
        await api.post('/advertisement', {
            title: form.value.title.trim(),
            image_url: form.value.image_url,
            link_url: form.value.link_url || null,
        })
        cancelForm()
        await fetchAds()
    } catch (err) {
        formError.value = err?.response?.data?.error ?? 'Erreur lors de la soumission.'
    } finally {
        submitting.value = false
    }
}

function cancelForm() {
    showForm.value = false
    form.value = { title: '', image_url: '', link_url: '' }
    formError.value = ''
}

function planPrice(planId) {
    return plans.value.find(p => p.id === planId)?.price_cents ?? 0
}

async function pay(id) {
    const planId = selectedPlan.value[id]
    if (!planId) return
    paying.value = id
    try {
        const { data } = await api.post(`/advertisement/${id}/checkout`, { plan_id: planId })
        window.location.href = data.url
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur paiement.')
        paying.value = null
    }
}

async function fetchAds() {
    try {
        const { data } = await api.get('/user/advertisements')
        ads.value = data ?? []
        fetchError.value = ''
    } catch (err) {
        fetchError.value = err?.response?.data?.error ?? 'Erreur lors du chargement de vos publicités.'
    }
}

onMounted(async () => {
    if (route.query.success === '1') {
        successMessage.value = 'Paiement confirmé ! Votre publicité est maintenant en ligne.'
    }
    try {
        const { data } = await api.get('/user/subscription')
        isPremium.value = data.active === true
    } catch (err) {
        console.error('fetchSubscription error:', err)
    }
    try {
        const { data } = await api.get('/advertisement/plans')
        plans.value = data ?? []
    } catch (err) {
        console.error('fetchPlans error:', err)
    }
    await fetchAds()
    loading.value = false
})
</script>
