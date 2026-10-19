# Specification: GitMap Efficiency and Coding Guidelines Enhancement

## Status
- **Task ID**: `gitmap-efficiency-and-coding-guidelines-enhancement`
- **Specification Version**: 1.0.0
- **Status**: APPROVED
- **Date**: 2026-10-19

---

## 1. User Request (Verbatim)

```text
# GitMap Efficiency and Coding Guidelines Enhancement: high priority instruction, non-negotiable task

Git pull done by gitmap needs to be a lot more efficient than a regular one. Because let's say I run the gitmap here, and it has all the updates. If within, let's say, one second or two seconds, or let's say we can have a window of five seconds. If within the five seconds window, there is no change, the same Git pull is asked. It will say that this is the Git. Let's put it as three seconds. Three seconds as default settings, but it can be changed as well. So add help text regarding this, and also it can be changed from the settings UI, so make sure you update the settings UI as well. Now, anytime we have the Git pull done, the CWD is fine. There is no need to mention CWD is a Git repo running pull. No need to. So you already mentioned where you are running it in the first tree view. That should be fine. Pulling changes, and then when it goes into the remote changes, it was checked two, three seconds ago. Three seconds, let's put it as three seconds ago. It already has run the cache. The second thing would be it would always check the last cache that is there. Even though it'd be cache, we count to not. So let's not cache with the cache, because even though the cache is the last cache, but it's still, there could be items that could be changed or merged in between. But if it's not possible, then we can actually add the hash to check what is the last cache that we have received. And I think there should be a pull Database (Database (Db)) or table already. That should actually keep this information. So confirm that and improve this practice. So if it has no changes, it should say, "No changes. Very good. Nice. That is the last cache that has been updated." Things like that. But yes, if it has changed, then it would share something like what you have done, what is the last cache, which branch, and also how many changes, dirty, modified, and untracked. That is really nice. In one liner. Really appreciated that one. So rather than checking Git status later on, if in the Git pull we have the updated and untracked files, it should show up here like this automatically, because these are the untracked files, so that the CLI does not have to be in this same time. So the task creation, I asked you before, the task creation should also give us the database path where it is. It should also give us the root Database (Database (Db)) information and some commands if we want to query the root Database (Database (Db)) as well inside this task JSON list. And make sure that these agent tasks can be easy to import, export as a JSON, YAML files. Add that feature as well. That should be inside its help. And also add this help to the LLM train section as well. LLM train, Muse should also give the gitmap skills to Muse. Give a path or file location. That way the gitmap skill can be learned by Muse as well. You can work on it. So with the gitmap find type search, always give the first line as how many found, what is the time it took. I know it's truncated, so probably I don't see the first line. But if it is not done, then do that like the gitmap AUM search, which shows nicely. That is a nice way to deal with the situation. The heat for the coding guideline is your steps B6 does not start with the heatmap. Can you please tell me heatmap is not created every decision that AI is taking exception? And also, can you tell me you have updated for the heatmap, right? You have C heatmap and heatmap, both uses the dry code you have taken from dry version of the code. I do see that you have some printing in place, update section. Rather than having that, why don't you have a framework that would automatically calculate this and print out this stuff so that we can reuse in the future? Can you do that? If you can, create that terminal package as generic as possible so that any code base can use it. It would have a header, things like that, and can be reused very easily. And also at the same time, I want you to take those own packages like the terminal coloring, the use red. I asked you that progress bar, this needs to be in a dry and separate generic package that anyone can use. Take that, and also the worker process, worker design pattern. Yes. So take all of these as a separate package, so that any code can use. Put this also inside the coding guideline as a Golang package section, right? So add that and make sure the code compiles in there. And also update its usability inside the coding guideline spec folder, you think? Check if we have any dry code section. If not, then we have to add it. The dry code section folder. Based on that, we need to have these principles mentioned. And also in the coding guideline, can we improve the error manage and how we have done it? Can we update a little bit inside the coding guideline, the error manage folder inside the spec, so that it is a lot more enhanced. And make sure we should not remove anything, any kind of solution.
```

---

## 2. Screenshot Forensics & Gap Analysis

### 2.1 GitMap Pull Redundant CWD Notice Removal & Caching
- **Screenshot Evidence**: `assets/screenshots/gitmap-efficiency-and-coding-guidelines-enhancement-01.png`
- **Existing Output**:
  ```text
  d:\work\gitmap > gitmap pull
  -> gitmap pull (cwd: D:\work\gitmap)
  -> cwd is a git repo — running plain `git pull` here
  ↓ pulling remote changes...
  ⚠ gitmap: up-to-date (8cc6949) [main] (dirty: +1 staged, +3 modified, +1 untracked)
  ```
- **Identified Defects**:
  1. `-> cwd is a git repo — running plain git pull here` is redundant since the CWD context was already established on the preceding line.
  2. Repeated `gitmap pull` calls within seconds re-run the complete external git network pull process without checking recent pull telemetry.
- **Target Architecture**:
  1. Remove the redundant `-> cwd is a git repo...` print in `cli/cmdpull/pull_dispatch_routes.go`.
  2. Query `store.GetRecentRepoPullHistory(repoPath, 1)` in `gitmap-pull.db`.
  3. If last pull occurred within TTL window (default 3s, configurable via setting `pull_cache_window_seconds`), check local working tree dirty state via `gitutil.InspectDirtyState(repoPath)`.
  4. If local state has no new dirty changes and remote fetch was satisfied $\le 3$s ago, return cached response with elapsed timing indication (e.g. `⚡ gitmap: cached (checked 2s ago) - up-to-date (8cc6949) [main]`).
  5. Add `pull_cache_window_seconds` setting constant to `cli/constants/constants_settings.go` and add Pull Cache Window input card to `src/pages/Settings.tsx`.

