package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/ignore"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runIgnore(args []string) error {
	checkHelp(constants.CmdIgnore, args)
	if err := ignore.RunIgnore(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func runIgnoreRm(args []string) error {
	checkHelp(constants.CmdIgnoreRm, args)
	if err := ignore.RunIgnoreRm(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}
