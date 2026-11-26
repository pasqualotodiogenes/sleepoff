<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useSleepoffSiteState, type SleepoffSectionId } from '~/composables/useSleepoffSiteState'

interface NavRoute {
  id: SleepoffSectionId
  label: string
  href: string
}

const ROUTES: NavRoute[] = [
  { id: '01', label: '01 / EMULADOR', href: '#emulador' },
  { id: '02', label: '02 / SETUP', href: '#setup' },
  { id: '03', label: '03 / ARQUITETURA', href: '#arquitetura' },
  { id: '04', label: '04 / CLI', href: '#cli' },
]

const {
  activeSectionId,
  isEmulatorFocused,
  toggleCommandPalette,
  toggleEmulatorFocus,
} = useSleepoffSiteState()

const isMobileMenuOpen = ref(false)

function navigateToSection(id: SleepoffSectionId, href: string) {
  activeSectionId.value = id
  isMobileMenuOpen.value = false

  if (import.meta.client) {
    const target = document.querySelector(href)
    if (target) {
      target.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }
}

function handleConsoleToggle() {
  isMobileMenuOpen.value = false
  toggleEmulatorFocus(true)
  if (import.meta.client) {
    const target = document.querySelector('#emulador')
    if (target) {
      target.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }
}

function handleOpenPalette() {
  isMobileMenuOpen.value = false
  toggleCommandPalette(true)
}

function handleEscapeKey(event: KeyboardEvent) {
  if (event.key === 'Escape' && isMobileMenuOpen.value) {
    isMobileMenuOpen.value = false
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscapeKey)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscapeKey)
})
</script>

