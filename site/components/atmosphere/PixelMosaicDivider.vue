<script setup lang="ts">
import { computed } from 'vue'
import { useSleepoffSiteState } from '~/composables/useSleepoffSiteState'

interface Props {
  variant?: 'dark-to-onyx' | 'onyx-to-void'
}

interface MosaicCell {
  id: string
  row: number
  col: number
  tone: 'void' | 'panel' | 'onyx' | 'graphite' | 'accent'
  pulseDelay: string
  desktopOnly: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'dark-to-onyx',
})

const ROWS = 3
const COLS = 36

const { isPanicActive } = useSleepoffSiteState()

function deterministicHash(row: number, col: number, seed: number): number {
  const value = Math.imul(row + 11, 73856093) ^ Math.imul(col + 23, 19349663) ^ Math.imul(seed + 7, 83492791)
  return Math.abs(value) % 100
}

const cells = computed<MosaicCell[]>(() => {
  const seed = props.variant === 'dark-to-onyx' ? 17 : 53
  const list: MosaicCell[] = []

  for (let row = 0; row < ROWS; row++) {
    const stepRow = props.variant === 'dark-to-onyx' ? row : ROWS - 1 - row
    for (let col = 0; col < COLS; col++) {
      const h = deterministicHash(stepRow, col, seed)
      let tone: MosaicCell['tone'] = 'void'

      if (stepRow === 0) {
        if (h >= 97) tone = 'accent'
        else if (h >= 88) tone = 'onyx'
        else if (h >= 66) tone = 'panel'
        else tone = 'void'
      } else if (stepRow === 1) {
        if (h >= 95) tone = 'accent'
        else if (h >= 86) tone = 'graphite'
        else if (h >= 58) tone = 'onyx'
        else if (h >= 28) tone = 'panel'
        else tone = 'void'
      } else {
        if (h >= 96) tone = 'accent'
        else if (h >= 78) tone = 'graphite'
        else if (h >= 34) tone = 'onyx'
        else if (h >= 12) tone = 'panel'
        else tone = 'void'
      }

      list.push({
        id: `${row}-${col}`,
        row,
        col,
        tone,
        pulseDelay: `${((col * 190 + row * 310) % 2400)}ms`,
        desktopOnly: col >= 18,
      })
    }
  }

  return list
})

function cellClass(cell: MosaicCell): string {
  const base = cell.desktopOnly ? 'hidden md:block' : 'block'

  if (cell.tone === 'accent') {
    const color = isPanicActive.value ? 'bg-[#ef4444]/85' : 'bg-[#5266eb]/85'
    return `${base} ${color} mosaic-accent-pulse`
  }
  if (cell.tone === 'graphite') return `${base} bg-[#1e1e2a]`
  if (cell.tone === 'onyx') return `${base} bg-[#171721]`
  if (cell.tone === 'panel') return `${base} bg-[#0e1016]`
  return `${base} bg-[#06070a]`
}
</script>

<template>
  <div
    class="w-full overflow-hidden border-y border-[#1e1e2a] bg-[#06070a]"
    aria-hidden="true"
  >
    <div
      class="grid w-full grid-cols-[repeat(18,minmax(0,1fr))] divide-x divide-[#1e1e2a]/35 md:grid-cols-[repeat(36,minmax(0,1fr))]"
    >
      <span
        v-for="cell in cells"
        :key="cell.id"
        class="h-4 border-b border-[#1e1e2a]/35 transition-colors duration-300 sm:h-5"
        :class="cellClass(cell)"
        :style="cell.tone === 'accent' ? { animationDelay: cell.pulseDelay } : undefined"
      />
    </div>
  </div>
</template>

<style scoped>
@keyframes mosaic-pulse {
  0%, 100% {
    opacity: 0.55;
  }
  50% {
    opacity: 1;
  }
}

.mosaic-accent-pulse {
  animation: mosaic-pulse 3.2s ease-in-out infinite;
}

@media (prefers-reduced-motion: reduce) {
  .mosaic-accent-pulse {
    animation: none !important;
    opacity: 0.85;
  }
}
</style>
