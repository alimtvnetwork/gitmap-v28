package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PullFailureSummary holds concise failure information for a repository.
type PullFailureSummary struct {
	RepoID    int64
	RepoName  string
	RepoPath  string
	RemoteURL string
	ErrorText string
	ErrorType string
}

// PromptInteractiveBatchFix displays an interactive prompt for remediating failed repositories.
func PromptInteractiveBatchFix(failures []PullFailureSummary) error {
	hasFailures := len(failures) > 0
	if !hasFailures || !isInteractiveTerminal() {
		return nil
	}
	for {
		printRemediationPromptMenu()
		choice := readPromptChoice()
		if isChoiceHelp(choice) {
			renderRemediationHelpGuide()
			continue
		}

		return dispatchRemediationChoice(choice, failures)
	}
}

func printRemediationPromptMenu() {
	fmt.Printf("\n  %s[?]%s %sRemediate failed repositories?%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset)
	fmt.Println("    [1/a] Resolve all failed repositories at once (auto-merge)")
	fmt.Println("    [2/s] Step through one by one")
	fmt.Println("    [3/r] Rebase all failed repositories (pull --rebase)")
	fmt.Println("    [h/?/help] Help & explanations of resolution options")
	fmt.Println("    [q/n] Skip / Exit")
	fmt.Printf("  %sChoice [1/2/3/h/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func isChoiceHelp(choice string) bool {
	return choice == "h" || choice == "help" || choice == "?"
}

func renderRemediationHelpGuide() {
	renderRemediationHelpHeader()
	renderAutoMergeExplanation()
	renderStepThroughExplanation()
	renderRebaseExplanation()
}

func renderRemediationHelpHeader() {
	fmt.Printf("\n  %s╔══════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║%s  %s📖 Remediation Options & Safety Guide%s                              %s║%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func renderAutoMergeExplanation() {
	fmt.Printf("  %s[1/a] Auto-Merge (Batch Resolution):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Operation: Attempts git pull with automatic conflict detection and")
	fmt.Println("      credential auto-conversion (HTTPS -> SSH) across all failed repos.")
	fmt.Println("    • Safety Guarantee: Non-destructive. Stashes dirty worktrees when needed.")
	fmt.Println("      Never force-pushes or discards uncommitted code.")
	fmt.Println("    • Risk Profile: Low. If merge conflicts arise, repository is safely")
	fmt.Println("      aborted and flagged for manual resolution.")
	fmt.Println()
}

func renderStepThroughExplanation() {
	fmt.Printf("  %s[2/s] Step-Through (Interactive Resolution):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Operation: Interactively iterates through each failure one by one,")
	fmt.Println("      allowing choice of auto-merge, rebase, skip, or quit per repository.")
	fmt.Println("    • Safety Guarantee: Complete operator oversight for every action.")
	fmt.Println("    • Risk Profile: Lowest. Recommended when repositories require different")
	fmt.Println("      handling strategies.")
	fmt.Println()
}

func renderRebaseExplanation() {
	fmt.Printf("  %s[3/r] Rebase (git pull --rebase):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Operation: Replays local unpushed commits on top of incoming remote")
	fmt.Println("      commits to maintain a linear git history.")
	fmt.Println("    • Safety Guarantee: Automatically aborts (git rebase --abort) on failure,")
	fmt.Println("      restoring repository to original working state.")
	fmt.Println("    • Risk Profile: Moderate. Rewrites local commit history. Not suited")
	fmt.Println("      for repositories with complex branch merges.")
	fmt.Println()
}

func readPromptChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "q"
	}

	return strings.ToLower(strings.TrimSpace(line))
}

func dispatchRemediationChoice(choice string, failures []PullFailureSummary) error {
	if isChoiceResolveAll(choice) {
		return RunBatchRemediateAll(failures)
	}
	if isChoiceStepByStep(choice) {
		return RunStepByStepRemediate(failures)
	}
	if isChoiceRebaseAll(choice) {
		return RunBatchRemediateRebase(failures)
	}

	return nil
}

func isChoiceResolveAll(choice string) bool {
	return choice == "1" || choice == "a" || choice == "all" || choice == "y" || choice == "yes"
}

func isChoiceStepByStep(choice string) bool {
	return choice == "2" || choice == "s" || choice == "step"
}

func isChoiceRebaseAll(choice string) bool {
	return choice == "3" || choice == "r" || choice == "rebase"
}

// RunBatchRemediateAll resolves all failed repositories.
func RunBatchRemediateAll(failures []PullFailureSummary) error {
	fmt.Printf("\n  %sResolving %d failed repositories...%s\n\n",
		constants.ColorBold, len(failures), constants.ColorReset)
	successCount := 0
	for _, f := range failures {
		if remediateSingleRepo(f) {
			successCount++
		}
	}
	printRemediationBatchSummary(successCount, len(failures))

	return nil
}

func printRemediationBatchSummary(successCount, total int) {
	fmt.Printf("\n  %s✓ Batch remediation complete: %d/%d resolved.%s\n\n",
		constants.ColorGreen, successCount, total, constants.ColorReset)
}

// RunStepByStepRemediate presents each failed repository to the user one by one.
func RunStepByStepRemediate(failures []PullFailureSummary) error {
	reader := bufio.NewReader(os.Stdin)
	successCount := 0
	for i, f := range failures {
		shouldContinue, isFixed := promptAndRemediateSingle(reader, i, len(failures), f)
		if isFixed {
			successCount++
		}
		if !shouldContinue {
			break
		}
	}
	printStepRemediationSummary(successCount, len(failures))

	return nil
}

func promptAndRemediateSingle(reader *bufio.Reader, idx, total int, f PullFailureSummary) (bool, bool) {
	printSingleRepoFailureHeader(idx+1, total, f)
	if cloner.IsMissingRepo(f.RepoPath) {
		return PromptMissingRepoAction(reader, f)
	}

	fmt.Printf("  %s?%s Remediate this repository? [Y/r/n/q] (Y=auto-merge, r=rebase, n=skip, q=quit): ", constants.ColorCyan, constants.ColorReset)
	line, _ := reader.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))
	if ans == "q" || ans == "quit" {
		return false, false
	}
	if ans == "n" || ans == "no" {
		return true, false
	}
	if ans == "r" || ans == "rebase" {
		return true, remediateSingleRepoRebase(f)
	}

	return true, remediateSingleRepo(f)
}

