package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// UpdateRepoTransportWithStatus updates remote transport and returns whether a row was modified.
func (db *DB) UpdateRepoTransportWithStatus(repoPath, sshURL, transport string) (bool, error) {
	res, err := db.execUpdateTransport(repoPath, sshURL, transport)
	if err != nil {
		return false, err
	}

	rows, errRows := res.RowsAffected()
	if errRows != nil {
		return false, errRows
	}

	hasUpdated := rows > 0

	return hasUpdated, nil
}

func (db *DB) execUpdateTransport(repoPath, sshURL, transport string) (sql.Result, error) {
	cleanPath := NormalizeStoragePath(repoPath)
	cleanURL := strings.TrimSpace(sshURL)
	cleanTransport := cleanTransportOrDefault(transport)

	return ExecWrapper(db.conn, constants.SQLUpdateRepoTransport, cleanURL, cleanTransport, cleanPath).Destruct()
}

// UpdateRepoTransport updates the remote SSH URL and identified transport for a repository
// identified by its absolute path. It canonicalizes the storage path and timestamp.
func (db *DB) UpdateRepoTransport(repoPath, sshURL, transport string) error {
	_, err := db.UpdateRepoTransportWithStatus(repoPath, sshURL, transport)

	return err
}

// GetRepoPrimaryURLByStrategy queries the primary URL and transport using a typed lookup strategy.
func (db *DB) GetRepoPrimaryURLByStrategy(strategy constants.RepoTransportQueryType, target string) (int64, string, string, error) {
	if strategy == constants.RepoTransportQueryByPath {
		cleanPath := NormalizeStoragePath(target)

		id, url, err := db.GetRepoPrimaryURLByPath(cleanPath)

		// The path lookup already selects the URL using the record's
		// identified transport; it does not resolve transport separately.
		return id, url, "", err
	}

	return 0, "", "", appfault.NewValidationError(fmt.Sprintf("unsupported transport query strategy: %s", strategy))
}

func cleanTransportOrDefault(transport string) string {
	clean := strings.ToLower(strings.TrimSpace(transport))
	if clean == "" {
		return constants.TransportSSH
	}

	return clean
}
