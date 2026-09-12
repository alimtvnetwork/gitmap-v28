package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

const (
	antigravityWindowsUrl = "https://antigravity.google/download/windows-x64/Antigravity-x64.exe"
	antigravityLinuxUrl   = "https://antigravity.google/download/linux-x64/Antigravity.tar.gz"
)

func getAntigravityDesktopDownloadUrl(osName string) string {
	if osName == "windows" {

		return antigravityWindowsUrl
	}

	return antigravityLinuxUrl
}

func copyDownloadStream(resp *http.Response, destFile string) error {
	out, err := os.Create(destFile)
	if err != nil {

		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}

func downloadFileToDest(url, destPath string) error {
	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Get(url)
	if err != nil {

		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("download failed with HTTP %d from %s", resp.StatusCode, url)
	}

	return copyDownloadStream(resp, destPath)
}

func recordAntigravityDesktopInstalled(installPath string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {

		return
	}
	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("antigravity", installPath, "installer")
}
