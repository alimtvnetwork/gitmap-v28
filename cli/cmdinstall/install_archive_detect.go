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

func cleanArchiveExtPath(path string) string {
	low := strings.ToLower(path)
	if idx := strings.Index(low, "?"); idx != -1 {
		return low[:idx]
	}
	return low
}

func isArchiveExtension(path string) bool {
	cleaned := cleanArchiveExtPath(path)
	for _, ext := range archiveExts {
		if strings.HasSuffix(cleaned, ext) {
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

func matchDesktopOrIcon(p string, desktop, icon *string) {
	if *desktop == "" && strings.HasSuffix(p, ".desktop") {
		*desktop = p
	}
	if *icon == "" && (strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".svg")) {
		*icon = p
	}
}

func findDesktopAndIcon(dir string) (string, string) {
	var desktopPath, iconPath string
	_ = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			matchDesktopOrIcon(p, &desktopPath, &iconPath)
		}
		return nil
	})
	return desktopPath, iconPath
}

func evaluateBinaryMatch(p string, info fs.FileInfo, appName string, best *string) error {
	if isCandidateMatch(p, info, appName) {
		*best = p
		return filepath.SkipAll
	}
	if *best == "" && (isElfBinary(p) || isExecutableFile(info)) {
		*best = p
	}
	return nil
}

func findCandidateBinary(dir, appName string) string {
	var bestCandidate string
	_ = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Size() > 0 {
			return evaluateBinaryMatch(p, info, appName, &bestCandidate)
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

func resolvePackageStrategy(res ArchiveInspectionResult) ArchiveInstallStrategyType {
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
