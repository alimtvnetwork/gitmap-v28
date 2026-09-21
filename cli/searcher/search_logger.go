// Package searcher — search_logger.go records search queries into SearchSplitDB.
package searcher

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SearchLogEntry captures search execution metadata for persistence.
type SearchLogEntry struct {
	CategoryCode string
	QueryText    string
	RegexPattern string
	SearchType   string
	CallerIp     string
	IsAiCaller   bool
	DurationMs   int
	ResultCount  int
}

// LogSearchQuery records a search query and its metrics into SearchSplitDB.
func LogSearchQuery(entry SearchLogEntry) error {
	category := resolveCategoryCode(entry.CategoryCode, entry.IsAiCaller)
	return store.RecordSearchQuery(
		category,
		entry.QueryText,
		entry.RegexPattern,
		entry.SearchType,
		entry.CallerIp,
		entry.IsAiCaller,
		entry.DurationMs,
		entry.ResultCount,
	)
}

func resolveCategoryCode(category string, isAiCaller bool) string {
	if category != "" {
		return category
	}
	if isAiCaller {
		return "ai"
	}
	return "user"
}
