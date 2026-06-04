<template>
    <UserLayout>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1
                    class="text-3xl font-bold text-gray-800"
                    style="font-family: var(--font-family-title)"
                >
                    Projets
                </h1>
                <p class="text-sm text-gray-400 mt-1">
                    Rejoignez des projets collaboratifs d'upcycling
                </p>
            </div>
            <button
                v-if="canCreate"
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un projet
            </button>
        </div>

        <div v-if="canCreate && myCreated.length > 0" class="bg-white rounded-2xl shadow-sm p-6 mb-6">
            <h2 class="font-semibold text-gray-800 mb-4" style="font-family: var(--font-family-title)">
                Mes projets créés
            </h2>
            <div class="space-y-3">
                <div
                    v-for="p in myCreated"
                    :key="p.id"
                    class="flex items-center justify-between p-4 rounded-xl border border-gray-100 hover:border-gray-200 transition-colors"
                >
                    <div class="min-w-0 flex-1">
                        <div class="flex items-center gap-2 mb-0.5">
                            <span class="font-medium text-gray-800 truncate">{{ p.title }}</span>
                            <span :class="statusClass(p.status)" class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0">
                                {{ statusLabel(p.status) }}
                            </span>
                        </div>
                        <p v-if="p.status === 'rejected' && p.rejection_reason" class="text-xs text-red-500 mt-0.5">
                            Raison : {{ p.rejection_reason }}
                        </p>
                        <p class="text-xs text-gray-400">{{ p.members_count }} membre(s)</p>
                    </div>
                    <div class="flex items-center gap-2 flex-shrink-0 ml-4">
                        <button
                            v-if="p.status === 'pending' || p.status === 'rejected'"
                            @click="openEdit(p)"
                            class="p-1.5 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-colors"
                            title="Modifier"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                            </svg>
                        </button>
                        <button
                            @click="confirmDelete(p)"
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
        </div>

        <div class="flex gap-3 mb-6">
            <div class="flex-1 relative">
                <svg class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 1 0 5 11a6 6 0 0 0 12 0z"/>
                </svg>
                <input
                    v-model="search"
                    @input="debouncedSearch"
                    type="text"
                    placeholder="Rechercher un projet…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
            <select
                v-model="filterStatus"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Tous les statuts</option>
                <option value="open">Ouvert</option>
                <option value="in_progress">En cours</option>
                <option value="completed">Terminé</option>
            </select>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-sm text-gray-400">
            Chargement…
        </div>

        <div
            v-else-if="publicProjects.length === 0"
            class="bg-white rounded-2xl shadow-sm p-12 text-center"
        >
            <svg class="w-12 h-12 text-gray-200 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
            <p class="text-gray-400 text-sm">Aucun projet disponible pour le moment.</p>
        </div>

        <div v-else class="grid grid-cols-3 gap-5">
            <div
                v-for="project in filteredProjects"
                :key="project.id"
                class="bg-white rounded-2xl shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow cursor-pointer"
                @click="openDetail(project)"
            >
                <div class="flex items-start justify-between gap-2">
                    <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ project.title }}</h3>
                    <span :class="statusClass(project.status)" class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0">
                        {{ statusLabel(project.status) }}
                    </span>
                </div>
                <p class="text-xs text-gray-500 line-clamp-2">{{ project.description }}</p>

                <div class="flex items-center gap-3 mt-auto pt-2 border-t border-gray-50">
                    <div class="flex items-center gap-2 flex-1">
                        <button
                            @click.stop="vote(project, 1)"
                            :class="project.user_vote === 1 ? 'text-secondary font-semibold' : 'text-gray-400 hover:text-secondary'"
                            class="flex items-center gap-0.5 text-xs transition-colors"
                            title="J'aime"
                        >
                            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"/>
                            </svg>
                            {{ project.up_votes }}
                        </button>
                        <button
                            @click.stop="vote(project, -1)"
                            :class="project.user_vote === -1 ? 'text-red-500 font-semibold' : 'text-gray-400 hover:text-red-400'"
                            class="flex items-center gap-0.5 text-xs transition-colors"
                            title="Je n'aime pas"
                        >
                            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.095c.5 0 .905-.405.905-.905 0-.714.211-1.412.608-2.006L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5"/>
                            </svg>
                            {{ project.down_votes }}
                        </button>
                    </div>
                    <span class="text-xs text-gray-400">{{ project.members_count }} membre(s)</span>
                    <span v-if="project.is_registered" class="text-xs text-secondary font-medium">Membre</span>
                    <span v-else-if="isFull(project)" class="text-xs text-gray-400">Complet</span>
                </div>
            </div>
        </div>

        <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />

        <div
            v-if="selected"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="selected = null"
        >
            <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="flex items-center justify-between px-6 py-4 border-b border-gray-100">
                    <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ selected.title }}
                    </h2>
                    <div class="flex items-center gap-1">
                        <button
                            v-if="canDeleteProject(selected)"
                            @click="confirmDelete(selected); selected = null"
                            title="Supprimer le projet"
                            class="p-1.5 rounded-lg text-gray-300 hover:text-red-500 hover:bg-red-50 transition-colors"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </button>
                        <button
                            v-if="canReportProject(selected)"
                            @click="openReport(selected)"
                            title="Signaler le projet"
                            class="p-1.5 rounded-lg text-gray-300 hover:text-orange-500 hover:bg-orange-50 transition-colors"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"/>
                            </svg>
                        </button>
                        <button
                            @click="selected = null"
                            class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors"
                        >
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                            </svg>
                        </button>
                    </div>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div v-if="selected.start_date">
                            <p class="text-xs text-gray-400 mb-0.5">Début</p>
                            <p class="font-medium text-gray-800">{{ selected.start_date.slice(0, 10) }}</p>
                        </div>
                        <div v-if="selected.end_date">
                            <p class="text-xs text-gray-400 mb-0.5">Fin</p>
                            <p class="font-medium text-gray-800">{{ selected.end_date.slice(0, 10) }}</p>
                        </div>
                        <div v-if="selected.location">
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ selected.location }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Membres</p>
                            <p class="font-medium text-gray-800">
                                {{ selected.members_count }}
                                <span v-if="selected.capacity"> / {{ selected.capacity }}</span>
                            </p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Créé par</p>
                            <p class="font-medium text-gray-800">{{ selected.creator_name }}</p>
                        </div>
                    </div>
                    <div v-if="selected.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-600 leading-relaxed">{{ selected.description }}</p>
                    </div>
                    <div v-if="selected.is_registered" class="flex items-center gap-2 p-3 rounded-xl bg-secondary/10">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        <p class="text-sm text-secondary font-medium">Vous êtes membre de ce projet</p>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
                    <button
                        @click="selected = null"
                        class="px-4 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Fermer
                    </button>
                    <button
                        v-if="canRegister && selected.is_registered"
                        @click="toggleRegistration(selected); selected = null"
                        class="px-4 py-2 rounded-xl border border-red-200 text-red-600 text-sm font-medium hover:bg-red-50 transition-colors"
                    >
                        Quitter le projet
                    </button>
                    <button
                        v-else-if="canRegister && !selected.is_registered && !isFull(selected)"
                        @click="toggleRegistration(selected); selected = null"
                        class="px-4 py-2 rounded-xl bg-secondary text-white text-sm font-medium hover:bg-secondary-dark transition-colors"
                    >
                        Rejoindre le projet
                    </button>
                    <span v-else-if="isFull(selected)" class="px-4 py-2 text-sm text-gray-400">Complet</span>
                </div>
            </div>
        </div>

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

        <ReportModal
            v-model="showReport"
            content-type="project"
            :content-id="reportTarget?.id"
            :content-title="reportTarget?.title"
        />
    </UserLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import UserLayout from '@/Layouts/UserLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import ReportModal from '@/Components/ReportModal.vue'
