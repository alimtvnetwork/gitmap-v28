package cmdlemon

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DaemonSettingKey is the canonical setting key for daemon enablement.
const DaemonSettingKey = constants.SettingDaemonEnabled

var daemonSettingAliases = []string{
	"gitmap-daemon.enabled",
	"gitmap_daemon.enabled",
	"daemon_enabled",
	"daemon.enable",
}

// IsDaemonEnabled returns true if the gitmap-daemon is enabled in settings (default: true).
func IsDaemonEnabled() bool {
	dbConn, errOpen := store.OpenDefault()
	if errOpen == nil {
		defer dbConn.Close()
		val := strings.TrimSpace(dbConn.GetSetting(DaemonSettingKey))
		if len(val) > 0 {
			return !isFalsySetting(val)
		}
		for _, alias := range daemonSettingAliases {
			aliasVal := strings.TrimSpace(dbConn.GetSetting(alias))
			if len(aliasVal) > 0 {
				return !isFalsySetting(aliasVal)
			}
		}
	}

	val, isFound := config.GetVariable("global", DaemonSettingKey)
	if isFound && len(strings.TrimSpace(val)) > 0 {
		return !isFalsySetting(val)
	}

	for _, alias := range daemonSettingAliases {
		aliasVal, hasAlias := config.GetVariable("global", alias)
		if hasAlias && len(strings.TrimSpace(aliasVal)) > 0 {
			return !isFalsySetting(aliasVal)
		}
	}

	return true
}

// SetDaemonEnabled persists the daemon enabled state to DB and config.
func SetDaemonEnabled(isEnabled bool) error {
	strVal := "true"
	if !isEnabled {
		strVal = "false"
	}

	_ = config.SetVariable("global", DaemonSettingKey, strVal)

	dbConn, errOpen := store.OpenDefault()
	if errOpen != nil {
		return errOpen
	}
	defer dbConn.Close()

	return dbConn.SetSetting(DaemonSettingKey, strVal)
}

func isFalsySetting(val string) bool {
	norm := strings.ToLower(strings.TrimSpace(val))
	return norm == "false" || norm == "0" || norm == "off" || norm == "disabled" || norm == "no"
}
