package cmdschedule

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runScheduleDebug prints deep diagnostics and split database info for a schedule.
func runScheduleDebug(args []string) error {
	if len(args) == 0 {
		return apperror.New("schedule", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify schedule name: gitmap schedule debug <name>",
		})
	}

	name := args[0]
	slug := store.ScheduleSlug(name)
	sdb, err := store.OpenScheduleSplitDB(slug)
	if err != nil {
		return err
	}
	defer sdb.Close()

	return printScheduleDiagnostics(name, slug, sdb)
}

func printScheduleDiagnostics(name, slug string, sdb *store.ScheduleSplitDB) error {
	cfg, _ := sdb.GetConfig()
	runs, err := sdb.GetRuns(10)
	if err != nil {
		return err
	}

	dbPath := store.ScheduleDBPath(slug)
	dbSize := getFileSizeStr(dbPath)

	printScheduleHeader(name, slug, dbPath, dbSize)
	printScheduleConfigDetails(cfg)
	printScheduleLogSummary(runs)

	return nil
}

func getFileSizeStr(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return "0 KB"
	}

	return fmt.Sprintf("%.2f KB", float64(fi.Size())/1024.0)
}

func printScheduleHeader(name, slug, path, size string) {
	fmt.Println()
	fmt.Printf("  %s Schedule Diagnostics: %s (%s)\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, name, slug)
	fmt.Printf("    Database Path : %s\n", path)
	fmt.Printf("    Database Size : %s\n", size)
}

func printScheduleConfigDetails(cfg *store.ScheduleConfig) {
	if cfg == nil {
		fmt.Println("    Config        : Not initialized")

		return
	}

	fmt.Printf("    Interval      : %s\n", cfg.IntervalVal)
	fmt.Printf("    Enabled       : %v\n", cfg.IsEnabled)
	fmt.Printf("    Startup Hook  : %v\n", cfg.IsStartup)
	if cfg.MacroName != "" {
		fmt.Printf("    Macro Target  : %s\n", cfg.MacroName)
	}
	if cfg.CommandLine != "" {
		fmt.Printf("    Command Line  : %s\n", cfg.CommandLine)
	}
}

func printScheduleLogSummary(runs []store.ScheduleRunRecord) {
	fmt.Printf("    Recent Runs   : %d\n", len(runs))
	if len(runs) == 0 {
		fmt.Println("    Execution Log : No recorded runs yet.")
		fmt.Println()

		return
	}

	last := runs[0]
	statusStr := "SUCCESS"
	if last.IsSuccess == false {
		statusStr = fmt.Sprintf("FAIL(%d)", last.ExitCode)
	}

	fmt.Printf("    Last Run At   : %s\n", last.StartedAt)
	fmt.Printf("    Last Duration : %d ms\n", last.DurationMS)
	fmt.Printf("    Last Status   : %s\n", statusStr)
	if last.Output != "" {
		fmt.Printf("    Last Output   : %s\n", truncateOutputExcerpt(last.Output))
	}
	fmt.Println()
}

func truncateOutputExcerpt(out string) string {
	excerpt := strings.TrimSpace(out)
	if len(excerpt) > 80 {
		return excerpt[:77] + "..."
	}

	return excerpt
}

// runScheduleEdit modifies configuration for an existing schedule.
func runScheduleEdit(args []string) error {
	if len(args) == 0 {
		return apperror.New("schedule", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify schedule name: gitmap schedule edit <name> [flags]",
		})
	}

	name := args[0]
	slug := store.ScheduleSlug(name)
	sdb, err := store.OpenScheduleSplitDB(slug)
	if err != nil {
		return err
	}
	defer sdb.Close()

	return applyScheduleEdit(sdb, name, slug, args[1:])
}

func applyScheduleEdit(sdb *store.ScheduleSplitDB, name, slug string, flags []string) error {
	cfg, err := sdb.GetConfig()
	if err != nil || cfg == nil {
		return apperror.New("schedule", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Schedule %q not found for editing.", name),
		})
	}

	updateScheduleConfigFromFlags(cfg, flags)
	if saveErr := sdb.SaveConfig(*cfg); saveErr != nil {
		return saveErr
	}

	fmt.Printf("%s Updated schedule %q (%s)\n",
		constants.ColorGreen+"✓"+constants.ColorReset, name, slug)

	return nil
}

func updateScheduleConfigFromFlags(cfg *store.ScheduleConfig, flags []string) {
	for i := 0; i < len(flags); i++ {
		f := flags[i]
		if strings.HasPrefix(f, "--interval=") {
			cfg.IntervalVal = strings.TrimPrefix(f, "--interval=")
		} else if strings.HasPrefix(f, "--delay=") {
			cfg.DelayVal = strings.TrimPrefix(f, "--delay=")
		} else if f == "--enable" {
			cfg.IsEnabled = true
		} else if f == "--disable" {
			cfg.IsEnabled = false
		}
	}
}
