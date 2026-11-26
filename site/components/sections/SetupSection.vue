<script setup lang="ts">
import { createReusableTemplate } from '@vueuse/core'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { ClientOnly, Icon } from '#components'
import AnimatedTabs from '~/components/ui/animated-tabs/AnimatedTabs.vue'
import type { AnimatedTab } from '~/components/ui/animated-tabs/types'
import EncryptedText from '~/components/ui/encrypted-text/EncryptedText.vue'

type PlatformId = 'windows-app' | 'go-install' | 'linux' | 'macos' | 'windows-zip'

interface PlatformSetupItem {
  id: PlatformId
  railLabel: string
  shortLabel: string
  badge: string
  title: string
  description: string
  highlights: string[]
  commandLabel: string
  command: string
  secondaryNote?: string
  downloadFile?: string
  secondaryLinkLabel: string
  secondaryLinkUrl?: string
}

const REPO_URL = 'https://github.com/pasqualotodiogenes/sleepoff'
const RELEASE_BASE = `${REPO_URL}/releases/latest`

const platforms: PlatformSetupItem[] = [
  {
    id: 'windows-app',
    railLabel: 'App Windows (Instalador)',
    shortLabel: 'App Windows',
    badge: 'RECOMENDADO // SEM PRECISAR DE TERMINAL',
    title: 'Instalador Gráfico Windows (sleepoff-setup.exe)',
    description:
      'Instale como qualquer aplicativo tradicional no Windows. Cria atalhos no Menu Iniciar e Área de Trabalho, abre com dois cliques, exibe ícone na bandeja do sistema (System Tray) com alertas e verifica atualizações semanais.',
    highlights: [
      'Menu Iniciar, Área de Trabalho e execução por duplo clique (sem terminal)',
      'Bandeja do sistema (System Tray) com Pausar/Retomar, Cancelar e notificações de 10s',
    ],
    commandLabel: 'Opcional — iniciar após instalar via Win+R ou terminal',
    command: 'sleepoff 30',
    downloadFile: 'sleepoff-setup.exe',
    secondaryLinkLabel: 'Ver release v1.0.3 no GitHub',
  },
  {
    id: 'go-install',
    railLabel: 'Go (go install)',
    shortLabel: 'Go install',
    badge: 'GO 1.25+ // COMPILAÇÃO DIRETA',
    title: 'Instalação via Toolchain Go',
    description:
      'Compila o binário estático diretamente do repositório oficial em Go 1.25+ e instala o executável sleepoff no seu $GOPATH/bin.',
    highlights: [
      'Compilação limpa diretamente do módulo GitHub oficial',
      'Compatível com Windows, Linux e macOS (CGO_ENABLED=0)',
    ],
    commandLabel: 'Comando completo (Go 1.25+)',
    command: 'go install github.com/pasqualotodiogenes/sleepoff@latest',
    secondaryNote: 'Certifique-se de que $(go env GOPATH)/bin esteja incluído no seu PATH.',
    secondaryLinkLabel: 'Ver código-fonte no GitHub',
    secondaryLinkUrl: REPO_URL,
  },
  {
    id: 'linux',
    railLabel: 'Linux (Binário Estático)',
    shortLabel: 'Linux',
    badge: 'AMD64 & ARM64 // ZERO DEPENDÊNCIAS',
    title: 'Binário Estático para Linux',
    description:
      'Executável único sem bibliotecas dinâmicas. Em sessões desktop usa systemctl poweroff via systemd-logind / D-Bus sem exigir sudo.',
    highlights: [
      'Lock de instância única via flock em $XDG_RUNTIME_DIR/sleepoff.lock',
      'Áudio de pânico compatível com pw-play, paplay, aplay e terminal bell',
    ],
    commandLabel: 'x86_64 (amd64) — baixar, extrair e mover para /usr/local/bin',
    command:
      'curl -fsSL https://github.com/pasqualotodiogenes/sleepoff/releases/latest/download/sleepoff_linux_amd64.tar.gz -o sleepoff.tar.gz && tar -xzf sleepoff.tar.gz && sudo mv sleepoff /usr/local/bin/',
    secondaryNote: 'Para ARM64 (aarch64), utilize sleepoff_linux_arm64.tar.gz.',
    downloadFile: 'sleepoff_linux_amd64.tar.gz',
    secondaryLinkLabel: 'Ver todos os artefatos Linux',
  },
  {
    id: 'macos',
    railLabel: 'macOS (Binário Estático)',
    shortLabel: 'macOS',
    badge: 'APPLE SILICON (ARM64) & INTEL (AMD64)',
    title: 'Binário Estático para macOS',
    description:
      'Binário nativo para Apple Silicon e Intel. Executa o desligamento via Apple Events (osascript), solicitando permissão de Automação na primeira execução.',
    highlights: [
      'Alarme sonoro nativo via afplay com fallback do sistema',
      'Lock atômico via flock em /tmp/sleepoff.lock contra timers duplicados',
    ],
    commandLabel: 'Apple Silicon (arm64) — baixar, extrair e mover para /usr/local/bin',
    command:
      'curl -fsSL https://github.com/pasqualotodiogenes/sleepoff/releases/latest/download/sleepoff_darwin_arm64.tar.gz -o sleepoff.tar.gz && tar -xzf sleepoff.tar.gz && sudo mv sleepoff /usr/local/bin/',
    secondaryNote: 'Para Macs Intel (x86_64), utilize sleepoff_darwin_amd64.tar.gz.',
    downloadFile: 'sleepoff_darwin_arm64.tar.gz',
    secondaryLinkLabel: 'Ver todos os artefatos macOS',
  },
  {
    id: 'windows-zip',
    railLabel: 'Windows Portátil (.zip)',
    shortLabel: 'Win .zip',
    badge: 'PORTÁTIL // DUPLO CLIQUE OU POWERSHELL',
    title: 'Pacote Portátil Windows (.zip)',
    description:
      'Extraia o arquivo e dê dois cliques diretos em sleepoff.exe (MousetrapHelpText habilitado) ou execute pelo PowerShell sem instalador.',
    highlights: [
      'Funciona com duplo clique direto ou argumentos no PowerShell',
      'Motor nativo winmm.dll + CreateMutexW sem privilégios de administrador',
    ],
    commandLabel: 'PowerShell — baixar e extrair sleepoff_windows_amd64.zip',
    command:
      'Invoke-WebRequest -Uri "https://github.com/pasqualotodiogenes/sleepoff/releases/latest/download/sleepoff_windows_amd64.zip" -OutFile "sleepoff.zip"; Expand-Archive -Path "sleepoff.zip" -DestinationPath ".\\sleepoff" -Force',
    secondaryNote: 'Para Windows ARM64, utilize sleepoff_windows_arm64.zip nos releases.',
    downloadFile: 'sleepoff_windows_amd64.zip',
    secondaryLinkLabel: 'Ver checksums e releases',
  },
]

