<template>
    <SalarieLayout>

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Modération du forum
            </h1>
            <p class="text-sm text-gray-400 mt-1">Consultez et supprimez les sujets et réponses</p>
        </div>

        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center gap-3">
                <div class="relative flex-1">
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 1 0 5 11a6 6 0 0 0 12 0z"/>
                    </svg>
                    <input
                        v-model="search"
                        @input="debouncedFetch"
                        type="text"
                        placeholder="Rechercher"
                        class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                    />
                </div>
            </div>

            <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">Chargement…</div>

            <div v-else-if="topics.length === 0" class="py-12 text-center">
                <svg class="w-12 h-12 text-gray-200 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                </svg>
                <p class="text-gray-400 text-sm">Aucun sujet pour le moment.</p>
            </div>

            <div v-else class="space-y-3 p-6">
                <div
                    v-for="topic in topics"
                :key="topic.id"
                class="bg-white rounded-2xl shadow-sm px-5 py-4 flex items-center gap-4 hover:shadow-md transition-shadow"
            >
                <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
                    <svg class="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                    </svg>
                </div>
                <div class="flex-1 min-w-0 cursor-pointer" @click="openTopic(topic)">
                    <h3 class="font-semibold text-gray-800 text-sm truncate">{{ topic.title }}</h3>
                    <p class="text-xs text-gray-400 mt-0.5">
                        Par {{ topic.author_name }} · {{ formatDate(topic.created_at) }}
                    </p>
                </div>
                <div class="flex items-center gap-3 flex-shrink-0">
                    <span class="text-xs text-gray-400 flex items-center gap-1">
                        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                        </svg>
                        {{ topic.replies_count }}
                    </span>
                    <button
                        v-if="canDeleteTopic(topic)"
                        @click="toDeleteTopic = topic"
                        class="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        title="Supprimer le sujet"
                    >
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                </div>
            </div>
            </div>

            <Pagination
                v-if="total > LIMIT"
                :page="page"
                :total="total"
                :limit="LIMIT"
                @update:page="changePage"
            />
        </div>

        <div
            v-if="selected || selectedLoading || selectedError"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="closeTopic"
        >
            <div class="bg-white rounded-2xl shadow-xl w-full max-w-2xl max-h-[85vh] flex flex-col overflow-hidden">
                <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-100 flex-shrink-0">
                    <h2 class="flex-1 font-semibold text-gray-800 text-sm truncate" style="font-family: var(--font-family-title)">
                        {{ selected?.title ?? '…' }}
                    </h2>
                    <button
                        v-if="selected && canDeleteTopic(selected)"
                        @click="toDeleteTopic = selected; closeTopic()"
                        class="p-1.5 rounded-lg text-red-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                        title="Supprimer le sujet"
                    >
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                    <button @click="closeTopic" class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div v-if="selectedLoading" class="text-center text-gray-400 text-sm py-6">Chargement…</div>
                    <div v-else-if="selectedError" class="text-center text-red-400 text-sm py-6">{{ selectedError }}</div>

                    <template v-else-if="selected">
                        <div
                            v-for="post in selected.posts"
                            :key="post.id"
                            class="flex items-start justify-between gap-3 bg-gray-50 rounded-xl p-4"
                        >
                            <div class="flex-1 min-w-0">
                                <div class="flex items-center gap-2 mb-2">
                                    <div class="w-7 h-7 rounded-full bg-primary flex items-center justify-center text-white text-xs font-semibold flex-shrink-0">
                                        {{ post.author_name?.[0] ?? '?' }}
                                    </div>
                                    <span class="text-sm font-medium text-gray-800">{{ post.author_name }}</span>
                                    <span class="text-xs text-gray-400">· {{ formatDate(post.created_at) }}</span>
                                </div>
                                <p class="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap pl-9">{{ post.body }}</p>
                            </div>
                            <button
                                v-if="canDeletePost(post)"
                                @click="deletePost(post)"
                                class="flex-shrink-0 p-1.5 text-red-400 hover:text-red-600 hover:bg-red-100 rounded-lg transition-colors"
                                title="Supprimer ce message"
                            >
                                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                </svg>
                            </button>
                        </div>
                        <div v-if="!selected.posts?.length" class="text-center text-gray-400 text-sm py-4">
                            Aucun message dans ce sujet.
                        </div>
                    </template>
                </div>
            </div>
        </div>

        <div v-if="toDeleteTopic" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/40" @click="toDeleteTopic = null"/>
            <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6">
                <h3 class="font-semibold text-gray-800 mb-2">Supprimer le sujet ?</h3>
                <p class="text-sm text-gray-500 mb-5">« {{ toDeleteTopic.title }} » et tous ses messages seront définitivement supprimés.</p>
                <div class="flex gap-3">
                    <button @click="toDeleteTopic = null" class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors">Annuler</button>
                    <button @click="deleteTopic" :disabled="deleting" class="flex-1 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors disabled:opacity-60">
                        {{ deleting ? 'Suppression…' : 'Supprimer' }}
                    </button>
                </div>
            </div>
        </div>

    </SalarieLayout>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import SalarieLayout from '@/Layouts/SalarieLayout.vue'
