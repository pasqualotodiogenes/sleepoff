//go:build windows

package windowsux

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	AppUserModelID = "pasqualotodiogenes.sleepoff"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	shell32                      = syscall.NewLazyDLL("shell32.dll")
	messageBoxW                  = user32.NewProc("MessageBoxW")
	registerApplicationRestart   = kernel32.NewProc("RegisterApplicationRestart")
	setCurrentProcessAppModelIDW = shell32.NewProc("SetCurrentProcessExplicitAppUserModelID")
)

func ApplyProcessIdentity() error {
	appID, err := syscall.UTF16PtrFromString(AppUserModelID)
	if err != nil {
		return err
	}

	if ret, _, callErr := setCurrentProcessAppModelIDW.Call(uintptr(unsafe.Pointer(appID))); ret != 0 {
		return fmt.Errorf("SetCurrentProcessExplicitAppUserModelID: %w", callErr)
	}

	if ret, _, callErr := registerApplicationRestart.Call(0, 0); ret != 0 {
		return fmt.Errorf("RegisterApplicationRestart: %w", callErr)
	}

	return nil
}

func ShowAlreadyRunningMessage() {
	ShowInfo("sleepoff já está aberto", "Use a instância atual pelo terminal, pelo atalho ou pela bandeja do sistema.")
}

func ShowUpdateInstallingMessage(version string) {
	ShowInfo("Atualização encontrada", "O sleepoff vai abrir o instalador da versão "+version+" agora.")
}

func OpenURL(url string) error {
	return exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start()
}

func ShowInfo(title, message string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)
	const mbOK = 0x00000000
	const mbIconInfo = 0x00000040
	const mbSetForeground = 0x00010000
	messageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(mbOK|mbIconInfo|mbSetForeground),
	)
}
