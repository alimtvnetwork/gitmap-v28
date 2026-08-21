# Specification 16 - Part 3: Terminal UI Glyphs, Cross-Platform Emojis & Parallel Clone Progress

## 1. Terminal Emoji & Glyph Rendering Architecture

### 1.1 Root Cause of Broken Glyphs in PowerShell
As shown in user screenshots 1, 4, and 5:
- Screenshot 1 shows `[] parallel clone enabled: 14 workers` where an emoji (e.g. `⚡` or `📦`) was rendered as an empty square `[]` due to font fallback or legacy encoding.
- Screenshot 3 shows `v clean` where `v` was an aggressive ASCII replacement for `✔`.
- Screenshot 4 shows `[] not found` where `✖` or `⚠` was corrupted into an empty square.
- Screenshot 5 shows artifact sections emitting bracketed text strings `[chart]`, `[dna]`, `[tree]`, `[wand]`, `[db]`, `[OK]`, `[tag]`, `[find]` instead of clean visual symbols.

### 1.2 Universal Glyph Table & Safe Fallbacks
We introduce a tiered glyph matrix that detects console capabilities (Windows Terminal, PowerShell 7, ConEmu, VSCode Terminal, Linux/macOS) and selects high-compatibility Unicode symbols that render consistently without broken boxes:

| Intent | Primary UTF-8 Symbol | Secondary Console Symbol | ASCII Fallback |
|---|---|---|---|
| Success / Clean | `✔` (U+2714) | `√` (U+221A) | `[OK]` |
| Missing / Failed | `✖` (U+2716) | `×` (U+00D7) | `[X]` |
| Warning / Caution | `▲` (U+25B2) | `!` | `[!]` |
| Folder / Scope | `📂` (U+1F4C2) | `»` (U+00BB) | `[DIR]` |
| Database | `◆` (U+25C6) | `*` | `[DB]` |
| Network / Parallel | `⚡` (U+26A1) | `~` | `[//]` |
| Bullet / Item | `•` (U+2022) | `-` | `*` |
| Active / Processing | `●` (U+25CF) | `o` | `(o)` |

### 1.3 Windows Console Encoding Enforcement
At application startup, the runtime initializes:
1. `SetConsoleOutputCP(CP_UTF8)` on Windows via system syscalls.
2. Console capability probe verifying `WT_SESSION` (Windows Terminal), `TERM_PROGRAM`, and standard stream TTY status.
3. Elimination of hardcoded ASCII bracket strings (`[dna]`, `[tree]`, `[wand]`) in `scanoutput.go` and `clonetermblock.go`.

---

## 2. Parallel Clone Progress Visualization

### 2.1 Problem Analysis (Screenshot 1)
Screenshot 1 shows the current defect during `gitmap clone gitmap.json`:
```text
[ 1/14] Cloning alim-status-sample-v2 ...[ 2/14] Cloning cat-my-v12 ...[ 3/14] Cloning coding-guidelines-v24 ...
```
All worker goroutines write directly to stdout concurrently without mutual exclusion or line-feed separation, causing ANSI cursor codes and strings to smash together horizontally.

### 2.2 Re-architected Progress Display
The parallel cloner MUST route all worker status updates through a synchronized UI coordinator:

1. **Option A: Synchronized Sequential Log Mode (CI / non-interactive / redirect-friendly)**
   Each worker emits complete, atomic lines with worker ID, percentage/index, repository name, status, and elapsed duration:
   ```text
   PS D:\work> gitmap clone gitmap.json
   Existing repos detected — safe-pull enabled automatically.
   ⚡ Parallel clone active: 14 workers (14 repositories)

   [ 1/14] 📂 Cloning alim-status-sample-v2 ... ✔ done (1.2s)
   [ 2/14] 📂 Cloning cat-my-v12 ............. ✔ done (1.8s)
   [ 3/14] 📂 Cloning coding-guidelines-v24 .. ✔ done (0.9s)
   [ 4/14] 📂 Updating gitmap-v28 (existing) . ✔ up-to-date (0.4s)
   ...
   ────────────────────────────────────────────────────────────
   ✔ Clone complete: 14/14 succeeded · 1 existing updated · Elapsed: 4.2s
   ```

2. **Option B: Multi-Row Interactive Dynamic Terminal (TTY Mode)**
   Using an isolated multi-row terminal printer where each active worker occupies a designated row that updates in-place, closing cleanly upon completion.

### 2.3 Existing vs. Fresh Repository Handling
- **Fresh Repository**: Executes `git clone <url> <path>` with depth and branch settings. Status shows `✔ cloned (<duration>)`.
- **Existing Repository**:
  - Automatically activates `--safe-pull` mode.
  - Checks if remote matches; if matching, executes `git pull --ff-only` (or rebase if configured).
  - If worktree is dirty, safely skips or stashes with a clear warning: `▲ skipped (dirty worktree)`.
  - Status shows `✔ safe-pulled (<duration>)` or `✔ up-to-date (<duration>)`.
