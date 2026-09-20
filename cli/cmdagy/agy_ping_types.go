package cmdagy

// AgyPingExecutableCheck holds result of looking up Antigravity IDE executable.
type AgyPingExecutableCheck struct {
	IsFound bool   `json:"isFound"`
	Path    string `json:"path,omitempty"`
	Error   string `json:"error,omitempty"`
}

// AgyPingProcessCheck holds result of checking if Antigravity IDE is running.
type AgyPingProcessCheck struct {
	IsRunning bool   `json:"isRunning"`
	PID       int    `json:"pid,omitempty"`
	Name      string `json:"name,omitempty"`
	Error     string `json:"error,omitempty"`
}

// AgyFilesystemHealth holds health status of Antigravity config and brain dirs.
type AgyFilesystemHealth struct {
	IsAccessible     bool   `json:"isAccessible"`
	BrainDir         string `json:"brainDir,omitempty"`
	ProjectsDir      string `json:"projectsDir,omitempty"`
	ConversationsDir string `json:"conversationsDir,omitempty"`
	Error            string `json:"error,omitempty"`
}

// AgyPingWorkspaceCheck holds workspace-specific conversation status.
type AgyPingWorkspaceCheck struct {
	TargetWorkspace string `json:"targetWorkspace"`
	ConvID          string `json:"convId,omitempty"`
	ConvStatus      string `json:"convStatus"`
	StepCount       int    `json:"stepCount,omitempty"`
	TranscriptPath  string `json:"transcriptPath,omitempty"`
	Error           string `json:"error,omitempty"`
}

// AgyPingQueueCheck holds verification queue status.
type AgyPingQueueCheck struct {
	HasQueue          bool   `json:"hasQueue"`
	ActivePromptCount int    `json:"activePromptCount"`
	QueuedCount       int    `json:"queuedCount"`
	LastUpdated       string `json:"lastUpdated,omitempty"`
	Error             string `json:"error,omitempty"`
}

// AgyPingReport aggregates all Antigravity IDE ping checks.
type AgyPingReport struct {
	Timestamp  string                 `json:"timestamp"`
	Executable AgyPingExecutableCheck `json:"executable"`
	Process    AgyPingProcessCheck    `json:"process"`
	Filesystem AgyFilesystemHealth    `json:"filesystem"`
	Workspace  AgyPingWorkspaceCheck  `json:"workspace"`
	Queue      AgyPingQueueCheck      `json:"queue"`
	IsHealthy  bool                   `json:"isHealthy"`
}
