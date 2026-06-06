<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Catégories</h1>
        </div>

        <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 items-center">
            <span class="text-xs text-gray-400 ml-auto">{{ categories.length }} catégorie(s)</span>
            <button
                @click="openCreate"
                class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors whitespace-nowrap"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Ajouter une catégorie
            </button>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <div v-else class="bg-white rounded-xl shadow-sm overflow-hidden">
            <table class="w-full text-sm">
                <thead>
                    <tr class="bg-primary">
                        <th class="text-left text-white font-medium px-5 py-3">Nom</th>
                        <th class="text-left text-white font-medium px-5 py-3">Description</th>
                        <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr
                        v-for="(cat, i) in categories"
                        :key="cat.id"
                        :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']"
                    >
                        <td class="px-5 py-3 font-medium text-gray-800">{{ cat.name }}</td>
                        <td class="px-5 py-3 text-gray-500 text-xs">{{ cat.description || '-' }}</td>
                        <td class="px-5 py-3">
                            <div class="flex items-center justify-end gap-1">
                                <button
                                    @click="openEdit(cat)"
                                    class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors"
                                    title="Modifier"
                                >
                                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                                    </svg>
                                </button>
                                <button
                                    @click="confirmDelete(cat)"
                                    class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors"
                                    title="Supprimer"
                                >
                                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                    </svg>
                                </button>
                            </div>
                        </td>
                    </tr>
                    <tr v-if="categories.length === 0">
                        <td colspan="3" class="px-5 py-12 text-center text-gray-400 text-sm">
                            Aucune catégorie pour le moment.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Modal création / édition -->
        <div v-if="modal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="modal = false"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-4">
                    {{ editTarget ? 'Modifier la catégorie' : 'Nouvelle catégorie' }}
                </h3>
                <div class="space-y-3">
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Nom <span class="text-red-400">*</span></label>
                        <input
                            v-model="form.name"
                            type="text"
                            placeholder="Ex : Mobilier, Vêtements…"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea
                            v-model="form.description"
                            rows="3"
                            placeholder="Description optionnelle…"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                        />
                    </div>
                    <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
                </div>
                <div class="flex gap-3 mt-5">
                    <button @click="modal = false" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitForm" :disabled="submitting" class="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ submitting ? 'Enregistrement…' : (editTarget ? 'Mettre à jour' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal confirmation suppression -->
        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer la catégorie ?</h3>
                <p class="text-sm text-gray-500 mb-5">
                    « {{ toDelete.name }} » sera supprimée. Les annonces associées perdront leur catégorie.
                </p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="deleteCategory" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const categories = ref([])
const loading = ref(true)
const error = ref('')
const modal = ref(false)
const editTarget = ref(null)
const form = ref({ name: '', description: '' })
const formError = ref('')
const submitting = ref(false)
const toDelete = ref(null)
const deleting = ref(false)

async function fetchCategories() {
    try {
        const { data } = await api.get('/admin/categories')
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {
        error.value = 'Impossible de charger les catégories.'
    } finally {
        loading.value = false
    }
}

function openCreate() {
    editTarget.value = null
    form.value = { name: '', description: '' }
    formError.value = ''
    modal.value = true
}

function openEdit(cat) {
    editTarget.value = cat
    form.value = { name: cat.name, description: cat.description ?? '' }
    formError.value = ''
    modal.value = true
}

async function submitForm() {
    if (!form.value.name.trim()) {
        formError.value = 'Le nom est requis.'
        return
    }
    submitting.value = true
    formError.value = ''
    try {
        const payload = {
            name: form.value.name.trim(),
            description: form.value.description.trim() || null,
        }
        if (editTarget.value) {
            await api.put(`/admin/category/${editTarget.value.id}`, payload)
        } else {
            await api.post('/admin/categories', payload)
        }
        modal.value = false
        await fetchCategories()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

function confirmDelete(cat) {
    toDelete.value = cat
}

async function deleteCategory() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/category/${toDelete.value.id}`)
        toDelete.value = null
        await fetchCategories()
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

onMounted(fetchCategories)
</script>
