package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runCloneAudit implements `gitmap clone --audit`. Read-only: it parses
// the source manifest, computes the planned `git clone` / `git pull`
// command for every record, and prints a diff-style summary. Never
// invokes git, never writes outside stdout. Direct-URL invocations are
// rejected so the audit always operates on a manifest the user can
// inspect later.
func runCloneAudit(cf CloneFlags) error {
	source := resolveCloneShorthand(cf.Source)
	if isDirectURL(source) {
		fmt.Fprint(os.Stderr, constants.ErrCloneAuditDirectURL)
		return apperror.NewSimple("fatal error", "E9000")
	}

	report, err := cloner.PlanCloneAudit(source, cf.TargetDir)
	if err != nil {
		return apperror.NewSimple(constants.ErrCloneAuditLoad, "E9000")
	}
	if printErr := report.Print(os.Stdout); printErr != nil {
		return apperror.NewSimple(constants.ErrCloneAuditLoad, "E9000")
	}
	return nil
}
