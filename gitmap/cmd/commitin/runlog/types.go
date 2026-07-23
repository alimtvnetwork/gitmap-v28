package runlog

// SourceCommitRow is the persistence-layer projection of one walked
// commit. Date strings are pre-formatted RFC3339 so the writer never
// touches a `time.Time` (keeps the SQL layer pure-string).
type SourceCommitRow struct {
	OrderIndex           int
	Sha                  string
	AuthorName           string
	AuthorEmail          string
	AuthorDateRFC3339    string
	CommitterDateRFC3339 string
	OriginalMessage      string
	Files                []string
}

// RewrittenRow is the persistence-layer projection of one
// RewrittenCommit row. Outcome MUST be a constants.CommitInOutcome*
// literal; NewSha empty when Outcome != Created.
type RewrittenRow struct {
	NewSha               string
	SourceSha            string // for ShaMap insert
	FinalMessage         string
	AuthorName           string
	AuthorEmail          string
	AuthorDateRFC3339    string
	CommitterDateRFC3339 string
	Outcome              string
}
