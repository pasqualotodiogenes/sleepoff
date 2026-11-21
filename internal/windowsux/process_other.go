//go:build !windows && !linux && !darwin

package windowsux

func ApplyProcessIdentity() error { return nil }

func ShowAlreadyRunningMessage() {}

func ShowUpdateInstallingMessage(_ string) {}

func ShowInfo(_, _ string) {}

func OpenURL(_ string) error { return nil }
