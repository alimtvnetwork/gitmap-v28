// Package cmd — cg_version.go manages coding-guidelines metadata inside version.json.
package cmdcg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	cgStatusHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	cgActiveBadgeStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	cgMissingBadgeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
	cgDateStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
)

// CGMetadata represents the coding guidelines section inside version.json.
type CGMetadata struct {
	Version     string `json:"version"`
	InstalledAt string `json:"installed_at,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ReadCGMetadata reads the coding-guidelines section from version.json in a repo.
func ReadCGMetadata(repoPath string) (*CGMetadata, error) {
	raw, errRead := loadVersionJSON(repoPath)
	if errRead != nil {
		return nil, errRead
	}

	meta := extractCGMetadata(raw)
	if meta != nil {
		return meta, nil
	}

	return nil, apperror.New("ReadCGMetadata", "E_CG_NOT_INSTALLED", map[string]any{"path": repoPath})
}

func loadVersionJSON(repoPath string) (map[string]json.RawMessage, error) {
	vPath := filepath.Join(repoPath, "version.json")
	data, errRead := os.ReadFile(vPath)
	if errRead != nil {
		return nil, errRead
	}

	var raw map[string]json.RawMessage
	errUnmarshal := json.Unmarshal(data, &raw)

	return raw, errUnmarshal
}

func extractCGMetadata(raw map[string]json.RawMessage) *CGMetadata {
	keys := []string{"coding-guidelines", "coding_guidelines", "codingguideline"}
	for _, key := range keys {
		if meta := unmarshalCGKey(raw, key); meta != nil {
			return meta
		}
	}

	return nil
}

func unmarshalCGKey(raw map[string]json.RawMessage, key string) *CGMetadata {
	rawMeta, hasKey := raw[key]
	if !hasKey {
		return nil
	}

	var meta CGMetadata
	if err := json.Unmarshal(rawMeta, &meta); err == nil {
		return &meta
	}

	return nil
}

// WriteCGMetadata updates or inserts the coding-guidelines section inside version.json.
func WriteCGMetadata(repoPath string, meta CGMetadata) error {
	vPath := filepath.Join(repoPath, "version.json")
	raw := readRawVersionMap(vPath)

	populateDefaultCGMetadata(&meta)
	raw["coding-guidelines"] = meta
	data, errMarshal := json.MarshalIndent(raw, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(vPath, data, 0644)
}

func readRawVersionMap(vPath string) map[string]any {
	var raw map[string]any
	if data, errRead := os.ReadFile(vPath); errRead == nil {
		_ = json.Unmarshal(data, &raw)
	}

	if raw == nil {
		raw = make(map[string]any)
	}

	return raw
}

func populateDefaultCGMetadata(meta *CGMetadata) {
	if meta.InstalledAt == "" {
		meta.InstalledAt = time.Now().UTC().Format(time.RFC3339)
	}

	if meta.Status == "" {
		meta.Status = "active"
	}
}

// PrintCGStatus prints a formatted status table of coding guidelines across repos.
func PrintCGStatus(repos []string) {
	fmt.Println(cgStatusHeaderStyle.Render("📋 Coding Guidelines Multi-Repo Status:"))

	for _, repo := range repos {
		printSingleRepoStatus(repo)
	}
}

func printSingleRepoStatus(repo string) {
	name := filepath.Base(repo)
	meta, err := ReadCGMetadata(repo)
	if err != nil {
		fmt.Println(renderStatusMissingRow(name))

		return
	}

	fmt.Println(renderStatusActiveRow(name, meta))
}

func renderStatusActiveRow(name string, meta *CGMetadata) string {
	badge := cgActiveBadgeStyle.Width(12).Render("✔ active")
	ver := cgVersionStyle.Width(12).Render(meta.Version)
	date := cgDateStyle.Render(meta.InstalledAt)

	return fmt.Sprintf("  %-24s %s %s %s", name, badge, ver, date)
}

func renderStatusMissingRow(name string) string {
	badge := cgMissingBadgeStyle.Width(12).Render("✖ missing")
	ver := cgDateStyle.Width(12).Render("-")
	date := cgDateStyle.Render("-")

	return fmt.Sprintf("  %-24s %s %s %s", name, badge, ver, date)
}

// PrintCGVersion prints the coding guidelines version for repos.
func PrintCGVersion(repos []string) {
	for _, repo := range repos {
		printSingleRepoVersion(repo)
	}
}

func printSingleRepoVersion(repo string) {
	name := filepath.Base(repo)
	meta, err := ReadCGMetadata(repo)
	if err != nil {
		fmt.Printf("%s: %s\n", name, cgMissingBadgeStyle.Render("not installed"))

		return
	}

	fmt.Printf("%s: %s\n", name, cgVersionStyle.Render(strings.TrimSpace(meta.Version)))
}
