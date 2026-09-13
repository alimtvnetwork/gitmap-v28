package cmdinstall

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var archiveExts = []string{
	".tar.gz", ".tar.xz", ".tar.bz2", ".tar.zst",
	".tgz", ".txz", ".tbz2", ".tzst",
	".tar", ".zip", ".gz",
}

func isArchiveExtension(path string) bool {
	low := strings.ToLower(path)
	for _, ext := range archiveExts {
		if strings.HasSuffix(low, ext) {
			return true
		}
	}
	return false
}

func isElfBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	header := make([]byte, 4)
	_, readErr := io.ReadFull(f, header)
	if readErr != nil {
		return false
	}
	return bytes.Equal(header, []byte{0x7f, 'E', 'L', 'F'})
}

func isExecutableFile(info fs.FileInfo) bool {
	if info.IsDir() {
		return false
	}
	return info.Mode()&0111 != 0
}

func findScriptInDir(dir string) string {
	candidates := []string{"install.sh", "setup.sh", "install", "setup"}
	for _, name := range candidates {
		p := filepath.Join(dir, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

func findMakefileInDir(dir string) (bool, bool) {
	_, errMake := os.Stat(filepath.Join(dir, "Makefile"))
	hasMake := errMake == nil
	_, errConf := os.Stat(filepath.Join(dir, "configure"))
	hasConf := errConf == nil
	return hasMake, hasConf
}

func findDesktopAndIcon(dir string) (string, string) {
	var desktopPath, iconPath string
	_ = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if desktopPath == "" && strings.HasSuffix(p, ".desktop") {
			desktopPath = p
		}
		if iconPath == "" && (strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".svg")) {
			iconPath = p
		}
		return nil
	})
	return desktopPath, iconPath
}

func findCandidateBinary(dir, appName string) string {
	var bestCandidate string
	_ = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Size() == 0 {
			return nil
		}
		if isCandidateMatch(p, info, appName) {
			bestCandidate = p
			return filepath.SkipAll
		}
		if bestCandidate == "" && (isElfBinary(p) || isExecutableFile(info)) {
			bestCandidate = p
		}
		return nil
	})
	return bestCandidate
}

func isCandidateMatch(p string, info fs.FileInfo, appName string) bool {
	if !isElfBinary(p) && !isExecutableFile(info) {
		return false
	}
	name := strings.ToLower(filepath.Base(p))
	return name == strings.ToLower(appName)
}

func inspectExtractedPackage(extractDir, baseName string) ArchiveInspectionResult {
	res := ArchiveInspectionResult{SourceRoot: extractDir, SuggestedName: baseName}
	res.DesktopPath, res.IconPath = findDesktopAndIcon(extractDir)
	res.HasDesktopFile = res.DesktopPath != ""
	res.HasMakefile, res.HasConfigure = findMakefileInDir(extractDir)
	res.ScriptPath = findScriptInDir(extractDir)
	res.BinaryPath = findCandidateBinary(extractDir, baseName)
	res.Strategy = resolvePackageStrategy(res)
	return res
}

func resolvePackageStrategy(res ArchiveInspectionResult) ArchiveInstallStrategy {
	if res.BinaryPath != "" {
		return StrategyBinaryApp
	}
	if res.ScriptPath != "" {
		return StrategyScript
	}
	if res.HasMakefile || res.HasConfigure {
		return StrategySource
	}
	return StrategyUnknown
}
