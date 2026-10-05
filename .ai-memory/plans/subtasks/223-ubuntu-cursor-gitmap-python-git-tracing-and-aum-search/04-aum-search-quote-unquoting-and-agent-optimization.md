# Subtask 223.4: AUM Search Quote Unquoting & Agent Optimization

> **Subtask ID:** 223.4  
> **Target File:** `.ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/04-aum-search-quote-unquoting-and-agent-optimization.md`  
> **Parent Plan:** [.ai-memory/plans/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md](../../223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md)  
> **Spec Reference:** [02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md](../../../../02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md)  
> **Status:** Pending  
> **Target Area:** `cli/cmdautomation/search.go`, `cli/cmdautomation/automation_cmd.go`, `cli/cmdautomation/search_worker.go`  

---

## 1. Objective

Enhance `gitmap aum search` to automatically strip shell-escaped quotes and unwrapped delimiters from input search patterns. This ensures AI agents (Antigravity, Cursor, Claude, Copilot) passing quoted arguments (such as `"\"all\""` or `"'exact phrase'"`) find matches seamlessly, eliminating false-negative zero results and optimizing agent task speed.

---

## 2. Problem Statement

1. **Agent Quote Delimiter Pollution:** When autonomous agents invoke CLI tools via structured JSON tool-callers or shell runners, strings are frequently serialized with redundant quotation escapes:
   - Example agent call: `gitmap aum search "\"all\""`
   - Resulting Go argument: `args[0] = "\"all\""`
2. **Literal Matching Failure in Workers:** `search_worker.go` executes byte-level searches using `bytes.Contains(data, ctx.patBytes)`. Because the byte slice includes the literal quotation marks (`\"all\"`), files containing the word `all` without quotes fail to match, returning `0 hit(s)` across the entire repository.
3. **Agent Confusion & Hallucination Risk:** When searches return 0 hits, AI agents frequently assume required functions, constants, or files do not exist, triggering redundant code rewrites or exploratory stalls.

---

## 3. Implementation Details

### Step 3.1: Pattern Normalization Algorithm in `cli/cmdautomation/search.go`

Implement `cleanSearchPattern(pattern string) string`:
1. Trim leading and trailing whitespace.
2. In a loop, strip matching symmetric quotation pairs:
   - Outer escaped double quotes: `\"...\"`
   - Outer unescaped double quotes: `"..."`
   - Outer escaped single quotes: `\'...\'`
   - Outer unescaped single quotes: `'...'`
   - Outer backticks: `` `...` ``
3. Unescape internal quotes:
   - Replace `\"` with `"`
   - Replace `\'` with `'`
4. Boundary preservation: If the input consists solely of a single quotation mark (e.g. searching for the character `"`), preserve it to allow literal syntax searches.

### Step 3.2: Integration into CLI Handler (`cli/cmdautomation/automation_cmd.go`)

- In `runSearchCmd(cmd *cobra.Command, args []string) error`:
  - Sanitize `args[0]` immediately before assigning to `searchOpts.Pattern`:
    ```go
    searchOpts.Pattern = cleanSearchPattern(args[0])
    ```

### Step 3.3: Integration into Core Search Validator (`cli/cmdautomation/search.go`)

- In `validateSearchOptions(opts *SearchOptions) *apperror.AppError`:
  - Normalize `opts.Pattern = cleanSearchPattern(opts.Pattern)`.
  - Validate that the cleaned pattern is non-empty.

### Step 3.4: Benchmarking & Search Tool Cross-Verification

- Compare AUM search results and latency against external tools across polyglot repositories:
  - Benchmark exact string search (`gitmap aum search "func RunSearch"`).
  - Benchmark case-insensitive search (`gitmap aum search -i "runsearch"`).
  - Benchmark regex search (`gitmap aum search -r "type\s+\w+\s+struct"`).
- Document performance parity and memory allocation characteristics.

---

## 4. Verification & Testing Plan

1. **Unit Testing (`cli/cmdautomation/search_quote_test.go`):**
   - Author comprehensive table-driven tests verifying `cleanSearchPattern`:
     - Input: `"\"all\""` -> Output: `all`
     - Input: `"'RunSearch'"` -> Output: `RunSearch`
     - Input: `"\"\"nested\"\""` -> Output: `nested`
     - Input: ``"`backticks`"`` -> Output: `backticks`
     - Input: `"hello\"world"` -> Output: `hello"world`
     - Input: `"\""` -> Output: `"`
2. **End-to-End CLI Verification:**
   - Execute: `gitmap aum search "\"RunSearch\""`
   - Assert that matches in `cli/cmdautomation/search.go` are discovered identically to `gitmap aum search RunSearch`.
3. **Agent Search Benchmark:**
   - Verify execution time across all 8,400+ repository files completes in under 600ms.

---

## 5. Execution Checklist

- [ ] Implement `cleanSearchPattern` in `cli/cmdautomation/search.go`.
- [ ] Connect `cleanSearchPattern` to `validateSearchOptions` in `cli/cmdautomation/search.go`.
- [ ] Connect `cleanSearchPattern` to `runSearchCmd` in `cli/cmdautomation/automation_cmd.go`.
- [ ] Author table-driven unit tests in `cli/cmdautomation/search_quote_test.go`.
- [ ] Test CLI directly with quoted and unquoted inputs.
- [ ] Verify zero linter errors and zero regressions in `cli/cmdautomation/automation_test.go`.

---

## 6. Acceptance Criteria

- [ ] **AC-223.4-1:** `cleanSearchPattern` correctly strips outer double quotes, single quotes, escaped quotes, and backticks.
- [ ] **AC-223.4-2:** Running `gitmap aum search "\"RunSearch\""` returns positive matches in `cli/cmdautomation/search.go`.
- [ ] **AC-223.4-3:** Single-character syntax searches (such as `"` or `'`) retain their literal character without being stripped to an empty string.
- [ ] **AC-223.4-4:** Search performance remains sub-second across 8,000+ files.
- [ ] **AC-223.4-5:** Code strictly complies with zero nested ifs and early guard clause rules.
