package constants

// Help section headers and flag-line help strings, grouped by command domain.
// Extracted from constants_cli.go to keep that file under the 200-line guideline.
const (
	// Scan flags help section.
	HelpScanFlags             = "Scan flags:"
	HelpConfig                = "  " + ColorCyan + ColorCyan + "--config" + ColorReset + ColorReset + " <path>     Config file (default: ./data/config.json)"
	HelpMode                  = "  " + ColorCyan + ColorCyan + "--mode" + ColorReset + ColorReset + " ssh|https    Clone URL style (default: https)"
	HelpOutput                = "  " + ColorCyan + ColorCyan + "--output" + ColorReset + ColorReset + " csv|json|terminal  Output format (default: terminal)"
	HelpOutputPath            = "  " + ColorCyan + ColorCyan + "--output-path" + ColorReset + ColorReset + " <dir> Output directory (default: .gitmap/output)"
	HelpOutFile               = "  " + ColorCyan + ColorCyan + "--out-file" + ColorReset + ColorReset + " <path>   Exact output file path"
	HelpScanFlagGitHubDesktop = "  " + ColorCyan + ColorCyan + "--github-desktop" + ColorReset + ColorReset + "    Add repos to GitHub Desktop"
	HelpOpen                  = "  " + ColorCyan + ColorCyan + "--open" + ColorReset + ColorReset + "              Open output folder after scan"
	HelpQuiet                 = "  " + ColorCyan + ColorCyan + "--quiet" + ColorReset + ColorReset + "             Suppress clone help section (for CI/scripted use)"

	// Clone flags help section.
	HelpCloneFlags = "Clone flags:"
	HelpTargetDir  = "  " + ColorCyan + ColorCyan + "--target-dir" + ColorReset + ColorReset + " <dir>  Base directory for clones (default: .)"
	HelpSafePull   = "  " + ColorCyan + ColorCyan + "--safe-pull" + ColorReset + ColorReset + "         Pull existing repos with retry + unlock diagnostics (auto-enabled)"
	HelpVerbose    = "  " + ColorCyan + ColorCyan + "--verbose" + ColorReset + ColorReset + "           Write detailed debug log to a timestamped file"

	// Release flags help section.
	HelpReleaseFlags  = "Release flags:"
	HelpAssets        = "  " + ColorCyan + ColorCyan + "--assets" + ColorReset + ColorReset + " <path>     Directory or file to attach to the release"
	HelpCommit        = "  " + ColorCyan + ColorCyan + "--commit" + ColorReset + ColorReset + " <sha>      Create release from a specific commit"
	HelpRelBranch     = "  " + ColorCyan + ColorCyan + "--branch" + ColorReset + ColorReset + " <name>     Create release from latest commit of a branch"
	HelpBump          = "  " + ColorCyan + ColorCyan + "--bump" + ColorReset + ColorReset + " major|minor|patch  Auto-increment from latest released version"
	HelpDraft         = "  " + ColorCyan + ColorCyan + "--draft" + ColorReset + ColorReset + "             Create an unpublished draft release"
	HelpDryRun        = "  " + ColorCyan + ColorCyan + "--dry-run" + ColorReset + ColorReset + "           Preview release steps without executing"
	HelpCompressFlag  = "  " + ColorCyan + ColorCyan + "--compress" + ColorReset + ColorReset + "          Wrap assets in .zip (Windows) or .tar.gz archives"
	HelpChecksumsFlag = "  " + ColorCyan + ColorCyan + "--checksums" + ColorReset + ColorReset + "         Generate SHA256 checksums.txt for assets"
)
