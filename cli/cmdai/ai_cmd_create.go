package cmdai

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	createCmd = &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"new", "scaffold", "gen"},
		Short:   "Scaffold a new AI automation script skeleton",
		RunE:    runCreateCmd,
	}

	createOpts    CreateScriptOptions
	createTypeStr string
)

func runCreateCmd(cmd *cobra.Command, args []string) error {
	isEmpty := len(args) == 0
	if isEmpty {
		return apperror.NewValidationError("script name is required (e.g. 'my-linter')")
	}

	createOpts.Name = args[0]
	createOpts.Type = ScriptTemplateType(createTypeStr)

	return toError(RunAiCreate(createOpts))
}

func initCreateCmd() {
	AiCmd.AddCommand(createCmd)
	initCreateFlags()
}

func initCreateFlags() {
	createCmd.Flags().StringVarP(&createTypeStr, "type", "t", "linter", "Script archetype (linter, fixer, auditor, generator, util)")
	createCmd.Flags().StringVarP(&createOpts.Description, "desc", "d", "", "Short description of the script")
	createCmd.Flags().BoolVarP(&createOpts.IsParallel, "parallel", "p", false, "Include multi-worker parallel execution support")
	createCmd.Flags().BoolVar(&createOpts.HasFixMode, "fix", false, "Include --fix mode argument and application logic")
	createCmd.Flags().BoolVar(&createOpts.IsDryRun, "dry-run", false, "Preview generated script without writing to disk")
	createCmd.Flags().BoolVarP(&createOpts.IsForce, "force", "f", false, "Overwrite existing script file if already present")
}
