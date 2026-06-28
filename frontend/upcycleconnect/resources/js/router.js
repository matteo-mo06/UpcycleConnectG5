import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import Login from './Pages/Login.vue'
import Accueil from './Pages/User/Accueil.vue'
import Annonces from './Pages/User/Annonces.vue'
import Depot    from './Pages/User/Depot.vue'
import UserEvents from './Pages/User/Events.vue'
import Projets from './Pages/User/Projets.vue'
import Forum from './Pages/User/Forum.vue'
import Formations from './Pages/User/Formations.vue'
import Conseils from './Pages/User/Conseils.vue'
import Parametres from './Pages/User/Parametres.vue'
import Score from './Pages/User/Score.vue'
import Calendrier from './Pages/User/Calendrier.vue'
import Dashboard from './Pages/Admin/Dashboard.vue'
import Users from './Pages/Admin/Users.vue'
import Roles from './Pages/Admin/Roles.vue'
import Listings from './Pages/Admin/Listings.vue'
import Events from './Pages/Admin/Events.vue'
import AdminFormations from './Pages/Admin/Formations.vue'
import AdminScore from './Pages/Admin/Score.vue'
import Reports from './Pages/Admin/Reports.vue'
import Lockers from './Pages/Admin/Lockers.vue'
import AdminProjets from './Pages/Admin/Projets.vue'
import Categories from './Pages/Admin/Categories.vue'
import DemandesPro from './Pages/Admin/DemandesPro.vue'
import Revenus from './Pages/Admin/Revenus.vue'
import AdminPublicite from './Pages/Admin/Publicite.vue'
import AdminSubscriptions from './Pages/Admin/Subscriptions.vue'
import Paiement from './Pages/User/Paiement.vue'
import PaiementConfirmation from './Pages/User/PaiementConfirmation.vue'
import PaiementFormation from './Pages/User/PaiementFormation.vue'
import CGU from './Pages/CGU.vue'
import PolitiqueConfidentialite from './Pages/PolitiqueConfidentialite.vue'
import Error from './Pages/Error.vue'
import StripeReturn from './Pages/StripeReturn.vue'
import StripeRefresh from './Pages/StripeRefresh.vue'

const ArtisanDashboard  = () => import('./Pages/Artisan/Dashboard.vue')
const ArtisanAnnonces   = () => import('./Pages/Artisan/Annonces.vue')
const ArtisanDepot      = () => import('./Pages/Artisan/Depot.vue')
const ArtisanEvenements = () => import('./Pages/Artisan/Evenements.vue')
const ArtisanFormations = () => import('./Pages/Artisan/Formations.vue')
const ArtisanProjets    = () => import('./Pages/Artisan/Projets.vue')
const ArtisanCalendrier = () => import('./Pages/Artisan/Calendrier.vue')
const ArtisanScore      = () => import('./Pages/Artisan/Score.vue')
const ArtisanAbonnement  = () => import('./Pages/Artisan/Abonnement.vue')
const ArtisanPublicites  = () => import('./Pages/Artisan/Publicites.vue')

