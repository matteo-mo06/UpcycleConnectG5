<template>
    <SalarieLayout>

        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Formations</h1>
                <p class="text-sm text-gray-400 mt-1">Gérez vos formations et suivez leur statut</p>
            </div>
            <button
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer une formation
            </button>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">Chargement…</div>

        <template v-else>
            <div class="flex items-center gap-3 mb-4">
                <select v-model="filterStatus" class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600 bg-white">
                    <option value="">Tous les statuts</option>
                    <option value="pending">En attente</option>
                    <option value="approved">Approuvée</option>
                    <option value="rejected">Rejetée</option>
                </select>
                <span class="text-xs text-gray-400">{{ filteredFormations.length }} formation(s)</span>
            </div>

            <div v-if="filteredFormations.length === 0" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">
                {{ filterStatus ? 'Aucune formation pour ce statut.' : 'Vous n\'avez créé aucune formation.' }}
            </div>

            <div v-else class="grid grid-cols-3 gap-5">
                <div
                    v-for="f in filteredFormations"
                    :key="f.id"
                    @click="selectedFormation = f"
                    class="bg-white rounded-2xl overflow-hidden hover:shadow-md transition-shadow cursor-pointer flex flex-col border border-gray-100"
                >
                    <div class="h-32 bg-secondary/10 flex items-center justify-center">
                        <svg class="w-10 h-10 text-secondary/40" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                        </svg>
                    </div>
                    <div class="p-4 flex flex-col gap-2 flex-1">
                        <div class="flex items-center gap-2">
                            <span :class="levelClass(f.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">{{ levelLabel(f.level) }}</span>
                            <span v-if="f.duration_hours" class="text-xs text-gray-400">{{ f.duration_hours }}h</span>
                        </div>
                        <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ f.title }}</h3>
                        <p v-if="f.description" class="text-xs text-gray-500 line-clamp-2">{{ f.description }}</p>
                        <p v-if="f.status === 'rejected' && f.rejection_reason" class="text-xs text-red-500 line-clamp-1">{{ f.rejection_reason }}</p>
                        <div class="mt-auto pt-2 flex items-center justify-between">
                            <span :class="statusClass(f.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">{{ statusLabel(f.status) }}</span>
                            <span v-if="f.price" class="text-xs font-semibold text-primary">{{ f.price.toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' }) }}</span>
                            <button
                                @click.stop="toDelete = f"
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
        </template>

        <div v-if="selectedFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="selectedFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selectedFormation.title }}</h3>
                    <button @click="selectedFormation = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto flex-1">
                    <div class="flex gap-2 flex-wrap">
                        <span :class="statusClass(selectedFormation.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                            {{ statusLabel(selectedFormation.status) }}
                        </span>
                        <span class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">{{ levelLabel(selectedFormation.level) }}</span>
                    </div>
                    <p v-if="selectedFormation.status === 'rejected' && selectedFormation.rejection_reason" class="text-xs text-red-500">
                        Raison : {{ selectedFormation.rejection_reason }}
                    </p>
                    <p v-if="selectedFormation.description" class="leading-relaxed text-gray-600">{{ selectedFormation.description }}</p>
                    <div class="grid grid-cols-2 gap-3 text-xs text-gray-500">
                        <div v-if="selectedFormation.date"><span class="font-medium text-gray-700">Date :</span> {{ formatDateTime(selectedFormation.date) }}</div>
                        <div v-if="selectedFormation.location"><span class="font-medium text-gray-700">Lieu :</span> {{ selectedFormation.location }}</div>
                        <div v-if="selectedFormation.duration_hours"><span class="font-medium text-gray-700">Durée :</span> {{ selectedFormation.duration_hours }}h</div>
                        <div v-if="selectedFormation.capacity"><span class="font-medium text-gray-700">Capacité :</span> {{ selectedFormation.capacity }} places</div>
                        <div><span class="font-medium text-gray-700">Inscriptions :</span> {{ selectedFormation.inscription_count ?? 0 }}<span v-if="selectedFormation.capacity"> / {{ selectedFormation.capacity }}</span></div>
                        <div v-if="selectedFormation.category_name"><span class="font-medium text-gray-700">Thématique :</span> {{ selectedFormation.category_name }}</div>
                    </div>
                </div>
                <div v-if="selectedFormation.status !== 'approved'" class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button
                        @click="openEditFromDetail(selectedFormation)"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
                    >Modifier</button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer la formation ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimée.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteFormation" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
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
                        {{ editTarget ? 'Modifier la formation' : 'Créer une formation' }}
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
                            placeholder="Titre de la formation"/>
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea v-model="form.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de la formation…"/>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Niveau</label>
                            <select v-model="form.level"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                                <option value="beginner">Débutant</option>
                                <option value="intermediate">Intermédiaire</option>
                                <option value="advanced">Avancé</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Durée (heures)</label>
                            <input v-model.number="form.duration_hours" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Ex: 3"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date & heure <span class="text-red-400">*</span></label>
                            <input v-model="form.date" type="datetime-local" :min="minDateTime" required
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Lieu <span class="text-red-400">*</span></label>
                            <input v-model="form.location" type="text" required
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse ou salle"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Capacité</label>
                            <input v-model.number="form.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Thématique</label>
                            <select v-model="form.id_category"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                                <option value="">Aucune</option>
                                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Prix (€) <span class="text-gray-300">(laisser vide si gratuit)</span></label>
                            <input v-model.number="form.price" type="number" min="0" step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Ex: 25.00"/>
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

    </SalarieLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SalarieLayout from '@/Layouts/SalarieLayout.vue'
