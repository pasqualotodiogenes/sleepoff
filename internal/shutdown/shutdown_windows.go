//go:build windows

// Package shutdown lida com a lógica de desligamento do sistema.
package shutdown

import (
	"os/exec"
	"syscall"
)

// beepDll é a referência ao kernel32.dll para o Beep nativo do Windows.
var beepDll = syscall.NewLazyDLL("kernel32.dll")
var beepProc = beepDll.NewProc("Beep")

// Beep emite um som pelo speaker do sistema (Windows via kernel32.dll Beep).
func Beep(frequency, duration int) {
	go beepProc.Call(uintptr(frequency), uintptr(duration))
}

// Execute realiza o desligamento real do sistema.
// CUIDADO: Isso desliga o computador!
func Execute(dryRun bool) error {
	if dryRun {
		// Modo teste, não faz nada
		return nil
	}
	cmd := exec.Command("shutdown", "/s", "/t", "0")
	return cmd.Run()
}

// Shutdown é o alias idiomático para Execute.
func Shutdown(dryRun bool) error {
	return Execute(dryRun)
}
