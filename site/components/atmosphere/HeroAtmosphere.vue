<script setup lang="ts">
import { computed } from 'vue'
import DitherShader from '~/components/ui/dither-shader/DitherShader.vue'
import ParticlesBg from '~/components/ui/particles-bg/ParticlesBg.vue'
import { useSleepoffSiteState } from '~/composables/useSleepoffSiteState'

interface RingTick {
  id: number
  x1: string
  y1: string
  x2: string
  y2: string
  major: boolean
}

const SVG_SIGNAL = `<svg xmlns="http://www.w3.org/2000/svg" width="320" height="180" viewBox="0 0 320 180">
  <defs>
    <radialGradient id="rg" cx="50%" cy="40%" r="68%">
      <stop offset="0%" stop-color="#f2f2f2"/>
      <stop offset="48%" stop-color="#7c7c7c"/>
      <stop offset="100%" stop-color="#050505"/>
    </radialGradient>
    <pattern id="grid" width="32" height="32" patternUnits="userSpaceOnUse">
      <path d="M 32 0 L 0 0 0 32" fill="none" stroke="#5c5c5c" stroke-width="1"/>
    </pattern>
  </defs>
  <rect width="320" height="180" fill="url(#rg)"/>
  <rect width="320" height="180" fill="url(#grid)" opacity="0.45"/>
</svg>`

const DITHER_SOURCE_URI = `data:image/svg+xml;utf8,${encodeURIComponent(SVG_SIGNAL)}`

function buildRingTicks(count: number, innerR: number, outerR: number, majorEvery: number): RingTick[] {
  const cx = 400
  const cy = 400
  const ticks: RingTick[] = []
  for (let i = 0; i < count; i++) {
    const angle = (i / count) * Math.PI * 2 - Math.PI / 2
    const major = i % majorEvery === 0
    const rStart = major ? innerR - 8 : innerR
    const rEnd = major ? outerR + 6 : outerR
    ticks.push({
      id: i,
      x1: (cx + Math.cos(angle) * rStart).toFixed(2),
      y1: (cy + Math.sin(angle) * rStart).toFixed(2),
      x2: (cx + Math.cos(angle) * rEnd).toFixed(2),
      y2: (cy + Math.sin(angle) * rEnd).toFixed(2),
      major,
    })
  }
  return ticks
}

const INNER_RING_TICKS = buildRingTicks(72, 228, 240, 6)
const OUTER_RING_TICKS = buildRingTicks(36, 314, 324, 3)

const { isPanicActive } = useSleepoffSiteState()

const particleColor = computed(() => (isPanicActive.value ? '#ef4444' : '#c3c3cc'))
const lineColor = computed(() => (isPanicActive.value ? '#ef4444' : '#5266eb'))
const lineOpacity = computed(() => (isPanicActive.value ? 0.28 : 0.16))
const ditherSecondaryColor = computed(() => (isPanicActive.value ? '#ef4444' : '#5266eb'))
const ringStrokePrimary = computed(() => (isPanicActive.value ? '#ef4444' : '#5266eb'))
const ringStrokeSecondary = computed(() => (isPanicActive.value ? '#ef4444' : '#70707d'))
</script>

<template>
  <div
    class="pointer-events-none absolute inset-0 select-none overflow-hidden"
    aria-hidden="true"
  >
    <!-- Layer 1: Panic-Reactive DitherShader -->
    <ClientOnly>
      <div
        class="pointer-events-none absolute inset-0 transition-opacity duration-300"
        :class="isPanicActive ? 'opacity-[0.34]' : 'opacity-[0.14]'"
      >
        <DitherShader
          :src="DITHER_SOURCE_URI"
          dither-mode="bayer"
          color-mode="duotone"
          primary-color="#06070a"
          :secondary-color="ditherSecondaryColor"
          :animated="isPanicActive"
          :animation-speed="0.04"
          :grid-size="4"
          :threshold="0.52"
          :contrast="1.15"
          object-fit="cover"
          class="h-full w-full"
        />
      </div>
      <template #fallback>
        <div class="pointer-events-none absolute inset-0 bg-[#06070a]" />
      </template>
    </ClientOnly>

    <!-- Layer 2: Graduated Concentric Hero Rings (getartcraft motif) -->
    <div
      class="pointer-events-none absolute inset-0 flex items-center justify-center transition-opacity duration-300"
      :class="[
        isPanicActive ? 'hero-rings-panic opacity-[0.22]' : 'opacity-[0.09]',
      ]"
    >
      <svg
        viewBox="0 0 800 800"
        fill="none"
        class="h-[680px] w-[680px] max-w-none sm:h-[820px] sm:w-[820px]"
      >
        <!-- Crosshair Axes -->
        <line x1="40" y1="400" x2="760" y2="400" :stroke="ringStrokeSecondary" stroke-width="1" stroke-dasharray="4 8" />
        <line x1="400" y1="40" x2="400" y2="760" :stroke="ringStrokeSecondary" stroke-width="1" stroke-dasharray="4 8" />

        <!-- Concentric Circles -->
        <circle cx="400" cy="400" r="150" :stroke="ringStrokeSecondary" stroke-width="1" />
        <circle cx="400" cy="400" r="240" :stroke="ringStrokePrimary" stroke-width="1" />
        <circle cx="400" cy="400" r="320" :stroke="ringStrokeSecondary" stroke-width="1" stroke-dasharray="2 6" />

        <!-- Graduated Tick Marks -->
        <line
          v-for="tick in INNER_RING_TICKS"
          :key="`inner-${tick.id}`"
          :x1="tick.x1"
          :y1="tick.y1"
          :x2="tick.x2"
          :y2="tick.y2"
          :stroke="tick.major ? ringStrokePrimary : ringStrokeSecondary"
          :stroke-width="tick.major ? 1.5 : 1"
        />
        <line
          v-for="tick in OUTER_RING_TICKS"
          :key="`outer-${tick.id}`"
          :x1="tick.x1"
          :y1="tick.y1"
          :x2="tick.x2"
          :y2="tick.y2"
          :stroke="ringStrokeSecondary"
          :stroke-width="tick.major ? 1.25 : 0.75"
        />
      </svg>
    </div>

    <!-- Layer 3: Constellation Particles -->
    <ClientOnly>
      <ParticlesBg
        :connect-lines="true"
        :quantity="64"
        :staticity="45"
        :ease="55"
        :color="particleColor"
        :line-color="lineColor"
        :line-opacity="lineOpacity"
        :line-distance="125"
        class="pointer-events-none absolute inset-0 h-full w-full"
      />
      <template #fallback>
        <div class="pointer-events-none absolute inset-0" />
      </template>
    </ClientOnly>

    <!-- Bottom Fade Seam -->
    <div
      class="pointer-events-none absolute inset-x-0 bottom-0 h-28 bg-gradient-to-t from-[#06070a] to-transparent"
    />
  </div>
</template>

<style scoped>
@keyframes hero-panic-pulse {
  0%, 100% {
    opacity: 0.2;
    transform: scale(1);
  }
  50% {
    opacity: 0.32;
    transform: scale(1.012);
  }
}

.hero-rings-panic {
  animation: hero-panic-pulse 1.2s ease-in-out infinite;
}

@media (prefers-reduced-motion: reduce) {
  .hero-rings-panic {
    animation: none !important;
    transform: none !important;
  }
}
</style>
