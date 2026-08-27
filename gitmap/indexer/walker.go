package indexer

import (
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/worker"
)

const maxFileSize = 300 * 1024 // 300KB

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

// Walk traverses the directory, schedules indexing for changed files
func (w *Walker) Walk(ctx context.Context, workers int) error {
	fileChan := make(chan FileInfo, 100)

	// Create generic worker pool
	pool := worker.NewPool(workers, func(c context.Context, input FileInfo) (bool, error) {
		return w.processFile(c, input)
	})

	// Start worker pool
	results := pool.Run(ctx, fileChan)

	// Consume results in a separate goroutine to avoid blocking
	go func() {
		for range results {
			// In a real application, handle errors here
		}
	}()

	err := filepath.WalkDir(w.RepoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			if !w.ForceDot && strings.HasPrefix(name, ".") && path != w.RepoPath {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(w.RepoPath, path)
		if err != nil {
			return nil
		}

		isBig := info.Size() > maxFileSize
		writeTime := info.ModTime().Unix()

		// Delta check
		var lastWriteTime int64
		err = w.RepoDB.QueryRowContext(ctx, "SELECT WriteTime FROM RepoFile WHERE RelativePath = ?", relPath).Scan(&lastWriteTime)
		if err == sql.ErrNoRows || lastWriteTime < writeTime {
			fileChan <- FileInfo{
				AbsolutePath: path,
				RelativePath: relPath,
				IsBig:        isBig,
				WriteTime:    writeTime,
			}
		}

		return nil
	})

	close(fileChan)
	return err
}

func (w *Walker) processFile(ctx context.Context, info FileInfo) (bool, error) {
	if !info.IsBig {
		b, err := os.ReadFile(info.AbsolutePath)
		if err == nil {
			info.Content = string(b)
		}
	}

	now := time.Now().Unix()
	
	// Upsert file into SQLite
	query := `
		INSERT INTO RepoFile (RelativePath, AbsolutePath, Content, IsBig, WriteTime, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(RelativePath) DO UPDATE SET
			Content=excluded.Content,
			IsBig=excluded.IsBig,
			WriteTime=excluded.WriteTime,
			UpdatedAt=excluded.UpdatedAt;
	`
	_, err := w.RepoDB.ExecContext(ctx, query,
		info.RelativePath, info.AbsolutePath, info.Content, info.IsBig, info.WriteTime, now, now,
	)

	return err == nil, err
}
