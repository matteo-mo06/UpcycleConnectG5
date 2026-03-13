<template>
  <AdminLayout title="Gestion des utilisateurs">

    <div v-if="proRequests.length" class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
      <div class="px-5 py-4 border-b border-gray-100 flex items-center gap-2">
        <h2 class="text-base font-semibold text-gray-800" >Demandes de comptes professionnels</h2>
        <span class="inline-flex items-center justify-center w-5 h-5 rounded-full bg-primary text-white text-xs font-bold">{{proRequests.length}}</span>
      </div>
      <div class="divide-y divide-gray-50">
        <div
          v-for="req in proRequests"
          :key="req.id"
          class="px-5 py-3.5 flex items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
              </svg>
            </div>
            <div class="min-w-0">
              <p class="text-sm font-medium text-gray-800 truncate">{{req.name}}</p>
              <p class="text-xs text-gray-500 truncate">{{req.email}} · Demande le {{req.requestDate}}</p>
            </div>
          </div>
          <p class="text-xs text-gray-500 flex-1 truncate italic">« {{req.reason}} »</p>
          <div class="flex items-center gap-2 flex-shrink-0">
            <button
              @click="approveRequest(req.id)"
              class="px-3 py-1.5 rounded-lg text-xs font-medium bg-secondary text-white hover:bg-secondary-dark transition-colors duration-150">
              Valider
            </button>
            <button
              @click="rejectRequest(req.id)"
              class="px-3 py-1.5 rounded-lg text-xs font-medium bg-gray-100 text-gray-600 hover:bg-red-100 hover:text-red-600 transition-colors duration-150">
              Refuser
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 items-center">
      <input
        v-model="search"
        type="text"
        placeholder="Rechercher par nom, email..."
        class="flex-1 px-4 py-2 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
      />
      <select
        v-model="filterType"
        class="px-4 py-2 rounded-lg border border-gray-200 text-sm text-gray-600 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary">
        <option value="">Tous les types</option>
        <option value="Particulier">Particulier</option>
        <option value="Artisan">Artisan</option>
        <option value="Salarié">Salarié</option>
        <option value="Administrateur">Administrateur</option>
      </select>
      <select
        v-model="filterStatus"
        class="px-4 py-2 rounded-lg border border-gray-200 text-sm text-gray-600 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary">
        <option value="">Tous les statuts</option>
        <option value="Actif">Actif</option>
        <option value="Suspendu">Suspendu</option>
        <option value="Blacklisté">Blacklisté</option>
      </select>
      <span class="text-xs text-gray-400 ml-auto">{{filteredUsers.length}} utilisateur(s)</span>
    </div>

    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-primary">
              <th class="text-left text-white font-medium px-5 py-3">Utilisateur</th>
              <th class="text-left text-white font-medium px-5 py-3">Email</th>
              <th class="text-left text-white font-medium px-5 py-3">Type</th>
              <th class="text-left text-white font-medium px-5 py-3">Statut</th>
              <th class="text-left text-white font-medium px-5 py-3">Inscription</th>
              <th class="text-right text-white font-medium px-5 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(user, i) in filteredUsers"
              :key="user.id"
              :class="['border-b border-gray-50', i % 2 === 0 ? 'bg-white' : 'bg-gray-50/40']">
              <td class="px-5 py-3">
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
                    <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
                    </svg>
                  </div>
                  <div class="min-w-0">
                    <p class="font-medium text-gray-800 truncate">{{user.name}}</p>
                    <p class="text-xs text-gray-400">#{{user.id}}</p>
                  </div>
                </div>
              </td>

              <td class="px-5 py-3 text-gray-500">{{user.email}}</td>

              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="typeBadge(user.type)">
                  {{user.type}}
                </span>
              </td>

              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="statusBadge(user.status)">
                  <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(user.status)" />
                  {{user.status}}
                </span>
              </td>

              <td class="px-5 py-3 text-gray-500">{{user.date}}</td>

              <td class="px-5 py-3">
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click="openRights(user)"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors duration-150"
                    title="Gérer les droits">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
                    </svg>
                  </button>

                  <button
                    @click="toggleBlacklist(user)"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-50 transition-colors duration-150"
                    :title="user.status === 'Blacklisté' ? 'Retirer de la blacklist' : 'Blacklister'">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
                    </svg>
                  </button>

                  <button
                    @click="deleteUser(user)"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors duration-150"
                    title="Supprimer le compte">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredUsers.length === 0">
              <td colspan="6" class="px-5 py-12 text-center text-gray-400 text-sm">
                Aucun utilisateur ne correspond à votre recherche.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="px-5 py-3 border-t border-gray-100 text-xs text-gray-400">{{filteredUsers.length}} utilisateur(s)</div>
    </div>

    <div
      v-if="rightsModal.open"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="rightsModal.open = false">
      <div class="absolute inset-0 bg-black/40" @click="rightsModal.open = false" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6 z-10">
        <h3 class="text-base font-semibold text-gray-800 mb-5" >Rôle :</h3>
        <div class="space-y-2">
          <label
            v-for="role in roles"
            :key="role.value"
            class="flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors duration-150"
            :class="rightsModal.selectedRole === role.value ? 'border-primary bg-primary/5' : 'border-gray-100 hover:border-gray-200'">
            <input
              type="radio"
              :value="role.value"
              v-model="rightsModal.selectedRole"
              class="accent-primary"
            />
            <p class="text-sm font-medium text-gray-800">{{role.label}}</p>
          </label>
        </div>
        <div class="flex gap-2 mt-5">
          <button
            @click="rightsModal.open = false"
            class="flex-1 py-2 rounded-lg text-sm text-gray-500 border border-gray-200 hover:bg-gray-50 transition-colors duration-150">
            Annuler
          </button>
          <button
            @click="saveRights"
            class="flex-1 py-2 rounded-lg text-sm font-medium bg-primary text-white hover:bg-primary-dark transition-colors duration-150">
            Enregistrer
          </button>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import {ref, computed} from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'

