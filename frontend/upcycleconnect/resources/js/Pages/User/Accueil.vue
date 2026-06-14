<template>
    <UserLayout @openDepot="showDepot = true">

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                Espace particulier
            </h1>
        </div>

        <!-- Annonces mises en avant -->
        <div class="mb-6">
            <div class="flex items-center justify-between mb-3">
                <div>
                    <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">Dernières annonces</h2>
                    <p class="text-xs text-gray-400 mt-0.5">Les annonces récemment publiées sur la plateforme</p>
                </div>
                <RouterLink to="/annonces" class="text-xs font-medium text-primary hover:underline">Voir tout</RouterLink>
            </div>

            <div v-if="announcesLoading" class="grid grid-cols-4 gap-4">
                <div v-for="i in 4" :key="i" class="bg-white rounded-2xl shadow-sm h-48 animate-pulse" />
            </div>

            <div v-else-if="featuredAnnounces.length === 0" class="bg-white rounded-2xl shadow-sm p-10 text-center text-sm text-gray-400">
                Aucune annonce disponible pour le moment.
            </div>

            <div v-else class="grid grid-cols-4 gap-4">
                <RouterLink
                    v-for="ann in featuredAnnounces"
                    :key="ann.id"
                    :to="`/annonces/${ann.id}`"
                    class="bg-white rounded-2xl shadow-sm overflow-hidden hover:shadow-md transition-shadow group"
                >
                    <div class="h-32 bg-gray-100 overflow-hidden">
                        <img v-if="ann.first_photo" :src="ann.first_photo" :alt="ann.title"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                        <div v-else class="w-full h-full flex items-center justify-center text-gray-300">
                            <svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                            </svg>
                        </div>
                    </div>
                    <div class="p-3">
                        <p class="text-sm font-medium text-gray-800 truncate">{{ ann.title }}</p>
                        <p class="text-xs text-gray-400 mt-0.5 truncate">{{ ann.city || ann.address }}</p>
                        <div class="flex items-center justify-between mt-2">
                            <span class="text-xs font-semibold text-primary">
                                {{ ann.price > 0 ? ann.price.toFixed(2) + ' €' : 'Gratuit' }}
                            </span>
                            <span class="px-1.5 py-0.5 rounded text-xs font-medium"
                                :class="ann.type === 'don' ? 'bg-secondary/10 text-secondary' : 'bg-primary/10 text-primary'">
                                {{ ann.type === 'don' ? 'Don' : 'Vente' }}
                            </span>
                        </div>
                    </div>
                </RouterLink>
            </div>
        </div>

        <div class="grid grid-cols-5 gap-6">

            <!-- Actions rapides -->
            <div class="col-span-2 bg-white rounded-2xl shadow-sm p-5" data-tutorial="quick-actions">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Actions rapides</h2>
                <p class="text-xs text-gray-400 mb-4">Accès direct à vos actions</p>
                <div class="grid grid-cols-2 gap-3">
                    <button
                        @click="showDepot = true"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-primary text-white hover:bg-primary-dark transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                        </svg>
                        <span class="text-sm font-medium text-center">Déposer un objet</span>
                    </button>

                    <button
                        @click="showDepot = true"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-primary text-primary hover:bg-primary/5 transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                        </svg>
                        <span class="text-sm font-medium text-center">+ Publier une annonce</span>
                    </button>

                    <RouterLink
                        to="/formations"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-secondary text-white hover:bg-secondary-dark transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                        </svg>
                        <span class="text-sm font-medium text-center">+ Voir les formations</span>
                    </RouterLink>

                    <RouterLink
                        to="/projets"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-secondary text-secondary hover:bg-secondary/5 transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                        </svg>
                        <span class="text-sm font-medium text-center">+ Rejoindre un projet</span>
                    </RouterLink>
                </div>
            </div>

            <!-- Vitrine partenaires -->
            <div class="col-span-3 bg-white rounded-2xl shadow-sm p-5 flex flex-col">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">Partenaires & publicités</h2>
                <p class="text-xs text-gray-400 mb-4">Ils soutiennent UpcycleConnect</p>
                <div class="flex-1 flex flex-col items-center justify-center gap-3 py-8 border-2 border-dashed border-gray-200 rounded-xl">
                    <svg class="w-10 h-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-2 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
                    </svg>
                    <p class="text-sm font-medium text-gray-400">Aucun partenaire pour le moment</p>
                    <p class="text-xs text-gray-300 text-center max-w-xs">Cet espace est réservé à nos partenaires et annonceurs. Contactez-nous pour en savoir plus.</p>
                </div>
            </div>

        </div>

        <CreateAnnouncementModal v-model="showDepot" />

        <OnboardingTutorial
            v-if="showTutorial"
            :steps="tutorialSteps"
            @done="finishTutorial"
        />

    </UserLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import UserLayout from '@/Layouts/UserLayout.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import OnboardingTutorial from '@/Components/OnboardingTutorial.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const auth = useAuthStore()
const showDepot = ref(false)
const showTutorial = ref(false)
const featuredAnnounces = ref([])
const announcesLoading = ref(true)

const tutorialSteps = [
    {
        target: 'sidebar',
        placement: 'right',
        title: 'Bienvenue sur UpcycleConnect !',
        description: 'Ce menu latéral vous donne accès à toutes les sections de la plateforme. Faisons un rapide tour ensemble.',
    },
    {
        target: 'nav-annonces',
        placement: 'right',
        title: 'Annonces',
        description: 'Parcourez les objets disponibles à l\'adoption ou à la vente, et publiez vos propres annonces pour donner une seconde vie à vos objets.',
    },
    {
        target: 'nav-depot',
        placement: 'right',
        title: 'Dépôt d\'objet',
        description: 'Une fois votre objet vendu, demandez l\'attribution d\'un casier. L\'acheteur recevra un code d\'accès pour venir le récupérer.',
    },
    {
        target: 'nav-projets',
        placement: 'right',
        title: 'Projets collaboratifs',
        description: 'Rejoignez ou créez des projets de création et de réparation avec d\'autres membres de la communauté.',
    },
    {
        target: 'nav-forum',
        placement: 'right',
        title: 'Forum',
        description: 'Échangez avec la communauté, posez vos questions, partagez vos astuces et participez aux discussions autour de l\'upcycling.',
    },
    {
        target: 'nav-formations',
        placement: 'right',
        title: 'Formations',
        description: 'Accédez à des formations pratiques pour apprendre à réparer, transformer et revaloriser vos objets.',
    },
    {
        target: 'nav-evenements',
        placement: 'right',
        title: 'Événements',
        description: 'Découvrez et inscrivez-vous aux ateliers, marchés et rencontres organisés près de chez vous.',
    },
    {
        target: 'nav-conseils',
        placement: 'right',
        title: 'Conseils',
        description: 'Retrouvez des guides et astuces pour prendre soin de vos objets, les réparer et leur donner une nouvelle vie.',
    },
    {
        target: 'quick-actions',
        placement: 'left',
        title: 'Actions rapides',
        description: 'Depuis votre accueil, accédez directement aux actions les plus courantes sans passer par les menus.',
    },
]

async function finishTutorial() {
    showTutorial.value = false
    try {
        await api.post('/user/tutorial-done')
        auth.setTutorialDone()
    } catch {}
}

onMounted(async () => {
    if (!auth.user?.tutorial_done) {
        showTutorial.value = true
    }
    try {
        const { data } = await api.get('/announcements?limit=4')
        featuredAnnounces.value = (data.data ?? []).slice(0, 4)
    } catch {}
    announcesLoading.value = false
})
</script>
