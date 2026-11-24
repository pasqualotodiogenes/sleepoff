<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

// --- Tipos & Estados ---
type Mode = 'shell' | 'tui'
type TuiState = 'SPLASH' | 'MENU' | 'CUSTOM_INPUT' | 'RUNNING' | 'CONFIRMATION' | 'DONE_TEST' | 'CANCELLED'

interface Mission {
  id: string
  title: string
  desc: string
  completed: boolean
}

// Missões do modo guiado "Aprenda na Prática"
const missions = ref<Mission[]>([
  { id: 'start', title: '1. Iniciar Timer', desc: 'Digite sleepoff 10s --dry-run ou escolha no menu', completed: false },
  { id: 'pause', title: '2. Pausar / Retomar', desc: 'Pressione [p] ou [Espaço] durante a contagem', completed: false },
  { id: 'add_time', title: '3. Somar +5 min', desc: 'Pressione [+] para estender o temporizador', completed: false },
  { id: 'panic', title: '4. Abortar Alarme', desc: 'Pressione [c] ou [Esc] durante o aviso de pânico', completed: false }
])

const isMuted = ref(false)
const mode = ref<Mode>('shell')
const tuiState = ref<TuiState>('SPLASH')
const isDryRun = ref(false)

// Shell state
const shellInput = ref('')
const shellInputRef = ref<HTMLInputElement | null>(null)
const terminalContainerRef = ref<HTMLDivElement | null>(null)
const terminalBodyRef = ref<HTMLDivElement | null>(null)
const isTerminalFocused = ref(false)

const shellHistory = ref<Array<{ cmd?: string; output?: string; type?: 'input' | 'output' | 'error' }>>([
  {
    type: 'output',
    output: 'sleepoff CLI v1.0.0 (Go 1.25+, CGO_ENABLED=0, 100% Nativo)\nLinux 6.11 x86_64 • TUI Charm Bubble Tea • Single-Instance Mutex\nDigite `sleepoff` para o menu TUI, `10s -d` para teste rápido, ou `--help`.'
  }
])
const historyIndex = ref(-1)
const commandHistory = ref<string[]>([])

// TUI state
const selectedMenuIdx = ref(0)
const customMinutes = ref('45')
const totalSeconds = ref(900)
const remainingSeconds = ref(900)
const isPaused = ref(false)
const eventsLog = ref<string[]>([])
const panicCountdown = ref(10)
const cancelPending = ref(false)
let cancelPendingTimeout: any = null

// Start and finish schedule times (faithful to view.go)
const startTimeStr = ref('')
const finishTimeStr = ref('')

function updateScheduleTimes(secs: number) {
  const now = new Date()
  startTimeStr.value = now.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
  const finish = new Date(now.getTime() + secs * 1000)
  finishTimeStr.value = finish.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
}

// Timers
let runningInterval: any = null
let panicInterval: any = null
let splashTimeout: any = null

// Audio Synthesis
let audioCtx: AudioContext | null = null
const isAudioPlaying = ref(false)
let audioVisualizerTimer: any = null

function triggerAudioVisualizer(duration = 0.15) {
  isAudioPlaying.value = true
  clearTimeout(audioVisualizerTimer)
  audioVisualizerTimer = setTimeout(() => {
    isAudioPlaying.value = false
  }, duration * 1000 + 50)
}

function getAudioContext(): AudioContext | null {
  if (!import.meta.client) return null
  if (!audioCtx) {
    const AudioCtx = window.AudioContext || (window as any).webkitAudioContext
    if (AudioCtx) {
      audioCtx = new AudioCtx()
    }
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
    // Silencioso se navegador restringir
  }
}

// Áudio Oficial de Contagem Regressiva (countdown.wav)
let countdownAudioTemplate: HTMLAudioElement | null = null

function preloadCountdownSound() {
  if (!import.meta.client) return
  try {
    const config = useRuntimeConfig()
    const soundUrl = `${config.app.baseURL}sounds/countdown.wav`
    if (!countdownAudioTemplate) {
      countdownAudioTemplate = new Audio(soundUrl)
      countdownAudioTemplate.volume = 0.85
      countdownAudioTemplate.preload = 'auto'
    }
  } catch {}
}

function playOfficialCountdownSound() {
  if (isMuted.value || !import.meta.client) return
  try {
    preloadCountdownSound()
    if (!countdownAudioTemplate) return
    const audioClone = countdownAudioTemplate.cloneNode() as HTMLAudioElement
    audioClone.volume = 0.85
    const playPromise = audioClone.play()
    if (playPromise !== undefined) {
      playPromise
        .then(() => {
          triggerAudioVisualizer(0.2)
        })
        .catch(() => {
          const freq = 600 + (10 - panicCountdown.value) * 60
          playBeep(freq, 0.15, 'sawtooth')
        })
    }
  } catch {
    const freq = 600 + (10 - panicCountdown.value) * 60
    playBeep(freq, 0.15, 'sawtooth')
  }
}

function toggleMute() {
  isMuted.value = !isMuted.value
  if (!isMuted.value) {
    getAudioContext()
    playBeep(880, 0.08)
  }
}

// Opções de Menu idênticas a config.go
const menuOptions = [
  { title: '15 min', desc: 'Tempo rápido', minutes: 15 },
  { title: '30 min', desc: 'Meia hora', minutes: 30 },
  { title: '45 min', desc: 'Episódio de série', minutes: 45 },
  { title: '1 hora', desc: 'Sessão completa', minutes: 60 },
  { title: '1h 30min', desc: 'Filme curto', minutes: 90 },
  { title: '2 horas', desc: 'Filme longo', minutes: 120 },
  { title: 'Personalizado', desc: 'Definir minutos', minutes: -1 }
]

function addLogEvent(msg: string) {
  const time = new Date().toTimeString().split(' ')[0]
  eventsLog.value.unshift(`[${time}] ${msg}`)
  if (eventsLog.value.length > 5) {
    eventsLog.value.pop()
  }
}

function formatRemaining(sec: number): string {
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  const s = sec % 60
  if (h > 0) {
    return `${h}h ${String(m).padStart(2, '0')}m ${String(s).padStart(2, '0')}s`
  }
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

function formatPanicTime(sec: number): string {
  return `00:${String(Math.max(0, sec)).padStart(2, '0')}`
}

async function fireConfettiCelebration() {
  if (import.meta.client) {
    try {
      const confettiModule = await import('canvas-confetti')
      const confetti = confettiModule.default || confettiModule
      confetti({
        particleCount: 120,
        spread: 80,
        origin: { y: 0.6 },
        colors: ['#7c3aed', '#06b6d4', '#10b981', '#f59e0b', '#ec4899']
      })
    } catch {}
  }
}

function triggerMission(id: string) {
  const m = missions.value.find(item => item.id === id)
  if (m && !m.completed) {
    m.completed = true
    playBeep(1200, 0.15)
    const allCompleted = missions.value.every(item => item.completed)
    if (allCompleted) {
      fireConfettiCelebration()
    }
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (terminalBodyRef.value) {
      terminalBodyRef.value.scrollTop = terminalBodyRef.value.scrollHeight
    }
  })
}

function focusTerminal() {
  isTerminalFocused.value = true
  getAudioContext()
  preloadCountdownSound()
  if (mode.value === 'shell') {
    shellInputRef.value?.focus()
  }
}

// Navegação entre abas preservando processos ativos
function switchTab(target: Mode) {
  mode.value = target
  if (target === 'shell') {
    nextTick(() => {
      shellInputRef.value?.focus()
      scrollToBottom()
    })
  } else {
    // Se o TUI nunca foi iniciado ou estava parado, abre o menu
    if (tuiState.value === 'SPLASH' || tuiState.value === 'DONE_TEST' || tuiState.value === 'CANCELLED') {
      startTui(-1, false)
    }
  }
}