const SalarieDashboard  = () => import('./Pages/Salarie/Dashboard.vue')
const SalarieFormations = () => import('./Pages/Salarie/Formations.vue')
const SalarieEvenements = () => import('./Pages/Salarie/Evenements.vue')
const SalarieConseils   = () => import('./Pages/Salarie/Conseils.vue')
const SalarieForum      = () => import('./Pages/Salarie/Forum.vue')
const SalarieCalendrier = () => import('./Pages/Salarie/Calendrier.vue')
const SalarieAbonnement = () => import('./Pages/Salarie/Abonnement.vue')

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/', redirect: '/accueil' },
        { path: '/login', component: Login },
        {
            path: '/abonnement',
            redirect: to => {
                const auth = useAuthStore()
                if (auth.isArtisan) return { path: '/artisan/abonnement', query: to.query }
                if (auth.isSalarie) return { path: '/salarie/abonnement', query: to.query }
                return '/accueil'
            },
        },

        { path: '/cgu', component: CGU },
        { path: '/politique-confidentialite', component: PolitiqueConfidentialite },

        { path: '/accueil',     component: Accueil,     meta: { requiresAuth: true } },
        { path: '/annonces',    component: Annonces,    meta: { requiresAuth: true } },
        { path: '/depot',       component: Depot,       meta: { requiresAuth: true } },
        { path: '/evenements',  component: UserEvents,  meta: { requiresAuth: true } },
        { path: '/projets',     component: Projets,     meta: { requiresAuth: true } },
        { path: '/forum',       component: Forum,       meta: { requiresAuth: true } },
        { path: '/formations',  component: Formations,  meta: { requiresAuth: true } },
        { path: '/conseils',    component: Conseils,    meta: { requiresAuth: true } },
        { path: '/parametres',  component: Parametres,  meta: { requiresAuth: true } },
        { path: '/score',       component: Score,       meta: { requiresAuth: true } },
        { path: '/calendrier',       component: Calendrier,            meta: { requiresAuth: true } },
        { path: '/paiement/:id',             component: Paiement,           meta: { requiresAuth: true } },
        { path: '/paiement-formation/:id',   component: PaiementFormation,  meta: { requiresAuth: true } },
        { path: '/paiement-confirmation',    component: PaiementConfirmation, meta: { requiresAuth: true } },

        { path: '/admin/dashboard', component: Dashboard, meta: { requiresAdmin: true } },
        { path: '/admin/users', component: Users, meta: { requiresAdmin: true } },
        { path: '/admin/roles', component: Roles, meta: { requiresAdmin: true } },
        { path: '/admin/annonces', component: Listings, meta: { requiresAdmin: true } },
        { path: '/admin/events', component: Events, meta: { requiresAdmin: true } },
        { path: '/admin/formations', component: AdminFormations, meta: { requiresAdmin: true } },
        { path: '/admin/reports', component: Reports, meta: { requiresAdmin: true } },
        { path: '/admin/lockers', component: Lockers, meta: { requiresAdmin: true } },
        { path: '/admin/score', component: AdminScore, meta: { requiresAdmin: true } },
        { path: '/admin/projets', component: AdminProjets, meta: { requiresAdmin: true } },
        { path: '/admin/categories', component: Categories, meta: { requiresAdmin: true } },
        { path: '/admin/demandes-pro', component: DemandesPro, meta: { requiresAdmin: true } },
        { path: '/admin/revenus', component: Revenus, meta: { requiresAdmin: true } },
        { path: '/admin/abonnements', component: AdminSubscriptions, meta: { requiresAdmin: true } },
        { path: '/admin/publicite', component: AdminPublicite, meta: { requiresAdmin: true } },

        { path: '/artisan/dashboard',  component: ArtisanDashboard,  meta: { requiresArtisan: true } },
        { path: '/artisan/annonces',   component: ArtisanAnnonces,   meta: { requiresArtisan: true } },
        { path: '/artisan/depot',      component: ArtisanDepot,      meta: { requiresArtisan: true } },
        { path: '/artisan/evenements', component: ArtisanEvenements, meta: { requiresArtisan: true } },
        { path: '/artisan/formations', component: ArtisanFormations, meta: { requiresArtisan: true } },
        { path: '/artisan/projets',    component: ArtisanProjets,    meta: { requiresArtisan: true } },
        { path: '/artisan/calendrier', component: ArtisanCalendrier, meta: { requiresArtisan: true } },
        { path: '/artisan/score',      component: ArtisanScore,      meta: { requiresArtisan: true } },
        { path: '/artisan/abonnement',  component: ArtisanAbonnement,  meta: { requiresArtisan: true } },
        { path: '/artisan/publicites',  component: ArtisanPublicites,  meta: { requiresArtisan: true } },

        { path: '/salarie/dashboard',  component: SalarieDashboard,  meta: { requiresSalarie: true } },
        { path: '/salarie/formations', component: SalarieFormations, meta: { requiresSalarie: true } },
        { path: '/salarie/evenements', component: SalarieEvenements, meta: { requiresSalarie: true } },
        { path: '/salarie/conseils',   component: SalarieConseils,   meta: { requiresSalarie: true } },
        { path: '/salarie/forum',      component: SalarieForum,      meta: { requiresSalarie: true } },
        { path: '/salarie/calendrier', component: SalarieCalendrier, meta: { requiresSalarie: true } },
        { path: '/salarie/abonnement', component: SalarieAbonnement, meta: { requiresSalarie: true } },

        { path: '/stripe/return',  component: StripeReturn,  meta: { requiresAuth: true } },
        { path: '/stripe/refresh', component: StripeRefresh, meta: { requiresAuth: true } },
        { path: '/error/:code', component: Error },
        { path: '/:pathMatch(.*)*', component: Error },
    ],
})

router.beforeEach((to) => {
    const auth = useAuthStore()

    if (to.meta.requiresAdmin && !auth.isAdmin) return auth.isLoggedIn ? '/accueil' : '/login'
    if (to.meta.requiresArtisan && !auth.isArtisan) return auth.isLoggedIn ? '/accueil' : '/login'
    if (to.meta.requiresSalarie && !auth.isSalarie) return auth.isLoggedIn ? '/accueil' : '/login'
    if (to.meta.requiresAuth && !auth.isLoggedIn) return '/login'

    if (auth.isLoggedIn && !auth.user?.tutorial_done && to.path.startsWith('/artisan/')) {
        return '/accueil'
    }
})

export default router
