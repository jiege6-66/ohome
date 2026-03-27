package service

import (
	"errors"
	"ohome/global"
	"ohome/model"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestBuildQuarkTVLoginStatus(t *testing.T) {
	now := time.Now()
	snapshot := quarkTVConfigSnapshot{
		RefreshToken:    modelConfigWithUpdatedAt("refresh", now.Add(-2*time.Minute)),
		QueryToken:      modelConfigWithUpdatedAt("query", now),
		HasRefreshToken: true,
		HasQueryToken:   true,
	}

	status := buildQuarkTVLoginStatus(snapshot)
	if !status.Configured {
		t.Fatalf("expected configured to be true")
	}
	if !status.Pending {
		t.Fatalf("expected pending to be true")
	}
	if status.UpdatedAt == nil || !status.UpdatedAt.Equal(now) {
		t.Fatalf("expected latest updatedAt to be query token time, got %#v", status.UpdatedAt)
	}
}

func TestClassifyQuarkTVPollError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{name: "pending chinese", err: errors.New("等待扫码确认"), want: "pending"},
		{name: "expired english", err: errors.New("query token expired"), want: "expired"},
		{name: "generic", err: errors.New("network down"), want: "error"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyQuarkTVPollError(tt.err); got != tt.want {
				t.Fatalf("classifyQuarkTVPollError() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsQuarkTVAccessTokenInvalid(t *testing.T) {
	cases := []struct {
		name     string
		envelope quarkTVResponseEnvelope
		want     bool
	}{
		{
			name: "errno match",
			envelope: quarkTVResponseEnvelope{
				Status: -1,
				Errno:  10001,
			},
			want: true,
		},
		{
			name: "message match",
			envelope: quarkTVResponseEnvelope{
				ErrorInfo: "access token invalid",
			},
			want: true,
		},
		{
			name: "other error",
			envelope: quarkTVResponseEnvelope{
				Errno:     1,
				ErrorInfo: "other error",
			},
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := isQuarkTVAccessTokenInvalid(tt.envelope); got != tt.want {
				t.Fatalf("isQuarkTVAccessTokenInvalid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuarkTVAccessTokenCache(t *testing.T) {
	clearCachedQuarkTVAccessToken()
	cacheQuarkTVAccessToken("refresh-a", "access-a", 120)

	if token, ok := getCachedQuarkTVAccessToken("refresh-a"); !ok || token != "access-a" {
		t.Fatalf("expected cached token to be returned, got token=%q ok=%v", token, ok)
	}
	if _, ok := getCachedQuarkTVAccessToken("refresh-b"); ok {
		t.Fatalf("expected mismatched refresh token to miss cache")
	}

	clearCachedQuarkTVAccessToken()
	if _, ok := getCachedQuarkTVAccessToken("refresh-a"); ok {
		t.Fatalf("expected cache to be cleared")
	}
}

func TestQuarkTVConfigHelpers(t *testing.T) {
	oldDB := global.DB
	defer func() {
		global.DB = oldDB
	}()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&model.Config{}); err != nil {
		t.Fatalf("migrate config model: %v", err)
	}
	global.DB = db

	saved, err := upsertQuarkTVConfig(quarkTVDeviceIDConfigKey, "device-1")
	if err != nil {
		t.Fatalf("upsertQuarkTVConfig() error = %v", err)
	}
	if saved.Value != "device-1" {
		t.Fatalf("expected saved value device-1, got %q", saved.Value)
	}

	fetched, exists, err := getOptionalConfigByKey(quarkTVDeviceIDConfigKey)
	if err != nil {
		t.Fatalf("getOptionalConfigByKey() error = %v", err)
	}
	if !exists || fetched.Value != "device-1" {
		t.Fatalf("expected saved config to exist, got exists=%v value=%q", exists, fetched.Value)
	}

	if err := deleteQuarkTVConfigByKey(quarkTVDeviceIDConfigKey); err != nil {
		t.Fatalf("deleteQuarkTVConfigByKey() error = %v", err)
	}
	if _, exists, err := getOptionalConfigByKey(quarkTVDeviceIDConfigKey); err != nil || exists {
		t.Fatalf("expected config to be deleted, exists=%v err=%v", exists, err)
	}
}

func modelConfigWithUpdatedAt(value string, updatedAt time.Time) model.Config {
	return model.Config{
		Value: value,
		CommonModel: model.CommonModel{
			UpdatedAt: updatedAt,
		},
	}
}
