package updater

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultUpdaterToken = "ohome-local-updater"
	defaultManifestURL  = "https://github.com/leftScience/ohome/releases/latest/download/server.json"
)

func DetectDeployMode() DeployMode {
	return DeployModeBinary
}

func ManifestURL() string {
	value := strings.TrimSpace(viper.GetString("update.manifestUrl"))
	if value == "" {
		return defaultManifestURL
	}
	return value
}

func DefaultChannel() string {
	value := strings.TrimSpace(viper.GetString("update.channel"))
	if value == "" {
		return "stable"
	}
	return value
}

func UpdaterToken() string {
	value := strings.TrimSpace(viper.GetString("update.updater.token"))
	if value == "" {
		return defaultUpdaterToken
	}
	return value
}

func UpdaterListenAddr() string {
	value := strings.TrimSpace(viper.GetString("update.updater.listenAddr"))
	if value == "" {
		return "127.0.0.1:18091"
	}
	return value
}

func UpdaterBaseURL() string {
	value := strings.TrimSpace(viper.GetString("update.updater.baseUrl"))
	if value == "" {
		return "http://127.0.0.1:18091"
	}
	return strings.TrimRight(value, "/")
}

func HealthURLForMode(_ DeployMode) string {
	configured := strings.TrimSpace(viper.GetString("update.healthUrl"))
	if configured != "" {
		return configured
	}
	port := strings.TrimSpace(viper.GetString("server.port"))
	if port == "" {
		port = "18090"
	}
	return fmt.Sprintf("http://127.0.0.1:%s/api/v1/public/discovery", port)
}
