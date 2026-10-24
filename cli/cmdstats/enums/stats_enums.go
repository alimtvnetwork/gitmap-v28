package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// StatsOutputFormatType specifies presentation output formats for stats.
type StatsOutputFormatType string

const (
	StatsOutputFormatTable StatsOutputFormatType = "table"
	StatsOutputFormatJSON  StatsOutputFormatType = "json"
	StatsOutputFormatTUI   StatsOutputFormatType = "tui"
	StatsOutputFormatCSV   StatsOutputFormatType = "csv"
)

// IsValid reports whether the format matches a recognized output variant.
func (f StatsOutputFormatType) IsValid() bool {
	switch f {
	case StatsOutputFormatTable, StatsOutputFormatJSON,
		StatsOutputFormatTUI, StatsOutputFormatCSV:
		return true
	default:
		return false
	}
}

// String returns the string representation of the output format.
func (f StatsOutputFormatType) String() string {
	return string(f)
}

// IsTable reports whether the output format is tabular text.
func (f StatsOutputFormatType) IsTable() bool {
	return f == StatsOutputFormatTable
}

// IsJSON reports whether the output format is JSON.
func (f StatsOutputFormatType) IsJSON() bool {
	return f == StatsOutputFormatJSON
}

// IsTUI reports whether the output format is interactive TUI.
func (f StatsOutputFormatType) IsTUI() bool {
	return f == StatsOutputFormatTUI
}

// IsCSV reports whether the output format is comma-separated values.
func (f StatsOutputFormatType) IsCSV() bool {
	return f == StatsOutputFormatCSV
}

// HasFormatting reports whether the format has valid defined formatting.
func (f StatsOutputFormatType) HasFormatting() bool {
	return f.IsValid()
}

// ParseStatsOutputFormatType parses a raw string into a StatsOutputFormatType.
func ParseStatsOutputFormatType(s string) (StatsOutputFormatType, *apperror.AppError) {
	val := StatsOutputFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatsOutputFormatType", "E1000", fmt.Sprintf("invalid stats output format: %q", s)).WithContext("value", s)
}

// ParseStatsOutputFormat is an alias for ParseStatsOutputFormatType.
func ParseStatsOutputFormat(s string) (StatsOutputFormatType, *apperror.AppError) {
	return ParseStatsOutputFormatType(s)
}

// AllStatsOutputFormatTypes returns all supported output formats.
func AllStatsOutputFormatTypes() []StatsOutputFormatType {
	return []StatsOutputFormatType{
		StatsOutputFormatTable,
		StatsOutputFormatJSON,
		StatsOutputFormatTUI,
		StatsOutputFormatCSV,
	}
}

// AllStatsOutputFormats is an alias for AllStatsOutputFormatTypes.
func AllStatsOutputFormats() []StatsOutputFormatType {
	return AllStatsOutputFormatTypes()
}

// StatsMetricType specifies an aggregation metric for stats evaluation.
type StatsMetricType string

const (
	StatsMetricTotalRuns      StatsMetricType = "total_runs"
	StatsMetricSuccessCount   StatsMetricType = "success_count"
	StatsMetricFailCount      StatsMetricType = "fail_count"
	StatsMetricFailRate       StatsMetricType = "fail_rate"
	StatsMetricAvgDurationMs  StatsMetricType = "avg_duration_ms"
	StatsMetricMinDurationMs  StatsMetricType = "min_duration_ms"
	StatsMetricMaxDurationMs  StatsMetricType = "max_duration_ms"
	StatsMetricTotalCommands  StatsMetricType = "total_commands"
	StatsMetricUniqueCommands StatsMetricType = "unique_commands"
	StatsMetricLastUsed       StatsMetricType = "last_used"
)

// IsValid reports whether the metric is a recognized variant.
func (m StatsMetricType) IsValid() bool {
	switch m {
	case StatsMetricTotalRuns, StatsMetricSuccessCount, StatsMetricFailCount,
		StatsMetricFailRate, StatsMetricAvgDurationMs, StatsMetricMinDurationMs,
		StatsMetricMaxDurationMs, StatsMetricTotalCommands,
		StatsMetricUniqueCommands, StatsMetricLastUsed:
		return true
	default:
		return false
	}
}

