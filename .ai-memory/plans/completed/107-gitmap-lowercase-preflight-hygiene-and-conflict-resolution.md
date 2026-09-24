# Completed Plan 107: GitMap Lowercase Pre-Flight Hygiene, Working Tree Conflict Resolution & Automated Push

Spec Reference: [02-spec/21-app/156-gitmap-lowercase-preflight-hygiene-and-conflict-resolution.md](../../02-spec/21-app/156-gitmap-lowercase-preflight-hygiene-and-conflict-resolution.md)  
Issue Reference: [02-spec/22-app-issues/43-gitmap-lowercase-unresolved-conflicts-and-missing-push-rca.md](../../02-spec/22-app-issues/43-gitmap-lowercase-unresolved-conflicts-and-missing-push-rca.md)  
Execution Summary: Completed in 2 orchestration loops across 4 subtask domains with comprehensive pre-flight verification, conflict detection, discard prompt, full disclosure, auto-push, and minor version bump to v6.338.0.

## User Request (Verbatim)

```text
https://prnt.sc/n_O_kOGNEJqZ


PS D:\work\presentations-repos\hiltrax> gitmap lowercase "*.md"

⚡ GitMap Lowercase File Renamer
  ● Mode:        Git Repository (2-step git mv)
  ● Working Dir: D:\work\presentations-repos\hiltrax
  ● Filter:      *.md
  ● Scanned:     1043 files in directory
  ● Matched:     18 uppercase file(s)


════════════════════════════════════════════════════════════════
⚡ Pre-Flight Verification: Git Lowercase Renamer
════════════════════════════════════════════════════════════════
  ● Found: 18 uppercase file(s) to rename:
    • SKILL.md -> skill.md
    • SKILL.md -> skill.md
    • SKILL.md -> skill.md
    • SKILL.md -> skill.md
    • AGENTS.md -> agents.md
    • LLM.md -> llm.md

...
Okay. So here the big problem is that even if you fix those uppercase to lowercase, the first thing you created is the Git conflict, and you didn't resolve the Git conflict. Okay, and it looks like that there are too many Git merge happen, conflict resolved. Why? Because you are doing this change, you should make sure that it is done properly. Okay, if there is something pending, you make sure that you discard it, you don't care it, and you confirm with the user. You list out all the files that you are dealing with, make sure there is no hidden stuff. And then when you start, you make sure you commit and resolve and push to the Git, which is missing from your task. Is it understood? Can you please fix it and bump the minor version and release it
```

---

## Consolidated Subtasks & Outcomes

### Subtask 01: Pre-Flight Working Tree Hygiene & Git Conflict Guard
- **Target Files:** `cli/cmd/lowercasefix_git.go`, `cli/cmd/lowercasefix_ops.go`
- **Delivered:**
  - Added `checkGitWorkingTreeStatus()` using `git status --porcelain` before executing renames.
  - Detected unresolved merge conflicts (`UU`, `AA`, `UD`, `DU`, `DD`, `AU`, `UA`) and halted immediately with clear conflict resolution guidance.
  - Audited pending uncommitted working tree changes (`M `, ` M`, `A `, `??`).

### Subtask 02: Full Relative Path Disclosure & Pending Dirt Discard Confirmation
- **Target Files:** `cli/cmd/lowercasefix_confirm.go`, `cli/cmd/lowercasefix_ops.go`, `cli/cmd/lowercasefix.go`, `cli/cmd/lowercasefix_types.go`
- **Delivered:**
  - Removed hardcoded 8-file truncation limit in `renderPreflightFound` so all candidate files are disclosed without hidden items.
  - Displayed relative path mapping for each file (e.g. `path/to/SKILL.md → skill.md`) to distinguish identical filenames across subdirectories.
  - Added `--discard-pending` and `-f, --force` flags.
  - If uncommitted changes exist without `--discard-pending`, prompts user interactively with file list and offers to discard them via `git reset --hard HEAD` and `git clean -fd`.

### Subtask 03: Atomic Commit & Remote Push
- **Target Files:** `cli/cmd/lowercasefix_commit.go`, `cli/cmd/lowercasefix_report.go`, `cli/cmd/lowercasefix_ops.go`, `cli/cmd/lowercasefix.go`
- **Delivered:**
  - Updated `commitRenames()` to invoke `executeGitPush()` automatically pushing to `origin <branch>` (with fallback to `git push`).
  - Added `--no-push` flag to optionally bypass automatic remote push.
  - Reflected push status in pre-flight confirmation box and execution summary report table.

### Subtask 04: Minor Version Bump & Release Orchestration
- **Target Files:** `version.json`, `cli/constants/constants.go`, `package.json`, `readme.md`, `changelog.md`, `98-changelog.md`
- **Delivered:**
  - Bumped version from `6.337.0` to `6.338.0` via repository release tooling.
  - Synchronized version constants and changelog entries.
