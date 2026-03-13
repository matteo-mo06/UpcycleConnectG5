<template>
  <AdminLayout title="Rôles & Permissions">

    <div class="flex items-center justify-between mb-6">
      <p class="text-sm text-gray-500">
        {{roles.length}} rôle(s) ·
        {{roles.filter(r => !r.system).length}} personnalisé(s)
      </p>
      <button
        @click="openCreate"
        class="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors duration-150">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
        </svg>
        Nouveau rôle
      </button>
    </div>

    <div class="grid grid-cols-2 gap-5">
      <div
        v-for="role in roles"
        :key="role.id"
        class="bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <div class="flex items-center gap-2.5">
            <div
              class="flex-shrink-0 w-9 h-9 rounded-lg flex items-center justify-center"
              :style="{backgroundColor: role.color + '20'}">
              <div :style="{color: role.color}" v-html="role.icon" />
            </div>
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <h2 class="text-sm font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{role.name}}</h2>
                <span
                  v-if="role.system"
                  class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-500">
                  Système
                </span>
              </div>
              <p class="text-xs text-gray-400 truncate">{{role.userCount}} utilisateur(s)</p>
            </div>
          </div>
          <div class="flex items-center gap-1 flex-shrink-0">
            <button
              @click="openEdit(role)"
              title="Modifier"
              class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors duration-150">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
              </svg>
            </button>
            <button
              v-if="!role.system"
              @click="deleteRole(role)"
              title="Supprimer"
              class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors duration-150">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>

        <div class="px-5 py-4">
          <p class="text-xs text-gray-500 mb-4">{{role.description}}</p>
          <div class="space-y-3">
            <div
              v-for="category in permissionCategories"
              :key="category.id">
              <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">{{category.label}}</p>
              <div class="flex gap-1.5">
                <span
                  v-for="perm in category.permissions"
                  :key="perm.id"
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium transition-colors"
                  :class="role.permissions.includes(perm.id)
                    ? 'bg-secondary/15 text-secondary-dark'
                    : 'bg-gray-100 text-gray-300 line-through'">
                  <svg v-if="role.permissions.includes(perm.id)" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                  </svg>
                  {{perm.label}}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="modal.open"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="modal.open = false">
      <div class="absolute inset-0 bg-black/40" @click="modal.open = false" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">

        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
          <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
            {{modal.isEdit ? 'Modifier le rôle' : 'Nouveau rôle'}}
          </h3>
          <button
            @click="modal.open = false"
            class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="px-6 py-5 overflow-y-auto space-y-5">

          <div class="space-y-3">
            <div>
              <label class="block text-xs font-medium text-gray-500 mb-1.5">Nom du rôle</label>
              <input
                v-model="modal.name"
                type="text"
                placeholder="Ex : Modérateur"
                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                :disabled="modal.isEdit && modal.system"
              />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-500 mb-1.5">Description</label>
              <textarea
                v-model="modal.description"
                rows="2"
                placeholder="Décrivez les responsabilités de ce rôle…"
                class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
              />
            </div>
          </div>

          <div>
            <p class="text-xs font-medium text-gray-500 mb-3">Permissions</p>
            <div class="space-y-4">
              <div
                v-for="category in permissionCategories"
                :key="category.id">
                <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">{{category.label}}</p>
                <div class="space-y-1.5">
                  <label
                    v-for="perm in category.permissions"
                    :key="perm.id"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg border cursor-pointer transition-colors duration-150"
                    :class="modal.permissions.includes(perm.id)
                      ? 'border-primary/30 bg-primary/5'
                      : 'border-gray-100 hover:border-gray-200'">
                    <input
                      type="checkbox"
                      :value="perm.id"
                      v-model="modal.permissions"
                      class="accent-primary w-3.5 h-3.5"
                    />
                    <div class="min-w-0">
                      <p class="text-sm text-gray-800">{{perm.label}}</p>
                      <p class="text-xs text-gray-400">{{perm.description}}</p>
                    </div>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2 flex-shrink-0">
          <button
            @click="modal.open = false"
            class="px-4 py-2 rounded-lg text-sm text-gray-500 border border-gray-200 hover:bg-gray-50 transition-colors duration-150">
            Annuler
          </button>
          <button
            @click="saveRole"
            class="px-4 py-2 rounded-lg text-sm font-medium bg-primary text-white hover:bg-primary-dark transition-colors duration-150">
            {{modal.isEdit ? 'Enregistrer' : 'Créer le rôle'}}
          </button>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import {ref} from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'

