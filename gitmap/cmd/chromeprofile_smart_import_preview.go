package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func renderProfileCandidatesTable(target string, candidates []DiscoveredProfileCandidate) {
	fmt.Printf("\n\033[1;96mChrome Snapshot Inspection & Import Preview (%s)\033[0m\n\n", target)
	for i, c := range candidates {
		printCandidateRow(i+1, c)
	}

	fmt.Printf("Total profiles inspected: %d\n\n", len(candidates))
	fmt.Println("Usage hints:")
	fmt.Println("  Import all:            gitmap chrome import-all <path>")
	fmt.Println("  Import single email:   gitmap chrome profile import <email>")
	fmt.Println("  Import with limit:     gitmap chrome import-all <path> --limit 1")
	fmt.Println("  Import with exclusion: gitmap chrome import-all <path> --except <pattern>")
	fmt.Println()
}

func printCandidateRow(idx int, c DiscoveredProfileCandidate) {
	emailStr := c.Email
	if emailStr == "" {
		emailStr = "(none)"
	}
	dispStr := c.DisplayName
	if dispStr == "" {
		dispStr = "(none)"
	}

	sourceLabel := c.SourcePath
	if c.IsFromZip {
		sourceLabel = fmt.Sprintf("%s [archive entry: %s]", c.ZipArchive, c.ProfileDirName)
	}

	fmt.Printf("  [%d] \033[1;97m%s\033[0m (source: %s)\n", idx, c.ProfileDirName, sourceLabel)
	fmt.Printf("      • Profile: %-12s | Display: %-14s | Email: %s\n", c.ProfileDirName, dispStr, emailStr)
	fmt.Printf("      • Bookmarks: %-10d | Extensions: %-11d\n", c.BookmarksCount, c.ExtensionsCount)

	printCandidateAction(c)
	fmt.Println()
}

func printCandidateAction(c DiscoveredProfileCandidate) {
	if c.TargetDestination.IsNew {
		fmt.Printf("      \033[1;92m+ Action: Will CREATE NEW profile %q\033[0m (%s)\n",
			c.TargetDestination.Dir, c.TargetDestination.Action)

		return
	}

	fmt.Printf("      \033[1;93m↷ Action: Will UPDATE profile %q\033[0m (%s)\n",
		c.TargetDestination.Dir, c.TargetDestination.Action)
}

func renderProfileCandidatesJSON(candidates []DiscoveredProfileCandidate) error {
	raw, err := json.MarshalIndent(candidates, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(raw))

	return nil
}

func isPreflightInspectArg(arg string) bool {
	clean := strings.ToLower(strings.TrimSpace(arg))

	return clean == "ls" || clean == "list" || clean == "inspect" || clean == "check" || clean == "scan" || clean == "status" || clean == "preview"
}

func sortCandidates(c []DiscoveredProfileCandidate) {
	sort.Slice(c, func(i, j int) bool {
		return c[i].ProfileDirName < c[j].ProfileDirName
	})
}
