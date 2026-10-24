package constants

// CLI command and alias literals.
const (
	CmdCargo      = "cargo"
	CmdCargoAlias = "crg"
	CmdInstall    = "install"
	CmdIn         = "in"
	CmdStatus     = "status"
	CmdSt         = "st"
	CmdInfo       = "info"
	CmdRust       = "rust"
)

// CLI flag literals.
const (
	FlagLongInstall  = "--install"
	FlagShortInstall = "-i"
	FlagLongHelp     = "--help"
	FlagShortHelp    = "-h"
	FlagLongYes      = "--yes"
	FlagShortYes     = "-y"
	FlagLongVersion  = "--version"
	FlagShortVersion = "-V"
	FlagLongRelease  = "--release"
	FlagShortRelease = "-r"
	FlagLongCheck    = "--check"
	FlagDoubleDash   = "--"
	FlagInstall      = "install"
	FlagHelp         = "help"
)

// Binary names.
const (
	BinCargo  = "cargo"
	BinRustc  = "rustc"
	BinRustup = "rustup"
	BinClang  = "clang"
)

// File and path tokens.
const (
	DirCargo        = ".cargo"
	DirBin          = "bin"
	ExtExe          = ".exe"
	PathUsrBin      = "/usr/bin"
	PathUsrLocalBin = "/usr/local/bin"
)

// UI display names and titles.
const (
	DisplayNameCargo   = "Cargo"
	DisplayNameRustc   = "Rustc"
	DisplayNameRustup  = "Rustup"
	StatusHeaderTitle  = "Cargo & Rust Toolchain Status:"
	StatusInstalled    = "installed"
	StatusNotInstalled = "Not installed"
	PromptToInstall    = "  To install Cargo, run: "
	KeyColumnWidth     = 14
)

// Platform installation suggestions.
const (
	WingetInstallRust      = "    winget install Rustlang.Rustup\n    (or: choco install rust)\n\n"
	ChocoInstallRust       = "choco install rust"
	BrewInstallRust        = "    brew install rust\n\n"
	AptInstallCargo        = "    sudo apt install -y cargo\n\n"
	SuggestOfficialRustup  = "    https://rustup.rs/\n    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh\n\n"
	SuggestMissingHeader   = "✗ Cargo (Rust toolchain) is not installed or not found in PATH."
	SuggestGitMapInstall   = "gitmap install cargo"
	SuggestGitMapInCargo   = "gitmap in cargo"
	SuggestGitMapInstallRust = "gitmap install rust"
)

// Dracula theme Lipgloss color hex codes.
const (
	ColorDraculaPurple = "#bd93f9"
	ColorDraculaWhite  = "#f8f8f2"
	ColorDraculaGreen  = "#50fa7b"
	ColorDraculaRed    = "#ff5555"
	ColorDraculaGray   = "#a6adc8"
	ColorDraculaCyan   = "#8be9fd"
)

// Error codes and messages.
const (
	ErrCodeMissingCargo   = "E7100"
	MissingCargoErrorMsg  = "cargo is not installed or not in PATH. Run 'gitmap install cargo' to install it."
)
