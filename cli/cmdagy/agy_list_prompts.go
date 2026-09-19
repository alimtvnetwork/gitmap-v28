package cmdagy

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	listPromptsAllProjects   bool
	listPromptsProjectsCount int
	listPromptsProjectPrefix string
	listPromptsJSON          bool
)

var agyListPromptsCmd = &cobra.Command{
	Use:   "list-prompts [N]",
	Short: "List prompt history per project or across projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyListPrompts(args)
		if appErr != nil {
			return appErr
		}

		return nil
	},
}

func init() {
	agyListPromptsCmd.Flags().BoolVar(&listPromptsAllProjects, "all-projects", false, "List prompts across all projects")
	agyListPromptsCmd.Flags().IntVar(&listPromptsProjectsCount, "projects", 0, "Showcase prompt changes across last N projects in VS Code")
	agyListPromptsCmd.Flags().StringVar(&listPromptsProjectPrefix, "project", "", "Filter prompts by project name prefix")
	agyListPromptsCmd.Flags().BoolVarP(&listPromptsJSON, "json", "j", false, "Output prompt list as structured JSON")
}

// RunListPromptsTopLevelCLI executes agy list-prompts from top-level gitmap aliases.
func RunListPromptsTopLevelCLI(args []string) error {
	runArgs := append([]string{"list-prompts"}, args...)
	AgyCmd.SetArgs(runArgs)

	return AgyCmd.Execute()
}

func runAgyListPrompts(args []string) *apperror.AppError {
	limit := parsePromptLimit(args)
	if listPromptsProjectsCount > 0 {
		return launchVSCodeForLastProjects(listPromptsProjectsCount, limit)
	}

	err := displayPromptListing(limit)
	if err != nil {
		return apperror.WrapSimple(err, "display prompt listing")
	}

	return nil
}

func parsePromptLimit(args []string) int {
	cleanArgs := filterOutListPromptFlags(args)
	if len(cleanArgs) == 0 {
		return 50
	}

	if strings.ToLower(cleanArgs[0]) == "all" {
		return 10000
	}

	val, err := strconv.Atoi(cleanArgs[0])
	if err != nil || val < 1 {
		return 50
	}

	return val
}

func filterOutListPromptFlags(args []string) []string {
	var clean []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			clean = append(clean, a)
		}
	}

	return clean
}
