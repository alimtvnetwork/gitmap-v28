package cmdai

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunAiList filters and renders the registered AI script catalog and frequent command history.
func RunAiList(opts ScriptListOptions) *apperror.AppError {
	if opts.IsFrequentOnly {
		return RunAiFrequent(25, opts.IsCopyClipboard)
	}

	_ = RunAiFrequent(10, opts.IsCopyClipboard)

	scripts := AllScripts()
	filtered := filterScripts(scripts, opts)

	if opts.IsJsonFormat {
		return renderScriptsJson(filtered)
	}

	renderScriptsTable(filtered)

	return nil
}

func filterScripts(scripts []ScriptMetadata, opts ScriptListOptions) []ScriptMetadata {
	var results []ScriptMetadata

	for _, s := range scripts {
		isMatch := isScriptMatching(s, opts)
		if isMatch {
			results = append(results, s)
		}
	}

	return results
}

func isScriptMatching(s ScriptMetadata, opts ScriptListOptions) bool {
	isCatMatch := matchesCategory(s, opts.CategoryFilter)
	isFixMatch := matchesFixOnly(s, opts.IsFixOnly)
	isQueryMatch := matchesQuery(s, opts.SearchQuery)

	return isCatMatch && isFixMatch && isQueryMatch
}

func matchesCategory(s ScriptMetadata, category string) bool {
	isEmptyCat := len(category) == 0
	if isEmptyCat {
		return true
	}

	clean := strings.ToLower(strings.TrimSpace(category))

	return strings.ToLower(string(s.Category)) == clean
}

func matchesFixOnly(s ScriptMetadata, isFixOnly bool) bool {
	if isFixOnly {
		return s.HasFixMode
	}

	return true
}

func matchesQuery(s ScriptMetadata, query string) bool {
	isEmptyQuery := len(query) == 0
	if isEmptyQuery {
		return true
	}

	clean := strings.ToLower(strings.TrimSpace(query))
	isNumMatch := strings.Contains(strings.ToLower(s.Number), clean)
	isSlugMatch := strings.Contains(strings.ToLower(s.Slug), clean)
	isFileMatch := strings.Contains(strings.ToLower(s.Filename), clean)
	isDescMatch := strings.Contains(strings.ToLower(s.Description), clean)
	isAliasMatch := matchesAliasQuery(s.Aliases, clean)

	return isNumMatch || isSlugMatch || isFileMatch || isDescMatch || isAliasMatch
}

func matchesAliasQuery(aliases []string, query string) bool {
	for _, a := range aliases {
		isMatch := strings.Contains(strings.ToLower(a), query)
		if isMatch {
			return true
		}
	}

	return false
}
