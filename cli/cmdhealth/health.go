package cmdhealth

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const healthHelp = `health — daily workstation health aggregation.

Description:
  Read-only rollup of the morning checks: dirty repos, failing
  pipelines, fleet node count, and stale databases. Points at the
  worst item first with the command that digs deeper.

Usage:
  gitmap health [--workdir <dir>]
  gitmap hlth [--workdir <dir>]

Flags:
  --workdir <dir>   repo workdir (default: $GITMAP_WORKDIR or ~/workspace/repos)
  --help, -h        show this help

Example:
  gitmap health
`

// RunHealth implements `gitmap health` (alias: hlth).
func RunHealth(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(healthHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("health", flag.ContinueOnError)
	workdirFlag := fs.String("workdir", "", "repo workdir")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(healthHelp)
		return apperror.WrapSimple(parseErr, "health.flags")
	}
	workdir := resolveHealthWorkdir(*workdirFlag)
	facts := collectHealthFacts(workdir)

	fmt.Println("workstation health:")
	fmt.Printf("  repos:     %d (%d dirty", facts.repoTotal, facts.repoDirty)
	if len(facts.dirtyNames) > 0 {
		fmt.Printf(": %s", strings.Join(facts.dirtyNames, ", "))
	}
	fmt.Println(")")
	if facts.pipelineNote != "" {
		fmt.Printf("  pipelines: %s\n", facts.pipelineNote)
	} else {
		fmt.Printf("  pipelines: %d failing", facts.pipelineRed)
		if len(facts.redNames) > 0 {
			fmt.Printf(": %s", strings.Join(facts.redNames, ", "))
		}
		fmt.Println()
	}
	if facts.fleetNote == "node registry unreachable" {
		fmt.Println("  fleet:     registry unreachable")
	} else {
		fmt.Printf("  fleet:     %d node(s) registered (%s)\n", facts.fleetNodes, facts.fleetNote)
	}
	fmt.Printf("  dbs:       %d stale", facts.staleDBs)
	if len(facts.staleNames) > 0 {
		fmt.Printf(": %s", strings.Join(facts.staleNames, ", "))
	}
	fmt.Println()
	fmt.Printf("  worst:     %s\n", worstItemPointer(facts))

	return nil
}

// worstItemPointer picks the most urgent follow-up command.
func worstItemPointer(facts healthFacts) string {
	if facts.pipelineRed > 0 {
		return fmt.Sprintf("pe triage -> %s (%d failing)", facts.redNames[0], facts.pipelineRed)
	}
	if facts.staleDBs > 0 {
		return "db doctor --repair (stale databases)"
	}
	if facts.repoDirty > 0 {
		return fmt.Sprintf("dirty working trees (%d) — review before pull/push", facts.repoDirty)
	}

	return "nothing urgent — all checks green"
}

// resolveHealthWorkdir mirrors the shared workdir resolution.
func resolveHealthWorkdir(flagVal string) string {
	if strings.TrimSpace(flagVal) != "" {
		return flagVal
	}
	if envVal := strings.TrimSpace(os.Getenv("GITMAP_WORKDIR")); envVal != "" {
		return envVal
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "."
	}

	return filepath.Join(home, "workspace", "repos")
}
