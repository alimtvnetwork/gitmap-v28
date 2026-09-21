package pipelinedb

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// OpenPipelineSplitDb opens or initializes the split SQLite database for a repo.
func OpenPipelineSplitDb(repoSlug string) (*PipelineSplitDb, error) {
	dbPath := ResolvePipelineDbPath(repoSlug)
	dir := filepath.Dir(dbPath)
	_ = os.MkdirAll(dir, 0755)
	cleanLegacyDbIfPresentInDir(dir)

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open pipeline split db "+repoSlug)
	}

	return initPipelineSplitConn(conn, repoSlug, dbPath)
}

func cleanLegacyDbIfPresentInDir(dir string) {
	legacy := filepath.Join(dir, "pipeline.db")
	if isFileExisting(legacy) {
		_ = os.Remove(legacy)
	}
}

func initPipelineSplitConn(conn *sql.DB, repoSlug, dbPath string) (*PipelineSplitDb, error) {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "configure pipeline split db "+repoSlug)
	}

	p := &PipelineSplitDb{conn: conn, RepoSlug: repoSlug, Path: dbPath}
	if err := p.setupSchema(); err != nil {
		_ = conn.Close()

		return resetAndReopenPipelineSplitDb(repoSlug, dbPath)
	}

	return p, nil
}

func wipeDbFiles(dbPath string) {
	_ = os.Remove(dbPath)
	_ = os.Remove(dbPath + "-wal")
	_ = os.Remove(dbPath + "-shm")
}

func openCleanConn(dbPath, repoSlug string) (*sql.DB, error) {
	newConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "reopen pipeline split db "+repoSlug)
	}
	if cfgErr := store.ConfigureSQLiteConn(newConn); cfgErr != nil {
		_ = newConn.Close()

		return nil, apperror.WrapSimple(cfgErr, "reconfigure pipeline split db "+repoSlug)
	}

	return newConn, nil
}

func resetAndReopenPipelineSplitDb(repoSlug, dbPath string) (*PipelineSplitDb, error) {
	wipeDbFiles(dbPath)
	newConn, err := openCleanConn(dbPath, repoSlug)
	if err != nil {
		return nil, err
	}
	p := &PipelineSplitDb{conn: newConn, RepoSlug: repoSlug, Path: dbPath}
	if schemaErr := p.setupSchema(); schemaErr != nil {
		_ = newConn.Close()

		return nil, schemaErr
	}

	return p, nil
}

// OpenPipelineSplitDB is an alias to OpenPipelineSplitDb.
var OpenPipelineSplitDB = OpenPipelineSplitDb

// Close closes the underlying SQLite connection.
func (p *PipelineSplitDb) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}
