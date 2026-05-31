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
                    type="text"
                    placeholder="Rechercher un conseil…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
            <select
                v-model="filterTag"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Tous les thèmes</option>
                <option v-for="tag in tags" :key="tag" :value="tag">
                    {{ tag }}
                </option>
            </select>
        </div>

        <div
            v-if="conseils.length === 0"
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
                v-for="conseil in filteredConseils"
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
                    {{ conseil.summary }}
                </p>
                <div
                    class="flex items-center gap-2 mt-auto pt-2 border-t border-gray-50 flex-wrap"
                >
                    <span
                        v-for="tag in conseil.tags"
                        :key="tag"
                        class="px-2 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-medium"
                    >
                        {{ tag }}
                    </span>
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
                    <h2
                        class="font-semibold text-gray-800"
                        style="font-family: var(--font-family-title)"
                    >
                        {{ selected.title }}
                    </h2>
                    <button
                        @click="selected = null"
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
                    <div class="flex flex-wrap gap-2">
                        <span
                            v-for="tag in selected.tags"
                            :key="tag"
                            class="px-2 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-medium"
                        >
                            {{ tag }}
                        </span>
                    </div>
                    <p class="text-sm text-gray-600 whitespace-pre-line">
                        {{ selected.content }}
                    </p>
                    <p class="text-xs text-gray-400 text-center py-4">
                        Le contenu complet sera chargé depuis l'API.
                    </p>
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
import { ref, computed } from "vue";
import UserLayout from "@/Layouts/UserLayout.vue";
import api from "@/api.js";

const conseils = ref([]);
const tags = ref([]);
const search = ref("");
const filterTag = ref("");
const selected = ref(null);

const filteredConseils = computed(() => {
    return conseils.value.filter((c) => {
        const matchSearch =
            !search.value ||
            c.title.toLowerCase().includes(search.value.toLowerCase());
        const matchTag =
            !filterTag.value || (c.tags ?? []).includes(filterTag.value);
        return matchSearch && matchTag;
    });
});

function openDetail(conseil) {
    selected.value = conseil;
}

// Les conseils et thèmes seront chargés depuis l'API
// onMounted(async () => {
//     const { data } = await api.get('/conseils')
//     conseils.value = data
//     tags.value = [...new Set(data.flatMap(c => c.tags ?? []))]
// })
</script>
