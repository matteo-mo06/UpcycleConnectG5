<template>
    <AdminLayout>

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Conseils</h1>
            <button
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un conseil
            </button>
        </div>

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                <select
                    v-model="filterCategory"
                    class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
                >
                    <option value="">Toutes les catégories</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap ml-auto">{{ filtered.length }} conseil(s)</span>
            </div>

            <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">Chargement…</div>

            <div v-else-if="filtered.length === 0" class="py-12 text-center text-gray-400 text-sm">
                Aucun conseil publié pour le moment.
            </div>

            <div v-else class="grid grid-cols-3 gap-5 p-6">
                <div
                    v-for="c in filtered"
                :key="c.id"
                class="bg-white rounded-2xl overflow-hidden border border-gray-100 flex flex-col hover:shadow-md transition-shadow cursor-pointer"
                @click="openDetail(c)"
            >
                <div class="h-32 bg-secondary/10 flex items-center justify-center overflow-hidden">
                    <img v-if="covers[c.id]" :src="covers[c.id]" class="w-full h-full object-cover" />
                    <svg v-else class="w-10 h-10 text-secondary/40" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                    </svg>
                </div>
                <div class="p-5 flex flex-col gap-2 flex-1">
                    <div class="flex items-center justify-between gap-2">
                        <div class="flex items-center gap-1.5 flex-wrap min-w-0">
                            <span class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0" :class="priorityBadge(c.priority).class">{{ priorityBadge(c.priority).label }}</span>
                            <span v-if="c.is_expired" class="px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500 flex-shrink-0">Expiré</span>
                            <span v-if="c.category_name" class="px-2 py-0.5 rounded-full text-xs font-medium bg-primary/10 text-primary truncate">{{ c.category_name }}</span>
                        </div>
                        <span class="text-xs text-gray-400 flex-shrink-0">{{ c.created_at?.slice(0, 10) ?? '-' }}</span>
                    </div>
                    <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ c.title }}</h3>
                    <p v-if="c.description" class="text-xs text-gray-500 line-clamp-3 flex-1">{{ c.description }}</p>
                    <div class="mt-auto pt-3 flex items-center gap-2 justify-between border-t border-gray-50">
                        <span class="text-xs text-gray-400 truncate">{{ c.creator_name ?? 'Inconnu' }}</span>
                        <div class="flex items-center gap-2 flex-shrink-0">
                            <button
                                @click.stop="openEdit(c)"
                                class="p-1.5 text-primary hover:bg-primary/10 rounded-lg transition-colors"
                                title="Modifier"
                            >
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                                </svg>
                            </button>
                            <button
                                @click.stop="toDelete = c"
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
            </div>
        </div>

        <div class="mt-10">
            <h2 class="text-xl font-bold text-gray-800 mb-1" style="font-family: var(--font-family-title)">
                Conseils supprimés
            </h2>
            <p class="text-sm text-gray-400 mb-4">Historique des conseils supprimés</p>

            <div v-if="deletedConseils.length === 0" class="bg-white rounded-2xl shadow-sm p-8 text-center">
                <p class="text-gray-400 text-sm">Aucun conseil supprimé.</p>
            </div>

            <div v-else class="space-y-3">
                <div
                    v-for="c in deletedConseils"
                    :key="c.id"
                    @click="selectedDeleted = c"
                    class="bg-white rounded-2xl shadow-sm px-5 py-4 flex items-center gap-4 cursor-pointer hover:shadow-md transition-shadow"
                >
                    <div class="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center flex-shrink-0">
                        <svg class="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </div>
                    <div class="flex-1 min-w-0">
                        <h3 class="font-semibold text-gray-500 text-sm truncate">{{ c.title }}</h3>
                        <p class="text-xs text-gray-400 mt-0.5">
                            Par {{ c.creator_name ?? 'Inconnu' }} · Supprimé le {{ c.deleted_at?.slice(0, 10) }}
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="selected" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="selected = null">
            <div class="bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[85vh] flex flex-col overflow-hidden">
                <div class="flex items-center justify-between px-6 py-4 border-b border-gray-100 flex-shrink-0">
                    <div class="flex-1 min-w-0">
                        <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selected.title }}</h2>
                        <p v-if="selected.category_name" class="text-xs text-primary mt-0.5">{{ selected.category_name }}</p>
                    </div>
                    <button @click="selected = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors ml-3 flex-shrink-0">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <p v-if="selected.description" class="text-sm text-gray-600 whitespace-pre-line break-words">{{ selected.description }}</p>
                    <p v-else class="text-sm text-gray-400 italic">Aucun contenu.</p>
                    <div v-if="docsLoading" class="text-xs text-gray-400">Chargement des photos…</div>
                    <div v-else-if="selectedDocs.length" class="flex gap-2 flex-wrap">
                        <a v-for="doc in selectedDocs" :key="doc.id" :href="doc.link" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border border-gray-100">
                            <img :src="doc.link" class="w-full h-full object-cover" />
                        </a>
                    </div>
                    <div class="flex items-center gap-3 text-xs text-gray-400 pt-2 border-t border-gray-50">
                        <span v-if="selected.is_expired" class="px-2 py-0.5 rounded-full bg-gray-100 text-gray-500 font-medium">Expiré</span>
                        <span v-else-if="selected.date_advice">Valable jusqu'au {{ selected.date_advice.slice(0, 10) }}</span>
                        <span class="ml-auto">{{ selected.creator_name }}</span>
                        <span>{{ selected.created_at?.slice(0, 10) }}</span>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button @click="openEdit(selected); selected = null" class="px-4 py-1.5 text-sm text-primary border border-primary/30 rounded-lg hover:bg-primary/5 transition-colors">Modifier</button>
                    <button @click="selected = null" class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">Fermer</button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le conseil ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera déplacé dans les conseils supprimés.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteConseil" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="formModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal = false"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ editTarget ? 'Modifier le conseil' : 'Nouveau conseil' }}
                    </h3>
                    <button @click="formModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input v-model="form.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre du conseil"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Catégorie</label>
                        <select v-model="form.id_category"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                            <option value="">Aucune catégorie</option>
                            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Priorité</label>
                        <select v-model.number="form.priority"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                            <option :value="1">Haute</option>
                            <option :value="2">Moyenne</option>
                            <option :value="3">Basse</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Valable jusqu'au (optionnel)</label>
                        <input v-model="form.date_advice" type="date"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea v-model="form.description" rows="5"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Contenu du conseil…"/>
                    </div>
                    <div v-if="!editTarget">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Photos <span class="text-gray-300 font-normal">(3 max)</span></label>
                        <div
                            class="border-2 border-dashed border-gray-200 rounded-xl p-5 text-center cursor-pointer hover:border-primary/50 transition-colors"
                            @dragover.prevent
                            @drop.prevent="handleDrop"
                            @click="$refs.fileInput.click()"
                        >
                            <svg class="w-7 h-7 text-gray-300 mx-auto mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"/>
                            </svg>
                            <p class="text-xs text-gray-400">Glissez vos photos ici</p>
                            <span class="mt-1 inline-block text-xs text-primary underline">Parcourir</span>
                            <input ref="fileInput" type="file" accept="image/*" multiple class="hidden" @change="handleFileSelect" />
                        </div>
                        <div v-if="photos.length" class="flex gap-2 mt-3 flex-wrap">
                            <div v-for="(photo, i) in photos" :key="i" class="relative w-14 h-14 rounded-lg overflow-hidden border border-gray-200">
                                <img :src="photo.preview" class="w-full h-full object-cover" />
                                <button
                                    @click="removePhoto(i)"
                                    class="absolute top-0.5 right-0.5 w-4 h-4 bg-black/60 rounded-full flex items-center justify-center text-white"
                                >
                                    <svg class="w-2.5 h-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                                    </svg>
                                </button>
                            </div>
                        </div>
                    </div>
                    <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex justify-end gap-2">
                    <button @click="formModal = false" class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitForm" :disabled="submitting"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ submitting ? 'Enregistrement…' : (editTarget ? 'Enregistrer' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="selectedDeleted" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="selectedDeleted = null">
            <div class="bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[85vh] flex flex-col overflow-hidden">
                <div class="flex items-start gap-3 px-6 py-4 border-b border-gray-100 flex-shrink-0">
                    <div class="flex-1 min-w-0">
                        <h2 class="font-semibold text-gray-600 text-sm truncate" style="font-family: var(--font-family-title)">{{ selectedDeleted.title }}</h2>
                        <p class="text-xs text-gray-400 mt-0.5">Par {{ selectedDeleted.creator_name ?? 'Inconnu' }} · Supprimé le {{ selectedDeleted.deleted_at?.slice(0, 10) }}</p>
                    </div>
                    <span class="flex-shrink-0 px-2 py-1 rounded-full bg-gray-100 text-gray-400 text-xs font-medium">Supprimé</span>
                    <button @click="selectedDeleted = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors flex-shrink-0">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div class="flex items-center gap-1.5 flex-wrap">
                        <span class="px-2 py-0.5 rounded-full text-xs font-medium" :class="priorityBadge(selectedDeleted.priority).class">{{ priorityBadge(selectedDeleted.priority).label }}</span>
                        <span v-if="selectedDeleted.category_name" class="px-2 py-0.5 rounded-full text-xs font-medium bg-primary/10 text-primary">{{ selectedDeleted.category_name }}</span>
                        <span v-if="selectedDeleted.date_advice" class="text-xs text-gray-400">Valable jusqu'au {{ selectedDeleted.date_advice.slice(0, 10) }}</span>
                    </div>
                    <p v-if="selectedDeleted.description" class="text-sm text-gray-600 whitespace-pre-line break-words">{{ selectedDeleted.description }}</p>
                    <p v-else class="text-sm text-gray-400 italic">Aucun contenu.</p>
                    <div v-if="selectedDeleted.documents?.length" class="flex gap-2 flex-wrap">
                        <a v-for="doc in selectedDeleted.documents" :key="doc.id" :href="doc.link" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border border-gray-100">
                            <img :src="doc.link" class="w-full h-full object-cover" />
                        </a>
                    </div>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const conseils = ref([])
