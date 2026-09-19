package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	namingCmd = &cobra.Command{
		Use:     "naming [dir]",
		Aliases: []string{"bool-naming", "conventions"},
		Short:   "Audit boolean comparison anti-patterns and affirmative naming conventions",
		RunE:    runNamingCmd,
	}

	namingOpts NamingOptions
)

func runNamingCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		namingOpts.Dir = args[0]
	}
	monad := RunNamingAudit(namingOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	res := monad.Value
	renderNamingResult(res, namingOpts.IsJson)
	return nil
}

func renderNamingResult(res NamingResult, isJson bool) {
	if isJson {
		printNamingJson(res)
		return
	}
	printNamingTerminal(res)
}

func printNamingJson(res NamingResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printNamingTerminal(res NamingResult) {
	fmt.Printf("\n%s[Boolean & Naming Convention Audit]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Scanned Files: %d\n", res.ScannedFiles)
	fmt.Printf("  Violations:    %d\n", len(res.Violations))
	fmt.Printf("  Duration:      %s\n", res.Duration)

	if len(res.Violations) > 0 {
		printNamingViolations(res.Violations)
		return
	}
	fmt.Printf("%s✅ PASS: All files adhere to affirmative booleans and naming conventions.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func printNamingViolations(vios []NamingViolation) {
	fmt.Printf("\n%s❌ Naming Violations:%s\n", constants.ColorRed, constants.ColorReset)
	limit := 10
	if len(vios) < limit {
		limit = len(vios)
	}
	for i := 0; i < limit; i++ {
		v := vios[i]
		fmt.Printf("  %s%s:%d%s [%s] %s\n", constants.ColorCyan, v.File, v.LineNumber,
			constants.ColorReset, v.Kind, v.LineContent)
	}
	if len(vios) > limit {
		fmt.Printf("  ...and %d more violation(s)\n", len(vios)-limit)
	}
	fmt.Println()
}

func initNamingFlags() {
	namingCmd.Flags().StringSliceVarP(&namingOpts.Extensions, "ext", "e", nil, "Filter by file extensions")
	namingCmd.Flags().BoolVar(&namingOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
