// Package installer — prompt_rollback.go logs rollback actions on failure.
package installer

import "fmt"

func FormatPromptRollbackNote(repoPath string, err error) string {
	return fmt.Sprintf("Installation in %s failed: %v. Working directory was left intact.", repoPath, err)
}