const permissionCategories = [
  {
    id: 'general',
    label: 'Général',
    permissions: [
      {id: 'browse', label: 'Consulter la plateforme', description: 'Accéder aux annonces, événements et profils publics'},
      {id: 'messaging', label: 'Messagerie interne', description: 'Envoyer et recevoir des messages entre utilisateurs'},
    ],
 },
  {
    id: 'listings',
    label: 'Annonces',
    permissions: [
      {id: 'listings_create', label: 'Publier des annonces', description: 'Créer et gérer ses propres annonces'},
      {id: 'listings_pro', label: 'Annonces professionnelles', description: 'Publier des annonces avec badge professionnel'},
    ],
 },
  {
    id: 'events',
    label: 'Événements',
    permissions: [
      {id: 'events_create', label: 'Créer des événements', description: 'Proposer des événements (soumis à validation)'},
      {id: 'events_register', label: "S'inscrire aux événements", description: "Participer aux événements de la plateforme"},
    ],
 },
  {
    id: 'moderation',
    label: 'Modération',
    permissions: [
      {id: 'report', label: 'Signaler du contenu', description: 'Soumettre des signalements à la modération'},
      {id: 'moderate', label: 'Modérer du contenu', description: 'Traiter les signalements et masquer des contenus'},
    ],
 },
  {
    id: 'admin',
    label: 'Administration',
    permissions: [
      {id: 'backoffice', label: 'Accéder au back-office', description: 'Accès au tableau de bord administrateur'},
      {id: 'users_manage', label: 'Gérer les utilisateurs', description: 'Modifier les comptes, rôles et statuts'},
      {id: 'roles_manage', label: 'Gérer les rôles', description: 'Créer, modifier et supprimer des rôles'},
    ],
 },
]

const roles = ref([
  {
    id: 1,
    name: 'Particulier',
    description: 'Accès standard à la plateforme. Peut publier des annonces et échanger avec d\'autres membres.',
    permissions: ['browse', 'messaging', 'listings_create', 'events_register', 'report'],
    userCount: 892,
    system: true,
    color: '#6366f1',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>`,
 },
  {
    id: 2,
    name: 'Artisan',
    description: 'Compte professionnel vérifié. Peut publier des annonces avec badge artisan et accéder aux fonctionnalités pro.',
    permissions: ['browse', 'messaging', 'listings_create', 'listings_pro', 'events_register', 'report'],
    userCount: 234,
    system: true,
    color: '#8dc734',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>`,
 },
  {
    id: 3,
    name: 'Salarié',
    description: 'Employé d\'un espace partenaire. Peut créer et organiser des événements soumis à validation.',
    permissions: ['browse', 'messaging', 'listings_create', 'events_create', 'events_register', 'report', 'moderate'],
    userCount: 17,
    system: true,
    color: '#f59e0b',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>`,
 },
  {
    id: 4,
    name: 'Administrateur',
    description: 'Accès complet à la plateforme et au back-office. Gère les utilisateurs, les rôles et la modération.',
    permissions: ['browse', 'messaging', 'listings_create', 'listings_pro', 'events_create', 'events_register', 'report', 'moderate', 'backoffice', 'users_manage', 'roles_manage'],
    userCount: 3,
    system: true,
    color: '#c46d68',
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>`,
 },
])

const modal = ref({
  open: false,
  isEdit: false,
  editId: null,
  system: false,
  name: '',
  description: '',
  permissions: [],
})

function openCreate() {
  modal.value = {
    open: true,
    isEdit: false,
    editId: null,
    system: false,
    name: '',
    description: '',
    permissions: ['browse'],
 }
}

function openEdit(role) {
  modal.value = {
    open: true,
    isEdit: true,
    editId: role.id,
    system: role.system,
    name: role.name,
    description: role.description,
    permissions: [...role.permissions],
 }
}

function saveRole() {
  if (!modal.value.name.trim()) return

  if (modal.value.isEdit) {
    const target = roles.value.find(r => r.id === modal.value.editId)
    target.description = modal.value.description
    target.permissions = [...modal.value.permissions]
    if (!target.system) target.name = modal.value.name
 } else {
    roles.value.push({
      id: Date.now(),
      name: modal.value.name,
      description: modal.value.description,
      permissions: [...modal.value.permissions],
      userCount: 0,
      system: false,
      color: '#64748b',
      icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>`,
   })
 }

  modal.value.open = false
}

function deleteRole(role) {
  if (!confirm(`Supprimer le rôle "${role.name}" ?`)) return
  roles.value = roles.value.filter(r => r.id !== role.id)
}
</script>
