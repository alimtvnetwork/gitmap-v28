package cmdapps

// AppScope designates system-wide vs user-specific installations.
type AppScope string

const (
	// ScopeSystem indicates a system-wide installed application.
	ScopeSystem AppScope = "system"
	// ScopeUser indicates a user-scoped application installation.
	ScopeUser AppScope = "user"
)

// PackageManager represents the detected managing package manager.
type PackageManager string

const (
	// ManagerApt represents Debian/Ubuntu APT.
	ManagerApt PackageManager = "apt"
	// ManagerDpkg represents Debian DPKG.
	ManagerDpkg PackageManager = "dpkg"
	// ManagerSnap represents Canonical Snap.
	ManagerSnap PackageManager = "snap"
	// ManagerFlatpak represents Flatpak.
	ManagerFlatpak PackageManager = "flatpak"
	// ManagerWinget represents Windows Package Manager (winget).
	ManagerWinget PackageManager = "winget"
	// ManagerReg represents Windows Registry Uninstall entry.
	ManagerReg PackageManager = "registry"
	// ManagerNpm represents Global Node Package Manager.
	ManagerNpm PackageManager = "npm"
	// ManagerCustom represents a standalone custom installed binary.
	ManagerCustom PackageManager = "custom"
)

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
