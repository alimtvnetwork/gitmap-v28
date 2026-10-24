package cmdplans

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderStatusPill(status enums.PlanStatusType) string {
	if status.IsCompleted() {
		return termout.TokyoOKPill(constants.PillCompleted)
	}
	if status.IsInProgress() || status.IsActive() {
		return termout.TokyoTag(constants.TagActive, constants.TagInProgress)
	}

	return termout.TokyoWarnPill(constants.PillPending)
}

func planColumnsConfig() []termout.ColumnConfig {
	return []termout.ColumnConfig{
		{Title: constants.ColID, MinWidth: 5, Alignment: termout.AlignRight},
		{Title: constants.ColTitle, MinWidth: 26, MaxWidth: 52, WrapMode: termout.WrapEllipsis},
		{Title: constants.ColStatus, MinWidth: 14, Alignment: termout.AlignCenter},
		{Title: constants.ColFilePath, MinWidth: 24, WrapMode: termout.WrapClip},
	}
}

func buildPlanRow(item PlanItem) termout.Row {
	idStr := constants.DashID
	if item.ID > 0 {
		idStr = strconv.Itoa(item.ID)
	}

	return termout.Row{
		Cells: []string{
			idStr,
			item.Title,
			renderStatusPill(item.Status),
			item.FilePath,
		},
	}
}

func buildPlanRows(items []PlanItem) []termout.Row {
	rows := make([]termout.Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, buildPlanRow(item))
	}

	return rows
}

func buildTableConfig(items []PlanItem) termout.TableConfig {
	return termout.TableConfig{
		Columns:     planColumnsConfig(),
		Rows:        buildPlanRows(items),
		HasBorders:  true,
		BorderStyle: termout.BorderRounded,
		RowStripe:   true,
		HeaderColor: constants.ColorHeader,
		BorderColor: constants.ColorBorder,
		StripeColor: constants.ColorStripe,
	}
}

func printPlansTable(title string, items []PlanItem) {
	if len(items) == 0 {
		return
	}
	fmt.Printf("\n  %s\n", title)
	termout.PrintTable(buildTableConfig(items))
}

func printOverflowGuidance(summary PlansSummary) {
	if summary.HasMorePending {
		fmt.Printf("\n  %s %s\n",
			termout.TokyoTag(constants.TagNotice, constants.TagOverflow),
			fmt.Sprintf("Showing %d of %d pending plans. %s",
				summary.DisplayedLimit, summary.PendingCount, summary.PendingOverflow))
	}
	if summary.HasMoreCompleted {
		fmt.Printf("  %s %s\n",
			termout.TokyoTag(constants.TagInfo, constants.TagCompleted),
			fmt.Sprintf("Showing %d of %d completed plans. Use 'gitmap %s %s' or 'gitmap %s %d' to view all.",
				summary.DisplayedLimit, summary.CompletedCount, constants.CmdPlans, constants.FlagAll, constants.CmdPlans, summary.CompletedCount))
	}
}

func formatPlanMetrics(summary PlansSummary) string {
	return fmt.Sprintf("  %s   %s   %s",
		termout.TokyoMetricPill(constants.PillCompleted, summary.CompletedCount),
		termout.TokyoMetricPill(constants.PillPending, summary.PendingCount),
		termout.TokyoMetricPill(constants.PillTotal, summary.TotalCount),
	)
}

// RenderTTY prints an aesthetically formatted Tokyo Night overview of plans.
func RenderTTY(report PlansReport) {
	fmt.Printf("\n%s\n\n", termout.TokyoHeader(constants.HeaderGitMapPlans))
	fmt.Println(formatPlanMetrics(report.Summary))
	printPlansTable(constants.TitleCompletedPlans, report.Completed)
	printPlansTable(constants.TitlePendingPlans, report.Pending)
	printOverflowGuidance(report.Summary)
	fmt.Println()
}
