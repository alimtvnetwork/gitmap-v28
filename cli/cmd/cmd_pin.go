package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"strings"
)

// runPinCLI handles 'gitmap pin [installer] <target> <version>'.
func RunPinCLI(args []string) error {
	cleanArgs := filterLeadingPinKeywords(args)
	if isPinHelp(cleanArgs) {
		printPinHelp()
		return nil
	}

	if len(cleanArgs) < 2 {
		return apperror.NewValidationError("missing arguments: expected gitmap pin <target> <version>")
	}

	target := cleanArgs[0]
	version := cleanArgs[1]
	return cmdinstaller.PinVersion(target, version)
}

func isPinHelp(tokens []string) bool {
	if len(tokens) == 0 {
		return true
	}
	first := tokens[0]
	return first == "-h" || first == "--help" || first == "help"
}

func filterLeadingPinKeywords(args []string) []string {
	tokens := args
	for len(tokens) > 0 {
		first := strings.ToLower(tokens[0])
		isKeyword := first == "installer" || first == "object" || first == "version"
		if !isKeyword {
			break
		}
		tokens = tokens[1:]
	}
	return tokens
}

// runUnpinCLI handles 'gitmap unpin [installer] <target>'.
func RunUnpinCLI(args []string) error {
	cleanArgs := filterLeadingPinKeywords(args)
	if isPinHelp(cleanArgs) {
		printUnpinHelp()
		return nil
	}

	target := cleanArgs[0]
	return cmdinstaller.UnpinVersion(target)
}

func printPinHelp() {
	fmt.Println(`gitmap pin: Version Pinning Manager
========================================================================

USAGE:
  gitmap pin <target> <version>
  gitmap pin installer <name> <version>
  gitmap installer pin <name> <version>

EXAMPLES:
  gitmap pin gitmap v6.402.0
  gitmap pin agm v1.5.0
  gitmap pin installer docker 24.0.5
  gitmap installer ls`)
}

func printUnpinHelp() {
	fmt.Println(`gitmap unpin: Version Unpinning Manager
========================================================================

USAGE:
  gitmap unpin <target>
  gitmap unpin installer <name>
  gitmap installer unpin <name>

EXAMPLES:
  gitmap unpin gitmap
  gitmap unpin agm`)
}
