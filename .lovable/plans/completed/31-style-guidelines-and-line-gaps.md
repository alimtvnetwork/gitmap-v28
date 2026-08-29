# Master Audit: Style Guidelines, Formatting & Line-Gaps (v2.2.0)

## Executive Summary

- **Theme:** Comprehensive source code newline formatting, blank line before `if`, blank line after `}`, blank line before `return`, vertical breathing room around multiline struct calls, flattened conditionals (depth 0), function sizing, file caps, and automated quality gate enforcement.
- **Created Date:** 2026-08-30
- **Completed Date:** 2026-08-30
- **Status:** `COMPLETED`

---

## 1. Architectural Rules & Standards

1. **Rule 1 (Blank Line BEFORE Control Structures):**
   - Exactly one blank line precedes `if`, `for`, `switch`, `while`, `try` when preceded by any statement or assignment.
2. **Rule 2 (Blank Line AFTER Closing Brace `}`):**
   - Exactly one blank line follows closing braces `}` when followed by subsequent statements or invocations.
3. **Rule 3 (Blank Line BEFORE `return` / `throw`):**
   - Clean vertical separation before exit points in multi-line blocks.
4. **Rule 8 (No Leading Blank Line in Function Bodies):**
   - Functions start immediately on line 1 without empty gap after `{`.
5. **Rule 9 (Universal File Hygiene):**
   - Unix LF (`\n`), UTF-8 (no BOM), single trailing newline.

---

## 2. Micro-Batch Refactoring Manifest

40 files refactored across frontend components, hooks, pages, and utility modules with 160 vertical newline insertions:

- `src/components/docs/CloneNextCommandBuilder.tsx`
- `src/components/docs/CodeBlock.tsx`
- `src/components/docs/CommandPalette.tsx`
- `src/components/docs/CopyPaletteButton.tsx`
- `src/components/docs/DocsTooltip.tsx`
- `src/components/docs/SpecPage.tsx`
- `src/components/docs/TabOrderMap.tsx`
- `src/components/docs/TerminalDemo.tsx`
- `src/components/docs/commandsMarkdown.ts`
- `src/components/troubleshooting/TroubleshootingIssueCard.tsx`
- `src/components/projects/ProjectDetailDialog.tsx`
- `src/components/terminal/TerminalView.tsx`
- `src/components/ui/carousel.tsx`
- `src/components/ui/chart.tsx`
- `src/components/ui/form.tsx`
- `src/components/ui/sidebar.tsx`
- `src/hooks/use-toast.ts`
- `src/hooks/useTheme.ts`
- `src/lib/changelogTags.ts`
- `src/lib/clipboard.ts`
- `src/lib/theme.ts`
- `src/pages/Cd.tsx`
- `src/pages/ChromeProfileSpec.tsx`
- `src/pages/CloneNextCommand.tsx`
- `src/pages/Commands.tsx`
- `src/pages/DesignSystem.tsx`
- `src/pages/FlagReference.tsx`
- `src/pages/GenericCLI.tsx`
- `src/pages/PostMortems.tsx`
- `src/pages/ProjectDetection.tsx`
- `src/pages/Projects.tsx`
- `src/pages/Release.tsx`
- `src/pages/ReleaseVersion.tsx`
- `src/pages/ScanCloneFlags.tsx`
- `src/pages/Setup.tsx`
- `src/pages/SpecIndex.tsx`
- `src/pages/Troubleshooting.tsx`
- `src/test/chip-contrast.test.tsx`
- `src/test/new-command-pages.test.ts`
- `src/types/helpJson.ts`

---

## 3. Verification & CI Gates

- `python linter-scripts/check-newline-styling.py` -> 0 violations.
- `python linter-scripts/check-markdown-header-spacing.py` -> all files OK.
- `python linter-scripts/check-nested-ifs.py` -> 0 violations.
- `npm run build` -> exit 0 (built in 4.87s).
- `npm test` -> 9 passed (9 suites, 96 tests).
- `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` -> **23/23 quality gates passed (exit 0)**.
