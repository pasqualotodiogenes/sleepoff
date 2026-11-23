<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent } from 'vue'

// Importando componentes Inspira UI
import ShimmerButton from '~/components/ui/shimmer-button/ShimmerButton.vue'
import InteractiveHoverButton from '~/components/ui/interactive-hover-button/InteractiveHoverButton.vue'
import CardSpotlight from '~/components/ui/card-spotlight/CardSpotlight.vue'
import FlickeringGrid from '~/components/ui/flickering-grid/FlickeringGrid.vue'
import ParticlesBg from '~/components/ui/particles-bg/ParticlesBg.vue'
import Meteors from '~/components/ui/meteors/Meteors.vue'
import FlipWords from '~/components/ui/flip-words/FlipWords.vue'
import BentoGrid from '~/components/ui/bento-grid/BentoGrid.vue'
import BentoGridItem from '~/components/ui/bento-grid/BentoGridItem.vue'
import GlareCard from '~/components/ui/glare-card/GlareCard.vue'

// Carregamento sob demanda do emulador (Item 22 do txt.txt - Dividir JavaScript)
const SleepoffTerminal = defineAsyncComponent(() => import('~/components/SleepoffTerminal.vue'))

// Estado para abas de instalação com detecção dinâmica do SO do visitante
const activeOs = ref<'windows' | 'macos' | 'linux'>('linux')
const copiedCommand = ref<string | null>(null)

function detectOs() {
  if (!import.meta.client) return
  const ua = (navigator.userAgent || navigator.platform || '').toLowerCase()
  if (ua.includes('linux') || ua.includes('x11')) {
    activeOs.value = 'linux'
  } else if (ua.includes('mac') || ua.includes('darwin') || ua.includes('os x')) {
    activeOs.value = 'macos'
  } else if (ua.includes('win')) {
    activeOs.value = 'windows'
  }
}

onMounted(() => {
  detectOs()
})

function copyToClipboard(text: string, id: string) {
  if (import.meta.client && navigator.clipboard) {
    navigator.clipboard.writeText(text)
    copiedCommand.value = id
    setTimeout(() => {
      if (copiedCommand.value === id) {
        copiedCommand.value = null
      }
    }, 2000)
  }
}
</script>

