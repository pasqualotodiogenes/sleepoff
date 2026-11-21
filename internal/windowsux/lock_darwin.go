//go:build darwin

package windowsux

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

type SingleInstanceLock struct {
	file *os.File
	path string
}

// Lock obtém o lock de instância única no macOS via POSIX Flock em /tmp/sleepoff.lock.
func Lock() (*SingleInstanceLock, bool, error) {
	const path = "/tmp/sleepoff.lock"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, false, fmt.Errorf("abrir lock file %q: %w", path, err)
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		_ = file.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
			return nil, true, nil
		}
		return nil, false, fmt.Errorf("flock %q: %w", path, err)
	}

	return &SingleInstanceLock{file: file, path: path}, false, nil
}

// AcquireSingleInstance é compatível com chamadas existentes.
func AcquireSingleInstance() (*SingleInstanceLock, bool, error) {
	return Lock()
}

func (l *SingleInstanceLock) Release() {
	if l == nil || l.file == nil {
		return
	}
	_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	_ = l.file.Close()
	l.file = nil
}
