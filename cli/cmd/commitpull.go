package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runCommitPull executes `gitmap commit-pull` / `cpull` / `pull-commits`.
// Replays commits chronologically while creating Pull Requests for all merges and releases.
func runCommitPull(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			printCommitPullHelp()
			cliexit.HandleError(nil, 0)
			return nil
		}
	}

	if len(args) == 1 && args[0] == "--tree" {
		printCommitPullTree()
		return nil
	}

	adjusted := append([]string{"--pr=merges"}, args...)
	return runCommitIn(adjusted)
}

func printCommitPullHelp() {
	fmt.Println(`gitmap commit-pull (cpull, pull-commits): PR Replay & Migration Engine
========================================================================

OVERVIEW:
  Chronologically replays commits from multiple source repositories (or a range
  such as v2..v28) into a target repository, automatically creating Pull Requests
  for every branch merge and version release tag.

DECLARATIVE CONFIG JSON (--config / -c):
  Run an entire multi-repository migration with a single config JSON file that
  auto-creates the target repository, imports state templates & variables with
  hash deduplication, strips unwanted lines (starts_with, ends_with, contains,
  regex), replaces generic 'Changes' titles with $files.2.names, and enters --cd:
    gitmap commit-pull --config .ai-memory/temp/commit-pull-config.json

MULTI-REPO RANGE EXPANSION (V2 TO V28):
  You can pass dynamic range expansions to automatically pull and convert
  all 28 repository histories into a single repository:
    gitmap commit-pull "D:\target" "https://github.com/alimtvnetwork/gitmap-v{2..28}" --tree

STATE TEMPLATES (--seo-template):
  Loads pre-compiled templates and variables from gitmap-templates.db:
    gitmap templates import .ai-memory/temp/seo-templates.json
    gitmap commit-pull "D:\target" v2..v28 --seo-template seo

FLAGS:
  --config, -c <file.json>  Declarative migration config (imports, skippers, title rules)
  --tree                    Display the commit/PR branch dependency tree
  --seo-template <cat>      Load pre-compiled templates from state DB category
  --pr merges               Simulate feature branches and PR merges (default)
  --cd                      Change directory into target repository upon completion
  --final-sync              Synchronize target repository with IDE/Desktop tools
  --exclude <glob>          Exclude specific files or directories
  --dry-run                 Simulate without modifying repositories`)
}

func printCommitPullTree() {
	fmt.Println(constants.ColorCyan + `
  ┌── Migration & Commit-Pull Replay Tree ────────────────────────────────┐
  │ Target: Target Mainline (main)                                        │
  │ PR Engine: Simulated Feature Branches (--pr merges)                   │
  │ Templates: Pre-Compiled State DB Engine (gitmap-templates.db)         │
  └───────────────────────────────────────────────────────────────────────┘` + constants.ColorReset)
	fmt.Print(`
  * [v1.0.0-legacy] PR #1: git-repo-navigator (legacy origin)
  │ \
  │  * feat(nav): core scanner and discovery engine
  │  * feat(nav): multi-threaded filesystem indexing
  │ /
  * [v2.0.0-legacy] PR #2: gitmap-v2 (split sqlite architecture)
  │ \
  │  * feat(db): three-tier split SQLite schema (.gitmap/data/)
  │  * feat(ui): terminal rendering & termpad table formatting
  │ /
  * ... [PR #3 .. PR #27: Sequential V3 to V27 Architecture Iterations]
  * [v28.0.0-current] PR #28: gitmap-v28 (modern AGY & AI self-healing)
  │ \
  │  * feat(agy): prompt injection, Lapp/Wapp aggregation
  │  * feat(pl): dynamic pipeline-ai waiting & 4-part RCA
  │  * feat(templates): state DB pre-compiled variables & title synthesis
  │ /
  * (main) HEAD: Complete 28-Repository Replay & Release Stack` + "\n\n")
}
