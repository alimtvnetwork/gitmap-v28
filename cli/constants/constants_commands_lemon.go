package constants

// gitmap:cmd top-level
// GitMap daemon commands (formerly Lemon daemon).
const (
	// CmdDaemon is the primary top-level namespace for the GitMap background daemon.
	CmdDaemon       = "daemon"
	CmdGitmapDaemon = "gitmap-daemon"

	// CmdLemon is the backward-compatible namespace for the Lemon email daemon and command service.
	CmdLemon      = "lemon"
	CmdLemonAlias = "lm"

	// SettingDaemonEnabled is the settings key used to switch the daemon on or off.
	SettingDaemonEnabled = "daemon.enabled"
)
