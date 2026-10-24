package cmdapps

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdapps/enums"

// AppScopeType designates system-wide vs user-specific installations.
// AppScopeType aliases the strongly-typed enum from the enums package.
type AppScopeType = enums.AppScopeType

const (
	ScopeSystem = enums.ScopeSystem
	ScopeUser   = enums.ScopeUser
)

// AppScope is an alias for AppScopeType for backward compatibility.
type AppScope = enums.AppScopeType

// PackageManagerType represents the detected managing package manager.
// PackageManagerType aliases the strongly-typed enum from the enums package.
type PackageManagerType = enums.PackageManagerType

const (
	ManagerApt     = enums.ManagerApt
	ManagerDpkg    = enums.ManagerDpkg
	ManagerSnap    = enums.ManagerSnap
	ManagerFlatpak = enums.ManagerFlatpak
	ManagerWinget  = enums.ManagerWinget
	ManagerReg     = enums.ManagerReg
	ManagerNpm     = enums.ManagerNpm
	ManagerCustom  = enums.ManagerCustom
)

// PackageManager is an alias for PackageManagerType for backward compatibility.
type PackageManager = enums.PackageManagerType

// InstalledApp represents an application discovered on the host system.
type InstalledApp struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Exec        string         `json:"exec"`
	Icon        string         `json:"icon,omitempty"`
	Package     string         `json:"package,omitempty"`
	Version     string         `json:"version,omitempty"`
	Manager     PackageManager `json:"manager"`
	Scope       AppScope       `json:"scope"`
	DesktopFile string         `json:"desktopFile,omitempty"`
	IsRemovable bool           `json:"isRemovable"`
}

// AppListResponse is the top-level envelope for gitmap apps list --json.
type AppListResponse struct {
	Success  bool           `json:"success"`
	Count    int            `json:"count"`
	Platform string         `json:"platform"`
	Data     []InstalledApp `json:"data"`
	Error    string         `json:"error,omitempty"`
}

// AppUninstallResponse is the top-level envelope for gitmap apps uninstall --json.
type AppUninstallResponse struct {
	Success      bool     `json:"success"`
	AppID        string   `json:"appId"`
	Package      string   `json:"package,omitempty"`
	Purged       bool     `json:"purged"`
	RemovedFiles []string `json:"removedFiles"`
	CachesReset  []string `json:"cachesReset"`
	DurationMs   int64    `json:"durationMs"`
	Error        string   `json:"error,omitempty"`
}

// ListOptions controls the filtering for listing apps.
type ListOptions struct {
	IsJson    bool
	IsSystem  bool
	IsUser    bool
	IsAll     bool
	FilterStr string
}

// UninstallOptions controls the uninstallation execution.
type UninstallOptions struct {
	IsPurge  bool
	IsForce  bool
	IsDryRun bool
	IsJson   bool
}
