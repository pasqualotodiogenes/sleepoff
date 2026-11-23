package model

import (
	"fmt"
	"strings"

	"github.com/pasqualotodiogenes/sleepoff/internal/buildinfo"
	"github.com/pasqualotodiogenes/sleepoff/internal/config"
	"github.com/pasqualotodiogenes/sleepoff/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

const cardWidth = 60

// Renderiza barra superior padronizada estilo Claude Code / GitHub CLI
func renderTopBar(title string, statusBadge string, dryRun bool) string {
	left := ui.StyleBadgeApp.Render("sleepoff") + " " + ui.StyleBadgeVersion.Render("v"+buildinfo.VersionString())
	if title != "" {
		left += ui.StyleDim.Render(" ── ") + ui.StyleSectionHeader.Render(title)
	}

	right := ""
	if statusBadge != "" {
		right += statusBadge + " "
	}
	if dryRun {
		right += ui.StyleBadgeDryRun.Render("DRY-RUN")
	}

	totalTextLen := lipgloss.Width(left) + lipgloss.Width(right)
	spaces := cardWidth - 4 - totalTextLen
	if spaces < 2 {
		spaces = 2
	}

	return left + strings.Repeat(" ", spaces) + right
}

func renderDivider() string {
	return ui.StyleDivider.Render(strings.Repeat("─", cardWidth-4))
}

// --- VIEW PRINCIPAL ---
func (m Model) View() string {
	if m.Quitting {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ui.ColBorderMuted).
			Padding(1, 3).
			Render(ui.StyleBadgeApp.Render("sleepoff") + " finalizado com sucesso. Até mais! 👋")
		return "\n" + box + "\n\n"
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

// --- 1. SPLASH SCREEN ---
func (m Model) viewSplash() string {
	progress := float64(m.SplashFrame) / 100.0
	if progress > 1.0 {
		progress = 1.0
	}

	var b strings.Builder
	b.WriteString(renderTopBar("", "", m.DryRun) + "\n")
	b.WriteString(renderDivider() + "\n\n")

	logo := []string{
		"███████╗██╗     ███████╗███████╗██████╗  ██████╗ ",
		"██╔════╝██║     ██╔════╝██╔════╝██╔══██╗██╔═══██╗",
		"███████╗██║     █████╗  █████╗  ██████╔╝██║   ██║",
		"╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ ██║   ██║",
		"███████║███████╗███████╗███████╗██║     ╚██████╔╝",
		"╚══════╝╚══════╝╚══════╝╚══════╝╚═╝      ╚═════╝ ",
	}

	logoStyle := lipgloss.NewStyle().
		Foreground(config.ColPrimary).
		Bold(true)

	for _, line := range logo {
		b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(logoStyle.Render(line)) + "\n")
	}
	b.WriteString("\n")

	barLen := 36
	filled := int(progress * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
	barStyled := lipgloss.NewStyle().Foreground(ui.ColAccentCyan).Render(bar)
	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(barStyled) + "\n\n")

	statusMsg := "● Carregando temporizador..."
	if progress > 0.85 {
		statusMsg = "✔ Pronto!"
	}
	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(ui.StyleDim.Render(statusMsg)) + "\n")

	return ui.StyleCard.Width(cardWidth).Render(b.String())
}

