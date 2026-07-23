// Package constants — Chrome profile copy/export/import command IDs,
// help text, messages, and exit codes for `gitmap chrome-profile-copy`
// and friends. Spec: spec/04-generic-cli/40-chrome-profile-copy.md.
package constants

// gitmap:cmd top-level
// Chrome profile command IDs and short aliases.
const (
	CmdChromeProfileCopy        = "chrome-profile-copy"
	CmdChromeProfileCopyAlias   = "cpc"
	CmdChromeProfileExport      = "chrome-profile-export"
	CmdChromeProfileExportAlias = "cpe"
	CmdChromeProfileImport      = "chrome-profile-import"
	CmdChromeProfileImportAlias = "cpi"
	CmdChromeProfileList        = "chrome-profile-list"
	CmdChromeProfileListAlias   = "cpl"
	CmdChromeProfileListAlias2  = "chrome-profiles"
	CmdChromeProfileDelete      = "chrome-profile-delete"
	CmdChromeProfileDeleteAlias = "cpd"
	CmdChromeProfileMerge       = "chrome-profile-merge"
	CmdChromeProfileMergeAlias  = "cpm"
)

// gitmap:cmd top-level
// Chrome umbrella command (v6.69.0) + subcommands.
const (
	CmdChrome                 = "chrome"
	SubCmdChromeBackup        = "backup"           // gitmap:cmd skip
	SubCmdChromeRestore       = "restore"          // gitmap:cmd skip
	SubCmdChromeDiff          = "diff"             // gitmap:cmd skip
	SubCmdChromeExportBookmrk = "export-bookmarks" // gitmap:cmd skip
	SubCmdChromeWhich         = "which"            // gitmap:cmd skip
)

// Chrome profile merge knobs.
const (
	ChromeMergeWhatAll        = "all"
	ChromeMergeWhatSettings   = "settings"
	ChromeMergeWhatBookmarks  = "bookmarks"
	ChromeMergeWhatExtensions = "extensions"
	ChromePreferencesFile     = "Preferences"
	ChromeSecurePrefsFile     = "Secure Preferences"
	ChromeBookmarksFile       = "Bookmarks"
)

// Chrome merge user-facing messages.
const (
	MsgChromeMergeStart    = "\n\033[1;96m▸ chrome-profile-merge\033[0m  \033[1m%s\033[0m → \033[1m%s\033[0m  \033[2;37m(what=%s)\033[0m\n  \033[2;37msource     \033[0m %s\n  \033[2;37mdestination\033[0m %s\n"
	MsgChromeMergeStepHdr  = "\n\033[1;94m• %s\033[0m\n"
	MsgChromeMergeSummary  = "\n\033[1;92m✓ merge complete\033[0m  added=\033[1m%d\033[0m skipped=\033[1m%d\033[0m overwrote=\033[1m%d\033[0m\n"
	MsgChromeMergePrompt   = "  conflict: %s\n    [k]eep destination, [o]verwrite with source, [a]ll-keep, [A]ll-overwrite, [q]uit: "
	MsgChromeMergeDryRun   = "  \033[2;37m(dry-run; no changes written)\033[0m\n"
	MsgChromeMergeDryAdd   = "  \033[1;92m+ add\033[0m       %s\n"
	MsgChromeMergeDryOver  = "  \033[1;93m~ overwrite\033[0m %s\n"
	MsgChromeMergeDryKeep  = "  \033[2;37m= keep\033[0m      %s\n"
	ErrChromeMergeUsage    = "chrome-profile-merge: ERROR <src> and <dst> are required\n  usage: gitmap chrome-profile-merge <src> <dst> [--what all|settings|bookmarks|extensions] [--yes|--force] [--dry-run]\n"
	ErrChromeMergeUnknown  = "chrome-profile-merge: ERROR --what=%q unknown (use: all|settings|bookmarks|extensions)\n"
	HelpChromeProfileMerge = "  chrome-profile-merge (cpm) <src> <dst> [--what all|settings|bookmarks|extensions] [--yes|--force] Merge selected pieces of one profile into another"
)

// Chrome profile help-line entries surfaced by `gitmap help`.
const (
	HelpChromeProfileCopy   = "  chrome-profile-copy (cpc) <src> <dst>   Copy a Chrome profile (bookmarks, extensions, prefs, flags) into an offline profile"
	HelpChromeProfileExport = "  chrome-profile-export (cpe) <name> [out] Export profile to JSON (default: .gitmap/chrome/<name>.json)"
	HelpChromeProfileImport = "  chrome-profile-import (cpi) <file> [name] Import a Chrome profile from a JSON export"
	HelpChromeProfileList   = "  chrome-profile-list (cpl)               List Chrome profiles known to gitmap"
)

