<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Formations</h1>
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

                <div class="px-6 py-4 border-b border-gray-100 flex flex-wrap items-center gap-3">
                    <div class="relative flex-1">
                        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 115 11a6 6 0 0112 0z"/>
                        </svg>
                        <input v-model="search" type="text" placeholder="Rechercher une formation…"
                            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                    </div>
                    <select v-model="filterLevel"
                        class="text-sm border border-gray-200 rounded-xl px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600">
                        <option value="">Tous les niveaux</option>
                        <option value="beginner">Débutant</option>
                        <option value="intermediate">Intermédiaire</option>
                        <option value="advanced">Avancé</option>
                    </select>
                    <select v-model="filterStatus"
                        class="text-sm border border-gray-200 rounded-xl px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/30 text-gray-600">
                        <option value="">Tous les statuts</option>
                        <option value="pending">En attente</option>
                        <option value="approved">Approuvées</option>
                        <option value="rejected">Rejetées</option>
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
                                <th class="text-left text-white font-medium px-5 py-3">Formation</th>
                                <th class="text-left text-white font-medium px-5 py-3">Date & lieu</th>
                                <th class="text-left text-white font-medium px-5 py-3">Niveau</th>
                                <th class="text-left text-white font-medium px-5 py-3">Inscrits</th>
                                <th class="text-left text-white font-medium px-5 py-3">Statut</th>
                                <th class="text-right text-white font-medium px-5 py-3">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(f, i) in formations" :key="f.id"
                                :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/50']">
                                <td class="px-5 py-3">
                                    <p class="font-medium text-gray-800 truncate max-w-52">{{ f.title }}</p>
                                    <p class="text-xs text-gray-500 truncate">{{ f.creator_name }}</p>
                                </td>
                                <td class="px-5 py-3">
                                    <p class="text-gray-800 text-xs font-medium">{{ f.date ? f.date.slice(0,10) : '-' }}</p>
                                    <p class="text-gray-500 text-xs truncate">{{ f.city || f.address || '-' }}</p>
                                </td>
                                <td class="px-5 py-3">
                                    <span :class="levelClass(f.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                        {{ levelLabel(f.level) }}
                                    </span>
                                </td>
                                <td class="px-5 py-3 text-xs text-gray-700 font-medium">
                                    {{ f.inscription_count }}<span v-if="f.capacity">/{{ f.capacity }}</span>
                                </td>
                                <td class="px-5 py-3">
                                    <span :class="statusClass(f.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                        {{ statusLabel(f.status) }}
                                    </span>
                                </td>
                                <td class="px-5 py-3">
                                    <div class="flex items-center justify-end gap-1">
                                        <button v-if="f.status === 'pending'" @click="approveFormation(f)"
                                            title="Approuver"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-secondary hover:bg-secondary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                            </svg>
                                        </button>
                                        <button v-if="f.status === 'pending'" @click="openRejectModal(f)"
                                            title="Rejeter"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                            </svg>
                                        </button>
                                        <button @click="openDetail(f)" title="Voir le détail"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                        </button>
                                        <button @click="confirmDelete(f)" title="Supprimer"
                                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                            <tr v-if="formations.length === 0">
                                <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
                                    Aucune formation ne correspond à vos filtres.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="px-6 py-3 border-t border-gray-100 text-xs text-gray-400">
                    {{ total }} formation(s) au total
                </div>
            </div>

            <Pagination v-if="total > 20" :page="page" :total="total" :limit="20" @update:page="changePage"/>

            <div class="mt-10">
                <h2 class="text-xl font-bold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Formations supprimées</h2>
                <p class="text-sm text-gray-400 mb-4">Historique des formations supprimées</p>
                <div v-if="deletedFormations.length === 0" class="bg-white rounded-xl shadow-sm p-8 text-center">
                    <p class="text-gray-400 text-sm">Aucune formation supprimée.</p>
                </div>
                <div v-else class="space-y-3">
                    <div
                        v-for="f in deletedFormations"
                        :key="f.id"
                        @click="selectedDeletedFormation = f"
                        class="bg-white rounded-xl shadow-sm px-5 py-4 flex items-center gap-4 cursor-pointer hover:shadow-md transition-shadow">
                        <div class="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center flex-shrink-0">
                            <svg class="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </div>
                        <div class="flex-1 min-w-0">
                            <h3 class="font-semibold text-gray-500 text-sm truncate">{{ f.title }}</h3>
                            <p class="text-xs text-gray-400 mt-0.5">Par {{ f.creator_name }} · Supprimée le {{ f.deleted_at?.slice(0, 10) }}</p>
                        </div>
                    </div>
                </div>
            </div>

        </template>

        <div v-if="detailFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="detailFormation = null">
            <div class="absolute inset-0 bg-black/40" @click="detailFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4 flex-shrink-0">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ detailFormation.title }}</h3>
                        <p class="text-xs text-gray-500 mt-0.5">Par {{ detailFormation.creator_name }}</p>
                    </div>
                    <button @click="detailFormation = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Statut</p>
                            <span :class="statusClass(detailFormation.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ statusLabel(detailFormation.status) }}
                            </span>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Niveau</p>
                            <span :class="levelClass(detailFormation.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                                {{ levelLabel(detailFormation.level) }}
                            </span>
                        </div>
                        <div v-if="detailFormation.date">
                            <p class="text-xs text-gray-400 mb-0.5">Date</p>
                            <p class="font-medium text-gray-800">{{ formatDateTime(detailFormation.date) }}</p>
                        </div>
                        <div v-if="detailFormation.address">
                            <p class="text-xs text-gray-400 mb-0.5">Adresse</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.address }}</p>
                        </div>
                        <div v-if="detailFormation.city">
                            <p class="text-xs text-gray-400 mb-0.5">Ville</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.postal ? detailFormation.postal + ' ' : '' }}{{ detailFormation.city }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Inscriptions</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.inscription_count }}<span v-if="detailFormation.capacity">/{{ detailFormation.capacity }}</span></p>
                        </div>
                        <div v-if="detailFormation.duration_hours">
                            <p class="text-xs text-gray-400 mb-0.5">Durée</p>
                            <p class="font-medium text-gray-800">{{ detailFormation.duration_hours }}h</p>
                        </div>
                    </div>
                    <div v-if="detailFormation.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ detailFormation.description }}</p>
                    </div>
                    <div v-if="detailFormation.objectives">
                        <p class="text-xs text-gray-400 mb-1">Objectifs pédagogiques</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ detailFormation.objectives }}</p>
                    </div>
                    <div v-if="detailFormation.prerequisites">
                        <p class="text-xs text-gray-400 mb-1">Prérequis</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ detailFormation.prerequisites }}</p>
                    </div>
                    <div v-if="detailFormation.syllabus">
                        <p class="text-xs text-gray-400 mb-1">Syllabus</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ detailFormation.syllabus }}</p>
                    </div>
                    <button
                        @click="downloadSyllabus(detailFormation)"
                        class="flex items-center gap-2 text-xs font-medium text-primary hover:text-primary-dark transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                        </svg>
                        Télécharger le syllabus
                    </button>
                    <div v-if="detailFormation.status === 'rejected' && detailFormation.rejection_reason">
                        <p class="text-xs text-gray-400 mb-1">Raison du rejet</p>
                        <p class="text-sm text-red-600">{{ detailFormation.rejection_reason }}</p>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Pièces jointes</p>
                        <p v-if="docsLoading" class="text-xs text-gray-400">Chargement…</p>
                        <p v-else-if="documents.length === 0" class="text-xs text-gray-400">Aucune pièce jointe.</p>
                        <div v-else class="flex gap-2 flex-wrap">
                            <template v-for="doc in documents" :key="doc.id">
                                <video v-if="isVideo(doc.link)" :src="doc.link" controls class="w-24 h-24 rounded-lg border border-gray-100 object-cover" />
                                <a v-else :href="doc.link" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border border-gray-100">
                                    <img :src="doc.link" class="w-full h-full object-cover" />
                                </a>
                            </template>
                        </div>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <div class="flex items-center justify-between mb-2">
                            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide">Étapes</p>
                            <button @click="showAddStep = !showAddStep" class="p-1 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
                            </button>
                        </div>
                        <div v-if="showAddStep" class="flex gap-2 mb-3">
                            <input v-model="newStep.title" placeholder="Titre de l'étape" @keyup.enter="addStep" class="flex-1 px-2 py-1.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                            <button @click="addStep" class="px-2 py-1.5 bg-primary text-white text-xs rounded-lg hover:bg-primary-dark transition-colors">Ajouter</button>
                        </div>
                        <p v-if="steps.length === 0" class="text-xs text-gray-400">Aucune étape renseignée.</p>
                        <div v-else class="space-y-1">
                            <div v-for="(s, i) in steps" :key="s.id" class="flex items-center gap-2 py-1.5 px-2 rounded-lg hover:bg-gray-50 group">
                                <template v-if="editingStepId === s.id">
                                    <div class="flex-1 flex flex-col gap-1">
                                        <input v-model="editStep.title" class="w-full px-2 py-1 text-xs border border-primary/40 rounded focus:outline-none" @keyup.enter="saveEditStep(s)"/>
                                        <input v-model="editStep.description" placeholder="Description (optionnelle)" class="w-full px-2 py-1 text-xs border border-gray-200 rounded focus:outline-none text-gray-500"/>
                                    </div>
                                    <div class="flex gap-0.5 flex-shrink-0">
                                        <button @click="saveEditStep(s)" class="p-1 text-secondary hover:bg-secondary/10 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg></button>
                                        <button @click="editingStepId = null" class="p-1 text-gray-400 hover:bg-gray-100 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg></button>
                                    </div>
                                </template>
                                <template v-else>
                                    <button
                                        @click="cycleStepStatus(s)"
                                        class="flex-shrink-0 w-5 h-5 rounded-full border-2 flex items-center justify-center transition-all"
                                        :class="{ 'border-gray-300 text-gray-300': s.status === 'todo', 'border-primary text-primary': s.status === 'in_progress', 'border-secondary bg-secondary text-white': s.status === 'done' }"
                                    >
                                        <svg v-if="s.status === 'done'" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg>
                                        <span v-else-if="s.status === 'in_progress'" class="w-1.5 h-1.5 rounded-full bg-primary"></span>
                                    </button>
                                    <div class="flex-1 min-w-0">
                                        <p class="text-sm text-gray-700 truncate" :class="{ 'line-through text-gray-400': s.status === 'done' }">{{ s.title }}</p>
                                        <p v-if="s.description" class="text-xs text-gray-500 truncate">{{ s.description }}</p>
                                    </div>
                                    <span class="px-1.5 py-0.5 rounded text-xs flex-shrink-0" :class="{ 'bg-gray-100 text-gray-500': s.status === 'todo', 'bg-primary/10 text-primary': s.status === 'in_progress', 'bg-secondary/10 text-secondary': s.status === 'done' }">{{ stepStatusLabel(s.status) }}</span>
                                    <div class="flex gap-0.5 opacity-0 group-hover:opacity-100 flex-shrink-0 transition-opacity">
                                        <button @click="moveStep(i, -1)" :disabled="i === 0" class="p-1 text-gray-300 hover:text-primary hover:bg-primary/10 rounded transition-colors disabled:opacity-30 disabled:cursor-default"><svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M5 15l7-7 7 7"/></svg></button>
                                        <button @click="moveStep(i, 1)" :disabled="i === steps.length - 1" class="p-1 text-gray-300 hover:text-primary hover:bg-primary/10 rounded transition-colors disabled:opacity-30 disabled:cursor-default"><svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/></svg></button>
                                        <button @click="startEditStep(s)" class="p-1 text-gray-300 hover:text-primary hover:bg-primary/10 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg></button>
                                        <button @click="deleteStep(s.id)" class="p-1 text-gray-300 hover:text-red-500 hover:bg-red-50 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg></button>
                                    </div>
                                </template>
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
                        <button @click="openEdit(detailFormation)"
                            class="px-3 py-1.5 text-sm font-medium text-primary border border-primary/30 rounded-lg hover:bg-primary/5 transition-colors">
                            Modifier
                        </button>
                        <button v-if="detailFormation.status === 'pending'" @click="approveFormation(detailFormation); detailFormation = null"
                            class="px-3 py-1.5 text-sm font-medium text-secondary border border-secondary/30 rounded-lg hover:bg-secondary/5 transition-colors">
                            Approuver
                        </button>
                        <button v-if="detailFormation.status === 'pending'" @click="openRejectModal(detailFormation); detailFormation = null"
                            class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                            Rejeter
                        </button>
                    </div>
                    <button @click="confirmDelete(detailFormation); detailFormation = null"
                        class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
                        Supprimer
                    </button>
                </div>
            </div>
        </div>

        <div v-if="rejectModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="rejectModal = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Rejeter la formation</h3>
                <p class="text-sm text-gray-500 mb-4">« {{ rejectModal.title }} »</p>
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Raison (optionnelle)</label>
                    <textarea v-model="rejectReason" rows="3"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                        placeholder="Expliquez pourquoi la formation est rejetée…"/>
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
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer la formation ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimée avec toutes ses inscriptions.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50">Annuler</button>
                    <button @click="deleteFormation" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="formModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal = false"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ formMode === 'edit' ? 'Modifier la formation' : 'Créer une formation' }}
                    </h3>
                    <button @click="formModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4 overflow-y-auto flex-1">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Titre <span class="text-red-400">*</span></label>
                        <input v-model="form.title" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Titre de la formation"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea v-model="form.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Description de la formation…"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Objectifs pédagogiques</label>
                        <textarea v-model="form.objectives" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Ce que les participants sauront faire à l'issue de la formation…"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Prérequis</label>
                        <textarea v-model="form.prerequisites" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Connaissances ou matériel nécessaires avant de participer…"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Syllabus (programme détaillé)</label>
                        <textarea v-model="form.syllabus" rows="4"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Déroulé détaillé de la formation, étape par étape…"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Adresse <span class="text-red-400">*</span></label>
                        <input v-model="form.address" type="text"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                            placeholder="Adresse ou salle"/>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Niveau</label>
                            <select v-model="form.level"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                                <option value="beginner">Débutant</option>
                                <option value="intermediate">Intermédiaire</option>
                                <option value="advanced">Avancé</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Durée (heures)</label>
                            <input v-model.number="form.duration_hours" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Ex: 3"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Date & heure <span class="text-red-400">*</span></label>
                            <input v-model="form.date" type="datetime-local" :min="formMode === 'create' ? minDateTime : undefined" required
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Ville</label>
                            <input v-model="form.city" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Paris"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Code postal</label>
                            <input v-model="form.postal" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="75001"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Capacité</label>
                            <input v-model.number="form.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Prix (€) <span class="text-gray-300">(laisser vide si gratuit)</span></label>
                            <input v-model.number="form.price" type="number" min="0" step="0.01"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Ex: 25.00"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Thématique</label>
                            <select v-model="form.id_category"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary">
                                <option value="">Aucune</option>
                                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                            </select>
                        </div>
                    </div>
                    <div v-if="formMode === 'create'">
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
                    <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2 flex-shrink-0">
                    <button @click="formModal = false" class="px-4 py-1.5 text-sm text-gray-600 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button @click="submitForm" :disabled="submitting"
                        class="px-4 py-1.5 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ submitting ? 'Enregistrement…' : (formMode === 'edit' ? 'Enregistrer' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="selectedDeletedFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="selectedDeletedFormation = null">
            <div class="absolute inset-0 bg-black/40" @click="selectedDeletedFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[85vh] flex flex-col overflow-hidden">
                <div class="px-6 py-4 border-b border-gray-100 flex items-start justify-between gap-4 flex-shrink-0">
                    <div class="min-w-0">
                        <h3 class="font-semibold text-gray-600" style="font-family: var(--font-family-title)">{{ selectedDeletedFormation.title }}</h3>
                        <p class="text-xs text-gray-400 mt-0.5">Par {{ selectedDeletedFormation.creator_name }} · Supprimée le {{ selectedDeletedFormation.deleted_at?.slice(0, 10) }}</p>
                    </div>
                    <span class="flex-shrink-0 px-2 py-1 rounded-full bg-gray-100 text-gray-400 text-xs font-medium">Supprimée</span>
                    <button @click="selectedDeletedFormation = null" class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p class="text-xs text-gray-400 mb-0.5">Niveau</p>
                            <span :class="levelClass(selectedDeletedFormation.level)" class="px-2 py-0.5 rounded-full text-xs font-medium">{{ levelLabel(selectedDeletedFormation.level) }}</span>
                        </div>
                        <div v-if="selectedDeletedFormation.date">
                            <p class="text-xs text-gray-400 mb-0.5">Date</p>
                            <p class="font-medium text-gray-800">{{ formatDateTime(selectedDeletedFormation.date) }}</p>
                        </div>
                        <div v-if="selectedDeletedFormation.address">
                            <p class="text-xs text-gray-400 mb-0.5">Adresse</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedFormation.address }}</p>
                        </div>
                        <div v-if="selectedDeletedFormation.city">
                            <p class="text-xs text-gray-400 mb-0.5">Ville</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedFormation.postal ? selectedDeletedFormation.postal + ' ' : '' }}{{ selectedDeletedFormation.city }}</p>
                        </div>
                        <div v-if="selectedDeletedFormation.duration_hours">
                            <p class="text-xs text-gray-400 mb-0.5">Durée</p>
                            <p class="font-medium text-gray-800">{{ selectedDeletedFormation.duration_hours }}h</p>
                        </div>
                    </div>
                    <div v-if="selectedDeletedFormation.description">
                        <p class="text-xs text-gray-400 mb-1">Description</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedDeletedFormation.description }}</p>
                    </div>
                    <div v-if="selectedDeletedFormation.objectives">
                        <p class="text-xs text-gray-400 mb-1">Objectifs pédagogiques</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedDeletedFormation.objectives }}</p>
                    </div>
                    <div v-if="selectedDeletedFormation.prerequisites">
                        <p class="text-xs text-gray-400 mb-1">Prérequis</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedDeletedFormation.prerequisites }}</p>
                    </div>
                    <div v-if="selectedDeletedFormation.syllabus">
                        <p class="text-xs text-gray-400 mb-1">Syllabus</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedDeletedFormation.syllabus }}</p>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Pièces jointes</p>
                        <p v-if="!selectedDeletedFormation.documents?.length" class="text-xs text-gray-400">Aucune pièce jointe.</p>
                        <div v-else class="flex gap-2 flex-wrap">
                            <template v-for="doc in selectedDeletedFormation.documents" :key="doc.id">
                                <video v-if="isVideo(doc.link)" :src="doc.link" controls class="w-24 h-24 rounded-lg border border-gray-100 object-cover" />
                                <a v-else :href="doc.link" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border border-gray-100">
                                    <img :src="doc.link" class="w-full h-full object-cover" />
                                </a>
                            </template>
                        </div>
                    </div>

                    <div class="border-t border-gray-100 pt-4">
                        <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Étapes</p>
                        <p v-if="!selectedDeletedFormation.steps?.length" class="text-xs text-gray-400">Aucune étape renseignée.</p>
                        <div v-else class="space-y-1">
                            <div v-for="s in selectedDeletedFormation.steps" :key="s.id" class="flex items-center gap-2 py-1.5 px-2 rounded-lg">
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
            </div>
        </div>

    </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const formations = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const error = ref('')
