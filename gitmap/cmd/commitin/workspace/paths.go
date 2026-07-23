package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Paths is the resolved on-disk layout for a single commit-in run.
// All fields are absolute paths. Anchored at <sourceRoot>/.gitmap so
// the workspace travels with the source repo.
type Paths struct {
	SourceRoot   string // resolved absolute path of <source>
	GitmapRoot   string // <sourceRoot>/.gitmap
	CommitInRoot string // <sourceRoot>/.gitmap/commit-in
	ProfilesDir  string // <sourceRoot>/.gitmap/commit-in/profiles
	TempRoot     string // <sourceRoot>/.gitmap/temp
	LockFile     string // <sourceRoot>/.gitmap/commit-in.lock
	DbFile       string // <sourceRoot>/.gitmap/<dbName>
}

// EnsureWorkspace creates every directory the run needs. Idempotent:
// re-running on an existing workspace mutates nothing on disk.
// Maps to spec §3.1 stage 03 (`EnsureWorkspace`). On failure the
// caller exits with constants.CommitInExitDbFailed (per stage table).
func EnsureWorkspace(sourceRoot string) (*Paths, error) {
	abs, err := filepath.Abs(sourceRoot)
	if err != nil {
		return nil, fmt.Errorf("absolutize source: %w", err)
	}
	p := buildPaths(abs)
	if mkErr := makeDirs(p); mkErr != nil {
		return nil, mkErr
	}
	return p, nil
}

// buildPaths assembles the Paths struct. Pure function — no syscalls.
func buildPaths(sourceRoot string) *Paths {
	gitmapRoot := filepath.Join(sourceRoot, constants.GitMapDir)
	commitInRoot := filepath.Join(gitmapRoot, constants.CommitInDirRoot)
	return &Paths{
		SourceRoot:   sourceRoot,
		GitmapRoot:   gitmapRoot,
		CommitInRoot: commitInRoot,
		ProfilesDir:  filepath.Join(gitmapRoot, constants.CommitInDirProfiles),
		TempRoot:     filepath.Join(gitmapRoot, constants.CommitInDirTemp),
		LockFile:     filepath.Join(commitInRoot, constants.CommitInLockFileName),
		DbFile:       filepath.Join(gitmapRoot, constants.DBFile),
	}
}

// makeDirs creates every directory listed in Paths. Mode 0o755 mirrors
// the rest of the codebase (see store/store.go).
func makeDirs(p *Paths) error {
	for _, dir := range []string{p.GitmapRoot, p.CommitInRoot, p.ProfilesDir, p.TempRoot} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	return nil
}
