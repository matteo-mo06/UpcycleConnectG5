<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Projets</h1>
        </div>

        <div v-if="loading" class="py-12 text-center text-sm text-gray-400">Chargement…</div>
        <div v-else-if="error" class="py-12 text-center text-sm text-red-500">{{ error }}</div>

        <template v-else>

            <div class="grid grid-cols-4 gap-5 mb-8">
                <div v-for="stat in stats" :key="stat.label"
                    class="bg-white rounded-xl shadow-sm p-5 flex items-center gap-4">
                    <div :class="['flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center', stat.bgClass]">
                        <div :class="stat.iconClass" v-html="stat.icon"/>
                    </div>
                    <div class="min-w-0">
                        <p class="text-2xl font-bold text-gray-800 leading-none">{{ stat.value }}</p>
                        <p class="text-sm text-gray-500 mt-1 truncate">{{ stat.label }}</p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-2xl shadow-sm overflow-hidden">

                <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                    <div class="relative flex-1">
                        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                        </svg>
                        <input v-model="search" @input="onSearchInput" type="text" placeholder="Rechercher un projet…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                    </div>
                    <select v-model="filterStatus"
                        class="text-sm border border-gray-200 rounded-xl px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option value="pending">En attente</option>
                        <option value="open">Ouvert</option>
                        <option value="in_progress">En cours</option>
                        <option value="completed">Terminé</option>
                        <option value="rejected">Rejeté</option>
                    </select>
                    <button
                        @click="openCreate"
                        class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-xl hover:bg-primary-dark transition-colors whitespace-nowrap">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                        </svg>
                        Créer
                    </button>

                    <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} résultat(s)</span>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr class="bg-primary">
                                <th class="text-left text-white font-medium px-5 py-3">Projet</th>
                                <th class="text-left text-white font-medium px-5 py-3">Dates</th>
                                <th class="text-left text-white font-medium px-5 py-3">Lieu</th>
                                <th class="text-left text-white font-medium px-5 py-3">Membres</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, i) in projects" :key="p.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                                <td class="px-5 py-3">
                                    <p class="font-medium text-gray-800 truncate max-w-52">{{ p.title }}</p>
                                    <p class="text-xs text-gray-500 truncate">{{ p.creator_name }}</p>
                                </td>
                                <td class="px-5 py-3">
                                    <p class="text-gray-800 text-xs font-medium">{{ p.start_date ? p.start_date.slice(0,10) : '-' }}</p>
                                    <p class="text-gray-500 text-xs">{{ p.end_date ? p.end_date.slice(0,10) : '-' }}</p>
                                </td>
                                <td class="px-5 py-3 text-xs text-gray-500 truncate max-w-36">
                                    {{ p.location ?? '-' }}
                                </td>
                                <td class="px-5 py-3 text-xs text-gray-700 font-medium">
                                    {{ p.members_count }}<span v-if="p.capacity">/{{ p.capacity }}</span>
                                </td>
                                <td class="px-5 py-3">
                                    <span :class="statusClass(p.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                        {{ statusLabel(p.status) }}
                                    </span>
                                </td>
                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button v-if="p.status === 'pending'" @click="approveProject(p)"
                                            title="Approuver"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                            </svg>
                                        </button>
                                        <button v-if="p.status === 'pending'" @click="openRejectModal(p)"
                                            title="Rejeter"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                            </svg>
                                        </button>
                                        <button @click="openDetail(p)" title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>
                                        <button @click="confirmDelete(p)" title="Supprimer"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                            <tr v-if="projects.length === 0">
                                <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucun projet ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-6 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ total }} projet(s) au total
                </div>
            </div>

            <Pagination v-if="total > 20" :page="page" :total="total" :limit="20" @update:page="changePage"/>

            <div class="mt-10">
                <h2 class="text-xl font-bold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Projets supprimés</h2>
                <p class="text-sm text-gray-400 mb-4">Historique des projets supprimés</p>
                <div v-if="deletedProjects.length === 0" class="bg-white rounded-xl shadow-sm p-8 text-center">
                    <p class="text-gray-400 text-sm">Aucun projet supprimé.</p>
                </div>
                <div v-else class="space-y-3">
                    <div
                        v-for="p in deletedProjects"
                        :key="p.id"
                        @click="selectedDeletedProject = p"
                        class="bg-white rounded-xl shadow-sm px-5 py-4 flex items-center gap-4 cursor-pointer hover:shadow-md transition-shadow">
                        <div class="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center flex-shrink-0">
                            <svg class="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </div>
                        <div class="flex-1 min-w-0">
                            <h3 class="font-semibold text-gray-500 text-sm truncate">{{ p.title }}</h3>
                            <p class="text-xs text-gray-400 mt-0.5">Par {{ p.creator_name }} · Supprimé le {{ p.deleted_at?.slice(0, 10) }}</p>
                        </div>
                    </div>
                </div>
            </div>

        </template>

        <div v-if="detailProject" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="detailProject = null">
            <div class="absolute inset-0 bg-black/40" @click="detailProject = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4 flex-shrink-0">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ detailProject.title }}</h3>
                        <p class="text-xs text-gray-500 mt-0.5">Par {{ detailProject.creator_name }}</p>
                    </div>
                    <button @click="detailProject = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span :class="statusClass(detailProject.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ statusLabel(detailProject.status) }}
                            </span>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Membres</p>
                            <p class="font-medium text-gray-800">{{ detailProject.members_count }}<span v-if="detailProject.capacity">/{{ detailProject.capacity }}</span></p>
                        </div>
                        <div v-if="detailProject.start_date">
                            <p class="text-xs text-gray-400 mb-0.5">Début</p>
                            <p class="font-medium text-gray-800">{{ detailProject.start_date.slice(0,10) }}</p>
                        </div>
                        <div v-if="detailProject.end_date">
                            <p class="text-xs text-gray-400 mb-0.5">Fin</p>
                            <p class="font-medium text-gray-800">{{ detailProject.end_date.slice(0,10) }}</p>
                        </div>
                        <div v-if="detailProject.location">
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ detailProject.location }}</p>
                        </div>
                    </div>
                    <div v-if="detailProject.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ detailProject.description }}</p>
                    </div>
                    <div v-if="detailProject.status === 'rejected' && detailProject.rejection_reason">
                        <p class="text-xs text-gray-400 mb-1">Raison du rejet</p>
                        <p class="text-sm text-red-600">{{ detailProject.rejection_reason }}</p>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Pièces jointes</p>
                        <p v-if="docsLoading" class="text-xs text-gray-400">Chargement…</p>
                        <p v-else-if="documents.length === 0" class="text-xs text-gray-400">Aucune pièce jointe.</p>
                        <div v-else class="flex gap-2 flex-wrap">
                            <div v-for="doc in documents" :key="doc.id" class="relative w-24 h-24">
                                <video v-if="isVideo(doc.link)" :src="doc.link" controls class="w-full h-full rounded-lg border border-gray-100 object-cover" />
                                <a v-else :href="doc.link" target="_blank" class="block w-full h-full rounded-lg overflow-hidden border border-gray-100">
                                    <img :src="doc.link" class="w-full h-full object-cover" />
                                </a>
                                <button
                                    @click="deleteDocument(doc)"
                                    class="absolute top-1 right-1 w-5 h-5 bg-black/60 rounded-full flex items-center justify-center text-white hover:bg-red-500 transition-colors"
                                    title="Supprimer"
                                >
                                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                                    </svg>
                                </button>
                            </div>
                        </div>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Participants</p>
                        <div class="space-y-1">
                            <p v-if="participantsLoading" class="text-xs text-gray-400">Chargement…</p>
                            <template v-else>
                                <p v-if="participants.length === 0" class="text-xs text-gray-400">Aucun participant</p>
                                <div v-for="p in participants" :key="p.id" class="flex items-center gap-2 py-1">
                                    <img v-if="p.avatar_url" :src="p.avatar_url" class="w-6 h-6 rounded-full object-cover" />
                                    <span v-else class="w-6 h-6 rounded-full bg-gray-200 text-xs flex items-center justify-center font-medium text-gray-500">
                                        {{ p.first_name[0] }}{{ p.last_name[0] }}
                                    </span>
                                    <span class="text-sm text-gray-700">{{ p.first_name }} {{ p.last_name }}</span>
                                </div>
                            </template>
                        </div>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-between gap-2 flex-shrink-0">
                    <div class="flex gap-2">
                        <button @click="openEdit(detailProject)"
                            class="px-3 py-1.5 text-sm font-medium text-primary border border-primary/30 rounded-lg hover:bg-primary/5 transition-colors">
                            Modifier
                        </button>
                        <button v-if="detailProject.status === 'pending'" @click="approveProject(detailProject); detailProject = null"
                            class="px-3 py-1.5 text-sm font-medium text-secondary border border-secondary/30 rounded-lg hover:bg-secondary/5 transition-colors">
                            Approuver
                        </button>
                        <button v-if="detailProject.status === 'pending'" @click="openRejectModal(detailProject); detailProject = null"
                            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Rejeter
                        </button>
                    </div>
                    <button @click="confirmDelete(detailProject); detailProject = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Supprimer
                    </button>
                </div>
            </div>
        </div>

        <div v-if="rejectModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="rejectModal = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Rejeter le projet</h3>
                <p class="text-sm text-gray-500 mb-4">« {{ rejectModal.title }} »</p>
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Raison (optionnelle)</label>
                    <textarea v-model="rejectReason" rows="3"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                        placeholder="Expliquez pourquoi le projet est rejeté…"/>
                </div>
                <div class="flex gap-3">
                    <button @click="rejectModal = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50">Annuler</button>
                    <button @click="submitReject" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600">Rejeter</button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le projet ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimé avec toutes ses inscriptions.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50">Annuler</button>
                    <button @click="deleteProject" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="formModal.open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal.open = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ formModal.mode === 'edit' ? 'Modifier le projet' : 'Créer un projet' }}</h3>
                    <button
                        @click="formModal.open = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input
                            v-model="projectForm.title"
                            type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre du projet" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea
                            v-model="projectForm.description"
                            rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description du projet…" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Lieu</label>
                        <input
                            v-model="projectForm.location"
                            type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Adresse ou ville" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Date de début</label>
                            <input
                                v-model="projectForm.startDate"
                                type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Date de fin</label>
                            <input
                                v-model="projectForm.endDate"
                                type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary" />
                        </div>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Capacité</label>
                        <input
                            v-model.number="projectForm.capacity"
                            type="number"
                            min="1"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="0 = illimitée" />
                    </div>
                    <div v-if="formModal.mode === 'create'">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Photos et vidéos <span class="text-gray-300 font-normal">(5 max)</span></label>
                        <div
                            class="border-2 border-dashed border-gray-200 rounded-xl p-5 text-center cursor-pointer hover:border-primary/50 transition-colors"
                            @dragover.prevent
                            @drop.prevent="handleDrop"
                            @click="$refs.fileInput.click()"
                        >
                            <svg class="w-7 h-7 text-gray-300 mx-auto mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"/>
                            </svg>
                            <p class="text-xs text-gray-400">Glissez vos fichiers ici</p>
                            <span class="mt-1 inline-block text-xs text-primary underline">Parcourir</span>
                            <input ref="fileInput" type="file" accept="image/*,video/*" multiple class="hidden" @change="handleFileSelect" />
                        </div>
                        <div v-if="photos.length" class="flex gap-2 mt-3 flex-wrap">
                            <div v-for="(photo, i) in photos" :key="i" class="relative w-14 h-14 rounded-lg overflow-hidden border border-gray-200">
                                <video v-if="photo.file.type.startsWith('video/')" :src="photo.preview" class="w-full h-full object-cover" />
                                <img v-else :src="photo.preview" class="w-full h-full object-cover" />
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
                </div>

                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2 flex-shrink-0">
                    <button
                        @click="formModal.open = false"
                        class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="saveProject"
                        :disabled="saving"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ saving ? 'Enregistrement…' : (formModal.mode === 'edit' ? 'Enregistrer' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="selectedDeletedProject" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="selectedDeletedProject = null">
            <div class="absolute inset-0 bg-black/40" @click="selectedDeletedProject = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[85vh] flex flex-col overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4 flex-shrink-0">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-600" style="font-family: var(--font-family-title)">{{ selectedDeletedProject.title }}</h3>
                        <p class="text-xs text-gray-400 mt-0.5">Par {{ selectedDeletedProject.creator_name }} · Supprimé le {{ selectedDeletedProject.deleted_at?.slice(0, 10) }}</p>
                    </div>
                    <span class="flex-shrink-0 px-2 py-1 rounded-full bg-gray-100 text-gray-400 text-xs font-medium">Supprimé</span>
                    <button @click="selectedDeletedProject = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div v-if="selectedDeletedProject.start_date">
                            <p class="text-xs text-gray-400 mb-0.5">Début</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedProject.start_date.slice(0, 10) }}</p>
                        </div>
                        <div v-if="selectedDeletedProject.end_date">
                            <p class="text-xs text-gray-400 mb-0.5">Fin</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedProject.end_date.slice(0, 10) }}</p>
                        </div>
                        <div v-if="selectedDeletedProject.location">
                            <p class="text-xs text-gray-400 mb-0.5">Lieu</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedProject.location }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Membres</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedProject.members_count }}<span v-if="selectedDeletedProject.capacity">/{{ selectedDeletedProject.capacity }}</span></p>
                        </div>
                    </div>
                    <div v-if="selectedDeletedProject.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed">{{ selectedDeletedProject.description }}</p>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Pièces jointes</p>
                        <p v-if="!selectedDeletedProject.documents?.length" class="text-xs text-gray-400">Aucune pièce jointe.</p>
                        <div v-else class="flex gap-2 flex-wrap">
                            <template v-for="doc in selectedDeletedProject.documents" :key="doc.id">
                                <video v-if="isVideo(doc.link)" :src="doc.link" controls class="w-24 h-24 rounded-lg border border-gray-100 object-cover" />
                                <a v-else :href="doc.link" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border border-gray-100">
                                    <img :src="doc.link" class="w-full h-full object-cover" />
                                </a>
                            </template>
                        </div>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Matériaux</p>
                        <p v-if="!selectedDeletedProject.materials?.length" class="text-xs text-gray-400">Aucun matériau renseigné.</p>
                        <table v-else class="w-full text-xs">
                            <thead>
                                <tr class="text-gray-400 border-b border-gray-100">
                                    <th class="text-left pb-1 font-medium">Nom</th>
                                    <th class="text-right pb-1 font-medium">Qté</th>
                                    <th class="text-left pb-1 font-medium pl-2">Unité</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="m in selectedDeletedProject.materials" :key="m.id" class="border-b border-gray-50 last:border-0">
                                    <td class="py-1.5 pr-2 text-gray-700">{{ m.name }}</td>
                                    <td class="py-1.5 text-right text-gray-600">{{ m.quantity }}</td>
                                    <td class="py-1.5 pl-2 text-gray-500">{{ m.unit ?? '-' }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Étapes</p>
                        <p v-if="!selectedDeletedProject.steps?.length" class="text-xs text-gray-400">Aucune étape renseignée.</p>
                        <div v-else class="space-y-1">
                            <div v-for="s in selectedDeletedProject.steps" :key="s.id" class="flex items-center gap-2 py-1.5 px-2 rounded-lg">
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
                            </div>
                        </div>
                    </div>
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

const projects = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const error = ref('')
const search = ref('')
const filterStatus = ref('')
const detailProject = ref(null)
const deletedProjects = ref([])
const selectedDeletedProject = ref(null)
const rejectModal = ref(null)
const rejectReason = ref('')
const toDelete = ref(null)
const deleting = ref(false)
const statsData = ref({ total: 0, pending: 0, open: 0, in_progress: 0, completed: 0, rejected: 0 })
const formModal = ref({ open: false })
const projectForm = ref({ title: '', description: '', location: '', startDate: '', endDate: '', capacity: null })
const saving = ref(false)
const participants = ref([])
const participantsLoading = ref(false)
const documents = ref([])
const docsLoading = ref(false)

function isVideo(link) {
    return /\.(mp4|webm|ogg|mov)(\?|$)/i.test(link)
}

const photos = ref([])

function handleFileSelect(e) {
    addFiles(Array.from(e.target.files))
    e.target.value = ''
}

function handleDrop(e) {
    addFiles(Array.from(e.dataTransfer.files).filter(f => f.type.startsWith('image/') || f.type.startsWith('video/')))
}

function addFiles(files) {
    for (const file of files) {
        if (photos.value.length >= 5) break
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

async function loadDocuments(id) {
    docsLoading.value = true
    try {
        const { data } = await api.get(`/projects/${id}/documents`)
        documents.value = data ?? []
    } catch {
        documents.value = []
    } finally {
        docsLoading.value = false
    }
}

async function deleteDocument(doc) {
    try {
        await api.delete(`/documents/${doc.id}`)
        documents.value = documents.value.filter(d => d.id !== doc.id)
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function loadParticipants(url) {
    participantsLoading.value = true
    try {
        const { data } = await api.get(url)
        participants.value = data ?? []
    } catch {
        participants.value = []
    } finally {
        participantsLoading.value = false
    }
}

const stats = computed(() => [
    {
        label: 'Total projets',
        value: statsData.value.total,
        bgClass: 'bg-primary/10',
        iconClass: 'w-6 h-6 text-primary',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
    },
    {
        label: 'En attente',
        value: statsData.value.pending,
        bgClass: 'bg-yellow-100',
        iconClass: 'w-6 h-6 text-yellow-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Ouverts',
        value: statsData.value.open,
        bgClass: 'bg-green-100',
        iconClass: 'w-6 h-6 text-green-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Rejetés',
        value: statsData.value.rejected,
        bgClass: 'bg-red-100',
        iconClass: 'w-6 h-6 text-red-500',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
])

function statusLabel(status) {
    return {
        pending: 'En attente',
        open: 'Ouvert',
        in_progress: 'En cours',
        completed: 'Terminé',
        rejected: 'Rejeté',
    }[status] ?? status
}

function statusClass(status) {
    return {
        pending: 'bg-yellow-100 text-yellow-700',
        open: 'bg-green-100 text-green-700',
        in_progress: 'bg-primary/10 text-primary',
        completed: 'bg-gray-100 text-gray-500',
        rejected: 'bg-red-100 text-red-600',
    }[status] ?? 'bg-gray-100 text-gray-500'
}

function openDetail(p) {
    participants.value = []
    participantsLoading.value = false
    documents.value = []
    docsLoading.value = false
    detailProject.value = p
    loadParticipants(`/projects/${p.id}/members`)
    loadDocuments(p.id)
}

function confirmDelete(p) {
    toDelete.value = p
}

function openRejectModal(p) {
    rejectModal.value = p
    rejectReason.value = ''
}

async function approveProject(p) {
    try {
        await api.patch(`/admin/project/${p.id}/approve`)
        await Promise.all([fetchProjects(), fetchStats()])
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de l\'approbation.')
    }
}

async function submitReject() {
    if (!rejectModal.value) return
    try {
        await api.patch(`/admin/project/${rejectModal.value.id}/reject`, { reason: rejectReason.value || null })
        rejectModal.value = null
        await Promise.all([fetchProjects(), fetchStats()])
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors du rejet.')
    }
}

async function deleteProject() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/project/${toDelete.value.id}`)
        projects.value = projects.value.filter(p => p.id !== toDelete.value.id)
        total.value = Math.max(0, total.value - 1)
        toDelete.value = null
        await Promise.all([fetchStats(), fetchDeletedProjects()])
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function fetchProjects() {
    loading.value = true
    error.value = ''
    try {
        const { data } = await api.get('/admin/projects', {
            params: { page: page.value, limit: 20, search: search.value, status: filterStatus.value }
        })
        projects.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (e) {
        error.value = 'Erreur lors du chargement des projets.'
    } finally {
        loading.value = false
    }
}

async function fetchStats() {
    try {
        const { data } = await api.get('/admin/projects/stats')
        statsData.value = data
    } catch {}
}

function openCreate() {
    projectForm.value = { title: '', description: '', location: '', startDate: '', endDate: '', capacity: null }
    photos.value = []
    formModal.value = { open: true, mode: 'create', id: null }
}

function openEdit(p) {
    projectForm.value = {
        title: p.title ?? '',
        description: p.description ?? '',
        location: p.location ?? '',
        startDate: p.start_date ? p.start_date.slice(0, 10) : '',
        endDate: p.end_date ? p.end_date.slice(0, 10) : '',
        capacity: p.capacity ?? null,
    }
    formModal.value = { open: true, mode: 'edit', id: p.id }
    detailProject.value = null
}

async function saveProject() {
    if (!projectForm.value.title.trim()) return
    saving.value = true
    try {
        const photoURLs = photos.value.length ? await uploadPhotos() : []
        const payload = {
            title: projectForm.value.title,
            description: projectForm.value.description || null,
            location: projectForm.value.location || null,
            start_date: projectForm.value.startDate || null,
            end_date: projectForm.value.endDate || null,
            capacity: projectForm.value.capacity || null,
            photo_urls: photoURLs,
        }
        if (formModal.value.mode === 'edit') {
            await api.put(`/admin/project/${formModal.value.id}`, payload)
        } else {
            await api.post('/admin/projects', payload)
            page.value = 1
        }
        formModal.value.open = false
        photos.value = []
        await Promise.all([fetchProjects(), fetchStats()])
    } catch {
        alert('Erreur lors de la création du projet.')
    } finally {
        saving.value = false
    }
}

async function fetchDeletedProjects() {
    try {
        const { data } = await api.get('/admin/project/deleted')
        deletedProjects.value = data ?? []
    } catch {
        deletedProjects.value = []
    }
}

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(() => { page.value = 1; fetchProjects() }, 300)
}

watch(page, fetchProjects)
watch(filterStatus, () => { page.value = 1; fetchProjects() })

onMounted(() => Promise.all([fetchProjects(), fetchStats(), fetchDeletedProjects()]))
</script>
