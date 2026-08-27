package searcher

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/lazyregex"
)

// SearchExact searches for the exact query in content and returns results.
func SearchExact(content, query, absPath, relPath string) []SearchResult {
	var results []SearchResult
	if query == "" {
		return results
	}
	
	offset := 0
	for {
		idx := strings.Index(content[offset:], query)
		if idx == -1 {
			break
		}
		
		start := offset + idx
		end := start + len(query)
		
		results = append(results, SearchResult{
			MatchedText:   query,
			StartPosition: start,
			EndPosition:   end,
			FilePath:      absPath,
			RelativePath:  relPath,
		})
		
		offset = end
	}
	return results
}

// SearchRegex searches for regex pattern in content and returns results.
func SearchRegex(content string, lz *lazyregex.LazyRegexp, absPath, relPath string) []SearchResult {
	var results []SearchResult
	matches := lz.Re().FindAllStringIndex(content, -1)
	
	for _, m := range matches {
		if len(m) == 2 {
			start := m[0]
			end := m[1]
			results = append(results, SearchResult{
				MatchedText:   content[start:end],
				StartPosition: start,
				EndPosition:   end,
				FilePath:      absPath,
				RelativePath:  relPath,
			})
		}
	}
	return results
}
