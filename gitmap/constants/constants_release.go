package constants

// Setup section headers.
const (
	SetupSectionDiff  = "Diff Tool"
	SetupSectionMerge = "Merge Tool"
	SetupSectionAlias = "Aliases"
	SetupSectionCred  = "Credential Helper"
	SetupSectionCore  = "Core Settings"
	SetupSectionComp  = "■ Shell Completion —"
	SetupGlobalFlag   = "--global"
)

// Release messages.
const (
	MsgReleaseStart         = "\n  " + ColorMagenta + "🚀 Creating release " + ColorWhite + "%s" + ColorReset + ColorMagenta + "..." + ColorReset + "\n"
	MsgReleaseBranch        = "  " + ColorGreen + "✓ Created branch" + ColorReset + " " + ColorCyan + "%s" + ColorReset + "\n"
	MsgReleaseTag           = "  " + ColorGreen + "🏷  Created tag" + ColorReset + " " + ColorMagenta + "%s" + ColorReset + "\n"
	MsgReleasePushed        = "  " + ColorGreen + "📤 Pushed" + ColorReset + " branch and tag to " + ColorCyan + "origin" + ColorReset + "\n"
	MsgReleaseMeta          = "  " + ColorGreen + "📝 Release metadata written to" + ColorReset + " " + ColorCyan + "%s" + ColorReset + "\n"
	MsgReleaseMetaCommitted = "  " + ColorGreen + "✓ Committed release metadata on" + ColorReset + " " + ColorCyan + "%s" + ColorReset + "\n"
	MsgReleaseLatest        = "  " + ColorGreen + "⭐ Marked" + ColorReset + " " + ColorMagenta + "%s" + ColorReset + " " + ColorGreen + "as latest release" + ColorReset + "\n"
	MsgReleaseAttach        = "  " + ColorGreen + "📎 Attached" + ColorReset + " " + ColorCyan + "%s" + ColorReset + "\n"
	MsgReleaseChangelog     = "  " + ColorGreen + "📰 Using CHANGELOG.md as release body" + ColorReset + "\n"
	MsgReleaseReadme        = "  " + ColorGreen + "📖 Attached README.md" + ColorReset + "\n"
	MsgReleaseDryRun        = "  " + ColorDim + "[dry-run] %s" + ColorReset + "\n"
	MsgReleaseComplete      = "\n  " + ColorMagenta + "🎉 ── Release " + ColorWhite + "%s" + ColorMagenta + " complete ── 🎉" + ColorReset + "\n"
	MsgReleaseBranchStart   = "\n  " + ColorCyan + "Completing release from " + ColorReset + "%s...\n"
	MsgReleaseBranchPending = "\n  " + ColorYellow + "→ On release branch %s with no tag — completing pending release..." + ColorReset + "\n"
	MsgReleaseVersionRead   = "  → Version from %s: " + ColorCyan + "%s" + ColorReset + "\n"
	MsgReleaseBumpResult    = "  → Bumped %s → " + ColorCyan + "%s" + ColorReset + "\n"
	MsgReleaseNotes         = "  → Release notes: %s\n"
	MsgReleaseSwitchedBack  = "  " + ColorGreen + "✓ Switched back to" + ColorReset + " %s\n"
	MsgReleasePendingNone   = "  No pending release branches found."
	MsgReleasePendingFound  = "\n  Found %d pending release branch(es).\n"
	MsgReleasePendingFailed = "  " + ColorRed + "✗ Failed to release %s: %v" + ColorReset + "\n"
	ReleaseBranchPrefix     = "release/"
	ChangelogFile           = "CHANGELOG.md"
	ReadmeFile              = "README.md"
	ReleaseTagPrefix        = "Release "
	FlagDescNotes           = "Release notes or title for the release"
)

// Bare-release auto-bump messages (v3.19.0).
//
// When `gitmap release` / `gitmap r` is run with no version and no --bump,
// gitmap reads the last release from .gitmap/release/latest.json, bumps the
// MINOR segment, and prompts the user. -y skips the prompt.
const (
	MsgReleaseAutoBumpHeader  = "\n  Auto-bump: %s → %s (minor)\n"
	MsgReleaseAutoBumpPrompt  = "  Proceed with this release? [y/N]: "
	MsgReleaseAutoBumpYes     = "  → -y supplied; proceeding without prompt.\n"
	MsgReleaseAutoBumpAborted = "  ✗ Auto-bump aborted by user.\n"
)

