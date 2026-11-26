<script setup lang="ts">
import { Icon } from '#components'
import { useSleepoffSiteState } from '~/composables/useSleepoffSiteState'

const { toggleCommandPalette } = useSleepoffSiteState()

function scrollToTop() {
  if (!import.meta.client) return
  const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  window.scrollTo({ top: 0, behavior: reducedMotion ? 'auto' : 'smooth' })
}
</script>

<template>
  <footer class="relative z-10 border-t border-[#1e1e2a] bg-[#171721] text-[#c3c3cc]">
    <div class="mx-auto max-w-6xl border-x border-[#1e1e2a]">
      <!-- Top 1px Grid Matrix -->
      <div class="grid grid-cols-1 divide-y divide-[#1e1e2a] md:grid-cols-3 md:divide-x md:divide-y-0">
        <!-- Brand & License -->
        <div class="p-6 sm:p-8">
          <div class="flex items-center gap-2.5">
            <span class="inline-flex size-6 items-center justify-center border border-[#5266eb] bg-[#5266eb]/15 font-mono text-xs font-bold text-[#ededf3]">
              S_
            </span>
            <span class="font-display text-lg font-bold tracking-tight text-[#ededf3]">
              sleepoff
            </span>
            <span class="font-mono-hud border border-[#272735] bg-[#0e1016] px-2 py-0.5 text-[10px] text-[#5266eb]">
              v1.0.3
            </span>
          </div>
          <p class="mt-3 text-xs leading-relaxed text-[#c3c3cc]">
            Agendador de desligamento em Go: aplicativo instalável com System Tray no Windows e CLI estática para Windows, Linux e macOS.
          </p>
          <p class="mt-4 font-mono text-[11px] uppercase tracking-[0.06em] text-[#70707d]">
            Licença MIT · Código Aberto
          </p>
        </div>

        <!-- Official Repository Links -->
        <div class="flex flex-col justify-between p-6 sm:p-8">
          <span class="font-mono-hud text-[11px] text-[#70707d]">
            REPOSITÓRIO &amp; RELEASES
          </span>
          <div class="mt-3 flex flex-col divide-y divide-[#1e1e2a] border border-[#1e1e2a] bg-[#0e1016]">
            <a
              href="https://github.com/pasqualotodiogenes/sleepoff"
              target="_blank"
              rel="noopener noreferrer"
              class="flex min-h-[44px] items-center justify-between px-4 py-2 font-mono text-xs text-[#ededf3] transition-colors hover:bg-[#1e1e2a]"
            >
              <span>Código-Fonte (GitHub)</span>
              <Icon name="simple-icons:github" class="size-4 shrink-0 text-[#70707d]" />
            </a>
            <a
              href="https://github.com/pasqualotodiogenes/sleepoff/releases/latest"
              target="_blank"
              rel="noopener noreferrer"
              class="flex min-h-[44px] items-center justify-between px-4 py-2 font-mono text-xs text-[#ededf3] transition-colors hover:bg-[#1e1e2a]"
            >
              <span>Releases v1.0.3 (Setup &amp; Binários)</span>
              <Icon name="heroicons:arrow-down-tray" class="size-4 shrink-0 text-[#5266eb]" />
            </a>
            <a
              href="https://github.com/pasqualotodiogenes/sleepoff/issues"
              target="_blank"
              rel="noopener noreferrer"
              class="flex min-h-[44px] items-center justify-between px-4 py-2 font-mono text-xs text-[#c3c3cc] transition-colors hover:bg-[#1e1e2a] hover:text-[#ededf3]"
            >
              <span>Reportar Bug / Issue</span>
              <span aria-hidden="true">↗</span>
            </a>
          </div>
        </div>

        <!-- Tactical Shortcuts & Actions -->
        <div class="flex flex-col justify-between p-6 sm:p-8">
          <div>
            <span class="font-mono-hud text-[11px] text-[#70707d]">
              NAVEGAÇÃO TÁTICA
            </span>
            <div class="mt-3 space-y-2.5 font-mono text-xs">
              <div class="flex items-center justify-between border border-[#1e1e2a] bg-[#0e1016] px-3.5 py-2.5">
                <span class="text-[#c3c3cc]">Alternar / focar console</span>
                <kbd class="px-2 py-0.5 text-xs">~</kbd>
              </div>
              <div class="flex items-center justify-between border border-[#1e1e2a] bg-[#0e1016] px-3.5 py-2.5">
                <span class="text-[#c3c3cc]">Paleta de comandos</span>
                <kbd class="px-2 py-0.5 text-xs">⌘K</kbd>
              </div>
            </div>
          </div>

          <div class="mt-4 flex flex-wrap gap-2">
            <button
              type="button"
              class="inline-flex min-h-[44px] flex-1 cursor-pointer items-center justify-center border border-[#272735] bg-[#272735] px-4 py-2 font-mono text-xs font-semibold text-[#ededf3] transition-colors hover:border-[#5266eb]"
              @click="toggleCommandPalette(true)"
            >
              Abrir ⌘K
            </button>
            <button
              type="button"
              class="inline-flex min-h-[44px] cursor-pointer items-center justify-center border border-[#1e1e2a] bg-[#0e1016] px-4 py-2 font-mono text-xs text-[#c3c3cc] transition-colors hover:border-[#5266eb] hover:text-[#ededf3]"
              @click="scrollToTop"
            >
              ↑ Topo
            </button>
          </div>
        </div>
      </div>

      <!-- Bottom Status Strip -->
      <div class="flex flex-col items-start justify-between gap-2 border-t border-[#1e1e2a] bg-[#0e1016] px-6 py-4 font-mono text-[11px] text-[#70707d] sm:flex-row sm:items-center sm:px-8">
        <span>SLEEPOFF // GO 1.25+ // CGO_ENABLED=0 // MIT LICENSE</span>
        <span>github.com/pasqualotodiogenes/sleepoff</span>
      </div>
    </div>
  </footer>
</template>