const search = ref('')
const filterStatus = ref('')
const filterLevel = ref('')
const detailFormation = ref(null)
const rejectModal = ref(null)
const rejectReason = ref('')
const toDelete = ref(null)
const deleting = ref(false)

const categories = ref([])
const formModal = ref(false)
const formMode = ref('create')
const editId = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', date: '', address: '', city: '', postal: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '', price: null, objectives: '', prerequisites: '', syllabus: '' })

const minDateTime = computed(() => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 1)
    return now.toISOString().slice(0, 16)
})

const statCounts = ref({ total: 0, pending: 0, approved: 0, rejected: 0 })
const participants = ref([])
const participantsLoading = ref(false)
const steps = ref([])
const showAddStep = ref(false)
const newStep = ref({ title: '' })
const editingStepId = ref(null)
const editStep = ref({ title: '', description: '' })
const documents = ref([])
const docsLoading = ref(false)
const deletedFormations = ref([])
const selectedDeletedFormation = ref(null)

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

async function loadSteps(id) {
    try {
        const { data } = await api.get(`/formations/${id}/steps`)
        steps.value = data ?? []
    } catch {
        steps.value = []
    }
}

function stepStatusLabel(s) {
    return { todo: 'À faire', in_progress: 'En cours', done: 'Fait' }[s] ?? s
}

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
        const { data } = await api.get(`/formations/${id}/documents`)
        documents.value = data ?? []
    } catch {
        documents.value = []
    } finally {
        docsLoading.value = false
    }
}

