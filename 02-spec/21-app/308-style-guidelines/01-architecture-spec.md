# Architecture Spec: 308 - Coding Style, Vertical Newline Rhythm & Return New Line Concept

## 1. Executive Summary & Intent
This specification formalizes compliance with `01-prompts/15-cg-execute/10-style-guidelines.md` and repository guidelines (`02-spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md`).
The Return New Line and block spacing guidelines mandate:
1. Blank line before every `return` statement (unless it is the only statement in a function body).
2. Blank line after every closing brace `}` if followed by further statements.
3. Zero nested `if` statements and affirmative boolean evaluations.

## 2. Spacing Transformation Rules
```mermaid
flowchart TD
    subgraph SpacingRules [Newline Rhythm Rules]
        A[Code Statement] --> B[Blank Line]
        B --> C[return result]
        D[Closing Brace }] --> E[Blank Line]
        E --> F[Following Statement]
    end
```

## 3. Disjoint Micro-Batch Allocation
- **Worker 01 (AGYPrompts Return Spacing):**
  - File: `src/pages/AGYPrompts.tsx`
  - Action: Add a blank line before `return` in `filteredPrompts` (line 253).
- **Worker 02 (Database Closing Brace Spacing):**
  - File: `src/pages/Database.tsx`
  - Action: Add blank lines after closing `}` on lines 140, 454, 607, 610, 621, 640, 643, 646, 746.
- **Worker 03 (MergeAIPage Return Spacing):**
  - File: `src/pages/MergeAIPage.tsx`
  - Action: Add blank lines before `return` on lines 200, 214, 219, 224.

## 4. Verification & Quality Gates
- `FastGate`: `python 03-ai-scripts/50-fastgate.py` must report clean.
- `Newline Styling`: `python linter-scripts/check-newline-styling.py --all` must report PASS (0 violations).
- `Forbidden Strings`: `python linter-scripts/check-forbidden-strings.py` must report PASS.
