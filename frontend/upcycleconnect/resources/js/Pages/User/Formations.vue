<template>
    <UserLayout>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                    Formations
                </h1>
                <p class="text-sm text-gray-400 mt-1">Apprenez les techniques d'upcycling avec nos formations</p>
            </div>
            <button
                v-if="canCreate"
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer une formation
            </button>
        </div>

        <!-- Section mes formations créées -->
        <div v-if="canCreate && myCreated.length > 0" class="bg-white rounded-2xl shadow-sm p-6 mb-6">
            <h2 class="font-semibold text-gray-800 mb-4" style="font-family: var(--font-family-title)">
                Mes formations créées
            </h2>
            <div class="space-y-3">
                <div
                    v-for="f in myCreated"
                    :key="f.id"
                    class="flex items-center justify-between p-4 rounded-xl border border-gray-100 hover:border-gray-200 transition-colors">
                    <div class="min-w-0 flex-1">
                        <div class="flex items-center gap-2 mb-0.5">
                            <p class="font-medium text-gray-800 truncate">{{ f.title }}</p>
                            <span :class="statusClass(f.status)" class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0">
                                {{ statusLabel(f.status) }}
                            </span>
                        </div>
                        <p v-if="f.status === 'rejected' && f.rejection_reason" class="text-xs text-red-500 mt-0.5">
                            Raison : {{ f.rejection_reason }}
                        </p>
                        <p class="text-xs text-gray-400">{{ f.date ? f.date.slice(0,10) : 'Date non définie' }} · {{ levelLabel(f.level) }}</p>
                    </div>
                    <div class="flex items-center gap-2 flex-shrink-0 ml-4">
                        <button
                            v-if="f.status !== 'approved'"
                            @click="openEdit(f)"
                            class="p-1.5 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-colors"
                            title="Modifier">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                            </svg>
                        </button>
                        <button
                            @click="confirmDelete(f)"
                            class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                            title="Supprimer">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Section validation (manage_formations) -->
        <div v-if="canManage" class="bg-white rounded-2xl shadow-sm p-6 mb-6">
            <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">
                Formations en attente de validation
            </h2>
            <p class="text-xs text-gray-400 mb-4">Formations soumises par les créateurs, à approuver ou rejeter</p>

            <div v-if="loadingPending" class="text-center py-8 text-gray-400 text-sm">Chargement…</div>
            <p v-else-if="pendingFormations.length === 0" class="text-sm text-gray-400 text-center py-6">
                Aucune formation en attente de validation.
            </p>
            <div v-else class="space-y-3">
                <div
                    v-for="f in pendingFormations"
                    :key="f.id"
                    class="flex items-center justify-between p-4 rounded-xl border border-gray-100 hover:border-gray-200 transition-colors">
                    <div class="min-w-0 flex-1">
                        <p class="font-medium text-gray-800 truncate">{{ f.title }}</p>
                        <p class="text-xs text-gray-400 mt-0.5">
                            Par {{ f.creator_name }} · {{ levelLabel(f.level) }}
                            <span v-if="f.date"> · {{ f.date.slice(0,10) }}</span>
                        </p>
                        <p v-if="f.description" class="text-xs text-gray-500 mt-1 line-clamp-1">{{ f.description }}</p>
                    </div>
                    <div class="flex items-center gap-2 flex-shrink-0 ml-4">
                        <button
                            @click="openRejectModal(f)"
                            class="px-3 py-1.5 text-xs font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Rejeter
                        </button>
                        <button
                            @click="approveFormation(f)"
                            class="px-3 py-1.5 text-xs font-semibold bg-secondary text-white rounded-lg hover:bg-secondary-dark transition-colors">
                            Approuver
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Filtres + liste formations publiques -->
        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        @input="onSearchInput"
                        type="text"
                        placeholder="Rechercher une formation…"
                        class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                </div>
                <select v-model="filterLevel" @change="resetAndFetch"
                    class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Tous les niveaux</option>
                    <option value="beginner">Débutant</option>
                    <option value="intermediate">Intermédiaire</option>
                    <option value="advanced">Avancé</option>
                </select>
                <select v-if="categories.length" v-model="filterCategory" @change="resetAndFetch"
                    class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Toutes les thématiques</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} résultat(s)</span>
            </div>

            <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>

            <div v-else-if="formations.length === 0" class="py-12 text-center">
                <svg class="w-12 h-12 text-gray-200 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                <p class="text-gray-400 text-sm">Aucune formation disponible pour le moment.</p>
            </div>

            <div v-else class="grid grid-cols-3 gap-5 p-6">
                <div
                    v-for="f in formations"
                    :key="f.id"
                    class="bg-gray-50 rounded-2xl overflow-hidden hover:shadow-md transition-shadow cursor-pointer flex flex-col border border-gray-100"
                    @click="openDetail(f)">
                    <div class="h-32 bg-secondary/10 flex items-center justify-center">
                        <svg class="w-10 h-10 text-secondary/40" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                        </svg>
                    </div>
                    <div class="p-4 flex flex-col gap-2 flex-1">
                        <div class="flex items-center gap-2">
                            <span :class="levelClass(f.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ levelLabel(f.level) }}
                            </span>
                            <span v-if="f.duration_hours" class="text-xs text-gray-400">{{ f.duration_hours }}h</span>
                        </div>
                        <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ f.title }}</h3>
                        <p class="text-xs text-gray-500 line-clamp-2">{{ f.description }}</p>
                        <div class="mt-auto pt-2 flex items-center justify-between">
                            <span class="text-xs text-gray-400">Par {{ f.creator_name }}</span>
                            <span v-if="f.is_registered" class="text-xs text-secondary font-medium">Inscrit</span>
                            <span v-else-if="f.capacity && f.inscription_count >= f.capacity" class="text-xs text-gray-400">Complet</span>
                        </div>
                    </div>
                </div>
            </div>

            <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />
        </div>

        <!-- Modal détail formation -->
        <div v-if="detailFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="detailFormation = null">
            <div class="absolute inset-0 bg-black/40" @click="detailFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ detailFormation.title }}</h3>
                        <p class="text-xs text-gray-500 mt-0.5">Animée par {{ detailFormation.formateur_name ?? detailFormation.creator_name }}</p>
                    </div>
                    <button @click="detailFormation = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Niveau</p>
                            <span :class="levelClass(detailFormation.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ levelLabel(detailFormation.level) }}
                            </span>
                        </div>
                        <div v-if="detailFormation.duration_hours">
                            <p class="text-xs text-gray-400 mb-0.5">Durée</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.duration_hours }}h</p>
                        </div>
                        <div v-if="detailFormation.date">
                            <p class="text-xs text-gray-400 mb-0.5">Date</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.date.slice(0,10) }}</p>
                            <p class="text-xs text-gray-500">{{ detailFormation.date.slice(11,16) }}</p>
                        </div>
                        <div v-if="detailFormation.location">
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.location }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Inscriptions</p>
                            <p class="font-medium text-gray-800">
                                {{ detailFormation.inscription_count }}
                                <span v-if="detailFormation.capacity"> / {{ detailFormation.capacity }}</span>
                                participants
                            </p>
                        </div>
                        <div v-if="detailFormation.category_name">
                            <p class="text-xs text-gray-400 mb-0.5">Thématique</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.category_name }}</p>
                        </div>
                    </div>
                    <div v-if="detailFormation.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailFormation.description }}</p>
                    </div>
                    <div v-if="detailFormation.is_registered" class="flex items-center gap-2 p-3 rounded-xl bg-secondary/10">
                        <svg class="w-4 h-4 text-secondary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        <p class="text-sm text-secondary font-medium">Vous êtes inscrit à cette formation</p>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button
                        v-if="detailFormation.is_registered"
                        @click="toggleRegistration(detailFormation); detailFormation = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Se désinscrire
                    </button>
                    <button
                        v-else-if="!isFull(detailFormation)"
                        @click="toggleRegistration(detailFormation); detailFormation = null"
                        class="px-3 py-1.5 text-sm font-semibold bg-secondary text-white rounded-lg hover:bg-secondary-dark transition-colors">
                        S'inscrire
                    </button>
                    <span v-else class="px-3 py-1.5 text-sm text-gray-400">Complet</span>
                </div>
            </div>
        </div>

        <!-- Modal rejet -->
        <div v-if="rejectModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="rejectModal = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Rejeter la formation</h3>
                <p class="text-sm text-gray-500 mb-4">« {{ rejectModal.title }} »</p>
                <div class="mb-4">
                    <label class="block text-xs text-gray-400 mb-1">Raison (optionnelle)</label>
                    <textarea
                        v-model="rejectReason"
                        rows="3"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                        placeholder="Expliquez pourquoi la formation est rejetée…"/>
                </div>
                <div class="flex gap-3">
                    <button @click="rejectModal = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="submitReject" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors">Rejeter</button>
                </div>
            </div>
        </div>

        <!-- Modal confirmation suppression -->
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

        <!-- Modal création / édition -->
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
                            <label class="block text-xs text-gray-400 mb-1">Date & heure</label>
                            <input v-model="form.date" type="datetime-local" :min="minDateTime"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Lieu</label>
                            <input v-model="form.location" type="text"
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
                                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                                    {{ cat.name }}
                                </option>
                            </select>
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
    </UserLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
