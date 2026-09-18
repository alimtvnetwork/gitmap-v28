package cmdagy

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
)

var (
	agyFixDetailed     bool
	agyFixNoRelease    bool
	agyFixCustomPrompt string
	agyFixNoClipboard  bool
	agyFixOutputFile   string
	agyFixDryRun       bool
	agyFixForce        bool
	agyFixAll          bool
	agyFixProjects     int
	agyFixLimit        int
	agyFixResetBatch   bool
	agyFixNoInject     bool
)

// agyFixPipelineCmd represents the agy fix-pipeline CLI command.
var agyFixPipelineCmd = &cobra.Command{
	Use:     "fix-pipeline [repo]",
	Aliases: []string{"fix", "fp", "pipeline-fix", "fixpipeline", "aef", "agy-errors-fix"},
	Short:   "Extract latest pipeline error logs and CI/CD fix prompt into clipboard and active temp prompt for Antigravity IDE",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyFixPipelineCLI(args)
	},
}

func init() {
	agyFixPipelineCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helptext.PrintWithMode("agy-fix-pipeline", render.PrettyAuto)
	})
	initAgyFixBehaviorFlags()
	initAgyFixBatchFlags()
}

func initAgyFixBehaviorFlags() {
	agyFixPipelineCmd.Flags().BoolVarP(&agyFixDetailed, "detailed", "v", false, "Include verbose passing lines in error logs")
	agyFixPipelineCmd.Flags().BoolVar(&agyFixNoRelease, "no-release", false, "Use CI/CD fix prompt without automated release")
	agyFixPipelineCmd.Flags().StringVarP(&agyFixCustomPrompt, "prompt", "p", "", "Path to custom prompt template")
	agyFixPipelineCmd.Flags().BoolVar(&agyFixNoClipboard, "no-clipboard", false, "Skip writing to system clipboard")
	agyFixPipelineCmd.Flags().StringVar(&agyFixOutputFile, "file", "", "Optional destination file path for prompt payload")
	agyFixPipelineCmd.Flags().BoolVarP(&agyFixDryRun, "dry-run", "d", false, "Preview payload statistics without saving or copying")
	agyFixPipelineCmd.Flags().BoolVarP(&agyFixForce, "force", "f", false, "Force resending even if previously sent")
}

func initAgyFixBatchFlags() {
	agyFixPipelineCmd.Flags().BoolVar(&agyFixAll, "all", false, "Scan all tracked repositories for failing pipelines")
	agyFixPipelineCmd.Flags().IntVar(&agyFixProjects, "projects", 0, "Number of projects to batch-fix")
	agyFixPipelineCmd.Flags().IntVar(&agyFixLimit, "limit", 3, "Limit number of projects to batch-fix")
	agyFixPipelineCmd.Flags().BoolVar(&agyFixResetBatch, "reset-batch", false, "Reset multi-project batch cursor to project 1")
	agyFixPipelineCmd.Flags().BoolVar(&agyFixNoInject, "no-inject", false, "Skip direct Antigravity CLI/IDE injection")
}

func resolveTargetRepoArg(args []string) string {
	opts := parseAgyFixArgs(args)

	return opts.Repo
}

// RunAgyFixPipelineCLI parses arguments and executes the pipeline fix feed assembly.
func RunAgyFixPipelineCLI(args []string) error {
	return RunPipelineFixAgyCLI(args)
}

// RunPipelineFixAgyCLI is the unified entrypoint for pipeline fix errors agy / aef.
func RunPipelineFixAgyCLI(args []string) error {
	if hasHelpFlag(args) {
		helptext.PrintWithMode("agy-fix-pipeline", render.PrettyAuto)

		return nil
	}

	opts := parseAgyFixArgs(args)
	if opts.IsAll || opts.ProjectsCount > 0 || opts.IsResetBatch {
		return RunAgyFixMultiProjectBatch(opts)
	}

	return executeSingleAgyFix(opts)
}