// Multi-repo scan-dir release messages (v3.19.0).
//
// When `gitmap r` is run from a directory containing many git repos (the
// cwd itself is NOT a git repo), gitmap walks the tree, keeps only repos
// that have a prior release manifest, computes a minor bump per repo, and
// prompts ONCE before releasing them all.
const (
	MsgReleaseScanHeader  = "\n  Auto-bump %d repo(s) with prior releases:\n"
	MsgReleaseScanRow     = "    • %s   %s → %s\n"
	MsgReleaseScanPrompt  = "\n  Proceed with all releases? [y/N]: "
	MsgReleaseScanYes     = "\n  → -y supplied; proceeding without prompt.\n"
	MsgReleaseScanAborted = "  ✗ Multi-repo release aborted by user.\n"
	MsgReleaseScanRunning = "\n  ── Releasing %s → %s ──\n"
	MsgReleaseScanFail    = "  ✗ Release failed for %s: %v\n"
	MsgReleaseScanPartial = "\n  ⚠ %d of %d release(s) failed.\n"
	MsgReleaseScanDone    = "\n  ✓ All %d release(s) complete.\n"
)

// Release orphaned metadata messages.
const (
	MsgReleaseOrphanedMeta    = "  ⚠ Release metadata exists for %s but no tag or branch was found.\n"
	MsgReleaseOrphanedPrompt  = "  → Do you want to remove the release JSON and proceed? (y/N): "
	MsgReleaseOrphanedRemoved = "  ✓ Removed orphaned release metadata for %s\n"
	ErrReleaseOrphanedRemove  = "failed to remove release metadata at %s: %w (operation: delete)"
	ErrReleaseAborted         = "release aborted by user"
)

// Self-release messages.
const (
	MsgSelfReleaseSwitch      = "\n  → Self-release: switching to %s\n"
	MsgSelfReleaseReturn      = "  ✓ Returned to %s\n"
	MsgSelfReleaseSameDir     = "\n  → Self-release: already in source repo %s\n"
	MsgSelfReleasePromptPath  = "  → Enter gitmap source repo path: "
	MsgSelfReleaseSavedPath   = "  ✓ Saved gitmap source repo path: %s\n"
	MsgSelfReleaseInvalidPath = "  ✗ Invalid gitmap source repo path: %s\n"
	ErrSelfReleaseExec        = "could not resolve executable path at %s: %w (operation: resolve)"
	ErrSelfReleaseNoRepo      = "could not locate gitmap source repository"
)

// Install hint constants (printed after release for gitmap repos).
const (
	GitmapRepoPrefix     = "github.com/alimtvnetwork/gitmap-v27"
	GitmapRepoOwner      = "github.com/alimtvnetwork/"
	GitmapRepoNamePrefix = "gitmap-v"
	MsgInstallHintHeader = `
  📦 Install gitmap %s
  ─────────────────────
`
	MsgInstallHintWindows = `
  🪟  Windows · PowerShell
     irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/install.ps1 | iex
`
	// Trailing blank line (the second \n after the curl command) ensures
	// the shell prompt (PS1) lands on its own visually-separated line
	// instead of sitting flush under the install one-liner. Matches the
	// auto-register message convention (releaseautoregister.go line 48).
	MsgInstallHintUnix = `
  🐧  Linux / macOS
     curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/install.sh | sh

`
)

// Release-body pinned-install snippet.
//
// AppendPinnedInstallSnippet renders this block into the GitHub release
// body so anyone copying from the release page installs EXACTLY that
// tag — no "latest" lookup, no -v<N> sibling-repo discovery.
//
// Spec: spec/07-generic-release/08-pinned-version-install-snippet.md
const (
	ReleaseSnippetMarker   = "<!-- gitmap-pinned-install-snippet:%s -->"
	ReleaseSnippetTemplate = "<!-- gitmap-pinned-install-snippet:%s -->\n" +
		"## Install this exact version (%s)\n\n" +
		"Copy-paste these snippets to install **this exact tag**. " +
		"They skip the GitHub `latest` lookup and the versioned-repo discovery probe.\n\n" +
		"**Windows (PowerShell)**\n" +
		"```powershell\n" +
		"$ver = '%s'\n" +
		"$installer = irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/install.ps1\n" +
		"& ([scriptblock]::Create($installer)) -Version $ver -NoDiscovery\n" +
		"```\n\n" +
		"**Linux / macOS (bash)**\n" +
		"```bash\n" +
		"curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/install.sh \\\n" +
		"  | bash -s -- --version %s --no-discovery\n" +
		"```\n"
)

// Release rollback messages.
const (
	MsgRollbackStart  = "\n  ⚠ Push failed — rolling back local branch and tag...\n"
	MsgRollbackBranch = "  ✓ Deleted local branch %s\n"
	MsgRollbackTag    = "  ✓ Deleted local tag %s\n"
	MsgRollbackDone   = "  ✓ Rollback complete. No changes remain.\n"
	MsgRollbackWarn   = "  ⚠ Rollback warning (%s): %v\n"
)
