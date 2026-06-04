<template>
    <AdminLayout title="Upcycling Score">

        <div class="grid grid-cols-3 gap-5 mb-8">
            <div v-for="stat in stats" :key="stat.label" class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                    <div :class="stat.iconClass" v-html="stat.icon" />
                </div>
                <div class="min-w-0">
                    <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                    <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                </div>
            </div>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <div v-if="!loading && !error">
            <div class="bg-white rounded-xl shadow-sm overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Description</th>
                                <th class="text-left text-white font-medium px-5 py-3">Identifiant</th>
                                <th class="text-left text-white font-medium px-5 py-3">Points</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(action, i) in actions"
                                :key="action.action_type"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']"
                            >
                                <td class="px-5 py-3 text-gray-700">{{ action.description }}</td>
                                <td class="px-5 py-3 font-mono text-xs text-gray-400">{{ action.action_type }}</td>
                                <td class="px-5 py-3">
                                    <span
                                        :class="action.points > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-600'"
                                        class="px-2 py-0.5 rounded-full text-xs font-semibold"
                                    >
                                        {{ action.points > 0 ? '+' : '' }}{{ action.points }} pts
                                    </span>
                                </td>
                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end">
                                        <button
                                            @click="openEdit(action)"
                                            title="Modifier les points"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors"
                                        >
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ actions.length }} action(s) configurée(s)
                </div>
            </div>
        </div>

        <div v-if="toEdit" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toEdit = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-1">Modifier les points</h3>
                <p class="text-sm text-gray-500 mb-4">{{ toEdit.description }}</p>
                <input
                    v-model.number="editPoints"
                    type="number"
                    class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:border-primary focus:ring-primary mb-4"
                    @keydown.enter="saveEdit"
                />
                <div class="flex gap-3">
                    <button
                        @click="toEdit = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Annuler
                    </button>
                    <button
                        @click="saveEdit"
                        :disabled="saving"
                        class="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-60"
                    >
                        {{ saving ? 'Enregistrement…' : 'Enregistrer' }}
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

const actions = ref([])
const loading = ref(true)

const stats = computed(() => [
    {
        label: 'Actions configurées',
        value: actions.value.length,
        bgClass: 'bg-blue-100',
        iconClass: 'text-blue-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>`,
    },
    {
        label: 'Récompenses',
        value: actions.value.filter(a => a.points > 0).length,
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/></svg>`,
    },
    {
        label: 'Pénalités',
        value: actions.value.filter(a => a.points < 0).length,
        bgClass: 'bg-red-100',
        iconClass: 'text-red-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"/></svg>`,
    },
])
const error = ref('')
const toEdit = ref(null)
const editPoints = ref(0)
const saving = ref(false)

async function fetchActions() {
    try {
        const { data } = await api.get('/admin/score-actions')
        actions.value = data ?? []
    } catch {
        error.value = 'Impossible de charger la configuration.'
    }
}

onMounted(async () => {
    await fetchActions()
    loading.value = false
})

function openEdit(action) {
    toEdit.value = action
    editPoints.value = action.points
}

async function saveEdit() {
    if (!toEdit.value) return
    saving.value = true
    try {
        await api.put(`/admin/score-action/${toEdit.value.action_type}`, { points: editPoints.value })
        toEdit.value = null
        await fetchActions()
    } catch {
        alert('Erreur lors de la mise à jour.')
    } finally {
        saving.value = false
    }
}
</script>
