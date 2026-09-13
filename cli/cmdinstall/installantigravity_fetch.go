package cmdinstall

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	antigravityWindowsUrl   = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/windows-x64/Antigravity-x64.exe"
	antigravityLinuxUrl     = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-x64/Antigravity.tar.gz"
	antigravityLinuxArmUrl  = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-arm/Antigravity.tar.gz"
	antigravityDarwinUrl    = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-x64/Antigravity.dmg"
	antigravityDarwinArmUrl = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-arm/Antigravity.dmg"
)

func getAntigravityDesktopDownloadUrl(osName string) string {
	if osName == "windows" {
		return antigravityWindowsUrl
	}
	if osName == "darwin" && runtime.GOARCH == "arm64" {
		return antigravityDarwinArmUrl
	}
	if osName == "darwin" {
		return antigravityDarwinUrl
	}
	if runtime.GOARCH == "arm64" || runtime.GOARCH == "arm" {
		return antigravityLinuxArmUrl
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
	fmt.Printf("Downloading Antigravity from Google Cloud Storage...\n")
	client := &http.Client{Timeout: 600 * time.Second}
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
	recordToolInDatabases("antigravity", installPath, "installer")
	recordToolInDatabases("agy", installPath, "installer")
}
