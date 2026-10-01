package cmdcpar

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunCPAR executes commit-push-all-repos across all dirty repositories.
func RunCPAR(args []string) *apperror.AppError {
	opts := parseCPAROptions(args)
	dirtyList := CollectDirtyRepositories()
	if len(dirtyList) == 0 {
		fmt.Printf("%s✓ All repositories are clean.%s\n", constants.ColorGreen, constants.ColorReset)
		return nil
	}
	if isReviewAborted(dirtyList, opts) {
		fmt.Println("Aborted by user.")
		return nil
	}
	executeCPARAcrossDirty(dirtyList, opts)
	return nil
}

func isReviewAborted(dirtyList []DirtyRepoSummary, opts cparOptions) bool {
	if !opts.isReview {
		return false
	}
	isProceed := handleReviewFlow(dirtyList, opts)
	return !isProceed
}

func parseCPAROptions(args []string) cparOptions {
	opts := cparOptions{commitMsg: "chore: commit pending changes"}
	for _, a := range args {
		opts = applyOption(opts, strings.ToLower(a))
	}
	return opts
}

func applyOption(opts cparOptions, a string) cparOptions {
	if a == "-y" || a == "--yes" {
		opts.isAutoYes = true
	} else if a == "-r" || a == "--review" {
		opts.isReview = true
	} else if a == "-co" || a == "--commit-only" {
		opts.isCommitOnly = true
	} else if strings.HasPrefix(a, "-m=") {
		opts.commitMsg = a[3:]
	}
	return opts
}

func CollectDirtyRepositories() []DirtyRepoSummary {
	records := resolveAllRecords()
	var dirty []DirtyRepoSummary
	for _, r := range records {
		diag := gitutil.InspectDirtyState(r.RelativePath)
		if diag.IsDirty {
			dirty = appendDirty(dirty, r, diag)
		}
	}
	return dirty
}

func appendDirty(dirty []DirtyRepoSummary, r model.ScanRecord, diag gitutil.DirtyDiagnosis) []DirtyRepoSummary {
	return append(dirty, DirtyRepoSummary{
		RepoName:  r.RepoName,
		RepoPath:  r.RelativePath,
		Diagnosis: diag,
	})
}

func handleReviewFlow(dirtyList []DirtyRepoSummary, opts cparOptions) bool {
	renderDirtyReviewTable(dirtyList)
	if opts.isAutoYes {
		return true
	}
	return promptReviewConsent(opts.isCommitOnly)
}

func resolveAllRecords() []model.ScanRecord {
	records := queryStoreRecords()
	if len(records) > 0 {
		return records
	}
	return []model.ScanRecord{{RelativePath: ".", RepoName: "."}}
}

func queryStoreRecords() []model.ScanRecord {
	s, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer s.Close()
	records, _ := s.ListRepos()
	return records
}

func renderDirtyReviewTable(dirty []DirtyRepoSummary) {
	fmt.Printf("\n%s  === PENDING COMMITS REVIEW (%d repositories) ===%s\n",
		constants.ColorCyan, len(dirty), constants.ColorReset)
	for _, d := range dirty {
		printDirtyRow(d)
	}
	fmt.Println()
}

func printDirtyRow(d DirtyRepoSummary) {
	fmt.Printf("    %-30s | modified: %d | untracked: %d | staged: %d\n",
		d.RepoName, d.Diagnosis.ModifiedCount, d.Diagnosis.UntrackedCount, d.Diagnosis.StagedCount)
}

func promptReviewConsent(isCommitOnly bool) bool {
	action := "commit and push"
	if isCommitOnly {
		action = "commit (without push)"
	}
	fmt.Printf("Proceed to %s across these repositories? [Y/n]: ", action)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	low := strings.ToLower(strings.TrimSpace(text))
	return isConsentGranted(low)
}

func isConsentGranted(low string) bool {
	if low == "" || low == "y" || low == "yes" {
		return true
	}
	return false
}

func executeCPARAcrossDirty(dirty []DirtyRepoSummary, opts cparOptions) {
	fmt.Printf("\n%s  ▶ Executing CPAR across %d dirty repo(s)...%s\n",
		constants.ColorCyan, len(dirty), constants.ColorReset)
	for _, d := range dirty {
		commitAndMaybePush(d, opts)
	}
	fmt.Println()
}

func commitAndMaybePush(d DirtyRepoSummary, opts cparOptions) {
	_ = exec.Command("git", "-C", d.RepoPath, "add", "-A").Run()
	cmd := exec.Command("git", "-C", d.RepoPath, "commit", "-m", opts.commitMsg)
	cmd.Env = gitutil.BuildSafeGitEnv()
	if err := cmd.Run(); err != nil {
		fmt.Printf("      %-30s %scommit failed%s\n", d.RepoName, constants.ColorRed, constants.ColorReset)
		return
	}
	handlePush(d, opts)
}

func handlePush(d DirtyRepoSummary, opts cparOptions) {
	if opts.isCommitOnly {
		fmt.Printf("      %-30s %scommitted (commit-only)%s\n", d.RepoName, constants.ColorGreen, constants.ColorReset)
		return
	}
	pushCmd := exec.Command("git", "-C", d.RepoPath, "push")
	pushCmd.Env = gitutil.BuildSafeGitEnv()
	if pushErr := pushCmd.Run(); pushErr != nil {
		fmt.Printf("      %-30s %scommitted (push failed)%s\n", d.RepoName, constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("      %-30s %scommitted & pushed%s\n", d.RepoName, constants.ColorGreen, constants.ColorReset)
}
