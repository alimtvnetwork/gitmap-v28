package cmdfleet

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const nodesJumpHelp = `nodes jump — print the ssh command to reach a fleet node.

Description:
  Looks up <node> in the fleet node registry (SSH connections) and
  prints the ssh command that lands you in the node's workspace.
  Read-only: this command never opens an ssh session itself.

Usage:
  gitmap nodes jump <node> [--workspace <dir>]
  gitmap nj <node> [--workspace <dir>]

Flags:
  --workspace <dir>   remote workspace to cd into (default: ~/workspace/repos)
  --help, -h          show this help

Example:
  gitmap nodes jump build-ubuntu-02
`

// RunNodesJump implements `gitmap nodes jump <node>` (alias: nj).
func RunNodesJump(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(nodesJumpHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("nodes jump", flag.ContinueOnError)
	workspace := fs.String("workspace", "~/workspace/repos", "remote workspace directory")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(nodesJumpHelp)
		return apperror.WrapSimple(parseErr, "nodes.jump.flags")
	}
	rest := fs.Args()
	if len(rest) == 0 || strings.TrimSpace(rest[0]) == "" {
		fmt.Print(nodesJumpHelp)
		return appfault.NewValidationError("missing required argument: node")
	}
	token := strings.TrimSpace(rest[0])

	nodes, note, loadErr := loadFleetNodes()
	if loadErr != nil {
		return loadErr
	}
	if note != "" {
		fmt.Printf("nodes jump: %s\n", note)
		return nil
	}
	node, isFound := findFleetNode(nodes, token)
	if !isFound {
		var names []string
		for _, n := range nodes {
			names = append(names, n.Name)
		}
		return appfault.NewNotFoundError(fmt.Sprintf("node %q not in fleet registry (known: %s)", token, strings.Join(names, ", ")))
	}

	target := sshTarget(node)
	ws := strings.TrimSpace(*workspace)
	sshCmd := fmt.Sprintf("ssh %s -t \"cd %s && exec $SHELL -l\"", target, ws)
	fmt.Printf("node:      %s (%s)\n", node.Name, target)
	fmt.Printf("workspace: %s\n", ws)
	fmt.Printf("command:   %s\n", sshCmd)
	if isTerminalStdout() {
		fmt.Println("hint: paste the command above to connect (nothing was executed).")
	}

	return nil
}

// isTerminalStdout reports whether stdout is a terminal.
func isTerminalStdout() bool {
	info, statErr := os.Stdout.Stat()
	if statErr != nil {
		return false
	}

	return (info.Mode() & os.ModeCharDevice) != 0
}
