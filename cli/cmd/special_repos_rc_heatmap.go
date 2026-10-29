package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/heatmap"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func init() {
	hookRCHeatmapDispatch()
}

func hookRCHeatmapDispatch() {
	if len(os.Args) < 2 {
		return
	}
	arg1 := strings.ToLower(strings.TrimSpace(os.Args[1]))
	if isRCHeatmapSubcommand(arg1, os.Args) {
		executeAndExitRCHeatmap(os.Args[3:])
	}
	if arg1 == "rc-heatmap" || arg1 == "rc_heatmap" {
		executeAndExitRCHeatmap(os.Args[2:])
	}
}

func isRCHeatmapSubcommand(arg1 string, args []string) bool {
	if (arg1 != "rc" && arg1 != "repo-cache") || len(args) < 3 {
		return false
	}
	arg2 := strings.ToLower(strings.TrimSpace(args[2]))
	return arg2 == "heatmap" || arg2 == "rc-heatmap" || arg2 == "hm"
}

func executeAndExitRCHeatmap(args []string) {
	if err := RunRCHeatmapCommand(args); err != nil {
		cliexit.HandleError(err, 1)
	}
	cliexit.Exit(0)
}

type rcCommandTokens struct {
	subcommand string
	flags      []string
	hasSubcmd  bool
}

func isRCHeatmapHelpOrEmpty(args []string) bool {
	if len(args) == 0 {
		return true
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))

	return first == "help" || first == "--help" || first == "-h"
}

func parseRcCommandTokens(args []string) rcCommandTokens {
	var tokens rcCommandTokens
	for _, a := range args {
		tokens = consumeRcArg(tokens, a)
	}

	return tokens
}

func consumeRcArg(tokens rcCommandTokens, arg string) rcCommandTokens {
	if strings.HasPrefix(arg, "-") || tokens.hasSubcmd {
		tokens.flags = append(tokens.flags, arg)

		return tokens
	}

	tokens.subcommand = strings.ToLower(strings.TrimSpace(arg))
	tokens.hasSubcmd = true

	return tokens
}

// RunRCHeatmapCommand dispatches gitmap rc heatmap or fallback operations.
func RunRCHeatmapCommand(args []string) error {
	if isRCHeatmapHelpOrEmpty(args) {
		return heatmap.DispatchHeatmapOverview(context.Background(), os.Stdout, true)
	}

	tokens := parseRcCommandTokens(args)
	if !tokens.hasSubcmd || tokens.subcommand == "help" {
		return heatmap.DispatchHeatmapOverview(context.Background(), os.Stdout, true)
	}

	return dispatchWithSubcommandFirst(tokens)
}

func dispatchWithSubcommandFirst(tokens rcCommandTokens) error {
	enrichedFlags := enrichRCHeatmapFlags(tokens.flags)
	dispatched := append([]string{"rc", tokens.subcommand}, enrichedFlags...)
	err := heatmap.DispatchHeatmapSubsystem(context.Background(), dispatched, os.Stdout)
	if err == nil {
		return nil
	}

	return executeFallbackRCHeatmap(append([]string{tokens.subcommand}, enrichedFlags...))
}

func hasFlagKey(args []string, longKey, shortKey string) bool {
	for _, a := range args {
		trimmed := strings.TrimLeft(a, "-")
		key, _, _ := strings.Cut(trimmed, "=")
		norm := strings.ToLower(key)
		if norm == longKey || norm == shortKey {
			return true
		}
	}

	return false
}

func enrichRCHeatmapFlags(flags []string) []string {
	hasRepo := hasFlagKey(flags, "repo", "r")
	hasAi := hasFlagKey(flags, "ai", "a")
	enriched := append([]string{}, flags...)
	if !hasRepo {
		enriched = append(enriched, "--repo", AutoDetectRepoSlug())
	}

	if !hasAi {
		enriched = append(enriched, "--ai", AutoDetectAiName())
	}

	return enriched
}

func executeFallbackRCHeatmap(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printRCHeatmapHelp()

		return nil
	}

	switch strings.ToLower(args[0]) {
	case "sync":
		return runFallbackRCSync(args[1:])
	case "status":
		return runFallbackRCStatus(args[1:])
	case "commit":
		return runFallbackRCCommit(args[1:])
	case "tree", "live", "watch", "in-flight", "inflight":
		return heatmap.DispatchRcHeatmap(context.Background(), args, os.Stdout)
	}

	printRCHeatmapHelp()

	return nil
}

func runFallbackRCSync(args []string) error {
	repoSlug := extractFlagValue(args, "repo", AutoDetectRepoSlug())
	aiName := extractFlagValue(args, "ai", AutoDetectAiName())
	taskSlug := extractFlagValue(args, "task", "rc-heatmap-sync")
	fmt.Printf("[rc-heatmap] Synchronizing %s (%s)...\n", repoSlug, aiName)
	if err := heatmap.RunRcHeatmapSyncWithConflictResolution(aiName, repoSlug, taskSlug); err != nil {
		return err
	}
	fmt.Printf("[rc-heatmap] Successfully synchronized and committed %s.\n", repoSlug)
	return nil
}

