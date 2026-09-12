package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

var rscSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
var rscSkipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c"))
var rscHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Bold(true)

func runReleaseScanCommits(args []string) error {
	isAll := parseRscArgs(args)
	cwd, head, err := resolveCwdAndHead()
	if err != nil {
		return apperror.WrapSimple(err, "runReleaseScanCommits")
	}

	return executeAndPersistScan(cwd, head, isAll)
}

func resolveCwdAndHead() (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", apperror.WrapSimple(err, "resolveCwdAndHead")
	}

	head, err := getGitHead(cwd)
	if err != nil {
		return "", "", apperror.WrapSimple(err, "resolveCwdAndHead")
	}

	return cwd, head, nil
}

func executeAndPersistScan(cwd, head string, isAll bool) error {
	commits, err := fetchCommits(cwd, isAll)
	if err != nil {
		return apperror.WrapSimple(err, "executeAndPersistScan")
	}

	actions, err := release.ExecuteCommitActions(cwd, commits)
	if err != nil {
		return apperror.WrapSimple(err, "executeAndPersistScan")
	}

	printScanCommitsSummary(actions)

	return release.WriteLastScannedCommit(cwd, head)
}

func parseRscArgs(args []string) bool {
	for _, a := range args {
		if a == "--all" || a == "-a" {
			return true
		}
	}

	return false
}

func getGitHead(cwd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "getGitHead")
	}

	return strings.TrimSpace(string(out)), nil
}

func fetchCommits(cwd string, isAll bool) ([]release.ParsedCommit, error) {
	rangeStr := buildCommitRange(cwd, isAll)
	cmdArgs := []string{"log", "--oneline"}
	if rangeStr != "" {
		cmdArgs = append(cmdArgs, rangeStr)
	}

	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "fetchCommits")
	}

	return parseGitLogLines(string(out)), nil
}

func buildCommitRange(cwd string, isAll bool) string {
	if isAll {
		return ""
	}

	lastHash, err := release.ReadLastScannedCommit(cwd)
	if err != nil || lastHash == "" {
		return ""
	}

	return lastHash + "..HEAD"
}

func parseGitLogLines(logOut string) []release.ParsedCommit {
	var commits []release.ParsedCommit
	lines := strings.Split(logOut, "\n")
	for _, line := range lines {
		if commit, isParsed := parseGitLogLine(line); isParsed {
			commits = append(commits, commit)
		}
	}

	return commits
}

func parseGitLogLine(line string) (release.ParsedCommit, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return release.ParsedCommit{}, false
	}

	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) < 2 {
		return release.ParsedCommit{}, false
	}

	ver, isFound := release.ParseVersionFromCommit(parts[1])
	if !isFound {
		return release.ParsedCommit{}, false
	}

	return release.ParsedCommit{Hash: parts[0], Message: parts[1], Version: ver}, true
}

func printScanCommitsSummary(actions []release.ScanCommitAction) {
	fmt.Println(rscHeaderStyle.Render(fmt.Sprintf("\n--- Found %d version bump commits ---", len(actions))))
	for _, a := range actions {
		fmt.Printf("Commit %s (version %s):\n", a.CommitHash, a.Version)
		printActionLine("Branch release/"+a.Version, a.IsBranchCreated, a.IsBranchSkipped)
		printActionLine("Tag "+a.Version, a.IsTagCreated, a.IsTagSkipped)
	}

	fmt.Println(rscHeaderStyle.Render("Done."))
}

func printActionLine(name string, isCreated, isSkipped bool) {
	if isCreated {
		fmt.Printf("  %s %s created\n", rscSuccessStyle.Render("✓"), name)

		return
	}

	if isSkipped {
		fmt.Printf("  %s %s skipped (already exists)\n", rscSkipStyle.Render("~"), name)

		return
	}
}
