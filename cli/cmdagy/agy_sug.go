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

var (
	sugIntervalStr string
	sugDryRun      bool
)

// SUGWatchConfig models the persisted watch list configuration.
type SUGWatchConfig struct {
	ProjectTargets  []string `json:"projectTargets"`
	IntervalSec     int      `json:"intervalSeconds"`
	ClearOnShutdown *bool    `json:"clearOnShutdown,omitempty"`
	UpdatedAt       string   `json:"updatedAt"`
}

// IsClearOnShutdown reports whether projects list should clear when shutdown starts.
func (c *SUGWatchConfig) IsClearOnShutdown() bool {
	if c.ClearOnShutdown == nil {
		return true
	}
	return *c.ClearOnShutdown
}

// AgySUGCmd monitors designated projects and initiates OS shutdown once all pipelines pass.
var AgySUGCmd = &cobra.Command{
	Use:     "shutdown-until-green [command]",
	Aliases: []string{"sug", "shutdown-until"},
	Short:   "Monitor designated projects and trigger OS shutdown when all CI/CD pipelines turn green",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSUGCLI(args)
	},
}

func init() {
	AgySUGCmd.Flags().StringVarP(&sugIntervalStr, "time", "t", "5m", "Polling interval (minimum 2m, default 5m)")
	AgySUGCmd.Flags().BoolVarP(&sugDryRun, "dry-run", "n", false, "Simulate shutdown without power off")
	AgySUGCmd.Flags().BoolP("once", "1", false, "Run single verification pass and exit")
	AgySUGCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		RenderAgySugHelp()
	})
	AgySUGCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{
				"ls\tList monitored projects and watch loop status",
				"status\tCheck whether watch loop is running or idle",
				"run\tStart monitoring loop until green and then shut down",
				"watch\tAlias for run",
				"ui\tLaunch local dark-mode monitoring dashboard",
				"add-projects\tAdd project directories, repo names, or URLs to watch list",
				"rm\tRemove projects from watch list",
				"agy-running-projects\tAdd all currently running Antigravity projects",
				"help\tShow help and usage guide",
			}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	AgyCmd.AddCommand(AgySUGCmd)
}

// RunSUGCLI routes shutdown-until-green subcommands with robust normalization.
func RunSUGCLI(args []string) error {
	if len(args) == 0 || checkSUGHelp(args[0]) {
		printSUGHelp()
		return nil
	}
	subcmd, rest := normalizeSUGArgs(args)
	return routeSUGSubcommand(subcmd, rest)
}

func normalizeSUGArgs(args []string) (string, []string) {
	if len(args) == 0 {
		return "help", nil
	}
	first := strings.ToLower(args[0])
	if isRunningProjectsToken(first) && len(args) > 1 && strings.EqualFold(args[1], "projects") {
		return "agy-running-projects", args[2:]
	}
	if isRunningProjectsToken(first) {
		return "agy-running-projects", args[1:]
	}
	if first == "watch" && len(args) > 1 && strings.EqualFold(args[1], "ui") {
		return "ui", args[2:]
	}
	return first, args[1:]
}

func isRunningProjectsToken(token string) bool {
	switch token {
	case "agy-running-projects", "running-projects", "agy-running", "running", "arp", "rp":
		return true
	default:
		return false
	}
}

func routeSUGSubcommand(subcmd string, rest []string) error {
	switch subcmd {
	case "ls", "list", "status", "st":
		return listSUGProjects()
	case "add-projects", "add", "/add", "ap":
		if len(rest) > 0 && checkSUGHelp(rest[0]) {
			RenderAgySugHelp()
			return nil
		}
		return addSUGProjects(rest)
	case "rm", "/rm", "remove", "del", "delete":
		if len(rest) > 0 && checkSUGHelp(rest[0]) {
			RenderAgySugHelp()
			return nil
		}
		return removeSUGProjects(rest)
	case "agy-running-projects", "running-projects":
		return setSUGRunningProjects()
	case "run", "watch", "w":
		return runSUGLoop(rest)
	case "ui", "web":
		return RunSUGUI(rest)
	case "help", "--help", "-h":
		RenderAgySugHelp()
		return nil
	default:
		return handleDirectTargetOrFallback(subcmd, rest)
	}
}

