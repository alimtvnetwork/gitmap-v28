# Master Plan: Multi-User Project Switching, GitHub CLI Auth Status & Git Config Manager

- **Task Slug:** `242-multi-user-project-switching-gh-auth-status-and-git-config-manager`
- **Execution Mode:** Multi-Agent Autonomous (v6 Protocol: A=2, H=2, zero solo execution)
- **Primary Goal:** Implement GitHub CLI authorization readout on launch, `gitmap user info` public inspector, multi-user Git profile switching per project, and global/local Git configuration commands.

---

## Subtask Breakdown

| ID | Code | Subtask Title | Assigned Role | Status |
| :--- | :--- | :--- | :--- | :--- |
| 1 | `Task-01` | GitHub CLI Authorization Status Detection & Launch Readout | Worker 01 | PENDING |
| 2 | `Task-02` | User Info Command & Git Configuration Inspector (Global & Local) | Worker 02 | PENDING |
| 3 | `Task-03` | Multi-User Profile Store & Per-Project Auto-Switching Engine | Worker 01 | PENDING |
| 4 | `Task-04` | GitMap User CLI Subcommands, Project Configuration & Release Ceremony | Worker 02 | PENDING |

---

## Dependencies & Sequencing
- **Phase 1 (Core Engines):**
  - Worker 01 implements GitHub CLI auth detection (`cli/usercontext/gh_auth.go`) and launch readout hook (`cli/cmd/rootusagecompact.go`).
  - Worker 02 implements Git identity inspector (`cli/usercontext/git_identity.go`) and `gitmap user info` command (`cli/cmd/user_info.go`).
- **Phase 2 (Profiles & Switching):**
  - Worker 01 extends model & store (`cli/model/git_profile.go`, `cli/store/git_profile_store.go`, `cli/usercontext/project_binder.go`).
  - Worker 02 wires CLI subcommands in `cli/cmd/user_cmd.go` (`list`, `switch`, `use`, `project bind`, `config global/local`).
- **Phase 3 (Verification & Release):**
  - Run FastGate pre-commit linters: `gitmap py 03-ai-scripts/50-fastgate.py`.
  - Run package unit tests: `go test -v ./usercontext/... ./cmd/...`.
  - Bump minor version: `03-ai-scripts/37-bump-version.py -t minor`.
  - Commit & push: `gitmap cpf "user - multi-user project switching, github cli auth status, and git config manager"`.
  - Monitor CI/CD telemetry via `gitmap pe -t --ai` until green.
