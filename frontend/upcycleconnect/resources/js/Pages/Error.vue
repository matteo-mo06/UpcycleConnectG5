<template>
    <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 p-6">
        <div class="text-center max-w-md">
            <p class="text-8xl font-bold text-primary/20 mb-2">{{ code }}</p>
            <h1
                class="text-2xl font-bold text-gray-800 mb-2"
                style="font-family: var(--font-family-title)"
            >
                {{ title }}
            </h1>
            <p class="text-sm text-gray-400 mb-8">{{ description }}</p>
            <div class="flex gap-3 justify-center">
                <button
                    @click="$router.back()"
                    class="px-5 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-100 transition-colors"
                >
                    Retour
                </button>
                <RouterLink
                    to="/accueil"
                    class="px-5 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors"
                >
                    Accueil
                </RouterLink>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

const code = computed(() => route.params.code ?? "404");

const title = computed(() => {
    switch (code.value) {
        case "403": return "Accès refusé";
        case "404": return "Page introuvable";
        case "500": return "Erreur serveur";
        default:    return "Une erreur est survenue";
    }
});

const description = computed(() => {
    switch (code.value) {
        case "403": return "Vous n'avez pas les droits pour accéder à cette page.";
        case "404": return "La page que vous cherchez n'existe pas ou a été déplacée.";
        case "500": return "Le serveur a rencontré une erreur inattendue. Réessayez dans quelques instants.";
        default:    return "Quelque chose s'est mal passé.";
    }
});
</script>
