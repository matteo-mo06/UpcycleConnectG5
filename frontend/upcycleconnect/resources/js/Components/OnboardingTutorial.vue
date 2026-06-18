<template>
    <Teleport to="body">
        <div class="fixed inset-0 z-[9999]" style="pointer-events: none">

            <svg class="absolute inset-0 w-full h-full" style="pointer-events: none">
                <path :d="overlayPath" fill="rgba(0,0,0,0.6)" fill-rule="evenodd" pointer-events="fill" />
                <rect
                    :x="rect.left" :y="rect.top"
                    :width="rect.width" :height="rect.height"
                    :rx="rect.radius" :ry="rect.radius"
                    fill="none" stroke="rgba(255,255,255,0.55)" stroke-width="2"
                    pointer-events="none"
                />
            </svg>

            <div class="absolute" :style="holeStyle" style="pointer-events: all" />

            <div
                class="absolute bg-white rounded-2xl shadow-xl p-5 w-72 transition-[top,left] duration-300"
                :style="tooltipStyle"
                style="pointer-events: all"
            >
                <div class="flex items-center justify-between mb-3">
                    <div class="flex gap-1.5">
                        <span
                            v-for="(_, i) in steps" :key="i"
                            :class="['w-2 h-2 rounded-full transition-colors duration-200', i === currentStep ? 'bg-primary' : 'bg-gray-200']"
                        />
                    </div>
                    <span class="text-xs text-gray-400">{{ currentStep + 1 }} / {{ steps.length }}</span>
                </div>

                <h3 class="font-semibold text-gray-800 mb-2" style="font-family: var(--font-family-title)">
                    {{ steps[currentStep].title }}
                </h3>
                <p class="text-sm text-gray-500 leading-relaxed">{{ steps[currentStep].description }}</p>

                <div class="flex gap-2 mt-4">
                    <button
                        v-if="currentStep > 0"
                        @click="prev"
                        class="flex-1 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >{{ t('accueil.tutorialNav.prev') }}</button>
                    <div v-else class="flex-1" />
                    <button
                        @click="next"
                        class="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors"
                    >{{ currentStep === steps.length - 1 ? t('accueil.tutorialNav.finish') : t('accueil.tutorialNav.next') }}</button>
                </div>
            </div>

        </div>
    </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({ steps: { type: Array, required: true } })
const { t } = useI18n({ useScope: 'global' })
const emit = defineEmits(['done'])

const currentStep = ref(0)
const rect = ref({ top: 0, left: 0, width: 0, height: 0, radius: 0 })

const PAD = 8
const TOOLTIP_W = 288
const GAP = 16

function updateRect() {
    const el = document.querySelector(`[data-tutorial="${props.steps[currentStep.value].target}"]`)
    if (!el) return
    const r = el.getBoundingClientRect()
    const rawRadius = parseFloat(getComputedStyle(el).borderRadius) || 0
    rect.value = {
        top:    r.top    - PAD,
        left:   r.left   - PAD,
        width:  r.width  + PAD * 2,
        height: r.height + PAD * 2,
        radius: rawRadius + PAD,
    }
}

function roundedRectPath(x, y, w, h, r) {
    r = Math.min(r, w / 2, h / 2)
    return `M${x + r},${y} H${x + w - r} A${r},${r} 0 0 1 ${x + w},${y + r} V${y + h - r} A${r},${r} 0 0 1 ${x + w - r},${y + h} H${x + r} A${r},${r} 0 0 1 ${x},${y + h - r} V${y + r} A${r},${r} 0 0 1 ${x + r},${y} Z`
}

const overlayPath = computed(() => {
    const vw = window.innerWidth
    const vh = window.innerHeight
    const { top, left, width, height, radius } = rect.value
    const outer = `M0,0 H${vw} V${vh} H0 Z`
    const inner = roundedRectPath(left, top, width, height, radius)
    return `${outer} ${inner}`
})

const holeStyle = computed(() => ({
    top:    rect.value.top    + 'px',
    left:   rect.value.left   + 'px',
    width:  rect.value.width  + 'px',
    height: rect.value.height + 'px',
}))

const tooltipStyle = computed(() => {
    const { top, left, width, height } = rect.value
    const placement = props.steps[currentStep.value].placement ?? 'right'
    const vw = window.innerWidth
    const vh = window.innerHeight
    let tx, ty

    if (placement === 'right') {
        tx = left + width + GAP
        ty = top + height / 2 - 100
    } else if (placement === 'left') {
        tx = left - TOOLTIP_W - GAP
        ty = top + height / 2 - 100
    } else if (placement === 'bottom') {
        tx = left + width / 2 - TOOLTIP_W / 2
        ty = top + height + GAP
    } else {
        tx = left + width / 2 - TOOLTIP_W / 2
        ty = top - GAP - 220
    }

    tx = Math.max(8, Math.min(tx, vw - TOOLTIP_W - 8))
    ty = Math.max(8, Math.min(ty, vh - 240))
    return { left: tx + 'px', top: ty + 'px', width: TOOLTIP_W + 'px' }
})

function blockInput(e) {
    e.stopPropagation()
    e.preventDefault()
}

watch(currentStep, () => nextTick(updateRect))

onMounted(() => {
    nextTick(updateRect)
    window.addEventListener('resize', updateRect)
    document.addEventListener('keydown', blockInput, true)
})

onUnmounted(() => {
    window.removeEventListener('resize', updateRect)
    document.removeEventListener('keydown', blockInput, true)
})

function next() {
    if (currentStep.value < props.steps.length - 1) currentStep.value++
    else emit('done')
}
function prev() {
    if (currentStep.value > 0) currentStep.value--
}
</script>
