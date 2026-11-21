//go:build !windows && !linux && !darwin

package shutdown

import (
	"fmt"
	"os"
)

// Execute fora dos sistemas suportados retorna erro caso não seja dryRun.
func Execute(dryRun bool) error {
	if dryRun {
		return nil
	}
	return fmt.Errorf("shutdown real suportado apenas no Windows, Linux e macOS")
}

// Shutdown é o alias idiomático para Execute.
func Shutdown(dryRun bool) error {
	return Execute(dryRun)
}

// Beep emite o terminal bell (\a) em plataformas não mapeadas.
func Beep(_ int, _ int) {
	fmt.Fprint(os.Stderr, "\a")
}
