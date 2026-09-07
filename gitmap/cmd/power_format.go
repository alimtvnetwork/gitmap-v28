package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/power"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func formatTimeoutMinutes(m int) string {
	if m == 0 {
		return "Never"
	}

	return fmt.Sprintf("%d minutes", m)
}

func printPowerStatus(current, dbSetting power.Settings) {
	fmt.Printf("▶ OS Power & Sleep Status (%s)\n", current.Platform)
	fmt.Printf("  • Display Timeout: %s\n", formatTimeoutMinutes(current.DisplayTimeoutMinutes))
	fmt.Printf("  • Sleep Timeout:   %s\n", formatTimeoutMinutes(current.SleepTimeoutMinutes))
	if current.IsNeverSleep {
		fmt.Printf("  • Mode:            Never-Sleep (inhibited)\n")
	} else {
		fmt.Printf("  • Mode:            Standard timeouts\n")
	}
	if dbSetting.Source != "" {
		fmt.Printf("  • SQLite Profile:  %s\n", dbSetting.Source)
	}
}

func printPowerHistoryTable(records []store.PowerHistoryRecord) {
	if len(records) == 0 {
		fmt.Println("  (no power history records found)")

		return
	}

	fmt.Printf("%-6s %-12s %-10s %-10s %-20s %s\n", "ID", "ACTION", "DISPLAY", "SLEEP", "TIMESTAMP", "NOTES")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, r := range records {
		disp := formatTimeoutMinutes(r.DisplayTimeoutMinutes)
		sleep := formatTimeoutMinutes(r.SleepTimeoutMinutes)
		fmt.Printf("%-6d %-12s %-10s %-10s %-20s %s\n",
			r.ID, r.Action, disp, sleep, r.CreatedAt, r.Notes)
	}
}
