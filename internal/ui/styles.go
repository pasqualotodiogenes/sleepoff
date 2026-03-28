// Package ui contém estilos lipgloss reutilizáveis.
package ui

import (
	"github.com/pasqualotodiogenes/sleepoff/internal/config"

	"github.com/charmbracelet/lipgloss"
)

// --- ESTILOS GLOBAIS ---
var (
	// Sem bordas - design limpo
	StyleContainer = lipgloss.NewStyle().
			Padding(1, 2)

	StyleTitle = lipgloss.NewStyle().
			Foreground(config.ColPrimary).
			Bold(true)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(config.ColDim).
			Italic(true)

	StyleStatusRunning = lipgloss.NewStyle().
				Foreground(config.ColSuccess).
				Bold(true)

	StyleStatusPaused = lipgloss.NewStyle().
				Foreground(config.ColWarn).
				Bold(true)

	StyleLogInfo = lipgloss.NewStyle().Foreground(config.ColDim)
	StyleLogWarn = lipgloss.NewStyle().Foreground(config.ColWarn)
	StyleLogCrit = lipgloss.NewStyle().Foreground(config.ColErr)

	StyleDim = lipgloss.NewStyle().Foreground(config.ColDim)

	StyleTime = lipgloss.NewStyle().
			Foreground(config.ColText).
			Bold(true)

	StyleTimeBig = lipgloss.NewStyle().
			Foreground(config.ColPrimary).
			Bold(true)

	StyleHelp = lipgloss.NewStyle().
			Foreground(config.ColDim)

	StylePanic = lipgloss.NewStyle().
			Foreground(config.ColErr).
			Bold(true)

	StylePanicSubtle = lipgloss.NewStyle().
				Foreground(config.ColWarn)
)
