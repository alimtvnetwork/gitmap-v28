package cmdui

// SettingsData models application configuration settings.
type SettingsData struct {
	Theme         string            `json:"theme"`
	DefaultRemote string            `json:"defaultRemote"`
	ClusterPort   int               `json:"clusterPort"`
	AutoDeployKey bool              `json:"isAutoDeployKey"`
	Attributes    map[string]string `json:"attributes"`
}

// CommitinOptions defines options for commit operations.
type CommitinOptions struct {
	Message    string `json:"message"`
	Direction  string `json:"direction"`
	IsAmend    bool   `json:"isAmend"`
	IsPush     bool   `json:"isPush"`
	TargetNode string `json:"targetNode,omitempty"`
}

// NodeSummary represents an SSH cluster node.
type NodeSummary struct {
	NodeId   string `json:"nodeId"`
	Alias    string `json:"alias"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	IsOnline bool   `json:"isOnline"`
	OS       string `json:"os"`
}

// RemoteFileReadReq requests a remote file content.
type RemoteFileReadReq struct {
	NodeAlias string `json:"nodeAlias"`
	FilePath  string `json:"filePath"`
}

// RemoteFileContent represents file content with syntax.
type RemoteFileContent struct {
	NodeAlias string `json:"nodeAlias"`
	FilePath  string `json:"filePath"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error,omitempty"`
}

// RemoteFileSaveReq holds edits to be saved back to a remote host.
type RemoteFileSaveReq struct {
	NodeAlias string `json:"nodeAlias"`
	FilePath  string `json:"filePath"`
	Content   string `json:"content"`
}

// CustomInstallerItem models a custom multi-OS installer definition.
type CustomInstallerItem struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	TargetOS       string `json:"targetOS"`
	InstallCommand string `json:"installCommand"`
	VerifyCommand  string `json:"verifyCommand"`
}
