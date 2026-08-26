// Package cmd — cg_version.go manages coding-guidelines metadata inside version.json.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/charmbracelet/lipgloss"
)

// CGMetadata represents the coding guidelines section inside version.json.
type CGMetadata struct {
	Version     string `json:"version"`
	InstalledAt string `json:"installed_at,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ReadCGMetadata reads the coding-guidelines section from version.json in a repo.
func ReadCGMetadata(repoPath string) (*CGMetadata, error) {
	vPath := filepath.Join(repoPath, "version.json")
	data, errRead := os.ReadFile(vPath)
	if errRead != nil {
		return nil, errRead
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	for _, key := range []string{"coding-guidelines", "coding_guidelines", "codingguideline"} {
		if rawMeta, ok := raw[key]; ok {
			var meta CGMetadata
			if err := json.Unmarshal(rawMeta, &meta); err == nil {
				return &meta, nil
			}
		}
	}

	return nil, apperror.New("ReadCGMetadata", "E_CG_NOT_INSTALLED", map[string]any{"path": repoPath})
}

// WriteCGMetadata updates or inserts the coding-guidelines section inside version.json.
func WriteCGMetadata(repoPath string, meta CGMetadata) error {
	vPath := filepath.Join(repoPath, "version.json")
	var raw map[string]any

	if data, errRead := os.ReadFile(vPath); errRead == nil {
		_ = json.Unmarshal(data, &raw)
	}
	if raw == nil {
		raw = make(map[string]any)
	}

	if meta.InstalledAt == "" {
		meta.InstalledAt = time.Now().UTC().Format(time.RFC3339)
	}
	if meta.Status == "" {
		meta.Status = "active"
	}

	raw["coding-guidelines"] = meta
	data, errMarshal := json.MarshalIndent(raw, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(vPath, data, 0644)
}

// PrintCGStatus prints a formatted status table of coding guidelines across repos.
func PrintCGStatus(repos []string) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	fmt.Println(titleStyle.Render("Coding Guidelines Status:"))

	for _, repo := range repos {
		name := filepath.Base(repo)
		meta, err := ReadCGMetadata(repo)
		if err != nil {
			fmt.Printf("  %-25s : Not Installed\n", name)
		} else {
			fmt.Printf("  %-25s : Installed (%s) - %s [%s]\n", name, meta.Version, meta.InstalledAt, meta.Status)
		}
	}
}

// PrintCGVersion prints the coding guidelines version for repos.
func PrintCGVersion(repos []string) {
	for _, repo := range repos {
		name := filepath.Base(repo)
		meta, err := ReadCGMetadata(repo)
		if err != nil {
			fmt.Printf("%s: not installed\n", name)
		} else {
			fmt.Printf("%s: %s\n", name, strings.TrimSpace(meta.Version))
		}
	}
}
