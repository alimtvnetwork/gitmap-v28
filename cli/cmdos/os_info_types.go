package cmdos

// OSInfoReport contains comprehensive system and operating system profile information.
type OSInfoReport struct {
	OSType       string `json:"osType"`
	OSVersion    string `json:"osVersion"`
	Architecture string `json:"architecture"`
	Hostname     string `json:"hostname"`
	NumCPU       int    `json:"numCpu"`
	Kernel       string `json:"kernel,omitempty"`
	Platform     string `json:"platform"`
}