import Pagination from '@/Components/Pagination.vue'
import api from '@/api.js'
import { formatDate } from '@/utils.js'
import { useAuthStore } from '@/stores/auth.js'

const auth = useAuthStore()

const LIMIT = 20

const topics = ref([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const total = ref(0)

const selected = ref(null)
const selectedLoading = ref(false)
const selectedError = ref('')
const toDeleteTopic = ref(null)
const deleting = ref(false)

function canDeleteTopic(topic) {
    return topic?.id_author === auth.user?.id || auth.isAdmin || auth.hasPermission('moderate_forum')
}

function canDeletePost(post) {
    return post.id_author === auth.user?.id || auth.isAdmin || auth.hasPermission('moderate_forum')
}

let debounceTimer = null
function debouncedFetch() {
    clearTimeout(debounceTimer)
    page.value = 1
    debounceTimer = setTimeout(fetchTopics, 300)
}

function changePage(p) {
    page.value = p
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function fetchTopics() {
    loading.value = true
    try {
        const params = { page: page.value, limit: LIMIT }
        if (search.value) params.search = search.value
        const { data } = await api.get('/forum/topics', { params })
        topics.value = data.data ?? []
        total.value = data.total ?? 0
    } catch (err) {
        console.error('fetchTopics error:', err)
        topics.value = []
        total.value = 0
    } finally {
        loading.value = false
    }
}

async function openTopic(topic) {
    selected.value = null
    selectedError.value = ''
    selectedLoading.value = true
    try {
        const { data } = await api.get(`/forum/topics/${topic.id}`)
        selected.value = data
    } catch {
        selectedError.value = 'Impossible de charger ce sujet.'
    } finally {
        selectedLoading.value = false
    }
}

function closeTopic() {
    selected.value = null
    selectedLoading.value = false
    selectedError.value = ''
}

async function deleteTopic() {
    if (!toDeleteTopic.value) return
    deleting.value = true
    try {
        await api.delete(`/forum/topics/${toDeleteTopic.value.id}`)
        topics.value = topics.value.filter(t => t.id !== toDeleteTopic.value.id)
        total.value = Math.max(total.value - 1, 0)
        toDeleteTopic.value = null
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    } finally {
        deleting.value = false
    }
}

async function deletePost(post) {
    if (!selected.value) return
    try {
        const { data } = await api.delete(`/forum/topics/${selected.value.id}/posts/${post.id}`)
        selected.value = data
        const t = topics.value.find(t => t.id === data.id)
        if (t) t.replies_count = Math.max((data.posts?.length ?? 1) - 1, 0)
    } catch (e) {
        alert(e.response?.data?.error ?? 'Erreur lors de la suppression.')
    }
}

watch(page, fetchTopics)
onMounted(fetchTopics)
</script>
