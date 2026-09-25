// Package cmdssh — ssh_deploy_keys_types.go defines models for fleet SSH key deployment.
package cmdssh

// DeployKeysNodeResult stores the deployment outcome for a single fleet node.
type DeployKeysNodeResult struct {
	Alias       string `json:"alias"`
	IPAddress   string `json:"ip_address"`
	KeysAdded   int    `json:"keys_added"`
	KeysTotal   int    `json:"keys_total"`
	IsOnline    bool   `json:"is_online"`
	ErrorMsg    string `json:"error,omitempty"`
}

// DeployKeysSummary captures overall mesh public key synchronization metrics.
type DeployKeysSummary struct {
	LocalKeysCollected   int                    `json:"local_keys_collected"`
	RemoteKeysCollected  int                    `json:"remote_keys_collected"`
	UniqueKeysIdentified int                    `json:"unique_keys_identified"`
	NodesTargeted        int                    `json:"nodes_targeted"`
	NodesSucceeded       int                    `json:"nodes_succeeded"`
	NodeResults          []DeployKeysNodeResult `json:"node_results"`
	IsDryRun             bool                   `json:"is_dry_run"`
}
