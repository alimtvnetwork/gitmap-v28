# Subtask 01: SQLite `SpecialRepository` Split-DB, One-Time Scan Prompt, `gitmap cd rs`/`rc` & `gitmap rs`/`rc` Put Commands

## Allowed Target Files (Strict Disjoint Bounding Box)
- `cli/store/special_repos_split_db.go`
- `cli/cmd/special_repos_cli.go`
- `cli/cmd/special_repos_ops.go`
- `cli/cmd/cd.go`
- `cli/cmd/roottooling.go`
- `cli/cmd/special_repos_test.go`

## Deliverables
1. SQLite Split-DB (`gitmap-special-repos.db` in `store.BinaryDataDir()`) with `SpecialRepository` (`ShortKey` = `"rs"` / `"rc"`, `DefaultName`, `ConfiguredName`, `Category`, `LocalPath`, `RemoteURL`, `IsPromptAnswered`, `UserDecision`, `PromptedAt`) and `SpecialRepoFolderSeq` (`SpecialKey`, `RepoName`, `FolderPrefix`, `SeqNumber`).
2. `gitmap cd rs` (`repo-secrets`) and `gitmap cd rc` (`repo-cache` / `repo-storage`) shortcut resolution in `cli/cmd/cd.go`.
3. `gitmap rs` (`repo-secrets`) and `gitmap rc` (`repo-cache` / `repo-storage`) with `file`, `folder`, `text`, `ls`, `init`, `scan-check`, and `help` subcommands (`XX-<repo-name>/01-<slug>.ext` auto-sequencing + automatic `git add`, `git commit`, and `git push`).
4. First-scan one-time check & prompt (`CheckSpecialReposOnScan`) integrated so each special repository is prompted at most once and persisted in `SpecialRepository`.
