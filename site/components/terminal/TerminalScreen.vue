<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import {
  MENU_OPTIONS,
  formatPanicTime,
  formatRemaining,
  type EmulatorMode,
  type ShellHistoryEntry,
  type TuiState,
} from '../../composables/useSleepoffEmulator'

interface CancelStats {
  elapsedMins: number
  elapsedSecs: number
  totalMins: number
  pct: number
  msg: string
  asciiBar: string
}

interface Props {
  mode: EmulatorMode
  tuiState: TuiState
  isDryRun: boolean
  isFocused: boolean
  shellInput: string
  shellHistory: ShellHistoryEntry[]
  selectedMenuIdx: number
  customMinutes: string
  totalSeconds: number
  remainingSeconds: number
  isPaused: boolean
  eventsLog: string[]
  panicCountdown: number
  cancelPending: boolean
  startTimeStr: string
  finishTimeStr: string
  progressPercent: number
  cancelSummaryStats: CancelStats
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'focus'): void
  (e: 'update:shellInput', value: string): void
  (e: 'update:customMinutes', value: string): void
  (e: 'submit-shell', cmd: string): void
  (e: 'quick-command', cmd: string): void
  (e: 'select-menu', idx: number): void
  (e: 'set-tui-state', state: TuiState): void
  (e: 'submit-custom'): void
  (e: 'toggle-pause'): void
  (e: 'add-five'): void
  (e: 'subtract-five'): void
  (e: 'cancel'): void
  (e: 'trigger-panic'): void
  (e: 'exit-tui'): void
  (e: 'start-15m'): void
  (e: 'register-elements', els: { shellInput: HTMLInputElement | null; body: HTMLDivElement | null }): void
}>()

const localShellInput = ref<HTMLInputElement | null>(null)
const localBody = ref<HTMLDivElement | null>(null)

onMounted(() => {
  emit('register-elements', {
    shellInput: localShellInput.value,
    body: localBody.value,
  })
})

watch([localShellInput, localBody], ([inputEl, bodyEl]) => {
  emit('register-elements', {
    shellInput: inputEl,
    body: bodyEl,
  })
})
</script>

