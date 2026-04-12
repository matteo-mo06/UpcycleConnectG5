<template>
  <div class="min-h-screen bg-bg flex items-center justify-center p-4">

    <div class="bg-white rounded-2xl shadow-md w-full max-w-sm px-8 py-10">

      <!-- Logo + titre -->
      <div class="flex flex-col items-center mb-7">
        <div class="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center mb-3">
          <svg class="w-8 h-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
        </div>
        <h1 class="text-xl font-bold text-primary" style="font-family: var(--font-family-title)">
          UpcycleConnect
        </h1>
        <p class="text-xs text-gray-400 tracking-widest uppercase mt-1">
          Donnez une seconde vie
        </p>
      </div>

      <!-- Tabs -->
      <div class="flex bg-[#f0ebe4] rounded-xl p-1 mb-6">
        <button
          @click="tab = 'login'"
          :class="[
            'flex-1 py-2 rounded-lg text-sm font-medium transition-all duration-150',
            tab === 'login'
              ? 'bg-primary text-white shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          ]">
          Connexion
        </button>
        <button
          @click="tab = 'register'"
          :class="[
            'flex-1 py-2 rounded-lg text-sm font-medium transition-all duration-150',
            tab === 'register'
              ? 'bg-primary text-white shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          ]">
          Inscription
        </button>
      </div>

      <!-- Formulaire Connexion -->
      <form v-if="tab === 'login'" @submit.prevent="submitLogin" class="space-y-4">

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            Adresse e-mail
          </label>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
            </svg>
            <input
              v-model="loginForm.email"
              type="email"
              required
              placeholder="vous@exemple.com"
              class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
          </div>
        </div>

        <div>
          <div class="flex items-center justify-between mb-1.5">
            <label class="text-sm font-medium text-gray-700">Mot de passe</label>
            <button type="button" class="text-xs text-primary hover:underline">
              Mot de passe oublié ?
            </button>
          </div>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
            <input
              v-model="loginForm.password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="Votre mot de passe"
              class="w-full pl-10 pr-10 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
              <svg v-if="!showPassword" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
              </svg>
              <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
              </svg>
            </button>
          </div>
        </div>

        <p v-if="loginError" class="text-xs text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ loginError }}
        </p>

        <button
          type="submit"
          :disabled="loginLoading"
          class="w-full py-2.5 rounded-lg bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors duration-150 disabled:opacity-60 mt-2">
          {{ loginLoading ? 'Connexion...' : 'Se connecter' }}
        </button>

      </form>

      <!-- Formulaire Inscription -->
      <form v-if="tab === 'register'" @submit.prevent="submitRegister" class="space-y-4">

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Prénom</label>
            <input
              v-model="registerForm.first_name"
              type="text"
              required
              placeholder="Jean"
              class="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Nom</label>
            <input
              v-model="registerForm.last_name"
              type="text"
              required
              placeholder="Dupont"
              class="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">Adresse e-mail</label>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
            </svg>
            <input
              v-model="registerForm.email"
              type="email"
              required
              placeholder="vous@exemple.com"
              class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">Mot de passe</label>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
            <input
              v-model="registerForm.password"
              type="password"
              required
              placeholder="Minimum 8 caractères"
              class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary"
            />
          </div>
        </div>

        <p v-if="registerError" class="text-xs text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ registerError }}
        </p>
        <p v-if="registerSuccess" class="text-xs text-green-700 bg-green-50 px-3 py-2 rounded-lg">
          {{ registerSuccess }}
        </p>

        <button
          type="submit"
          :disabled="registerLoading"
          class="w-full py-2.5 rounded-lg bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors duration-150 disabled:opacity-60">
          {{ registerLoading ? 'Création...' : 'Créer mon compte' }}
        </button>

      </form>

      <!-- Mentions légales -->
      <p class="text-xs text-center text-gray-400 mt-6 leading-relaxed">
        En continuant, vous acceptez nos
        <a href="#" class="underline hover:text-gray-600">conditions d'utilisation</a>
        et notre
        <a href="#" class="underline hover:text-gray-600">politique de confidentialité</a>.
      </p>

    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import api from '@/api.js'

const auth   = useAuthStore()
const router = useRouter()

const tab = ref('login')

// --- Connexion ---
const loginForm    = ref({ email: '', password: '' })
const loginError   = ref('')
const loginLoading = ref(false)
const showPassword = ref(false)

async function submitLogin() {
  loginError.value   = ''
  loginLoading.value = true
  try {
    await auth.login(loginForm.value.email, loginForm.value.password)
    router.push('/admin/dashboard')
  } catch (e) {
    loginError.value = e.response?.data?.error ?? 'Email ou mot de passe incorrect'
  } finally {
    loginLoading.value = false
  }
}

// --- Inscription ---
const registerForm    = ref({ first_name: '', last_name: '', email: '', password: '' })
const registerError   = ref('')
const registerSuccess = ref('')
const registerLoading = ref(false)

async function submitRegister() {
  registerError.value   = ''
  registerSuccess.value = ''
  registerLoading.value = true
  try {
    await api.post('/auth/register', {
      email:         registerForm.value.email,
      password_user: registerForm.value.password,
      first_name:    registerForm.value.first_name,
      last_name:     registerForm.value.last_name,
    })
    registerSuccess.value = 'Compte créé ! Vous pouvez vous connecter.'
    registerForm.value = { first_name: '', last_name: '', email: '', password: '' }
    setTimeout(() => { tab.value = 'login' }, 1500)
  } catch (e) {
    registerError.value = e.response?.data?.error ?? 'Erreur lors de la création du compte'
  } finally {
    registerLoading.value = false
  }
}
</script>
