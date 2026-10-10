# Engineering Subtask Plan: Tree Export Flags, Wildcard Glob Interception & Format Deduction

**Subtask ID:** `01-tree-flags-and-glob-fix`  
**Subtask Code:** `Task-01-TreeFlags`  
**Parent Plan:** `gitmap-command-implementation` ([gitmap-command-implementation.md](../../gitmap-command-implementation.md))  
**Run Number:** 83  
**Assigned Agent Role:** Implementation Worker 01 (Downstream Execution)  
**Spec Reference:**  
- [01-architecture-spec.md](../../../../02-spec/21-app/gitmap-command-implementation/01-architecture-spec.md)  
- [02-component-spec.md](../../../../02-spec/21-app/gitmap-command-implementation/02-component-spec.md)  
**Status:** `READY FOR IMPLEMENTATION`  

---

## 1. Targeted File Inventory & Disjoint Boundaries

To guarantee multi-agent isolation and eliminate concurrency conflicts:
- **Owned Subtask & Spec Files (Author 01):**
  * `02-spec/21-app/gitmap-command-implementation/01-architecture-spec.md`
  * `.ai-memory/plans/subtasks/gitmap-command-implementation/01-tree-flags-and-glob-fix.md`
- **Downstream Owned Implementation Files (Worker 01):**
  * `cli/cmd/folder/filter.go`
  * `cli/cmd/folder/parse_args.go`
  * `cli/cmd/folder/folder.go`
  * `cli/cmd/folder/folder_test.go`
- **Strictly Prohibited Actions:**
  * **TOTAL BAN ON GIT COMMANDS:** Subagents MUST NOT execute `git add`, `git commit`, `git status`, `git push`, `git diff`, or `git checkout`.
  * **Strictly Relative Paths Only:** All file references in documentation, code comments, and test fixtures must use forward-slash paths relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero absolute paths or `file:///` URIs.
  * **Search Exclusively via GitMap:** Use `gitmap aum search`, `gitmap find`, and `gitmap cat`. Strict ban on `grep`, `ripgrep`, `rg`, `Select-String`, `findstr`.
  * **No Build / No Test in Spec Phase:** Spec authoring requires zero build or test execution.

---

## 2. Objective & Scope

Worker 01 will implement the following core capabilities and bug fixes in the `cli/cmd/folder` package:
1. **Positional Wildcard Glob Interception:** Intercept wildcard patterns (`*`, `?`) passed as non-flag arguments (e.g. `gitmap tree "*.go"`), preventing Windows path syntax errors, defaulting `opts.TargetDir` to `.`, and routing patterns to `FilterConfig.IncludeGlobs` and `FilterConfig.Extensions`.
2. **Export Flags (`--file` / `-f`):** Support `--file <path>` and `-f <path>` as first-class CLI flags alongside legacy `-o` and `--out`.
3. **Automatic Format Deduction:** When outputting to a file without an explicit format flag, deduce format from the file extension (`.json` -> `FormatJson`, `.yaml`/`.yml` -> `FormatYaml`, `.txt`/`.tree` -> `FormatTree`, `.md` -> `FormatMd`).
4. **Tree Visual Layout Preservation:** Ensure `.txt` and `.tree` output retains the full hierarchical ASCII/Unicode branch tree (`├── `, `└── `) rather than flattening.
5. **Direct File Persistence:** Ensure `writeOrPrintOutput` safely creates missing parent directories before writing.
6. **Comprehensive Unit Testing:** Add rigorous unit tests in `cli/cmd/folder/folder_test.go` verifying wildcard arguments, export flags, format deduction, and filesystem operations.

---

## 3. Operational Step-by-Step Implementation Instructions

### Step 1: Update `cli/cmd/folder/filter.go`
* **Objective:** Extend `FilterConfig` with `IncludeGlobs` and implement inclusion matching for file metadata while guaranteeing directory descent remains unblocked.
* **Exact Modifications:**
  1. In `FilterConfig` struct, add the `IncludeGlobs` field:
     ```go
     type FilterConfig struct {
         ExceptGlobs  []string
         IncludeGlobs []string
         Extensions   []string
         MaxDepth     int
         OnlyText     bool
         OnlyBinary   bool
     }
     ```
  2. Implement helper `matchAnyIncludeGlob(normPath, base string, globs []string) bool`:
     ```go
     func matchAnyIncludeGlob(normPath, base string, globs []string) bool {
         for _, g := range globs {
             if matchGlob(g, normPath, base, false) {
                 return true
             }
         }
         return false
     }
     ```
  3. Update `IsMetaAllowed(meta *FileMeta) bool`:
     ```go
     func (fc *FilterConfig) IsMetaAllowed(meta *FileMeta) bool {
         if fc.OnlyText && meta.IsBinary {
             return false
         }

         if fc.OnlyBinary && !meta.IsBinary {
             return false
         }

         if len(fc.Extensions) > 0 && !hasMatchingExtension(meta.Extension, fc.Extensions) {
             return false
         }

         if len(fc.IncludeGlobs) > 0 && !matchAnyIncludeGlob(meta.Path, meta.Filename, fc.IncludeGlobs) {
             return false
         }

         return true
     }
     ```
  4. **Directory Invariant:** In `handleDirSkip(rel string, filter FilterConfig)`, ensure directory skipping continues to evaluate **only** `ExceptGlobs` and `MaxDepth`. Do **NOT** filter directories with `IncludeGlobs`, allowing recursive descent into all child folders.

