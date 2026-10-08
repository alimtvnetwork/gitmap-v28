# Subtask 03: Chore Commits (`cpc`) & Help Menu Full Forms

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md`  
> **Associated Specs:**  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/01-architecture-spec.md`  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/02-component-spec.md`  
> **Owned Files:**  
> - `cli/constants/constants_cli.go`  
> - `cli/cmd/commit_push.go`  
> - `cli/cmd/commit_help_menu.go`  
> - `cli/cmd/roottooling.go`  
> - `cli/cmd/nodes_cmd.go`  
> - `cli/cmd/sends_cmd.go`  
> - `cli/helptext/commit-push-chore.md`  
> - `cli/helptext/catalog.go`  
> - `cli/helptext/commit-push.md`  
> - `cli/helptext/help.md`  

---

## 1. Objectives

- [ ] 1. Define CLI constants in `cli/constants/constants_cli.go`:
  - `CmdCommitPushChore = "commit-push-chore"`
  - `CmdCommitPushChoreAlias = "cpc"`
- [ ] 2. Implement chore commit command handler in `cli/cmd/commit_push.go`:
  - Implement `runCommitPushChore(args []string) error` with help argument detection (`isCommitPushHelpArg`), validation for empty messages (`Usage: gitmap commit-push-chore "<what chore was done>"`), and `"Chore: "` message prefixing.
  - Wire into the top-level command switch in `cli/cmd/roottooling.go` under `case constants.CmdCommitPushChore, constants.CmdCommitPushChoreAlias:`.
  - Wire into multi-node commit dispatchers in `cli/cmd/nodes_cmd.go` and `cli/cmd/sends_cmd.go`.
- [ ] 3. Update CLI help menu in `cli/cmd/commit_help_menu.go`:
  - Add `gitmap cpc "<message>" (commit-push-chore)` to `UsageLines`.
  - Update `buildCommitAIWorkflowSection()` with full forms:
    - `commit, cm <msg>`: commit (flat git commit with auto-stage)
    - `cpf <msg>`: commit-push-feature (stage all, commit feat, and push)
    - `cpb <msg>`: commit-push-bug (stage all, commit bugfix, and push)
    - `cpc <msg>`: commit-push-chore (stage all, commit chore/maintenance, and push)
    - `cpr <msg>`: commit-push-release (stage all, commit release chore, and push)
    - `pcp <msg>`: pull-commit-push (pull latest rebase, stage, commit, and push)
    - `pas`: pull-all-ssh (pull all repositories using SSH transport)
- [ ] 4. Author standalone help markdown document in `cli/helptext/commit-push-chore.md`:
  - Provide complete usage documentation, comparison with manual git commands (`git add -A && git commit -m "Chore: ..." && git push`), aliases (`gitmap cpc`), and real-world examples.
- [ ] 5. Register and cross-reference documentation in `cli/helptext/`:
  - Register `commit-push-chore` in `cli/helptext/catalog.go`.
  - Update `cli/helptext/commit-push.md`, `cli/helptext/commit.md`, and `cli/helptext/help.md` to reference `cpc` and document full forms.
- [ ] 6. Verify compilation and linter compliance:
  - `go test -v ./cli/cmd/...`
  - `python linter-scripts/check-nested-ifs.py`
  - `python linter-scripts/check-enum-and-boolean.py`
  - `python linter-scripts/check-relative-paths.py`

---

## 2. Types & Constants

```go
package constants

const (
	// CmdCommitPushChore stages all changes, commits with a "Chore: " prefix, and pushes.
	CmdCommitPushChore      = "commit-push-chore"
	CmdCommitPushChoreAlias = "cpc"
)
```

```go
package cmd

// runCommitPushChore commits with a "Chore: " prefix.
func runCommitPushChore(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushChore, []string{"--help"})
		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-chore \"<what chore was done>\"", "E9000")
	}

	commitMessage := "Chore: " + strings.Join(args, " ")
	return executeCommitPush(commitMessage)
}
```

---

## 3. Acceptance Criteria

- **AC-CPC-001 (Command Registration):** Running `gitmap commit-push-chore --help` or `gitmap cpc --help` renders help text from `cli/helptext/commit-push-chore.md`.
- **AC-CPC-002 (Commit Prefixing):** Executing `gitmap cpc "cleanup unused styles"` creates a Git commit with subject `Chore: cleanup unused styles` and pushes to remote.
- **AC-CPC-003 (Help Menu Full Forms):** Running `gitmap commit --help` or `gitmap help commit` displays the updated help menu explicitly listing the full forms:
  - `cpf` -> commit-push-feature
  - `cpb` -> commit-push-bug
  - `cpc` -> commit-push-chore
  - `cpr` -> commit-push-release
  - `pcp` -> pull-commit-push
  - `pas` -> pull-all-ssh
- **AC-CPC-004 (Empty Message Guard):** Running `gitmap cpc` without arguments returns exit code with error message `Usage: gitmap commit-push-chore "<what chore was done>"`.
- **AC-CPC-005 (Hyphen Syntax Compliance):** Coding guidelines and skills enforce hyphen syntax for chore messages (e.g. `gitmap cpc "cli - update help text"`; ban on duplicate colons like `gitmap cpc "Chore: update"`).

---

## 4. Verification Instructions

1. Run unit tests:
   ```bash
   go test -v ./cli/cmd/...
   ```
2. Verify help text output:
   ```bash
   gitmap cpc --help
   gitmap commit --help
   ```
3. Run repository linters:
   ```bash
   python linter-scripts/check-nested-ifs.py
   python linter-scripts/check-enum-and-boolean.py
   python linter-scripts/check-relative-paths.py
   ```
