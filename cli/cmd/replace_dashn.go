package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// looksLikeDashN matches strings of the form "-1", "-23", etc.
func looksLikeDashN(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}

	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return true
}

// mustParseDashN converts "-N" to the integer N.
func mustParseDashN(s string) int {
	n, isParsed := parseDashNDigits(s)
	if !isParsed || n < 1 {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceBadN, s)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}

	return n
}

func parseDashNDigits(s string) (int, bool) {
	n := 0
	for i := 1; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
		if n > constants.ReplaceMaxDashN {
			return 0, false
		}
	}

	return n, true
}
