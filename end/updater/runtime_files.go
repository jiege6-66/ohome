package updater

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadArtifact(tempDir string, taskID string, url string) (string, error) {
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		return "", err
	}
	tmpFile, err := os.CreateTemp(tempDir, taskID+"-*.tar.gz")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	resp, err := http.Get(strings.TrimSpace(url))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("下载更新包失败: %s", resp.Status)
	}
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func verifySHA256File(path string, expected string) error {
	expected = strings.ToLower(strings.TrimSpace(expected))
	if expected == "" {
		return fmt.Errorf("缺少更新包 SHA256")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	if actual != expected {
		return fmt.Errorf("更新包校验失败: 期望 %s，实际 %s", expected, actual)
	}
	return nil
}

func extractServerArchive(archivePath string, releaseDir string) error {
	if err := os.RemoveAll(releaseDir); err != nil {
		return err
	}
	if err := os.MkdirAll(releaseDir, 0o755); err != nil {
		return err
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	serverPath := filepath.Join(releaseDir, "server")
	foundServer := false
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if filepath.Base(header.Name) != "server" {
				continue
			}
			out, err := os.OpenFile(serverPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tarReader); err != nil {
				out.Close()
				return err
			}
			if err := out.Close(); err != nil {
				return err
			}
			foundServer = true
		}
	}
	if !foundServer {
		return fmt.Errorf("更新包中缺少 server 可执行文件")
	}
	return ensureExecutable(serverPath)
}

func ensureExecutable(path string) error {
	return os.Chmod(path, 0o755)
}

func copyExecutable(src string, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	output, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		return err
	}
	if err := output.Close(); err != nil {
		return err
	}
	return ensureExecutable(dst)
}

func resolveCurrentReleasePath(linkPath string) (string, error) {
	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func setCurrentReleaseLink(linkPath string, target string) error {
	if strings.TrimSpace(target) == "" {
		return fmt.Errorf("目标版本目录不能为空")
	}
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		return err
	}
	tmpPath := linkPath + ".tmp"
	_ = os.Remove(tmpPath)
	if err := os.Symlink(filepath.Clean(target), tmpPath); err != nil {
		return err
	}
	if err := removeLinkOrFile(linkPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, linkPath)
}

func rollbackCurrentRelease(linkPath string, previousReleasePath string) error {
	if strings.TrimSpace(previousReleasePath) == "" {
		return fmt.Errorf("缺少上一版本目录")
	}
	return setCurrentReleaseLink(linkPath, previousReleasePath)
}

func removeLinkOrFile(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s 不是符号链接，拒绝覆盖目录", path)
	}
	return os.Remove(path)
}
