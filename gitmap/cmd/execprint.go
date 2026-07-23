package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// execInRepo runs a git command inside a single repo directory.
func execInRepo(rec model.ScanRecord, gitArgs []string) bool {
	cmd := exec.Command(constants.GitBin, gitArgs...)
	cmd.Dir = rec.AbsolutePath
	cmd.Stdout = nil
	cmd.Stderr = nil

	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	printExecResult(rec.RepoName, output, err)

	return err == nil
}

// printExecResult prints the success or failure line for one repo.
func printExecResult(name, output string, err error) {
	if err == nil {
		fmt.Printf(constants.ExecSuccessFmt, constants.ColorGreen, truncate(name, 22), constants.ColorReset)
	} else {
		fmt.Printf(constants.ExecFailFmt, constants.ColorYellow, truncate(name, 22), constants.ColorReset)
	}

	printExecOutput(output)
}

// printExecOutput prints indented command output lines.
func printExecOutput(output string) {
	if len(output) == 0 {
		return
	}
	for _, line := range strings.Split(output, "\n") {
		fmt.Printf(constants.ExecOutputLineFmt, constants.ColorDim, line, constants.ColorReset)
	}
}

// printExecBanner shows the command header.
func printExecBanner(gitArgs []string, count int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.ExecBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.ExecBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.ExecBannerBottom, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s"+constants.ExecCommandFmt+"%s\n", constants.ColorWhite, strings.Join(gitArgs, " "), constants.ColorReset)
	fmt.Printf("  %s"+constants.ExecRepoCountFmt+"%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	fmt.Println()
}

// printExecSummary shows final totals.
func printExecSummary(succeeded, failed, missing, total int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.ExecSummaryRule, constants.ColorReset)
	parts := buildExecSummaryParts(succeeded, failed, missing, total)
	line := strings.Join(parts, constants.SummaryJoinSep)
	fmt.Printf("  %s\n\n", line)
}

// buildExecSummaryParts assembles exec summary line segments.
func buildExecSummaryParts(succeeded, failed, missing, total int) []string {
	parts := []string{fmt.Sprintf(constants.SummaryReposFmt, total)}
	parts = appendSummaryPart(parts, succeeded, constants.ColorGreen, constants.SummarySucceededFmt)
	parts = appendSummaryPart(parts, failed, constants.ColorYellow, constants.SummaryFailedFmt)
	parts = appendSummaryPart(parts, missing, constants.ColorYellow, constants.SummaryMissingFmt)

	return parts
}