import UserLayout from '@/Layouts/UserLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'
import { useAuthStore } from '@/stores/auth.js'

const auth = useAuthStore()

const formations = ref([])
const myCreated = ref([])
const pendingFormations = ref([])
const categories = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const loadingPending = ref(false)
const search = ref('')
const filterLevel = ref('')
const filterCategory = ref('')
const detailFormation = ref(null)
const rejectModal = ref(null)
const rejectReason = ref('')
const toDelete = ref(null)
const deleting = ref(false)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', date: '', location: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '' })

const canCreate = computed(() => auth.hasPermission('create_formation'))
const canManage = computed(() => auth.hasPermission('manage_formations'))
const canRegister = computed(() => auth.hasPermission('register_formation'))

const minDateTime = computed(() => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 1)
    return now.toISOString().slice(0, 16)
})

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(resetAndFetch, 300)
}

function resetAndFetch() {
    page.value = 1
    fetchFormations()
}

function levelLabel(level) {
    return { beginner: 'Débutant', intermediate: 'Intermédiaire', advanced: 'Avancé' }[level] ?? level
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

function isFull(f) {
    return f.capacity > 0 && f.inscription_count >= f.capacity
}

function openDetail(f) {
    detailFormation.value = f
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', date: '', location: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '' }
    formError.value = ''
    formModal.value = true
}

function openEdit(f) {
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
    }
    formError.value = ''
    formModal.value = true
}

