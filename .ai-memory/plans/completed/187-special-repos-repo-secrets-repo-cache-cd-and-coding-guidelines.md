# Plan 187: Special Default Repositories (`repo-secrets` = `rs` & `repo-cache` = `rc`), One-Time Scan Prompt, `gitmap cd rs`/`rc`, Auto-Commit/Push Put Commands & Coding Guidelines Update

## 1. Objective & Scope
Implement `02-spec/21-app/91-special-repos-repo-secrets-repo-cache-cd-and-coding-guidelines.md` across GitMap and `coding-guidelines`, release a minor version (`v6.380.0`), and verify with `gitmap pe`.

## 2. Implementation Summary
- **SQLite Split-DB (`store.SpecialReposSplitDB`)**:
  - Implemented `SpecialReposDBFileName = "gitmap-special-repos.db"` stored in `BinaryDataDir()`.
  - Tables: `SpecialRepository` and `SpecialRepoFolderSeq`.
  - Monotonic 2-digit folder and file allocation (`01-gitmap`, `02-coding-guidelines`) with disk and DB alignment.
  - One-time scan prompt state persistence (`IsPromptAnswered = 1`, `UserDecision`) so subsequent scans never repeat the question once answered.
- **Fast Navigation (`gitmap cd rs` & `gitmap cd rc`)**:
  - Integrated into `cli/cmd/cd.go` (`isSpecialRepoCDAlias`), navigating directly to `repo-secrets` or `repo-cache` via `WriteShellHandoff(targetPath)`.
- **Special Repo Put Commands (`gitmap rs` & `gitmap rc`)**:
  - `file <filepath>`: Copies `<filepath>` to `<special-repo>/XX-<repo-name>/NN-<filename>` with automatic `git add`, `git commit`, and `git push`.
  - `folder <folderpath>`: Recursively copies folder into `<special-repo>/XX-<repo-name>/NN-<foldername>` with automatic commit and push.
  - `text "<content>"`: Derives kebab-case slug and writes `<special-repo>/XX-<repo-name>/NN-<slug>.ext` with automatic commit and push.
  - `ls`: Lists all `XX-<repo-name>/NN-<slug>` items on disk and in JSON mode.
  - `init`: Creates and verifies git repositories for `repo-secrets` and `repo-cache`.
  - `scan-check`: Checks and persists one-time scan decisions in SQLite.
- **Unified Settings & Web UI**:
  - Added `special_repos.secrets_name` and `special_repos.cache_name` to `gitmap settings` and Web UI (`cli/cmdui/ui_assets.go`).
- **Help Text & V2 Prompts**:
  - Authored `cli/helptext/repo-secrets.md`, `cli/helptext/repo-cache.md`, registered in `cli/helptext/catalog.go`, and authored `01-prompts/special-repos-secrets-and-cache.md`.
- **Coding Guidelines Sync**:
  - Pulled `d:\work\coding-guidelines`, updated `agents.md`, `01-prompts/v2/00-folder-structure/01-canonical-folder-structure.md`, `01-prompts/v2/04-coding-standards/01-coding-guidelines.md`, authored `02-spec/02-coding-guidelines/01-cross-language/30-special-repos-secrets-and-cache.md`, committed, and pushed (`fdef19d1`).

## 3. Verification & CI Gates
- Unit tests: `TestSpecialReposNormalizeAndSchema`, `TestDeriveTextSlugAndFilename`, `TestIsSpecialRepoCDAlias`, `TestAllocateSequencedItemPath` all passed.
- All 5 CI linters passed with 0 violations (`gofmt`, boolean, relative paths, error management, nested ifs).
- Minor version bump: `v6.380.0`.
