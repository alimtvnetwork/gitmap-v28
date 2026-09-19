package cmdautomation

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type cleanTotals struct {
	freedBytes   int64
	deletedCount int
}

// RunCleanArtifacts discovers and removes build and test artifacts.
func RunCleanArtifacts(opts CleanArtifactsOptions) CleanArtifactsResultMonad {
	start := time.Now()
	items := discoverCleanArtifacts(opts)
	processed := executeArtifactCleaning(items, opts.IsDryRun)
	res := aggregateCleanArtifacts(processed, start)
	return result.Ok(res)
}

func resolveArtifactBaseDir(dir string) string {
	hasDir := len(dir) > 0
	if hasDir {
		return filepath.Clean(dir)
	}
	return "."
}

func isSkipDirectory(name string) bool {
	isGit := name == ".git" || name == ".gitmap"
	isVendor := name == "vendor" || name == "node_modules"
	return isGit || isVendor
}

func isBinaryArtifact(ext string, name string) bool {
	isExe := ext == ".exe" || ext == ".syso" || ext == ".dll"
	isShared := ext == ".so" || ext == ".dylib" || ext == ".bin"
	isGitmap := name == "gitmap.exe" || name == "gitmap"
	return isExe || isShared || isGitmap
}

func isPycacheArtifact(name string, ext string, isDir bool) bool {
	isPycacheDir := isDir && (name == "__pycache__" || name == ".pytest_cache")
	isPyc := ext == ".pyc" || ext == ".pyo" || ext == ".pyd"
	isCoverage := name == ".coverage" || name == "coverage.out" || name == "test_output"
	return isPycacheDir || isPyc || isCoverage
}

func isTempArtifact(name string, ext string, isDir bool) bool {
	isTempDir := isDir && (name == "tmp" || name == "temp" || name == ".tmp")
	isTempExt := ext == ".tmp" || ext == ".swp" || ext == ".bak" || ext == ".log"
	isDsStore := name == ".DS_Store" || name == "Thumbs.db"
	return isTempDir || isTempExt || isDsStore
}

func classifyArtifact(name string, ext string, isDir bool, opts CleanArtifactsOptions) string {
	canCleanBin := opts.IsAll || opts.IsCleanBinaries
	if canCleanBin && isBinaryArtifact(ext, name) {
		return "binary"
	}
	canCleanPy := opts.IsAll || opts.IsCleanPycache
	if canCleanPy && isPycacheArtifact(name, ext, isDir) {
		return "pycache"
	}
	canCleanTmp := opts.IsAll || opts.IsCleanTemp
	if canCleanTmp && isTempArtifact(name, ext, isDir) {
		return "temp"
	}
	return ""
}

func checkGitTracked(relPath string, baseDir string) bool {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", relPath)
	cmd.Dir = baseDir
	err := cmd.Run()
	return err == nil
}

func processArtifactDeletion(item CleanArtifactItem, isDryRun bool) CleanArtifactItem {
	if isDryRun {
		item.IsDeleted = false
		return item
	}
	err := os.RemoveAll(item.Path)
	item.IsDeleted = (err == nil)
	return item
}

func executeArtifactCleaning(items []CleanArtifactItem, isDryRun bool) []CleanArtifactItem {
	processed := make([]CleanArtifactItem, 0, len(items))
	for _, it := range items {
		updated := processArtifactDeletion(it, isDryRun)
		processed = append(processed, updated)
	}
	return processed
}

func computeCleanTotals(items []CleanArtifactItem) cleanTotals {
	var tot cleanTotals
	for _, it := range items {
		if it.IsDeleted {
			tot.deletedCount++
			tot.freedBytes += it.SizeBytes
		}
	}
	return tot
}

func aggregateCleanArtifacts(items []CleanArtifactItem, start time.Time) CleanArtifactsResult {
	tot := computeCleanTotals(items)
	return CleanArtifactsResult{
		TotalFound:   len(items),
		TotalDeleted: tot.deletedCount,
		FreedBytes:   tot.freedBytes,
		FreedMB:      float64(tot.freedBytes) / (1024 * 1024),
		Items:        items,
		Duration:     time.Since(start),
		IsSuccess:    true,
	}
}

func discoverCleanArtifacts(opts CleanArtifactsOptions) []CleanArtifactItem {
	baseDir := resolveArtifactBaseDir(opts.Dir)
	var items []CleanArtifactItem
	walkFn := buildArtifactWalkFunc(&items, baseDir, opts)
	_ = filepath.WalkDir(baseDir, walkFn)
	return items
}

func buildArtifactWalkFunc(items *[]CleanArtifactItem, baseDir string, opts CleanArtifactsOptions) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		return inspectArtifactEntry(items, path, d, baseDir, opts)
	}
}

func inspectArtifactEntry(items *[]CleanArtifactItem, path string, d fs.DirEntry, baseDir string, opts CleanArtifactsOptions) error {
	name := d.Name()
	if d.IsDir() && isSkipDirectory(name) {
		return filepath.SkipDir
	}
	category := classifyArtifact(name, filepath.Ext(name), d.IsDir(), opts)
	if len(category) == 0 {
		return nil
	}
	*items = append(*items, createArtifactItem(path, d, baseDir, category))
	if d.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

func createArtifactItem(path string, d fs.DirEntry, baseDir string, category string) CleanArtifactItem {
	rel, _ := filepath.Rel(baseDir, path)
	relPath := filepath.ToSlash(rel)
	info, _ := d.Info()
	var size int64
	if info != nil {
		size = info.Size()
	}
	isTracked := checkGitTracked(relPath, baseDir)
	return CleanArtifactItem{
		Path:         relPath,
		SizeBytes:    size,
		Category:     category,
		IsGitTracked: isTracked,
	}
}
