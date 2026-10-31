// Package cmd - rootutility_user.go: user-identity dispatch entries (change-user, change-login, autobot).
// Split from rootutility.go on 31-Oct-26 so rootutility.go stays under the 300-line file-size gate.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmduser"
)

func utilityUserEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"change-user", "cu"}, runChangeUserUtility},
		// WHY no "cl" alias (31-Oct-26): "cl" is already the
		// changelog alias, and the dispatch registry test fails on
		// the collision. change-login keeps its full name only.
		{[]string{"change-login"}, runChangeLoginUtility},
		{[]string{"autobot", "bot"}, runAutobotUtility},
		{[]string{"restrict", "restriction"}, runRestrictUtility},
	}
}

func runChangeUserUtility() error {
	checkHelp("change-user", argsTail())

	_, err := cmduser.DispatchUser("change-user")

	return err
}

func runChangeLoginUtility() error {
	checkHelp("change-login", argsTail())

	_, err := cmduser.DispatchUser("change-login")

	return err
}

func runAutobotUtility() error {
	checkHelp("autobot", argsTail())

	_, err := cmduser.DispatchUser("autobot")

	return err
}
