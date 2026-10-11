# Subtask 02: Merge-AI Staging Engine, Collision Sequencing & Commit Telemetry Tracking (Task-03)

> **Subtask ID:** `Subtask-02`  
> **Parent Plan:** `.ai-memory/plans/repository-creation-and-folder-structure-instructions.md`  
> **Target Path:** `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/02-merge-ai-engine-and-manifest-tracking.md`  
> **Status:** `PENDING`  
> **Assigned Tasks:** `Task-03` (Amalgamation Staging Engine, Collision Sequencing & Metadata Tracking)  
> **Assigned Subsystem:** `cli/cmdmergeai/`  
> **Related Architecture Spec:** `02-spec/21-app/repository-creation-and-folder-structure-instructions/01-architecture-spec.md`  
> **Owned Implementation Files:**  
> - `cli/cmdmergeai/merge_ai_core.go`  
> - `cli/cmdmergeai/merge_ai_staging.go`  
> - `cli/cmdmergeai/merge_ai_inspect.go`  
> - `cli/cmdmergeai/merge_ai_types.go`  
> - `cli/cmdmergeai/merge_ai_test.go`  

---

## 1. Objectives & Scope

This subtask implements the core amalgamation staging engine, chronological sorting, collision sequencing algorithm, and commit metadata extraction pipeline in `cli/cmdmergeai`:

1. **Chronological Source Precedence (`cli/cmdmergeai/merge_ai_core.go`):**
   Ensure all resolved source repositories are inspected and sorted chronologically by `LastCommitDate` in ascending order (earliest to latest) before staging begins. This establishes an objective historical baseline where older repositories provide foundation code (`01_`) and newer repositories represent iterative modifications (`02_`, `03_`).
2. **The Zero Replacement Invariant (`cli/cmdmergeai/merge_ai_staging.go`):**
   Enforce zero-replacement collision handling across all staged file trees. Under no circumstances may an existing file in the destination tree be overwritten. Conflicting files must be segregated using two-digit zero-padded prefixes (`01_<filename>`, `02_<filename>`, `03_<filename>`).
3. **Collision Detection & Two-Phase Renaming (`cli/cmdmergeai/merge_ai_staging.go`):**
   - For unique files: Copy directly to `<destDir>/<relPath>`, record in `placedOrigin`, and increment `uniqueFilesDirectlyPlaced`.
   - On first collision: Dynamically rename the previously placed canonical file to `<dir>/01_<basename>`, copy the incoming colliding file to `<dir>/02_<basename>`, and register both variants in `collisionMap`.
   - On subsequent collisions (3rd, 4th, etc.): Copy incoming file to `<dir>/%02d_<basename>` and append variant to existing `collisionMap[relPath].Variants`.
4. **Rich Commit Metadata Extraction (`cli/cmdmergeai/merge_ai_inspect.go`):**
   Execute Git plumbing commands (`git rev-parse`, `git log -1`) inside each source repository staging directory to extract:
   - 8-character short commit SHA and 40-character full commit SHA.
   - Author name and email formatted as `"Name <email>"`.
   - Commit timestamp in RFC 3339 format.
   - Commit subject line (message).
   - Contributed file count and active branch name.
5. **Granular File Provenance Mapping (`cli/cmdmergeai/merge_ai_staging.go`):**
   Construct a comprehensive `fileProvenance` map indexing every placed file in the destination repository to its source repository URL, source slug, origin commit SHA, and affirmative collision/placement booleans (`hasCollision`, `isDirectPlacement`).
6. **Automatic Git Index Staging (`cli/cmdmergeai/merge_ai_core.go`):**
   Execute `git add -A` in the destination workspace upon successful staging to stage all placed and sequenced files, manifest, and instructions into the Git index for immediate single-commit completion.
7. **Comprehensive Unit Testing (`cli/cmdmergeai/merge_ai_test.go`):**
   Write automated tests verifying chronological sorting, unique file placement, two-way collisions, three-way collisions, Zero Replacement Invariant enforcement, commit metadata extraction, and provenance mapping.

