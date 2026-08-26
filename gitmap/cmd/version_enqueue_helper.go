// Package cmd — version_enqueue_helper.go updates what-to-read.md with versioning docs.
package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

const versioningQueueEntry = "- `.lovable/versioning.md`, why: Explains that `version.json` is the canonical Single Source of Truth (SSOT) and documents component inheritance."

// EnqueueVersioningInWhatToRead appends the versioning documentation link to target what-to-read file if missing.
func EnqueueVersioningInWhatToRead(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Optional file
		}
		return err
	}

	content := string(data)
	if strings.Contains(content, "versioning.md") {
		return nil // Already enqueued
	}

	sectionHeader := "## Before writing code"
	if strings.Contains(content, sectionHeader) {
		newContent := strings.Replace(content, sectionHeader, sectionHeader+"\n\n"+versioningQueueEntry, 1)
		return os.WriteFile(filePath, []byte(newContent), 0644)
	}

	// Append at the end if section not found
	content = strings.TrimSpace(content) + "\n\n" + versioningQueueEntry + "\n"
	return os.WriteFile(filePath, []byte(content), 0644)
}

// EnqueueVersioningDocs enqueues versioning documentation in both .lovable and root what-to-read files.
func EnqueueVersioningDocs(repoPath string, cfg VersionInstallConfig) {
	lovableWTR := filepath.Join(repoPath, cfg.WhatToReadDoc)
	_ = EnqueueVersioningInWhatToRead(lovableWTR)

	rootWTR := filepath.Join(repoPath, cfg.RootWhatToReadDoc)
	_ = EnqueueVersioningInWhatToRead(rootWTR)
}
