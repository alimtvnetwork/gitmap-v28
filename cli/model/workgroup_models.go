package model

// WorkDirectory represents a registered parent workspace root.
type WorkDirectory struct {
	WorkDirectoryId  int64  `json:"workDirectoryId"`
	AbsolutePath     string `json:"absolutePath"`
	Label            string `json:"label"`
	IsDefault        bool   `json:"isDefault"`
	HasActiveBinding bool   `json:"hasActiveBinding"`
	RepoCount        int    `json:"repoCount,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// FolderWorkDir binds a specific directory tree to a target WorkDirectory.
type FolderWorkDir struct {
	FolderWorkDirId int64  `json:"folderWorkDirId"`
	FolderPath      string `json:"folderPath"`
	WorkDirectoryId int64  `json:"workDirectoryId"`
	IsRecursive     bool   `json:"isRecursive"`
	Priority        int    `json:"priority"`
	IsActive        bool   `json:"isActive"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// FolderGroup represents a named organizational collection of folders.
type FolderGroup struct {
	FolderGroupId int64  `json:"folderGroupId"`
	GroupName     string `json:"groupName"`
	FolderPath    string `json:"folderPath"`
	Description   string `json:"description"`
	IsDefault     bool   `json:"isDefault"`
	IsActive      bool   `json:"isActive"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// FolderGroupSummary represents an aggregate view of a folder group with item count.
type FolderGroupSummary struct {
	GroupName   string `json:"groupName"`
	FolderCount int    `json:"folderCount"`
	CreatedAt   string `json:"createdAt"`
}

// SshKeyAssociation maps an SSH key to a folder, group, workspace, or repository.
type SshKeyAssociation struct {
	SshKeyAssociationId int64  `json:"sshKeyAssociationId"`
	SshKeyId            int64  `json:"sshKeyId"`
	AssociationType     string `json:"associationType"`
	TargetIdentifier    string `json:"targetIdentifier"`
	FolderPath          string `json:"folderPath,omitempty"`
	GroupName           string `json:"groupName,omitempty"`
	WorkDirectoryId     *int64 `json:"workDirectoryId,omitempty"`
	RepoId              *int64 `json:"repoId,omitempty"`
	GitAuthorName       string `json:"gitAuthorName,omitempty"`
	GitAuthorEmail      string `json:"gitAuthorEmail,omitempty"`
	IsActive            bool   `json:"isActive"`
	IsAutoInjected      bool   `json:"isAutoInjected"`
	Priority            int    `json:"priority"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

// EffectiveFolderCredentials holds the fully resolved credential payload for a git operation.
type EffectiveFolderCredentials struct {
	HasCustomKey    bool   `json:"hasCustomKey"`
	SshKeyId        int64  `json:"sshKeyId"`
	KeyName         string `json:"keyName"`
	PrivatePath     string `json:"privatePath"`
	PublicKey       string `json:"publicKey"`
	Fingerprint     string `json:"fingerprint"`
	AssociationType string `json:"associationType"`
	GitAuthorName   string `json:"gitAuthorName,omitempty"`
	GitAuthorEmail  string `json:"gitAuthorEmail,omitempty"`
	IsAutoInjected  bool   `json:"isAutoInjected"`
}
