# Subtask 03: Merge-AI CLI Help, Flag Parsing & Parameterized Examples Specification

> **Subtask ID:** Task-03  
> **Parent Plan:** `.ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md`  
> **Spec Reference:** `02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/02-component-and-cli-spec.md`  
> **Assigned Role:** Merge-AI & CLI Parameterization Engineer  
> **Status:** PENDING  
> **Owned Files:**  
> - `.ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/03-merge-ai-parameterized-examples-spec.md`  

---

## 1. Deliverables & Implementation Objectives

This subtask implements the comprehensive CLI interface enhancements for `gitmap merge-ai` (alias `ma`), resolving the acute deficiency where developers and AI agents lacked parameterized examples, named parameter flag parsing, and help flag interception.

### Key Deliverables:
1. **Help Flag Interception & ANSI Renderer (`cli/cmdmergeai/merge_ai_help.go`):**
   - Intercept `-h`, `--help`, `help`, or zero arguments.
   - Print a beautifully formatted PTerm / ANSI terminal manual detailing command syntax, flag descriptions, collision architecture (`01_` base, `02_` incoming), and concrete examples across all 5 input modalities.
   - Return exit code 0 when help is requested.
2. **Multi-Modal Argument Parser (`cli/cmdmergeai/merge_ai_core.go`):**
   - Refactor `parseMergeAIArgs` to support:
     - Named flags: `--destination <dir>`, `-d <dir>`, `--source <url>`, `-s <url>`, `--url <url>`, `--config <file.json>`, `-c <file.json>`.
     - Comma-delimited strings: `--source "repo1,repo2,repo3"`.
     - Text manifest file: `sources.txt` containing newline-delimited URLs.
     - JSON configuration file: `config.json` defining `destination` and `sources` array.
     - Positional backwards-compatibility: `gitmap merge-ai <destination> <source1> <source2>...`.
3. **Embedded Help Documentation (`cli/helpdoc/merge-ai.md`):**
   - Author rich Markdown documentation for GitMap's embedded help catalog (`gitmap help merge-ai`).
4. **Comprehensive Go Unit Tests (`cli/cmdmergeai/merge_ai_test.go`):**
   - Test flag parsing for all 5 modalities.
   - Test help interception with `-h` and `--help`.
   - Test negative validation cases (missing destination, missing sources, invalid config path).

---

## 2. Technical Execution Steps

### Step 1: Author Help Display Module (`cli/cmdmergeai/merge_ai_help.go`)
- Create `RenderMergeAIHelp()` using `pterm` styling.
- Display:
  - Command usage banner (`gitmap merge-ai [destination] [sources...] [flags]`).
  - Flag reference table (`-d`, `-s`, `-u`, `-c`, `-h`).
  - The 5 Parameterized Input Modalities with real-world examples.
  - Overview of `01_`/`02_` collision sequencing and `instruction.md` generation.

### Step 2: Refactor `parseMergeAIArgs` in `cli/cmdmergeai/merge_ai_core.go`
- Introduce state-machine / slice flag scanner:
  - Scan for help tokens (`-h`, `--help`, `help`). If detected, invoke `RenderMergeAIHelp()` and return a sentinel `ErrHelpRequested` or clean exit.
  - Parse `--destination` / `-d` into `destTarget`.
  - Collect repeated `--source` / `-s` / `--url` / `-u` arguments into `sources` slice.
  - Parse `--config` / `-c` or a lone `.json` argument into `MergeAIConfigFile`.
  - Expand any `.txt` manifest into individual repository sources.
  - Split any comma-delimited strings into discrete source items.
  - Fall back to standard positional parsing if no named flags are detected.

### Step 3: Register in Help Catalog (`cli/helpdoc/`)
- Create `cli/helpdoc/merge-ai.md` with complete reference material.
- Register `"merge-ai"` and `"ma"` in `cli/helpdoc/catalog.go`.

