<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Annonces</h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-5 gap-5 mb-8">
                <div
                    v-for="stat in stats"
                    :key="stat.label"
                    class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                    <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                        <div :class="stat.iconClass" v-html="stat.icon" />
                    </div>
                    <div class="min-w-0">
                        <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                        <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
                <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h2 class="text-base font-semibold text-gray-800">Annonces en attente de validation</h2>
                    <span class="text-xs text-gray-400">{{ pendingListings.length }} annonce(s) à traiter</span>
                </div>
                <div class="divide-y divide-gray-50">
                    <div
                        v-for="listing in pendingListings"
                        :key="listing.id"
                        class="px-5 py-3.5 flex items-center justify-between hover:bg-gray-50/60 transition-colors">
                        <div class="flex-1 mr-4">
                            <div class="flex items-center gap-2">
                                <p class="font-medium text-gray-800 text-sm truncate">{{ listing.title }}</p>
                                <span
                                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                                    :class="typeBadge(listing.type)">
                                    {{ listing.type }}
                                </span>
                            </div>
                            <p class="text-xs text-gray-500 mt-0.5">
                                Par <span class="font-medium text-gray-700">{{ listing.author }}</span>
                                · {{ listing.category }} · {{ listing.createdAt }}
                            </p>
                        </div>
                        <div class="flex-shrink-0 flex items-center gap-1">
                            <button
                                @click="openDetail(listing)"
                                title="Voir"
                                class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                </svg>
                            </button>
                            <button
                                @click="approveAnnouncement(listing)"
                                title="Approuver"
                                class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                </svg>
                            </button>
                            <button
                                @click="rejectAnnouncement(listing)"
                                title="Rejeter"
                                class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                </svg>
                            </button>
                            <button
                                @click="confirmDelete(listing)"
                                title="Supprimer"
                                class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                </svg>
                            </button>
                        </div>
                    </div>
                    <div v-if="pendingListings.length === 0" class="px-5 py-8 text-center text-gray-400 text-sm">
                        Aucune annonce en attente de validation
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm overflow-hidden">

                <div class="px-5 py-4 border-b border-gray-100 flex items-center gap-3">
                    <div class="relative flex-1">
                        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                        </svg>
                        <input
                            v-model="search"
                            type="text"
                            placeholder="Rechercher par titre ou auteur…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>

                    <select
                        v-model="filterType"
                        @change="resetAndFetch"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les types</option>
                        <option value="don">Don</option>
                        <option value="vente">Vente</option>
                    </select>

                    <select
                        v-model="filterCategory"
                        @change="resetAndFetch"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Toutes les catégories</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                    </select>

                    <select
                        v-model="filterStatus"
                        @change="resetAndFetch"
                        class="text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option>En attente</option>
                        <option>Active</option>
                        <option>Vendu</option>
                        <option>Refusée</option>
                        <option>Expirée</option>
                        <option>Supprimée</option>
                    </select>

                    <button
                        @click="openCreate"
                        class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors whitespace-nowrap">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                        </svg>
                        Créer
                    </button>

                    <span class="text-xs text-gray-400 whitespace-nowrap">
                        {{ total }} résultat(s)
                    </span>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Annonce</th>
                                <th class="text-left text-white font-medium px-5 py-3">Type</th>
                                <th class="text-left text-white font-medium px-5 py-3">Catégorie</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-left text-white font-medium px-5 py-3">Dates</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="(listing, i) in listings"
                                :key="listing.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                                <td class="px-5 py-3">
                                    <div class="flex items-center gap-2">
                                        <svg
                                            v-if="listing.featured"
                                            class="w-3.5 h-3.5 text-amber-400 flex-shrink-0"
                                            viewBox="0 0 24 24"
                                            fill="currentColor">
                                            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                                        </svg>
                                        <div class="min-w-0">
                                            <p class="font-medium text-gray-800 truncate max-w-48">{{ listing.title }}</p>
                                            <p class="text-xs text-gray-500 truncate">{{ listing.author }}</p>
                                        </div>
                                    </div>
                                </td>

                                <td class="px-5 py-3">
                                    <span
                                        class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                                        :class="typeBadge(listing.type)">
                                        {{ listing.type }}
                                    </span>
                                </td>

                                <td class="px-5 py-3 text-gray-500 text-xs">{{ listing.category }}</td>

                                <td class="px-5 py-3">
                                    <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                        <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="statusDot(listing.status)" />
                                        <span :class="statusText(listing.status)">{{ listing.status }}</span>
                                    </span>
                                </td>

                                <td class="px-5 py-3 text-xs">
                                    <p class="text-gray-500">Créé le {{ listing.createdAt }}</p>
                                    <p v-if="listing.deletedAt" class="text-red-400 mt-0.5">Supprimé le {{ listing.deletedAt }}</p>
                                </td>

                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button
                                            @click="openDetail(listing)"
                                            title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>
                                        <template v-if="listing.status === 'En attente'">
                                            <button
                                                @click="approveAnnouncement(listing)"
                                                title="Approuver"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                </svg>
                                            </button>
                                            <button
                                                @click="rejectAnnouncement(listing)"
                                                title="Rejeter"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                </svg>
                                            </button>
                                        </template>
                                        <template v-else>
                                            <button
                                                @click="toggleFeatured(listing)"
                                                :title="listing.featured ? 'Retirer la mise en avant' : 'Mettre en avant'"
                                                :class="[
                                                    'p-1.5 rounded-lg transition-colors',
                                                    listing.featured
                                                        ? 'text-amber-400 hover:text-amber-600 hover:bg-amber-50'
                                                        : 'text-gray-400 hover:text-amber-400 hover:bg-amber-50'
                                                ]">
                                                <svg class="w-4 h-4" viewBox="0 0 24 24" :fill="listing.featured ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                                                </svg>
                                            </button>
                                            <button
                                                @click="openEdit(listing)"
                                                title="Modifier l'annonce"
                                                class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                                                </svg>
                                            </button>
                                        </template>
                                        <button
                                            @click="confirmDelete(listing)"
                                            title="Supprimer l'annonce"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>

                            <tr v-if="listings.length === 0">
                                <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucune annonce ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ total }} annonce(s) au total
                </div>
            </div>
            <Pagination v-if="total > 20" :page="page" :total="total" :limit="20" @update:page="changePage" />

        </template>

        <div
            v-if="detailListing"
            class="fixed inset-0 z-50 flex items-center justify-center p-4"
            @click.self="detailListing = null">
            <div class="absolute inset-0 bg-black/40" @click="detailListing = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4">
                    <div class="min-w-0">
                        <div class="flex items-center gap-2">
                            <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                                {{ detailListing.title }}
                            </h3>
                            <span
                                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                                :class="typeBadge(detailListing.type)">
                                {{ detailListing.type }}
                            </span>
                        </div>
                        <p class="text-xs text-gray-500 mt-0.5">{{ detailListing.category }}</p>
                    </div>
                    <button
                        @click="detailListing = null"
                        class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Auteur</p>
                            <p class="font-medium text-gray-800">{{ detailListing.author }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span class="inline-flex items-center gap-1.5 text-xs font-medium">
                                <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(detailListing.status)" />
                                <span :class="statusText(detailListing.status)">{{ detailListing.status }}</span>
                            </span>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Date de création</p>
                            <p class="text-gray-700">{{ detailListing.createdAt }}</p>
                            <template v-if="detailListing.deletedAt">
                                <p class="text-xs text-gray-400 mb-0.5 mt-2">Date de suppression</p>
                                <p class="text-red-500">{{ detailListing.deletedAt }}</p>
                            </template>
                        </div>
                        <div v-if="detailListing.featured && detailListing.featuredUntil">
                            <p class="text-xs text-gray-400 mb-0.5">En avant jusqu'au</p>
                            <p class="text-amber-600 font-medium">{{ detailListing.featuredUntil.slice(0, 10) }}</p>
                        </div>
                    </div>

                    <div>
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailListing.description }}</p>
                    </div>

                    <div v-if="detailListing.status === 'Refusée' && detailListing.rejectionReason" class="p-3 rounded-xl bg-red-50 border border-red-100">
                        <p class="text-xs font-medium text-red-600 mb-0.5">Motif de refus</p>
                        <p class="text-sm text-red-700">{{ detailListing.rejectionReason }}</p>
                    </div>

                    <div v-if="detailListing.tags && detailListing.tags.length">
                        <p class="text-xs text-gray-400 mb-1.5">Tags</p>
                        <div class="flex gap-1.5">
                            <span
                                v-for="tag in detailListing.tags"
                                :key="tag"
                                class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs">
                                {{ tag }}
                            </span>
                        </div>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-between items-center">
                    <button
                        @click="openEdit(detailListing)"
                        class="px-3 py-1.5 text-sm font-medium text-primary border border-primary/30 rounded-lg hover:bg-primary/5 transition-colors">
                        Modifier
                    </button>
                    <button
                        v-if="detailListing.status !== 'En attente'"
                        @click="toggleFeatured(detailListing); detailListing = null"
                        :class="[
                            'flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg border transition-colors',
                            detailListing.featured
                                ? 'border-amber-300 text-amber-600 hover:bg-amber-50'
                                : 'border-gray-200 text-gray-600 hover:bg-gray-50'
                        ]">
                        <svg class="w-4 h-4" viewBox="0 0 24 24" :fill="detailListing.featured ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                        {{ detailListing.featured ? 'Retirer la mise en avant' : 'Mettre en avant' }}
                    </button>
                    <button
                        @click="confirmDelete(detailListing); detailListing = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Supprimer
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
                    <button
                        @click="toDelete = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="deleteListing"
                        :disabled="deleting"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="rejectModal.open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="rejectModal.open = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-1">Refuser l'annonce</h3>
                <p class="text-sm text-gray-500 mb-4">« {{ rejectModal.listing?.title }} »</p>
                <div class="mb-4">
                    <label class="block text-xs text-gray-400 mb-1">Motif de refus (optionnel)</label>
                    <textarea
                        v-model="rejectModal.reason"
                        rows="3"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-300 focus:border-red-400 resize-none"
                        placeholder="Expliquez pourquoi cette annonce est refusée…" />
                </div>
                <div class="flex gap-3">
                    <button @click="rejectModal.open = false" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="confirmRejectAnnouncement" :disabled="rejectModal.loading" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ rejectModal.loading ? 'Refus…' : 'Refuser' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="formModal.open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal.open = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ formModal.mode === 'create' ? 'Créer une annonce' : 'Modifier l\'annonce' }}
                    </h3>
                    <button
                        @click="formModal.open = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4 max-h-[70vh] overflow-y-auto">
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input
                            v-model="listingForm.title"
                            type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre de l'annonce" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Type</label>
                            <select
                                v-model="listingForm.type"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-700">
                                <option value="don">Don</option>
                                <option value="vente">Vente</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Catégorie</label>
                            <select
                                v-model="listingForm.idCategory"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-700">
                                <option value="">- Aucune -</option>
                                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                            </select>
                        </div>
                    </div>
                    <div>
                        <label class="block text-xs text-gray-400 mb-1">Description</label>
                        <textarea
                            v-model="listingForm.description"
                            rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de l'annonce…" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div class="col-span-2">
                            <label class="block text-xs text-gray-400 mb-1">Adresse</label>
                            <input
                                v-model="listingForm.address"
                                type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Ville</label>
                            <input
                                v-model="listingForm.city"
                                type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Ville" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Code postal</label>
                            <input
                                v-model="listingForm.postal"
                                type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="75000" />
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Date de disponibilité <span class="text-red-400">*</span></label>
                            <input
                                v-model="listingForm.availDate"
                                type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div v-if="listingForm.type === 'vente'">
                            <label class="block text-xs text-gray-400 mb-1">Prix (€)</label>
                            <input
                                v-model.number="listingForm.price"
                                type="number"
                                min="0.01"
                                step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0.00" />
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Condition</label>
                            <select
                                v-model="listingForm.condition"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-700">
                                <option value="">-</option>
                                <option value="neuf">Neuf</option>
                                <option value="tres_bon">Très bon état</option>
                                <option value="bon">Bon état</option>
                                <option value="acceptable">Acceptable</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-xs text-gray-400 mb-1">Statut</label>
                            <select
                                v-model="listingForm.state"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-gray-700">
                                <option value="En attente">En attente</option>
                                <option value="Active">Active</option>
                                <option value="Refusée">Refusée</option>
                                <option value="Vendu">Vendu</option>
                                <option value="Expirée">Expirée</option>
                                <option value="Supprimée">Supprimée</option>
                            </select>
                        </div>
                    </div>
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
                    <button
                        @click="formModal.open = false"
                        class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="saveListing"
                        :disabled="saving"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ saving ? 'Enregistrement…' : (formModal.mode === 'create' ? 'Créer' : 'Enregistrer') }}
                    </button>
                </div>
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const listings = ref([])
const total = ref(0)
const page = ref(1)
const categories = ref([])
const loading = ref(true)
const error = ref('')
const statsData = ref({})
const search = ref('')
const filterType = ref('')
const filterCategory = ref('')
const filterStatus = ref('')
const detailListing = ref(null)
const toDelete = ref(null)
const rejectModal = ref({ open: false, listing: null, reason: '', loading: false })
const deleting = ref(false)
const formModal = ref({ open: false, mode: 'create', id: null })
const listingForm = ref({ title: '', type: 'don', idCategory: '', description: '', address: '', city: '', postal: '', availDate: '', price: 0, condition: '', state: 'Active' })
const saving = ref(false)

