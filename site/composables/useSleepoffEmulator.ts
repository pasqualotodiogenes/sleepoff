import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRuntimeConfig } from '#imports'
import { useSleepoffSiteState } from './useSleepoffSiteState'

export type EmulatorMode = 'shell' | 'tui'
export type TuiState =
  | 'SPLASH'
  | 'MENU'
  | 'CUSTOM_INPUT'
  | 'RUNNING'
  | 'CONFIRMATION'
  | 'DONE_TEST'
  | 'CANCELLED'

export interface EmulatorMission {
  id: 'start' | 'pause' | 'add_time' | 'panic'
  title: string
  desc: string
  completed: boolean
}

export interface ShellHistoryEntry {
  cmd?: string
  output?: string
  type: 'input' | 'output' | 'error'
}

export interface MenuOption {
  title: string
  desc: string
  minutes: number
}

export interface HeaderStatus {
  label: string
  tone: 'running' | 'done' | 'panic' | 'neutral'
}

export const MENU_OPTIONS: MenuOption[] = [
  { title: '15 min', desc: 'Tempo rápido', minutes: 15 },
  { title: '30 min', desc: 'Meia hora', minutes: 30 },
  { title: '45 min', desc: 'Episódio de série', minutes: 45 },
  { title: '1 hora', desc: 'Sessão completa', minutes: 60 },
  { title: '1h 30min', desc: 'Filme curto', minutes: 90 },
  { title: '2 horas', desc: 'Filme longo', minutes: 120 },
  { title: 'Personalizado', desc: 'Definir minutos (1–1440)', minutes: -1 },
]

