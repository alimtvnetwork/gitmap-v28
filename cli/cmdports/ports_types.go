package cmdports

// PortEntry captures the state and diagnostic metadata of an inspected network port.
type PortEntry struct {
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	ProcessName    string `json:"process_name"`
	PID            int    `json:"pid"`
	State          string `json:"state"`
	FirewallStatus string `json:"firewall_status"`
	Recommendation string `json:"recommendation"`
}

// PortsOptions configures port inspection behaviors and output formatting.
type PortsOptions struct {
	TargetPort   int  `json:"target_port"`
	CommonOnly   bool `json:"common_only"`
	FirewallOnly bool `json:"firewall_only"`
	JSONOutput   bool `json:"json_output"`
}

