// Package model define o state machine e tipos do app.
package model

import (
	"fmt"
	"io"
	"time"

	"github.com/pasqualotodiogenes/sleepoff/internal/config"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- ESTADOS ---
type State int

const (
	StateSplash State = iota // Tela de splash inicial
	StateMenu
	StateCustomInput
	StateRunning
	StateConfirmation
	StateDone
)

// --- LOG ENTRY ---
type LogEntry struct {
	Time    string
	Message string
	Level   string // "INFO", "WARN", "CRIT"
}

// --- ITEM DO MENU ---
type menuItem struct {
	title, desc string
	minutes     int
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

// Delegate customizado (sem bordas)
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	title := i.title
	desc := i.desc

	if isSelected {
		title = lipgloss.NewStyle().
			Foreground(config.ColPrimary).
			Bold(true).
			Render("▸ " + title)
		desc = lipgloss.NewStyle().
			Foreground(config.ColSecondary).
			Render("  " + desc)
	} else {
		title = lipgloss.NewStyle().
			Foreground(config.ColText).
			Render("  " + title)
		desc = lipgloss.NewStyle().
			Foreground(config.ColDim).
			Render("  " + desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

// --- MODEL PRINCIPAL ---
type Model struct {
	State          State
	List           list.Model
	Input          textinput.Model
	Progress       progress.Model
	TotalDuration  time.Duration
	Remaining      time.Duration
	FinishTime     time.Time
	StartTime      time.Time
	Paused         bool
	Logs           []LogEntry
	Quitting       bool
	Width, Height  int
	PanicCountdown int // Segundos restantes na tela de pânico (10 -> 0)
	PanicDeadline  time.Time
	DryRun         bool // Se true, não executa shutdown real
	ShutdownError  string

	// Animações
	SplashFrame int       // Frame atual do splash (0-100)
	SplashStart time.Time // Quando começou o splash
	FadeOut     bool      // Se está fazendo fade out

	// Confirmação de cancelamento (2x para sair)
	CancelPending     bool      // Se já pressionou uma vez
	CancelPendingTime time.Time // Quando pressionou pela primeira vez
}

// --- CONSTRUCTOR ---
func New() Model {
	var items []list.Item
	for _, opt := range config.DefaultTimeOptions {
		items = append(items, menuItem{title: opt.Title, desc: opt.Desc, minutes: opt.Minutes})
	}

	// Lista customizada sem bordas
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(config.ColPrimary).
		Bold(true).
		Padding(0, 0, 0, 2)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(config.ColSecondary).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(config.ColText).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(config.ColDim).
		Padding(0, 0, 0, 2)

	l := list.New(items, delegate, 40, 16)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	ti := textinput.New()
	ti.Placeholder = "ex: 25"
	ti.Focus()
	ti.CharLimit = 4
	ti.Width = 10

	return Model{
		State:       StateSplash,
		SplashStart: time.Now(),
		List:        l,
		Input:       ti,
		Progress:    progress.New(progress.WithDefaultGradient()),
	}
}

// NewWithDuration cria um model já com duração definida (pra CLI)
func NewWithDuration(d time.Duration, dryRun bool) Model {
	m := New()
	m.State = StateRunning
	m.TotalDuration = d
	m.Remaining = m.TotalDuration
	m.StartTime = time.Now()
	m.FinishTime = time.Now().Add(m.Remaining)
	m.DryRun = dryRun

	// Log formatado inteligente
	if d.Seconds() < 60 {
		m.AddLog(fmt.Sprintf("Timer: %ds", int(d.Seconds())), "INFO")
	} else {
		m.AddLog(fmt.Sprintf("Timer: %dm", int(d.Minutes())), "INFO")
	}
	return m
}

// --- INIT ---
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd())
}

// --- TICK ---
type TickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// --- HELPERS ---
func (m *Model) AddLog(msg, level string) {
	entry := LogEntry{
		Time:    time.Now().Format("15:04:05"),
		Message: msg,
		Level:   level,
	}
	m.Logs = append(m.Logs, entry)
	if len(m.Logs) > 5 {
		m.Logs = m.Logs[1:]
	}
}