### Step 4: Add Unit Tests (`cli/cmdmergeai/merge_ai_test.go`)
- Test cases:
  - `TestParseMergeAIArgs_Positional`
  - `TestParseMergeAIArgs_NamedFlags`
  - `TestParseMergeAIArgs_CommaSeparated`
  - `TestParseMergeAIArgs_TextManifest`
  - `TestParseMergeAIArgs_JSONConfig`
  - `TestParseMergeAIArgs_HelpInterception`
  - `TestParseMergeAIArgs_MissingDestination`

---

## 3. Concrete Verification Commands & Test Matrix

| Step | Target Subsystem | Verification Command | Expected Outcome | Failure Remediation |
|---|---|---|---|---|
| **V-01** | **Help Flag `--help`** | `go run ./cli/cmd/gitmap merge-ai --help` | Exit code 0, displays formatted PTerm help manual with all 5 modalities | Check help interception before argument length check in `RunMergeAI` |
| **V-02** | **Help Flag Shorthand `-h`** | `go run ./cli/cmd/gitmap ma -h` | Exit code 0, displays formatted PTerm help manual | Check alias routing and short flag handling |
| **V-03** | **Zero Arguments Help** | `go run ./cli/cmd/gitmap merge-ai` | Exit code 0 (or friendly usage error), displays usage and examples | Route zero args to help display |
| **V-04** | **Named Flags Syntax** | `go run ./cli/cmd/gitmap merge-ai -d ./test-dest -s https://github.com/org/a.git -s https://github.com/org/b.git --dry-run` | Successfully parses destination `./test-dest` and 2 sources | Check named flag accumulation loop |
| **V-05** | **Comma-Separated Sources** | `go run ./cli/cmd/gitmap ma -d ./test-dest -s "repo-a,repo-b,repo-c" --dry-run` | Successfully splits into 3 discrete sources | Check comma splitting logic |
| **V-06** | **JSON Config Parsing** | `go run ./cli/cmd/gitmap merge-ai --config test-config.json --dry-run` | Successfully loads destination and sources from JSON | Verify JSON unmarshaling into `MergeAIConfigFile` |
| **V-07** | **Go Unit Tests** | `go test -v ./cli/cmdmergeai/...` | 100% passing tests, 0 failures | Fix failing assertions in unit tests |
| **V-08** | **Embedded Help Catalog** | `go run ./cli/cmd/gitmap help merge-ai` | Prints markdown documentation from `cli/helpdoc/merge-ai.md` | Check catalog registration in `catalog.go` |

---

## 4. Acceptance Criteria Traceability Matrix

| Spec Criterion | Description | Target Step | Status |
|---|---|---|---|
| **AC-06** | `merge-ai` Help Interception (`-h`, `--help`, `help`) | V-01, V-02, V-03 | Pending Implementation |
| **AC-07** | Named Parameter Flags (`-d`, `--destination`, `-s`, `--source`, `--url`, `--config`) | V-04, V-06 | Pending Implementation |
| **AC-08** | Support for All 5 Input Modalities (Positional, Named, Comma, Text, JSON) | V-04, V-05, V-06, V-07 | Pending Implementation |
| **AC-09** | Collision Staging Architecture (`01_`/`02_`, manifest, instruction.md) | V-07 (Unit tests) | Pending Implementation |

---

## 5. Evidence Checks & Quality Gates

1. **Exit Code Fidelity:** Running `gitmap merge-ai --help` MUST exit with code 0.
2. **Zero Code Deletion in Staging:** Collision tests must verify that `01_<file>` contains the exact bytes of the base repo and `02_<file>` contains the exact bytes of the incoming repo.
3. **Strict Relative Git Paths:** All generated documentation and tests must reference paths relative to repository root (`cli/cmdmergeai/...`).
4. **Clean Code Style:** All Go code must follow coding guidelines (blank lines around structs, positive boolean conventions, structured error wrapping).
