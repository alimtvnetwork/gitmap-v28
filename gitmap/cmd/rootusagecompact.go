package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// compactGroup maps a group key to its header and compact line.
type compactGroup struct {
	Header  string
	Compact string
}

// compactGroups returns the ordered list of all groups.
func compactGroups() []compactGroup {
	return []compactGroup{
		{constants.HelpGroupScanning, constants.CompactScanning},
		{constants.HelpGroupCloning, constants.CompactCloning},
		{constants.HelpGroupGitOps, constants.CompactGitOps},
		{constants.HelpGroupNavigation, constants.CompactNavigation},
		{constants.HelpGroupRelease, constants.CompactRelease},
		{constants.HelpGroupReleaseInfo, constants.CompactRelInfo},
		{constants.HelpGroupData, constants.CompactData},
		{constants.HelpGroupHistory, constants.CompactHistory},
		{constants.HelpGroupAmendGroup, constants.CompactAmend},
		{constants.HelpGroupProject, constants.CompactProject},
		{constants.HelpGroupSSH, constants.CompactSSH},
		{constants.HelpGroupZip, constants.CompactZip},
		{constants.HelpGroupEnvTools, constants.CompactEnvTools},
		{constants.HelpGroupTasks, constants.CompactTasks},
		{constants.HelpGroupVisualize, constants.CompactVisualize},
		{constants.HelpGroupCommitXfer, constants.CompactCommitXfer},
		{constants.HelpGroupUtilities, constants.CompactUtilities},
	}
}

// printUsageCompact prints a minimal command list, optionally filtered by group.
func printUsageCompact() {
	filter := resolveCompactFilter()

	fmt.Printf(constants.UsageHeaderFmt, constants.Version)

	if len(filter) > 0 {
		printCompactFiltered(filter)

		return
	}

	printCompactAll()
}

// resolveCompactFilter extracts the group filter from os.Args (skips flags).
func resolveCompactFilter() string {
	for _, arg := range os.Args[2:] {
		if !strings.HasPrefix(arg, "-") {
			return strings.ToLower(arg)
		}
	}

	return ""
}

// printCompactAll prints all groups in compact mode.
func printCompactAll() {
	fmt.Println(constants.HelpUsage)
	fmt.Println()

	for _, g := range compactGroups() {
		fmt.Println(g.Header)
		fmt.Println(g.Compact)
	}

	fmt.Println()
	fmt.Println(constants.HelpGroupHint)
}

// printCompactFiltered prints only groups matching the filter keyword.
func printCompactFiltered(filter string) {
	matched := false

	for _, g := range compactGroups() {
		headerLower := strings.ToLower(g.Header)
		if !strings.Contains(headerLower, filter) {
			continue
		}

		fmt.Println(g.Header)
		fmt.Println(g.Compact)
		matched = true
	}

	if !matched {
		fmt.Printf(constants.CompactNoMatchFmt, filter)
		fmt.Println()
		printCompactAll()
	}
}

// printHelpGroups lists all available group names for quick reference.
func printHelpGroups() {
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)
	fmt.Println("  Available help groups:")
	fmt.Println()

	for _, key := range constants.HelpGroupKeys {
		fmt.Printf("    %s\n", key)
	}

	fmt.Println()
	fmt.Println("  Usage: gitmap help --compact <group>")
}