---

## 2. Detailed Technical Deliverables

### 2.1 Chronological Sorting Engine (`cli/cmdmergeai/merge_ai_core.go`)

In `RunMergeAI`, sort sources chronologically after inspection:

```go
// Sort chronologically (earliest to latest) to establish baseline precedence
sort.Slice(sources, func(i, j int) bool {
    return sources[i].LastCommitDate.Before(sources[j].LastCommitDate)
})
```

This guarantees deterministic sequencing independent of CLI argument order:
- `Source 1` (earliest commit date) $\rightarrow$ prefix `01_`
- `Source 2` (middle commit date) $\rightarrow$ prefix `02_`
- `Source 3` (latest commit date) $\rightarrow$ prefix `03_`

---

### 2.2 Collision Sequencing & Zero Replacement Staging (`cli/cmdmergeai/merge_ai_staging.go`)

#### 2.2.1 Core Staging Dispatcher (`stageAndSequenceFiles`)
Iterate through sorted sources and process each file in `src.FolderTree`:

```go
func stageAndSequenceFiles(destDir string, sources []SourceRepoConfig) ([]FileCollisionRecord, map[string]FileProvenanceRecord, int, error) {
	collisionMap := make(map[string]*FileCollisionRecord)
	placedOrigin := make(map[string]FileCollisionVariant)
	provenanceMap := make(map[string]FileProvenanceRecord)
	uniquePlaced := 0

	for seqIdx, src := range sources {
		seqPrefix := fmt.Sprintf("%02d", seqIdx+1)
		sourceSlug := filepath.Base(src.LocalStaging)
		headSha, _ := gitExec(src.LocalStaging, "rev-parse", "HEAD")
		shortSha := truncateCommitSha(headSha)

		for _, relPath := range src.FolderTree {
			errStage := stageSingleFile(destDir, sources, src, relPath, seqPrefix, sourceSlug, shortSha, collisionMap, placedOrigin, provenanceMap, &uniquePlaced)
			if errStage != nil {
				return nil, nil, 0, errStage
			}
		}
	}

	var collisionList []FileCollisionRecord
	for _, c := range collisionMap {
		collisionList = append(collisionList, *c)
	}

	sort.Slice(collisionList, func(i, j int) bool {
		return collisionList[i].CanonicalPath < collisionList[j].CanonicalPath
	})

	return collisionList, provenanceMap, uniquePlaced, nil
}
```

#### 2.2.2 Single File Staging Logic (`stageSingleFile`)
Evaluate disk existence and collision registry:
1. **Known Collision:** If `relPath` already exists in `collisionMap`, call `handleKnownCollision`.
2. **Unique File Placement:** If `destPath` does not exist on disk, call `handleUniqueFile`.
3. **First Collision:** If `destPath` exists on disk but is not yet in `collisionMap`, call `handleFirstCollision`.

#### 2.2.3 First Collision Handler (`handleFirstCollision`)
Enforces the **Zero Replacement Invariant**:
- The existing file at `destPath` (originally placed from Source 1) is **never overwritten**.
- Instead, it is renamed to `<dir>/01_<basename>`.
- The incoming file from Source 2 is copied to `<dir>/02_<basename>`.
- Both variants are registered in `collisionMap[relPath]`.
- Provenance map entries are updated for both `01_` and `02_` files.
- `uniquePlaced` counter is decremented.

