package updater

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	artifactFormatTarGz = "tar.gz"
	artifactFormatZip   = "zip"
)

func PlatformExecutableName(base string, goos string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return ""
	}
	if strings.EqualFold(strings.TrimSpace(goos), "windows") {
		return base + ".exe"
	}
	return base
}

func ServerExecutableName() string {
	return PlatformExecutableName("server", runtime.GOOS)
}

func ServerExecutableNameForGOOS(goos string) string {
	return PlatformExecutableName("server", goos)
}

func LauncherExecutableName() string {
	return PlatformExecutableName("launcher", runtime.GOOS)
}

func ArtifactKeyForPlatform(goos string, goarch string) (string, error) {
	goos = strings.TrimSpace(goos)
	goarch = strings.TrimSpace(goarch)
	if goos == "" || goarch == "" {
		return "", fmt.Errorf("平台信息不完整")
	}
	switch goos {
	case "linux", "windows", "darwin":
	default:
		return "", fmt.Errorf("当前平台 %s/%s 暂不支持在线更新", goos, goarch)
	}
	switch goarch {
	case "amd64", "arm64":
	default:
		return "", fmt.Errorf("当前架构 %s/%s 暂不支持在线更新", goos, goarch)
	}
	return goos + "-" + goarch, nil
}

func normalizeArtifactFormat(format string, url string) string {
	normalized := strings.ToLower(strings.TrimSpace(format))
	switch normalized {
	case artifactFormatTarGz, artifactFormatZip:
		return normalized
	}

	lowerURL := strings.ToLower(strings.TrimSpace(url))
	switch {
	case strings.HasSuffix(lowerURL, ".zip"):
		return artifactFormatZip
	case strings.HasSuffix(lowerURL, ".tar.gz"), strings.HasSuffix(lowerURL, ".tgz"):
		return artifactFormatTarGz
	default:
		return artifactFormatTarGz
	}
}
