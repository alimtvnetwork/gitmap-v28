// Package cmd — macro_root_dispatch.go enables dynamic root-level execution of saved macros.
package cmd

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func dispatchMacroDynamic(command string, shouldAudit bool, auditID int64, auditStart time.Time) bool {
	if isExcludedRootCommand(command) {
		return false
	}

	m, loadErr := macro.LoadMacro(command)
	if loadErr != nil || m == nil {
		return false
	}

	cmdmacro.ExecuteDynamicMacro(command)
	finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

	return true
}

func isExcludedRootCommand(command string) bool {
	return strings.HasPrefix(command, "-") || command == ""
}
