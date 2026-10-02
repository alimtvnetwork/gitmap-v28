package cmdignore

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ParseIgnoreInterval parses duration strings into time.Duration.
func ParseIgnoreInterval(val string) (time.Duration, error) {
	return config.ParseIgnoreInterval(val)
}

// RunIgnoreConfig handles inspection and updating of gitignore cache settings.
func RunIgnoreConfig(args []string) error {
	if len(args) > 0 && strings.EqualFold(args[0], "set") {
		return runIgnoreConfigSet(args[1:])
	}
	return printIgnoreConfigDetails()
}

func runIgnoreConfigSet(args []string) error {
	val := extractIntervalValue(args)
	if val == "" {
		return apperror.NewSimple("Usage: gitmap ignore config set interval <duration>", "E1030")
	}
	dur, err := ParseIgnoreInterval(val)
	if err != nil {
		fmt.Printf("  ✗ Error: %s\n", err.Error())
		return apperror.WrapSimple(err, "invalid interval")
	}
	return persistAndPrintInterval(val, dur)
}

func extractIntervalValue(args []string) string {
	if len(args) == 0 {
		return ""
	}
	if strings.ToLower(args[0]) == "interval" && len(args) > 1 {
		return args[1]
	}
	return args[0]
}

func persistAndPrintInterval(val string, dur time.Duration) error {
	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "open store database")
	}
	defer db.Close()
	if setErr := config.SetGitIgnoreCheckInterval(db, val); setErr != nil {
		return apperror.WrapSimple(setErr, "persist interval")
	}
	printConfigSetSuccess(val, dur)
	return nil
}

func printConfigSetSuccess(val string, dur time.Duration) {
	splitPath := store.ResolveSplitDbPath("gitignore", "cache", "")
	fmt.Printf("\n  %s✓ GitIgnore check interval successfully updated to: %s%s\n",
		constants.ColorGreen, val, constants.ColorReset)
	fmt.Printf("    Audit frequency : Repositories checked within %s will be bypassed during pull.\n", val)
	fmt.Printf("    Split-DB target : %s\n\n", splitPath)
}

func printIgnoreConfigDetails() error {
	intervalStr := resolveActiveInterval()
	splitPath := store.ResolveSplitDbPath("gitignore", "cache", "")
	cacheStats := loadCacheStats()
	printConfigHeader(intervalStr, splitPath, cacheStats)
	printConfigExamples()
	return nil
}

func resolveActiveInterval() string {
	db, err := store.OpenDefault()
	if err != nil {
		return config.DefaultGitIgnoreCheckInterval
	}
	defer db.Close()
	return config.GetGitIgnoreCheckInterval(db)
}

type ignoreCacheStats struct {
	TotalCount      int
	CleanCount      int
	RemediatedCount int
	LastAuditTime   string
}

func loadCacheStats() ignoreCacheStats {
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return ignoreCacheStats{LastAuditTime: "Never"}
	}
	defer db.Close()
	records, err := db.ListCachedChecks()
	if err != nil || len(records) == 0 {
		return ignoreCacheStats{LastAuditTime: "Never"}
	}
	return computeCacheStats(records)
}

func computeCacheStats(records []store.GitIgnoreRecord) ignoreCacheStats {
	stats := ignoreCacheStats{TotalCount: len(records), LastAuditTime: "Never"}
	latestTs := aggregateStatsRecords(records, &stats)
	if latestTs > 0 {
		stats.LastAuditTime = formatAuditTimestamp(latestTs)
	}
	return stats
}

func aggregateStatsRecords(records []store.GitIgnoreRecord, stats *ignoreCacheStats) int64 {
	var latestTs int64
	for _, r := range records {
		tallyStatus(r.Status, stats)
		if r.LastCheckedAt > latestTs {
			latestTs = r.LastCheckedAt
		}
	}
	return latestTs
}

func tallyStatus(status string, stats *ignoreCacheStats) {
	if status == "clean" {
		stats.CleanCount++
	} else if status == "remediated" {
		stats.RemediatedCount++
	}
}

func formatAuditTimestamp(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	ago := formatTimeAgo(time.Since(t))
	return fmt.Sprintf("%s (%s)", t.Format("2006-01-02 15:04:05 UTC"), ago)
}

