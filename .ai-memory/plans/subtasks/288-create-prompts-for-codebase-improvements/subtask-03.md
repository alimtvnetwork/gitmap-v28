# Subtask 288.3: Develop Prompt 2: Find 100 Improvements in the Codebase

- **Parent Plan:** [.ai-memory/plans/pending/288-create-prompts-for-codebase-improvements.md](../../pending/288-create-prompts-for-codebase-improvements.md)
- **Status:** Completed
- **Completed At:** 2026-10-22
- **File Boundaries:**
  - `01-prompts/30-improvements/02-find-100-improvements-in-codebase-coding-guidelines.md`

## Objectives & Outcomes
1. Created `01-prompts/30-improvements/02-find-100-improvements-in-codebase-coding-guidelines.md`.
2. Structured the prompt across 10 distinct guideline categories (10 items each = 100 items).
3. Enforced the mandatory 7-field reporting schema:
   - File Path & Line Range (strictly relative)
   - Function / Symbol Name
   - Existing Guideline Violation (citing official rules)
   - Context & Architectural Rationale
   - What to Improve & How to Improve
   - Concrete Code Refactoring (Before vs. After)
   - Verification Criterion
4. Embedded the 6 core coding standards: Error Return Sovereignty, 15-line function limit, zero nested ifs, affirmative booleans, Split-DB repository encapsulation, and strongly typed `*Type` enums.
5. Added grounded example item and strict non-negotiable execution rules.
