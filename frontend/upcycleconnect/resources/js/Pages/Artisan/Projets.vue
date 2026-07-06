<template>
    <ArtisanLayout>

        <div class="mb-6 flex items-center justify-between">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Mes projets</h1>
            <button
                @click="openCreate"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un projet
            </button>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">Chargement…</div>

        <template v-else>
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
                        placeholder="Rechercher un projet…"
                        class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                </div>
                <select v-model="filterStatus" @change="resetAndFetch" class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30">
                    <option value="">Tous les statuts</option>
                    <option value="pending">En attente</option>
                    <option value="rejected">Rejeté</option>
                    <option value="open">Ouvert</option>
                    <option value="in_progress">En cours</option>
                    <option value="completed">Terminé</option>
                </select>
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ total }} projet(s)</span>
            </div>

            <div v-if="projects.length === 0" class="py-12 text-center text-gray-400 text-sm">
                {{ filterStatus || search ? 'Aucun projet pour ce filtre.' : 'Vous n\'avez créé aucun projet.' }}
            </div>

            <div v-else class="grid grid-cols-3 gap-5 p-6">
                <div
                    v-for="p in projects"
                    :key="p.id"
                    @click="selectedProject = p"
                    class="bg-white rounded-2xl shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow cursor-pointer border border-gray-100"
                >
                    <div class="flex items-start justify-between gap-2">
                        <h3 class="font-semibold text-gray-800 text-sm leading-snug">{{ p.title }}</h3>
                        <span :class="statusClass(p.status)" class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0">
                            {{ statusLabel(p.status) }}
                        </span>
                    </div>
                    <p v-if="p.description" class="text-xs text-gray-500 line-clamp-2">{{ p.description }}</p>
                    <p v-if="p.status === 'rejected' && p.rejection_reason" class="text-xs text-red-500 line-clamp-1">{{ p.rejection_reason }}</p>
                    <div class="flex items-center gap-3 mt-auto pt-2 border-t border-gray-50">
                        <div class="flex items-center gap-2 flex-1">
                            <button
                                @click.stop="vote(p, 1)"
                                :class="p.user_vote === 1 ? 'text-secondary font-semibold' : 'text-gray-400 hover:text-secondary'"
                                class="flex items-center gap-0.5 text-xs transition-colors"
                                title="J'aime"
                            >
                                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"/>
                                </svg>
                                {{ p.up_votes ?? 0 }}
                            </button>
                            <button
                                @click.stop="vote(p, -1)"
                                :class="p.user_vote === -1 ? 'text-red-500 font-semibold' : 'text-gray-400 hover:text-red-400'"
                                class="flex items-center gap-0.5 text-xs transition-colors"
                                title="Je n'aime pas"
                            >
                                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.095c.5 0 .905-.405.905-.905 0-.714.211-1.412.608-2.006L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5"/>
                                </svg>
                                {{ p.down_votes ?? 0 }}
                            </button>
                        </div>
                        <span class="text-xs text-gray-400">{{ p.members_count }} membre(s)</span>
                        <button
                            @click.stop="toDelete = p"
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

            <Pagination v-if="total > 15" :page="page" :total="total" :limit="15" @update:page="changePage" />
        </div>
        </template>

        <div v-if="selectedProject" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="selectedProject = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl flex flex-col max-h-[90vh]">
                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{ selectedProject.title }}</h3>
                    <button @click="selectedProject = null" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-3 text-sm text-gray-700 overflow-y-auto flex-1">
                    <div>
                        <span :class="statusClass(selectedProject.status)" class="px-2 py-0.5 rounded-full text-xs font-medium">
                            {{ statusLabel(selectedProject.status) }}
                        </span>
                    </div>
                    <p v-if="selectedProject.status === 'rejected' && selectedProject.rejection_reason" class="text-xs text-red-500">
                        Raison : {{ selectedProject.rejection_reason }}
                    </p>
                    <p v-if="selectedProject.description" class="leading-relaxed text-gray-600">{{ selectedProject.description }}</p>
                    <div class="grid grid-cols-2 gap-3 text-xs text-gray-500">
                        <div v-if="selectedProject.start_date">
                            <span class="font-medium text-gray-700">Début :</span>
                            {{ selectedProject.start_date.slice(0, 10) }}
                        </div>
                        <div v-if="selectedProject.end_date">
                            <span class="font-medium text-gray-700">Fin :</span>
                            {{ selectedProject.end_date.slice(0, 10) }}
                        </div>
                        <div v-if="selectedProject.location"><span class="font-medium text-gray-700">Lieu :</span> {{ selectedProject.location }}</div>
                        <div v-if="selectedProject.capacity"><span class="font-medium text-gray-700">Capacité :</span> {{ selectedProject.capacity }} membres</div>
                        <div><span class="font-medium text-gray-700">Membres :</span> {{ selectedProject.members_count }}</div>
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
                        <div class="flex items-center justify-between mb-2">
                            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide">Matériaux</p>
                            <button @click="showAddMaterial = !showAddMaterial" class="p-1 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
                            </button>
                        </div>
                        <div v-if="showAddMaterial" class="flex gap-2 mb-3">
                            <input v-model="newMaterial.name" placeholder="Nom" class="flex-1 px-2 py-1.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                            <input v-model.number="newMaterial.quantity" type="number" min="0" step="0.01" placeholder="Qté" class="w-16 px-2 py-1.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30"/>
                            <select v-model="newMaterial.unit" class="w-24 px-2 py-1.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 bg-white">
                              <option value="">Unité</option>
                              <option>pièce</option><option>kg</option><option>g</option><option>m</option><option>cm</option><option>L</option><option>mL</option><option>paire</option><option>lot</option>
                            </select>
                            <button @click="addMaterial" class="px-2 py-1.5 bg-primary text-white text-xs rounded-lg hover:bg-primary-dark transition-colors">Ajouter</button>
                        </div>
                        <p v-if="materials.length === 0" class="text-xs text-gray-400">Aucun matériau renseigné.</p>
                        <table v-else class="w-full text-xs">
                            <thead>
                                <tr class="text-gray-400 border-b border-gray-100">
                                    <th class="text-left pb-1 font-medium">Nom</th>
                                    <th class="text-right pb-1 font-medium">Qté</th>
                                    <th class="text-left pb-1 font-medium pl-2">Unité</th>
                                    <th class="w-16"></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="m in materials" :key="m.id" class="border-b border-gray-50 last:border-0">
                                    <template v-if="editingMaterialId === m.id">
                                        <td class="py-1.5 pr-1"><input v-model="editMaterial.name" class="w-full px-2 py-1 text-xs border border-primary/40 rounded focus:outline-none"/></td>
                                        <td class="py-1.5 px-1 w-16"><input v-model.number="editMaterial.quantity" type="number" min="0" step="0.01" class="w-full px-2 py-1 text-xs border border-primary/40 rounded focus:outline-none text-right"/></td>
                                        <td class="py-1.5 px-1 w-20"><select v-model="editMaterial.unit" class="w-full px-2 py-1 text-xs border border-primary/40 rounded focus:outline-none bg-white">
                                          <option value="">-</option>
                                          <option>pièce</option><option>kg</option><option>g</option><option>m</option><option>cm</option><option>L</option><option>mL</option><option>paire</option><option>lot</option>
                                        </select></td>
                                        <td class="py-1.5 pl-1"><div class="flex gap-0.5">
                                            <button @click="saveEditMaterial(m)" class="p-1 text-secondary hover:bg-secondary/10 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg></button>
                                            <button @click="editingMaterialId = null" class="p-1 text-gray-400 hover:bg-gray-100 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg></button>
                                        </div></td>
                                    </template>
                                    <template v-else>
                                        <td class="py-1.5 pr-2 text-gray-700">{{ m.name }}</td>
                                        <td class="py-1.5 text-right text-gray-600">{{ m.quantity }}</td>
                                        <td class="py-1.5 pl-2 text-gray-500">{{ m.unit ?? '-' }}</td>
                                        <td class="py-1.5 pl-1"><div class="flex gap-0.5">
                                            <button @click="startEditMaterial(m)" class="p-1 text-gray-300 hover:text-primary hover:bg-primary/10 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg></button>
                                            <button @click="deleteMaterial(m.id)" class="p-1 text-gray-300 hover:text-red-500 hover:bg-red-50 rounded transition-colors"><svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg></button>
                                        </div></td>
                                    </template>
                                </tr>
                            </tbody>
                        </table>
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
                <div v-if="selectedProject.status === 'pending' || selectedProject.status === 'rejected'" class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-2 justify-end">
                    <button
                        @click="openEditFromDetail(selectedProject)"
                        class="px-4 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
                    >Modifier</button>
                </div>
            </div>
        </div>

        <div v-if="toDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDelete = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le projet ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDelete.title }} » sera définitivement supprimé.</p>
                <div class="flex gap-3">
                    <button @click="toDelete = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteProject" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
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
                        {{ editTarget ? 'Modifier le projet' : 'Créer un projet' }}
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
                            placeholder="Titre du projet"/>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea v-model="form.description" rows="3"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                            placeholder="Décrivez votre projet…"/>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Date de début</label>
                            <input v-model="form.start_date" type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Date de fin</label>
                            <input v-model="form.end_date" type="date"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Lieu</label>
                            <input v-model="form.location" type="text"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="Adresse ou lieu"/>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">Capacité</label>
                            <input v-model.number="form.capacity" type="number" min="1"
                                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                                placeholder="0 = illimitée"/>
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
                    <p v-if="formError && !premiumError" class="text-xs text-red-500">{{ formError }}</p>
                    <div v-if="premiumError" class="text-xs text-red-600 bg-red-50 rounded-lg p-3">
                        {{ formError }}
                        <router-link to="/artisan/abonnement" class="block mt-1 font-semibold underline text-primary">Voir les offres premium →</router-link>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 space-y-2 flex-shrink-0">
                    <p v-if="limits?.projects && !editTarget" class="text-xs text-gray-400 text-center">
                        {{ limits.projects.used }}/{{ limits.projects.max }} projets {{ limits.is_premium ? 'créés ce mois-ci' : 'créés' }}
                    </p>
                    <div class="flex justify-end gap-2">
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
        </div>

    </ArtisanLayout>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { usePolling } from '@/composables/usePolling.js'