const [DefinePlatformPanel, ReusePlatformPanel] = createReusableTemplate()
const activePlatform = ref<PlatformId>('windows-app')
const copiedKey = ref<string | null>(null)
let copyTimer: ReturnType<typeof setTimeout> | null = null

const selectedItem = computed(
  () => platforms.find((item) => item.id === activePlatform.value) ?? platforms[0]!,
)
const mobileTabs = computed<AnimatedTab[]>(() =>
  platforms.map((item) => ({ title: item.shortLabel, value: item.id })),
)

function handleMobileTabUpdate(val: string) {
  if (platforms.some((item) => item.id === val)) activePlatform.value = val as PlatformId
}

async function copyCommand(text: string, key: string) {
  if (!import.meta.client) return
  try {
    if (navigator.clipboard?.writeText) await navigator.clipboard.writeText(text)
    else {
      const el = Object.assign(document.createElement('textarea'), { value: text })
      el.style.cssText = 'position:fixed;opacity:0'
      document.body.appendChild(el)
      el.select()
      document.execCommand('copy')
      document.body.removeChild(el)
    }
  } catch { /* Keep visual feedback responsive */ }
  copiedKey.value = key
  if (copyTimer) clearTimeout(copyTimer)
  copyTimer = setTimeout(() => { if (copiedKey.value === key) copiedKey.value = null }, 2200)
}