---

### Step 2: Update `cli/cmd/folder/parse_args.go`
* **Objective:** Support `--file` / `-f`, detect wildcard patterns in positional arguments, route globs to `FilterConfig`, and correct format deduction for `.txt` and `.tree`.
* **Exact Modifications:**
  1. Add wildcard pattern detector:
     ```go
     func hasWildcardPattern(arg string) bool {
         return strings.ContainsAny(arg, "*?")
     }
     ```
  2. In `ParseArgs(args []string)`, track explicit format flags (`hasExplicitFormat bool`) and handle `--file` / `-f`:
     ```go
     var hasExplicitFormat bool

     for i := 0; i < len(args); i++ {
         arg := args[i]
         switch {
         case arg == "--tree" || arg == "-tree":
             opts.Format = FormatTree
             hasExplicitFormat = true
         case arg == "--md" || arg == "-md" || arg == "--markdown":
             opts.Format = FormatMd
             hasExplicitFormat = true
         case arg == "--json" || arg == "-json":
             opts.Format = FormatJson
             hasExplicitFormat = true
         case arg == "--yaml" || arg == "-yaml" || arg == "--yml":
             opts.Format = FormatYaml
             hasExplicitFormat = true
         case arg == "--flat" || arg == "-flat" || arg == "--list":
             opts.Format = FormatFlat
             hasExplicitFormat = true
         case (arg == "--file" || arg == "-f" || arg == "-o" || arg == "--out") && i+1 < len(args):
             opts.OutFile = args[i+1]
             i++
         ...
     ```
  3. In `resolveNonFlags(opts Options, nonFlags []string, hasExplicitFormat bool)`:
     - Iterate through `nonFlags`:
       * If `hasWildcardPattern(arg)`:
         - Append `arg` to `opts.Filter.IncludeGlobs`.
         - If `strings.HasPrefix(arg, "*.")`:
           - Extract extension: `ext := strings.TrimPrefix(arg, "*.")`.
           - If `ext != ""` and `!strings.ContainsAny(ext, "*?")`:
             - Append `ext` to `opts.Filter.Extensions`.
         - Continue without assigning `arg` to `opts.TargetDir` or `opts.OutFile`.
       * If `!hasWildcardPattern(arg)`:
         - If `opts.TargetDir == "."` && target not explicitly assigned yet:
           - Set `opts.TargetDir = arg`.
         - Else if `opts.OutFile == ""`:
           - Set `opts.OutFile = arg`.
           - If `!hasExplicitFormat`:
             - `opts.Format = deduceFormatFromExt(opts.OutFile, opts.Format)`.
     - If `opts.OutFile != ""` && `!hasExplicitFormat`:
       - `opts.Format = deduceFormatFromExt(opts.OutFile, opts.Format)`.
     - Guarantee `opts.TargetDir` remains `.` if no directory was assigned.
  4. In `deduceFormatFromExt(outFile string, fallback OutputFormat) OutputFormat`:
     - Update `.txt` and `.tree` to return `FormatTree` (preserving visual tree layout):
     ```go
     func deduceFormatFromExt(outFile string, fallback OutputFormat) OutputFormat {
         ext := strings.ToLower(filepath.Ext(outFile))
         switch ext {
         case ".json":
             return FormatJson
         case ".yaml", ".yml":
             return FormatYaml
         case ".md":
             return FormatMd
         case ".txt", ".tree":
             return FormatTree
         default:
             return fallback
         }
     }
     ```

---

### Step 3: Update `cli/cmd/folder/folder.go`
* **Objective:** Ensure atomic directory creation before file writes and verify output dispatch.
* **Exact Modifications:**
  1. In `writeOrPrintOutput(content, outFile string) error`:
     ```go
     func writeOrPrintOutput(content, outFile string) error {
         if outFile != "" {
             dir := filepath.Dir(outFile)
             if dir != "." && dir != "" {
                 if err := os.MkdirAll(dir, 0755); err != nil {
                     return err
                 }
             }
             return os.WriteFile(outFile, []byte(content), 0644)
         }

         fmt.Print(content)
         return nil
     }
     ```
  2. Verify that `Run(args []string)` and `ScanDirectory` operate cleanly with the updated `FilterConfig`.

---

