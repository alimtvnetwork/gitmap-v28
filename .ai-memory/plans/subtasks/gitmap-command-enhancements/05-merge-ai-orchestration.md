# Subtask 05: Multi-Repo AI Merge Orchestration
Parent Task: gitmap-command-enhancements
Status: COMPLETED

## Objective
Implement `gitmap merge-ai` / `gitmap ma` with single-commit staging, chronological commit sorting, `01_`/`02_` collision sequencing (zero overwrites/deletions), `merge-ai-manifest.json`, and root `instruction.md` consolidation checklist.

## Target Files
- `cli/cmdmergeai/merge_ai_types.go`
- `cli/cmdmergeai/merge_ai_manifest.go`
- `cli/cmdmergeai/merge_ai_core.go`
- `cli/cmdmergeai/merge_ai_test.go`

## Verification
Unit tests pass:
- `TestParseMergeAIArgs`
- `TestStageAndSequenceFilesCollision`
- `TestWriteManifestAndInstruction`
