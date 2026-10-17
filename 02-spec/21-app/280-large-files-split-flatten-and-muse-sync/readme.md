# Spec 280: Large File Splitting, Git LFS Conversion, Repository Flattening & Muse Sync

## Overview
Comprehensive architecture and toolchain for managing oversized files, Git history bloat, repository flattening, and Muse LLM training sync across:
1. **Large Files Splitting Engine (`gitmap split large files`):** Detects tracked/untracked files $> 10\text{ MB}$, creates multipart zip chunks (first `.zip`, then `.z01`, `.z02`...), applies smart compression (best for data/text, STORE for images/videos/audio), removes bloated originals from Git history, and prompts with yellow irreversible warning.
2. **Git LFS Conversion Suite (`gitmap lfs convert`):** Automates identification of files $> 5\text{ MB}$, updates `.gitattributes`, migrates history to Git LFS pointers, and purges raw blobs.
3. **Repository Flattening Engine (`gitmap repo flatten`):** Creates an incremental numbered successor repo, performs a single squashed flat commit, archives previous repo (`<repo>-marked-for-deletion`), swaps to the canonical name, and provides GitMap `undo` commands.
4. **Muse Training Sync Update:** Updates `gitmap llm train muse` to prioritize `repo-secrets` for active training while synchronizing `repo-cache` for future transition.
5. **Helptext & Web UI Help:** CLI help docs (`split-large-files.md`, `lfs-convert.md`, `flatten.md`), terminal `howto.go`, and web dashboard troubleshooting registries.

## Specifications
- [01-scope-and-plan.md](01-scope-and-plan.md)
- [02-acceptance-criteria-and-evidence.md](02-acceptance-criteria-and-evidence.md)
