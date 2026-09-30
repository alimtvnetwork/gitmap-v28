# Plan 199: JSON Envelope Variables, WorkDirectory Object, OS Password CLI, and Repo-Secrets Sequence (Completed)

- **Status:** COMPLETED
- **Version:** `v6.416.0`
- **Spec Reference:** [02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md](../../../02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md)
- **Loops / Steps Taken:** 18 steps across Phase 1 Planning and Phase 2 Execution

---

## 1. Consolidated Subtasks & Verified Outcomes

### Subtask 01: JSON Envelope Variable Interpolation & Polymorphic `workDirectory` Object
- **Target Files:**
  - `cli/jsonenvelope/variables.go`
  - `cli/jsonenvelope/variables_test.go`
  - `cli/jsonenvelope/envelope.go`
- **Outcome:**
  - Implemented `ExpandVariables()`, `MergeVariables()`, and `ResolveStringVariable()` supporting `${varName}`, `${variables.varName}`, and `$variables.varName` interpolation with safe JSON string escaping for Windows paths (`\\`).
  - Added `WorkDirectoryConfig` struct (`path`, `defaultPath`, `isApplied`, `isEnforced`, `variables`) and custom `UnmarshalJSON` / `MarshalJSON` on `EnvelopeAttributes` so `workDirectory` can be either a scalar string or a structured object.

### Subtask 02: Cross-Platform OS Password Management (`gitmap os change-password`)
- **Target Files:**
  - `cli/cmdos/change_password_cmd.go`
  - `cli/cmdos/change_password_cmd_test.go`
  - `cli/cmdos/os.go`
  - `cli/cmdos/exports.go`
  - `cli/cmdos/os_help_modern.go`
  - `cli/cmd/roottooling.go`
- **Outcome:**
  - Implemented `RunChangePasswordCLI` supporting `gitmap os change-password [user] [password]`, `gitmap os passwd`, and root `gitmap change-password`.
  - Defaults to current OS user when omitted, prompts for confirmation unless `-y`/`--yes` is passed, prompts securely for password when omitted, and dispatches across Windows (`net user`), Linux/Ubuntu (`chpasswd`), and macOS (`dscl`).
  - Documented command and cross-platform examples in `gitmap os --help`.

### Subtask 03: Interactive SSH Join Password Prompting & Salted RSA Vault Storage
- **Target Files:**
  - `cli/cmdssh/sshjoin_enroll.go`
  - `cli/cmdssh/sshjoin_cmd.go`
- **Outcome:**
  - Added `promptPasswordIfInteractive` in `ExecuteSSHJoinEnrollment` so `gitmap ssh join {user}@ip <alias>` prompts for the SSH password if omitted on an interactive terminal (`Enter SSH password for <target> (leave blank to skip password vault): `).
  - If entered, encrypts with salted RSA (`EncryptSSHPassword`) and stores in `installation.db`; if left blank, skips password storage cleanly.

### Subtask 04: `repo-secrets/01-gitmap` Sequence Overhaul (`01-` to `09-`), Main Machine Info & Link Sync
- **Target Files:**
  - `D:\work\repo-secrets\01-gitmap\01-commit-pull-config.json`
  - `D:\work\repo-secrets\01-gitmap\02-git-profiles.json`
  - `D:\work\repo-secrets\01-gitmap\03-git-setup.json`
  - `D:\work\repo-secrets\01-gitmap\04-ssh-nodes.json`
  - `D:\work\repo-secrets\01-gitmap\05-import-ssh-nodes.ps1`
  - `D:\work\repo-secrets\01-gitmap\06-vmpass.json`
  - `D:\work\repo-secrets\01-gitmap\07-seo-templates.json`
  - `D:\work\repo-secrets\01-gitmap\08-seo-templates-cli.json`
  - `D:\work\repo-secrets\01-gitmap\09-test-gitmap-recreate.ps1`
  - `D:\work\repo-secrets\01-gitmap\readme.md`
  - `D:\work\repo-secrets\readme.md`
- **Outcome:**
  - Removed duplicate `00-commit-pull-config.json` and renamed all files in `01-gitmap/` to a strict two-digit lowercase sequence (`01-` through `09-`) with zero underscores.
  - Updated `04-ssh-nodes.json` with `${keyPath}`, `${adminUser}`, `${linuxUser}`, and `${workDir}` variables, structured `workDirectory` object, cross-platform OS import note, and complete `mainMachine` (`192.168.1.20`) record.
  - Updated all internal links and verified both `05-import-ssh-nodes.ps1` and `09-test-gitmap-recreate.ps1` end-to-end.
  - Committed and pushed changes to `repo-secrets`.

### Subtask 05: Version Bump to `v6.416.0` & Release
- **Target Files:**
  - `version.json`
  - `package.json`
  - `cli/constants/constants.go`
  - `changelog.md`
- **Outcome:**
  - Bumped version to `v6.416.0` and synchronized all version manifests and changelogs.
