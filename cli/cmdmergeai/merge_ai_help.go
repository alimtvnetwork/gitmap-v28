// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/pterm/pterm"
)

// RenderMergeAIHelp renders the formatted PTerm CLI manual for gitmap merge-ai.
func RenderMergeAIHelp() {
	headerBox := pterm.DefaultBox.WithTitle(pterm.Cyan("GitMap merge-ai — Multi-Repo AI Merge Orchestrator"))
	pterm.Println()
	headerBox.Println("Amalgamates multiple git repositories into a unified staging repository,\n" +
		"intelligently skips redundant file URLs, resolves chronological precedence,\n" +
		"sequences conflicting paths (01_<file>, 02_<file>), and generates\n" +
		"an AI consolidation manifest and checklist.")

	fmt.Printf("\n%s\n", pterm.Bold.Sprint("USAGE:"))
	fmt.Printf("  gitmap merge-ai [destination] [manifest.txt|sources...] [flags]\n")
	fmt.Printf("  gitmap ma [destination] [manifest.txt] [flags]\n")

	fmt.Printf("\n%s\n", pterm.Bold.Sprint("INPUT MODALITIES & EXAMPLES:"))
	fmt.Printf("  %s %s\n", pterm.Yellow("1. Manifest file with default target:"),
		pterm.Gray("gitmap merge-ai repo-cache/merge-ai-icon-studio.txt"))
	fmt.Printf("  %s %s\n", pterm.Yellow("2. Target repo + manifest file:"),
		pterm.Gray("gitmap merge-ai icon-studio-v1 merge-ai-icon-studio.txt"))
	fmt.Printf("  %s %s\n", pterm.Yellow("3. Multi-word destination name:"),
		pterm.Gray(`gitmap merge-ai "Icon Studio v1" merge-ai-icon-studio.txt`))
	fmt.Printf("  %s %s\n", pterm.Yellow("4. Named parameter flags:"),
		pterm.Gray("gitmap merge-ai -d icon-studio-v1 -s merge-ai-icon-studio.txt"))
	fmt.Printf("  %s %s\n", pterm.Yellow("5. Comma-delimited repo sources:"),
		pterm.Gray("gitmap merge-ai icon-studio-v1 https://github.com/org/a,https://github.com/org/b"))
	fmt.Printf("  %s %s\n", pterm.Yellow("6. JSON configuration file:"),
		pterm.Gray("gitmap merge-ai --config merge-config.json"))

	fmt.Printf("\n%s\n", pterm.Bold.Sprint("FLAGS:"))
	fmt.Printf("  %-24s %s\n", "-d, --destination", "Target repository directory or slug (default: icon-studio-v1)")
	fmt.Printf("  %-24s %s\n", "-s, --source, --url", "Source manifest file (.txt), repo URL, or comma-separated list")
	fmt.Printf("  %-24s %s\n", "-c, --config", "Path to JSON config defining destination and sources")
	fmt.Printf("  %-24s %s\n", "--dry-run", "Simulate inspection, collision detection, and skips without modifying disk")
	fmt.Printf("  %-24s %s\n", "-h, --help", "Display this formatted manual")

	fmt.Printf("\n%s\n", pterm.Bold.Sprint("COLLISION & MANIFEST STRATEGY:"))
	fmt.Printf("  • Redundant file URLs belonging to an included repo are automatically skipped.\n")
	fmt.Printf("  • Non-conflicting unique files are placed directly at their canonical path.\n")
	fmt.Printf("  • Conflicting files are sequenced chronologically: 01_<file>, 02_<file>, ...\n")
	fmt.Printf("  • Generated artifacts:\n")
	fmt.Printf("      - %s: Structured machine-readable schema & file collision registry\n", pterm.Cyan("merge-ai-manifest.json"))
	fmt.Printf("      - %s: Action checklist for an AI agent to consolidate code\n", pterm.Cyan("instruction.md"))

	fmt.Printf("\n%s gitmap %s\n\n", pterm.Gray("GitMap CLI"), pterm.Cyan(constants.Version))
}