export function formatRemaining(sec: number): string {
  const safe = Math.max(0, Math.floor(sec))
  const h = Math.floor(safe / 3600)
  const m = Math.floor((safe % 3600) / 60)
  const s = safe % 60
  if (h > 0) {
    return `${h}h ${String(m).padStart(2, '0')}m ${String(s).padStart(2, '0')}s`
  }
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

export function formatPanicTime(sec: number): string {
  return `00:${String(Math.max(0, Math.floor(sec))).padStart(2, '0')}`
}

export function parseDurationString(raw: string): number | null {
  const s = raw.trim().toLowerCase()
  if (/^\d+$/.test(s)) {
    const mins = parseInt(s, 10)
    return mins > 0 ? mins * 60 : null
  }
  const match = s.match(/^(?:(\d+(?:\.\d+)?)h)?(?:(\d+(?:\.\d+)?)m)?(?:(\d+(?:\.\d+)?)s)?$/)
  if (!match || match[0] === '') return null
  const hours = parseFloat(match[1] || '0')
  const minutes = parseFloat(match[2] || '0')
  const seconds = parseFloat(match[3] || '0')
  const total = Math.round(hours * 3600 + minutes * 60 + seconds)
  return total > 0 ? total : null
}

const INITIAL_SHELL_BANNER: ShellHistoryEntry = {
  type: 'output',
  output: [
    'sleepoff v1.0.3 (Go 1.25.6 • CGO_ENABLED=0 • Bubble Tea v1.3.10)',
    'Modo seguro (--dry-run) ativo para testes no navegador — nenhum comando real afeta o seu PC.',
    '',
    'Dica: clique em "1. Iniciar Timer" abaixo ou digite: sleepoff 15m --dry-run',
  ].join('\n'),
}

export function useSleepoffEmulator() {
  const siteState = useSleepoffSiteState()
  const {
    isPanicActive,
    isEmulatorOpen,
    isEmulatorFocused,
    isCommandPaletteOpen,
    isMuted,
    requestedQuickCommand,
    setPanicState,
    toggleEmulatorFocus,
  } = siteState

  const missions = ref<EmulatorMission[]>([
    { id: 'start', title: '1. Iniciar Timer', desc: 'Roda sleepoff 15m --dry-run', completed: false },
    { id: 'pause', title: '2. Pausar / Retomar', desc: 'Pressione [p] ou [Espaço] rodando', completed: false },
    { id: 'add_time', title: '3. Somar +5 min', desc: 'Pressione [+] para estender +5m', completed: false },
    { id: 'panic', title: '4. Abortar Alarme', desc: 'Pressione [c] ou [Esc] no Pânico 10s', completed: false },
  ])

  const mode = ref<EmulatorMode>('shell')
  const tuiState = ref<TuiState>('SPLASH')
  const isDryRun = ref(true)

  const shellInput = ref('')
  const shellInputRef = ref<HTMLInputElement | null>(null)
  const terminalContainerRef = ref<HTMLDivElement | null>(null)
  const terminalBodyRef = ref<HTMLDivElement | null>(null)

  const shellHistory = ref<ShellHistoryEntry[]>([INITIAL_SHELL_BANNER])
  const commandHistory = ref<string[]>([])
  const historyIndex = ref(-1)

  const selectedMenuIdx = ref(0)
  const customMinutes = ref('45')
  const totalSeconds = ref(900)
  const remainingSeconds = ref(900)
  const isPaused = ref(false)
  const eventsLog = ref<string[]>([])
  const panicCountdown = ref(10)
  const cancelPending = ref(false)
  const startTimeStr = ref('--:--')
  const finishTimeStr = ref('--:--')
  const isAudioPlaying = ref(false)

  let runningInterval: ReturnType<typeof setInterval> | null = null
  let panicInterval: ReturnType<typeof setInterval> | null = null
  let splashTimeout: ReturnType<typeof setTimeout> | null = null
  let cancelPendingTimeout: ReturnType<typeof setTimeout> | null = null
  let audioVisualizerTimer: ReturnType<typeof setTimeout> | null = null
  let audioCtx: AudioContext | null = null
  let countdownAudioTemplate: HTMLAudioElement | null = null
  let activeCountdownAudio: HTMLAudioElement | null = null
  let pendingStartSeconds = -1

  const activeStartingMissionId = computed(() => {
    const first = missions.value.find(m => !m.completed)
    return first ? first.id : null
  })

  const completedMissionsCount = computed(() => missions.value.filter(m => m.completed).length)
  const allMissionsCompleted = computed(() => missions.value.every(m => m.completed))

  const headerStatus = computed<HeaderStatus>(() => {
    if (tuiState.value === 'CONFIRMATION') {
      return { label: `! PÂNICO ${Math.max(0, panicCountdown.value)}s`, tone: 'panic' }
    }
    if (mode.value === 'tui' && tuiState.value === 'RUNNING') {
      return isPaused.value
        ? { label: 'PAUSADO', tone: 'neutral' }
        : { label: '• RODANDO', tone: 'running' }
    }
    if (mode.value === 'tui' && (tuiState.value === 'DONE_TEST' || tuiState.value === 'CANCELLED')) {
      return { label: '✓ FEITO', tone: 'done' }
    }
    if (allMissionsCompleted.value) {
      return { label: '✓ FEITO', tone: 'done' }
    }
    return mode.value === 'tui'
      ? { label: 'PRONTO', tone: 'neutral' }
      : { label: 'SHELL', tone: 'neutral' }
  })

  const progressPercent = computed(() => {
    if (totalSeconds.value <= 0) return 0
    const elapsed = totalSeconds.value - remainingSeconds.value
    return Math.min(100, Math.max(0, Math.round((elapsed / totalSeconds.value) * 100)))
  })

  const cancelSummaryStats = computed(() => {
    const elapsed = Math.max(0, totalSeconds.value - remainingSeconds.value)
    const totalMins = Math.floor(totalSeconds.value / 60)
    const elapsedMins = Math.floor(elapsed / 60)
    const elapsedSecs = elapsed % 60
    const pct = totalSeconds.value > 0 ? Math.min(100, Math.round((elapsed / totalSeconds.value) * 100)) : 0

    let msg = 'Desistiu rápido — timer abortado com segurança.'
    if (pct >= 80) msg = 'QUASE LÁ! Tão perto do final!'
    else if (pct >= 50) msg = 'PASSOU DA METADE! Respeitável.'
    else if (pct >= 20) msg = 'Cancelado antes do encerramento.'

    const barLen = 28
    const filled = Math.min(barLen, Math.round((pct * barLen) / 100))
    const asciiBar = '█'.repeat(filled) + '░'.repeat(Math.max(0, barLen - filled))
    return { elapsedMins, elapsedSecs, totalMins, pct, msg, asciiBar }
  })

  function prefersReducedMotion(): boolean {
    if (!import.meta.client || typeof window.matchMedia !== 'function') return false
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches
  }

  function updateScheduleTimes(secs: number) {
    const now = new Date()
    startTimeStr.value = now.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
    const finish = new Date(now.getTime() + secs * 1000)
    finishTimeStr.value = finish.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
  }

  function addLogEvent(msg: string) {
    const time = new Date().toTimeString().split(' ')[0] || '00:00:00'
    eventsLog.value.unshift(`[${time}] ${msg}`)
    if (eventsLog.value.length > 5) eventsLog.value.pop()
  }

  function triggerAudioVisualizer(duration = 0.15) {
    isAudioPlaying.value = true
    if (audioVisualizerTimer) clearTimeout(audioVisualizerTimer)
    audioVisualizerTimer = setTimeout(() => {
      isAudioPlaying.value = false
    }, duration * 1000 + 50)
  }

  function getAudioContext(): AudioContext | null {
    if (!import.meta.client) return null
    if (!audioCtx) {
      const Ctx = window.AudioContext || (window as unknown as { webkitAudioContext?: typeof AudioContext }).webkitAudioContext
      if (Ctx) audioCtx = new Ctx()
    }
    if (audioCtx && audioCtx.state === 'suspended') {
      audioCtx.resume().catch(() => {})
    }
    return audioCtx
  }

  function playBeep(freq = 800, duration = 0.08, type: OscillatorType = 'sine') {
    if (isMuted.value || !import.meta.client) return
    try {
      const ctx = getAudioContext()
      if (!ctx) return
      const osc = ctx.createOscillator()
      const gain = ctx.createGain()
      osc.type = type
      osc.frequency.setValueAtTime(freq, ctx.currentTime)
      gain.gain.setValueAtTime(0.12, ctx.currentTime)
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + duration)
      osc.connect(gain)
      gain.connect(ctx.destination)
      osc.start()
      osc.stop(ctx.currentTime + duration)
      triggerAudioVisualizer(duration)
    } catch {
      // Browser autoplay restriction ignored
    }
  }

  function preloadCountdownSound() {
    if (!import.meta.client) return
    try {
      const config = useRuntimeConfig()
      const base = (config.app.baseURL || '/sleepoff/').replace(/\/?$/, '/')
      const soundUrl = `${base}sounds/countdown.wav`
      if (!countdownAudioTemplate) {
        countdownAudioTemplate = new Audio(soundUrl)
        countdownAudioTemplate.volume = 0.85
        countdownAudioTemplate.preload = 'auto'
      }
    } catch {
      // Ignore if audio unavailable
    }
  }

  function stopPanicAudio() {
    if (activeCountdownAudio) {
      try {
        activeCountdownAudio.pause()
        activeCountdownAudio.currentTime = 0
      } catch {
        // Ignore cleanup failure
      }
      activeCountdownAudio = null
    }
    isAudioPlaying.value = false
  }

  function playOfficialCountdownSound() {
    if (isMuted.value || !import.meta.client) return
    try {
      preloadCountdownSound()
      if (!countdownAudioTemplate) {
        playBeep(600 + (10 - panicCountdown.value) * 60, 0.15, 'sawtooth')
        return
      }
      stopPanicAudio()
      const audioClone = countdownAudioTemplate.cloneNode() as HTMLAudioElement
      audioClone.volume = 0.85
      activeCountdownAudio = audioClone
      const promise = audioClone.play()
      if (promise !== undefined) {
        promise
          .then(() => triggerAudioVisualizer(0.2))
          .catch(() => playBeep(600 + (10 - panicCountdown.value) * 60, 0.15, 'sawtooth'))
      }
    } catch {
      playBeep(600 + (10 - panicCountdown.value) * 60, 0.15, 'sawtooth')
    }
  }

  function toggleMute() {
    isMuted.value = !isMuted.value
    if (isMuted.value) {
      stopPanicAudio()
    } else {
      getAudioContext()
      playBeep(880, 0.08)
    }
  }

  async function fireConfettiCelebration() {
    if (!import.meta.client || prefersReducedMotion()) return
    try {
      const confettiModule = await import('canvas-confetti')
      const confetti = confettiModule.default || confettiModule
      confetti({
        particleCount: 90,
        spread: 72,
        origin: { y: 0.65 },
        colors: ['#5266eb', '#ededf3', '#22c55e'],
      })
    } catch {
      // Non-critical visual enhancement
    }
  }

  function triggerMission(id: EmulatorMission['id']) {
    const target = missions.value.find(item => item.id === id)
    if (target && !target.completed) {
      target.completed = true
      playBeep(1160, 0.12)
      if (missions.value.every(item => item.completed)) {
        fireConfettiCelebration()
      }
    }
  }

  function clearTimers() {
    if (runningInterval) {
      clearInterval(runningInterval)
      runningInterval = null
    }
    if (panicInterval) {
      clearInterval(panicInterval)
      panicInterval = null
    }
    if (splashTimeout) {
      clearTimeout(splashTimeout)
      splashTimeout = null
    }
    if (cancelPendingTimeout) {
      clearTimeout(cancelPendingTimeout)
      cancelPendingTimeout = null
    }
  }

  function scrollToBottom() {
    nextTick(() => {
      if (terminalBodyRef.value) {
        terminalBodyRef.value.scrollTop = terminalBodyRef.value.scrollHeight
      }
    })
  }

  function scrollEmulatorIntoViewIfNeeded() {
    if (!import.meta.client) return
    nextTick(() => {
      const el = document.getElementById('emulador') || terminalContainerRef.value
      if (!el) return
      const rect = el.getBoundingClientRect()
      const offscreen = rect.top < 40 || rect.bottom > window.innerHeight - 20
      if (offscreen) {
        el.scrollIntoView({
          behavior: prefersReducedMotion() ? 'auto' : 'smooth',
          block: 'center',
        })
      }
    })
  }

  function focusTerminal() {
    isEmulatorOpen.value = true
    isEmulatorFocused.value = true
    getAudioContext()
    preloadCountdownSound()
    if (mode.value === 'shell') {
      nextTick(() => shellInputRef.value?.focus())
    }
  }

  function setConsoleOpen(open: boolean) {
    if (open) {
      isEmulatorOpen.value = true
      isEmulatorFocused.value = true
      scrollEmulatorIntoViewIfNeeded()
      getAudioContext()
      preloadCountdownSound()
      if (mode.value === 'shell') {
        nextTick(() => shellInputRef.value?.focus())
      }
      return
    }
    isEmulatorFocused.value = false
    isEmulatorOpen.value = false
    shellInputRef.value?.blur()
  }

  function toggleConsoleOpen() {
    setConsoleOpen(!isEmulatorOpen.value)
  }

  function toggleConsoleWithShortcut() {
    if (!isEmulatorOpen.value) {
      setConsoleOpen(true)
      return
    }
    toggleEmulatorFocus()
    if (isEmulatorFocused.value) {
      scrollEmulatorIntoViewIfNeeded()
      getAudioContext()
      preloadCountdownSound()
      if (mode.value === 'shell') {
        nextTick(() => shellInputRef.value?.focus())
      }
    } else {
      shellInputRef.value?.blur()
    }
  }

  function switchTab(target: EmulatorMode) {
    mode.value = target
    isEmulatorOpen.value = true
    if (target === 'shell') {
      nextTick(() => {
        shellInputRef.value?.focus()
        scrollToBottom()
      })
    } else if (
      tuiState.value === 'SPLASH' ||
      tuiState.value === 'DONE_TEST' ||
      tuiState.value === 'CANCELLED'
    ) {
      startTui(-1, true, true)
    }
  }

  function startTui(seconds: number, dry = true, skipSplash = false) {
    clearTimers()
    stopPanicAudio()
    setPanicState(false, 0)
    mode.value = 'tui'
    isDryRun.value = dry
    cancelPending.value = false
    pendingStartSeconds = seconds
    playBeep(600, 0.08)

    if (skipSplash) {
      if (seconds > 0) launchTimer(seconds)
      else {
        tuiState.value = 'MENU'
        selectedMenuIdx.value = 0
      }
      return
    }

    tuiState.value = 'SPLASH'
    splashTimeout = setTimeout(() => {
      pendingStartSeconds = -1
      if (seconds > 0) launchTimer(seconds)
      else {
        tuiState.value = 'MENU'
        selectedMenuIdx.value = 0
      }
    }, 550)
  }

  function launchTimer(seconds: number, dryOverride?: boolean) {
    clearTimers()
    stopPanicAudio()
    setPanicState(false, 0)
    if (typeof dryOverride === 'boolean') isDryRun.value = dryOverride
    mode.value = 'tui'
    totalSeconds.value = seconds
    remainingSeconds.value = seconds
    isPaused.value = false
    cancelPending.value = false
    eventsLog.value = []
    updateScheduleTimes(seconds)

    const formatted = formatRemaining(seconds)
    addLogEvent(`Timer iniciado: ${formatted}${isDryRun.value ? ' (--dry-run)' : ''}`)
    triggerMission('start')
    playBeep(900, 0.09)

    if (seconds <= 10) {
      startPanic(seconds)
      return
    }

    tuiState.value = 'RUNNING'
    runningInterval = setInterval(() => {
      if (isPaused.value) return
      if (remainingSeconds.value > 10) {
        remainingSeconds.value -= 1
      } else if (remainingSeconds.value > 0) {
        remainingSeconds.value -= 1
        startPanic(10)
      } else {
        finishDone()
      }
    }, 1000)
  }

  function startPanic(startSeconds = 10) {
    if (runningInterval) {
      clearInterval(runningInterval)
      runningInterval = null
    }
    if (panicInterval) {
      clearInterval(panicInterval)
      panicInterval = null
    }
    mode.value = 'tui'
    isEmulatorOpen.value = true
    tuiState.value = 'CONFIRMATION'
    panicCountdown.value = startSeconds
    cancelPending.value = false
    setPanicState(true, startSeconds)
    addLogEvent(`! Alarme de pânico ativado (${startSeconds}s)`)
    playOfficialCountdownSound()

    panicInterval = setInterval(() => {
      if (panicCountdown.value > 1) {
        panicCountdown.value -= 1
        setPanicState(true, panicCountdown.value)
        playOfficialCountdownSound()
      } else {
        if (panicInterval) {
          clearInterval(panicInterval)
          panicInterval = null
        }
        finishDone()
      }
    }, 1000)
  }

  function triggerPanicDemo() {
    isDryRun.value = true
    if (tuiState.value !== 'RUNNING' || totalSeconds.value < 10) {
      totalSeconds.value = 10
      remainingSeconds.value = 10
      updateScheduleTimes(10)
    }
    startPanic(10)
  }

  function finishDone() {
    clearTimers()
    stopPanicAudio()
    setPanicState(false, 0)
    tuiState.value = 'DONE_TEST'
    addLogEvent(isDryRun.value ? 'Simulação --dry-run concluída!' : 'Contagem finalizada!')
    playBeep(420, 0.3, 'triangle')
  }

  function togglePause() {
    if (mode.value !== 'tui' || tuiState.value !== 'RUNNING') return
    isPaused.value = !isPaused.value
    addLogEvent(isPaused.value ? 'Timer pausado pelo usuário' : 'Timer retomado')
    triggerMission('pause')
    playBeep(isPaused.value ? 550 : 850, 0.08)
  }

  function addFiveMinutes() {
    if (mode.value !== 'tui' || tuiState.value !== 'RUNNING') return
    remainingSeconds.value += 300
    totalSeconds.value += 300
    updateScheduleTimes(remainingSeconds.value)
    addLogEvent('+5 min adicionados (+300s)')
    triggerMission('add_time')
    playBeep(1000, 0.08)
  }

  function subtractFiveMinutes() {
    if (mode.value !== 'tui' || tuiState.value !== 'RUNNING') return
    // Go internal/model/update.go:269 enforces 1-minute floor (remaining - 5m >= 1m)
    if (remainingSeconds.value - 300 >= 60) {
      remainingSeconds.value -= 300
      totalSeconds.value = Math.max(remainingSeconds.value, totalSeconds.value - 300)
      updateScheduleTimes(remainingSeconds.value)
      addLogEvent('-5 min subtraídos (-300s)')
      playBeep(600, 0.08)
    } else {
      addLogEvent('Limite mínimo de 1 minuto atingido')
    }
  }

  function handleCancel() {
    if (mode.value !== 'tui') return

    if (tuiState.value === 'CONFIRMATION') {
      clearTimers()
      stopPanicAudio()
      setPanicState(false, 0)
      tuiState.value = 'CANCELLED'
      addLogEvent('Alarme abortado em modo --dry-run!')
      triggerMission('panic')
      playBeep(360, 0.2)
      return
    }

    if (tuiState.value === 'RUNNING') {
      if (cancelPending.value) {
        clearTimers()
        stopPanicAudio()
        setPanicState(false, 0)
        cancelPending.value = false
        tuiState.value = 'CANCELLED'
        addLogEvent('Timer cancelado pelo usuário')
        playBeep(360, 0.2)
      } else {
        cancelPending.value = true
        playBeep(700, 0.08)
        if (cancelPendingTimeout) clearTimeout(cancelPendingTimeout)
        cancelPendingTimeout = setTimeout(() => {
          cancelPending.value = false
        }, 2500)
      }
    }
  }

  function exitTuiToShell() {
    clearTimers()
    stopPanicAudio()
    setPanicState(false, 0)
    cancelPending.value = false
    mode.value = 'shell'
    tuiState.value = 'SPLASH'
    shellHistory.value.push({
      type: 'output',
      output: '[sleepoff encerrado — retornou ao prompt shell]',
    })
    scrollToBottom()
    nextTick(() => shellInputRef.value?.focus())
  }

  function handleShellCommand(rawCmd: string) {
    const cmd = rawCmd.trim()
    if (!cmd) return

    commandHistory.value.push(cmd)
    historyIndex.value = -1
    shellHistory.value.push({ type: 'input', cmd })

    const parts = cmd.split(/\s+/)
    let base = (parts[0] || '').toLowerCase()
    if (base === './sleepoff' || base === '.\\sleepoff' || base === 'sleepoff.exe' || base === './sleepoff.exe') {
      base = 'sleepoff'
    }

    if (base === 'clear' || base === 'cls') {
      shellHistory.value = []
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base === 'panic-demo') {
      shellInput.value = ''
      triggerPanicDemo()
      return
    }

    if (base === 'help') {
      shellHistory.value.push({
        type: 'output',
        output: [
          'Comandos do emulador sleepoff v1.0.3:',
          '  sleepoff                   Abre o menu interativo TUI',
          '  sleepoff 15m --dry-run     Inicia timer de 15 min em modo seguro',
          '  sleepoff <tempo>           Formatos: 30, 30m, 90s, 1h, 1h30m',
          '  panic-demo                 Simula diretamente a tela de Pânico 10s',
          '  sleepoff --version         Exibe a versão oficial (v1.0.3)',
          '  sleepoff --help            Exibe a ajuda oficial do CLI Go',
          '  ls | pwd | whoami | date | uname | clear',
        ].join('\n'),
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base === 'ls' || base === 'dir') {
      shellHistory.value.push({
        type: 'output',
        output: 'go.mod  go.sum  main.go  main_test.go  README.md  internal/  packaging/  site/',
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base === 'pwd') {
      shellHistory.value.push({ type: 'output', output: '/home/user/sleepoff' })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base === 'whoami') {
      shellHistory.value.push({ type: 'output', output: 'user' })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base === 'date') {
      shellHistory.value.push({ type: 'output', output: new Date().toString() })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base === 'uname') {
      shellHistory.value.push({
        type: 'output',
        output: 'Linux sleepoff-node 6.11.0-x86_64 GNU/Linux (Windows/macOS também suportados)',
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (base !== 'sleepoff') {
      const directDuration = parseDurationString(base)
      if (directDuration !== null) {
        const rest = parts.slice(1)
        if (rest.includes('-d')) {
          shellHistory.value.push({
            type: 'error',
            output: 'Use --dry-run (flag oficial do sleepoff v1.0.3; não há flag -d).',
          })
          shellInput.value = ''
          scrollToBottom()
          return
        }
        shellInput.value = ''
        startTui(directDuration, rest.includes('--dry-run'), true)
        return
      }

      shellHistory.value.push({
        type: 'error',
        output: `Comando não reconhecido: "${base}". Digite "sleepoff 15m --dry-run", "help" ou use os botões abaixo.`,
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    const args = parts.slice(1)
    if (args.includes('-d')) {
      shellHistory.value.push({
        type: 'error',
        output: 'Aviso: no sleepoff v1.0.3 a flag oficial é --dry-run (sem atalho -d). Ex.: sleepoff 15m --dry-run',
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    const isDry = args.includes('--dry-run')
    const cleanArgs = args.filter(a => a !== '--dry-run')

    if (cleanArgs.includes('--help') || cleanArgs.includes('-h')) {
      shellHistory.value.push({
        type: 'output',
        output: [
          'sleepoff v1.0.3 — Timer de desligamento automático',
          'Uso:',
          '  sleepoff                   Menu interativo',
          '  sleepoff <tempo>           Inicia timer direto (30, 30m, 90s, 1h, 1h30m)',
          'Opções:',
          '  --dry-run                  Modo teste (não desliga de verdade)',
          '  -h, --help                 Mostra esta ajuda',
          '  -v, --version              Mostra a versão',
        ].join('\n'),
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (cleanArgs.includes('--version') || cleanArgs.includes('-v')) {
      shellHistory.value.push({
        type: 'output',
        output: 'sleepoff version 1.0.3 (Go 1.25.6, CGO_ENABLED=0, Windows/Linux/macOS)',
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (tuiState.value === 'RUNNING' || tuiState.value === 'CONFIRMATION') {
      shellHistory.value.push({
        type: 'error',
        output: `sleepoff: trava Single-Instance ativa — já existe um timer em TUI (${formatRemaining(remainingSeconds.value)}).`,
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    if (cleanArgs.length === 0) {
      shellInput.value = ''
      startTui(-1, isDry, true)
      return
    }

    const timeArg = cleanArgs[0] || ''
    const parsedSecs = parseDurationString(timeArg)
    if (!parsedSecs) {
      shellHistory.value.push({
        type: 'error',
        output: `Duração inválida: "${timeArg}". Formatos aceitos: 30, 30m, 90s, 1h, 1h30m.`,
      })
      shellInput.value = ''
      scrollToBottom()
      return
    }

    shellInput.value = ''
    startTui(parsedSecs, isDry, true)
  }

  function runQuickCommand(cmd: string) {
    focusTerminal()
    shellInput.value = cmd
    handleShellCommand(cmd)
  }

  function selectMenuOption(idx: number) {
    selectedMenuIdx.value = idx
    const opt = MENU_OPTIONS[idx]
    if (!opt) return
    if (opt.minutes === -1) {
      tuiState.value = 'CUSTOM_INPUT'
    } else {
      launchTimer(opt.minutes * 60)
    }
  }

  function submitCustomMinutes() {
    const val = parseInt(customMinutes.value, 10)
    if (!Number.isNaN(val) && val >= 1 && val <= 1440) {
      launchTimer(val * 60)
    }
  }

  function resetMissions() {
    clearTimers()
    stopPanicAudio()
    setPanicState(false, 0)
    missions.value.forEach(m => {
      m.completed = false
    })
    mode.value = 'shell'
    tuiState.value = 'SPLASH'
    isDryRun.value = true
    isPaused.value = false
    cancelPending.value = false
  }

  function clickMission(mission: EmulatorMission) {
    focusTerminal()

    // Ensure any pending splash resolves immediately
    if (tuiState.value === 'SPLASH' && pendingStartSeconds > 0) {
      const secs = pendingStartSeconds
      pendingStartSeconds = -1
      launchTimer(secs)
    }

    if (mission.id === 'start') {
      if (tuiState.value === 'RUNNING') {
        switchTab('tui')
        triggerMission('start')
      } else {
        if (tuiState.value === 'CONFIRMATION') {
          clearTimers()
          stopPanicAudio()
          setPanicState(false, 0)
        }
        startTui(900, true, true)
      }
      return
    }

    if (mission.id === 'pause') {
      if (mode.value !== 'tui') mode.value = 'tui'
      if (tuiState.value !== 'RUNNING') {
        launchTimer(900, true)
      }
      togglePause()
      return
    }

    if (mission.id === 'add_time') {
      if (mode.value !== 'tui') mode.value = 'tui'
      if (tuiState.value !== 'RUNNING') {
        launchTimer(900, true)
      }
      addFiveMinutes()
      return
    }

    if (mission.id === 'panic') {
      if (mode.value !== 'tui') mode.value = 'tui'
      if (tuiState.value === 'CONFIRMATION') {
        handleCancel()
      } else {
        triggerPanicDemo()
      }
    }
  }

  function onKeyDown(e: KeyboardEvent) {
    const activeEl = document.activeElement
    const tag = activeEl?.tagName?.toUpperCase() || ''
    const isAnyInput =
      tag === 'INPUT' ||
      tag === 'TEXTAREA' ||
      tag === 'SELECT' ||
      (activeEl instanceof HTMLElement && activeEl.isContentEditable)
    const isShellInput = Boolean(shellInputRef.value && activeEl === shellInputRef.value)
    const isOuterInput = isAnyInput && !isShellInput

    // Global [~] shortcut: allow toggling even while #terminal-shell-input is focused,
    // but never hijack CommandPalette or unrelated external inputs.
    if (
      !isCommandPaletteOpen.value &&
      !isOuterInput &&
      !e.ctrlKey &&
      !e.metaKey &&
      !e.altKey &&
      (e.key === '~' || e.code === 'Backquote')
    ) {
      e.preventDefault()
      toggleConsoleWithShortcut()
      return
    }

    if (isOuterInput || !isEmulatorFocused.value) return

    if (mode.value === 'shell') {
      if (e.key === 'Tab') {
        e.preventDefault()
        const val = shellInput.value.trim()
        const candidates = [
          'sleepoff 15m --dry-run',
          'sleepoff 30m --dry-run',
          'sleepoff 1h30m',
          'panic-demo',
          'sleepoff --help',
          'sleepoff --version',
          'clear',
        ]
        const found = candidates.find(c => c.startsWith(val))
        if (found) shellInput.value = found
        return
      }
      if (e.ctrlKey && e.key.toLowerCase() === 'l') {
        e.preventDefault()
        shellHistory.value = []
        return
      }
      if (e.key === 'ArrowUp') {
        if (commandHistory.value.length > 0) {
          if (historyIndex.value === -1) historyIndex.value = commandHistory.value.length - 1
          else if (historyIndex.value > 0) historyIndex.value -= 1
          shellInput.value = commandHistory.value[historyIndex.value] || ''
        }
        e.preventDefault()
      } else if (e.key === 'ArrowDown') {
        if (historyIndex.value !== -1) {
          if (historyIndex.value < commandHistory.value.length - 1) {
            historyIndex.value += 1
            shellInput.value = commandHistory.value[historyIndex.value] || ''
          } else {
            historyIndex.value = -1
            shellInput.value = ''
          }
        }
        e.preventDefault()
      }
      return
    }

    if (tuiState.value === 'MENU') {
      if (e.key === 'ArrowUp' || e.key === 'k') {
        selectedMenuIdx.value = (selectedMenuIdx.value - 1 + MENU_OPTIONS.length) % MENU_OPTIONS.length
        playBeep(750, 0.05)
        e.preventDefault()
      } else if (e.key === 'ArrowDown' || e.key === 'j') {
        selectedMenuIdx.value = (selectedMenuIdx.value + 1) % MENU_OPTIONS.length
        playBeep(750, 0.05)
        e.preventDefault()
      } else if (e.key === 'Enter') {
        selectMenuOption(selectedMenuIdx.value)
        e.preventDefault()
      } else if (e.key === 'q' || e.key === 'Escape') {
        exitTuiToShell()
        e.preventDefault()
      }
    } else if (tuiState.value === 'CUSTOM_INPUT') {
      if (e.key === 'Enter') {
        submitCustomMinutes()
        e.preventDefault()
      } else if (e.key === 'Escape') {
        tuiState.value = 'MENU'
        e.preventDefault()
      }
    } else if (tuiState.value === 'RUNNING') {
      if (e.key === 'p' || e.key === ' ' || e.code === 'Space') {
        togglePause()
        e.preventDefault()
      } else if (e.key === '+' || e.key === '=') {
        addFiveMinutes()
        e.preventDefault()
      } else if (e.key === '-' || e.key === '_') {
        subtractFiveMinutes()
        e.preventDefault()
      } else if (e.key === 'c' || e.key === 'q' || e.key === 'Escape') {
        handleCancel()
        e.preventDefault()
      }
    } else if (tuiState.value === 'CONFIRMATION') {
      if (e.key === 'c' || e.key === 'Escape' || e.key === 'q') {
        handleCancel()
        e.preventDefault()
      }
    } else if (tuiState.value === 'DONE_TEST' || tuiState.value === 'CANCELLED') {
      if (e.key === 'Enter' || e.key === 'q' || e.key === 'Escape') {
        exitTuiToShell()
        e.preventDefault()
      }
    }
  }

  function handleDocumentClick(e: MouseEvent) {
    const target = e.target as HTMLElement | null
    if (!target || !target.isConnected || target.closest('.emulator-focus-btn')) {
      return
    }
    if (terminalContainerRef.value && terminalContainerRef.value.contains(target)) {
      isEmulatorFocused.value = true
    } else {
      isEmulatorFocused.value = false
    }
  }

  watch(isEmulatorFocused, (focused) => {
    if (!import.meta.client) return
    if (focused) {
      getAudioContext()
      preloadCountdownSound()
      if (mode.value === 'shell' && isEmulatorOpen.value) {
        nextTick(() => shellInputRef.value?.focus())
      }
    } else {
      shellInputRef.value?.blur()
    }
  })

  watch(shellInputRef, (el) => {
    if (!import.meta.client) return
    if (el && isEmulatorOpen.value && isEmulatorFocused.value && mode.value === 'shell') {
      el.focus()
    }
  })

  watch(
    requestedQuickCommand,
    req => {
      if (!req) return
      focusTerminal()
      scrollEmulatorIntoViewIfNeeded()
      if (req.command === 'panic-demo') {
        triggerPanicDemo()
      } else {
        if (tuiState.value === 'RUNNING' || tuiState.value === 'CONFIRMATION') {
          clearTimers()
          stopPanicAudio()
          setPanicState(false, 0)
        }
        runQuickCommand(req.command)
      }
    },
    { flush: 'post' },
  )

  onMounted(() => {
    if (!import.meta.client) return
    window.addEventListener('keydown', onKeyDown)
    document.addEventListener('click', handleDocumentClick)
    const unlockAudio = () => {
      getAudioContext()
      preloadCountdownSound()
      window.removeEventListener('pointerdown', unlockAudio)
      window.removeEventListener('keydown', unlockAudio)
    }
    window.addEventListener('pointerdown', unlockAudio, { passive: true })
    window.addEventListener('keydown', unlockAudio, { passive: true })
  })

  onUnmounted(() => {
    if (import.meta.client) {
      window.removeEventListener('keydown', onKeyDown)
      document.removeEventListener('click', handleDocumentClick)
    }
    clearTimers()
    stopPanicAudio()
    setPanicState(false, 0)
  })

  return {
    // State
    isPanicActive,
    isEmulatorOpen,
    isEmulatorFocused,
    isMuted,
    missions,
    mode,
    tuiState,
    isDryRun,
    shellInput,
    shellInputRef,
    terminalContainerRef,
    terminalBodyRef,
    shellHistory,
    selectedMenuIdx,
    customMinutes,
    totalSeconds,
    remainingSeconds,
    isPaused,
    eventsLog,
    panicCountdown,
    cancelPending,
    startTimeStr,
    finishTimeStr,
    isAudioPlaying,
    activeStartingMissionId,
    completedMissionsCount,
    allMissionsCompleted,
    headerStatus,
    progressPercent,
    cancelSummaryStats,
    // Actions
    focusTerminal,
    toggleConsoleOpen,
    toggleConsoleWithShortcut,
    switchTab,
    toggleMute,
    startTui,
    launchTimer,
    startPanic,
    triggerPanicDemo,
    togglePause,
    addFiveMinutes,
    subtractFiveMinutes,
    handleCancel,
    exitTuiToShell,
    handleShellCommand,
    runQuickCommand,
    selectMenuOption,
    submitCustomMinutes,
    clickMission,
    resetMissions,
    formatRemaining,
    formatPanicTime,
  }
}
