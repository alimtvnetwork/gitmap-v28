package cmdplans

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"
)

func parseFormatFlag(arg string, opts *PlansOptions) bool {
	if slices.Contains([]string{constants.FlagJSONShort, constants.FlagJSON}, arg) {
		opts.IsJSON = true
		opts.Format = enums.PlansOutputFormatJSON
		return true
	}
	if slices.Contains([]string{constants.FlagYAML, constants.FlagYML}, arg) {
		opts.IsYAML = true
		opts.Format = enums.PlansOutputFormatYAML
		return true
	}

	return false
}

func parseFilterFlag(arg string, opts *PlansOptions) bool {
	if slices.Contains([]string{constants.FlagAllShort, constants.FlagAll}, arg) {
		opts.IsAll = true
		opts.FilterMode = enums.PlansFilterModeAll
		return true
	}
	if slices.Contains([]string{constants.FlagPendingShort, constants.FlagPending}, arg) {
		opts.IsPendingOnly = true
		opts.FilterMode = enums.PlansFilterModePending
		return true
	}
	if slices.Contains([]string{constants.FlagCompletedShort, constants.FlagCompleted}, arg) {
		opts.IsCompletedOnly = true
		opts.FilterMode = enums.PlansFilterModeCompleted
		return true
	}

	return false
}

func parsePlanArg(arg string, opts *PlansOptions) bool {
	if parseFormatFlag(arg, opts) || parseFilterFlag(arg, opts) {
		return true
	}

	return parsePositionalLimit(arg, opts)
}

func parsePositionalLimit(arg string, opts *PlansOptions) bool {
	if n, err := strconv.Atoi(arg); err == nil && n > 0 {
		opts.Limit = n
		return true
	}

	return false
}

func parseLimitOption(args []string, idx int, opts *PlansOptions) (int, bool) {
	arg := args[idx]
	if !slices.Contains([]string{constants.FlagLimitShort, constants.FlagLimit}, arg) || idx+1 >= len(args) {
		return idx, false
	}

	n, err := strconv.Atoi(args[idx+1])
	if err != nil || n <= 0 {
		return idx, false
	}

	opts.Limit = n

	return idx + 1, true
}

func parsePlansOptions(args []string) (PlansOptions, bool) {
	opts := PlansOptions{
		Limit:      constants.DefaultLimit,
		RepoPath:   constants.DefaultRepoPath,
		Format:     enums.PlansOutputFormatTable,
		FilterMode: enums.PlansFilterModeAll,
	}
	for i := 0; i < len(args); i++ {
		if slices.Contains([]string{constants.FlagHelpShort, constants.FlagHelp}, args[i]) {
			fmt.Print(constants.UsagePlans)
			return opts, false
		}
		if nextIdx, hasLimit := parseLimitOption(args, i, &opts); hasLimit {
			i = nextIdx
			continue
		}
		parsePlanArg(args[i], &opts)
	}

	return opts, true
}

func slicePlans(items []PlanItem, limit int, isAll bool) []PlanItem {
	if isAll || limit <= 0 || len(items) <= limit {
		return items
	}

	return items[:limit]
}

func buildSummary(compTotal, pendTotal, compShown, pendShown, limit int, isAll bool) PlansSummary {
	hasMorePending := !isAll && pendTotal > pendShown
	hasMoreCompleted := !isAll && compTotal > compShown
	overflowCmd := fmt.Sprintf("Use 'gitmap %s %d' or 'gitmap %s %s' to view all.",
		constants.CmdPlans, pendTotal, constants.CmdPlans, constants.FlagAll)

	return PlansSummary{
		CompletedCount:   compTotal,
		PendingCount:     pendTotal,
		TotalCount:       compTotal + pendTotal,
		DisplayedLimit:   limit,
		HasMorePending:   hasMorePending,
		HasMoreCompleted: hasMoreCompleted,
		PendingOverflow:  overflowCmd,
	}
}

func filterPlanReport(opts PlansOptions, comp, pend []PlanItem) ([]PlanItem, []PlanItem) {
	activeComp := comp
	activePend := pend
	if opts.IsPendingOnly {
		activeComp = nil
	}
	if opts.IsCompletedOnly {
		activePend = nil
	}

	return slicePlans(activeComp, opts.Limit, opts.IsAll), slicePlans(activePend, opts.Limit, opts.IsAll)
}

func dispatchReportOutput(opts PlansOptions, report PlansReport) error {
	if opts.Format.IsJSON() || opts.IsJSON {
		return RenderJSON(report)
	}
	if opts.Format.IsYAML() || opts.IsYAML {
		return RenderYAML(report)
	}
	RenderTTY(report)

	return nil
}

func executePlansReport(opts PlansOptions, compTotal, pendTotal []PlanItem) error {
	compShown, pendShown := filterPlanReport(opts, compTotal, pendTotal)
	summary := buildSummary(len(compTotal), len(pendTotal), len(compShown), len(pendShown), opts.Limit, opts.IsAll)
	report := PlansReport{
		Summary:   summary,
		Completed: compShown,
		Pending:   pendShown,
	}

	return dispatchReportOutput(opts, report)
}

// RunPlans executes the `gitmap plans` command workflow.
func RunPlans(args []string) error {
	opts, isProceed := parsePlansOptions(args)
	if !isProceed {
		return nil
	}

	compTotal, pendTotal, err := DiscoverPlans(opts.RepoPath)
	if err != nil {
		return apperror.WrapSimple(err, "DiscoverPlans")
	}

	return executePlansReport(opts, compTotal, pendTotal)
}
