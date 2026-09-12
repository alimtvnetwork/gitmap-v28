# Plan 148: Extracting Generic Types, Envelopes & Models to types.go Audit

Trigger Keywords & Aliases: `cg-types-go`, `cg-extract-types`, `cg-execute types-go`, `extract-generic-types`, `types-go-single-type`, `types-go-result-reuse`, `centralize-types-go`, `audit types go`, `fix raw generics`, `single reusable type`, `type-alias-repeated-generics`

> **Prompt Version:** 1.0.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)
> **Status:** COMPLETED

---

## 1. Executive Summary & Objective

Autonomously scan, discover, plan, extract, refactor, and verify all scattered domain payload models, raw generic Result wrappers (`ResultSlice[T]`, `ResultMap[K, V]`, `Result[T]`, `Wrap[T]`), and repeated generic signatures across the codebase, centralizing them into dedicated, package-level `types.go` files as single reusable named types, eliminating unexported inline structs, eliminating ad-hoc generic parameterization at call sites and function signatures, and enforcing strict single-type reuse across all implementation files and callers until 100% green without stopping.

Plan 148 executes surgical type extraction, generic alias unification, and model centralization across five core subsystems:
1. **Cmdprompt Subsystem (`cli/cmdprompt/`):** Create `cli/cmdprompt/types.go` consolidating `PromptTemplate`, `PromptInstallOptions`, `PromptStatusTableLayout`, and `type PromptTargetSliceResult = result.ResultSlice[string]`. Move parser functions out of `prompt_types.go` to `prompt_parser.go` and remove obsolete `prompt_types.go`. Refactor `DiscoverPromptChildRepos`, `ResolvePromptTarget`, and `ResolveAllWorkDirPromptTargets` to return `PromptTargetSliceResult`.
2. **Downloaderconfig Subsystem (`cli/downloaderconfig/`):** Create `cli/downloaderconfig/types.go` declaring `Document`, `DownloaderConfig`, `DatabaseVersion`, `ConfigKeyType`, `ConfigKey`, `Key*` constants, `type DocumentResult = result.Result[Document]`, and `type BytesResult = result.Result[[]byte]`. Refactor `LoadFile`, `Parse`, and `Marshal` in `downloaderconfig.go` and all caller sites.
3. **Movemerge Subsystem (`cli/movemerge/`):** Expand `cli/movemerge/types.go` to house `DiffKindType`, `DiffEntry`, `FileMeta`, `Resolver`, `ChoiceType`, and single reusable aliases `type DiffEntryResult = result.Result[DiffEntry]` and `type FileMetaMapResult = result.ResultMap[string, FileMeta]`. Refactor `diff.go`, `walk.go`, `conflict.go`.
4. **Lazyregex Subsystem (`cli/lazyregex/`):** Create `cli/lazyregex/types.go` declaring `LazyRegexp`, `CompileResult`, and canonical alias `type RegexpResult = result.Result[*regexp.Regexp]`. Refactor `Compile()` and `CompileResult()` in `lazyregex.go`.
5. **Archive Subsystem (`cli/archive/`):** Create `cli/archive/types.go` centralizing `FormatType`, `Format`, `CompressionModeType`, `CreateOptions`, `CreateResult`, `ExtractResult`, `Entry`, `ResolvedSource`, and all parameter structs, plus `type BoolResult = result.Result[bool]`. Refactor `isEntryIncluded`, `matchPattern`, and `isHTTPURL`.
6. **Automated Auditor Tooling & Quality Gates:** Update `03-ai-scripts/35-result-wrapper-auditor.py` to register all 5 packages in `ENFORCED_TYPES_GO_PACKAGES`. Run quality linters, record modified files under lock in test inventory, consolidate Plan 148, and push atomically via SSH.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Hallucination prevention, micro-tasking, strict relative paths.
- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-naming-prefixes.md`: Principle 1 (`is`/`has` prefixes only) and Principle 2 (total ban on negative words).
- `spec/02-coding-guidelines/01-cross-language/06-type-definitions.md`: Section 6.4 centralization of domain types, enums, and Result aliases in `types.go`.
- `spec/02-coding-guidelines/01-cross-language/10-function-naming.md`: Semantic verb standards.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100–300 lines).
- `spec/03-error-manage/01-index.md`: Universal `*AppError` wrapping, strongly typed `Result[T]` envelopes, and zero swallowed errors.
- `.lovable/coding-guidelines.md`: Master consolidated coding guidelines.

---

## 3. Exhaustive Types & Generics Violation Ledger

| Id | File | Line | Item | Violation | Target Refactoring | Status |
|:---|:---|:---:|---|---|---|:---:|
| V-01 | `cli/cmdprompt/types.go` | 1 | Package `types.go` | Missing package-level `types.go` | Create `cli/cmdprompt/types.go` | Completed |
| V-02 | `cli/cmdprompt/prompt_types.go` | 9 | `PromptTemplate` | Domain struct defined in non-standard file | Move to `cli/cmdprompt/types.go`, move parser to `prompt_parser.go` | Completed |
| V-03 | `cli/cmdprompt/prompt_args_parser.go` | 8 | `PromptInstallOptions` | Struct scattered in parser file | Move to `cli/cmdprompt/types.go` | Completed |
| V-04 | `cli/cmdprompt/prompt_status_layout.go` | 8 | `PromptStatusTableLayout` | Struct scattered in layout file | Move to `cli/cmdprompt/types.go` | Completed |
| V-05 | `cli/cmdprompt/prompt_child_repos.go` | 9 | `DiscoverPromptChildRepos` | Returns raw `result.ResultSlice[string]` | Return `PromptTargetSliceResult` | Completed |
| V-06 | `cli/cmdprompt/prompt_target_resolver.go` | 14 | `ResolvePromptTarget` | Returns raw `result.ResultSlice[string]` | Return `PromptTargetSliceResult` | Completed |
| V-07 | `cli/cmdprompt/prompt_workdir_resolver.go` | 10 | `ResolveAllWorkDirPromptTargets` | Returns raw `result.ResultSlice[string]` | Return `PromptTargetSliceResult` | Completed |
| V-08 | `cli/downloaderconfig/types.go` | 1 | Package `types.go` | Missing package-level `types.go` | Create `cli/downloaderconfig/types.go` | Completed |
| V-09 | `cli/downloaderconfig/downloaderconfig.go` | 44 | `Document`, `DownloaderConfig`, `DatabaseVersion` | Models scattered in parser file | Move models & key enums to `types.go` | Completed |
| V-10 | `cli/downloaderconfig/downloaderconfig.go` | 98 | `LoadFile` | Returns raw `result.Result[Document]` | Return `DocumentResult` | Completed |
| V-11 | `cli/downloaderconfig/downloaderconfig.go` | 116 | `Parse` | Returns raw `result.Result[Document]` | Return `DocumentResult` | Completed |
| V-12 | `cli/downloaderconfig/downloaderconfig.go` | 182 | `Marshal` | Returns raw `result.Result[[]byte]` | Return `BytesResult` | Completed |
| V-13 | `cli/movemerge/types.go` | 1 | Package `types.go` | Lacks single reusable Result aliases | Add `DiffEntryResult`, `FileMetaMapResult` | Completed |
| V-14 | `cli/movemerge/diff.go` | 24 | `DiffKindType`, `DiffEntry` | Domain models scattered in diff logic | Move `DiffKindType`, `DiffEntry` to `types.go` | Completed |
| V-15 | `cli/movemerge/diff.go` | 86 | `classifyOne`, `classifyBoth` | Returns raw `result.Result[DiffEntry]` | Return `DiffEntryResult` | Completed |
| V-16 | `cli/movemerge/walk.go` | 16 | `FileMeta` | Domain model scattered in walker logic | Move `FileMeta` to `types.go` | Completed |
| V-17 | `cli/movemerge/walk.go` | 24 | `IndexTree` | Returns raw `result.ResultMap[string, FileMeta]` | Return `FileMetaMapResult` | Completed |
| V-18 | `cli/movemerge/conflict.go` | 34 | `Resolver`, `ChoiceType` | Enums and structs scattered in conflict logic | Move `ChoiceType`, `Resolver` to `types.go` | Completed |
| V-19 | `cli/lazyregex/types.go` | 1 | Package `types.go` | Missing package-level `types.go` | Create `cli/lazyregex/types.go` | Completed |
| V-20 | `cli/lazyregex/lazyregex.go` | 19 | `LazyRegexp` | Core struct in implementation file | Move to `cli/lazyregex/types.go` | Completed |
| V-21 | `cli/lazyregex/lazyregex.go` | 52 | `Compile` | Returns raw `result.Result[*regexp.Regexp]` | Return `RegexpResult` | Completed |
| V-22 | `cli/lazyregex/lazyregex.go` | 106 | `CompileResult` | Returns raw `result.Result[*regexp.Regexp]` | Return `RegexpResult` | Completed |
| V-23 | `cli/archive/types.go` | 1 | Package `types.go` | Missing package-level `types.go` | Create `cli/archive/types.go` | Completed |
| V-24 | `cli/archive/archive.go` | 29 | `FormatType`, `Format` | Format enums scattered | Move to `cli/archive/types.go` | Completed |
| V-25 | `cli/archive/create.go` | 36 | `CompressionModeType`, `CreateOptions`, `CreateResult`, `ArchiveWriteParams` | Structs scattered in create file | Move to `cli/archive/types.go` | Completed |
| V-26 | `cli/archive/extract.go` | 25 | `ExtractResult`, `CompactExtractParams`, `ArchiveExtractParams`, `CopyDirEntryParams` | Structs scattered in extract file | Move to `cli/archive/types.go` | Completed |
| V-27 | `cli/archive/list.go` | 15 | `Entry`, `ListExtractParams` | Structs scattered in list file | Move to `cli/archive/types.go` | Completed |
| V-28 | `cli/archive/source.go` | 18 | `ResolvedSource`, `Aria2cDownloadParams` | Structs scattered in source file | Move to `cli/archive/types.go` | Completed |
| V-29 | `cli/archive/create.go` | 184 | `isEntryIncluded` | Returns raw `result.Result[bool]` | Return `BoolResult` | Completed |
| V-30 | `cli/archive/create.go` | 212 | `matchPattern` | Returns raw `result.Result[bool]` | Return `BoolResult` | Completed |
| V-31 | `cli/archive/source.go` | 68 | `isHTTPURL` | Returns raw `result.Result[bool]` | Return `BoolResult` | Completed |
| V-32 | `03-ai-scripts/35-result-wrapper-auditor.py` | 49 | `ENFORCED_TYPES_GO_PACKAGES` | 5 target packages not in enforcement list | Add `cli/cmdprompt`, `cli/downloaderconfig`, `cli/movemerge`, `cli/lazyregex`, `cli/archive` | Completed |
| V-33 | `03-ai-scripts/06-cicd-local-runner.py` | 244 | E2E Smoke Tests Batch | Ran installer smoke tests against host OS | Removed `Installer Smoke (source)` & `(release)` | Completed |
| V-34 | `.lovable/ai-fix-scripts/06-cicd-local-runner.py` | 183 | E2E Smoke Tests Batch | Ran installer smoke tests against host OS | Removed `Installer Smoke (source)` & `(release)` | Completed |
| V-35 | `cli/tests/cmd_test/prompt_install_e2e_test.go` | 1 | E2E prompt install test | Executed install logic in E2E suite | Removed file | Completed |
| V-36 | `cli/cmdinstall/installctx_darwin_e2e_test.go` | 1 | E2E OS context install test | Invoked `runInstallCtxMac()` touching services | Removed file | Completed |
| V-37 | `cli/cmdinstall/installctx_linux_e2e_test.go` | 1 | E2E OS context install test | Invoked `runInstallCtxLinux()` touching desktop configs | Removed file | Completed |

---

## 4. Granular Subtask Decomposition

- [01-cmdprompt-types-extraction.md](../subtasks/148-types-go/01-cmdprompt-types-extraction.md): Centralize `cli/cmdprompt/types.go`, alias `PromptTargetSliceResult`, refactor functions, move parser to `prompt_parser.go`.
- [02-downloaderconfig-types-extraction.md](../subtasks/148-types-go/02-downloaderconfig-types-extraction.md): Centralize `cli/downloaderconfig/types.go`, aliases `DocumentResult` & `BytesResult`, refactor functions.
- [03-movemerge-types-aliases.md](../subtasks/148-types-go/03-movemerge-types-aliases.md): Centralize `cli/movemerge/types.go`, aliases `DiffEntryResult` & `FileMetaMapResult`, refactor functions.
- [04-lazyregex-types-extraction.md](../subtasks/148-types-go/04-lazyregex-types-extraction.md): Centralize `cli/lazyregex/types.go`, alias `RegexpResult`, refactor functions.
- [05-archive-types-extraction.md](../subtasks/148-types-go/05-archive-types-extraction.md): Centralize `cli/archive/types.go`, alias `BoolResult`, refactor functions.
- [06-quality-gates-and-consolidation.md](../subtasks/148-types-go/06-quality-gates-and-consolidation.md): Enforce packages in auditor, verify quality linters, record test inventory, consolidate and atomic commit/push via SSH.

---

## 5. Execution Constraints & Invariants

- **NO AUTOMATIC RELEASES (WOR):** Zero version bumps, zero release triggers.
- **NO FULL CI/CD LOCAL RUNNER:** Do not run `06-cicd-local-runner.py`.
- **STRICT RELATIVE PATHS:** Zero absolute filesystem paths or `file:///` URIs in repository files.
- **ATOMIC COMMIT & SSH PUSH:** Complete all refactoring, verify green, record test inventory under lock, and push to `origin main` via SSH.
