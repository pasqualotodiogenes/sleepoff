<script setup lang="ts">
import { ClientOnly, Icon } from '#components'
import EncryptedText from '~/components/ui/encrypted-text/EncryptedText.vue'
import { useSleepoffSiteState } from '~/composables/useSleepoffSiteState'

const { triggerCommandInEmulator } = useSleepoffSiteState()

interface CliEntry {
  token: string
  description: string
  isPanicAction?: boolean
}

const durationCommands: CliEntry[] = [
  { token: 'sleepoff', description: 'Abre o menu TUI interativo (15m, 30m, 45m, 1h, 1h30m, 2h ou Personalizado)' },
  { token: 'sleepoff 30', description: 'Inteiro sem sufixo — interpretado automaticamente como 30 minutos' },
  { token: 'sleepoff 30m', description: 'Minutos explícitos com sufixo m' },
  { token: 'sleepoff 90s', description: 'Segundos explícitos com sufixo s' },
  { token: 'sleepoff 1h30m', description: 'Formato composto Go time.ParseDuration (1 hora e 30 minutos)' },
]

const cliFlags: CliEntry[] = [
  { token: '--dry-run', description: 'Modo seguro (sem flag curta): simula timer, alarme e pânico 10s sem desligar' },
  { token: '-h / --help', description: 'Exibe banner, exemplos de duração, flags e controles da TUI' },
  { token: '-v / --version', description: 'Exibe versão v1.0.3, hash de commit e metadados de compilação' },
]

const tuiControls: CliEntry[] = [
  { token: 'p / Espaço', description: 'Pausar ou retomar cronômetro em execução (StateRunning)' },
  { token: '+ / =', description: 'Adicionar 5 minutos ao tempo restante e duração total' },
  { token: '- / _', description: 'Subtrair 5 minutos (mantém piso mínimo de 1 minuto)' },
  { token: 'c / q / Esc (2x)', description: 'Cancelar timer ativo (duas ativações dentro de 2 segundos)' },
  { token: 'c / Esc (1x)', description: 'Abortar desligamento imediatamente na Tela de Pânico (10s)', isPanicAction: true },
  { token: '↑ ↓ / j k + Enter', description: 'Navegar e confirmar presets no Menu Principal' },
  { token: 'Ctrl + C', description: 'Interrupção imediata de emergência via sinal do terminal' },
]

function runInEmulator(cmd: string) {
  triggerCommandInEmulator(cmd)
  if (!import.meta.client) return
  const target = document.querySelector('#emulador')
  if (!target) return
  const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  target.scrollIntoView({ behavior: reducedMotion ? 'auto' : 'smooth', block: 'start' })
}
</script>

