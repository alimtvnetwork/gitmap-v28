// Package cmdclone — clone_auth_inject.go injects authentication tokens into HTTPS URLs.
package cmdclone

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// InjectTokenIntoHTTPS formats an HTTPS URL with an embedded authentication token.
func InjectTokenIntoHTTPS(rawURL, token string) string {
	if token == "" || !strings.HasPrefix(strings.ToLower(rawURL), constants.PrefixHTTPS) {
		return rawURL
	}

	rest := rawURL[len(constants.PrefixHTTPS):]
	if at := strings.Index(rest, "@"); at >= 0 {
		rest = rest[at+1:]
	}

	return fmt.Sprintf("%s%s@%s", constants.PrefixHTTPS, token, rest)
}
