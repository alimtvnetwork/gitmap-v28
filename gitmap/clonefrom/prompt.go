package clonefrom

import (
	"fmt"
	"os"
	"strings"
)

// confirmYesNo prompts on stderr and reads a y/N answer from stdin.
// This is isolated here to avoid circular dependencies with cmd/.
func confirmYesNo(prompt string) bool {
	fmt.Fprintf(os.Stderr, "%s [y/N] ", prompt)
	var resp string
	if _, err := fmt.Fscanln(os.Stdin, &resp); err != nil {
		return false
	}
	resp = strings.ToLower(strings.TrimSpace(resp))

	return resp == "y" || resp == "yes"
}
