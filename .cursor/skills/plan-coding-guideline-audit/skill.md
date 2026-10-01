---
name: plan-coding-guideline-audit
description: Plan a structured coding guideline audit across repository codebases against 02-spec/02-coding-guidelines/.
---

# Plan Coding Guideline Audit

Autonomously plans a comprehensive audit of repository codebases against the master coding guidelines.

## Audit Areas

1. **Boolean Conventions & Affirmative Parameters/Fields:** `is` and `has` prefixes only, implicit positive checks, no mixed polarity, affirmative parameter/field naming (no single-letter `v bool` or bare `stop bool`, `defined bool` -> `isStopOnFail bool`, `isStopped bool`, `isDefined bool`).
2. **Control Flow:** Maximum nesting depth 1, guard clauses, flattened conditionals.
3. **Naming & Types:** Enum `Type` suffixes, PascalCase types, no underscores in Go.
4. **Error Handling & Result Containers:** Single Result return containers (`appfault.ResultMap[K, V]`, `appfault.ResultSlice[T]`, `appfault.Result[T]`), `*appfault.AppError` returns, pointer-attached null safety (`*Result[T]` with line-1 `if r == nil` guards), method composition, 4 core predicates (`IsCountOtherThan`, `IsEmpty`, `HasRecord`, `IsDefined`), mandatory `types.go` single reusable type definitions for domain structs and Result aliases, no swallowed errors, typed error responses.
5. **Code Metrics:** Functions <= 15 lines, files <= 100 lines coding, blank line padding.

## Output

Generates structured audit logs and phased remediation plans in `.ai-memory/plans/pending/` with subtask micro-batches.

---

## Fast File Discovery & Reading Toolchain (GitMap AUM Primary, Python Fallback)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (e.g. `gitmap find "*.go" -ext "go"`, `gitmap find "01*"`)
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<pattern>" --lang go --limit 50`
- **Sub-Millisecond Folder Explorer & Reader:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> --ext .go --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Codebase Topology:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`