// Transição TUI
function startTui(seconds: number, dry = false) {
  mode.value = 'tui'
  isDryRun.value = dry
  tuiState.value = 'SPLASH'
  cancelPending.value = false
  playBeep(600, 0.1)

  clearTimeout(splashTimeout)
  splashTimeout = setTimeout(() => {
    if (seconds > 0) {
      launchTimer(seconds)
    } else {
      tuiState.value = 'MENU'
      selectedMenuIdx.value = 0
    }
  }, 1200)
}

function launchTimer(seconds: number) {
  totalSeconds.value = seconds
  remainingSeconds.value = seconds
  isPaused.value = false
  cancelPending.value = false
  eventsLog.value = []
  updateScheduleTimes(seconds)

  const formatted = formatRemaining(seconds)
  addLogEvent(`Timer iniciado: ${formatted}${isDryRun.value ? ' (DRY-RUN)' : ''}`)
  triggerMission('start')
  playBeep(900, 0.1)

  // Se o tempo solicitado for 10s ou menos, entra direto na contagem de pânico!
  if (seconds <= 10) {
    startPanic(seconds)
    return
  }

  tuiState.value = 'RUNNING'
  clearInterval(runningInterval)
  runningInterval = setInterval(() => {
    if (isPaused.value) return
    if (remainingSeconds.value > 10) {
      remainingSeconds.value--
    } else if (remainingSeconds.value > 0) {
      remainingSeconds.value--
      startPanic(10)
    } else {
      finishDone()
    }
  }, 1000)
}

function startPanic(startSeconds = 10) {
  clearInterval(runningInterval)
  tuiState.value = 'CONFIRMATION'
  panicCountdown.value = startSeconds
  cancelPending.value = false
  playOfficialCountdownSound()

  clearInterval(panicInterval)
  panicInterval = setInterval(() => {
    if (panicCountdown.value > 1) {
      panicCountdown.value--
      playOfficialCountdownSound()
    } else {
      clearInterval(panicInterval)
      finishDone()
    }
  }, 1000)
}

function finishDone() {
  clearInterval(runningInterval)
  clearInterval(panicInterval)
  tuiState.value = 'DONE_TEST'
  addLogEvent('Desligamento simulado concluído!')
  playBeep(400, 0.4, 'triangle')
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
  if (remainingSeconds.value > 300) {
    remainingSeconds.value -= 300
    totalSeconds.value = Math.max(remainingSeconds.value, totalSeconds.value - 300)
    updateScheduleTimes(remainingSeconds.value)
    addLogEvent('-5 min subtraídos (-300s)')
    playBeep(600, 0.08)
  }
}

// Confirmação dupla de cancelamento durante o timer (fiel ao m.CancelPending do Go)
function handleCancel() {
  if (mode.value !== 'tui') return

  // No modo pânico de 10s: cancelamento IMEDIATO no primeiro aperto (fiel ao updateConfirmation do Go)
  if (tuiState.value === 'CONFIRMATION') {
    clearInterval(panicInterval)
    clearInterval(runningInterval)
    tuiState.value = 'CANCELLED'
    triggerMission('panic')
    addLogEvent('Cancelado durante o pânico!')
    playBeep(350, 0.25)
    return
  }

  // Durante o timer regular: requer confirmação dupla em 2.5 segundos
  if (tuiState.value === 'RUNNING') {
    if (cancelPending.value) {
      // Segundo toque -> cancela de verdade
      clearTimeout(cancelPendingTimeout)
      cancelPending.value = false
      clearInterval(runningInterval)
      clearInterval(panicInterval)
      tuiState.value = 'CANCELLED'
      triggerMission('panic')
      addLogEvent('Operação cancelada pelo usuário')
      playBeep(350, 0.25)
    } else {
      // Primeiro toque -> ativa aviso visual
      cancelPending.value = true
      playBeep(700, 0.08)
      clearTimeout(cancelPendingTimeout)
      cancelPendingTimeout = setTimeout(() => {
        cancelPending.value = false
      }, 2500)
    }
  }
}

function exitTuiToShell() {
  clearInterval(runningInterval)
  clearInterval(panicInterval)
  clearTimeout(splashTimeout)
  clearTimeout(cancelPendingTimeout)
  cancelPending.value = false
  mode.value = 'shell'
  tuiState.value = 'SPLASH'
  shellHistory.value.push({
    type: 'output',
    output: '[sleepoff finalizado — Retornado ao terminal Bash]'
  })
  scrollToBottom()
  nextTick(() => {
    shellInputRef.value?.focus()
  })
}

// Parser de duração composto fiel ao time.ParseDuration do Go
function parseDurationString(raw: string): number | null {
  const s = raw.trim().toLowerCase()
  // Caso puro número: minutos (ex: "30" -> 1800s)
  if (/^\d+$/.test(s)) {
    const mins = parseInt(s, 10)
    return mins > 0 ? mins * 60 : null
  }
  // Caso composto: ex: 1h30m, 2h, 45m, 90s, 10s, 1h15m30s
  const durationRegex = /^(?:(\d+(?:\.\d+)?)h)?(?:(\d+(?:\.\d+)?)m)?(?:(\d+(?:\.\d+)?)s)?$/
  const match = s.match(durationRegex)
  if (!match || match[0] === '') return null
  const hours = parseFloat(match[1] || '0')
  const minutes = parseFloat(match[2] || '0')
  const seconds = parseFloat(match[3] || '0')
  const total = Math.round(hours * 3600 + minutes * 60 + seconds)
  return total > 0 ? total : null
}