const deletedConseils = ref([])
const categories = ref([])
const loading = ref(true)
const toDelete = ref(null)
const deleting = ref(false)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const filterCategory = ref('')
const photos = ref([])
const form = ref({ title: '', description: '', id_category: '', priority: 2, date_advice: '' })
const selected = ref(null)
const selectedDeleted = ref(null)
const selectedDocs = ref([])
const docsLoading = ref(false)
const covers = ref({})

const filtered = computed(() => {
    if (!filterCategory.value) return conseils.value
    return conseils.value.filter(c => c.id_category === filterCategory.value)
})

const priorityBadges = {
    1: { label: 'Haute', class: 'bg-primary/10 text-primary' },
    2: { label: 'Moyenne', class: 'bg-amber-100 text-amber-700' },
    3: { label: 'Basse', class: 'bg-gray-100 text-gray-500' },
}

function priorityBadge(p) {
    return priorityBadges[p] ?? priorityBadges[2]
}

async function loadCovers() {
    const results = await Promise.allSettled(
        conseils.value.map(c =>
            api.get(`/conseils/${c.id}/documents`).then(({ data }) => ({ id: c.id, docs: data ?? [] }))
        )
    )
    const map = {}
    for (const r of results) {
        if (r.status === 'fulfilled' && r.value.docs.length) {
            map[r.value.id] = r.value.docs[0].link
        }
    }
    covers.value = map
}

