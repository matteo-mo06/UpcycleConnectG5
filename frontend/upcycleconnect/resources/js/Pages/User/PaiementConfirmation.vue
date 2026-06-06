<template>
    <UserLayout>
        <div class="max-w-lg mx-auto py-20 px-4 text-center">

            <div v-if="status === 'succeeded'" class="flex flex-col items-center gap-4">
                <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center">
                    <svg class="w-8 h-8 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                    </svg>
                </div>
                <h1 class="text-2xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Paiement réussi !</h1>
                <p class="text-gray-500">Votre achat a bien été enregistré. Vous retrouverez l'annonce dans vos acquisitions.</p>
                <router-link to="/annonces" class="mt-4 px-6 py-3 bg-primary text-white rounded-xl font-medium hover:bg-primary-dark transition-colors">
                    Retour aux annonces
                </router-link>
            </div>

            <div v-else-if="status === 'processing'" class="flex flex-col items-center gap-4">
                <div class="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center">
                    <svg class="w-8 h-8 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                </div>
                <h1 class="text-2xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Paiement en cours…</h1>
                <p class="text-gray-500">Votre paiement est en cours de traitement. Vous recevrez une confirmation sous peu.</p>
            </div>

            <div v-else-if="status !== 'loading'" class="flex flex-col items-center gap-4">
                <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
                    <svg class="w-8 h-8 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                </div>
                <h1 class="text-2xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Paiement échoué</h1>
                <p class="text-gray-500">{{ errorMessage ?? 'Une erreur est survenue lors du paiement.' }}</p>
                <button @click="router.back()" class="mt-4 px-6 py-3 bg-primary text-white rounded-xl font-medium hover:bg-primary-dark transition-colors">
                    Réessayer
                </button>
            </div>

        </div>
    </UserLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { loadStripe } from '@stripe/stripe-js'
import UserLayout from '@/Layouts/UserLayout.vue'

const router = useRouter()
const status = ref('loading')
const errorMessage = ref(null)

onMounted(async () => {
    const stripe = await loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)
    const clientSecret = new URLSearchParams(window.location.search).get('payment_intent_client_secret')

    if (!clientSecret) {
        status.value = 'error'
        return
    }

    const { paymentIntent, error } = await stripe.retrievePaymentIntent(clientSecret)

    if (error) {
        status.value = 'error'
        errorMessage.value = error.message
    } else {
        status.value = paymentIntent.status
    }
})
</script>
