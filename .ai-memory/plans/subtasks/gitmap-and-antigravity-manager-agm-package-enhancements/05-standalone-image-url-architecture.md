# Subtask 05: Standalone Image URL & Lightweight Prompt Architecture

**Parent Task:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Subtask ID:** `subtask-05`  
**Status:** `IN_PROGRESS`  
**Owned Files:**  
- `02-spec/21-app/gitmap-and-antigravity-manager-agm-package-enhancements/03-image-url-architecture.md`
- `cli/cmdagy/agy_image_assets.go`

## Objective
Design and implement the standalone image URL and content-addressable storage architecture to eliminate multi-megabyte inline base64 image blobs from prompt files, terminal logs, and AI context windows.

## Implementation Details
1. Architect the Content-Addressable Asset Server (`http://127.0.0.1:18188/assets/<sha256>.<ext>`):
   - Project-scoped storage in `.ai-memory/assets/` and machine-global cache in `~/.gitmap/assets/`.
   - Dual-reference markdown structure:
     ```markdown
     ![Screenshot](http://127.0.0.1:18188/assets/img-a1b2c3d4.png)
     <!-- asset: .ai-memory/assets/img-a1b2c3d4.png -->
     ```
2. Implement `cli/cmdagy/agy_image_assets.go`:
   - `ExtractAndStoreBase64Images(prompt string, targetDir string) (string, []string, error)`:
     Extracts all inline `data:image/...;base64,...` URIs, decodes binary data, calculates SHA-256 hash, saves to `.ai-memory/assets/img-<sha256>.<ext>`, and rewrites the prompt with clean standalone URLs.
3. Author detailed architecture spec `03-image-url-architecture.md`.
