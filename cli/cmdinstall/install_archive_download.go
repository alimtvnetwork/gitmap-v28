package cmdinstall

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fspath"
)

func isRemoteURL(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	return strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
}

func checkUrlScheme(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return apperror.NewSimple("unsupported URL scheme: "+u.Scheme, "E_INVALID_URL")
	}
	if u.Host == "" {
		return apperror.NewSimple("URL host cannot be empty", "E_INVALID_URL")
	}
	return nil
}

func validateRemoteArchiveURL(rawURL string) error {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return apperror.NewSimple("archive URL cannot be empty", "E_EMPTY_URL")
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return apperror.WrapSimple(err, "archive.parseUrl")
	}
	return checkUrlScheme(parsed)
}

func fallbackArchiveFilename(name string) string {
	if name == "" || name == "." || name == "/" {
		return "archive-download.tar.gz"
	}
	return name
}

func filenameFromArchiveURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "archive-download.tar.gz"
	}
	base := path.Base(parsed.Path)
	if name, ok := antigravityHubFilename(parsed.Path, base); ok {
		return name
	}
	return fallbackArchiveFilename(base)
}

func antigravityHubFilename(urlPath, base string) (string, bool) {
	if !strings.Contains(urlPath, "antigravity-hub/") {
		return "", false
	}
	parts := strings.Split(urlPath, "antigravity-hub/")
	if len(parts) <= 1 {
		return "", false
	}
	subParts := strings.Split(parts[1], "/")
	if len(subParts) <= 1 {
		return "", false
	}
	tag := subParts[0]
	ext := archiveExtension(base)
	prefix := strings.TrimSuffix(base, ext)
	return fallbackArchiveFilename(fmt.Sprintf("%s-%s%s", prefix, tag, ext)), true
}

func archiveExtension(base string) string {
	if strings.HasSuffix(base, ".tar.gz") {
		return ".tar.gz"
	}
	return path.Ext(base)
}

func getGitmapHomeDownloadDir() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", false
	}
	dir := filepath.Join(home, ".gitmap", "downloads")
	return dir, os.MkdirAll(dir, 0755) == nil
}

func resolveArchiveDownloadCacheDir() string {
	if dir, hasHomeDir := getGitmapHomeDownloadDir(); hasHomeDir {
		return dir
	}
	fallback := fspath.RepoTempDir("downloads")
	_ = os.MkdirAll(fallback, 0755)
	return fallback
}

func validateFileSize(size int64) (bool, string) {
	if size <= 1024 {
		return false, "cached archive is smaller than 1KB minimum threshold"
	}
	return true, ""
}

func CheckCachedArchive(filePath string) ArchiveCacheValidationResult {
	info, err := os.Stat(filePath)
	if err != nil {
		return ArchiveCacheValidationResult{IsValid: false, ErrorMessage: "cached archive does not exist"}
	}
	if hasMinSize, sizeErr := validateFileSize(info.Size()); !hasMinSize {
		return ArchiveCacheValidationResult{Size: info.Size(), ErrorMessage: sizeErr}
	}
	if hasMagic := VerifyArchiveHeaderMagic(filePath); !hasMagic {
		return ArchiveCacheValidationResult{Size: info.Size(), ErrorMessage: "invalid archive magic header"}
	}
	return ArchiveCacheValidationResult{IsValid: true, Size: info.Size(), HasValidMagic: true}
}

func writeHttpStreamToFile(r io.Reader, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return apperror.WrapSimple(err, "file.create")
	}
	defer out.Close()
	if _, copyErr := io.Copy(out, r); copyErr != nil {
		return apperror.WrapSimple(copyErr, "file.copy")
	}
	return nil
}

func downloadViaGoHttp(rawURL, destPath string) error {
	client := &http.Client{Timeout: 600 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return apperror.WrapSimple(err, "http.get")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return apperror.NewSimple(fmt.Sprintf("HTTP download status: %d", resp.StatusCode), "E_HTTP_STATUS")
	}
	return writeHttpStreamToFile(resp.Body, destPath)
}

