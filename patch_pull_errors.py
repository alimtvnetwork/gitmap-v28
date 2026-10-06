import re

with open('cli/store/pull_split_db_errors.go', 'r') as f:
    content = f.read()

# Modify sqlCreatePullErrors
content = re.sub(
    r'CREATE TABLE IF NOT EXISTS pull_errors \([\s\S]*?\);',
    '''CREATE TABLE IF NOT EXISTS pull_errors (
    error_id        TEXT PRIMARY KEY,
    repo_slug       TEXT NOT NULL,
    repo_path       TEXT NOT NULL,
    node_id         TEXT NOT NULL,
    node_version    TEXT NOT NULL,
    error_type      TEXT NOT NULL,
    error_text      TEXT NOT NULL,
    created_at      DATETIME NOT NULL
);''',
    content
)

# Modify execInsertPullError
content = re.sub(
    r'func \(s \*PullSplitDB\) execInsertPullError\(rec PullErrorRecord\) error \{[\s\S]*?return nil\n\}',
    '''func (s *PullSplitDB) execInsertPullError(rec PullErrorRecord) error {
	query := `INSERT OR REPLACE INTO pull_errors
		(error_id, repo_slug, repo_path, node_id, node_version, error_type, error_text, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.conn.Exec(query,
		rec.ErrorID, rec.RepoSlug, rec.RepoPath, rec.NodeID, rec.NodeVersion,
		rec.ErrorType, rec.ErrorText,
		rec.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return apperror.WrapSimple(err, "pull_errors.insert")
	}

	repoDBPath := filepath.Join(rec.RepoPath, ".gitmap", "pull_errors.db")
	_ = os.MkdirAll(filepath.Dir(repoDBPath), 0755)
	repoConn, repoErr := sql.Open("sqlite", repoDBPath)
	if repoErr == nil {
		defer repoConn.Close()
		_, _ = repoConn.Exec(`CREATE TABLE IF NOT EXISTS RepoPullErrorDetails (
			error_id TEXT PRIMARY KEY,
			stack_trace TEXT,
			remediation_cmd TEXT
		);`)
		_, _ = repoConn.Exec(`INSERT OR REPLACE INTO RepoPullErrorDetails (error_id, stack_trace, remediation_cmd) VALUES (?, ?, ?)`,
			rec.ErrorID, rec.StackTrace, rec.RemediationCmd)
	}

	return nil
}''',
    content
)

# Modify QueryLatestPullErrors to do federated read
content = re.sub(
    r'func \(s \*PullSplitDB\) QueryLatestPullErrors\(repoSlug string, limit int\) \(\[\]PullErrorRecord, error\) \{[\s\S]*?return scanPullErrorRows\(rows\)\n\}',
    '''func (s *PullSplitDB) QueryLatestPullErrors(repoSlug string, limit int) ([]PullErrorRecord, error) {
	if s.conn == nil {
		return nil, apperror.NewValidationError("database connection is nil")
	}

	limit = normalizePullErrorLimit(limit)
	query, args := buildPullErrorQuery(repoSlug, limit)
	rows, err := s.conn.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_errors.query")
	}
	defer rows.Close()

	records, err := scanPullErrorRows(rows)
	if err != nil {
		return records, err
	}

	for i := range records {
		records[i] = enrichPullErrorRecord(records[i])
	}
	return records, nil
}

func enrichPullErrorRecord(rec PullErrorRecord) PullErrorRecord {
	repoDBPath := filepath.Join(rec.RepoPath, ".gitmap", "pull_errors.db")
	repoConn, err := sql.Open("sqlite", repoDBPath)
	if err == nil {
		defer repoConn.Close()
		row := repoConn.QueryRow(`SELECT stack_trace, remediation_cmd FROM RepoPullErrorDetails WHERE error_id = ?`, rec.ErrorID)
		var stack, rem sql.NullString
		if row.Scan(&stack, &rem) == nil {
			if stack.Valid {
				rec.StackTrace = stack.String
			}
			if rem.Valid {
				rec.RemediationCmd = rem.String
			}
		}
	}
	return rec
}''',
    content
)

# Fix buildPullErrorQuery fields
content = re.sub(
    r'SELECT error_id, repo_slug, repo_path, node_id, node_version, error_type, error_text, stack_trace, remediation_cmd, created_at',
    'SELECT error_id, repo_slug, repo_path, node_id, node_version, error_type, error_text, created_at',
    content
)

# Fix scanSinglePullError
content = re.sub(
    r'func scanSinglePullError\(rows \*sql\.Rows\) \(PullErrorRecord, bool\) \{[\s\S]*?return rec, true\n\}',
    '''func scanSinglePullError(rows *sql.Rows) (PullErrorRecord, bool) {
	var rec PullErrorRecord
	var rawCreatedAt any

	err := rows.Scan(
		&rec.ErrorID, &rec.RepoSlug, &rec.RepoPath, &rec.NodeID, &rec.NodeVersion,
		&rec.ErrorType, &rec.ErrorText, &rawCreatedAt,
	)
	if err != nil {
		return rec, false
	}

	rec.CreatedAt = parseFlexibleDBValue(rawCreatedAt)
	return rec, true
}''',
    content
)

if '"os"' not in content:
    content = re.sub(
        r'import \(\n',
        'import (\n\t"os"\n\t"path/filepath"\n',
        content
    )

with open('cli/store/pull_split_db_errors.go', 'w') as f:
    f.write(content)
