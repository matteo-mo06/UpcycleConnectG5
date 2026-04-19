import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import Login    from './Pages/Login.vue'
import Home     from './Pages/Home.vue'
import Accueil  from './Pages/User/Accueil.vue'
import Annonces from './Pages/User/Annonces.vue'
import Depot    from './Pages/User/Depot.vue'
import Dashboard from './Pages/Admin/Dashboard.vue'
import Users     from './Pages/Admin/Users.vue'
import Roles     from './Pages/Admin/Roles.vue'
import Listings  from './Pages/Admin/Listings.vue'
import Events    from './Pages/Admin/Events.vue'
import Reports   from './Pages/Admin/Reports.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',      component: Home  },
    { path: '/login', component: Login },

    { path: '/accueil',  component: Accueil,  meta: { requiresAuth: true } },
    { path: '/annonces', component: Annonces, meta: { requiresAuth: true } },
    { path: '/depot',    component: Depot,    meta: { requiresAuth: true } },

    { path: '/admin/dashboard', component: Dashboard, meta: { requiresAdmin: true } },
    { path: '/admin/users',     component: Users,     meta: { requiresAdmin: true } },
    { path: '/admin/roles',     component: Roles,     meta: { requiresAdmin: true } },
    { path: '/admin/listings',  component: Listings,  meta: { requiresAdmin: true } },
    { path: '/admin/events',    component: Events,    meta: { requiresAdmin: true } },
    { path: '/admin/reports',   component: Reports,   meta: { requiresAdmin: true } },

    { path: '/:pathMatch(.*)*', redirect: '/accueil' },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAdmin && !auth.isAdmin) return '/login'
  if (to.meta.requiresAuth && !auth.isLoggedIn) return '/login'
})

export default router
