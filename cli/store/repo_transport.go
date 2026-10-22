package store

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// UpdateRepoTransport updates the remote SSH URL and identified transport for a repository
// identified by its absolute path. It canonicalizes the storage path and timestamp.
func (db *DB) UpdateRepoTransport(repoPath, sshURL, transport string) error {
	cleanPath := NormalizeStoragePath(repoPath)
	cleanURL := strings.TrimSpace(sshURL)
	cleanTransport := cleanTransportOrDefault(transport)

	_, err := ExecWrapper(db.conn, constants.SQLUpdateRepoTransport, cleanURL, cleanTransport, cleanPath).Destruct()

	return err
}

func cleanTransportOrDefault(transport string) string {
	clean := strings.ToLower(strings.TrimSpace(transport))
	if clean == "" {
		return constants.TransportSSH
	}

	return clean
}
