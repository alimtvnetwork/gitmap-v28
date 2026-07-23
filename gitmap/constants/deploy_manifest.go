package constants

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

// deployManifestJSON is the SINGLE SOURCE OF TRUTH for deploy-target folder
// names. Loaders in run.ps1, run.sh, and gitmap/scripts/install.sh read the
// same JSON file from disk; Go embeds it at build time so the binary stays
// self-contained. Renaming the deploy folder ONLY requires editing
// gitmap/constants/deploy-manifest.json — no other source change.
//
//go:embed deploy-manifest.json
var deployManifestJSON []byte

// DeployManifest mirrors the on-disk JSON layout. Keep field names in sync
// with gitmap/constants/deploy-manifest.json.
type DeployManifest struct {
	SchemaVersion    int      `json:"schemaVersion"`
	AppSubdir        string   `json:"appSubdir"`
	LegacyAppSubdirs []string `json:"legacyAppSubdirs"`
	BinaryName       struct {
		Windows string `json:"windows"`
		Unix    string `json:"unix"`
	} `json:"binaryName"`
	SourceRepoSubdir string `json:"sourceRepoSubdir"`
}

// Manifest is the parsed deploy manifest, loaded once at package init.
var Manifest DeployManifest

// GitMapSubdir, GitMapCliSubdir are populated from Manifest at init().
// They remain `var` (not `const`) so the manifest is the only edit point.
//
//nolint:gochecknoglobals // populated from embedded manifest, single source of truth
var (
	// GitMapSubdir is the SOURCE-REPO subdirectory name (<RepoRoot>/gitmap/...).
	GitMapSubdir string
	// GitMapCliSubdir is the DEPLOY-TARGET subdirectory name
	// (<DeployRoot>/gitmap-cli/gitmap.exe). Renamed from "gitmap" in v3.6.0.
	// Sourced from deploy-manifest.json — never hardcode this string.
	GitMapCliSubdir string
	// LegacyAppSubdirs lists deploy-folder names from prior schema versions.
	// Migration code uses this to detect and rename old layouts.
	LegacyAppSubdirs []string
)

func init() {
	if err := json.Unmarshal(deployManifestJSON, &Manifest); err != nil {
		fmt.Fprintf(os.Stderr, "constants: failed to parse embedded deploy-manifest.json: %v\n", err)
		// Fall back to v3.13.x defaults so the binary stays usable.
		Manifest.AppSubdir = "gitmap-cli"
		Manifest.SourceRepoSubdir = "gitmap"
		Manifest.LegacyAppSubdirs = []string{"gitmap"}
	}
	GitMapSubdir = Manifest.SourceRepoSubdir
	GitMapCliSubdir = Manifest.AppSubdir
	LegacyAppSubdirs = Manifest.LegacyAppSubdirs
}
