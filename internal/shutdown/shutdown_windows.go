//go:build windows

// Package shutdown lida com a lógica de desligamento do sistema.
package shutdown

import (
	_ "embed"
	"os/exec"
	"syscall"
	"unsafe"
)

//go:embed countdown.wav
var countdownWavData []byte

// beepDll é a referência ao kernel32.dll para o Beep nativo do Windows.
var beepDll = syscall.NewLazyDLL("kernel32.dll")
var beepProc = beepDll.NewProc("Beep")

var winmmDll = syscall.NewLazyDLL("winmm.dll")
var playSoundProc = winmmDll.NewProc("PlaySoundW")

const (
	sndAsync      = 0x0001
	sndNoDefault  = 0x0002
	sndMemory     = 0x0004
)

// Beep emite o som oficial do contador via winmm PlaySoundW (em memória),
// com fallback para o speaker do sistema via kernel32.dll Beep.
func Beep(frequency, duration int) {
	if len(countdownWavData) > 0 {
		go playSoundProc.Call(
			uintptr(unsafe.Pointer(&countdownWavData[0])),
			0,
			uintptr(sndMemory|sndAsync|sndNoDefault),
		)
		return
	}
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
