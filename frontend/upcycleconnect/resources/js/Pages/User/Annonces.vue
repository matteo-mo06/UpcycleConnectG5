<template>
    <UserLayout @openDepot="showDepot = true">

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Annonces</h1>
            <button
                v-if="auth.hasPermission('create_announcement')"
                @click="showDepot = true"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-lg hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Publier
            </button>
        </div>

        <div class="flex items-center gap-4 mb-6">
            <div class="flex gap-1 bg-gray-100 rounded-lg p-1">
                <button
                    v-for="tab in tabs"
                    :key="tab.value"
                    @click="switchTab(tab.value)"
                    :class="[
                        'px-4 py-1.5 rounded-md text-sm font-medium transition-colors',
                        activeTab === tab.value ? 'bg-white text-gray-800 shadow-sm' : 'text-gray-500 hover:text-gray-700'
                    ]"
                >{{ tab.label }}</button>
            </div>

            <template v-if="activeTab === 'all'">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        @input="debouncedFetch"
                        type="text"
                        placeholder="Rechercher une annonce…"
                        class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                    />
                </div>
                <select v-model="filterType" @change="resetAndFetch" class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600">
                    <option value="">Tous les types</option>
                    <option value="vente">Vente</option>
                    <option value="don">Don</option>
                </select>
                <select v-model="filterCategory" @change="resetAndFetch" class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600">
                    <option value="">Toutes les catégories</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
            </template>
        </div>

        <div v-if="loading" class="text-center py-12 text-gray-400 text-sm">Chargement…</div>

        <div v-else-if="announcements.length === 0" class="text-center py-12 text-gray-400 text-sm">
            {{ activeTab === 'mine' ? 'Vous n\'avez aucune annonce.' : 'Aucune annonce disponible.' }}
        </div>

        <div v-else class="grid grid-cols-3 gap-5">
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

                    <div class="mt-3 flex items-center gap-2">
                        <span v-if="a.condition" class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">{{ conditionLabel(a.condition) }}</span>
                    </div>

                    <p class="text-xs text-gray-400 mt-2">Disponible le {{ formatDate(a.availability_date) }}</p>

                    <div class="mt-3 flex items-center justify-between">
                        <div class="flex items-center gap-1.5">
                            <div class="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-white text-xs font-semibold">
                                {{ a.author_name?.[0] ?? '?' }}
                            </div>
                            <span class="text-xs text-gray-500">{{ a.author_name || 'Inconnu' }}</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <template v-if="activeTab === 'mine'">
                                <span v-if="a.state === 'En attente'" class="px-2 py-0.5 bg-orange-100 text-orange-600 rounded-full text-xs font-medium">En attente de validation</span>
                                <span v-else-if="a.state === 'Active'" class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium">Validée</span>
                                <span v-else-if="a.state === 'Refusée'" class="px-2 py-0.5 bg-red-100 text-red-600 rounded-full text-xs font-medium" :title="a.rejection_reason ?? ''">Refusée{{ a.rejection_reason ? ' - voir détail' : '' }}</span>
                                <span v-else-if="a.state === 'Vendu'" class="px-2 py-0.5 bg-gray-100 text-gray-500 rounded-full text-xs font-medium">{{ a.type === 'vente' ? 'Vendu' : 'Donné' }}</span>
                            </template>

                            <button
                                v-if="activeTab === 'mine' && a.state === 'Vendu' && a.request === 0 && !a.locker_number"
                                @click.stop="requestDeposit(a)"
                                class="p-1.5 text-blue-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                                title="Déposer l'objet dans un casier"
                            >
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                                </svg>
                            </button>
                            <span
                                v-if="activeTab === 'mine' && a.state === 'Vendu' && a.request === 1 && !a.locker_number"
                                class="px-2 py-0.5 bg-yellow-100 text-yellow-700 rounded-full text-xs font-medium"
                                title="En attente d'assignation de casier"
                            >En attente de casier</span>
                            <span
                                v-if="activeTab === 'mine' && a.state === 'Vendu' && a.locker_number"
                                class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium"
                            >Casier N° {{ a.locker_number }}</span>

                            <span
                                v-if="activeTab === 'acquisitions' && a.locker_number"
                                class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium"
                            >Casier N° {{ a.locker_number }}</span>

                            <button
                                v-if="activeTab === 'acquisitions' && a.invoice_path"
                                @click.stop="downloadInvoice(a)"
                                class="p-1.5 text-blue-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                                title="Télécharger la facture"
                            >
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 10v6m0 0l-3-3m3 3l3-3M3 17V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
                                </svg>
                            </button>

                            <button
                                v-if="canDelete(a)"
                                @click.stop="confirmDelete(a)"
                                class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                title="Supprimer"
                            >
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                </svg>
                            </button>

                            <button
                                v-if="activeTab === 'all' && a.author_id !== auth.user?.id"
                                @click.stop="handleAcquire(a)"
                                class="px-3 py-1 bg-secondary text-white text-xs font-medium rounded-lg hover:bg-secondary-dark transition-colors"
                            >
                                {{ a.type === 'vente' ? 'Acheter' : 'Je le veux' }}
                            </button>

                            <button
                                v-if="activeTab === 'all' && a.author_id !== auth.user?.id"
                                @click.stop="openReport(a)"
                                class="px-3 py-1 border border-gray-200 text-gray-500 text-xs font-medium rounded-lg hover:border-red-300 hover:text-red-500 transition-colors"
                            >
                                Signaler
                            </button>

                        </div>
                    </div>
                </div>
            </div>
        </div>

        <Pagination v-if="activeTab === 'all' && total > 12" :page="page" :total="total" :limit="12" @update:page="changePage" />

        <CreateAnnouncementModal v-model="showDepot" @created="onCreated" />
        <ReportModal v-model="showReport" contentType="announcement" :contentId="reportTarget?.id" :contentTitle="reportTarget?.title ?? ''" />

        <div v-if="selected" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="closeDetail">
            <div class="absolute inset-0 bg-black/40" @click="closeDetail" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selected.title }}</h3>
                    <button @click="closeDetail" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>

                <div v-if="selectedPhotos.length" class="relative bg-gray-100">
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

                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto">
                    <div class="flex gap-2">
                        <span :class="['px-2 py-0.5 rounded-full text-xs font-semibold text-white', selected.type === 'vente' ? 'bg-primary' : 'bg-secondary']">{{ selected.type === 'vente' ? 'Vente' : 'Don' }}</span>
                        <span v-if="selected.condition" class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">{{ conditionLabel(selected.condition) }}</span>
                    </div>
                    <p class="leading-relaxed">{{ selected.description }}</p>
                    <div class="grid grid-cols-2 gap-4 text-xs text-gray-500">
                        <div><span class="font-medium text-gray-700">Ville :</span> {{ selected.city }}</div>
                        <div><span class="font-medium text-gray-700">Prix :</span> {{ selected.price ? Number(selected.price).toFixed(2) + ' €' : 'Gratuit' }}</div>
                        <div><span class="font-medium text-gray-700">Disponible :</span> {{ formatDate(selected.availability_date) }}</div>
                        <div><span class="font-medium text-gray-700">Par :</span> {{ selected.author_name || 'Inconnu' }}</div>
                    </div>
                    <div v-if="activeTab === 'mine' && selected.state === 'Refusée' && selected.rejection_reason" class="p-3 rounded-xl bg-red-50 border border-red-100">
                        <p class="text-xs font-medium text-red-600 mb-0.5">Motif de refus</p>
                        <p class="text-sm text-red-700">{{ selected.rejection_reason }}</p>
                    </div>

                    <div v-if="activeTab === 'all' && selected.author_id !== auth.user?.id" class="pt-2 flex justify-end">
                        <button @click="openReport(selected)" class="text-xs text-gray-400 hover:text-red-500 transition-colors">Signaler ce contenu</button>
                    </div>

                    <div v-if="activeTab === 'acquisitions' && selected.locker_number" class="pt-2 border-t border-gray-100 space-y-2">
                        <p class="text-xs text-gray-500"><span class="font-medium text-gray-700">Casier :</span> N° {{ selected.locker_number }}</p>
                        <template v-if="selected.access_code">
                            <p class="text-xs text-gray-400">Scannez ce code pour ouvrir le casier</p>
                            <svg ref="barcodeEl" class="mx-auto" />
                        </template>
                        <p v-else class="text-xs text-yellow-600">Code d'accès non encore disponible</p>
                    </div>
                </div>
                <div v-if="canEdit(selected)" class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button
                        @click="openEdit(selected)"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
                    >Modifier</button>
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

        <!-- Modale de modification -->
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

    </UserLayout>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