import api from '@/api.js'
import { useAuthStore } from '@/stores/auth.js'

const auth = useAuthStore()

const projects = ref([])
const myCreated = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const search = ref('')
const filterStatus = ref('')
const selected = ref(null)
const toDelete = ref(null)
const deleting = ref(false)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', start_date: '', end_date: '', location: '', capacity: null })
const showReport = ref(false)
const reportTarget = ref(null)

const canCreate = computed(() => auth.hasPermission('create_project'))
const canRegister = computed(() => auth.hasPermission('register_project'))

function canDeleteProject(p) {
    return p.id_creator === auth.user?.id || auth.isAdmin || auth.hasPermission('moderate_projects')
}
function canReportProject(p) {
    return p.id_creator !== auth.user?.id
}

const publicProjects = computed(() =>
    projects.value.filter(p => ['open', 'in_progress', 'completed'].includes(p.status))
)

const filteredProjects = computed(() =>
    publicProjects.value.filter(p => !filterStatus.value || p.status === filterStatus.value)
)

let debounceTimer = null
function debouncedSearch() {
    clearTimeout(debounceTimer)
    page.value = 1
    debounceTimer = setTimeout(fetchProjects, 300)
}

function statusLabel(status) {
    return ({ open: 'Ouvert', in_progress: 'En cours', completed: 'Terminé', pending: 'En attente', rejected: 'Rejeté' }[status] ?? status)
}

