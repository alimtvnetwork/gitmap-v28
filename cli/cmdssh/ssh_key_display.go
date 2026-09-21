package cmdssh

import "strings"

func hasRawFlag(args []string) bool {
	for _, a := range args {
		isRaw := a == "--raw" || a == "-r"
		if isRaw {
			return true
		}
	}

	return false
}

func formatDisplayPublicKey(pubKey string, isRaw bool) string {
	return strings.TrimSpace(pubKey)
}