### Step 4: Update `cli/cmd/folder/folder_test.go`
* **Objective:** Author comprehensive unit tests covering all new features and edge cases.
* **Exact Test Cases to Implement:**
  1. `TestParseArgs_FileFlags`:
     - Test `-f "a.txt"`, `--file "a.json"`, `-f "a.yaml"`, `-f "a.yml"`.
     - Verify `opts.OutFile` is populated correctly.
     - Verify `opts.Format` is deduced correctly (`FormatTree`, `FormatJson`, `FormatYaml`).
  2. `TestParseArgs_WildcardPositionalArgs`:
     - Test input `[]string{"*.go"}`:
       * `opts.TargetDir` == `"."`.
       * `opts.Filter.IncludeGlobs` contains `"*.go"`.
       * `opts.Filter.Extensions` contains `"go"`.
       * `opts.OutFile` == `""`.
     - Test input `[]string{"src", "*.ts"}`:
       * `opts.TargetDir` == `"src"`.
       * `opts.Filter.IncludeGlobs` contains `"*.ts"`.
       * `opts.Filter.Extensions` contains `"ts"`.
     - Test input `[]string{"*.go", "-f", "out.txt"}`:
       * `opts.TargetDir` == `"."`.
       * `opts.OutFile` == `"out.txt"`.
       * `opts.Format` == `FormatTree`.
  3. `TestParseArgs_ExplicitFormatPrecedence`:
     - Test input `[]string{"--flat", "-f", "out.txt"}`:
       * `opts.Format` == `FormatFlat` (override honored).
     - Test input `[]string{"--json", "-f", "out.txt"}`:
       * `opts.Format` == `FormatJson` (override honored).
  4. `TestDeduceFormatFromExt`:
     - Validate `.json` -> `FormatJson`.
     - Validate `.yaml` and `.yml` -> `FormatYaml`.
     - Validate `.txt` and `.tree` -> `FormatTree`.
     - Validate `.md` -> `FormatMd`.
     - Validate unknown extension `.unknown` -> fallback.
  5. `TestScanDirectory_IncludeGlobs`:
     - Using `createTestWorkspace(t)`:
       * Set `FilterConfig{IncludeGlobs: []string{"*.go"}}`.
       * Assert only `.go` files are returned.
       * Assert directory tree traversal descended into `src/core/util.go`.
       * Assert `.ts`, `.md`, `.js`, `.png` are excluded.
  6. `TestWriteOrPrintOutput_NestedDirectories`:
     - Create temp folder, invoke `writeOrPrintOutput("test content", filepath.Join(tempDir, "sub", "dir", "tree.txt"))`.
     - Verify file exists and content matches.

---

## 4. Verification Matrix & Quality Gates

| Gate ID | Target Feature | Validation Command / Test Case | Success Assertion |
| :--- | :--- | :--- | :--- |
| **VG-01** | Wildcard Positional Glob | `TestParseArgs_WildcardPositionalArgs` | Positional `"*.go"` sets `TargetDir="."`, populates `IncludeGlobs` & `Extensions`; no Win32 syntax crash. |
| **VG-02** | Text Tree Export (`-f "a.txt"`) | `TestParseArgs_FileFlags` & `TestDeduceFormatFromExt` | `opts.OutFile=="a.txt"`, `Format==FormatTree`; output preserves `├──` branch formatting. |
| **VG-03** | JSON Export (`-f "a.json"`) | `TestParseArgs_FileFlags` & `TestDeduceFormatFromExt` | `opts.OutFile=="a.json"`, `Format==FormatJson`; valid JSON output. |
| **VG-04** | YAML Export (`-f "a.yaml"`) | `TestParseArgs_FileFlags` & `TestDeduceFormatFromExt` | `opts.OutFile=="a.yaml"`, `Format==FormatYaml`; valid YAML output. |

---

## 5. Definition of Done for Worker 01

- [ ] `FilterConfig` in `filter.go` includes `IncludeGlobs []string`.
- [ ] `IsMetaAllowed` checks `IncludeGlobs` for files without impeding directory traversal in `handleDirSkip`.
- [ ] `hasWildcardPattern` detects `*` and `?` in `parse_args.go`.
- [ ] Non-flag arguments containing wildcards default `opts.TargetDir` to `.` and populate `IncludeGlobs` and `Extensions`.
- [ ] `--file` and `-f` flags are parsed in `parse_args.go` alongside `-o` and `--out`.
- [ ] Format deduction maps `.json` -> `FormatJson`, `.yaml`/`.yml` -> `FormatYaml`, `.txt`/`.tree` -> `FormatTree`, `.md` -> `FormatMd`.
- [ ] Explicit format flags (`--tree`, `--flat`, `--json`, etc.) override extension deduction.
- [ ] `writeOrPrintOutput` in `folder.go` creates parent directories before writing.
- [ ] Unit tests in `folder_test.go` cover all gates (VG-01 to VG-04) and pass with 0 failures.
- [ ] Positive boolean naming conventions strictly followed (`isAllowed`, `hasWildcard`, `isMatch`).
- [ ] All file references use forward-slash relative Git paths (`02-spec/...`, `.ai-memory/...`, `cli/...`).
- [ ] Zero git commands executed.