func formatTimeAgo(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	return fmt.Sprintf("%dd ago", days)
}

func printConfigHeader(intervalStr, splitPath string, stats ignoreCacheStats) {
	fmt.Printf("\n%s  === GITMAP IGNORE CONFIGURATION ===%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Check Frequency  : %s\n", formatFrequencyLabel(intervalStr))
	fmt.Printf("  Split-DB Path    : %s\n", splitPath)
	fmt.Printf("  Cache Status     : Active (WAL Mode)\n")
	fmt.Printf("  Cached Repos     : %d total (%d clean, %d remediated)\n",
		stats.TotalCount, stats.CleanCount, stats.RemediatedCount)
	fmt.Printf("  Last Audit Time  : %s\n", stats.LastAuditTime)
}

func formatFrequencyLabel(interval string) string {
	dur, err := ParseIgnoreInterval(interval)
	if err != nil || dur == 0 {
		return interval
	}

	days := dur / (24 * time.Hour)
	hasExactDays := days >= 1 && dur%(24*time.Hour) == 0
	if !hasExactDays {
		return interval
	}

	if days == 1 {
		return fmt.Sprintf("%s (1 day)", interval)
	}

	return fmt.Sprintf("%s (%d days)", interval, days)
}

func printConfigExamples() {
	fmt.Printf("\n  Examples & Usage:\n")
	fmt.Printf("    gitmap ignore config set interval 12h   # Check once every 12 hours\n")
	fmt.Printf("    gitmap ignore config set interval 1d    # Check once a day (default)\n")
	fmt.Printf("    gitmap ignore config set interval 30m   # Check every 30 minutes\n")
	fmt.Printf("    gitmap ignore config set interval 0     # Disable cache (always inspect repos)\n")
	fmt.Printf("    gitmap ignore cache                     # View cached repo statuses\n")
	fmt.Printf("    gitmap ignore cache clear               # Invalidate all cached entries\n\n")
}

// RunIgnoreCache displays cached repository records or handles invalidation.
func RunIgnoreCache(args []string) error {
	if isCacheClearCommand(args) {
		return executeClearCache()
	}
	if isJSONOutput(args) {
		return renderCacheJSON()
	}
	return renderCacheTable()
}

func isCacheClearCommand(args []string) bool {
	for _, a := range args {
		lower := strings.ToLower(a)
		if lower == "clear" || lower == "--force" || lower == "-f" {
			return true
		}
	}
	return false
}

func isJSONOutput(args []string) bool {
	for _, a := range args {
		if strings.EqualFold(a, "--json") {
			return true
		}
	}
	return false
}

func executeClearCache() error {
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "open split db")
	}
	defer db.Close()
	return performClearCache(db)
}

func performClearCache(db *store.GitIgnoreSplitDB) error {
	records, _ := db.ListCachedChecks()
	count := len(records)
	if err := db.InvalidateAll(); err != nil {
		return apperror.WrapSimple(err, "clear cache")
	}
	printClearSuccess(count)
	return nil
}

func printClearSuccess(count int) {
	fmt.Printf("\n  %s✓ GitIgnore split-DB cache cleared (%d entries invalidated).%s\n",
		constants.ColorGreen, count, constants.ColorReset)
	fmt.Printf("    Next 'gitmap pull-all' or 'gitmap ignore action' will perform full inspection.\n\n")
}

func renderCacheJSON() error {
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "open split db")
	}
	defer db.Close()
	records, err := db.ListCachedChecks()
	if err != nil {
		return apperror.WrapSimple(err, "list cached checks")
	}
	data, _ := json.MarshalIndent(records, "", "  ")
	fmt.Println(string(data))
	return nil
}

func renderCacheTable() error {
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "open split db")
	}
	defer db.Close()
	return displayCacheRecords(db)
}

func displayCacheRecords(db *store.GitIgnoreSplitDB) error {
	records, err := db.ListCachedChecks()
	if err != nil {
		return apperror.WrapSimple(err, "list cached checks")
	}
	ttl := resolveActiveTTL()
	printCacheTableHeader(len(records), ttl)
	validCount, expiredCount := printCacheRows(records, ttl)
	fmt.Printf("\n  Summary: %d valid, %d expired.\n\n", validCount, expiredCount)
	return nil
}

