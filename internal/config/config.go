// Package config define constantes globais, cores e configurações do app.
package config

import "github.com/charmbracelet/lipgloss"

// --- CORES (Tema Premium Dark) ---
var (
	ColPrimary   = lipgloss.Color("#7C3AED") // Roxo vibrante
	ColSecondary = lipgloss.Color("#06B6D4") // Cyan
	ColText      = lipgloss.Color("#E2E8F0") // Cinza claro
	ColDim       = lipgloss.Color("#64748B") // Cinza médio
	ColSuccess   = lipgloss.Color("#22C55E") // Verde
	ColWarn      = lipgloss.Color("#F59E0B") // Laranja
	ColErr       = lipgloss.Color("#EF4444") // Vermelho
	ColBg        = lipgloss.Color("#0F172A") // Azul escuro
)

// --- OPÇÕES DE TEMPO PRÉ-DEFINIDAS ---
type TimeOption struct {
	Title   string
	Desc    string
	Minutes int
}

var DefaultTimeOptions = []TimeOption{
	{Title: "15 min", Desc: "Tempo rápido", Minutes: 15},
	{Title: "30 min", Desc: "Meia hora", Minutes: 30},
	{Title: "45 min", Desc: "Episódio de série", Minutes: 45},
	{Title: "1 hora", Desc: "Sessão completa", Minutes: 60},
	{Title: "1h 30min", Desc: "Filme curto", Minutes: 90},
	{Title: "2 horas", Desc: "Filme longo", Minutes: 120},
	{Title: "Personalizado", Desc: "Definir minutos", Minutes: -1},
}

// --- APP INFO ---
const (
	AppName = "sleepoff"
	AppDesc = "Timer de desligamento automático"
)
