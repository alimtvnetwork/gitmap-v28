package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"path/filepath"
)

func populateRepoErrorDetails(rec *InternalErrorRecord) {
	if rec.RepoPath == "" {
		return
	}

	repoDBPath := filepath.Join(rec.RepoPath, ".gitmap", "errors.db")
	repoConn, err := sql.Open("sqlite", repoDBPath)
	if err != nil {
		return
	}
	defer repoConn.Close()

	fetchAndApplyRepoErrorDetails(repoConn, rec)
}

func fetchAndApplyRepoErrorDetails(conn *sql.DB, rec *InternalErrorRecord) {
	row := conn.QueryRow(`SELECT Details, SourceFile, ContextJson, StackTrace, Notes, Comments FROM RepoErrorDB WHERE InternalErrorLogId = ?`, rec.ID)
	var det, src, ctx, stack, notes, comm sql.NullString
	if errScan := row.Scan(&det, &src, &ctx, &stack, &notes, &comm); errScan != nil {
		return
	}

	applyNullStringFields(rec, det, src, ctx, stack, notes, comm)
}

func applyNullStringFields(rec *InternalErrorRecord, det, src, ctx, stack, notes, comm sql.NullString) {
	if det.Valid {
		rec.Details = det.String
	}
	if src.Valid {
		rec.SourceFile = src.String
	}
	if ctx.Valid {
		rec.ContextJson = ctx.String
	}
	if stack.Valid {
		rec.StackTrace = stack.String
	}
	if notes.Valid {
		rec.Notes = notes.String
	}
	if comm.Valid {
		rec.Comments = comm.String
	}
}
