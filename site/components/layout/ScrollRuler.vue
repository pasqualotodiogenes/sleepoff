<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'
import AnimatedCircularProgressBar from '~/components/ui/animated-circular-progressbar/AnimatedCircularProgressBar.vue'
import { useSleepoffSiteState, type SleepoffSectionId } from '~/composables/useSleepoffSiteState'

interface SectionTarget {
  id: SleepoffSectionId
  selector: string
}

const SECTION_TARGETS: SectionTarget[] = [
  { id: '01', selector: '#emulador' },
  { id: '02', selector: '#setup' },
  { id: '03', selector: '#arquitetura' },
  { id: '04', selector: '#cli' },
]

const TICKS = Array.from({ length: 21 }, (_, index) => ({
  index,
  isMajor: index % 5 === 0,
}))

const { scrollProgress, activeSectionId, isPanicActive } = useSleepoffSiteState()

const paddedProgress = computed(() => {
  const clamped = Math.max(0, Math.min(100, Math.round(scrollProgress.value)))
  return `${String(clamped).padStart(3, '0')}%`
})

const markerTopStyle = computed(() => {
  const clamped = Math.max(0, Math.min(100, scrollProgress.value))
  return { top: `${clamped}%` }
})

let observer: IntersectionObserver | null = null
let rafHandle: number | null = null

function updateScrollMetrics() {
  if (!import.meta.client) return
  const doc = document.documentElement
  const maxScroll = Math.max(0, doc.scrollHeight - window.innerHeight)
  const raw = maxScroll > 0 ? Math.round((window.scrollY / maxScroll) * 100) : 0
  scrollProgress.value = Math.max(0, Math.min(100, raw))

  if (window.scrollY <= 40) {
    activeSectionId.value = '01'
  }
}

function handleScroll() {
  if (rafHandle !== null) return
  rafHandle = window.requestAnimationFrame(() => {
    rafHandle = null
    updateScrollMetrics()
  })
}

function setupIntersectionObserver() {
  if (!import.meta.client || typeof IntersectionObserver === 'undefined') return

  const sectionElements: Array<{ id: SleepoffSectionId; el: Element }> = []
  for (const target of SECTION_TARGETS) {
    const el = document.querySelector(target.selector)
    if (el) {
      sectionElements.push({ id: target.id, el })
    }
  }

  if (sectionElements.length === 0) return

  observer = new IntersectionObserver(
    (entries) => {
      const visible = entries
        .filter((entry) => entry.isIntersecting)
        .sort((a, b) => b.intersectionRatio - a.intersectionRatio)

      const topEntry = visible[0]
      if (!topEntry) return

      const matched = sectionElements.find((item) => item.el === topEntry.target)
      if (matched) {
        activeSectionId.value = matched.id
      }
    },
    {
      rootMargin: '-20% 0px -55% 0px',
      threshold: [0, 0.1, 0.25, 0.5],
    },
  )

  for (const item of sectionElements) {
    observer.observe(item.el)
  }
}

function scrollToTop() {
  if (!import.meta.client) return
  const prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  window.scrollTo({
    top: 0,
    behavior: prefersReduced ? 'auto' : 'smooth',
  })
}

onMounted(() => {
  updateScrollMetrics()
  setupIntersectionObserver()
  window.addEventListener('scroll', handleScroll, { passive: true })
  window.addEventListener('resize', handleScroll, { passive: true })
})

onBeforeUnmount(() => {
  if (rafHandle !== null) {
    window.cancelAnimationFrame(rafHandle)
    rafHandle = null
  }
  if (observer) {
    observer.disconnect()
    observer = null
  }
  window.removeEventListener('scroll', handleScroll)
  window.removeEventListener('resize', handleScroll)
})
</script>

