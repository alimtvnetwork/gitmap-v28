package constants

// Cross-command CLI flag values.
//
// Each command parses its own `flag.NewFlagSet`, so these strings are
// reused safely (e.g. --json appears for `list-versions`, `list-releases`,
// `amend-list`, etc.). Centralized here so renames stay consistent.
const (
	FlagOpenValue   = "--open"
	FlagJSON        = "--json"
	FlagLimit       = "--limit"
	FlagSource      = "--source"
	FlagAllRepos    = "--all-repos" // v3.20.0: list-releases multi-repo batch view
	FlagCompact     = "--compact"
	FlagGroups      = "--groups"
	FlagFilter      = "--filter" // v5.42.0: `gitmap help --filter <q>` search
	FlagFilterShort = "-f"
)
