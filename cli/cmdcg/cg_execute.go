// Package cmdcg provides streaming AI execution prompts for on-the-fly model learning.
package cmdcg

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/charmbracelet/lipgloss"
)

var (
	executeHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9")).Padding(0, 1)
	executeTipStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Italic(true)
)

// RunCGExecutePrompt streams either the V6 or V7 master AI prompt to terminal stdout.
func RunCGExecutePrompt(version string, args []string) error {
	normalizedVer := strings.ToLower(strings.TrimSpace(version))
	if normalizedVer != "v6" && normalizedVer != "v7" {
		return apperror.NewValidationError(fmt.Sprintf("unsupported execute prompt version: %q (expected v6 or v7)", version))
	}

	isRaw := slices.Contains(args, "--raw") || slices.Contains(args, "-raw")

	promptPath, err := resolveExecutePromptPath(normalizedVer)
	if err != nil {
		return err
	}

	contentBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return apperror.WrapSimple(err, "failed to read execute prompt file")
	}

	if isRaw {
		fmt.Print(string(contentBytes))
		return nil
	}

	printExecuteBanner(normalizedVer, promptPath)
	fmt.Println(string(contentBytes))

	return nil
}

func resolveExecutePromptPath(version string) (string, error) {
	fileName := fmt.Sprintf("02-execute-parent-task-with-n-steps-%s.md", version)
	relCandidate := filepath.Join("01-prompts", "14-execute", fileName)

	if _, err := os.Stat(relCandidate); err == nil {
		return relCandidate, nil
	}

	sourceRoot, err := cmdsync.ResolveSourceRoot()
	if err != nil {
		return "", apperror.NewNotFoundError(fmt.Sprintf("execute prompt %s not found", fileName))
	}

	upstreamCandidate := filepath.Join(sourceRoot, relCandidate)
	if _, statErr := os.Stat(upstreamCandidate); statErr == nil {
		return upstreamCandidate, nil
	}

	return "", apperror.NewNotFoundError(fmt.Sprintf("execute prompt %s not found", fileName))
}

func printExecuteBanner(version string, promptPath string) {
	title := fmt.Sprintf("=== GitMap Master AI Execution Standard (%s) ===", strings.ToUpper(version))
	fmt.Println(executeHeaderStyle.Render(title))
	fmt.Printf("Source: %s\n", promptPath)
	fmt.Println(executeTipStyle.Render("► Learn this execution standard and follow through during task execution:\n"))
}
