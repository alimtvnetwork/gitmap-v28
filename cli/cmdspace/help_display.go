package cmdspace

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdisplay"
)

// spaceHelpTheme returns the shared theme for space help displays: cyan
// headers, green commands, bright descriptions, dim hints.
func spaceHelpTheme() *helpdisplay.Theme {
	return helpdisplay.NewTheme(nil).
		WithHeaderColor(helpdisplay.NewColor(constants.ColorCyan)).
		WithCommandColor(helpdisplay.NewColor(constants.ColorGreen)).
		WithDescriptionColor(helpdisplay.NewColor(constants.ColorWhite)).
		WithHintColor(helpdisplay.NewColor(constants.ColorDim))
}

// SpaceHelpDisplay builds the HelpDisplay for `gitmap space --help`.
// Subcommands common and backup-branch are CommandHelpers with examples;
// each carries subHelpers so `--help <sub>` renders the sub-display.
func SpaceHelpDisplay() helpdisplay.Displayer {
	return helpdisplay.NewHelpDisplay(
		"Usage: gitmap space <subcommand> [flags]",
		spaceHelpGroups(),
		spaceHelpSuggestions(),
		spaceHelpTheme(),
	)
}

// spaceHelpGroups assembles the top-level space help groups.
func spaceHelpGroups() []helpdisplay.CommandHelpGroup {
	return []helpdisplay.CommandHelpGroup{
		helpdisplay.NewCommandHelpGroup("Subcommands", spaceSubcommandHelpers(), nil),
		helpdisplay.NewCommandHelpGroup("Flags (with common)", spaceCommonFlagHelpers(), nil),
		helpdisplay.NewCommandHelpGroup("Flags (with backup-branch)", spaceBackupBranchFlagHelpers(), nil),
		helpdisplay.NewCommandHelpGroup("Examples", spaceExampleHelpers(), nil),
	}
}

// spaceSubcommandHelpers lists the space subcommands.
func spaceSubcommandHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		spaceCommonHelper(),
		spaceBackupBranchHelper(),
	}
}

// spaceCommonHelper describes `gitmap space common`; its subHelpers mirror
// the `space common --help` sub-display entries.
func spaceCommonHelper() helpdisplay.CommandHelper {
	return helpdisplay.NewCommandHelper(
		"common",
		"Apply the curated common baselines (.gitignore, .gitattributes, .prettierignore, .prettierrc) and run 'git lfs install --local' — same logic as 'gitmap commons'.",
		"gitmap space common --dry-run",
		"",
		concatHelpers(spaceCommonFlagHelpers(), spaceCommonExampleHelpers())...,
	)
}

// spaceBackupBranchHelper describes `gitmap space backup-branch`; its
// subHelpers mirror the `space backup-branch --help` sub-display entries.
func spaceBackupBranchHelper() helpdisplay.CommandHelper {
	return helpdisplay.NewCommandHelper(
		"backup-branch",
		"Create a backup branch (backup/<slug>) from the current HEAD for the given task string, then push it to origin.",
		`gitmap space backup-branch "CLI help displayer overhaul"`,
		"",
		concatHelpers(spaceBackupBranchFlagHelpers(), spaceBackupBranchExampleHelpers())...,
	)
}

// concatHelpers merges helper slices into one subHelpers list.
func concatHelpers(slices ...[]helpdisplay.CommandHelper) []helpdisplay.CommandHelper {
	merged := make([]helpdisplay.CommandHelper, 0)
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}

// spaceCommonFlagHelpers lists the `space common` flags.
func spaceCommonFlagHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--dry-run, -n", "Print planned additions without touching disk", "gitmap space common --dry-run", ""),
		helpdisplay.NewCommandHelper("--force, -f", "Overwrite conflicting JSON values in .prettierrc", "gitmap space common --force", ""),
	}
}

// spaceBackupBranchFlagHelpers lists the `space backup-branch` flags.
func spaceBackupBranchFlagHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--no-push", "Skip pushing the new branch to origin", "", ""),
		helpdisplay.NewCommandHelper("--force", "Recreate the branch at current HEAD if it already exists", "", ""),
	}
}

// spaceExampleHelpers lists the top-level space examples.
func spaceExampleHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("gitmap space common", "Apply the curated baselines in the current repo", "", ""),
		helpdisplay.NewCommandHelper("gitmap space common --dry-run", "Preview planned additions without touching disk", "", ""),
		helpdisplay.NewCommandHelper(`gitmap space backup-branch "CLI help displayer overhaul"`, "Snapshot HEAD into backup/<slug> and push to origin", "", ""),
		helpdisplay.NewCommandHelper(`gitmap space backup-branch "hotfix" --no-push`, "Snapshot HEAD without pushing", "", ""),
	}
}

