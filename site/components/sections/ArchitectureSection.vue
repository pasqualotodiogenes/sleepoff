<script setup lang="ts">
import { ClientOnly, Icon } from '#components'
import CardSpotlight from '~/components/ui/card-spotlight/CardSpotlight.vue'
import EncryptedText from '~/components/ui/encrypted-text/EncryptedText.vue'
import NumberTicker from '~/components/ui/number-ticker/NumberTicker.vue'

interface TelemetryMetric {
  value: number
  suffix: string
  label: string
}

const telemetryMetrics: TelemetryMetric[] = [
  { value: 0, suffix: '', label: 'Dependências runtime externas' },
  { value: 50, suffix: ' ms', label: 'Tick determinístico sem drift' },
  { value: 10, suffix: ' s', label: 'Janela de pânico com alarme' },
  { value: 1, suffix: '', label: 'Instância ativa por máquina' },
]
</script>

<template>
  <section
    id="arquitetura"
    data-section="03"
    class="relative z-10 border-b border-[#1e1e2a] bg-[#06070a]"
  >
    <div class="mx-auto max-w-6xl border-x border-[#1e1e2a] px-4 py-16 sm:px-8 sm:py-20 lg:px-12">
      <!-- Header -->
      <div class="mb-10 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <div class="mb-3 flex flex-wrap items-center gap-2.5">
            <ClientOnly>
              <EncryptedText
                text="03 / ARQUITETURA"
                class="font-mono-hud text-xs text-[#5266eb]"
                encrypted-class="text-[#70707d]"
                revealed-class="text-[#5266eb]"
              />
              <template #fallback>
                <span class="font-mono-hud text-xs text-[#5266eb]">03 / ARQUITETURA</span>
              </template>
            </ClientOnly>
            <span class="font-mono-hud border border-[#1e1e2a] bg-[#171721] px-2 py-0.5 text-[10px] text-[#c3c3cc]">
              03 / ZERO BLOAT
            </span>
          </div>
          <h2 class="font-display text-ms-4xl sm:text-ms-5xl font-bold tracking-tighter text-[#ededf3]">
            Engenharia nativa, <span class="font-serif-accent">sem</span> desperdício
          </h2>
        </div>
        <p class="max-w-md text-ms-sm leading-relaxed text-[#c3c3cc]">
          Construído em Go 1.25+ com primitivas nativas de sistema operacional para Windows, Linux e macOS.
        </p>
      </div>

      <!-- Telemetry Strip -->
      <div class="mb-8 grid grid-cols-2 gap-px border border-[#1e1e2a] bg-[#1e1e2a] lg:grid-cols-4">
        <div
          v-for="metric in telemetryMetrics"
          :key="metric.label"
          class="bg-[#0e1016] p-4 sm:p-5"
        >
          <div class="font-display text-ms-3xl sm:text-ms-4xl font-bold tracking-tight text-[#ededf3]">
            <ClientOnly>
              <NumberTicker
                :value="metric.value"
                :decimal-places="0"
                class="text-[#ededf3]"
              />
              <template #fallback>
                <span class="tabular-nums text-[#ededf3]">{{ metric.value }}</span>
              </template>
            </ClientOnly>
            <span v-if="metric.suffix" class="font-mono text-base font-medium text-[#5266eb]">
              {{ metric.suffix }}
            </span>
          </div>
          <p class="mt-1 font-mono text-[11px] uppercase tracking-[0.06em] text-[#70707d]">
            {{ metric.label }}
          </p>
        </div>
      </div>

      <!-- 4 CardSpotlight Grid (1px hairline dividers, unclipped previews) -->
      <div class="grid grid-cols-1 gap-px border border-[#1e1e2a] bg-[#1e1e2a] md:grid-cols-2">
        <!-- Card 1: Binário Estático & App Nativo -->
        <CardSpotlight
          class="h-full border-0 bg-[#0e1016] p-5 sm:p-6"
          slot-class="flex h-full flex-col justify-between gap-5"
        >
          <div class="min-h-[11rem] h-auto w-full overflow-x-auto border border-[#1e1e2a] bg-[#06070a] p-4 font-mono text-xs">
            <div class="mb-3 flex items-center justify-between gap-2 border-b border-[#1e1e2a] pb-2 text-[#c3c3cc]">
              <span class="whitespace-nowrap text-[#5266eb]">$ CGO_ENABLED=0 go build -ldflags="-s -w"</span>
              <span class="shrink-0 text-[10px] text-[#70707d]">STATIC</span>
            </div>
            <div class="space-y-2 text-[11px] leading-relaxed">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">Windows Desktop:</span>
                <span class="text-[#ededf3]">sleepoff-setup.exe + lxn/walk Tray</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">Linux / macOS:</span>
                <span class="text-[#ededf3]">Binário único (amd64 / arm64)</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">MousetrapHelpText:</span>
                <span class="text-[#5266eb]">Duplo clique habilitado</span>
              </div>
            </div>
          </div>

          <div>
            <div class="mb-2 flex items-center gap-2">
              <Icon name="heroicons:cpu-chip" class="size-4 shrink-0 text-[#5266eb]" />
              <h3 class="font-display text-ms-lg font-bold tracking-tight text-[#ededf3]">
                Binário Estático &amp; App Nativo
              </h3>
            </div>
            <p class="text-ms-sm leading-relaxed text-[#c3c3cc]">
              Compilação Go pura — de 0 dependências a instalador com System Tray, atalhos no Menu Iniciar e atualização automática semanal no Windows.
            </p>
          </div>
        </CardSpotlight>

        <!-- Card 2: Single-Instance Lock (unclipped vertical height) -->
        <CardSpotlight
          class="h-full border-0 bg-[#0e1016] p-5 sm:p-6"
          slot-class="flex h-full flex-col justify-between gap-5"
        >
          <div class="min-h-[11rem] h-auto w-full overflow-x-auto border border-[#1e1e2a] bg-[#06070a] p-4 font-mono text-xs">
            <div class="mb-3 flex items-center justify-between gap-2 border-b border-[#1e1e2a] pb-2">
              <span class="flex items-center gap-1.5 text-[#ededf3]">
                <Icon name="heroicons:lock-closed" class="size-3.5 shrink-0 text-[#5266eb]" />
                <span>Mutex / flock Atômico</span>
              </span>
              <span class="border border-[#5266eb]/50 bg-[#5266eb]/15 px-1.5 py-0.5 text-[10px] text-[#ededf3]">
                LOCK ACTIVE
              </span>
            </div>
            <div class="space-y-2 text-[11px] leading-relaxed">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">Win32:</span>
                <span class="text-[#ededf3]">CreateMutexW("sleepoff-single-instance")</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">Linux:</span>
                <span class="text-[#c3c3cc]">flock(LOCK_EX|LOCK_NB) → $XDG_RUNTIME_DIR/sleepoff.lock</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">macOS:</span>
                <span class="text-[#c3c3cc]">flock(LOCK_EX|LOCK_NB) → /tmp/sleepoff.lock</span>
              </div>
              <div class="border-t border-[#1e1e2a] pt-2 text-[#5266eb]">
                ● Impede dois timers concorrentes na mesma sessão
              </div>
            </div>
          </div>

          <div>
            <div class="mb-2 flex items-center gap-2">
              <Icon name="heroicons:lock-closed" class="size-4 shrink-0 text-[#5266eb]" />
              <h3 class="font-display text-ms-lg font-bold tracking-tight text-[#ededf3]">
                Single-Instance Lock
              </h3>
            </div>
            <p class="text-ms-sm leading-relaxed text-[#c3c3cc]">
              <code class="text-xs text-[#ededf3]">CreateMutexW</code> no Windows e <code class="text-xs text-[#ededf3]">flock</code> em
              <code class="text-xs text-[#ededf3]">$XDG_RUNTIME_DIR/sleepoff.lock</code> no Linux/macOS — impede dois timers concorrentes.
            </p>
          </div>
        </CardSpotlight>

        <!-- Card 3: Modo Seguro (--dry-run) (unclipped vertical height) -->
        <CardSpotlight
          class="h-full border-0 bg-[#0e1016] p-5 sm:p-6"
          slot-class="flex h-full flex-col justify-between gap-5"
        >
          <div class="min-h-[11rem] h-auto w-full overflow-x-auto border border-[#1e1e2a] bg-[#06070a] p-4 font-mono text-xs">
            <div class="mb-3 flex items-center justify-between gap-2 border-b border-[#1e1e2a] pb-2">
              <span class="text-[#ededf3]">$ sleepoff 15m --dry-run</span>
              <span class="border border-[#272735] bg-[#171721] px-1.5 py-0.5 text-[10px] text-[#5266eb]">
                SIMULATION
              </span>
            </div>
            <div class="space-y-2 text-[11px] leading-relaxed">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#c3c3cc]">1. StateRunning + Tick 50ms</span>
                <span class="text-[#5266eb]">EXECUTADO</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#c3c3cc]">2. Tela de Pânico (10s) + countdown.wav</span>
                <span class="text-[#5266eb]">EXECUTADO</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#c3c3cc]">3. Syscall shutdown / poweroff</span>
                <span class="border border-[#5266eb]/50 bg-[#5266eb]/10 px-1.5 py-0.5 text-[#ededf3]">
                  BLOQUEADO (--dry-run)
                </span>
              </div>
            </div>
          </div>

          <div>
            <div class="mb-2 flex items-center gap-2">
              <Icon name="heroicons:shield-check" class="size-4 shrink-0 text-[#5266eb]" />
              <h3 class="font-display text-ms-lg font-bold tracking-tight text-[#ededf3]">
                Modo Seguro (<code class="font-mono text-base">--dry-run</code>)
              </h3>
            </div>
            <p class="text-ms-sm leading-relaxed text-[#c3c3cc]">
              Testa todo o fluxo de timer, alarme sonoro e tela de pânico de 10 segundos sem desligar o computador.
            </p>
          </div>
        </CardSpotlight>

        <!-- Card 4: Alarme Sonoro + TUI Reativa (50ms) -->
        <CardSpotlight
          class="h-full border-0 bg-[#0e1016] p-5 sm:p-6"
          slot-class="flex h-full flex-col justify-between gap-5"
        >
          <div class="min-h-[11rem] h-auto w-full overflow-x-auto border border-[#1e1e2a] bg-[#06070a] p-4 font-mono text-xs">
            <div class="mb-3 flex items-center justify-between gap-2 border-b border-[#1e1e2a] pb-2">
              <span class="text-[#ededf3]">//go:embed countdown.wav</span>
              <span class="text-[10px] text-[#5266eb]">TICK 50ms</span>
            </div>
            <div class="space-y-2 text-[11px] leading-relaxed">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">Precisão temporal:</span>
                <span class="text-[#ededf3]">time.Until(FinishTime) sem drift</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[#70707d]">Drivers de áudio:</span>
                <span class="text-[#c3c3cc]">winmm.dll · pw-play/aplay · afplay</span>
              </div>
              <div class="flex flex-wrap items-center justify-between gap-2 border-t border-[#1e1e2a] pt-2">
                <span class="text-[#70707d]">Atalhos TUI:</span>
                <span class="text-[#ededf3]">p · + · - · c/Esc</span>
              </div>
            </div>
          </div>

          <div>
            <div class="mb-2 flex items-center gap-2">
              <Icon name="heroicons:speaker-wave" class="size-4 shrink-0 text-[#5266eb]" />
              <h3 class="font-display text-ms-lg font-bold tracking-tight text-[#ededf3]">
                Alarme Sonoro + TUI Reativa (50ms)
              </h3>
            </div>
            <p class="text-ms-sm leading-relaxed text-[#c3c3cc]">
              <code class="text-xs text-[#ededf3]">countdown.wav</code> embarcado no binário, cálculo por <code class="text-xs text-[#ededf3]">time.Until()</code> sem drift de clock e atalhos <kbd class="px-1 text-[11px]">p</kbd>, <kbd class="px-1 text-[11px]">+</kbd>, <kbd class="px-1 text-[11px]">-</kbd>, <kbd class="px-1 text-[11px]">c/Esc</kbd>.
            </p>
          </div>
        </CardSpotlight>
      </div>
    </div>
  </section>
</template>