<template>
  <div class="flex flex-col bg-[#06070a]">
    <!-- Screen Output Area -->
    <div
      ref="localBody"
      @click="emit('focus')"
      class="p-4 sm:p-6 min-h-[280px] sm:min-h-[360px] max-h-[480px] overflow-y-auto overflow-x-hidden font-mono text-sm leading-relaxed text-[#ededf3] select-text flex flex-col justify-between cursor-text"
    >
      <!-- SHELL MODE -->
      <div v-if="mode === 'shell'" class="space-y-3 flex-1 flex flex-col justify-between min-w-0">
        <div class="space-y-2.5 min-w-0">
          <div v-for="(item, idx) in shellHistory" :key="idx" class="space-y-1 min-w-0">
            <div v-if="item.type === 'input'" class="flex flex-wrap items-center gap-1 font-mono text-xs sm:text-sm">
              <span class="text-[#c3c3cc] font-bold">user@sleepoff</span>
              <span class="text-[#70707d]">:</span>
              <span class="text-[#5266eb] font-bold">~</span>
              <span class="text-[#70707d] mr-1">$</span>
              <span class="text-[#ededf3] font-semibold break-all">{{ item.cmd }}</span>
            </div>
            <pre
              v-else
              class="font-mono text-xs whitespace-pre-wrap break-words leading-relaxed"
              :class="item.type === 'error' ? 'text-[#ededf3] border-l-2 border-[#5266eb] pl-2.5' : 'text-[#c3c3cc]'"
            >{{ item.output }}</pre>
          </div>

          <!-- Active Prompt -->
          <div class="flex items-center flex-wrap font-mono text-xs sm:text-sm pt-1 cursor-text" @click="emit('focus')">
            <span class="text-[#c3c3cc] font-bold">user@sleepoff</span>
            <span class="text-[#70707d]">:</span>
            <span class="text-[#5266eb] font-bold">~</span>
            <span class="text-[#70707d] mr-2">$</span>
            <span class="text-[#ededf3] whitespace-pre-wrap break-all">{{ shellInput }}</span>
            <span
              class="inline-block w-2 h-4 bg-[#5266eb] align-middle ml-0.5"
              :class="isFocused ? 'animate-pulse' : 'opacity-40'"
            />
            <input
              ref="localShellInput"
              id="terminal-shell-input"
              aria-label="Linha de comando do emulador sleepoff"
              :value="shellInput"
              type="text"
              autocomplete="off"
              autocapitalize="off"
              spellcheck="false"
              @input="emit('update:shellInput', ($event.target as HTMLInputElement).value)"
              @focus="emit('focus')"
              @keydown.enter.prevent="emit('submit-shell', shellInput)"
              class="absolute opacity-0 w-0 h-0 pointer-events-none"
            />
          </div>
        </div>

        <div class="mt-auto pt-3 border-t border-[#1e1e2a] flex flex-wrap items-center justify-between gap-2 text-xs font-mono select-none">
          <div class="flex flex-wrap items-center gap-2 text-[#70707d]">
            <span class="px-1.5 py-0.5 bg-[#171721] text-[#c3c3cc] font-bold text-[10px] border border-[#272735] uppercase">
              SHELL
            </span>
            <span class="text-[11px] text-[#c3c3cc]">
              Dica: digite <code class="text-[#ededf3] bg-[#171721] px-1 border border-[#272735]">sleepoff 15m --dry-run</code> ou use os atalhos abaixo.
            </span>
          </div>
        </div>
      </div>

      <!-- TUI MODE -->
      <div v-else class="flex flex-col items-center justify-center flex-1 my-auto text-center w-full min-w-0">
        <!-- SPLASH -->
        <div v-if="tuiState === 'SPLASH'" class="space-y-4 py-4 w-full max-w-lg mx-auto">
          <div class="flex items-center justify-between border-b border-[#1e1e2a] pb-2 font-mono text-xs">
            <span class="text-[#ededf3] font-bold">sleepoff <span class="text-[#70707d]">v1.0.3</span></span>
            <span v-if="isDryRun" class="px-2 py-0.5 text-[10px] bg-[#171721] text-[#5266eb] border border-[#5266eb] font-bold">
              --DRY-RUN
            </span>
          </div>
          <pre class="text-[8px] sm:text-xs font-mono font-bold text-[#ededf3] tracking-tighter select-none leading-tight overflow-x-auto max-w-full text-left mx-auto">
 ███████╗██╗     ███████╗███████╗██████╗  ██████╗ 
 ██╔════╝██║     ██╔════╝██╔════╝██╔══██╗██╔═══██╗
 ███████╗██║     █████╗  █████╗  ██████╔╝██║   ██║
 ╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ ██║   ██║
 ███████║███████╗███████╗███████╗██║     ╚██████╔╝
 ╚══════╝╚══════╝╚══════╝╚══════╝╚═╝      ╚═════╝ </pre>
          <div class="text-xs text-[#c3c3cc] font-mono">
            Timer de desligamento • v1.0.3 (Go 1.25.6 • Bubble Tea)
          </div>
        </div>

        <!-- MENU -->
        <div v-else-if="tuiState === 'MENU'" class="w-full max-w-md mx-auto space-y-3 text-left py-2">
          <div class="flex items-center justify-between border-b border-[#1e1e2a] pb-2 font-mono text-xs">
            <div class="flex items-center gap-2">
              <span class="text-[#ededf3] font-bold">sleepoff v1.0.3</span>
              <span class="text-[#70707d]">/</span>
              <span class="text-[#c3c3cc] uppercase">SELECIONE O TEMPO</span>
            </div>
            <span v-if="isDryRun" class="px-2 py-0.5 text-[10px] bg-[#171721] text-[#5266eb] border border-[#5266eb] font-bold">
              --DRY-RUN
            </span>
          </div>

          <div class="border border-[#1e1e2a] divide-y divide-[#1e1e2a] font-mono text-xs">
            <button
              v-for="(opt, idx) in MENU_OPTIONS"
              :key="opt.title"
              type="button"
              @click.stop="emit('select-menu', idx)"
              class="w-full min-h-[44px] px-3 py-2 flex items-center justify-between text-left transition-colors cursor-pointer"
              :class="
                selectedMenuIdx === idx
                  ? 'bg-[#171721] text-[#ededf3] font-bold border-l-2 border-l-[#5266eb]'
                  : 'bg-[#0e1016] text-[#c3c3cc] hover:bg-[#171721] hover:text-[#ededf3]'
              "
            >
              <span class="flex items-center gap-2.5">
                <span class="text-[#5266eb] font-bold w-3">{{ selectedMenuIdx === idx ? '▶' : ' ' }}</span>
                <span>{{ opt.title }}</span>
              </span>
              <span class="text-[11px] text-[#70707d]">{{ opt.desc }}</span>
            </button>
          </div>
        </div>

        <!-- CUSTOM INPUT -->
        <div v-else-if="tuiState === 'CUSTOM_INPUT'" class="w-full max-w-sm mx-auto space-y-4 py-4 text-center font-mono">
          <div class="flex items-center justify-between border-b border-[#1e1e2a] pb-2 text-xs">
            <span class="text-[#ededf3] font-bold">TEMPO PERSONALIZADO</span>
            <span v-if="isDryRun" class="px-2 py-0.5 text-[10px] bg-[#171721] text-[#5266eb] border border-[#5266eb] font-bold">
              --DRY-RUN
            </span>
          </div>

          <p class="text-xs text-[#c3c3cc]">Quantos minutos antes de desligar? (1 a 1440)</p>
          <div class="flex items-center justify-center gap-2">
            <span class="text-[#5266eb] font-bold">❯</span>
            <input
              id="tui-custom-minutes-input"
              aria-label="Minutos personalizados para contagem regressiva"
              :value="customMinutes"
              type="number"
              min="1"
              max="1440"
              @input="emit('update:customMinutes', ($event.target as HTMLInputElement).value)"
              @keydown.enter.prevent="emit('submit-custom')"
              @keydown.esc.prevent="emit('set-tui-state', 'MENU')"
              class="w-28 min-h-[44px] bg-[#0e1016] border border-[#5266eb] px-3 py-2 text-center text-lg font-bold text-[#ededf3] outline-none"
            />
            <span class="text-[#70707d] font-bold">min</span>
          </div>

          <div class="flex flex-wrap justify-center gap-2 pt-2 border-t border-[#1e1e2a]">
            <button
              type="button"
              @click.stop="emit('set-tui-state', 'MENU')"
              class="min-h-[44px] min-w-[44px] px-4 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#c3c3cc] text-xs cursor-pointer border border-[#272735]"
            >
              Voltar (Esc)
            </button>
            <button
              type="button"
              @click.stop="emit('submit-custom')"
              class="min-h-[44px] min-w-[44px] px-5 py-2 bg-[#5266eb] hover:opacity-90 text-[#ededf3] font-bold text-xs cursor-pointer"
            >
              Iniciar (Enter)
            </button>
          </div>
        </div>

        <!-- RUNNING / PAUSED -->
        <div v-else-if="tuiState === 'RUNNING'" class="w-full max-w-lg mx-auto space-y-4 py-2">
          <div class="flex flex-wrap items-center justify-between gap-2 border-b border-[#1e1e2a] pb-2.5 font-mono text-xs">
            <div class="flex items-center gap-2">
              <span class="text-[#ededf3] font-bold">sleepoff v1.0.3</span>
              <span class="text-[#70707d]">/</span>
              <span class="text-[#c3c3cc]">{{ formatRemaining(totalSeconds) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span
                class="px-2 py-0.5 text-[10px] font-bold border uppercase"
                :class="
                  isPaused
                    ? 'border-[#272735] bg-[#171721] text-[#c3c3cc]'
                    : 'border-[#22c55e] bg-[#06070a] text-[#22c55e]'
                "
              >
                {{ isPaused ? 'PAUSADO' : '• RODANDO' }}
              </span>
              <span v-if="isDryRun" class="px-2 py-0.5 text-[10px] bg-[#171721] text-[#5266eb] border border-[#5266eb] font-bold">
                --DRY-RUN
              </span>
            </div>
          </div>

          <div class="py-2 text-center font-mono">
            <div class="text-4xl sm:text-6xl font-bold text-[#ededf3] tracking-wider tabular-nums">
              {{ formatRemaining(remainingSeconds) }}
            </div>
            <div class="text-xs text-[#70707d] mt-2 flex flex-wrap items-center justify-center gap-x-3 gap-y-1">
              <span>Início: <strong class="text-[#c3c3cc]">{{ startTimeStr }}</strong></span>
              <span>/</span>
              <span>Previsão: <strong class="text-[#ededf3]">{{ finishTimeStr }}</strong></span>
              <span>/</span>
              <span>Progresso: <strong class="text-[#5266eb]">{{ progressPercent }}%</strong></span>
            </div>
          </div>

          <!-- 1px Hairline Progress Bar -->
          <div class="w-full h-2.5 bg-[#171721] border border-[#272735] overflow-hidden">
            <div
              class="h-full bg-[#5266eb] transition-all duration-300"
              :style="{ width: `${progressPercent}%` }"
            />
          </div>

          <div
            v-if="cancelPending"
            class="p-2.5 bg-[#171721] border border-[#5266eb] text-[#ededf3] text-xs font-bold text-center font-mono"
          >
            Confirmação dupla: pressione <kbd class="px-1.5 py-0.5 mx-0.5">c</kbd> ou <kbd class="px-1.5 py-0.5 mx-0.5">Esc</kbd> novamente para cancelar.
          </div>

          <div class="bg-[#0e1016] p-3 border border-[#1e1e2a] text-left text-xs font-mono space-y-1">
            <div class="text-[#70707d] font-bold border-b border-[#1e1e2a] pb-1 text-[11px] uppercase tracking-wider">
              Eventos em tempo real
            </div>
            <div v-for="(ev, idx) in eventsLog" :key="idx" class="text-[#c3c3cc] text-[11px] truncate">
              <span class="text-[#5266eb] mr-1.5">│</span>{{ ev }}
            </div>
          </div>
        </div>

        <!-- CONFIRMATION (10s PANIC) -->
        <div
          v-else-if="tuiState === 'CONFIRMATION'"
          class="w-full max-w-lg mx-auto space-y-4 py-4 px-4 sm:px-6 border border-[#ef4444] bg-[#ef4444]/10"
        >
          <div class="flex items-center justify-between border-b border-[#ef4444]/60 pb-2 font-mono text-xs">
            <span class="text-[#ef4444] font-bold uppercase tracking-wider">! PÂNICO 10s — DESLIGAMENTO IMINENTE</span>
            <span v-if="isDryRun" class="px-2 py-0.5 text-[10px] bg-[#06070a] text-[#ededf3] border border-[#ef4444] font-bold">
              --DRY-RUN
            </span>
          </div>

          <div class="text-center py-1 font-mono">
            <div class="text-[#ef4444] text-xs sm:text-sm font-bold tracking-wider">
              DESLIGANDO EM {{ String(Math.max(0, panicCountdown)).padStart(2, '0') }} SEGUNDOS
            </div>
            <div class="text-5xl sm:text-7xl font-bold text-[#ef4444] tracking-widest tabular-nums my-2">
              {{ formatPanicTime(panicCountdown) }}
            </div>
            <p class="text-xs text-[#c3c3cc]">
              Modo simulado ativo — aperte [C] ou Esc para abortar o alarme agora.
            </p>
          </div>

          <div class="w-full h-2.5 bg-[#06070a] border border-[#ef4444] overflow-hidden">
            <div
              class="h-full bg-[#ef4444] transition-all duration-300"
              :style="{ width: `${(Math.max(0, panicCountdown) / 10) * 100}%` }"
            />
          </div>

          <div class="pt-2 border-t border-[#ef4444]/50">
            <button
              type="button"
              @click.stop="emit('cancel')"
              class="w-full min-h-[44px] py-3 px-4 bg-[#ef4444] hover:opacity-90 text-[#ededf3] font-mono font-bold text-xs sm:text-sm uppercase tracking-wider cursor-pointer flex items-center justify-center gap-2"
            >
              <kbd class="px-2 py-0.5 bg-[#06070a] text-[#ededf3] border border-[#272735] text-xs font-bold">C</kbd>
              <span>ABORTAR ALARME AGORA (ou Esc)</span>
            </button>
          </div>
        </div>

        <!-- DONE_TEST -->
        <div v-else-if="tuiState === 'DONE_TEST'" class="w-full max-w-xl mx-auto space-y-4 py-3 text-center">
          <div class="inline-flex items-center gap-2 px-3 py-1 border border-[#22c55e] bg-[#06070a] text-[#22c55e] font-mono text-xs font-bold">
            ✓ FEITO — SIMULAÇÃO CONCLUÍDA
          </div>
          <pre class="text-[7.5px] sm:text-xs font-mono font-bold text-[#ededf3] leading-tight select-none overflow-x-auto max-w-full text-left mx-auto">
 ╔═══════════════════════════════════════════════════════╗
 ║              TIMER FINALIZADO (--DRY-RUN)             ║
 ╚═══════════════════════════════════════════════════════╝</pre>
          <p class="text-xs text-[#c3c3cc] font-mono max-w-md mx-auto">
            O computador NÃO foi desligado (modo seguro --dry-run). Em produção, o sistema executaria o desligamento limpo.
          </p>
          <div class="flex flex-wrap justify-center gap-2 pt-1">
            <button
              type="button"
              @click.stop="emit('exit-tui')"
              class="min-h-[44px] min-w-[44px] px-4 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] text-xs font-mono font-bold border border-[#272735] cursor-pointer"
            >
              Voltar ao Shell (Enter)
            </button>
            <button
              type="button"
              @click.stop="emit('start-15m')"
              class="min-h-[44px] min-w-[44px] px-4 py-2 bg-[#5266eb] hover:opacity-90 text-[#ededf3] text-xs font-mono font-bold cursor-pointer"
            >
              Rodar 15m (--dry-run)
            </button>
          </div>
        </div>

        <!-- CANCELLED -->
        <div v-else-if="tuiState === 'CANCELLED'" class="w-full max-w-xl mx-auto space-y-4 py-3 text-center">
          <div class="inline-flex items-center gap-2 px-3 py-1 border border-[#22c55e] bg-[#06070a] text-[#22c55e] font-mono text-xs font-bold">
            ✓ FEITO — ALARME ABORTADO COM SEGURANÇA
          </div>
          <div class="text-xs font-mono text-[#c3c3cc] px-3 py-2 border border-[#1e1e2a] bg-[#0e1016] max-w-md mx-auto">
            Alarme interrompido em modo seguro (--dry-run). Nenhum comando de energia foi executado.
          </div>

          <div class="text-xs font-mono max-w-full overflow-x-auto">
            <span class="text-[#70707d]">Progresso: </span>
            <span class="text-[#5266eb] font-bold">[{{ cancelSummaryStats.asciiBar }}]</span>
            <span class="text-[#ededf3] font-bold ml-1.5">{{ cancelSummaryStats.pct }}%</span>
          </div>

          <div class="inline-block bg-[#0e1016] border border-[#1e1e2a] p-3 text-left font-mono text-xs text-[#c3c3cc] min-w-[240px] max-w-full">
            <div class="flex justify-between gap-4 py-0.5">
              <span class="text-[#70707d]">Tempo decorrido:</span>
              <span class="text-[#ededf3] font-bold">{{ String(cancelSummaryStats.elapsedMins).padStart(2, '0') }}m {{ String(cancelSummaryStats.elapsedSecs).padStart(2, '0') }}s</span>
            </div>
            <div class="flex justify-between gap-4 py-0.5">
              <span class="text-[#70707d]">Tempo total:</span>
              <span class="text-[#ededf3] font-bold">{{ String(cancelSummaryStats.totalMins).padStart(2, '0') }}m</span>
            </div>
          </div>

          <div class="text-xs font-bold text-[#c3c3cc] font-mono">
            {{ cancelSummaryStats.msg }}
          </div>

          <div class="flex flex-wrap justify-center gap-2 pt-1">
            <button
              type="button"
              @click.stop="emit('exit-tui')"
              class="min-h-[44px] min-w-[44px] px-4 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] text-xs font-mono font-bold border border-[#272735] cursor-pointer"
            >
              Voltar ao Shell (Enter)
            </button>
            <button
              type="button"
              @click.stop="emit('start-15m')"
              class="min-h-[44px] min-w-[44px] px-4 py-2 bg-[#5266eb] hover:opacity-90 text-[#ededf3] text-xs font-mono font-bold cursor-pointer"
            >
              Rodar 15m (--dry-run)
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Contextual Action Controls Strip (1px hairline border) -->
    <div class="px-3 sm:px-4 py-2.5 bg-[#0e1016] border-t border-[#1e1e2a] flex flex-wrap items-center justify-between gap-2 text-xs select-none">
      <!-- Shell Quick Actions -->
      <template v-if="mode === 'shell'">
        <span class="text-[#70707d] font-mono-hud text-[10px]">Atalhos rápidos:</span>
        <div class="flex flex-wrap items-center gap-1.5">
          <button
            type="button"
            @click.stop="emit('quick-command', 'sleepoff 15m --dry-run')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#5266eb] hover:opacity-90 text-[#ededf3] font-mono font-bold text-xs cursor-pointer"
          >
            sleepoff 15m --dry-run
          </button>
          <button
            type="button"
            @click.stop="emit('quick-command', 'sleepoff 30m --dry-run')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#c3c3cc] border border-[#272735] font-mono text-xs cursor-pointer"
          >
            30m
          </button>
          <button
            type="button"
            @click.stop="emit('quick-command', 'sleepoff 1h --dry-run')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#c3c3cc] border border-[#272735] font-mono text-xs cursor-pointer"
          >
            1h
          </button>
          <button
            type="button"
            @click.stop="emit('quick-command', 'sleepoff --dry-run')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer"
          >
            Menu TUI
          </button>
          <button
            type="button"
            @click.stop="emit('trigger-panic')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer"
          >
            Simular Pânico (10s)
          </button>
          <button
            type="button"
            @click.stop="emit('quick-command', 'sleepoff --help')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#70707d] hover:text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer"
          >
            --help
          </button>
          <button
            type="button"
            @click.stop="emit('quick-command', 'clear')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#70707d] hover:text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer"
          >
            clear
          </button>
        </div>
      </template>

      <!-- Running TUI Controls -->
      <template v-else-if="tuiState === 'RUNNING'">
        <div class="flex flex-wrap items-center gap-1.5">
          <button
            type="button"
            @click.stop="emit('toggle-pause')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer flex items-center gap-1.5"
          >
            <kbd class="px-1.5 py-0.5 text-[10px] font-bold">p</kbd>
            <span>{{ isPaused ? 'Retomar' : 'Pausar' }}</span>
          </button>
          <button
            type="button"
            @click.stop="emit('add-five')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer flex items-center gap-1.5"
          >
            <kbd class="px-1.5 py-0.5 text-[10px] font-bold">+</kbd>
            <span>+5min</span>
          </button>
          <button
            type="button"
            @click.stop="emit('subtract-five')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#c3c3cc] border border-[#272735] font-mono text-xs cursor-pointer flex items-center gap-1.5"
          >
            <kbd class="px-1.5 py-0.5 text-[10px] font-bold">-</kbd>
            <span>-5min</span>
          </button>
          <button
            type="button"
            @click.stop="emit('cancel')"
            class="min-h-[44px] min-w-[44px] px-3 py-2 font-mono text-xs cursor-pointer border flex items-center gap-1.5"
            :class="cancelPending ? 'bg-[#5266eb] text-[#ededf3] border-[#5266eb] font-bold' : 'bg-[#171721] hover:bg-[#1e1e2a] text-[#c3c3cc] border-[#272735]'"
          >
            <kbd class="px-1.5 py-0.5 text-[10px] font-bold">c</kbd>
            <span>{{ cancelPending ? 'Confirmar Cancelamento' : 'Cancelar (2x)' }}</span>
          </button>
        </div>

        <button
          type="button"
          @click.stop="emit('trigger-panic')"
          class="min-h-[44px] min-w-[44px] px-3.5 py-2 bg-[#5266eb] hover:opacity-90 text-[#ededf3] font-mono font-bold text-xs cursor-pointer"
        >
          Simular Pânico (10s)
        </button>
      </template>

      <!-- Panic Controls -->
      <template v-else-if="tuiState === 'CONFIRMATION'">
        <span class="text-[#ef4444] font-mono font-bold text-xs uppercase tracking-wider">
          ! PÂNICO 10s ATIVO — ABORTE PARA TESTAR
        </span>
        <button
          type="button"
          @click.stop="emit('cancel')"
          class="min-h-[44px] min-w-[44px] px-4 py-2 bg-[#ef4444] hover:opacity-90 text-[#ededf3] font-mono font-bold text-xs cursor-pointer"
        >
          [C] ABORTAR ALARME
        </button>
      </template>

      <!-- Menu / Exit Fallback Strip -->
      <template v-else>
        <span class="text-[#70707d] font-mono text-xs">
          Use ↑/↓ + Enter para escolher ou volte para o shell.
        </span>
        <button
          type="button"
          @click.stop="emit('exit-tui')"
          class="min-h-[44px] min-w-[44px] px-3.5 py-2 bg-[#171721] hover:bg-[#1e1e2a] text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer"
        >
          Voltar ao Shell (q)
        </button>
      </template>
    </div>
  </div>
</template>
