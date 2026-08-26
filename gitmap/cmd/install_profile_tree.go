// Package cmd — install_profile_tree.go manages profile composition hierarchies and tree rendering.
package cmd

import (
	"strings"
)

// ToolEntry represents an individual tool included in a profile composition.
type ToolEntry struct {
	Slug        string
	Description string
}

// ProfileComposition represents a profile configuration and its inheritance chain.
type ProfileComposition struct {
	Name        string
	Alias       string
	Description string
	Base        *ProfileComposition
	Tools       []ToolEntry
}

// buildUbuntuBasicProfile constructs the base Ubuntu toolchain profile.
func buildUbuntuBasicProfile() ProfileComposition {
	return ProfileComposition{
		Name:        "ubuntu-basic",
		Alias:       "ub",
		Description: "Basic Ubuntu toolchain and essential utilities",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "curl", Description: "Command line tool for transferring data with URLs"},
			{Slug: "git", Description: "Fast, scalable, distributed revision control system"},
			{Slug: "build-essential", Description: "Meta-package for compiling software (gcc, make)"},
			{Slug: "zsh", Description: "Z shell with modern command line features"},
		},
	}
}

// buildUbuntuVscodeProfile constructs the Ubuntu VSCode profile composition.
func buildUbuntuVscodeProfile() ProfileComposition {
	baseProfile := buildUbuntuBasicProfile()
	return ProfileComposition{
		Name:        "ubuntu+vscode",
		Alias:       "ub+code",
		Description: "Ubuntu basic plus Visual Studio Code editor",
		Base:        &baseProfile,
		Tools: []ToolEntry{
			{Slug: "code", Description: "Visual Studio Code editor binaries and extensions"},
		},
	}
}

// buildUbuntuSmallDevProfile constructs the Ubuntu small dev profile composition.
func buildUbuntuSmallDevProfile() ProfileComposition {
	baseProfile := buildUbuntuVscodeProfile()
	return ProfileComposition{
		Name:        "ubuntu+small-dev",
		Alias:       "ub+sdev",
		Description: "Ubuntu lightweight development suite with runtime interpreters",
		Base:        &baseProfile,
		Tools: []ToolEntry{
			{Slug: "golang", Description: "Go programming language runtime and tools"},
			{Slug: "nodejs", Description: "Node.js JavaScript runtime and npm"},
		},
	}
}

// buildUbuntuDevProfile constructs the full Ubuntu developer profile composition.
func buildUbuntuDevProfile() ProfileComposition {
	baseProfile := buildUbuntuSmallDevProfile()
	return ProfileComposition{
		Name:        "ubuntu+dev",
		Alias:       "ub+dev",
		Description: "Full Ubuntu developer workstation suite",
		Base:        &baseProfile,
		Tools: []ToolEntry{
			{Slug: "docker", Description: "Container runtime and management engine"},
			{Slug: "python3", Description: "Python 3 runtime, pip, and venv tooling"},
		},
	}
}

// resolveProfileTree locates and constructs the composition hierarchy for a given slug.
func resolveProfileTree(slug string) (ProfileComposition, bool) {
	normalizedSlug := strings.ToLower(strings.TrimSpace(slug))
	switch normalizedSlug {
	case "ubuntu-basic", "ubuntu+basic", "ub":
		return buildUbuntuBasicProfile(), true
	case "ubuntu+vscode", "ubuntu-vscode", "ub+code", "vscode":
		return buildUbuntuVscodeProfile(), true
	case "ubuntu+small-dev", "ubuntu-small-dev", "ub+sdev", "small-dev":
		return buildUbuntuSmallDevProfile(), true
	case "ubuntu+dev", "ubuntu-dev", "ub+dev", "dev":
		return buildUbuntuDevProfile(), true
	default:
		return ProfileComposition{}, false
	}
}

// profileToTreeNode converts a profile hierarchy into an installer tree node structure.
func profileToTreeNode(profile ProfileComposition) InstallerTreeNode {
	rootNode := InstallerTreeNode{
		Title:       profile.Name,
		Description: profile.Description,
	}
	if profile.Base != nil {
		rootNode.Children = append(rootNode.Children, profileToTreeNode(*profile.Base))
	}
	return appendToolNodes(rootNode, profile.Tools)
}

// appendToolNodes attaches tool entries to the given parent tree node.
func appendToolNodes(parent InstallerTreeNode, tools []ToolEntry) InstallerTreeNode {
	for _, toolEntry := range tools {
		parent.Children = append(parent.Children, InstallerTreeNode{
			Title:       toolEntry.Slug,
			Description: toolEntry.Description,
		})
	}
	return parent
}

// printProfileTree renders a profile composition tree to stdout.
func printProfileTree(profile ProfileComposition) {
	rootNode := profileToTreeNode(profile)
	printInstallerTree(rootNode, "", true)
}

// printProfileInstallSummary resolves a profile slug and prints its tree summary.
func printProfileInstallSummary(slug string) {
	profile, hasProfile := resolveProfileTree(slug)
	if !hasProfile {
		printInstallSummaryHeader(slug)
		return
	}
	printInstallSummaryHeader(profile.Name)
	printProfileTree(profile)
}
