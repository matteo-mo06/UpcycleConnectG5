<template>
    <ArtisanLayout>

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Annonces</h1>
            <button
                @click="showCreate = true"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-xl hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Publier
            </button>
        </div>

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
                        placeholder="Rechercher une annonce…"
                        class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                </div>
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ announcements.length }} annonce(s)</span>
            </div>

            <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">Chargement…</div>

            <div v-else-if="announcements.length === 0" class="py-12 text-center text-gray-400 text-sm">
                {{ search ? 'Aucune annonce pour cette recherche.' : 'Vous n\'avez aucune annonce publiée.' }}
            </div>

            <div v-else class="grid grid-cols-3 gap-5 p-6">
            <div
                v-for="a in announcements"
                :key="a.id"
                @click="openDetail(a)"
                class="bg-white rounded-2xl shadow-sm overflow-hidden hover:shadow-md transition-shadow cursor-pointer"
            >
                <div class="relative h-40 bg-gray-100 flex items-center justify-center">
                    <img v-if="a.first_photo" :src="a.first_photo" class="w-full h-full object-cover" :alt="a.title" />
                    <svg v-else class="w-10 h-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"/>
                    </svg>
                    <span :class="['absolute top-2 right-2 px-2 py-0.5 rounded-full text-xs font-semibold text-white', a.type === 'vente' ? 'bg-primary' : 'bg-secondary']">
                        {{ a.type === 'vente' ? 'VENTE' : 'DON' }}
                    </span>
                </div>

                <div class="p-4">
                    <h3 class="font-semibold text-gray-800 truncate">{{ a.title }}</h3>
                    <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ a.description }}</p>
                    <p class="text-xs text-gray-400 mt-2">Disponible le {{ formatDate(a.availability_date) }}</p>

                    <div class="mt-3 flex items-center justify-between">
                        <div class="flex items-center gap-1.5 flex-wrap">
                            <span v-if="a.state === 'En attente'" class="px-2 py-0.5 bg-orange-100 text-orange-600 rounded-full text-xs font-medium">En attente</span>
                            <span v-else-if="a.state === 'Active'" class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium">Active</span>
                            <span v-else-if="a.state === 'Refusée'" class="px-2 py-0.5 bg-red-100 text-red-600 rounded-full text-xs font-medium">Refusée</span>
                            <span v-else-if="a.state === 'Vendu'" class="px-2 py-0.5 bg-gray-100 text-gray-500 rounded-full text-xs font-medium">{{ a.type === 'vente' ? 'Vendu' : 'Donné' }}</span>
                            <span v-if="a.is_featured === 1" class="inline-flex items-center gap-1 px-2 py-0.5 bg-amber-100 text-amber-600 rounded-full text-xs font-medium">
                                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
                                En avant
                            </span>
                            <span v-else-if="a.featured_requested_at && a.is_featured !== 1" class="px-2 py-0.5 bg-amber-50 text-amber-500 rounded-full text-xs font-medium">File d'attente</span>
                            <button
                                v-if="a.state === 'Vendu' && a.request === 0 && !a.locker_number"
                                @click.stop="requestDeposit(a)"
                                class="p-1.5 text-blue-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                                title="Déposer dans un casier"
                            >
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                                </svg>
                            </button>
                            <span v-if="a.state === 'Vendu' && a.request === 1 && !a.locker_number" class="px-2 py-0.5 bg-yellow-100 text-yellow-700 rounded-full text-xs font-medium">Casier en attente</span>
                            <span v-if="a.state === 'Vendu' && a.locker_number" class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium">Casier N° {{ a.locker_number }}</span>
                        </div>
                        <button
                            v-if="a.state === 'Active'"
                            @click.stop="toDelete = a"
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

        <div v-if="selected" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="closeDetail">
            <div class="absolute inset-0 bg-black/40" @click="closeDetail" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selected.title }}</h3>
                    <button @click="closeDetail" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>

                <div v-if="selectedPhotos.length" class="relative bg-gray-100 flex-shrink-0">
                    <img :src="selectedPhotos[photoIndex]" class="w-full h-56 object-cover" :alt="selected.title" />
                    <template v-if="selectedPhotos.length > 1">
                        <button @click="photoIndex = (photoIndex - 1 + selectedPhotos.length) % selectedPhotos.length"
                            class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/40 hover:bg-black/60 rounded-full flex items-center justify-center text-white transition-colors">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>
                        </button>
                        <button @click="photoIndex = (photoIndex + 1) % selectedPhotos.length"
                            class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/40 hover:bg-black/60 rounded-full flex items-center justify-center text-white transition-colors">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/></svg>
                        </button>
                        <div class="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
                            <button v-for="(_, i) in selectedPhotos" :key="i" @click="photoIndex = i"
                                :class="['w-2 h-2 rounded-full transition-colors', i === photoIndex ? 'bg-white' : 'bg-white/50']" />
                        </div>
                    </template>
                </div>

                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto flex-1">
                    <div class="flex gap-2 flex-wrap">
                        <span :class="['px-2 py-0.5 rounded-full text-xs font-semibold text-white', selected.type === 'vente' ? 'bg-primary' : 'bg-secondary']">{{ selected.type === 'vente' ? 'Vente' : 'Don' }}</span>
                        <span v-if="selected.condition" class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">{{ conditionLabel(selected.condition) }}</span>
                        <span v-if="selected.state === 'En attente'" class="px-2 py-0.5 bg-orange-100 text-orange-600 rounded-full text-xs font-medium">En attente</span>
                        <span v-else-if="selected.state === 'Active'" class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium">Active</span>
                        <span v-else-if="selected.state === 'Refusée'" class="px-2 py-0.5 bg-red-100 text-red-600 rounded-full text-xs font-medium">Refusée</span>
                        <span v-else-if="selected.state === 'Vendu'" class="px-2 py-0.5 bg-gray-100 text-gray-500 rounded-full text-xs font-medium">{{ selected.type === 'vente' ? 'Vendu' : 'Donné' }}</span>
                    </div>
                    <p class="leading-relaxed">{{ selected.description }}</p>
                    <div class="grid grid-cols-2 gap-4 text-xs text-gray-500">
                        <div><span class="font-medium text-gray-700">Ville :</span> {{ selected.city }}</div>
                        <div><span class="font-medium text-gray-700">Prix :</span> {{ selected.price ? Number(selected.price).toFixed(2) + ' €' : 'Gratuit' }}</div>
                        <div><span class="font-medium text-gray-700">Disponible :</span> {{ formatDate(selected.availability_date) }}</div>
                    </div>
                    <div v-if="selected.state === 'Refusée' && selected.rejection_reason" class="p-3 rounded-xl bg-red-50 border border-red-100">
                        <p class="text-xs font-medium text-red-600 mb-0.5">Motif de refus</p>
                        <p class="text-sm text-red-700">{{ selected.rejection_reason }}</p>
                    </div>
                </div>

                <div v-if="selected.state === 'Active' && isPremium" class="px-6 py-4 border-t border-gray-100 flex-shrink-0">
                    <div v-if="selected.is_featured === 1" class="flex items-center justify-between">
                        <span class="text-xs text-amber-600">
                            En avant jusqu'au {{ selected.featured_until ? selected.featured_until.slice(0, 10) : '—' }}
                        </span>
                        <button
                            @click="cancelFeature"
                            :disabled="featureLoading"
                            class="px-3 py-1.5 text-xs text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-60"
                        >Retirer</button>
                    </div>
                    <div v-else-if="selected.featured_requested_at" class="flex items-center justify-between gap-3">
                        <span class="text-xs text-amber-600">
                            En file d'attente{{ featureResult?.estimated_slot ? ' — slot estimé le ' + featureResult.estimated_slot.slice(0, 10) : '' }}
                        </span>
                        <div class="flex gap-2">
                            <button
                                @click="confirmFeature"
                                :disabled="featureLoading"
                                class="px-3 py-1.5 text-xs font-medium bg-amber-400 text-white rounded-lg hover:bg-amber-500 transition-colors disabled:opacity-60"
                            >{{ featureLoading ? '…' : 'Confirmer' }}</button>
                            <button
                                @click="cancelFeature"
                                :disabled="featureLoading"
                                class="px-3 py-1.5 text-xs text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-60"
                            >Annuler</button>
                        </div>
                    </div>
                    <div v-else class="flex items-center justify-between gap-3">
                        <span class="text-xs text-gray-400">Mise en avant · 7 jours · 1/mois</span>
                        <button
                            @click="requestFeature"
                            :disabled="featureLoading"
                            class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium border border-amber-300 text-amber-600 rounded-lg hover:bg-amber-50 transition-colors disabled:opacity-60"
                        >
                            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                            </svg>
                            {{ featureLoading ? '…' : 'Mettre en avant' }}
                        </button>
                    </div>
                    <p v-if="featureResult?.error" class="text-xs text-red-600 mt-2">{{ featureResult.error }}</p>
                    <p v-if="featureResult?.status === 'active'" class="text-xs text-green-700 mt-2">Annonce mise en avant jusqu'au {{ featureResult.featured_until?.slice(0, 10) }}.</p>
                </div>

                <div v-if="canEdit(selected) || selected.state === 'Active'" class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button
                        v-if="selected.state === 'Active'"
                        @click="toDelete = selected; closeDetail()"
                        class="px-4 py-2 text-sm text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors"
                    >Supprimer</button>
                    <button
                        v-if="canEdit(selected)"
                        @click="openEdit(selected)"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
                    >Modifier</button>
                </div>
            </div>
        </div>

        <div v-if="editModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="editModal = false"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800 text-lg" style="font-family: var(--font-family-title)">Modifier l'annonce</h3>
                    <button @click="editModal = false" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Type de dépôt</label>
                        <div class="flex gap-2">
                            <button
                                v-for="t in [{ label: 'Don', value: 'don' }, { label: 'Vente', value: 'vente' }]"
                                :key="t.value"
                                @click="editForm.type = t.value; if (t.value === 'don') editForm.price = 0"
                                :class="['px-4 py-2 rounded-full text-sm font-medium border transition-colors', editForm.type === t.value ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-200 hover:border-primary hover:text-primary']"
                            >{{ t.label }}</button>
                        </div>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Titre</label>
                        <input v-model="editForm.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Catégorie</label>
                        <select v-model="editForm.id_category"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                            <option value="">Sélectionner une catégorie…</option>
                            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">État de l'objet</label>
                        <div class="flex gap-2">
                            <button
                                v-for="c in [{ label: 'Neuf', value: 'neuf' }, { label: 'Bon état', value: 'bon_etat' }, { label: 'Usé', value: 'use' }]"
                                :key="c.value"
                                @click="editForm.condition = c.value"
                                :class="['px-4 py-2 rounded-full text-sm font-medium border transition-colors', editForm.condition === c.value ? 'bg-secondary text-white border-secondary' : 'bg-white text-gray-600 border-gray-200 hover:border-secondary hover:text-secondary-dark']"
                            >{{ c.label }}</button>
                        </div>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea v-model="editForm.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none" />
                    </div>
                    <div class="grid grid-cols-2 gap-3">
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-gray-700 mb-1">Adresse</label>
                            <input v-model="editForm.address" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Ville</label>
                            <input v-model="editForm.city" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Code postal</label>
                            <input v-model="editForm.postal" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                    </div>
                    <div v-if="editForm.type === 'vente'">
                        <label class="block text-sm font-medium text-gray-700 mb-1">Prix (€)</label>
                        <input v-model.number="editForm.price" type="number" min="0.01" step="0.01"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Disponible à partir du</label>
                        <input v-model="editForm.availability_date" type="date"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                    </div>
                    <p v-if="editError" class="text-sm text-red-600">{{ editError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex justify-end gap-2">
                    <button @click="editModal = false" class="px-4 py-2 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="submitEdit" :disabled="editSubmitting"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ editSubmitting ? 'Enregistrement…' : 'Enregistrer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer l'annonce ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimée.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteAnnouncement" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <CreateAnnouncementModal v-model="showCreate" @created="fetchAnnouncements" />

    </ArtisanLayout>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import api from '@/api.js'
import { formatDate, conditionLabel } from '@/utils.js'

const announcements = ref([])
const categories = ref([])
const search = ref('')
const loading = ref(true)
const isPremium = ref(false)
const route = useRoute()
const router = useRouter()
const showCreate = ref(false)
const selected = ref(null)
const selectedPhotos = ref([])
const photoIndex = ref(0)
const toDelete = ref(null)
const deleting = ref(false)
const featureLoading = ref(false)
const featureResult = ref(null)
const editModal = ref(false)
const editSubmitting = ref(false)
const editError = ref('')
const editTarget = ref(null)
const editForm = ref({ type: 'don', title: '', id_category: '', condition: 'bon_etat', description: '', address: '', city: '', postal: '', price: 0, availability_date: '' })

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(fetchAnnouncements, 300)
}

async function fetchAnnouncements() {
    loading.value = true
    try {
        const params = {}
        if (search.value) params.search = search.value
        const { data } = await api.get('/user/announcements', { params })
        announcements.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {
        announcements.value = []
    } finally {
        loading.value = false
    }
}

async function silentFetch() {
    try {
        const params = {}
        if (search.value) params.search = search.value
        const { data } = await api.get('/user/announcements', { params })
        announcements.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {}
}

let pollInterval = null

async function openDetail(a) {
    selected.value = a
    photoIndex.value = 0
    featureResult.value = null
    try {
        const { data } = await api.get(`/announcements/${a.id}/documents`)
        selectedPhotos.value = data.map(d => d.link).filter(Boolean)
    } catch {
        selectedPhotos.value = a.first_photo ? [a.first_photo] : []
    }
}

function closeDetail() {
    selected.value = null
    selectedPhotos.value = []
    photoIndex.value = 0
}

function canEdit(a) {
    return a?.state === 'En attente' || a?.state === 'Refusée'
}

function openEdit(a) {
    editTarget.value = a
    editForm.value = {
        type: a.type ?? 'don',
        title: a.title ?? '',
        id_category: a.id_category ?? '',
        condition: a.condition ?? 'bon_etat',
        description: a.description ?? '',
        address: a.address ?? '',
        city: a.city ?? '',
        postal: a.postal ?? '',
        price: a.price ?? 0,
        availability_date: a.availability_date ? a.availability_date.slice(0, 10) : '',
    }
    editError.value = ''
    closeDetail()
    editModal.value = true
}

async function submitEdit() {
    if (!editForm.value.title.trim()) { editError.value = 'Le titre est requis.'; return }
    editSubmitting.value = true
    editError.value = ''
    try {
        await api.patch(`/user/announcement/${editTarget.value.id}`, {
            title: editForm.value.title,
            type: editForm.value.type,
            id_category: editForm.value.id_category,
            condition: editForm.value.condition,
            description: editForm.value.description,
            address: editForm.value.address,
            city: editForm.value.city,
            postal: editForm.value.postal,
            price: editForm.value.price,
            availability_date: editForm.value.availability_date,
        })
        editModal.value = false
        await fetchAnnouncements()
    } catch (e) {
        editError.value = e.response?.data?.error ?? 'Erreur lors de la modification.'
    } finally {
        editSubmitting.value = false
    }
}

async function deleteAnnouncement() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/announcement/${toDelete.value.id}`)
        announcements.value = announcements.value.filter(a => a.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

async function requestDeposit(a) {
    try {
        await api.post(`/announcements/${a.id}/deposit-request`)
        a.request = 1
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la demande de dépôt.')
    }
}

async function fetchCategories() {
    try {
        const { data } = await api.get('/categories')
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {}
}

async function requestFeature() {
    if (!selected.value) return
    featureLoading.value = true
    featureResult.value = null
    try {
        const { data } = await api.post(`/announcements/${selected.value.id}/feature`)
        featureResult.value = data
        await fetchAnnouncements()
        const refreshed = announcements.value.find(a => a.id === selected.value.id)
        if (refreshed) selected.value = refreshed
    } catch (e) {
        featureResult.value = { error: e.response?.data?.error ?? 'Erreur serveur' }
    } finally {
        featureLoading.value = false
    }
}

async function confirmFeature() {
    if (!selected.value) return
    featureLoading.value = true
    featureResult.value = null
    try {
        const { data } = await api.post(`/announcements/${selected.value.id}/feature/confirm`)
        featureResult.value = data
        await fetchAnnouncements()
        const refreshed = announcements.value.find(a => a.id === selected.value.id)
        if (refreshed) selected.value = refreshed
    } catch (e) {
        featureResult.value = e.response?.data ?? { error: 'Erreur serveur' }
    } finally {
        featureLoading.value = false
    }
}

async function cancelFeature() {
    if (!selected.value) return
    featureLoading.value = true
    featureResult.value = null
    try {
        await api.delete(`/announcements/${selected.value.id}/feature`)
        await fetchAnnouncements()
        const refreshed = announcements.value.find(a => a.id === selected.value.id)
        if (refreshed) selected.value = refreshed
    } catch {
        featureResult.value = { error: 'Erreur lors de l\'annulation.' }
    } finally {
        featureLoading.value = false
    }
}

async function fetchLimits() {
    try {
        const { data } = await api.get('/user/limits')
        isPremium.value = data.is_premium === true
    } catch {}
}

onMounted(() => {
    Promise.all([fetchAnnouncements(), fetchCategories(), fetchLimits()])
    if (route.query.publish === '1') {
        showCreate.value = true
        router.replace({ query: { ...route.query, publish: undefined } })
    }
    pollInterval = setInterval(silentFetch, 2000)
})

onUnmounted(() => {
    if (pollInterval) { clearInterval(pollInterval); pollInterval = null }
})
</script>
