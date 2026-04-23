<template>
    <UserLayout>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1
                    class="text-3xl font-bold text-gray-800"
                    style="font-family: var(--font-family-title)"
                >
                    Projets
                </h1>
                <p class="text-sm text-gray-400 mt-1">
                    Rejoignez des projets collaboratifs d'upcycling
                </p>
            </div>
            <button
                disabled
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium opacity-50 cursor-not-allowed"
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
                Créer un projet
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
                    placeholder="Rechercher un projet…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
            <select
                v-model="filterStatus"
                class="px-3 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
                <option value="">Tous les statuts</option>
                <option value="open">Ouvert</option>
                <option value="in_progress">En cours</option>
                <option value="completed">Terminé</option>
            </select>
        </div>

        <div
            v-if="projects.length === 0"
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
                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"
                />
            </svg>
            <p class="text-gray-400 text-sm">
                Aucun projet disponible pour le moment.
            </p>
        </div>

        <div v-else class="grid grid-cols-3 gap-5">
            <div
                v-for="project in filteredProjects"
                :key="project.id"
                class="bg-white rounded-2xl shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow cursor-pointer"
                @click="openDetail(project)"
            >
                <div class="flex items-start justify-between gap-2">
                    <h3
                        class="font-semibold text-gray-800 text-sm leading-snug"
                    >
                        {{ project.title }}
                    </h3>
                    <span
                        :class="statusClass(project.status)"
                        class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0"
                    >
                        {{ statusLabel(project.status) }}
                    </span>
                </div>
                <p class="text-xs text-gray-500 line-clamp-2">
                    {{ project.description }}
                </p>
                <div
                    class="flex items-center justify-between mt-auto pt-2 border-t border-gray-50"
                >
                    <span class="text-xs text-gray-400"
                        >{{ project.members_count }} membre(s)</span
                    >
                    <button
                        @click.stop="joinProject(project)"
                        :disabled="project.loading"
                        class="px-3 py-1 rounded-lg bg-secondary text-white text-xs font-medium hover:bg-secondary-dark transition-colors disabled:opacity-50"
                    >
                        Rejoindre
                    </button>
                </div>
            </div>
        </div>

        <div
            v-if="selected"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="selected = null"
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
                <div class="px-6 py-5 space-y-4">
                    <p class="text-sm text-gray-600">
                        {{ selected.description }}
                    </p>
                    <div class="flex items-center gap-4 text-xs text-gray-400">
                        <span>{{ selected.members_count }} membre(s)</span>
                        <span>Créé par {{ selected.creator }}</span>
                    </div>
                </div>
                <div
                    class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3"
                >
                    <button
                        @click="selected = null"
                        class="px-4 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Fermer
                    </button>
                    <button
                        @click="joinProject(selected)"
                        class="px-4 py-2 rounded-xl bg-secondary text-white text-sm font-medium hover:bg-secondary-dark transition-colors"
                    >
                        Rejoindre le projet
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

const projects = ref([]);
const search = ref("");
const filterStatus = ref("");
const selected = ref(null);

const filteredProjects = computed(() => {
    return projects.value.filter((p) => {
        const matchSearch =
            !search.value ||
            p.title.toLowerCase().includes(search.value.toLowerCase());
        const matchStatus =
            !filterStatus.value || p.status === filterStatus.value;
        return matchSearch && matchStatus;
    });
});

function statusLabel(status) {
    return (
        { open: "Ouvert", in_progress: "En cours", completed: "Terminé" }[
            status
        ] ?? status
    );
}

function statusClass(status) {
    return (
        {
            open: "bg-secondary/10 text-secondary",
            in_progress: "bg-primary/10 text-primary",
            completed: "bg-gray-100 text-gray-500",
        }[status] ?? "bg-gray-100 text-gray-500"
    );
}

function openDetail(project) {
    selected.value = project;
}

async function joinProject(project) {
    // Rejoindre un projet via l'API
    project.loading = true;
    try {
        await api.post(`/projects/${project.id}/join`);
    } catch {}
    project.loading = false;
}

// Les projets seront chargés depuis l'API
// onMounted(async () => {
//     const { data } = await api.get('/projects')
//     projects.value = data
// })
</script>
