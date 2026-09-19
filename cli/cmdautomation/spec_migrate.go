package cmdautomation

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunSpecMigrate re-sequences spec file prefixes and updates cross-references.
func RunSpecMigrate(opts SpecMigrateOptions) SpecMigrateResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	specDir := resolveSpecDir(root, opts.Dir)
	specFiles := collectSpecFiles(specDir)

	hasRange := opts.FromNum > 0 && opts.ToNum > 0
	if !hasRange {
		return result.Ok(buildSpecOverviewResult(specFiles, start))
	}
	changeMonad := planSpecMigration(specDir, opts.FromNum, opts.ToNum)
	isFail := changeMonad.IsFailure()
	if isFail {
		return result.Fail[SpecMigrateResult](changeMonad.Err)
	}
	res := executeSpecMigration(root, changeMonad.Value, len(specFiles), opts.IsDryRun)
	res.Duration = time.Since(start)
	return result.Ok(res)
}

func resolveSpecDir(root, customDir string) string {
	base := root
	if len(customDir) > 0 {
		base = customDir
	}
	appDir := filepath.Join(base, "02-spec", "21-app")
	if _, err := os.Stat(appDir); err == nil {
		return appDir
	}
	specDir := filepath.Join(base, "02-spec")
	if _, err := os.Stat(specDir); err == nil {
		return specDir
	}
	return base
}

func collectSpecFiles(specDir string) []string {
	var files []string
	entries, err := os.ReadDir(specDir)
	if err != nil {
		return files
	}
	for _, e := range entries {
		isMd := strings.HasSuffix(strings.ToLower(e.Name()), ".md")
		if !e.IsDir() && isMd {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files
}

func buildSpecOverviewResult(specFiles []string, start time.Time) SpecMigrateResult {
	return SpecMigrateResult{
		TotalSpecs: len(specFiles),
		Duration:   time.Since(start),
		IsSuccess:  true,
	}
}

func planSpecMigration(specDir string, fromNum, toNum int) result.Result[SpecMigrateChange] {
	fromPrefix := fmt.Sprintf("%02d-", fromNum)
	if fromNum >= 100 {
		fromPrefix = fmt.Sprintf("%d-", fromNum)
	}
	sourceFile := findSpecFileWithPrefix(specDir, fromPrefix)
	if sourceFile == "" {
		err := apperror.NewValidationError(fmt.Sprintf("no spec file found with prefix %d", fromNum))
		return result.Fail[SpecMigrateChange](err)
	}
	targetFile := computeNewSpecName(sourceFile, fromNum, toNum)
	return result.Ok(SpecMigrateChange{
		SourceFile: filepath.Join(specDir, sourceFile),
		TargetFile: filepath.Join(specDir, targetFile),
	})
}

func findSpecFileWithPrefix(specDir, prefix string) string {
	entries, err := os.ReadDir(specDir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		hasPrefix := strings.HasPrefix(e.Name(), prefix)
		if !e.IsDir() && hasPrefix {
			return e.Name()
		}
	}
	return ""
}

func computeNewSpecName(sourceFile string, fromNum, toNum int) string {
	oldPrefix := fmt.Sprintf("%02d-", fromNum)
	if fromNum >= 100 {
		oldPrefix = fmt.Sprintf("%d-", fromNum)
	}
	newPrefix := fmt.Sprintf("%02d-", toNum)
	if toNum >= 100 {
		newPrefix = fmt.Sprintf("%d-", toNum)
	}
	base := strings.TrimPrefix(sourceFile, oldPrefix)
	return newPrefix + base
}

func executeSpecMigration(root string, change SpecMigrateChange, totalSpecs int, isDryRun bool) SpecMigrateResult {
	oldBase := filepath.Base(change.SourceFile)
	newBase := filepath.Base(change.TargetFile)
	updatedCount := updateCrossReferences(root, oldBase, newBase, isDryRun)
	change.References = updatedCount
	if !isDryRun {
		_ = os.Rename(change.SourceFile, change.TargetFile)
	}
	return SpecMigrateResult{
		TotalSpecs:   totalSpecs,
		Migrated:     []SpecMigrateChange{change},
		FilesUpdated: updatedCount,
		IsSuccess:    true,
	}
}

func updateCrossReferences(root, oldBase, newBase string, isDryRun bool) int {
	files := collectSpecMigrateFiles(root)
	updated := 0
	for _, f := range files {
		hasRef := replaceInFileIfPresent(f, oldBase, newBase, isDryRun)
		if hasRef {
			updated++
		}
	}
	return updated
}

func collectSpecMigrateFiles(root string) []string {
	var files []string
	searchDirs := []string{"02-spec", ".ai-memory", "cli"}
	for _, d := range searchDirs {
		p := filepath.Join(root, d)
		files = append(files, walkMarkdownFiles(p)...)
	}
	return files
}

func replaceInFileIfPresent(path, oldStr, newStr string, isDryRun bool) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	content := string(data)
	if !strings.Contains(content, oldStr) {
		return false
	}
	if !isDryRun {
		mod := strings.ReplaceAll(content, oldStr, newStr)
		_ = os.WriteFile(path, []byte(mod), 0o644)
	}
	return true
}
