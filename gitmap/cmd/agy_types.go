// Package cmd — agy_types.go provides data structures and helpers for Antigravity projects.
package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type AgyProject struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	ProjectResources *AgyProjectResources   `json:"projectResources,omitempty"`
	Settings         map[string]interface{} `json:"settings,omitempty"`
	UpdatedAt        string                 `json:"updatedAt,omitempty"`
	IsWorkspaceOnly  bool                   `json:"isWorkspaceOnly,omitempty"`
}

type AgyProjectResources struct {
	Resources []AgyResource `json:"resources,omitempty"`
}

type AgyResource struct {
	GitFolder *AgyGitFolder `json:"gitFolder,omitempty"`
}

type AgyGitFolder struct {
	FolderURI     string `json:"folderUri,omitempty"`
	DefaultBranch string `json:"defaultBranch,omitempty"`
}

func (p *AgyProject) GetPath() string {
	if p.ProjectResources == nil || len(p.ProjectResources.Resources) == 0 {
		return ""
	}
	gf := p.ProjectResources.Resources[0].GitFolder
	if gf == nil {
		return ""
	}
	return parseFolderUri(gf.FolderURI)
}

func (p *AgyProject) GetBranch() string {
	if p.ProjectResources == nil || len(p.ProjectResources.Resources) == 0 {
		return ""
	}
	gf := p.ProjectResources.Resources[0].GitFolder
	if gf == nil {
		return ""
	}
	return gf.DefaultBranch
}

func parseFolderUri(rawUri string) string {
	if rawUri == "" {
		return ""
	}
	cleanUri := strings.TrimPrefix(rawUri, "file:///")
	cleanUri = strings.TrimPrefix(cleanUri, "file://")
	decoded, err := url.PathUnescape(cleanUri)
	if err != nil {
		decoded = cleanUri
	}
	if len(decoded) >= 2 && decoded[1] == ':' {
		return filepath.FromSlash(decoded)
	}
	if strings.HasPrefix(rawUri, "file:///") {
		return filepath.FromSlash("/" + decoded)
	}
	return filepath.FromSlash(decoded)
}

func shortProjectId(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8]
}

func formatRelativeTime(rfc3339Str string) string {
	if rfc3339Str == "" {
		return "—"
	}
	t, err := time.Parse(time.RFC3339Nano, rfc3339Str)
	if err != nil {
		t, err = time.Parse(time.RFC3339, rfc3339Str)
		if err != nil {
			return "—"
		}
	}
	return formatTimeDuration(time.Since(t), t)
}

func formatTimeDuration(d time.Duration, t time.Time) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if days < 30 {
		return fmt.Sprintf("%dd ago", days)
	}
	return t.Format("2006-01-02")
}
