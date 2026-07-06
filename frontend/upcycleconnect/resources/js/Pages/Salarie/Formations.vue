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

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex flex-wrap items-center gap-3">
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
                <select v-model="filterStatus" @change="resetAndFetch" class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Tous les statuts</option>
                    <option value="pending">En attente</option>
                    <option value="approved">Approuvée</option>
                    <option value="rejected">Rejetée</option>
                </select>
                <select v-model="filterLevel" @change="resetAndFetch" class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Tous les niveaux</option>
                    <option value="beginner">Débutant</option>
                    <option value="intermediate">Intermédiaire</option>
                    <option value="advanced">Avancé</option>
                </select>
                <select v-if="categories.length" v-model="filterCategory" @change="resetAndFetch" class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Toutes les thématiques</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} formation(s)</span>
            </div>

            <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">Chargement…</div>

            <div v-else-if="formations.length === 0" class="py-12 text-center text-gray-400 text-sm">
                {{ filterStatus || filterLevel || filterCategory || search ? 'Aucune formation pour ce filtre.' : 'Vous n\'avez créé aucune formation.' }}
            </div>

            <div v-else class="grid grid-cols-3 gap-5 p-6">
                <div
                    v-for="f in formations"
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

            <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />
        </div>

        <div v-if="selectedFormation" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="selectedFormation = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl flex flex-col max-h-[90vh]">
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
                    <p v-if="selectedFormation.description" class="leading-relaxed text-gray-600 whitespace-pre-line break-words">{{ selectedFormation.description }}</p>
                    <div class="grid grid-cols-2 gap-3 text-xs text-gray-500">
                        <div v-if="selectedFormation.date"><span class="font-medium text-gray-700">Date :</span> {{ formatDateTime(selectedFormation.date) }}</div>
                        <div v-if="selectedFormation.address"><span class="font-medium text-gray-700">Adresse :</span> {{ selectedFormation.address }}</div>
                        <div v-if="selectedFormation.city"><span class="font-medium text-gray-700">Ville :</span> {{ selectedFormation.postal ? selectedFormation.postal + ' ' : '' }}{{ selectedFormation.city }}</div>
                        <div v-if="selectedFormation.duration_hours"><span class="font-medium text-gray-700">Durée :</span> {{ selectedFormation.duration_hours }}h</div>
                        <div v-if="selectedFormation.capacity"><span class="font-medium text-gray-700">Capacité :</span> {{ selectedFormation.capacity }} places</div>
                        <div><span class="font-medium text-gray-700">Inscriptions :</span> {{ selectedFormation.inscription_count ?? 0 }}<span v-if="selectedFormation.capacity"> / {{ selectedFormation.capacity }}</span></div>
                        <div v-if="selectedFormation.category_name"><span class="font-medium text-gray-700">Thématique :</span> {{ selectedFormation.category_name }}</div>
                    </div>
                    <div v-if="selectedFormation.objectives">
                        <p class="text-xs text-gray-400 mb-1">Objectifs pédagogiques</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedFormation.objectives }}</p>
                    </div>
                    <div v-if="selectedFormation.prerequisites">
                        <p class="text-xs text-gray-400 mb-1">Prérequis</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedFormation.prerequisites }}</p>
                    </div>
                    <div v-if="selectedFormation.syllabus">
                        <p class="text-xs text-gray-400 mb-1">Syllabus</p>
                        <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-line break-words">{{ selectedFormation.syllabus }}</p>
                    </div>
                    <button
                        @click="downloadSyllabus(selectedFormation)"
                        class="flex items-center gap-2 text-xs font-medium text-primary hover:text-primary-dark transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                        </svg>
                        Télécharger le syllabus
                    </button>

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
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ editTarget ? 'Modifier la formation' : 'Créer une formation' }}
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
                            <input v-model="form.date" type="datetime-local" :min="minDateTime" required
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
                    <div>
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
                        {{ submitting ? 'Enregistrement…' : (editTarget ? 'Enregistrer' : 'Créer') }}
                    </button>
                </div>
            </div>
        </div>

    </SalarieLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SalarieLayout from '@/Layouts/SalarieLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const route = useRoute()

const formations = ref([])
const categories = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(true)
const search = ref('')
const filterStatus = ref('')
const filterLevel = ref('')
const filterCategory = ref('')
const toDelete = ref(null)
const deleting = ref(false)
const selectedFormation = ref(null)
const participants = ref([])
const participantsLoading = ref(false)
const steps = ref([])
const showAddStep = ref(false)
const newStep = ref({ title: '' })
const editingStepId = ref(null)
const editStep = ref({ title: '', description: '' })
const photos = ref([])
const documents = ref([])
const docsLoading = ref(false)

