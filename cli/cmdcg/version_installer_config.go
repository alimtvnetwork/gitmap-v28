// Package cmd — version_installer_config.go defines configurable paths for version.json installation.
package cmdcg

// VersionInstallConfig specifies destination paths for version installation.
type VersionInstallConfig struct {
	DocDir            string
	MemoryDir         string
	WhatToReadDoc     string
	RootWhatToReadDoc string
	InitialVersion    string
}

// DefaultVersionInstallConfig returns default configuration pointing to .ai-memory directory.
func DefaultVersionInstallConfig(initialVersion string) VersionInstallConfig {
	if initialVersion == "" {
		initialVersion = "1.0.0"
	}

	return VersionInstallConfig{
		DocDir:            ".ai-memory",
		MemoryDir:         ".ai-memory/memory/learned",
		WhatToReadDoc:     ".ai-memory/what-to-read.md",
		RootWhatToReadDoc: "what-to-read.md",
		InitialVersion:    initialVersion,
	}
}
