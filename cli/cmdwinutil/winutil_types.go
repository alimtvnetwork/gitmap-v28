package cmdwinutil

// WinRemovalResult records the outcome of a Windows utility removal operation.
type WinRemovalResult struct {
	Target       string   `json:"target"`
	IsSuccess    bool     `json:"isSuccess"`
	IsDryRun     bool     `json:"isDryRun"`
	ItemsRemoved int      `json:"itemsRemoved"`
	Actions      []string `json:"actions"`
	Warnings     []string `json:"warnings,omitempty"`
}

// EdgeRemovalOptions configures Microsoft Edge uninstallation.
type EdgeRemovalOptions struct {
	HasKeepWebView2 bool `json:"hasKeepWebView2"`
	IsDryRun        bool `json:"isDryRun"`
}

// WinUtilOptions configures general WinUtil command flags.
type WinUtilOptions struct {
	IsDryRun        bool `json:"isDryRun"`
	HasAutoYes      bool `json:"hasAutoYes"`
	HasKeepWebView2 bool `json:"hasKeepWebView2"`
	IsJSON          bool `json:"isJson"`
}
