package cmdworkgroup

// WorkgroupRecord represents a registered named workgroup.
type WorkgroupRecord struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	RootDirectory string `json:"root_directory"`
	SSHKeyPath    string `json:"ssh_key_path"`
	Description   string `json:"description"`
	GitUserName   string `json:"git_user_name"`
	GitUserEmail  string `json:"git_user_email"`
	IsDefault     bool   `json:"is_default"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// FolderBindingRecord represents a path-specific SSH key binding.
type FolderBindingRecord struct {
	ID           int64  `json:"id"`
	FolderPath   string `json:"folder_path"`
	WorkgroupID  int64  `json:"workgroup_id"`
	SSHKeyPath   string `json:"ssh_key_path"`
	GitUserName  string `json:"git_user_name"`
	GitUserEmail string `json:"git_user_email"`
	CreatedAt    string `json:"created_at"`
}

// WorkgroupOptions holds parsed CLI flags and arguments for workgroup commands.
type WorkgroupOptions struct {
	Subcommand  string
	Name        string
	Path        string
	SSHKeyPath  string
	Description string
	UserName    string
	UserEmail   string
	Host        string
	IsDefault   bool
	IsJSON      bool
}

// EffectiveWorkgroupCredentials represents resolved credentials for a path.
type EffectiveWorkgroupCredentials struct {
	WorkgroupName  string `json:"workgroup_name"`
	RootDirectory  string `json:"root_directory"`
	SSHKeyPath     string `json:"ssh_key_path"`
	GitUserName    string `json:"git_user_name"`
	GitUserEmail   string `json:"git_user_email"`
	HasSSHKey      bool   `json:"has_ssh_key"`
	HasGitIdentity bool   `json:"has_git_identity"`
	IsInherited    bool   `json:"is_inherited"`
	HasActiveGroup bool   `json:"has_active_group"`
}
