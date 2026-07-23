package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runCloneAudit implements `gitmap clone --audit`. Read-only: it parses
// the source manifest, computes the planned `git clone` / `git pull`
// command for every record, and prints a diff-style summary. Never
// invokes git, never writes outside stdout. Direct-URL invocations are
// rejected so the audit always operates on a manifest the user can
// inspect later.
func runCloneAudit(cf CloneFlags) {
	source := resolveCloneShorthand(cf.Source)
	if isDirectURL(source) {
		fmt.Fprint(os.Stderr, constants.ErrCloneAuditDirectURL)
		os.Exit(1)
	}

	report, err := cloner.PlanCloneAudit(source, cf.TargetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneAuditLoad, source, err)
		os.Exit(1)
	}
	if printErr := report.Print(os.Stdout); printErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneAuditLoad, source, printErr)
		os.Exit(1)
	}
}