const stats = computed(() => [
    {
        label: 'Total annonces',
        value: statsData.value.total ?? '-',
        bgClass: 'bg-red-100',
        iconClass: 'text-red-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>`,
    },
    {
        label: 'Actives',
        value: statsData.value.active ?? '-',
        bgClass: 'bg-green-100',
        iconClass: 'text-green-600',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Mises en avant',
        value: statsData.value.featured ?? '-',
        bgClass: 'bg-amber-100',
        iconClass: 'text-amber-500',
        icon: `<svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>`,
    },
    {
        label: 'En attente',
        value: statsData.value.pending ?? '-',
        bgClass: 'bg-orange-100',
        iconClass: 'text-orange-500',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Signalées',
        value: statsData.value.reported ?? '-',
        bgClass: 'bg-red-100',
        iconClass: 'text-red-400',
        icon: `<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"/></svg>`,
    },
])

const pendingListings = ref([])

function resetAndFetch() {
    page.value = 1
    fetchListings()
}

function changePage(p) { page.value = p }

async function fetchListings() {
    try {
        const params = { page: page.value, limit: 20 }
        if (search.value) params.search = search.value
        if (filterType.value) params.type = filterType.value
        if (filterStatus.value) params.status = filterStatus.value
        if (filterCategory.value) params.id_category = filterCategory.value
        const { data: announcementsData } = await api.get('/admin/announcements', { params })
        const catMap = Object.fromEntries(categories.value.map(c => [c.id, c.name]))
        listings.value = announcementsData.data.map(a => ({
            id: a.id,
            title: a.title,
            author: a.author_name || '-',
            type: a.type === 'vente' ? 'Vente' : 'Don',
            category: catMap[a.id_category] ?? '-',
            status: a.state ?? 'Active',
            featured: a.is_featured === 1,
            featuredUntil: a.featured_until ?? null,
            date: a.availability_date?.slice(0, 10) ?? '-',
            createdAt: a.created_at?.slice(0, 10) ?? '-',
            deletedAt: a.deleted_at ? a.deleted_at.slice(0, 16).replace('T', ' ') : null,
            description: a.description ?? '-',
            rejectionReason: a.rejection_reason ?? null,
            tags: [],
            idCategory: a.id_category ?? '',
            address: a.address ?? '',
            city: a.city ?? '',
            postal: a.postal ?? '',
            price: a.price ?? 0,
            condition: a.condition ?? '',
        }))
        total.value = announcementsData.total
    } catch {
        error.value = 'Impossible de charger les annonces.'
    }
}

