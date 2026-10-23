package cmdcommit

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ensureGitIdentity fails fast with an actionable validation error when no
// git user.name/user.email is configured, instead of letting git commit die
// with "Author identity unknown" streamed above a stack trace. A validation
// error renders cleanly (no stack) in the global error handler.
func ensureGitIdentity() *apperror.AppError {
	name, _ := execGitOutputCP("config", "user.name")
	email, _ := execGitOutputCP("config", "user.email")
	if strings.TrimSpace(name) != "" && strings.TrimSpace(email) != "" {
		return nil
	}

	return apperror.NewValidationError("git identity is not configured: user.name and/or user.email is missing. Set it with: git config --global user.name \"Your Name\" and git config --global user.email \"you@example.com\" (or inspect it with: gitmap whoami)")
}