func resolveActiveTTL() time.Duration {
	sdb, err := store.OpenDefault()
	if err != nil {
		return 24 * time.Hour
	}
	defer sdb.Close()
	return config.GetGitIgnoreTTL(sdb)
}

func printCacheTableHeader(total int, ttl time.Duration) {
	fmt.Printf("\n%s  === GITMAP IGNORE SPLIT-DB CACHE (%d entries) ===%s\n",
		constants.ColorCyan, total, constants.ColorReset)
	fmt.Printf("  TTL Setting: %s\n\n", ttl.String())
	fmt.Printf("  %-20s %-12s %-11s %-9s %-17s %s\n",
		"REPO SLUG", "STATUS", "REMEDIATED", "DURATION", "LAST CHECKED", "CACHE STATUS")
	fmt.Printf("  %s\n", strings.Repeat("-", 84))
}

func printCacheRows(records []store.GitIgnoreRecord, ttl time.Duration) (int, int) {
	var validCount, expiredCount int
	for _, r := range records {
		isValid := isRecordValid(r, ttl)
		if isValid {
			validCount++
		} else {
			expiredCount++
		}
		printSingleCacheRow(r, isValid)
	}
	return validCount, expiredCount
}

func isRecordValid(r store.GitIgnoreRecord, ttl time.Duration) bool {
	if !r.IsActive || ttl <= 0 {
		return false
	}
	elapsed := time.Since(time.Unix(r.LastCheckedAt, 0))
	return elapsed < ttl
}

func printSingleCacheRow(r store.GitIgnoreRecord, isValid bool) {
	ago := formatTimeAgo(time.Since(time.Unix(r.LastCheckedAt, 0)))
	statusTag := formatStatusTag(isValid)
	fmt.Printf("  %-20s %-12s %-11d %-9s %-17s %s\n",
		truncateString(r.RepoSlug, 20),
		r.Status,
		r.RemediatedCount,
		fmt.Sprintf("%dms", r.DurationMs),
		ago,
		statusTag,
	)
}

func formatStatusTag(isValid bool) string {
	if isValid {
		return fmt.Sprintf("%s[VALID]%s", constants.ColorGreen, constants.ColorReset)
	}
	return fmt.Sprintf("%s[EXPIRED]%s", constants.ColorYellow, constants.ColorReset)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ParsedFlags stores decomposed CLI flag arguments.
type ParsedFlags struct {
	Duration    time.Duration
	HasDuration bool
	Remaining   []string
}

// ParseIgnoreFlags extracts --interval, -i, and --force flags.
func ParseIgnoreFlags(args []string) (time.Duration, bool, []string, error) {
	pf := &ParsedFlags{}
	for i := 0; i < len(args); i++ {
		nextI, err := processFlagArg(args, i, pf)
		if err != nil {
			return 0, false, nil, err
		}
		i = nextI
	}
	return pf.Duration, pf.HasDuration, pf.Remaining, nil
}

func processFlagArg(args []string, i int, pf *ParsedFlags) (int, error) {
	arg := args[i]
	if arg == "--force" {
		pf.Duration = 0
		pf.HasDuration = true
		return i, nil
	}
	if arg == "--interval" || arg == "-i" {
		return parseIntervalFlag(args, i, pf)
	}
	pf.Remaining = append(pf.Remaining, arg)
	return i, nil
}

func parseIntervalFlag(args []string, i int, pf *ParsedFlags) (int, error) {
	if i+1 >= len(args) {
		return i, fmt.Errorf("flag %s requires duration value", args[i])
	}
	dur, err := ParseIgnoreInterval(args[i+1])
	if err != nil {
		return i, err
	}
	pf.Duration = dur
	pf.HasDuration = true
	return i + 1, nil
}

// ResolveEffectiveTTL returns the effective duration considering CLI overrides or DB setting.
func ResolveEffectiveTTL(args []string, db *store.DB) (time.Duration, []string, error) {
	dur, hasDuration, remaining, err := ParseIgnoreFlags(args)
	if err != nil {
		return 0, remaining, err
	}
	if hasDuration {
		return dur, remaining, nil
	}
	return config.GetGitIgnoreTTL(db), remaining, nil
}
