<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useSleepoffSiteState, type SleepoffSectionId } from '~/composables/useSleepoffSiteState'

interface PaletteItem {
  id: string
  category: string
  label: string
  hint: string
  shortcut?: string
  keywords: string
  run: () => void | Promise<void>
}

const GO_INSTALL_CMD = 'go install github.com/pasqualotodiogenes/sleepoff@latest'
const WINDOWS_SETUP_URL =
  'https://github.com/pasqualotodiogenes/sleepoff/releases/latest/download/sleepoff-setup.exe'

const {
  isCommandPaletteOpen,
  activeSectionId,
  toggleCommandPalette,
  toggleEmulatorFocus,
  triggerCommandInEmulator,
} = useSleepoffSiteState()

const query = ref('')
const activeIndex = ref(0)
const copiedGoInstall = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
let copyTimeout: ReturnType<typeof setTimeout> | null = null

function scrollToSelector(selector: string, sectionId?: SleepoffSectionId) {
  if (sectionId) activeSectionId.value = sectionId
  if (!import.meta.client) return
  const target = document.querySelector(selector)
  if (!target) return
  const prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  target.scrollIntoView({ behavior: prefersReduced ? 'auto' : 'smooth', block: 'start' })
}

async function copyGoInstallCommand() {
  if (import.meta.client && navigator.clipboard) {
    try {
      await navigator.clipboard.writeText(GO_INSTALL_CMD)
    } catch {
      // Ignore clipboard denial in restricted contexts
    }
  }
  copiedGoInstall.value = true
  if (copyTimeout) clearTimeout(copyTimeout)
  copyTimeout = setTimeout(() => {
    copiedGoInstall.value = false
  }, 2000)
}

function jump(selector: string, sectionId: SleepoffSectionId) {
  toggleCommandPalette(false)
  scrollToSelector(selector, sectionId)
}

const allItems = computed<PaletteItem[]>(() => [
  {
    id: 'nav-01', category: 'NAVEGAÇÃO', label: '01 / Emulador Tático',
    hint: '#emulador', shortcut: '01', keywords: 'emulador terminal console tatico 01',
    run: () => jump('#emulador', '01'),
  },
  {
    id: 'nav-02', category: 'NAVEGAÇÃO', label: '02 / Setup (App & Terminal)',
    hint: '#setup', shortcut: '02', keywords: 'setup instalar app windows linux macos download 02',
    run: () => jump('#setup', '02'),
  },
  {
    id: 'nav-03', category: 'NAVEGAÇÃO', label: '03 / Arquitetura Zero Bloat',
    hint: '#arquitetura', shortcut: '03', keywords: 'arquitetura zero bloat lock dry-run mutex flock 03',
    run: () => jump('#arquitetura', '03'),
  },
  {
    id: 'nav-04', category: 'NAVEGAÇÃO', label: '04 / Referência CLI',
    hint: '#cli', shortcut: '04', keywords: 'cli comandos flags atalhos ajuda 04',
    run: () => jump('#cli', '04'),
  },
  {
    id: 'emu-toggle', category: 'CONSOLE', label: '~ Focar / Alternar Console Tático',
    hint: 'Alterna foco no emulador integrado', shortcut: '~', keywords: 'console focar alternar tilde ~ emulador',
    run: () => {
      toggleCommandPalette(false)
      toggleEmulatorFocus()
      scrollToSelector('#emulador', '01')
    },
  },
  {
    id: 'emu-15m', category: 'EMULADOR', label: 'Executar sleepoff 15m --dry-run no Emulador',
    hint: 'Missão 1 · Simulação segura de 15 minutos', shortcut: 'RUN', keywords: 'executar sleepoff 15m dry-run simular timer missao 1',
    run: () => {
      toggleCommandPalette(false)
      triggerCommandInEmulator('sleepoff 15m --dry-run')
      scrollToSelector('#emulador', '01')
    },
  },
  {
    id: 'emu-panic', category: 'EMULADOR', label: 'Simular Tela de Pânico (10s) no Emulador',
    hint: 'Dispara contagem regressiva de alerta de 10s', shortcut: '10s', keywords: 'panico pânico alerta 10s abortar som alarme',
    run: () => {
      toggleCommandPalette(false)
      triggerCommandInEmulator('panic-demo')
      scrollToSelector('#emulador', '01')
    },
  },
  {
    id: 'copy-go', category: 'INSTALAÇÃO', label: `Copiar ${GO_INSTALL_CMD}`,
    hint: copiedGoInstall.value ? 'COPIADO ✓' : 'Área de transferência',
    shortcut: copiedGoInstall.value ? 'COPIADO ✓' : 'COPY',
    keywords: 'copiar go install toolchain github pasqualotodiogenes',
    run: () => { void copyGoInstallCommand() },
  },
  {
    id: 'download-win', category: 'DOWNLOAD', label: 'Baixar App Windows (sleepoff-setup.exe)',
    hint: 'Instalador oficial com bandeja do sistema e atalhos', shortcut: '.EXE',
    keywords: 'baixar download windows exe instalador app bandeja',
    run: () => {
      toggleCommandPalette(false)
      if (import.meta.client) window.open(WINDOWS_SETUP_URL, '_blank', 'noopener,noreferrer')
    },
  },
])

const filteredItems = computed(() => {
  const term = query.value.trim().toLowerCase()
  if (!term) return allItems.value
  return allItems.value.filter(
    (item) =>
      item.label.toLowerCase().includes(term) ||
      item.category.toLowerCase().includes(term) ||
      item.hint.toLowerCase().includes(term) ||
      item.keywords.toLowerCase().includes(term),
  )
})

