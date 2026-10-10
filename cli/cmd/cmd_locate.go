package cmd

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdautomation"

func RunLocateTopLevel(args []string) error {
	return cmdautomation.DispatchAutomation(append([]string{"locate"}, args...))
}
