package macro

import "time"

// Macro represents a named sequence of shell commands.
type Macro struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	TotalSteps  int         `json:"total_steps"`
	Tags        string      `json:"tags,omitempty"`
	Steps       []MacroStep `json:"steps"`
}

// MacroStep represents a single executable command within a macro.
type MacroStep struct {
	ID              int64  `json:"id"`
	MacroID         int64  `json:"macro_id"`
	StepNum         int    `json:"step_num"`
	CommandLine     string `json:"command_line"`
	WorkingDir      string `json:"working_dir,omitempty"`
	ContinueOnError bool   `json:"continue_on_error"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
}

// ExecOptions holds runtime options for macro execution.
type ExecOptions struct {
	DryRun  bool
	Verbose bool
}