let debounce = null
watch(search, () => {
    clearTimeout(debounce)
    debounce = setTimeout(() => { page.value = 1; fetchListings() }, 300)
})
watch(page, fetchListings)

async function fetchPendingListings() {
    try {
        const { data: pendingData } = await api.get('/admin/announcements', { params: { status: 'En attente', limit: 100 } })
        const catMap = Object.fromEntries(categories.value.map(c => [c.id, c.name]))
        pendingListings.value = (pendingData.data ?? []).map(a => ({
            id: a.id,
            title: a.title,
            author: a.author_name || '-',
            type: a.type === 'vente' ? 'Vente' : 'Don',
            category: catMap[a.id_category] ?? '-',
            status: a.state ?? 'En attente',
            createdAt: a.created_at?.slice(0, 10) ?? '-',
            rawCreatedAt: a.created_at ?? '',
            date: a.availability_date?.slice(0, 10) ?? '-',
            description: a.description ?? '-',
            rejectionReason: a.rejection_reason ?? null,
            idCategory: a.id_category ?? '',
            address: a.address ?? '',
            city: a.city ?? '',
            postal: a.postal ?? '',
            price: a.price ?? 0,
            condition: a.condition ?? '',
            featured: a.is_featured === 1,
            featuredUntil: a.featured_until ?? null,
            tags: [],
        })).sort((a, b) => b.rawCreatedAt.localeCompare(a.rawCreatedAt))
    } catch {}
}

