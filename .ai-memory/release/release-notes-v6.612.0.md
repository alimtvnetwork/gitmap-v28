# GitMap v6.612.0

## What's Changed in v6.612.0

- **Flexible Path Parsing Engine (`cli/macro/expand_flexible.go`):**
  - Added `CleanFlexiblePath` to automatically parse and normalize four distinct agent and markdown path patterns:
    1. Bracketed Windows paths: `@[<drive>:\path\to\media.png]`
    2. Markdown embed tags with file URI: `![Attached Image 1](file:///<drive>:/path/to/media.png)`
    3. Numbered list with repo-relative path and trailing comment: `1. @[assets/media.png] # comment`
    4. Numbered list with parent repo-cache path and trailing comment: `1. @[../repo-cache/assets/media.png] # comment`
  - Strips markdown list bullets (`1. `, `- `, `* `, `+ `), surrounding quotes, and trailing shell/code comments (`# ...`, `// ...`).
  - Resolves relative paths against working directory, repository root (`gitutil.RepoRoot`), and `../repo-cache/` cache hierarchy.
- **Safe Terminal Image Asset Cat (`cli/cmdmacro`):**
  - Upgraded `gitmap cat` to detect image file extensions (`.png`, `.jpg`, `.jpeg`, `.gif`, `.webp`, `.svg`, `.ico`, `.bmp`, `.tiff`) and prevent dumping binary bytes to standard output.
  - Renders a styled TokyoNight card displaying image format, dimensions (`<width> x <height> px`), file size, and disk verification status.
  - For SVG vector files, renders the metadata card followed by clean, formatted SVG text content.
- **Skills & Coding Guidelines Synchronization:**
  - Updated `.agents/skills/gitmap/skill.md` with explicit command runbooks and real-world examples for flexible path parsing and safe image inspection.
  - Updated `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` with media inspection guidance using `gitmap cat`.
  - Updated `.agents/skills/coding-guidelines/skill.md` and `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md` under Anti-Pattern 4 with concrete examples.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.612.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.612.0/install.sh | sh
```