import { useRoute } from 'vue-router'
import ArtisanLayout from '@/Layouts/ArtisanLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'

const route = useRoute()

const projects = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(true)
const search = ref('')
const toDelete = ref(null)
const deleting = ref(false)
const selectedProject = ref(null)
const formModal = ref(false)
const editTarget = ref(null)
const submitting = ref(false)
const formError = ref('')
const premiumError = ref(false)
const limits = ref(null)
const form = ref({ title: '', description: '', start_date: '', end_date: '', location: '', capacity: null })
const filterStatus = ref('')
const materials = ref([])
const steps = ref([])
const showAddMaterial = ref(false)
const showAddStep = ref(false)
const newMaterial = ref({ name: '', quantity: 0, unit: '' })
const newStep = ref({ title: '' })
const editingMaterialId = ref(null)
const editMaterial = ref({ name: '', quantity: 0, unit: '' })
const editingStepId = ref(null)
const editStep = ref({ title: '', description: '' })
const participants = ref([])
const participantsLoading = ref(false)
const photos = ref([])
const documents = ref([])
const docsLoading = ref(false)

function isVideo(link) {
    return /\.(mp4|webm|ogg|mov)(\?|$)/i.test(link)
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

function statusLabel(status) {
    return { open: 'Ouvert', in_progress: 'En cours', completed: 'Terminé', pending: 'En attente', rejected: 'Rejeté' }[status] ?? status
}

function statusClass(status) {
    return {
        open: 'bg-secondary/10 text-secondary',
        in_progress: 'bg-primary/10 text-primary',
        completed: 'bg-gray-100 text-gray-500',
        pending: 'bg-yellow-100 text-yellow-700',
        rejected: 'bg-red-100 text-red-600',
    }[status] ?? 'bg-gray-100 text-gray-500'
}

function openCreate() {
    editTarget.value = null
    form.value = { title: '', description: '', start_date: '', end_date: '', location: '', capacity: null }
    photos.value = []
    formError.value = ''
    formModal.value = true
}

watch(selectedProject, async (p) => {
    materials.value = []
    steps.value = []
    documents.value = []
    docsLoading.value = false
    showAddMaterial.value = false
    showAddStep.value = false
    editingMaterialId.value = null
    editingStepId.value = null
    participants.value = []
    participantsLoading.value = false
    if (!p) return
    loadParticipants(`/projects/${p.id}/members`)
    loadDocuments(p.id)
    try {
        const [mRes, sRes] = await Promise.all([
            api.get(`/projects/${p.id}/materials`),
            api.get(`/projects/${p.id}/steps`),
        ])
        materials.value = mRes.data ?? []
        steps.value = sRes.data ?? []
    } catch {}
})

async function addMaterial() {
    if (!newMaterial.value.name.trim()) return
    try {
        const { data } = await api.post(`/projects/${selectedProject.value.id}/materials`, {
            name: newMaterial.value.name,
            quantity: newMaterial.value.quantity || 0,
            unit: newMaterial.value.unit || null,
        })
        materials.value.push(data)
        newMaterial.value = { name: '', quantity: 0, unit: '' }
        showAddMaterial.value = false
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

function startEditMaterial(m) {
    editingMaterialId.value = m.id
    editMaterial.value = { name: m.name, quantity: m.quantity, unit: m.unit ?? '' }
}

async function saveEditMaterial(m) {
    try {
        await api.patch(`/projects/${selectedProject.value.id}/materials/${m.id}`, {
            name: editMaterial.value.name,
            quantity: editMaterial.value.quantity || 0,
            unit: editMaterial.value.unit || null,
        })
        const idx = materials.value.findIndex(x => x.id === m.id)
        if (idx !== -1) materials.value[idx] = { ...materials.value[idx], name: editMaterial.value.name, quantity: editMaterial.value.quantity, unit: editMaterial.value.unit || null }
        editingMaterialId.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function deleteMaterial(id) {
    try {
        await api.delete(`/projects/${selectedProject.value.id}/materials/${id}`)
        materials.value = materials.value.filter(m => m.id !== id)
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

function stepStatusLabel(s) {
    return { todo: 'À faire', in_progress: 'En cours', done: 'Fait' }[s] ?? s
}

async function addStep() {
    if (!newStep.value.title.trim()) return
    try {
        const { data } = await api.post(`/projects/${selectedProject.value.id}/steps`, {
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
        await api.patch(`/projects/${selectedProject.value.id}/steps/${s.id}/status`, { status: next })
        s.status = next
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

async function deleteStep(id) {
    try {
        await api.delete(`/projects/${selectedProject.value.id}/steps/${id}`)
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
        await api.patch(`/projects/${selectedProject.value.id}/steps/${s.id}`, {
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
            api.patch(`/projects/${selectedProject.value.id}/steps/${a.id}`, { title: a.title, description: a.description ?? null, step_order: b.step_order }),
            api.patch(`/projects/${selectedProject.value.id}/steps/${b.id}`, { title: b.title, description: b.description ?? null, step_order: a.step_order }),
        ])
        const newSteps = [...steps.value]
        newSteps[target] = { ...a, step_order: b.step_order }
        newSteps[index] = { ...b, step_order: a.step_order }
        steps.value = newSteps
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur')
    }
}

function openEditFromDetail(p) {
    selectedProject.value = null
    editTarget.value = p
    form.value = {
        title: p.title,
        description: p.description ?? '',
        start_date: p.start_date ? p.start_date.slice(0, 10) : '',
        end_date: p.end_date ? p.end_date.slice(0, 10) : '',
        location: p.location ?? '',
        capacity: p.capacity ?? null,
    }
    photos.value = []
    formError.value = ''
    formModal.value = true
}

async function vote(p, value) {
    try {
        if (p.user_vote === value) {
            await api.delete(`/projects/${p.id}/vote`)
        } else {
            await api.put(`/projects/${p.id}/vote`, { value })
        }
        await fetchProjects()
    } catch {}
}

let searchDebounce = null
function onSearchInput() {
    clearTimeout(searchDebounce)
    searchDebounce = setTimeout(resetAndFetch, 300)
}

function resetAndFetch() {
    page.value = 1
    fetchProjects()
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
    fetchProjects()
}

async function fetchProjects(silent = false) {
    if (!silent) loading.value = true
    try {
        const params = { page: page.value, limit: 15 }
        if (search.value) params.search = search.value
        if (filterStatus.value) params.status = filterStatus.value
        const { data } = await api.get('/user/my-projects', { params })
        projects.value = data.data ?? []
        total.value = data.total ?? 0
    } catch {
        projects.value = []
        total.value = 0
    } finally {
        if (!silent) loading.value = false
    }
}

async function deleteProject() {
    if (!toDelete.value) return
    deleting.value = true
    try {
        await api.delete(`/projects/${toDelete.value.id}`)
        toDelete.value = null
        await fetchProjects()
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
    premiumError.value = false
    try {
        const photoURLs = photos.value.length ? await uploadPhotos() : []
        const payload = {
            title: form.value.title,
            description: form.value.description || null,
            start_date: form.value.start_date || null,
            end_date: form.value.end_date || null,
            location: form.value.location || null,
            capacity: form.value.capacity || null,
            photo_urls: photoURLs,
        }
        if (editTarget.value) {
            await api.patch(`/projects/${editTarget.value.id}`, payload)
        } else {
            await api.post('/projects', payload)
        }
        formModal.value = false
        photos.value = []
        await fetchProjects()
    } catch (e) {
        premiumError.value = e.response?.status === 403
        formError.value = e.response?.data?.error ?? 'Erreur lors de l\'enregistrement.'
    } finally {
        submitting.value = false
    }
}

usePolling(fetchProjects)

onMounted(async () => {
    try {
        const { data } = await api.get('/user/limits')
        limits.value = data
    } catch {}
    const openId = route.query.open
    if (!openId) return
    await fetchProjects()
    const found = projects.value.find(p => String(p.id) === String(openId))
    if (found) selectedProject.value = found
})
</script>
