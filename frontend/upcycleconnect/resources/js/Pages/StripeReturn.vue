<template>
    <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 p-6">
        <div class="text-center max-w-md">
            <div class="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                </svg>
            </div>
            <h1 class="text-2xl font-bold text-gray-800 mb-2" style="font-family: var(--font-family-title)">
                Compte configuré
            </h1>
            <p class="text-sm text-gray-400 mb-8">
                Votre compte de paiement Stripe a bien été configuré. Vous pouvez maintenant publier des annonces de vente.
            </p>
            <RouterLink
                to="/accueil"
                class="px-5 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                Retour à l'accueil
            </RouterLink>
        </div>
    </div>
</template>
