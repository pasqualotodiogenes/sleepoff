//go:build !windows
// +build !windows

package shutdown

import "fmt"

// Beep é no-op fora do Windows.
func Beep(_ int, _ int) {}

// Execute fora do Windows só permite dry-run.
func Execute(dryRun bool) error {
	if dryRun {
		return nil
	}
	return fmt.Errorf("shutdown real suportado apenas no Windows")
}
