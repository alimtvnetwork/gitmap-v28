package orchestrator

import (
	"fmt"
	"io"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/workspace"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func resolveSource(raw *commitin.RawArgs, stderr io.Writer) (*workspace.SourceHandle, int) {
	fmt.Fprintf(stderr, constants.CommitInMsgPhaseResolveSource, raw.Source)
	src, err := workspace.EnsureSource(raw.Source)
	if err != nil {
		fmt.Fprint(stderr, err.Error())
		return nil, constants.CommitInExitSourceUnusable
	}
	return src, constants.CommitInExitOk
}

func ensureWorkspace(sourceRoot string, stderr io.Writer) (*workspace.Paths, int) {
	paths, err := workspace.EnsureWorkspace(sourceRoot)
	if err != nil {
		fmt.Fprintf(stderr, "commit-in: workspace: %v\n", err)
		return nil, constants.CommitInExitDbFailed
	}
	return paths, constants.CommitInExitOk
}

func acquireLock(paths *workspace.Paths, stderr io.Writer) (*workspace.LockHandle, int) {
	lock, err := workspace.AcquireLock(paths)
	if err != nil {
		fmt.Fprint(stderr, err.Error())
		return nil, constants.CommitInExitLockBusy
	}
	return lock, constants.CommitInExitOk
}

func openAndMigrate(paths *workspace.Paths, stderr io.Writer) (dbCloser, int) {
	db, err := store.OpenAt(paths.DbFile)
	if err != nil {
		fmt.Fprintf(stderr, constants.CommitInErrDbMigrate, err)
		return nil, constants.CommitInExitDbFailed
	}
	if err := db.Migrate(); err != nil {
		_ = db.Close()
		fmt.Fprintf(stderr, constants.CommitInErrDbMigrate, err)
		return nil, constants.CommitInExitDbFailed
	}
	return db, constants.CommitInExitOk
}

// loadProfile resolves --profile / --default to a Profile (may be nil)
// and computes the layered Resolved settings. Missing-but-requested
// profile is fatal; the default-profile case is best-effort (missing
// returns nil so defaults+CLI still apply).
func loadProfile(raw *commitin.RawArgs, paths *workspace.Paths, _ dbCloser, stderr io.Writer) (profile.Resolved, *profile.Profile, int) {
	prof, code := pickProfile(raw, paths, stderr)
	if code != constants.CommitInExitOk {
		return profile.Resolved{}, nil, code
	}
	cli := buildCliOverrides(raw)
	return profile.Resolve(cli, prof), prof, constants.CommitInExitOk
}

func pickProfile(raw *commitin.RawArgs, paths *workspace.Paths, stderr io.Writer) (*profile.Profile, int) {
	if raw.ProfileName != "" {
		p, err := profile.LoadFromDisk(paths.SourceRoot, raw.ProfileName)
		if err != nil {
			fmt.Fprintf(stderr, constants.CommitInErrProfileMissing, raw.ProfileName)
			return nil, constants.CommitInExitProfileMissing
		}
		return p, constants.CommitInExitOk
	}
	if raw.UseDefaultProfile {
		p, err := profile.LoadFromDisk(paths.SourceRoot, "default")
		if err == nil {
			return p, constants.CommitInExitOk
		}
	}
	return nil, constants.CommitInExitOk
}

// runIDDir matches workspace.CloneInputs's <runId> subdir naming so
// finalize.CleanupTemp targets the exact directory we staged into.
func runIDDir(runID int64) string { return strconv.FormatInt(runID, 10) }
