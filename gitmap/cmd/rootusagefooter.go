package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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
	if isShortFooterRequested() {
		printUsageFooterShort()

		return
	}

	printUsageFooterLong()
}

func isShortFooterRequested() bool {
	for _, arg := range os.Args {
		if arg == "--short-footer" || arg == "--short" {
			return true
		}
	}

	if os.Getenv("GITMAP_FOOTER") == "short" {
		return true
	}

	return false
}

func printUsageFooterLong() {
	printGitmapIdentityBlockLong()

	cwd, err := os.Getwd()
	if err != nil || !isFooterGitRepo(cwd) {
		return
	}

	printCurrentRepoIdentityBlock(cwd)
}

func printUsageFooterShort() {
	printGitmapIdentityBlockShort()
}

func printGitmapIdentityBlock() {
	if isShortFooterRequested() {
		printGitmapIdentityBlockShort()

		return
	}

	printGitmapIdentityBlockLong()
}

func printGitmapIdentityBlockLong() {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + footerRule + constants.ColorReset)
	fmt.Println("  " + constants.ColorMagenta + "gitmap binary" + constants.ColorReset)

	name := "gitmap"
	gitURL := resolveGitmapRepoURL()
	version := resolveGitmapVersion()
	sha := resolveGitmapCommitSHA()
	dbPath := resolveDatabasePath()
	installedPath := resolveInstalledBinaryPath()

	fmt.Printf("  %s● Name:%s           %s%s%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorWhite, name, constants.ColorReset)

	if len(gitURL) > 0 {
		fmt.Printf("  %s● Git URL:%s        %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorCyan, gitURL, constants.ColorReset)
	}

	fmt.Printf("  %s● Version:%s        %s%s%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorWhite, version, constants.ColorReset)

	if len(sha) > 0 {
		fmt.Printf("  %s● Commit SHA:%s     %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, sha, constants.ColorReset)
	}

	if len(dbPath) > 0 {
		fmt.Printf("  %s● Database:%s       %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorWhite, dbPath, constants.ColorReset)
	}

	if len(installedPath) > 0 {
		fmt.Printf("  %s● Installed path:%s %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorWhite, installedPath, constants.ColorReset)
	}

	printFooterBuildDate()
	fmt.Println()
}

func printGitmapIdentityBlockShort() {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + footerRule + constants.ColorReset)
	fmt.Println("  " + constants.ColorMagenta + "gitmap binary" + constants.ColorReset)

	version := resolveGitmapVersion()
	sha := resolveGitmapCommitSHA()

	fmt.Printf("  %s● Version:%s        %s%s%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorWhite, version, constants.ColorReset)

	if len(sha) > 0 {
		fmt.Printf("  %s● Commit SHA:%s     %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, sha, constants.ColorReset)
	}

	fmt.Println()
}

func resolveGitmapRepoURL() string {
	if len(BuildRepo) > 0 {
		return BuildRepo
	}

	srcURL := resolveRepoURLFromDir(gitmapSourceDir())
	if len(srcURL) > 0 {
		return srcURL
	}

	cwdURL := resolveRepoURLFromCWD()
	if len(cwdURL) > 0 {
		return cwdURL
	}

	return fmt.Sprintf("https://github.com/%s/%s", constants.UpdateRepoOwner, constants.UpdateCurrentRepoSlug)
}

func resolveRepoURLFromDir(dir string) string {
	if len(dir) == 0 {
		return ""
	}

	return captureGit(dir, "config", "--get", "remote.origin.url")
}

func resolveRepoURLFromCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	u := captureGit(cwd, "config", "--get", "remote.origin.url")
	if len(u) > 0 && strings.Contains(u, "gitmap") {
		return u
	}

	return ""
}

func resolveGitmapCommitSHA() string {
	if len(BuildCommit) > 0 {
		return BuildCommit
	}

	srcSHA := resolveCommitSHAFromDir(gitmapSourceDir())
	if len(srcSHA) > 0 {
		return srcSHA
	}

	cwdSHA := resolveCommitSHAFromCWD()
	if len(cwdSHA) > 0 {
		return cwdSHA
	}

	return readVersionJSONCommitSHA()
}

func resolveCommitSHAFromDir(dir string) string {
	if len(dir) == 0 {
		return ""
	}

	return captureGit(dir, "rev-parse", "HEAD")
}

func resolveCommitSHAFromCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	if !isFooterGitRepo(cwd) {
		return ""
	}

	u := captureGit(cwd, "config", "--get", "remote.origin.url")
	if !strings.Contains(u, "gitmap") {
		return ""
	}

	return captureGit(cwd, "rev-parse", "HEAD")
}

func readVersionJSONCommitSHA() string {
	candidatePaths := collectVersionJSONCandidatePaths()
	for _, p := range candidatePaths {
		sha := extractSHAFromVersionJSON(p)
		if len(sha) > 0 {
			return sha
		}
	}

	return ""
}

func collectVersionJSONCandidatePaths() []string {
	paths := []string{
		"version.json",
		filepath.Join(gitmapSourceDir(), "version.json"),
	}

	binPath := resolveInstalledBinaryPath()
	if len(binPath) == 0 {
		return paths
	}

	binDir := filepath.Dir(binPath)

	return append(paths,
		filepath.Join(binDir, "version.json"),
		filepath.Join(binDir, "..", "version.json"),
	)
}

func extractSHAFromVersionJSON(path string) string {
	if len(path) == 0 {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	var parsed struct {
		LastCommitSha string `json:"LastCommitSha"`
		Git           struct {
			Sha string `json:"sha"`
		} `json:"git"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return ""
	}

	if len(parsed.Git.Sha) > 0 {
		return parsed.Git.Sha
	}

	return parsed.LastCommitSha
}

func resolveInstalledBinaryPath() string {
	selfPath, err := os.Executable()
	if err != nil {
		return ""
	}

	resolved, err := filepath.EvalSymlinks(selfPath)
	if err != nil {
		resolved = selfPath
	}

	absPath, err := filepath.Abs(resolved)
	if err != nil {
		absPath = resolved
	}

	return absPath
}

func resolveDatabasePath() string {
	return store.DefaultDBPath()
}

func resolveGitmapVersion() string {
	v := constants.Version
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}

	return v
}

func printFooterBuildDate() {
	if len(BuildDate) > 0 {
		fmt.Printf("  %s● Built:%s          %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorDim, BuildDate, constants.ColorReset)
	}
}

func printFooterDatabase() {
	dbPath := resolveDatabasePath()
	if len(dbPath) > 0 {
		fmt.Printf("  %s● Database:%s       %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorWhite, dbPath, constants.ColorReset)
	}
}

func printCurrentRepoIdentityBlock(cwd string) {
	fmt.Println("  " + constants.ColorCyan + footerRule + constants.ColorReset)
	fmt.Println("  " + constants.ColorCyan + "current repo" + constants.ColorReset)
	emitIdentityRows(IdentityRowParams{
		Dir: cwd,
	})
	fmt.Println()
}

// IdentityRowParams encapsulates parameters for rendering an identity block.
type IdentityRowParams struct {
	Dir            string
	RepoOverride   string
	BranchOverride string
	ShaOverride    string
}

func resolveLocalRepoName(dir string) string {
	top := captureGit(dir, "rev-parse", "--show-toplevel")
	if len(top) > 0 {
		return filepath.Base(top)
	}

	return filepath.Base(dir)
}

func resolveLatestBranch(dir string) string {
	latest := captureGit(dir, "for-each-ref", "--sort=-committerdate", "refs/heads/", "--format=%(refname:short) (%(committerdate:relative))", "--count=1")
	if len(latest) > 0 {
		return latest
	}

	return captureGit(dir, "rev-parse", "--abbrev-ref", "HEAD")
}

func resolveOpenPRCount(dir string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "pr", "list", "--state", "open", "--limit", "100", "--json", "number", "--jq", "length")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "0"
	}

	count := strings.TrimSpace(string(out))
	if len(count) == 0 {
		return "0"
	}

	return count
}

func resolveBranchDirtyStatus(dir string) string {
	out := captureGit(dir, "status", "--porcelain")
	if len(out) == 0 {
		return "clean"
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	return fmt.Sprintf("dirty (%d changed)", len(lines))
}

func formatSyncCounts(ahead, behind string) string {
	if ahead == "0" && behind == "0" {
		return "up to date"
	}
	if ahead != "0" && behind == "0" {
		return fmt.Sprintf("ahead %s", ahead)
	}
	if ahead == "0" && behind != "0" {
		return fmt.Sprintf("behind %s", behind)
	}
	return fmt.Sprintf("ahead %s, behind %s", ahead, behind)
}

func resolveBranchSyncStatus(dir string) string {
	counts := captureGit(dir, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if len(counts) == 0 {
		return ""
	}

	parts := strings.Fields(counts)
	if len(parts) != 2 {
		return ""
	}

	return formatSyncCounts(parts[0], parts[1])
}

func resolveCurrentBranchInfo(dir string) string {
	dirtyStatus := resolveBranchDirtyStatus(dir)
	syncStatus := resolveBranchSyncStatus(dir)
	if len(syncStatus) > 0 {
		return dirtyStatus + " · " + syncStatus
	}

	return dirtyStatus
}

// emitIdentityRows prints Repo/Git URL/Branch/Latest branch/PR count/Current branch info/Last commit/Commit SHA rows for dir.
func emitIdentityRows(params IdentityRowParams) {
	repoName := resolveLocalRepoName(params.Dir)
	if len(repoName) > 0 {
		fmt.Printf("  %s● Repo:%s                %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorWhite, repoName, constants.ColorReset)
	}

	gitURL := firstNonEmptyVar(params.RepoOverride, captureGit(params.Dir, "config", "--get", "remote.origin.url"))
	if len(gitURL) > 0 {
		fmt.Printf("  %s● Git URL:%s             %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorCyan, gitURL, constants.ColorReset)
	}

	branch := firstNonEmptyVar(params.BranchOverride, captureGit(params.Dir, "rev-parse", "--abbrev-ref", "HEAD"))
	if len(branch) > 0 {
		fmt.Printf("  %s● Branch:%s              %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorGreen, branch, constants.ColorReset)
	}

	latestBranch := resolveLatestBranch(params.Dir)
	if len(latestBranch) > 0 {
		fmt.Printf("  %s● Latest branch:%s       %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorWhite, latestBranch, constants.ColorReset)
	}

	prCount := resolveOpenPRCount(params.Dir)
	fmt.Printf("  %s● PR count (open):%s     %s%s%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorYellow, prCount, constants.ColorReset)

	branchInfo := resolveCurrentBranchInfo(params.Dir)
	if len(branchInfo) > 0 {
		fmt.Printf("  %s● Current branch info:%s %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorWhite, branchInfo, constants.ColorReset)
	}

	if commit := captureGit(params.Dir, "log", "-1", "--format=%h · %s · %cr"); len(commit) > 0 {
		fmt.Printf("  %s● Last commit:%s         %s%s%s\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, commit, constants.ColorReset)
	}

	sha := firstNonEmptyVar(params.ShaOverride, captureGit(params.Dir, "rev-parse", "HEAD"))
	if len(sha) > 0 {
		fmt.Printf("  %s● Commit SHA:%s          %s%s%s\n",
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
