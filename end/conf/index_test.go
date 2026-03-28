package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfigEnvOverrideDevelopMode(t *testing.T) {
	baseDir := t.TempDir()
	configDir := filepath.Join(baseDir, "conf")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	config := "mode:\n  develop: true\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv(BaseDirEnv, baseDir)
	t.Setenv("MODE_DEVELOP", "false")
	t.Cleanup(func() {
		appBaseDir = ""
		viper.Reset()
	})

	InitConfig()

	if viper.GetBool("mode.develop") {
		t.Fatal("expected MODE_DEVELOP env to override config file to false")
	}
}
