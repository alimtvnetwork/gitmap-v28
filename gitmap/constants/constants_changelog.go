package constants

// gitmap:cmd top-level
// Changelog generate command.
const (
	CmdChangelogGen      = "changelog-generate"
	CmdChangelogGenAlias = "cg"
)

// Changelog generate help text.
const (
	HelpChangelogGen = "  changelog-gen (cg)  Auto-generate changelog from commits between tags (--from, --to, --write)"
)

// Changelog generate flag descriptions.
const (
	FlagDescFrom  = "Start tag (older). Defaults to second-latest tag"
	FlagDescTo    = "End tag or HEAD. Defaults to latest tag"
	FlagDescWrite = "Prepend output to CHANGELOG.md instead of printing"
)

// Changelog generate git arguments.
const (
	ChangelogGenFormat      = "--format=%s"
	ChangelogGenNoMerges    = "--no-merges"
	ChangelogGenSortFlag    = "--sort"
	ChangelogGenSortVersion = "-version:refname"
)

// Changelog generate messages.
const (
	MsgChangelogGenHeader  = "\n  Changelog: %s → %s\n\n"
	MsgChangelogGenEmpty   = "  No commits found between %s and %s.\n"
	MsgChangelogGenWritten = "  ✓ Prepended changelog to %s\n"
	MsgChangelogGenPreview = "  Preview (use --write to save):\n\n"
)

// Changelog generate errors.
const (
	ErrChangelogGenCommits     = "failed to list commits between %s and %s: %v"
	ErrChangelogGenTags        = "failed to list tags: %v"
	ErrChangelogGenTagNotFound = "tag %s not found locally"
	ErrChangelogGenNoTags      = "no version tags found — create a release first"
	ErrChangelogGenWrite       = "failed to write changelog at %s: %v (operation: write)"
	ErrChangelogGenRead        = "failed to read existing changelog at %s: %v (operation: read)"
)