// String returns the string representation of the metric.
func (m StatsMetricType) String() string {
	return string(m)
}

// IsTotalRuns reports whether the metric is total executions.
func (m StatsMetricType) IsTotalRuns() bool {
	return m == StatsMetricTotalRuns
}

// IsSuccessCount reports whether the metric is successful execution count.
func (m StatsMetricType) IsSuccessCount() bool {
	return m == StatsMetricSuccessCount
}

// IsFailCount reports whether the metric is failed execution count.
func (m StatsMetricType) IsFailCount() bool {
	return m == StatsMetricFailCount
}

// IsFailRate reports whether the metric is failure rate percentage.
func (m StatsMetricType) IsFailRate() bool {
	return m == StatsMetricFailRate
}

// IsAvgDurationMs reports whether the metric is average duration in ms.
func (m StatsMetricType) IsAvgDurationMs() bool {
	return m == StatsMetricAvgDurationMs
}

// IsMinDurationMs reports whether the metric is minimum duration in ms.
func (m StatsMetricType) IsMinDurationMs() bool {
	return m == StatsMetricMinDurationMs
}

// IsMaxDurationMs reports whether the metric is maximum duration in ms.
func (m StatsMetricType) IsMaxDurationMs() bool {
	return m == StatsMetricMaxDurationMs
}

// IsTotalCommands reports whether the metric is overall command executions.
func (m StatsMetricType) IsTotalCommands() bool {
	return m == StatsMetricTotalCommands
}

// IsUniqueCommands reports whether the metric is unique distinct command count.
func (m StatsMetricType) IsUniqueCommands() bool {
	return m == StatsMetricUniqueCommands
}

// IsLastUsed reports whether the metric is the last executed timestamp.
func (m StatsMetricType) IsLastUsed() bool {
	return m == StatsMetricLastUsed
}

// HasDuration reports whether the metric represents a duration value.
func (m StatsMetricType) HasDuration() bool {
	return m.IsAvgDurationMs() || m.IsMinDurationMs() || m.IsMaxDurationMs()
}

// HasCount reports whether the metric represents a numerical counter.
func (m StatsMetricType) HasCount() bool {
	return m.IsTotalRuns() || m.IsSuccessCount() || m.IsFailCount() ||
		m.IsTotalCommands() || m.IsUniqueCommands()
}

// ParseStatsMetricType parses a raw string into a StatsMetricType.
func ParseStatsMetricType(s string) (StatsMetricType, *apperror.AppError) {
	val := StatsMetricType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatsMetricType", "E1000", fmt.Sprintf("invalid stats metric: %q", s)).WithContext("value", s)
}

// ParseStatsMetric is an alias for ParseStatsMetricType.
func ParseStatsMetric(s string) (StatsMetricType, *apperror.AppError) {
	return ParseStatsMetricType(s)
}

// AllStatsMetricTypes returns all supported aggregation metrics.
func AllStatsMetricTypes() []StatsMetricType {
	return []StatsMetricType{
		StatsMetricTotalRuns,
		StatsMetricSuccessCount,
		StatsMetricFailCount,
		StatsMetricFailRate,
		StatsMetricAvgDurationMs,
		StatsMetricMinDurationMs,
		StatsMetricMaxDurationMs,
		StatsMetricTotalCommands,
		StatsMetricUniqueCommands,
		StatsMetricLastUsed,
	}
}

// AllStatsMetrics is an alias for AllStatsMetricTypes.
func AllStatsMetrics() []StatsMetricType {
	return AllStatsMetricTypes()
}

// StatsTimeRangeType specifies the time window for aggregating stats.
type StatsTimeRangeType string

const (
	StatsTimeRangeAll   StatsTimeRangeType = "all"
	StatsTimeRangeToday StatsTimeRangeType = "today"
	StatsTimeRange7D    StatsTimeRangeType = "7d"
	StatsTimeRange30D   StatsTimeRangeType = "30d"
	StatsTimeRange90D   StatsTimeRangeType = "90d"
)

