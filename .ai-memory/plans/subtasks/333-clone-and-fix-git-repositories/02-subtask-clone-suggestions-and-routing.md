# Subtask Plan 02: Clone Routing Disambiguation & Suggestion Engine

- **Task Code:** `Task-02`
- **Subtask Slug:** `02-subtask-clone-suggestions-and-routing`
- **Parent Plan:** `333-clone-and-fix-git-repositories` (`.ai-memory/plans/pending/333-clone-and-fix-git-repositories.md`)
- **Specification Reference:** `02-spec/21-app/333-clone-and-fix-git-repositories/02-clone-suggestion-and-disambiguation-spec.md`
- **Status:** PENDING
- **Target Repository:** `gitmap-v28`
- **Owned Files:**
  - `cli/cmd/roottooling_util.go`
  - `cli/cmd/rootcore.go`
  - `cli/constants/constants_clonefixrepo.go`
  - `cli/cmdclone/clone.go`
  - `cli/cmdclone/clone_manifest.go`
  - `cli/cmdclone/gh_resolver.go`
  - `cli/cmdclone/clonefixrepo.go`

---

## 1. Objectives & Scope

This subtask implements the technical changes specified in `02-spec/21-app/333-clone-and-fix-git-repositories/02-clone-suggestion-and-disambiguation-spec.md` to guarantee that GitMap can clone repositories by bare name or slug, disambiguate overloaded verbs, provide dedicated aliases, and present elegant Catppuccin suggestions when lookups fail:

1. **Disambiguate `cfr` and Add Dedicated `crf` / `crfb` Aliases:**
   - In `cli/cmd/roottooling_util.go`, fix `runCfrDisambiguated` and `isCfrCloneRequest`:
     - If `!cmdpull.IsGitRepoCWD()`, the user is not in a git repository; route directly to `cmdclone.RunCloneFixRepo`.
     - If the first argument is an alphanumeric slug or repo name (e.g. `yt-manage-v1`), route to `cmdclone.RunCloneFixRepo`.
   - In `cli/cmd/rootcore.go` and `cli/constants/constants_clonefixrepo.go`, register `crf` as a dedicated alias for `clone-fix-repo`, and `crfb` as a dedicated alias for `clone-fix-repo-pub`.

2. **Differentiate File Manifests from Slugs in `dispatchCloneExecution`:**
   - In `cli/cmdclone/clone.go` and `cli/cmdclone/clone_manifest.go`, refactor `dispatchCloneExecution`:
     - If `isDirectURL(cf.Source)`: execute direct clone via `executeDirectClone`.
     - If `isManifestSource(cf.Source)` (`fsutil.FileExists(cf.Source)` or shorthand `json`/`csv`/`text`): execute manifest clone flow via `executeManifestCloneFlow`.
     - If neither: resolve slug via `ResolveRepoSlugForProfile(cf.Source, cf.Profile)` or `ResolveRepoSlug(cf.Source)`.
     - If resolved to a Git URL: execute direct clone via `executeDirectClone(buildDirectCloneParams(cf))`.
     - If resolution fails: query local SQLite cache (`store.DB`) and GitHub CLI (`gh repo list`), render Catppuccin suggestion card via `diag.RenderBox`, and return clean `*appfault.AppError` without attempting `os.Open` on the slug!

3. **Enforce Leaf Error Sovereignty & Suggestions in `clonefixrepo.go`:**
   - In `cli/cmdclone/clonefixrepo.go`, refactor `validateAndPrepareCFR` to return an error instead of invoking `cliexit.HandleError(nil, constants.ExitCloneFixRepoBadFlag)` with exit code 6.
   - When URL is missing or unresolvable, query SQLite/GitHub, display Catppuccin suggestions, and return `*appfault.AppError`.

4. **Preserve Coding Guidelines & Invariants:**
   - Max 15 lines per function, <= 8 lines preferred.
   - Zero nested `if` statements (guard clauses with early returns).
   - Positive booleans only (`is*`, `has*`).
   - Strict relative paths across all files and evidence.

