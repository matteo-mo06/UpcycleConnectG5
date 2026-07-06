<template>
    <AdminLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">Rôles & Permissions</h1>
        </div>

        <div class="flex items-center justify-between mb-6">
            <p class="text-sm text-gray-500">{{roles.length}} rôle(s)</p>
            <button
                @click="openCreate"
                class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
                Créer un rôle
            </button>
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
                    <div class="flex items-center gap-1">
                        <button
                            @click="openEdit(role)"
                            title="Modifier les permissions"
                            class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors duration-150">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                            </svg>
                        </button>
                        <button
                            v-if="role.name !== 'admin'"
                            @click="confirmDelete(role)"
                            title="Supprimer le rôle"
                            class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors duration-150">
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                            </svg>
                        </button>
                    </div>
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
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden flex flex-col max-h-[90vh]">

                <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                    <div>
                        <h3 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                            Permissions : {{modal.roleName}}
                        </h3>
                        <p v-if="isAdminRole" class="text-xs text-amber-600 mt-1">L'administrateur possède toutes les permissions par défaut.</p>
                    </div>
                    <button
                        @click="modal.open = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="px-6 py-5 overflow-y-auto space-y-5 flex-1">
                    <div
                        v-for="cat in permissionCategories"
                        :key="cat.domain">
                        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">{{cat.domain}}</p>
                        <p v-if="cat.domain === 'espaces'" class="text-xs text-gray-400 mb-2">
                            Un seul espace accessible à la fois (ou aucun).
                        </p>
                        <div class="space-y-1.5">
                            <template v-if="cat.domain === 'espaces'">
                                <label
                                    class="flex items-center gap-3 px-3 py-2 rounded-lg border cursor-pointer transition-colors duration-150"
                                    :class="(!isAdminRole && !cat.perms.some(p => modal.permIds.includes(p.id)))
                                        ? 'border-primary/30 bg-primary/5'
                                        : 'border-gray-100 hover:border-gray-200'">
                                    <input
                                        type="radio"
                                        name="espace-access"
                                        :checked="!isAdminRole && !cat.perms.some(p => modal.permIds.includes(p.id))"
                                        :disabled="isAdminRole"
                                        @change="selectSpacePermission(null, cat.perms)"
                                        class="accent-primary w-3.5 h-3.5"
                                    />
                                    <p class="text-sm text-gray-800">Aucun accès spécifique</p>
                                </label>
                                <label
                                    v-for="perm in cat.perms"
                                    :key="perm.id"
                                    class="flex items-center gap-3 px-3 py-2 rounded-lg border cursor-pointer transition-colors duration-150"
                                    :class="(isAdminRole || modal.permIds.includes(perm.id))
                                        ? 'border-primary/30 bg-primary/5'
                                        : 'border-gray-100 hover:border-gray-200'">
                                    <input
                                        type="radio"
                                        name="espace-access"
                                        :checked="isAdminRole || modal.permIds.includes(perm.id)"
                                        :disabled="isAdminRole"
                                        @change="selectSpacePermission(perm.id, cat.perms)"
                                        class="accent-primary w-3.5 h-3.5"
                                    />
                                    <p class="text-sm text-gray-800">{{perm.name}}</p>
                                </label>
                            </template>
                            <template v-else>
                                <label
                                    v-for="perm in cat.perms"
                                    :key="perm.id"
                                    class="flex items-center gap-3 px-3 py-2 rounded-lg border cursor-pointer transition-colors duration-150"
                                    :class="(isAdminRole || modal.permIds.includes(perm.id))
                                        ? 'border-primary/30 bg-primary/5'
                                        : 'border-gray-100 hover:border-gray-200'">
                                    <input
                                        type="checkbox"
                                        :checked="isAdminRole || modal.permIds.includes(perm.id)"
                                        :disabled="isAdminRole"
                                        @change="togglePerm(perm.id)"
                                        class="accent-primary w-3.5 h-3.5"
                                    />
                                    <p class="text-sm text-gray-800">{{perm.name}}</p>
                                </label>
                            </template>
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

        <div v-if="formModal.open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="formModal.open = false" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-3" style="font-family: var(--font-family-title)">Créer un rôle</h3>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nom du rôle <span class="text-red-400">*</span></label>
                <input
                    v-model="formModal.name"
                    type="text"
                    placeholder="Nom du rôle"
                    class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary mb-4"
                    @keyup.enter="saveForm"
                />
                <p v-if="formModal.error" class="text-xs text-red-500 mb-3">{{ formModal.error }}</p>
                <div class="flex gap-3">
                    <button
                        @click="formModal.open = false"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="saveForm"
                        :disabled="formModal.saving"
                        class="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-60">
                        {{ formModal.saving ? 'Enregistrement…' : 'Créer' }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="deleteModal.role" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="deleteModal.role = null" />
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le rôle ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ deleteModal.role.name }} » sera définitivement supprimé.</p>
                <p v-if="deleteModal.error" class="text-xs text-red-500 mb-3">{{ deleteModal.error }}</p>
                <div class="flex gap-3">
                    <button
                        @click="deleteModal.role = null"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">
                        Annuler
                    </button>
                    <button
                        @click="deleteRole"
                        :disabled="deleteModal.deleting"
                        class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleteModal.deleting ? 'Suppression…' : 'Supprimer' }}
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

const formModal = ref({ open: false, name: '', saving: false, error: '' })
const deleteModal = ref({ role: null, deleting: false, error: '' })

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

const isAdminRole = computed(() => modal.value.roleName === 'admin')

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

function openCreate() {
    formModal.value = { open: true, name: '', saving: false, error: '' }
}

async function saveForm() {
    const name = formModal.value.name.trim()
    if (!name) { formModal.value.error = 'Le nom est requis.'; return }
    formModal.value.saving = true
    formModal.value.error = ''
    try {
        const { data } = await api.post('/admin/roles', { name })
        roles.value.push({ ...data, permIds: [] })
        formModal.value.open = false
    } catch (e) {
        formModal.value.error = e.response?.data?.error ?? 'Erreur lors de la création.'
    } finally {
        formModal.value.saving = false
    }
}

function confirmDelete(role) {
    deleteModal.value = { role, deleting: false, error: '' }
}

async function deleteRole() {
    deleteModal.value.deleting = true
    deleteModal.value.error = ''
    try {
        await api.delete(`/admin/role/${deleteModal.value.role.id}`)
        roles.value = roles.value.filter(r => r.id !== deleteModal.value.role.id)
        deleteModal.value.role = null
    } catch (e) {
        deleteModal.value.error = e.response?.data?.error ?? 'Erreur lors de la suppression.'
    } finally {
        deleteModal.value.deleting = false
    }
}

function togglePerm(permId) {
    if (isAdminRole.value) return
    const idx = modal.value.permIds.indexOf(permId)
    if (idx === -1) modal.value.permIds.push(permId)
    else modal.value.permIds.splice(idx, 1)
}

function selectSpacePermission(permId, groupPerms) {
    if (isAdminRole.value) return
    const groupIds = groupPerms.map(p => p.id)
    modal.value.permIds = modal.value.permIds.filter(id => !groupIds.includes(id))
    if (permId) modal.value.permIds.push(permId)
}

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
    if (isAdminRole.value) return
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
