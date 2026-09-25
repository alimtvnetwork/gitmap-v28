---
name: cg-file-size-and-function-reduction
description: Autonomously audits, decomposes, and refactors files over 100 lines and functions over 8–15 lines across the codebase using a two-part decomposition strategy, wrapper objects, and strict whitespace preservation without intermediate test or build runs.
---

# Skill: File Size & Function Size Reduction (`cg-file-size-and-function-reduction`)

This skill governs autonomous auditing, decomposition, and refactoring of bloated source code files exceeding the **100-line standard file cap** (recommended <= 80 lines) and functions exceeding **8–15 lines** across polyglot codebases (Go, TypeScript, Python, PHP, Rust, C#).

---

## 1. Mandatory Architectural Invariants

### 1. The Two-Part Decomposition Strategy
Never attempt to split a file without first decomposing its functions:
- **Part 1 (Function Level):** Identify functions > 15 lines and decompose them into focused, single-responsibility helper functions targeting **<= 8 lines** (hard cap: 15 lines).
- **Part 2 (File Level):** If the file remains > 100 lines after function decomposition, group cohesive clusters of helper functions and extract them into dedicated sibling files (e.g. `<name>_validator.go`, `<name>_helpers.go`, `<name>_converter.go`).

### 2. Strict Allowed Exceptions (Exempt from 100-Line Limit)
The following file categories are explicitly permitted to exceed 100 lines:
- **Data & Schema Files:** JSON files (`*.json`), YAML/TOML configuration manifests.
- **Type Definition Files:** Dedicated types files (`types.go`, `*.d.ts`, `types.ts`, `models.go`, `schema.ts`) that centralize domain structs or discriminated unions.
- **Lookup & Map Files:** Static lookup tables, dictionary mappings, registries, and transition matrices (`map.go`, `registry.go`, `lookup.ts`, `routes.go`).
- **Variables & Constants Files:** Centralized constant packages, enum registrations, and configuration variable dictionaries (`vars.go`, `consts.go`, `constants.ts`, `vars.ts`).

All logic, services, controllers, handlers, utilities, hooks, CLI commands, and scripts MUST remain strictly under 100 lines.

3. **Total Ban on Line-Compression Cheating:**
- **Zero Whitespace Stripping:** NEVER delete blank lines, squash vertical spacing, or compress `if/else` statements onto a single line to reduce line count.
- **No Inline Compound Condition Cramming:** NEVER cram variable declarations, type assertions, and compound conditions into the `if` header (e.g. `if v, isString := rawMap[key].(string); isString && len(v) > 0 {`). Put assignments on separate lines, evaluate booleans before `if`, and keep `if` conditions simple with one variable.
- **Zero Magic Strings & Constant Returns:** Never hardcode raw string literals (`"unknown"`, `"version"`) into function logic or return statements. Merge them into named constants/slices and return defined constants (`return VersionUnknown`).
- **Return New Line Concept (Mandatory):**
  - Exactly **ONE blank line BEFORE** every `return`, `throw`, or `break` statement.
  - Exactly **ONE blank line AFTER** every closing curly brace `}` of an `if`, `for`, `switch`, or helper block.
- Files must be reduced by modular decomposition into separate files, NOT by code compression.

### 4. Wrapper Objects & Parameter Reduction Pattern

- **Return Wrapper Struct:** If an extracted helper returns 2 or more related values (beyond standard error/Result), encapsulate them into a dedicated named struct or type.
- **Parameter Structs:** If an extracted function requires more than 2 parameters, combine them into an options or params struct (`*Params` / `*Options`).

### 5. Boolean Principles During Extraction

- **Implicit Positive Checks Only:** NEVER write `if isReady == true`. ALWAYS write `if isReady`.
- **No Negative Polarity:** NEVER combine a positive and negative condition (`if isA && !isB`). Split into separate guard clauses!
- **Positive Naming Only:** All booleans MUST use affirmative prefixes: `is*` or `has*` only (e.g., `isFail`, `isInvalid`, `isPaused`).
- **Zero Nested If Statements:** Every extracted function MUST maintain a nesting depth <= 1 using early guard returns.

### 6. Execution & Build Policy

- **NO INTERMEDIATE TEST RUNNING:** NEVER run unit test suites (`go test ./...`, `npm test`, `pytest`) during routine refactoring turns.
- **NO INTERMEDIATE BUILD CHECKING:** DO NOT execute build commands (`go build`, `npm run build`) after individual file edits.
- **FINAL STEP BUILD VERIFICATION ONLY:** Verify compilation strictly at the final step after all file extractions and import adjustments are completed.

---

## 2. Code Patterns: Before vs After

### Pattern 1: Monolithic Function Decomposition (Go)

```go
// ❌ ANTI-PATTERN: 32 lines, monolithic, mixed responsibilities
func HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    if req.Email == "" || req.Password == "" {
        http.Error(w, "missing required fields", http.StatusBadRequest)
        return
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        http.Error(w, "hashing failed", http.StatusInternalServerError)
        return
    }
    user := &User{Email: req.Email, PasswordHash: string(hash)}
    if err := db.CreateUser(user); err != nil {
        http.Error(w, "db error", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}
```

```go
// ✅ REQUIRED: Decomposed into small functions (<= 8 lines each) + early guards
func HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
    req, err := decodeRegisterRequest(r.Body)
    if err != nil {
        writeErrorResponse(w, err, http.StatusBadRequest)
        return
    }

    if isInvalid := req.Validate(); isInvalid {
        writeErrorResponse(w, errInvalidInput, http.StatusBadRequest)
        return
    }

    processRegistration(w, req)
}

func decodeRegisterRequest(body io.Reader) (RegisterRequest, error) {
    var req RegisterRequest
    err := json.NewDecoder(body).Decode(&req)

    return req, err
}

func processRegistration(w http.ResponseWriter, req RegisterRequest) {
    hash, err := hashPassword(req.Password)
    if err != nil {
        writeErrorResponse(w, err, http.StatusInternalServerError)
        return
    }

    saveAndRespondUser(w, req.Email, hash)
}
```

---

### Pattern 2: Multi-Value Return Wrapper Object (TypeScript)

```ts
// ❌ ANTI-PATTERN: returning tuple/multiple raw values
function parseFileHeader(content: string): [string, number, boolean] {
  const lines = content.split("\n");
  const title = lines[0].replace("# ", "");
  const size = content.length;
  const isDraft = lines[1]?.includes("draft") ?? false;

  return [title, size, isDraft];
}
```

```ts
// ✅ REQUIRED: encapsulate in clean named wrapper type
export interface FileHeaderMeta {
  title: string;
  contentLength: number;
  isDraft: boolean;
}

export function parseFileHeader(content: string): FileHeaderMeta {
  const lines = content.split("\n");
  const title = lines[0].replace("# ", "");
  const isDraft = Boolean(lines[1]?.includes("draft"));

  return {
    title,
    contentLength: content.length,
    isDraft,
  };
}
```

### Pattern 3: No Inline Compound Init Cramming (Multi-Line Separation & Single-Variable Guard)

```go
// ❌ BANNED ANTI-PATTERN:
// 1. Cramming type assertion assignment and compound condition into one line.
// 2. Hardcoding magic strings ("Version", "version", "unknown") inline.
// 3. Returning raw fallback literal instead of a defined constant.
func extractVersionValue(rawMap map[string]interface{}) string {
    for _, key := range []string{"Version", "version"} {
        if v, isString := rawMap[key].(string); isString && len(v) > 0 {
            return v
        }
    }

    return "unknown"
}

// ✅ MANDATORY CLEAN PATTERN:
// 1. Zero magic strings: extract lookup keys and defaults into constants.
// 2. Merge repeated/related strings into reusable collections (versionKeys).
// 3. Assignment on its own dedicated line.
// 4. Affirmative boolean (hasContent) pre-evaluated BEFORE the if statement.
// 5. Clean vertical breathing room (blank line before if).
// 6. Dead-simple if statement evaluating exactly ONE variable.
// 7. Return defined constant (VersionUnknown) instead of raw magic string literal.
const (
    VersionUnknown  = "unknown"
    versionKeyUpper = "Version"
    versionKeyLower = "version"
)

var versionKeys = []string{versionKeyUpper, versionKeyLower}

func extractVersionValue(rawMap map[string]interface{}) string {
    for _, key := range versionKeys {
        v, isString := rawMap[key].(string)
        hasContent := isString && len(v) > 0

        if hasContent {
            return v
        }
    }

    return VersionUnknown
}
```

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

---

## 3. Routine Execution Checklist

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) after individual file edits. Build verification is checked ONLY at the final step.
- [ ] **NO LINE COMPRESSION (TOTAL BAN):** NEVER remove blank lines, merge statements, or compress `if/else` blocks to reduce line count.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## 4. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!
