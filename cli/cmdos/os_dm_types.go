package cmdos

// DMType enumerates supported display managers.
type DMType string

const (
	DMTypeGDM3    DMType = "gdm3"
	DMTypeLightDM DMType = "lightdm"
	DMTypeSDDM    DMType = "sddm"
	DMTypeUnknown DMType = "unknown"
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
