package cmdos

// LocalOSProbe holds platform probe metrics.
type LocalOSProbe struct {
	OSType           string
	OSGroup          string
	OSVersion        string
	BuildVersion     string
	Kernel           string
	OSVersionDetails OSVersionInfo
}

// OSVersionInfo provides typed, structured operating system and kernel version metrics.
type OSVersionInfo struct {
	Type            string `json:"type"`
	IsWindowsServer bool   `json:"isWindowsServer"`
	IsWindowsPro    bool   `json:"isWindowsPro"`
	IsWindowsHome   bool   `json:"isWindowsHome"`
	Version         string `json:"version"`
	Major           int    `json:"major"`
	Minor           int    `json:"minor"`
	Build           int    `json:"build"`
	Release         string `json:"release"`
	Edition         string `json:"edition"`
	KernelVersion   string `json:"kernelVersion,omitempty"`
}

// OSInfoReport contains comprehensive system and operating system profile information.
type OSInfoReport struct {
	OSType           string         `json:"osType"`
	OSGroup          string         `json:"osGroup"`
	OSVersion        string         `json:"osVersion"`
	BuildVersion     string         `json:"buildVersion"`
	Architecture     string         `json:"architecture"`
	Platform         string         `json:"platform"`
	Hostname         string         `json:"hostname"`
	NumCPU           int            `json:"numCpu"`
	Kernel           string         `json:"kernel,omitempty"`
	GitPath          string         `json:"gitPath,omitempty"`
	BashPath         string         `json:"bashPath,omitempty"`
	PowerShellPath   string         `json:"powerShellPath,omitempty"`
	HasGit           bool           `json:"hasGit"`
	HasBash          bool           `json:"hasBash"`
	HasPowerShell    bool           `json:"hasPowerShell"`
	OSVersionDetails *OSVersionInfo `json:"osVersionDetails,omitempty"`
}
