package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ohome/buildinfo"
	"ohome/conf"
	"ohome/global"
	"ohome/updater"
)

func main() {
	conf.InitConfig()
	global.Logger = conf.InitLogger()

	store := updater.NewStore()
	embeddedServerPath, embeddedServerVersionPath, err := bundledServerPaths()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	supervisor := updater.NewRuntimeSupervisor(store, conf.AppBaseDir(), embeddedServerPath)
	supervisor.SetUnexpectedExitHook(func(err error) {
		if global.Logger != nil {
			global.Logger.Errorf("server exited unexpectedly: %v", err)
		}
		os.Exit(1)
	})

	currentVersion, currentReleasePath, err := supervisor.BootstrapCurrentRelease(readEmbeddedServerVersion(embeddedServerVersionPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := supervisor.StartCurrentServer(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := ensureInitialState(store, currentVersion, currentReleasePath); err != nil && global.Logger != nil {
		global.Logger.Warnf("bootstrap updater state failed: %v", err)
	}

	manager := updater.NewManager(supervisor)
	apiServer := updater.NewAPIServerWithManager(manager)
	apiErrCh := make(chan error, 1)
	go func() {
		apiErrCh <- apiServer.Run()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case <-ctx.Done():
	case err := <-apiErrCh:
		if err != nil && global.Logger != nil {
			global.Logger.Errorf("launcher api server exited: %v", err)
		}
		_ = supervisor.StopCurrentServer(10 * time.Second)
		if err != nil {
			os.Exit(1)
		}
		return
	}
	_ = supervisor.StopCurrentServer(10 * time.Second)
}

func ensureInitialState(store *updater.Store, currentVersion string, currentReleasePath string) error {
	state, err := store.LoadState()
	if err != nil {
		return err
	}
	if strings.TrimSpace(state.CurrentVersion) == currentVersion &&
		strings.TrimSpace(state.CurrentReleasePath) == filepath.Clean(currentReleasePath) &&
		strings.TrimSpace(state.RuntimeVersion) == buildinfo.CleanRuntimeVersion() &&
		state.DeployMode == updater.DetectDeployMode() {
		return nil
	}
	state.CurrentVersion = strings.TrimSpace(currentVersion)
	state.CurrentReleasePath = filepath.Clean(currentReleasePath)
	state.RuntimeVersion = buildinfo.CleanRuntimeVersion()
	state.DeployMode = updater.DetectDeployMode()
	return store.SaveState(state)
}

func bundledServerPaths() (string, string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	baseDir := filepath.Dir(execPath)
	return filepath.Join(baseDir, "bin", updater.ServerExecutableName()),
		filepath.Join(baseDir, "bin", "server.version"),
		nil
}

func readEmbeddedServerVersion(path string) string {
	payload, err := os.ReadFile(path)
	if err != nil {
		return buildinfo.CleanVersion()
	}
	version := strings.TrimSpace(string(payload))
	if version == "" {
		return buildinfo.CleanVersion()
	}
	return version
}