```go
func handleFirstCollision(
	destDir string,
	sources []SourceRepoConfig,
	src SourceRepoConfig,
	relPath, destPath, seqPrefix, sourceSlug, headSha string,
	collisionMap map[string]*FileCollisionRecord,
	placedOrigin map[string]FileCollisionVariant,
	provenanceMap map[string]FileProvenanceRecord,
	uniquePlaced *int,
) error {
	orig := resolveOriginalVariant(sources, relPath, placedOrigin)
	dirName := filepath.Dir(relPath)
	baseName := filepath.Base(relPath)

	origSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", orig.Sequence, baseName)))
	origSeqPath := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", orig.Sequence, baseName))
	_ = os.Rename(destPath, origSeqPath)
	if *uniquePlaced > 0 {
		*uniquePlaced--
	}
	delete(placedOrigin, relPath)

	// Update provenance for the renamed 01 variant
	delete(provenanceMap, relPath)
	provenanceMap[origSeqRel] = FileProvenanceRecord{
		CanonicalPath:     relPath,
		SequencePrefix:    orig.Sequence,
		SourceRepoSlug:    orig.SourceRepo,
		CommitSha:         orig.Commit,
		HasCollision:      true,
		IsDirectPlacement: false,
	}

	newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
	newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
	errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), newSeqDest)
	if errCopy != nil {
		return errCopy
	}

	// Record provenance for the new sequenced variant
	provenanceMap[newSeqRel] = FileProvenanceRecord{
		CanonicalPath:     relPath,
		SequencePrefix:    seqPrefix,
		SourceRepoSlug:    sourceSlug,
		SourceRepoURL:     src.RepoURL,
		CommitSha:         headSha,
		HasCollision:      true,
		IsDirectPlacement: false,
	}

	collisionMap[relPath] = &FileCollisionRecord{
		CanonicalPath:  relPath,
		Status:         "pending_ai_consolidation",
		IsConsolidated: false,
		Variants: []FileCollisionVariant{
			{
				Sequence:   orig.Sequence,
				File:       origSeqRel,
				SourceRepo: orig.SourceRepo,
				Commit:     orig.Commit,
			},
			{
				Sequence:   seqPrefix,
				File:       newSeqRel,
				SourceRepo: sourceSlug,
				Commit:     headSha,
			},
		},
	}
	return nil
}
```

#### 2.2.4 Subsequent Collision Handler (`handleKnownCollision`)
When a 3rd or subsequent source conflicts on an already collided path:
- Copy incoming file to `<dir>/%02d_<basename>` (e.g. `03_package.json`).
- Append new `FileCollisionVariant` to `col.Variants`.
- Record entry in `provenanceMap`.

---

### 2.3 Commit Metadata Extraction (`cli/cmdmergeai/merge_ai_inspect.go`)

Extract rich commit telemetry from each source repository:

```go
func extractSourceCommitMetadata(repoURL, stagingDir string, contributedFiles int) SourceCommitMetadata {
	slug := filepath.Base(stagingDir)
	headSha, _ := gitExec(stagingDir, "rev-parse", "HEAD")
	shortSha := truncateCommitSha(headSha)
	branch, _ := gitExec(stagingDir, "rev-parse", "--abbrev-ref", "HEAD")
	author, _ := gitExec(stagingDir, "log", "-1", "--format=%an <%ae>")
	dateStr, _ := gitExec(stagingDir, "log", "-1", "--format=%cI")
	msg, _ := gitExec(stagingDir, "log", "-1", "--format=%s")

	parsedDate, errDate := time.Parse(time.RFC3339, dateStr)
	if errDate != nil {
		parsedDate = time.Now().UTC()
	}

	return SourceCommitMetadata{
		RepoSlug:              slug,
		RepoURL:               repoURL,
		Branch:                branch,
		HeadCommitSha:         headSha,
		ShortCommitSha:        shortSha,
		Author:                author,
		CommitDate:            parsedDate,
		CommitMessage:         msg,
		TotalFilesContributed: contributedFiles,
	}
}
```

---

### 2.4 Manifest Assembly & Git Index Staging (`cli/cmdmergeai/merge_ai_core.go`)

Assemble the complete `MergeAIManifest` including `commitMetadata` and `fileProvenance`:

```go
manifest := assembleManifest(resolvedDest, sources, collisions, commitMetaList, provenanceMap, uniqueCount, skippedSources)

if errManifest := WriteMergeAIManifest(resolvedDest.LocalPath, manifest); errManifest != nil {
    return fmt.Errorf("failed to write merge-ai-manifest.json: %w", errManifest)
}

if errInstruction := WriteMergeAIInstruction(resolvedDest.LocalPath, collisions, skippedSources...); errInstruction != nil {
    return fmt.Errorf("failed to write instruction.md: %w", errInstruction)
}

// Stage all changes into git index
_, _ = gitExec(resolvedDest.LocalPath, "add", "-A")
```

