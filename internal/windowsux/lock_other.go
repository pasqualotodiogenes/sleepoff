//go:build !windows && !linux && !darwin

package windowsux

type SingleInstanceLock struct{}

// Lock fallback para plataformas não mapeadas.
func Lock() (*SingleInstanceLock, bool, error) {
	return &SingleInstanceLock{}, false, nil
}

func AcquireSingleInstance() (*SingleInstanceLock, bool, error) {
	return Lock()
}

func (l *SingleInstanceLock) Release() {}
