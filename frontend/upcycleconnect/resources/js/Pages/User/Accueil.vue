<template>
    <UserLayout @openDepot="showDepot = true">

        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-800" style="font-family: var(--font-family-title)">
                {{ t('accueil.title') }}
            </h1>
        </div>

        <div class="mb-6">
            <div class="flex items-center justify-between mb-3">
                <div>
                    <h2 class="font-semibold text-gray-800" style="font-family: var(--font-family-title)">
                        {{ t('accueil.announcements.sectionTitle') }}
                    </h2>
                    <p class="text-xs text-gray-400 mt-0.5">{{ t('accueil.announcements.sectionSubtitle') }}</p>
                </div>
                <RouterLink to="/annonces" class="text-xs font-medium text-primary hover:underline">
                    {{ t('accueil.announcements.seeAll') }}
                </RouterLink>
            </div>

            <div v-if="announcesLoading" class="grid grid-cols-4 gap-4">
                <div v-for="i in 4" :key="i" class="bg-white rounded-2xl shadow-sm h-48 animate-pulse" />
            </div>

            <div v-else-if="featuredAnnounces.length === 0" class="bg-white rounded-2xl shadow-sm p-10 text-center text-sm text-gray-400">
                {{ t('accueil.announcements.empty') }}
            </div>

            <div v-else class="grid grid-cols-4 gap-4">
                <RouterLink
                    v-for="ann in featuredAnnounces"
                    :key="ann.id"
                    :to="`/annonces?highlight=${ann.id}`"
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
                                {{ ann.price > 0 ? ann.price.toFixed(2) + ' €' : t('accueil.announcements.free') }}
                            </span>
                            <span class="px-1.5 py-0.5 rounded text-xs font-medium"
                                :class="ann.type === 'don' ? 'bg-secondary/10 text-secondary' : 'bg-primary/10 text-primary'">
                                {{ ann.type === 'don' ? t('accueil.announcements.gift') : t('accueil.announcements.sale') }}
                            </span>
                        </div>
                    </div>
                </RouterLink>
            </div>
        </div>

        <div class="grid grid-cols-5 gap-6">

            <div class="col-span-2 bg-white rounded-2xl shadow-sm p-5" data-tutorial="quick-actions">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">
                    {{ t('accueil.actions.sectionTitle') }}
                </h2>
                <p class="text-xs text-gray-400 mb-4">{{ t('accueil.actions.sectionSubtitle') }}</p>
                <div class="grid grid-cols-2 gap-3">
                    <button
                        @click="router.push('/annonces?tab=mine')"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-primary text-white hover:bg-primary-dark transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                        </svg>
                        <span class="text-sm font-medium text-center">{{ t('accueil.actions.deposit') }}</span>
                    </button>

                    <button
                        @click="router.push('/annonces?publish=1')"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-primary text-primary hover:bg-primary/5 transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                        </svg>
                        <span class="text-sm font-medium text-center">{{ t('accueil.actions.publish') }}</span>
                    </button>

                    <RouterLink
                        to="/formations"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl bg-secondary text-white hover:bg-secondary-dark transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                        </svg>
                        <span class="text-sm font-medium text-center">{{ t('accueil.actions.formations') }}</span>
                    </RouterLink>

                    <RouterLink
                        to="/projets"
                        class="flex flex-col items-center justify-center gap-2 p-4 rounded-xl border-2 border-secondary text-secondary hover:bg-secondary/5 transition-colors min-h-24"
                    >
                        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                        </svg>
                        <span class="text-sm font-medium text-center">{{ t('accueil.actions.projects') }}</span>
                    </RouterLink>
                </div>
            </div>

            <div class="col-span-3 bg-white rounded-2xl shadow-sm p-5 flex flex-col">
                <h2 class="font-semibold text-gray-800 mb-1" style="font-family: var(--font-family-title)">
                    {{ t('accueil.partners.sectionTitle') }}
                </h2>
                <p class="text-xs text-gray-400 mb-4">{{ t('accueil.partners.sectionSubtitle') }}</p>
                <div class="flex-1 flex flex-col items-center justify-center gap-3 py-8 border-2 border-dashed border-gray-200 rounded-xl">
                    <svg class="w-10 h-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-2 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
                    </svg>
                    <p class="text-sm font-medium text-gray-400">{{ t('accueil.partners.empty') }}</p>
                    <p class="text-xs text-gray-300 text-center max-w-xs">{{ t('accueil.partners.emptyDetail') }}</p>
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
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import UserLayout from '@/Layouts/UserLayout.vue'
import CreateAnnouncementModal from '@/Components/CreateAnnouncementModal.vue'
import OnboardingTutorial from '@/Components/OnboardingTutorial.vue'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const { t, tm, rt } = useI18n()

const auth = useAuthStore()
const router = useRouter()
const showDepot = ref(false)
const showTutorial = ref(false)
const featuredAnnounces = ref([])
const announcesLoading = ref(true)

const TUTORIAL_TARGETS = ['sidebar', 'nav-annonces', 'nav-depot', 'nav-projets', 'nav-forum', 'nav-formations', 'nav-evenements', 'nav-calendrier', 'nav-conseils', 'quick-actions']

const tutorialSteps = computed(() =>
    tm('accueil.tutorial.steps').map((step, i) => ({
        target: TUTORIAL_TARGETS[i],
        placement: 'right',
        title: rt(step.title),
        description: rt(step.description),
    }))
)

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
