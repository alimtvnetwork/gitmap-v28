package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var sugIntervalStr string

// SUGWatchConfig models the persisted watch list configuration.
type SUGWatchConfig struct {
	ProjectTargets []string `json:"projectTargets"`
	IntervalSec    int      `json:"intervalSeconds"`
	UpdatedAt      string   `json:"updatedAt"`
}

// AgySUGCmd monitors registered projects and initiates OS shutdown once all pipelines pass.
var AgySUGCmd = &cobra.Command{
	Use:     "shutdown-until-green [command]",
	Aliases: []string{"sug"},
	Short:   "Monitor designated projects and trigger OS shutdown when all CI/CD pipelines turn green",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSUGCLI(args)
	},
}

func init() {
	AgySUGCmd.Flags().StringVarP(&sugIntervalStr, "time", "t", "5m", "Polling interval (minimum 2m, default 5m)")
	AgyCmd.AddCommand(AgySUGCmd)
}

// RunSUGCLI routes shutdown-until-green subcommands.
func RunSUGCLI(args []string) error {
	if len(args) == 0 || checkSUGHelp(args[0]) {
		printSUGHelp()
		return nil
	}
	return routeSUGSubcommand(strings.ToLower(args[0]), args[1:])
}

func routeSUGSubcommand(subcmd string, rest []string) error {
	switch subcmd {
	case "ls", "list":
		return listSUGProjects()
	case "add-projects", "add":
		return addSUGProjects(rest)
	case "rm", "remove", "del":
		return removeSUGProjects(rest)
	case "agy-running-projects", "running-projects":
		return setSUGRunningProjects()
	case "run":
		return runSUGLoop(rest)
	default:
		return apperror.NewValidationError("unknown sug subcommand: " + subcmd)
	}
}

func checkSUGHelp(token string) bool {
	t := strings.ToLower(token)
	return t == "help" || t == "--help" || t == "-h"
}

func printSUGHelp() {
	fmt.Println()
	fmt.Printf("  %sAntigravity Shutdown Until Green (sug)%s\n\n", constants.ColorCyan, constants.ColorReset)
	printSUGUsageHelp()
	printSUGSubcommandsHelp()
}

func printSUGUsageHelp() {
	fmt.Println("  Usage:")
	fmt.Println("    gitmap agy shutdown-until-green [command] [-t <duration>]")
	fmt.Println("    gitmap agy sug <ls|run|add-projects|rm|agy-running-projects> [-t 5m]")
	fmt.Println()
}

func printSUGSubcommandsHelp() {
	fmt.Println("  Subcommands:")
	fmt.Println("    ls                      List projects in shutdown-watch list")
	fmt.Println("    add-projects <targets>  Add project(s) to watch list")
	fmt.Println("    rm <targets>            Remove project(s) from watch list")
	fmt.Println("    agy-running-projects    Set watch list to all currently running AGY projects")
	fmt.Println("    run [-t 5m]             Start the watch loop and shut down when green")
	fmt.Println("    help                    Show this help synopsis")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    -t, --time              Check interval (default 5m, minimum 2m)")
	fmt.Println()
}

func resolveSUGConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".gitmap", "sug_watch_list.json")
	}
	dir := filepath.Join(home, ".gemini", "antigravity")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "sug_watch_list.json")
}

func loadSUGConfig() SUGWatchConfig {
	path := resolveSUGConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return SUGWatchConfig{ProjectTargets: []string{}, IntervalSec: 300}
	}
	var cfg SUGWatchConfig
	if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
		return SUGWatchConfig{ProjectTargets: []string{}, IntervalSec: 300}
	}
	return cfg
}

func saveSUGConfig(cfg SUGWatchConfig) error {
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal sug config")
	}
	path := resolveSUGConfigPath()
	if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write sug config")
	}
	return nil
}

func listSUGProjects() error {
	cfg := loadSUGConfig()
	fmt.Println()
	fmt.Printf("  %s%s SHUTDOWN-UNTIL-GREEN WATCH LIST (%d projects) %s%s\n",
		constants.ColorCyan, "╔════", len(cfg.ProjectTargets), "════╗", constants.ColorReset)
	if len(cfg.ProjectTargets) == 0 {
		fmt.Printf("  %sWatch list is empty. Use 'add-projects' or 'agy-running-projects'.%s\n\n",
			constants.ColorDim, constants.ColorReset)
		return nil
	}
	for i, p := range cfg.ProjectTargets {
		fmt.Printf("  [%d] %s\n", i+1, p)
	}
	fmt.Println()
	return nil
}

func addSUGProjects(targets []string) error {
	if len(targets) == 0 {
		return apperror.NewValidationError("add-projects requires at least one project target")
	}
	cfg := loadSUGConfig()
	added := appendNewSUGTargets(&cfg, targets)
	if err := saveSUGConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("  ✔ Added %d project(s) to shutdown watch list (total: %d)\n", added, len(cfg.ProjectTargets))
	return nil
}

func appendNewSUGTargets(cfg *SUGWatchConfig, targets []string) int {
	seen := buildSUGSeenMap(cfg.ProjectTargets)
	added := 0
	for _, t := range targets {
		clean := strings.Trim(t, ",")
		if clean != "" && !seen[strings.ToLower(clean)] {
			seen[strings.ToLower(clean)] = true
			cfg.ProjectTargets = append(cfg.ProjectTargets, clean)
			added++
		}
	}
	return added
}

func buildSUGSeenMap(projects []string) map[string]bool {
	seen := make(map[string]bool)
	for _, p := range projects {
		seen[strings.ToLower(p)] = true
	}
	return seen
}

func removeSUGProjects(targets []string) error {
	if len(targets) == 0 {
		return apperror.NewValidationError("rm requires at least one project target")
	}
	cfg := loadSUGConfig()
	cfg.ProjectTargets = filterKeptSUGProjects(cfg.ProjectTargets, targets)
	if err := saveSUGConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("  ✔ Updated shutdown watch list (remaining: %d project(s))\n", len(cfg.ProjectTargets))
	return nil
}

func filterKeptSUGProjects(existing, toRemove []string) []string {
	removeMap := buildSUGSeenMap(toRemove)
	var kept []string
	for _, p := range existing {
		if !removeMap[strings.ToLower(strings.Trim(p, ","))] {
			kept = append(kept, p)
		}
	}
	return kept
}

func setSUGRunningProjects() error {
	running, err := DiscoverRunningProjects()
	if err != nil {
		return err
	}
	targets := extractRunningProjectNames(running)
	cfg := SUGWatchConfig{ProjectTargets: targets, IntervalSec: 300}
	if saveErr := saveSUGConfig(cfg); saveErr != nil {
		return saveErr
	}
	fmt.Printf("  ✔ Registered %d running project(s) into shutdown watch list\n", len(targets))
	return nil
}

func extractRunningProjectNames(running []RunningProjectRecord) []string {
	var targets []string
	for _, r := range running {
		targets = append(targets, r.ProjectName)
	}
	return targets
}

func runSUGLoop(args []string) error {
	cfg := loadSUGConfig()
	if len(cfg.ProjectTargets) == 0 {
		return apperror.NewValidationError("shutdown watch list is empty; add projects before running")
	}
	interval := parseSUGInterval(args)
	return ExecuteSUGWatch(cfg.ProjectTargets, interval)
}

func parseSUGInterval(args []string) time.Duration {
	raw := sugIntervalStr
	for i, a := range args {
		if (a == "-t" || a == "--time") && i+1 < len(args) {
			raw = args[i+1]
		}
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d < 2*time.Minute {
		return 2 * time.Minute
	}
	return d
}
