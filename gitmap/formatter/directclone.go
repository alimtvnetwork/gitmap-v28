// Package formatter — directclone.go generates plain direct clone scripts.
package formatter

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// WriteDirectCloneScript writes a plain PS1 with one HTTPS git clone per line.
func WriteDirectCloneScript(w io.Writer, records []model.ScanRecord) error {
	return writeDirectCloneScriptTemplate(w, records, "direct-clone.ps1.tmpl", false)
}

// WriteDirectCloneSSHScript writes a plain PS1 with one SSH git clone per line.
func WriteDirectCloneSSHScript(w io.Writer, records []model.ScanRecord) error {
	return writeDirectCloneScriptTemplate(w, records, "direct-clone-ssh.ps1.tmpl", true)
}

// writeDirectCloneScriptTemplate renders a direct clone template.
func writeDirectCloneScriptTemplate(w io.Writer, records []model.ScanRecord, templateName string, useSSH bool) error {
	tmpl, err := loadTemplate(templateName)
	if err != nil {
		return err
	}

	data := CloneData{
		Repos: buildDirectCloneEntries(records, useSSH),
	}

	return tmpl.Execute(w, data)
}

// buildDirectCloneEntries builds direct clone template entries.
// Per-repo Transport (identified from origin) wins over the scan-wide
// useSSH default — an SSH-origin repo must not be silently rewritten
// to HTTPS, otherwise private remotes trigger browser auth at clone
// time. useSSH still acts as the fallback when Transport is unset.
func buildDirectCloneEntries(records []model.ScanRecord, useSSH bool) []RepoEntry {
	entries := make([]RepoEntry, 0, len(records))
	for _, r := range records {
		url := pickDirectCloneURL(r, useSSH)
		if len(url) == 0 {
			url = cloneURL(r)
		}

		entries = append(entries, RepoEntry{
			Branch: r.Branch,
			URL:    url,
			Path:   backslashPath(r.RelativePath),
		})
	}

	return entries
}

// pickDirectCloneURL honors the per-repo identified transport first,
// then falls back to the scan-wide useSSH preference.
func pickDirectCloneURL(r model.ScanRecord, useSSH bool) string {
	if r.Transport == "ssh" && len(r.SSHUrl) > 0 {
		return r.SSHUrl
	}
	if r.Transport == "https" && len(r.HTTPSUrl) > 0 {
		return r.HTTPSUrl
	}
	if useSSH {
		return r.SSHUrl
	}

	return r.HTTPSUrl
}