function handleShellCommand(rawCmd: string) {
  const cmd = rawCmd.trim()
  if (!cmd) return

  commandHistory.value.push(cmd)
  historyIndex.value = -1
  shellHistory.value.push({ type: 'input', cmd })

  const parts = cmd.split(/\s+/)
  let base = parts[0].toLowerCase()

  // Normalização de invocação: ./sleepoff, .\sleepoff, sleepoff.exe -> sleepoff
  if (base === './sleepoff' || base === '.\\sleepoff' || base === 'sleepoff.exe' || base === './sleepoff.exe') {
    base = 'sleepoff'
  }

  if (base === 'clear' || base === 'cls') {
    shellHistory.value = []
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'help') {
    shellHistory.value.push({
      type: 'output',
      output: `Comandos disponíveis no emulador:
  sleepoff                   Abre menu TUI interativo
  sleepoff <tempo>           Inicia timer direto (ex: 30m, 1h30m, 10s)
  sleepoff --dry-run         Simula execução sem desligar
  sleepoff --version         Exibe versão do binário
  ls                         Lista arquivos do workspace
  pwd                        Diretório corrente
  whoami                     Usuário autenticado
  uptime                     Tempo de atividade do host
  date                       Exibe data/hora do sistema
  uname -a                   Informações do kernel e arquitetura
  clear                      Limpa a tela do terminal`
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'ls' || base === 'dir') {
    shellHistory.value.push({
      type: 'output',
      output: `go.mod  go.sum  main.go  main_test.go  README.md  ARCHITECTURE.md
internal/  packaging/  scripts/  site/  sleepoff*`
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'pwd') {
    shellHistory.value.push({
      type: 'output',
      output: '/home/user/sleepoff'
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'whoami') {
    shellHistory.value.push({
      type: 'output',
      output: 'user'
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'uptime') {
    shellHistory.value.push({
      type: 'output',
      output: ' 04:20:00 up 42 days, 13:37,  1 user,  load average: 0.12, 0.08, 0.05'
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'cat') {
    const target = parts[1] || ''
    if (target.toLowerCase().includes('readme')) {
      shellHistory.value.push({
        type: 'output',
        output: `# sleepoff
Sleep Timer Inteligente para Windows e Multiplataforma.
Construído em Go 1.25+ com Bubble Tea, áudio embedado e zero bloat.`
      })
    } else {
      shellHistory.value.push({
        type: 'output',
        output: `cat: ${target || 'arquivo'}: Arquivo protegido ou emulador estático.`
      })
    }
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'date') {
    shellHistory.value.push({
      type: 'output',
      output: new Date().toString()
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base === 'uname') {
    shellHistory.value.push({
      type: 'output',
      output: 'Linux sleepoff-node 6.11.0-arch1-1-x86_64 SMP PREEMPT_DYNAMIC GNU/Linux'
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (base !== 'sleepoff') {
    // Atalho direto: aceita durações como "10s", "30m", "15m -d" sem precisar digitar "sleepoff" primeiro
    const directDuration = parseDurationString(base)
    if (directDuration !== null) {
      const isDry = parts.slice(1).includes('--dry-run') || parts.slice(1).includes('-d')
      shellInput.value = ''
      startTui(directDuration, isDry)
      return
    }

    shellHistory.value.push({
      type: 'error',
      output: `Comando não reconhecido: "${base}". Digite "sleepoff", "help" ou selecione os atalhos abaixo.`
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  const args = parts.slice(1)
  const isDry = args.includes('--dry-run') || args.includes('-d')
  const cleanArgs = args.filter(a => a !== '--dry-run' && a !== '-d')

  if (cleanArgs.includes('--help') || cleanArgs.includes('-h')) {
    shellHistory.value.push({
      type: 'output',
      output: `sleepoff v1.0.0 — Timer de desligamento automático
Comandos:
  sleepoff                   Abre menu interativo
  sleepoff <tempo>           Inicia timer direto (ex: 15m, 1h30m, 10s)
  sleepoff --dry-run, -d     Modo teste (não desliga de verdade)
  sleepoff --version, -v     Exibe versão do binário
  clear                      Limpa a tela`
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (cleanArgs.includes('--version') || cleanArgs.includes('-v')) {
    shellHistory.value.push({
      type: 'output',
      output: 'sleepoff version 1.0.0 (Go 1.25.6, CGO_ENABLED=0, multiplatform)'
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  // Se já há um timer rodando no Tab 2: single-instance check!
  if (tuiState.value === 'RUNNING' || tuiState.value === 'CONFIRMATION') {
    shellHistory.value.push({
      type: 'error',
      output: `sleepoff: outro timer já está em execução no Tab 2 (${formatRemaining(remainingSeconds.value)} restantes). Single-Instance lock ativo!`
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  if (cleanArgs.length === 0) {
    shellInput.value = ''
    startTui(-1, isDry)
    return
  }

  // Parse de duração composto
  const timeArg = cleanArgs[0]
  const parsedSecs = parseDurationString(timeArg)

  if (!parsedSecs) {
    shellHistory.value.push({
      type: 'error',
      output: `Formato de duração inválido: "${timeArg}". Formatos aceitos: 30, 30m, 90s, 10s, 1h, 1h30m.`
    })
    shellInput.value = ''
    scrollToBottom()
    return
  }

  shellInput.value = ''
  startTui(parsedSecs, isDry)
}

function runQuickCommand(cmd: string) {
  focusTerminal()
  shellInput.value = cmd
  handleShellCommand(cmd)
}

// Atalhos globais e navegação de teclado focada
function onKeyDown(e: KeyboardEvent) {
  // Ignora teclas se o usuário está digitando fora do terminal
  const activeEl = document.activeElement
  const isOuterInput = activeEl && (activeEl.tagName === 'INPUT' || activeEl.tagName === 'TEXTAREA') && activeEl !== shellInputRef.value
  if (isOuterInput) return

  // Se o foco não está no terminal e a tecla não é de controle, não captura
  if (!isTerminalFocused.value) return

  if (mode.value === 'shell') {
    if (e.key === 'Tab') {
      e.preventDefault()
      const val = shellInput.value.trim()
      const candidates = ['sleepoff', 'sleepoff 10s --dry-run', 'sleepoff 15m --dry-run', 'sleepoff 30m', 'sleepoff 1h30m', 'sleepoff --help', 'clear']
      const match = candidates.find(c => c.startsWith(val))
      if (match) {
        shellInput.value = match
      }
      return
    }

    if (e.ctrlKey && e.key.toLowerCase() === 'l') {
      e.preventDefault()
      shellHistory.value = []
      return
    }

    if (e.key === 'ArrowUp') {
      if (commandHistory.value.length > 0) {
        if (historyIndex.value === -1) {
          historyIndex.value = commandHistory.value.length - 1
        } else if (historyIndex.value > 0) {
          historyIndex.value--
        }
        shellInput.value = commandHistory.value[historyIndex.value] || ''
      }
      e.preventDefault()
    } else if (e.key === 'ArrowDown') {
      if (historyIndex.value !== -1) {
        if (historyIndex.value < commandHistory.value.length - 1) {
          historyIndex.value++
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

  // Teclado no modo TUI
  if (mode.value === 'tui') {
    if (tuiState.value === 'MENU') {
      if (e.key === 'ArrowUp' || e.key === 'k') {
        selectedMenuIdx.value = (selectedMenuIdx.value - 1 + menuOptions.length) % menuOptions.length
        playBeep(750, 0.05)
        e.preventDefault()
      } else if (e.key === 'ArrowDown' || e.key === 'j') {
        selectedMenuIdx.value = (selectedMenuIdx.value + 1) % menuOptions.length
        playBeep(750, 0.05)
        e.preventDefault()
      } else if (e.key === 'Enter') {
        const sel = menuOptions[selectedMenuIdx.value]
        if (sel.minutes === -1) {
          tuiState.value = 'CUSTOM_INPUT'
        } else {
          launchTimer(sel.minutes * 60)
        }
        e.preventDefault()
      } else if (e.key === 'q' || e.key === 'Escape') {
        exitTuiToShell()
        e.preventDefault()
      }
    } else if (tuiState.value === 'CUSTOM_INPUT') {
      if (e.key === 'Enter') {
        const val = parseInt(customMinutes.value)
        if (!isNaN(val) && val > 0) {
          launchTimer(val * 60)
        }
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
      } else if (e.key === '-') {
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
}

function handleDocumentClick(e: MouseEvent) {
  if (terminalContainerRef.value && terminalContainerRef.value.contains(e.target as Node)) {
    isTerminalFocused.value = true
  } else {
    isTerminalFocused.value = false
  }
}

// Interação guiada pelas Missões "Aprenda na Prática"
function clickMission(mission: Mission) {
  focusTerminal()
  if (mission.id === 'start') {
    if (tuiState.value === 'RUNNING' || tuiState.value === 'CONFIRMATION') {
      switchTab('tui')
    } else {
      runQuickCommand('sleepoff 10s --dry-run')
    }
  } else if (mission.id === 'pause') {
    if (mode.value !== 'tui') switchTab('tui')
    if (tuiState.value !== 'RUNNING') {
      launchTimer(900)
      setTimeout(togglePause, 200)
    } else {
      togglePause()
    }
  } else if (mission.id === 'add_time') {
    if (mode.value !== 'tui') switchTab('tui')
    if (tuiState.value !== 'RUNNING') {
      launchTimer(900)
      setTimeout(addFiveMinutes, 200)
    } else {
      addFiveMinutes()
    }
  } else if (mission.id === 'panic') {
    if (mode.value !== 'tui') switchTab('tui')
    if (tuiState.value !== 'RUNNING' && tuiState.value !== 'CONFIRMATION') {
      startPanic(10)
    } else {
      handleCancel()
    }
  }
}

onMounted(() => {
  if (import.meta.client) {
    window.addEventListener('keydown', onKeyDown)
    document.addEventListener('click', handleDocumentClick)

    // Unlocking de áudio no primeiro gesto do usuário em qualquer lugar da tela
    const unlockAudio = () => {
      getAudioContext()
      preloadCountdownSound()
      window.removeEventListener('pointerdown', unlockAudio)
      window.removeEventListener('keydown', unlockAudio)
    }
    window.addEventListener('pointerdown', unlockAudio, { passive: true })
    window.addEventListener('keydown', unlockAudio, { passive: true })
  }
})

onUnmounted(() => {
  if (import.meta.client) {
    window.removeEventListener('keydown', onKeyDown)
    document.removeEventListener('click', handleDocumentClick)
  }
  clearInterval(runningInterval)
  clearInterval(panicInterval)
  clearTimeout(splashTimeout)
  clearTimeout(cancelPendingTimeout)
})

const progressPercent = computed(() => {
  if (totalSeconds.value <= 0) return 0
  const elapsed = totalSeconds.value - remainingSeconds.value
  return Math.min(100, Math.max(0, Math.round((elapsed / totalSeconds.value) * 100)))
})

// Estatísticas fieis ao showCancelSummary() de main.go
const cancelSummaryStats = computed(() => {
  const elapsed = Math.max(0, totalSeconds.value - remainingSeconds.value)
  const totalMins = Math.floor(totalSeconds.value / 60)
  const elapsedMins = Math.floor(elapsed / 60)
  const elapsedSecs = elapsed % 60
  const pct = totalSeconds.value > 0 ? Math.min(100, Math.round((elapsed / totalSeconds.value) * 100)) : 0

  let emoji = '😅'
  let msg = 'Desistiu rápido hein?'
  if (pct >= 80) {
    emoji = '🏆'
    msg = 'QUASE LÁ! Tão perto do final!'
  } else if (pct >= 50) {
    emoji = '💪'
    msg = 'PASSOU DA METADE! Respeitável.'
  } else if (pct >= 20) {
    emoji = '🤔'
    msg = 'Bem... pelo menos tentou.'
  }

  const barLen = 30
  const filled = Math.min(barLen, Math.round((pct * barLen) / 100))
  const asciiBar = '█'.repeat(filled) + '░'.repeat(Math.max(0, barLen - filled))

  return {
    elapsedMins,
    elapsedSecs,
    totalMins,
    pct,
    emoji,
    msg,
    asciiBar
  }
})
</script>

<template>
  <div
    ref="terminalContainerRef"
    @click="focusTerminal"
    class="relative w-full max-w-5xl mx-auto min-h-[560px] sm:min-h-[580px] rounded-2xl border border-slate-700/80 bg-slate-950/95 backdrop-blur-xl shadow-2xl overflow-hidden flex flex-col transition-all duration-300 ring-1 ring-white/5"
    :class="{ 'ring-1 ring-violet-500/60 shadow-violet-950/50': isTerminalFocused }"
  >

    <!-- Header Terminal Bar (Estilo macOS / Alacritty com Indicadores de Alta Precisão) -->
    <div class="relative z-10 flex items-center justify-between px-3 sm:px-4 py-2.5 bg-slate-900/90 border-b border-slate-800/90 select-none min-w-0">
      <div class="flex items-center space-x-2 sm:space-x-3 min-w-0">
        <!-- macOS Window Action Dots com acabamento polido e interatividade -->
        <div class="flex items-center space-x-2 shrink-0">
          <button
            @click.stop="runQuickCommand('clear')"
            title="Fechar / Limpar terminal (clear)"
            class="w-3 h-3 rounded-full bg-[#ff5f56] border border-[#e0443e]/80 shadow-sm inline-block hover:brightness-110 active:scale-90 transition-all cursor-pointer"
            aria-label="Limpar terminal"
          />
          <button
            @click.stop="exitTuiToShell"
            title="Minimizar TUI / Voltar ao Bash"
            class="w-3 h-3 rounded-full bg-[#ffbd2e] border border-[#dea123]/80 shadow-sm inline-block hover:brightness-110 active:scale-90 transition-all cursor-pointer"
            aria-label="Voltar ao Bash"
          />
          <button
            @click.stop="switchTab(mode === 'shell' ? 'tui' : 'shell')"
            title="Alternar entre Bash e TUI"
            class="w-3 h-3 rounded-full bg-[#27c93f] border border-[#1aab29]/80 shadow-sm inline-block hover:brightness-110 active:scale-90 transition-all cursor-pointer"
            aria-label="Alternar abas"
          />
        </div>

        <!-- Terminal Tabs com status dinâmico e atalhos -->
        <div class="flex items-center space-x-1.5 font-mono text-xs ml-1 sm:ml-2">
          <button
            @click.stop="switchTab('shell')"
            aria-label="1: bash — terminal bash interativo"
            :class="mode === 'shell' ? 'bg-slate-800 text-cyan-300 font-bold border-b-2 border-cyan-400 shadow-sm' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'"
            class="px-2.5 py-1.5 rounded-t-lg transition-all flex items-center space-x-1.5 cursor-pointer shrink-0 min-h-[36px]"
          >
            <Icon name="heroicons:command-line" class="w-3.5 h-3.5 text-cyan-400" />
            <span>1: bash</span>
          </button>
          <button
            @click.stop="switchTab('tui')"
            aria-label="2: sleepoff — interface Bubble Tea TUI"
            :class="mode === 'tui' ? 'bg-slate-800 text-violet-300 font-bold border-b-2 border-violet-400 shadow-sm' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'"
            class="px-2.5 py-1.5 rounded-t-lg transition-all flex items-center space-x-1.5 cursor-pointer min-w-0 min-h-[36px]"
          >
            <Icon name="heroicons:cpu-chip" class="w-3.5 h-3.5 text-violet-400 shrink-0" />
            <span class="truncate">2: sleepoff</span>
            <span
              v-if="tuiState === 'RUNNING'"
              class="inline-flex items-center text-[10px] text-emerald-400 font-bold ml-1 px-1.5 py-0.5 bg-emerald-950/80 rounded border border-emerald-700/50 shrink-0"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse mr-1" />
              {{ formatRemaining(remainingSeconds) }}
            </span>
            <span
              v-else-if="tuiState === 'CONFIRMATION'"
              class="inline-flex items-center text-[10px] text-red-400 font-bold ml-1 px-1.5 py-0.5 bg-red-950/80 rounded border border-red-700/50 shrink-0"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-red-400 animate-ping mr-1" />
              {{ formatPanicTime(panicCountdown) }}
            </span>
          </button>
        </div>
      </div>

      <!-- Título Central Sutil -->
      <div class="hidden lg:flex items-center space-x-1.5 text-[11px] font-mono text-slate-500 select-none">
        <Icon name="heroicons:command-line" class="w-3.5 h-3.5 text-slate-500" />
        <span>sleepoff • Charm Bubble Tea v1.3.10</span>
      </div>

      <!-- Controles de Áudio & Sair -->
      <div class="flex items-center space-x-1.5 sm:space-x-3 text-xs shrink-0">
        <!-- Áudio Visualizer SVG animado reativo aos bipes -->
        <div
          class="hidden sm:flex items-center space-x-1 px-2.5 py-1 rounded-lg bg-slate-800/70 border border-slate-700/60"
          :title="isMuted ? 'Áudio desativado' : (isAudioPlaying ? 'Bipe emitindo som' : 'Aguardando som')"
        >
          <span
            class="w-1 rounded-full transition-all duration-100"
            :class="[
              isMuted ? 'h-1 bg-slate-600' : (isAudioPlaying ? 'h-4 bg-cyan-400 animate-pulse' : 'h-1.5 bg-slate-500')
            ]"
          />
          <span
            class="w-1 rounded-full transition-all duration-100"
            :class="[
              isMuted ? 'h-1 bg-slate-600' : (isAudioPlaying ? 'h-5 bg-violet-400 animate-pulse' : 'h-2.5 bg-slate-500')
            ]"
          />
          <span
            class="w-1 rounded-full transition-all duration-100"
            :class="[
              isMuted ? 'h-1 bg-slate-600' : (isAudioPlaying ? 'h-3.5 bg-emerald-400 animate-pulse' : 'h-1.5 bg-slate-500')
            ]"
          />
          <span
            class="w-1 rounded-full transition-all duration-100"
            :class="[
              isMuted ? 'h-1 bg-slate-600' : (isAudioPlaying ? 'h-4.5 bg-cyan-300 animate-pulse' : 'h-2 bg-slate-500')
            ]"
          />
          <span class="ml-1.5 text-[10px] font-mono tracking-wider" :class="isAudioPlaying ? 'text-cyan-300 font-bold' : 'text-slate-400'">
            {{ isMuted ? 'MUTE' : (isAudioPlaying ? 'BEEP' : 'AUDIO') }}
          </span>
        </div>

        <button
          @click.stop="toggleMute"
          :aria-label="isMuted ? 'Ativar efeitos sonoros do terminal' : 'Mutar efeitos sonoros do terminal'"
          class="flex items-center space-x-1.5 px-2.5 py-1.5 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 transition-colors cursor-pointer min-h-[36px] border border-slate-700/50"
          :title="isMuted ? 'Ativar Efeitos Sonoros' : 'Mutar Efeitos Sonoros'"
        >
          <Icon :name="isMuted ? 'heroicons:speaker-x-mark' : 'heroicons:speaker-wave'" class="w-4 h-4 text-cyan-400" />
          <span class="hidden sm:inline text-xs font-mono">{{ isMuted ? 'Mudo' : 'Som' }}</span>
        </button>

        <button
          v-if="mode === 'tui'"
          @click.stop="exitTuiToShell"
          aria-label="Sair da TUI e retornar ao shell"
          class="px-2.5 py-1.5 rounded-lg bg-rose-500/20 text-rose-300 hover:bg-rose-500/30 border border-rose-500/30 font-mono transition-colors cursor-pointer text-xs min-h-[36px]"
        >
          <span>Sair (q)</span>
        </button>
      </div>
    </div>

    <!-- Terminal Body / Visualizador -->
    <div
      ref="terminalBodyRef"
      @click="focusTerminal"
      class="p-4 sm:p-6 min-h-[260px] sm:min-h-[400px] max-h-[520px] overflow-y-auto font-mono text-sm leading-relaxed text-slate-200 select-text flex flex-col justify-between cursor-text"
    >
      
      <!-- 1. MODO SHELL INTERATIVO -->
      <div v-if="mode === 'shell'" class="space-y-3 flex-1 flex flex-col justify-between">
        <div class="space-y-2">
          <!-- History items -->
          <div v-for="(item, idx) in shellHistory" :key="idx" class="space-y-1">
            <div v-if="item.type === 'input'" class="text-cyan-400 flex items-center space-x-1.5 font-mono text-xs sm:text-sm">
              <span class="text-emerald-400 font-bold">user@sleepoff</span>
              <span class="text-slate-500">:</span>
              <span class="text-cyan-400 font-bold">~</span>
              <span class="text-violet-400 font-bold mr-1">$</span>
              <span class="text-slate-100 font-semibold">{{ item.cmd }}</span>
            </div>
            <pre v-else :class="item.type === 'error' ? 'text-rose-400 font-mono text-xs' : 'text-slate-300 font-mono text-xs whitespace-pre-wrap leading-relaxed'">{{ item.output }}</pre>
          </div>

          <!-- Active Prompt Line: Authentic Terminal Prompt with Blinking Block Cursor -->
          <div class="flex items-center flex-wrap font-mono text-xs sm:text-sm pt-1 cursor-text" @click="focusTerminal">
            <span class="text-emerald-400 font-bold">user@sleepoff</span>
            <span class="text-slate-500">:</span>
            <span class="text-cyan-400 font-bold">~</span>
            <span class="text-violet-400 font-bold mr-2">$</span>
            <span class="text-white font-mono whitespace-pre">{{ shellInput }}</span>
            <span
              class="inline-block w-2.5 h-4.5 bg-cyan-400 align-middle ml-0.5"
              :class="isTerminalFocused ? 'animate-pulse' : 'opacity-40'"
            />
            <input
              ref="shellInputRef"
              id="terminal-shell-input"
              aria-label="Linha de comando do terminal sleepoff"
              v-model="shellInput"
              type="text"
              autofocus
              @focus="isTerminalFocused = true"
              @blur="isTerminalFocused = false"
              @keydown.enter.prevent="handleShellCommand(shellInput)"
              class="absolute opacity-0 w-0 h-0 pointer-events-none"
            />
          </div>
        </div>

        <!-- Shell statusline hint -->
        <div class="mt-auto pt-3 border-t border-slate-900/90 flex flex-wrap items-center justify-between gap-2 text-xs font-mono select-none">
          <div class="flex items-center space-x-2 text-slate-400">
            <span class="px-1.5 py-0.5 rounded bg-violet-950/80 text-violet-300 font-bold text-[10px] border border-violet-800/40">BASH</span>
            <span class="text-slate-400 text-[11px]">Dica: digite <code class="text-cyan-300">sleepoff 10s --dry-run</code> ou use os atalhos abaixo:</span>
          </div>
        </div>
      </div>

      <!-- 2. MODO TUI (BUBBLE TEA FIDELIDADE TOTAL AO GO VIEW.GO) -->
      <div v-else class="flex flex-col items-center justify-center flex-1 my-auto text-center w-full">
        
        <!-- Estado: SPLASH -->
        <div v-if="tuiState === 'SPLASH'" class="space-y-4 animate-fade-in py-6 w-full max-w-lg mx-auto">
          <div class="flex items-center justify-between border-b border-slate-800 pb-2 font-mono text-xs">
            <span class="text-violet-300 font-bold">sleepoff <span class="text-slate-500">v1.0.0</span></span>
            <span v-if="isDryRun" class="px-2 py-0.5 rounded text-[10px] bg-cyan-950/80 text-cyan-300 border border-cyan-800/60 font-bold">DRY-RUN</span>
          </div>
          <pre class="text-[10px] sm:text-xs font-mono font-bold text-violet-400 tracking-tighter select-none leading-tight inline-block text-left">
 ███████╗██╗     ███████╗███████╗██████╗  ██████╗ 
 ██╔════╝██║     ██╔════╝██╔════╝██╔══██╗██╔═══██╗
 ███████╗██║     █████╗  █████╗  ██████╔╝██║   ██║
 ╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ ██║   ██║
 ███████║███████╗███████╗███████╗██║     ╚██████╔╝
 ╚══════╝╚══════╝╚══════╝╚══════╝╚═╝      ╚═════╝ </pre>
          <div class="text-cyan-400 font-semibold tracking-wide text-xs">
            Timer de desligamento automático • v1.0.0 (Go 1.25.6)
          </div>
          <div class="text-xs text-slate-500 font-mono animate-pulse">
            Iniciando máquina de estados TEA (The Elm Architecture)...
          </div>
        </div>

        <!-- Estado: MENU (Fiel ao viewMenu() do Go) -->
        <div v-else-if="tuiState === 'MENU'" class="w-full max-w-md mx-auto space-y-3.5 text-left py-2">
          <div class="flex items-center justify-between border-b border-slate-800 pb-2.5 font-mono text-xs">
            <div class="flex items-center space-x-2">
              <span class="text-violet-300 font-bold">sleepoff</span>
              <span class="text-slate-500 text-[10px]">v1.0.0</span>
              <span class="text-slate-600">──</span>
              <span class="text-slate-200 font-bold">SELECIONE O TEMPO</span>
            </div>
            <span v-if="isDryRun" class="px-2 py-0.5 rounded text-[10px] bg-cyan-950/80 text-cyan-300 border border-cyan-800/60 font-bold">
              DRY-RUN
            </span>
          </div>

          <div class="space-y-1.5 font-mono text-xs">
            <div
              v-for="(opt, idx) in menuOptions"
              :key="idx"
              @click="selectedMenuIdx = idx; opt.minutes === -1 ? tuiState = 'CUSTOM_INPUT' : launchTimer(opt.minutes * 60)"
              :class="[
                'p-2.5 rounded-xl border transition-all cursor-pointer flex items-center justify-between',
                selectedMenuIdx === idx
                  ? 'border-violet-500/80 bg-violet-950/40 text-cyan-300 font-bold shadow-[0_0_15px_rgba(124,58,237,0.25)] translate-x-1'
                  : 'border-slate-800/80 bg-slate-900/40 text-slate-400 hover:border-slate-700 hover:text-slate-200'
              ]"
            >
              <div class="flex items-center space-x-3">
                <span :class="selectedMenuIdx === idx ? 'text-cyan-400 font-bold' : 'text-slate-600'">
                  {{ selectedMenuIdx === idx ? '▶' : ' ' }}
                </span>
                <span>{{ opt.title }}</span>
              </div>
              <span class="text-[11px] text-slate-500">{{ opt.desc }}</span>
            </div>
          </div>

          <div class="border-t border-slate-800/80 pt-2.5 text-[11px] text-center text-slate-500 font-mono">
            <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300">↑/↓</kbd> navegar • <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300">Enter</kbd> confirmar • <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300">q</kbd> sair
          </div>
        </div>

        <!-- Estado: CUSTOM INPUT (Fiel ao viewCustomInput() do Go) -->
        <div v-else-if="tuiState === 'CUSTOM_INPUT'" class="w-full max-w-sm mx-auto space-y-4 py-6 text-center font-mono">
          <div class="flex items-center justify-between border-b border-slate-800 pb-2.5 text-xs">
            <div class="flex items-center space-x-2">
              <span class="text-violet-300 font-bold">sleepoff</span>
              <span class="text-slate-600">──</span>
              <span class="text-slate-200 font-bold">TEMPO PERSONALIZADO</span>
            </div>
            <span v-if="isDryRun" class="px-2 py-0.5 rounded text-[10px] bg-cyan-950/80 text-cyan-300 border border-cyan-800/60 font-bold">DRY-RUN</span>
          </div>

          <div class="text-xs text-slate-300">Quantos minutos você deseja antes de desligar?</div>
          <div class="flex items-center justify-center space-x-2 pt-2">
            <span class="text-cyan-400 font-bold">❯</span>
            <input
              id="tui-custom-minutes-input"
              aria-label="Minutos personalizados para contagem regressiva"
              v-model="customMinutes"
              type="number"
              min="1"
              max="1440"
              class="w-32 bg-slate-900 border border-violet-500 rounded-xl px-4 py-2 text-center text-xl font-bold text-cyan-300 outline-none shadow-[0_0_15px_rgba(124,58,237,0.2)]"
            />
            <span class="text-slate-400 font-bold">min</span>
          </div>
          <div class="text-[11px] text-slate-500">Valores aceitos: 1 a 1440 minutos</div>

          <div class="flex justify-center space-x-3 pt-3 border-t border-slate-800/80">
            <button
              @click="tuiState = 'MENU'"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs cursor-pointer border border-slate-700"
            >
              Voltar (Esc)
            </button>
            <button
              @click="launchTimer(parseInt(customMinutes) * 60)"
              class="px-5 py-2 rounded-xl bg-violet-600 hover:bg-violet-500 text-white font-bold text-xs shadow-lg shadow-violet-600/30 cursor-pointer"
            >
              Iniciar (Enter)
            </button>
          </div>
        </div>

        <!-- Estado: RUNNING (Fiel ao viewRunning() do Go com Claude Code TUI styling) -->
        <div v-else-if="tuiState === 'RUNNING'" class="w-full max-w-lg mx-auto space-y-4 py-2">
          <!-- Top bar padronizada -->
          <div class="flex items-center justify-between border-b border-slate-800 pb-2.5 font-mono text-xs">
            <div class="flex items-center space-x-2">
              <span class="text-violet-300 font-bold">sleepoff</span>
              <span class="text-slate-500 text-[10px]">v1.0.0</span>
              <span class="text-slate-600">──</span>
              <span class="text-slate-300 font-semibold">{{ formatRemaining(totalSeconds) }}</span>
            </div>
            <div class="flex items-center space-x-2">
              <span
                class="px-2 py-0.5 rounded text-[10px] font-bold flex items-center space-x-1"
                :class="isPaused ? 'bg-amber-950/80 text-amber-300 border border-amber-800/60' : 'bg-emerald-950/80 text-emerald-300 border border-emerald-800/60'"
              >
                <span class="w-1.5 h-1.5 rounded-full" :class="isPaused ? 'bg-amber-400' : 'bg-emerald-400 animate-pulse'"></span>
                <span>{{ isPaused ? '❚❚ PAUSADO' : '● RODANDO' }}</span>
              </span>
              <span v-if="isDryRun" class="px-2 py-0.5 rounded text-[10px] bg-cyan-950/80 text-cyan-300 border border-cyan-800/60 font-bold">
                DRY-RUN
              </span>
            </div>
          </div>

          <!-- Cronômetro Digital Centrado -->
          <div class="py-2 text-center font-mono">
            <div class="text-5xl sm:text-6xl font-black text-sky-400 tracking-wider drop-shadow-[0_0_25px_rgba(56,189,248,0.25)]">
              {{ formatRemaining(remainingSeconds) }}
            </div>
            <div class="text-xs text-slate-400 mt-2 flex items-center justify-center space-x-3">
              <span>Início: <strong class="text-slate-200">{{ startTimeStr }}</strong></span>
              <span class="text-slate-600">•</span>
              <span>Previsão: <strong class="text-cyan-300">{{ finishTimeStr }}</strong></span>
              <span class="text-slate-600">•</span>
              <span>Progresso: <strong class="text-violet-300">{{ progressPercent }}%</strong></span>
            </div>
          </div>

          <!-- Barra de Progresso com Gradiente -->
          <div class="w-full h-3 bg-slate-900 rounded-full overflow-hidden p-0.5 border border-slate-800">
            <div
              class="h-full rounded-full bg-gradient-to-r from-violet-600 via-indigo-500 to-cyan-400 transition-all duration-500"
              :style="{ width: `${progressPercent}%` }"
            ></div>
          </div>

          <!-- Alerta de cancelamento pendente -->
          <div
            v-if="cancelPending"
            class="p-2.5 rounded-xl bg-rose-950/90 border border-rose-500 text-rose-300 text-xs font-bold animate-pulse shadow-lg shadow-rose-950/50 text-center font-mono"
          >
            ⚠ Pressione <kbd class="px-1.5 py-0.5 rounded bg-rose-900 text-white font-mono">c</kbd> ou <kbd class="px-1.5 py-0.5 rounded bg-rose-900 text-white font-mono">Esc</kbd> novamente para confirmar o cancelamento!
          </div>

          <!-- Histórico de Eventos Timeline Compacta -->
          <div class="bg-slate-900/60 rounded-xl p-3 border border-slate-800/80 text-left text-xs font-mono space-y-1">
            <div class="text-slate-500 font-bold border-b border-slate-800/60 pb-1 text-[11px]">EVENTOS:</div>
            <div v-for="(ev, idx) in eventsLog" :key="idx" class="text-slate-300 truncate text-[11px] flex items-center space-x-2">
              <span class="text-slate-600">│</span>
              <span>{{ ev }}</span>
            </div>
          </div>
        </div>

        <!-- Estado: CONFIRMATION (Tela de Pânico dos Últimos 10s — Fiel ao viewPanic() do Go) -->
        <div v-else-if="tuiState === 'CONFIRMATION'" class="w-full max-w-lg mx-auto space-y-4 py-4 px-6 border-2 border-red-500/80 rounded-2xl bg-red-950/30 shadow-[0_0_50px_rgba(239,68,68,0.35)] animate-pulse">
          <div class="flex items-center justify-between border-b border-red-800/80 pb-2.5 font-mono text-xs">
            <span class="text-red-400 font-bold flex items-center space-x-1.5">
              <span>⚠</span>
              <span>DESLIGAMENTO IMINENTE</span>
            </span>
            <span v-if="isDryRun" class="px-2 py-0.5 rounded text-[10px] bg-red-900/60 text-red-200 border border-red-700/60 font-bold">
              DRY-RUN
            </span>
          </div>

          <div class="text-center py-2">
            <div class="text-red-400 font-mono text-xs sm:text-sm font-bold tracking-wider mb-2">
              &gt;&gt;&gt; DESLIGANDO EM {{ String(Math.max(0, panicCountdown)).padStart(2, '0') }} SEGUNDOS &lt;&lt;&lt;
            </div>
            <div class="text-6xl sm:text-7xl font-black text-red-500 font-mono tracking-widest my-2 drop-shadow-[0_0_20px_rgba(239,68,68,0.6)]">
              {{ formatPanicTime(panicCountdown) }}
            </div>
            <div class="text-xs text-slate-300 mt-2 font-mono">
              Salve seus arquivos. O sistema será desligado agora!
            </div>
          </div>

          <!-- Barra de Progresso Vermelha ASCII idêntica ao Go -->
          <div class="w-full h-3 bg-red-950/80 rounded-full overflow-hidden p-0.5 border border-red-800/80">
            <div
              class="h-full rounded-full bg-gradient-to-r from-red-600 to-amber-500 transition-all duration-300"
              :style="{ width: `${(Math.max(0, panicCountdown) / 10) * 100}%` }"
            ></div>
          </div>

          <div class="border-t border-red-800/80 pt-3">
            <button
              @click="handleCancel"
              class="w-full py-3.5 rounded-xl bg-red-600 hover:bg-red-500 text-white font-mono font-black text-xs sm:text-sm uppercase tracking-wider shadow-lg shadow-red-600/50 cursor-pointer active:scale-95 transition-all flex items-center justify-center space-x-2"
            >
              <kbd class="px-2 py-0.5 rounded bg-red-800 text-white font-bold text-xs border border-red-700">c</kbd>
              <span>ABORTAR DESLIGAMENTO IMEDIATAMENTE (ou Esc)</span>
            </button>
          </div>
        </div>

        <!-- Estado: DONE_TEST (Paridade Total com showFinalSuccess() de main.go) -->
        <div v-else-if="tuiState === 'DONE_TEST'" class="w-full max-w-xl mx-auto space-y-4 py-4 text-center">
          <pre class="text-[10px] sm:text-xs font-mono font-bold text-cyan-400 leading-tight select-none inline-block text-left">
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║      ████████╗██╗███╗   ███╗███████╗██████╗                   ║
   ║      ╚══██╔══╝██║████╗ ████║██╔════╝██╔══██╗                  ║
   ║         ██║   ██║██╔████╔██║█████╗  ██████╔╝                  ║
   ║         ██║   ██║██║╚██╔╝██║██╔══╝  ██╔══██╗                  ║
   ║         ██║   ██║██║ ╚═╝ ██║███████╗██║  ██║                  ║
   ║         ╚═╝   ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝                  ║
   ║                                                               ║
   ║                  ✓  FINALIZADO (TESTE)                        ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝</pre>
          <div class="text-emerald-400 font-bold text-sm font-mono">
            ✓ TESTE DE DESLIGAMENTO CONCLUÍDO COM SUCESSO!
          </div>
          <div class="space-y-1 text-xs text-slate-400 font-mono">
            <div>ℹ O computador <strong class="text-slate-200">NÃO</strong> foi desligado (modo simulado).</div>
            <div>ℹ No modo de produção, a syscall de energia desligaria seu sistema com segurança.</div>
          </div>
          <div class="pt-2 flex justify-center space-x-3">
            <button
              @click="exitTuiToShell"
              class="px-5 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-mono font-bold cursor-pointer transition-all border border-slate-700"
            >
              Voltar ao Shell (Enter)
            </button>
            <button
              @click="startTui(10, true)"
              class="px-5 py-2.5 rounded-xl bg-violet-600 hover:bg-violet-500 text-white text-xs font-mono font-bold cursor-pointer transition-all shadow-md shadow-violet-600/30"
            >
              ⚡ Testar Novamente (10s)
            </button>
          </div>
        </div>

        <!-- Estado: CANCELLED (Paridade Total com showCancelSummary() de main.go) -->
        <div v-else-if="tuiState === 'CANCELLED'" class="w-full max-w-xl mx-auto space-y-4 py-4 text-center">
          <pre class="text-[10px] sm:text-xs font-mono font-bold text-rose-500 leading-tight select-none inline-block text-left">
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║      ████████╗██╗███╗   ███╗███████╗██████╗                   ║
   ║      ╚══██╔══╝██║████╗ ████║██╔════╝██╔══██╗                  ║
   ║         ██║   ██║██╔████╔██║█████╗  ██████╔╝                  ║
   ║         ██║   ██║██║╚██╔╝██║██╔══╝  ██╔══██╗                  ║
   ║         ██║   ██║██║ ╚═╝ ██║███████╗██║  ██║                  ║
   ║         ╚═╝   ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝                  ║
   ║                                                               ║
   ║                     ⏹  CANCELADO                              ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝</pre>

          <!-- Barra de Progresso ASCII idêntica ao Go -->
          <div class="text-xs font-mono">
            <span class="text-slate-400">Progresso: </span>
            <span class="text-cyan-400 font-bold">[{{ cancelSummaryStats.asciiBar }}]</span>
            <span class="text-white font-bold ml-2">{{ cancelSummaryStats.pct }}%</span>
          </div>

          <!-- Box com estatísticas exatas do Go -->
          <div class="inline-block bg-slate-900/80 border border-slate-800 rounded-xl p-3 text-left font-mono text-xs text-slate-300 min-w-[280px]">
            <div class="flex justify-between py-0.5">
              <span class="text-slate-500">Tempo decorrido:</span>
              <span class="text-cyan-300 font-bold">{{ String(cancelSummaryStats.elapsedMins).padStart(2, '0') }}m {{ String(cancelSummaryStats.elapsedSecs).padStart(2, '0') }}s</span>
            </div>
            <div class="flex justify-between py-0.5">
              <span class="text-slate-500">Tempo total:</span>
              <span class="text-violet-300 font-bold">{{ String(cancelSummaryStats.totalMins).padStart(2, '0') }}m</span>
            </div>
          </div>

          <!-- Frase Motivacional idêntica ao Go -->
          <div class="text-sm font-bold text-violet-300 font-mono flex items-center justify-center space-x-2">
            <span>{{ cancelSummaryStats.emoji }}</span>
            <span>{{ cancelSummaryStats.msg }}</span>
          </div>

          <div class="pt-2 flex justify-center space-x-3">
            <button
              @click="exitTuiToShell"
              class="px-5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-mono font-bold cursor-pointer transition-all border border-slate-700"
            >
              Voltar ao Shell (Enter)
            </button>
            <button
              @click="startTui(10, true)"
              class="px-5 py-2 rounded-xl bg-violet-600 hover:bg-violet-500 text-white text-xs font-mono font-bold cursor-pointer transition-all shadow-md shadow-violet-600/30"
            >
              ⚡ Testar Novamente (10s)
            </button>
          </div>
        </div>

      </div>
    </div>

    <!-- Barra de Ações Rápidas & Atalhos (Contextual e 100% Livre de Botões Mortos) -->
    <div class="px-3 sm:px-4 py-2.5 bg-slate-900/95 border-t border-slate-800/80 flex flex-wrap items-center justify-between gap-2.5 text-xs select-none">
      
      <!-- Se estiver no Shell: Paleta de Iniciação Rápida -->
      <template v-if="mode === 'shell'">
        <div class="flex items-center space-x-2 text-slate-400 font-mono text-xs">
          <span class="text-slate-400 font-bold uppercase tracking-wider text-[10px]">Atalhos Rápidos:</span>
        </div>
        <div class="flex items-center flex-wrap gap-1.5 sm:gap-2">
          <button
            @click.stop="runQuickCommand('sleepoff 10s --dry-run')"
            class="px-3 py-1.5 rounded-lg bg-violet-600 hover:bg-violet-500 text-white font-bold transition-all font-mono text-xs shadow-md shadow-violet-600/30 cursor-pointer flex items-center space-x-1 active:scale-95"
            title="Teste rápido de 10 segundos com pânico e alarme"
          >
            <span>⚡ Testar 10s (Dry-Run)</span>
          </button>
          <button
            @click.stop="runQuickCommand('sleepoff')"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-cyan-300 font-mono border border-slate-700 transition-colors cursor-pointer"
            title="Abrir Menu Interativo TUI"
          >
            sleepoff (TUI)
          </button>
          <button
            @click.stop="runQuickCommand('sleepoff 15m --dry-run')"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-mono border border-slate-700 transition-colors cursor-pointer"
          >
            15m --dry-run
          </button>
          <button
            @click.stop="runQuickCommand('sleepoff 30m')"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-emerald-300 font-mono border border-slate-700 transition-colors cursor-pointer"
          >
            30m
          </button>
          <button
            @click.stop="runQuickCommand('sleepoff --help')"
            class="px-2 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-mono border border-slate-700 transition-colors cursor-pointer"
          >
            --help
          </button>
          <button
            @click.stop="runQuickCommand('clear')"
            class="px-2 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-mono border border-slate-700 transition-colors cursor-pointer"
          >
            clear
          </button>
        </div>
      </template>

      <!-- Se estiver no TUI Rodando: Controles Ativos do Timer -->
      <template v-else-if="tuiState === 'RUNNING'">
        <div class="flex items-center flex-wrap gap-1.5 sm:gap-2">
          <span class="text-slate-400 font-bold text-[10px] uppercase tracking-wider mr-1">Controles:</span>
          <button
            @click.stop="togglePause"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 transition-colors cursor-pointer border border-slate-700 font-mono text-xs"
          >
            <kbd class="font-bold text-violet-400">p</kbd> {{ isPaused ? 'Retomar' : 'Pausar' }}
          </button>
          <button
            @click.stop="addFiveMinutes"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 transition-colors cursor-pointer border border-slate-700 font-mono text-xs"
          >
            <kbd class="font-bold text-cyan-400">+</kbd> +5min
          </button>
          <button
            @click.stop="subtractFiveMinutes"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 transition-colors cursor-pointer border border-slate-700 font-mono text-xs"
          >
            <kbd class="font-bold text-cyan-400">-</kbd> -5min
          </button>
          <button
            @click.stop="handleCancel"
            :class="cancelPending ? 'bg-rose-600 text-white animate-pulse border-rose-400' : 'bg-slate-800 text-rose-300 hover:bg-slate-700 border-slate-700'"
            class="px-2.5 py-1.5 rounded-lg transition-colors cursor-pointer border font-mono text-xs"
          >
            <kbd class="font-bold" :class="cancelPending ? 'text-white' : 'text-rose-400'">c</kbd> {{ cancelPending ? 'Confirmar Cancelamento' : 'Cancelar (2x)' }}
          </button>
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click.stop="startPanic(10)"
            class="px-3 py-1.5 rounded-lg bg-amber-500/20 text-amber-300 border border-amber-500/30 hover:bg-amber-500/30 font-bold transition-all font-mono text-xs cursor-pointer"
          >
            ⚠️ Simular Pânico (10s)
          </button>
        </div>
      </template>

      <!-- Se estiver na tela de pânico: Ação Crítica de Abortar -->
      <template v-else-if="tuiState === 'CONFIRMATION'">
        <div class="flex items-center space-x-2">
          <span class="text-red-400 font-bold animate-pulse font-mono text-xs">⚠ PÂNICO: O COMPUTADOR IRÁ DESLIGAR!</span>
        </div>
        <div>
          <button
            @click.stop="handleCancel"
            class="px-4 py-1.5 rounded-lg bg-red-600 hover:bg-red-500 text-white font-bold transition-all font-mono text-xs shadow-md shadow-red-600/40 cursor-pointer"
          >
            [c] Abortar Desligamento Imediato
          </button>
        </div>
      </template>

      <!-- Se estiver no Menu TUI -->
      <template v-else-if="tuiState === 'MENU'">
        <div class="text-slate-400 font-mono text-xs">
          Use <kbd class="px-1 py-0.5 rounded bg-slate-800 text-slate-300">↑</kbd> <kbd class="px-1 py-0.5 rounded bg-slate-800 text-slate-300">↓</kbd> para navegar e <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300">Enter</kbd> para selecionar.
        </div>
        <div>
          <button
            @click.stop="exitTuiToShell"
            class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-mono text-xs border border-slate-700 cursor-pointer"
          >
            Voltar ao Bash (q)
          </button>
        </div>
      </template>

      <!-- Se estiver nos estados de saída -->
      <template v-else>
        <div class="text-slate-400 font-mono text-xs">
          Pressione <kbd class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300">Enter</kbd> ou clique para voltar ao terminal.
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click.stop="exitTuiToShell"
            class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 font-mono text-xs border border-slate-700 cursor-pointer"
          >
            Voltar ao Shell
          </button>
        </div>
      </template>
    </div>

    <!-- MODO GUIADO: "Aprenda na Prática" (Missões Interativas Clicáveis em Tempo Real) -->
    <div class="relative z-10 px-3 sm:px-4 py-3 bg-slate-950/90 border-t border-slate-800/80">
      <div class="flex items-center justify-between mb-2.5">
        <span class="text-xs font-bold text-cyan-400 uppercase tracking-wider flex items-center gap-1.5 min-w-0">
          <Icon name="heroicons:academic-cap" class="w-4 h-4 text-cyan-400 shrink-0" />
          <span class="truncate">Aprenda na Prática</span>
          <span class="hidden sm:inline font-normal text-slate-400 truncate">— Missões Clicáveis</span>
        </span>
        <div class="flex items-center space-x-3">
          <div class="w-20 sm:w-32 h-1.5 rounded-full bg-slate-800 overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-violet-500 to-emerald-400 transition-all duration-500 rounded-full"
              :style="{ width: `${(missions.filter(m => m.completed).length / missions.length) * 100}%` }"
            />
          </div>
          <span class="text-[11px] font-mono shrink-0" :class="missions.every(m => m.completed) ? 'text-emerald-400 font-bold' : 'text-slate-400'">
            {{ missions.filter(m => m.completed).length }}/{{ missions.length }} Feito
          </span>
        </div>
      </div>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-2">
        <div
          v-for="mission in missions"
          :key="mission.id"
          @click="clickMission(mission)"
          :class="[
            'p-2.5 rounded-xl border text-xs transition-all duration-300 flex items-start space-x-2.5 cursor-pointer',
            mission.completed
              ? 'border-emerald-500/50 bg-emerald-950/30 text-emerald-300 shadow-sm shadow-emerald-900/20'
              : 'border-slate-800/80 bg-slate-900/40 text-slate-400 hover:border-violet-500/50 hover:bg-slate-900/70 hover:text-slate-200'
          ]"
          :title="mission.completed ? 'Missão concluída!' : 'Clique para executar esta missão'"
        >
          <div class="shrink-0 mt-0.5">
            <span
              v-if="mission.completed"
              class="inline-flex items-center justify-center w-4 h-4 rounded-full bg-emerald-500/20 text-emerald-400 border border-emerald-500/40"
            >
              <Icon name="heroicons:check" class="w-3 h-3" />
            </span>
            <span
              v-else
              class="inline-flex items-center justify-center w-4 h-4 rounded-full bg-slate-800 text-slate-500 border border-slate-700"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
            </span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <span class="font-bold truncate" :class="mission.completed ? 'text-emerald-300' : 'text-slate-200'">{{ mission.title }}</span>
              <span
                v-if="mission.completed"
                class="ml-1 text-[9px] uppercase px-1.5 py-0.2 rounded bg-emerald-900/60 text-emerald-300 border border-emerald-700/60 font-mono shrink-0"
              >
                Feito
              </span>
            </div>
            <div class="text-[10px] text-slate-400 line-clamp-2 leading-tight mt-0.5">{{ mission.desc }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
