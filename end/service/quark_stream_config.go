package service

import (
	"ohome/model"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	quarkStreamWebProxyModeKey        = "quark_fs_web_proxy_mode"
	quarkStreamConfigConcurrencyKey   = "quark.stream.concurrency"
	quarkStreamConfigPartSizeMBKey    = "quark.stream.partSizeMB"
	quarkStreamConfigChunkRetriesKey  = "quark.stream.chunkMaxRetries"
	defaultQuarkStreamWebProxyMode    = "native_proxy"
	defaultQuarkStreamConcurrency     = 3
	defaultQuarkStreamPartSizeMB      = 10
	defaultQuarkStreamChunkMaxRetries = 3
)

var quarkStreamConfigCache = struct {
	mu     sync.RWMutex
	loaded bool
	value  quarkStreamConfig
}{}

type quarkStreamConfig struct {
	WebProxyMode    string
	Concurrency     int
	PartSize        int64
	ChunkMaxRetries int
}

func loadQuarkStreamConfig() quarkStreamConfig {
	quarkStreamConfigCache.mu.RLock()
	if quarkStreamConfigCache.loaded {
		cfg := quarkStreamConfigCache.value
		quarkStreamConfigCache.mu.RUnlock()
		return cfg
	}
	quarkStreamConfigCache.mu.RUnlock()

	cfg, err := loadQuarkStreamConfigFromStore()
	if err != nil {
		cfg = defaultQuarkStreamConfig()
		logQuarkWarnf("[quarkFs:config] load runtime config failed, fallback to defaults err=%v", err)
	}

	quarkStreamConfigCache.mu.Lock()
	quarkStreamConfigCache.loaded = true
	quarkStreamConfigCache.value = cfg
	quarkStreamConfigCache.mu.Unlock()
	return cfg
}

func GetQuarkStreamWebProxyMode() string {
	return loadQuarkStreamConfig().WebProxyMode
}

func invalidateQuarkStreamConfigCache() {
	quarkStreamConfigCache.mu.Lock()
	quarkStreamConfigCache.loaded = false
	quarkStreamConfigCache.value = quarkStreamConfig{}
	quarkStreamConfigCache.mu.Unlock()
}

func isQuarkStreamConfigKey(key string) bool {
	switch strings.TrimSpace(key) {
	case quarkStreamWebProxyModeKey:
		return true
	default:
		return false
	}
}

func loadQuarkStreamConfigFromStore() (quarkStreamConfig, error) {
	configs, err := configDao.GetByKeys(quarkStreamConfigKeys())
	if err != nil {
		return quarkStreamConfig{}, err
	}
	return buildQuarkStreamConfig(configs), nil
}

func buildQuarkStreamConfig(configs []model.Config) quarkStreamConfig {
	values := make(map[string]string, len(configs))
	for _, cfg := range configs {
		values[strings.TrimSpace(cfg.Key)] = strings.TrimSpace(cfg.Value)
	}

	partSizeMB := parseQuarkConfigInt(
		viper.GetString(quarkStreamConfigPartSizeMBKey),
		defaultQuarkStreamPartSizeMB,
		1,
	)

	return quarkStreamConfig{
		WebProxyMode: normalizeQuarkWebProxyMode(values[quarkStreamWebProxyModeKey]),
		Concurrency: parseQuarkConfigInt(
			viper.GetString(quarkStreamConfigConcurrencyKey),
			defaultQuarkStreamConcurrency,
			1,
		),
		PartSize: int64(partSizeMB) * 1024 * 1024,
		ChunkMaxRetries: parseQuarkConfigInt(
			viper.GetString(quarkStreamConfigChunkRetriesKey),
			defaultQuarkStreamChunkMaxRetries,
			0,
		),
	}
}

func defaultQuarkStreamConfig() quarkStreamConfig {
	return quarkStreamConfig{
		WebProxyMode:    defaultQuarkStreamWebProxyMode,
		Concurrency:     defaultQuarkStreamConcurrency,
		PartSize:        int64(defaultQuarkStreamPartSizeMB) * 1024 * 1024,
		ChunkMaxRetries: defaultQuarkStreamChunkMaxRetries,
	}
}

func parseQuarkConfigInt(raw string, fallback int, minValue int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return fallback
	}
	if value < minValue {
		return fallback
	}
	return value
}

func normalizeQuarkWebProxyMode(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "redirect", "302", "302_redirect", "direct":
		return "302_redirect"
	default:
		return defaultQuarkStreamWebProxyMode
	}
}

func quarkStreamConfigKeys() []string {
	return []string{
		quarkStreamWebProxyModeKey,
	}
}
