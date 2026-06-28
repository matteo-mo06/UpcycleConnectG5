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
                v-for="project in projects"
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
            <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg flex flex-col max-h-[90vh]">
                <div class="flex items-center justify-between px-6 py-4 border-b border-gray-100 flex-shrink-0">
                    <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ selected.title }}
                    </h2>
                    <div class="flex items-center gap-1">
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
                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
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

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Matériaux</p>
                        <p v-if="materials.length === 0" class="text-xs text-gray-400">Aucun matériau renseigné.</p>
                        <table v-else class="w-full text-xs">
                            <thead>
                                <tr class="text-gray-400 border-b border-gray-100">
                                    <th class="text-left pb-1 font-medium">Nom</th>
                                    <th class="text-right pb-1 font-medium">Qté</th>
                                    <th class="text-left pb-1 font-medium pl-2">Unité</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="m in materials" :key="m.id" class="border-b border-gray-50 last:border-0">
                                    <td class="py-1.5 pr-2 text-gray-700">{{ m.name }}</td>
                                    <td class="py-1.5 text-right text-gray-600">{{ m.quantity }}</td>
                                    <td class="py-1.5 pl-2 text-gray-500">{{ m.unit ?? '-' }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Étapes</p>
                        <p v-if="steps.length === 0" class="text-xs text-gray-400">Aucune étape renseignée.</p>
                        <div v-else class="space-y-1">
                            <div v-for="s in steps" :key="s.id" class="flex items-center gap-2 py-1.5 px-2 rounded-lg">
                                <span
                                    class="flex-shrink-0 w-5 h-5 rounded-full border-2 flex items-center justify-center"
                                    :class="{ 'border-gray-300 text-gray-300': s.status === 'todo', 'border-primary text-primary': s.status === 'in_progress', 'border-secondary bg-secondary text-white': s.status === 'done' }"
                                >
                                    <svg v-if="s.status === 'done'" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg>
                                    <span v-else-if="s.status === 'in_progress'" class="w-1.5 h-1.5 rounded-full bg-primary"></span>
                                </span>
                                <div class="flex-1 min-w-0">
                                    <p class="text-sm text-gray-700 truncate" :class="{ 'line-through text-gray-400': s.status === 'done' }">{{ s.title }}</p>
                                    <p v-if="s.description" class="text-xs text-gray-500 truncate">{{ s.description }}</p>
                                </div>
                                <span class="px-1.5 py-0.5 rounded text-xs flex-shrink-0" :class="{ 'bg-gray-100 text-gray-500': s.status === 'todo', 'bg-primary/10 text-primary': s.status === 'in_progress', 'bg-secondary/10 text-secondary': s.status === 'done' }">{{ stepStatusLabel(s.status) }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3 flex-shrink-0">
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

        <ReportModal
            v-model="showReport"
            content-type="project"
            :content-id="reportTarget?.id"
            :content-title="reportTarget?.title"
        />
    </UserLayout>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
import UserLayout from '@/Layouts/UserLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import ReportModal from '@/Components/ReportModal.vue'
import api from '@/api.js'
import { useAuthStore } from '@/stores/auth.js'

const auth = useAuthStore()

const projects = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const search = ref('')
const filterStatus = ref('')
const selected = ref(null)
const showReport = ref(false)
const reportTarget = ref(null)
const materials = ref([])
const steps = ref([])

const canRegister = computed(() => auth.hasPermission('register_project'))

function canReportProject(p) {
    return p.id_creator !== auth.user?.id
}

const publicProjects = computed(() =>
    projects.value.filter(p => ['open', 'in_progress', 'completed'].includes(p.status))
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
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors du vote.')
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

async function fetchProjects(silent = false) {
    if (!silent) loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        const { data } = await api.get('/projects', { params })
        projects.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (err) {
        console.error('fetchProjects error:', err)
    } finally {
        loading.value = false
    }
}

function stepStatusLabel(s) {
    return { todo: 'À faire', in_progress: 'En cours', done: 'Fait' }[s] ?? s
}

watch(page, fetchProjects)

watch(selected, async (p) => {
    materials.value = []
    steps.value = []
    if (!p) return
    try {
        const [mRes, sRes] = await Promise.all([
            api.get(`/projects/${p.id}/materials`),
            api.get(`/projects/${p.id}/steps`),
        ])
        materials.value = mRes.data ?? []
        steps.value = sRes.data ?? []
    } catch (err) {
        console.error('project details error:', err)
    }
})

usePolling(fetchProjects)
</script>
