# Master Plan: Commit Commands Overhaul

## Overview

Overhaul the commit-in, commit-left, and commit-right command trio to support advanced history manipulation, pull request integration, and enhanced terminal rendering.

## Architecture & Scope

1. **Help Text Enhancement**:
   - Update constants_commithelp.go or equivalent to include detailed examples for commit-in, commit-left, commit-right.
   - Add examples showing how to rewrite commit history.
   - Add examples demonstrating how to add 20+ co-authors (e.g., using Co-authored-by: footers).
2. **PR Integration (--pr flag)**:
   - Add --pr flag to commit-in, commit-left, commit-right (values: ll, 	ags,
elease).
   - Implement the PR workflow: create branch, push, create PR via GitHub API (gh pr create), merge PR, and release.
   - Integrate this flag into the command options struct.
3. **Terminal UI Enhancement**:
   - Update the terminal output for replaying commits to include left padding and color-coding.
   - Use github.com/charmbracelet/lipgloss for styling (e.g., dim tags, colorized action names).
4. **Commit Message Templating (Rewrite Generic Messages)**:
   - Provide a feature to intercept generic commit messages (e.g., "Changes", "Lovable update", "Work in progress").
   - By default, map these to something more descriptive (e.g., "chore: apply automated updates").
   - Read from gitmap settings to allow user-defined template overrides.
5. **Previous Commit URL Toggle**:
   - When copying commits, stop appending the original commit URL to the message body by default.
   - Add a configuration key in gitmap settings (e.g., CommitReplayKeepUrl = false) to toggle this behavior.
6. **Release**:
   - Final step: Bump minor version, update
eadme.md, update architecture map, and commit.

## Subtasks Mapping

See .lovable/plans/subtasks/01-commit-commands-overhaul/ for detailed subtasks.
