package enums

// SQLitePragmaType defines strongly-typed SQLite configuration PRAGMAs.
type SQLitePragmaType string

const (
	PragmaJournalWAL        SQLitePragmaType = "PRAGMA journal_mode = WAL;"
	PragmaSynchronousNormal SQLitePragmaType = "PRAGMA synchronous = NORMAL;"
	PragmaBusyTimeout5s     SQLitePragmaType = "PRAGMA busy_timeout = 5000;"
	PragmaForeignKeysOn     SQLitePragmaType = "PRAGMA foreign_keys = ON;"
)

// IsValid reports whether the pragma is a known standard configuration pragma.
func (p SQLitePragmaType) IsValid() bool {
	switch p {
	case PragmaJournalWAL, PragmaSynchronousNormal, PragmaBusyTimeout5s, PragmaForeignKeysOn:
		return true
	default:
		return false
	}
}

// String returns the raw pragma SQL statement.
func (p SQLitePragmaType) String() string {
	return string(p)
}

// Query returns the pragma query execution string.
func (p SQLitePragmaType) Query() string {
	return string(p)
}

// DefaultPragmas returns the slice of standard pragmas applied to SQLite connections.
func DefaultPragmas() []SQLitePragmaType {
	return []SQLitePragmaType{
		PragmaBusyTimeout5s,
		PragmaJournalWAL,
		PragmaSynchronousNormal,
		PragmaForeignKeysOn,
	}
}
