# Subtask 04: GitMap User CLI Subcommands, Project Configuration & Release Ceremony

- **Parent Task:** `242-multi-user-project-switching-gh-auth-status-and-git-config-manager`
- **Subtask Code:** `Task-04`
- **Assigned Worker:** Worker 02
- **Status:** COMPLETED

---

## Scope & Implementation Details
1. Extend `cli/cmd/user_cmd.go` and `cli/cmd/user_help_menu.go`:
   - Route Git user subcommands alongside OS user commands:
     - `gitmap user info` / `user-info`: Calls `runUserInfo(args)`.
     - `gitmap user list` / `ls`: Calls `runUserList(args)` (listing Git profiles, active user, project user).
     - `gitmap user switch <alias>` / `use <alias>`: Calls `runUserSwitch(args)` (setting local/global and project binding).
     - `gitmap user add <alias>`: If `--name` or `--email` provided (without `--password`), registers Git user profile.
     - `gitmap user project [bind|unbind|status]`: Calls `runUserProject(args)`.
     - `gitmap user config [global|local]`: Calls `runUserConfig(args)` to inspect or set name/email.
     - `gitmap user sync`: Calls `runUserSync(args)` to apply bound user identity.
2. Update `.agents/skills/gitmap/SKILL.md` and `.agents/skills/coding-guidelines/skill.md`:
   - Document `gitmap user info`, `gitmap user switch`, `gitmap user project`, and launch authorization readout.
3. Verification & Quality Gates:
   - Run FastGate: `gitmap py 03-ai-scripts/50-fastgate.py`.
   - Run unit tests: `go test -v ./usercontext/... ./cmd/...`.
   - Run AST linters: `check-relative-paths.py`, `check-nested-ifs.py`, `check-boolean-guidelines.py`.
4. Release Ceremony:
   - Minor bump via `03-ai-scripts/37-bump-version.py -t minor`.
   - Commit & push via `gitmap cpf "user - multi-user project switching, github cli auth status, and git config manager"`.
   - Tag release `v6.513.0` and monitor CI/CD via `gitmap pe -t --ai` until green.

---

## Acceptance Criteria
- [ ] All `gitmap user` subcommands execute cleanly.
- [ ] Backward compatibility with OS user commands is preserved.
- [ ] Zero lint violations across all modified and created files.
- [ ] Release ceremony completes and telemetry monitoring passes.
