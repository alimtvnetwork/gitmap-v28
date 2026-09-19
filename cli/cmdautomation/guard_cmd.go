package cmdautomation

import (
	"github.com/spf13/cobra"
)

var (
	guardCmd = &cobra.Command{
		Use:     "guard [dir]",
		Aliases: []string{"file-guard", "size-guard", "blob-guard"},
		Short:   "Audit repository file sizes, common binaries, and large JSONs",
		RunE:    runGuardCmd,
	}

	guardOpts GuardOptions
)

func runGuardCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		guardOpts.Dir = args[0]
	}
	res, err := RunFileSizeGuard(guardOpts)
	if err != nil {
		return err
	}
	if guardOpts.AsJson {
		return renderJSON(res)
	}
	renderGuardResult(res)
	return nil
}

func initGuardFlags() {
	guardCmd.Flags().IntVarP(&guardOpts.MaxFileKb, "max-kb", "k", 500, "Maximum allowed file size in KB")
	guardCmd.Flags().IntVar(&guardOpts.MaxJsonKb, "max-json-kb", 500, "Maximum allowed JSON size before auto-exclusion")
	guardCmd.Flags().BoolVarP(&guardOpts.Interactive, "interactive", "i", true, "Prompt interactively when binaries are found")
	guardCmd.Flags().BoolVar(&guardOpts.AutoExclude, "auto-exclude", false, "Automatically exclude detected binaries without prompting")
	guardCmd.Flags().BoolVar(&guardOpts.AsJson, "json", false, "Output results as machine-readable JSON")
}
