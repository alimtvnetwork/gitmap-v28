package cmdssh

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ClusterConfigJSON defines the JSON schema mapping to 01-config.json.
type ClusterConfigJSON struct {
	Control map[string]string `json:"control"`
	Nodes   map[string]string `json:"nodes"`
	User    ClusterUserJSON   `json:"user"`
}

// ClusterUserJSON defines credentials for cluster nodes.
type ClusterUserJSON struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func checkFileStat(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return apperror.WrapNotFound(err, fmt.Sprintf("config file not found: %s", path))
	}
	if info.IsDir() {
		return apperror.NewValidationError(fmt.Sprintf("config path is a directory: %s", path))
	}
	return nil
}

func readConfigFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "readConfigFile")
	}
	return data, nil
}

func parseClusterConfigBytes(data []byte) (*ClusterConfigJSON, error) {
	hasContent := len(strings.TrimSpace(string(data))) > 0
	if !hasContent {
		return nil, apperror.NewValidationError("cluster config file is empty")
	}
	var cfg ClusterConfigJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, apperror.NewValidationError(fmt.Sprintf("invalid cluster config json: %v", err))
	}
	return &cfg, nil
}

// ValidateClusterUser validates username and password presence.
func ValidateClusterUser(user ClusterUserJSON) error {
	hasName := len(strings.TrimSpace(user.Name)) > 0
	if !hasName {
		return apperror.NewValidationError("user.name is required in cluster config")
	}
	hasPassword := len(strings.TrimSpace(user.Password)) > 0
	if !hasPassword {
		return apperror.NewValidationError("user.password is required in cluster config")
	}
	return nil
}

// LoadClusterConfigFile loads and validates a cluster config JSON file.
func LoadClusterConfigFile(path string) (*ClusterConfigJSON, error) {
	if err := checkFileStat(path); err != nil {
		return nil, err
	}
	data, err := readConfigFile(path)
	if err != nil {
		return nil, err
	}
	return parseClusterConfigBytes(data)
}
