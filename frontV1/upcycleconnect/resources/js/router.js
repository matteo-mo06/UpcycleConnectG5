import { createRouter, createWebHistory } from 'vue-router'
import Home from './Pages/Home.vue'
import Dashboard from './Pages/Admin/Dashboard.vue'
import Users from './Pages/Admin/Users.vue'
import Roles from './Pages/Admin/Roles.vue'
import Listings from './Pages/Admin/Listings.vue'
import Events from './Pages/Admin/Events.vue'
import Reports from './Pages/Admin/Reports.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/admin/dashboard', component: Dashboard },
    { path: '/admin/users', component: Users },
    { path: '/admin/roles', component: Roles },
    { path: '/admin/listings', component: Listings },
    { path: '/admin/events', component: Events },
    { path: '/admin/reports', component: Reports },
  ],
})
