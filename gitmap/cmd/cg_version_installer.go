// Package cmd — cg_version_installer.go creates canonical version.json and documentation.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

const versioningDocMarkdown = `# Canonical Repository Versioning Specification (Single Source of Truth)

## 1. Overview & Principle
This repository adheres to a strict **Single Source of Truth** versioning protocol:
- The canonical version of record is located exclusively in **version.json** at the repository root.
- **Sub-Component Inheritance**: Sub-components (e.g. backend, frontend, CLI) can declare "version": "inherit" to dynamically inherit the top-level repository version, or specify their own explicit version.
- **Every sub-codebase must import/read version.json** rather than hardcoding detached version strings.

## 2. Manifest Schema
- Schema Description: Canonical Single Source of Truth for repository versioning.
- Root version defines the global release version.
- Sub-components specify "inherit" to use top-level version or define an explicit SemVer string.
`

const versioningMemoryMarkdown = `# Learned: Versioning Single Source of Truth (SSOT)

- Always read version.json at the repository root as the sole source of truth.
- Sub-components declare "inherit" to use the root version.
- Never parse versions from test mocks or commit messages.
`

// InstallVersionJSON initializes version.json, memory documentation, and enqueues in what-to-read.md.
func InstallVersionJSON(repoPath string, cfg VersionInstallConfig, isDryRun bool) error {
	vPath := filepath.Join(repoPath, "version.json")
	if isDryRun {
		fmt.Printf("  [DRY-RUN] Would install version.json (v%s) in: %s\n", cfg.InitialVersion, repoPath)
		return nil
	}

	manifest := buildInitialManifest(vPath, cfg.InitialVersion)
	data, errMarshal := json.MarshalIndent(manifest, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	if errWrite := os.WriteFile(vPath, append(data, '\n'), 0644); errWrite != nil {
		return errWrite
	}

	writeVersioningDocsAndMemory(repoPath, cfg)
	EnqueueVersioningDocs(repoPath, cfg)
	return nil
}

func buildInitialManifest(vPath string, initialVersion string) model.RepositoryVersionManifest {
	var manifest model.RepositoryVersionManifest
	manifest.SchemaDescription = "Canonical Single Source of Truth for repository versioning. Every tool, script, and AI agent must read and update this file exclusively."
	manifest.Documentation = "docs/versioning.md"
	manifest.InheritanceRules = "Sub-components specify 'inherit' to use top-level version or define an explicit SemVer string."
	manifest.Version = initialVersion

	manifest.Backend = &model.ComponentVersion{
		Version:     "inherit",
		Status:      "active",
		Description: "Backend service inheriting root version",
	}
	manifest.Frontend = &model.ComponentVersion{
		Version:     "inherit",
		Status:      "active",
		Description: "Frontend application inheriting root version",
	}

	data, err := os.ReadFile(vPath)
	if err == nil {
		_ = json.Unmarshal(data, &manifest)
	}
	if manifest.Version == "" {
		manifest.Version = initialVersion
	}
	return manifest
}

func writeVersioningDocsAndMemory(repoPath string, cfg VersionInstallConfig) {
	docPath := filepath.Join(repoPath, cfg.DocDir, "versioning.md")
	_ = os.MkdirAll(filepath.Dir(docPath), 0755)
	_ = os.WriteFile(docPath, []byte(versioningDocMarkdown), 0644)

	memPath := filepath.Join(repoPath, cfg.MemoryDir, "01-versioning-ssot.md")
	_ = os.MkdirAll(filepath.Dir(memPath), 0755)
	_ = os.WriteFile(memPath, []byte(versioningMemoryMarkdown), 0644)
}
