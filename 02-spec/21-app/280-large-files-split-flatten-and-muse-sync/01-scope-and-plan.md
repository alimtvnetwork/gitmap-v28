# Spec 280: Scope & Technical Plan

## User Request (Verbatim)
```text
PS D:\work\story-telling> gitmap cfr https://github.com/alimtvnetwork/story-telling-v1
Cloning story-telling-v1 into story-telling...
  [clone] target free, cloning directly into D:\work\story-telling\story-telling

  ▸ git clone  https://[REDACTED_GH_TOKEN]@github.com/alimtvnetwork/story-telling-v1
    target  D:\work\story-telling\story-telling
    exec    git clone https://[REDACTED_GH_TOKEN]@github.com/alimtvnetwork/story-telling-v1 D:\work\story-telling\story-telling
Cloning into 'D:\work\story-telling\story-telling'...
...
fatal: cannot create directory at 'story/audio/chhotto-lal-murgi': No space left on device
warning: Clone succeeded, but checkout failed.
You can inspect what was checked out with 'git status'
and retry with 'git restore --source=HEAD :/'
```

## Architectural Scope & Analysis

### 1. Root Cause Analysis: Storytelling Large File Clone Failure
- **Root Cause:** The `alimtvnetwork/story-telling-v1` repository contains 30 complete story and learning audio narrations, image illustrations, and multi-segment mp3 files. The raw Git packfile exceeds 194.43 MiB with hundreds of binary blobs. During `git clone`, unpacking these massive binary files simultaneously on Windows systems triggers `No space left on device` and file lock / checkout aborts.
- **Muse Context:** In `repo-secrets/12-muse-ai/`, Muse executed an audio-first pipeline producing dozens of high-quality `.mp3` and `.png` assets. Storing large binary assets directly in Git commit history without splitting or LFS causes repository bloat and clone failures.

### 2. Large File Splitting Engine (`gitmap split large files`)
- **Threshold:** Default $> 10\text{ MB}$ (configurable via `--size <N>MB`).
- **Multipart Zip Architecture:**
  - File is packaged into multipart zip chunks.
  - Chunk 1 is named `<filename>.zip`. Subsequent parts are `<filename>.z01`, `<filename>.z02`, etc., compatible with 7-Zip, WinRAR, and standard unzippers.
  - Compression mode heuristic:
    - Text, code, documents, data: `zip.Deflate` (Best compression).
    - Pre-compressed media (`.mp3`, `.wav`, `.png`, `.jpg`, `.mp4`, `.mov`): `zip.Store` (No recompression, fast chunking).
- **Git History Rewrite:**
  - Removes the raw large binary from Git index and history.
  - Adds split parts to Git tracking.
- **User Safety:**
  - Dry run mode (`--dry-run`) shows candidate files, sizes, projected chunks, and compression ratios.
  - Execution requires interactive user confirmation or `-y` / `--yes` flag.
  - Displays terminal warning in bold yellow:
    `WARNING: This operation rewrites Git history and splits files into multipart zip archives. It is irreversible.`

### 3. Git LFS Conversion (`gitmap lfs convert`)
- **Threshold:** Default $> 5\text{ MB}$ (configurable via `--size <N>MB`).
- **LFS Tracking & Migration:**
  - Identifies tracked files exceeding 5MB.
  - Configures `.gitattributes` tracking patterns.
  - Migrates Git commit history to LFS pointers using `git lfs migrate import` or manual blob rewriting.
  - Provides dry-run preview and interactive `-y` warning.

### 4. Repository Flattening (`gitmap repo flatten`)
- **Purpose:** Compresses a bloated repository history into a single clean snapshot commit.
- **Workflow:**
  1. Creates temporary staging directory `<repo>-flat-v1`.
  2. Copies all working tree files (excluding `.git`, node_modules, build artifacts).
  3. Initializes a fresh git repository with a single flat squashed commit (`flat: initial squashed repository snapshot`).
  4. Archives the bloated original repository to `<repo>-marked-for-deletion`.
  5. Renames `<repo>-flat-v1` to `<repo>`.
  6. Provides clear GitMap undo command: `gitmap repo undo-flatten <repo>` restoring the original directory.

### 5. Muse LLM Training Source Synchronization
- Updates `defaultResolveTrainingPath` in `cli/cmd/llm/llm_train_muse.go` to prioritize `repo-secrets` for active training operations, while supporting `repo-cache` synchronization.
- Resolves undefined `decisionLog` error identified during CI pipeline audit (`gitmap pe`).
