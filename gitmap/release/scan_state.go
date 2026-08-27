package release

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type scanState struct {
	LastCommit string `json:"last_commit"`
}

func ReadLastScannedCommit(repoDir string) (string, error) {
	path := filepath.Join(repoDir, ".gitmap", "commit_scan_state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", apperror.Wrap(err, "read_scan_state", map[string]any{"path": path})
	}
	var state scanState
	if err := json.Unmarshal(data, &state); err != nil {
		return "", apperror.Wrap(err, "parse_scan_state", map[string]any{"path": path})
	}
	return state.LastCommit, nil
}

func WriteLastScannedCommit(repoDir, commitHash string) error {
	path := filepath.Join(repoDir, ".gitmap", "commit_scan_state.json")
	state := scanState{LastCommit: commitHash}
	data, _ := json.Marshal(state)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return apperror.Wrap(err, "mkdir_scan_state", map[string]any{"path": path})
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return apperror.Wrap(err, "write_scan_state", map[string]any{"path": path})
	}
	return nil
}
