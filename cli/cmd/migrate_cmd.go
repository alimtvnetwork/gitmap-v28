package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunMigrate handles the `gitmap migrate` subcommand.
func RunMigrate(args []string) error {
	sub := "graph"
	if len(args) > 0 {
		sub = args[0]
	}

	switch sub {
	case "graph", "preflight", "plan", "status":
		return renderMigrationPreflightGraph()
	default:
		return renderMigrationPreflightGraph()
	}
}

func renderMigrationPreflightGraph() error {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║     GITMAP REPOSITORY REMAPPING & STACKED CONSOLIDATION PREFLIGHT GRAPH      ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  %sTarget Repository:%s   %salimtvnetwork/gitmap-v28%s (Consolidated Core)\n",
		constants.ColorWhite, constants.ColorReset, constants.ColorGreen+"\033[1m", constants.ColorReset)
	fmt.Printf("  %sSource Repositories:%s %salimtvnetwork/git-repo-navigator%s, %salimtvnetwork/gitmap-v2%s\n\n",
		constants.ColorWhite, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset)

	fmt.Printf("  %sREMAP TOPOLOGY & STACKED PR EXECUTION GRAPH:%s\n\n", constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  %salimtvnetwork/gitmap-v28 (main)%s\n", constants.ColorGreen+"\033[1m", constants.ColorReset)
	fmt.Printf("  ├── %s[Phase 1] Remote History Remap & Isolation%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("  │   ├── git remote add legacy-navigator https://github.com/alimtvnetwork/git-repo-navigator.git\n")
	fmt.Printf("  │   ├── git remote add intermediate-v2 https://github.com/alimtvnetwork/gitmap-v2.git\n")
	fmt.Printf("  │   └── git fetch --all --tags\n")
	fmt.Printf("  │\n")
	fmt.Printf("  ├── %s[Phase 2] Consolidated Base Branch%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("  │   └── %smerge/consolidated-v28-foundation%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  │\n")
	fmt.Printf("  ├── %s[Phase 3] Stacked PR Pipeline & Approval Ceremony%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("  │   ├── %sPR #1:%s Core Engine & High-Speed Cloner Sync\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  │   │   ├── Preserves author commit history from git-repo-navigator\n")
	fmt.Printf("  │   │   └── Approval Gate: CI Test Matrix Green\n")
	fmt.Printf("  │   │\n")
	fmt.Printf("  │   ├── %sPR #2:%s Split SQLite Database Architecture\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  │   │   ├── Migrates gitmap.db, installation.db, repodb schemas\n")
	fmt.Printf("  │   │   └── Approval Gate: Zero Schema Drift Verified\n")
	fmt.Printf("  │   │\n")
	fmt.Printf("  │   ├── %sPR #3:%s AGY Prompt Automation & Multi-Node Cluster Engine\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  │   │   ├── Terminal prompt watcher, injection engine, and cache pruner\n")
	fmt.Printf("  │   │   └── Approval Gate: E2E Live Discovery Verified\n")
	fmt.Printf("  │   │\n")
	fmt.Printf("  │   └── %sPR #4:%s Complete 80-Template Catalog & Rise Up Asia LLC Sponsorship\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  │       ├── 20 Verification, 20 UI/UX, 20 Sponsor, 20 PR description templates\n")
	fmt.Printf("  │       └── Approval Gate: Full Catalog E2E Passed\n")
	fmt.Printf("  │\n")
	fmt.Printf("  └── %s[Phase 4] Automated Release Orchestration%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("      ├── Version Bump: %sv28.0.0%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("      ├── Release Tag: git tag -a v28.0.0 -m \"Release v28.0.0: Consolidated Architecture\"\n")
	fmt.Printf("      └── Package & Documentation Synchronization\n\n")

	planPath := filepath.Join(".ai-memory", "temp", "migration-gitmap-plan.md")
	fmt.Printf("  %s Full markdown migration plan available at: %s%s%s\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, constants.ColorYellow, planPath, constants.ColorReset)

	return nil
}
