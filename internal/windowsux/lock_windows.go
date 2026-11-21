//go:build windows

package windowsux

import (
	"syscall"
	"unsafe"
)

const SingleInstanceMutexName = "sleepoff-single-instance"

var (
	kernel32Lock = syscall.NewLazyDLL("kernel32.dll")
	createMutexW = kernel32Lock.NewProc("CreateMutexW")
	closeHandle  = kernel32Lock.NewProc("CloseHandle")
)

type SingleInstanceLock struct {
	handle uintptr
}

// Lock obtém o lock de instância única no Windows usando Win32 Mutex nativo.
func Lock() (*SingleInstanceLock, bool, error) {
	name, err := syscall.UTF16PtrFromString(SingleInstanceMutexName)
	if err != nil {
		return nil, false, err
	}

	handle, _, callErr := createMutexW.Call(0, 0, uintptr(unsafe.Pointer(name)))
	if handle == 0 {
		return nil, false, callErr
	}

	if errno, ok := callErr.(syscall.Errno); ok && errno == syscall.ERROR_ALREADY_EXISTS {
		_ = syscall.CloseHandle(syscall.Handle(handle))
		return nil, true, nil
	}

	return &SingleInstanceLock{handle: handle}, false, nil
}

// AcquireSingleInstance é compatível com chamadas legadas do pacote.
func AcquireSingleInstance() (*SingleInstanceLock, bool, error) {
	return Lock()
}

func (l *SingleInstanceLock) Release() {
	if l == nil || l.handle == 0 {
		return
	}
	closeHandle.Call(l.handle)
	l.handle = 0
}
