package cmdignore

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RunSeeIssues implements c gitignore issues (ig).
func RunSeeIssues(args []string) *apperror.AppError {
	records := resolveTargetRepos()
	isZero := len(records) == 0
	if isZero {
		fmt.Println("No repositories to check.")
		return nil
	}
	return printFoundIssues(records)
}

func printFoundIssues(records []model.ScanRecord) *apperror.AppError {
	issues := scanReposForIgnoreIssues(records)
	isClean := len(issues) == 0
	if isClean {
		fmt.Println("No gitignore issues found.")
		return nil
	}
	fmt.Printf("%sFound %d issues.%s\n", constants.ColorYellow, len(issues), constants.ColorReset)
	for _, issue := range issues {
		fmt.Printf(" - %s\n", issue.RepoName)
	}
	return nil
}
