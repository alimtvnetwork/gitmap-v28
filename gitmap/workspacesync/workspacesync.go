package workspacesync

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

type AntigravityConfig struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	ProjectResources ProjectResources `json:"projectResources"`
	Settings         map[string]any   `json:"settings"`
	UpdatedAt        string           `json:"updatedAt"`
	IsWorkspaceOnly  bool             `json:"isWorkspaceOnly"`
}

type ProjectResources struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	GitFolder GitFolder `json:"gitFolder"`
}

type GitFolder struct {
	FolderURI     string `json:"folderUri"`
	DefaultBranch string `json:"defaultBranch"`
}

func SyncAll(repoPath string, repoName string) {
	fmt.Printf("  " + constants.ColorDim + "→ sync:" + constants.ColorReset)

	pmRes := "[vsc: skipped]"
	pairs := []vscodepm.Pair{{
		RootPath: repoPath,
		Name:     repoName,
		Paths:    []string{repoPath},
		Tags:     []string{"gitmap"},
	}}
	_, err := vscodepm.SyncMode(pairs, vscodepm.MergeModeUnion)
	if err == nil {
		pmRes = "[vsc: " + constants.ColorGreen + "ok" + constants.ColorReset + "]"
	}

	dtRes := "[desktop: skipped]"
	records := []model.ScanRecord{{AbsolutePath: repoPath}}
	dtSummary := desktop.AddRepos(records)
	if dtSummary.Added > 0 {
		dtRes = "[desktop: " + constants.ColorGreen + "ok" + constants.ColorReset + "]"
	}

	agyRes := "[agy: skipped]"
	if SyncAntigravity(repoPath, repoName) {
		agyRes = "[agy: " + constants.ColorGreen + "ok" + constants.ColorReset + "]"
	}

	fmt.Printf(" %s %s %s\n", pmRes, dtRes, agyRes)
}

func SyncAntigravity(repoPath, repoName string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configDir := filepath.Join(home, ".gemini", "config", "projects")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return false
	}

	uriPath := filepath.ToSlash(repoPath)
	fileURI := "file:///" + url.PathEscape(uriPath)
	if len(uriPath) > 1 && uriPath[1] == ':' {
		fileURI = "file:///" + string(uriPath[0]) + "%3A" + url.PathEscape(uriPath[2:])
	} else if len(uriPath) > 0 && uriPath[0] == '/' {
		fileURI = "file://" + url.PathEscape(uriPath)
	}

	uuid := findExistingProjectID(configDir, fileURI)
	if uuid != "" {
		return true
	}
	uuid = generateUUID()
	nowStr := time.Now().UTC().Format(time.RFC3339Nano)

	cfg := AntigravityConfig{
		ID:   uuid,
		Name: repoName,
		ProjectResources: ProjectResources{
			Resources: []Resource{
				{
					GitFolder: GitFolder{
						FolderURI:     fileURI,
						DefaultBranch: "main",
					},
				},
			},
		},
		Settings:        make(map[string]any),
		UpdatedAt:       nowStr,
		IsWorkspaceOnly: false,
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return false
	}

	outPath := filepath.Join(configDir, uuid+".json")
	if err := os.WriteFile(outPath, b, 0644); err != nil {
		return false
	}

	return true
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	dst := make([]byte, 36)
	hex.Encode(dst, b[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], b[10:])
	return string(dst)
}

func findExistingProjectID(configDir, fileURI string) string {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		if id := checkEntryURI(configDir, entry.Name(), fileURI); id != "" {
			return id
		}
	}
	return ""
}

func checkEntryURI(configDir, fileName, fileURI string) string {
	data, readErr := os.ReadFile(filepath.Join(configDir, fileName))
	if readErr != nil {
		return ""
	}
	var cfg AntigravityConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	if len(cfg.ProjectResources.Resources) > 0 && strings.EqualFold(cfg.ProjectResources.Resources[0].GitFolder.FolderURI, fileURI) {
		return cfg.ID
	}
	return ""
}
