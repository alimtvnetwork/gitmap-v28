// Package cmdssh — remote_dispatch.go: `gitmap remote` subcommand routing.
//
// Moved from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D):
// remote-update dispatch delegates to RunSSHUpdateCLI, which is cmdssh's concern.
package cmdssh

import "fmt"

// RunRemote implements `gitmap remote`: prints usage on help flags,
// routes update/up/u to RunSSHUpdateCLI, and passes anything else
// through to RunSSHUpdateCLI for fleet-wide remote updates.
func RunRemote(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Println("\nUsage: gitmap remote update [--node <target>] [package]")
		fmt.Println("       gitmap update --remote <target> [package]")
		fmt.Println("       gitmap agm update --remote <target>")
		return nil
	}
	if args[0] == "update" || args[0] == "up" || args[0] == "u" {
		return RunSSHUpdateCLI(args[1:])
	}
	return RunSSHUpdateCLI(args)
}
