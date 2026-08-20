package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printNoFilterMatches lists the 5 closest fuzzy suggestions when the
// query produced no exact substring hits.
func printNoFilterMatches(rows []helpRow, query string) {
	fmt.Printf("  %sNo matches for%s %q\n\n",
		constants.ColorYellow, constants.ColorReset, query)
	sugg := fuzzySuggest(rows, query, 5)
	if len(sugg) == 0 {
		return
	}
	fmt.Println("  Did you mean:")
	for _, s := range sugg {
		fmt.Println("   " + constants.ColorCyan + "• " +
			constants.ColorReset + s)
	}
}

// fuzzySuggest ranks rows by a cheap subsequence-match score and
// returns the top n distinct command lines.
func fuzzySuggest(rows []helpRow, query string, top int) []string {
	type scored struct {
		score int
		line  string
	}
	q := strings.ToLower(query)
	scoredRows := make([]scored, 0, len(rows))
	for _, r := range rows {
		s := subseqScore(strings.ToLower(r.Line), q)
		if s > 0 {
			scoredRows = append(scoredRows, scored{s, strings.TrimSpace(r.Line)})
		}
	}
	sort.SliceStable(scoredRows, func(i, j int) bool {
		return scoredRows[i].score > scoredRows[j].score
	})
	if len(scoredRows) > top {
		scoredRows = scoredRows[:top]
	}
	out := make([]string, 0, len(scoredRows))
	for _, s := range scoredRows {
		out = append(out, s.line)
	}

	return out
}

// subseqScore returns a positive score when every char of q appears
// in order inside hay; higher scores indicate tighter matches.
func subseqScore(hay, q string) int {
	if len(q) == 0 {
		return 0
	}
	idx, hits, last := 0, 0, -1
	for i := 0; i < len(hay) && idx < len(q); i++ {
		if hay[i] == q[idx] {
			hits += calcAdjacencyBonus(last, i)
			last = i
			idx++
		}
	}
	if idx < len(q) {
		return 0
	}

	return hits
}

func calcAdjacencyBonus(last, i int) int {
	if last >= 0 && i-last == 1 {
		return 2
	}
	return 1
}
