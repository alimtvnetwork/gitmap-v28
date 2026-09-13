package cmdinstall

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

const (
	maxDownloadAttempts     = 3
	antigravityWindowsUrl   = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/windows-x64/Antigravity-x64.exe"
	antigravityLinuxUrl     = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-x64/Antigravity.tar.gz"
	antigravityLinuxArmUrl  = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-arm/Antigravity.tar.gz"
	antigravityDarwinUrl    = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-x64/Antigravity.dmg"
	antigravityDarwinArmUrl = "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-arm/Antigravity.dmg"
)

func resolveBuildVersion(version string) string {
	if version != "" {
		return version
	}
	return AntigravityDefaultVersion
}

func resolveBuildID(buildID string) string {
	if buildID != "" {
		return buildID
	}
	return AntigravityDefaultBuildID
}

func buildAntigravityURL(version, buildID, platform, artifact string) string {
	v := resolveBuildVersion(version)
	b := resolveBuildID(buildID)
	return fmt.Sprintf("%s/%s-%s/%s/%s", AntigravityBaseURL, v, b, platform, artifact)
}

func resolveDownloadArch() string {
	if runtime.GOARCH == "arm64" || runtime.GOARCH == "arm" {
		return "arm64"
	}
	return "x64"
}

func getAntigravityDesktopDownloadUrl(osName string) string {
	arch := resolveDownloadArch()
	platformID, artifact, appErr := resolvePlatformArtifact(osName, arch)
	if appErr != nil {
		return ""
	}
	return buildAntigravityURL(AntigravityDefaultVersion, AntigravityDefaultBuildID, platformID, artifact)
}

func resolveDownloadDest(url, destPath string) string {
	if destPath != "" {
		return destPath
	}
	targetDir := tempdir.RepoTempDir("antigravity-install")
	return filepath.Join(targetDir, filepath.Base(url))
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

func executeSingleDownload(client *http.Client, url, destPath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}
	return copyDownloadStream(resp, destPath)
}

func hasValidDownloadSize(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Size() > 0
}

func verifyDownloadedFile(targetFile string, attempt int) bool {
	if !hasValidDownloadSize(targetFile) {
		fmt.Printf("[WARN ] Attempt %d produced empty file, retrying\n", attempt)
		_ = os.Remove(targetFile)
		return false
	}
	return true
}

func attemptDownloadStep(client *http.Client, url, targetFile string, attempt int) bool {
	fmt.Printf("[STEP ] Downloading Antigravity (attempt %d/%d)...\n", attempt, maxDownloadAttempts)
	err := executeSingleDownload(client, url, targetFile)
	if err != nil {
		fmt.Printf("[WARN ] Attempt %d failed: %v\n", attempt, err)
		return false
	}
	return verifyDownloadedFile(targetFile, attempt)
}

func applyDownloadBackoff(attempt int) {
	if attempt < maxDownloadAttempts {
		time.Sleep(time.Duration(attempt) * time.Second)
	}
}

func createDownloadFailedError(url string) *apperror.AppError {
	return apperror.NewWithDetails(
		"downloadFileWithRetry",
		ErrDownloadFailed,
		fmt.Sprintf("failed to download from %s after %d attempts", url, maxDownloadAttempts),
		"installer",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		map[string]any{"url": url, "attempts": maxDownloadAttempts},
	)
}

func downloadFileWithRetry(url, destPath string) *apperror.AppError {
	targetFile := resolveDownloadDest(url, destPath)
	client := &http.Client{Timeout: 600 * time.Second}
	for attempt := 1; attempt <= maxDownloadAttempts; attempt++ {
		if isSuccess := attemptDownloadStep(client, url, targetFile, attempt); isSuccess {
			fmt.Printf("[OK   ] Download complete: %s\n", targetFile)
			return nil
		}
		applyDownloadBackoff(attempt)
	}
	return createDownloadFailedError(url)
}

func downloadFileToDest(url, destPath string) error {
	appErr := downloadFileWithRetry(url, destPath)
	if appErr != nil {
		return appErr
	}
	return nil
}

func recordAntigravityDesktopInstalled(installPath string) {
	recordToolInDatabases("antigravity", installPath, "installer")
	recordToolInDatabases("agy", installPath, "installer")
}
