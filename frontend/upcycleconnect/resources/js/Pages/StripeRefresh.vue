<template>
    <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 p-6">
        <div class="text-center max-w-md">
            <h1 class="text-2xl font-bold text-gray-800 mb-2" style="font-family: var(--font-family-title)">
                Lien expiré
            </h1>
            <p class="text-sm text-gray-400 mb-8">
                Le lien d'inscription Stripe a expiré. Cliquez ci-dessous pour en générer un nouveau.
            </p>
            <button
                @click="refresh"
                :disabled="loading"
                class="px-5 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
            >
                {{ loading ? 'Chargement…' : 'Continuer l\'inscription' }}
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '@/api.js'

const loading = ref(false)

async function refresh() {
    loading.value = true
    try {
        const { data } = await api.post('/user/stripe/connect/onboarding')
        window.location.href = data.url
    } catch {
        loading.value = false
    }
}
</script>
