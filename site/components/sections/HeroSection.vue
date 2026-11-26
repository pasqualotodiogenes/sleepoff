<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '#components'
import BlurReveal from '~/components/ui/blur-reveal/BlurReveal.vue'
import InteractiveHoverButton from '~/components/ui/interactive-hover-button/InteractiveHoverButton.vue'
import { useSleepoffSiteState } from '~/composables/useSleepoffSiteState'

const { isPanicActive, panicSecondsRemaining, triggerCommandInEmulator } = useSleepoffSiteState()

const hudStatusText = computed(() => {
  if (isPanicActive.value) {
    return `ALERTA // PÂNICO ATIVO (${panicSecondsRemaining.value}s)`
  }
  return 'RELEASE v1.0.3 // APP DESKTOP & TERMINAL TÁTICO'
})

function scrollToElement(selector: string) {
  if (!import.meta.client) return
  const target = document.querySelector(selector)
  if (!target) return
  const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  target.scrollIntoView({ behavior: reducedMotion ? 'auto' : 'smooth', block: 'start' })
}

function handleTryDryRun() {
  triggerCommandInEmulator('sleepoff 15m --dry-run')
  scrollToElement('#emulador')
}
</script>

<template>
  <section
    id="hero"
    class="relative z-10 border-b border-[#1e1e2a] bg-transparent"
  >
    <div class="mx-auto max-w-6xl border-x border-[#1e1e2a] px-4 py-14 sm:px-8 sm:py-20 lg:px-12 lg:py-24">
      <!-- Eyebrow HUD -->
      <div class="mb-6 flex flex-wrap items-center justify-between gap-3 border border-[#1e1e2a] bg-[#0e1016]/90 px-3.5 py-2">
        <div class="flex items-center gap-2.5">
          <span
            class="inline-block size-2 shrink-0 rounded-none transition-colors"
            :class="isPanicActive ? 'bg-[#ef4444]' : 'bg-[#5266eb]'"
            aria-hidden="true"
          />
          <span
            class="font-mono-hud text-[11px] sm:text-xs"
            :class="isPanicActive ? 'text-[#ef4444]' : 'text-[#c3c3cc]'"
          >
            {{ hudStatusText }}
          </span>
        </div>
        <span class="font-mono-hud text-[11px] text-[#70707d]">
          WIN · LINUX · MACOS // MIT
        </span>
      </div>

      <!-- Display Headline (3 typographic voices) -->
      <div class="max-w-4xl">
        <h1 class="font-display text-ms-5xl sm:text-ms-6xl lg:text-ms-7xl font-bold tracking-tighter text-[#ededf3]">
          Desligue seu PC na hora certa — sem
          <span class="font-serif-accent text-[#ededf3] underline decoration-[#5266eb] decoration-2 underline-offset-4">
            surpresas</span>.
        </h1>

        <BlurReveal
          :delay="0.06"
          :duration="0.45"
          blur="8px"
          :y-offset="10"
          class="mt-6"
        >
          <p class="max-w-3xl text-ms-base sm:text-ms-lg leading-relaxed text-[#c3c3cc]">
            O <strong class="font-semibold text-[#ededf3]">sleepoff</strong> não é apenas para quem usa terminal.
            No Windows você instala como um <strong class="font-semibold text-[#ededf3]">aplicativo comum</strong>
            (<code class="border border-[#272735] bg-[#171721] px-1.5 py-0.5 text-xs text-[#ededf3]">sleepoff-setup.exe</code>),
            abre com dois cliques e acompanha tudo pela bandeja do sistema. Para quem prefere CLI, entrega um binário
            único em Go com zero dependências.
          </p>

          <!-- CTAs -->
          <div class="mt-8 flex flex-col items-stretch gap-3 sm:flex-row sm:items-center">
            <InteractiveHoverButton
              text="Testar Agora (15m Dry-Run)"
              dot-color="bg-[#5266eb]"
              active-text-color="text-[#ededf3]"
              class="min-h-[44px] w-full border-[#5266eb] bg-[#5266eb]/15 px-6 font-mono text-xs uppercase tracking-[0.06em] text-[#ededf3] sm:w-auto"
              aria-label="Testar sleepoff 15m --dry-run no emulador interativo"
              @click="handleTryDryRun"
            >
              <template #icon>
                <Icon name="heroicons:play" class="size-4 shrink-0 text-[#5266eb]" />
              </template>
              <template #active-icon>
                <Icon name="heroicons:play" class="size-4 shrink-0 text-[#ededf3]" />
              </template>
            </InteractiveHoverButton>

            <a
              href="#setup"
              class="hero-secondary-cta inline-flex min-h-[44px] w-full items-center justify-center gap-2 rounded-none border border-[#272735] bg-[#171721] px-6 py-2.5 font-mono text-xs font-semibold uppercase tracking-[0.08em] text-[#ededf3] transition-colors sm:w-auto"
            >
              <Icon name="heroicons:arrow-down-tray" class="size-4 shrink-0 text-[#5266eb]" />
              <span>Instalar App / CLI</span>
            </a>
          </div>
        </BlurReveal>
      </div>

      <!-- Two Audience Tracks in 1px Grid -->
      <div class="mt-12 grid grid-cols-1 divide-y divide-[#1e1e2a] border border-[#1e1e2a] bg-[#0e1016]/95 md:grid-cols-2 md:divide-x md:divide-y-0">
        <!-- Track 1: Installable Desktop App -->
        <article class="flex min-w-0 flex-col justify-between p-5 sm:p-6">
          <div>
            <div class="mb-3 flex items-center justify-between gap-2">
              <span class="font-mono-hud border border-[#5266eb]/50 bg-[#5266eb]/10 px-2 py-0.5 text-[11px] text-[#ededf3]">
                MODO APP // SEM TERMINAL
              </span>
              <Icon name="simple-icons:windows" class="size-4 shrink-0 text-[#70707d]" />
            </div>
            <h2 class="font-display text-ms-xl font-bold tracking-tight text-[#ededf3]">
              Aplicativo Desktop no Windows
            </h2>
            <p class="mt-2 text-ms-sm leading-relaxed text-[#c3c3cc]">
              Instalador gráfico <code class="text-xs text-[#ededf3]">sleepoff-setup.exe</code> com atalhos no Menu Iniciar
              e Área de Trabalho. Abra com dois cliques, escolha <span class="text-[#ededf3]">15 min</span>,
              <span class="text-[#ededf3]">30 min</span> ou <span class="text-[#ededf3]">1 hora</span> no menu guiado,
              acompanhe pelo ícone nativo na bandeja do sistema (<span class="text-[#ededf3]">System Tray</span>) com
              notificações e conte com alarme sonoro de 10s e atualização automática semanal.
            </p>
          </div>
          <div class="mt-5 border-t border-[#1e1e2a] pt-4">
            <a
              href="#setup"
              class="inline-flex min-h-[44px] items-center gap-2 font-mono text-xs font-semibold uppercase tracking-[0.06em] text-[#5266eb] hover:underline"
            >
              <span>Baixar sleepoff-setup.exe (v1.0.3)</span>
              <span aria-hidden="true">→</span>
            </a>
          </div>
        </article>

        <!-- Track 2: Terminal CLI & TUI -->
        <article class="flex min-w-0 flex-col justify-between p-5 sm:p-6">
          <div>
            <div class="mb-3 flex items-center justify-between gap-2">
              <span class="font-mono-hud border border-[#272735] bg-[#171721] px-2 py-0.5 text-[11px] text-[#c3c3cc]">
                MODO CLI // WIN · LINUX · MACOS
              </span>
              <Icon name="heroicons:command-line" class="size-4 shrink-0 text-[#70707d]" />
            </div>
            <h2 class="font-display text-ms-xl font-bold tracking-tight text-[#ededf3]">
              Terminal Tático &amp; Binário Único
            </h2>
            <p class="mt-2 text-ms-sm leading-relaxed text-[#c3c3cc]">
              Compilação estática em Go (<code class="text-xs text-[#ededf3]">CGO_ENABLED=0</code>), zero runtimes ou
              DLLs externas. Dispare <code class="text-xs text-[#ededf3]">sleepoff 30</code> ou
              <code class="text-xs text-[#ededf3]">sleepoff 15m --dry-run</code> direto no terminal, ajuste o tempo em
              execução (<kbd class="px-1.5 py-0.5 text-[11px]">p</kbd>, <kbd class="px-1.5 py-0.5 text-[11px]">+</kbd>,
              <kbd class="px-1.5 py-0.5 text-[11px]">-</kbd>) e proteja com lock nativo de instância única.
            </p>
          </div>
          <div class="mt-5 border-t border-[#1e1e2a] pt-4">
            <button
              type="button"
              class="inline-flex min-h-[44px] cursor-pointer items-center gap-2 font-mono text-xs font-semibold uppercase tracking-[0.06em] text-[#ededf3] hover:text-[#5266eb]"
              @click="handleTryDryRun"
            >
              <span>Rodar sleepoff 15m --dry-run no console</span>
              <span aria-hidden="true">→</span>
            </button>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
@media (hover: hover) {
  .hero-secondary-cta:hover {
    border-color: #5266eb;
    background-color: #1e1e2a;
  }
}
</style>
