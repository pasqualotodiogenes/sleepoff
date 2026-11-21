/**
 * sleepoff — Interactive TUI Simulator
 * Faithful JavaScript reimplementation of the Bubble Tea state machine.
 * Zero external frameworks. Pure Web Audio API & DOM rendering.
 */

(function () {
  'use strict';

  // --- Paleta de Cores (config.go) ---
  const COLORS = {
    primary: '#7C3AED',
    secondary: '#06B6D4',
    text: '#E2E8F0',
    dim: '#64748B',
    success: '#22C55E',
    warn: '#F59E0B',
    err: '#EF4444',
    cancelAccent: '#FF6B6B',
  };

  // --- Estados da State Machine ---
  const State = {
    SPLASH: 'SPLASH',
    MENU: 'MENU',
    CUSTOM_INPUT: 'CUSTOM_INPUT',
    RUNNING: 'RUNNING',
    CONFIRMATION: 'CONFIRMATION', // Pânico de 10s
    DONE_TEST: 'DONE_TEST',       // Tela finalizada (teste)
    CANCELLED: 'CANCELLED',       // Tela de cancelamento com stats
    GOODBYE: 'GOODBYE'           // Saída amigável
  };

  // --- Opções de Tempo Padrão ---
  const TIME_OPTIONS = [
    { title: '15 min', desc: 'Tempo rápido', minutes: 15 },
    { title: '30 min', desc: 'Meia hora', minutes: 30 },
    { title: '45 min', desc: 'Episódio de série', minutes: 45 },
    { title: '1 hora', desc: 'Sessão completa', minutes: 60 },
    { title: '1h 30min', desc: 'Filme curto', minutes: 90 },
    { title: '2 horas', desc: 'Filme longo', minutes: 120 },
    { title: 'Personalizado', desc: 'Definir minutos', minutes: -1 }
  ];

  // --- Logo ASCII ---
  const LOGO_ASCII = [
    "███████╗██╗     ███████╗███████╗██████╗  ██████╗ ███████╗███████╗",
    "██╔════╝██║     ██╔════╝██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔════╝",
    "███████╗██║     █████╗  █████╗  ██████╔╝██║   ██║█████╗  █████╗  ",
    "╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ ██║   ██║██╔══╝  ██╔══╝  ",
    "███████║███████╗███████╗███████╗██║     ╚██████╔╝██║     ██║     ",
    "╚══════╝╚══════╝╚══════╝╚══════╝╚═╝      ╚═════╝ ╚═╝     ╚═╝     "
  ];

  // --- Web Audio API Synthesizer ---
  class AudioEngine {
    constructor() {
      this.ctx = null;
      this.muted = false;
    }

    init() {
      if (!this.ctx && typeof window.AudioContext !== 'undefined') {
        const AudioCtx = window.AudioContext || window.webkitAudioContext;
        this.ctx = new AudioCtx();
      }
      if (this.ctx && this.ctx.state === 'suspended') {
        this.ctx.resume();
      }
    }

    beep(freq, durationMs, type = 'sine') {
      if (this.muted) return;
      try {
        this.init();
        if (!this.ctx) return;

        const osc = this.ctx.createOscillator();
        const gain = this.ctx.createGain();

        osc.type = type;
        osc.frequency.setValueAtTime(freq, this.ctx.currentTime);

        const now = this.ctx.currentTime;
        const durSec = durationMs / 1000;

        // Envelope linear anti-clique (Attack-Decay)
        gain.gain.setValueAtTime(0, now);
        gain.gain.linearRampToValueAtTime(0.2, now + 0.02);
        gain.gain.setValueAtTime(0.2, now + durSec - 0.02);
        gain.gain.linearRampToValueAtTime(0.001, now + durSec);

        osc.connect(gain);
        gain.connect(this.ctx.destination);

        osc.start(now);
        osc.stop(now + durSec);
      } catch (e) {
        console.warn('Audio play error:', e);
      }
    }

    toggleMute() {
      this.muted = !this.muted;
      return this.muted;
    }
  }

  // --- Simulador TUI ---
  class SleepoffSimulator {
    constructor(screenEl) {
      this.screen = screenEl;
      this.audio = new AudioEngine();

      this.state = State.SPLASH;
      this.selectedIndex = 0;
      this.customInput = '';

      // Timer variables
      this.totalDurationMs = 0;
      this.remainingMs = 0;
      this.startTime = null;
      this.finishTime = null;
      this.paused = false;
      this.logs = [];

      // Cancelamento duplo
      this.cancelPending = false;
      this.cancelPendingTimeout = null;

      // Pânico
      this.panicCountdown = 10;
      this.panicDeadline = null;
      this.lastBeepCountdown = -1;

      // Splash animation
      this.splashStart = null;
      this.splashDurationMs = 2000;
      this.splashFrame = 0;

      // Loop timer
      this.timerInterval = null;

      this.startSplash();
    }

    // Helper de estilização ANSI / HTML
    span(text, color, bold = false, italic = false) {
      const styles = [];
      if (color) styles.push(`color: ${color}`);
      if (bold) styles.push('font-weight: 700');
      if (italic) styles.push('font-style: italic');
      return `<span style="${styles.join('; ')}">${this.escapeHtml(text)}</span>`;
    }

    escapeHtml(str) {
      return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
    }

    addLog(msg, level = 'INFO') {
      const now = new Date();
      const timeStr = [
        String(now.getHours()).padStart(2, '0'),
        String(now.getMinutes()).padStart(2, '0'),
        String(now.getSeconds()).padStart(2, '0')
      ].join(':');

      this.logs.push({ time: timeStr, message: msg, level });
      if (this.logs.length > 5) {
        this.logs.shift();
      }
    }

    startSplash() {
      this.state = State.SPLASH;
      this.splashStart = Date.now();
      this.splashFrame = 0;

      if (this.timerInterval) clearInterval(this.timerInterval);

      this.timerInterval = setInterval(() => {
        this.tickSplash();
      }, 30);
    }

    tickSplash() {
      const elapsed = Date.now() - this.splashStart;
      this.splashFrame = Math.min(100, Math.floor((elapsed / this.splashDurationMs) * 100));

      if (this.splashFrame >= 100) {
        clearInterval(this.timerInterval);
        this.state = State.MENU;
        this.render();
        return;
      }

      this.render();
    }

    startTimer(minutes) {
      if (this.timerInterval) clearInterval(this.timerInterval);

      this.state = State.RUNNING;
      this.totalDurationMs = minutes * 60 * 1000;
      this.remainingMs = this.totalDurationMs;
      this.startTime = new Date();
      this.finishTime = new Date(Date.now() + this.remainingMs);
      this.paused = false;
      this.logs = [];
      this.cancelPending = false;

      this.addLog(`Timer de ${minutes}min iniciado`, 'INFO');

      this.timerInterval = setInterval(() => {
        this.tickRunning();
      }, 100);

      this.render();
    }

    tickRunning() {
      if (this.state === State.RUNNING) {
        if (!this.paused) {
          const now = Date.now();
          this.remainingMs = Math.max(0, this.finishTime.getTime() - now);

          if (this.remainingMs <= 0) {
            this.triggerPanic();
            return;
          }
        }
        this.render();
      } else if (this.state === State.CONFIRMATION) {
        this.tickPanic();
      }
    }

    triggerPanic() {
      this.state = State.CONFIRMATION;
      this.panicCountdown = 10;
      this.panicDeadline = Date.now() + 10000;
      this.lastBeepCountdown = 10;

      this.addLog('Tempo esgotado!', 'CRIT');
      // Beep de 800Hz / 300ms no timeout exato do config/update.go
      this.audio.beep(800, 300, 'square');

      this.render();
    }

    jumpToPanic() {
      if (this.timerInterval) clearInterval(this.timerInterval);

      this.totalDurationMs = 15 * 60 * 1000;
      this.startTime = new Date(Date.now() - (15 * 60 * 1000 - 10000));
      this.finishTime = new Date(Date.now() + 10000);
      this.paused = false;
      this.logs = [
        { time: '00:00:01', message: 'Timer de 15min iniciado', level: 'INFO' },
        { time: '00:14:50', message: 'Reta final', level: 'WARN' }
      ];

      this.triggerPanic();

      this.timerInterval = setInterval(() => {
        this.tickPanic();
      }, 100);
    }

    tickPanic() {
      const remainingSec = Math.max(0, Math.ceil((this.panicDeadline - Date.now()) / 1000));
      this.panicCountdown = remainingSec;

      if (this.panicCountdown <= 0) {
        clearInterval(this.timerInterval);
        this.state = State.DONE_TEST;
        this.render();
        return;
      }

      if (this.panicCountdown > 0 && this.panicCountdown !== this.lastBeepCountdown) {
        // Frequência crescente conforme aproxima: 600 + (10 - countdown) * 50 Hz
        const freq = 600 + (10 - this.panicCountdown) * 50;
        this.audio.beep(freq, 100, 'sine');
        this.lastBeepCountdown = this.panicCountdown;
      }

      this.render();
    }

    handleCancel() {
      if (this.state === State.RUNNING) {
        if (this.cancelPending) {
          // Segundo clique dentro do prazo: Cancela de vez!
          if (this.cancelPendingTimeout) clearTimeout(this.cancelPendingTimeout);
          clearInterval(this.timerInterval);
          this.state = State.CANCELLED;
          this.render();
        } else {
          // Primeiro clique: Arma aviso
          this.cancelPending = true;
          this.render();
          this.cancelPendingTimeout = setTimeout(() => {
            this.cancelPending = false;
            this.render();
          }, 2000);
        }
      } else if (this.state === State.CONFIRMATION) {
        // Na tela de pânico, qualquer cancelamento sai direto
        clearInterval(this.timerInterval);
        this.state = State.CANCELLED;
        this.render();
      }
    }

    togglePause() {
      if (this.state !== State.RUNNING) return;
      this.paused = !this.paused;
      if (this.paused) {
        this.remainingMs = Math.max(0, this.finishTime.getTime() - Date.now());
        this.addLog('Pausado', 'WARN');
      } else {
        this.finishTime = new Date(Date.now() + this.remainingMs);
        this.addLog('Retomado', 'INFO');
      }
      this.render();
    }

    adjustTime(deltaMinutes) {
      if (this.state !== State.RUNNING) return;
      const deltaMs = deltaMinutes * 60 * 1000;
      this.remainingMs = Math.max(10000, this.remainingMs + deltaMs);
      this.totalDurationMs = Math.max(this.remainingMs, this.totalDurationMs + deltaMs);
      this.finishTime = new Date(Date.now() + this.remainingMs);

      const sign = deltaMinutes > 0 ? '+' : '';
      this.addLog(`${sign}${deltaMinutes}min adicionado`, deltaMinutes > 0 ? 'INFO' : 'WARN');
      this.render();
    }

    restart() {
      if (this.timerInterval) clearInterval(this.timerInterval);
      if (this.cancelPendingTimeout) clearTimeout(this.cancelPendingTimeout);
      this.selectedIndex = 0;
      this.customInput = '';
      this.logs = [];
      this.cancelPending = false;
      this.startSplash();
    }

    // --- Tratamento de Teclado ---
    handleKeyDown(e) {
      // Inicia áudio no primeiro gesto do usuário
      this.audio.init();

      const key = e.key;

      if (this.state === State.SPLASH) {
        // Pressionar qualquer tecla no splash pula direto para o menu
        clearInterval(this.timerInterval);
        this.state = State.MENU;
        this.render();
        return;
      }

      if (this.state === State.MENU) {
        if (key === 'ArrowUp' || key === 'k') {
          e.preventDefault();
          this.selectedIndex = (this.selectedIndex - 1 + TIME_OPTIONS.length) % TIME_OPTIONS.length;
          this.render();
        } else if (key === 'ArrowDown' || key === 'j') {
          e.preventDefault();
          this.selectedIndex = (this.selectedIndex + 1) % TIME_OPTIONS.length;
          this.render();
        } else if (key === 'Enter') {
          e.preventDefault();
          const opt = TIME_OPTIONS[this.selectedIndex];
          if (opt.minutes === -1) {
            this.state = State.CUSTOM_INPUT;
            this.customInput = '';
            this.render();
          } else {
            this.startTimer(opt.minutes);
          }
        } else if (key === 'q' || key === 'Q') {
          e.preventDefault();
          this.state = State.GOODBYE;
          this.render();
        }
        return;
      }

      if (this.state === State.CUSTOM_INPUT) {
        if (key === 'Escape') {
          e.preventDefault();
          this.state = State.MENU;
          this.render();
        } else if (key === 'Enter') {
          e.preventDefault();
          const mins = parseInt(this.customInput, 10);
          if (mins && mins > 0) {
            this.startTimer(mins);
          }
        } else if (key === 'Backspace') {
          e.preventDefault();
          this.customInput = this.customInput.slice(0, -1);
          this.render();
        } else if (/^[0-9]$/.test(key)) {
          e.preventDefault();
          if (this.customInput.length < 4) {
            this.customInput += key;
            this.render();
          }
        }
        return;
      }

      if (this.state === State.RUNNING) {
        if (key === 'p' || key === 'P' || key === ' ') {
          e.preventDefault();
          this.togglePause();
        } else if (key === '+' || key === '=') {
          e.preventDefault();
          this.adjustTime(5);
        } else if (key === '-' || key === '_') {
          e.preventDefault();
          this.adjustTime(-5);
        } else if (key === 'c' || key === 'C' || key === 'q' || key === 'Q' || key === 'Escape') {
          e.preventDefault();
          this.handleCancel();
        }
        return;
      }

      if (this.state === State.CONFIRMATION) {
        if (key === 'c' || key === 'C' || key === 'Escape') {
          e.preventDefault();
          this.handleCancel();
        }
        return;
      }

      if (this.state === State.CANCELLED || this.state === State.DONE_TEST || this.state === State.GOODBYE) {
        if (key === 'Enter' || key === 'r' || key === 'R' || key === 'Escape') {
          e.preventDefault();
          this.restart();
        }
      }
    }

    // --- Renderização dos Estados ---
    render() {
      let html = '';
      switch (this.state) {
        case State.SPLASH:
          html = this.renderSplash();
          break;
        case State.MENU:
          html = this.renderMenu();
          break;
        case State.CUSTOM_INPUT:
          html = this.renderCustomInput();
          break;
        case State.RUNNING:
          html = this.renderRunning();
          break;
        case State.CONFIRMATION:
          html = this.renderPanic();
          break;
        case State.CANCELLED:
          html = this.renderCancelled();
          break;
        case State.DONE_TEST:
          html = this.renderDoneTest();
          break;
        case State.GOODBYE:
          html = this.renderGoodbye();
          break;
      }
      this.screen.innerHTML = html;
    }

    // 1. Splash Screen
    renderSplash() {
      const progress = this.splashFrame / 100.0;
      let opacity = 1.0;
      if (progress < 0.3) {
        opacity = progress / 0.3;
      } else if (progress > 0.7) {
        opacity = 1.0 - ((progress - 0.7) / 0.3);
      }

      let color = COLORS.primary;
      if (opacity < 0.3) color = COLORS.dim;
      else if (opacity < 0.6) color = COLORS.text;

      let out = '\n\n';
      LOGO_ASCII.forEach(line => {
        out += `<div style="text-align: center; color: ${color}; font-weight: bold;">${this.escapeHtml(line)}</div>`;
      });
      out += '\n';

      const barLen = 40;
      const filled = Math.min(barLen, Math.floor(progress * barLen));
      const bar = '█'.repeat(filled) + '░'.repeat(barLen - filled);

      out += `<div style="text-align: center; color: ${COLORS.secondary}; font-weight: bold;">${bar}</div>\n`;

      const msg = progress > 0.7 ? 'Pronto!' : 'Iniciando...';
      out += `<div style="text-align: center; color: ${COLORS.dim};">${msg}</div>`;

      return out;
    }

    // 2. Menu Principal
    renderMenu() {
      let out = '\n';
      out += this.span('  sleepoff ', COLORS.primary, true);
      out += this.span('v1.0.0\n', COLORS.dim, false, true);
      out += this.span('  Timer de desligamento automático\n\n', COLORS.dim);

      TIME_OPTIONS.forEach((item, idx) => {
        const isSelected = idx === this.selectedIndex;
        if (isSelected) {
          out += this.span(`  ▸ ${item.title}\n`, COLORS.primary, true);
          out += this.span(`    ${item.desc}\n`, COLORS.secondary);
        } else {
          out += this.span(`    ${item.title}\n`, COLORS.text);
          out += this.span(`    ${item.desc}\n`, COLORS.dim);
        }
      });

      out += '\n' + this.span('  ↑/↓ navegar • enter selecionar • q sair', COLORS.dim);
      return out;
    }

    // 3. Custom Input
    renderCustomInput() {
      let out = '\n';
      out += this.span('  Tempo personalizado\n\n', COLORS.primary, true);
      out += '  Digite os minutos:\n\n';
      out += '  > ' + this.span(this.customInput || ' ', COLORS.secondary, true);
      out += '<span class="cursor" style="border-left: 2px solid #06B6D4; animation: blink 1s infinite;">&nbsp;</span>\n\n';
      out += this.span('  enter confirmar • esc voltar', COLORS.dim);
      return out;
    }

    // 4. Running (Timer em andamento)
    renderRunning() {
      let out = '\n';
      out += this.span('  sleepoff\n', COLORS.primary, true);
      out += this.span('  Desligamento programado\n\n', COLORS.dim);

      const totalSec = Math.floor(this.remainingMs / 1000);
      const mins = Math.floor(totalSec / 60);
      const secs = totalSec % 60;
      const timeStr = `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
      out += `  ` + this.span(timeStr, COLORS.primary, true) + '\n\n';

      // Barra de progresso com porcentagem
      const pct = this.totalDurationMs > 0 ? (this.totalDurationMs - this.remainingMs) / this.totalDurationMs : 0;
      const pctPercent = Math.min(100, Math.max(0, Math.floor(pct * 100)));
      const barLen = 32;
      const filled = Math.min(barLen, Math.floor(pct * barLen));
      const bar = '█'.repeat(filled) + '░'.repeat(barLen - filled);

      out += `  ${this.span(bar, COLORS.secondary)} ${this.span(pctPercent + '%', COLORS.dim)}\n\n`;

      // Status
      if (this.paused) {
        out += '  ' + this.span('⏸ Pausado\n', COLORS.warn, true);
      } else {
        out += '  ' + this.span('▶ Rodando\n', COLORS.success, true);
      }

      // Horários
      const fmtTime = (d) => {
        if (!d) return '--:--';
        return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
      };
      out += this.span(`  Início: ${fmtTime(this.startTime)}  •  Fim: ${fmtTime(this.finishTime)}\n\n`, COLORS.dim);

      // Logs
      if (this.logs.length > 0) {
        out += this.span('  ' + '─'.repeat(32) + '\n', COLORS.dim);
        this.logs.forEach(l => {
          let col = COLORS.dim;
          if (l.level === 'WARN') col = COLORS.warn;
          if (l.level === 'CRIT') col = COLORS.err;
          out += `  ${this.span(l.time, COLORS.dim)} ${this.span(l.message, col)}\n`;
        });
        out += '\n';
      }

      // Aviso de cancelamento pendente
      if (this.cancelPending) {
        out += this.span('  ⚠ Pressione C novamente para sair\n\n', COLORS.cancelAccent, true);
      }

      out += this.span('  [p] pausar  [c] cancelar  [+/-] ajustar tempo', COLORS.dim);
      return out;
    }

    // 5. Pânico (StateConfirmation)
    renderPanic() {
      let out = '\n\n';
      out += this.span('  ⏰\n\n', COLORS.warn);
      out += this.span('  Tempo esgotado\n\n', COLORS.warn);

      out += '  Desligando em ' + this.span(String(this.panicCountdown), COLORS.err, true) + ' segundos...\n\n';

      const barLen = 20;
      const filled = Math.min(barLen, Math.max(0, Math.floor((this.panicCountdown * barLen) / 10)));
      const bar = '█'.repeat(filled) + '░'.repeat(barLen - filled);
      out += '  ' + this.span(bar, COLORS.warn) + '\n\n';

      out += this.span('  [c] cancelar agora', COLORS.dim);
      return out;
    }

    // 6. Cancelado
    renderCancelled() {
      const elapsedMs = this.totalDurationMs - this.remainingMs;
      const elapsedSec = Math.floor(elapsedMs / 1000);
      const elapsedMins = Math.floor(elapsedSec / 60);
      const elapsedSecsMod = elapsedSec % 60;
      const totalMins = Math.round(this.totalDurationMs / 60000);
      const percentage = this.totalDurationMs > 0 ? Math.min(100, Math.floor((elapsedMs / this.totalDurationMs) * 100)) : 0;

      let out = `
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
   ╚═══════════════════════════════════════════════════════════════╝
`;
      out = this.span(out, COLORS.err);

      const barLen = 40;
      const filled = Math.min(barLen, Math.floor((percentage * barLen) / 100));
      const bar = '█'.repeat(filled) + '░'.repeat(barLen - filled);
      out += `\n   Progresso: ${this.span(bar, COLORS.success)} ${percentage}%\n\n`;

      out += this.span('   ┌─────────────────────────────────────┐\n', COLORS.dim);
      out += `   │  ${this.span('Tempo decorrido: ', COLORS.dim)}${this.span(String(elapsedMins).padStart(2, '0') + 'm ' + String(elapsedSecsMod).padStart(2, '0') + 's', COLORS.success)}               │\n`;
      out += `   │  ${this.span('Tempo total:     ', COLORS.dim)}${this.span(String(totalMins).padStart(2, '0') + 'm', COLORS.secondary)}                    │\n`;
      out += this.span('   └─────────────────────────────────────┘\n\n', COLORS.dim);

      let emoji = '😅';
      let msg = 'Desistiu rápido hein?';
      if (percentage >= 80) {
        emoji = '🏆';
        msg = 'QUASE LÁ! Tão perto do final!';
      } else if (percentage >= 50) {
        emoji = '💪';
        msg = 'PASSOU DA METADE! Respeitável.';
      } else if (percentage >= 20) {
        emoji = '🤔';
        msg = 'Bem... pelo menos tentou.';
      }

      out += `   ${emoji}  ${this.span(msg, COLORS.primary, true)}\n\n`;
      out += this.span('   ─────────────────────────────────────────────────────────────────\n', COLORS.dim);
      out += this.span('   [enter / r] reiniciar simulador', COLORS.dim);

      return out;
    }

    // 7. Finalizado (Teste)
    renderDoneTest() {
      let out = `
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
   ╚═══════════════════════════════════════════════════════════════╝
`;
      out = this.span(out, COLORS.success);
      out += '\n\n';
      out += this.span('   ℹ  O computador NÃO foi desligado.\n', COLORS.dim);
      out += this.span('   ℹ  Use sem --dry-run para desligar de verdade no terminal.\n\n', COLORS.dim);
      out += this.span('   ─────────────────────────────────────────────────────────────────\n', COLORS.dim);
      out += this.span('   [enter / r] reiniciar simulador', COLORS.dim);

      return out;
    }

    // 8. Goodbye
    renderGoodbye() {
      let out = `
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║       ██████╗ ██╗   ██╗███████╗██╗                            ║
   ║       ██╔══██╗╚██╗ ██╔╝██╔════╝██║                            ║
   ║       ██████╔╝ ╚████╔╝ █████╗  ██║                            ║
   ║       ██╔══██╗  ╚██╔╝  ██╔══╝  ╚═╝                            ║
   ║       ██████╔╝   ██║   ███████╗██╗                            ║
   ║       ╚═════╝    ╚═╝   ╚══════╝╚═╝                            ║
   ║                                                               ║
   ║                     👋 Até mais!                              ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝
`;
      out = this.span(out, COLORS.primary);
      out += '\n\n';
      out += this.span('   ─────────────────────────────────────────────────────────────────\n', COLORS.dim);
      out += this.span('   [enter / r] reiniciar simulador', COLORS.dim);
      return out;
    }
  }

  // --- Inicialização e Bindings DOM ---
  document.addEventListener('DOMContentLoaded', () => {
    const screenEl = document.getElementById('terminal-screen');
    if (!screenEl) return;

    const simulator = new SleepoffSimulator(screenEl);

    // Audio toggle
    const audioBtn = document.getElementById('btn-audio');
    if (audioBtn) {
      audioBtn.addEventListener('click', () => {
        const isMuted = simulator.audio.toggleMute();
        if (isMuted) {
          audioBtn.classList.add('muted');
          audioBtn.innerHTML = '🔕 Mudo';
        } else {
          audioBtn.classList.remove('muted');
          audioBtn.innerHTML = '🔔 Som On';
          simulator.audio.beep(800, 100);
        }
      });
    }

    // Teclado global
    window.addEventListener('keydown', (e) => {
      // Ignora se estiver focado em um input convencional fora da TUI
      if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
      simulator.handleKeyDown(e);
    });

    // Toolbar Buttons
    const bindBtn = (id, action) => {
      const el = document.getElementById(id);
      if (el) {
        el.addEventListener('click', (e) => {
          e.preventDefault();
          simulator.audio.init();
          action();
        });
      }
    };

    bindBtn('btn-nav-up', () => simulator.handleKeyDown({ key: 'ArrowUp', preventDefault: () => {} }));
    bindBtn('btn-nav-down', () => simulator.handleKeyDown({ key: 'ArrowDown', preventDefault: () => {} }));
    bindBtn('btn-nav-enter', () => simulator.handleKeyDown({ key: 'Enter', preventDefault: () => {} }));
    bindBtn('btn-toggle-pause', () => simulator.togglePause());
    bindBtn('btn-add-5m', () => simulator.adjustTime(5));
    bindBtn('btn-sub-5m', () => simulator.adjustTime(-5));
    bindBtn('btn-cancel', () => simulator.handleCancel());
    bindBtn('btn-jump-panic', () => simulator.jumpToPanic());
    bindBtn('btn-restart', () => simulator.restart());

    // Seletor de Tabs de Instalação
    const tabBtns = document.querySelectorAll('.tab-btn');
    const tabPanes = document.querySelectorAll('.tab-pane');

    tabBtns.forEach(btn => {
      btn.addEventListener('click', () => {
        const target = btn.getAttribute('data-tab');
        tabBtns.forEach(b => b.classList.remove('active'));
        tabPanes.forEach(p => p.classList.remove('active'));

        btn.classList.add('active');
        const activePane = document.getElementById(`tab-${target}`);
        if (activePane) activePane.classList.add('active');
      });
    });

    // Botões de Cópia em 1-Clique
    const copyBtns = document.querySelectorAll('.btn-copy');
    copyBtns.forEach(btn => {
      btn.addEventListener('click', async () => {
        const codeText = btn.getAttribute('data-clipboard');
        if (!codeText) return;

        try {
          await navigator.clipboard.writeText(codeText);
          const origHtml = btn.innerHTML;
          btn.classList.add('copied');
          btn.innerHTML = 'Copiado! ✓';
          setTimeout(() => {
            btn.classList.remove('copied');
            btn.innerHTML = origHtml;
          }, 2000);
        } catch (err) {
          console.error('Falha ao copiar:', err);
        }
      });
    });
  });
})();
