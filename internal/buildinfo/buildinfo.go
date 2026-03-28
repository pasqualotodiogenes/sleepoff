package buildinfo

import (
	"runtime/debug"
	"strings"
	"sync"
)

const (
	RepoURL   = "https://github.com/pasqualotodiogenes/sleepoff"
	RepoRef   = "github.com/pasqualotodiogenes/sleepoff"
	Publisher = "Diogenes Pasqualoto"
	Copyright = "Copyright (c) 2026 Diogenes Pasqualoto"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"

	versionOnce sync.Once
	versionText string
)

func VersionString() string {
	versionOnce.Do(func() {
		versionText = normalizeVersion(Version)
		if versionText != "" && versionText != "dev" {
			return
		}

		if info, ok := debug.ReadBuildInfo(); ok {
			candidate := normalizeVersion(info.Main.Version)
			if candidate != "" && candidate != "(devel)" {
				versionText = candidate
				return
			}
		}

		if versionText == "" || versionText == "(devel)" {
			versionText = "dev"
		}
	})

	return versionText
}

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	return strings.TrimPrefix(v, "v")
}