function statusClass(status) {
    return ({
        open: 'bg-secondary/10 text-secondary',
        in_progress: 'bg-primary/10 text-primary',
        completed: 'bg-gray-100 text-gray-500',
        pending: 'bg-yellow-100 text-yellow-700',
        rejected: 'bg-red-100 text-red-600',
    }[status] ?? 'bg-gray-100 text-gray-500')
}

function isFull(p) {
    return p.capacity > 0 && p.members_count >= p.capacity
}

function openDetail(project) {
    selected.value = project
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', start_date: '', end_date: '', location: '', capacity: null }
    formError.value = ''
    formModal.value = true
}

function openEdit(p) {
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

function confirmDelete(p) {
    toDelete.value = p
}

function openReport(project) {
    reportTarget.value = project
    showReport.value = true
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function vote(project, value) {
    try {
        if (project.user_vote === value) {
            await api.delete(`/projects/${project.id}/vote`)
        } else {
            await api.put(`/projects/${project.id}/vote`, { value })
        }
        await fetchProjects()
    } catch {
    }
}

async function toggleRegistration(project) {
    if (!canRegister.value) return
    try {
        if (project.is_registered) {
            await api.delete(`/projects/${project.id}/join`)
            project.is_registered = false
            project.members_count = Math.max(0, project.members_count - 1)
        } else {
            await api.post(`/projects/${project.id}/join`)
            project.is_registered = true
            project.members_count++
        }
    } catch (err) {
        alert(err.response?.data?.error ?? 'Erreur lors de l\'inscription.')
    }
}

async function deleteProject() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/projects/${toDelete.value.id}`)
        myCreated.value = myCreated.value.filter(p => p.id !== toDelete.value.id)
        projects.value = projects.value.filter(p => p.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

async function submitForm() {
    if (!form.value.title.trim()) {
        formError.value = 'Le titre est requis.'
        return
    }
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
        await fetchMyCreated()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

async function fetchProjects() {
    loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        const { data } = await api.get('/projects', { params })
        projects.value = data.data ?? []
        total.value = data.total ?? 0
    } catch {
    } finally {
        loading.value = false
    }
}

async function fetchMyCreated() {
    if (!canCreate.value) return
    try {
        const { data } = await api.get('/user/my-projects')
        myCreated.value = data ?? []
    } catch {
    }
}

watch(page, fetchProjects)

onMounted(async () => {
    await Promise.all([fetchProjects(), fetchMyCreated()])
})
</script>