func printSingleRepoFailureHeader(idx, total int, f PullFailureSummary) {
	fmt.Printf("\n  [%d/%d] Repository: %s%s%s\n", idx, total, constants.ColorBold, f.RepoName, constants.ColorReset)
	fmt.Printf("    Path:  %s\n", f.RepoPath)
	fmt.Printf("    Error: %s\n", f.ErrorText)
}

func printStepRemediationSummary(successCount, total int) {
	fmt.Printf("\n  %s✓ Step-by-step remediation complete: %d/%d resolved.%s\n\n",
		constants.ColorGreen, successCount, total, constants.ColorReset)
}

func remediateSingleRepo(f PullFailureSummary) bool {
	if f.RepoPath == "" {
		return false
	}

	if cloner.IsMissingRepo(f.RepoPath) {
		return RemediateMissingRepo(f)
	}

	fmt.Printf("  • %s%s%s\n", constants.ColorBold, f.RepoName, constants.ColorReset)
	if isCredentialOrHTTPSFailure(f.ErrorText) || isSameOwnerHTTPSRepo(f.RepoPath) {
		remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName)
	}
	fmt.Printf("    %s Attempting auto-merge...\n", constants.TreeBranch)
	out, err := execPullAutoMerge(f.RepoPath)
	if err == nil {
		clearResolvedPullError(f.RepoName, f.RepoPath)
		fmt.Printf("    %s %s✓%s Successfully updated.\n",
			constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
		return true
	}

	if isCredentialOrHTTPSFailure(string(out)) {
		remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName)
		outRetry, errRetry := execPullAutoMerge(f.RepoPath)
		if errRetry == nil {
			clearResolvedPullError(f.RepoName, f.RepoPath)
			fmt.Printf("    %s %s✓%s Successfully updated via SSH.\n",
				constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
			return true
		}
		out = outRetry
	}

	return handleRemediationFailure(f.RepoPath, f.RepoName, string(out))
}

