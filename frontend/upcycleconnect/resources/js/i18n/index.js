import { createI18n } from 'vue-i18n'
import fr from './locales/fr.js'
import en from './locales/en.js'

const savedLocale = localStorage.getItem('locale') ?? 'fr'

const i18n = createI18n({
    legacy: false,
    locale: savedLocale,
    fallbackLocale: 'fr',
    messages: { fr, en },
})

export default i18n
