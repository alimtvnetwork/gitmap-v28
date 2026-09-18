package cmdagy

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listPromptsAllProjects bool
	listPromptsProjectsCount int
	listPromptsProjectPrefix string
	listPromptsJSON          bool
)

var agyListPromptsCmd = &cobra.Command{
	Use:   "list-prompts [N]",
	Short: "List prompt history per project or across projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyListPrompts(args)
	},
}

func init() {
	agyListPromptsCmd.Flags().BoolVar(&listPromptsAllProjects, "all-projects", false, "List prompts across all projects")
	agyListPromptsCmd.Flags().IntVar(&listPromptsProjectsCount, "projects", 0, "Showcase prompt changes across last N projects in VS Code")
	agyListPromptsCmd.Flags().StringVar(&listPromptsProjectPrefix, "project", "", "Filter prompts by project name prefix")
	agyListPromptsCmd.Flags().BoolVarP(&listPromptsJSON, "json", "j", false, "Output prompt list as structured JSON")
}

func runAgyListPrompts(args []string) error {
	limit := parsePromptLimit(args)
	if listPromptsProjectsCount > 0 {
		return launchVSCodeForLastProjects(listPromptsProjectsCount, limit)
	}

	return displayPromptListing(limit)
}

func parsePromptLimit(args []string) int {
	if len(args) == 0 {
		return 50
	}
	if strings.ToLower(args[0]) == "all" {
		return 10000
	}
	val, err := strconv.Atoi(args[0])
	if err != nil || val < 1 {
		return 50
	}

	return val
}
