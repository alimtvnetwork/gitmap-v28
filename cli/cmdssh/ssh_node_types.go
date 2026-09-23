package cmdssh

// NodeVersionInfo captures gitmap version and status on a remote SSH node.
type NodeVersionInfo struct {
	Alias       string `json:"alias"`
	Host        string `json:"host"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	Version     string `json:"version"`
	IsInstalled bool   `json:"isInstalled"`
	IsOnline    bool   `json:"isOnline"`
	ErrorMsg    string `json:"errorMsg,omitempty"`
}

// SSHFleetUpdateOptions defines flags for fleet-wide binary updates.
type SSHFleetUpdateOptions struct {
	Target   string
	Pkg      string
	Except   string
	Exclude  string
	IsDryRun bool
	IsForce  bool
}