func handleDirectTargetOrFallback(subcmd string, rest []string) error {
	resolved, isValid := ValidateSUGTarget(subcmd)
	if isValid {
		fmt.Printf("  ✔ Direct project target identified: %s\n", resolved)
		cfg := loadSUGConfig()
		appendNewSUGTargets(&cfg, []string{resolved})
		_ = saveSUGConfig(cfg)
		return ExecuteSUGWatch([]string{resolved}, parseSUGInterval(rest, isSUGDryRun(rest)), isSUGDryRun(rest), hasSUGOnceFlag(rest))
	}
	PrintSUGTargetValidationDiagnostic(subcmd)
	return apperror.NewValidationError("unknown sug subcommand or invalid target: " + subcmd)
}

func checkSUGHelp(token string) bool {
	t := strings.ToLower(token)
	return t == "help" || t == "--help" || t == "-h"
}

func printSUGHelp() {
	RenderAgySugHelp()
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
	st, isRunning := LoadSUGRuntimeStatus()

	fmt.Println()
	fmt.Printf("  %s%s SHUTDOWN-UNTIL-GREEN WATCH LIST (%d projects) %s%s\n",
		constants.ColorCyan, "╔════", len(cfg.ProjectTargets), "════╗", constants.ColorReset)

	if isRunning {
		fmt.Printf("  %s● Watch loop is ACTIVE%s (PID: %d, interval: %ds, started: %s)\n\n",
			constants.ColorGreen, constants.ColorReset, st.PID, st.IntervalSec, st.StartedAt)
	} else {
		fmt.Printf("  %s○ Watch loop is IDLE%s\n\n", constants.ColorYellow, constants.ColorReset)
	}

	if len(cfg.ProjectTargets) == 0 {
		fmt.Printf("  %sWatch list is empty. Use 'add-projects <target>' or 'agy-running-projects'.%s\n\n",
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
	var validTargets []string
	for _, t := range targets {
		resolved, isValid := ValidateSUGTarget(t)
		if isValid {
			validTargets = append(validTargets, resolved)
		} else {
			PrintSUGTargetValidationDiagnostic(t)
		}
	}
	if len(validTargets) == 0 {
		return apperror.NewValidationError("no valid targets provided; see diagnostic suggestions above")
	}

	cfg := loadSUGConfig()
	added := appendNewSUGTargets(&cfg, validTargets)
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
	if len(targets) == 0 {
		fmt.Println("  ℹ No active running AGY projects found with pending queues or prompts.")
		return nil
	}
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
	isDryRun := isSUGDryRun(args)
	isOnce := hasSUGOnceFlag(args)
	interval := parseSUGInterval(args, isDryRun)
	return ExecuteSUGWatch(cfg.ProjectTargets, interval, isDryRun, isOnce)
}

func parseSUGInterval(args []string, isDryRun bool) time.Duration {
	raw := sugIntervalStr
	for i, a := range args {
		if (a == "-t" || a == "--time") && i+1 < len(args) {
			raw = args[i+1]
		}
	}
	d, err := time.ParseDuration(raw)
	minInterval := 2 * time.Minute
	if isDryRun {
		minInterval = 10 * time.Second
	}
	if err != nil || d < minInterval {
		return minInterval
	}
	return d
}

func isSUGDryRun(args []string) bool {
	if sugDryRun {
		return true
	}
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "--dry-run" || low == "-n" || low == "--dry" {
			return true
		}
	}
	return false
}

func hasSUGOnceFlag(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "--once" || low == "-1" {
			return true
		}
	}
	return false
}