onMounted(() => {
  const ua = (navigator.userAgent || '').toLowerCase()
  if (ua.includes('mac') || ua.includes('darwin')) activePlatform.value = 'macos'
  else if (ua.includes('linux') && !ua.includes('android')) activePlatform.value = 'linux'
  else activePlatform.value = 'windows-app'
})

onUnmounted(() => {
  if (copyTimer) clearTimeout(copyTimer)
})
</script>

<template>
  <DefinePlatformPanel>
    <div class="flex min-w-0 flex-col justify-between gap-6">
      <div>
        <span class="font-mono-hud inline-block border border-[#272735] bg-[#171721] px-2.5 py-1 text-[10px] sm:text-[11px] text-[#5266eb]">{{ selectedItem.badge }}</span>
        <h3 class="font-display mt-3 text-ms-xl sm:text-ms-2xl font-bold tracking-tight text-[#ededf3]">{{ selectedItem.title }}</h3>
        <p class="mt-2 text-ms-sm leading-relaxed text-[#c3c3cc]">{{ selectedItem.description }}</p>
        <ul class="mt-4 space-y-1.5 border-l border-[#5266eb] pl-3.5 text-xs text-[#c3c3cc]">
          <li v-for="point in selectedItem.highlights" :key="point">{{ point }}</li>
        </ul>
        <div class="mt-5 border border-[#1e1e2a] bg-[#06070a]">
          <div class="flex flex-wrap items-center justify-between gap-2 border-b border-[#1e1e2a] bg-[#171721] px-3.5 py-2">
            <span class="font-mono text-[11px] text-[#70707d]">{{ selectedItem.commandLabel }}</span>
            <button
              type="button"
              class="inline-flex min-h-[44px] cursor-pointer items-center gap-2 border px-3.5 py-2 font-mono text-xs font-semibold uppercase tracking-[0.06em] transition-colors"
              :class="copiedKey === selectedItem.id ? 'border-[#5266eb] bg-[#5266eb] text-[#ededf3]' : 'border-[#272735] bg-[#272735] text-[#ededf3] hover:border-[#5266eb]'"
              :aria-label="`Copiar comando de instalação para ${selectedItem.railLabel}`"
              @click="copyCommand(selectedItem.command, selectedItem.id)"
            >
              <Icon v-if="copiedKey !== selectedItem.id" name="heroicons:clipboard-document" class="size-4 shrink-0" />
              <span>{{ copiedKey === selectedItem.id ? 'COPIADO ✓' : 'COPIAR' }}</span>
            </button>
          </div>
          <div class="w-full min-w-0 px-3.5 py-3.5 font-mono text-xs leading-relaxed text-[#ededf3] sm:overflow-x-auto">
            <code class="block whitespace-normal break-all select-all sm:inline sm:break-normal sm:whitespace-nowrap">{{ selectedItem.command }}</code>
          </div>
        </div>
        <p v-if="selectedItem.secondaryNote" class="mt-2.5 font-mono text-[11px] text-[#70707d]">{{ selectedItem.secondaryNote }}</p>
      </div>
      <div class="flex flex-col gap-2.5 border-t border-[#1e1e2a] pt-5 sm:flex-row sm:flex-wrap sm:items-center sm:gap-3">
        <a
          v-if="selectedItem.downloadFile"
          :href="`${RELEASE_BASE}/download/${selectedItem.downloadFile}`"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex min-h-[44px] w-full items-center justify-center gap-2 border border-[#5266eb] bg-[#5266eb] px-5 py-2.5 font-mono text-xs font-semibold uppercase tracking-[0.06em] text-[#ededf3] transition-colors hover:bg-[#5266eb]/90 sm:w-auto"
        >
          <Icon name="heroicons:arrow-down-tray" class="size-4 shrink-0" />
          <span>Baixar {{ selectedItem.downloadFile }}</span>
        </a>
        <a
          :href="selectedItem.secondaryLinkUrl || RELEASE_BASE"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex min-h-[44px] w-full items-center justify-center gap-2 border border-[#272735] bg-[#171721] px-5 py-2.5 font-mono text-xs font-semibold uppercase tracking-[0.06em] text-[#c3c3cc] transition-colors hover:border-[#5266eb] hover:text-[#ededf3] sm:w-auto"
        >
          <Icon name="simple-icons:github" class="size-4 shrink-0" />
          <span>{{ selectedItem.secondaryLinkLabel }}</span>
        </a>
      </div>
    </div>
  </DefinePlatformPanel>

  <section id="setup" data-section="02" class="relative z-10 border-b border-[#1e1e2a] bg-[#06070a]">
    <div class="mx-auto max-w-6xl border-x border-[#1e1e2a] px-4 py-16 sm:px-8 sm:py-20 lg:px-12">
      <div class="mb-10 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <div class="mb-3">
            <ClientOnly>
              <EncryptedText
                text="02 / SETUP"
                class="font-mono-hud text-xs text-[#5266eb]"
                encrypted-class="text-[#70707d]"
                revealed-class="text-[#5266eb]"
              />
              <template #fallback>
                <span class="font-mono-hud text-xs text-[#5266eb]">02 / SETUP</span>
              </template>
            </ClientOnly>
          </div>
          <h2 class="font-display text-ms-4xl sm:text-ms-5xl font-bold tracking-tighter text-[#ededf3]">
            Instale como <span class="font-serif-accent">aplicativo</span> ou no terminal
          </h2>
        </div>
        <p class="max-w-md text-ms-sm leading-relaxed text-[#c3c3cc]">
          Artefatos oficiais da release <code class="text-xs text-[#ededf3]">v1.0.3</code>. Zero runtimes externos, zero telemetria.
        </p>
      </div>

      <!-- Desktop: Vertical Rail + Panel -->
      <div class="hidden border border-[#1e1e2a] bg-[#0e1016] md:grid md:grid-cols-12">
        <div role="tablist" aria-label="Plataformas de instalação do sleepoff" class="border-r border-[#1e1e2a] md:col-span-4">
          <button
            v-for="(item, index) in platforms"
            :key="item.id"
            type="button"
            role="tab"
            :aria-selected="activePlatform === item.id"
            class="flex min-h-[58px] w-full cursor-pointer items-center justify-between border-b border-[#1e1e2a] px-5 py-4 text-left transition-colors last:border-b-0"
            :class="activePlatform === item.id ? 'bg-[#5266eb] text-[#ededf3]' : 'setup-rail-btn bg-[#0e1016] text-[#c3c3cc]'"
            @click="activePlatform = item.id"
          >
            <div class="flex min-w-0 items-center gap-3">
              <span class="font-mono text-xs opacity-80">{{ activePlatform === item.id ? '▶' : `0${index + 1}` }}</span>
              <span class="font-display text-sm font-semibold">{{ item.railLabel }}</span>
            </div>
            <span class="font-mono text-[10px] uppercase tracking-[0.08em] opacity-80">
              {{ item.id === 'windows-app' ? 'APP' : 'CLI' }}
            </span>
          </button>
        </div>
        <div role="tabpanel" class="min-w-0 p-6 lg:p-8 md:col-span-8">
          <ReusePlatformPanel />
        </div>
      </div>

      <!-- Mobile: AnimatedTabs -->
      <div class="md:hidden">
        <AnimatedTabs :model-value="activePlatform" :tabs="mobileTabs" @update:model-value="handleMobileTabUpdate">
          <template #content>
            <div class="min-w-0 border border-[#1e1e2a] bg-[#0e1016] p-4">
              <ReusePlatformPanel />
            </div>
          </template>
        </AnimatedTabs>
      </div>
    </div>
  </section>
</template>

<style scoped>
@media (hover: hover) {
  .setup-rail-btn:hover {
    background-color: #171721;
    color: #ededf3;
  }
}
</style>
