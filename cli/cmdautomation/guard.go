package cmdautomation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunFileSizeGuard scans files, auditing sizes, common binaries, and large JSONs.
func RunFileSizeGuard(opts GuardOptions) (GuardResult, *apperror.AppError) {
	start := time.Now()
	res := initGuardResult()
	targetDir := resolveGuardTargetDir(opts.Dir)
	exclusions := loadActiveExclusions()

	walkErr := filepath.Walk(targetDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && isExcludedDir(info.Name()) {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			auditSingleFile(p, info, opts, exclusions, &res)
		}
		return nil
	})

	if walkErr != nil {
		ctx := map[string]any{"dir": targetDir, "err": walkErr.Error()}
		return res, apperror.New("file_size_guard", "E_WALK_FAILED", ctx)
	}

	res.Duration = time.Since(start)
	return res, nil
}

func initGuardResult() GuardResult {
	return GuardResult{
		OversizedFiles:     []OversizedFile{},
		BinaryFiles:        []OversizedFile{},
		ExcludedLargeJsons: []OversizedFile{},
	}
}

func resolveGuardTargetDir(dir string) string {
	if strings.TrimSpace(dir) == "" {
		return "."
	}
	return dir
}

func loadActiveExclusions() []string {
	entries, err := ListExclusions()
	if err != nil || len(entries) == 0 {
		return nil
	}
	patterns := make([]string, len(entries))
	for i, e := range entries {
		patterns[i] = e.Pattern
	}
	return patterns
}

func auditSingleFile(p string, info os.FileInfo, opts GuardOptions, exclusions []string, res *GuardResult) {
	res.TotalFiles++
	rel := filepath.ToSlash(p)

	if IsPathExcluded(rel, exclusions) {
		res.ExcludedCount++
		return
	}

	sz := info.Size()
	if checkAndHandleBinary(rel, sz, opts, res) {
		return
	}
	if checkAndHandleLargeJson(rel, sz, opts, res) {
		return
	}
	checkAndHandleOversized(rel, sz, opts, res)
}

func checkAndHandleBinary(rel string, sz int64, opts GuardOptions, res *GuardResult) bool {
	ext := filepath.Ext(rel)
	isBin := IsBinaryExtension(ext)
	if !isBin {
		data, err := os.ReadFile(rel)
		isBin = (err == nil && HasBinaryContent(data))
	}
	if !isBin {
		return false
	}

	fileItem := OversizedFile{Path: rel, SizeBytes: sz, Kind: "binary"}
	res.BinaryFiles = append(res.BinaryFiles, fileItem)

	if isPromptBinaryExclude(opts) {
		promptAndPersistBinaryExclusion(rel, sz, res)
	}
	return true
}

func isPromptBinaryExclude(opts GuardOptions) bool {
	if opts.AutoExclude {
		return true
	}
	if !opts.Interactive {
		return false
	}
	return isInteractiveStdin()
}

func promptAndPersistBinaryExclusion(rel string, sz int64, res *GuardResult) {
	kb := sz / 1024
	fmt.Printf("\n%s[Binary File Detected]%s %s (%d KB)\n", constants.ColorYellow, constants.ColorReset, rel, kb)
	fmt.Print("  Exclude this file from future searches? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	if ans == "y" || ans == "yes" {
		_ = AddExclusion(rel, "binary")
		res.ExcludedCount++
		fmt.Printf("  %s✔ Excluded from future searches.%s\n", constants.ColorGreen, constants.ColorReset)
	}
}

func checkAndHandleLargeJson(rel string, sz int64, opts GuardOptions, res *GuardResult) bool {
	if filepath.Ext(rel) != ".json" {
		return false
	}
	maxKb := opts.MaxJsonKb
	if maxKb <= 0 {
		maxKb = 500
	}
	limitBytes := int64(maxKb * 1024)
	if sz <= limitBytes {
		return false
	}

	item := OversizedFile{Path: rel, SizeBytes: sz, MaxAllowedBytes: limitBytes, Kind: "large_json"}
	res.ExcludedLargeJsons = append(res.ExcludedLargeJsons, item)
	res.ExcludedCount++
	return true
}

func checkAndHandleOversized(rel string, sz int64, opts GuardOptions, res *GuardResult) {
	if IsAllowedLargeWaiver(rel) {
		return
	}
	maxKb := opts.MaxFileKb
	if maxKb <= 0 {
		maxKb = 500
	}
	limitBytes := int64(maxKb * 1024)
	if sz > limitBytes {
		item := OversizedFile{Path: rel, SizeBytes: sz, MaxAllowedBytes: limitBytes, Kind: "oversized_file"}
		res.OversizedFiles = append(res.OversizedFiles, item)
	}
}
