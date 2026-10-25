# Subtask 05: Verification and Release Ceremony

- **Parent Task:** 327 (`heatmap-status-improvement-and-ai-confusion-handling`)
- **Subtask Sequence:** 05
- **Subtask Path:** `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/05-subtask-verification-and-release-ceremony.md`
- **Assigned Subagent:** Subagent 3
- **Created At:** 2026-10-25
- **Status:** `PENDING`

---

## 1. Context & Objectives

This subtask governs end-to-end verification of all deliverables in Task 327 across Subagents 1, 2, and 3, followed by the execution of the minor version release ceremony (`v6.611.0`).

### 1.1 In-Scope Deliverables
1. **Repo Cache Sequence Index:** `RepoCacheSequenceIndex` in root SQLite database caching `RepoSlug` $\rightarrow$ `XX-<repo>` folder names.
2. **Normalized AI Confusion Database:** 1-to-N normalized `AiConfusion` and `AiConfusionFile` schema.
3. **Autonomous Confusion Sync & Status:** `gitmap rc ai-confusion status` and relative file path querying (`gitmap rc ai-confusion view <file>`).
4. **Pull Remediation & Non-Repo Sanitization:** `isValidRepoNameCandidate` rejecting commit messages and ghost repositories from failure lists.
5. **Pull Output Collapsing:** Clean collapsing of up-to-date repositories with one-liner command hint.
6. **Minor Version Release Ceremony:** Bump to `v6.611.0`.

---

## 2. Verification Protocol

### 2.1 Code & Schema Verification
- Review Go code diffs for compliance with GitMap coding rules:
  - <= 15 lines per function.
  - Zero nested `if` statements.
  - Positive booleans (`is*` / `has*`).
  - Blank line spacing rules.
  - Leaf error sovereignty.

### 2.2 Functional Acceptance Checklist
- [ ] Root SQLite sequence cache table exists and supports upsert and query.
- [ ] AI confusion database establishes 1-to-N file records with start/end lines and code snippets.
- [ ] Non-repo folder classification rejects commit message strings and phantom directories.
- [ ] Pull failure conversion filters out records where directory does not exist and remote is empty.
- [ ] Semantic release ceremony files updated cleanly.

---

## 3. Release Ceremony Protocol (`v6.611.0`)

1. **Version Files:**
   - Update `package.json` to `"6.611.0"`.
   - Update `version.json` to `"6.611.0"`.
   - Ensure `cli/constants/version.go` references `"6.611.0"`.
2. **Changelog:**
   - Append release highlights to `changelog.md` under `## [6.611.0] - 2026-10-25`.
   - Create release record in `.ai-memory/release/v6.611.0.md`.
3. **Semantic Release Commit:**
   - Commit changes using `gitmap cpr "RELEASE - bump version to v6.611.0"`.