<template>
  <!-- Desktop Vertical Graduated Ruler -->
  <aside
    aria-label="Indicador vertical de leitura e seção ativa"
    class="fixed top-[49px] right-0 bottom-0 z-30 hidden w-10 select-none flex-col items-center justify-between border-l border-[#1e1e2a] bg-[#06070a]/80 py-4 backdrop-blur-[2px] lg:flex"
  >
    <!-- Active Section Tag -->
    <div class="flex flex-col items-center gap-1.5">
      <span
        class="h-1.5 w-1.5 rounded-none transition-colors"
        :class="isPanicActive ? 'bg-[#ef4444]' : 'bg-[#5266eb]'"
        aria-hidden="true"
      />
      <span
        class="writing-vertical font-mono-hud text-[10px] tracking-[0.14em] transition-colors"
        :class="isPanicActive ? 'text-[#ef4444]' : 'text-[#70707d]'"
      >
        SEC // {{ activeSectionId }}
      </span>
    </div>

    <!-- Graduated Ruler Track -->
    <div
      class="relative my-4 flex w-full flex-1 flex-col items-end justify-between pr-0"
      aria-hidden="true"
    >
      <div
        v-for="tick in TICKS"
        :key="tick.index"
        class="h-[1px] transition-colors"
        :class="
          tick.isMajor
            ? 'w-4 bg-[#70707d]/65'
            : 'w-2 bg-[#272735]'
        "
      />

      <!-- Live Position Marker -->
      <div
        class="pointer-events-none absolute right-0 left-0 -translate-y-1/2 transition-[top] duration-100"
        :style="markerTopStyle"
      >
        <div class="flex items-center justify-end gap-0.5">
          <span
            class="h-1.5 w-1.5 rounded-none"
            :class="isPanicActive ? 'bg-[#ef4444]' : 'bg-[#5266eb]'"
          />
          <span
            class="h-[2px] w-5"
            :class="isPanicActive ? 'bg-[#ef4444]' : 'bg-[#5266eb]'"
          />
        </div>
      </div>
    </div>

    <!-- Live Percentage & Quick Top -->
    <button
      type="button"
      class="ruler-top-btn flex min-h-[44px] w-full cursor-pointer flex-col items-center justify-center gap-0.5 border-t border-[#1e1e2a] pt-2 font-mono text-[10px] tabular-nums transition-colors"
      :class="isPanicActive ? 'text-[#ef4444]' : 'text-[#c3c3cc]'"
      :aria-label="`Voltar ao topo da página (progresso atual: ${paddedProgress})`"
      title="Voltar ao topo"
      @click="scrollToTop"
    >
      <span>{{ paddedProgress }}</span>
      <span class="text-[10px] text-[#70707d]" aria-hidden="true">↑</span>
    </button>
  </aside>

  <!-- Mobile Floating Circular Progress Ring Button -->
  <button
    type="button"
    class="fixed right-4 bottom-4 z-30 flex min-h-[44px] min-w-[44px] cursor-pointer items-center justify-center rounded-[2px] border bg-[#0e1016]/95 p-0.5 transition-colors lg:hidden"
    :class="isPanicActive ? 'border-[#ef4444]' : 'border-[#272735]'"
    :aria-label="`Voltar ao topo da página (progresso: ${paddedProgress})`"
    title="Voltar ao topo"
    @click="scrollToTop"
  >
    <div class="relative flex items-center justify-center">
      <AnimatedCircularProgressBar
        :value="scrollProgress"
        :gauge-primary-color="isPanicActive ? '#ef4444' : '#5266eb'"
        gauge-secondary-color="#1e1e2a"
        :circle-stroke-width="8"
        :show-percentage="false"
        class="size-11"
      />
      <span
        class="pointer-events-none absolute font-mono text-sm font-bold"
        :class="isPanicActive ? 'text-[#ef4444]' : 'text-[#ededf3]'"
        aria-hidden="true"
      >
        ↑
      </span>
    </div>
  </button>
</template>

<style scoped>
.writing-vertical {
  writing-mode: vertical-rl;
  transform: rotate(180deg);
}

@media (hover: hover) {
  .ruler-top-btn:hover {
    background-color: #0e1016;
    color: #ededf3;
  }
}
</style>