onMounted(async () => {
    try {
        const [{ data: catsData }, { data: announcementStats }] = await Promise.all([
            api.get('/admin/categories'),
            api.get('/admin/announcements/stats'),
        ])
        categories.value = catsData ?? []
        statsData.value = announcementStats
        await Promise.all([fetchPendingListings(), fetchListings()])
    } catch {
        error.value = 'Impossible de charger les annonces.'
    } finally {
        loading.value = false
    }
})

function openDetail(listing) {
    detailListing.value = listing
}

function confirmDelete(listing) {
    toDelete.value = listing
}

async function toggleFeatured(listing) {
    try {
        await api.put(`/admin/announcement/${listing.id}/feature`, { active: !listing.featured })
        await fetchListings()
    } catch {
        alert('Erreur lors de la mise à jour de la mise en avant.')
    }
}

async function deleteListing() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/announcement/${toDelete.value.id}`)
        if (detailListing.value?.id === toDelete.value.id) detailListing.value = null
        toDelete.value = null
        await Promise.all([fetchListings(), fetchPendingListings()])
    } catch {
        alert('Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function openCreate() {
    listingForm.value = { title: '', type: 'don', idCategory: '', description: '', address: '', city: '', postal: '', availDate: '', price: 0, condition: '', state: 'Active' }
    formModal.value = { open: true, mode: 'create', id: null }
}

function openEdit(listing) {
    listingForm.value = {
        title: listing.title,
        type: listing.type === 'Vente' ? 'vente' : 'don',
        idCategory: listing.idCategory ?? '',
        description: listing.description === '-' ? '' : listing.description,
        address: listing.address ?? '',
        city: listing.city ?? '',
        postal: listing.postal ?? '',
        availDate: listing.date !== '-' ? listing.date : '',
        price: listing.price ?? 0,
        condition: listing.condition ?? '',
        state: listing.status ?? 'Active',
    }
    formModal.value = { open: true, mode: 'edit', id: listing.id }
    detailListing.value = null
}

async function saveListing() {
    if (!listingForm.value.title.trim() || !listingForm.value.availDate) return
    saving.value = true
    try {
        const catMap = Object.fromEntries(categories.value.map(c => [c.id, c.name]))
        const payload = {
            title: listingForm.value.title,
            type: listingForm.value.type,
            id_category: listingForm.value.idCategory,
            description: listingForm.value.description,
            address: listingForm.value.address,
            city: listingForm.value.city,
            postal: listingForm.value.postal,
            availability_date: listingForm.value.availDate,
            price: listingForm.value.type === 'vente' ? listingForm.value.price : 0,
            condition: listingForm.value.condition,
            state: listingForm.value.state,
        }
        if (formModal.value.mode === 'create') {
            await api.post('/admin/announcements', payload)
            page.value = 1
        } else {
            await api.put(`/admin/announcement/${formModal.value.id}`, payload)
        }
        formModal.value.open = false
        await fetchListings()
    } catch {
        alert('Erreur lors de l\'enregistrement.')
    } finally {
        saving.value = false
    }
}

async function approveAnnouncement(listing) {
    try {
        await api.patch(`/admin/announcement/${listing.id}/approve`)
        await Promise.all([fetchListings(), fetchPendingListings()])
    } catch {
        alert('Erreur lors de l\'approbation.')
    }
}

function rejectAnnouncement(listing) {
    rejectModal.value = { open: true, listing, reason: '', loading: false }
}

async function confirmRejectAnnouncement() {
    if (!rejectModal.value.listing) return
    rejectModal.value.loading = true
    try {
        await api.patch(`/admin/announcement/${rejectModal.value.listing.id}/reject`, { reason: rejectModal.value.reason })
        rejectModal.value.open = false
        await Promise.all([fetchListings(), fetchPendingListings()])
    } catch {
        alert('Erreur lors du rejet.')
    } finally {
        rejectModal.value.loading = false
    }
}

function typeBadge(type) {
    const map = { 'Don': 'bg-green-100 text-green-700', 'Vente': 'bg-blue-100 text-blue-700' }
    return map[type] ?? 'bg-gray-100 text-gray-600'
}

function statusDot(status) {
    const map = { 'Active': 'bg-green-500', 'En attente': 'bg-orange-400', 'Expirée': 'bg-gray-400', 'Supprimée': 'bg-gray-400', 'Refusée': 'bg-red-400' }
    return map[status] ?? 'bg-gray-300'
}

function statusText(status) {
    const map = { 'Active': 'text-green-700', 'En attente': 'text-orange-600', 'Expirée': 'text-gray-500', 'Supprimée': 'text-gray-500', 'Refusée': 'text-red-600' }
    return map[status] ?? 'text-gray-500'
}
</script>
