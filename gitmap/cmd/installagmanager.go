package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type agManagerAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type agManagerRelease struct {
	Assets []agManagerAsset `json:"assets"`
}

func runInstallAgManager() {
	fmt.Println("Fetching latest release for Antigravity-Manager...")
	assetURL, err := getAgManagerAssetURL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching release: %v\n", err)
		return
	}
	fmt.Printf("Downloading %s...\n", assetURL)
	performAgManagerDownloadAndInstall(assetURL)
}

func performAgManagerDownloadAndInstall(assetURL string) {
	tmpPath, err := downloadAgManagerFile(assetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err)
		return
	}
	fmt.Printf("Installing %s...\n", filepath.Base(tmpPath))
	if err := executeAgManagerInstaller(tmpPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing: %v\n", err)
		return
	}
	fmt.Println(constants.ColorGreen + "✓" + constants.ColorReset + " Antigravity Manager installed successfully.")
}

func getAgManagerAssetURL() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/lbjlaq/Antigravity-Manager/releases/latest")
	if err != nil { return "", err }
	defer resp.Body.Close()
	var release agManagerRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return matchAgManagerAsset(release.Assets)
}

func matchAgManagerAsset(assets []agManagerAsset) (string, error) {
	osStr, archStr := runtime.GOOS, runtime.GOARCH
	for _, asset := range assets {
		if strings.HasSuffix(asset.Name, ".sig") || strings.HasSuffix(asset.Name, "updater.json") {
			continue
		}
		if matchAssetOS(asset.Name, osStr, archStr) {
			return asset.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("no suitable asset found for %s %s", osStr, archStr)
}

func matchAssetOS(name, osStr, archStr string) bool {
	n := strings.ToLower(name)
	switch osStr {
	case "windows":
		return (strings.HasSuffix(n, ".exe") || strings.HasSuffix(n, ".msi")) && matchArch(n, archStr)
	case "darwin":
		return strings.HasSuffix(n, ".dmg") && matchArch(n, archStr)
	case "linux":
		return (strings.HasSuffix(n, ".deb") || strings.HasSuffix(n, ".appimage")) && matchArch(n, archStr)
	}
	return false
}

func matchArch(n, archStr string) bool {
	if archStr == "arm64" {
		return strings.Contains(n, "aarch64") || strings.Contains(n, "arm64")
	}
	if archStr == "amd64" {
		return strings.Contains(n, "x64") || strings.Contains(n, "amd64") || strings.Contains(n, "x86_64")
	}
	return false
}

func downloadAgManagerFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil { return "", err }
	defer resp.Body.Close()
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	tmpPath := filepath.Join(os.TempDir(), name)
	out, err := os.Create(tmpPath)
	if err != nil { return "", err }
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return tmpPath, err
}

func executeAgManagerInstaller(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		if strings.HasSuffix(path, ".msi") { cmd = exec.Command("msiexec", "/i", path, "/qn") } else { cmd = exec.Command(path, "/S") }
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		if strings.HasSuffix(strings.ToLower(path), ".deb") { cmd = exec.Command("sudo", "dpkg", "-i", path) } else { os.Chmod(path, 0755); cmd = exec.Command(path) }
	}
	if cmd != nil {
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}
	return nil
}
