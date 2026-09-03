# 122 — Antigravity Empty Conversations Auditor & Pruner

## Overview

**Module Number:** 122  
**Version:** 1.0.0  
**Updated:** 2026-09-03  
**Status:** Production-Ready  
**AI Confidence:** Production-Ready  
**Ambiguity Score:** None  

---

## Purpose

Over extended agentic workflows, dozens of ephemeral Antigravity project workspaces accumulate in `~/.gemini/config/projects/*.json` that never initiated an active session, or failed during initialization before any user prompts occurred. This specification formalizes the conversation inspection algorithm, the tabular audit report (`agy ls show-projects-with-empty-conversations`), and the automated pruning command (`agy remove-projects-with-empty-conversations`).

---

## Technical Data Model & Discovery

### 1. File & Database Architecture
- **Projects**: Stored as JSON files in `~/.gemini/config/projects/<project-id>.json`. Contains `id`, `name`, and workspace URI in `projectResources.resources[].gitFolder.folderUri`.
- **Conversations**: Stored as SQLite databases in `~/.gemini/antigravity/conversations/<conversation-id>.db`.
- **Workspace Linking**: In `<conversation-id>.db`, the table `trajectory_metadata_blob` stores serialized metadata containing the workspace `file:///` path URI.

### 2. Definition of an Empty Conversation Project
A project is classified as an **Empty Conversation Project** if:
1. It has zero associated conversation databases in `~/.gemini/antigravity/conversations/`, OR
2. All associated conversation databases have `steps <= 2` and `user_steps == 0` (failed or aborted initialization sessions).

Projects with active conversations (`steps > 2` or `user_steps > 0`) are classified as **Active** and preserved.

---

## Command Signatures & Behaviors

### 1. `gitmap agy ls show-projects-with-empty-conversations`
- **Aliases**: `show-proects-with-empty-conversations`, `empty-conversations`, `--empty-conversations`, `empty-convs`
- **Output**:
  ```text
  ── Antigravity Projects with Empty Conversations ──

  Found 50 project(s) with empty or zero conversations (out of 59 total):

  PROJECT ID                             NAME                     WORKSPACE PATH                             CONVS  STATUS
  ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  098304c1-b79b-4c7e-a223-b4b318b8a6cd   prompts-connect-v3       D:\wp-work\riseup-asia\02-prompts\prom...  0      No Convs
  ```
- Directly below findings, prints remediation recipes.

### 2. `gitmap agy remove-projects-with-empty-conversations [flags]`
- **Aliases**: `rm-empty-conversations`, `clean-empty-conversations`, `prune-empty-conversations`
- **Flags**:
  - `--except <spec>` / `-e`: Comma-separated list or file path (`.csv` / `.txt`) containing IDs, names, paths, or aliases to exclude from removal.
  - `--dry-run` / `-d`: Previews target projects without deleting files.
  - `--yes` / `-y` / `-f` / `--force`: Non-interactive bypass of confirmation prompt.
- **Confirmation Prompt**:
  ```text
  Are you sure you want to remove 50 Antigravity project(s) with empty conversations? [y/N]:
  ```
- **Execution**: Deletes the corresponding `<project-id>.json` files and reports the remaining active project count.

---

## Acceptance Criteria

### Scenario 1: Whitelisted Preservation
- **Given** 50 empty projects, including `prompts-connect-v3` and `extendcore`
- **When** `gitmap agy remove-projects-with-empty-conversations --except "prompts-connect-v3, extendcore" --dry-run` is run
- **Then** the CLI targets 48 projects for removal, explicitly sparing the two whitelisted projects.

### Scenario 2: Safe Deletion Confirmation
- **Given** detected empty projects
- **When** the removal command is invoked without `--yes`
- **Then** execution blocks on interactive confirmation `[y/N]` and only deletes files upon affirmative user input.

---

## Cross-References

- Cross-Platform Duplicate Audit: [`121-cross-platform-duplicate-audit-and-remediation.md`](./121-cross-platform-duplicate-audit-and-remediation.md)
- CLI Help Document: [`../../gitmap/helptext/agy-empty-conversations.md`](../../gitmap/helptext/agy-empty-conversations.md)