func isCredentialOrHTTPSFailure(msg string) bool {
	return isWincredmanFailure(msg) || isCacheFailure(msg) || isAuthFailure(msg) ||
		strings.Contains(strings.ToLower(msg), "credential")
}

func isSameOwnerHTTPSRepo(repoPath string) bool {
	originOut, err := exec.Command("git", "-C", repoPath, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return false
	}
	origin := strings.TrimSpace(string(originOut))
	if !strings.HasPrefix(strings.ToLower(origin), "https://github.com/") {
		return false
	}
	owner := extractGitHubOwnerFromURL(origin)
	if owner == "" {
		return false
	}
	return hasSSHRepoForSameOwner(owner)
}

func extractGitHubOwnerFromURL(url string) string {
	trimmed := strings.TrimSpace(url)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "https://github.com/") {
		rest := trimmed[len("https://github.com/"):]
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) >= 1 {
			return strings.TrimSpace(parts[0])
		}
	}
	if strings.HasPrefix(lower, "git@github.com:") {
		rest := trimmed[len("git@github.com:"):]
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) >= 1 {
			return strings.TrimSpace(parts[0])
		}
	}
	return ""
}

func hasSSHRepoForSameOwner(owner string) bool {
	if strings.EqualFold(owner, "alimtvnetwork") {
		return true
	}
	db, err := store.OpenDefault()
	if err != nil {
		return false
	}
	defer db.Close()
	repos, err := db.ListRepos()
	if err != nil {
		return false
	}
	prefix := "git@github.com:" + strings.ToLower(owner) + "/"
	for _, r := range repos {
		if strings.HasPrefix(strings.ToLower(r.SSHUrl), prefix) ||
			strings.HasPrefix(strings.ToLower(r.DiscoveredURL), prefix) {
			return true
		}
	}
	return false
}

func remediateCredentialAndConvertToSSH(repoPath, repoName string) {
	_ = exec.Command("git", "config", "--global", "credential.helper", "manager").Run()
	_ = exec.Command("git", "config", "--global", "credential.credentialStore", "wincredman").Run()

	changed, oldURL, newURL, err := ApplyTransportFlag(repoPath, true, false)
	if err == nil && changed {
		fmt.Printf("    %s Converted remote to SSH: %s -> %s\n", constants.TreeBranch, oldURL, newURL)
		updateRepoTransportInDB(repoPath, newURL)
	}
}

func updateRepoTransportInDB(repoPath, sshURL string) {
	db, err := store.OpenDefault()
	if err != nil || db.Conn() == nil {
		return
	}
	defer db.Close()
	_, _ = db.Conn().Exec(
		`UPDATE Repo SET SshUrl = ?, DiscoveredUrl = ?, IdentifiedTransport = 'ssh' WHERE AbsolutePath = ?`,
		sshURL, sshURL, repoPath,
	)
}

func clearResolvedPullError(repoName, repoPath string) {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return
	}
	defer db.Close()
	_ = db.ClearPullErrorsForRepo(repoName, repoPath)
}

func execPullAutoMerge(repoPath string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoPath, "pull", "--progress", "--no-rebase", "--no-edit", "--autostash")
	cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd.CombinedOutput()
}

