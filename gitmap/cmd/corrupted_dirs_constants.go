package cmd

// Constants defining pattern matching criteria for corrupted directories.
const (
	CorruptedPatternQuickInstaller = "quick installer"
	CorruptedPatternInstaller      = "gitmap installer"
	CorruptedPatternDefault        = "Default:"
	CorruptedPatternPrompt         = "Choose install folder"
	CorruptedPatternPath           = "Install path"
	CorruptedTildeDir              = "~"
	AnsiEscapePrefix               = "\x1b["
	AnsiEscapeOctal                = "\033["
)

// CorruptedDirInfo describes a detected corrupted installation directory.
type CorruptedDirInfo struct {
	Path           string   `json:"path"`
	Name           string   `json:"name"`
	Reason         string   `json:"reason"`
	IsLiteralTilde bool     `json:"isLiteralTilde"`
	HasFiles       bool     `json:"hasFiles"`
	Files          []string `json:"files,omitempty"`
}

// CleanOptions specifies configuration for corrupted directory cleanup.
type CleanOptions struct {
	IsDryRun bool `json:"isDryRun"`
	IsForce  bool `json:"isForce"`
}

// CleanResult summarizes the outcome of a cleanup operation.
type CleanResult struct {
	DetectedDirs   []CorruptedDirInfo `json:"detectedDirs"`
	RemovedDirs    []string           `json:"removedDirs"`
	RecoveredFiles []string           `json:"recoveredFiles"`
	EscapedFrom    string             `json:"escapedFrom,omitempty"`
	EscapedTo      string             `json:"escapedTo,omitempty"`
}
