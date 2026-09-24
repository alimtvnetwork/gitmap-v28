# Plan 96: Pull Batch Abort, CFR GitHub Resolver, Asset Downloader, Color Contrast, and Token Fleet Management

> **Plan Status:** Completed  
> **Traceability IDs:** Task-01 .. Task-10  
> **Spec Reference:** [02-spec/21-app/147-pull-abort-cfr-gh-resolver-asset-downloader-and-contrast.md](../../../02-spec/21-app/147-pull-abort-cfr-gh-resolver-asset-downloader-and-contrast.md)  
> **Issues Covered:** [02-spec/22-app-issues/39-pull-batch-abort-and-search-untracked-directory-crash.md](../../../02-spec/22-app-issues/39-pull-batch-abort-and-search-untracked-directory-crash.md), [02-spec/22-app-issues/40-cfr-short-name-clone-failure-and-missing-gh-resolution.md](../../../02-spec/22-app-issues/40-cfr-short-name-clone-failure-and-missing-gh-resolution.md)

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 95-01** | `cmdpull`, `cmd/search.go` | Graceful pull abort on quit/skip, no `E9000` stack trace; untracked search directory fallback | Completed |
| **Subtask 95-02** | `cmdclone` | CFR short-name slug resolution via GitHub CLI & SQLite `repodb` cache; Desktop & VS Code sync | Completed |
| **Subtask 95-03** | `cmdasset` / CLI | `gitmap asset download-prnt <url>` command extracting LightShot images to filesystem & brain | Completed |
| **Subtask 95-04** | `cmdtoken` / `cmdssh` | `gitmap token add|remove|list|deploy` for multi-node SSH token deployment | Completed |
| **Subtask 95-05** | Terminal & Install UI | High-contrast color replacement for low-contrast dark blue accents across CLI & install scripts | Completed |
| **Subtask 95-06** | Prompts Sync | Folder 21 E2E testing prompt 2-half N-step lifecycle, and 2-agent mandate updates in prompts | Completed |

---

## 2. Ingested Screenshots Reference

- `![Screenshot 1](assets/screenshots/8BnieUEKidCr.png)` — Search failure in untracked directory
- `![Screenshot 2](assets/screenshots/_N8xDMM-6ylG.png)` — Dark blue contrast and pull quit failure
- `![Screenshot 3](assets/screenshots/0J1Th8lwMNIL.png)` — Scripts fixer and installation UI
- `![Screenshot 4](assets/screenshots/w66J-xV1NN3E.png)` — Pull batch failure on 62 repos

---

## 3. Verified Outcomes

1. **Pull Batch Abort**: Removed fatal `[E9000:EXECUTION]` stack dump when quitting or skipping dirty repository remediation in `cli/cmdpull/pull.go`.
2. **Untracked Directory Search**: Added fallback in `cli/cmd/search.go` when run outside a tracked repository, querying global repositories or searching current filesystem.
3. **CFR Short-Name GitHub CLI Resolution**: Created `cli/cmdclone/gh_resolver.go` to resolve bare repo slugs (e.g., `pwp-mobile`) via `gh repo list` / local cache, and sync clones with GitHub Desktop and VS Code `projects.json`.
4. **LightShot & Asset Downloader**: Added `cli/cmdasset/asset_prnt.go` and `gitmap asset download-prnt` / `gitmap prnt` command.
5. **Token Fleet Distribution**: Implemented `gitmap token` and `cli/cmdssh/ssh_token_deploy.go` for multi-node SSH git credentials sync.
6. **High-Contrast Terminal Styling**: Replaced dark blue accents with high-contrast cyan `#8be9fd` and bright white `#f8f8f2` in `cli/cmd/reconcile_prompt.go` and `install.ps1`.
7. **Prompt & E2E Guidelines**: Synchronized guidelines in `01-prompts/` across both repositories without backporting code to coding-guidelines.
