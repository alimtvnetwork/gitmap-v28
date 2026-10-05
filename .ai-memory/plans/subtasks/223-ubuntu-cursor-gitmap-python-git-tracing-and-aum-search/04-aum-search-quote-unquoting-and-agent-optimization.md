# Subtask 223.4: AUM Search Quote Unquoting & Agent Optimization

> **Subtask ID:** 223.4  
> **Target File:** `.ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/04-aum-search-quote-unquoting-and-agent-optimization.md`  
> **Parent Plan:** [.ai-memory/plans/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md](../../223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md)  
> **Spec Reference:** [02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md](../../../../02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md)  
> **Status:** Completed  
> **Target Area:** `cli/cmdautomation/search.go`, `cli/cmdautomation/automation_cmd.go`, `cli/cmdautomation/search_worker.go`, `cli/cmdautomation/search_regex_helper.go`, `cli/cmdautomation/search_path_resolve.go`  

---

## 1. Objective

Enhance `gitmap aum search` to automatically strip shell-escaped quotes, unwrap delimiters, normalize escaped alternation pipes (`\|` -> `|`), auto-promote regex mode when regex syntax is present, and provide tolerant file path and extension resolution. This ensures AI agents (Antigravity, Cursor, Claude, Copilot) passing quoted arguments (such as `"\"all\""`), escaped pipes (such as `"timer\|Clock\|8:47"`), or omitted extensions (`FormRunner.t` -> `FormRunner.tsx`) find matches seamlessly, eliminating false-negative zero results and optimizing agent task speed.

---

## 2. Problem Statement

1. **Agent Quote Delimiter Pollution:** When autonomous agents invoke CLI tools via structured JSON tool-callers or shell runners, strings are frequently serialized with redundant quotation escapes:
   - Example agent call: `gitmap aum search "\"all\""`
   - Resulting Go argument: `args[0] = "\"all\""`
2. **Literal Matching Failure in Workers:** `search_worker.go` executes byte-level searches using `bytes.Contains(data, ctx.patBytes)`. Because the byte slice includes the literal quotation marks (`\"all\"`), files containing the word `all` without quotes fail to match, returning `0 hit(s)` across the entire repository.
3. **Escaped Pipe (`\|`) Regexp Mismatch:** In Go RE2, `\|` matches a literal pipe `|`, not alternation. Agents passing `"Candidate Response\|CANDIDATE RESPONSE"` or `"timer\|Clock\|8:47"` expect alternation.
4. **Omitted `-r` Flag:** Agents often pass regex patterns without `-r`, resulting in literal search failures.
5. **Path / Extension Typos:** Agents often pass paths missing `.tsx` or truncated extensions (`FormRunner.t`), resulting in 0 files walked.

---

## 3. Implementation Details

### Step 3.1: Pattern Normalization Algorithm in `cli/cmdautomation/search.go`

Implement `cleanSearchPattern(pattern string) string`:
1. Trim leading and trailing whitespace.
2. In a loop, strip matching symmetric quotation pairs (`\"...\"`, `"..."`, `\'...\'`, `'...'`, `` `...` ``).
3. Unescape internal quotes (`\"` -> `"`, `\'` -> `'`).
4. Boundary preservation: If the input consists solely of a single quotation mark, preserve it to allow literal syntax searches.

### Step 3.2: Regex Normalization & Auto-Promotion (`cli/cmdautomation/search_regex_helper.go`)

1. `normalizeRegexPattern(pattern)`: Converts unescaped `\|` to `|` while preserving `\\|`.
2. `hasRegexPatternSyntax(pattern)`: Detects `\|`, `.*`, `.+`, `|`, `[`, `]`, `\d`, `\w`, `\s`, anchors.
3. `autoPromoteRegex(opts)`: Auto-promotes `opts.IsRegex = true` and normalizes pattern if syntax is detected.

### Step 3.3: Tolerant Path Resolution (`cli/cmdautomation/search_path_resolve.go`)

1. Resolves exact files directly (`resolveSearchCandidateFiles`).
2. Checks extension variations (`.tsx`, `.ts`, `.jsx`, `.js`, `.go`, `.py`, etc.).
3. Checks prefix matches in parent folder (`FormRunner.t` -> `FormRunner.tsx`).

---

## 4. Execution Checklist

- [x] Implement `cleanSearchPattern` in `cli/cmdautomation/search.go`.
- [x] Connect `cleanSearchPattern` to `validateSearchOptions` in `cli/cmdautomation/search.go`.
- [x] Connect `cleanSearchPattern` to `runSearchCmd` in `cli/cmdautomation/automation_cmd.go`.
- [x] Author table-driven unit tests in `cli/cmdautomation/search_quote_test.go`.
- [x] Implement `normalizeRegexPattern`, `hasRegexPatternSyntax`, and `autoPromoteRegex` in `cli/cmdautomation/search_regex_helper.go`.
- [x] Implement `resolveSearchCandidateFiles` in `cli/cmdautomation/search_path_resolve.go`.
- [x] Author table-driven scenario tests in `cli/cmdautomation/search_regex_test.go`.
- [x] Verify zero linter errors and zero regressions.

---

## 5. Acceptance Criteria

- [x] **AC-223.4-1:** `cleanSearchPattern` correctly strips outer double quotes, single quotes, escaped quotes, and backticks.
- [x] **AC-223.4-2:** Running `gitmap aum search "\"RunSearch\""` returns positive matches in `cli/cmdautomation/search.go`.
- [x] **AC-223.4-3:** Single-character syntax searches (such as `"` or `'`) retain their literal character without being stripped to an empty string.
- [x] **AC-223.4-4:** Escaped alternation pipes (`\|`) are normalized to `|`, enabling both explicit and auto-promoted regex matching.
- [x] **AC-223.4-5:** Code strictly complies with zero nested ifs and early guard clause rules.
