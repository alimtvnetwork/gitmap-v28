// Package formatter — template.go provides shared template loading via go:embed.
package formatter

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// RepoEntry is the data passed into PowerShell templates for each repo.
type RepoEntry struct {
	Name   string
	Branch string
	URL    string
	Path   string
}

// CloneData is the top-level data for clone.ps1.tmpl.
type CloneData struct {
	Repos []RepoEntry
}

// DesktopData is the top-level data for desktop.ps1.tmpl.
type DesktopData struct {
	Repos []RepoEntry
}

// loadTemplate parses an embedded template file by name.
func loadTemplate(name string) (*template.Template, error) {
	content, err := templateFS.ReadFile("templates/" + name)
	if err != nil {
		return nil, err
	}

	return template.New(name).Parse(string(content))
}

// backslashPath converts forward slashes to backslashes for PowerShell paths.
func backslashPath(p string) string {
	return strings.ReplaceAll(p, "/", "\\")
}
