package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	relPathsCmd = &cobra.Command{
		Use:     "relative-paths [dir]",
		Aliases: []string{"rel-paths", "paths", "links"},
		Short:   "Audit and sanitize forbidden absolute filesystem paths in documentation and code",
		RunE:    runRelPathsCmd,
	}

	relPathOpts RelPathOptions
)

func runRelPathsCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		relPathOpts.Dir = args[0]
	}
	monad := RunRelPathAudit(relPathOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	res := monad.Value
	renderRelPathsResult(res, relPathOpts.IsJson)
	return nil
}

func renderRelPathsResult(res RelPathResult, isJson bool) {
	if isJson {
		printRelPathsJson(res)
		return
	}
	printRelPathsTerminal(res)
}

func printRelPathsJson(res RelPathResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printRelPathsTerminal(res RelPathResult) {
	fmt.Printf("\n%s[Relative Path & Absolute URI Audit]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Scanned Files: %d\n", res.ScannedFiles)
	fmt.Printf("  Violations:    %d\n", len(res.Violations))
	fmt.Printf("  Fixed:         %d\n", res.FixedCount)
	fmt.Printf("  Duration:      %s\n", res.Duration)

	if len(res.Violations) > 0 {
		printPathViolations(res.Violations)
		return
	}
	fmt.Printf("%s✅ PASS: No forbidden absolute paths found.%s\n\n", constants.ColorGreen, constants.ColorReset)
}

func printPathViolations(vios []PathViolation) {
	fmt.Printf("\n%s❌ Absolute Path Violations:%s\n", constants.ColorRed, constants.ColorReset)
	limit := 10
	if len(vios) < limit {
		limit = len(vios)
	}
	for i := 0; i < limit; i++ {
		v := vios[i]
		fmt.Printf("  %s%s:%d%s %s\n", constants.ColorCyan, v.File, v.LineNumber, constants.ColorReset, v.RawPath)
	}
	if len(vios) > limit {
		fmt.Printf("  ...and %d more violation(s)\n", len(vios)-limit)
	}
	fmt.Println()
}

func initRelPathsFlags() {
	relPathsCmd.Flags().BoolVarP(&relPathOpts.IsFixMode, "fix", "f", false, "Auto-fix recognized path patterns")
	relPathsCmd.Flags().StringSliceVarP(&relPathOpts.Extensions, "ext", "e", nil, "Filter by file extensions (e.g. .md, .ts)")
	relPathsCmd.Flags().BoolVar(&relPathOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
