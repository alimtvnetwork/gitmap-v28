package cmdplans

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderStatusPill(status PlanStatusType) string {
	if status == StatusCompleted {
		return termout.TokyoOKPill(string(StatusCompleted))
	}
	if status == StatusInProgress {
		return termout.TokyoTag("ACTIVE", "IN PROGRESS")
	}

	return termout.TokyoWarnPill(string(StatusPending))
}

func planColumnsConfig() []termout.ColumnConfig {
	return []termout.ColumnConfig{
		{Title: "ID", MinWidth: 5, Alignment: termout.AlignRight},
		{Title: "Title", MinWidth: 26, MaxWidth: 52, WrapMode: termout.WrapEllipsis},
		{Title: "Status", MinWidth: 14, Alignment: termout.AlignCenter},
		{Title: "File Path", MinWidth: 24, WrapMode: termout.WrapClip},
	}
}

func buildPlanRow(item PlanItem) termout.Row {
	idStr := "-"
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

func printPlansTable(title string, items []PlanItem) {
	if len(items) == 0 {
		return
	}
	fmt.Printf("\n  %s\n", title)
	cfg := termout.TableConfig{
		Columns:     planColumnsConfig(),
		Rows:        buildPlanRows(items),
		HasBorders:  true,
		BorderStyle: termout.BorderRounded,
		RowStripe:   true,
		HeaderColor: "#7aa2f7",
		BorderColor: "#3b4261",
		StripeColor: "#1f2335",
	}
	termout.PrintTable(cfg)
}

func printOverflowGuidance(summary PlansSummary) {
	if summary.HasMorePending {
		fmt.Printf("\n  %s %s\n",
			termout.TokyoTag("NOTICE", "OVERFLOW"),
			fmt.Sprintf("Showing %d of %d pending plans. %s",
				summary.DisplayedLimit, summary.PendingCount, summary.PendingOverflow))
	}
	if summary.HasMoreCompleted {
		fmt.Printf("  %s %s\n",
			termout.TokyoTag("INFO", "COMPLETED"),
			fmt.Sprintf("Showing %d of %d completed plans. Use 'gitmap plans --all' or 'gitmap plans %d' to view all.",
				summary.DisplayedLimit, summary.CompletedCount, summary.CompletedCount))
	}
}

// RenderTTY prints an aesthetically formatted Tokyo Night overview of plans.
func RenderTTY(report PlansReport) {
	fmt.Println()
	fmt.Println(termout.TokyoHeader("GitMap Plans — Repository Plan & Task Inventory"))
	fmt.Println()

	metrics := fmt.Sprintf("  %s   %s   %s",
		termout.TokyoMetricPill("Completed", report.Summary.CompletedCount),
		termout.TokyoMetricPill("Pending", report.Summary.PendingCount),
		termout.TokyoMetricPill("Total", report.Summary.TotalCount),
	)
	fmt.Println(metrics)

	printPlansTable("Completed Plans", report.Completed)
	printPlansTable("Pending / Queued Plans", report.Pending)
	printOverflowGuidance(report.Summary)
	fmt.Println()
}
