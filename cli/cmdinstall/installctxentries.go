package cmdinstall

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/enums/ctxmodetype"
)

// ctxEntry describes a single right-click menu item (or category).
// A category is an entry with empty Args and non-nil Children.
type ctxEntry struct {
	KeyName string // numeric prefix preserves menu order, e.g. "10_release_next"
	MUIVerb string // visible label
	Args    []string
	Mode    ctxmodetype.VariantType
	Exe     string // override executable; empty => use the gitmap binary
	Icon    string // Windows: per-entry Icon registry value. Empty => no Icon written.
	//                  Supports the constants.CtxIconExeToken ("{exe}") placeholder,
	//                  which is substituted with the resolved gitmap binary path at
	//                  registry-write time (see leafCommands / categoryCommands).
	Extended bool       // Windows: Shift+right-click only (HKCU "Extended" REG_SZ); macOS/Linux: prepend confirm prompt
	Children []ctxEntry // non-nil => this is a submenu
}

// ctxMenu returns the full menu tree for directory context.
func ctxMenu() []ctxEntry {
	entries := ctxDirectoryEntries()
	entries = append(entries, ctxCategoryEntries()...)
	entries = append(entries, ctxTailEntries()...)

	return entries
}

// ctxBackgroundMenu returns the menu tree for right-clicking in folder background.
// Omits Show log which requires an active repository target.
func ctxBackgroundMenu() []ctxEntry {
	entries := []ctxEntry{
		{KeyName: "01_scan_here", MUIVerb: "Scan here", Args: []string{constants.CmdScan}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "02_repos_ui", MUIVerb: "Open Repositories Web UI", Args: []string{"ui"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "03_installer_ui", MUIVerb: "Open Installer Manager UI", Args: []string{"ui", "installer"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "04_settings_ui", MUIVerb: "GitMap Settings (UI)", Args: []string{"ui", "settings"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
	}
	entries = append(entries, ctxCategoryEntries()...)
	entries = append(entries, ctxTailEntries()...)

	return entries
}

func ctxDirectoryEntries() []ctxEntry {
	return []ctxEntry{
		{KeyName: "01_show_log", MUIVerb: "Show log", Args: []string{"show-log", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "02_diff", MUIVerb: "Diff (Working Tree)", Args: []string{"diff", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "03_branch_compare", MUIVerb: "Branch comparison", Args: []string{"show-log", "--branch-compare", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "04_commit", MUIVerb: "Commit to main", Args: []string{"commit", "-b", "main", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "05_pull", MUIVerb: "Pull", Args: []string{"pull", `"%1"`}, Mode: ctxmodetype.Terminal},
		{KeyName: "06_push", MUIVerb: "Push", Args: []string{"push", `"%1"`}, Mode: ctxmodetype.Terminal},
	}
}

func ctxCategoryEntries() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_scan", MUIVerb: "Scan", Children: scanChildren()},
		{KeyName: "20_clone", MUIVerb: "Clone", Children: cloneChildren()},
		{KeyName: "30_release", MUIVerb: "Release", Children: releaseChildren()},
		{KeyName: "40_repos", MUIVerb: "Repositories", Children: repoChildren()},
		{KeyName: "50_visibility", MUIVerb: "Visibility", Children: visibilityChildren()},
		{KeyName: "60_install", MUIVerb: "Install", Children: installChildren()},
		{KeyName: "70_settings", MUIVerb: "Settings (UI)", Args: []string{"ui", "settings"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
	}
}

func ctxTailEntries() []ctxEntry {
	return []ctxEntry{
		{KeyName: "80_antigravity", MUIVerb: "Open project with Antigravity", Args: []string{constants.CmdAg}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "90_terminal", MUIVerb: constants.MsgCtxOpenTerminalLbl, Mode: ctxmodetype.Prefill},
		{KeyName: "91_docs", MUIVerb: constants.MsgCtxDocsLbl, Args: []string{constants.CmdDocs}, Mode: ctxmodetype.Terminal},
		{KeyName: "92_help", MUIVerb: "Help (filter…)", Args: []string{constants.CmdHelp}, Mode: ctxmodetype.Prefill, Icon: constants.CtxIconGitmap},
	}
}

// ctxFileMenu returns the context menu entries registered under
// HKCU\Software\Classes\*\shell\gitmap for right-clicking any file.
func ctxFileMenu() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_diff", MUIVerb: "Diff", Args: []string{"diff", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "20_diff_later", MUIVerb: "Diff later", Args: []string{"diff", "--stash-path", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "30_diff_prev", MUIVerb: "Diff with previous version", Args: []string{"diff", "HEAD~1", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "40_show_log", MUIVerb: "Show log", Args: []string{"show-log", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "50_blame", MUIVerb: "Blame", Args: []string{"blame", `"%1"`}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "60_commit", MUIVerb: "Commit to main", Args: []string{"commit", "-b", "main"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
	}
}

func scanChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_scan_here", MUIVerb: "Scan here", Args: []string{constants.CmdScan}, Mode: ctxmodetype.Terminal},
		{KeyName: "20_rescan", MUIVerb: "Rescan", Args: []string{constants.CmdRescan}, Mode: ctxmodetype.Terminal},
		{KeyName: "30_find_next", MUIVerb: "Find next", Args: []string{constants.CmdFindNext}, Mode: ctxmodetype.Terminal},
	}
}

func cloneChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_clone_next", MUIVerb: "Clone-next here", Args: []string{constants.CmdCloneNext}, Mode: ctxmodetype.Terminal},
		{KeyName: "20_pull", MUIVerb: "Pull", Args: []string{constants.CmdPull}, Mode: ctxmodetype.Terminal},
		{KeyName: "30_pull_all", MUIVerb: "Pull all (every tracked repo)", Args: []string{constants.CmdPullAll}, Mode: ctxmodetype.Terminal, Extended: true},
	}
}

func releaseChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_release", MUIVerb: "Release current", Args: []string{constants.CmdRelease}, Mode: ctxmodetype.Terminal},
		{KeyName: "20_release_next", MUIVerb: "Release next (bump minor)", Args: []string{constants.CmdRelease, constants.FlagBumpDash, constants.BumpMinor}, Mode: ctxmodetype.Terminal},
		{KeyName: "30_release_pull", MUIVerb: "Pull then release", Args: []string{constants.CmdReleasePull}, Mode: ctxmodetype.Terminal},
		{KeyName: "40_release_pending", MUIVerb: "Release pending", Args: []string{constants.CmdReleasePending}, Mode: ctxmodetype.Terminal},
		{KeyName: "50_list_releases", MUIVerb: "List releases", Args: []string{constants.CmdListReleases}, Mode: ctxmodetype.Terminal},
		{KeyName: "60_list_versions", MUIVerb: "List versions", Args: []string{constants.CmdListVersions}, Mode: ctxmodetype.Terminal},
	}
}

func visibilityChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_public", MUIVerb: "Make public", Args: []string{constants.CmdMakePublic}, Mode: ctxmodetype.Terminal},
		{KeyName: "20_private", MUIVerb: "Make private", Args: []string{constants.CmdMakePrivate}, Mode: ctxmodetype.Terminal},
	}
}

func repoChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "10_web_ui", MUIVerb: "Open Repositories Web UI", Args: []string{"ui"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "20_list", MUIVerb: "List Tracked Repositories", Args: []string{"list"}, Mode: ctxmodetype.Terminal},
		{KeyName: "30_status", MUIVerb: "Status All Repositories", Args: []string{"status-all"}, Mode: ctxmodetype.Terminal},
		{KeyName: "40_pull_all", MUIVerb: "Pull All Repositories", Args: []string{"pull-all"}, Mode: ctxmodetype.Terminal},
	}
}

func installChildren() []ctxEntry {
	return []ctxEntry{
		{KeyName: "01_installer_ui", MUIVerb: "Open Installer Manager UI", Args: []string{"ui", "installer"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "05_install_ctx", MUIVerb: "Install Context Menu", Args: []string{"install", "context-menu"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "10_wizard", MUIVerb: "Install Wizard (Interactive)", Args: []string{"install", "--wizard"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
		{KeyName: "20_dev", MUIVerb: "Install Dev Workstation Profile", Args: []string{"install", "profile", "dev"}, Mode: ctxmodetype.Terminal},
		{KeyName: "30_tools", MUIVerb: "Install Development Tools (ls)", Args: []string{"install", "list"}, Mode: ctxmodetype.Terminal},
		{KeyName: "40_agy", MUIVerb: "Install Antigravity (AGY)", Args: []string{"install", "antigravity"}, Mode: ctxmodetype.Terminal},
		{KeyName: "50_clean_pwsh", MUIVerb: "Clean Duplicate PowerShell", Args: []string{"context-menu", "clean-powershell"}, Mode: ctxmodetype.Terminal},
	}
}
