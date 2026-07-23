package store

import (
	"testing"
	"time"
)

// seedDeps creates the minimum CommitInRun + RewrittenCommit chain so
// the replay-map tests have valid FK targets. Inserts:
//
//   - 1 Profile (id=1) so CommitInRun has a parent.
//   - 1 CommitInRun  (id=1).
//   - 1 InputRepo    (id=1) per CommitInRun.
//   - 2 SourceCommit (id=1, id=2) so the "tagged vs non-tagged"
//     test can attach two different tag rows to two commits.
//   - 2 RewrittenCommit rows pointing at the two SourceCommits.
//
// All FK enums are resolved via the standard mirror-table pattern so
// migration drift surfaces here BEFORE the replay-map test runs.
func seedDeps(t *testing.T, db *DB) {
	t.Helper()
	now := time.Now().Format(time.RFC3339)
	mustExec(t, db, `INSERT INTO Profile (ProfileId, Name, SourceRepoPath, PayloadJson)
		VALUES (1, 'p', '/abs/src', '{}')`)
	mustExec(t, db, `INSERT INTO CommitInRun
		(CommitInRunId, RunStatusId, SourceRepoPath, WasSourceFreshlyInit,
		 ProfileId, StartedAt)
		VALUES (1,
		 (SELECT RunStatusId FROM RunStatus WHERE Name='Running'),
		 '/abs/src', 0, 1, ?)`, now)
	mustExec(t, db, `INSERT INTO InputRepo
		(InputRepoId, CommitInRunId, OrderIndex, OriginalRef,
		 ResolvedPath, InputKindId)
		VALUES (1, 1, 1, 'inp', '/tmp/inp',
		 (SELECT InputKindId FROM InputKind WHERE Name='LocalFolder'))`)
	for i := 1; i <= 2; i++ {
		mustExec(t, db, `INSERT INTO SourceCommit
			(SourceCommitId, InputRepoId, OrderIndex, SourceSha,
			 AuthorName, AuthorEmail, AuthorDate, CommitterDate, OriginalMessage)
			VALUES (?, 1, ?, ?, 'A', 'a@x', ?, ?, 'm')`,
			i, i, "src-sha-"+itoa(i), now, now)
		mustExec(t, db, `INSERT INTO RewrittenCommit
			(RewrittenCommitId, CommitInRunId, SourceCommitId, NewSha,
			 FinalMessage, AppliedAuthorName, AppliedAuthorEmail,
			 AppliedAuthorDate, AppliedCommitterDate, CommitOutcomeId)
			VALUES (?, 1, ?, ?, 'm', 'A', 'a@x', ?, ?,
			 (SELECT CommitOutcomeId FROM CommitOutcome WHERE Name='Created'))`,
			i, i, "new-sha-"+itoa(i), now, now)
	}
}

// mustInsertReplay inserts one CommitInReplayMap row for the given
// rewritten-commit id with the supplied tag name + version flag,
// using TagReplayOutcome=Created.
func mustInsertReplay(t *testing.T, db *DB, runID, rewID int64, tagName string, isVersion bool) {
	t.Helper()
	v := 0
	if isVersion {
		v = 1
	}
	q := `INSERT INTO CommitInReplayMap
		(CommitInRunId, RewrittenCommitId, SourceTagName, SourceTagSha,
		 SourceCommitSha, DestTagSha, DestCommitSha,
		 MirroredReleaseBranch, IsVersionTag, TagReplayOutcomeId)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,
		 (SELECT TagReplayOutcomeId FROM TagReplayOutcome WHERE Name='Created'))`
	branch := ""
	if isVersion {
		branch = "release/" + tagName
	}
	mustExec(t, db, q, runID, rewID, tagName, "tag-sha-"+tagName,
		"src-sha-"+itoa(int(rewID)), "dest-tag-"+tagName,
		"new-sha-"+itoa(int(rewID)), nullable(branch), v)
}

func mustExec(t *testing.T, db *DB, query string, args ...any) {
	t.Helper()
	if _, err := db.conn.Exec(query, args...); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

func scanInt(t *testing.T, db *DB, q string, args ...any) int {
	t.Helper()
	var n int
	if err := db.conn.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatalf("scan int %q: %v", q, err)
	}
	return n
}

// nullable maps "" -> nil so SQL NULL is stored (mirrors the runlog
// helper of the same name; duplicated here to avoid a cross-package
// test-only dependency).
func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// itoa avoids pulling strconv into every call site; tests stay terse.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	const digits = "0123456789"
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	return string(buf[i:])
}
