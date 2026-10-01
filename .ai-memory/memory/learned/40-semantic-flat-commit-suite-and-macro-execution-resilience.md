# Learned Memory 40: Semantic Flat Commit Suite & Macro Execution Platform Resilience

## 1. Context & Architectural Problem

Developers needed a 1-step auto-staging flat git commit command (`gitmap commit`, `cm`, `commit-all`, `ca`) that bypasses mandatory semantic prefixes like `feat:` or `fix:`. During initial attempts, a routing collision in `rootcore.go` shadowed the command, and a macro failure on Windows PowerShell demonstrated that `rm test` steps abort with terminating `ItemNotFoundException` if runtime binaries are not actively synchronized with platform-adaptive shims.

---

## 2. Key Decisions & Technical Principles

1. **Auto-Staging Flat Commit (`gitmap commit` / `cm`):**
   - Automatically executes `git add -A` to stage modifications, deletions, and new untracked files before committing.
   - Joins positional unquoted tokens into a single commit message.
   - Supports `--push` / `-p` for immediate remote synchronization.
   - Supports `--dry-run` / `-n` for safe preview via `git status --short`.
2. **Dispatch Consolidation in RootCore:**
   - Centralize all commit aliases in `cli/cmd/rootcore.go` using constants (`CmdCommit`, `CmdCommitAlias`, `CmdCommitAlias2`, `CmdCommitAlias3`).
   - Remove legacy passthrough `commit_cli.go` to eliminate shadow routing.
3. **Macro Platform Resilience on Windows:**
   - On Windows, `rm` is an alias for PowerShell's `Remove-Item` cmdlet, which throws terminating errors if the target is absent.
   - `cli/macro/safe_rm.go` transforms removal commands into idempotent PowerShell scripts that check `Test-Path` before invoking `Remove-Item`.
   - The compiled binary at both `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` and `./gitmap.exe` must be updated and kept current.

---

## 3. Reference Files

- Spec: `02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md`
- Command Spec: `.ai-memory/spec/commands/13-semantic-flat-commit-suite.md`
- Plan: `.ai-memory/plans/completed/63-semantic-flat-commit-suite.md`
- Subtasks: `.ai-memory/plans/subtasks/63-semantic-flat-commit-suite/`
- Issue RCA: `02-spec/22-app-issues/60-rootcore-commit-dispatch-shadowing-and-stale-binary-rca.md` & `.ai-memory/issues/26-rootcore-commit-dispatch-shadowing-and-stale-binary-rca.md`
