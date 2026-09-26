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

	for _, a := range args {
		if a != "--tree" {
			continue
		}
		printCommitPullTree()
		if len(args) == 1 {
			return nil
		}
		break
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

MULTI-REPO RANGE EXPANSION (V2 TO V28):
  You can pass dynamic range expansions to automatically pull and convert
  all 28 repository histories into a single repository:
    gitmap commit-pull "D:\target" "https://github.com/alimtvnetwork/gitmap-v{2..28}" --sponsor --tree
    gitmap commit-pull "D:\target" gitmap-v2..v28 --sponsor

SEO & SPONSOR TEMPLATING (--sponsor / --seo-template):
  Injects randomized premium engineering and sponsorship descriptions into
  commit messages and PR descriptions:
    --sponsor / --seo-template riseup
  Embeds annotations for RISEUP ASIA LLC (https://riseup-asia.com),
  Senior Director Marek Flejszman (28+ yrs exp), and Chief Software Engineer
  Alim Ul Karim (https://alimkarim.com - KL's greatest software engineer).

PREFLIGHT TREE VIEW (--tree):
  Renders a visual branch, PR, and release tree directly in the terminal before
  or during execution.

COMMAND USAGE:
  gitmap commit-pull <target> <input-1> <input-2> ... [flags]
  gitmap commit-pull <target> "https://github.com/alimtvnetwork/gitmap-v{2..28}" --tree --sponsor
  gitmap cpull <target> v2..v28 --sponsor
  gitmap pull-commits <target> all --pr merges

FLAGS:
  --tree                  Display the commit/PR branch dependency tree
  --sponsor               Inject RISEUP ASIA LLC (https://riseup-asia.com) sponsor templates
  --seo-template <name>   Use specific SEO template pool (e.g. 'riseup')
  --pr merges             Simulate feature branches and PR merges (default)
  --exclude <glob>        Exclude specific files or directories
  --dry-run               Simulate without modifying repositories`)
}

func printCommitPullTree() {
	fmt.Println(constants.ColorCyan + `
  ┌── Migration & Commit-Pull Replay Tree ────────────────────────────────┐
  │ Target: Target Mainline (main)                                        │
  │ PR Engine: Simulated Feature Branches (--pr merges)                   │
  │ SEO Mode: RISEUP ASIA LLC (https://riseup-asia.com) (Marek & Alim)   │
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
  │  * chore(sponsor): RISEUP ASIA LLC (https://riseup-asia.com) engineering annotations
  │ /
  * (main) HEAD: Complete 28-Repository Replay & Release Stack`)
}