---

## 2. Target Files & Component Map

| Target File Path | Purpose | Modification Scope |
| :--- | :--- | :--- |
| `cli/cmd/roottooling_util.go` | Top-level tooling router | Refactor `runCfrDisambiguated` and `isCfrCloneRequest`; add slug and git cwd guards |
| `cli/cmd/rootcore.go` | Core command dispatch table | Register `crf` and `crfb` aliases in `coreCloneExtEntries()` |
| `cli/constants/constants_clonefixrepo.go` | Clone-fix-repo constant tokens | Add `CmdCloneFixRepoShort` (`"crf"`) and `CmdCloneFixRepoPubShort` (`"crfb"`) |
| `cli/cmdclone/clone.go` | Clone CLI orchestrator & dispatcher | Refactor `dispatchCloneExecution` to branch between URL, manifest, and slug |
| `cli/cmdclone/clone_manifest.go` | Manifest-based cloner | Isolate manifest parsing so non-existent files are never passed to `cloner.LoadRecords` |
| `cli/cmdclone/gh_resolver.go` | Slug and GitHub repo resolver | Add `buildRepoSuggestionGroup` and integrate `diag.RenderBox` Catppuccin cards |
| `cli/cmdclone/clonefixrepo.go` | Clone-fix-repo pipeline | Refactor `validateAndPrepareCFR` to return `*appfault.AppError` and render suggestions |

---

## 3. Atomic Execution Checklist

### Phase 1: Disambiguation & Dedicated Aliases (`crf`, `crfb`, `cfr`)
- [x] **1.1**: In `cli/constants/constants_clonefixrepo.go`, define:
  ```go
  CmdCloneFixRepoShort    = "crf"
  CmdCloneFixRepoPubShort = "crfb"
  ```
- [x] **1.2**: In `cli/cmd/rootcore.go`, update `coreCloneExtEntries()`:
  - Add `constants.CmdCloneFixRepoShort` to the `constants.CmdCloneFixRepo` entry.
  - Add `constants.CmdCloneFixRepoPubShort` to the `constants.CmdCloneFixRepoPub` entry.
- [x] **1.3**: In `cli/cmd/roottooling_util.go`:
  - Update `runCfrDisambiguated(args []string) error`:
    - Add guard: `if !cmdpull.IsGitRepoCWD() { return cmdclone.RunCloneFixRepo(args) }`.
    - Check `isCfrCloneRequest(args)`.
  - Update `isCfrCloneRequest(args []string) bool`:
    - Check `isLikelyURLOrGit(first)` and `isLikelyRepoSlug(first)`.
  - Implement `isLikelyRepoSlug(s string) bool` helper (ensuring no spaces, valid slug characters, no leading flags).
  - Ensure all new/modified functions in `roottooling_util.go` are $\le 15$ lines with zero nested `if` statements.

### Phase 2: Manifest vs Slug Clone Routing (`dispatchCloneExecution`)
- [ ] **2.1**: In `cli/cmdclone/clone.go`:
  - Add `isManifestSource(source string) bool` helper checking `fsutil.FileExists(source)` and shorthand tokens (`json`, `csv`, `text`).
  - Refactor `dispatchCloneExecution(cf CloneFlags) error`:
    - Guard 1: `isMultiCloneEnabled(cf)` $\rightarrow$ `runCloneMulti(cf)`.
    - Guard 2: `isDirectURL(cf.Source)` $\rightarrow$ `executeDirectCloneFlow(cf)`.
    - Guard 3: `isManifestSource(cf.Source)` $\rightarrow$ `executeManifestCloneFlow(cf)`.
    - Guard 4: Default $\rightarrow$ `executeSlugCloneFlow(cf)`.
