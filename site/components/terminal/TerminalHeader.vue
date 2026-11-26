<script setup lang="ts">
import type { EmulatorMode, HeaderStatus, TuiState } from '../../composables/useSleepoffEmulator'

interface Props {
  mode: EmulatorMode
  tuiState: TuiState
  isDryRun: boolean
  isMuted: boolean
  isAudioPlaying: boolean
  isFocused: boolean
  isOpen: boolean
  headerStatus: HeaderStatus
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'switch-tab', mode: EmulatorMode): void
  (e: 'toggle-mute'): void
  (e: 'toggle-focus'): void
  (e: 'exit-tui'): void
}>()
</script>

<template>
  <header
    class="border-b bg-[#0e1016] select-none transition-colors"
    :class="headerStatus.tone === 'panic' ? 'border-[#ef4444]' : 'border-[#1e1e2a]'"
  >
    <div class="flex flex-wrap items-stretch justify-between">
      <!-- Cell 1: Index + Mode Switcher -->
      <div class="flex flex-wrap items-stretch border-b sm:border-b-0 border-[#1e1e2a] divide-x divide-[#1e1e2a] min-w-0">
        <div class="flex items-center gap-2 px-3 py-2 bg-[#171721] font-mono-hud text-[11px] text-[#c3c3cc]">
          <span class="text-[#5266eb] font-bold">01</span>
          <span class="text-[#70707d]">/</span>
          <span>CONSOLE</span>
        </div>

        <button
          type="button"
          @click.stop="emit('switch-tab', 'tui')"
          aria-label="Alternar para interface TUI interativa"
          class="min-h-[44px] min-w-[44px] px-3.5 py-2 font-mono text-xs uppercase tracking-wider transition-colors cursor-pointer flex items-center gap-1.5"
          :class="
            mode === 'tui'
              ? 'bg-[#171721] text-[#ededf3] font-bold border-b-2 border-b-[#5266eb]'
              : 'bg-[#0e1016] text-[#70707d] hover:text-[#ededf3] hover:bg-[#171721]'
          "
        >
          <span>TUI</span>
        </button>

        <button
          type="button"
          @click.stop="emit('switch-tab', 'shell')"
          aria-label="Alternar para linha de comando SHELL"
          class="min-h-[44px] min-w-[44px] px-3.5 py-2 font-mono text-xs uppercase tracking-wider transition-colors cursor-pointer flex items-center gap-1.5"
          :class="
            mode === 'shell'
              ? 'bg-[#171721] text-[#ededf3] font-bold border-b-2 border-b-[#5266eb]'
              : 'bg-[#0e1016] text-[#70707d] hover:text-[#ededf3] hover:bg-[#171721]'
          "
        >
          <span>SHELL</span>
        </button>
      </div>

      <!-- Cell 2: Binary Version + Dry-Run Flag -->
      <div class="hidden lg:flex items-center px-4 border-l border-r border-[#1e1e2a] font-mono text-xs text-[#70707d]">
        <span class="text-[#c3c3cc]">sleepoff v1.0.3</span>
        <span class="ml-1.5 text-[#5266eb]">({{ isDryRun ? '--dry-run' : 'simulado' }})</span>
      </div>

      <!-- Cell 3: Status Pill + Audio + [~] Focar -->
      <div class="flex flex-wrap items-stretch divide-x divide-[#1e1e2a] border-l border-[#1e1e2a] ml-auto">
        <!-- Semantic Status Pill -->
        <div class="flex items-center px-3 min-h-[44px] bg-[#0e1016]">
          <span
            class="px-2 py-1 font-mono text-[11px] font-bold uppercase tracking-wider border rounded-none whitespace-nowrap"
            :class="[
              headerStatus.tone === 'running' && 'border-[#22c55e] text-[#22c55e] bg-[#06070a]',
              headerStatus.tone === 'done' && 'border-[#22c55e] text-[#22c55e] bg-[#06070a]',
              headerStatus.tone === 'panic' && 'border-[#ef4444] text-[#ef4444] bg-[#06070a]',
              headerStatus.tone === 'neutral' &&
                (mode === 'tui'
                  ? 'border-[#5266eb] text-[#ededf3] bg-[#171721]'
                  : 'border-[#272735] text-[#c3c3cc] bg-[#171721]'),
            ]"
          >
            {{ headerStatus.label }}
          </span>
        </div>

        <!-- Audio Mute / Unmute Toggle -->
        <button
          type="button"
          @click.stop="emit('toggle-mute')"
          :aria-label="isMuted ? 'Ativar áudio do emulador' : 'Silenciar áudio do emulador'"
          :title="isMuted ? 'Ativar áudio' : 'Silenciar áudio'"
          class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#0e1016] hover:bg-[#171721] text-[#c3c3cc] hover:text-[#ededf3] font-mono text-xs transition-colors cursor-pointer flex items-center justify-center gap-2"
        >
          <svg
            v-if="isMuted"
            xmlns="http://www.w3.org/2000/svg"
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="square"
            stroke-linejoin="miter"
            aria-hidden="true"
            class="text-[#70707d]"
          >
            <path d="M11 5 6 9H2v6h4l5 4V5Z" />
            <line x1="22" y1="9" x2="16" y2="15" />
            <line x1="16" y1="9" x2="22" y2="15" />
          </svg>
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="square"
            stroke-linejoin="miter"
            aria-hidden="true"
            :class="isAudioPlaying ? 'text-[#5266eb]' : 'text-[#ededf3]'"
          >
            <path d="M11 5 6 9H2v6h4l5 4V5Z" />
            <path d="M15.54 8.46a5 5 0 0 1 0 7.07" />
            <path d="M19.07 4.93a10 10 0 0 1 0 14.14" />
          </svg>
          <span class="hidden sm:inline uppercase text-[11px] tracking-wider">
            {{ isMuted ? 'MUDO' : isAudioPlaying ? 'SOM' : 'ÁUDIO' }}
          </span>
        </button>

        <!-- Global [~] Focus / Toggle Button -->
        <button
          type="button"
          @click.stop="emit('toggle-focus')"
          :aria-pressed="isFocused && isOpen"
          aria-label="Focar ou alternar console tático usando a tecla til"
          class="min-h-[44px] min-w-[44px] px-3 py-2 font-mono text-[11px] uppercase tracking-wider transition-colors cursor-pointer flex items-center gap-1.5"
          :class="
            isFocused && isOpen
              ? 'bg-[#5266eb] text-[#ededf3] font-bold'
              : 'bg-[#0e1016] hover:bg-[#171721] text-[#c3c3cc] hover:text-[#ededf3]'
          "
        >
          <kbd class="px-1.5 py-0.5 text-[11px] font-bold bg-[#06070a] text-[#ededf3] border border-[#272735]">~</kbd>
          <span>{{ isFocused && isOpen ? 'FOCADO' : 'FOCAR' }}</span>
        </button>

        <!-- Exit TUI to Shell when inside TUI -->
        <button
          v-if="mode === 'tui'"
          type="button"
          @click.stop="emit('exit-tui')"
          aria-label="Sair da TUI e voltar para o shell"
          class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#c3c3cc] hover:text-[#ededf3] font-mono text-[11px] uppercase tracking-wider transition-colors cursor-pointer"
        >
          [q] SAIR
        </button>
      </div>
    </div>
  </header>
</template>
