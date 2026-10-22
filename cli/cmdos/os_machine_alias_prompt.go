// Package cmdos — os_machine_alias_prompt.go handles interactive terminal prompts for machine identity configuration.
package cmdos

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var osStdinReader io.Reader = os.Stdin
var isTerminalFn = defaultIsTerminal

func isInteractiveTerminal() bool {
	return isTerminalFn()
}

// isTerminalInteractive reports whether stdin is an interactive terminal.
func isTerminalInteractive() bool {
	return isInteractiveTerminal()
}

func defaultIsTerminal() bool {
	if os.Getenv("GITMAP_NON_INTERACTIVE") == "1" {
		return false
	}

	return term.IsTerminal(int(os.Stdin.Fd()))
}

func resolveDefaultMachineName(mode string) string {
	_ = mode
	host, err := os.Hostname()
	if err != nil || len(strings.TrimSpace(host)) == 0 {
		return "worker-01"
	}

	return strings.TrimSpace(host)
}

// promptMachineNameFromReader reads one line from r and returns the
// sanitized machine-name token, or "" when the line is blank.
func promptMachineNameFromReader(r io.Reader) (string, error) {
	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}

	cleaned := strings.TrimSpace(line)
	if len(cleaned) == 0 {
		return "", nil
	}

	return sanitizeMachineToken(cleaned), nil
}

func promptForMachineName(mode string) (string, error) {
	defaultVal := resolveDefaultMachineName(mode)
	fmt.Printf("  You forgot to pass the %s. Please enter %s [%s]: ", mode, mode, defaultVal)
	cleaned, err := promptMachineNameFromReader(osStdinReader)
	if err != nil {
		return "", err
	}

	if len(cleaned) == 0 {
		return defaultVal, nil
	}

	return cleaned, nil
}

// resolveMachineAliasArgs resolves target args for the alias set flow.
func resolveMachineAliasArgs(args []string, isJSON bool) ([]string, error) {
	return resolveMachineSetArgs("alias", args, isJSON)
}

func resolveMachineSetArgs(mode string, args []string, isJSON bool) ([]string, error) {
	if len(args) > 0 {
		return args, nil
	}
	if isJSON || !isInteractiveTerminal() {
		return nil, apperror.NewSimple("missing target name/alias (e.g. gitmap "+mode+" set dev-win-01 -y)", "E_MISSING_MACHINE_NAME")
	}

	return promptInteractiveMachineArg(mode)
}

func promptInteractiveMachineArg(mode string) ([]string, error) {
	prompted, err := promptForMachineName(mode)
	if err != nil {
		return nil, apperror.WrapSimple(err, "prompt machine name")
	}
	if len(prompted) == 0 {
		return nil, apperror.NewSimple("machine name cannot be empty", "E_EMPTY_MACHINE_NAME")
	}

	return []string{prompted}, nil
}
