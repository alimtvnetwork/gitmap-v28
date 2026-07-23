package orchestrator

import (
	"database/sql"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/finalize"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/workspace"
)

// dbCloser is the subset of *store.DB the orchestrator needs. Defined
// as an interface so tests can swap the SQLite-backed implementation
// for an in-memory fake.
type dbCloser interface {
	Conn() *sql.DB
	Close() error
}

// runContext bundles every long-lived handle a single commit-in run
// needs. Cleanup releases the lock + closes the DB; it is safe to
// call multiple times (idempotent).
type runContext struct {
	Raw      *commitin.RawArgs
	Source   *workspace.SourceHandle
	Paths    *workspace.Paths
	Lock     *workspace.LockHandle
	DB       dbCloser
	Resolved profile.Resolved
	RunID    int64
	TempDir  string
	Counters finalize.Counters
	// inputRepoIds caches InputRepo PKs keyed by ResolvedInput.OrderIndex
	// so persistSource emits exactly one InputRepo row per staged input.
	inputRepoIds map[int]int64
	// aborted is flipped to true by conflictCheck when the user picks
	// Abort under Prompt mode. Pipeline loops MUST treat this as a hard
	// stop and return CommitInExitConflictAborted.
	aborted bool
}

func newContext(raw *commitin.RawArgs, src *workspace.SourceHandle, paths *workspace.Paths, lock *workspace.LockHandle, db dbCloser, resolved profile.Resolved, runID int64) *runContext {
	return &runContext{
		Raw:      raw,
		Source:   src,
		Paths:    paths,
		Lock:     lock,
		DB:       db,
		Resolved: resolved,
		RunID:    runID,
		TempDir:  filepath.Join(paths.TempRoot, runIDDir(runID)),
		Counters: finalize.Counters{RunId: runID},
	}
}

// Cleanup releases the lock, closes the DB, and removes the temp run
// dir unless --keep-temp was set. Safe to call multiple times.
func (c *runContext) Cleanup() {
	if c == nil {
		return
	}
	finalize.CleanupTemp(c.TempDir, c.Raw.IsKeepTemp)
	if c.DB != nil {
		_ = c.DB.Close()
		c.DB = nil
	}
	if c.Lock != nil {
		c.Lock.Release()
		c.Lock = nil
	}
}
