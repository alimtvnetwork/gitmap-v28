package gitignoreagm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type cliOptions struct {
	targetDir       string
	isAutoYes       bool
	isAllRepos      bool
	isCommitEnabled bool
	isHelp          bool
}

// FilterReposNeedingRemediation returns the subset of repoDirs that have unignored or tracked resume task files.
func FilterReposNeedingRemediation(repoDirs []string) []string {
	seen := make(map[string]bool, len(repoDirs))
	var affected []string
	for _, dir := range repoDirs {
		clean := filepath.Clean(dir)
		if isEligibleAffectedRepo(clean, seen) {
			seen[clean] = true
			affected = append(affected, clean)
		}
	}
	return affected
}

func isEligibleAffectedRepo(cleanDir string, seen map[string]bool) bool {
	if cleanDir == "" || seen[cleanDir] {
		return false
	}
	return HasUnignoredResumeTask(cleanDir)
}

// CheckAndPromptRepos detects affected repos during scan or pull-all and prompts the user to remediate.
func CheckAndPromptRepos(repoDirs []string, isQuiet bool, isAutoYes bool) int {
	affected := FilterReposNeedingRemediation(repoDirs)
	if len(affected) == 0 {
		return 0
	}
	if isAutoYes {
		return RemediateBatch(affected, true, isQuiet)
	}
	if isQuiet || !isInteractiveTerminal() {
		printNonInteractiveNotice(len(affected), isQuiet)
		return 0
	}
	if !promptUserConfirmation(affected) {
		return 0
	}
	return RemediateBatch(affected, true, false)
}

func printNonInteractiveNotice(count int, isQuiet bool) {
	if isQuiet {
		return
	}
	fmt.Printf("  ℹ Detected %s in %d repo(s). Run 'gitmap gitignore agm' to untrack and ignore.\n", PrimaryIgnoreEntry, count)
}

func promptUserConfirmation(affected []string) bool {
	fmt.Printf("\n  ⚠ Detected %s in %d repository(ies):\n", PrimaryIgnoreEntry, len(affected))
	printAffectedRepoPreview(affected)
	fmt.Printf("  ? Delete from Git, add %s to .gitignore, and commit? [Y/n]: ", PrimaryIgnoreEntry)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "" || ans == "y" || ans == "yes"
}

func printAffectedRepoPreview(affected []string) {
	limit := 5
	for i, repo := range affected {
		if i >= limit {
			fmt.Printf("    ... and %d more\n", len(affected)-limit)
			return
		}
		fmt.Printf("    • %s\n", filepath.Base(repo))
	}
}

