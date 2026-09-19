package cmdautomation

import (
	"github.com/spf13/cobra"
)

var (
	sequenceCmd = &cobra.Command{
		Use:     "sequence [dir]",
		Aliases: []string{"seq", "titles", "seq-audit"},
		Short:   "Audit markdown sequence numbering gaps and H1 header alignment",
		RunE:    runSequenceCmd,
	}

	seqOpts SequenceAuditorOptions
)

func runSequenceCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		seqOpts.Dir = args[0]
	}
	res, err := RunSequenceAuditor(seqOpts)
	if err != nil {
		return err
	}
	if seqOpts.AsJson {
		return renderJSON(res)
	}
	renderSequenceAuditorResult(res)
	return nil
}

func initSequenceFlags() {
	sequenceCmd.Flags().BoolVarP(&seqOpts.IsFixMode, "fix", "f", false, "Automatically fix mismatched H1 headers to file prefix numbers")
	sequenceCmd.Flags().BoolVar(&seqOpts.AsJson, "json", false, "Output results as machine-readable JSON")
}