### 2.2 Untracked & Modified Files Inline Display in Pull Output
- **Screenshot Evidence**: `assets/screenshots/gitmap-efficiency-and-coding-guidelines-enhancement-02.png`
- **Existing State**:
  The user had to run a separate `git status -s` command to see which files are modified or untracked:
  ```text
  M .gitignore
  M cli/cmd/roottooling.go
  M cli/cmdlemon/lemon_cmd.go
  M cli/cmdlemon/lemon_cmd_service.go
  M cli/constants/constants_commands_lemon.go
  ?? cli/cmdlemon/daemon_settings.go
  ```
- **Target Architecture**:
  In `cli/cmdpull/pull_summary_single.go`, when `diag.IsDirty` is true:
  Directly render the dirty file entries (staged, modified, untracked) below the one-liner status, indented or formatted cleanly so the user never needs to execute a redundant `git status -s`.

### 2.3 Agent Task Init Database Info, Commands & Import/Export
- **Screenshot Evidence**: `assets/screenshots/gitmap-efficiency-and-coding-guidelines-enhancement-03.png`
- **Existing State**:
  `gitmap agent task init` JSON output only includes:
  `budget`, `parentTaskId`, `taskSlug`, `taskName`, `rootDbPath`, `runDirectory`, `status`.
- **Target Architecture**:
  1. Add `databasePath` (relative path to `agent-task.db`), `rootDbPath` (relative path to `gitmap.db`), `rootDbInfo` object, and a list of actionable `queryCommands` (e.g. `sqlite3 <dbPath> "SELECT * FROM tasks;"`, `gitmap agent subtask list --task-id <id>`).
  2. Implement `gitmap agent task export [--format json|yaml] [--task-id <id>] [--out <file>]` and `gitmap agent task import [--format json|yaml] --file <file>` with full JSON & YAML support (`gopkg.in/yaml.v3`).
  3. Register commands in `cli/cmdagent/agent_cmd.go` and document in CLI help text.

### 2.4 LLM Train Muse Skills Discovery
- **Target Architecture**:
  In `cli/cmd/llm/llm_train_muse.go` and `cli/cmd/llm/llm_train_muse_ui.go`:
  Add step 5 to `printTraversalChecklist` exposing the relative path to GitMap autonomous skills (`.agents/skills/gitmap/SKILL.md`) and command `gitmap cat .agents/skills/gitmap/SKILL.md`, allowing Muse to directly discover and ingest GitMap operational capabilities.

### 2.5 `gitmap find` First-Line Match Count & Timing
- **Screenshot Evidence**: `assets/screenshots/gitmap-efficiency-and-coding-guidelines-enhancement-04.png`, `-05.png`
- **Existing State**:
  `gitmap find` prints matched paths without a header, causing tool output truncations to obscure the hit count and execution duration.
- **Target Architecture**:
  In `cli/cmdfind/find_files.go`:
  Measure elapsed time with `time.Now()` and print the header on the very first line before matches:
  `Search completed: %d hit(s) in %v` (matching `gitmap aum search` format).

### 2.6 Generic Terminal Framework, Worker Pool & Coding Guidelines
- **Target Architecture**:
  1. Implement `cli/termgeneric/` (or `pkg/termgeneric/`):
     - `header.go`: Declarative box header renderer (`╭─ ... ──╮`) with dynamic padding and color styling.
     - `color.go`: Clean ANSI coloring helpers (Red, Green, Cyan, Yellow, Magenta, Dim, Reset).
     - `progress.go`: Thread-safe spinner and progress bar with deterministic width calculation and clean erase (`\r\033[2K`).
     - `worker_pool.go`: Generic concurrency worker pool pattern (`WorkerPool[T, R]`) with bounded queue capacity, context cancellation, and error aggregation.
  2. Document in `02-spec/02-coding-guidelines/14-dry-code-architecture/`:
     Author comprehensive specification for Golang DRY package standards, terminal UI framework, and worker concurrency patterns.
  3. Enhance `02-spec/03-error-manage/`:
     Add structured error handling conventions and CLI exit status mapping without removing any existing solutions.

---

## 3. Implementation Subtasks

1. `subtask-01-pull-caching-and-untracked-display`:
   - Remove redundant `-> cwd is a git repo` notice.
   - Implement 3s cache check via `gitmap-pull.db`.
   - Render modified & untracked files directly in single repo pull summary.
   - Add cache window setting to constants and `src/pages/Settings.tsx`.

2. `subtask-02-agent-task-db-info-import-export-and-find-header`:
   - Enhance `gitmap agent task init` JSON output with `databasePath`, `rootDbInfo`, and `queryCommands`.
   - Implement `gitmap agent task export` and `gitmap agent task import` (JSON/YAML).
   - Expose GitMap skills in `cli/cmd/llm/llm_train_muse.go`.
   - Add first-line hit count and timing header to `gitmap find`.

3. `subtask-03-generic-terminal-worker-package-and-coding-guidelines`:
   - Implement `cli/termgeneric` package (box header, ANSI colors, progress/spinner, worker pool).
   - Author `02-spec/02-coding-guidelines/14-dry-code-architecture/readme.md`.
   - Enhance `02-spec/03-error-manage/` with structured propagation notes.
