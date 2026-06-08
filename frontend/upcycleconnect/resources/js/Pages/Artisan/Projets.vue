<template>
    <ArtisanLayout>

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Mes projets</h1>
            <button
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un projet
            </button>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">Chargement…</div>

        <template v-else>
            <div class="flex items-center gap-3 mb-4">
                <select v-model="filterStatus" class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600 bg-white">
                    <option value="">Tous les statuts</option>
                    <option value="pending">En attente</option>
                    <option value="rejected">Rejeté</option>
                    <option value="open">Ouvert</option>
                    <option value="in_progress">En cours</option>
                    <option value="completed">Terminé</option>
                </select>
                <span class="text-xs text-gray-400">{{ filteredProjects.length }} projet(s)</span>
            </div>

            <div v-if="filteredProjects.length === 0" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">
                {{ filterStatus ? 'Aucun projet pour ce statut.' : 'Vous n\'avez créé aucun projet.' }}
            </div>

            <div v-else class="grid grid-cols-3 gap-5">
                <div
                    v-for="p in filteredProjects"
                    :key="p.id"
                    @click="selectedProject = p"
                    class="bg-white rounded-2xl shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow cursor-pointer border border-gray-100"
                >
                    <div class="flex items-start justify-between gap-2">
                        <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ p.title }}</h3>
                        <span :class="statusClass(p.status)" class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0">
                            {{ statusLabel(p.status) }}
                        </span>
                    </div>
                    <p v-if="p.description" class="text-xs text-gray-500 line-clamp-2">{{ p.description }}</p>
                    <p v-if="p.status === 'rejected' && p.rejection_reason" class="text-xs text-red-500 line-clamp-1">{{ p.rejection_reason }}</p>
                    <div class="flex items-center gap-3 mt-auto pt-2 border-t border-gray-50">
                        <div class="flex items-center gap-2 flex-1">
                            <span class="flex items-center gap-0.5 text-xs text-gray-400">
                                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"/>
                                </svg>
                                {{ p.up_votes ?? 0 }}
                            </span>
                            <span class="flex items-center gap-0.5 text-xs text-gray-400">
                                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.095c.5 0 .905-.405.905-.905 0-.714.211-1.412.608-2.006L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5"/>
                                </svg>
                                {{ p.down_votes ?? 0 }}
                            </span>
                        </div>
                        <span class="text-xs text-gray-400">{{ p.members_count }} membre(s)</span>
                        <button
                            @click.stop="toDelete = p"
                            class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                            title="Supprimer"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </template>

        <!-- Modale de détail -->
        <div v-if="selectedProject" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="selectedProject = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selectedProject.title }}</h3>
                    <button @click="selectedProject = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto flex-1">
                    <div>
                        <span :class="statusClass(selectedProject.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                            {{ statusLabel(selectedProject.status) }}
                        </span>
                    </div>
                    <p v-if="selectedProject.status === 'rejected' && selectedProject.rejection_reason" class="text-xs text-red-500">
                        Raison : {{ selectedProject.rejection_reason }}
                    </p>
                    <p v-if="selectedProject.description" class="leading-relaxed text-gray-600">{{ selectedProject.description }}</p>
                    <div class="grid grid-cols-2 gap-3 text-xs text-gray-500">
                        <div v-if="selectedProject.start_date">
                            <span class="font-medium text-gray-700">Début :</span>
                            {{ selectedProject.start_date.slice(0, 10) }}
                        </div>
                        <div v-if="selectedProject.end_date">
                            <span class="font-medium text-gray-700">Fin :</span>
                            {{ selectedProject.end_date.slice(0, 10) }}
                        </div>
                        <div v-if="selectedProject.location"><span class="font-medium text-gray-700">Lieu :</span> {{ selectedProject.location }}</div>
                        <div v-if="selectedProject.capacity"><span class="font-medium text-gray-700">Capacité :</span> {{ selectedProject.capacity }} membres</div>
                        <div><span class="font-medium text-gray-700">Membres :</span> {{ selectedProject.members_count }}</div>
                    </div>
                </div>
                <div v-if="selectedProject.status === 'pending' || selectedProject.status === 'rejected'" class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button
                        @click="openEditFromDetail(selectedProject)"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
                    >Modifier</button>
                </div>
            </div>
        </div>

        <!-- Confirmation de suppression -->
        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le projet ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimé.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteProject" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Modale de création / modification -->
        <div v-if="formModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal = false"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ editTarget ? 'Modifier le projet' : 'Créer un projet' }}
                    </h3>
                    <button @click="formModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4 max-h-[60vh] overflow-y-auto">
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input v-model="form.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre du projet"/>
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea v-model="form.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Décrivez votre projet…"/>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date de début</label>
                            <input v-model="form.start_date" type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date de fin</label>
                            <input v-model="form.end_date" type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Lieu</label>
                            <input v-model="form.location" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse ou lieu"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Capacité</label>
                            <input v-model.number="form.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée"/>
                        </div>
                    </div>
                    <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button @click="formModal = false" class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitForm" :disabled="submitting"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ submitting ? 'Enregistrement…' : (editTarget ? 'Mettre à jour' : 'Soumettre') }}
                    </button>
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

