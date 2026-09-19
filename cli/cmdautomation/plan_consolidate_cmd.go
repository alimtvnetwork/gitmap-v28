package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	planConsolidateCmd = &cobra.Command{
		Use:     "plan-consolidate [dir]",
		Aliases: []string{"consolidate-plans", "consolidate"},
		Short:   "Cluster completed plans and subtasks into milestone summaries",
		RunE:    runPlanConsolidateCmd,
	}

	planConsolidateOpts PlanConsolidateOptions
)

func runPlanConsolidateCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		planConsolidateOpts.Dir = args[0]
	}
	monad := RunPlanConsolidate(planConsolidateOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderPlanConsolidateResult(res, planConsolidateOpts.IsJson, planConsolidateOpts.IsDryRun)
	return nil
}

func renderPlanConsolidateResult(res PlanConsolidateResult, isJson, isDryRun bool) {
	if isJson {
		printPlanConsolidateJson(res)
		return
	}
	printPlanConsolidateTerminal(res, isDryRun)
}

func printPlanConsolidateJson(res PlanConsolidateResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	isSuccess := err == nil
	if isSuccess {
		fmt.Println(string(bytes))
	}
}

func printPlanConsolidateTerminal(res PlanConsolidateResult, isDryRun bool) {
	printPlanConsolidateHeader(res, isDryRun)
	printPlanConsolidateClusters(res.Clusters)
	printPlanConsolidateStatus(res, isDryRun)
}

func printPlanConsolidateHeader(res PlanConsolidateResult, isDryRun bool) {
	title := "PLAN MEMORY CONSOLIDATION"
	if isDryRun {
		title = "PLAN MEMORY CONSOLIDATION (DRY RUN)"
	}
	fmt.Printf("\n%s[%s]%s\n", constants.ColorBold, title, constants.ColorReset)
	fmt.Printf("  Total Plans:      %d\n", res.TotalPlans)
	fmt.Printf("  Completed Plans:  %d\n", res.CompletedPlans)
	fmt.Printf("  Clusters Found:   %d\n", res.ClustersCount)
	fmt.Printf("  Subtasks Cleaned: %d\n", res.SubtasksCleaned)
	fmt.Printf("  Duration:         %s\n", res.Duration)
}

func printPlanConsolidateClusters(clusters []PlanCluster) {
	hasClusters := len(clusters) > 0
	if !hasClusters {
		return
	}
	fmt.Printf("\n%sConsolidated Milestone Clusters:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, c := range clusters {
		fmt.Printf("  • %s: %s (%d plans merged)\n", c.Id, c.Title, len(c.MergedPlans))
		for _, p := range c.MergedPlans {
			fmt.Printf("      - %s\n", p)
		}
	}
}

func printPlanConsolidateStatus(res PlanConsolidateResult, isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n%sℹ️ DRY RUN COMPLETE: Plans evaluated without disk modifications.%s\n\n",
			constants.ColorCyan, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s✅ SUCCESS: Plans index synchronized and milestone clusters consolidated.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(planConsolidateCmd)
	initPlanConsolidateFlags()
}

func initPlanConsolidateFlags() {
	planConsolidateCmd.Flags().IntVarP(&planConsolidateOpts.Threshold, "threshold", "t", 5, "Minimum completed plans threshold to cluster")
	planConsolidateCmd.Flags().BoolVarP(&planConsolidateOpts.IsDryRun, "dry-run", "d", false, "Preview consolidation without modifying files")
	planConsolidateCmd.Flags().BoolVar(&planConsolidateOpts.IsForce, "force", false, "Bypass confirmation prompts")
	planConsolidateCmd.Flags().BoolVar(&planConsolidateOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