function confirmDelete(f) {
    toDelete.value = f
}

async function deleteFormation() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/formation/${toDelete.value.id}`)
        myCreated.value = myCreated.value.filter(f => f.id !== toDelete.value.id)
        formations.value = formations.value.filter(f => f.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

async function toggleRegistration(f) {
    if (!canRegister.value) return
    f.loading = true
    try {
        if (f.is_registered) {
            await api.delete(`/user/formation/${f.id}/unregister`)
            f.is_registered = false
            f.inscription_count = Math.max(0, f.inscription_count - 1)
        } else {
            await api.post(`/user/formation/${f.id}/register`)
            f.is_registered = true
            f.inscription_count++
        }
    } catch (err) {
        alert(err.response?.data?.error ?? 'Erreur lors de l\'inscription.')
    } finally {
        f.loading = false
    }
}

async function approveFormation(f) {
    try {
        await api.patch(`/formations/${f.id}/approve`)
        pendingFormations.value = pendingFormations.value.filter(p => p.id !== f.id)
        await fetchFormations()
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de l\'approbation.')
    }
}

function openRejectModal(f) {
    rejectModal.value = f
    rejectReason.value = ''
}

async function submitReject() {
    if (!rejectModal.value) return
    try {
        await api.patch(`/formations/${rejectModal.value.id}/reject`, { reason: rejectReason.value || null })
        pendingFormations.value = pendingFormations.value.filter(p => p.id !== rejectModal.value.id)
        rejectModal.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors du rejet.')
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
            date: form.value.date || null,
            location: form.value.location || null,
            capacity: form.value.capacity || null,
            level: form.value.level,
            duration_hours: form.value.duration_hours || null,
            id_category: form.value.id_category || null,
        }
        if (editTarget.value) {
            await api.patch(`/formations/${editTarget.value.id}`, payload)
        } else {
            await api.post('/formations', payload)
        }
        formModal.value = false
        await fetchMyCreated()
    } catch (e) {
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function fetchFormations(silent = false) {
    if (!silent) loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        if (filterLevel.value) params.level = filterLevel.value
        if (filterCategory.value) params.id_category = filterCategory.value
        const { data } = await api.get('/formations', { params })
        formations.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (e) {
        console.error('fetchFormations error:', e)
    } finally {
        if (!silent) loading.value = false
    }
}

async function fetchMyCreated() {
    if (!canCreate.value) return
    try {
        const { data } = await api.get('/user/my-formations')
        myCreated.value = data ?? []
    } catch (e) {
        console.error('fetchMyCreated error:', e)
    }
}

async function fetchPending(silent = false) {
    if (!canManage.value) return
    if (!silent) loadingPending.value = true
    try {
        const { data } = await api.get('/formations/pending')
        pendingFormations.value = data.data ?? []
    } catch (e) {
        console.error('fetchPending error:', e)
    } finally {
        if (!silent) loadingPending.value = false
    }
}

async function fetchCategories() {
    try {
        const { data } = await api.get('/categories')
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {}
}

async function fetchAll(silent = false) {
    await Promise.all([
        fetchFormations(silent),
        fetchMyCreated(),
        fetchPending(silent),
        fetchCategories(),
    ])
}

usePolling(fetchAll)
</script>
