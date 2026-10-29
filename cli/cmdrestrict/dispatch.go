package cmdrestrict

import "fmt"

const (
	actionStatus = "status"
	actionAdd    = "add"
	actionEdit   = "edit"
	actionMode   = "mode"
)

// Run prints or updates named commit restrictions.
func Run(args []string) error {
	if len(args) == 0 {
		return printStatus()
	}
	switch args[0] {
	case actionStatus:
		return printStatus()
	case actionAdd:
		return addRestriction(args[1:])
	case actionEdit:
		return editRestriction(args[1:])
	case actionMode:
		return runMode(args[1:])
	default:
		return fmt.Errorf("unknown restrict action %s", args[0])
	}
}
