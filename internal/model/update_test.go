package model

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func keyRunes(r string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(r)}
}

func TestSplashTransitionsToMenu(t *testing.T) {
	m := New()
	m.State = StateSplash
	m.SplashStart = time.Now().Add(-3 * time.Second)

	updated, cmd := m.Update(TickMsg(time.Now()))
	got := updated.(Model)

	if got.State != StateMenu {
		t.Fatalf("state = %v, want %v", got.State, StateMenu)
	}
	if cmd != nil {
		t.Fatal("expected nil cmd after splash transition")
	}
}

func TestStartTimerSchedulesTick(t *testing.T) {
	m := New()
	updated, cmd := m.startTimer(10)

	if updated.State != StateRunning {
		t.Fatalf("state = %v, want %v", updated.State, StateRunning)
	}
	if cmd == nil {
		t.Fatal("expected non-nil tick cmd")
	}
}

func TestPausedAdjustmentsKeepRemainingSynced(t *testing.T) {
	m := NewWithDuration(30*time.Minute, true)

	updated, _ := m.updateRunning(keyRunes("p"))
	m = updated.(Model)
	before := m.Remaining

	updated, _ = m.updateRunning(keyRunes("+"))
	m = updated.(Model)
	if m.Remaining != before+5*time.Minute {
		t.Fatalf("remaining after + = %v, want %v", m.Remaining, before+5*time.Minute)
	}

	updated, _ = m.updateRunning(keyRunes("-"))
	m = updated.(Model)
	if m.Remaining != before {
		t.Fatalf("remaining after - = %v, want %v", m.Remaining, before)
	}
}

func TestCancelInConfirmationReturnsCancelFlow(t *testing.T) {
	m := NewWithDuration(1*time.Minute, true)
	m.State = StateConfirmation
	m.PanicCountdown = 5
	m.PanicDeadline = time.Now().Add(5 * time.Second)

	updated, cmd := m.updateConfirmation(keyRunes("c"))
	got := updated.(Model)

	if got.State != StateRunning {
		t.Fatalf("state = %v, want %v", got.State, StateRunning)
	}
	if !got.Quitting {
		t.Fatal("expected quitting=true")
	}
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd")
	}
}

func TestExpiredPanicGoesToDone(t *testing.T) {
	m := NewWithDuration(1*time.Minute, true)
	m.State = StateConfirmation
	m.PanicCountdown = 1
	m.PanicDeadline = time.Now().Add(-100 * time.Millisecond)

	updated, cmd := m.Update(TickMsg(time.Now()))
	got := updated.(Model)

	if got.State != StateDone {
		t.Fatalf("state = %v, want %v", got.State, StateDone)
	}
	if !got.Quitting {
		t.Fatal("expected quitting=true")
	}
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd")
	}
}
