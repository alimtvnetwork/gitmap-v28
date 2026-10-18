package cmdclone

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ParseManifestFile parses a JSON manifest from disk into normalized RepoCacheEntry objects.
func ParseManifestFile(manifestPath string) ([]RepoCacheEntry, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	rawItems, parseErr := parseRawManifestItems(data, manifestPath)
	if parseErr != nil {
		return nil, parseErr
	}

	entries := make([]RepoCacheEntry, 0, len(rawItems))
	for _, item := range rawItems {
		cloneURL := resolveItemCloneURL(item)
		if len(cloneURL) == 0 {
			continue
		}
		owner, repo := FormatOwnerRepo(cloneURL)
		switch {
		case len(item.RepoName) > 0:
			repo = item.RepoName
		case len(item.Name) > 0:
			repo = item.Name
		}
		if len(item.Owner) > 0 {
			owner = item.Owner
		}
		isSSH := isSSHProtocol(cloneURL)
		isHTTPS := isHTTPSProtocol(cloneURL)
		entries = append(entries, RepoCacheEntry{
			RepoName:     repo,
			Owner:        owner,
			CloneUrl:     cloneURL,
			IsSSH:        isSSH,
			IsHTTPS:      isHTTPS,
			Branch:       item.Branch,
			RelativePath: item.RelativePath,
		})
	}

	return entries, nil
}

func resolveItemCloneURL(item genericManifestItem) string {
	if len(item.URL) > 0 {
		return item.URL
	}
	if len(item.CloneUrl) > 0 {
		return item.CloneUrl
	}
	if len(item.HTTPSUrl) > 0 {
		return item.HTTPSUrl
	}
	if len(item.SSHUrl) > 0 {
		return item.SSHUrl
	}
	if len(item.DiscoveredURL) > 0 {
		return item.DiscoveredURL
	}
	return item.RemoteURL
}

func isSSHProtocol(url string) bool {
	low := strings.ToLower(url)
	return strings.HasPrefix(low, "git@") || strings.HasPrefix(low, "ssh://")
}

func isHTTPSProtocol(url string) bool {
	low := strings.ToLower(url)
	return strings.HasPrefix(low, "https://") || strings.HasPrefix(low, "http://")
}

// GetOrSyncManifestEntries gets cached entries if valid, or reads disk, syncs to DB, and returns them.
func GetOrSyncManifestEntries(db *sql.DB, manifestPath string) ([]RepoCacheEntry, error) {
	fileHash, fileMtime, err := ComputeFileHashAndMtime(manifestPath)
	if err != nil {
		return nil, err
	}

	entries, isHit, err := GetCachedManifestEntries(db, manifestPath, fileHash, fileMtime)
	if err == nil && isHit {
		return entries, nil
	}

	parsedEntries, parseErr := ParseManifestFile(manifestPath)
	if parseErr != nil {
		return nil, parseErr
	}

	_ = StoreManifestEntries(db, manifestPath, fileHash, fileMtime, parsedEntries)
	return parsedEntries, nil
}

func parseRawManifestItems(data []byte, manifestPath string) ([]genericManifestItem, error) {
	var rawItems []genericManifestItem
	if err := json.Unmarshal(data, &rawItems); err == nil {
		return rawItems, nil
	}

	var wrapper genericManifestWrapper
	if wrapErr := json.Unmarshal(data, &wrapper); wrapErr != nil {
		return nil, fmt.Errorf("unable to parse manifest JSON %s: %w", manifestPath, wrapErr)
	}

	switch {
	case len(wrapper.Records) > 0:
		return wrapper.Records, nil
	case len(wrapper.Repositories) > 0:
		return wrapper.Repositories, nil
	case len(wrapper.Items) > 0:
		return wrapper.Items, nil
	case len(wrapper.Repos) > 0:
		return wrapper.Repos, nil
	default:
		return nil, nil
	}
}

func upsertManifestRecord(tx *sql.Tx, manifestPath, fileHash string, fileMtime int64, totalRepos int, now string) (int64, error) {
	upsertManifestQuery := `
	INSERT INTO RepoCacheManifest (ManifestPath, FileHash, FileMtime, TotalRepos, HasValidData, CreatedAt, UpdatedAt)
	VALUES (?, ?, ?, ?, 1, ?, ?)
	ON CONFLICT(ManifestPath) DO UPDATE SET
		FileHash = excluded.FileHash,
		FileMtime = excluded.FileMtime,
		TotalRepos = excluded.TotalRepos,
		HasValidData = 1,
		UpdatedAt = excluded.UpdatedAt
	RETURNING RepoCacheManifestId;
	`
	var manifestId int64
	err := tx.QueryRow(upsertManifestQuery, manifestPath, fileHash, fileMtime, totalRepos, now, now).Scan(&manifestId)
	if err == nil {
		return manifestId, nil
	}

	return executeFallbackManifestUpsert(tx, manifestPath, fileHash, fileMtime, totalRepos, now)
}

func executeFallbackManifestUpsert(tx *sql.Tx, manifestPath, fileHash string, fileMtime int64, totalRepos int, now string) (int64, error) {
	fallbackUpsert := `
	INSERT INTO RepoCacheManifest (ManifestPath, FileHash, FileMtime, TotalRepos, HasValidData, CreatedAt, UpdatedAt)
	VALUES (?, ?, ?, ?, 1, ?, ?)
	ON CONFLICT(ManifestPath) DO UPDATE SET
		FileHash = excluded.FileHash,
		FileMtime = excluded.FileMtime,
		TotalRepos = excluded.TotalRepos,
		HasValidData = 1,
		UpdatedAt = excluded.UpdatedAt;
	`
	if _, execErr := tx.Exec(fallbackUpsert, manifestPath, fileHash, fileMtime, totalRepos, now, now); execErr != nil {
		return 0, execErr
	}

	var manifestId int64
	if scanErr := tx.QueryRow(`SELECT RepoCacheManifestId FROM RepoCacheManifest WHERE ManifestPath = ?`, manifestPath).Scan(&manifestId); scanErr != nil {
		return 0, scanErr
	}

	return manifestId, nil
}
