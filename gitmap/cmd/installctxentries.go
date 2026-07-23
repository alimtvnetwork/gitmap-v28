package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

// ctxEntry describes a single right-click menu item (or category).
// A category is an entry with empty Args and non-nil Children.
type ctxEntry struct {
	KeyName string // numeric prefix preserves menu order, e.g. "10_release_next"
	MUIVerb string // visible label
	Args    []string
	Mode    constants.CtxMode
	Exe     string // override executable; empty => use the gitmap binary
	Icon    string // Windows: per-entry Icon registry value. Empty => no Icon written.
	//                  Supports the constants.CtxIconExeToken ("{exe}") placeholder,
	//                  which is substituted with the resolved gitmap binary path at
	//                  registry-write time (see leafCommands / categoryCommands).
	Extended bool       // Windows: Shift+right-click only (HKCU "Extended" REG_SZ); macOS/Linux: prepend confirm prompt
	Children []ctxEntry // non-nil => this is a submenu
}

// ctxMenu returns the full menu tree. Single source of truth for both
// install and uninstall — uninstall walks the same tree and deletes
// each generated key.
func ctxMenu() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_scan", MUIVerb: "Scan", Children: scanChildren()},
		{KeyName: "20_clone", MUIVerb: "Clone", Children: cloneChildren()},
		{KeyName: "30_release", MUIVerb: "Release", Children: releaseChildren()},
		{KeyName: "40_repos", MUIVerb: "Repos", Children: repoChildren()},
		{KeyName: "50_visibility", MUIVerb: "Visibility", Children: visibilityChildren()},
		{KeyName: "60_tools", MUIVerb: "Tools", Children: toolsChildren()},
		{KeyName: "70_git", MUIVerb: "Git", Children: gitChildren()},
		{KeyName: "90_terminal", MUIVerb: constants.MsgCtxOpenTerminalLbl, Mode: constants.CtxModePrefill},
		{KeyName: "91_docs", MUIVerb: constants.MsgCtxDocsLbl, Args: []string{constants.CmdDocs}, Mode: constants.CtxModeSilent},
		{KeyName: "90_terminal", MUIVerb: constants.MsgCtxOpenTerminalLbl, Mode: constants.CtxModePrefill},
		{KeyName: "91_docs", MUIVerb: constants.MsgCtxDocsLbl, Args: []string{constants.CmdDocs}, Mode: constants.CtxModeSilent},
		{KeyName: "92_help", MUIVerb: "Help (filter…)", Args: []string{constants.CmdHelp}, Mode: constants.CtxModePrefill, Icon: constants.CtxIconGitmap},
	}
}

func scanChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_scan_here", MUIVerb: "Scan here", Args: []string{constants.CmdScan}, Mode: constants.CtxModeTerminal},
		{KeyName: "20_rescan", MUIVerb: "Rescan", Args: []string{constants.CmdRescan}, Mode: constants.CtxModeTerminal},
		{KeyName: "30_find_next", MUIVerb: "Find next", Args: []string{constants.CmdFindNext}, Mode: constants.CtxModeSilent},
	}
}

func cloneChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_clone_next", MUIVerb: "Clone-next here", Args: []string{constants.CmdCloneNext}, Mode: constants.CtxModeTerminal},
		{KeyName: "20_pull", MUIVerb: "Pull", Args: []string{constants.CmdPull}, Mode: constants.CtxModeTerminal},
		// Pull-all is a multi-repo batch op. Hidden behind Shift+right-click on Windows
		// (Extended verb); on macOS/Linux it surfaces with a "(all repos)" label so the
		// fan-out is obvious. Implemented as runPullAll => `gitmap pull --all`.
		{KeyName: "30_pull_all", MUIVerb: "Pull all (every tracked repo)", Args: []string{constants.CmdPullAll}, Mode: constants.CtxModeTerminal, Extended: true},
	}
}

func releaseChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_release", MUIVerb: "Release current", Args: []string{constants.CmdRelease}, Mode: constants.CtxModeTerminal},
		{KeyName: "20_release_next", MUIVerb: "Release next (bump minor)", Args: []string{constants.CmdRelease, constants.FlagBumpDash, constants.BumpMinor}, Mode: constants.CtxModeTerminal},
		{KeyName: "30_release_pull", MUIVerb: "Pull then release", Args: []string{constants.CmdReleasePull}, Mode: constants.CtxModeTerminal},
		{KeyName: "40_release_pending", MUIVerb: "Release pending", Args: []string{constants.CmdReleasePending}, Mode: constants.CtxModeSilent},
		{KeyName: "50_list_releases", MUIVerb: "List releases", Args: []string{constants.CmdListReleases}, Mode: constants.CtxModeSilent},
		{KeyName: "60_list_versions", MUIVerb: "List versions", Args: []string{constants.CmdListVersions}, Mode: constants.CtxModeSilent},
	}
}

func repoChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_go", MUIVerb: "Go projects", Args: []string{constants.CmdGoRepos}, Mode: constants.CtxModeSilent},
		{KeyName: "20_node", MUIVerb: "Node projects", Args: []string{constants.CmdNodeRepos}, Mode: constants.CtxModeSilent},
		{KeyName: "30_react", MUIVerb: "React projects", Args: []string{constants.CmdReactRepos}, Mode: constants.CtxModeSilent},
		{KeyName: "40_cpp", MUIVerb: "C++ projects", Args: []string{constants.CmdCppRepos}, Mode: constants.CtxModeSilent},
		{KeyName: "50_csharp", MUIVerb: "C# projects", Args: []string{constants.CmdCsharpRepos}, Mode: constants.CtxModeSilent},
	}
}

func visibilityChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_public", MUIVerb: "Make public", Args: []string{constants.CmdMakePublic}, Mode: constants.CtxModeTerminal},
		{KeyName: "20_private", MUIVerb: "Make private", Args: []string{constants.CmdMakePrivate}, Mode: constants.CtxModeTerminal},
	}
}

func toolsChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_fix_repo", MUIVerb: "Fix repo", Args: []string{constants.CmdFixRepo}, Mode: constants.CtxModeTerminal},
		{KeyName: "20_diff", MUIVerb: "Diff", Args: []string{constants.CmdDiff}, Mode: constants.CtxModeTerminal},
		{KeyName: "30_history", MUIVerb: "History", Args: []string{constants.CmdHistory}, Mode: constants.CtxModeTerminal},
		{KeyName: "40_update", MUIVerb: "Update", Args: []string{constants.CmdUpdate}, Mode: constants.CtxModeTerminal},
	}
}

// gitChildren returns the raw-git submenu. These entries shell out to
// `git` directly (Exe override) so users can inspect repository state
// in a terminal viewer without launching gitmap. All are Terminal mode
// because the output is multi-line and worth reading.
func gitChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_history", MUIVerb: constants.CtxGitHistoryLabel, Args: constants.CtxGitHistoryArgs, Exe: constants.CtxExeGit, Mode: constants.CtxModeTerminal},
		{KeyName: "20_diff", MUIVerb: constants.CtxGitDiffLabel, Args: constants.CtxGitDiffArgs, Exe: constants.CtxExeGit, Mode: constants.CtxModeTerminal},
		{KeyName: "30_log", MUIVerb: constants.CtxGitLogLabel, Args: constants.CtxGitLogArgs, Exe: constants.CtxExeGit, Mode: constants.CtxModeTerminal},
		{KeyName: "40_status", MUIVerb: constants.CtxGitStatusLabel, Args: constants.CtxGitStatusArgs, Exe: constants.CtxExeGit, Mode: constants.CtxModeSilent},
	}
}
