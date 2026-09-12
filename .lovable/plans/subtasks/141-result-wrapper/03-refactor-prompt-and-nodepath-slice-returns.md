# Subtask 03: Refactor Prompt Targets & NodePath Slice Returns

Parent Plan: [141-result-wrapper-and-slice-returns.md](../../pending/141-result-wrapper-and-slice-returns.md)

## Goals
1. Refactor `cli/cmdprompt/prompt_child_repos.go`:
   - `DiscoverPromptChildRepos(rootDir string) result.ResultSlice[string]`
2. Refactor `cli/cmdprompt/prompt_target_resolver.go`:
   - `ResolvePromptTarget(target string) result.ResultSlice[string]`
3. Refactor `cli/cmdprompt/prompt_workdir_resolver.go`:
   - `ResolveAllWorkDirPromptTargets() result.ResultSlice[string]`
4. Refactor `cli/db/nodepath.go`:
   - `ListPathAliases(db *sql.DB, nodeId string) result.ResultSlice[NodePathAlias]`
5. Refactor `cli/cluster/pathalias.go`:
   - `ParseSetPathAliasArg(raw string) result.ResultSlice[AliasEntry]`
6. Modernize callers in:
   - `cli/cmd/ct.go`
   - `cli/cmdcg/cg_prompts.go`
   - `cli/cmdcg/cg_version_cmd.go`
   - `cli/cmdprompt/prompt_status_cmd.go`
   - `cli/cmdprompt/prompt_version_cmd.go`
   - `cli/cmdprompt/prompt_target_test.go`
   - `cli/cmdprompt/prompt_workdir_resolver.go`

## Acceptance Criteria
- [x] All 5 target functions return `result.ResultSlice[...]`.
- [x] Callers modernized with single object inspection.
- [x] Subtask completed and logged.