import { useRoute, useRouter } from 'vue-router'
import UserLayout from '@/Layouts/UserLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import ReportModal from '@/Components/ReportModal.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'
import { formatDate, conditionLabel } from '@/utils.js'
import JsBarcode from 'jsbarcode'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const announcements = ref([])
const loading = ref(true)
const search = ref('')
const filterType = ref('')
const filterCategory = ref('')
const categories = ref([])
const page = ref(1)
const total = ref(0)
const showDepot = ref(false)
const selected = ref(null)
const selectedPhotos = ref([])
const photoIndex = ref(0)
const toDelete = ref(null)
const deleting = ref(false)
const reportTarget = ref(null)
const showReport = ref(false)

const validTabs = ['all', 'mine', 'acquisitions']
const activeTab = ref(validTabs.includes(route.query.tab) ? route.query.tab : 'all')

const tabs = [
    { label: 'Toutes', value: 'all' },
    { label: 'Mes annonces', value: 'mine' },
    { label: 'Mes acquisitions', value: 'acquisitions' },
]

let debounceTimer = null
function debouncedFetch() {
    clearTimeout(debounceTimer)
    page.value = 1
    debounceTimer = setTimeout(fetchAnnouncements, 300)
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

function resetAndFetch() {
    page.value = 1
    fetchAnnouncements()
}

function switchTab(tab) {
    activeTab.value = tab
    search.value = ''
    filterType.value = ''
    filterCategory.value = ''
    page.value = 1
    router.replace({ query: { tab } })
    fetchAnnouncements()
}

async function fetchAnnouncements(silent = false) {
    if (!silent) loading.value = true
    try {
        if (activeTab.value === 'mine') {
            const { data } = await api.get('/user/announcements')
            announcements.value = data
            total.value = 0
        } else if (activeTab.value === 'acquisitions') {
            const { data } = await api.get('/user/acquisitions')
            announcements.value = data
            total.value = 0
        } else {
            const params = { page: page.value, limit: 12 }
            if (search.value) params.search = search.value
            if (filterType.value) params.type = filterType.value
            if (filterCategory.value) params.id_category = filterCategory.value
            const { data } = await api.get('/announcements', { params })
            announcements.value = data.data ?? []
            total.value = data.total ?? 0
        }
    } catch {
        announcements.value = []
        total.value = 0
    } finally {
        loading.value = false
    }
}

async function openDetail(a) {
    selected.value = a
    photoIndex.value = 0
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

function canDelete(a) {
    if (a?.state !== 'Active') return false
    if (a.author_id === auth.user?.id) return true
    return auth.isAdmin || auth.hasPermission('manage_announcements')
}

function confirmDelete(a) {
    toDelete.value = a
}

async function deleteAnnouncement() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        const isOwner = toDelete.value.author_id === auth.user?.id
        if (isOwner) {
            await api.delete(`/user/announcement/${toDelete.value.id}`)
        } else if (auth.isAdmin) {
            await api.delete(`/admin/announcement/${toDelete.value.id}`)
        } else {
            await api.delete(`/announcements/${toDelete.value.id}`)
        }
        announcements.value = announcements.value.filter(a => a.id !== toDelete.value.id)
        toDelete.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function handleAcquire(a) {
    if (a.type === 'vente' && a.price > 0) {
        router.push(`/paiement/${a.id}`)
    } else {
        claimAnnouncement(a)
    }
}

async function claimAnnouncement(a) {
    try {
        await api.post(`/announcements/${a.id}/claim`)
        announcements.value = announcements.value.filter(x => x.id !== a.id)
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de l\'acquisition.')
    }
}

function openReport(a) {
    reportTarget.value = { id: a.id, title: a.title }
    showReport.value = true
}

async function downloadInvoice(a) {
    try {
        const { data } = await api.get(`/user/announcement/${a.id}/invoice`, { responseType: 'blob' })
        const url = URL.createObjectURL(data)
        const link = document.createElement('a')
        link.href = url
        link.download = `facture-${a.id}.pdf`
        link.click()
        URL.revokeObjectURL(url)
    } catch {
        alert('Facture non disponible.')
    }
}

async function requestDeposit(a) {
    try {
        const { data } = await api.post(`/announcements/${a.id}/deposit-request`)
        a.request = 1
        if (data?.locker_number) {
            a.locker_number = data.locker_number
        }
    } catch (e) {
        const msg = e.response?.data?.error ?? 'Erreur lors de la demande de dépôt.'
        alert(msg)
    }
}

const editModal = ref(false)
const editSubmitting = ref(false)
const editError = ref('')
const editTarget = ref(null)
const editForm = ref({ type: 'don', title: '', id_category: '', condition: 'bon_etat', description: '', address: '', city: '', postal: '', price: 0, availability_date: '' })

function canEdit(a) {
    return activeTab.value === 'mine' && (a?.state === 'En attente' || a?.state === 'Refusée')
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

const barcodeEl = ref(null)

watch(() => selected.value?.access_code, (code) => {
    if (!code) return
    nextTick(() => {
        if (barcodeEl.value) {
            JsBarcode(barcodeEl.value, code, { format: 'CODE128', width: 2, height: 50, displayValue: true, fontSize: 12 })
        }
    })
})

function onCreated() {
    switchTab('mine')
}


watch(page, () => { if (activeTab.value === 'all') fetchAnnouncements() })

async function fetchCategories() {
    try {
        const { data } = await api.get('/categories')
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    } catch {}
}

onMounted(() => fetchCategories())
usePolling(fetchAnnouncements, 2000, () => activeTab.value === 'mine')
</script>
