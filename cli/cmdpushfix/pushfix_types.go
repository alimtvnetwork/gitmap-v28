package cmdpushfix

// PushFixOptions represents the user-configurable flags for gitmap push-fix.
type PushFixOptions struct {
	Remote   string
	Branch   string
	IsDryRun bool
	IsForce  bool
	IsYes    bool
	IsSSH    bool
	IsHTTPS  bool
}

// PushFixState captures the runtime git and environment diagnosis state.
type PushFixState struct {
	RepoDir            string
	RemoteName         string
	RemoteURL          string
	Branch             string
	UnpushedCount      int
	HasUpstream        bool
	IsWorkingTreeClean bool
	Transport          string
}

// PushFixResult represents the execution outcome of the recovery pipeline.
type PushFixResult struct {
	IsSuccess        bool
	RemediationsDone []string
	PushedSHA        string
	SummaryMessage   string
}
