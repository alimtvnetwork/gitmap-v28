// Package cmdagy — agy_clean_cache_cmd.go defines CLI commands for cache clear and retention.
package cmdagy

import (
	"strconv"

	"github.com/spf13/cobra"
)

var (
	agyCleanDryRun      bool
	agyCleanPreflight   bool
	agyCleanYes         bool
	agyCleanForce       bool
	agyCleanNoKill      bool
	agyCleanJSON        bool
	agyCleanIncludeTemp bool
	agyCleanKeep        int
)

var agyCleanCacheCmd = &cobra.Command{
	Use:     "clean-cache [keep-count]",
	Aliases: []string{"cache-clear", "cleancache", "clean_cache", "cc"},
	Short:   "Clean Antigravity cache stores and prune conversations with retention",
	RunE: func(cmd *cobra.Command, args []string) error {
		keep := resolveKeepCount(args, agyCleanKeep)
		opts := buildCleanCacheOpts(keep, agyCleanDryRun, agyCleanPreflight)
		return ExecuteCleanCache(opts)
	},
}

var agyCacheClearKeepOneCmd = &cobra.Command{
	Use:     "cache-clear-keep-one",
	Aliases: []string{"ccko", "cache-clear-keep-1", "cc-keep-one"},
	Short:   "Clean Antigravity cache keeping only the single most recent conversation",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := buildCleanCacheOpts(1, agyCleanDryRun, agyCleanPreflight)
		return ExecuteCleanCache(opts)
	},
}

var agyCacheClearKeepFiveCmd = &cobra.Command{
	Use:     "cache-clear-keep-five",
	Aliases: []string{"cckf", "cache-clear-keep-5", "cc-keep-five"},
	Short:   "Clean Antigravity cache keeping the top 5 most recent conversations",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := buildCleanCacheOpts(5, agyCleanDryRun, agyCleanPreflight)
		return ExecuteCleanCache(opts)
	},
}

func init() {
	bindCleanCacheFlags(agyCleanCacheCmd)
	bindCleanCacheFlags(agyCacheClearKeepOneCmd)
	bindCleanCacheFlags(agyCacheClearKeepFiveCmd)
}

func bindCleanCacheFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&agyCleanDryRun, "dry-run", "d", false, "Preview cache targets and processes without deleting")
	cmd.Flags().BoolVar(&agyCleanPreflight, "preflight", false, "Simulate cleanup and retention analysis (alias: --pre, --precheck)")
	cmd.Flags().BoolVar(&agyCleanPreflight, "pre", false, "Simulate cleanup and retention analysis")
	cmd.Flags().BoolVar(&agyCleanPreflight, "precheck", false, "Simulate cleanup and retention analysis")
	cmd.Flags().BoolVarP(&agyCleanYes, "yes", "y", false, "Proceed with cleanup without interactive prompt")
	cmd.Flags().BoolVarP(&agyCleanForce, "force", "f", false, "Force terminate processes and proceed without confirmation")
	cmd.Flags().BoolVar(&agyCleanNoKill, "no-kill", false, "Skip terminating running processes before cleaning")
	cmd.Flags().BoolVar(&agyCleanJSON, "json", false, "Output results in JSON format")
	cmd.Flags().BoolVar(&agyCleanIncludeTemp, "include-temp", false, "Include system/user temp directory in cleanup")
	cmd.Flags().IntVarP(&agyCleanKeep, "keep", "k", 10, "Number of conversations to retain (default: 10)")
}

func resolveKeepCount(args []string, defaultKeep int) int {
	if len(args) == 0 {
		return defaultKeep
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val >= 0 {
		return val
	}
	return defaultKeep
}

func buildCleanCacheOpts(keep int, dryRun, preflight bool) CleanCacheOptions {
	return CleanCacheOptions{
		DryRun:      dryRun,
		Preflight:   preflight,
		Force:       agyCleanForce,
		Yes:         agyCleanYes || agyCleanForce,
		NoKill:      agyCleanNoKill,
		JSON:        agyCleanJSON,
		IncludeTemp: agyCleanIncludeTemp,
		Keep:        keep,
	}
}
