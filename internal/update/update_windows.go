//go:build windows

package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pasqualotodiogenes/sleepoff/internal/buildinfo"
)

const (
	releaseAPIURL        = "https://api.github.com/repos/pasqualotodiogenes/sleepoff/releases/latest"
	installerAssetName   = "sleepoff-setup.exe"
	checkInterval        = 7 * 24 * time.Hour
	stateRelativeDir     = "sleepoff"
	stateFileName        = "update-state.json"
	installedRelativeExe = "Programs\\sleepoff\\sleepoff.exe"
)

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type releaseResponse struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

type stateFile struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type Result struct {
	Checked         bool
	UpdateAvailable bool
	LatestVersion   string
	InstallerPath   string
}

func CheckAndPrepare(force bool) (Result, error) {
	var result Result

	if buildinfo.VersionString() == "dev" {
		return result, nil
	}

	if !force && !isInstalledBinary() {
		return result, nil
	}

	statePath, err := statePath()
	if err != nil {
		return result, err
	}

	state, _ := loadState(statePath)
	if !force && time.Since(state.LastCheckedAt) < checkInterval {
		return result, nil
	}

	if err := saveState(statePath, stateFile{LastCheckedAt: time.Now()}); err != nil {
		return result, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	latest, err := fetchLatestRelease(ctx)
	if err != nil {
		return result, err
	}

	result.Checked = true
	result.LatestVersion = normalizeVersion(latest.TagName)
	if !isVersionNewer(result.LatestVersion, buildinfo.VersionString()) {
		return result, nil
	}

	installerURL := ""
	for _, asset := range latest.Assets {
		if asset.Name == installerAssetName {
			installerURL = asset.BrowserDownloadURL
			break
		}
	}
	if installerURL == "" {
		return result, fmt.Errorf("installer asset %q not found in latest release", installerAssetName)
	}

	installerPath, err := downloadInstaller(ctx, installerURL, result.LatestVersion)
	if err != nil {
		return result, err
	}

	result.UpdateAvailable = true
	result.InstallerPath = installerPath
	return result, nil
}

func LaunchInstaller(installerPath string) error {
	cmd := exec.Command(installerPath)
	return cmd.Start()
}

func fetchLatestRelease(ctx context.Context) (releaseResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releaseAPIURL, nil)
	if err != nil {
		return releaseResponse{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "sleepoff/"+buildinfo.VersionString())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return releaseResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return releaseResponse{}, fmt.Errorf("github release check failed: %s (%s)", resp.Status, strings.TrimSpace(string(body)))
	}

	var parsed releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return releaseResponse{}, err
	}
	return parsed, nil
}

func downloadInstaller(ctx context.Context, url, version string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "sleepoff/"+buildinfo.VersionString())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("installer download failed: %s", resp.Status)
	}

	targetPath := filepath.Join(os.TempDir(), "sleepoff-update-"+version+".exe")
	file, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}

	return targetPath, nil
}

func statePath() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	stateDir := filepath.Join(baseDir, stateRelativeDir)
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(stateDir, stateFileName), nil
}

func loadState(path string) (stateFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return stateFile{}, err
	}
	var state stateFile
	if err := json.Unmarshal(content, &state); err != nil {
		return stateFile{}, err
	}
	return state, nil
}

func saveState(path string, state stateFile) error {
	content, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func isInstalledBinary() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return false
	}

	expected := filepath.Join(localAppData, installedRelativeExe)
	return strings.EqualFold(filepath.Clean(exePath), filepath.Clean(expected))
}

func normalizeVersion(v string) string {
	return strings.TrimPrefix(strings.TrimSpace(v), "v")
}

func isVersionNewer(latest, current string) bool {
	latestParts := splitVersion(latest)
	currentParts := splitVersion(current)

	for i := 0; i < len(latestParts) || i < len(currentParts); i++ {
		var latestValue, currentValue int
		if i < len(latestParts) {
			latestValue = latestParts[i]
		}
		if i < len(currentParts) {
			currentValue = currentParts[i]
		}
		if latestValue > currentValue {
			return true
		}
		if latestValue < currentValue {
			return false
		}
	}

	return false
}

func splitVersion(value string) []int {
	value = normalizeVersion(value)
	rawParts := strings.Split(value, ".")
	parts := make([]int, 0, len(rawParts))
	for _, part := range rawParts {
		if part == "" {
			parts = append(parts, 0)
			continue
		}
		var parsed int
		fmt.Sscanf(part, "%d", &parsed)
		parts = append(parts, parsed)
	}
	return parts
}
