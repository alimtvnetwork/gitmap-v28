# Subtask 06: Magic Constants Elimination & Nested Ifs Flattening

Slug: 06-magic-constants-and-nested-ifs
Parent Plan: 01-coding-guideline-fixes
Status: pending

## Objective
Eliminate all magic strings and numbers into centralized constants packages with typed enums, and flatten nested `if` statements with guard clauses and early returns.

## Concrete Execution Steps (25 Steps)

1. `gitmap/cmd/rootflags.go:120`: Extract magic string `"all"` to `constants.FlagScopeAll`.
2. `gitmap/cmd/rootflags.go:145`: Extract magic number `50` to `constants.DefaultPaginationLimit`.
3. `gitmap/cmd/clone.go:85`: Replace `"depth"` with `constants.FlagCloneDepth`.
4. `gitmap/cmd/code.go:60`: Replace `"code"` with `constants.EditorVSCode`.
5. `gitmap/cmd/pull.go:110`: Replace `"origin"` with `constants.DefaultRemote`.
6. `gitmap/cmd/sync.go:95`: Replace `"main"`, `"master"` with `constants.DefaultBranchMain`, `constants.DefaultBranchMaster`.
7. `gitmap/cmd/history.go:140`: Replace magic limit `100` with `constants.HistoryMaxRecords`.
8. `gitmap/cmd/pendingclear.go:75`: Replace `"completed"` with `constants.PendingTaskStatusCompleted`.
9. `gitmap/cmd/sshgen.go:90`: Replace `"ed25519"` with `constants.SSHKeyTypeEd25519`.
10. `gitmap/cmd/whoami.go:45`: Replace `"user.name"`, `"user.email"` with `constants.GitConfigKeyUserName`, `constants.GitConfigKeyUserEmail`.
11. `gitmap/release/metadata.go:80`: Replace `"v"` with `constants.SemverPrefixV`.
12. `gitmap/startup/add.go:65`: Replace `"registry"` with `constants.StartupBackendRegistry`.
13. `src/pages/Projects.tsx:177`: Flatten nested dialog condition `if (!open) setSelectedProject(null)` into explicit positive handler.
14. `src/pages/SpecIndex.tsx:58`: Flatten `if (!q) return sections` into guard clause.
15. `src/components/docs/TabOrderMap.tsx:66`: Flatten `if (s.pointerEvents === "none" && node === el) { if (...) ... }`.
16. `src/components/docs/TabOrderMap.tsx:132`: Flatten `if (el.labels && el.labels.length > 0) { if (text) ... }`.
17. `src/components/docs/TabOrderMap.tsx:179`: Flatten `if (labelId) { if (ref?.textContent) ... }`.
18. `src/components/docs/TabOrderMap.tsx:288`: Flatten `if (!target) { if (!active) ... }`.
19. `src/components/docs/TerminalDemo.tsx:41`: Flatten `if (isPaused || isFinished) { if (isFinished) ... }`.
20. `src/components/ui/carousel.tsx:59`: Flatten `if (!api) { if (event.key === "ArrowLeft") ... }`.
21. `src/components/ui/chart.tsx:142`: Flatten `if (labelFormatter) { if (!value) ... }`.
22. `src/components/ui/sidebar.tsx:66`: Flatten `if (setOpenProp) { if (event.key === ...) ... }`.
23. `src/hooks/use-toast.ts:90`: Flatten `if (toastId) { if (action.toastId === undefined) ... }`.
24. `gitmap/cmd/installctx_harness_test.go:110`: Flatten nested if with guard clause.
25. `gitmap/cmd/codingguidelines_test.go:20`: Flatten nested if with early return.

## Target Verification Files
- `gitmap/constants/*`
- `gitmap/cmd/*`
- `src/pages/*`
- `src/components/*`
