package committransfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Replay walks plan.Commits in order and applies each one to the
// target. The source working dir is restored to its original ref on
// exit (success or failure) via deferred cleanup.
//
// In dry-run mode no writes happen — the caller has already printed
// the plan; Replay just returns a result with everything counted as
// skipped (dry-run is treated like "would replay 0").
func Replay(plan ReplayPlan, opts Options) (ReplayResult, error) {
	if opts.DryRun {
		return ReplayResult{}, nil
	}
	stopGuard := installInterruptGuard(plan.SourceDir, plan.SourceHEAD, opts.LogPrefix)
	defer stopGuard()
	defer func() { _ = checkoutRef(plan.SourceDir, plan.SourceHEAD) }()

	res := ReplayResult{}
	for i, commit := range plan.Commits {
		if commit.SkipCause != "" {
			tallySkip(&res, commit.SkipCause)

			continue
		}
		newSHA, emptyAfterSnapshot, err := replayOne(plan, commit, opts)
		if err != nil {
			return res, fmt.Errorf("commit %d/%d (%s): %w",
				i+1, len(plan.Commits), commit.ShortSHA, err)
		}
		if newSHA == "" {
			res.SkippedEmpty++
			if emptyAfterSnapshot {
				// Surface which source SHA disappeared so the user
				// can reconcile against `git log`. Issue:
				// .lovable/memory/issues/2026-05-09-commit-transfer-count-mismatch.md
				fmt.Fprintf(os.Stdout,
					"%s [%d/%d] %s → empty (snapshot tree unchanged on target)\n",
					opts.LogPrefix, i+1, len(plan.Commits), commit.ShortSHA)
			}

			continue
		}
		res.Replayed++
		res.NewSHAs = append(res.NewSHAs, newSHA)
	}

	return res, nil
}

// replayOne is the per-commit step: checkout in source, snapshot copy
// into target, stage, commit (preserving source author). The returned
// emptyAfterSnapshot flag distinguishes "intentional NoCommit skip"
// (false) from "snapshot produced no diff" (true) so the caller can
// log the latter — silent empty-skips were the root cause behind the
// 212→150 commit-count mismatch (issue 2026-05-09).
func replayOne(plan ReplayPlan, commit SourceCommit, opts Options) (string, bool, error) {
	if err := checkoutDetached(plan.SourceDir, commit.SHA); err != nil {
		return "", false, fmt.Errorf("checkout source %s: %w", commit.ShortSHA, err)
	}
	if err := snapshotCopy(plan.SourceDir, plan.TargetDir, opts); err != nil {
		return "", false, fmt.Errorf("snapshot copy: %w", err)
	}
	if opts.NoCommit {
		return "", false, nil
	}
	if err := addAll(plan.TargetDir); err != nil {
		return "", false, fmt.Errorf("git add -A target: %w", err)
	}
	if !hasStagedChanges(plan.TargetDir) {
		return "", true, nil
	}
	sha, err := commitWithEnv(plan.TargetDir, commit.Cleaned, commit.Author, commit.AuthorAt)

	return sha, false, err
}

// tallySkip routes a skip cause into the right counter.
func tallySkip(res *ReplayResult, cause string) {
	if cause == "already-replayed" {
		res.SkippedReplayed++

		return
	}
	if isDropSkip(cause) {
		res.SkippedDrop++

		return
	}
	res.SkippedEmpty++
}

// snapshotCopy walks source and copies each file into target, skipping
// .git/ (always) and node_modules/ (unless opts.IncludeNodeMod). When
// opts.Mirror is set, target-only files NOT present in source are
// removed before the copy.
func snapshotCopy(source, target string, opts Options) error {
	wanted := map[string]struct{}{}
	walkErr := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(source, path)
		if relErr != nil {
			return relErr
		}
		if shouldSkipPath(rel, opts) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
		if info.IsDir() {
			return nil
		}
		wanted[rel] = struct{}{}

		return copyOne(path, filepath.Join(target, rel), info)
	})
	if walkErr != nil {
		return walkErr
	}
	if opts.Mirror {
		return mirrorPrune(target, wanted, opts)
	}

	return nil
}

// shouldSkipPath returns true for paths the snapshot must ignore.
func shouldSkipPath(rel string, opts Options) bool {
	if rel == "." {
		return false
	}
	first := strings.SplitN(filepath.ToSlash(rel), "/", 2)[0]
	if first == ".git" && !opts.IncludeVCS {
		return true
	}
	if first == "node_modules" && !opts.IncludeNodeMod {
		return true
	}

	return false
}

// copyOne is a thin wrapper around io.Copy with mode preservation.
func copyOne(src, dst string, info os.FileInfo) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)

	return err
}

// mirrorPrune deletes files under target that are not in wanted. Only
// runs when opts.Mirror is set (spec §4 caveat).
func mirrorPrune(target string, wanted map[string]struct{}, opts Options) error {
	return filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == target {
			return err
		}
		rel, relErr := filepath.Rel(target, path)
		if relErr != nil {
			return relErr
		}
		if shouldSkipPath(rel, opts) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
		if info.IsDir() {
			return nil
		}
		if _, keep := wanted[rel]; !keep {
			return os.Remove(path)
		}

		return nil
	})
}