// Chrome profile messages and errors.
const (
	MsgChromeProfileCopyStart   = "\n\033[1;96m▸ chrome-profile-copy\033[0m  \033[1m%s\033[0m → \033[1m%s\033[0m\n  \033[2;37msource     \033[0m %s\n  \033[2;37mdestination\033[0m %s\n"
	MsgChromeProfileCopyDone    = "\n\033[1;92m✓ copy complete\033[0m  \033[1m%d\033[0m files in \033[1m%s\033[0m\n"
	MsgChromeProfileLockSummary = "\033[1;93m⚠ skipped %d volatile Chrome lock file(s)\033[0m  \033[2;37m(held by Chrome/extension; safe to ignore)\033[0m\n"
	MsgChromeProfileNextSteps   = "\n\033[1;94mNext steps\033[0m\n  \033[2;37mundo  \033[0m \033[1;96mgitmap chrome-profile-delete %s --yes\033[0m\n  \033[2;37mredo  \033[0m \033[1;96mgitmap chrome-profile-copy %s %s\033[0m\n  \033[2;37mverify\033[0m \033[1;96mgitmap chrome-profile-list\033[0m\n"
	MsgChromeProfileExportOk    = "chrome-profile-export: wrote %s (%d bytes)\n"
	MsgChromeProfileExportCSV   = "chrome-profile-export: csv  %s (%d bytes)\n"
	MsgChromeProfileArtifactsHd = "\n\033[1;94mArtifacts\033[0m\n"
	MsgChromeProfileArtifactRow = "  \033[2;37m%-5s\033[0m \033[1;96m%s\033[0m\n"
	MsgChromeProfileArtifactNA  = "(skipped)"
	MsgChromeProfileDBSynced    = "\033[1;92m✓ db synced\033[0m  profile \033[1m%s\033[0m\n"
	MsgChromeProfileDBWarn      = "  \033[1;93m⚠\033[0m chrome-profile: db sync failed: %v\n"
	MsgChromeProfileImportOk    = "chrome-profile-import: imported %s into profile %q\n"
	MsgChromeProfileImportCSV   = "chrome-profile-import: csv source detected — restoring extension IDs + known preferences (bookmarks omitted)\n"
	ErrChromeProfileChromeOpen  = "chrome-profile-copy: ERROR Chrome is still running.\n  Close every Chrome window and tray/background Chrome process, then rerun:\n  gitmap chrome-profile-copy %s %s\n"
	MsgChromeProfileListEmpty   = "chrome-profile-list: no profiles found at %s\n"
	MsgChromeProfileListHdr     = "Chrome profiles (%s):\n"
	MsgChromeProfileListDBHdr   = "Tracked in gitmap DB:\n"
	MsgChromeProfileListDBRow   = "  - %-30s  exports=%d  last=%s\n"
	MsgChromeProfileSkipChrome  = "  \033[2;37mhint: close Chrome before copying — open sessions may corrupt the destination.\033[0m\n"
	WarnChromeProfileSkipLock   = "  \033[2;37m· skipped volatile Chrome lock file: %s\033[0m\n"
	WarnChromeProfileCheckOpen  = "  \033[1;93m⚠\033[0m chrome-profile-copy: could not check whether Chrome is running: %v\n"
	MsgChromeProfileDeleteOk    = "chrome-profile-delete: removed profile %q (%d artifacts)\n"
	MsgChromeProfileDeleteRm    = "  rm %s\n"
	MsgChromeProfileDeleteAbort = "chrome-profile-delete: aborted — re-run with --yes to confirm\n"

	ErrChromeProfileUsageCopy   = "chrome-profile-copy: ERROR <src> and <dst> are required\n  usage: gitmap chrome-profile-copy <src-profile> <dst-profile>\n"
	ErrChromeProfileUsageExport = "chrome-profile-export: ERROR <name> is required\n  usage: gitmap chrome-profile-export <name> [out.json]\n"
	ErrChromeProfileUsageImport = "chrome-profile-import: ERROR <file> is required\n  usage: gitmap chrome-profile-import <file.json|file.csv> [dst-profile]\n"
	ErrChromeProfileUsageDelete = "chrome-profile-delete: ERROR <name> is required\n  usage: gitmap chrome-profile-delete <name> [--yes]\n"
	ErrChromeProfileSrcMissing  = "chrome-profile-copy: ERROR source profile %q not found at %s\n"
	ErrChromeProfileCopyFailed  = "chrome-profile-copy: ERROR copy failed\n  source profile: %s\n  destination profile: %s\n  source path: %s\n  destination path: %s\n  failed entry: %s\n  operation: %s\n  cause: %v\n  hint: close Chrome completely, then retry. If it still fails, check the listed file permissions.\n"
	ErrChromeProfileExportFail  = "chrome-profile-export: ERROR %v\n"
	ErrChromeProfileImportFail  = "chrome-profile-import: ERROR %v\n"
	ErrChromeProfileDeleteFail  = "chrome-profile-delete: ERROR %v\n"
	ErrChromeProfileNotInDB     = "chrome-profile-delete: ERROR profile %q not found in gitmap DB\n"
)

