# Subtask 148.1: Cmdprompt Types Extraction & Result Centralization

> **Parent Plan:** [148-types-go-extraction-and-generic-result-centralization.md](../../pending/148-types-go-extraction-and-generic-result-centralization.md)
> **Target Subsystem:** `cli/cmdprompt/`
> **Status:** COMPLETED

## Scope of Work

1. Create `cli/cmdprompt/types.go`:
   - Declare `PromptTemplate` struct with JSON tags.
   - Declare `PromptInstallOptions` struct with affirmative booleans (`IsDryRun`, `IsAll`).
   - Declare `PromptStatusTableLayout` struct.
   - Declare canonical alias `type PromptTargetSliceResult = result.ResultSlice[string]`.
2. Move markdown parsing functions from `prompt_types.go` to `prompt_parser.go`:
   - `parsePromptMarkdown`, `parseFrontmatterLines`, `parseFrontmatterLine`, `assignPromptBody`, `finalizePromptTemplate`.
   - Safely remove `prompt_types.go`.
3. Refactor function signatures in `cli/cmdprompt/`:
   - `DiscoverPromptChildRepos(rootDir string) PromptTargetSliceResult` in `prompt_child_repos.go`.
   - `ResolvePromptTarget(target string) PromptTargetSliceResult` in `prompt_target_resolver.go`.
   - `ResolveAllWorkDirPromptTargets() PromptTargetSliceResult` in `prompt_workdir_resolver.go`.
4. Verify call sites and tests in `cli/cmdprompt/` compile cleanly with `go test ./cli/cmdprompt/...`.
