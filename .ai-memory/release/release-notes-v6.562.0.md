# GitMap v6.562.0

## What's Changed in v6.562.0

- **Large Files Splitting Engine (`gitmap split large files`):** Discovers files $> 10\text{ MB}$, creates multipart zip chunks (`<file>.zip`, `<file>.z01`), applies STORE for media and DEFLATE for data/code, removes raw files from Git history, and warns before rewriting.
- **Git LFS Conversion Suite (`gitmap lfs convert`):** Automated scanner finding files $> 5\text{ MB}$, updates `.gitattributes`, migrates history to Git LFS pointers, and previews via `--dry-run`.
- **Repository Flattening Engine (`gitmap repo flatten`):** Compresses bloated repositories into a clean successor (`<repo>-flat-v1`) with a single flat squashed snapshot commit, archives the previous repository to `<repo>-marked-for-deletion`, and swaps to the canonical name.
- **Flatten Undo Manager (`gitmap repo undo-flatten`):** Safely restores the previous repository archive from `<repo>-marked-for-deletion`.
- **Muse Training Priority Sync:** Prioritizes `repo-secrets` checkouts for active LLM training, restored missing `decisionLog` struct/methods, and flattened nested if statements.
- **CLI Help & Troubleshooting Docs:** Authored `cli/helpdoc/split-large-files.md`, `cli/helpdoc/lfs-convert.md`, `cli/helpdoc/flatten.md`, updated `cli/cmdhelp/howto.go`, and added a troubleshooting issue card for 'Large Files / Git Clone Failure (No space left on device)'.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.562.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.562.0/install.sh | sh
```