func isInteractiveTerminal() bool {
	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") != "" {
		return false
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// RemediateBatch applies RemediateRepo across all target repositories and returns the number of remediated repos.
func RemediateBatch(repoDirs []string, isCommitEnabled bool, isQuiet bool) int {
	remediatedCount := 0
	for _, dir := range repoDirs {
		res, err := RemediateRepo(dir, isCommitEnabled)
		if err == nil && (res.WasUntracked || res.WasFileDeleted || res.WasIgnored || res.WasCommitted) {
			remediatedCount++
			printRepoRemediationStatus(res, isQuiet)
		}
	}
	return remediatedCount
}

func printRepoRemediationStatus(res RepoRemediationResult, isQuiet bool) {
	if isQuiet {
		return
	}
	name := filepath.Base(res.RepoPath)
	if res.WasCommitted {
		fmt.Printf("  ✓ [%s] Untracked %s, added to .gitignore, and committed\n", name, PrimaryIgnoreEntry)
		return
	}
	fmt.Printf("  ✓ [%s] Added %s to .gitignore\n", name, PrimaryIgnoreEntry)
}

// RunCLI executes the `gitmap gitignore [agm|agy]` command.
func RunCLI(args []string) error {
	opts := parseCLIArgs(args)
	if opts.isHelp {
		printGitignoreHelp()
		return nil
	}
	absTarget, err := filepath.Abs(opts.targetDir)
	if err != nil {
		absTarget = opts.targetDir
	}
	if IsGitRepository(absTarget) && !opts.isAllRepos {
		return runSingleRepoCLI(absTarget, opts.isCommitEnabled)
	}
	return runWorkspaceReposCLI(absTarget, opts)
}

func runSingleRepoCLI(repoDir string, isCommitEnabled bool) error {
	res, err := RemediateRepo(repoDir, isCommitEnabled)
	if err != nil {
		return err
	}
	printSingleRepoSummary(res)
	return nil
}

func printSingleRepoSummary(res RepoRemediationResult) {
	name := filepath.Base(res.RepoPath)
	if !res.WasUntracked && !res.WasFileDeleted && !res.WasIgnored && !res.WasCommitted {
		fmt.Printf("  ✓ [%s] %s is already in .gitignore and clean.\n", name, PrimaryIgnoreEntry)
		return
	}
	printRepoRemediationStatus(res, false)
}

func runWorkspaceReposCLI(rootDir string, opts cliOptions) error {
	repos := discoverGitRepositories(rootDir)
	if len(repos) == 0 {
		fmt.Printf("  ⚠ No git repositories found under %s\n", rootDir)
		return nil
	}
	targets := selectWorkspaceTargets(repos, opts.isAllRepos)
	if len(targets) == 0 {
		fmt.Printf("  ✓ Scanned %d repositories under %s — no unignored %s files found.\n", len(repos), rootDir, PrimaryIgnoreEntry)
		return nil
	}
	count := RemediateBatch(targets, opts.isCommitEnabled, false)
	fmt.Printf("  ✓ Completed gitignore AGM remediation across %d/%d repository(ies).\n", count, len(targets))
	return nil
}

func selectWorkspaceTargets(repos []string, isAllRepos bool) []string {
	if isAllRepos {
		return repos
	}
	return FilterReposNeedingRemediation(repos)
}

func discoverGitRepositories(rootDir string) []string {
	var repos []string
	if IsGitRepository(rootDir) {
		repos = append(repos, rootDir)
	}
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return repos
	}
	for _, entry := range entries {
		repos = appendDiscoveredEntry(repos, rootDir, entry)
	}
	return repos
}

func appendDiscoveredEntry(repos []string, rootDir string, entry os.DirEntry) []string {
	if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
		return repos
	}
	childPath := filepath.Join(rootDir, entry.Name())
	if IsGitRepository(childPath) {
		return append(repos, childPath)
	}
	return repos
}

func parseCLIArgs(args []string) cliOptions {
	opts := cliOptions{targetDir: ".", isCommitEnabled: true}
	for _, arg := range args {
		applyCLIArg(&opts, arg)
	}
	return opts
}

func applyCLIArg(opts *cliOptions, arg string) {
	lower := strings.ToLower(strings.TrimSpace(arg))
	switch lower {
	case "-h", "--help", "help":
		opts.isHelp = true
	case "-y", "--yes":
		opts.isAutoYes = true
	case "-a", "--all", "all":
		opts.isAllRepos = true
	case "--no-commit":
		opts.isCommitEnabled = false
	case "agm", "agy", "antigravity", "resume", "ignore":
		// Subcommand selector token; default mode is AGM resume task ignore.
	default:
		if !strings.HasPrefix(arg, "-") {
			opts.targetDir = arg
		}
	}
}

func printGitignoreHelp() {
	fmt.Println("Usage: gitmap gitignore [agm|agy] [path] [flags]")
	fmt.Println("")
	fmt.Println("Untrack, delete, and ignore .antigravity_resume_task.json (and antigravity-resume_task.json)")
	fmt.Println("across one or more Git repositories, automatically staging and committing .gitignore.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gitmap gitignore agm                  # Untrack, ignore, and commit in current repo or workspace")
	fmt.Println("  gitmap gitignore agm D:\\work --all    # Ensure .gitignore entry across all repos in D:\\work")
	fmt.Println("  gitmap gitignore agy -y               # Non-interactive remediation and commit")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -a, --all         Apply .gitignore entry to all discovered repositories even if file is absent")
	fmt.Println("  -y, --yes         Automatically confirm all actions without prompting")
	fmt.Println("      --no-commit   Update .gitignore and untrack from index without creating a git commit")
	fmt.Println("  -h, --help        Show this help message")
}