func executeAria2cDownload(rawURL, destPath string) error {
	ariaPath, err := exec.LookPath("aria2c")
	if err != nil {
		return apperror.WrapSimple(err, "aria2c.lookup")
	}
	dir := filepath.Dir(destPath)
	file := filepath.Base(destPath)
	cmd := exec.Command(ariaPath, "-s", "16", "-x", "16", "-k", "1M", "--allow-overwrite=true", "--dir", dir, "--out", file, rawURL)
	if runErr := cmd.Run(); runErr != nil {
		return apperror.WrapSimple(runErr, "aria2c.exec")
	}
	return nil
}

func executeCurlDownload(rawURL, destPath string) error {
	curlPath, err := exec.LookPath("curl")
	if err != nil {
		return apperror.WrapSimple(err, "curl.lookup")
	}
	cmd := exec.Command(curlPath, "-fL", "--retry", "2", "-o", destPath, rawURL)
	if runErr := cmd.Run(); runErr != nil {
		return apperror.WrapSimple(runErr, "curl.exec")
	}
	return nil
}

func tryAriaOrCurlDownload(rawURL, destPath string) bool {
	if err := executeAria2cDownload(rawURL, destPath); err == nil {
		return true
	}
	return executeCurlDownload(rawURL, destPath) == nil
}

func ensureDownloadParentDir(destPath string) error {
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.downloadDir")
	}
	return nil
}

func downloadMultiTier(rawURL, destPath string) error {
	if hasFastDownload := tryAriaOrCurlDownload(rawURL, destPath); hasFastDownload {
		return nil
	}
	return downloadViaGoHttp(rawURL, destPath)
}

func resolveArchiveDownloadPath(params ArchiveDownloadParams) string {
	if params.DestPath != "" {
		return params.DestPath
	}
	cacheDir := resolveArchiveDownloadCacheDir()
	fileName := filenameFromArchiveURL(params.URL)
	return filepath.Join(cacheDir, fileName)
}

func tryReuseCachedArchive(cachedPath string) (bool, string) {
	cacheVal := CheckCachedArchive(cachedPath)
	if cacheVal.IsValid {
		fmt.Printf("Reusing valid download from cache: %s\n", cachedPath)
		return true, cachedPath
	}
	return false, ""
}

func purgeArchiveIfForced(cachedPath string, isDownloadMust bool) {
	if isDownloadMust {
		_ = os.Remove(cachedPath)
	}
}

func verifyDownloadedArchive(cachedPath string) error {
	val := CheckCachedArchive(cachedPath)
	if !val.IsValid {
		_ = os.Remove(cachedPath)
		return apperror.NewSimple("download validation failed: "+val.ErrorMessage, "E_INVALID_ARCHIVE")
	}
	return nil
}

func executeArchiveDownloadFlow(params ArchiveDownloadParams, targetPath string) (string, error) {
	if err := ensureDownloadParentDir(targetPath); err != nil {
		return "", err
	}
	if err := downloadMultiTier(params.URL, targetPath); err != nil {
		return "", err
	}
	if err := verifyDownloadedArchive(targetPath); err != nil {
		return "", err
	}
	return targetPath, nil
}

func checkEligibleCachedArchive(targetPath string, isDownloadMust bool) (bool, string) {
	if isDownloadMust {
		return false, ""
	}
	return tryReuseCachedArchive(targetPath)
}

func FetchOrReuseArchive(params ArchiveDownloadParams) (string, error) {
	if err := validateRemoteArchiveURL(params.URL); err != nil {
		return "", err
	}
	targetPath := resolveArchiveDownloadPath(params)
	if hasCached, cached := checkEligibleCachedArchive(targetPath, params.IsDownloadMust); hasCached {
		return cached, nil
	}
	purgeArchiveIfForced(targetPath, params.IsDownloadMust)
	return executeArchiveDownloadFlow(params, targetPath)
}
