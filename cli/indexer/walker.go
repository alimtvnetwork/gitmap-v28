package indexer

import (
	"bytes"
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/worker"
)

const (
	maxFileSize   = 300 * 1024 // 300KB
	probeMaxBytes = 8192       // 8KB binary sniffer probe
)

var excludedDirMap = map[string]bool{
	".git":         true,
	"node_modules": true,
	".venv":        true,
	"dist":         true,
	"build":        true,
	"bin":          true,
	"vendor":       true,
	".gemini":      true,
	"coverage":     true,
	"tmp":          true,
	"__pycache__":  true,
	".turbo":       true,
}

type FileInfo struct {
	AbsolutePath string
	RelativePath string
	IsBig        bool
	WriteTime    int64
	Content      string
}

type Walker struct {
	RepoPath string
	RepoDB   *sql.DB
	ForceDot bool
}

func NewWalker(repoPath string, repoDB *sql.DB, forceDot bool) *Walker {
	return &Walker{
		RepoPath: repoPath,
		RepoDB:   repoDB,
		ForceDot: forceDot,
	}
}

func scanWriteTimes(rows *sql.Rows, times map[string]int64) result.ResultMap[string, int64] {
	for rows.Next() {
		var relPath string
		var writeTime int64
		if err := rows.Scan(&relPath, &writeTime); err != nil {
			appErr := apperror.WrapSimple(err, "scan write time row")

			return result.FailMap[string, int64](appErr)
		}

		times[relPath] = writeTime
	}

	if err := rows.Err(); err != nil {
		appErr := apperror.WrapSimple(err, "iterate write time rows")

		return result.FailMap[string, int64](appErr)
	}

	return result.OkMap(times)
}

func loadExistingWriteTimes(ctx context.Context, db *sql.DB) result.ResultMap[string, int64] {
	times := make(map[string]int64)
	if db == nil {
		return result.OkMap(times)
	}

	rows, err := db.QueryContext(ctx, "SELECT RelativePath, WriteTime FROM RepoFile")
	if err != nil {
		appErr := apperror.WrapSimple(err, "query repo file write times")

		return result.FailMap[string, int64](appErr)
	}

	defer rows.Close()

	return scanWriteTimes(rows, times)
}

func isFileModified(cachedTimes map[string]int64, relPath string, writeTime int64) bool {
	lastWriteTime, hasCached := cachedTimes[relPath]
	if !hasCached {
		return true
	}

	return lastWriteTime < writeTime
}

func isExcludedDir(name string) bool {
	return excludedDirMap[name]
}

func (w *Walker) isDotDirSkipped(name string, path string) bool {
	if w.ForceDot {
		return false
	}

	if path == w.RepoPath {
		return false
	}

	return strings.HasPrefix(name, ".")
}

func (w *Walker) handleDirSkip(d fs.DirEntry, path string) error {
	name := d.Name()
	if isExcludedDir(name) {
		return filepath.SkipDir
	}

	if w.isDotDirSkipped(name, path) {
		return filepath.SkipDir
	}

	return nil
}

func isBinaryContent(head []byte) bool {
	probeLimit := len(head)
	if probeLimit > probeMaxBytes {
		probeLimit = probeMaxBytes
	}

	return bytes.IndexByte(head[:probeLimit], 0) != -1
}

func readFileContent(path string) (string, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	if isBinaryContent(b) {
		return "", true
	}

	return string(b), false
}

func resolveFileInfo(info FileInfo) FileInfo {
	if info.IsBig {
		return info
	}

	content, isBinary := readFileContent(info.AbsolutePath)
	if isBinary {
		info.IsBig = true
		info.Content = ""

		return info
	}

	info.Content = content

	return info
}

const sqlUpsertRepoFile = `
INSERT INTO RepoFile (RelativePath, AbsolutePath, Content, IsBig, WriteTime, CreatedAt, UpdatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RelativePath) DO UPDATE SET
	Content=excluded.Content,
	IsBig=excluded.IsBig,
	WriteTime=excluded.WriteTime,
	UpdatedAt=excluded.UpdatedAt;`

func (w *Walker) upsertRepoFile(ctx context.Context, info FileInfo) (bool, error) {
	now := time.Now().Unix()
	_, err := w.RepoDB.ExecContext(ctx, sqlUpsertRepoFile,
		info.RelativePath, info.AbsolutePath, info.Content, info.IsBig, info.WriteTime, now, now,
	)

	return err == nil, err
}

func (w *Walker) processFile(ctx context.Context, info FileInfo) (bool, error) {
	resolved := resolveFileInfo(info)

	return w.upsertRepoFile(ctx, resolved)
}

func (w *Walker) inspectAndQueueFile(path string, d fs.DirEntry, cachedTimes map[string]int64, fileChan chan<- FileInfo) error {
	info, err := d.Info()
	if err != nil {
		return nil
	}

	relPath, err := filepath.Rel(w.RepoPath, path)
	if err != nil {
		return nil
	}

	writeTime := info.ModTime().Unix()
	if isFileModified(cachedTimes, relPath, writeTime) {
		fileChan <- FileInfo{
			AbsolutePath: path,
			RelativePath: relPath,
			IsBig:        info.Size() > maxFileSize,
			WriteTime:    writeTime,
		}
	}

	return nil
}

func (w *Walker) buildWalkFn(cachedTimes map[string]int64, fileChan chan<- FileInfo) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return w.handleDirSkip(d, path)
		}

		return w.inspectAndQueueFile(path, d, cachedTimes, fileChan)
	}
}

func startDrain(results <-chan worker.Result[bool]) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for range results {
		}

		close(done)
	}()

	return done
}

// Walk traverses the directory, schedules indexing for changed files
func (w *Walker) Walk(ctx context.Context, workers int) error {
	timesRes := loadExistingWriteTimes(ctx, w.RepoDB)
	if timesRes.IsFailure() {
		return timesRes.AppError()
	}

	fileChan := make(chan FileInfo, 100)
	pool := worker.NewPool(workers, func(c context.Context, input FileInfo) (bool, error) {
		return w.processFile(c, input)
	})
	done := startDrain(pool.Run(ctx, fileChan))
	walkErr := filepath.WalkDir(w.RepoPath, w.buildWalkFn(timesRes.Data, fileChan))
	close(fileChan)
	<-done

	return walkErr
}