// --- 2. MENU INTERATIVO COMPACTO (ESTILO CLAUDE CODE / GUM) ---
func (m Model) viewMenu() string {
	var b strings.Builder

	b.WriteString(renderTopBar("", "", m.DryRun) + "\n")
	b.WriteString(renderDivider() + "\n\n")

	b.WriteString(ui.StyleSectionHeader.Render("Selecione o tempo de desligamento:") + "\n\n")

	// Renderização direta, compacta e elegante das opções
	curIdx := m.List.Index()
	for idx, opt := range config.DefaultTimeOptions {
		isSelected := idx == curIdx
		if isSelected {
			selector := lipgloss.NewStyle().Foreground(ui.ColAccentCyan).Bold(true).Render("❯ ")
			title := lipgloss.NewStyle().Foreground(ui.ColSelectedFg).Bold(true).Render(fmt.Sprintf("%-15s", opt.Title))
			desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render(opt.Desc)
			content := selector + title + "  " + desc

			styledLine := lipgloss.NewStyle().
				Background(ui.ColSelectedBg).
				Width(cardWidth - 4).
				Render(" " + content)
			b.WriteString(styledLine + "\n")
		} else {
			selector := "  "
			title := lipgloss.NewStyle().Foreground(config.ColText).Render(fmt.Sprintf("%-15s", opt.Title))
			desc := lipgloss.NewStyle().Foreground(config.ColDim).Render(opt.Desc)
			b.WriteString("  " + selector + title + "  " + desc + "\n")
		}
	}

	b.WriteString("\n" + renderDivider() + "\n")
	footer := ui.StyleKey.Render("↑/↓") + " " + ui.StyleKeyDesc.Render("navegar") + "  •  " +
		ui.StyleKey.Render("enter") + " " + ui.StyleKeyDesc.Render("confirmar") + "  •  " +
		ui.StyleKey.Render("q") + " " + ui.StyleKeyDesc.Render("sair")
	b.WriteString(footer)

	return ui.StyleCard.Width(cardWidth).Render(b.String())
}

// --- 3. INPUT DE TEMPO PERSONALIZADO ---
func (m Model) viewCustomInput() string {
	var b strings.Builder

	b.WriteString(renderTopBar("TEMPO PERSONALIZADO", "", m.DryRun) + "\n")
	b.WriteString(renderDivider() + "\n\n")

	b.WriteString(ui.StyleSectionHeader.Render("Quantos minutos antes de desligar?") + "\n\n")

	inputLine := lipgloss.NewStyle().
		Foreground(ui.ColSelectedFg).
		Bold(true).
		Render("❯ " + m.Input.View() + " min")
	b.WriteString(inputLine + "\n\n")

	b.WriteString(ui.StyleDim.Render("Valores aceitos: 1 a 1440 minutos (ex: 15, 30, 45, 90, 120)") + "\n\n")

	b.WriteString(renderDivider() + "\n")
	footer := ui.StyleKey.Render("enter") + " " + ui.StyleKeyDesc.Render("confirmar") + "  •  " +
		ui.StyleKey.Render("esc") + " " + ui.StyleKeyDesc.Render("voltar ao menu")
	b.WriteString(footer)

	return ui.StyleCard.Width(cardWidth).Render(b.String())
}

