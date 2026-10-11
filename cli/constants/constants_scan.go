package constants

// DefaultScanExcludeDirs defines canonical directories automatically ignored by scanner and discovery.
var DefaultScanExcludeDirs = []string{
	".oh-my-zsh",
	"oh-my-zsh",
	"ohmyzsh",
	".ohmyzsh",
	"omz",
	".omz",
	"omizssh",
	"node_modules",
	"vendor",
	".cache",
	"cache",
	"caches",
	".venv",
	"venv",
	"__pycache__",
	".cargo",
	".rustup",
	".git",
	".terraform",
	".next",
	".turbo",
	"dist",
	"build",
	"bin",
	"obj",
	"target",
	"appdata",
	".cursor",
	".vscode",
	".vscode-insiders",
	".vscode-server",
	".antigravity",
	".antigravity_tools",
	".gemini",
	"plugins",
	"plugin",
	"uv",
	".uv",
	"pip",
	".pip",
	".npm",
	".pnpm",
	".pnpm-store",
	"pnpm-store",
	".gradle",
	".m2",
	".dotnet",
	".nuget",
	"temp",
	"tmp",
	".temp",
	".tmp",
	".local",
	".config",
	"sdists",
	"wheels",
	"site-packages",
}

// Scan flag names, aliases, and descriptions.
const (
	FlagScanForceInclude      = "force-include"
	FlagScanForceIncludeAlias = "fi"
	FlagDescScanForceInclude  = "Comma-separated directories to include despite default exclusions (or 'all' to include everything)"
	FlagAddForce              = "force"
	FlagAddForceAlias         = "f"
	FlagDescAddForce          = "Force tracking of directory even if matching default exclusions"
)
