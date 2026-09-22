// Package cmdclone — clone_auth_store.go manages thread-safe global access token storage.
package cmdclone

import "sync"

var (
	globalTokenMu sync.RWMutex
	globalToken   string
)

// SetGlobalAccessToken saves an access token for reuse across all repositories in the session.
func SetGlobalAccessToken(token string) {
	globalTokenMu.Lock()
	defer globalTokenMu.Unlock()

	globalToken = token
}

// GetGlobalAccessToken retrieves the cached session token if available.
func GetGlobalAccessToken() (string, bool) {
	globalTokenMu.RLock()
	defer globalTokenMu.RUnlock()

	isDefined := len(globalToken) > 0

	return globalToken, isDefined
}

// ClearGlobalAccessToken clears the stored access token.
func ClearGlobalAccessToken() {
	globalTokenMu.Lock()
	defer globalTokenMu.Unlock()

	globalToken = ""
}
