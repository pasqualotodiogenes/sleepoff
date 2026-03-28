//go:build !windows

package windowsux

type SingleInstanceLock struct{}

func ApplyProcessIdentity() error { return nil }

func AcquireSingleInstance() (*SingleInstanceLock, bool, error) {
	return &SingleInstanceLock{}, false, nil
}

func (l *SingleInstanceLock) Release() {}

func ShowAlreadyRunningMessage() {}

func ShowUpdateInstallingMessage(_ string) {}

func ShowInfo(_, _ string) {}

func OpenURL(_ string) error { return nil }
