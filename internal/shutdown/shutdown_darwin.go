//go:build darwin

package shutdown

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute realiza o desligamento real no macOS via Apple Events (System Events).
func Execute(dryRun bool) error {
	if dryRun {
		return nil
	}
	cmd := exec.Command("osascript", "-e", `tell app "System Events" to shut down`)
	return cmd.Run()
}

// Shutdown é o alias idiomático para Execute.
func Shutdown(dryRun bool) error {
	return Execute(dryRun)
}

// Beep emite o som Tink nativo do macOS via afplay, com fallback para o terminal bell (\a).
func Beep(_ int, _ int) {
	go func() {
		cmd := exec.Command("afplay", "/System/Library/Sounds/Tink.aiff")
		if err := cmd.Run(); err != nil {
			fmt.Fprint(os.Stderr, "\a")
		}
	}()
}
