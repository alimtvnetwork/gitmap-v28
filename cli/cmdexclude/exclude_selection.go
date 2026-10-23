package cmdexclude

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ParseNumberedSelection parses numeric sequences and ranges into sorted 1-based indices.
func ParseNumberedSelection(input string, maxCount int) ([]int, error) {
	clean := strings.TrimSpace(input)
	if clean == "" {
		return nil, nil
	}

	if strings.EqualFold(clean, "all") {
		return buildAllIndices(maxCount), nil
	}

	if strings.EqualFold(clean, "none") {
		return []int{}, nil
	}

	return parseTokenizedSelection(clean, maxCount)
}

func buildAllIndices(maxCount int) []int {
	results := make([]int, maxCount)
	for i := 0; i < maxCount; i++ {
		results[i] = i + 1
	}

	return results
}

func parseTokenizedSelection(clean string, maxCount int) ([]int, error) {
	tokens := strings.Split(clean, ",")
	selected := make(map[int]bool)

	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t == "" {
			continue
		}
		if err := processSelectionToken(t, maxCount, selected); err != nil {
			return nil, err
		}
	}

	return collectSortedIndices(selected), nil
}

func processSelectionToken(t string, maxCount int, selected map[int]bool) error {
	hasDash := strings.Contains(t, "-")
	if hasDash {
		return parseRangeToken(t, maxCount, selected)
	}

	idx, err := strconv.Atoi(t)
	if err != nil || idx < 1 || idx > maxCount {
		return apperror.NewValidationError(fmt.Sprintf("invalid index: %s (must be between 1 and %d)", t, maxCount))
	}

	selected[idx] = true

	return nil
}

func parseRangeToken(t string, maxCount int, selected map[int]bool) error {
	parts := strings.SplitN(t, "-", 2)
	start, errStart := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, errEnd := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errStart != nil || errEnd != nil || start < 1 || end > maxCount || start > end {
		return apperror.NewValidationError(fmt.Sprintf("invalid range: %s (valid: 1-%d with start <= end)", t, maxCount))
	}

	for i := start; i <= end; i++ {
		selected[i] = true
	}

	return nil
}

func collectSortedIndices(selected map[int]bool) []int {
	results := make([]int, 0, len(selected))
	for idx := range selected {
		results = append(results, idx)
	}
	sort.Ints(results)

	return results
}
