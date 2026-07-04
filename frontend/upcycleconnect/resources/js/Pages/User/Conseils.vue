<template>
    <UserLayout>
        <div class="mb-6">
            <h1
                class="text-3xl font-bold text-gray-800"
                style="font-family: var(--font-family-title)"
            >
                Conseils
            </h1>
            <p class="text-sm text-gray-400 mt-1">
                Astuces et bonnes pratiques pour votre démarche upcycling
            </p>
        </div>

        <div class="flex gap-3 mb-6">
            <div class="flex-1 relative">
                <svg
                    class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.8"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M21 21l-4.35-4.35M17 11A6 6 0 1 0 5 11a6 6 0 0 0 12 0z"
                    />
                </svg>
                <input
                    v-model="search"
                    @input="debouncedFetch"
                    type="text"
                    placeholder="Rechercher un conseil…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
            <select
                v-model="filterCategory"
                @change="debouncedFetch"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Tous les thèmes</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                </option>
            </select>
        </div>

        <div v-if="loading" class="bg-white rounded-2xl shadow-sm p-12 text-center text-gray-400 text-sm">
            Chargement…
        </div>

        <div
            v-else-if="conseils.length === 0"
            class="bg-white rounded-2xl shadow-sm p-12 text-center"
        >
            <svg
                class="w-12 h-12 text-gray-200 mx-auto mb-3"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                />
            </svg>
            <p class="text-gray-400 text-sm">
                Aucun conseil disponible pour le moment.
            </p>
        </div>

        <div v-else class="grid grid-cols-3 gap-5">
            <div
                v-for="conseil in conseils"
                :key="conseil.id"
                class="bg-white rounded-2xl shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow cursor-pointer"
                @click="openDetail(conseil)"
            >
                <div
                    class="w-10 h-10 rounded-xl bg-secondary/10 flex items-center justify-center flex-shrink-0"
                >
                    <svg
                        class="w-5 h-5 text-secondary"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="1.8"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                        />
                    </svg>
                </div>
                <h3 class="font-semibold text-gray-800 text-sm leading-snug">
                    {{ conseil.title }}
                </h3>
                <p class="text-xs text-gray-500 line-clamp-3">
                    {{ conseil.description }}
                </p>
                <div class="flex items-center gap-2 mt-auto pt-2 border-t border-gray-50">
                    <span
                        v-if="conseil.category_name"
                        class="px-2 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-medium"
                    >
                        {{ conseil.category_name }}
                    </span>
                    <span class="text-xs text-gray-400 ml-auto">{{ conseil.creator_name }}</span>
                </div>
            </div>
        </div>

        <div
            v-if="selected"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="selected = null"
        >
            <div
                class="bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[85vh] flex flex-col overflow-hidden"
            >
                <div
                    class="flex items-center justify-between px-6 py-4 border-b border-gray-100 flex-shrink-0"
                >
                    <div class="flex-1 min-w-0">
                        <h2
                            class="font-semibold text-gray-800"
                            style="font-family: var(--font-family-title)"
                        >
                            {{ selected.title }}
                        </h2>
                        <p v-if="selected.category_name" class="text-xs text-primary mt-0.5">{{ selected.category_name }}</p>
                    </div>
                    <button
                        @click="selected = null"
                        class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors ml-3 flex-shrink-0"
                    >
                        <svg
                            class="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>
                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <p v-if="selected.description" class="text-sm text-gray-600 whitespace-pre-line">
                        {{ selected.description }}
                    </p>
                    <p v-else class="text-sm text-gray-400 italic">Aucun contenu.</p>

                    <div v-if="docsLoading" class="text-xs text-gray-400">Chargement des photos…</div>
                    <div v-else-if="selectedDocs.length" class="flex gap-2 flex-wrap">
                        <a
                            v-for="doc in selectedDocs"
                            :key="doc.id"
                            :href="doc.link"
                            target="_blank"
                            class="block w-24 h-24 rounded-lg overflow-hidden border border-gray-100"
                        >
                            <img :src="doc.link" class="w-full h-full object-cover" />
                        </a>
                    </div>

                    <div class="flex items-center justify-between text-xs text-gray-400 pt-2 border-t border-gray-50">
                        <span>{{ selected.creator_name }}</span>
                        <span>{{ selected.created_at?.slice(0, 10) }}</span>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex-shrink-0">
                    <button
                        @click="selected = null"
                        class="w-full py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        </div>
    </UserLayout>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import UserLayout from '@/Layouts/UserLayout.vue'
import api from '@/api.js'

const conseils = ref([])
const categories = ref([])
const loading = ref(false)
const search = ref('')
const filterCategory = ref('')
const selected = ref(null)
const selectedDocs = ref([])
const docsLoading = ref(false)

let debounceTimer = null
function debouncedFetch() {
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(fetchConseils, 300)
}

async function fetchConseils() {
    loading.value = true
    try {
        const params = {}
        if (search.value) params.search = search.value
        if (filterCategory.value) params.id_category = filterCategory.value
        const { data } = await api.get('/conseils', { params })
        conseils.value = data?.data ?? data ?? []
    } catch {
        conseils.value = []
    } finally {
        loading.value = false
    }
}

async function openDetail(conseil) {
    selected.value = conseil
    selectedDocs.value = []
    docsLoading.value = true
    try {
        const { data } = await api.get(`/conseils/${conseil.id}/documents`)
        selectedDocs.value = data ?? []
    } catch {
        selectedDocs.value = []
    } finally {
        docsLoading.value = false
    }
}

onMounted(() => {
    api.get('/categories').then(({ data }) => {
        categories.value = Array.isArray(data) ? data : (data.data ?? [])
    })
})

</script>
