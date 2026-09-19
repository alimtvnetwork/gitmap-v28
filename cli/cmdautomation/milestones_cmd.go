package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	milestonesCmd = &cobra.Command{
		Use:     "milestones [milestone-id]",
		Aliases: []string{"milestone", "consolidate-milestones"},
		Short:   "Query GitHub milestone issues and format structured release notes",
		RunE:    runMilestonesCmd,
	}

	milestonesOpts MilestonesOptions
)

func runMilestonesCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		milestonesOpts.MilestoneId = args[0]
	}
	monad := RunMilestones(milestonesOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderMilestonesResult(res, milestonesOpts.IsJson)
	return nil
}

func renderMilestonesResult(res MilestonesResult, isJson bool) {
	if isJson {
		printMilestonesJson(res)
		return
	}
	printMilestonesTerminal(res)
}

func printMilestonesJson(res MilestonesResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	isSuccess := err == nil
	if isSuccess {
		fmt.Println(string(bytes))
	}
}

func printMilestonesTerminal(res MilestonesResult) {
	printMilestonesHeader(res)
	printMilestonesIssueSummary(res.Items)
	printMilestonesNotes(res.ReleaseNotes)
}

func printMilestonesHeader(res MilestonesResult) {
	fmt.Printf("\n%s[Milestone Consolidator & Release Notes]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Milestone:     %s\n", res.MilestoneId)
	fmt.Printf("  Total Items:   %d\n", res.TotalItems)
	fmt.Printf("  Closed Items:  %d\n", res.ClosedCount)
	fmt.Printf("  Open Items:    %d\n", res.OpenCount)
	fmt.Printf("  Duration:      %s\n", res.Duration)
}

func printMilestonesIssueSummary(items []MilestoneIssue) {
	fmt.Printf("\n%sConsolidated Milestone Items:%s\n", constants.ColorCyan, constants.ColorReset)
	limit := 10
	if len(items) < limit {
		limit = len(items)
	}
	for i := 0; i < limit; i++ {
		item := items[i]
		fmt.Printf("  [%s] #%d %s (%s)\n", item.Category, item.Number, item.Title, item.State)
	}
	hasMore := len(items) > limit
	if hasMore {
		fmt.Printf("  ...and %d more item(s)\n", len(items)-limit)
	}
}

func printMilestonesNotes(notes string) {
	hasNotes := notes != ""
	if hasNotes {
		fmt.Printf("\n%sGenerated Release Notes:%s\n", constants.ColorGreen, constants.ColorReset)
		fmt.Println(notes)
		return
	}
	fmt.Printf("\n%sNo closed items found to generate release notes.%s\n\n",
		constants.ColorYellow, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(milestonesCmd)
	initMilestonesFlags()
}

func initMilestonesFlags() {
	milestonesCmd.Flags().StringVarP(&milestonesOpts.OutputPath, "out", "o", "", "Output file path for generated release notes")
	milestonesCmd.Flags().BoolVar(&milestonesOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
