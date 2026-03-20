package quarksearch

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type stubPlugin struct {
	name    string
	results []SearchResult
	err     error
}

func (p stubPlugin) Name() string { return p.name }

func (p stubPlugin) Search(_ context.Context, _ *http.Client, _ string, _ map[string]any) ([]SearchResult, error) {
	if p.err != nil {
		return nil, p.err
	}
	return append([]SearchResult(nil), p.results...), nil
}

func TestEngineSearchAllSourceDegradesTelegramFailureWithPluginResults(t *testing.T) {
	tgErr := errors.New("telegram timeout")
	engine := NewEngine()
	engine.plugins = map[string]Plugin{
		"stub": stubPlugin{
			name: "stub",
			results: []SearchResult{
				{
					UniqueID: "plugin-1",
					Title:    "plugin result",
					Content:  "plugin result",
					Links: []Link{
						{Type: "quark", URL: "https://pan.quark.cn/s/test1"},
					},
				},
			},
		},
	}
	engine.telegramSearchFn = func(context.Context, *http.Client, Settings, string, []string) ([]SearchResult, error) {
		return nil, tgErr
	}

	resp, err := engine.Search(context.Background(), Settings{
		Channels:       []string{"pikpak_share_channel"},
		EnabledPlugins: []string{"stub"},
	}, Request{
		Keyword:    "test",
		SourceType: "all",
	})
	if err != nil {
		t.Fatalf("expected degraded success, got error: %v", err)
	}
	if got := len(resp.MergedByType["quark"]); got != 1 {
		t.Fatalf("expected 1 merged result, got %d", got)
	}
}

func TestEngineSearchAllSourceDegradesTelegramFailureWithEmptyPluginResults(t *testing.T) {
	engine := NewEngine()
	engine.plugins = map[string]Plugin{
		"stub": stubPlugin{name: "stub"},
	}
	engine.telegramSearchFn = func(context.Context, *http.Client, Settings, string, []string) ([]SearchResult, error) {
		return nil, errors.New("telegram timeout")
	}

	resp, err := engine.Search(context.Background(), Settings{
		Channels:       []string{"pikpak_share_channel"},
		EnabledPlugins: []string{"stub"},
	}, Request{
		Keyword:    "test",
		SourceType: "all",
	})
	if err != nil {
		t.Fatalf("expected empty success response, got error: %v", err)
	}
	if resp.Total != 0 {
		t.Fatalf("expected empty response total, got %d", resp.Total)
	}
}

func TestEngineSearchAllSourceReturnsErrorWhenAllEnabledSourcesFail(t *testing.T) {
	tgErr := errors.New("telegram timeout")
	pluginErr := errors.New("plugin timeout")
	engine := NewEngine()
	engine.plugins = map[string]Plugin{
		"stub": stubPlugin{name: "stub", err: pluginErr},
	}
	engine.telegramSearchFn = func(context.Context, *http.Client, Settings, string, []string) ([]SearchResult, error) {
		return nil, tgErr
	}

	_, err := engine.Search(context.Background(), Settings{
		Channels:       []string{"pikpak_share_channel"},
		EnabledPlugins: []string{"stub"},
	}, Request{
		Keyword:    "test",
		SourceType: "all",
	})
	if !errors.Is(err, tgErr) {
		t.Fatalf("expected telegram error when all enabled sources fail, got %v", err)
	}
}

func TestEngineSearchTelegramOnlyStillReturnsTelegramError(t *testing.T) {
	tgErr := errors.New("telegram timeout")
	engine := NewEngine()
	engine.telegramSearchFn = func(context.Context, *http.Client, Settings, string, []string) ([]SearchResult, error) {
		return nil, tgErr
	}

	_, err := engine.Search(context.Background(), Settings{
		Channels: []string{"pikpak_share_channel"},
	}, Request{
		Keyword:    "test",
		SourceType: "tg",
	})
	if !errors.Is(err, tgErr) {
		t.Fatalf("expected telegram-only error, got %v", err)
	}
}
