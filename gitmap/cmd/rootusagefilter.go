package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// resolveFilterQuery extracts the value of --filter / -f from os.Args.
func resolveFilterQuery() string {
	args := os.Args[2:]
	for i, a := range args {
		if (a == constants.FlagFilter || a == constants.FlagFilterShort) && i+1 < len(args) {
			return strings.TrimSpace(args[i+1])
		}
		if a == constants.FlagFilter || a == constants.FlagFilterShort {
			return ""
		}
		if v, ok := strings.CutPrefix(a, constants.FlagFilter+"="); ok {
			return strings.TrimSpace(v)
		}
		if v, ok := strings.CutPrefix(a, constants.FlagFilterShort+"="); ok {
			return strings.TrimSpace(v)
		}
	}

	return ""
}

// printUsageFiltered renders only rows whose group or command line
// contains the (case-insensitive) query.
func printUsageFiltered(query string) {
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)

	if len(query) == 0 {
		fmt.Println("  " + constants.ColorYellow +
			"--filter requires a query (e.g. `gitmap help --filter ssh`)" +
			constants.ColorReset)

		return
	}

	rows := allHelpRows()
	hits := filterRows(rows, query)
	if len(hits) == 0 {
		printNoFilterMatches(rows, query)

		return
	}

	renderFilteredGroups(hits, query)
	printFilterRecapBanner(hits, query)
	printUsageFooter()
}

// printFilterRecapBanner repeats the matched command lines in a tight
// block at the bottom of filtered help.
func printFilterRecapBanner(hits []helpRow, query string) {
	if len(hits) == 0 {
		return
	}
	const cap = 10
	bar := strings.Repeat("─", 12)
	fmt.Println()
	fmt.Printf("  %s%s matches for %q %s%s\n",
		constants.ColorMagenta, bar, query, bar, constants.ColorReset)
	shown := hits
	if len(shown) > cap {
		shown = shown[:cap]
	}
	for _, r := range shown {
		fmt.Println(highlight(strings.TrimRight(r.Line, "\n"), query))
	}
	if len(hits) > cap {
		fmt.Printf("  %s… +%d more (refine with a tighter --filter)%s\n",
			constants.ColorDim, len(hits)-cap, constants.ColorReset)
	}
	fmt.Println()
}

func filterRows(rows []helpRow, query string) []helpRow {
	needle := strings.ToLower(query)
	out := make([]helpRow, 0, 16)
	for _, r := range rows {
		hay := strings.ToLower(r.Group + " " + r.Line)
		if strings.Contains(hay, needle) {
			out = append(out, r)
		}
	}

	return out
}

func renderFilteredGroups(hits []helpRow, query string) {
	groupOrder := make([]string, 0, 8)
	byGroup := make(map[string][]string)
	for _, r := range hits {
		if _, seen := byGroup[r.Group]; !seen {
			groupOrder = append(groupOrder, r.Group)
		}
		byGroup[r.Group] = append(byGroup[r.Group], highlight(r.Line, query))
	}

	fmt.Printf("  %sMatches for%s %q  (%d found)\n\n",
		constants.ColorCyan, constants.ColorReset, query, len(hits))
	for _, g := range groupOrder {
		fmt.Println(colorGroupHeader(g))
		fmt.Println()
		for _, ln := range byGroup[g] {
			fmt.Println(ln)
		}
		fmt.Println()
	}
}

// highlight wraps every case-insensitive occurrence of query in a
// bold yellow ANSI marker so matches pop in the rendered output.
func highlight(line, query string) string {
	if len(query) == 0 {
		return line
	}
	lowLine := strings.ToLower(line)
	lowQ := strings.ToLower(query)
	var out strings.Builder
	idx := 0
	for {
		hit := strings.Index(lowLine[idx:], lowQ)
		if hit < 0 {
			out.WriteString(line[idx:])

			break
		}
		out.WriteString(line[idx : idx+hit])
		out.WriteString(constants.ColorYellow)
		out.WriteString(line[idx+hit : idx+hit+len(query)])
		out.WriteString(constants.ColorReset)
		idx += hit + len(query)
	}

	return out.String()
}
