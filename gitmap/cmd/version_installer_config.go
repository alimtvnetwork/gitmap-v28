// Package cmd — version_installer_config.go defines configurable paths for version.json installation.
package cmd

// VersionInstallConfig specifies destination paths for version installation.
type VersionInstallConfig struct {
	DocDir            string
	MemoryDir         string
	WhatToReadDoc     string
	RootWhatToReadDoc string
	InitialVersion    string
}

// DefaultVersionInstallConfig returns default configuration pointing to .lovable directory.
func DefaultVersionInstallConfig(initialVersion string) VersionInstallConfig {
	if initialVersion == "" {
		initialVersion = "1.0.0"
	}
	return VersionInstallConfig{
		DocDir:            ".lovable",
		MemoryDir:         ".lovable/memory/learned",
		WhatToReadDoc:     ".lovable/what-to-read.md",
		RootWhatToReadDoc: "what-to-read.md",
		InitialVersion:    initialVersion,
	}
}
