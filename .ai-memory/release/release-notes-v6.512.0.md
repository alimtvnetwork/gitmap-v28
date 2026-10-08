# GitMap v6.512.0

## What's Changed in v6.512.0

- **FastGate Staged-Only Pre-Commit Runner (`03-ai-scripts/50-fastgate.py`)**:
  - Scopes AST linters (relative paths, nested ifs, boolean naming) exclusively to staged files in <0.05s (<1.5s target).
  - Eliminates slow full-repo scans prior to commits while providing `--full` fallback mode.
- **Central Suggestion Engine (`cli/suggestion/`)**:
  - Unified domain models (`Suggestion`, `SuggestionGroup`, `SuggestionCategory`, `ActionType`).
  - 4-tier resolution engine: Tier 1 Canonical Aliases, Tier 2 Levenshtein distance $\le 3$, Tier 3 Subsequence/Prefix match, Tier 4 Semantic Synonyms.
  - Catppuccin Macchiato rounded box, compact inline, and structured JSON renderers.
  - Replaces all ad-hoc `fmt.Printf("Did you mean: ...")` console statements across commands.
- **Global Error Interceptor & AppError Integration (`cli/cmd/rootsuggestion.go`, `cli/apperror/`)**:
  - Integrated `WithSuggestions(...)` and `HasSuggestions()` on `AppError`.
  - Global unknown command interceptor displaying Catppuccin suggestion cards and recording telemetry.
- **Muse Semantic Help Text Clustering & Markdown Parsing (`cli/termhelp/`, `cli/cmd/roothelp_clusters.go`)**:
  - 5 Canonical Semantic Clusters: Core & Repo Operations, Release & Commits, Fleet & Remote SSH, AI & Automation, System, OS & Developer Tooling.
  - Data-driven Markdown parser (`termhelp.FromMarkdown`) converting Markdown documentation into terminal cards.
- **Skills & Coding Guidelines Synchronized**:
  - Updated `.agents/skills/gitmap/SKILL.md` with FastGate runner and central suggestion engine commands.
  - Updated `.agents/skills/coding-guidelines/skill.md` with suggestion object standards and pre-commit fastgate protocols.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.512.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.512.0/install.sh | sh
```