const stats = computed(() => [
    {
        label: 'Total formations',
        value: statCounts.value.total,
        bgClass: 'bg-primary/10',
        iconClass: 'w-6 h-6 text-primary',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>`,
    },
    {
        label: 'En attente',
        value: statCounts.value.pending,
        bgClass: 'bg-yellow-100',
        iconClass: 'w-6 h-6 text-yellow-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Approuvées',
        value: statCounts.value.approved,
        bgClass: 'bg-green-100',
        iconClass: 'w-6 h-6 text-green-600',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
    {
        label: 'Rejetées',
        value: statCounts.value.rejected,
        bgClass: 'bg-red-100',
        iconClass: 'w-6 h-6 text-red-500',
        icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    },
])

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

function openDetail(f) {
    participants.value = []
    participantsLoading.value = false
    steps.value = []
    showAddStep.value = false
    editingStepId.value = null
    newStep.value = { title: '' }
    documents.value = []
    docsLoading.value = false
    detailFormation.value = f
    loadParticipants(`/formations/${f.id}/participants`)
    loadSteps(f.id)
    loadDocuments(f.id)
}

async function addStep() {
    if (!newStep.value.title.trim()) return
    try {
        const { data } = await api.post(`/formations/${detailFormation.value.id}/steps`, {
            title: newStep.value.title,
            step_order: steps.value.length,
        })
        steps.value.push(data)
        newStep.value = { title: '' }
        showAddStep.value = false
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function cycleStepStatus(s) {
    const next = { todo: 'in_progress', in_progress: 'done', done: 'todo' }[s.status]
    try {
        await api.patch(`/formations/${detailFormation.value.id}/steps/${s.id}/status`, { status: next })
        s.status = next
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function deleteStep(id) {
    try {
        await api.delete(`/formations/${detailFormation.value.id}/steps/${id}`)
        steps.value = steps.value.filter(s => s.id !== id)
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

function startEditStep(s) {
    editingStepId.value = s.id
    editStep.value = { title: s.title, description: s.description ?? '' }
}

async function saveEditStep(s) {
    if (!editStep.value.title.trim()) return
    try {
        await api.patch(`/formations/${detailFormation.value.id}/steps/${s.id}`, {
            title: editStep.value.title,
            description: editStep.value.description || null,
            step_order: s.step_order,
        })
        const idx = steps.value.findIndex(x => x.id === s.id)
        if (idx !== -1) steps.value[idx] = { ...steps.value[idx], title: editStep.value.title, description: editStep.value.description || null }
        editingStepId.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function moveStep(index, dir) {
    const target = index + dir
    if (target < 0 || target >= steps.value.length) return
    const a = steps.value[index]
    const b = steps.value[target]
    try {
        await Promise.all([
            api.patch(`/formations/${detailFormation.value.id}/steps/${a.id}`, { title: a.title, description: a.description ?? null, step_order: b.step_order }),
            api.patch(`/formations/${detailFormation.value.id}/steps/${b.id}`, { title: b.title, description: b.description ?? null, step_order: a.step_order }),
        ])
        const newSteps = [...steps.value]
        newSteps[target] = { ...a, step_order: b.step_order }
        newSteps[index] = { ...b, step_order: a.step_order }
        steps.value = newSteps
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

function confirmDelete(f) {
    toDelete.value = f
}

function openCreate() {
    form.value = { title: '', description: '', date: '', address: '', city: '', postal: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '', price: null, objectives: '', prerequisites: '', syllabus: '' }
    formError.value = ''
    formMode.value = 'create'
    editId.value = null
    photos.value = []
    formModal.value = true
}

function openEdit(f) {
    form.value = {
        title: f.title ?? '',
        description: f.description ?? '',
        date: f.date ? f.date.slice(0, 16).replace(' ', 'T') : '',
        address: f.address ?? '',
        city: f.city ?? '',
        postal: f.postal ?? '',
        capacity: f.capacity ?? null,
        level: f.level ?? 'beginner',
        duration_hours: f.duration_hours ?? null,
        id_category: f.id_category ?? '',
        price: f.price ?? null,
        objectives: f.objectives ?? '',
        prerequisites: f.prerequisites ?? '',
        syllabus: f.syllabus ?? '',
    }
    formError.value = ''
    formMode.value = 'edit'
    editId.value = f.id
    detailFormation.value = null
    formModal.value = true
}

async function submitForm() {
    if (!form.value.title.trim()) { formError.value = 'Le titre est requis.'; return }
    if (!form.value.date) { formError.value = 'La date est requise.'; return }
    if (!form.value.address?.trim()) { formError.value = 'L\'adresse est requise.'; return }
    submitting.value = true
    formError.value = ''
    try {
        const photoURLs = photos.value.length ? await uploadPhotos() : []
        const payload = {
            title: form.value.title,
            description: form.value.description || null,
            date: form.value.date || null,
            address: form.value.address || null,
            city: form.value.city || null,
            postal: form.value.postal || null,
            photo_urls: photoURLs,
            capacity: form.value.capacity || null,
            level: form.value.level,
            duration_hours: form.value.duration_hours || null,
            id_category: form.value.id_category || null,
            price: form.value.price || null,
            objectives: form.value.objectives || null,
            prerequisites: form.value.prerequisites || null,
            syllabus: form.value.syllabus || null,
        }
        if (formMode.value === 'edit') {
            await api.put(`/admin/formation/${editId.value}`, payload)
        } else {
            await api.post('/admin/formations', payload)
        }
        formModal.value = false
        photos.value = []
        await Promise.all([fetchFormations(), fetchStats()])
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

async function downloadSyllabus(f) {
    try {
        const { data } = await api.get(`/formations/${f.id}/syllabus-pdf`, { responseType: 'blob' })
        const url = URL.createObjectURL(data)
        const link = document.createElement('a')
        link.href = url
        link.download = `syllabus-${f.id}.pdf`
        link.click()
        URL.revokeObjectURL(url)
    } catch {
        alert('Syllabus non disponible.')
    }
}

function openRejectModal(f) {
    rejectModal.value = f
    rejectReason.value = ''
}

async function approveFormation(f) {
    try {
        await api.patch(`/formations/${f.id}/approve`)
        await Promise.all([fetchFormations(), fetchStats()])
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de l\'approbation.')
    }
}

async function submitReject() {
    if (!rejectModal.value) return
    try {
        await api.patch(`/formations/${rejectModal.value.id}/reject`, { reason: rejectReason.value || null })
        rejectModal.value = null
        await Promise.all([fetchFormations(), fetchStats()])
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors du rejet.')
    }
}

async function deleteFormation() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/admin/formation/${toDelete.value.id}`)
        formations.value = formations.value.filter(f => f.id !== toDelete.value.id)
        total.value = Math.max(0, total.value - 1)
        toDelete.value = null
        await Promise.all([fetchStats(), fetchDeletedFormations()])
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

function changePage(p) {
    page.value = p
}

async function fetchFormations() {
    loading.value = true
    error.value = ''
    try {
        const { data } = await api.get('/admin/formations', {
            params: { page: page.value, limit: 20, search: search.value, status: filterStatus.value, level: filterLevel.value }
        })
        formations.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (e) {
        error.value = 'Erreur lors du chargement des formations.'
    } finally {
        loading.value = false
    }
}

async function fetchStats() {
    try {
        const [all, pending, approved, rejected] = await Promise.all([
            api.get('/admin/formations', { params: { limit: 1 } }),
            api.get('/admin/formations', { params: { limit: 1, status: 'pending' } }),
            api.get('/admin/formations', { params: { limit: 1, status: 'approved' } }),
            api.get('/admin/formations', { params: { limit: 1, status: 'rejected' } }),
        ])
        statCounts.value = {
            total:    all.data.total ?? 0,
            pending:  pending.data.total ?? 0,
            approved: approved.data.total ?? 0,
            rejected: rejected.data.total ?? 0,
        }
    } catch {}
}

async function fetchDeletedFormations() {
    try {
        const { data } = await api.get('/admin/formation/deleted')
        deletedFormations.value = data ?? []
    } catch {
        deletedFormations.value = []
    }
}

watch([page, filterStatus, filterLevel], fetchFormations)
watch(search, () => { page.value = 1; fetchFormations() })

onMounted(() => { fetchFormations(); fetchStats(); fetchCategories(); fetchDeletedFormations() })
</script>
