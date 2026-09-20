package prclean

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/prdb"
	"github.com/pterm/pterm"
)

// RunPRClean inspects merged PR branches in SQLite, renders a pre-flight summary table,
// and deletes branches upon confirmation (or immediately if isYes is true).
func RunPRClean(repoRoot string, isYes bool) *apperror.AppError {
	slug := prdb.SanitizeRepoSlug(filepath.Base(repoRoot))
	dbRes := prdb.OpenPrSplitDb(slug, repoRoot)
	if dbRes.IsFailure() {
		return dbRes.Err
	}
	db := dbRes.Value
	defer db.Close()

	branchesRes := db.ListMergedPrBranches()
	if branchesRes.IsFailure() {
		return branchesRes.Err
	}

	return processCleanBranches(db, repoRoot, branchesRes.Value, isYes)
}

func processCleanBranches(db *prdb.PrSplitDb, repoRoot string, branches []prdb.PrBranchRecord, isYes bool) *apperror.AppError {
	if len(branches) == 0 {
		pterm.Info.Println("No merged PR branches eligible for removal.")

		return nil
	}

	renderPreflightTable(branches)
	if !isYes && !confirmPrune(len(branches)) {
		pterm.Warning.Println("PR branch clean canceled by user.")

		return nil
	}

	return executePruneBranches(db, repoRoot, branches)
}

func renderPreflightTable(branches []prdb.PrBranchRecord) {
	pterm.DefaultHeader.WithFullWidth().Println("PR Branch Cleanup: Pre-Flight Report")
	data := pterm.TableData{
		{"Branch Name", "Type", "Status", "Merged"},
	}
	for _, b := range branches {
		data = append(data, []string{
			pterm.Yellow(b.BranchName),
			b.BranchType,
			pterm.Green("Merged"),
			formatMergedTime(b.MergedAt),
		})
	}
	_ = pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	fmt.Println()
}

func formatMergedTime(t int64) string {
	if t <= 0 {
		return "N/A"
	}

	return fmt.Sprintf("%d (epoch)", t)
}

func confirmPrune(count int) bool {
	prompt := fmt.Sprintf("Remove %d closed/merged PR branches? [y/N]: ", count)
	fmt.Print(pterm.LightMagenta(prompt))
	var response string
	_, _ = fmt.Scanln(&response)
	resp := strings.TrimSpace(strings.ToLower(response))

	return resp == "y" || resp == "yes"
}

func executePruneBranches(db *prdb.PrSplitDb, repoRoot string, branches []prdb.PrBranchRecord) *apperror.AppError {
	deletedCount := 0
	for _, b := range branches {
		cmd := exec.Command("git", "-C", repoRoot, "branch", "-D", b.BranchName)
		if out, err := cmd.CombinedOutput(); err != nil {
			pterm.Warning.Printf("Could not delete branch %s: %s\n", b.BranchName, strings.TrimSpace(string(out)))
		} else {
			deletedCount++
			pterm.Success.Printf("Deleted branch: %s\n", b.BranchName)
		}
		_ = db.MarkPrBranchDeleted(b.BranchName)
	}
	pterm.Success.Printf("Successfully cleaned %d PR branches.\n", deletedCount)

	return nil
}
