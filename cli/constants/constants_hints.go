package constants

// Hint header.
const MsgHintHeader = "\nHints:\n"

// Hint format.
const MsgHintRowFmt = "  → %-35s %s\n"

// Hint messages for project-repos commands (go-repos, node-repos, etc.).
const (
	HintCDRepo     = "gitmap cd <repo-name>"
	HintCDRepoDesc = "Navigate to a repo"

	HintGroupCreate     = "gitmap g create <name>"
	HintGroupCreateDesc = "Create a group"

	HintLsType     = "gitmap ls go"
	HintLsTypeDesc = "List only Go projects"

	HintGroupAdd     = "gitmap g add <group> <slug>"
	HintGroupAddDesc = "Add repos to a group"

	HintPullGroup     = "gitmap g pull"
	HintPullGroupDesc = "Pull repos in active group"

	HintGroupShow     = "gitmap g show <name>"
	HintGroupShowDesc = "Show repos in a group"

	HintGroupDelete     = "gitmap g delete <name>"
	HintGroupDeleteDesc = "Delete a group"

	HintLsGroups     = "gitmap ls groups"
	HintLsGroupsDesc = "List all groups"

	HintGPull     = "gitmap g pull"
	HintGPullDesc = "Pull active group repos"

	HintGStatus     = "gitmap g status"
	HintGStatusDesc = "Show active group status"

	HintGExec     = "gitmap g exec <cmd>"
	HintGExecDesc = "Run git across active group"

	HintGClear     = "gitmap g clear"
	HintGClearDesc = "Clear active group"

	HintCDSetDefault     = "gitmap cd set-default <name> <path>"
	HintCDSetDefaultDesc = "Set a default repo path"

	HintCDRepos     = "gitmap cd repos"
	HintCDReposDesc = "Browse all repos interactively"

	HintMGUsage     = "gitmap mg g1,g2"
	HintMGUsageDesc = "Select multiple groups"

	HintZGCreate     = "gitmap z create <name>"
	HintZGCreateDesc = "Create a zip group"

	HintZGAdd     = "gitmap z add <group> <path>"
	HintZGAddDesc = "Add files to a zip group"

	HintZGShow     = "gitmap z show <name>"
	HintZGShowDesc = "Show items in a zip group"

	HintZGDelete     = "gitmap z delete <name>"
	HintZGDeleteDesc = "Delete a zip group"

	HintZGRelease     = "gitmap r v1.0.0 --zip-group <name>"
	HintZGReleaseDesc = "Include zip group in release"

	HintAliasSet     = "gitmap a set <alias> <slug>"
	HintAliasSetDesc = "Create a repo alias"

	HintAliasList     = "gitmap a list"
	HintAliasListDesc = "List all aliases"

	HintAliasSuggest     = "gitmap a suggest"
	HintAliasSuggestDesc = "Auto-suggest aliases"

	HintAliasUse     = "gitmap pull -A <alias>"
	HintAliasUseDesc = "Use alias with any command"

	HintAliasRemove     = "gitmap a remove <alias>"
	HintAliasRemoveDesc = "Remove an alias"
)
