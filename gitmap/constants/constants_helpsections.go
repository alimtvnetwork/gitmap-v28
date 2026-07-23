package constants

// Help section headers and flag-line help strings, grouped by command domain.
// Extracted from constants_cli.go to keep that file under the 200-line guideline.
const (
	// Scan flags help section.
	HelpScanFlags             = "Scan flags:"
	HelpConfig                = "  --config <path>     Config file (default: ./data/config.json)"
	HelpMode                  = "  --mode ssh|https    Clone URL style (default: https)"
	HelpOutput                = "  --output csv|json|terminal  Output format (default: terminal)"
	HelpOutputPath            = "  --output-path <dir> Output directory (default: .gitmap/output)"
	HelpOutFile               = "  --out-file <path>   Exact output file path"
	HelpScanFlagGitHubDesktop = "  --github-desktop    Add repos to GitHub Desktop"
	HelpOpen                  = "  --open              Open output folder after scan"
	HelpQuiet                 = "  --quiet             Suppress clone help section (for CI/scripted use)"

	// Clone flags help section.
	HelpCloneFlags = "Clone flags:"
	HelpTargetDir  = "  --target-dir <dir>  Base directory for clones (default: .)"
	HelpSafePull   = "  --safe-pull         Pull existing repos with retry + unlock diagnostics (auto-enabled)"
	HelpVerbose    = "  --verbose           Write detailed debug log to a timestamped file"

	// Release flags help section.
	HelpReleaseFlags  = "Release flags:"
	HelpAssets        = "  --assets <path>     Directory or file to attach to the release"
	HelpCommit        = "  --commit <sha>      Create release from a specific commit"
	HelpRelBranch     = "  --branch <name>     Create release from latest commit of a branch"
	HelpBump          = "  --bump major|minor|patch  Auto-increment from latest released version"
	HelpDraft         = "  --draft             Create an unpublished draft release"
	HelpDryRun        = "  --dry-run           Preview release steps without executing"
	HelpCompressFlag  = "  --compress          Wrap assets in .zip (Windows) or .tar.gz archives"
	HelpChecksumsFlag = "  --checksums         Generate SHA256 checksums.txt for assets"
)
