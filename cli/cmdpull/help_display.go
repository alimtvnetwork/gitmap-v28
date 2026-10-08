package cmdpull

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdisplay"
)

// PullAllHelpDisplay builds the pull-all help display: the DRY source of
// truth for the pull-all help (spec 243.3 wave C), bound to helpdisplay.Displayer.
// Content is migrated faithfully from the pull-all help topic (helptext/pull-all.md);
// helptext/pull-all.md is generated from this builder.
func PullAllHelpDisplay() helpdisplay.Displayer {
	return helpdisplay.NewHelpDisplay(pullAllHeader(), pullAllGroups(), pullAllSuggestions(), pullAllTheme())
}

// pullAllHeader composes the display title plus the usage block.
func pullAllHeader() string {
	lines := []string{
		"Batch-pull every tracked repository in the catalog (gitmap pull-all)",
		"",
		"  Usage:",
		"    gitmap pull-all [flags]",
		"    gitmap pa [--probe] [--status | --table | --json]",
		"    gitmap pa --probe [-y]",
		"    gitmap pat [flags]",
		"    gitmap pull all table [flags]",
	}
	return strings.Join(lines, "\n")
}

// pullAllGroups assembles the alias, flag, and example sections; the
// forwarded-flags note closes the display as a hint on the last group.
func pullAllGroups() []helpdisplay.CommandHelpGroup {
	return []helpdisplay.CommandHelpGroup{
		helpdisplay.NewCommandHelpGroup("Aliases:", pullAllAliasCommands(), nil),
		helpdisplay.NewCommandHelpGroup("Flags:", pullAllFlagCommands(), nil),
		helpdisplay.NewCommandHelpGroup("Examples:", pullAllExampleCommands(), pullAllHints()),
	}
}

// pullAllAliasCommands lists the pull-all aliases.
func pullAllAliasCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("pa", "Batch pull every tracked repository (fast concise mode)", "", ""),
		helpdisplay.NewCommandHelper("pat", "Batch pull with the full post-pull status table", "", ""),
		helpdisplay.NewCommandHelper("pull-all-table", "Batch pull with the full post-pull status table", "", ""),
	}
}

// pullAllFlagCommands lists the pull-all flags.
func pullAllFlagCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--probe", "Probe GitHub account & workspace for companion repositories (repo-secrets & repo-cache) before pulling", "", ""),
		helpdisplay.NewCommandHelper("--status, --table", "Render the full post-pull repository status table", "", ""),
		helpdisplay.NewCommandHelper("--json", "Suppress interactive progress bar and emit JSON summary to stdout", "", ""),
		helpdisplay.NewCommandHelper("--ssh, -s", "Dispatch pull-all across SSH cluster fleet and Local VM concurrently", "", ""),
		helpdisplay.NewCommandHelper("--verbose", "Enable verbose logging", "", ""),
		helpdisplay.NewCommandHelper("-p, --parallel <N>", "Explicitly set concurrency worker count (default: 0 = auto; overrides presets)", "", ""),
		helpdisplay.NewCommandHelper("--auto-scale", "Dynamically adjust worker concurrency based on CPU core availability (default)", "", ""),
		helpdisplay.NewCommandHelper("--high-perf, --turbo", "Maximum concurrency preset when CPU pressure is low (up to 12 parallel workers)", "", ""),
		helpdisplay.NewCommandHelper("--low-cpu, --conservative", "Throttled concurrency preset when system is under high CPU pressure (1-2 workers)", "", ""),
		helpdisplay.NewCommandHelper("--only-available", "Skip repos whose latest probe reports no new tag", "", ""),
		helpdisplay.NewCommandHelper("--stop-on-fail", "Halt the batch after the first failure", "", ""),
	}
}

// pullAllExampleCommands lists the pull-all examples with variant hints.
func pullAllExampleCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("gitmap pa", "Default fast batch pull (shows only active/updated repos)", "", ""),
		helpdisplay.NewCommandHelper("gitmap pa --status", "Full post-pull status table", "gitmap pat", ""),
		helpdisplay.NewCommandHelper("gitmap pa --json", "Machine-readable JSON summary", "", ""),
		helpdisplay.NewCommandHelper("gitmap pa --ssh", "SSH fleet batch pull (fleet + local VM concurrently)", "", ""),
		helpdisplay.NewCommandHelper("gitmap pull-all --parallel 8 --only-available --stop-on-fail", "8-way parallel, only repos with new commits, stop on first failure", "", ""),
		helpdisplay.NewCommandHelper("gitmap pa --probe", "Probe companion repositories (repo-secrets & repo-cache) before pulling", "gitmap pa --probe -y", ""),
		helpdisplay.NewCommandHelper("gitmap pa --probe --status", "Probe companion repos and render the full status table", "gitmap pa --probe --json", ""),
		helpdisplay.NewCommandHelper("gitmap pa --auto-scale", "CPU auto-scaling and concurrency presets", "gitmap pa --high-perf", ""),
	}
}

// pullAllHints returns closing notes rendered after the example list.
func pullAllHints() []string {
	return []string{
		"All pull flags are forwarded verbatim; --all is injected automatically and is idempotent (passing it again is a no-op).",
	}
}

// pullAllSuggestions returns contextual suggestions for pull-all.
func pullAllSuggestions() []helpdisplay.Suggestion {
	return []helpdisplay.Suggestion{
		helpdisplay.NewSuggestion("Run gitmap scan first to populate the database (see scan help)", nil),
	}
}

// pullAllTheme wires the displayer to the existing help color roles.
func pullAllTheme() *helpdisplay.Theme {
	return helpdisplay.NewTheme(nil).
		WithHeaderColor(helpdisplay.NewColor(constants.ColorYellow)).
		WithCommandColor(helpdisplay.NewColor(constants.ColorGreen)).
		WithDescriptionColor(helpdisplay.NewColor(constants.ColorWhite)).
		WithHintColor(helpdisplay.NewColor(constants.ColorCyan))
}
