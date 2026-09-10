package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func isHtmlContent(body []byte) bool {
	str := strings.ToLower(strings.TrimSpace(string(body)))
	if strings.HasPrefix(str, "<!") || strings.HasPrefix(str, "<html") {

		return true
	}

	return false
}

func fetchAgyScriptBytes() ([]byte, error) {
	url := "https://antigravity.google/cli/install.sh"
	if runtime.GOOS == "windows" {
		url = "https://antigravity.google/cli/install.ps1"
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("HTTP %d from installer endpoint", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func saveTempAgyScript(body []byte) (string, error) {
	scriptName := "install-agy.sh"
	if runtime.GOOS == "windows" {
		scriptName = "install-agy.ps1"
	}
	tmp := filepath.Join(os.TempDir(), scriptName)
	if err := os.WriteFile(tmp, body, 0755); err != nil {

		return "", err
	}

	return tmp, nil
}

func downloadAndValidateAgyScript() (string, error) {
	body, err := fetchAgyScriptBytes()
	if err != nil {

		return "", err
	}
	if isHtmlContent(body) {

		return "", fmt.Errorf("installer endpoint returned HTML document instead of shell script")
	}

	return saveTempAgyScript(body)
}

func runAgyNpmFallback() error {
	cmd := exec.Command("npm", "install", "-g", "@google/antigravity")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	return cmd.Run()
}

func recordAgyInstalled(ver string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {

		return
	}
	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("antigravity", ver, "installer")
}