// IsValid reports whether the time range is a recognized variant.
func (t StatsTimeRangeType) IsValid() bool {
	switch t {
	case StatsTimeRangeAll, StatsTimeRangeToday, StatsTimeRange7D,
		StatsTimeRange30D, StatsTimeRange90D:
		return true
	default:
		return false
	}
}

// String returns the string representation of the time range.
func (t StatsTimeRangeType) String() string {
	return string(t)
}

// IsAll reports whether the time range covers all recorded history.
func (t StatsTimeRangeType) IsAll() bool {
	return t == StatsTimeRangeAll
}

// IsToday reports whether the time range is limited to today.
func (t StatsTimeRangeType) IsToday() bool {
	return t == StatsTimeRangeToday
}

// Is7D reports whether the time range is the past 7 days.
func (t StatsTimeRangeType) Is7D() bool {
	return t == StatsTimeRange7D
}

// Is30D reports whether the time range is the past 30 days.
func (t StatsTimeRangeType) Is30D() bool {
	return t == StatsTimeRange30D
}

// Is90D reports whether the time range is the past 90 days.
func (t StatsTimeRangeType) Is90D() bool {
	return t == StatsTimeRange90D
}

// HasBoundedTime reports whether the time range restricts the query interval.
func (t StatsTimeRangeType) HasBoundedTime() bool {
	return t != StatsTimeRangeAll && t.IsValid()
}

// ParseStatsTimeRangeType parses a raw string into a StatsTimeRangeType.
func ParseStatsTimeRangeType(s string) (StatsTimeRangeType, *apperror.AppError) {
	val := StatsTimeRangeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatsTimeRangeType", "E1000", fmt.Sprintf("invalid stats time range: %q", s)).WithContext("value", s)
}

// ParseStatsTimeRange is an alias for ParseStatsTimeRangeType.
func ParseStatsTimeRange(s string) (StatsTimeRangeType, *apperror.AppError) {
	return ParseStatsTimeRangeType(s)
}

// AllStatsTimeRangeTypes returns all supported time ranges.
func AllStatsTimeRangeTypes() []StatsTimeRangeType {
	return []StatsTimeRangeType{
		StatsTimeRangeAll,
		StatsTimeRangeToday,
		StatsTimeRange7D,
		StatsTimeRange30D,
		StatsTimeRange90D,
	}
}

// AllStatsTimeRanges is an alias for AllStatsTimeRangeTypes.
func AllStatsTimeRanges() []StatsTimeRangeType {
	return AllStatsTimeRangeTypes()
}

// StatsCategoryType specifies functional groupings of gitmap commands.
type StatsCategoryType string

const (
	StatsCategoryAll         StatsCategoryType = "all"
	StatsCategoryRepository  StatsCategoryType = "repository"
	StatsCategoryRemote      StatsCategoryType = "remote"
	StatsCategorySSH         StatsCategoryType = "ssh"
	StatsCategoryAutomation  StatsCategoryType = "automation"
	StatsCategoryMaintenance StatsCategoryType = "maintenance"
)

// IsValid reports whether the category is a recognized variant.
func (c StatsCategoryType) IsValid() bool {
	switch c {
	case StatsCategoryAll, StatsCategoryRepository, StatsCategoryRemote,
		StatsCategorySSH, StatsCategoryAutomation, StatsCategoryMaintenance:
		return true
	default:
		return false
	}
}

// String returns the string representation of the category.
func (c StatsCategoryType) String() string {
	return string(c)
}

// IsAll reports whether all command categories are included.
func (c StatsCategoryType) IsAll() bool {
	return c == StatsCategoryAll
}

// IsRepository reports whether repository-specific commands are selected.
func (c StatsCategoryType) IsRepository() bool {
	return c == StatsCategoryRepository
}

// IsRemote reports whether remote synchronization commands are selected.
func (c StatsCategoryType) IsRemote() bool {
	return c == StatsCategoryRemote
}

// IsSSH reports whether SSH and fleet operations commands are selected.
func (c StatsCategoryType) IsSSH() bool {
	return c == StatsCategorySSH
}

