// Package ui contém estilos lipgloss inspirados no padrão Claude Code e Charm CLI.
package ui

import (
	"github.com/pasqualotodiogenes/sleepoff/internal/config"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Borda sutil estilo Claude Code / GitHub CLI (Slate suave, nunca neon berrante)
	ColBorderMuted   = lipgloss.Color("#334155") // slate-700
	ColBorderActive  = lipgloss.Color("#6366F1") // indigo
	ColBorderAlert   = lipgloss.Color("#EF4444") // red-500
	ColSelectedBg    = lipgloss.Color("#1E293B") // slate-800
	ColSelectedFg    = lipgloss.Color("#38BDF8") // sky-400
	ColAccentCyan    = lipgloss.Color("#06B6D4") // cyan-500

	// Card Principal com bordas arredondadas elegantes
	StyleCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColBorderMuted).
			Padding(0, 2)

	// Card com borda de Alerta (Pânico)
	StyleCardPanic = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColBorderAlert).
			Padding(0, 2)

	// Container geral
	StyleContainer = lipgloss.NewStyle().
			Padding(1, 2)

	// Badges de Cabeçalho
	StyleBadgeApp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")). // violeta claro
			Bold(true)

	StyleBadgeVersion = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#64748B")) // slate-500

	StyleBadgeDryRun = lipgloss.NewStyle().
				Background(lipgloss.Color("#1E293B")).
				Foreground(lipgloss.Color("#38BDF8")).
				Bold(true).
				Padding(0, 1)

	// Badges de Status do Timer
	StyleStatusRunning = lipgloss.NewStyle().
				Background(lipgloss.Color("#064E3B")).
				Foreground(lipgloss.Color("#34D399")).
				Bold(true).
				Padding(0, 1)

	StyleStatusPaused = lipgloss.NewStyle().
				Background(lipgloss.Color("#78350F")).
				Foreground(lipgloss.Color("#FBBF24")).
				Bold(true).
				Padding(0, 1)

	StyleTitle = lipgloss.NewStyle().
			Foreground(config.ColPrimary).
			Bold(true)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(config.ColDim)

	StyleSectionHeader = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F1F5F9")).
				Bold(true)

	// Display de Tempo Digital
	StyleTimeBig = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#38BDF8")).
			Bold(true)

	StyleTime = lipgloss.NewStyle().
			Foreground(config.ColText).
			Bold(true)

	// Linha divisória
	StyleDivider = lipgloss.NewStyle().
			Foreground(ColBorderMuted)

	// Atalhos no rodapé
	StyleKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CBD5E1")).
			Bold(true)

	StyleKeyDesc = lipgloss.NewStyle().
			Foreground(config.ColDim)

	StyleHelp = lipgloss.NewStyle().
			Foreground(config.ColDim)

	// Logs
	StyleLogInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	StyleLogWarn = lipgloss.NewStyle().Foreground(config.ColWarn).Bold(true)
	StyleLogCrit = lipgloss.NewStyle().Foreground(config.ColErr).Bold(true)

	StyleDim = lipgloss.NewStyle().Foreground(config.ColDim)

	// Pânico
	StylePanic = lipgloss.NewStyle().
			Foreground(config.ColErr).
			Bold(true)

	StylePanicSubtle = lipgloss.NewStyle().
				Foreground(config.ColWarn)
)
