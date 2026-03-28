//go:build windows

package windowsux

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pasqualotodiogenes/sleepoff/internal/buildinfo"
	"github.com/pasqualotodiogenes/sleepoff/internal/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lxn/walk"
)

type Integration struct {
	program       *tea.Program
	onCheckUpdate func()

	startOnce sync.Once
	closeOnce sync.Once

	stateCh chan model.DesktopState
	stopCh  chan struct{}

	runtimeMu sync.Mutex
	runtime   *trayRuntime
}

type trayRuntime struct {
	window       *walk.MainWindow
	notifyIcon   *walk.NotifyIcon
	pauseAction  *walk.Action
	cancelAction *walk.Action
	state        model.DesktopState
	initialized  bool
}

func NewIntegration(onCheckUpdate func()) *Integration {
	return &Integration{
		onCheckUpdate: onCheckUpdate,
		stateCh:       make(chan model.DesktopState, 1),
		stopCh:        make(chan struct{}),
	}
}

func (i *Integration) AttachProgram(program *tea.Program) {
	i.program = program
	i.start()
}

func (i *Integration) UpdateDesktop(state model.DesktopState) {
	i.start()
	select {
	case i.stateCh <- state:
	default:
		select {
		case <-i.stateCh:
		default:
		}
		i.stateCh <- state
	}
}

func (i *Integration) Close() {
	i.closeOnce.Do(func() {
		close(i.stopCh)

		i.runtimeMu.Lock()
		defer i.runtimeMu.Unlock()
		if i.runtime != nil && i.runtime.window != nil {
			i.runtime.window.Synchronize(func() {
				walk.App().Exit(0)
			})
		}
	})
}

func (i *Integration) start() {
	i.startOnce.Do(func() {
		go i.run()
	})
}

func (i *Integration) run() {
	runtime.LockOSThread()

	mw, err := walk.NewMainWindow()
	if err != nil {
		return
	}
	mw.SetVisible(false)

	notifyIcon, err := walk.NewNotifyIcon(mw)
	if err != nil {
		mw.Dispose()
		return
	}

	icon, err := loadTrayIcon()
	if err == nil {
		_ = notifyIcon.SetIcon(icon)
	}
	_ = notifyIcon.SetToolTip("sleepoff")

	pauseAction := walk.NewAction()
	_ = pauseAction.SetText("Pause timer")
	_ = pauseAction.SetEnabled(false)
	pauseAction.Triggered().Attach(func() {
		if i.program != nil {
			i.program.Send(model.TogglePauseMsg{})
		}
	})
	_ = notifyIcon.ContextMenu().Actions().Add(pauseAction)

	cancelAction := walk.NewAction()
	_ = cancelAction.SetText("Cancel timer")
	_ = cancelAction.SetEnabled(false)
	cancelAction.Triggered().Attach(func() {
		if i.program != nil {
			i.program.Send(model.ForceCancelMsg{})
		}
	})
	_ = notifyIcon.ContextMenu().Actions().Add(cancelAction)

	_ = notifyIcon.ContextMenu().Actions().Add(walk.NewSeparatorAction())

	updateAction := walk.NewAction()
	_ = updateAction.SetText("Check for updates")
	updateAction.Triggered().Attach(func() {
		if i.onCheckUpdate != nil {
			go i.onCheckUpdate()
		}
	})
	_ = notifyIcon.ContextMenu().Actions().Add(updateAction)

	releasesAction := walk.NewAction()
	_ = releasesAction.SetText("Open latest release")
	releasesAction.Triggered().Attach(func() {
		_ = OpenURL(buildinfo.RepoURL + "/releases/latest")
	})
	_ = notifyIcon.ContextMenu().Actions().Add(releasesAction)

	_ = notifyIcon.ContextMenu().Actions().Add(walk.NewSeparatorAction())

	exitAction := walk.NewAction()
	_ = exitAction.SetText("Exit sleepoff")
	exitAction.Triggered().Attach(func() {
		if i.program != nil {
			i.program.Quit()
			return
		}
		walk.App().Exit(0)
	})
	_ = notifyIcon.ContextMenu().Actions().Add(exitAction)

	notifyIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		title, body := traySummary(i.snapshotState())
		if body == "" {
			body = "sleepoff está pronto para uso."
		}
		_ = notifyIcon.ShowInfo(title, body)
	})

	if err := notifyIcon.SetVisible(true); err != nil {
		notifyIcon.Dispose()
		mw.Dispose()
		return
	}

	i.runtimeMu.Lock()
	i.runtime = &trayRuntime{
		window:       mw,
		notifyIcon:   notifyIcon,
		pauseAction:  pauseAction,
		cancelAction: cancelAction,
	}
	i.runtimeMu.Unlock()

	go func() {
		for {
			select {
			case state := <-i.stateCh:
				mw.Synchronize(func() {
					i.applyState(state)
				})
			case <-i.stopCh:
				mw.Synchronize(func() {
					walk.App().Exit(0)
				})
				return
			}
		}
	}()

	mw.Run()
	notifyIcon.Dispose()
	mw.Dispose()
}

