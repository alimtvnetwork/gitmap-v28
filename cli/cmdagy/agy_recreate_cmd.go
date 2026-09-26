// Package cmdagy — agy_recreate_cmd.go defines Cobra commands for recreate-project and recreate.
package cmdagy

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	recreatePromptFlag  string
	recreateModelFlag   string
	recreateProfileFlag string
	recreateDryRunFlag  bool
)

var agyRecreateProjectCmd = &cobra.Command{
	Use:     "recreate-project [target...]",
	Aliases: []string{"recreate", "rp", "rec"},
	Short:   "Purge project cache & conversations, re-register in AGY, and start new read-memory conversation",
	Long: `Purge cache, prompt queues, and conversations for target project(s),
re-register them freshly in Antigravity, and spawn a new conversation
instructing the agent to read all files, specs, and memory to understand the project.

Targeting:
  If no target is specified, the current git repository / directory is automatically used.
  Multiple targets can be specified as comma-separated or space-separated sequences,
  project IDs, aliases, or directory paths.

Examples:
  gitmap agy recreate-project
  gitmap agy recreate
  gitmap agy recreate D:\test-gitmap\test-gitmap
  gitmap agy recreate 1, 2, 3
  gitmap agy recreate gitmap, test-gitmap
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRecreate(args)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeAgyRmArgs(toComplete)
	},
}

func initAgyRecreateCmd() {
	bindRecreateFlags(agyRecreateProjectCmd)
	AgyCmd.AddCommand(agyRecreateProjectCmd)
}

func bindRecreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&recreatePromptFlag, "prompt", "p", "", "Custom prompt text (defaults to Read Memory protocol)")
	cmd.Flags().StringVarP(&recreateModelFlag, "model", "m", "", "Model type (flash_lite, flash, pro)")
	cmd.Flags().StringVar(&recreateProfileFlag, "profile", "", "Profile to use for session")
	cmd.Flags().BoolVarP(&recreateDryRunFlag, "dry-run", "d", false, "Preview recreate steps without modifying")
}

func runAgyRecreate(args []string) error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "get projects dir path")
	}

	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load agy projects")
	}
	sortAgyProjects(projects, "name")

	targets, resolveErr := ResolveAgyRecreateTargets(args, projects)
	if resolveErr != nil {
		return resolveErr
	}

	opts := AgyRecreateOptions{
		CustomPrompt: recreatePromptFlag,
		Model:        recreateModelFlag,
		Profile:      recreateProfileFlag,
		IsDryRun:     recreateDryRunFlag,
	}
	return ExecuteAgyRecreate(targets, opts)
}
