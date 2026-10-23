# Execution Plan: Gitmap Performance and CI/CD Issue Resolution

- **Task Slug**: `gitmap-performance-and-cicd-issue-resolution`
- **Parent Task ID**: `task-20261019141450-gitmap-performance-and-cicd-issue-resolution`
- **Specification**: `02-spec/21-app/gitmap-performance-and-cicd-issue-resolution/01-architecture-spec.md`
- **Status**: COMPLETED
- **Budget**: 400 Steps
- **Capacity**: A = 3 agents, H = 2 hands, Concurrency = 6

---

## 1. Verbatim User Request

```text
# Gitmap Performance and CI/CD Issue Resolution: high priority instruction, non-negotiable task

The Gitmap PE issues, all the issues, and also make sure that GitmapTE gives the result faster. It takes way longer time and always give the final time how long it took to provide this output. Timing is important, how long it takes to generate the support. And try to use parallel processing as much as possible. For example, downloading, every downloading should be parallel. And inside the processing of the pipeline, it should be also lines. If you have multiple lines to process, you can also do those processing parallelly so that the processing happens very quick. And if you have the cache data, show the cache data very quick, and also at the end show how long it take to complete the request. Fix all the CI/CD issues for now. Try to understand the root cause, the root cause and then fix it.
```

---

## 2. Discrete Subtask Allocations

| Subtask ID | Title | Owned Files | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | Fix Legacy Refs in Readme Conventions Spec and Pass CI Script Tests | `02-spec/01-spec-authoring-guide/13-root-readme-conventions.md` | Worker 01 | DONE | `python .github/scripts/check-legacy-refs.py` (exit 0) |
| **Subtask-02** | Decompose `cli/constants/constants_install.go` Below 300-Line Cap | `cli/constants/constants_install.go`, `cli/constants/constants_install_tools.go`, `cli/constants/constants_install_pkgs.go` | Worker 02 | DONE | `constants_install*.go` all $\le 300$ lines, tests pass |
| **Subtask-03** | Decompose `cli/cmdinstall/install_packages.go` and Tests Below 300-Line Cap | `cli/cmdinstall/install_packages.go`, `cli/cmdinstall/install_packages_resolve.go`, `cli/cmdinstall/install_packages_test.go`, `cli/cmdinstall/install_packages_resolve_test.go` | Worker 03 | DONE | `52-file-size-check.py` (exit 0), `go test ./cmdinstall` (exit 0) |
| **Subtask-04** | Accelerate Pipeline Log Downloading and Fast Cache Return in Gitmap TE/PE | `cli/cmdpipeline/pipeline_logs_fetch.go`, `cli/cmdpipeline/pipeline_query_logs.go`, `cli/cmdpipeline/pipeline_sha_cache.go` | Worker 01 | DONE | Dynamic worker pool (up to 8 workers), direct `.git/HEAD` disk resolution (<0.1ms), tests pass |
| **Subtask-05** | Implement Parallel Line Processing and Chunked Log Parsing in Gitmap TE/PE | `cli/cmdpipeline/pipeline_extract_scan.go`, `cli/cmdpipeline/pipeline_extract_parallel.go`, `cli/cmdpipeline/pipeline_extract_clean.go` | Worker 02 | DONE | Multi-core chunked parsing in `pipeline_extract_parallel.go`, 100% test parity with sequential parsing |
| **Subtask-06** | Always Display Elapsed Execution Time in Gitmap TE/PE Output and Verify All Quality Gates | `cli/cmdpipeline/pipeline_logs_output.go`, `cli/cmdpipeline/pipeline_logs_entry.go`, `cli/cmdpipeline/pipeline_timing_test.go` | Worker 03 | DONE | Elapsed timing banner output (`⚡ [gitmap te] Request completed in ...ms (from cache)`), all tests pass |

---

## 3. Verification Contracts

- **Subtask-01**: `check-legacy-refs.py` finds 0 occurrences; `test_ci_scripts.py` passes cleanly (19/19 tests passed).
- **Subtask-02**: `constants_install.go` and sibling files each strictly $\le 300$ lines; all tests referencing constants compile.
- **Subtask-03**: `install_packages.go` and tests each strictly $\le 300$ lines; `52-file-size-check.py` passes.
- **Subtask-04**: Downloading runs concurrently via worker pool; cached SHA queries return immediately (<20ms) without subprocess spawns.
- **Subtask-05**: Large logs parsed across parallel goroutines with zero data races; structured output identical to sequential parse.
- **Subtask-06**: `gitmap te` outputs total elapsed duration at end (`⚡ [gitmap te] Request completed in ...ms (from cache)`); all tests pass.
