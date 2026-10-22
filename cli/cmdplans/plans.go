package cmdplans

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const plansUsage = `Usage: gitmap plans [N] [flags]
Aliases: gitmap plan

Discovers and displays completed and pending architectural plans in .ai-memory/plans/.

Flags:
  -n, --limit <N>     Limit number of plans shown per section (default: 10)
  -a, --all           Show all completed and pending plans without truncation
  -p, --pending       Show only pending / active plans
  -c, --completed     Show only completed plans
  -j, --json          Emit machine-readable JSON to stdout
      --yaml, --yml   Emit machine-readable YAML to stdout
  -h, --help          Print this usage help

Examples:
  gitmap plans
  gitmap plans 20
  gitmap plans --all
  gitmap plans --json
  gitmap plans --yaml
`

func parseFormatFlag(arg string, opts *PlansOptions) bool {
	if arg == "-j" || arg == "--json" {
		opts.IsJSON = true
		return true
	}
	if arg == "--yaml" || arg == "--yml" {
		opts.IsYAML = true
		return true
	}

	return false
}

func parseFilterFlag(arg string, opts *PlansOptions) bool {
	if arg == "-a" || arg == "--all" {
		opts.IsAll = true
		return true
	}
	if arg == "-p" || arg == "--pending" {
		opts.IsPendingOnly = true
		return true
	}
	if arg == "-c" || arg == "--completed" {
		opts.IsCompletedOnly = true
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
	if (arg == "-n" || arg == "--limit") && idx+1 < len(args) {
		if n, err := strconv.Atoi(args[idx+1]); err == nil && n > 0 {
			opts.Limit = n
			return idx + 1, true
		}
	}

	return idx, false
}

func parsePlansOptions(args []string) (PlansOptions, bool) {
	opts := PlansOptions{Limit: 10, RepoPath: "."}
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Print(plansUsage)
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
	overflowCmd := fmt.Sprintf("Use 'gitmap plans %d' or 'gitmap plans --all' to view all.", pendTotal)

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
	if opts.IsJSON {
		return RenderJSON(report)
	}
	if opts.IsYAML {
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
		return apperror.Wrap(err, "DiscoverPlans")
	}

	return executePlansReport(opts, compTotal, pendTotal)
}
