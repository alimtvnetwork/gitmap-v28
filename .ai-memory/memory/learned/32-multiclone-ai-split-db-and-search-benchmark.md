# Learned Memory 32: MultiClone, AI Split DB Command Tracking & Search Benchmark

## 1. MultiClone (`gitmap multiclone`, `mc`, `mutliclone`) Subsystem
- **Markdown & Paste Parsing**: Accepts multiline text pasted from browsers, chat terminals, or text documents enclosed in markdown fences (```` ``` ````) or raw lines.
- **URL & Shorthand Extraction**:
  - Full Git URLs (`https://...`, `http://...`, `git@...`, `ssh://...`)
  - Markdown links (`[title](url)`)
  - `owner/repo` shorthands (e.g. `ChrisTitusTech/linutil: ...`), stripping trailing descriptions and expanding to `https://github.com/owner/repo.git`.
- **Deduplication**: Case-insensitive normalization, stripping `.git` and trailing slashes so multiple references to the same repository clone exactly once.
- **Target Folder Mapping**: Optional destination folder (`-d <dir>` or positional folder arg) routes repos into `<dir>/<repo-name>`.
- **Two-Column Styled Help**: Rich terminal help menu formatted via `termhelp.RenderMenu` matching GitMap standards.

## 2. AI Split Database & Command Tracking
- **Dedicated Split Database**: `~/.gitmap/ai/instructions.db` (`store.AiInstructionSplitDB`).
- **Global Tracking Flag (`--ai`)**:
  - Automatically stripped from CLI arguments before command execution and recorded via environment variable `GITMAP_AI_TRACKING=1`.
  - On command completion, logs execution details (command line, duration, exit code, timestamp) to `AiExecutionHistory`.
- **Frequent Command Inspection (`gitmap ai ls` / `gitmap ai history`)**:
  - Displays top frequent AI commands with columns: `#`, `RUNS`, `SUCCESS`, `COMMAND`, `LAST RUN`.
- **Clipboard Sync (`--copy`)**:
  - `gitmap ai ls --copy` or `gitmap ai history --copy` copies frequent commands directly to the OS clipboard via `atotto/clipboard`.

## 3. Automation Search Benchmark vs Python
- **Side-by-Side Performance Verification**:
  - Go Native Search (`gitmap automation search` / `cmdautomation.RunSearch`): **1.5ms**
  - Python Script Search (`03-ai-scripts/12-fast-cached-grep.py`): **70.5ms**
  - Result: Go native automation search is **46.1x faster** with ZERO temporary disk bloat.
