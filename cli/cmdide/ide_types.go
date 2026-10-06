package cmdide

// IDE target name constants.
const (
	TargetVSCode      = "vscode"
	TargetCursor      = "cursor"
	TargetAntigravity = "antigravity"
	TargetDesktop     = "desktop"
	TargetAll         = "all"
)

// IDEOptions controls target selection and execution modes.
type IDEOptions struct {
	IsVSCodeTargeted      bool   `json:"isVSCodeTargeted"`
	IsCursorTargeted      bool   `json:"isCursorTargeted"`
	IsAntigravityTargeted bool   `json:"isAntigravityTargeted"`
	IsDesktopTargeted     bool   `json:"isDesktopTargeted"`
	IsDryRun              bool   `json:"isDryRun"`
	IsQuiet               bool   `json:"isQuiet"`
	IsJSON                bool   `json:"isJSON"`
	IsForceCreate         bool   `json:"isForceCreate"`
	TargetDirectory       string `json:"targetDirectory,omitempty"`
}

// IDESyncSummary captures registration results across all editors.
type IDESyncSummary struct {
	VSCodeAdded         int `json:"vscodeAdded"`
	VSCodeRetained      int `json:"vscodeRetained"`
	CursorAdded         int `json:"cursorAdded"`
	CursorRetained      int `json:"cursorRetained"`
	AntigravityAdded    int `json:"antigravityAdded"`
	AntigravityRetained int `json:"antigravityRetained"`
	DesktopAdded        int `json:"desktopAdded"`
	DesktopFailed       int `json:"desktopFailed"`
	TotalRepos          int `json:"totalRepos"`
}

// IDEStatusItem details detection and health for a single editor.
type IDEStatusItem struct {
	Name            string `json:"name"`
	IsInstalled     bool   `json:"isInstalled"`
	ExecutablePath  string `json:"executablePath"`
	ConfigPath      string `json:"configPath"`
	RegisteredCount int    `json:"registeredCount"`
}

// IDERepoRegistration describes repository membership across IDEs.
type IDERepoRegistration struct {
	Name                string `json:"name"`
	Path                string `json:"path"`
	IsVSCodeLinked      bool   `json:"isVSCodeLinked"`
	IsCursorLinked      bool   `json:"isCursorLinked"`
	IsAntigravityLinked bool   `json:"isAntigravityLinked"`
	IsDesktopLinked     bool   `json:"isDesktopLinked"`
}

// IDEActionResult reports single-action outcome for add or remove.
type IDEActionResult struct {
	Path                  string `json:"path"`
	Action                string `json:"action"`
	IsVSCodeAffected      bool   `json:"isVSCodeAffected"`
	IsCursorAffected      bool   `json:"isCursorAffected"`
	IsAntigravityAffected bool   `json:"isAntigravityAffected"`
	IsDesktopAffected     bool   `json:"isDesktopAffected"`
}