const route = useRoute()

const projects = ref([])
const loading = ref(true)
const toDelete = ref(null)
const deleting = ref(false)
const selectedProject = ref(null)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', start_date: '', end_date: '', location: '', capacity: null })
const filterStatus = ref('')
const filteredProjects = computed(() =>
    filterStatus.value ? projects.value.filter(p => p.status === filterStatus.value) : projects.value
)

function statusLabel(status) {
    return { open: 'Ouvert', in_progress: 'En cours', completed: 'Terminé', pending: 'En attente', rejected: 'Rejeté' }[status] ?? status
}

function statusClass(status) {
    return {
        open: 'bg-secondary/10 text-secondary',
        in_progress: 'bg-primary/10 text-primary',
        completed: 'bg-gray-100 text-gray-500',
        pending: 'bg-yellow-100 text-yellow-700',
        rejected: 'bg-red-100 text-red-600',
    }[status] ?? 'bg-gray-100 text-gray-500'
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', start_date: '', end_date: '', location: '', capacity: null }
    formError.value = ''
    formModal.value = true
}

function openEditFromDetail(p) {
    selectedProject.value = null
    editTarget.value = p
    form.value = {
        title: p.title,
        description: p.description ?? '',
        start_date: p.start_date ? p.start_date.slice(0, 10) : '',
        end_date: p.end_date ? p.end_date.slice(0, 10) : '',
        location: p.location ?? '',
        capacity: p.capacity ?? null,
    }
    formError.value = ''
    formModal.value = true
}

async function fetchProjects() {
    loading.value = true
    try {
        const { data } = await api.get('/user/my-projects')
        projects.value = data ?? []
    } catch {
        projects.value = []
    } finally {
        loading.value = false
    }
}

async function deleteProject() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/projects/${toDelete.value.id}`)
        projects.value = projects.value.filter(p => p.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

async function submitForm() {
    if (!form.value.title.trim()) { formError.value = 'Le titre est requis.'; return }
    submitting.value = true
    formError.value = ''
    try {
        const payload = {
            title: form.value.title,
            description: form.value.description || null,
            start_date: form.value.start_date || null,
            end_date: form.value.end_date || null,
            location: form.value.location || null,
            capacity: form.value.capacity || null,
        }
        if (editTarget.value) {
            await api.patch(`/projects/${editTarget.value.id}`, payload)
        } else {
            await api.post('/projects', payload)
        }
        formModal.value = false
        await fetchProjects()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await fetchProjects()
    const openId = route.query.open
    if (openId) {
        const found = projects.value.find(p => String(p.id) === String(openId))
        if (found) selectedProject.value = found
    }
})
</script>