<template>
  <main class="min-h-screen bg-slate-950 text-slate-100 selection:bg-violet-500 selection:text-white relative overflow-hidden">
    
    <!-- Inspira UI Background Animado: Partículas e Grid contidos no Hero (evita canvas de 4000px em toda a página) -->
    <ClientOnly>
      <div class="absolute top-0 left-0 right-0 h-[900px] pointer-events-none z-0 overflow-hidden">
        <ParticlesBg :quantity="12" :ease="90" :color="'#7C3AED'" class="absolute inset-0 opacity-30" />
        <FlickeringGrid
          :square-size="4"
          :grid-gap="32"
          :color="'#06B6D4'"
          :max-opacity="0.08"
          :flicker-chance="0.03"
          class="absolute inset-0 mask-radial-fade opacity-25"
        />
        <Meteors :number="6" class="opacity-30" />
      </div>
      <template #fallback>
        <div class="absolute top-0 left-0 right-0 h-[900px] bg-slate-950 pointer-events-none z-0" />
      </template>
    </ClientOnly>

    <!-- Top Navigation -->
    <header class="relative z-10 border-b border-slate-800/80 bg-slate-950/60 backdrop-blur-md">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-violet-600 to-cyan-400 flex items-center justify-center font-black text-white text-base shadow-lg shadow-violet-600/30">
            S
          </div>
          <span class="font-bold tracking-tight text-lg text-slate-100">sleepoff</span>
          <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-violet-950/80 text-violet-300 border border-violet-800/60">
            v1.0.0
          </span>
        </div>

        <nav class="flex items-center space-x-3 sm:space-x-6 text-sm font-medium text-slate-400">
          <div class="hidden md:flex items-center space-x-6">
            <a href="#demo" class="hover:text-cyan-400 transition-colors">Terminal</a>
            <a href="#install" class="hover:text-cyan-400 transition-colors">Instalação</a>
            <a href="#features" class="hover:text-cyan-400 transition-colors">Arquitetura</a>
            <a href="#docs" class="hover:text-cyan-400 transition-colors">Comandos</a>
          </div>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Repositório oficial do sleepoff no GitHub"
            class="flex items-center space-x-1.5 px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-slate-300 hover:text-white hover:border-slate-700 transition-all text-xs font-mono min-h-[36px]"
          >
            <Icon name="simple-icons:github" class="w-4 h-4 text-slate-300" />
            <span class="hidden sm:inline">GitHub</span>
          </a>
        </nav>
      </div>
    </header>

    <!-- 1. Hero Section -->
    <section class="relative z-10 pt-16 sm:pt-24 pb-12 px-4 sm:px-6 max-w-5xl mx-auto text-center">
      <div class="inline-flex items-center space-x-2 px-3 py-1.5 rounded-full bg-slate-900/80 border border-slate-800 text-xs font-mono text-cyan-400 mb-6 shadow-sm">
        <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
        <span>Go 1.23 • CGO_ENABLED=0 • 100% Nativo</span>
      </div>

      <h1 class="text-4xl sm:text-6xl font-extrabold tracking-tight text-white mb-6">
        Timer de desligamento <br class="hidden sm:inline" />
        <span class="inline-block">
          <ClientOnly>
            <FlipWords
              :words="['inteligente', 'seguro', 'determinístico', 'ultra-leve', 'sem bloat']"
              class="bg-gradient-to-r from-violet-400 via-indigo-300 to-cyan-400 bg-clip-text text-transparent font-extrabold"
            />
            <template #fallback>
              <span class="bg-gradient-to-r from-violet-400 via-indigo-300 to-cyan-400 bg-clip-text text-transparent font-extrabold">
                inteligente
              </span>
            </template>
          </ClientOnly>
        </span>
      </h1>

      <p class="max-w-2xl mx-auto text-base sm:text-lg text-slate-400 mb-10 leading-relaxed">
        Controle o desligamento do computador pelo terminal ou com um clique. TUI com Bubble Tea, tela de pânico com contagem regressiva sonora e single-instance lock nativo.
      </p>

      <div class="flex flex-wrap items-center justify-center gap-4">
        <a href="#install" aria-label="Baixar Binário / Setup — Ir para seção de instalação">
          <ShimmerButton
            shimmer-color="#c084fc"
            shimmer-size="0.08em"
            shimmer-duration="2.2s"
            background="linear-gradient(135deg, #7c3aed 0%, #4f46e5 100%)"
            border-radius="12px"
            class="font-bold text-sm text-white shadow-xl shadow-violet-600/30 flex items-center space-x-2 cursor-pointer min-h-[44px] px-6"
          >
            <Icon name="heroicons:arrow-down-tray" class="w-4 h-4 mr-1 text-white" />
            <span>Baixar Binário / Setup</span>
          </ShimmerButton>
        </a>

        <a href="#demo" aria-label="Experimentar no Terminal — Ir para emulador interativo">
          <InteractiveHoverButton
            text="Experimentar no Terminal"
            dot-color="bg-cyan-400"
            active-text-color="text-slate-950"
            class="min-h-[44px] px-6 text-sm font-semibold border-cyan-500/30 bg-slate-900/90 text-cyan-300 hover:border-cyan-400"
          >
            <template #icon>
              <Icon name="heroicons:command-line" class="w-4 h-4 mr-1 text-cyan-400" />
            </template>
            <template #active-icon>
              <Icon name="heroicons:command-line" class="w-4 h-4 mr-1 text-slate-950" />
            </template>
          </InteractiveHoverButton>
        </a>
      </div>
    </section>

    <!-- 2. Terminal Emulador Real + Tutorial Interativo (ClientOnly com Skeleton Estruturado - Item 01 txt.txt) -->
    <section id="demo" class="relative z-10 px-4 sm:px-6 py-8">
      <div class="max-w-5xl mx-auto">
        <div class="text-center mb-6">
          <h2 class="text-2xl font-bold tracking-tight text-slate-100">
            Emulador Interativo Web
          </h2>
          <p class="text-xs sm:text-sm text-slate-300">
            Teste exatamente a mesma máquina de estados do Bubble Tea sem instalar nada agora.
          </p>
        </div>

        <div>
          <ClientOnly>
            <SleepoffTerminal />
            <template #fallback>
              <div class="w-full max-w-5xl mx-auto min-h-[580px] sm:min-h-[600px] rounded-2xl border border-slate-700/70 bg-slate-950/95 overflow-hidden shadow-2xl flex flex-col font-mono text-sm">
                <!-- Top bar skeleton -->
                <div class="h-11 px-4 border-b border-slate-800/90 bg-slate-900/90 flex items-center justify-between">
                  <div class="flex items-center space-x-2">
                    <span class="w-3 h-3 rounded-full bg-red-500/60" />
                    <span class="w-3 h-3 rounded-full bg-amber-500/60" />
                    <span class="w-3 h-3 rounded-full bg-emerald-500/60" />
                    <div class="h-6 w-24 rounded bg-slate-800 animate-pulse ml-3" />
                  </div>
                  <div class="h-5 w-20 rounded bg-slate-800/60 animate-pulse" />
                </div>
                <!-- Body skeleton -->
                <div class="p-6 flex-1 flex flex-col justify-between">
                  <div class="space-y-3">
                    <div class="flex items-center space-x-2">
                      <span class="text-emerald-400 font-bold text-xs">goddi@sleepoff</span>
                      <span class="text-slate-600">:</span>
                      <span class="text-cyan-400 font-bold text-xs">~</span>
                      <span class="text-violet-400 font-bold text-xs mr-1">$</span>
                      <div class="h-4 w-40 rounded bg-slate-800 animate-pulse" />
                    </div>
                    <div class="h-4 w-3/4 rounded bg-slate-900 animate-pulse" />
                    <div class="h-4 w-1/2 rounded bg-slate-900 animate-pulse" />
                  </div>
                  <div class="flex items-center justify-center text-slate-400 text-xs py-12">
                    <Icon name="heroicons:command-line" class="w-5 h-5 mr-2 animate-pulse text-cyan-400" />
                    <span>Iniciando ambiente Web TUI interativo...</span>
                  </div>
                </div>
                <!-- Mock missions bar -->
                <div class="px-4 py-3 bg-slate-950/90 border-t border-slate-800/60">
                  <div class="h-4 w-48 rounded bg-slate-900 animate-pulse mb-2" />
                  <div class="grid grid-cols-2 lg:grid-cols-4 gap-2">
                    <div class="h-12 rounded-xl bg-slate-900/40 border border-slate-800/80 animate-pulse" />
                    <div class="h-12 rounded-xl bg-slate-900/40 border border-slate-800/80 animate-pulse" />
                    <div class="h-12 rounded-xl bg-slate-900/40 border border-slate-800/80 animate-pulse" />
                    <div class="h-12 rounded-xl bg-slate-900/40 border border-slate-800/80 animate-pulse" />
                  </div>
                </div>
              </div>
            </template>
          </ClientOnly>
        </div>
      </div>
    </section>

    <!-- 3. Instalação Multiplataforma com Abas e Cópia em 1 Clique -->
    <section id="install" class="relative z-10 px-4 sm:px-6 py-16 max-w-5xl mx-auto">
      <div class="text-center mb-10">
        <h2 class="text-3xl font-extrabold text-white tracking-tight mb-3">
          Instalação em Poucos Segundos
        </h2>
        <p class="text-slate-400 text-sm max-w-lg mx-auto">
          Binários pré-compilados e auto-contidos. Sem dependências externas, sem runtimes pesados.
        </p>
      </div>

      <!-- OS Switcher Tabs -->
      <div class="flex justify-center mb-8 px-2">
        <div class="inline-flex p-1 sm:p-1.5 rounded-2xl bg-slate-900/90 border border-slate-800 shadow-inner max-w-full">
          <button
            @click="activeOs = 'windows'"
            aria-label="Selecionar instalador para Windows"
            :class="activeOs === 'windows' ? 'bg-gradient-to-r from-violet-600 to-indigo-600 text-white shadow-lg shadow-violet-600/30 scale-100' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'"
            class="px-3 sm:px-6 py-2 sm:py-2.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center space-x-1.5 sm:space-x-2 cursor-pointer min-h-[40px]"
          >
            <Icon name="simple-icons:windows" class="w-3.5 sm:w-4 h-3.5 sm:h-4" />
            <span>Windows</span>
          </button>
          <button
            @click="activeOs = 'macos'"
            aria-label="Selecionar instalador para macOS"
            :class="activeOs === 'macos' ? 'bg-gradient-to-r from-violet-600 to-indigo-600 text-white shadow-lg shadow-violet-600/30 scale-100' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'"
            class="px-3 sm:px-6 py-2 sm:py-2.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center space-x-1.5 sm:space-x-2 cursor-pointer min-h-[40px]"
          >
            <Icon name="simple-icons:apple" class="w-3.5 sm:w-4 h-3.5 sm:h-4" />
            <span>macOS</span>
          </button>
          <button
            @click="activeOs = 'linux'"
            aria-label="Selecionar instalador para Linux"
            :class="activeOs === 'linux' ? 'bg-gradient-to-r from-violet-600 to-indigo-600 text-white shadow-lg shadow-violet-600/30 scale-100' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'"
            class="px-3 sm:px-6 py-2 sm:py-2.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center space-x-1.5 sm:space-x-2 cursor-pointer min-h-[40px]"
          >
            <Icon name="simple-icons:linux" class="w-3.5 sm:w-4 h-3.5 sm:h-4" />
            <span>Linux</span>
          </button>
        </div>
      </div>

      <!-- Tab Windows -->
      <div v-if="activeOs === 'windows'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <CardSpotlight class="w-full min-w-0 p-6 rounded-2xl bg-slate-900/60 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="flex items-center space-x-3 mb-4">
              <div class="w-10 h-10 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center text-violet-400 shrink-0">
                <Icon name="heroicons:arrow-down-tray" class="w-5 h-5" />
              </div>
              <div class="min-w-0">
                <h3 class="font-bold text-white text-base truncate">Instalador Windows (GUI)</h3>
                <p class="text-xs text-slate-400 truncate">Instalação em 1 clique com Inno Setup</p>
              </div>
            </div>
            <p class="text-xs text-slate-300 leading-relaxed mb-6">
              Adiciona o executável ao PATH do sistema operacional e cria atalho no menu iniciar para abrir com duplo clique.
            </p>
          </div>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff/releases/latest"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Download sleepoff-setup.exe (GUI) — Instalador Windows"
            class="w-full text-center py-3 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-500 hover:to-indigo-500 text-white text-xs font-bold transition-all shadow-lg shadow-violet-600/30 flex items-center justify-center space-x-2 cursor-pointer hover:scale-[1.02] active:scale-[0.98] min-h-[44px]"
          >
            <Icon name="heroicons:arrow-down-tray" class="w-4 h-4 shrink-0" />
            <span class="truncate">Download sleepoff-setup.exe (GUI)</span>
          </a>
        </CardSpotlight>

        <CardSpotlight class="w-full min-w-0 p-6 rounded-2xl bg-slate-900/60 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="flex items-center space-x-3 mb-4">
              <div class="w-10 h-10 rounded-xl bg-cyan-600/20 border border-cyan-500/30 flex items-center justify-center text-cyan-400 shrink-0">
                <Icon name="heroicons:command-line" class="w-5 h-5" />
              </div>
              <div class="min-w-0">
                <h3 class="font-bold text-white text-base truncate">PowerShell / Winget / Scoop</h3>
                <p class="text-xs text-slate-400 truncate">Instalação rápida via terminal</p>
              </div>
            </div>
            <div class="relative bg-slate-950 p-3 rounded-xl border border-slate-800 font-mono text-xs text-cyan-300 mb-4 flex items-center justify-between min-w-0">
              <span class="truncate min-w-0 flex-1">irm https://git.new/sleepoff.ps1 | iex</span>
              <button
                @click="copyToClipboard('irm https://git.new/sleepoff.ps1 | iex', 'ps')"
                class="ml-2 text-slate-400 hover:text-white shrink-0 cursor-pointer p-1.5 rounded-lg hover:bg-slate-800 transition-colors min-w-[36px] min-h-[36px] flex items-center justify-center"
                aria-label="Copiar comando de instalação do PowerShell"
                title="Copiar comando"
              >
                <Icon :name="copiedCommand === 'ps' ? 'heroicons:check' : 'heroicons:clipboard-document'" class="w-4 h-4 text-cyan-400" />
              </button>
            </div>
          </div>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff/releases/latest"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Baixar sleepoff-windows-amd64.zip para Windows"
            class="w-full text-center py-3 rounded-xl bg-slate-800/90 hover:bg-slate-700/90 text-slate-200 hover:text-white text-xs font-bold transition-all border border-slate-700/60 flex items-center justify-center space-x-2 cursor-pointer hover:scale-[1.02] active:scale-[0.98] min-h-[44px]"
          >
            <Icon name="heroicons:archive-box" class="w-4 h-4 text-cyan-400 shrink-0" />
            <span class="truncate">Baixar sleepoff-windows-amd64.zip</span>
          </a>
        </CardSpotlight>
      </div>

      <!-- Tab macOS -->
      <div v-else-if="activeOs === 'macos'" class="max-w-2xl mx-auto w-full min-w-0">
        <CardSpotlight class="w-full min-w-0 p-6 rounded-2xl bg-slate-900/60 border border-slate-800">
          <div class="flex items-center space-x-3 mb-4">
            <div class="w-10 h-10 rounded-xl bg-cyan-600/20 border border-cyan-500/30 flex items-center justify-center text-cyan-400 shrink-0">
              <Icon name="simple-icons:apple" class="w-5 h-5" />
            </div>
            <div class="min-w-0">
              <h3 class="font-bold text-white text-base truncate">macOS (Apple Silicon & Intel)</h3>
              <p class="text-xs text-slate-400 truncate">Instale via Homebrew ou binário universal</p>
            </div>
          </div>
          <div class="space-y-3 font-mono text-xs mb-4 min-w-0">
            <div class="bg-slate-950 p-3 rounded-xl border border-slate-800 flex items-center justify-between text-cyan-300 min-w-0">
              <span class="truncate min-w-0 flex-1">brew install pasqualotodiogenes/tap/sleepoff</span>
              <button
                @click="copyToClipboard('brew install pasqualotodiogenes/tap/sleepoff', 'brew')"
                class="ml-2 text-slate-400 hover:text-white cursor-pointer shrink-0 p-1.5 rounded-lg hover:bg-slate-800 transition-colors min-w-[36px] min-h-[36px] flex items-center justify-center"
                aria-label="Copiar comando de instalação via Homebrew"
                title="Copiar comando"
              >
                <Icon :name="copiedCommand === 'brew' ? 'heroicons:check' : 'heroicons:clipboard-document'" class="w-4 h-4 text-cyan-400" />
              </button>
            </div>
            <div class="bg-slate-950 p-3 rounded-xl border border-slate-800 flex items-center justify-between text-slate-300 min-w-0">
              <span class="truncate min-w-0 flex-1">curl -fsSL https://git.new/sleepoff-darwin.sh | sh</span>
              <button
                @click="copyToClipboard('curl -fsSL https://git.new/sleepoff-darwin.sh | sh', 'darwin')"
                class="ml-2 text-slate-400 hover:text-white cursor-pointer shrink-0 p-1.5 rounded-lg hover:bg-slate-800 transition-colors min-w-[36px] min-h-[36px] flex items-center justify-center"
                aria-label="Copiar comando de instalação via curl no macOS"
                title="Copiar comando"
              >
                <Icon :name="copiedCommand === 'darwin' ? 'heroicons:check' : 'heroicons:clipboard-document'" class="w-4 h-4 text-cyan-400" />
              </button>
            </div>
          </div>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff/releases/latest"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Baixar sleepoff-darwin-universal.tar.gz para macOS"
            class="w-full text-center py-3 rounded-xl bg-slate-800/90 hover:bg-slate-700/90 text-slate-200 hover:text-white text-xs font-bold transition-all border border-slate-700/60 flex items-center justify-center space-x-2 cursor-pointer hover:scale-[1.02] active:scale-[0.98] min-h-[44px]"
          >
            <Icon name="heroicons:archive-box" class="w-4 h-4 text-cyan-400 shrink-0" />
            <span class="truncate">Baixar sleepoff-darwin-universal.tar.gz</span>
          </a>
        </CardSpotlight>
      </div>

      <!-- Tab Linux -->
      <div v-else class="max-w-2xl mx-auto w-full min-w-0">
        <CardSpotlight class="w-full min-w-0 p-6 rounded-2xl bg-slate-900/60 border border-slate-800">
          <div class="flex items-center space-x-3 mb-4">
            <div class="w-10 h-10 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center text-violet-400 shrink-0">
              <Icon name="simple-icons:linux" class="w-5 h-5" />
            </div>
            <div class="min-w-0">
              <h3 class="font-bold text-white text-base truncate">Linux (x86_64 / arm64)</h3>
              <p class="text-xs text-slate-400 truncate">Compatível com systemd, elogind e init.d</p>
            </div>
          </div>
          <div class="space-y-3 font-mono text-xs mb-4 min-w-0">
            <div class="bg-slate-950 p-3 rounded-xl border border-slate-800 flex items-center justify-between text-cyan-300 min-w-0">
              <span class="truncate min-w-0 flex-1">curl -fsSL https://git.new/sleepoff.sh | sudo sh</span>
              <button
                @click="copyToClipboard('curl -fsSL https://git.new/sleepoff.sh | sudo sh', 'linux-sh')"
                class="ml-2 text-slate-400 hover:text-white cursor-pointer shrink-0 p-1.5 rounded-lg hover:bg-slate-800 transition-colors min-w-[36px] min-h-[36px] flex items-center justify-center"
                aria-label="Copiar script de instalação curl no Linux"
                title="Copiar comando"
              >
                <Icon :name="copiedCommand === 'linux-sh' ? 'heroicons:check' : 'heroicons:clipboard-document'" class="w-4 h-4 text-cyan-400" />
              </button>
            </div>
            <div class="bg-slate-950 p-3 rounded-xl border border-slate-800 flex items-center justify-between text-slate-300 min-w-0">
              <span class="truncate min-w-0 flex-1">go install github.com/pasqualotodiogenes/sleepoff@latest</span>
              <button
                @click="copyToClipboard('go install github.com/pasqualotodiogenes/sleepoff@latest', 'go-inst')"
                class="ml-2 text-slate-400 hover:text-white cursor-pointer shrink-0 p-1.5 rounded-lg hover:bg-slate-800 transition-colors min-w-[36px] min-h-[36px] flex items-center justify-center"
                aria-label="Copiar comando go install"
                title="Copiar comando"
              >
                <Icon :name="copiedCommand === 'go-inst' ? 'heroicons:check' : 'heroicons:clipboard-document'" class="w-4 h-4 text-cyan-400" />
              </button>
            </div>
          </div>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff/releases/latest"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Baixar sleepoff-linux-amd64.tar.gz para Linux"
            class="w-full text-center py-3 rounded-xl bg-slate-800/90 hover:bg-slate-700/90 text-slate-200 hover:text-white text-xs font-bold transition-all border border-slate-700/60 flex items-center justify-center space-x-2 cursor-pointer hover:scale-[1.02] active:scale-[0.98] min-h-[44px]"
          >
            <Icon name="heroicons:archive-box" class="w-4 h-4 text-violet-400 shrink-0" />
            <span class="truncate">Baixar sleepoff-linux-amd64.tar.gz</span>
          </a>
        </CardSpotlight>
      </div>
    </section>

    <!-- 4. Arquitetura & Diferenciais de Engenharia -->
    <section id="features" class="relative z-10 px-4 sm:px-6 py-16 max-w-6xl mx-auto border-t border-slate-800/80">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-extrabold text-white tracking-tight mb-3">
          Engenharia Sólida & Zero Bloat
        </h2>
        <p class="text-slate-400 text-sm max-w-xl mx-auto">
          Construído seguindo os mais rigorosos fundamentos de sistemas operacionais e concorrência limpa em Go.
        </p>
      </div>

      <!-- Inspira UI BentoGrid Layout -->
      <BentoGrid class="max-w-5xl">
        <!-- Card 1: Destaque CGO_ENABLED=0 com BorderBeam -->
        <BentoGridItem class="md:col-span-2 relative overflow-hidden bg-slate-900/70 border-slate-800/80 hover:border-violet-500/40">
          <template #header>
            <div class="relative h-44 w-full rounded-xl overflow-hidden bg-gradient-to-br from-violet-950/60 via-slate-900 to-slate-950 border border-slate-800/60 flex items-center justify-center p-4">
              <div class="font-mono text-center">
                <div class="inline-block px-3 py-1 rounded bg-violet-900/40 border border-violet-700/50 text-cyan-300 text-xs font-bold mb-2">
                  CGO_ENABLED=0 go build -ldflags="-s -w"
                </div>
                <p class="text-[11px] text-slate-400">Zero DLLs externas • Auto-suficiente • Executável instantâneo</p>
              </div>
            </div>
          </template>
          <template #icon>
            <div class="w-10 h-10 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center text-violet-400 mb-2">
              <Icon name="heroicons:cpu-chip" class="w-5 h-5" />
            </div>
          </template>
          <template #title>
            <span class="text-lg font-bold text-white">Compilação Estática Pura</span>
          </template>
          <template #description>
            <span class="text-slate-400 leading-relaxed text-xs">
              Binário estático Go 1.23 puro sem qualquer amarração de dependências dinâmicas no SO. Portável, determinístico e com inicialização em menos de 10 milissegundos.
            </span>
          </template>
        </BentoGridItem>

        <!-- Card 2: Single-Instance Lock -->
        <BentoGridItem class="md:col-span-1 bg-slate-900/70 border-slate-800/80 hover:border-cyan-500/40">
          <template #header>
            <div class="h-44 w-full rounded-xl bg-gradient-to-br from-cyan-950/40 via-slate-900 to-slate-950 border border-slate-800/60 flex flex-col justify-center p-3.5 text-left font-mono text-[11px] space-y-1.5 overflow-hidden">
              <div class="flex items-center space-x-1.5 text-cyan-400 font-bold mb-1">
                <Icon name="heroicons:lock-closed" class="w-4 h-4" />
                <span>Mutex / flock Atômico</span>
              </div>
              <div class="text-slate-400 truncate">WIN: <span class="text-slate-200">CreateMutexW("sleepoff_lock")</span></div>
              <div class="text-slate-400 truncate">UNIX: <span class="text-slate-200">flock(fd, LOCK_EX|LOCK_NB)</span></div>
              <div class="text-emerald-400 pt-1.5 flex items-center space-x-1.5 border-t border-slate-800/80">
                <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                <span class="text-[10px]">Prevenção de processos duplicados</span>
              </div>
            </div>
          </template>
          <template #icon>
            <div class="w-10 h-10 rounded-xl bg-cyan-600/20 border border-cyan-500/30 flex items-center justify-center text-cyan-400 mb-2">
              <Icon name="heroicons:lock-closed" class="w-5 h-5" />
            </div>
          </template>
          <template #title>
            <span class="text-lg font-bold text-white">Single-Instance Lock</span>
          </template>
          <template #description>
            <span class="text-slate-400 leading-relaxed text-xs">
              Proteção atômica via Named Mutex (Windows) e file locks (Unix). Previne corridas de processos e execuções duplicadas acidentais.
            </span>
          </template>
        </BentoGridItem>

        <!-- Card 3: Modo Simulação Seguro -->
        <BentoGridItem class="md:col-span-1 bg-slate-900/70 border-slate-800/80 hover:border-emerald-500/40">
          <template #header>
            <div class="h-44 w-full rounded-xl bg-gradient-to-br from-emerald-950/40 via-slate-900 to-slate-950 border border-slate-800/60 flex flex-col justify-center p-3.5 text-left font-mono text-[11px] space-y-1.5 overflow-hidden">
              <div class="flex items-center space-x-1.5 text-emerald-400 font-bold mb-1">
                <Icon name="heroicons:shield-check" class="w-4 h-4" />
                <span>Interceptador de Hardware</span>
              </div>
              <div class="text-slate-400">FLAG: <span class="text-cyan-300 font-bold">--dry-run</span></div>
              <div class="text-slate-400">SYSCALL: <span class="text-rose-400 line-through">InitiateShutdown</span></div>
              <div class="text-emerald-400 pt-1.5 flex items-center space-x-1.5 border-t border-slate-800/80">
                <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                <span class="text-[10px]">100% Simulado (Zero risco)</span>
              </div>
            </div>
          </template>
          <template #icon>
            <div class="w-10 h-10 rounded-xl bg-emerald-600/20 border border-emerald-500/30 flex items-center justify-center text-emerald-400 mb-2">
              <Icon name="heroicons:shield-check" class="w-5 h-5" />
            </div>
          </template>
          <template #title>
            <span class="text-lg font-bold text-white">Modo Seguro (--dry-run)</span>
          </template>
          <template #description>
            <span class="text-slate-400 leading-relaxed text-xs">
              Simule todo o pipeline do timer, a tela de pânico e os bipes de aviso com total segurança, sem arriscar desligar a máquina.
            </span>
          </template>
        </BentoGridItem>

        <!-- Card 4: Terminal UI Nativa (Bubble Tea) -->
        <BentoGridItem class="md:col-span-2 bg-slate-900/70 border-slate-800/80 hover:border-violet-500/40">
          <template #header>
            <div class="h-44 w-full rounded-xl bg-gradient-to-br from-violet-950/40 via-slate-900 to-slate-950 border border-slate-800/60 flex flex-col justify-center p-4">
              <div class="flex items-center justify-between text-xs font-mono text-cyan-300 mb-3">
                <span class="px-2.5 py-1 rounded bg-slate-950 border border-violet-800/50 text-violet-300 font-bold">Init()</span>
                <span class="text-slate-400 font-bold">➔</span>
                <span class="px-2.5 py-1 rounded bg-slate-950 border border-cyan-800/50 text-cyan-300 font-bold">Update(Msg)</span>
                <span class="text-slate-400 font-bold">➔</span>
                <span class="px-2.5 py-1 rounded bg-slate-950 border border-indigo-800/50 text-indigo-300 font-bold">View()</span>
              </div>
              <div class="flex items-center justify-center space-x-2 text-[11px] font-mono text-slate-400">
                <span class="px-2 py-0.5 rounded bg-slate-900 border border-slate-800">Alt Screen Buffer</span>
                <span>•</span>
                <span class="px-2 py-0.5 rounded bg-slate-900 border border-slate-800">ANSI 256 Cores</span>
                <span>•</span>
                <span class="px-2 py-0.5 rounded bg-slate-900 border border-slate-800">Restauração Limpa</span>
              </div>
            </div>
          </template>
          <template #icon>
            <div class="w-10 h-10 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center text-violet-400 mb-2">
              <Icon name="heroicons:command-line" class="w-5 h-5" />
            </div>
          </template>
          <template #title>
            <span class="text-lg font-bold text-white">TUI Reativa e Fluida</span>
          </template>
          <template #description>
            <span class="text-slate-400 leading-relaxed text-xs">
              Interface rica para terminal baseada no padrão The Elm Architecture (TEA). Renderização em frame rate estável com alt screen buffer e restauração limpa do cursor.
            </span>
          </template>
        </BentoGridItem>
      </BentoGrid>
    </section>

    <!-- 5. Guia de Atalhos e Flags Fiel ao Código -->
    <section id="docs" class="relative z-10 px-4 sm:px-6 py-16 max-w-5xl mx-auto border-t border-slate-800/80">
      <div class="text-center mb-10">
        <h2 class="text-3xl font-extrabold text-white tracking-tight mb-3">
          Atalhos & Linha de Comando
        </h2>
        <p class="text-slate-400 text-sm max-w-md mx-auto">
          Mapeamento exato dos comandos e controles implementados no motor Go.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Tabela Comandos CLI -->
        <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-base font-bold text-white mb-4 flex items-center space-x-2">
            <Icon name="heroicons:command-line" class="w-5 h-5 text-violet-400" />
            <span>Comandos de Inicialização</span>
          </h3>
          <div class="space-y-3 font-mono text-xs">
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <span class="text-cyan-400 font-bold">sleepoff</span>
              <span class="text-slate-400">Abre menu TUI interativo</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <span class="text-cyan-400 font-bold">sleepoff 45m</span>
              <span class="text-slate-400">Inicia timer de 45 minutos</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <span class="text-cyan-400 font-bold">sleepoff 1h30m</span>
              <span class="text-slate-400">Aceita formato composto</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <span class="text-cyan-400 font-bold">sleepoff --dry-run</span>
              <span class="text-slate-400">Modo de simulação sem shutdown</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <span class="text-cyan-400 font-bold">sleepoff --help</span>
              <span class="text-slate-400">Exibe manual de uso</span>
            </div>
          </div>
        </div>

        <!-- Tabela Atalhos TUI -->
        <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-base font-bold text-white mb-4 flex items-center space-x-2">
            <Icon name="heroicons:cursor-arrow-rays" class="w-5 h-5 text-cyan-400" />
            <span>Controles do Timer Ativo</span>
          </h3>
          <div class="space-y-3 font-mono text-xs">
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-violet-300 font-bold">p</kbd>
                <span class="text-slate-400">ou</span>
                <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-violet-300 font-bold">Espaço</kbd>
              </div>
              <span class="text-slate-400">Pausar / Retomar cronômetro</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <kbd class="px-2 py-0.5 rounded bg-slate-800 text-cyan-300 font-bold">+</kbd>
              <span class="text-slate-400">Adicionar 5 minutos</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <kbd class="px-2 py-0.5 rounded bg-slate-800 text-cyan-300 font-bold">-</kbd>
              <span class="text-slate-400">Subtrair 5 minutos</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <kbd class="px-1.5 py-0.5 rounded bg-rose-950 border border-rose-800 text-rose-300 font-bold">2x c</kbd>
                <span class="text-slate-400">ou</span>
                <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-bold">Esc</kbd>
              </div>
              <span class="text-slate-400">Cancelar e desarmar</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
              <span class="text-amber-400 font-bold">Últimos 10s</span>
              <span class="text-slate-400">Alarme sonoro progressivo (pânico)</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-slate-800/80 py-10 px-4 sm:px-6 bg-slate-950/80 text-xs text-slate-400 text-center">
      <div class="max-w-5xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          sleepoff — Desenvolvido com Go e engenharia de precisão sob licença MIT.
        </div>
        <div class="flex items-center space-x-6">
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Código-Fonte do sleepoff no GitHub"
            class="hover:text-slate-300 transition-colors py-2 px-1 min-h-[36px] flex items-center"
          >
            Código-Fonte
          </a>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff/releases"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Releases e binários do sleepoff no GitHub"
            class="hover:text-slate-300 transition-colors py-2 px-1 min-h-[36px] flex items-center"
          >
            Releases
          </a>
          <a
            href="https://github.com/pasqualotodiogenes/sleepoff/issues"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Reportar Bug — abrir issue no GitHub"
            class="hover:text-slate-300 transition-colors py-2 px-1 min-h-[36px] flex items-center"
          >
            Reportar Bug
          </a>
        </div>
      </div>
    </footer>

  </main>
</template>
