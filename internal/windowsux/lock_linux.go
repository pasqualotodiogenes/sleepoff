//go:build linux

package windowsux

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type SingleInstanceLock struct {
	file *os.File
	path string
}

// Lock obtém o lock de instância única no Linux via POSIX Flock em $XDG_RUNTIME_DIR/sleepoff.lock
// com fallback para /tmp/sleepoff.lock.
func Lock() (*SingleInstanceLock, bool, error) {
	path := getLinuxLockPath()
	return acquireFlock(path)
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

func getLinuxLockPath() string {
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		if fi, err := os.Stat(xdg); err == nil && fi.IsDir() {
			return filepath.Join(xdg, "sleepoff.lock")
		}
	}
	return "/tmp/sleepoff.lock"
}

func acquireFlock(path string) (*SingleInstanceLock, bool, error) {
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
