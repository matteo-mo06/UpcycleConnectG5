<template>
    <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="close"
    >
        <div class="absolute inset-0 bg-black/40" @click="close" />
        <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">

            <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                <h2 class="font-semibold text-gray-800 text-lg" style="font-family: var(--font-family-title)">Signaler ce contenu</h2>
                <button @click="close" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                </button>
            </div>

            <div class="px-6 py-5 space-y-4">
                <p v-if="contentTitle" class="text-sm text-gray-500">« {{ contentTitle }} »</p>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Raison du signalement</label>
                    <textarea
                        v-model="reason"
                        rows="4"
                        placeholder="Décrivez pourquoi vous signalez ce contenu…"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                    />
                </div>

                <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
                <p v-if="success" class="text-sm text-green-600">Signalement envoyé.</p>
            </div>

            <div class="px-6 py-4 border-t border-gray-100 flex gap-3">
                <button @click="close" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                <button
                    @click="submit"
                    :disabled="loading || success"
                    class="flex-1 py-2 bg-red-500 text-white text-sm font-medium rounded-lg hover:bg-red-600 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    {{ loading ? 'Envoi…' : 'Signaler' }}
                </button>
            </div>

        </div>
    </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import api from '@/api.js'

const props = defineProps({
    modelValue: Boolean,
    contentType: String,
    contentId: String,
    contentTitle: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

const reason = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

watch(() => props.modelValue, (val) => {
    if (val) {
        reason.value = ''
        error.value = ''
        success.value = false
    }
})

function close() {
    if (!loading.value) emit('update:modelValue', false)
}

async function submit() {
    error.value = ''
    if (!reason.value.trim()) { error.value = 'La raison est requise.'; return }

    loading.value = true
    try {
        const body = { reason: reason.value }
        body[`id_${props.contentType}`] = props.contentId
        await api.post('/reports', body)
        success.value = true
        setTimeout(close, 1500)
    } catch (e) {
        error.value = e.response?.data?.error ?? 'Une erreur est survenue.'
    } finally {
        loading.value = false
    }
}
</script>
