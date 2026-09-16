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
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open pipeline split db "+repoSlug)
	}

	return initPipelineSplitConn(conn, repoSlug, dbPath)
}

func initPipelineSplitConn(conn *sql.DB, repoSlug, dbPath string) (*PipelineSplitDb, error) {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "configure pipeline split db "+repoSlug)
	}

	p := &PipelineSplitDb{conn: conn, RepoSlug: repoSlug, Path: dbPath}

	return p, p.setupSchema()
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
