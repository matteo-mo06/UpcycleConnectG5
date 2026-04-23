<template>
    <UserLayout>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1
                    class="text-3xl font-bold text-gray-800"
                    style="font-family: var(--font-family-title)"
                >
                    Forum
                </h1>
                <p class="text-sm text-gray-400 mt-1">
                    Échangez avec la communauté upcycling
                </p>
            </div>
            <button
                @click="showNewTopic = true"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg
                    class="w-4 h-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.8"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M12 4v16m8-8H4"
                    />
                </svg>
                Nouveau sujet
            </button>
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
                    type="text"
                    placeholder="Rechercher un sujet…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
            <select
                v-model="filterCategory"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Toutes les catégories</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                </option>
            </select>
        </div>

        <div
            v-if="topics.length === 0"
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
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
            </svg>
            <p class="text-gray-400 text-sm">
                Aucun sujet pour le moment. Soyez le premier à lancer une
                discussion !
            </p>
        </div>

        <div v-else class="space-y-3">
            <div
                v-for="topic in filteredTopics"
                :key="topic.id"
                class="bg-white rounded-2xl shadow-sm px-5 py-4 flex items-center gap-4 hover:shadow-md transition-shadow cursor-pointer"
                @click="openTopic(topic)"
            >
                <div
                    class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
                >
                    <svg
                        class="w-5 h-5 text-primary"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="1.8"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                        />
                    </svg>
                </div>
                <div class="flex-1 min-w-0">
                    <h3 class="font-semibold text-gray-800 text-sm truncate">
                        {{ topic.title }}
                    </h3>
                    <p class="text-xs text-gray-400 mt-0.5">
                        Par {{ topic.author }} · {{ topic.created_at }}
                    </p>
                </div>
                <div
                    class="flex items-center gap-4 flex-shrink-0 text-xs text-gray-400"
                >
                    <span class="flex items-center gap-1">
                        <svg
                            class="w-3.5 h-3.5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                            />
                        </svg>
                        {{ topic.replies_count }}
                    </span>
                    <span
                        v-if="topic.category"
                        class="px-2 py-0.5 rounded-full bg-secondary/10 text-secondary font-medium"
                        >{{ topic.category }}</span
                    >
                </div>
            </div>
        </div>

        <div
            v-if="showNewTopic"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="showNewTopic = false"
        >
            <div
                class="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden"
            >
                <div
                    class="flex items-center justify-between px-6 py-4 border-b border-gray-100"
                >
                    <h2
                        class="font-semibold text-gray-800"
                        style="font-family: var(--font-family-title)"
                    >
                        Nouveau sujet
                    </h2>
                    <button
                        @click="showNewTopic = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors"
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
                <div class="px-6 py-5 space-y-4">
                    <div>
                        <label
                            class="block text-xs font-medium text-gray-600 mb-1.5"
                            >Titre</label
                        >
                        <input
                            v-model="newTopic.title"
                            type="text"
                            placeholder="Titre du sujet…"
                            class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                        />
                    </div>
                    <div>
                        <label
                            class="block text-xs font-medium text-gray-600 mb-1.5"
                            >Catégorie</label
                        >
                        <select
                            v-model="newTopic.category_id"
                            class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                        >
                            <option value="">Choisir une catégorie</option>
                            <option
                                v-for="cat in categories"
                                :key="cat.id"
                                :value="cat.id"
                            >
                                {{ cat.name }}
                            </option>
                        </select>
                    </div>
                    <div>
                        <label
                            class="block text-xs font-medium text-gray-600 mb-1.5"
                            >Message</label
                        >
                        <textarea
                            v-model="newTopic.body"
                            rows="5"
                            placeholder="Décrivez votre sujet…"
                            class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                        />
                    </div>
                    <p v-if="topicError" class="text-xs text-red-500">
                        {{ topicError }}
                    </p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex gap-3">
                    <button
                        @click="showNewTopic = false"
                        class="flex-1 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Annuler
                    </button>
                    <button
                        @click="createTopic"
                        :disabled="topicLoading"
                        class="flex-1 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                    >
                        {{ topicLoading ? "Publication…" : "Publier" }}
                    </button>
                </div>
            </div>
        </div>

        <div
            v-if="openedTopic"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="openedTopic = null"
        >
            <div
                class="bg-white rounded-2xl shadow-xl w-full max-w-2xl max-h-[85vh] flex flex-col overflow-hidden"
            >
                <div
                    class="flex items-center justify-between px-6 py-4 border-b border-gray-100 flex-shrink-0"
                >
                    <h2
                        class="font-semibold text-gray-800 text-sm"
                        style="font-family: var(--font-family-title)"
                    >
                        {{ openedTopic.title }}
                    </h2>
                    <button
                        @click="openedTopic = null"
                        class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors"
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
                    <p class="text-xs text-gray-400 text-center">
                        Les réponses seront chargées depuis l'API.
                    </p>
                </div>
                <div
                    class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-3"
                >
                    <input
                        v-model="replyBody"
                        type="text"
                        placeholder="Votre réponse…"
                        class="flex-1 px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                        @keydown.enter="submitReply"
                    />
                    <button
                        @click="submitReply"
                        :disabled="replyLoading || !replyBody.trim()"
                        class="px-4 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                    >
                        Envoyer
                    </button>
                </div>
            </div>
        </div>
    </UserLayout>
</template>

<script setup>
import { ref, computed } from "vue";
import UserLayout from "@/Layouts/UserLayout.vue";
import api from "@/api.js";

const topics = ref([]);
const categories = ref([]);
const search = ref("");
const filterCategory = ref("");
const showNewTopic = ref(false);
const openedTopic = ref(null);
const replyBody = ref("");
const replyLoading = ref(false);
const topicLoading = ref(false);
const topicError = ref("");

const newTopic = ref({ title: "", category_id: "", body: "" });

const filteredTopics = computed(() => {
    return topics.value.filter((t) => {
        const matchSearch =
            !search.value ||
            t.title.toLowerCase().includes(search.value.toLowerCase());
        const matchCat =
            !filterCategory.value || t.category_id === filterCategory.value;
        return matchSearch && matchCat;
    });
});

function openTopic(topic) {
    openedTopic.value = topic;
    replyBody.value = "";
    // Les réponses du sujet seront chargées depuis l'API
    // fetchReplies(topic.id)
}

async function createTopic() {
    if (!newTopic.value.title.trim() || !newTopic.value.body.trim()) {
        topicError.value = "Le titre et le message sont requis.";
        return;
    }
    topicLoading.value = true;
    topicError.value = "";
    try {
        // await api.post('/forum/topics', newTopic.value)
        // fetchTopics()
        showNewTopic.value = false;
        newTopic.value = { title: "", category_id: "", body: "" };
    } catch {
        topicError.value = "Une erreur est survenue.";
    }
    topicLoading.value = false;
}

async function submitReply() {
    if (!replyBody.value.trim()) return;
    replyLoading.value = true;
    try {
        // await api.post(`/forum/topics/${openedTopic.value.id}/replies`, { body: replyBody.value })
        replyBody.value = "";
        // fetchReplies(openedTopic.value.id)
    } catch {}
    replyLoading.value = false;
}

// Les sujets et catégories seront chargés depuis l'API
// onMounted(async () => {
//     const [topicsData, catsData] = await Promise.all([
//         api.get('/forum/topics'),
//         api.get('/forum/categories'),
//     ])
//     topics.value = topicsData.data
//     categories.value = catsData.data
// })
</script>
