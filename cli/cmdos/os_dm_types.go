package cmdos

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos/enums"
)

// DMType enumerates supported display managers.
type DMType = enums.DMType

const (
	DMTypeGDM3    = enums.DMTypeGDM3
	DMTypeLightDM = enums.DMTypeLightDM
	DMTypeSDDM    = enums.DMTypeSDDM
	DMTypeUnknown = enums.DMTypeUnknown
)

// DMStatus holds display manager diagnostics and session states.
type DMStatus struct {
	Name             string `json:"name"`
	ServiceStatus    string `json:"service_status"`
	SessionType      string `json:"session_type"`
	IsWaylandEnabled bool   `json:"is_wayland_enabled"`
	IsAutoLoginSet   bool   `json:"is_autologin_set"`
	AutoLoginUser    string `json:"autologin_user"`
	ConfigFile       string `json:"config_file"`
}

// DMManager provides display manager operations.
type DMManager interface {
	GetStatus() (DMStatus, error)
	SetWayland(isEnabled bool) error
	RestartService() error
}
