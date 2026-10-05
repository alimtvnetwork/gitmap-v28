# Plan: 219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging

## User Request (Verbatim)

```text
In the Ubuntu system, I found these issues, but also at the same time, there is no stack trace so that you can remedy and find out what went wrong, and there is no solution. One more point. When we are in the Ubuntu machine, make sure that omizssh, we do not pick out automatically. And make sure that it is by default, automatically, it should not be catched by scan or something. But if the user adds it by default, manually add it to the scan, there should be specific force command or probably add a force for that repository mentioning this, and there should be command example for this. Then and only then this should be added. So make sure that you have a remedy for this. And also, I want you to connect back to the Ubuntu machine and try to see the root cause of it, why it happened, how it happened. Create scripts inside the repo secrets folder. Make sure that you implement those as well. And then you make sure that the solutions are integrated into the `gitmap`, and we make a final release bump and make a release. Also, write the root cause analysis, why it happened, how it happened, and also the, let's say, in the tree view, when you say next step, the `gitmap`, that actually go into another tree view, actually. So here, the next step, the option one, option two, that should be a subtree. Option one, option two. And next step, to clone this repository. The third section, the solution that you have provided, it feels a little bit buggy. So try to find the root cause. At the end, let me know the root cause, why it happened, how it happened, what was the solution, and `gitmap` should be able to fix these things. And also, if not, then it should give an error log and also hint the system to see the error log using error log methods. And you should also always do the error logs from your system and code so that it can be used for tracing, auditing, and things like that.
```

---

## 1. Objectives & Scope

1. **Task-01: Exclude Oh-My-Zsh & Default System Directories from Scanner / Discovery:**
   - Centralize default excluded directory names (`.oh-my-zsh`, `oh-my-zsh`, `ohmyzsh`, `.ohmyzsh`, `omz`, `.omz`, `omizssh`, `node_modules`, `vendor`, `.cache`, `.venv`, etc.) in `cli/constants/constants_scan.go`.
   - Integrate into `cli/scanner/scanner.go:buildExcludeSet`, `DefaultConfig` in `cli/model/record.go`, `cli/fsutil/child_repos.go`, and `cli/fsutil/recursive_top_level.go`.
   - Add `--force-include` to `gitmap scan` and `--force` to `gitmap add` with CLI command examples.
2. **Task-02: Cross-OS Backslash Path Normalization & Pre-Pull Auto-Healing:**
   - When loading repository paths on Unix/Linux, normalize backslashes (`\`) to forward slashes (`/`).
   - If a repository path exists with normalized separators or case differences, auto-heal the path in `gitmap.db` before pulling.
3. **Task-03: Pull Remediation Engine with Clone Fallback & Database-Only Removal:**
   - For missing repository directories, attempt path healing, then clone from known remote URLs (`HttpsUrl` / `SshUrl`), or offer `--db-only` removal instead of blindly executing `git pull`.
4. **Task-04: Pull Failure Tree Formatting & Hierarchical Subtrees:**
   - Refactor `renderSingleFailedItem` and `renderStructuredOptions` in `cli/cmdpull/pull_efficient_render.go` using box-drawing connectors (`├──`, `└──`, `│   `).
   - Fix dual options in `cli/cmdpull/pull_remediation_hint.go` for missing repositories.
5. **Task-05: Structured Error Logging, Stack Traces & Diagnostic Hints:**
   - Record pull errors with full timestamps, `NodeID`, and `StackTrace` in `gitmap-pull.db`.
   - Fix RFC3339 datetime parsing in `cli/store/pull_split_db_errors.go`.
   - Append structured error entries to `.gitmap/logs/pull-errors.log`.
   - Output terminal hint: `└── To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`.
6. **Task-06: Remote Ubuntu Node Diagnostic/Healing Scripts in `repo-secrets` & Release Ceremony:**
   - Create and execute `repo-secrets/05-scripts/heal-u1-pull-errors.py` (and `.sh`) on `u1`.
   - Document RCA-101 in `.ai-memory/cicd-issues/101-ubuntu-pull-all-remediation-and-omz-ignore-rca.md`.
   - Minor version bump, atomic commit, and release.

---

## 2. Canonical Specifications Reference

- Architecture Spec: [02-spec/21-app/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/01-architecture-spec.md](../../02-spec/21-app/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/01-architecture-spec.md)
- Component & SSH Spec: [02-spec/21-app/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/02-component-and-ssh-spec.md](../../02-spec/21-app/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/02-component-and-ssh-spec.md)

---

## 3. Subtask Decomposition

| Subtask ID | File | Focus | Owner |
|---|---|---|---|
| Subtask 01 | [subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/01-omz-scanner-exclusions-and-force-flag.md](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/01-omz-scanner-exclusions-and-force-flag.md) | Scanner exclusions, `--force-include`, `--force`, command examples | Worker 01 |
| Subtask 02 | [subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/02-cross-os-path-healing-and-missing-repo-pull.md](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/02-cross-os-path-healing-and-missing-repo-pull.md) | Cross-OS path normalization & auto-healing before pull | Worker 01 |
| Subtask 03 | [subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/03-pull-remediation-engine-and-clone-fallback.md](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/03-pull-remediation-engine-and-clone-fallback.md) | Remediation clone fallback & safe DB-only removal | Worker 01 |
| Subtask 04 | [subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/04-failure-tree-subtrees-and-dual-options.md](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/04-failure-tree-subtrees-and-dual-options.md) | Tree formatting with box connectors & dual options | Worker 02 |
| Subtask 05 | [subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/05-pull-error-logging-stacktrace-and-pe-hints.md](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/05-pull-error-logging-stacktrace-and-pe-hints.md) | Structured error logging, RFC3339 fix & hints | Worker 02 |
| Subtask 06 | [subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/06-remote-ubuntu-ssh-healing-script-in-repo-secrets.md](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/06-remote-ubuntu-ssh-healing-script-in-repo-secrets.md) | Remote healing script in `repo-secrets/05-scripts/` & verification | Worker 02 |
