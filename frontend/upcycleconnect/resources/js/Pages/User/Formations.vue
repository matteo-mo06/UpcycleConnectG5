<template>
    <UserLayout>
        <div class="mb-6">
            <h1
                class="text-3xl font-bold text-gray-800"
                style="font-family: var(--font-family-title)"
            >
                Formations
            </h1>
            <p class="text-sm text-gray-400 mt-1">
                Apprenez les techniques d'upcycling avec nos formations
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
                    placeholder="Rechercher une formation…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
            <select
                v-model="filterLevel"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Tous les niveaux</option>
                <option value="beginner">Débutant</option>
                <option value="intermediate">Intermédiaire</option>
                <option value="advanced">Avancé</option>
            </select>
            <select
                v-model="filterCategory"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Toutes les thématiques</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                </option>
            </select>
        </div>

        <div
            v-if="formations.length === 0"
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
                    d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
            </svg>
            <p class="text-gray-400 text-sm">
                Aucune formation disponible pour le moment.
            </p>
        </div>

        <div v-else class="grid grid-cols-3 gap-5">
            <div
                v-for="formation in filteredFormations"
                :key="formation.id"
                class="bg-white rounded-2xl shadow-sm overflow-hidden hover:shadow-md transition-shadow cursor-pointer flex flex-col"
                @click="openDetail(formation)"
            >
                <div
                    class="h-36 bg-secondary/10 flex items-center justify-center"
                >
                    <svg
                        class="w-10 h-10 text-secondary/40"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="1.5"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M15 10l4.553-2.069A1 1 0 0121 8.87v6.26a1 1 0 01-1.447.894L15 14M3 8a2 2 0 012-2h8a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V8z"
                        />
                    </svg>
                </div>
                <div class="p-4 flex flex-col gap-2 flex-1">
                    <div class="flex items-center gap-2">
                        <span
                            :class="levelClass(formation.level)"
                            class="px-2 py-0.5 rounded-full text-xs font-medium"
                        >
                            {{ levelLabel(formation.level) }}
                        </span>
                        <span class="text-xs text-gray-400">{{
                            formation.duration
                        }}</span>
                    </div>
                    <h3
                        class="font-semibold text-gray-800 text-sm leading-snug"
                    >
                        {{ formation.title }}
                    </h3>
                    <p class="text-xs text-gray-500 line-clamp-2">
                        {{ formation.description }}
                    </p>
                    <div class="mt-auto pt-2 flex items-center justify-between">
                        <span class="text-xs text-gray-400"
                            >Par {{ formation.author }}</span
                        >
                        <span
                            v-if="formation.enrolled"
                            class="text-xs text-secondary font-medium"
                            >Inscrit</span
                        >
                    </div>
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
                    <div class="flex items-center gap-3 text-sm text-gray-500">
                        <span
                            :class="levelClass(selected.level)"
                            class="px-2 py-0.5 rounded-full text-xs font-medium"
                            >{{ levelLabel(selected.level) }}</span
                        >
                        <span>{{ selected.duration }}</span>
                        <span>Par {{ selected.author }}</span>
                    </div>
                    <p class="text-sm text-gray-600">
                        {{ selected.description }}
                    </p>
                    <p class="text-xs text-gray-400 text-center py-4">
                        Le contenu de la formation sera chargé depuis l'API.
                    </p>
                </div>
                <div
                    class="px-6 py-4 border-t border-gray-100 flex-shrink-0 flex gap-3"
                >
                    <button
                        @click="selected = null"
                        class="flex-1 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Fermer
                    </button>
                    <button
                        @click="enroll(selected)"
                        :disabled="selected.enrolled || selected.loading"
                        class="flex-1 py-2 rounded-xl bg-secondary text-white text-sm font-medium hover:bg-secondary-dark transition-colors disabled:opacity-50"
                    >
                        {{ selected.enrolled ? "Déjà inscrit" : "S'inscrire" }}
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

const formations = ref([]);
const categories = ref([]);
const search = ref("");
const filterLevel = ref("");
const filterCategory = ref("");
const selected = ref(null);

const filteredFormations = computed(() => {
    return formations.value.filter((f) => {
        const matchSearch =
            !search.value ||
            f.title.toLowerCase().includes(search.value.toLowerCase());
        const matchLevel = !filterLevel.value || f.level === filterLevel.value;
        const matchCat =
            !filterCategory.value || f.category_id === filterCategory.value;
        return matchSearch && matchLevel && matchCat;
    });
});

function levelLabel(level) {
    return (
        {
            beginner: "Débutant",
            intermediate: "Intermédiaire",
            advanced: "Avancé",
        }[level] ?? level
    );
}

function levelClass(level) {
    return (
        {
            beginner: "bg-secondary/10 text-secondary",
            intermediate: "bg-primary/10 text-primary",
            advanced: "bg-red-100 text-red-500",
        }[level] ?? "bg-gray-100 text-gray-500"
    );
}

function openDetail(formation) {
    selected.value = formation;
}

async function enroll(formation) {
    formation.loading = true;
    try {
        // await api.post(`/formations/${formation.id}/enroll`)
        formation.enrolled = true;
    } catch {}
    formation.loading = false;
}

// Les formations et thématiques seront chargées depuis l'API
// onMounted(async () => {
//     const [formationsData, catsData] = await Promise.all([
//         api.get('/formations'),
//         api.get('/formations/categories'),
//     ])
//     formations.value = formationsData.data
//     categories.value = catsData.data
// })
</script>
