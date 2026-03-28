package model

import (
	"fmt"
	"strings"

	"github.com/pasqualotodiogenes/sleepoff/internal/buildinfo"
	"github.com/pasqualotodiogenes/sleepoff/internal/config"
	"github.com/pasqualotodiogenes/sleepoff/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

// --- VIEW PRINCIPAL ---
func (m Model) View() string {
	if m.Quitting {
		return "\n  Até mais! 👋\n\n"
	}

	switch m.State {
	case StateSplash:
		return m.viewSplash()
	case StateMenu:
		return m.viewMenu()
	case StateCustomInput:
		return m.viewCustomInput()
	case StateRunning:
		return m.viewRunning()
	case StateConfirmation:
		return m.viewPanic()
	}
	return ""
}

// --- SPLASH SCREEN ANIMADO ---
func (m Model) viewSplash() string {
	// Progresso da animação (0.0 a 1.0)
	progress := float64(m.SplashFrame) / 100.0

	// Logo ASCII
	logo := []string{
		"███████╗██╗     ███████╗███████╗██████╗  ██████╗ ███████╗███████╗",
		"██╔════╝██║     ██╔════╝██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔════╝",
		"███████╗██║     █████╗  █████╗  ██████╔╝██║   ██║█████╗  █████╗  ",
		"╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ ██║   ██║██╔══╝  ██╔══╝  ",
		"███████║███████╗███████╗███████╗██║     ╚██████╔╝██║     ██║     ",
		"╚══════╝╚══════╝╚══════╝╚══════╝╚═╝      ╚═════╝ ╚═╝     ╚═╝     ",
	}

	s := "\n\n\n"

	// Fade-in/out do logo
	opacity := 0.0
	if progress < 0.3 {
		// Fade in primeiro 30%
		opacity = progress / 0.3
	} else if progress > 0.7 {
		// Fade out últimos 30%
		opacity = 1.0 - ((progress - 0.7) / 0.3)
	} else {
		// Full brightness no meio
		opacity = 1.0
	}

	// Escolhe cor baseado na opacidade
	color := config.ColPrimary
	if opacity < 0.3 {
		color = config.ColDim
	} else if opacity < 0.6 {
		color = config.ColText
	}

	logoStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Align(lipgloss.Center)

	for _, line := range logo {
		s += logoStyle.Render(line) + "\n"
	}

	s += "\n\n"

	// Loading bar
	barLen := 40
	filled := int(progress * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)

	s += lipgloss.NewStyle().
		Foreground(config.ColSecondary).
		Align(lipgloss.Center).
		Render(bar) + "\n\n"

	// Mensagem
	msg := "Iniciando..."
	if progress > 0.7 {
		msg = "Pronto!"
	}

	s += ui.StyleDim.
		Align(lipgloss.Center).
		Render(msg)

	return s
}

// --- MENU ---
func (m Model) viewMenu() string {
	s := "\n"
	s += ui.StyleTitle.Render("  "+config.AppName) + " "
	s += ui.StyleSubtitle.Render("v"+buildinfo.VersionString()) + "\n"
	s += ui.StyleDim.Render("  "+config.AppDesc) + "\n\n"

	s += m.List.View() + "\n"

	s += ui.StyleHelp.Render("  ↑/↓ navegar • enter selecionar • q sair")

	return s
}

// --- INPUT PERSONALIZADO ---
func (m Model) viewCustomInput() string {
	s := "\n"
	s += ui.StyleTitle.Render("  Tempo personalizado") + "\n\n"
	s += "  Digite os minutos:\n\n"
	s += "  " + m.Input.View() + "\n\n"
	s += ui.StyleHelp.Render("  enter confirmar • esc voltar")
	return s
}

// --- TIMER RODANDO ---
func (m Model) viewRunning() string {
	s := "\n"

	// Header
	s += ui.StyleTitle.Render("  "+config.AppName) + "\n"
	s += ui.StyleDim.Render("  Desligamento programado") + "\n\n"

	// Tempo grande
	min := int(m.Remaining.Minutes())
	sec := int(m.Remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", min, sec)
	s += "  " + ui.StyleTimeBig.Render(timeStr) + "\n\n"

	// Barra de progresso
	pct := 0.0
	if m.TotalDuration.Seconds() > 0 {
		pct = (m.TotalDuration.Seconds() - m.Remaining.Seconds()) / m.TotalDuration.Seconds()
	}
	progWidth := m.Width - 6
	if progWidth < 20 {
		progWidth = 20
	}
	if progWidth > 60 {
		progWidth = 60
	}
	m.Progress.Width = progWidth
	s += "  " + m.Progress.ViewAs(pct) + "\n\n"

	// Status
	status := "▶ Rodando"
	statusStyle := ui.StyleStatusRunning
	if m.Paused {
		status = "⏸ Pausado"
		statusStyle = ui.StyleStatusPaused
	}
	s += "  " + statusStyle.Render(status) + "\n"

	// Horários
	s += ui.StyleDim.Render(fmt.Sprintf("  Início: %s  •  Fim: %s",
		m.StartTime.Format("15:04"),
		m.FinishTime.Format("15:04"),
	)) + "\n\n"

	// Logs (últimos eventos)
	if len(m.Logs) > 0 {
		s += ui.StyleDim.Render("  "+strings.Repeat("─", 30)) + "\n"
		for _, log := range m.Logs {
			logStyle := ui.StyleLogInfo
			switch log.Level {
			case "WARN":
				logStyle = ui.StyleLogWarn
			case "CRIT":
				logStyle = ui.StyleLogCrit
			}
			s += fmt.Sprintf("  %s %s\n",
				ui.StyleDim.Render(log.Time),
				logStyle.Render(log.Message),
			)
		}
		s += "\n"
	}

	// Aviso de cancelamento pendente
	if m.CancelPending {
		cancelWarning := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Render("  ⚠ Pressione C novamente para sair")
		s += cancelWarning + "\n\n"
	}

	// Controles
	s += ui.StyleHelp.Render("  [p] pausar  [c] cancelar  [+/-] ajustar tempo")

	return s
}

// --- TELA DE PÂNICO (sutil) ---
func (m Model) viewPanic() string {
	s := "\n\n"

	// Ícone sutil
	s += lipgloss.NewStyle().
		Foreground(config.ColWarn).
		Render("  ⏰") + "\n\n"

	s += ui.StylePanicSubtle.Render("  Tempo esgotado") + "\n\n"

	// Countdown
	remaining := int(m.PanicCountdown)
	if remaining < 0 {
		remaining = 0
	}

	s += "  Desligando em "
	s += ui.StylePanic.Render(fmt.Sprintf("%d", remaining))
	s += " segundos...\n\n"

	// Barra visual do countdown
	barLen := 20
	filled := (remaining * barLen) / 10
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
	s += "  " + ui.StylePanicSubtle.Render(bar) + "\n\n"

	s += ui.StyleHelp.Render("  [c] cancelar agora")

	return s
}
