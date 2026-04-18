<template>
  <AdminLayout title="Rôles & Permissions">

    <div class="flex items-center justify-between mb-6">
      <p class="text-sm text-gray-500">
        {{roles.length}} rôle(s)
      </p>
    </div>

    <div v-if="loading" class="text-sm text-gray-400 py-10 text-center">Chargement…</div>

    <div v-else class="grid grid-cols-2 gap-5">
      <div
        v-for="role in roles"
        :key="role.id"
        class="bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <div class="flex items-center gap-2.5">
            <div
              class="flex-shrink-0 w-9 h-9 rounded-lg flex items-center justify-center"
              :style="{backgroundColor: metaFor(role.name).color + '20'}">
              <div :style="{color: metaFor(role.name).color}" v-html="metaFor(role.name).icon" />
            </div>
            <h2 class="text-sm font-semibold text-gray-800" style="font-family: var(--font-family-title)">{{role.name}}</h2>
          </div>
          <button
            @click="openEdit(role)"
            title="Modifier les permissions"
            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors duration-150">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
            </svg>
          </button>
        </div>

        <div class="px-5 py-4">
          <div class="space-y-3">
            <div
              v-for="cat in permissionCategories"
              :key="cat.domain">
              <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">{{cat.domain}}</p>
              <div class="flex flex-wrap gap-1.5">
                <span
                  v-for="perm in cat.perms"
                  :key="perm.id"
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium transition-colors"
                  :class="role.permIds.includes(perm.id)
                    ? 'bg-secondary/15 text-secondary-dark'
                    : 'bg-gray-100 text-gray-300 line-through'">
                  <svg v-if="role.permIds.includes(perm.id)" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                  </svg>
                  {{perm.name}}
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
            Permissions — {{modal.roleName}}
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
          <div
            v-for="cat in permissionCategories"
            :key="cat.domain">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">{{cat.domain}}</p>
            <div class="space-y-1.5">
              <label
                v-for="perm in cat.perms"
                :key="perm.id"
                class="flex items-center gap-3 px-3 py-2 rounded-lg border cursor-pointer transition-colors duration-150"
                :class="modal.permIds.includes(perm.id)
                  ? 'border-primary/30 bg-primary/5'
                  : 'border-gray-100 hover:border-gray-200'">
                <input
                  type="checkbox"
                  :value="perm.id"
                  v-model="modal.permIds"
                  class="accent-primary w-3.5 h-3.5"
                />
                <p class="text-sm text-gray-800">{{perm.name}}</p>
              </label>
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
            @click="savePermissions"
            :disabled="saving"
            class="px-4 py-2 rounded-lg text-sm font-medium bg-primary text-white hover:bg-primary-dark transition-colors duration-150 disabled:opacity-50">
            {{saving ? 'Enregistrement…' : 'Enregistrer'}}
          </button>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '@/Layouts/AdminLayout.vue'
import api from '@/api.js'

const ROLE_META = {
  admin:   { color: '#c46d68', icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>` },
  artisan: { color: '#8dc734', icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>` },
  salarie: { color: '#f59e0b', icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>` },
  user:    { color: '#6366f1', icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>` },
}
const DEFAULT_META = { color: '#64748b', icon: `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" style="width:18px;height:18px"><path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>` }

function metaFor(name) {
  return ROLE_META[name] ?? DEFAULT_META
}

const loading = ref(true)
const saving = ref(false)
const roles = ref([])
const allPermissions = ref([])

const permissionCategories = computed(() => {
  const map = {}
  for (const p of allPermissions.value) {
    if (!map[p.domain]) map[p.domain] = []
    map[p.domain].push(p)
  }
  return Object.entries(map).map(([domain, perms]) => ({ domain, perms }))
})

const modal = ref({
  open: false,
  roleId: null,
  roleName: '',
  originalPermIds: [],
  permIds: [],
})

onMounted(async () => {
  try {
    const [{ data: permsData }, { data: rolesData }] = await Promise.all([
      api.get('/admin/permissions'),
      api.get('/admin/roles'),
    ])
    allPermissions.value = permsData ?? []

    roles.value = await Promise.all(
      (rolesData ?? []).map(async r => {
        const { data: rPerms } = await api.get(`/admin/role/${r.id}/permissions`)
        return { ...r, permIds: (rPerms ?? []).map(p => p.id) }
      })
    )
  } catch (e) {
    console.error('Roles fetch error:', e)
  } finally {
    loading.value = false
  }
})

function openEdit(role) {
  modal.value = {
    open: true,
    roleId: role.id,
    roleName: role.name,
    originalPermIds: [...role.permIds],
    permIds: [...role.permIds],
  }
}

async function savePermissions() {
  saving.value = true
  try {
    const orig = new Set(modal.value.originalPermIds)
    const next = new Set(modal.value.permIds)

    const toAdd = [...next].filter(id => !orig.has(id))
    const toRemove = [...orig].filter(id => !next.has(id))

    await Promise.all([
      ...toAdd.map(permId => api.post(`/admin/role/${modal.value.roleId}/permissions`, { permission_id: permId })),
      ...toRemove.map(permId => api.delete(`/admin/role/${modal.value.roleId}/permissions/${permId}`)),
    ])

    const role = roles.value.find(r => r.id === modal.value.roleId)
    if (role) role.permIds = [...modal.value.permIds]

    modal.value.open = false
  } catch (e) {
    console.error('Save permissions error:', e)
  } finally {
    saving.value = false
  }
}
</script>
