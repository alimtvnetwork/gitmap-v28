package cmdservice

// ServiceInfo holds cross-platform OS service metadata.
type ServiceInfo struct {
	Name        string `json:"name" yaml:"name"`
	Status      string `json:"status" yaml:"status"`
	IsEnabled   bool   `json:"is_enabled" yaml:"is_enabled"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	ExecPath    string `json:"exec_path,omitempty" yaml:"exec_path,omitempty"`
}

// ServiceDriver provides uniform lifecycle methods across Linux, Windows, and macOS.
type ServiceDriver interface {
	ListServices() ([]ServiceInfo, error)
	GetService(name string) (*ServiceInfo, error)
	StartService(name string) error
	StopService(name string) error
	CreateService(name, execPath, description string) error
	RemoveService(name string) error
}

// ServiceExportSchema defines format for multi-service export and import.
type ServiceExportSchema struct {
	ExportedAt string        `json:"exported_at" yaml:"exported_at"`
	Platform   string        `json:"platform" yaml:"platform"`
	Services   []ServiceInfo `json:"services" yaml:"services"`
}
