package macro

import "time"

// Macro represents a named sequence of shell commands.
type Macro struct {
	ID          int64       `json:"id" yaml:"id"`
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	CreatedAt   time.Time   `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" yaml:"updated_at"`
	TotalSteps  int         `json:"total_steps" yaml:"total_steps"`
	Tags        string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Steps       []MacroStep `json:"steps" yaml:"steps"`
}

// MacroStep represents a single executable command within a macro.
type MacroStep struct {
	ID              int64  `json:"id" yaml:"id"`
	MacroID         int64  `json:"macro_id" yaml:"macro_id"`
	StepNum         int    `json:"step_num" yaml:"step_num"`
	CommandLine     string `json:"command_line" yaml:"command_line"`
	WorkingDir      string `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`
	ContinueOnError bool   `json:"continue_on_error" yaml:"continue_on_error"`
	TimeoutSeconds  int    `json:"timeout_seconds" yaml:"timeout_seconds"`
}

// ExecOptions holds runtime options for macro execution.
type ExecOptions struct {
	DryRun   bool
	Verbose  bool
	JSON     bool
	YAML     bool
	FilePath string
}