func handleRemediationFailure(repoPath, repoName, out string) bool {
	relPath := formatRelativeOrCleanPath(repoPath)
	if relPath == "" || relPath == "." {
		relPath = repoName
	}

	isConflict := cloner.IsMergeConflictOutput(out) || strings.Contains(strings.ToLower(out), "conflict")
	if isConflict {
		cloner.SafeAbortMerge(repoPath)
		fmt.Printf("    %s %s✗%s Merge conflict detected (merge safely aborted to protect working tree).\n",
			constants.TreeCorner, constants.ColorRed, constants.ColorReset)
		fmt.Printf("        %s[hint]%s Resolve manually: git -C %s status\n",
			constants.ColorYellow, constants.ColorReset, relPath)

		return false
	}

	isDiverged := cloner.IsDivergedOutput(out)
	if isDiverged {
		cloner.SafeAbortMerge(repoPath)
		fmt.Printf("    %s %s✗%s Diverged branch detected (merge safely aborted to protect working tree).\n",
			constants.TreeCorner, constants.ColorRed, constants.ColorReset)
		fmt.Printf("        %s[hint]%s Resolve manually: git -C %s pull --rebase\n",
			constants.ColorYellow, constants.ColorReset, relPath)

		return false
	}

	cloner.SafeAbortMerge(repoPath)
	fmt.Printf("    %s %s✗%s Remediation failed.\n",
		constants.TreeCorner, constants.ColorRed, constants.ColorReset)
	fmt.Printf("        %s[hint]%s Resolve manually: git -C %s status\n",
		constants.ColorYellow, constants.ColorReset, relPath)

	return false
}

// RunBatchRemediateRebase resolves all failed repositories using pull --rebase.
func RunBatchRemediateRebase(failures []PullFailureSummary) error {
	fmt.Printf("\n  %sRebasing %d failed repositories...%s\n\n",
		constants.ColorBold, len(failures), constants.ColorReset)
	successCount := 0
	for _, f := range failures {
		if remediateSingleRepoRebase(f) {
			successCount++
		}
	}
	printRemediationBatchSummary(successCount, len(failures))

	return nil
}

func remediateSingleRepoRebase(f PullFailureSummary) bool {
	if f.RepoPath == "" {
		return false
	}

	if cloner.IsMissingRepo(f.RepoPath) {
		return RemediateMissingRepo(f)
	}

	fmt.Printf("  • %s%s%s\n", constants.ColorBold, f.RepoName, constants.ColorReset)
	fmt.Printf("    %s Attempting pull --rebase...\n", constants.TreeBranch)
	out, err := execPullRebase(f.RepoPath)
	if err == nil {
		fmt.Printf("    %s %s✓%s Successfully rebased and updated.\n",
			constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
		return true
	}

	return handleRebaseRemediationFailure(f.RepoPath, f.RepoName, string(out))
}

func execPullRebase(repoPath string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoPath, "pull", "--progress", "--rebase", "--autostash")
	cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd.CombinedOutput()
}

func handleRebaseRemediationFailure(repoPath, repoName, out string) bool {
	cloner.SafeAbortRebase(repoPath)
	relPath := formatRelativeOrCleanPath(repoPath)
	if relPath == "" || relPath == "." {
		relPath = repoName
	}

	isConflict := cloner.IsMergeConflictOutput(out) || strings.Contains(strings.ToLower(out), "conflict")
	if isConflict {
		fmt.Printf("    %s %s✗%s Rebase conflict detected (rebase safely aborted).\n",
			constants.TreeCorner, constants.ColorRed, constants.ColorReset)
		fmt.Printf("        %s[hint]%s Inspect rebase status: git -C %s status\n",
			constants.ColorYellow, constants.ColorReset, relPath)

		return false
	}

	fmt.Printf("    %s %s✗%s Rebase failed (rebase safely aborted).\n",
		constants.TreeCorner, constants.ColorRed, constants.ColorReset)
	fmt.Printf("        %s[hint]%s Inspect rebase status: git -C %s status\n",
		constants.ColorYellow, constants.ColorReset, relPath)

	return false
}

// RunBatchPullFix loads last recorded pull failures from SQLite and runs remediation.
func RunBatchPullFix(args []string) error {
	failures, err := loadRecentPullFailuresFromDB()
	if err != nil {
		return err
	}
	hasFailures := len(failures) > 0
	if !hasFailures {
		fmt.Printf("  %sℹ No recent pull failures found to remediate.%s\n",
			constants.ColorCyan, constants.ColorReset)
		return nil
	}
	fmt.Printf("  %sFound %d failed repository(ies) from recent runs.%s\n",
		constants.ColorCyan, len(failures), constants.ColorReset)

	return dispatchBatchFixArgs(args, failures)
}

