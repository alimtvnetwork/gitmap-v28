package cmdssh

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type vmPassCredentials struct {
	Windows struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"windows"`
	Ubuntu struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"ubuntu"`
}

func candidateVMPassPaths() []string {
	paths := []string{
		"vmpass.json",
		"../vmpass.json",
		"../../vmpass.json",
		filepath.Join(filepath.Dir(store.BinaryDataDir()), "vmpass.json"),
		"D:/work/repo-secrets/01-gitmap/vmpass.json",
		filepath.Join("..", "repo-secrets", "01-gitmap", "vmpass.json"),
		filepath.Join("..", "..", "repo-secrets", "01-gitmap", "vmpass.json"),
		filepath.Join("repo-secrets", "01-gitmap", "vmpass.json"),
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, ".gitmap", "vmpass.json"),
			filepath.Join(home, "vmpass.json"),
		)
	}
	return paths
}

func findVMPassFile() string {
	candidates := candidateVMPassPaths()
	for _, p := range candidates {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}
	return ""
}

func loadVMPassCreds() *vmPassCredentials {
	filePath := findVMPassFile()
	if filePath == "" {
		return nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	var creds vmPassCredentials
	if jsonErr := json.Unmarshal(data, &creds); jsonErr != nil {
		return nil
	}
	return &creds
}

func isWindowsTargetUserOrOS(user, osType string) bool {
	lowUser := strings.ToLower(user)
	if lowUser == "administrator" || lowUser == "admin" {
		return true
	}
	return isWindowsOS(osType)
}

// ResolveFallbackCredentials retrieves default cluster passwords from vmpass.json.
func ResolveFallbackCredentials(user, osType string) string {
	creds := loadVMPassCreds()
	if creds == nil {
		return ""
	}
	if isWindowsTargetUserOrOS(user, osType) && creds.Windows.Pass != "" {
		return creds.Windows.Pass
	}
	if creds.Ubuntu.Pass != "" {
		return creds.Ubuntu.Pass
	}
	return creds.Windows.Pass
}
