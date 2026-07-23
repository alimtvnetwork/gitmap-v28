package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// helpRow is a single command/flag line tagged with its sub-group
// header so filtered output keeps its context.
type helpRow struct {
	Group string
	Line  string
}

// allHelpRows returns every command + flag row rendered by the
// full `gitmap help` screen. Source of truth for `--filter` search.
//
// New groups must register here too — there is no automatic discovery.
func allHelpRows() []helpRow {
	rows := make([]helpRow, 0, 128)
	addGroup(&rows, constants.HelpGroupScanning,
		constants.HelpScan, constants.HelpRescan, constants.HelpList)
	addGroup(&rows, constants.HelpGroupCloning,
		constants.HelpClone, constants.HelpCloneNext,
		constants.HelpDesktopSync, constants.HelpGitHubDesktop)
	addGroup(&rows, constants.HelpGroupGitOps,
		constants.HelpPull, constants.HelpExec, constants.HelpStatus,
		constants.HelpWatch, constants.HelpHasAnyUpdates, constants.HelpLatestBr)
	addGroup(&rows, constants.HelpGroupNavigation,
		constants.HelpCD, constants.HelpGroup, constants.HelpMultiGroup,
		constants.HelpSf, constants.HelpAlias, constants.HelpDiffProfiles)
	addGroup(&rows, constants.HelpGroupRelease,
		constants.HelpRelease, constants.HelpReleasePull,
		constants.HelpReleaseSelf, constants.HelpReleaseBr, constants.HelpTempRelease)
	addGroup(&rows, constants.HelpGroupReleaseInfo,
		constants.HelpChangelog, constants.HelpChangelogGen,
		constants.HelpListVersions, constants.HelpListReleases,
		constants.HelpReleasePend, constants.HelpRevert,
		constants.HelpClearReleaseJSON, constants.HelpPrune)
	addGroup(&rows, constants.HelpGroupData,
		constants.HelpExport, constants.HelpImport, constants.HelpProfile,
		constants.HelpBookmark, constants.HelpRm, constants.HelpDBReset)
	addGroup(&rows, constants.HelpGroupHistory,
		constants.HelpHistory, constants.HelpHistoryReset,
		constants.HelpVersionHistory, constants.HelpStats)
	addGroup(&rows, constants.HelpGroupAmendGroup,
		constants.HelpAmend, constants.HelpAmendList)
	addGroup(&rows, constants.HelpGroupProject,
		constants.HelpGoRepos, constants.HelpNodeRepos,
		constants.HelpReactRepos, constants.HelpCppRepos, constants.HelpCsharpRepos)
	addGroup(&rows, constants.HelpGroupSSH, constants.HelpSSH)
	addGroup(&rows, constants.HelpGroupZip, constants.HelpZipGroup)
	addGroup(&rows, constants.HelpGroupEnvTools,
		constants.HelpEnv, constants.HelpInstall, constants.HelpUninstall)
	addGroup(&rows, constants.HelpGroupTasks,
		constants.HelpTask, constants.HelpPending, constants.HelpDoPending)
	addGroup(&rows, constants.HelpGroupVisualize, constants.HelpDashboard)
	addGroup(&rows, constants.HelpGroupCommitXfer,
		constants.HelpCommitRight, constants.HelpCommitLeft, constants.HelpCommitBoth)
	addGroup(&rows, constants.HelpGroupChromeProf,
		constants.HelpChromeProfileCopy, constants.HelpChromeProfileExport,
		constants.HelpChromeProfileImport, constants.HelpChromeProfileList,
		constants.HelpChromeProfileDelete)
	addGroup(&rows, constants.HelpGroupTemplates,
		constants.HelpAddIgnore, constants.HelpAddAttributes,
		constants.HelpAddLFSInstall, constants.HelpTemplatesInit,
		constants.HelpTemplatesList, constants.HelpTemplatesShow,
		constants.HelpTemplatesDiff, constants.HelpSync, constants.HelpCommons)
	addGroup(&rows, constants.HelpGroupUtilities,
		constants.HelpSetup, constants.HelpDoctor, constants.HelpUpdate,
		constants.HelpUpdateCleanup, constants.HelpVersion, constants.HelpCompletion,
		constants.HelpInteractive, constants.HelpDocs, constants.HelpHelpDash,
		constants.HelpGoMod, constants.HelpSEOWrite, constants.HelpLLMDocs,
		constants.HelpFixRepo, constants.HelpMakePublic, constants.HelpMakePrivate,
		constants.HelpCloneFixRepo, constants.HelpCloneFixRepoPub,
		constants.HelpCmdOpen, constants.HelpHelp)

	return rows
}

func addGroup(rows *[]helpRow, group string, lines ...string) {
	for _, ln := range lines {
		*rows = append(*rows, helpRow{Group: group, Line: ln})
	}
}

// resolveFilterQuery extracts the value of --filter / -f from os.Args.
// Accepts both --filter=foo and --filter foo forms. Returns "" if absent.
func resolveFilterQuery() string {
	args := os.Args[2:]
	for i, a := range args {
		if a == constants.FlagFilter || a == constants.FlagFilterShort {
			if i+1 < len(args) {
				return strings.TrimSpace(args[i+1])
			}

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
// contains the (case-insensitive) query. Matches are highlighted in
// yellow. When zero rows match, fuzzy suggestions are offered.
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
// block at the very bottom of filtered help so the user sees the hits
// without scrolling back up. Capped at 10 rows to stay terminal-sized.
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
			hits++
			if last >= 0 && i-last == 1 {
				hits++ // adjacency bonus
			}
			last = i
			idx++
		}
	}
	if idx < len(q) {
		return 0
	}

	return hits
}
