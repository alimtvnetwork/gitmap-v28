// Package cmd — prompt_defaults.go seeds built-in prompt templates on first run.
package cmd

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func getPromptsStorageDir() (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "user home dir")
	}

	dir := filepath.Join(homeDir, ".gemini", "config", "prompts")
	if mkErr := os.MkdirAll(dir, constants.DirPermission); mkErr != nil {
		return "", apperror.WrapSimple(mkErr, "mkdir prompts")
	}

	return dir, nil
}

const defaultCodeReviewPrompt = `---
title: "Code Review & Hygiene"
slug: "code-review"
version: "1.0.0"
description: "Autonomous code review and hygiene audit across repository files"
tags: ["review", "hygiene", "audit"]
---

# Autonomous Code Review & Hygiene
1. Audit all files for coding style and error handling.
2. Flatten nested if conditionals and enforce vertical line-gaps.
3. Validate strict typing and run automated linters.
`

const defaultCICDFixPrompt = `---
title: "CI/CD Fix Self-Loop"
slug: "ci-cd-fix"
version: "1.0.0"
description: "Diagnose and repair failed CI/CD pipeline runs using 4-part RCA"
tags: ["cicd", "rca", "automation"]
---

# CI/CD Fix Protocol
1. Locate failed CI/CD pipeline step and capture terminal error logs.
2. Formulate 4-part RCA (Symptoms, Root Cause, Fix Applied, Prevention).
3. Apply surgical fix, verify locally via runner script, and push.
`

func seedDefaultPromptTemplates() {
	dir, err := getPromptsStorageDir()
	if err != nil {
		return
	}

	reviewFile := filepath.Join(dir, "code-review.md")
	if _, statErr := os.Stat(reviewFile); os.IsNotExist(statErr) {
		_ = os.WriteFile(reviewFile, []byte(defaultCodeReviewPrompt), constants.FilePermission)
	}

	cicdFile := filepath.Join(dir, "ci-cd-fix.md")
	if _, statErr := os.Stat(cicdFile); os.IsNotExist(statErr) {
		_ = os.WriteFile(cicdFile, []byte(defaultCICDFixPrompt), constants.FilePermission)
	}
}

func init() {
	seedDefaultPromptTemplates()
}