watch(filteredItems, (items) => {
  if (activeIndex.value >= items.length) {
    activeIndex.value = Math.max(0, items.length - 1)
  }
})

watch(isCommandPaletteOpen, async (open) => {
  if (!open) return
  query.value = ''
  activeIndex.value = 0
  await nextTick()
  inputRef.value?.focus()
})

function handleGlobalKeydown(event: KeyboardEvent) {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    toggleCommandPalette()
    return
  }
  if (!isCommandPaletteOpen.value) return
  if (event.key === 'Escape') {
    event.preventDefault()
    toggleCommandPalette(false)
    return
  }
  const total = filteredItems.value.length
  if (total === 0) return
  if (event.key === 'ArrowDown') {
    event.preventDefault()
    activeIndex.value = (activeIndex.value + 1) % total
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    activeIndex.value = (activeIndex.value - 1 + total) % total
  } else if (event.key === 'Enter') {
    const selected = filteredItems.value[activeIndex.value]
    if (selected) {
      event.preventDefault()
      selected.run()
    }
  }
}

onMounted(() => window.addEventListener('keydown', handleGlobalKeydown))
onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
  if (copyTimeout) {
    clearTimeout(copyTimeout)
    copyTimeout = null
  }
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isCommandPaletteOpen"
      role="dialog"
      aria-modal="true"
      aria-label="Paleta de comandos do sleepoff"
      class="fixed inset-0 z-50 flex items-start justify-center bg-[#06070a]/80 px-3 pt-16 backdrop-blur-sm sm:pt-24"
      @click.self="toggleCommandPalette(false)"
    >
      <div class="w-full max-w-xl rounded-none border border-[#272735] bg-[#0e1016] text-[#ededf3]">
        <div class="flex items-center border-b border-[#1e1e2a] bg-[#06070a]">
          <span class="flex min-h-[48px] items-center border-r border-[#1e1e2a] px-3.5 font-mono text-xs text-[#5266eb]" aria-hidden="true">
            &gt;_
          </span>
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Ir para seção, executar comando ou copiar..."
            aria-label="Filtrar comandos"
            class="min-h-[48px] w-full min-w-0 bg-transparent px-3.5 font-mono text-xs text-[#ededf3] placeholder-[#70707d] focus:outline-none"
          />
          <button
            type="button"
            class="flex min-h-[48px] min-w-[48px] shrink-0 cursor-pointer items-center justify-center border-l border-[#1e1e2a] px-3 font-mono-hud text-[11px] text-[#70707d] active:bg-[#171721]"
            aria-label="Fechar paleta de comandos"
            @click="toggleCommandPalette(false)"
          >
            ESC
          </button>
        </div>

        <div role="listbox" aria-label="Ações disponíveis" class="max-h-[60vh] divide-y divide-[#1e1e2a] overflow-y-auto">
          <div v-if="filteredItems.length === 0" class="px-4 py-8 text-center font-mono text-xs text-[#70707d]">
            Nenhum comando corresponde a "{{ query }}".
          </div>

          <button
            v-for="(item, idx) in filteredItems"
            :key="item.id"
            type="button"
            role="option"
            :aria-selected="idx === activeIndex"
            class="palette-row flex min-h-[48px] w-full cursor-pointer items-center justify-between gap-3 px-4 py-2.5 text-left transition-colors"
            :class="idx === activeIndex ? 'bg-[#171721] text-[#ededf3]' : 'bg-[#0e1016] text-[#c3c3cc]'"
            @mouseenter="activeIndex = idx"
            @click="item.run()"
          >
            <div class="flex min-w-0 flex-1 items-center gap-3">
              <span
                class="inline-block h-2 w-2 shrink-0 rounded-none"
                :class="idx === activeIndex ? 'bg-[#5266eb]' : 'bg-[#272735]'"
                aria-hidden="true"
              />
              <div class="min-w-0 flex-1">
                <span class="block font-mono-hud text-[10px] text-[#70707d]">{{ item.category }}</span>
                <p class="truncate font-mono text-xs text-[#ededf3]">{{ item.label }}</p>
                <p
                  class="truncate text-[11px]"
                  :class="item.id === 'copy-go' && copiedGoInstall ? 'font-mono font-semibold text-[#22c55e]' : 'text-[#70707d]'"
                >
                  {{ item.hint }}
                </p>
              </div>
            </div>

            <span
              v-if="item.shortcut"
              class="shrink-0 rounded-[2px] border px-2 py-0.5 font-mono text-[10px]"
              :class="
                item.id === 'copy-go' && copiedGoInstall
                  ? 'border-[#22c55e]/60 bg-[#22c55e]/10 text-[#22c55e]'
                  : idx === activeIndex
                    ? 'border-[#5266eb] bg-[#5266eb]/15 text-[#ededf3]'
                    : 'border-[#272735] bg-[#171721] text-[#70707d]'
              "
            >
              {{ item.shortcut }}
            </span>
          </button>
        </div>

        <div class="flex items-center justify-between border-t border-[#1e1e2a] bg-[#06070a] px-4 py-2 font-mono-hud text-[10px] text-[#70707d]">
          <span>↑↓ NAVEGAR · ↵ EXECUTAR</span>
          <span>⌘K / CTRL+K</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
@media (hover: hover) {
  .palette-row:hover {
    background-color: #171721;
  }
}
</style>
