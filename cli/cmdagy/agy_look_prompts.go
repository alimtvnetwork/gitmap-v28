package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var (
	isLookPromptsJSON  bool
	isLookPromptsAll   bool
	isLookPromptsSSH   bool
	lookPromptsCompact int
	lookPromptsCount   int
	lookPromptsFile    string
)

var AgyLookPromptsCmd = &cobra.Command{
	Use:   "look-prompts",
	Short: "Snapshot of running and pending prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyLookPrompts(args)
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

func init() {
	AgyLookPromptsCmd.Flags().BoolVarP(&isLookPromptsJSON, "json", "j", false, "JSON output")
	AgyLookPromptsCmd.Flags().BoolVarP(&isLookPromptsAll, "all", "a", false, "All projects")
	AgyLookPromptsCmd.Flags().BoolVar(&isLookPromptsSSH, "ssh", false, "Use SSH to aggregate across cluster nodes")
	AgyLookPromptsCmd.Flags().StringVarP(&lookPromptsFile, "file", "f", "", "Output to file (default prompts-snapshot.json)")
	if flag := AgyLookPromptsCmd.Flags().Lookup("file"); flag != nil {
		flag.NoOptDefVal = "prompts-snapshot.json"
	}
	AgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "compact", 0, "Compact prompt words")
	AgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "words", 0, "Words (alias for compact)")
	AgyLookPromptsCmd.Flags().IntVarP(&lookPromptsCount, "count", "c", 5, "Number of prompts to show")
	AgyCmd.AddCommand(AgyLookPromptsCmd)
}

func runAgyLookPrompts(args []string) *apperror.AppError {
	snapshot, err := CollectPromptsSnapshot(isLookPromptsAll, lookPromptsCompact, lookPromptsCount, isLookPromptsSSH)
	if err != nil {
		return apperror.WrapSimple(err, "collect prompts snapshot")
	}

	if isLookPromptsJSON || lookPromptsFile != "" {
		return handlePromptsJSONOutput(snapshot, lookPromptsFile)
	}

	RenderPromptsTable(snapshot)
	return nil
}

func handlePromptsJSONOutput(snapshot *AgyPromptSnapshot, targetFile string) *apperror.AppError {
	if outErr := OutputPromptsJSON(snapshot, targetFile); outErr != nil {
		return apperror.WrapSimple(outErr, "output prompts JSON")
	}
	return nil
}
