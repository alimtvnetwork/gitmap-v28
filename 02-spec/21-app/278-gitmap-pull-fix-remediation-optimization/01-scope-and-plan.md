# Spec 278: Scope & Technical Plan

## User Request (Verbatim)
```text
PS D:\work\story-telling> gitmap cfr https://github.com/alimtvnetwork/story-telling-v1
Cloning story-telling-v1 into story-telling...
...
fatal: cannot create directory at 'story/audio/chhotto-lal-murgi': No space left on device
warning: Clone succeeded, but checkout failed.
Error: clone failed for https://[REDACTED_GH_TOKEN]@github.com/alimtvnetwork/story-telling-v1: exit status 128 (operation: git-clone)

PS D:\work\story-telling> gitmap reclone https://github.com/alimtvnetwork/story-telling-v1
gitmap reclone: execute failed: [E9000:EXECUTION] parse-manifest cfg.file: clone-now: unsupported file extension "" for D:\work\story-telling\https:\github.com\alimtvnetwork\story-telling-v1; supported extensions are .json, .csv, .txt (or pass --format json|csv|text) (at=cmdclone/clonenow.go:106)
```

## Architectural Scope & Analysis

### 1. GCM Windows Credential Store Fix
- **Root Cause:** `wincred` is an invalid token for Git Credential Manager. GCM supports `wincredman`, `dpapi`, `plaintext`, `none`. Setting `credential.credentialStore` to `wincred` caused `fatal: Unknown credential store 'wincred'` on every git invocation during pull or push.
- **Remediation:** Update `cli/cmdfixauth/fixcredential.go`, `cli/cmdpull/pull_remediation.go`, and `cli/gitutil/git_env.go` to use `wincredman` (or `dpapi` in headless SSH sessions).

### 2. Single Repo & URL Support in `gitmap reclone`
- **Root Cause:** `tryRunRepoReclone` in `cli/cmdclone/reporeclone.go` only checked if the argument was a local directory containing `.git`. When passed a remote URL like `https://github.com/alimtvnetwork/story-telling-v1`, it tried `filepath.Abs(url)`, failed `os.Stat`, and fell back to `clonenow.ParseFile` which crashed expecting a `.json`/`.csv`/`.txt` manifest.
- **Remediation:** Detect when the argument is a Git repository URL (`http://`, `https://`, `git@`, `ssh://`, or ending in `.git`). Extract the repository slug, detect if target folder exists in current working directory, wipe and re-clone cleanly.

### 3. Parallel Generic `fix all` Worker Pool
- **Root Cause:** `runFixAll` executed remediation recipes in a serial for-loop.
- **Remediation:** Implement a generic typed concurrent worker pool in Go (`cli/parallel/generic_pool.go` or `cli/cmdfix/`) that auto-scales based on available CPU logical cores (minimum 4, up to `runtime.NumCPU()`, e.g. 20 cores). Provide real-time spinning indicators and progress updates.

### 4. `gitmap pull-all` Speed & Post-Pull Resolution
- **Root Cause:** Concurrency was artificially capped at 6 even on 20-core machines, and stdout streaming from background child git processes polluted the terminal, causing ANSI escape truncation and lost bar colors.
- **Remediation:**
  - Auto-scale worker pool up to 20 based on CPU & RAM availability.
  - Buffer subprocess output during parallel pulls to preserve progress bar coloring.
  - Reorder post-pull interactive menu:
    - `[1/f/fix] Fix All (Parallel auto-remediation via gitmap fix --all)`
    - `[2/c/commit] Commit All (Local WIP commit on this machine only)`
    - `[3/t/stash] Stash All (git stash on all unresolved repos)`
    - `[4/s/step] Step through each repository individually (with tree-view diff & action preview)`
    - `[5/fleet] Fleet Multi-Repo Control: Run gitmap sends cp all "wip: save changes"`
    - `[6/q/skip] Skip and leave as-is`

### 5. `gitmap update` Daily Probe Throttle & Visual Polish
- **Requirements:**
  - Probe sibling repos at most once per 24 hours. Cache the probe result in a temp/state file (`.gitmap/data/update/probe-watermark.json` or OS temp).
  - Show version transition clearly with previous version in yellow/white.
  - Title casing (`Gitmap CLI Update`) with left-margin padding.

### 6. `repo-cache` Integration & Muse AI Memory Migration
- **Requirements:**
  - Auto-create and configure `repo-cache` repository on GitHub and locally in `d:\work\repo-cache`.
  - Migrate `12-muse-ai/` memory digest from `repo-secrets` to `repo-cache`.
  - Update `cli/cmd/llm/llm_train_muse.go` to train from `repo-cache` (with fallback to `repo-secrets`).
  - Support `XX-Muse-AI` directory structure.

### 7. Space Command & Subcommand Alias Hygiene
- **Requirements:**
  - Wire hyphenated aliases (`llm-train-muse`, `repo-export`, `repo-import`, `mr-list`, `mr-view`, `mr-merge`, `mr-checkout`, `scan-diff`, `scan-verify`) so users and scripts never depend on literal whitespace strings.

### 8. Heatmap Skill Documentation
- **Requirements:**
  - Document `gitmap heatmap`, `gitmap agent heatmap`, and `gitmap full-summary` in `.agents/skills/gitmap/SKILL.md` and repo guidelines so autonomous agents use GitMap heatmaps before modifying files.
