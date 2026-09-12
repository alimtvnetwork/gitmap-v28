// Package constants — constants_prompt_errors.go defines error constants.
package constants

const (
	ErrPromptInstallFailed = "failed to install Prompt Architect in %s: %v"
	ErrPromptNoPwsh        = "PowerShell is required on Windows to install Prompt Architect"
	ErrPromptNoBash        = "Bash and curl are required on Unix to install Prompt Architect"
)
