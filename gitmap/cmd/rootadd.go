package cmd

import (
	"fmt"
	"os"
)

// addUsage is the umbrella usage block printed when `gitmap add` is
// invoked with no subcommand or with an unknown one.
const addUsage = `Usage: gitmap add <subcommand>

Subcommands:
  ignore        Merge curated .gitignore template block into ./.gitignore
  attributes    Merge curated .gitattributes template block into ./.gitattributes
  lfs-install   Run 'git lfs install --local' and merge the lfs/common .gitattributes block

Examples:
  gitmap add ignore go node
  gitmap add ignore go node --dry-run
  gitmap add attributes go rust
  gitmap add lfs-install
`

// dispatchAdd routes `gitmap add <subcommand>` to its handler. Returns
// true when the top-level command was "add" so root.go knows the request
// was consumed (success or failure both count — failure exits inside the
// handler).
func dispatchAdd(command string) bool {
	if command != "add" {
		return false
	}
	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, addUsage)
		os.Exit(1)
	}

	sub, rest := os.Args[2], os.Args[3:]
	switch sub {
	case "ignore":
		runAddIgnore(rest)
	case "attributes":
		runAddAttributes(rest)
	case "lfs-install":
		runAddLFSInstall(rest)
	default:
		fmt.Fprintf(os.Stderr, "unknown 'add' subcommand: %s\n", sub)
		fmt.Fprint(os.Stderr, addUsage)
		os.Exit(1)
	}

	return true
}
