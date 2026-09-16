// Package cmdagy — agy_read_memory_prompt.go broadcasts the Read Memory protocol prompt to Antigravity projects.
package cmdagy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	agyAprmpExcept string
	agyAprmpPrompt string
	agyAprmpDryRun bool
	agyAprmpYes    bool
)

const defaultReadMemoryPrompt = "Execute enhanced Read Memory protocol. Defensively load memory, specs, constraints, and pending plans before taking action."

var agyAllProjectsReadMemoryCmd = &cobra.Command{
	Use:     "all-projects-read-memory-prompt",
	Aliases: []string{"aprmp", "read-memory-all", "all-read-memory"},
	Short:   "Broadcast Read Memory prompt to active projects",
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

	return processAgyPromptBroadcast(dirPath)
}

func processAgyPromptBroadcast(dirPath string) error {
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}

	targets, excluded := resolvePromptPartitions(projects)
	if len(targets) == 0 {
		fmt.Printf("%s No eligible active projects found to receive prompt.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)

		return nil
	}

	return executePromptBroadcast(targets, excluded)
}

func resolvePromptPartitions(projects []AgyProject) ([]AgyProject, []AgyProject) {
	tokens := parseAgyExceptTokens(agyAprmpExcept)

	return partitionPromptProjects(projects, tokens)
}
