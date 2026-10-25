# Subtask 03: RC Confusion Status Sync and File Context View

- **Subtask Slug:** `03-subtask-rc-confusion-status-sync-and-file-view`
- **Parent Plan:** `.ai-memory/plans/pending/327-heatmap-status-improvement-and-ai-confusion-handling.md`
- **Architecture Spec:** `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/02-component-and-cli-spec.md`
- **Target Subsystem:** `cli/store`, `cli/cmd`
- **Status:** `COMPLETED`
- **Exclusive Owned Files:**
  - `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/02-component-and-cli-spec.md`
  - `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/03-subtask-rc-confusion-status-sync-and-file-view.md`
  - `cli/store/repo_cache_sequence_store.go`
  - `cli/store/ai_confusion_store.go`
  - `cli/cmd/rc_ai_confusion.go`

---

## 1. Objectives

1. Implement `cli/store/repo_cache_sequence_store.go`:
   - Define `RepoCacheSequenceIndex` table in the root SQLite database (`gitmap.db` via `store.BinaryDataDir()`).
   - Implement `ResolveRepoCacheSequenceFolder(repoSlug string, repoCacheBaseDir string) (string, error)` with two-tier lookup: root DB cache hit, followed by disk directory scan and sequential number generation (`XX-<repoSlug>`).
   - Implement `ResolveRepoConfusionDbPath(repoSlug string) (string, error)` pointing to `<repo-cache>/XX-<repoSlug>/ai-confusion.db`.
2. Update `cli/store/ai_confusion_store.go`:
   - Enhance `InsertAiConfusion` to transactionally insert master `AiConfusion` and child `AiConfusionFile` rows.
   - Enhance `GetAiConfusionByUuidOrId` and `ListAiConfusions` to query and populate `rec.ReferencedFiles`.
   - Add `FindAiConfusionByFilePath(filePath string) (*AiConfusionRecord, error)` to lookup confusion episodes by exact or partial relative file paths.
3. Update `cli/cmd/rc_ai_confusion.go`:
   - Implement `gitmap rc ai-confusion status` (and default with no args) to trigger safe repo-cache sync and display active confusion status cards.
   - Implement `gitmap rc ai-confusion view <filepath>` (or `--file <filepath>`) to resolve relative file paths to confusion episodes and render complete agent context cards.
4. Enforce strict coding guidelines:
   - <= 15 line functions (<= 8 lines preferred).
   - Zero nested `if` statements.
   - Positive boolean naming (`is*`, `has*`).
   - Blank line after `}`, before `if`, and before `return`.
   - Leaf error sovereignty.
   - NO BUILD, NO TEST during editing.

---

## 2. Implementation Steps

### Step 1: Sequence Folder Store (`cli/store/repo_cache_sequence_store.go`)
- [x] Connect to root SQLite database `gitmap.db` via `BinaryDataDir()`.
- [x] Create table `RepoCacheSequenceIndex` with `(RepoSlug PRIMARY KEY, SequenceFolder, CachedPath, UpdatedAt)`.
- [x] Implement cached path lookup with filesystem existence verification.
- [x] Implement directory scanning for existing `*-<repoSlug>` prefixes.
- [x] Implement automatic sequence numbering (`maxSeq + 1`) formatted as `%02d-<repoSlug>`.
- [x] Implement `ResolveRepoConfusionDbPath` helper.

### Step 2: 1-to-N Transactional Persistence (`cli/store/ai_confusion_store.go`)
- [x] Wrap `InsertAiConfusion` in a SQLite transaction.
- [x] Insert into `AiConfusion` and obtain `AiConfusionId`.
- [x] Insert `rec.ReferencedFiles` into `AiConfusionFile`.
- [x] Populate `rec.ReferencedFiles` on single and list queries.
- [x] Implement `FindAiConfusionByFilePath(filePath string) (*AiConfusionRecord, error)`.

### Step 3: CLI Status & Context View (`cli/cmd/rc_ai_confusion.go`)
- [x] Route empty arguments and `"status"` subcommand to `runAiConfusionStatus`.
- [x] Run safe synchronization against repo-cache.
- [x] Query active confusion episodes and render formatted status cards.
- [x] Detect file paths in `runAiConfusionView` positional arguments and `--file` flags.
- [x] Call `FindAiConfusionByFilePath` and render detailed agent thinking, spec, task, and code snippet cards.

---

## 3. Verification Criteria

- [x] Root SQLite table `RepoCacheSequenceIndex` is created idempotently.
- [x] Sequence folders are properly cached and resolved.
- [x] Multiple referenced files are persisted and retrieved with 1-to-N relationships.
- [x] File path lookups accurately resolve confusion records across path variants.
- [x] Default invocation displays active confusion cards cleanly.
