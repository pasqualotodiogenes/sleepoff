<script setup lang="ts">
import TextScrollReveal from '~/components/ui/text-scroll-reveal/TextScrollReveal.vue'

const manifestoCopy =
  'Feito para quem domina o terminal e para quem só quer clicar em um aplicativo, agendar o desligamento e dormir tranquilo.'
</script>

<template>
  <section
    aria-label="Manifesto do sleepoff"
    class="relative z-10 border-y border-[#1e1e2a] bg-[#0e1016]"
  >
    <div class="mx-auto max-w-6xl border-x border-[#1e1e2a] px-4 pt-10 sm:px-8 lg:px-12">
      <div class="flex flex-wrap items-center justify-between gap-2 border-b border-[#1e1e2a] pb-4">
        <span class="font-mono-hud text-xs text-[#5266eb]">
          MANIFESTO // ACESSÍVEL &amp; PRECISO
        </span>
        <span class="font-mono-hud text-[11px] text-[#70707d]">
          APP DESKTOP + TUI BUBBLE TEA
        </span>
      </div>

      <TextScrollReveal
        :text="manifestoCopy"
        class="h-[120vh] md:h-[140vh]"
      />
    </div>
  </section>
</template>
