# Specification: AI Code Execution Restrictions & Letterly Prompt Standards

## 1. Centralized Restrictions Rule File
To prevent AI agents from writing ad-hoc Python SQLite queries, temporary scripts, or bypassing native GitMap commands, a centralized rule file is created:
`.agents/rules/ai_code_execution_restrictions.md`.

### 1.1 Strictly Banned Behaviors
1. **Banned:** Writing or executing ad-hoc Python SQLite queries (`python -c "import sqlite3..."` or `gitmap py -c "import sqlite3..."`).
   - **Mandatory Replacement:** `gitmap db-read <alias|path> "<sql>"`, `gitmap agy summaries`, `gitmap agent task list`.
2. **Banned:** Creating temporary Python scripts in `brain/.../scratch/` to query databases or parse logs.
   - **Mandatory Replacement:** Use native GitMap CLI commands directly.
3. **Banned:** Querying `conversation_summaries.db` with custom scripts.
   - **Mandatory Replacement:** `gitmap agy summaries` or `gitmap db-read agy "<sql>"`.
4. **Banned:** Reading codebase files during task planning or bug remediation without registering intent with GitMap Heatmap.
   - **Mandatory Replacement:** `gitmap heatmap touch <filepath> --reason "<intent>"`.

## 2. Letterly Prompt Standards
- Update `linter-scripts/check-letterly-format.py`:
  - Item 3 must strictly exclude `gitmap py` and `gitmap ps` from search tools.
  - Require Heatmap notification directive: `5. Log file analysis into GitMap Heatmap (gitmap heatmap touch <path> --reason "<why>")`.
- Propagate this contract across canonical prompts in `01-prompts/`.
