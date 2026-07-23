package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Build-time identity, injected via:
//
//	go build -ldflags "-X github.com/.../gitmap/cmd.BuildCommit=<sha> \
//	                  -X github.com/.../gitmap/cmd.BuildBranch=<branch> \
//	                  -X github.com/.../gitmap/cmd.BuildRepo=<origin-url> \
//	                  -X github.com/.../gitmap/cmd.BuildDate=<utc>"
//
// All four default to "" so unset values fall back to a runtime git probe
// against `constants.RepoPath` (the source repo baked in at link time).
var (
	BuildCommit = ""
	BuildBranch = ""
	BuildRepo   = ""
	BuildDate   = ""
)

const footerRule = "────────────────────────────────────────────────────────────"

// printUsageFooter renders TWO clearly separated identity blocks at the
// bottom of `gitmap` (no args) and `gitmap help`:
//
//  1. gitmap binary identity (which build is running) — magenta header
//  2. current repo identity (where you are right now)  — cyan header
//
// Blocks are separated by a blank line + thin rule so users can never
// confuse "what gitmap am I running" with "what repo am I sitting in".
func printUsageFooter() {
	printGitmapIdentityBlock()

	cwd, err := os.Getwd()
	if err != nil || !isFooterGitRepo(cwd) {
		return
	}
	if sameRepo(cwd, gitmapSourceDir()) {
		return // avoid duplicate block when cwd IS the gitmap source repo
	}
	printCurrentRepoIdentityBlock(cwd)
}

func printGitmapIdentityBlock() {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + footerRule + constants.ColorReset)
	fmt.Println("  " + constants.ColorMagenta + "gitmap binary" + constants.ColorReset)

	fmt.Printf("  %s● Version:%s     %sv%s%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorWhite, constants.Version, constants.ColorReset)

	src := gitmapSourceDir()
	emitIdentityRows(src, BuildRepo, BuildBranch, BuildCommit)
	if len(BuildDate) > 0 {
		fmt.Printf("  %s● Built:%s       %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorDim, BuildDate, constants.ColorReset)
	}
	fmt.Println()
}

func printCurrentRepoIdentityBlock(cwd string) {
	fmt.Println("  " + constants.ColorCyan + footerRule + constants.ColorReset)
	fmt.Println("  " + constants.ColorCyan + "current repo" + constants.ColorReset)
	emitIdentityRows(cwd, "", "", "")
	fmt.Println()
}

// emitIdentityRows prints Repo/Branch/Last commit/Commit SHA rows for dir,
// preferring the supplied build-time overrides when non-empty.
func emitIdentityRows(dir, repoOverride, branchOverride, shaOverride string) {
	repo := firstNonEmptyVar(repoOverride, captureGit(dir, "config", "--get", "remote.origin.url"))
	if len(repo) > 0 {
		fmt.Printf("  %s● Repo:%s        %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorCyan, repo, constants.ColorReset)
	}

	branch := firstNonEmptyVar(branchOverride, captureGit(dir, "rev-parse", "--abbrev-ref", "HEAD"))
	if len(branch) > 0 {
		fmt.Printf("  %s● Branch:%s      %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorGreen, branch, constants.ColorReset)
	}

	if commit := captureGit(dir, "log", "-1", "--format=%h · %s · %cr"); len(commit) > 0 {
		fmt.Printf("  %s● Last commit:%s %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, commit, constants.ColorReset)
	}

	sha := firstNonEmptyVar(shaOverride, captureGit(dir, "rev-parse", "HEAD"))
	if len(sha) > 0 {
		fmt.Printf("  %s● Commit SHA:%s  %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, sha, constants.ColorReset)
	}
}

// gitmapSourceDir returns the source repo baked into the binary, or "".
func gitmapSourceDir() string {
	if len(constants.RepoPath) == 0 {
		return ""
	}
	if _, err := os.Stat(filepath.Join(constants.RepoPath, ".git")); err != nil {
		return ""
	}

	return constants.RepoPath
}

// isFooterGitRepo reports whether dir (or any ancestor) is a git repo.
func isFooterGitRepo(dir string) bool {
	out := captureGit(dir, "rev-parse", "--is-inside-work-tree")

	return out == "true"
}

// sameRepo reports whether a and b resolve to the same git toplevel.
func sameRepo(a, b string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	ta := captureGit(a, "rev-parse", "--show-toplevel")
	tb := captureGit(b, "rev-parse", "--show-toplevel")

	return len(ta) > 0 && filepath.Clean(ta) == filepath.Clean(tb)
}

func firstNonEmptyVar(values ...string) string {
	for _, v := range values {
		if len(v) > 0 {
			return v
		}
	}

	return ""
}

// captureGit runs `git <args...>` in dir and returns trimmed stdout, or
// "" on any error. Stderr is discarded so the footer stays quiet when
// the directory is not a git repo.
//
// IMPORTANT: an empty dir is rejected up front. Without this guard, exec
// inherits the process CWD, which caused the "gitmap binary" footer to
// silently print the CURRENT repo's git identity whenever the source
// repo bake-in (constants.RepoPath) was missing — making the binary
// block indistinguishable from the "current repo" block (see v5.60.0).
func captureGit(dir string, args ...string) string {
	if len(dir) == 0 {
		return ""
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
