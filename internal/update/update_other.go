//go:build !windows

package update

type Result struct {
	Checked         bool
	UpdateAvailable bool
	LatestVersion   string
	InstallerPath   string
}

func CheckAndPrepare(_ bool) (Result, error) { return Result{}, nil }

func LaunchInstaller(_ string) error { return nil }
