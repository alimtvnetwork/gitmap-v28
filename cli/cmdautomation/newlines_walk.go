package cmdautomation

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunNormalizeNewlines walks files and applies polyglot newline normalization.
func RunNormalizeNewlines(opts NewlineOptions) (NewlineResult, *apperror.AppError) {
	start := time.Now()
	res := NewlineResult{}
	targets := resolveTargetPaths(opts.Paths)

	for _, p := range targets {
		processErr := processSingleFile(p, opts, &res)
		if processErr != nil {
			return res, processErr
		}
	}

	res.Duration = time.Since(start)
	return res, nil
}

func resolveTargetPaths(paths []string) []string {
	isEmpty := len(paths) == 0
	if isEmpty {
		return []string{"."}
	}
	return paths
}

func processSingleFile(path string, opts NewlineOptions, res *NewlineResult) *apperror.AppError {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	if info.IsDir() {
		return walkDirectoryNewlines(path, opts, res)
	}

	return normalizeFileOnDisk(path, opts, res)
}

func walkDirectoryNewlines(dir string, opts NewlineOptions, res *NewlineResult) *apperror.AppError {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() && isExcludedDir(info.Name()) {
			if info != nil && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			normalizeFileOnDisk(p, opts, res)
		}
		return nil
	})

	if err != nil {
		ctx := map[string]any{"dir": dir, "err": err.Error()}
		return apperror.New("normalize_walk", "E_WALK_FAILED", ctx)
	}
	return nil
}

func isExcludedDir(name string) bool {
	return name == ".git" || name == "node_modules" || name == "dist" ||
		name == "build" || name == ".pytest_cache" || name == "__pycache__"
}

func normalizeFileOnDisk(path string, opts NewlineOptions, res *NewlineResult) *apperror.AppError {
	ext := filepath.Ext(path)
	if !IsPolyglotTextExtension(ext) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil || HasBinaryContent(data) {
		return nil
	}

	res.ScannedFiles++
	cleaned, crlfs, isModified := NormalizeContent(data)
	res.CrlfCount += crlfs

	if isModified {
		res.ModifiedFiles++
		if opts.IsFixMode && !opts.IsDryRun {
			writeErr := os.WriteFile(path, cleaned, 0644)
			if writeErr != nil {
				ctx := map[string]any{"path": path, "err": writeErr.Error()}
				return apperror.New("write_file", "E_WRITE_FAILED", ctx)
			}
		}
	}
	return nil
}
