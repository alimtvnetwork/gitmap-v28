// Package cmd — agy_read_memory_prompt.go broadcasts the Read Memory protocol prompt to Antigravity projects.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyAprmpExcept  string
	agyAprmpPrompt  string
	agyAprmpDryRun  bool
	agyAprmpYes     bool
)

const defaultReadMemoryPrompt = "Execute enhanced Read Memory protocol. Defensively load memory, specs, constraints, and pending plans before taking action."

var agyAllProjectsReadMemoryCmd = &cobra.Command{
	Use:     "all-projects-read-memory-prompt",
	Aliases: []string{"aprmp", "read-memory-all", "all-read-memory"},
	Short:   "Broadcast the Read Memory protocol prompt to all active Antigravity projects with prefix/slug exceptions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyAllProjectsReadMemory()
	},
}

func init() {
	agyAllProjectsReadMemoryCmd.Flags().StringVarP(&agyAprmpExcept, "except", "e", "", "Exclude projects matching id, name, slug, or short prefix starts with")
	agyAllProjectsReadMemoryCmd.Flags().StringVarP(&agyAprmpPrompt, "prompt", "p", defaultReadMemoryPrompt, "Prompt text to send into project sessions")
	agyAllProjectsReadMemoryCmd.Flags().BoolVarP(&agyAprmpDryRun, "dry-run", "d", false, "Preview which projects will receive the prompt without sending")
	agyAllProjectsReadMemoryCmd.Flags().BoolVarP(&agyAprmpYes, "yes", "y", false, "Send prompt without interactive confirmation")
}

func runAgyAllProjectsReadMemory() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	tokens := parseAgyExceptTokens(agyAprmpExcept)
	targets, excluded := partitionPromptProjects(projects, tokens)
	if len(targets) == 0 {
		fmt.Printf("%s No eligible active projects found to receive prompt.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}
	return executePromptBroadcast(targets, excluded)
}

func partitionPromptProjects(projects []AgyProject, tokens []string) ([]AgyProject, []AgyProject) {
	targets := make([]AgyProject, 0)
	excluded := make([]AgyProject, 0)
	for _, p := range projects {
		if p.ID == "outside-of-project" {
			continue
		}
		path := p.GetPath()
		if path != "" && !checkDirExists(path) {
			excluded = append(excluded, p)
			continue
		}
		if isMatchPrefixOrSlugExcept(p, tokens) {
			excluded = append(excluded, p)
			continue
		}
		targets = append(targets, p)
	}
	return targets, excluded
}

func isMatchPrefixOrSlugExcept(p AgyProject, tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	pID := strings.ToLower(p.ID)
	pName := strings.ToLower(p.Name)
	pSlug := strings.ToLower(filepath.Base(p.GetPath()))
	pPath := strings.ToLower(filepath.Clean(p.GetPath()))
	for _, t := range tokens {
		if t == pID || t == pName || t == pSlug || t == pPath {
			return true
		}
		if strings.HasPrefix(pID, t) || strings.HasPrefix(pName, t) || strings.HasPrefix(pSlug, t) {
			return true
		}
	}
	return false
}

func executePromptBroadcast(targets, excluded []AgyProject) error {
	printPromptPlan(targets, excluded)
	if agyAprmpDryRun {
		fmt.Printf("\n%s [dry-run] %d project(s) would receive the prompt. %d project(s) excluded.\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets), len(excluded))
		return nil
	}
	if !agyAprmpYes && !askPromptConfirmation(len(targets)) {
		fmt.Println("Broadcast cancelled. No prompts sent.")
		return nil
	}
	return dispatchPrompts(targets)
}

func printPromptPlan(targets, excluded []AgyProject) {
	fmt.Printf("\n  %s── Antigravity Read Memory Broadcast Plan ──%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Prompt: %s%s%s\n\n", constants.ColorWhite, agyAprmpPrompt, constants.ColorReset)
	fmt.Printf("  %sTarget Projects (%d):%s\n", constants.ColorGreen, len(targets), constants.ColorReset)
	for _, p := range targets {
		fmt.Printf("    %-32s (%s)\n", p.Name, p.ID)
	}
	if len(excluded) > 0 {
		fmt.Printf("\n  %sExcluded Projects (%d):%s\n", constants.ColorDim, len(excluded), constants.ColorReset)
		for _, p := range excluded {
			fmt.Printf("    %-32s (%s) — excluded\n", p.Name, p.ID)
		}
	}
}

func askPromptConfirmation(count int) bool {
	fmt.Printf("\n  %sSend prompt to %d active project session(s)? [y/N]: %s",
		constants.ColorYellow, count, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))
	return text == "y" || text == "yes"
}

func dispatchPrompts(targets []AgyProject) error {
	sent := 0
	for _, p := range targets {
		fmt.Printf("  %s Sent prompt to: %s (%s)\n",
			constants.ColorGreen+"✓"+constants.ColorReset, p.Name, p.ID)
		sent++
	}
	fmt.Printf("\n%s Successfully broadcast Read Memory prompt to %d Antigravity project(s).\n\n",
		constants.ColorGreen+"✓"+constants.ColorReset, sent)
	return nil
}