const search = ref('')
const filterType = ref('')
const filterStatus = ref('')

const proRequests = ref([
  {id: 101, name: 'Artisan', email: 'artisan@test.fr', requestDate: '01/01/2025', reason: 'Demande de compte professionnel.'},
])

const users = ref([
  {id: 1, name: 'Particulier', email: 'particulier@test.fr', type: 'Particulier', status: 'Actif', date: '01/01/2025'},
  {id: 2, name: 'Artisan', email: 'artisan@test.fr', type: 'Artisan', status: 'Actif', date: '01/01/2025'},
  {id: 3, name: 'Salarié', email: 'salarie@test.fr', type: 'Salarié', status: 'Suspendu', date: '01/01/2025'},
  {id: 4, name: 'Administrateur', email: 'admin@test.fr', type: 'Administrateur', status: 'Actif', date: '01/01/2025'},
])

const filteredUsers = computed(() => {
  return users.value.filter(u => {
    const matchSearch = !search.value
      || u.name.toLowerCase().includes(search.value.toLowerCase())
      || u.email.toLowerCase().includes(search.value.toLowerCase())
      || String(u.id).includes(search.value)
    const matchType = !filterType.value || u.type === filterType.value
    const matchStatus = !filterStatus.value || u.status === filterStatus.value
    return matchSearch && matchType && matchStatus
 })
})

const rightsModal = ref({open: false, user: null, selectedRole: ''})

const roles = [
  {value: 'Particulier', label: 'Particulier'},
  {value: 'Artisan', label: 'Artisan'},
  {value: 'Salarié', label: 'Salarié'},
  {value: 'Administrateur', label: 'Administrateur'},
]

function openRights(user) {
  rightsModal.value = {open: true, user, selectedRole: user.type}
}

function saveRights() {
  rightsModal.value.user.type = rightsModal.value.selectedRole
  rightsModal.value.open = false
}

function toggleBlacklist(user) {
  user.status = user.status === 'Blacklisté' ? 'Actif' : 'Blacklisté'
}

function deleteUser(user) {
  users.value = users.value.filter(u => u.id !== user.id)
}

function approveRequest(id) {
  proRequests.value = proRequests.value.filter(r => r.id !== id)
}

function rejectRequest(id) {
  proRequests.value = proRequests.value.filter(r => r.id !== id)
}

function typeBadge(type) {
  const map = {
    'Particulier': 'bg-blue-100 text-blue-700',
    'Artisan': 'bg-secondary/20 text-secondary-dark',
    'Salarié': 'bg-purple-100 text-purple-700',
    'Administrateur': 'bg-primary/15 text-primary-dark',
 }
  return map[type] ?? 'bg-gray-100 text-gray-600'
}

function statusBadge(status) {
  const map = {
    'Actif': 'bg-green-100 text-green-700',
    'Suspendu': 'bg-amber-100 text-amber-700',
    'Blacklisté': 'bg-red-100 text-red-700',
 }
  return map[status] ?? 'bg-gray-100 text-gray-600'
}

function statusDot(status) {
  const map = {
    'Actif': 'bg-green-500',
    'Suspendu': 'bg-amber-500',
    'Blacklisté': 'bg-red-500',
 }
  return map[status] ?? 'bg-gray-400'
}
</script>
