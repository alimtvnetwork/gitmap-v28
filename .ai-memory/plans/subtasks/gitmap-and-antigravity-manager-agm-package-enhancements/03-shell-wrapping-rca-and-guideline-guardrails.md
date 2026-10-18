# Subtask 03: Shell Wrapping RCA, Guidelines & Native Git Commands

**Parent Task:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Subtask ID:** `subtask-03`  
**Status:** `IN_PROGRESS`  
**Owned Files:**  
- `.agents/skills/gitmap/SKILL.md`
- `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
- `.ai-memory/strictly-avoid.md`
- `cli/cmdgit/cmd_show.go`
- `cli/cmdgit/cmd_log.go`

## Objective
Identify why AI agents execute `gitmap ps "git ..."` or PowerShell commands unnecessarily instead of executing native GitMap commands or direct tool primitives. Implement prompt and skill refinements, strictly ban shell wrapping, and close GitMap's tooling gap by introducing native `gitmap show` and `gitmap log` commands.

## Implementation Details
1. Document the 5-part forensic RCA in `.ai-memory/strictly-avoid.md` and `02-spec/21-app/gitmap-and-antigravity-manager-agm-package-enhancements/01-architecture-spec.md`.
2. Add explicit ban in `.ai-memory/strictly-avoid.md`:
   - "TOTAL BAN on wrapping direct tools or standard git commands inside `gitmap ps` or `gitmap pwsh`."
3. Refine `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`:
   - Remove `gitmap ps` from search directives.
   - Clarify Worker Git Ban: mutating operations (`commit`, `push`, `checkout`, `add`) are banned; read-only inspection (`git show`, `git log`) is allowed when inspecting past revisions.
4. Implement native `gitmap show <rev> [-- <path>]` and `gitmap log [-n <N>]` in `cli/cmdgit/` so agents have first-class GitMap inspection tools.