async function openDetail(c) {
    selected.value = c
    selectedDocs.value = []
    docsLoading.value = true
    try {
        const { data } = await api.get(`/conseils/${c.id}/documents`)
        selectedDocs.value = data ?? []
    } catch {
        selectedDocs.value = []
    } finally {
        docsLoading.value = false
    }
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', id_category: '', priority: 2, date_advice: '' }
    photos.value = []
    formError.value = ''
    formModal.value = true
}

function openEdit(c) {
    editTarget.value = c
    form.value = {
        title: c.title,
        description: c.description ?? '',
        id_category: c.id_category ?? '',
        priority: c.priority ?? 2,
        date_advice: c.date_advice ? c.date_advice.slice(0, 10) : '',
    }
    formError.value = ''
    formModal.value = true
}

function handleFileSelect(e) {
    addFiles(Array.from(e.target.files))
    e.target.value = ''
}

function handleDrop(e) {
    addFiles(Array.from(e.dataTransfer.files).filter(f => f.type.startsWith('image/')))
}

function addFiles(files) {
    for (const file of files) {
        if (photos.value.length >= 3) break
        photos.value.push({ file, preview: URL.createObjectURL(file) })
    }
}

function removePhoto(i) {
    URL.revokeObjectURL(photos.value[i].preview)
    photos.value.splice(i, 1)
}

async function uploadPhotos() {
    const urls = []
    for (const photo of photos.value) {
        const fd = new FormData()
        fd.append('file', photo.file)
        const { data } = await api.post('/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
        urls.push(data.url)
    }
    return urls
}

async function fetchConseils() {
    loading.value = true
    try {
        const { data } = await api.get('/conseils')
        conseils.value = data?.data ?? data ?? []
    } catch {
        conseils.value = []
    } finally {
        loading.value = false
    }
}

async function fetchDeletedConseils() {
    try {
        const { data } = await api.get('/admin/advice/deleted')
        deletedConseils.value = data ?? []
    } catch {
        deletedConseils.value = []
    }
}

async function deleteConseil() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/conseils/${toDelete.value.id}`)
        conseils.value = conseils.value.filter(c => c.id !== toDelete.value.id)
        toDelete.value = null
        await Promise.all([loadCovers(), fetchDeletedConseils()])
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
        if (editTarget.value) {
            await api.patch(`/conseils/${editTarget.value.id}`, {
                title: form.value.title,
                description: form.value.description || null,
                id_category: form.value.id_category || null,
                priority: form.value.priority,
                date_advice: form.value.date_advice || null,
            })
        } else {
            const photoURLs = photos.value.length ? await uploadPhotos() : []
            await api.post('/conseils', {
                title: form.value.title,
                description: form.value.description || null,
                id_category: form.value.id_category || null,
                priority: form.value.priority,
                date_advice: form.value.date_advice || null,
                photo_urls: photoURLs,
            })
        }
        formModal.value = false
        photos.value = []
        await fetchConseils()
        await loadCovers()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await fetchConseils()
    await Promise.allSettled([
        loadCovers(),
        fetchDeletedConseils(),
        api.get('/categories').then(({ data }) => {
            categories.value = Array.isArray(data) ? data : (data.data ?? [])
        }),
    ])
})
</script>
