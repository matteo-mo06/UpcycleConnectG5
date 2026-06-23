<template>
    <UserLayout>
        <div class="max-w-lg mx-auto py-10 px-4">

            <button @click="router.back()" class="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700 mb-6">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/>
                </svg>
                Retour
            </button>

            <div v-if="loading" class="text-center py-20 text-gray-400">Chargement…</div>
            <div v-else-if="error" class="text-center py-20 text-red-500">{{ error }}</div>

            <template v-else>
                <div class="bg-white rounded-2xl shadow-sm p-6 mb-6">
                    <h1 class="text-xl font-bold text-gray-800 mb-1" style="font-family: var(--font-family-title)">
                        Paiement sécurisé
                    </h1>
                    <p class="text-sm text-gray-500 mb-4">{{ formation?.title }}</p>
                    <div class="space-y-2 border-t pt-4">
                        <div class="flex items-center justify-between border-t pt-2">
                            <span class="font-semibold text-gray-800">Total</span>
                            <span class="text-xl font-bold text-primary">{{ formatCents(priceCents) }}</span>
                        </div>
                    </div>
                </div>

                <div class="bg-white rounded-2xl shadow-sm p-6">
                    <div id="payment-element" class="mb-6" />

                    <p v-if="paymentError" class="text-sm text-red-500 mb-4">{{ paymentError }}</p>

                    <button
                        @click="submitPayment"
                        :disabled="paying || !stripeReady"
                        class="w-full py-3 bg-primary text-white rounded-xl font-medium hover:bg-primary-dark transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {{ paying ? 'Traitement en cours…' : `Payer ${formatCents(priceCents)}` }}
                    </button>

                    <p class="text-xs text-gray-400 text-center mt-3">
                        Paiement sécurisé par Stripe. Vos données bancaires ne sont jamais transmises à nos serveurs.
                    </p>
                </div>
            </template>

        </div>
    </UserLayout>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { loadStripe } from '@stripe/stripe-js'
import UserLayout from '@/Layouts/UserLayout.vue'
import api from '@/api.js'

const router = useRouter()
const route = useRoute()

const formationId = route.params.id

const loading = ref(true)
const error = ref(null)
const paying = ref(false)
const paymentError = ref(null)
const stripeReady = ref(false)
const formation = ref(null)
const priceCents = ref(0)

function formatCents(cents) {
    if (!cents) return '0,00 €'
    return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })
}

let stripe = null
let elements = null

onMounted(async () => {
    try {
        const [formRes, intentRes] = await Promise.all([
            api.get(`/formations/${formationId}`),
            api.post(`/pay/formation/${formationId}`),
        ])

        formation.value = formRes.data
        priceCents.value = intentRes.data.price_cents

        loading.value = false
        await nextTick()

        stripe = await loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)
        elements = stripe.elements({ clientSecret: intentRes.data.client_secret })
        const paymentElement = elements.create('payment')
        paymentElement.mount('#payment-element')
        paymentElement.on('ready', () => { stripeReady.value = true })

    } catch (e) {
        error.value = e?.response?.data?.error ?? 'Impossible de charger le paiement.'
        loading.value = false
    }
})

onUnmounted(() => {
    elements?.getElement('payment')?.destroy()
})

async function submitPayment() {
    if (!stripe || !elements) return
    paying.value = true
    paymentError.value = null

    const { error: stripeError } = await stripe.confirmPayment({
        elements,
        confirmParams: {
            return_url: `${window.location.origin}/paiement-confirmation?type=formation`,
        },
    })

    if (stripeError) {
        paymentError.value = stripeError.message
        paying.value = false
    }
}
</script>