- [ ] **2.2**: Implement `executeDirectCloneFlow(cf CloneFlags) error`:
  - Invoke `executeDirectClone(buildDirectCloneParams(cf))`.
  - Invoke `maybeExitOnCmdFaithfulMismatch()`.
  - Return error directly.
- [ ] **2.3**: Implement `executeSlugCloneFlow(cf CloneFlags) error`:
  - Attempt resolution via `ResolveRepoSlugForProfile(cf.Source, cf.Profile)`.
  - If resolved to direct URL, update `cf.Source = resolved` and return `executeDirectCloneFlow(cf)`.
  - If resolution fails, build suggestion group via `buildRepoSuggestionGroup(cf.Source)`.
  - Render Catppuccin suggestion box if suggestions exist.
  - Return `appfault.NewNotFoundError(...)` with attached suggestions.

### Phase 3: Suggestion Engine & Catppuccin Rendering
- [ ] **3.1**: In `cli/cmdclone/gh_resolver.go`:
  - Implement `buildRepoSuggestionGroup(target string) diag.SuggestionGroup`.
  - Implement `collectRepoSuggestions(target string) []string` querying `store.OpenDefault()` and `queryGitHubSuggestions(target)`.
  - Map candidate suggestions into `diag.Suggestion` structs with `Command: fmt.Sprintf("gitmap clone %s", name)`.
- [ ] **3.2**: Connect Catppuccin rendering via `diag.RenderBox(output.UIErr(), group)`.
- [ ] **3.3**: Ensure fallback returns structured `*appfault.AppError` so parent command dispatchers capture audit events.

### Phase 4: `clonefixrepo` Leaf Error Sovereignty
- [ ] **4.1**: In `cli/cmdclone/clonefixrepo.go`:
  - Update `validateAndPrepareCFR(f *cloneFixRepoFlags)` signature to return `(string, string, error)`.
  - Remove `cliexit.HandleError(nil, constants.ExitCloneFixRepoBadFlag)`.
  - If `len(f.url) == 0`: render suggestions from recent repos and return validation `*appfault.AppError`.
  - If `ResolveRepoSlug(f.url)` fails to produce a valid URL: render suggestions and return not-found `*appfault.AppError`.
- [ ] **4.2**: Update callers of `validateAndPrepareCFR` in `clonefixrepo.go` to bubble returned errors upward.

---

## 4. Coding Guidelines & Invariants Verification

- [ ] **Leaf Error Sovereignty:** No leaf functions call `cliexit.HandleError` or `os.Exit`.
- [ ] **Function Length:** Every new or edited function is $\le 15$ lines.
- [ ] **Conditionals:** Zero nested `if` statements. Early returns and guard clauses only.
- [ ] **Boolean Naming:** Positive booleans only (`is*` or `has*`).
- [ ] **Relative Paths:** Only relative paths used throughout all documentation, code, and plans.

---

## 5. Verification Evidence Contract & Quality Gates

The implementation of `Task-02` will be considered complete when the following verification steps succeed:

| Gate | Check / Command | Expected Outcome |
| :--- | :--- | :--- |
| **G1: Slug Clone Resolution** | `gitmap clone yt-manage-v1 --dry-run` | Successfully resolves repo to GitHub URL and outputs clone plan (no manifest error) |
| **G2: Unknown Slug Suggestion** | `gitmap clone non-existent-test-repo` | Displays Catppuccin `Repository Suggestions` box with closest matches and exits cleanly with not-found error |
| **G3: `cfr` Slug Routing** | Run `gitmap cfr yt-manage-v1` in non-git directory | Routes to `RunCloneFixRepo` without `Not a git repository` error |
| **G4: Dedicated `crf` Alias** | `gitmap crf --help` | Displays `clone-fix-repo` help menu |
| **G5: Dedicated `crfb` Alias** | `gitmap crfb --help` | Displays `clone-fix-repo-pub` help menu |
| **G6: Relative Path Linter** | `python linter-scripts/check-relative-paths.py` | 100% Pass, zero absolute path violations |
