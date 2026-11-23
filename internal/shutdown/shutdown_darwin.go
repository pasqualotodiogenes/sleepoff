//go:build darwin

package shutdown

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed countdown.wav
var countdownWavData []byte

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

// Beep emite o som oficial via afplay, com fallback para o som Tink / terminal bell (\a).
func Beep(_ int, _ int) {
	go func() {
		if len(countdownWavData) > 0 {
			tmp := filepath.Join(os.TempDir(), "sleepoff_countdown.wav")
			_ = os.WriteFile(tmp, countdownWavData, 0644)
			if err := exec.Command("afplay", tmp).Run(); err == nil {
				return
			}
		}
		cmd := exec.Command("afplay", "/System/Library/Sounds/Tink.aiff")
		if err := cmd.Run(); err != nil {
			fmt.Fprint(os.Stderr, "\a")
		}
	}()
}