// Chrome profile copy operation labels.
const (
	ChromeProfileCopyOpMkdir  = "create destination directory"
	ChromeProfileCopyOpStat   = "inspect source path"
	ChromeProfileCopyOpRead   = "read source file"
	ChromeProfileCopyOpWrite  = "write destination file"
	ChromeProfileCopyOpList   = "list source directory"
	ChromeProfileCopyOpCopy   = "copy profile"
	ChromeProfileCopyUnknown  = "(unknown)"
	ChromeProfileLockFileName = "LOCK"
	ChromeProfileLockReason   = "runtime-only Chrome lock file; Chrome recreates it"
	ChromeLocalStateFile      = "Local State"
	ChromeLocalStateTmpSuffix = ".gitmap.tmp"
	ChromeLocalStateBakSuffix = ".gitmap.bak"
	ChromeDefaultProfileDir   = "Default"
	ChromeProfileDirPrefix    = "Profile "
)

// Chrome process detection constants.
const (
	ChromeProcessTasklist     = "tasklist"
	ChromeProcessTasklistFlag = "/FI"
	ChromeProcessTasklistExpr = "IMAGENAME eq chrome.exe"
	ChromeProcessTasklistNH   = "/NH"
	ChromeProcessWindowsImage = "chrome.exe"
	ChromeProcessPgrep        = "pgrep"
	ChromeProcessPgrepExact   = "-x"
	ChromeProcessMacName      = "Google Chrome"
)

// ChromeProcessLinuxNames enumerates common Chromium process names.
var ChromeProcessLinuxNames = []string{"chrome", "google-chrome", "google-chrome-stable", "chromium", "chromium-browser"}

// Local State registration messages.
const (
	MsgChromeProfileRegistered = "\033[1;92m✓ registered\033[0m  \033[1m%s\033[0m in Chrome's profile picker (Local State)\n"
	MsgChromeProfileRegOnly    = "\n\033[1;96m▸ chrome-profile-copy --register-only\033[0m  \033[1m%s\033[0m\n  \033[2;37m(skipping file copy; refreshing Chrome Local State entry only)\033[0m\n"
	WarnChromeProfileRegister  = "  \033[1;93m⚠\033[0m chrome-profile-copy: could not register %q in Chrome's Local State: %v\n  \033[2;37m(profile files were copied — restart Chrome and add the profile manually if it does not appear)\033[0m\n"
	WarnChromeProfileBakRm     = "  \033[1;93m⚠\033[0m chrome-profile-copy: could not remove Local State backup %s: %v\n"
	HelpChromeProfileDelete    = "  chrome-profile-delete (cpd) <name> [--yes] Remove a profile + its stored artifacts from the gitmap DB"
)

// Chrome User Data subpaths copied by cpc. Excluded by design:
// Cookies, Login Data, History, Cache, GPUCache, sync tokens.
var ChromeProfileCopyEntries = []string{
	"Bookmarks",
	"Favicons",
	"Preferences",
	"Secure Preferences",
	"Extensions",
	"Local Extension Settings",
	"Extension Rules",
	"Extension State",
	"Sync Extension Settings",
	"Web Data",
	"Shortcuts",
	"TransportSecurity",
}

// Chrome User Data top-level files (siblings of profile dirs).
var ChromeUserDataTopLevel = []string{
	"Local State",
}

// Chrome profile exit codes.
const (
	ExitChromeProfileOk         = 0
	ExitChromeProfileUsage      = 6
	ExitChromeProfileNotFound   = 7
	ExitChromeProfileCopyFailed = 10
)