<template>
  <section
    id="cli"
    data-section="04"
    class="relative z-10 border-b border-[#1e1e2a] bg-[#06070a]"
  >
    <div class="mx-auto max-w-6xl border-x border-[#1e1e2a] px-4 py-16 sm:px-8 sm:py-20 lg:px-12">
      <!-- Header -->
      <div class="mb-8 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <div class="mb-3">
            <ClientOnly>
              <EncryptedText
                text="04 / CLI"
                class="font-mono-hud text-xs text-[#5266eb]"
                encrypted-class="text-[#70707d]"
                revealed-class="text-[#5266eb]"
              />
              <template #fallback>
                <span class="font-mono-hud text-xs text-[#5266eb]">04 / CLI</span>
              </template>
            </ClientOnly>
          </div>
          <h2 class="font-display text-ms-4xl sm:text-ms-5xl font-bold tracking-tighter text-[#ededf3]">
            Referência de <span class="font-serif-accent">comandos</span> &amp; controles
          </h2>
        </div>

        <!-- Quick Emulator Launchers -->
        <div class="flex flex-wrap items-center gap-2">
          <button
            type="button"
            class="inline-flex min-h-[44px] cursor-pointer items-center gap-2 border border-[#5266eb] bg-[#5266eb] px-4 py-2 font-mono text-xs font-semibold text-[#ededf3] transition-colors hover:bg-[#5266eb]/90"
            @click="runInEmulator('sleepoff 15m --dry-run')"
          >
            <Icon name="heroicons:play" class="size-3.5 shrink-0" />
            <span>Testar sleepoff 15m --dry-run</span>
          </button>
          <button
            type="button"
            class="inline-flex min-h-[44px] cursor-pointer items-center gap-2 border border-[#272735] bg-[#171721] px-4 py-2 font-mono text-xs text-[#c3c3cc] transition-colors hover:border-[#5266eb] hover:text-[#ededf3]"
            @click="runInEmulator('sleepoff')"
          >
            <span>sleepoff</span>
          </button>
          <button
            type="button"
            class="inline-flex min-h-[44px] cursor-pointer items-center gap-2 border border-[#272735] bg-[#171721] px-4 py-2 font-mono text-xs text-[#c3c3cc] transition-colors hover:border-[#5266eb] hover:text-[#ededf3]"
            @click="runInEmulator('sleepoff --help')"
          >
            <span>--help</span>
          </button>
        </div>
      </div>

      <!-- 1px Hairline Reference Matrix -->
      <div class="grid grid-cols-1 gap-px border border-[#1e1e2a] bg-[#1e1e2a] lg:grid-cols-2">
        <!-- Left Column: Commands & Flags -->
        <div class="flex flex-col divide-y divide-[#1e1e2a] bg-[#0e1016]">
          <div class="border-b border-[#1e1e2a] bg-[#171721] px-5 py-3.5">
            <h3 class="font-mono-hud text-xs text-[#ededf3]">
              FORMATOS DE DURAÇÃO &amp; FLAGS (v1.0.3)
            </h3>
          </div>

          <div
            v-for="item in durationCommands"
            :key="item.token"
            class="flex flex-col gap-1.5 px-5 py-3.5 sm:flex-row sm:items-center sm:justify-between sm:gap-4"
          >
            <code class="shrink-0 font-mono text-xs font-semibold text-[#5266eb]">
              {{ item.token }}
            </code>
            <span class="text-xs leading-relaxed text-[#c3c3cc] sm:text-right">
              {{ item.description }}
            </span>
          </div>

          <div class="border-t border-[#1e1e2a] bg-[#171721] px-5 py-2.5">
            <span class="font-mono-hud text-[11px] text-[#70707d]">FLAGS SUPORTADAS</span>
          </div>

          <div
            v-for="flag in cliFlags"
            :key="flag.token"
            class="flex flex-col gap-1.5 px-5 py-3.5 sm:flex-row sm:items-center sm:justify-between sm:gap-4"
          >
            <code class="shrink-0 font-mono text-xs font-semibold text-[#ededf3]">
              {{ flag.token }}
            </code>
            <span class="text-xs leading-relaxed text-[#c3c3cc] sm:text-right">
              {{ flag.description }}
            </span>
          </div>
        </div>

        <!-- Right Column: TUI Controls -->
        <div class="flex flex-col divide-y divide-[#1e1e2a] bg-[#0e1016]">
          <div class="border-b border-[#1e1e2a] bg-[#171721] px-5 py-3.5">
            <h3 class="font-mono-hud text-xs text-[#ededf3]">
              CONTROLES DE TECLADO NA TUI
            </h3>
          </div>

          <div
            v-for="ctrl in tuiControls"
            :key="ctrl.token"
            class="flex flex-col gap-2 px-5 py-3.5 sm:flex-row sm:items-center sm:justify-between sm:gap-4"
          >
            <kbd
              class="w-fit shrink-0 px-2.5 py-1 font-mono text-xs font-semibold"
              :class="
                ctrl.isPanicAction
                  ? 'border-[#ef4444]/60 bg-[#ef4444]/10 text-[#ef4444]'
                  : 'border-[#272735] bg-[#272735] text-[#ededf3]'
              "
            >
              {{ ctrl.token }}
            </kbd>
            <span class="text-xs leading-relaxed text-[#c3c3cc] sm:text-right">
              {{ ctrl.description }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