func runFallbackRCStatus(args []string) error {
	repoSlug := extractFlagValue(args, "repo", AutoDetectRepoSlug())
	aiName := extractFlagValue(args, "ai", AutoDetectAiName())
	tracking, err := heatmap.GetLocalTrackingStatus(repoSlug, aiName)
	if err != nil {
		return err
	}
	fmt.Printf("Heatmap RC Status: %s (%s)\n", repoSlug, aiName)
	fmt.Printf("  Sync Status:        %s\n", tracking.SyncStatus)
	fmt.Printf("  Is Dirty:           %v\n", tracking.IsDirty)
	fmt.Printf("  Uncommitted Changes: %d\n", tracking.UncommittedChangeCount)
	fmt.Printf("  Last Committed SHA: %s\n", tracking.LastCommittedSha)
	fmt.Printf("  Last Commit Time:   %s\n", tracking.LastCommitTimestamp)
	return nil
}

func runFallbackRCCommit(args []string) error {
	repoCachePath, err := ResolveRepoCachePath()
	if err != nil {
		return err
	}
	repoSlug := extractFlagValue(args, "repo", AutoDetectRepoSlug())
	aiName := extractFlagValue(args, "ai", AutoDetectAiName())
	msg := extractFlagValue(args, "msg", "manual sync")
	newSha, err := heatmap.CommitRcHeatmap(repoCachePath, repoSlug, aiName, msg)
	if err != nil {
		return err
	}
	fmt.Printf("[rc-heatmap] Committed %s -> %s\n", repoSlug, newSha)
	return nil
}

func extractFlagValue(args []string, key, fallback string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "--"+key+"=") {
			return strings.TrimPrefix(a, "--"+key+"=")
		}
		if (a == "--"+key || a == "-"+key) && i+1 < len(args) {
			return args[i+1]
		}
	}
	return fallback
}

// ResolveRepoCachePath returns the physical filesystem path to the repo-cache directory.
func ResolveRepoCachePath() (string, error) {
	if p := resolveRepoCacheFromSpecialDB(); p != "" {
		return p, nil
	}
	candidates := []string{`../repo-cache`, `../../repo-cache`, `./repo-cache`}
	for _, c := range candidates {
		if info, statErr := os.Stat(c); statErr == nil && info.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs, nil
		}
	}
	return filepath.Abs("../repo-cache")
}

func resolveRepoCacheFromSpecialDB() string {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}
	defer db.Close()

	rec, getErr := db.GetSpecialRepo("rc")
	if getErr != nil || rec.LocalPath == "" {
		return ""
	}
	info, statErr := os.Stat(rec.LocalPath)
	if statErr != nil || !info.IsDir() {
		return ""
	}
	return rec.LocalPath
}

// AutoDetectRepoSlug infers the active repository slug from Git remote or current directory.
func AutoDetectRepoSlug() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	if out, err := cmd.Output(); err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return sanitizeRepoSlug(cleanRepoRemoteURL(strings.TrimSpace(string(out))))
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "gitmap"
	}
	return sanitizeRepoSlug(filepath.Base(cwd))
}

func cleanRepoRemoteURL(rawURL string) string {
	trimmed := strings.TrimSuffix(strings.TrimSpace(rawURL), ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return trimmed
}

func sanitizeRepoSlug(raw string) string {
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	var b strings.Builder
	for _, r := range cleaned {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else if r == '_' || r == ' ' || r == '.' {
			b.WriteRune('-')
		}
	}
	res := strings.Trim(b.String(), "-")
	if res == "" {
		return "gitmap"
	}
	return res
}

// AutoDetectAiName inspects active runtime environment variables to identify the AI agent.
func AutoDetectAiName() string {
	envMap := []struct {
		envKey string
		aiName string
	}{
		{"ANTIGRAVITY_AGENT_ID", "antigravity"},
		{"ANTIGRAVITY_INSTANCE_ID", "antigravity"},
		{"MUSE_RUNNER_ID", "muse"},
		{"MUSE_VERSION", "muse"},
		{"CLOUD_CODE_PROJECT", "cloud-code"},
		{"CODEX_ENV", "codex"},
		{"COPILOT_AGENT", "codex"},
		{"GROK_RUNNER_ID", "grok"},
		{"CURSOR_AGENT", "cursor"},
	}
	for _, item := range envMap {
		if val := os.Getenv(item.envKey); strings.TrimSpace(val) != "" {
			return item.aiName
		}
	}
	return "antigravity"
}

// CommitAndPushHeatmapRC coordinates feature commit and push on repo-cache.
func CommitAndPushHeatmapRC(repoCachePath, repoSlug, taskSummary string) (string, error) {
	aiName := AutoDetectAiName()
	return heatmap.CommitRcHeatmap(repoCachePath, repoSlug, aiName, taskSummary)
}

func printRCHeatmapHelp() {
	fmt.Println("Usage: gitmap rc heatmap <command> [flags]")
	fmt.Println("\nCommands:")
	fmt.Println("  status  Show local tracking database synchronization status")
	fmt.Println("  tree    Render hierarchical execution and reasoning tree")
	fmt.Println("  live    Show in-flight file touches and real-time telemetry")
	fmt.Println("  sync    Reconcile local deltas, rebase remote commits, and commit to repo-cache")
	fmt.Println("  commit  Create feature commit on repo-cache repository")
	fmt.Println("\nFlags:")
	fmt.Println("  --repo  Target repository slug (defaults to ambient repo)")
	fmt.Println("  --ai    AI agent namespace (defaults to active AI)")
	fmt.Println("  --task  Parent task slug")
	fmt.Println("  --file  Filter by touched file path")
	fmt.Println("  --limit Maximum records to return (default: 25)")
	fmt.Println("  --msg   Custom commit message")
}
