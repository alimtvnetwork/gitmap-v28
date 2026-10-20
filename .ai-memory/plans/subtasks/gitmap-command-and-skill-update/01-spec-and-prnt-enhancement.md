# Subtask 01: Architecture Specs & `gitmap prnt` Batch Engine

> **Subtask ID:** Subtask-01  
> **Status:** DONE  

## Objectives
1. Author specifications in `02-spec/21-app/gitmap-command-and-skill-update/` (all <= 300 lines).
2. Enhance `cli/cmdasset/asset_prnt.go` (and modular helpers <= 300 lines):
   - Support multiple URLs in a single command.
   - Support `--extract-only` (`-e`, `-p`) to print direct image URLs without downloading.
   - Support `--json` (`-j`) structured envelope output.
   - Extend regex for `api.prntscr.com`, `i.imgur.com`, `image.prntscr.com`, and `img.lightshot.app`.
   - Implement concurrent downloading with error isolation.
3. Add unit tests in `cli/cmdasset/asset_prnt_test.go`.