watch(selectedFormation, async (val) => {
    steps.value = []
    showAddStep.value = false
    editingStepId.value = null
    documents.value = []
    docsLoading.value = false
    participants.value = []
    participantsLoading.value = false
    if (!val) return
    loadParticipants(`/formations/${val.id}/participants`)
    loadDocuments(val.id)
    try {
        const { data } = await api.get(`/formations/${val.id}/steps`)
        steps.value = data ?? []
    } catch {}
})

function isVideo(link) {
    return /\.(mp4|webm|ogg|mov)(\?|$)/i.test(link)
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

async function deleteDocument(doc) {
    try {
        await api.delete(`/documents/${doc.id}`)
        documents.value = documents.value.filter(d => d.id !== doc.id)
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

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

function stepStatusLabel(s) {
    return { todo: 'À faire', in_progress: 'En cours', done: 'Fait' }[s] ?? s
}

async function addStep() {
    if (!newStep.value.title.trim()) return
    try {
        const { data } = await api.post(`/formations/${selectedFormation.value.id}/steps`, {
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
        await api.patch(`/formations/${selectedFormation.value.id}/steps/${s.id}/status`, { status: next })
        s.status = next
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function deleteStep(id) {
    try {
        await api.delete(`/formations/${selectedFormation.value.id}/steps/${id}`)
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
        await api.patch(`/formations/${selectedFormation.value.id}/steps/${s.id}`, {
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
            api.patch(`/formations/${selectedFormation.value.id}/steps/${a.id}`, { title: a.title, description: a.description ?? null, step_order: b.step_order }),
            api.patch(`/formations/${selectedFormation.value.id}/steps/${b.id}`, { title: b.title, description: b.description ?? null, step_order: a.step_order }),
        ])
        const newSteps = [...steps.value]
        newSteps[target] = { ...a, step_order: b.step_order }
        newSteps[index] = { ...b, step_order: a.step_order }
        steps.value = newSteps
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

const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', description: '', date: '', address: '', city: '', postal: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '', price: null, objectives: '', prerequisites: '', syllabus: '' })

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

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
    fetchFormations()
}

function levelLabel(level) {
    return { beginner: 'Débutant', intermediate: 'Intermédiaire', advanced: 'Avancé' }[level] ?? level
}

function formatDateTime(dateStr) {
    if (!dateStr) return ''
    const d = new Date(dateStr)
    return d.toLocaleDateString('fr-FR') + ' à ' + d.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
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
    form.value = { title: '', description: '', date: '', address: '', city: '', postal: '', capacity: null, level: 'beginner', duration_hours: null, id_category: '', price: null, objectives: '', prerequisites: '', syllabus: '' }
    photos.value = []
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
        address: f.address ?? '',
        city: f.city ?? '',
        postal: f.postal ?? '',
        capacity: f.capacity ?? null,
        level: f.level,
        duration_hours: f.duration_hours ?? null,
        id_category: f.id_category ?? '',
        price: f.price ?? null,
        objectives: f.objectives ?? '',
        prerequisites: f.prerequisites ?? '',
        syllabus: f.syllabus ?? '',
    }
    photos.value = []
    formError.value = ''
    formModal.value = true
}

async function fetchFormations() {
    loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        if (filterStatus.value) params.status = filterStatus.value
        if (filterLevel.value) params.level = filterLevel.value
        if (filterCategory.value) params.id_category = filterCategory.value
        const { data } = await api.get('/user/my-formations', { params })
        formations.value = data.data ?? []
        total.value = data.total ?? 0
    } catch {
        formations.value = []
        total.value = 0
    } finally {
        loading.value = false
    }
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

async function deleteFormation() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/user/formation/${toDelete.value.id}`)
        toDelete.value = null
        await fetchFormations()
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
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
            capacity: form.value.capacity || null,
            level: form.value.level,
            duration_hours: form.value.duration_hours || null,
            id_category: form.value.id_category || null,
            price: form.value.price || null,
            objectives: form.value.objectives || null,
            prerequisites: form.value.prerequisites || null,
            syllabus: form.value.syllabus || null,
            photo_urls: photoURLs,
        }
        if (editTarget.value) {
            await api.patch(`/formations/${editTarget.value.id}`, payload)
        } else {
            await api.post('/formations', payload)
        }
        formModal.value = false
        photos.value = []
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
