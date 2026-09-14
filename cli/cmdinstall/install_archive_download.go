package cmdinstall

import (
	"bytes"
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
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
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
	return fallbackArchiveFilename(base)
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
	fallback := tempdir.RepoTempDir("downloads")
	_ = os.MkdirAll(fallback, 0755)
	return fallback
}

func isCompressedArchiveMagic(h []byte) bool {
	if len(h) < 4 {
		return false
	}
	isGzip := h[0] == 0x1f && h[1] == 0x8b
	isZip := h[0] == 0x50 && h[1] == 0x4b
	isBz2 := h[0] == 0x42 && h[1] == 0x5a && h[2] == 0x68
	isZst := bytes.Equal(h[:4], []byte{0x28, 0xb5, 0x2f, 0xfd})
	return isGzip || isZip || isBz2 || isZst
}

func isContainerOrExecutableMagic(h []byte) bool {
	if len(h) < 6 {
		return false
	}
	isXz := bytes.Equal(h[:6], []byte{0xfd, 0x37, 0x7a, 0x58, 0x5a, 0x00})
	is7z := bytes.Equal(h[:6], []byte{0x37, 0x7a, 0xbc, 0xaf, 0x27, 0x1c})
	isExe := h[0] == 'M' && h[1] == 'Z'
	isElf := bytes.Equal(h[:4], []byte{0x7f, 'E', 'L', 'F'})
	return isXz || is7z || isExe || isElf
}

func isDmgOrAppleMagic(h []byte) bool {
	if len(h) < 4 {
		return false
	}
	isMachO := (h[0] == 0xfe && h[1] == 0xed && h[2] == 0xfa && (h[3] == 0xce || h[3] == 0xcf)) ||
		(h[0] == 0xcf && h[1] == 0xfa && h[2] == 0xed && h[3] == 0xfe)
	isAppleXar := bytes.Equal(h[:4], []byte{0x78, 0x61, 0x72, 0x21})
	isZlib := h[0] == 0x78 && (h[1] == 0x01 || h[1] == 0x9c || h[1] == 0xda)
	return isMachO || isAppleXar || isZlib
}

func isTarArchiveMagic(h []byte) bool {
	if len(h) < 262 {
		return false
	}
	return bytes.HasPrefix(h[257:], []byte("ustar"))
}

func isSupportedArchiveHeader(h []byte) bool {
	return isCompressedArchiveMagic(h) || isContainerOrExecutableMagic(h) || isTarArchiveMagic(h) || isDmgOrAppleMagic(h)
}

func readArchiveHeaderBytes(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "archive.open")
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, readErr := io.ReadFull(f, buf)
	if readErr != nil && readErr != io.ErrUnexpectedEOF && readErr != io.EOF {
		return nil, apperror.WrapSimple(readErr, "archive.readHeader")
	}
	return buf[:n], nil
}

func isDmgKolyTrailer(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || stat.Size() < 512 {
		return false
	}
	trailer := make([]byte, 512)
	_, readErr := f.ReadAt(trailer, stat.Size()-512)
	return readErr == nil && bytes.Equal(trailer[:4], []byte("koly"))
}

func VerifyArchiveHeaderMagic(filePath string) bool {
	header, err := readArchiveHeaderBytes(filePath)
	if err != nil || len(header) < 4 {
		return false
	}
	if isSupportedArchiveHeader(header) {
		return true
	}
	return isDmgKolyTrailer(filePath)
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

