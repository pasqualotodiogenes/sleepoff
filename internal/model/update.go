package model

import (
	"fmt"
	"math"
	"time"

	"github.com/pasqualotodiogenes/sleepoff/internal/shutdown"

	tea "github.com/charmbracelet/bubbletea"
)

// --- UPDATE PRINCIPAL ---
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		// Q para sair no menu
		if m.State == StateMenu && msg.String() == "q" {
			m.Quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height

		availableWidth := msg.Width - 4
		availableHeight := msg.Height - 4

		if availableWidth < 30 {
			availableWidth = 30
		}
		if availableHeight < 10 {
			availableHeight = 10
		}

		m.List.SetWidth(availableWidth)
		m.List.SetHeight(availableHeight)
		m.Progress.Width = availableWidth - 6

		return m, tea.ClearScreen

	case TickMsg:
		// Animação do splash
		if m.State == StateSplash {
			elapsed := time.Since(m.SplashStart).Milliseconds()
			// Splash dura 2 segundos (2000ms)
			m.SplashFrame = int((elapsed * 100) / 2000)

			if m.SplashFrame >= 100 {
				// Transição para menu
				m.State = StateMenu
				return m, nil
			}
			return m, tickCmd()
		}

		// Durante o timer principal
		if m.State == StateRunning {
			// Reseta confirmação de cancelamento se passou 2 segundos
			if m.CancelPending && time.Since(m.CancelPendingTime) >= 2*time.Second {
				m.CancelPending = false
			}

			if !m.Paused {
				m.Remaining = time.Until(m.FinishTime)
				if m.Remaining <= 0 {
					m.Remaining = 0
					m.State = StateConfirmation
					m.PanicCountdown = 10 // 10 segundos de pânico
					m.PanicDeadline = time.Now().Add(10 * time.Second)
					m.AddLog("Tempo esgotado!", "CRIT")
					shutdown.Beep(800, 300)
					return m, tickCmd()
				}
			}
		}

		// Durante a tela de pânico
		if m.State == StateConfirmation {
			remaining := time.Until(m.PanicDeadline)
			previousCountdown := m.PanicCountdown

			if remaining <= 0 {
				m.PanicCountdown = 0
				if err := shutdown.Execute(m.DryRun); err != nil {
					m.ShutdownError = err.Error()
					m.AddLog("Falha no shutdown: "+err.Error(), "CRIT")
				}
				m.State = StateDone
				m.Quitting = true
				return m, tea.Quit
			}

			m.PanicCountdown = int(math.Ceil(remaining.Seconds()))
			if m.PanicCountdown < 0 {
				m.PanicCountdown = 0
			}

			if m.PanicCountdown > 0 && m.PanicCountdown != previousCountdown {
				// Beep crescente conforme aproxima
				freq := 600 + (10-m.PanicCountdown)*50
				shutdown.Beep(freq, 100)
			}
		}

		return m, tickCmd()
	}

	// State Machine
	switch m.State {
	case StateMenu:
		return m.updateMenu(msg)
	case StateCustomInput:
		return m.updateCustomInput(msg)
	case StateRunning:
		return m.updateRunning(msg)
	case StateConfirmation:
		return m.updateConfirmation(msg)
	}

	return m, cmd
}

// --- SUB-HANDLERS ---

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			i, ok := m.List.SelectedItem().(menuItem)
			if ok {
				if i.minutes == -1 {
					m.State = StateCustomInput
					return m, nil
				}
				return m.startTimer(i.minutes)
			}
		}
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) updateCustomInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			var val int
			fmt.Sscanf(m.Input.Value(), "%d", &val)
			if val > 0 {
				return m.startTimer(val)
			}
		case "esc":
			m.State = StateMenu
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) updateRunning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p", " ":
			m.Paused = !m.Paused
			if m.Paused {
				m.Remaining = time.Until(m.FinishTime)
				if m.Remaining < 0 {
					m.Remaining = 0
				}
				m.AddLog("Pausado", "WARN")
			} else {
				m.AddLog("Retomado", "INFO")
				m.FinishTime = time.Now().Add(m.Remaining)
			}
		case "c", "q", "esc":
			// Confirmação dupla: precisa pressionar 2x em 2 segundos
			if m.CancelPending && time.Since(m.CancelPendingTime) < 2*time.Second {
				// Segunda vez - cancela de verdade
				m.AddLog("Cancelado", "CRIT")
				m.Quitting = true
				return m, tea.Quit
			}
			// Primeira vez - só marca (aviso visual aparece na view)
			m.CancelPending = true
			m.CancelPendingTime = time.Now()
		case "+", "=":
			// Adiciona 5 minutos ao tempo final
			m.FinishTime = m.FinishTime.Add(5 * time.Minute)
			if m.Paused {
				m.Remaining += 5 * time.Minute
			}
			m.TotalDuration += 5 * time.Minute
			m.AddLog("+5 min", "INFO")
		case "-", "_":
			// Remove 5 minutos, mas mantém mínimo de 1 minuto restante
			remaining := time.Until(m.FinishTime)
			if m.Paused {
				remaining = m.Remaining
			}
			if remaining-5*time.Minute >= time.Minute {
				m.FinishTime = m.FinishTime.Add(-5 * time.Minute)
				if m.Paused {
					m.Remaining -= 5 * time.Minute
				}
				m.TotalDuration -= 5 * time.Minute
				m.AddLog("-5 min", "WARN")
			}
		}
	}
	return m, nil
}

func (m Model) updateConfirmation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "c" || msg.String() == "esc" {
			// Usuário cancelou na última hora!
			m.AddLog("Cancelado!", "INFO")
			m.State = StateRunning
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) startTimer(minutes int) (Model, tea.Cmd) {
	m.State = StateRunning
	m.TotalDuration = time.Duration(minutes) * time.Minute
	m.Remaining = m.TotalDuration
	m.StartTime = time.Now()
	m.FinishTime = time.Now().Add(m.Remaining)
	m.Paused = false
	m.CancelPending = false
	m.PanicCountdown = 0
	m.PanicDeadline = time.Time{}
	m.ShutdownError = ""
	m.AddLog(fmt.Sprintf("Timer: %d min", minutes), "INFO")
	return m, tickCmd()
}
