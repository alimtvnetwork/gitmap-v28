package cmddaemon

// DaemonExecReq defines command execution request over REST.
type DaemonExecReq struct {
	Command   string            `json:"command"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	TimeoutMs int64             `json:"timeoutMs,omitempty"`
}

// DaemonExecResp defines command execution response.
type DaemonExecResp struct {
	Status     string `json:"status"`
	Success    bool   `json:"success"`
	ExitCode   int    `json:"exitCode"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int64  `json:"durationMs"`
	Error      string `json:"error,omitempty"`
}

// DaemonStatusResp reports node health and daemon version.
type DaemonStatusResp struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	NodeAlias string `json:"nodeAlias"`
	OS        string `json:"os"`
	UptimeSec int64  `json:"uptimeSec"`
}
