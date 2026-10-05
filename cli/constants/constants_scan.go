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
