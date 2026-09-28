# SPEC-APP-88: Shutdown Until Green Enhancements, PE Path/Alias Resolution, PascalCase/Array Template Variables & Repo-Secrets Hygiene

## 1. Context & Purpose
This specification formalizes enhancements across four core subsystems:
1. **Shutdown Until Green (`gitmap sug` / `shutdown-until`):** Subcommand normalization for multi-word/spaced tokens (e.g. `agy-running projects`), help-keyword protection on `add-projects help`, target existence verification with diagnostic resolution guidance, direct path invocations, watch loop process monitoring (`status` / `watch` / `w`), and dark-mode local browser UI (`ui` / `web`).
2. **Pipeline Errors (`gitmap pe`):** Target-aware resolution allowing path-based (`d:/work/antigravity-manager`), alias-based (`antigravity-manager`), or Git URL based queries without requiring the user to `cd` into the target directory.
3. **Help & Navigation CLI Documentation:** Dedicated help group filtering (`gitmap help <group>`) with tab completion, explicit documentation distinguishing `search` (content walk), `aum search` (macro index), and `file-search` / `find` (filename indexing), and Antigravity Manager (`agm`) integration guidance.
4. **Templates Pre-Compilation & Secrets Hygiene (`d:\work\repo-secrets`):** Deterministic PascalCase variable resolution (`${VarName}`), array variable random selection (`${Var}`) and indexed lookup (`${Var[i]}`), directory normalization from `01-git-map` to `01-gitmap`, boolean flag positive standardization (`is*`, `has*`), and hermetic test script verification directly under `D:\test-gitmap`.

---

## 2. Technical Architecture & Invariants

### 2.1 SUG Subcommand Normalization & Target Verification
- **Token Normalization:** When routing tokens in `RunSUGCLI(args)`:
  - Spaced tokens `agy-running projects`, `running projects`, `arp`, `rp` MUST normalize to `agy-running-projects`.
  - Aliases for `add-projects`: `add-projects`, `add`, `/add`, `AP`, `ap`.
  - Aliases for `rm`: `rm`, `/rm`, `remove`, `del`, `delete`.
  - Aliases for `run`: `run`, `watch`, `w`.
  - Aliases for `status`: `status`, `ls`, `list`.
  - Aliases for `ui`: `ui`, `watch ui`, `web`.
- **Help Guard:** If the first target passed to `add-projects` or `add` is `help`, `-h`, or `--help`, the system MUST render `RenderAgySugHelp()` and MUST NOT insert `"help"` into the watch list.
- **Target Validation:** Before adding a target:
  1. Check if path exists locally (`os.Stat(target)`).
  2. Check if repo exists in GitMap database (`store.ListRepos()`).
  3. Check if target is an active AGY project (`getAllProjects()`).
  4. Check if target is a valid Git URL (`http://`, `https://`, `git@`, `ssh://`, `.git`).
  - If validation fails, abort target insertion and display actionable guidance:
    `gitmap list`
    `gitmap agy ls`
    `gitmap agy running-projects`
    `gitmap scan`
- **Direct Path Invocation:** If `gitmap sug <path>` or `gitmap shutdown-until <path>` is executed and `<path>` is a directory or recognized target, automatically register and inspect that target.
- **State & Status Tracking:**
  - Persist active watch process metadata in `~/.gemini/antigravity/sug_state.json`:
    - `pid`: Process ID running the watch loop.
    - `startedAt`: UTC timestamp.
    - `intervalSeconds`: Configured polling duration.
    - `status`: `"running"` or `"idle"`.
    - `projectTargets`: List of watched targets.
  - `gitmap sug status` checks if PID is still alive. If alive, reports `"RUNNING (PID: X)"`; otherwise marks `"IDLE"`.

### 2.2 PE Target-Aware Resolution (`gitmap pe [target]`)
- When `gitmap pe` is invoked with a target argument (folder path, repo alias, or Git URL):
  1. If directory: query `git config --get remote.origin.url` inside that directory. If no remote, use folder base name.
  2. If Git URL: parse owner/repo slug using `parseSlugFromGitURL`.
  3. If repo name or alias: query `store.DefaultDBPath()` table `Repo` for matching `Slug` or `AbsolutePath`.
- Pass resolved repo slug to pipeline cache evaluator and GH API query runners, allowing cross-repository pipeline error inspections from any working directory.

### 2.3 Help Group Filtering & Command Differentiation
- `gitmap help <group>` filters and displays only commands belonging to that specific category (e.g. `scanning`, `gitops`, `integrations`, `release`, `navigation`, `templates`, `cluster`, `installers`, `search-find`).
- Tab completion in `GetRootCompletionCmd()` completes all help groups when typing `gitmap help <Tab>`.
- Clarify command search taxonomy:
  - `gitmap search <text>`: Multi-repository content / code grep walk.
  - `gitmap aum search <text>`: Automation Unit Manager macro and automation script index search.
  - `gitmap find` / `sf` / `file-search`: Fast file path and filename search across indexed repositories.
  - `gitmap agm` / `antigravity-manager`: GUI desktop app and multi-node fleet manager for Antigravity.

### 2.4 Pre-Compiled Variables & Secrets Hygiene
- **PascalCase & Fallback Resolution:** Support `${VarName}` alongside `${var_name}` and `$var_name`. If exact key is missing, case-insensitively match against defined keys.
- **Array Variable Support:**
  - If variable value is formatted as a JSON array (e.g. `["A", "B", "C"]`) or slice:
    - `${Var}` randomly selects one element from the array.
    - `${Var[i]}` selects the element at index `i` (0-indexed).
- **Repo-Secrets (`d:\work\repo-secrets`):**
  - Rename directory `01-git-map` -> `01-gitmap`.
  - Standardize `seo-templates.json` with `RISEUP ASIA LLC`, `https://alimkarim.com`, and name variations (`MD. Alim Ul Karim`, `MD Alim Ul Karim`, `Alim Karim`).
  - Standardize boolean fields in `commit-pull-config.json` with positive prefix (`isRecreate`, `isTree`, `isFinalSync`, `isPushImmediate`, `isApply`, `isApplyCD`).
  - Commit and push all updates in `d:\work\repo-secrets` to `origin/main`.
- **Test Repo Script:**
  - Verification script strictly targets `D:\test-gitmap` (at root `D:\`, never under `D:\work\`).

---

## 3. Verification & Acceptance Criteria
1. `gitmap sug agy-running projects` and `gitmap sug arp` successfully route to `agy-running-projects`.
2. `gitmap sug add-projects help` displays the help menu and does not add `"help"` to the watch list.
3. Adding a non-existent target fails with clear diagnostic commands (`gitmap list`, `gitmap agy ls`, etc.).
4. `gitmap sug status` correctly reports whether the watch loop is running or idle with PID and target count.
5. `gitmap sug ui` starts the dark-mode dashboard on a local HTTP port.
6. `gitmap pe <path|alias|url>` resolves and checks pipeline status for target repositories outside cwd.
7. `gitmap help <group>` filters commands by category, and `gitmap help <Tab>` provides group suggestions.
8. PascalCase variables and array indexed/random expansions function deterministically.
9. `d:\work\repo-secrets` has `01-gitmap`, updated variables, and is cleanly committed and pushed to git.
10. `D:\test-gitmap` script runs and cleans up strictly at `D:\test-gitmap`.