func loadRecentPullFailuresFromDB() ([]PullFailureSummary, *apperror.AppError) {
	db, errDB := store.OpenPullSplitDB()
	if errDB != nil {
		return nil, apperror.WrapSimple(errDB, "open_pull_db")
	}
	defer db.Close()
	records, errQuery := db.QueryAllLatestPullErrors(50)
	if errQuery != nil {
		return nil, errQuery
	}

	return convertRecordsToFailures(records), nil
}

func convertRecordsToFailures(records []store.PullErrorRecord) []PullFailureSummary {
	var failures []PullFailureSummary
	seen := make(map[string]bool, len(records))
	for _, r := range records {
		hasPath := r.RepoPath != ""
		if hasPath && !seen[r.RepoPath] {
			seen[r.RepoPath] = true
			repoID, remoteURL := LookupRepoRemoteAndID(r.RepoPath)
			failures = append(failures, PullFailureSummary{
				RepoID:    repoID,
				RepoName:  r.RepoSlug,
				RepoPath:  r.RepoPath,
				RemoteURL: remoteURL,
				ErrorText: r.ErrorText,
				ErrorType: r.ErrorType,
			})
		}
	}

	return failures
}

func dispatchBatchFixArgs(args []string, failures []PullFailureSummary) error {
	if hasBatchFlag(args, "--rebase", "-r") {
		return RunBatchRemediateRebase(failures)
	}
	if hasBatchFlag(args, "--all", "-a", "--yes", "-y") {
		return RunBatchRemediateAll(failures)
	}
	if hasBatchFlag(args, "--step", "-s", "--one-by-one") {
		return RunStepByStepRemediate(failures)
	}

	return PromptInteractiveBatchFix(failures)
}

func hasBatchFlag(args []string, targets ...string) bool {
	for _, arg := range args {
		low := strings.ToLower(strings.TrimSpace(arg))
		for _, target := range targets {
			if low == target {
				return true
			}
		}
	}

	return false
}

// RemediateRepoByQuery resolves and remediates a single repository by name or path (including credential/SSH auto-fix).
func RemediateRepoByQuery(repoQuery string) (bool, error) {
	clean := strings.TrimSpace(repoQuery)
	if clean == "" {
		return false, nil
	}
	failures, _ := loadRecentPullFailuresFromDB()
	for _, f := range failures {
		if strings.EqualFold(f.RepoName, clean) || strings.EqualFold(f.RepoPath, clean) {
			if remediateSingleRepo(f) {
				return true, nil
			}
			return true, fmt.Errorf("remediation failed for %s", f.RepoName)
		}
	}
	records := resolveAllTrackedRecords()
	matched := FindBySlug(records, clean)
	if len(matched) == 1 {
		rec := matched[0]
		if isSameOwnerHTTPSRepo(rec.AbsolutePath) {
			f := PullFailureSummary{
				RepoName: rec.RepoName,
				RepoPath: rec.AbsolutePath,
			}
			if remediateSingleRepo(f) {
				return true, nil
			}
			return true, fmt.Errorf("remediation failed for %s", rec.RepoName)
		}
	}
	return false, nil
}

// HasPendingPullFailures reports whether gitmap-pull.db has recorded pull failures.
func HasPendingPullFailures() bool {
	failures, err := loadRecentPullFailuresFromDB()
	return err == nil && len(failures) > 0
}

// ExtractPullFailures collects failures from slice of repo states.
func ExtractPullFailures(states []*PullRepoState) []PullFailureSummary {
	var failures []PullFailureSummary
	for _, s := range states {
		if s.ErrorMsg != "" || s.Step == PullStepTypeError {
			repoID, remoteURL := LookupRepoRemoteAndID(s.RepoPath)
			failures = append(failures, PullFailureSummary{
				RepoID:    repoID,
				RepoName:  s.RepoName,
				RepoPath:  s.RepoPath,
				RemoteURL: remoteURL,
				ErrorText: s.ErrorMsg,
				ErrorType: string(s.Step),
			})
		}
	}

	return failures
}
