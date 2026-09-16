package cmdagy

import (
	"github.com/spf13/cobra"
)

var (
	agyRprpExcept string
	agyRprpPrompt string
	agyRprpDryRun bool
	agyRprpYes    bool
)

var agyReadAllProjectsWithReadPromptsCmd = &cobra.Command{
	Use:     "read-all-projects-with-read-prompts",
	Aliases: []string{"rprp", "rapwrp", "read-all-with-prompts"},
	Short:   "Discover local repos, register missing into Antigravity, and broadcast Read Memory prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyReadAllProjectsWithReadPrompts()
	},
}

func init() {
	agyReadAllProjectsWithReadPromptsCmd.Flags().StringVarP(&agyRprpExcept, "except", "e", "", "Exclude projects matching id, name, or prefix")
	agyReadAllProjectsWithReadPromptsCmd.Flags().StringVarP(&agyRprpPrompt, "prompt", "p", defaultReadMemoryPrompt, "Prompt text to send into project sessions")
	agyReadAllProjectsWithReadPromptsCmd.Flags().BoolVarP(&agyRprpDryRun, "dry-run", "d", false, "Preview without broadcasting")
	agyReadAllProjectsWithReadPromptsCmd.Flags().BoolVarP(&agyRprpYes, "yes", "y", false, "Send prompt without confirmation")
}

func runAgyReadAllProjectsWithReadPrompts() error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return err
	}
	discovered := discoverLocalRepos()
	syncMissingAgyProjects(discovered, dirPath)

	return broadcastReadMemoryToAll(dirPath)
}
