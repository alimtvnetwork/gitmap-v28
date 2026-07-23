package replay

import (
	"fmt"
	"strings"
	"time"
)

// Plan is the input to ApplyCommit. It carries the resolved file set,
// the final commit message, and the identity/date pair that MUST be
// replicated byte-for-byte (spec guardrail: "Replicate BOTH AuthorDate
// AND CommitterDate byte-for-byte").
type Plan struct {
	SourceRepoDir string    // staged input dir we're copying FROM
	TargetRepoDir string    // <source> repo we're committing INTO
	SourceSha     string    // the commit being replayed (used for git show)
	Files         []string  // POSIX paths to copy from source @ SourceSha
	Message       string    // post-pipeline commit message
	AuthorName    string    // GIT_AUTHOR_NAME for the new commit
	AuthorEmail   string    // GIT_AUTHOR_EMAIL for the new commit
	AuthorDate    time.Time // GIT_AUTHOR_DATE — copied verbatim
	CommitterDate time.Time // GIT_COMMITTER_DATE — copied verbatim
}

// Result reports the outcome of ApplyCommit. NewSha is the 40-char
// SHA of the new commit; empty when DryRun was true.
type Result struct {
	NewSha string
}

// ApplyCommit performs spec §3.1 stage 14 (`Commit`) using git
// plumbing so we can pin BOTH dates to RFC3339 strings. Steps:
//  1. For each file in Plan.Files, copy its blob @ SourceSha into the
//     target's index using `git update-index --add --cacheinfo`.
//  2. `git write-tree` to materialize the index.
//  3. `git commit-tree <tree> -p <HEAD?>` with both date env vars set.
//  4. `git update-ref HEAD <newSha>`.
//
// dryRun=true short-circuits before any mutation; Result.NewSha = "".
func ApplyCommit(p Plan, dryRun bool) (Result, error) {
	if dryRun {
		return Result{}, nil
	}
	if err := stageFiles(p); err != nil {
		return Result{}, fmt.Errorf("replay: stage files: %w", err)
	}
	tree, err := writeTree(p.TargetRepoDir)
	if err != nil {
		return Result{}, fmt.Errorf("replay: write-tree: %w", err)
	}
	parent, _ := readHead(p.TargetRepoDir)
	newSha, err := commitTree(p, tree, parent)
	if err != nil {
		return Result{}, fmt.Errorf("replay: commit-tree: %w", err)
	}
	if err := updateHead(p.TargetRepoDir, newSha); err != nil {
		return Result{}, fmt.Errorf("replay: update-ref: %w", err)
	}
	return Result{NewSha: newSha}, nil
}

// stageFiles copies each path's blob from SourceRepoDir@SourceSha into
// TargetRepoDir's index. Uses `git cat-file blob <sha>:path` to read
// then `git hash-object -w --stdin --path=...` against the target.
func stageFiles(p Plan) error {
	for _, rel := range p.Files {
		if err := copyOneFile(p, rel); err != nil {
			return fmt.Errorf("file %s: %w", rel, err)
		}
	}
	return nil
}

// copyOneFile reads one blob from the source and writes it into the
// target's object DB + index. Mode 100644 (regular file) — symlinks
// and submodules are out of scope for v1.
func copyOneFile(p Plan, rel string) error {
	blob, err := gitRunnerBytes(p.SourceRepoDir, "cat-file", "blob", p.SourceSha+":"+rel)
	if err != nil {
		return fmt.Errorf("cat-file: %w", err)
	}
	hash, err := hashObjectStdin(p.TargetRepoDir, blob)
	if err != nil {
		return fmt.Errorf("hash-object: %w", err)
	}
	if _, err := gitRunner(p.TargetRepoDir, "update-index", "--add", "--cacheinfo", "100644,"+hash+","+rel); err != nil {
		return fmt.Errorf("update-index: %w", err)
	}
	return nil
}

// writeTree materializes the index as a tree object and returns its SHA.
func writeTree(target string) (string, error) {
	out, err := gitRunner(target, "write-tree")
	return strings.TrimSpace(out), err
}

// readHead returns the SHA HEAD points at; empty string when the repo
// has no commits yet (initial commit case).
func readHead(target string) (string, error) {
	out, err := gitRunner(target, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// commitTree builds a commit object via `git commit-tree`. The env-var
// machinery lives in env.go so this stays under the 15-line cap.
func commitTree(p Plan, tree, parent string) (string, error) {
	args := []string{"commit-tree", tree, "-m", p.Message}
	if parent != "" {
		args = append(args, "-p", parent)
	}
	out, err := gitRunnerEnv(p.TargetRepoDir, commitEnv(p), args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// updateHead points HEAD at newSha (works for both attached and
// detached HEAD; uses the symbolic ref dereference).
func updateHead(target, newSha string) error {
	_, err := gitRunner(target, "update-ref", "HEAD", newSha)
	return err
}