<template>
  <header
    class="sticky top-0 z-40 w-full rounded-none border-b border-[#1e1e2a] bg-[#06070a]/95 backdrop-blur-sm"
  >
    <!-- Desktop 1px Cell Bar -->
    <div class="hidden w-full items-stretch divide-x divide-[#1e1e2a] lg:flex">
      <!-- Brand Cell -->
      <a
        href="#top"
        class="header-cell-hover flex min-h-[48px] shrink-0 items-center gap-2.5 px-4 transition-colors xl:px-5"
        aria-label="sleepoff — voltar ao topo"
      >
        <span class="font-mono text-sm font-semibold tracking-tight text-[#ededf3]">
          <span class="text-[#70707d]">[</span><span class="text-[#5266eb]">S_</span> sleepoff<span class="text-[#70707d]">]</span>
        </span>
        <span
          class="rounded-[2px] border border-[#272735] bg-[#0e1016] px-1.5 py-0.5 font-mono text-[10px] tracking-[0.06em] text-[#70707d]"
        >
          v1.0.3
        </span>
      </a>

      <!-- CLI Route Cells -->
      <nav
        aria-label="Navegação principal"
        class="flex min-w-0 flex-1 items-stretch divide-x divide-[#1e1e2a]"
      >
        <a
          v-for="route in ROUTES"
          :key="route.id"
          :href="route.href"
          :aria-current="activeSectionId === route.id ? 'location' : undefined"
          class="header-cell-hover relative flex min-h-[48px] shrink-0 items-center gap-2 px-3.5 font-mono-hud text-xs transition-colors xl:px-5"
          :class="
            activeSectionId === route.id
              ? 'bg-[#0e1016] text-[#ededf3]'
              : 'text-[#70707d]'
          "
          @click.prevent="navigateToSection(route.id, route.href)"
        >
          <span
            class="inline-block h-1.5 w-1.5 shrink-0 rounded-none transition-colors"
            :class="activeSectionId === route.id ? 'bg-[#5266eb]' : 'bg-[#272735]'"
            aria-hidden="true"
          />
          <span class="whitespace-nowrap">{{ route.label }}</span>
          <span
            v-if="activeSectionId === route.id"
            class="absolute inset-x-0 bottom-0 h-[2px] bg-[#5266eb]"
            aria-hidden="true"
          />
        </a>

        <!-- Structural Spacer Cell -->
        <div
          class="hidden min-w-0 flex-1 items-center justify-end px-4 xl:flex"
          aria-hidden="true"
        >
          <span class="truncate font-mono-hud text-[10px] text-[#70707d]">
            {{ isEmulatorFocused ? 'CONSOLE // ATIVO' : 'SYS // PRONTO' }}
          </span>
        </div>
      </nav>

      <!-- Command Palette Trigger Cell -->
      <button
        type="button"
        class="header-cell-hover flex min-h-[48px] shrink-0 cursor-pointer items-center gap-3 px-4 font-mono text-xs text-[#c3c3cc] transition-colors"
        aria-label="Abrir paleta de comandos (Atalho Command K ou Control K)"
        @click="handleOpenPalette"
      >
        <span class="whitespace-nowrap">&gt;_ Ir para...</span>
        <kbd class="px-1.5 py-0.5 font-mono text-[11px] leading-none text-[#ededf3]">
          ⌘K
        </kbd>
      </button>

      <!-- Primary Action Cell -->
      <a
        href="#setup"
        class="cobalt-cta flex min-h-[48px] shrink-0 items-center justify-center bg-[#5266eb] px-5 font-mono-hud text-xs font-semibold text-[#ededf3] transition-colors xl:px-6"
        @click.prevent="navigateToSection('02', '#setup')"
      >
        BAIXAR APP
      </a>
    </div>

    <!-- Mobile 1px Cell Bar (< lg) -->
    <div class="flex w-full items-stretch divide-x divide-[#1e1e2a] lg:hidden">
      <a
        href="#top"
        class="flex min-h-[44px] min-w-0 flex-1 items-center gap-2 px-3"
        aria-label="sleepoff — voltar ao topo"
      >
        <span class="truncate font-mono text-sm font-semibold tracking-tight text-[#ededf3]">
          <span class="text-[#70707d]">[</span><span class="text-[#5266eb]">S_</span> sleepoff<span class="text-[#70707d]">]</span>
        </span>
        <span
          class="hidden rounded-[2px] border border-[#272735] bg-[#0e1016] px-1.5 py-0.5 font-mono text-[10px] text-[#70707d] sm:inline-flex"
        >
          v1.0.3
        </span>
      </a>

      <button
        type="button"
        class="flex min-h-[44px] min-w-[44px] shrink-0 cursor-pointer items-center justify-center px-3 font-mono text-xs text-[#c3c3cc] active:bg-[#171721]"
        aria-label="Abrir paleta de comandos"
        @click="handleOpenPalette"
      >
        <kbd class="px-1.5 py-0.5 font-mono text-[11px] leading-none text-[#ededf3]">
          ⌘K
        </kbd>
      </button>

      <button
        type="button"
        class="flex min-h-[44px] min-w-[44px] shrink-0 cursor-pointer items-center justify-center gap-1.5 px-3 font-mono-hud text-[11px] text-[#ededf3] active:bg-[#171721]"
        :aria-expanded="isMobileMenuOpen"
        aria-controls="mobile-nav-drawer"
        aria-label="Alternar menu de navegação"
        @click="isMobileMenuOpen = !isMobileMenuOpen"
      >
        <span>{{ isMobileMenuOpen ? 'FECHAR' : 'MENU' }}</span>
        <svg
          class="h-4 w-4 text-[#c3c3cc]"
          viewBox="0 0 16 16"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          aria-hidden="true"
        >
          <path v-if="!isMobileMenuOpen" d="M2 4h12M2 8h12M2 12h12" />
          <path v-else d="M3 3l10 10M13 3L3 13" />
        </svg>
      </button>
    </div>

    <!-- Mobile 1px Bordered Drawer -->
    <nav
      v-if="isMobileMenuOpen"
      id="mobile-nav-drawer"
      aria-label="Menu móvel"
      class="w-full divide-y divide-[#1e1e2a] border-t border-[#1e1e2a] bg-[#06070a] lg:hidden"
    >
      <a
        v-for="route in ROUTES"
        :key="`mobile-${route.id}`"
        :href="route.href"
        class="flex min-h-[48px] w-full items-center justify-between px-4 font-mono-hud text-xs transition-colors"
        :class="
          activeSectionId === route.id
            ? 'bg-[#0e1016] text-[#ededf3]'
            : 'text-[#c3c3cc]'
        "
        @click.prevent="navigateToSection(route.id, route.href)"
      >
        <span>{{ route.label }}</span>
        <span
          class="inline-block h-2 w-2 rounded-none"
          :class="activeSectionId === route.id ? 'bg-[#5266eb]' : 'bg-[#272735]'"
          aria-hidden="true"
        />
      </a>

      <button
        type="button"
        class="flex min-h-[48px] w-full cursor-pointer items-center justify-between px-4 text-left font-mono-hud text-xs text-[#c3c3cc] active:bg-[#171721]"
        @click="handleConsoleToggle"
      >
        <span>~ Console</span>
        <kbd class="px-1.5 py-0.5 font-mono text-[11px] leading-none text-[#ededf3]">
          ~
        </kbd>
      </button>

      <a
        href="#setup"
        class="flex min-h-[48px] w-full items-center justify-center bg-[#5266eb] px-4 font-mono-hud text-xs font-semibold text-[#ededf3]"
        @click.prevent="navigateToSection('02', '#setup')"
      >
        BAIXAR APP
      </a>
    </nav>
  </header>
</template>

<style scoped>
@media (hover: hover) {
  .header-cell-hover:hover {
    background-color: #0e1016;
    color: #ededf3;
  }

  .cobalt-cta:hover {
    background-color: #4458da;
  }
}
</style>
