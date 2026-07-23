package constants

// Auto-derived tag detection for VS Code Project Manager sync (v3.40.0+).
//
// Each marker is a top-level file or directory inside a project's rootPath.
// When present, the corresponding tag is added (additively) to the entry's
// `tags` array in projects.json. User-added tags are NEVER removed.
//
// Detection is shallow (top-level only) and order-stable: results are
// emitted following AutoTagOrder so diffs stay clean across runs.

// Canonical tag identifiers — keep in sync with AutoTagMarkers / AutoTagOrder.
const (
	AutoTagGitmap = "gitmap"
	AutoTagGit    = "git"
	AutoTagNode   = "node"
	AutoTagGo     = "go"
	AutoTagPython = "python"
	AutoTagRust   = "rust"
	AutoTagDocker = "docker"
)

// AutoTagMarkers maps a top-level filesystem entry name to the tag it
// implies. Both files and directories qualify (.git can be either).
var AutoTagMarkers = map[string]string{
	".git":               AutoTagGit,
	"package.json":       AutoTagNode,
	"go.mod":             AutoTagGo,
	"pyproject.toml":     AutoTagPython,
	"requirements.txt":   AutoTagPython,
	"Cargo.toml":         AutoTagRust,
	"Dockerfile":         AutoTagDocker,
	"compose.yaml":       AutoTagDocker,
	"compose.yml":        AutoTagDocker,
	"docker-compose.yml": AutoTagDocker,
}

// AutoTagOrder is the canonical emission order. Tags not listed here are
// dropped (the detector never invents tags outside this list).
var AutoTagOrder = []string{
	AutoTagGit,
	AutoTagNode,
	AutoTagGo,
	AutoTagPython,
	AutoTagRust,
	AutoTagDocker,
}

// CLI flag for opting out of auto-tag detection during sync.
const (
	FlagNoAutoTags     = "no-auto-tags"
	FlagDescNoAutoTags = "skip auto-derived tags (git/node/go/...) when syncing VS Code Project Manager projects.json"
)

// Global tag-customization flags (v4.18.0+). Like
// `--vscode-sync-disabled`, these are stripped from argv at Run() and
// persisted into env vars so every code path that ultimately calls
// vscodepm.DetectTagsCustom inherits the rules without per-flagset
// plumbing. All three are repeatable AND accept comma-separated
// values, e.g. `--vscode-tag work --vscode-tag urgent` is equivalent
// to `--vscode-tag work,urgent`.
//
//	--vscode-tag <name>             always add <name> to every entry
//	--vscode-tag-skip <name>        never emit auto-detected <name>
//	--vscode-tag-marker <file>=<tag>  register marker→tag rule
//
// Env vars use ASCII unit separator (\x1f) between values so commas
// inside individual tokens stay intact.
const (
	FlagVSCodeTag           = "vscode-tag"
	FlagDescVSCodeTag       = "always add this tag to every VS Code Project Manager entry (repeatable; accepts comma-list)"
	FlagVSCodeTagSkip       = "vscode-tag-skip"
	FlagDescVSCodeTagSkip   = "drop this auto-detected tag from every VS Code Project Manager entry (repeatable; accepts comma-list)"
	FlagVSCodeTagMarker     = "vscode-tag-marker"
	FlagDescVSCodeTagMarker = "register a marker→tag rule, e.g. Gemfile=ruby (repeatable; accepts comma-list)"

	EnvVSCodeTagAdd       = "GITMAP_VSCODE_TAG_ADD"
	EnvVSCodeTagSkip      = "GITMAP_VSCODE_TAG_SKIP"
	EnvVSCodeTagMarker    = "GITMAP_VSCODE_TAG_MARKER"
	EnvVSCodeTagSeparator = "\x1f"
	TagMarkerKVSeparator  = "="
)
