package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	formatGoCmd = &cobra.Command{
		Use:     "format-go [dir]",
		Aliases: []string{"gofmt-ast", "fmt-go"},
		Short:   "AST-aware Go code formatter organizing imports and enforcing UTF-8 LF",
		RunE:    runFormatGoCmd,
	}

	formatGoOpts FormatGoOptions
)

func runFormatGoCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		formatGoOpts.Dir = args[0]
	}
	monad := RunFormatGo(formatGoOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderFormatGoResult(monad.Value, formatGoOpts.IsJson)
	return nil
}

func renderFormatGoResult(res FormatGoResult, isJson bool) {
	if isJson {
		printFormatGoJson(res)
		return
	}
	printFormatGoTerminal(res)
}

func printFormatGoJson(res FormatGoResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printFormatGoTerminal(res FormatGoResult) {
	printFormatGoHeader(res)
	if res.TotalFiles == 0 {
		fmt.Printf("%s✅ No Go source files discovered to format.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	if res.IsClean {
		fmt.Printf("%s✅ All %d Go file(s) are cleanly formatted with UTF-8 LF.%s\n\n",
			constants.ColorGreen, res.TotalFiles, constants.ColorReset)
		return
	}
	printFormatGoItems(res.Items)
}

func printFormatGoHeader(res FormatGoResult) {
	fmt.Printf("\n%s[Go AST Code Formatter & Encoding Normalizer]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Total Files:     %d\n", res.TotalFiles)
	fmt.Printf("  Formatted/Fixed: %s%d%s\n", constants.ColorGreen, res.FormattedCount, constants.ColorReset)
	fmt.Printf("  Violations:      %s%d%s\n", constants.ColorYellow, res.ViolationCount, constants.ColorReset)
	fmt.Printf("  Duration:        %s\n\n", res.Duration)
}

func printFormatGoItems(items []FormatGoItem) {
	idx := 1
	for _, it := range items {
		canPrint := it.IsFormatted || len(it.ErrorMessage) > 0
		if canPrint {
			printSingleFormatItem(idx, it)
			idx++
		}
	}
	fmt.Println()
}

func printSingleFormatItem(idx int, it FormatGoItem) {
	status := "NEEDS-FIX"
	color := constants.ColorYellow
	hasErr := len(it.ErrorMessage) > 0
	if hasErr {
		status = "ERROR"
		color = constants.ColorRed
	}
	fmt.Printf("  %3d. %s[%-9s]%s %s %s\n",
		idx, color, status, constants.ColorReset, it.Path, it.ErrorMessage)
}

func init() {
	AutomationCmd.AddCommand(formatGoCmd)
	initFormatGoFlags()
}

func initFormatGoFlags() {
	formatGoCmd.Flags().StringVar(&formatGoOpts.Dir, "dir", ".", "Directory or file to format")
	formatGoCmd.Flags().BoolVarP(&formatGoOpts.IsWrite, "write", "w", false, "Write formatted changes back to disk")
	formatGoCmd.Flags().BoolVarP(&formatGoOpts.IsCheck, "check", "c", false, "Check and return violations if files need formatting")
	formatGoCmd.Flags().BoolVar(&formatGoOpts.IsStaged, "staged", false, "Format only staged Go files")
	formatGoCmd.Flags().BoolVar(&formatGoOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
