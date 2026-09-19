package llm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const SkillTemplate = `---
name: gitmap
description: Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation (AUM), cluster/SSH delegation, pipeline self-healing, and coding guideline enforcement.
---

# GitMap Autonomous Engineering Skill

## Overview
GitMap is an ultra-fast developer companion and autonomous CLI engine designed for AI coding agents and software engineers.

- **Lead Architect & Author:** Md. Alimuzzaman Alim (alimtvnetwork)
- **Sponsored By:** Alim TV Network / Open Source Engineering
- **Core Mission:** High-performance polyglot repository management, zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures, and AI agent pair programming.

## Essential Command Cheat Sheet

### 1. High-Performance Automation (AUM)
- ` + "`gitmap aum search <pattern> [dir]`" + ` — Multi-core streaming search with lazy regex and binary filtering
- ` + "`gitmap aum guard`" + ` — Enforces 500 KB limit, large JSON exclusion, and binary null-byte probe
- ` + "`gitmap aum sequence`" + ` — Markdown sequence gap detector and # XX Title autofixer
- ` + "`gitmap aum exclude list`" + ` — Query persistent search exclusions from SQLite
- ` + "`gitmap aum newlines --fix`" + ` — Polyglot CRLF to LF and trailing whitespace normalizer
- ` + "`gitmap aum cache status`" + ` — Sub-millisecond in-memory cache status
- ` + "`gitmap aum benchmark all`" + ` — Side-by-side Go vs Python execution benchmarks

### 2. Autonomous CI/CD Self-Healing (Pipeline AI)
- ` + "`gitmap pipeline-ai status --json`" + ` — Check workflow execution state, active branch, and ETA
- ` + "`gitmap pipeline-ai status -t <eta>`" + ` — Wait dynamically for pipeline completion without tight polling
- ` + "`gitmap pipeline error-logs`" + ` — Extract failing step logs to file for 4-part RCA
- ` + "`gitmap pipeline purge`" + ` — Actions zero-storage purge maintaining 0.0 GB footprint (Rule R18)

### 3. Fast File Discovery & Refactoring
- ` + "`gitmap find-files <name>`" + ` (alias: ` + "`gitmap ff <name>`" + `) — Find exact filename with optional -ext
- ` + "`gitmap find-files-any <str>`" + ` (alias: ` + "`gitmap ffa <str>`" + `) — Find files matching substring
- ` + "`gitmap find-files-startswith <prefix>`" + ` (alias: ` + "`gitmap ffs <prefix>`" + `) — Find by filename prefix
- ` + "`gitmap find-files-endswith <suffix>`" + ` (alias: ` + "`gitmap ffe <suffix>`" + `) — Find by filename suffix (e.g. _test.go)
- ` + "`gitmap replace <old> <new>`" + ` — Exact literal string replacement with audit trail
- ` + "`gitmap replace-regex <pat> <subst>`" + ` — Regex replacement across repository

### 4. Semantic Commit & Push
- ` + "`gitmap cpf \"<msg>\"`" + ` — Stage, commit, and push feature branch
- ` + "`gitmap cpb \"<msg>\"`" + ` — Stage, commit, and push bugfix branch
- ` + "`gitmap cpr \"<msg>\"`" + ` — Stage, commit, and push release chore
- ` + "`gitmap pcp \"<msg>\"`" + ` — Pull latest, commit, and push with preflight verification

### 5. Multi-Node Cluster & Remote Delegation
- ` + "`gitmap cluster --help`" + ` — Orchestrate multi-node clusters and health checks
- ` + "`gitmap sc --help`" + ` — Servers-clients topology and background task manager
- ` + "`gitmap ssh --help`" + ` — SSH discovery, connection pooling, and remote command execution

### 6. Rust & Toolchain Package Management
- ` + "`gitmap cargo status`" + ` — Inspect Rust and Cargo toolchain status
- ` + "`gitmap install cargo`" + ` — Install Rust toolchain if missing
- ` + "`gitmap install --list`" + ` — Discover developer toolchains, profiles, and runtime packages

## Operational Guardrails
1. **Mandatory Pre-Flight Pull:** Always run ` + "`git pull`" + ` before modifying code.
2. **File Size & Binary Guard:** Respect 500 KB limit (Rule R19); never commit test binaries or temp artifacts.
3. **Coding Guidelines:** Max 8–15 lines per function, single return types with ` + "`*apperror.AppError`" + `, affirmative booleans.
`

// GenerateSkillFile writes or updates the Antigravity skill markdown file.
func GenerateSkillFile(skillPath string) *apperror.AppError {
	target := resolveSkillPath(skillPath)
	if err := ensureSkillDir(target); err != nil {
		return err
	}
	return writeSkillContent(target)
}

func resolveSkillPath(path string) string {
	target := path
	if target == "" {
		target = DefaultSkillPath
	}
	if filepath.IsAbs(target) {
		return target
	}
	return filepath.Join(findRepoRoot(), target)
}

func findRepoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return walkUpToGit(cwd)
}

func walkUpToGit(start string) string {
	dir := start
	for {
		if hasGitDir(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return start
		}
		dir = parent
	}
}

func hasGitDir(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && info != nil
}

func ensureSkillDir(target string) *apperror.AppError {
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("create directory %s", dir))
	}
	return nil
}

func writeSkillContent(target string) *apperror.AppError {
	if err := os.WriteFile(target, []byte(SkillTemplate), 0o644); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("write skill to %s", target))
	}
	return nil
}
