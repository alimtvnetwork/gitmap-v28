package constants

// Go release asset constants.
const (
	AssetsStagingDir = ".gitmap/release-assets"
	GitHubTokenEnv   = "GITHUB_TOKEN"
)

// Asset flag descriptions.
const (
	FlagDescBin         = "Cross-compile Go binaries and include in release assets"
	FlagDescTargets     = "Comma-separated cross-compile targets (e.g. windows/amd64,linux/arm64)"
	FlagDescListTargets = "Print resolved target matrix and exit"
)

// Asset help text.
const (
	HelpBin         = "  --bin, -b           Cross-compile Go binaries and include in release assets"
	HelpTargets     = "  --targets <list>    Cross-compile targets: windows/amd64,linux/arm64"
	HelpListTargets = "  --list-targets      Print resolved target matrix and exit"
)

// List-targets messages.
const (
	MsgListTargetsHeader = "Release targets (%d):\n"
	MsgListTargetsSource = "Source: %s\n\n"
	MsgListTargetsRow    = "  %s/%s\n"
)

// Asset messages.
const (
	MsgAssetDetected     = "  → Detected Go project: %s\n"
	MsgAssetCrossCompile = "\n  Cross-compiling %d target(s)...\n"
	MsgAssetBuilt        = "  ✓ Built %s (%s/%s)\n"
	MsgAssetBuildSummary = "  → Built %d/%d binaries successfully\n"
	MsgAssetUploaded     = "  " + ColorGreen + "✓ Uploaded" + ColorReset + " %s\n"
	MsgAssetUploadStart  = "\n  " + ColorCyan + "▲ Uploading %d asset(s) to GitHub..." + ColorReset + "\n"
	MsgAssetSkipped      = ""
	MsgAssetNoMain       = "  → No buildable main package found, skipping binaries\n"
	MsgAssetNoGoProject  = ""
	MsgAssetStagingClean = "  ✓ Cleaned up staging directory\n"
	MsgTokenFromSource   = "  " + ColorGreen + "🔑 Authenticated via " + ColorReset + ColorCyan + "%s" + ColorReset + "\n"
)

// Release-version snapshot install scripts (spec 105).
//
// The same `release-version.ps1` / `release-version.sh` source ships in
// two forms: the always-current generic script under gitmap/scripts/, and
// per-release snapshots uploaded as release assets with the version baked
// in. The naming convention is `release-version-<tag>.<ext>` so each
// release page can deep-link to a frozen, drift-proof installer.
const (
	ScriptReleaseVersionPS1 = "release-version.ps1"
	ScriptReleaseVersionSh  = "release-version.sh"

	ReleaseVersionSnapshotPS1Fmt = "release-version-%s.ps1"
	ReleaseVersionSnapshotShFmt  = "release-version-%s.sh"

	MsgReleaseScriptSnapshot = "  ✓ Generated release-script snapshot: %s\n"
	ErrReleaseScriptSnapshot = "Error: snapshot generation failed for %s: %v (operation: bake-version)\n"
)

// Asset dry-run messages.
const (
	MsgAssetDryRunHeader = "  [dry-run] Would cross-compile %d binaries:\n"
	MsgAssetDryRunBinary = "    → %s\n"
	MsgAssetDryRunUpload = "  [dry-run] Would upload %d assets\n"
)

// Asset error messages — Code Red: all file errors include exact path and reason.
const (
	ErrAssetBuildFailed = "Error: build failed for %s/%s: %s (operation: compile)\n"
	ErrAssetUploadFinal = "Error: upload failed for asset %s: %v (operation: upload)\n"
	ErrAssetNoToken     = ColorYellow + "  ⚠ No GitHub token available — skipping asset upload." + ColorReset + "\n" +
		"    " + ColorDim + "Tried: GITHUB_TOKEN, GH_TOKEN, `gh auth token`." + ColorReset + "\n" +
		"    " + ColorCyan + "Fix: run `gh auth login` once, or set GITHUB_TOKEN." + ColorReset + "\n"
	ErrAssetRemoteParse = ColorRed + "Error: could not parse remote origin: %v" + ColorReset + " (operation: resolve)\n"
)

// Retry constants.
const (
	RetryMaxAttempts   = 3
	RetryBaseDelayMs   = 1000
	RetryBackoffFactor = 2
)

// Retry HTTP status codes.
const (
	HTTPTooManyRequests = 429
	HTTPServerErrorMin  = 500
)

// Retry messages.
const (
	MsgRetryAttempt = "  ⟳ Retry %d/%d for %s (waiting %s)...\n"
	MsgRetrySuccess = "  ✓ Uploaded %s (attempt %d)\n"
)