func (i *Integration) applyState(state model.DesktopState) {
	i.runtimeMu.Lock()
	defer i.runtimeMu.Unlock()

	if i.runtime == nil || i.runtime.notifyIcon == nil {
		return
	}

	prev := i.runtime.state
	i.runtime.state = state
	i.runtime.initialized = true

	_ = i.runtime.notifyIcon.SetToolTip(truncateTooltip(trayTooltip(state)))

	running := state.State == model.StateRunning
	_ = i.runtime.pauseAction.SetEnabled(running)
	_ = i.runtime.cancelAction.SetEnabled(running || state.State == model.StateConfirmation)
	if running && state.Paused {
		_ = i.runtime.pauseAction.SetText("Resume timer")
	} else {
		_ = i.runtime.pauseAction.SetText("Pause timer")
	}

	if !shouldNotify(prev, state) {
		return
	}

	switch {
	case prev.State != model.StateRunning && state.State == model.StateRunning && !state.Paused && state.TotalDuration > 0:
		_ = i.runtime.notifyIcon.ShowInfo("Timer iniciado", fmt.Sprintf("sleepoff vai até %s.", state.FinishTime.Format("15:04")))
	case prev.State == model.StateRunning && !prev.Paused && state.Paused:
		_ = i.runtime.notifyIcon.ShowWarning("Timer pausado", fmt.Sprintf("Restante: %s.", formatRemaining(state.Remaining)))
	case prev.State == model.StateRunning && prev.Paused && !state.Paused:
		_ = i.runtime.notifyIcon.ShowInfo("Timer retomado", fmt.Sprintf("Restante: %s.", formatRemaining(state.Remaining)))
	case prev.State != model.StateConfirmation && state.State == model.StateConfirmation:
		_ = i.runtime.notifyIcon.ShowWarning("Desligamento em 10 segundos", "Clique no app ou use a bandeja para cancelar agora.")
	case prev.State == model.StateRunning && prev.Remaining > time.Minute && state.State == model.StateRunning && state.Remaining <= time.Minute && !state.Paused:
		_ = i.runtime.notifyIcon.ShowWarning("Menos de 1 minuto restante", "sleepoff está quase concluindo o timer.")
	case prev.State == model.StateRunning && state.State == model.StateDone && state.ShutdownError != "":
		_ = i.runtime.notifyIcon.ShowError("Falha no shutdown", state.ShutdownError)
	}
}

func (i *Integration) snapshotState() model.DesktopState {
	i.runtimeMu.Lock()
	defer i.runtimeMu.Unlock()
	if i.runtime == nil {
		return model.DesktopState{}
	}
	return i.runtime.state
}

func shouldNotify(prev, next model.DesktopState) bool {
	if prev.State == 0 && prev.TotalDuration == 0 && !prev.Paused && !prev.CancelPending && !prev.FinishTime.After(time.Time{}) {
		return next.State == model.StateRunning || next.State == model.StateConfirmation
	}

	if prev.State != next.State {
		return true
	}
	if prev.Paused != next.Paused {
		return true
	}
	if prev.Remaining > time.Minute && next.Remaining <= time.Minute && next.State == model.StateRunning {
		return true
	}
	return false
}

func trayTooltip(state model.DesktopState) string {
	switch state.State {
	case model.StateRunning:
		if state.Paused {
			return fmt.Sprintf("sleepoff - pausado (%s)", formatRemaining(state.Remaining))
		}
		return fmt.Sprintf("sleepoff - %s restante", formatRemaining(state.Remaining))
	case model.StateConfirmation:
		return fmt.Sprintf("sleepoff - desligando em %ds", state.PanicCountdown)
	case model.StateMenu:
		return "sleepoff - pronto"
	default:
		return "sleepoff"
	}
}

func traySummary(state model.DesktopState) (string, string) {
	switch state.State {
	case model.StateRunning:
		if state.Paused {
			return "sleepoff", "Timer pausado."
		}
		return "sleepoff", fmt.Sprintf("Timer ativo. Restante: %s.", formatRemaining(state.Remaining))
	case model.StateConfirmation:
		return "sleepoff", fmt.Sprintf("Desligando em %d segundos.", state.PanicCountdown)
	default:
		return "sleepoff", "Pronto para uso."
	}
}

func formatRemaining(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	totalSeconds := int(d.Round(time.Second).Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func truncateTooltip(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 63 {
		return value
	}
	return value[:60] + "..."
}

func loadTrayIcon() (*walk.Icon, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	return walk.NewIconExtractedFromFile(exePath, 0, 16)
}
