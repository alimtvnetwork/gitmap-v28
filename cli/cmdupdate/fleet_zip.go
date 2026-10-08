package cmdupdate

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func clearZipCache() {
	zipCacheMu.Lock()
	zipCache = make(map[string][]byte)
	zipCacheMu.Unlock()
}

func getCachedUpdateZip(pkg, osType string) ([]byte, error) {
	cacheKey := fmt.Sprintf("%s_%s", strings.ToLower(pkg), strings.ToLower(osType))
	zipCacheMu.Lock()
	defer zipCacheMu.Unlock()
	if data, ok := zipCache[cacheKey]; ok {
		return data, nil
	}
	data, err := CreateUpdateZipFn(pkg, osType)
	if err != nil {
		return nil, err
	}
	zipCache[cacheKey] = data
	return data, nil
}

func createUpdatePackageZip(pkg, osType string) ([]byte, error) {
	if data, ok := tryLoadAgmInstallerZip(pkg); ok {
		return data, nil
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	isWin := isWindowsOS(osType)

	binTarget := resolvePackageBinName(pkg, isWin)
	binData, err := locatePackageBinary(pkg)
	if err != nil {
		_ = zw.Close()
		return nil, apperror.WrapSimple(err, "locatePackageBinary")
	}
	if err := addZipFileEntry(zw, binTarget, binData, 0755); err != nil {
		_ = zw.Close()
		return nil, apperror.WrapSimple(err, "addZipFileEntry.bin")
	}
	appendAgmToZipIfIncluded(zw, pkg, isWin)
	if err := addLauncherScriptEntry(zw, isWin); err != nil {
		_ = zw.Close()
		return nil, apperror.WrapSimple(err, "addLauncherScriptEntry")
	}
	if err := zw.Close(); err != nil {
		return nil, apperror.WrapSimple(err, "zw.Close")
	}
	return buf.Bytes(), nil
}

func appendAgmToZipIfIncluded(zw *zip.Writer, pkg string, isWin bool) {
	if !isAgmIncludedInZip(pkg) {
		return
	}
	agmData, agmErr := locatePackageBinary("agm")
	if agmErr != nil || len(agmData) == 0 {
		return
	}
	agmTarget := resolvePackageBinName("agm", isWin)
	_ = addZipFileEntry(zw, agmTarget, agmData, 0755)
}

func resolvePackageBinName(pkg string, isWin bool) string {
	name := "gitmap"
	if isAgmPkg(pkg) {
		name = "agm"
	}
	if isWin {
		return name + ".exe"
	}
	return name
}

func isAgmIncludedInZip(pkg string) bool {
	return strings.EqualFold(pkg, "all")
}

func locatePackageBinary(pkg string) ([]byte, error) {
	if isAgmPkg(pkg) {
		return locateAgmBinary()
	}
	return locateGitmapBinary()
}

func readFileIfExists(path string) ([]byte, bool) {
	if path == "" {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return nil, false
	}
	return data, true
}

func locateGitmapBinary() ([]byte, error) {
	execPath, _ := os.Executable()
	if data, ok := readFileIfExists(execPath); ok {
		return data, nil
	}
	lp, _ := exec.LookPath("gitmap")
	if data, ok := readFileIfExists(lp); ok {
		return data, nil
	}
	return []byte("gitmap-payload-simulated"), nil
}

func agmExportZipPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".antigravity_tools", "update-export", "agm-update.zip")
}

func tryLoadAgmInstallerZip(pkg string) ([]byte, bool) {
	if !isAgmPkg(pkg) {
		return nil, false
	}
	data, err := loadOrExportAgmInstallerZip()
	if err == nil && len(data) > 4 && string(data[:2]) == "PK" {
		return data, true
	}

	return nil, false
}

func findAgmExecutable() (string, error) {
	if agm, err := exec.LookPath("agm"); err == nil {
		return agm, nil
	}

	return exec.LookPath("agm.exe")
}

func loadOrExportAgmInstallerZip() ([]byte, error) {
	path := agmExportZipPath()
	if data, ok := readFileIfExists(path); ok {
		return data, nil
	}
	agm, err := findAgmExecutable()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, agm, "update", "export-zip")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	data, ok := readFileIfExists(path)
	if !ok {
		return nil, fmt.Errorf("agm update export-zip did not write %s", path)
	}
	return data, nil
}

func locateAgmBinary() ([]byte, error) {
	lp, _ := exec.LookPath("agm")
	if data, ok := readFileIfExists(lp); ok {
		return data, nil
	}
	lpExe, _ := exec.LookPath("agm.exe")
	if data, ok := readFileIfExists(lpExe); ok {
		return data, nil
	}
	return []byte("agm-payload-simulated"), nil
}

func addZipFileEntry(zw *zip.Writer, name string, data []byte, mode os.FileMode) error {
	header := &zip.FileHeader{
		Name:   name,
		Method: zip.Deflate,
	}
	header.SetMode(mode)
	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func addLauncherScriptEntry(zw *zip.Writer, isWin bool) error {
	if isWin {
		script := "# Embedded GitMap Remote Installer\nExpand-Archive -Path $zipPath -DestinationPath $destDir -Force\n"
		return addZipFileEntry(zw, "install_remote.ps1", []byte(script), 0644)
	}
	script := "#!/bin/sh\nunzip -o \"$ZIP\" -d \"$DEST\"\n"
	return addZipFileEntry(zw, "install_remote.sh", []byte(script), 0755)
}
