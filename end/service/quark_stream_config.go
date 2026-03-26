package service

import (
	"ohome/model"
	"strconv"
	"strings"
	"sync"
)

const (
	quarkStreamWebProxyModeKey        = "quark_fs_web_proxy_mode"
	quarkStreamConcurrencyKey         = "quark_fs_concurrency"
	quarkStreamPartSizeMBKey          = "quark_fs_part_size_mb"
	quarkStreamChunkMaxRetriesKey     = "quark_fs_chunk_max_retries"
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
	case quarkStreamWebProxyModeKey, quarkStreamConcurrencyKey, quarkStreamPartSizeMBKey, quarkStreamChunkMaxRetriesKey:
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

	partSizeMB := parseQuarkConfigInt(values[quarkStreamPartSizeMBKey], defaultQuarkStreamPartSizeMB, 1)

	return quarkStreamConfig{
		WebProxyMode:    normalizeQuarkWebProxyMode(values[quarkStreamWebProxyModeKey]),
		Concurrency:     parseQuarkConfigInt(values[quarkStreamConcurrencyKey], defaultQuarkStreamConcurrency, 1),
		PartSize:        int64(partSizeMB) * 1024 * 1024,
		ChunkMaxRetries: parseQuarkConfigInt(values[quarkStreamChunkMaxRetriesKey], defaultQuarkStreamChunkMaxRetries, 0),
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
		quarkStreamConcurrencyKey,
		quarkStreamPartSizeMBKey,
		quarkStreamChunkMaxRetriesKey,
	}
}

func EnsureDefaultQuarkStreamConfigs() error {
	defaults := []struct {
		key    string
		name   string
		value  string
		remark string
	}{
		{
			key:    quarkStreamWebProxyModeKey,
			name:   "夸克播放代理模式",
			value:  defaultQuarkStreamWebProxyMode,
			remark: "夸克在线播放代理模式：native_proxy=本地代理，302_redirect=302直连",
		},
		{
			key:    quarkStreamConcurrencyKey,
			name:   "夸克播放并发数",
			value:  strconv.Itoa(defaultQuarkStreamConcurrency),
			remark: "夸克在线播放并发回源分片数，建议 2-4",
		},
		{
			key:    quarkStreamPartSizeMBKey,
			name:   "夸克播放分片大小MB",
			value:  strconv.Itoa(defaultQuarkStreamPartSizeMB),
			remark: "夸克在线播放每个分片大小，单位 MB",
		},
		{
			key:    quarkStreamChunkMaxRetriesKey,
			name:   "夸克播放分片重试次数",
			value:  strconv.Itoa(defaultQuarkStreamChunkMaxRetries),
			remark: "夸克在线播放单个分片失败后的最大重试次数",
		},
	}

	for _, item := range defaults {
		if err := ensureConfigDefault(item.key, item.name, item.value, item.remark, ""); err != nil {
			return err
		}
	}
	invalidateQuarkStreamConfigCache()
	return nil
}
