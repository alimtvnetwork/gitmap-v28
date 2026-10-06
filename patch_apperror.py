import re

def replace_error(file_path):
    with open(file_path, 'r') as f:
        content = f.read()

    # Replacements for pull_split_db_errors.go
    content = content.replace('func (s *PullSplitDB) EnsurePullErrorsTable() error', 'func (s *PullSplitDB) EnsurePullErrorsTable() *apperror.AppError')
    content = content.replace('func (s *PullSplitDB) InsertPullError(rec PullErrorRecord) error', 'func (s *PullSplitDB) InsertPullError(rec PullErrorRecord) *apperror.AppError')
    content = content.replace('func (s *PullSplitDB) execInsertPullError(rec PullErrorRecord) error', 'func (s *PullSplitDB) execInsertPullError(rec PullErrorRecord) *apperror.AppError')
    content = content.replace('func (s *PullSplitDB) QueryAllLatestPullErrors(limit int) ([]PullErrorRecord, error)', 'func (s *PullSplitDB) QueryAllLatestPullErrors(limit int) ([]PullErrorRecord, *apperror.AppError)')
    content = content.replace('func (s *PullSplitDB) QueryLatestPullErrors(repoSlug string, limit int) ([]PullErrorRecord, error)', 'func (s *PullSplitDB) QueryLatestPullErrors(repoSlug string, limit int) ([]PullErrorRecord, *apperror.AppError)')
    content = content.replace('func scanPullErrorRows(rows *sql.Rows) ([]PullErrorRecord, error)', 'func scanPullErrorRows(rows *sql.Rows) ([]PullErrorRecord, *apperror.AppError)')

    # Replacements for errors_split_db.go
    content = content.replace('func (s *ErrorsSplitDB) InitSchema() error', 'func (s *ErrorsSplitDB) InitSchema() *apperror.AppError')
    content = content.replace('func (s *ErrorsSplitDB) FederatedClearErrors() error', 'func (s *ErrorsSplitDB) FederatedClearErrors() *apperror.AppError')
    content = content.replace('func (s *ErrorsSplitDB) FederatedGetError(id int64) (*InternalErrorRecord, error)', 'func (s *ErrorsSplitDB) FederatedGetError(id int64) (*InternalErrorRecord, *apperror.AppError)')
    content = content.replace('func (s *ErrorsSplitDB) FederatedListErrors(limit int, unresolvedOnly bool) ([]InternalErrorRecord, error)', 'func (s *ErrorsSplitDB) FederatedListErrors(limit int, unresolvedOnly bool) ([]InternalErrorRecord, *apperror.AppError)')
    content = content.replace('func OpenErrorsSplitDB() (*ErrorsSplitDB, error)', 'func OpenErrorsSplitDB() (*ErrorsSplitDB, *apperror.AppError)')
    content = content.replace('func OpenErrorsSplitDBAt(dbPath string) (*ErrorsSplitDB, error)', 'func OpenErrorsSplitDBAt(dbPath string) (*ErrorsSplitDB, *apperror.AppError)')
    content = content.replace('func initErrorsSplitConn(conn *sql.DB, dbPath string) (*ErrorsSplitDB, error)', 'func initErrorsSplitConn(conn *sql.DB, dbPath string) (*ErrorsSplitDB, *apperror.AppError)')
    content = content.replace('func (s *ErrorsSplitDB) Close() error', 'func (s *ErrorsSplitDB) Close() *apperror.AppError')

    # Also fix rows.Err() return in scanPullErrorRows
    content = content.replace('return records, rows.Err()', 'if rows.Err() != nil { return records, apperror.WrapSimple(rows.Err(), "scan") }\n\treturn records, nil')

    with open(file_path, 'w') as f:
        f.write(content)

replace_error('cli/store/pull_split_db_errors.go')
replace_error('cli/store/errors_split_db.go')