---

## 3. Step-by-Step Implementation Sequence

```
┌──────────────────────────────────────────────────────────────┐
│  Phase 1: Inspection & Commit Telemetry Enhancement          │
├──────────────────────────────────────────────────────────────┤
│  1.1 Implement extractSourceCommitMetadata in                │
│      cli/cmdmergeai/merge_ai_inspect.go                      │
│  1.2 Capture short SHA, full SHA, author, date, message      │
│  1.3 Attach SourceCommitMetadata to inspectSingleSource      │
└──────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────┐
│  Phase 2: Chronological Sorting & Zero-Replacement Staging   │
├──────────────────────────────────────────────────────────────┤
│  2.1 Verify chronological sort in RunMergeAI                 │
│      (sort.Slice by LastCommitDate ascending)                │
│  2.2 Update stageAndSequenceFiles signature to return        │
│      (collisions, provenanceMap, uniqueCount, error)         │
│  2.3 Update handleFirstCollision to rename canonical file    │
│      to 01_ and copy incoming to 02_                         │
│  2.4 Update handleKnownCollision for 03_, 04_, etc.          │
│  2.5 Populate fileProvenance for all staged files            │
└──────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────┐
│  Phase 3: Manifest Assembly & Automated Staging              │
├──────────────────────────────────────────────────────────────┤
│  3.1 Update assembleManifest to attach commitMetadata &      │
│      fileProvenance                                          │
│  3.2 Stage staging tree into Git index (git add -A)          │
│  3.3 Add unit tests in cli/cmdmergeai/merge_ai_test.go:      │
│      - TestChronologicalSorting                              │
│      - TestZeroReplacementCollisionSequencing                │
│      - TestThreeWayCollisionSequencing                       │
│      - TestFileProvenanceMapping                             │
│      - TestSourceCommitMetadataExtraction                    │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. Verification Matrix & Acceptance Criteria

| Checkpoint | Target Component | Scenario / Test Case | Expected Outcome | Verification Metric |
| :--- | :--- | :--- | :--- | :--- |
| **AC-01** | Chronological Sorting | Sources with dates `[Aug 2026, May 2026, Oct 2026]` | Sorted to `[May 2026, Aug 2026, Oct 2026]` | Unit test assertion |
| **AC-02** | Zero Replacement Invariant | Source 1 and Source 2 share `readme.md` | Source 1 version becomes `01_readme.md`; Source 2 version becomes `02_readme.md`; neither is overwritten | Both files exist with exact source bytes |
| **AC-03** | Three-Way Collision | Sources 1, 2, 3 share `package.json` | Destination contains `01_package.json`, `02_package.json`, `03_package.json`; `canonicalPath: "package.json"` | `len(Variants) == 3` in manifest |
| **AC-04** | Unique File Placement | Source 1 has `src/unique.ts` | Placed at `src/unique.ts`; `hasCollision: false`, `isDirectPlacement: true` | Direct placement verified |
| **AC-05** | Commit Telemetry | Source repo with commit `HEAD` | `headCommitSha`, `shortCommitSha`, `author`, `commitDate`, `commitMessage` populated | Schema check in `merge-ai-manifest.json` |
| **AC-06** | File Provenance Map | Any staged file in destination | Resolves to source repo URL, commit SHA, and sequence prefix | `provenanceMap` contains valid entry |
| **AC-07** | Git Index Staging | End of staging execution | All files staged in Git index (`git status` shows all changes staged) | Exit code 0 from `git diff --cached --quiet` is false (changes staged) |
| **AC-08** | Positive Boolean Hygiene | Manifest JSON & Go structs | All boolean keys use positive formulation (`isConsolidated`, `hasCollision`, `isDirectPlacement`) | Zero negative boolean keys |
