package constants

// Antigravity Manager installation and update commands.
const (
	AgManagerWindowsInstallCmd = "irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1 | iex"
	AgManagerUnixInstallCmd    = "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash"
	AgManagerWindowsInstallURL = "https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1"
	AgManagerUnixInstallURL    = "https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh"
	AgManagerGitURL            = "https://github.com/alimtvnetwork/Antigravity-Manager.git"
)
