# Subtask 04: File Size Caps & React Component Decomposition

Slug: 04-file-size-caps-and-react-splits
Parent Plan: 01-coding-guideline-fixes
Status: pending

## Objective
Enforce file size limits ($\le 300$ lines max per file, $\le 100$ lines max per React component `.tsx`, $\le 120$ lines per struct/class).

## Concrete Execution Steps (25 Steps)

1. `gitmap/scanner/scanner.go` (469 lines): Split into:
   - `scanner_types.go`: Definitions for `RepoInfo`, `ScanOptions`, `dirJob`.
   - `scanner_sniff.go`: Git directory sniffing and worktree detection helpers.
   - `scanner_walker.go`: Concurrency walker and worker pool dispatch logic.
   - `scanner.go`: Clean entry point ($\le 120$ lines).
2. `gitmap/cmd/clone.go` (609 lines): Split into `clone_flags.go`, `clone_interactive.go`, and `clone_exec.go`.
3. `gitmap/cmd/cluster_ops.go` (496 lines): Split into `cluster_format.go`, `cluster_dispatch.go`, and `cluster_ops.go`.
4. `gitmap/store/store.go` (476 lines): Split into `store_connection.go`, `store_schema.go`, and `store.go`.
5. `gitmap/cmd/chromeprofile_merge.go` (453 lines): Split into `chromeprofile_rules.go` and `chromeprofile_merge.go`.
6. `gitmap/cmd/clonenext.go` (430 lines): Split into `clonenext_parse.go` and `clonenext_exec.go`.
7. `gitmap/cmd/rootflags.go` (426 lines): Split into `rootflags_core.go` and `rootflags_extended.go`.
8. `gitmap/cmd/code.go` (422 lines): Split into `code_detect.go` and `code_launch.go`.
9. `gitmap/cmd/pull.go` (411 lines): Split into `pull_args.go` and `pull_runner.go`.
10. `gitmap/cmd/selfuninstallparts.go` (400 lines): Split into platform-specific files.
11. `gitmap/cmd/selfinstall.go` (400 lines): Split into `selfinstall_setup.go` and `selfinstall_binary.go`.
12. `gitmap/cmd/installtools.go` (392 lines): Split into individual tool installers.
13. `gitmap/cmd/sync.go` (385 lines): Split into `sync_scan.go` and `sync_apply.go`.
14. `gitmap/cmd/clonefixrepo.go` (382 lines): Split into `clonefixrepo_prep.go` and `clonefixrepo_run.go`.
15. `gitmap/archive/extract.go` (344 lines): Split into `extract_archive.go` and `extract_filesystem.go`.
16. `src/pages/GenericCLI.tsx` (1107 lines): Extract sub-components to `src/components/docs/generic-cli/`:
   - `CLICommandList.tsx` ($\le 100$ lines)
   - `CLICommandFilter.tsx` ($\le 100$ lines)
   - `CLIPreviewPanel.tsx` ($\le 100$ lines)
17. `src/pages/Release.tsx` (628 lines): Extract sub-components to `src/components/docs/release/`:
   - `ReleaseHero.tsx` ($\le 100$ lines)
   - `ReleaseAssetTable.tsx` ($\le 100$ lines)
   - `ReleaseVersionList.tsx` ($\le 100$ lines)
18. `src/pages/Install.tsx` (530 lines): Extract sub-components to `src/components/docs/install/`:
   - `InstallHero.tsx` ($\le 100$ lines)
   - `InstallScriptTabs.tsx` ($\le 100$ lines)
   - `InstallQuickSnippet.tsx` ($\le 100$ lines)
19. `src/pages/ProjectDetection.tsx` (449 lines): Extract sub-components to `src/components/docs/project-detection/`:
   - `DetectionFilterBar.tsx` ($\le 100$ lines)
   - `DetectionProjectCard.tsx` ($\le 100$ lines)
20. `src/pages/CloneNextCommand.tsx` (427 lines): Extract sub-components to `src/components/docs/clone-next/`:
   - `CloneNextOptionsForm.tsx` ($\le 100$ lines)
   - `CloneNextPreviewCard.tsx` ($\le 100$ lines)
21. `src/pages/DesignSystem.tsx` (397 lines): Extract sub-components to `src/components/docs/design-system/`:
   - `ColorPaletteGrid.tsx` ($\le 100$ lines)
   - `TypographyShowcase.tsx` ($\le 100$ lines)
22. `src/pages/CloneNext.tsx` (370 lines): Extract `CloneNextGuide.tsx` and `CloneNextDemo.tsx`.
23. `src/pages/BatchActions.tsx` (331 lines): Extract `BatchActionTable.tsx` and `BatchActionControls.tsx`.
24. `src/components/docs/TabOrderMap.tsx` (490 lines): Extract `TabOrderOverlay.tsx` and `TabOrderKeyNav.tsx`.
25. `src/components/docs/CodeBlock.tsx` (382 lines): Extract `CodeBlockHeader.tsx` and `CodeBlockLines.tsx`.

## Target Verification Files
- `gitmap/scanner/*`
- `gitmap/cmd/*`
- `src/pages/*`
- `src/components/docs/*`
