package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
	"github.com/charmbracelet/lipgloss"
)

var rscSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
var rscSkipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c"))
var rscHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Bold(true)

func runReleaseScanCommits(args []string) {
		isAll := parseRscArgs(args)
	cwd, err := os.Getwd()
	if err != nil { os.Exit(1) }
	head, err := getGitHead(cwd)
	if err != nil { os.Exit(1) }
	
	commits, err := fetchCommits(cwd, isAll)
	if err != nil { os.Exit(1) }
	
	actions, err := release.ExecuteCommitActions(cwd, commits)
	if err != nil { os.Exit(1) }
	
	printScanCommitsSummary(actions)
	_ = release.WriteLastScannedCommit(cwd, head)
}

func parseRscArgs(args []string) bool {
	for _, a := range args {
		if a == "--all" || a == "-a" { return true }
	}
	return false
}

func getGitHead(cwd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil { return "", apperror.Wrap(err, "getGitHead", nil) }
	return strings.TrimSpace(string(out)), nil
}

func fetchCommits(cwd string, isAll bool) ([]release.ParsedCommit, error) {
	lastHash, _ := release.ReadLastScannedCommit(cwd)
	rangeStr := ""
	if !isAll && lastHash != "" { rangeStr = lastHash + "..HEAD" }
	cmdArgs := []string{"log", "--oneline"}
	if rangeStr != "" { cmdArgs = append(cmdArgs, rangeStr) }
	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil { return nil, apperror.Wrap(err, "fetchCommits", nil) }
	return parseGitLogLines(string(out)), nil
}

func parseGitLogLines(logOut string) []release.ParsedCommit {
	var commits []release.ParsedCommit
	lines := strings.Split(logOut, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" { continue }
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 { continue }
		if ver, isFound := release.ParseVersionFromCommit(parts[1]); isFound {
			commits = append(commits, release.ParsedCommit{Hash: parts[0], Message: parts[1], Version: ver})
		}
	}
	return commits
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
	} else if isSkipped {
		fmt.Printf("  %s %s skipped (already exists)\n", rscSkipStyle.Render("~"), name)
	}
}

