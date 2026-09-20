package committransfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

		if isPRRouteEligible(plan.SourceDir, commit.SHA, commit.Subject, opts.PRMode) {
			prRes := ProcessPR(plan, commit, opts)
			if prRes.IsFailure() {
				return res, apperror.Wrap(prRes.Err, fmt.Sprintf("commit %d/%d (%s)",
					i+1, len(plan.Commits), commit.ShortSHA), nil)
			}
			res.Replayed++
			res.NewSHAs = append(res.NewSHAs, prRes.Value)

			continue
		}

		newSHA, emptyAfterSnapshot, err := replayOne(plan, commit, opts)
		if err != nil {
			return res, apperror.Wrap(err, fmt.Sprintf("commit %d/%d (%s)",
				i+1, len(plan.Commits), commit.ShortSHA), nil)
		}

		if newSHA == "" && emptyAfterSnapshot {
			res.SkippedEmpty++
			fmt.Fprintf(os.Stdout,
				"%s %s %s → empty (snapshot tree unchanged on target)\n",
				"  "+pterm.Magenta(opts.LogPrefix), pterm.Gray(fmt.Sprintf("[%d/%d]", i+1, len(plan.Commits))), pterm.Cyan(commit.ShortSHA))
			continue
		}

		if newSHA == "" {
			res.SkippedEmpty++
			continue
		}

		res.Replayed++
		res.NewSHAs = append(res.NewSHAs, newSHA)
	}

	syncErr := FinalizeSnapshotSync(plan.SourceDir, plan.TargetDir, plan.SourceHEAD, opts.CommandName)
	if syncErr != nil && !opts.NoCommit {
		fmt.Fprintf(os.Stderr, "%s final snapshot sync notice: %v\n", opts.LogPrefix, syncErr)
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
		return "", false, apperror.Wrap(err, "checkout source", map[string]any{"sha": commit.ShortSHA})
	}

	if err := snapshotCopy(plan.SourceDir, plan.TargetDir, opts); err != nil {
		return "", false, apperror.WrapSimple(err, "snapshot copy")
	}

	if opts.NoCommit {
		return "", false, nil
	}

	if err := addAll(plan.TargetDir); err != nil {
		return "", false, apperror.WrapSimple(err, "git add target")
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

		if isSkippablePath(rel, opts) && info.IsDir() {
			return filepath.SkipDir
		}

		if isSkippablePath(rel, opts) {
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

// isSkippablePath returns true for paths the snapshot must ignore.
func isSkippablePath(rel string, opts Options) bool {
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

		if isSkippablePath(rel, opts) && info.IsDir() {
			return filepath.SkipDir
		}

		if isSkippablePath(rel, opts) {
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

// FinalizeSnapshotSync ensures target working tree exactly matches source repo snapshot.
func FinalizeSnapshotSync(sourceDir, targetDir, sourceHeadSha, cmdName string) error {
	opts := Options{CommandName: cmdName, Mirror: true}
	if err := snapshotCopy(sourceDir, targetDir, opts); err != nil {
		return apperror.WrapSimple(err, "finalizeSnapshotSync.copy")
	}

	return commitFinalSnapshotDiff(targetDir, sourceHeadSha, cmdName)
}

func commitFinalSnapshotDiff(targetDir, sourceHeadSha, cmdName string) error {
	if err := addAll(targetDir); err != nil {
		return apperror.WrapSimple(err, "finalizeSnapshotSync.addAll")
	}
	if !hasStagedChanges(targetDir) {
		LogSnapshotSynced(os.Stdout, "[sync]", targetDir)

		return nil
	}

	shortSha := sourceHeadSha
	if len(shortSha) > 7 {
		shortSha = shortSha[:7]
	}
	msg := fmt.Sprintf("chore(sync): synchronize final repository snapshot tree to match source %s\n\ngitmap-replay-final-snapshot: %s\ngitmap-replay-cmd: %s",
		shortSha, sourceHeadSha, cmdName)

	_, err := gitOut(targetDir, "commit", "-m", msg)
	if err == nil {
		LogSnapshotSynced(os.Stdout, "[sync]", targetDir)
	}

	return err
}
