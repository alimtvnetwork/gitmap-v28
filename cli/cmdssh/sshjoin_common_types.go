package cmdssh

// SSHCommonTarget holds parsed IP, alias, and connection state for batch join.
type SSHCommonTarget struct {
	FullIP     string `json:"fullIp"`
	Port       int    `json:"port"`
	Alias      string `json:"alias"`
	Username   string `json:"username"`
	IsSuccess  bool   `json:"isSuccess"`
	DetectedOS string `json:"detectedOs,omitempty"`
	OSVersion  string `json:"osVersion,omitempty"`
	OSArch     string `json:"osArch,omitempty"`
	ErrorMsg   string `json:"errorMsg,omitempty"`
}

// SSHCommonJoinResult holds the aggregated results of a batch common join.
type SSHCommonJoinResult struct {
	Username     string            `json:"username"`
	TotalCount   int               `json:"totalCount"`
	SuccessCount int               `json:"successCount"`
	FailureCount int               `json:"failureCount"`
	Targets      []SSHCommonTarget `json:"targets"`
}

// SSHCommonJoinOptions holds CLI arguments and flags for ssh-join-common.
type SSHCommonJoinOptions struct {
	Username        string
	RawIPs          []string
	Password        string
	Port            int
	IsJSON          bool
	DryRun          bool
	IsHelpRequested bool
}
