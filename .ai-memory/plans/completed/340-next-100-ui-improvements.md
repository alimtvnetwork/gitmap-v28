# Master Plan: Next 100 UI Improvements

- **Task Slug**: `340-next-100-ui-improvements`
- **Total Steps**: 400
- **Status**: IN_PROGRESS
- **Architecture Spec**: `02-spec/21-app/340-next-100-ui-improvements/00-master-audit-ledger.md`

---

## 1. Subtask Decomposition

| Subtask ID | File | Description | Target Files |
| :--- | :--- | :--- | :--- |
| **Subtask-01** | `01-tokens-contrast-and-depth.md` | Tokens, contrast calibration, 4-plane depth & card/dialog hairline lighting (Clusters 01–03) | `src/index.css`, `tailwind.config.ts`, `src/components/ui/card.tsx`, `src/components/ui/dialog.tsx` |
| **Subtask-02** | `02-hierarchy-typography-and-buttons.md` | Single-accent Von Restorff rule, Poppins/JetBrains Mono typography & 6-variant button system (Clusters 04–06) | `src/components/ui/button.tsx`, `src/pages/Index.tsx`, `src/pages/Pipeline.tsx`, `src/pages/Settings.tsx` |
| **Subtask-03** | `03-navigation-motion-and-hygiene.md` | Sticky 72px header, GPU motion, safe-region physics, and component deconstruction (Clusters 07–10) | `src/components/docs/DocsLayout.tsx`, `src/components/docs/DocsSidebar.tsx`, `cli/cmdui/ui_assets_theme.go` |

---

## 2. Invariants & Rules

1. **Strictly Relative Git Paths**: All links and file paths must be relative (e.g. `src/...`, `02-spec/...`).
2. **Zero Builds & Tests**: R1 rule prohibits running `npm run build`, `vite build`, or test suites during execution.
3. **Strict Function & Component Limits**: React components $\le 80$–100 lines; functions $\le 8$–15 lines.
4. **Positive Booleans Only**: Positive variable and predicate naming (`isReady`, `hasTokens`).
5. **Atomic GitMap Commit**: Commit with `gitmap cpf` using hyphen format upon wave completion.
