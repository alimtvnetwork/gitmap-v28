package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// helpJSONGroup is one super-group with its command rows, as emitted
// by `gitmap help --json`. Stable schema for scripting / IDE integrations.
type helpJSONGroup struct {
	Group string   `json:"group"`
	Lines []string `json:"lines"`
}

// helpJSONDoc is the top-level JSON payload returned by --json.
type helpJSONDoc struct {
	Version string          `json:"version"`
	Count   int             `json:"count"`
	Groups  []helpJSONGroup `json:"groups"`
	Query   string          `json:"query,omitempty"`
}

// printUsageJSON emits the full help registry as JSON. When --filter
// is also supplied, only matching rows are included (no ANSI color).
func printUsageJSON(query string) {
	rows := allHelpRows()
	if len(query) > 0 {
		rows = filterRows(rows, query)
	}

	doc := helpJSONDoc{Version: constants.Version, Query: query}
	byGroup := make(map[string]int)
	for _, r := range rows {
		line := stripANSI(strings.TrimRight(r.Line, "\n"))
		if idx, seen := byGroup[r.Group]; seen {
			doc.Groups[idx].Lines = append(doc.Groups[idx].Lines, line)

			continue
		}
		byGroup[r.Group] = len(doc.Groups)
		doc.Groups = append(doc.Groups, helpJSONGroup{
			Group: stripANSI(r.Group),
			Lines: []string{line},
		})
	}
	doc.Count = len(rows)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		fmt.Fprintf(os.Stderr, "gitmap: help --json encode failed: %v\n", err)
	}
}

// stripANSI removes ESC[...m color codes so JSON consumers get clean text.
func stripANSI(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			i = j

			continue
		}
		out.WriteByte(s[i])
	}

	return out.String()
}
