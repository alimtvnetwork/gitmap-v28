import re

with open('cli/store/errors_split_db.go', 'r') as f:
    content = f.read()

# Add RepoPath to InternalErrorRecord
content = content.replace(
    '	CreatedAt     string `json:"createdAt"`\n}',
    '	CreatedAt     string `json:"createdAt"`\n	RepoPath      string `json:"repoPath,omitempty"`\n}'
)

# Modify schema
content = re.sub(
    r'	sqlCreateInternalErrorLog = `CREATE TABLE IF NOT EXISTS InternalErrorLog \([\s\S]*?idxInternalErrorLog_IsResolved ON InternalErrorLog\(IsResolved\);`',
    '''	sqlCreateInternalErrorLog = `CREATE TABLE IF NOT EXISTS RootErrorIndex (
    InternalErrorLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    ErrorCode          TEXT NOT NULL DEFAULT '',
    ErrorType          TEXT NOT NULL DEFAULT 'general',
    Command            TEXT NOT NULL DEFAULT '',
    Message            TEXT NOT NULL,
    GitMapVersion      TEXT NOT NULL DEFAULT '',
    IsResolved         INTEGER NOT NULL DEFAULT 0,
    CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    RepoPath           TEXT NULL
);
CREATE INDEX IF NOT EXISTS IdxRootErrorIndex_CreatedAt ON RootErrorIndex(CreatedAt DESC);`

	sqlCreateRepoErrorDB = `CREATE TABLE IF NOT EXISTS RepoErrorDB (
    InternalErrorLogId INTEGER PRIMARY KEY,
    Details            TEXT NULL,
    SourceFile         TEXT NULL,
    ContextJson        TEXT NULL,
    StackTrace         TEXT NULL,
    Notes              TEXT NULL,
    Comments           TEXT NULL
);`''',
    content
)

# Ensure schema init doesn't break
content = re.sub(
    r'func \(s \*ErrorsSplitDB\) InitSchema\(\) error \{[\s\S]*?return nil\n\}',
    '''func (s *ErrorsSplitDB) InitSchema() error {
	if _, err := s.conn.Exec(sqlCreateInternalErrorLog); err != nil {
		return apperror.WrapSimple(err, "errors_split.init_schema")
	}
	if _, err := s.conn.Exec(sqlCreateFailedCommand); err != nil {
		return apperror.WrapSimple(err, "errors_split.init_failed_commands_schema")
	}

	return nil
}''',
    content
)

