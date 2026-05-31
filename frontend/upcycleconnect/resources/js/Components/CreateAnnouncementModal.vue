<template>
    <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="close"
    >
        <div class="absolute inset-0 bg-black/40" @click="close" />
        <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-xl overflow-hidden max-h-[90vh] flex flex-col">

            <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between flex-shrink-0">
                <h2 class="font-semibold text-gray-800 text-lg" style="font-family: var(--font-family-title)">Déposer un objet</h2>
                <button @click="close" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                </button>
            </div>

            <div class="overflow-y-auto flex-1 px-6 py-5 space-y-5">

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Type de dépôt</label>
                    <div class="flex gap-2">
                        <button
                            v-for="t in types"
                            :key="t.value"
                            @click="form.type = t.value; if (t.value === 'don') form.price = 0"
                            :class="[
                                'px-4 py-2 rounded-full text-sm font-medium border transition-colors',
                                form.type === t.value
                                    ? 'bg-primary text-white border-primary'
                                    : 'bg-white text-gray-600 border-gray-200 hover:border-primary hover:text-primary'
                            ]"
                        >{{ t.label }}</button>
                    </div>
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Titre</label>
                    <input
                        v-model="form.title"
                        type="text"
                        placeholder="Nom de votre objet"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">État de l'objet</label>
                    <div class="flex gap-2">
                        <button
                            v-for="c in conditions"
                            :key="c.value"
                            @click="form.condition = c.value"
                            :class="[
                                'px-4 py-2 rounded-full text-sm font-medium border transition-colors',
                                form.condition === c.value
                                    ? 'bg-secondary text-white border-secondary'
                                    : 'bg-white text-gray-600 border-gray-200 hover:border-secondary hover:text-secondary-dark'
                            ]"
                        >{{ c.label }}</button>
                    </div>
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                    <textarea
                        v-model="form.description"
                        rows="3"
                        placeholder="Décrivez votre objet : matière, dimensions, contexte…"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary resize-none"
                    />
                </div>

                <div class="grid grid-cols-2 gap-3">
                    <div class="col-span-2">
                        <label class="block text-sm font-medium text-gray-700 mb-1">Adresse</label>
                        <input
                            v-model="form.address"
                            type="text"
                            placeholder="Adresse de dépôt"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Ville</label>
                        <input
                            v-model="form.city"
                            type="text"
                            placeholder="Paris"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Code postal</label>
                        <input
                            v-model="form.postal"
                            type="text"
                            placeholder="75001"
                            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                        />
                    </div>
                </div>

                <div v-if="form.type === 'vente'">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Prix (€)</label>
                    <input
                        v-model.number="form.price"
                        type="number"
                        min="0.01"
                        step="0.01"
                        placeholder="Ex : 10.00"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Disponible à partir du</label>
                    <input
                        v-model="form.availability_date"
                        type="date"
                        class="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Photos de l'objet <span class="text-gray-400 font-normal">(3 max)</span></label>
                    <div
                        class="border-2 border-dashed border-gray-200 rounded-xl p-6 text-center cursor-pointer hover:border-primary/50 transition-colors"
                        @dragover.prevent
                        @drop.prevent="handleDrop"
                        @click="$refs.fileInput.click()"
                    >
                        <svg class="w-8 h-8 text-gray-300 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"/>
                        </svg>
                        <p class="text-sm text-gray-400">Glissez vos photos ici</p>
                        <button type="button" class="mt-2 text-xs text-primary underline">Parcourir</button>
                        <input ref="fileInput" type="file" accept="image/*" multiple class="hidden" @change="handleFileSelect" />
                    </div>

                    <div v-if="photos.length" class="flex gap-2 mt-3 flex-wrap">
                        <div v-for="(photo, i) in photos" :key="i" class="relative w-16 h-16 rounded-lg overflow-hidden border border-gray-200">
                            <img :src="photo.preview" class="w-full h-full object-cover" />
                            <button
                                @click="removePhoto(i)"
                                class="absolute top-0.5 right-0.5 w-4 h-4 bg-black/60 rounded-full flex items-center justify-center text-white"
                            >
                                <svg class="w-2.5 h-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>

                <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

            </div>

            <div class="px-6 py-4 border-t border-gray-100 flex-shrink-0">
                <button
                    @click="submit"
                    :disabled="loading"
                    class="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    {{ loading ? 'Publication…' : 'Publier' }}
                </button>
            </div>

        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import api from '@/api.js'

const props = defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue', 'created'])

const loading = ref(false)
const error = ref('')
const photos = ref([])

const form = reactive({
    type: 'don',
    title: '',
    condition: 'bon_etat',
    description: '',
    address: '',
    city: '',
    postal: '',
    price: 0,
    availability_date: '',
})

const types = [
    { label: 'Don', value: 'don' },
    { label: 'Vente', value: 'vente' },
]

const conditions = [
    { label: 'Neuf', value: 'neuf' },
    { label: 'Bon état', value: 'bon_etat' },
    { label: 'Usé', value: 'use' },
]

function close() {
    if (!loading.value) emit('update:modelValue', false)
}

function handleFileSelect(e) {
    addFiles(Array.from(e.target.files))
    e.target.value = ''
}

function handleDrop(e) {
    addFiles(Array.from(e.dataTransfer.files).filter(f => f.type.startsWith('image/')))
}

function addFiles(files) {
    for (const file of files) {
        if (photos.value.length >= 3) break
        photos.value.push({ file, preview: URL.createObjectURL(file) })
    }
}

function removePhoto(i) {
    URL.revokeObjectURL(photos.value[i].preview)
    photos.value.splice(i, 1)
}

async function uploadPhotos() {
    const urls = []
    for (const photo of photos.value) {
        const fd = new FormData()
        fd.append('file', photo.file)
        const { data } = await api.post('/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
        urls.push(data.url)
    }
    return urls
}

async function submit() {
    error.value = ''
    if (!form.title.trim()) { error.value = 'Le titre est requis.'; return }
    if (!form.availability_date) { error.value = 'La date de disponibilité est requise.'; return }
    if (form.type === 'vente' && (!form.price || form.price < 0.01)) { error.value = 'Le prix minimum pour une vente est 0.01 €.'; return }

    loading.value = true
    try {
        const photoURLs = photos.value.length ? await uploadPhotos() : []
        await api.post('/announcements', {
            ...form,
            photo_urls: photoURLs,
        })
        emit('created')
        emit('update:modelValue', false)
        resetForm()
    } catch (e) {
        error.value = e.response?.data?.error ?? 'Une erreur est survenue.'
    } finally {
        loading.value = false
    }
}

function resetForm() {
    Object.assign(form, { type: 'don', title: '', condition: 'bon_etat', description: '', address: '', city: '', postal: '', price: 0, availability_date: '' })
    photos.value.forEach(p => URL.revokeObjectURL(p.preview))
    photos.value = []
    error.value = ''
}
</script>
