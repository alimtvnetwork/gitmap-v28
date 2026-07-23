package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// formatTRRow renders a single temp-release list row.
func formatTRRow(r model.TempRelease) string {
	return fmt.Sprintf("%-28s %-12s %-5d %-10s %s",
		truncateStr(r.Branch, 28), r.VersionPrefix,
		r.SequenceNumber, shortSHA(r.CommitSha), r.CreatedAt)
}

// trPrefixGroup holds aggregated data for a version prefix.
type trPrefixGroup struct {
	prefix       string
	count        int
	seqRange     string
	latestCommit string
}

// groupTRByPrefix groups temp-release records by version prefix.
func groupTRByPrefix(records []model.TempRelease) []trPrefixGroup {
	prefixMap := make(map[string][]model.TempRelease)
	var order []string

	for _, r := range records {
		if _, exists := prefixMap[r.VersionPrefix]; !exists {
			order = append(order, r.VersionPrefix)
		}
		prefixMap[r.VersionPrefix] = append(prefixMap[r.VersionPrefix], r)
	}

	groups := make([]trPrefixGroup, 0, len(order))
	for _, prefix := range order {
		recs := prefixMap[prefix]
		sort.Slice(recs, func(i, j int) bool {
			return recs[i].SequenceNumber < recs[j].SequenceNumber
		})

		minSeq := recs[0].SequenceNumber
		maxSeq := recs[len(recs)-1].SequenceNumber
		seqRange := fmt.Sprintf("%d–%d", minSeq, maxSeq)
		if minSeq == maxSeq {
			seqRange = fmt.Sprintf("%d", minSeq)
		}

		groups = append(groups, trPrefixGroup{
			prefix:       prefix,
			count:        len(recs),
			seqRange:     seqRange,
			latestCommit: recs[len(recs)-1].CommitSha,
		})
	}

	return groups
}

// filterTRByPrefix filters records by version prefix substring.
func filterTRByPrefix(records []model.TempRelease, filter string) []model.TempRelease {
	if len(filter) == 0 {
		return records
	}

	var result []model.TempRelease
	lower := strings.ToLower(filter)

	for _, r := range records {
		if strings.Contains(strings.ToLower(r.VersionPrefix), lower) ||
			strings.Contains(strings.ToLower(r.Branch), lower) {
			result = append(result, r)
		}
	}

	return result
}