// --- 4. TIMER EM EXECUÇÃO (HERO DIGITAL) ---
func (m Model) viewRunning() string {
	var b strings.Builder

	statusBadge := ui.StyleStatusRunning.Render("● RODANDO")
	if m.Paused {
		statusBadge = ui.StyleStatusPaused.Render("❚❚ PAUSADO")
	}

	b.WriteString(renderTopBar("", statusBadge, m.DryRun) + "\n")
	b.WriteString(renderDivider() + "\n\n")

	// Cálculo do display digital do relógio
	totalSec := int(m.Remaining.Seconds())
	if totalSec < 0 {
		totalSec = 0
	}
	hours := totalSec / 3600
	mins := (totalSec % 3600) / 60
	secs := totalSec % 60

	var timeStr string
	if hours > 0 {
		timeStr = fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
	} else {
		timeStr = fmt.Sprintf("%02d:%02d", mins, secs)
	}

	bigTime := ui.StyleTimeBig.Render(timeStr)
	timeSub := ui.StyleDim.Render(fmt.Sprintf("Início: %s  •  Previsão: %s",
		m.StartTime.Format("15:04"),
		m.FinishTime.Format("15:04"),
	))

	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(bigTime) + "\n")
	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(timeSub) + "\n\n")

	// Barra de Progresso
	pct := 0.0
	if m.TotalDuration.Seconds() > 0 {
		pct = (m.TotalDuration.Seconds() - m.Remaining.Seconds()) / m.TotalDuration.Seconds()
		if pct < 0 {
			pct = 0
		}
		if pct > 1 {
			pct = 1
		}
	}

	barLen := 42
	filled := int(pct * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	progressBar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
	pctText := fmt.Sprintf(" %3.0f%%", pct*100)
	barLine := lipgloss.NewStyle().Foreground(ui.ColAccentCyan).Render(progressBar) +
		lipgloss.NewStyle().Foreground(config.ColDim).Render(pctText)

	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(barLine) + "\n\n")

	// Histórico de eventos (Timeline compacta estilo Claude Code)
	if len(m.Logs) > 0 {
		b.WriteString(ui.StyleDim.Render("EVENTOS") + "\n")
		startIdx := 0
		if len(m.Logs) > 3 {
			startIdx = len(m.Logs) - 3
		}
		for _, log := range m.Logs[startIdx:] {
			prefix := ui.StyleDim.Render("│ ")
			timeTag := ui.StyleDim.Render(log.Time)
			msgStyle := ui.StyleLogInfo
			if log.Level == "WARN" {
				msgStyle = ui.StyleLogWarn
			} else if log.Level == "CRIT" {
				msgStyle = ui.StyleLogCrit
			}
			b.WriteString(fmt.Sprintf("%s%s  %s\n", prefix, timeTag, msgStyle.Render(log.Message)))
		}
		b.WriteString("\n")
	}

	// Alerta de cancelamento pendente
	if m.CancelPending {
		alert := lipgloss.NewStyle().
			Foreground(config.ColErr).
			Bold(true).
			Render("⚠ Pressione C ou Esc novamente para cancelar!")
		b.WriteString(alert + "\n\n")
	}

	b.WriteString(renderDivider() + "\n")
	pauseKey := "p pausar"
	if m.Paused {
		pauseKey = "p retomar"
	}
	footer := ui.StyleKey.Render(pauseKey) + "  •  " +
		ui.StyleKey.Render("+ / - ajustar ±5m") + "  •  " +
		ui.StyleKey.Render("c cancelar (2x)")
	b.WriteString(footer)

	return ui.StyleCard.Width(cardWidth).Render(b.String())
}

// --- 5. TELA DE PÂNICO (ÚLTIMOS 10 SEGUNDOS) ---
func (m Model) viewPanic() string {
	var b strings.Builder

	left := lipgloss.NewStyle().Foreground(ui.ColBorderAlert).Bold(true).Render("⚠ DESLIGAMENTO IMINENTE")
	right := ""
	if m.DryRun {
		right = ui.StyleBadgeDryRun.Render("DRY-RUN")
	}
	spaces := cardWidth - 4 - lipgloss.Width(left) - lipgloss.Width(right)
	if spaces < 2 {
		spaces = 2
	}
	b.WriteString(left + strings.Repeat(" ", spaces) + right + "\n")
	b.WriteString(renderDivider() + "\n\n")

	remaining := int(m.PanicCountdown)
	if remaining < 0 {
		remaining = 0
	}

	alertMsg := fmt.Sprintf(">>>  DESLIGANDO EM %02d SEGUNDOS  <<<", remaining)
	b.WriteString(lipgloss.NewStyle().
		Width(cardWidth - 4).
		Align(lipgloss.Center).
		Foreground(ui.ColBorderAlert).
		Bold(true).
		Render(alertMsg) + "\n")

	subMsg := ui.StyleDim.Render("Salve seus arquivos. O sistema será desligado agora!")
	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(subMsg) + "\n\n")

	barLen := 42
	filled := (remaining * barLen) / 10
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
	barStyled := lipgloss.NewStyle().Foreground(ui.ColBorderAlert).Bold(true).Render(bar)
	b.WriteString(lipgloss.NewStyle().Width(cardWidth - 4).Align(lipgloss.Center).Render(barStyled) + "\n\n")

	b.WriteString(renderDivider() + "\n")
	abortBtn := lipgloss.NewStyle().
		Foreground(ui.ColBorderAlert).
		Bold(true).
		Render("[C] ou [Esc] ABORTAR DESLIGAMENTO IMEDIATAMENTE")
	b.WriteString(abortBtn)

	return ui.StyleCardPanic.Width(cardWidth).Render(b.String())
}