// IsAutomation reports whether automated pipeline and macro commands are selected.
func (c StatsCategoryType) IsAutomation() bool {
	return c == StatsCategoryAutomation
}

// IsMaintenance reports whether diagnostic and cleanup commands are selected.
func (c StatsCategoryType) IsMaintenance() bool {
	return c == StatsCategoryMaintenance
}

// HasCategoryFilter reports whether a specific category constraint is applied.
func (c StatsCategoryType) HasCategoryFilter() bool {
	return c != StatsCategoryAll && c.IsValid()
}

// ParseStatsCategoryType parses a raw string into a StatsCategoryType.
func ParseStatsCategoryType(s string) (StatsCategoryType, *apperror.AppError) {
	val := StatsCategoryType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatsCategoryType", "E1000", fmt.Sprintf("invalid stats category: %q", s)).WithContext("value", s)
}

// ParseStatsCategory is an alias for ParseStatsCategoryType.
func ParseStatsCategory(s string) (StatsCategoryType, *apperror.AppError) {
	return ParseStatsCategoryType(s)
}

// AllStatsCategoryTypes returns all supported categories.
func AllStatsCategoryTypes() []StatsCategoryType {
	return []StatsCategoryType{
		StatsCategoryAll,
		StatsCategoryRepository,
		StatsCategoryRemote,
		StatsCategorySSH,
		StatsCategoryAutomation,
		StatsCategoryMaintenance,
	}
}

// AllStatsCategories is an alias for AllStatsCategoryTypes.
func AllStatsCategories() []StatsCategoryType {
	return AllStatsCategoryTypes()
}

// StatsChartModeType specifies graphical visualization modes for statistics.
type StatsChartModeType string

const (
	StatsChartModeNone      StatsChartModeType = "none"
	StatsChartModeBar       StatsChartModeType = "bar"
	StatsChartModeSparkline StatsChartModeType = "sparkline"
	StatsChartModeHistogram StatsChartModeType = "histogram"
)

// IsValid reports whether the chart mode is a recognized variant.
func (c StatsChartModeType) IsValid() bool {
	switch c {
	case StatsChartModeNone, StatsChartModeBar,
		StatsChartModeSparkline, StatsChartModeHistogram:
		return true
	default:
		return false
	}
}

// String returns the string representation of the chart mode.
func (c StatsChartModeType) String() string {
	return string(c)
}

// IsNone reports whether chart rendering is disabled.
func (c StatsChartModeType) IsNone() bool {
	return c == StatsChartModeNone
}

// IsBar reports whether bar chart rendering is enabled.
func (c StatsChartModeType) IsBar() bool {
	return c == StatsChartModeBar
}

// IsSparkline reports whether sparkline trend rendering is enabled.
func (c StatsChartModeType) IsSparkline() bool {
	return c == StatsChartModeSparkline
}

// IsHistogram reports whether histogram distribution rendering is enabled.
func (c StatsChartModeType) IsHistogram() bool {
	return c == StatsChartModeHistogram
}

// HasVisualChart reports whether a visual graph or chart is requested.
func (c StatsChartModeType) HasVisualChart() bool {
	return c != StatsChartModeNone && c.IsValid()
}

// ParseStatsChartModeType parses a raw string into a StatsChartModeType.
func ParseStatsChartModeType(s string) (StatsChartModeType, *apperror.AppError) {
	val := StatsChartModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatsChartModeType", "E1000", fmt.Sprintf("invalid stats chart mode: %q", s)).WithContext("value", s)
}

// ParseStatsChartMode is an alias for ParseStatsChartModeType.
func ParseStatsChartMode(s string) (StatsChartModeType, *apperror.AppError) {
	return ParseStatsChartModeType(s)
}

// AllStatsChartModeTypes returns all supported chart visualization modes.
func AllStatsChartModeTypes() []StatsChartModeType {
	return []StatsChartModeType{
		StatsChartModeNone,
		StatsChartModeBar,
		StatsChartModeSparkline,
		StatsChartModeHistogram,
	}
}

// AllStatsChartModes is an alias for AllStatsChartModeTypes.
func AllStatsChartModes() []StatsChartModeType {
	return AllStatsChartModeTypes()
}