import api from '@/api.js'

const route = useRoute()

const formations = ref([])
const categories = ref([])
const loading = ref(true)
const toDelete = ref(null)
const deleting = ref(false)
const selectedFormation = ref(null)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', date: '', location: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '', price: null })
const filterStatus = ref('')
const filteredFormations = computed(() =>
    filterStatus.value ? formations.value.filter(f => f.status === filterStatus.value) : formations.value
)

const minDateTime = computed(() => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 1)
    return now.toISOString().slice(0, 16)
})

function levelLabel(level) {
    return { beginner: 'Débutant', intermediate: 'Intermédiaire', advanced: 'Avancé' }[level] ?? level
}

function formatDateTime(dateStr) {
    if (!dateStr) return ''
    const d = new Date(dateStr)
    return d.toLocaleDateString('fr-FR') + ' à ' + d.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

function endDateTime(dateStr, hours) {
    if (!dateStr || !hours) return null
    return new Date(new Date(dateStr).getTime() + hours * 3600000)
}

function levelClass(level) {
    return {
        beginner: 'bg-secondary/10 text-secondary',
        intermediate: 'bg-primary/10 text-primary',
        advanced: 'bg-red-100 text-red-500',
    }[level] ?? 'bg-gray-100 text-gray-500'
}

function statusLabel(status) {
    return { pending: 'En attente', approved: 'Approuvée', rejected: 'Rejetée' }[status] ?? status
}

function statusClass(status) {
    return {
        pending: 'bg-yellow-100 text-yellow-700',
        approved: 'bg-green-100 text-green-700',
        rejected: 'bg-red-100 text-red-600',
    }[status] ?? 'bg-gray-100 text-gray-500'
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', date: '', location: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '' }
    formError.value = ''
    formModal.value = true
}

function openEditFromDetail(f) {
    selectedFormation.value = null
    editTarget.value = f
    form.value = {
        title: f.title,
        description: f.description ?? '',
        date: f.date ? f.date.slice(0, 16) : '',
        location: f.location ?? '',
        capacity: f.capacity ?? null,
        level: f.level,
        duration_hours: f.duration_hours ?? null,
        id_category: f.id_category ?? '',
        price: f.price ?? null,
    }
    formError.value = ''
    formModal.value = true
}

async function fetchFormations() {
    loading.value = true
    try {
        const { data } = await api.get('/user/my-formations')
        formations.value = data ?? []
    } catch {
        formations.value = []
    } finally {
        loading.value = false
    }
}

async function deleteFormation() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/formation/${toDelete.value.id}`)
        formations.value = formations.value.filter(f => f.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

async function submitForm() {
    if (!form.value.title.trim()) { formError.value = 'Le titre est requis.'; return }
    if (!form.value.date) { formError.value = 'La date est requise.'; return }
    if (!form.value.location?.trim()) { formError.value = 'Le lieu est requis.'; return }
    submitting.value = true
    formError.value = ''
    try {
        const payload = {
            title: form.value.title,
            description: form.value.description || null,
            date: form.value.date || null,
            location: form.value.location || null,
            capacity: form.value.capacity || null,
            level: form.value.level,
            duration_hours: form.value.duration_hours || null,
            id_category: form.value.id_category || null,
            price: form.value.price || null,
        }
        if (editTarget.value) {
            await api.patch(`/formations/${editTarget.value.id}`, payload)
        } else {
            await api.post('/formations', payload)
        }
        formModal.value = false
        await fetchFormations()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

async function fetchCategories() {
    try {
        const { data } = await api.get('/categories')
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {}
}

onMounted(async () => {
    await Promise.all([fetchFormations(), fetchCategories()])
    const openId = route.query.open
    if (openId) {
        const found = formations.value.find(f => String(f.id) === String(openId))
        if (found) selectedFormation.value = found
    }
})
</script>
