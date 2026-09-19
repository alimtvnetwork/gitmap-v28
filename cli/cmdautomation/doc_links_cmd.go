package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	docLinksCmd = &cobra.Command{
		Use:     "doc-links [dir]",
		Aliases: []string{"check-links", "doc-paths"},
		Short:   "Validate markdown relative links across documentation and memory",
		RunE:    runDocLinksCmd,
	}

	docLinksOpts DocLinksOptions
)

func runDocLinksCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		docLinksOpts.Dir = args[0]
	}
	monad := RunDocLinks(docLinksOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderDocLinksResult(res, docLinksOpts.IsJson)
	return nil
}

func renderDocLinksResult(res DocLinksResult, isJson bool) {
	if isJson {
		printDocLinksJson(res)
		return
	}
	printDocLinksTerminal(res)
}

func printDocLinksJson(res DocLinksResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	isSuccess := err == nil
	if isSuccess {
		fmt.Println(string(bytes))
	}
}

func printDocLinksTerminal(res DocLinksResult) {
	printDocLinksHeader(res)
	printDocLinksViolations(res.Violations)
	printDocLinksStatus(res)
}

func printDocLinksHeader(res DocLinksResult) {
	fmt.Printf("\n%s[Markdown Document Link Integrity Linter]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Scanned Files: %d\n", res.ScannedFiles)
	fmt.Printf("  Total Links:   %d\n", res.TotalLinks)
	fmt.Printf("  Broken Links:  %d\n", res.BrokenLinks)
	fmt.Printf("  Fixed Links:   %d\n", res.FixedLinks)
	fmt.Printf("  Duration:      %s\n", res.Duration)
}

func printDocLinksViolations(vios []DocLinkViolation) {
	hasVios := len(vios) > 0
	if !hasVios {
		return
	}
	fmt.Printf("\n%sUnresolved Relative Links:%s\n", constants.ColorYellow, constants.ColorReset)
	limit := 10
	if len(vios) < limit {
		limit = len(vios)
	}
	for i := 0; i < limit; i++ {
		v := vios[i]
		fmt.Printf("  - %s:%d: %s -> %s\n", v.File, v.LineNumber, v.RawLink, v.Issue)
	}
	hasMore := len(vios) > limit
	if hasMore {
		fmt.Printf("  ...and %d more broken link(s)\n", len(vios)-limit)
	}
}

func printDocLinksStatus(res DocLinksResult) {
	if res.IsClean {
		fmt.Printf("\n%s✅ SUCCESS: All markdown links resolved cleanly on disk.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s⚠️ WARNING: Detected %d broken markdown link(s).%s\n\n",
		constants.ColorYellow, res.BrokenLinks, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(docLinksCmd)
	initDocLinksFlags()
}

func initDocLinksFlags() {
	docLinksCmd.Flags().BoolVarP(&docLinksOpts.IsFix, "fix", "f", false, "Autofix known outdated path references")
	docLinksCmd.Flags().BoolVar(&docLinksOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
