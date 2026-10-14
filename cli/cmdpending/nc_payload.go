package cmdpending
import (
	"strings"
	"time"
)
func buildNewCommandsPayload(catalog []NewCommandEntry, opts NewCommandsOptions) NewCommandsPayload {
	filtered := filterNewCommands(catalog, opts)
	return NewCommandsPayload{
		Timestamp:        time.Now().UTC(),
		TotalCommands:    len(catalog),
		FilteredCommands: len(filtered),
		Limit:            opts.Limit,
		Category:         opts.Category,
		Filter:           opts.Filter,
		Commands:         filtered,
	}
}

func filterNewCommands(catalog []NewCommandEntry, opts NewCommandsOptions) []NewCommandEntry {
	res := make([]NewCommandEntry, 0, len(catalog))
	for _, entry := range catalog {
		if isCommandMatched(entry, opts) {
			res = append(res, entry)
		}
		if len(res) >= opts.Limit {
			break
		}
	}
	return res
}

func isCommandMatched(entry NewCommandEntry, opts NewCommandsOptions) bool {
	if !isCategoryMatch(entry.Category, opts.Category) {
		return false
	}
	if !isTextFilterMatch(entry, opts.Filter) {
		return false
	}
	return true
}

func isCategoryMatch(entryCat, filterCat string) bool {
	trimmed := strings.TrimSpace(filterCat)
	if trimmed == "" {
		return true
	}
	return strings.EqualFold(entryCat, trimmed)
}

func isTextFilterMatch(entry NewCommandEntry, filter string) bool {
	query := strings.ToLower(strings.TrimSpace(filter))
	if query == "" {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Alias), query) {
		return true
	}
	return isEntryDetailMatch(entry, query)
}

func isEntryDetailMatch(entry NewCommandEntry, query string) bool {
	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Example), query) {
		return true
	}
	return false
}
