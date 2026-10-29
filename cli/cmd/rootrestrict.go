package cmd

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdrestrict"

func runRestrictUtility() error {
	checkHelp("restrict", argsTail())

	return cmdrestrict.Run(argsTail())
}
