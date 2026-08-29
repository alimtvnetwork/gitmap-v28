# Plan 22: Memory Consolidation — Completed Plans & Re-Sequence Milestones

## 1. Safety Backup & Rollback Recipe

- **Backup Branch Name:** `backup/plans-consolidation-20260829-203850`
- **Backup Commit SHA:** `cecf50a5355bedc8248f2cb336eed03e73343065`
- **Pushed to Origin:** `git push origin backup/plans-consolidation-20260829-203850` (Verified on origin)
- **Active Working Branch:** `main`

```bash

# Rollback Command (in case of accidental data loss or rollback need):

git reset --hard backup/plans-consolidation-20260829-203850
```

---

## 2. Inventory & Domain Clustering Mapping Table

| Source Files to Merge | Proposed Consolidated File | Domain / Epic Theme | Items & Specs Preserved | Status |
|---|---|---|---|:---:|
| `01-coding-guideline-fixes.md`, `04-cfr-cg-os-aware-coding-guidelines.md`, `04-cg-multirepo-and-status-dirty.md`, `16-nested-if-audit.md`, `17-boolean-and-naming-audit.md`, `17-booleans-and-complex-conditions-audit.md`, `01-coding-guideline-fixes/`, `04-cfr-cg-os-aware-coding-guidelines/`, `16-nested-if/`, `17-booleans/` | `01-coding-guidelines-and-boolean-refactoring.md` | Coding Guidelines & Boolean Quality | Affirmative prefixes (`is`, `has`), zero `== true`, 0 nested `if`s, MD022/MD032 spacing, LF/UTF-8 normalizations. | PENDING |
| `02-error-management-fixes.md`, `15-centralized-error-handling-and-exit-architecture.md`, `16-error-management-audit.md`, `02-error-management-fixes/` | `02-error-management-and-exit-architecture.md` | Centralized Error Architecture | AppError wrappers, cliexit handlers, universal response envelopes, structured metadata, and CI error linter. | PENDING |
| `01-ssh-login-and-join.md`, `02-ssh-aware-clone.md`, `03-installer-multios-cluster.md`, `06-cluster-command-delegation.md`, `02-ssh-aware-clone/`, `installer-multios-cluster/`, `ssh-login-and-join/` | `03-ssh-nodes-and-cluster-delegation.md` | SSH, Nodes & Cluster Management | SSH key generation, multi-host alias config, cluster node joining, broadcast/distribute commands, and remote login. | PENDING |
| `01-ui-and-macro-features.md`, `07-update-terminal-visualization.md`, `08-dashboard-recent-and-terminal-ui.md`, `ui-and-macro-features/` | `04-ui-terminal-and-dashboard-visualization.md` | UI, Terminal & Dashboard | Terminal column width, help text alignment, TUI tree views, dashboard live status, macro recorder and player. | PENDING |
| `01-bulk-visibility-mapub-mapri.md`, `02-chrome-profile-migration.md`, `03-reclone-transport-and-vscode-open.md`, `05-gitmap-improvements.md`, `05-lfs-smudge-fallback.md`, `05-mv-rm-resolver-replace-100-steps.md`, `05-workdir-pull-table-dirty-remedy.md`, `01-bulk-visibility-mapub-mapri/`, `03-reclone-transport-and-vscode-open/`, `chrome-profile-migration/`, `workdir-pull-table-dirty-remedy/` | `05-workspace-profile-and-repository-operations.md` | Workspace & Repo Operations | Repository moving/untracking (`mv`, `rm`), Chrome profile sync, Split SQLite DB operations, and LFS smudge fallback. | PENDING |
| `00-execution-plan.md`, `01-zsh-kube-consolidation.md`, `02-gitmap-installer.md`, `03-cg-and-macro-installer.md`, `06-prompt-architect-installer.md`, `15-relative-paths-audit.md`, `001-task.md` .. `255-task.md`, `15-relative-paths/`, `cg-and-macro-installer/`, `gitmap-installer/`, `prompt-architect-installer/`, `zsh-kube-consolidation/`, `subtasks/` | `06-installers-scaffolding-and-tooling-integrations.md` | Installers, Scaffolding & CI Tooling | Cross-platform setup (ZSH, Ubuntu, Windows), prompt architect engine, relative path enforcement, and AI fix scripts. | PENDING |

---

## 3. Subtasks Decomposition

- **Subtask 22.01:** Create consolidated milestone summaries (`01-` to `06-`) adhering to the standard milestone template with zero loss of specs.
- **Subtask 22.02:** Remove superseded micro-task files and empty subtask directories in `.lovable/plans/completed/`.
- **Subtask 22.03:** Re-sequence `.lovable/plans/completed/` to contiguous `01-` through `06-` using `01-file-manipulator.py`.
- **Subtask 22.04:** Synchronize `.lovable/plans/index.md` and `.lovable/memory/00-index.md` and verify with all relative path and spacing linters.
