package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// dpDiff holds a single repo that differs between profiles.
type dpDiff struct {
	Name  string `json:"name"`
	PathA string `json:"pathA"`
	PathB string `json:"pathB"`
	ModeA string `json:"modeA,omitempty"`
	ModeB string `json:"modeB,omitempty"`
}

// dpResult holds the full comparison outcome.
type dpResult struct {
	onlyInA   []model.ScanRecord
	onlyInB   []model.ScanRecord
	different []dpDiff
	same      []model.ScanRecord
}

// compareDPRepos categorizes repos from two profiles.
func compareDPRepos(reposA, reposB []model.ScanRecord) dpResult {
	mapB := indexByName(reposB)
	var result dpResult

	for _, a := range reposA {
		b, found := mapB[a.RepoName]
		if !found {
			result.onlyInA = append(result.onlyInA, a)

			continue
		}

		categorizeDPMatch(a, b, &result)
		delete(mapB, a.RepoName)
	}

	for _, b := range mapB {
		result.onlyInB = append(result.onlyInB, b)
	}

	return result
}

// categorizeDPMatch checks if two matching repos are same or different.
func categorizeDPMatch(a, b model.ScanRecord, result *dpResult) {
	if a.AbsolutePath == b.AbsolutePath && a.HTTPSUrl == b.HTTPSUrl {
		result.same = append(result.same, a)

		return
	}

	result.different = append(result.different, dpDiff{
		Name:  a.RepoName,
		PathA: a.AbsolutePath,
		PathB: b.AbsolutePath,
		ModeA: detectMode(a),
		ModeB: detectMode(b),
	})
}

// indexByName creates a map of RepoName → ScanRecord.
func indexByName(repos []model.ScanRecord) map[string]model.ScanRecord {
	m := make(map[string]model.ScanRecord, len(repos))

	for _, r := range repos {
		m[r.RepoName] = r
	}

	return m
}

// detectMode returns "https" or "ssh" based on the clone URL.
func detectMode(r model.ScanRecord) string {
	if len(r.SSHUrl) > 0 && len(r.HTTPSUrl) == 0 {
		return constants.ModeSSH
	}

	return constants.ModeHTTPS
}

// printDPOutput renders the comparison to stdout.
func printDPOutput(nameA, nameB string, result dpResult, showAll bool) {
	fmt.Printf(constants.MsgDPHeader, nameA, nameB)

	printDPOnlyIn(nameA, result.onlyInA)
	printDPOnlyIn(nameB, result.onlyInB)
	printDPDiffs(result.different)
	printDPSame(result.same, showAll)

	printDPSummary(result)
}

// printDPOnlyIn prints repos only in one profile.
func printDPOnlyIn(name string, repos []model.ScanRecord) {
	if len(repos) == 0 {
		return
	}

	fmt.Printf(constants.MsgDPOnlyInHeader, name)
	fmt.Println()

	for _, r := range repos {
		fmt.Printf(constants.MsgDPOnlyInRowFmt, r.RepoName, r.AbsolutePath)
	}
}

// printDPDiffs prints repos that differ between profiles.
func printDPDiffs(diffs []dpDiff) {
	if len(diffs) == 0 {
		return
	}

	fmt.Println(constants.MsgDPDiffHeader)

	for _, d := range diffs {
		fmt.Printf(constants.MsgDPDiffNameFmt, d.Name)
		fmt.Printf(constants.MsgDPDiffDetailFmt, "left:", d.PathA+" ("+d.ModeA+")")
		fmt.Printf(constants.MsgDPDiffDetailFmt, "right:", d.PathB+" ("+d.ModeB+")")
	}
}

// printDPSame prints identical repos (only with --all flag).
func printDPSame(repos []model.ScanRecord, showAll bool) {
	if !showAll {
		fmt.Printf(constants.MsgDPSameFmt, len(repos))

		return
	}

	fmt.Println(constants.MsgDPSameAllHeader)

	for _, r := range repos {
		fmt.Printf(constants.MsgDPSameRowFmt, r.RepoName, r.AbsolutePath)
	}
}

// printDPSummary prints the final summary line.
func printDPSummary(result dpResult) {
	fmt.Printf(constants.MsgDPSummaryFmt,
		len(result.onlyInA), len(result.onlyInB),
		len(result.different), len(result.same))
}
