<script setup lang="ts">
import { ClientOnly } from '#components'
import HeroAtmosphere from '~/components/atmosphere/HeroAtmosphere.vue'
import PixelMosaicDivider from '~/components/atmosphere/PixelMosaicDivider.vue'
import CommandPalette from '~/components/layout/CommandPalette.vue'
import ScrollRuler from '~/components/layout/ScrollRuler.vue'
import SiteHeader from '~/components/layout/SiteHeader.vue'
import ArchitectureSection from '~/components/sections/ArchitectureSection.vue'
import CliReferenceSection from '~/components/sections/CliReferenceSection.vue'
import HeroSection from '~/components/sections/HeroSection.vue'
import ManifestoSection from '~/components/sections/ManifestoSection.vue'
import SetupSection from '~/components/sections/SetupSection.vue'
import SiteFooter from '~/components/sections/SiteFooter.vue'
import SleepoffTerminal from '~/components/SleepoffTerminal.vue'
import EncryptedText from '~/components/ui/encrypted-text/EncryptedText.vue'
import { useSleepoffSiteState } from '~/composables/useSleepoffSiteState'

const { isPanicActive, panicSecondsRemaining, isEmulatorFocused, toggleEmulatorFocus } =
  useSleepoffSiteState()
</script>

<template>
  <div class="relative min-h-screen overflow-x-clip bg-[#06070a] text-[#ededf3]">
    <!-- Fixed Panic-Reactive Atmosphere (visible behind Hero and #emulador) -->
    <div
      class="pointer-events-none fixed inset-0 z-0 lg:pr-10"
      aria-hidden="true"
    >
      <HeroAtmosphere />
    </div>

    <!-- Global Tactical Navigation, Right-Edge Ruler & Command Palette -->
    <SiteHeader />
    <ScrollRuler />
    <CommandPalette />

    <!-- Main Editorial Grid (reserves w-10 on desktop for ScrollRuler) -->
    <main class="relative z-10 lg:pr-10">
      <HeroSection />

      <!-- 01 / EMULADOR -->
      <section
        id="emulador"
        data-section="01"
        class="relative z-10 border-b border-[#1e1e2a] bg-transparent"
      >
        <div class="mx-auto max-w-6xl border-x border-[#1e1e2a] px-4 py-14 sm:px-8 sm:py-16 lg:px-12">
          <div class="mb-8 flex flex-col gap-4 border-b border-[#1e1e2a] pb-6 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <div class="mb-3 flex flex-wrap items-center gap-2.5">
                <ClientOnly>
                  <EncryptedText
                    text="01 / EMULADOR"
                    class="font-mono-hud text-xs text-[#5266eb]"
                    encrypted-class="text-[#70707d]"
                    revealed-class="text-[#5266eb]"
                  />
                  <template #fallback>
                    <span class="font-mono-hud text-xs text-[#5266eb]">01 / EMULADOR</span>
                  </template>
                </ClientOnly>
                <span
                  class="font-mono-hud border px-2 py-0.5 text-[10px] transition-colors"
                  :class="
                    isPanicActive
                      ? 'border-[#ef4444]/70 bg-[#ef4444]/10 text-[#ef4444]'
                      : 'border-[#1e1e2a] bg-[#0e1016] text-[#c3c3cc]'
                  "
                >
                  {{ isPanicActive ? `PÂNICO ATIVO // ${panicSecondsRemaining}s` : 'MODO SEGURO // --DRY-RUN' }}
                </span>
              </div>
              <h2 class="font-display text-ms-4xl sm:text-ms-5xl font-bold tracking-tighter text-[#ededf3]">
                Console tático <span class="font-serif-accent">interativo</span>
              </h2>
            </div>

            <div class="flex flex-wrap items-center gap-3">
              <p class="max-w-md text-ms-sm leading-relaxed text-[#c3c3cc]">
                Teste o fluxo real do <code class="text-xs text-[#ededf3]">sleepoff</code> sem desligar o computador. Atalho
                <kbd class="px-1.5 py-0.5 text-[11px]">~</kbd> alterna foco ou recolhe o console.
              </p>
              <button
                type="button"
                class="emulator-focus-btn min-h-[44px] cursor-pointer border px-3.5 py-2 font-mono-hud text-[11px] transition-colors"
                :class="
                  isEmulatorFocused
                    ? 'border-[#5266eb] bg-[#5266eb]/15 text-[#ededf3]'
                    : 'border-[#272735] bg-[#0e1016] text-[#c3c3cc]'
                "
                @click="toggleEmulatorFocus()"
              >
                {{ isEmulatorFocused ? '[~] CONSOLE ATIVO' : '[~] FOCAR CONSOLE' }}
              </button>
            </div>
          </div>

          <ClientOnly>
            <SleepoffTerminal />
            <template #fallback>
              <div class="mx-auto flex w-full max-w-5xl flex-col rounded-none border border-[#1e1e2a] bg-[#0e1016] font-mono text-xs">
                <div class="flex min-h-[44px] items-center justify-between border-b border-[#1e1e2a] bg-[#171721] px-4">
                  <span class="font-mono-hud text-[#ededf3]">01 / CONSOLE // sleepoff v1.0.3 (--dry-run)</span>
                  <span class="font-mono-hud text-[#70707d]">PRONTO</span>
                </div>
                <div class="flex min-h-[360px] flex-col justify-between bg-[#06070a] p-6 text-[#c3c3cc]">
                  <div class="space-y-2">
                    <p><span class="text-[#5266eb]">$</span> sleepoff 15m --dry-run</p>
                    <p class="text-[#70707d]">Inicializando máquina de estados Bubble Tea...</p>
                  </div>
                  <p class="text-[#70707d]">Aguardando hidratação do emulador tático.</p>
                </div>
                <div class="grid grid-cols-2 gap-px border-t border-[#1e1e2a] bg-[#1e1e2a] lg:grid-cols-4">
                  <div class="bg-[#0e1016] px-4 py-3 text-[#70707d]">1. Iniciar Timer (15m)</div>
                  <div class="bg-[#0e1016] px-4 py-3 text-[#70707d]">2. Pausar / Retomar</div>
                  <div class="bg-[#0e1016] px-4 py-3 text-[#70707d]">3. Ajustar (+5min)</div>
                  <div class="bg-[#0e1016] px-4 py-3 text-[#70707d]">4. Abortar Alarme</div>
                </div>
              </div>
            </template>
          </ClientOnly>
        </div>
      </section>

      <PixelMosaicDivider variant="dark-to-onyx" />
      <ManifestoSection />
      <SetupSection />
      <ArchitectureSection />
      <PixelMosaicDivider variant="onyx-to-void" />
      <CliReferenceSection />
      <SiteFooter />
    </main>
  </div>
</template>

<style scoped>
@media (hover: hover) {
  .emulator-focus-btn:hover {
    border-color: #5266eb;
    color: #ededf3;
  }
}
</style>
