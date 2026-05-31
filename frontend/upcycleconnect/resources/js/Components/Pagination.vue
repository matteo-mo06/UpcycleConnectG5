<template>
    <div v-if="totalPages > 1" class="flex items-center justify-center gap-1 mt-6">
        <button
            @click="$emit('update:page', page - 1)"
            :disabled="page <= 1"
            class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors duration-150 disabled:opacity-40 disabled:cursor-not-allowed bg-gray-100 text-gray-600 hover:bg-gray-200"
        >
            ←
        </button>

        <template v-for="p in pages" :key="p">
            <button
                v-if="p !== '...'"
                @click="$emit('update:page', p)"
                :class="[
                    'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors duration-150',
                    p === page ? 'bg-primary text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                ]"
            >{{ p }}</button>
            <span v-else class="px-2 text-gray-400 text-sm select-none">…</span>
        </template>

        <button
            @click="$emit('update:page', page + 1)"
            :disabled="page >= totalPages"
            class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors duration-150 disabled:opacity-40 disabled:cursor-not-allowed bg-gray-100 text-gray-600 hover:bg-gray-200"
        >
            →
        </button>
    </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
    page: { type: Number, required: true },
    total: { type: Number, required: true },
    limit: { type: Number, required: true },
})

defineEmits(['update:page'])

const totalPages = computed(() => Math.ceil(props.total / props.limit))

const pages = computed(() => {
    const total = totalPages.value
    const cur = props.page
    if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)

    const result = [1]
    if (cur > 3) result.push('...')
    for (let i = Math.max(2, cur - 1); i <= Math.min(total - 1, cur + 1); i++) result.push(i)
    if (cur < total - 2) result.push('...')
    result.push(total)
    return result
})
</script>