// spaceHelpSuggestions returns follow-ups shown under the space help.
func spaceHelpSuggestions() []helpdisplay.Suggestion {
	return []helpdisplay.Suggestion{
		helpdisplay.NewSuggestion("Subcommand help: gitmap space --help <subcommand> (try: common, backup-branch)", nil),
	}
}

// SpaceCommonHelpDisplay builds the HelpDisplay for `gitmap space common --help`.
func SpaceCommonHelpDisplay() helpdisplay.Displayer {
	return helpdisplay.NewHelpDisplay(
		"Usage: gitmap space common [flags]",
		[]helpdisplay.CommandHelpGroup{
			helpdisplay.NewCommandHelpGroup("About", spaceCommonAboutHelpers(), nil),
			helpdisplay.NewCommandHelpGroup("Behavior", spaceCommonBehaviorHelpers(), nil),
			helpdisplay.NewCommandHelpGroup("Flags", spaceCommonFlagHelpers(), nil),
			helpdisplay.NewCommandHelpGroup("Examples", spaceCommonExampleHelpers(), nil),
		},
		nil,
		spaceHelpTheme(),
	)
}

// spaceCommonAboutHelpers describes what `space common` does.
func spaceCommonAboutHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("", "Applies the curated common baselines (.gitignore, .gitattributes, .prettierignore, .prettierrc) and runs 'git lfs install --local', all in one pass. Same logic as 'gitmap commons'.", "", ""),
	}
}

// spaceCommonBehaviorHelpers lists the `space common` behavior guarantees.
func spaceCommonBehaviorHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("", "Line-based targets append MISSING lines only; existing entries are preserved verbatim. Safe to re-run.", "", ""),
		helpdisplay.NewCommandHelper("", ".prettierrc is JSON key-union: missing keys added, existing kept unless --force is passed.", "", ""),
		helpdisplay.NewCommandHelper("", "Idempotent — a second run on an unchanged repo writes nothing.", "", ""),
	}
}

// spaceCommonExampleHelpers lists the `space common` examples.
func spaceCommonExampleHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("gitmap space common", "Apply the curated baselines in the current repo", "", ""),
		helpdisplay.NewCommandHelper("gitmap space common --dry-run", "Preview planned additions without touching disk", "", ""),
		helpdisplay.NewCommandHelper("gitmap space common --force", "Overwrite conflicting JSON values in .prettierrc", "", ""),
	}
}

// SpaceBackupBranchHelpDisplay builds the HelpDisplay for `gitmap space backup-branch --help`.
func SpaceBackupBranchHelpDisplay() helpdisplay.Displayer {
	return helpdisplay.NewHelpDisplay(
		`Usage: gitmap space backup-branch "<task string>" [flags]`,
		[]helpdisplay.CommandHelpGroup{
			helpdisplay.NewCommandHelpGroup("About", spaceBackupBranchAboutHelpers(), nil),
			helpdisplay.NewCommandHelpGroup("Behavior", spaceBackupBranchBehaviorHelpers(), nil),
			helpdisplay.NewCommandHelpGroup("Flags", spaceBackupBranchFlagHelpers(), nil),
			helpdisplay.NewCommandHelpGroup("Examples", spaceBackupBranchExampleHelpers(), nil),
		},
		nil,
		spaceHelpTheme(),
	)
}

// spaceBackupBranchAboutHelpers describes what `space backup-branch` does.
func spaceBackupBranchAboutHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("", "Creates a backup branch (backup/<slug>) from current HEAD. The slug is derived from the task string: lowercase, spaces/underscores become hyphens, only [a-z0-9-] kept, repeats collapsed, trimmed, max 60 chars.", "", ""),
	}
}

// spaceBackupBranchBehaviorHelpers lists the `space backup-branch` behavior guarantees.
func spaceBackupBranchBehaviorHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("", "Refuses to run when the working tree has uncommitted changes — backups snapshot HEAD only. There is no override.", "", ""),
		helpdisplay.NewCommandHelper("", "Refuses when backup/<slug> already exists unless --force is passed.", "", ""),
		helpdisplay.NewCommandHelper("", "Pushes backup/<slug> to origin by default; --no-push keeps it local.", "", ""),
		helpdisplay.NewCommandHelper("", "Prints backup/<slug> @ <short-sha> on success.", "", ""),
	}
}

// spaceBackupBranchExampleHelpers lists the `space backup-branch` examples.
func spaceBackupBranchExampleHelpers() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper(`gitmap space backup-branch "CLI help displayer overhaul"`, "Snapshot HEAD and push backup/<slug> to origin", "", ""),
		helpdisplay.NewCommandHelper(`gitmap space backup-branch "pre-release sweep" --no-push`, "Snapshot HEAD without pushing", "", ""),
		helpdisplay.NewCommandHelper(`gitmap space backup-branch "CLI help displayer overhaul" --force`, "Recreate backup/<slug> at current HEAD", "", ""),
	}
}