# Add Federated functions
federated_funcs = '''
func (s *ErrorsSplitDB) FederatedClearErrors() error {
	// Query all distinct repo paths
	rows, err := s.conn.Query(`SELECT DISTINCT RepoPath FROM RootErrorIndex WHERE RepoPath IS NOT NULL AND RepoPath != ''`)
	if err == nil {
		defer rows.Close()
		var paths []string
		for rows.Next() {
			var p string
			if err := rows.Scan(&p); err == nil {
				paths = append(paths, p)
			}
		}
		for _, p := range paths {
			repoDBPath := filepath.Join(p, ".gitmap", "errors.db")
			repoConn, err := sql.Open("sqlite", repoDBPath)
			if err == nil {
				_, _ = repoConn.Exec(`DELETE FROM RepoErrorDB`)
				repoConn.Close()
			}
		}
	}

	_, err = s.conn.Exec(`DELETE FROM RootErrorIndex`)
	if err != nil {
		return apperror.WrapSimple(err, "errors_split.clear")
	}
	return nil
}

func (s *ErrorsSplitDB) FederatedGetError(id int64) (*InternalErrorRecord, error) {
	row := s.conn.QueryRow(`SELECT InternalErrorLogId, ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, CreatedAt, COALESCE(RepoPath, '') FROM RootErrorIndex WHERE InternalErrorLogId = ?`, id)
	var rec InternalErrorRecord
	var isResolvedInt int

	err := row.Scan(
		&rec.ID, &rec.ErrorCode, &rec.ErrorType, &rec.Command, &rec.Message,
		&rec.GitMapVersion, &isResolvedInt, &rec.CreatedAt, &rec.RepoPath,
	)
	if err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.get")
	}

	rec.IsResolved = isResolvedInt > 0

	if rec.RepoPath != "" {
		repoDBPath := filepath.Join(rec.RepoPath, ".gitmap", "errors.db")
		repoConn, err := sql.Open("sqlite", repoDBPath)
		if err == nil {
			defer repoConn.Close()
			row := repoConn.QueryRow(`SELECT Details, SourceFile, ContextJson, StackTrace, Notes, Comments FROM RepoErrorDB WHERE InternalErrorLogId = ?`, id)
			var det, src, ctx, stack, notes, comm sql.NullString
			if row.Scan(&det, &src, &ctx, &stack, &notes, &comm) == nil {
				if det.Valid { rec.Details = det.String }
				if src.Valid { rec.SourceFile = src.String }
				if ctx.Valid { rec.ContextJson = ctx.String }
				if stack.Valid { rec.StackTrace = stack.String }
				if notes.Valid { rec.Notes = notes.String }
				if comm.Valid { rec.Comments = comm.String }
			}
		}
	}

	return &rec, nil
}

func (s *ErrorsSplitDB) FederatedListErrors(limit int, unresolvedOnly bool) ([]InternalErrorRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT InternalErrorLogId, ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, CreatedAt, COALESCE(RepoPath, '') FROM RootErrorIndex ORDER BY InternalErrorLogId DESC LIMIT ?`
	if unresolvedOnly {
		query = `SELECT InternalErrorLogId, ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, CreatedAt, COALESCE(RepoPath, '') FROM RootErrorIndex WHERE IsResolved = 1 ORDER BY InternalErrorLogId DESC LIMIT ?`
	}

	rows, err := s.conn.Query(query, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.list")
	}
	defer rows.Close()

	var records []InternalErrorRecord
	for rows.Next() {
		var rec InternalErrorRecord
		var isResolvedInt int
		if err := rows.Scan(&rec.ID, &rec.ErrorCode, &rec.ErrorType, &rec.Command, &rec.Message, &rec.GitMapVersion, &isResolvedInt, &rec.CreatedAt, &rec.RepoPath); err == nil {
			rec.IsResolved = isResolvedInt > 0
			records = append(records, rec)
		}
	}

	for i, rec := range records {
		if rec.RepoPath != "" {
			repoDBPath := filepath.Join(rec.RepoPath, ".gitmap", "errors.db")
			repoConn, err := sql.Open("sqlite", repoDBPath)
			if err == nil {
				row := repoConn.QueryRow(`SELECT Details, SourceFile, ContextJson, StackTrace, Notes, Comments FROM RepoErrorDB WHERE InternalErrorLogId = ?`, rec.ID)
				var det, src, ctx, stack, notes, comm sql.NullString
				if row.Scan(&det, &src, &ctx, &stack, &notes, &comm) == nil {
					if det.Valid { records[i].Details = det.String }
					if src.Valid { records[i].SourceFile = src.String }
					if ctx.Valid { records[i].ContextJson = ctx.String }
					if stack.Valid { records[i].StackTrace = stack.String }
					if notes.Valid { records[i].Notes = notes.String }
					if comm.Valid { records[i].Comments = comm.String }
				}
				repoConn.Close()
			}
		}
	}

	return records, nil
}

func LogInternalErrorRecord(rec InternalErrorRecord) {
	db, err := OpenErrorsSplitDB()
	if err != nil {
		return
	}
	defer db.Close()

	resolvedVal := 0
	if rec.IsResolved {
		resolvedVal = 1
	}

	res, err := db.conn.Exec(`INSERT INTO RootErrorIndex (ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, RepoPath) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rec.ErrorCode, rec.ErrorType, rec.Command, rec.Message, rec.GitMapVersion, resolvedVal, rec.RepoPath,
	)
	if err != nil {
		return
	}

	id, _ := res.LastInsertId()

	if rec.RepoPath != "" {
		repoDBPath := filepath.Join(rec.RepoPath, ".gitmap", "errors.db")
		_ = os.MkdirAll(filepath.Dir(repoDBPath), 0755)
		repoConn, err := sql.Open("sqlite", repoDBPath)
		if err == nil {
			defer repoConn.Close()
			_, _ = repoConn.Exec(sqlCreateRepoErrorDB)
			_, _ = repoConn.Exec(`INSERT INTO RepoErrorDB (InternalErrorLogId, Details, SourceFile, ContextJson, StackTrace, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?)`,
				id, rec.Details, rec.SourceFile, rec.ContextJson, rec.StackTrace, rec.Notes, rec.Comments,
			)
		}
	}
}
'''

content += "\n" + federated_funcs

with open('cli/store/errors_split_db.go', 'w') as f:
    f.write(content)
