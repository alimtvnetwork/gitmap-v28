# Subtask 01: Destination Target Resolution, Workspace Cleanup & AI Manifest Specification (Task-01 & Task-02)

> **Subtask ID:** `Subtask-01`  
> **Parent Plan:** `.ai-memory/plans/repository-creation-and-folder-structure-instructions.md`  
> **Target Path:** `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/01-workspace-cleanup-and-manifest-spec.md`  
> **Status:** `PENDING`  
> **Assigned Tasks:** `Task-01` (Destination Target Resolution & Staging Cleanup) & `Task-02` (Manifest Schema & Instruction Generator)  
> **Assigned Subsystems:** `cli/cmdresolver/`, `cli/cmdmergeai/`  
> **Related Architecture Spec:** `02-spec/21-app/repository-creation-and-folder-structure-instructions/01-architecture-spec.md`  
> **Owned Implementation Files:**  
> - `cli/cmdresolver/repo_feature.go`  
> - `cli/cmdresolver/repo_feature_workspace.go`  
> - `cli/cmdresolver/repo_feature_test.go`  
> - `cli/cmdmergeai/merge_ai_types.go`  
> - `cli/cmdmergeai/merge_ai_manifest.go`  
> - `cli/cmdmergeai/merge_ai_test.go`  

---

## 1. Objectives & Scope

This subtask delivers the foundational destination resolution, staging cleanup, manifest schema, and instruction generation subsystems for `gitmap merge-ai`:

1. **Universal Slug Normalization (`cli/cmdresolver/repo_feature.go` - Task-01):**
   Harden `cmdresolver.SanitizeSlug` to deterministically normalize arbitrary input expressions, including multi-word spaced titles, uppercase identifiers, versions, and symbols (e.g. `"Icon Studio 2026 V1"` $\rightarrow$ `"icon-studio-2026-v1"`).
2. **Bare Slug Target Resolution & Git Initialization (`cli/cmdresolver/repo_feature.go` - Task-01):**
   Implement automated bare-slug provisioning. When an input target matches no remote URL or existing filesystem path, compute `<workspaceRoot>/<canonicalSlug>`, create the directory, and initialize a new Git repository (`git init -b main`) with affirmative flags (`isNewRepo: true`, `isCloned: false`).
3. **Workspace Staging Cleanup Engine (`cli/cmdmergeai/merge_ai_staging.go` - Task-01):**
   Enforce safe directory sanitization in the resolved staging destination. Purge all preexisting stale files and subdirectories while strictly preserving the `.git` directory and metadata intact.
4. **Manifest Schema & Commit Metadata Modeling (`cli/cmdmergeai/merge_ai_types.go` - Task-02):**
   Define the typed Go contract for `MergeAIManifest`, extending attributes, destination parameters, sequence records, skipped sources, collision records, and adding structured commit metadata (`SourceCommitMetadata`) and granular file provenance tracking (`FileProvenanceRecord`) with strictly positive booleans (`isNewRepo`, `isCloned`, `hasCollision`, `isDirectPlacement`, `isConsolidated`, `isFile`).
5. **Deterministic Manifest Serialization (`cli/cmdmergeai/merge_ai_manifest.go` - Task-02):**
   Implement `WriteMergeAIManifest` to format and persist `merge-ai-manifest.json` into the destination repository root using 2-space indentation and RFC 3339 timestamps.
6. **Authoritative AI Instruction Generator (`cli/cmdmergeai/merge_ai_manifest.go` - Task-02):**
   Implement `WriteMergeAIInstruction` to generate `<destDir>/instruction.md`, dynamically formatting detected file collisions, skipped sources, non-negotiable consolidation rules, and a step-by-step actionable checklist for the downstream AI consolidation agent.

---

## 2. Detailed Technical Deliverables

### 2.1 Task-01: Slug Normalization & Target Resolution Subsystem

#### 2.1.1 Slug Normalization Pipeline (`SanitizeSlug`)
`cmdresolver.SanitizeSlug` must handle arbitrary project titles:
- Input: `"Icon Studio 2026 V1"`
- Transformations:
  1. Lowercase conversion: `"icon studio 2026 v1"`
  2. Separator replacement: `"icon-studio-2026-v1"`
  3. Non-alphanumeric stripping: removes all characters except `[a-z0-9-]`
  4. Redundant hyphen collapsing: collapses multiple hyphens into single hyphens
  5. Boundary trimming: removes leading/trailing hyphens
- Output: `"icon-studio-2026-v1"`

```go
// SanitizeSlug normalizes any target slug into a clean kebab-case identifier.
func SanitizeSlug(slug string) string {
	clean := strings.ToLower(strings.TrimSpace(slug))
	clean = strings.ReplaceAll(clean, " ", "-")
	clean = strings.ReplaceAll(clean, "_", "-")

	var b strings.Builder
	for _, r := range clean {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	clean = b.String()

	for strings.Contains(clean, "--") {
		clean = strings.ReplaceAll(clean, "--", "-")
	}

	trimmed := strings.Trim(clean, "-")
	if trimmed == "" {
		return "repo"
	}
	return trimmed
}
```

#### 2.1.2 Bare Slug Repository Provisioning (`resolveBareSlugFeature`)
When `cmdresolver.ResolveRepoFeature` evaluates an input that is neither a remote Git URL nor an existing directory on disk:
1. Normalize slug via `SanitizeSlug`.
2. Compute destination path under the canonical workspace root.
3. Ensure directory exists with permission mode `0755`.
4. If `<destPath>/.git` does not exist:
   - Execute `git init -b main` in `<destPath>`.
   - Set `isNewRepo = true`.
5. Return `ResolvedRepoFeature`:
   - `TargetType: TargetTypeBareSlug`
   - `CanonicalSlug: cleanSlug`
   - `LocalPath: destPath`
   - `IsNewRepo: true`
   - `IsCloned: false`

#### 2.1.3 Staging Directory Cleanup (`cleanDestStagingDirectory`)
Before copying files into the resolved staging repository, purge existing files from previous runs while preserving `.git`:
```go
func cleanDestStagingDirectory(destDir string) error {
	entries, err := os.ReadDir(destDir)
	if err != nil {
		return fmt.Errorf("failed to read staging directory: %w", err)
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		if errRemove := os.RemoveAll(filepath.Join(destDir, entry.Name())); errRemove != nil {
			return fmt.Errorf("failed to remove staging file %s: %w", entry.Name(), errRemove)
		}
	}
	return nil
}
```

---

### 2.2 Task-02: Manifest Specification & Instruction Generation Subsystem

#### 2.2.1 Data Struct Definitions (`cli/cmdmergeai/merge_ai_types.go`)
Extend `cli/cmdmergeai/merge_ai_types.go` to support commit metadata, provenance maps, and positive booleans:

```go
package cmdmergeai

import "time"

// SourceCommitMetadata captures granular commit telemetry per source repository.
type SourceCommitMetadata struct {
	RepoSlug              string    `json:"repoSlug"`
	RepoURL               string    `json:"repoUrl"`
	Branch                string    `json:"branch"`
	HeadCommitSha         string    `json:"headCommitSha"`
	ShortCommitSha        string    `json:"shortCommitSha"`
	Author                string    `json:"author"`
	CommitDate            time.Time `json:"commitDate"`
	CommitMessage         string    `json:"commitMessage"`
	TotalFilesContributed int       `json:"totalFilesContributed"`
}

// FileProvenanceRecord tracks the exact origin of a file placed in the destination tree.
type FileProvenanceRecord struct {
	CanonicalPath     string `json:"canonicalPath"`
	SequencePrefix    string `json:"sequencePrefix"`
	SourceRepoSlug    string `json:"sourceRepoSlug"`
	SourceRepoURL     string `json:"sourceRepoUrl"`
	CommitSha         string `json:"commitSha"`
	HasCollision      bool   `json:"hasCollision"`
	IsDirectPlacement bool   `json:"isDirectPlacement"`
}

// FileCollisionVariant models one sequenced instance of a conflicting file.
type FileCollisionVariant struct {
	Sequence      string `json:"sequence"`
	File          string `json:"file"`
	SourceRepo    string `json:"sourceRepo"`
	Commit        string `json:"commit"`
	CommitDate    string `json:"commitDate,omitempty"`
	CommitMessage string `json:"commitMessage,omitempty"`
}

// FileCollisionRecord models a conflicting path requiring AI consolidation.
type FileCollisionRecord struct {
	CanonicalPath  string                 `json:"canonicalPath"`
	Variants       []FileCollisionVariant `json:"variants"`
	Status         string                 `json:"status"` // "pending_ai_consolidation", "consolidated"
	IsConsolidated bool                   `json:"isConsolidated"`
}

// MergeAIManifest models the root merge-ai-manifest.json artifact.
type MergeAIManifest struct {
	Attributes struct {
		GeneratedAt       string `json:"generatedAt"`
		GitMapVersion     string `json:"gitMapVersion"`
		Tool              string `json:"tool"`
		CollisionStrategy string `json:"collisionStrategy"`
	} `json:"attributes"`
	Data struct {
		Destination struct {
			TargetType   string `json:"targetType"`
			Slug         string `json:"slug"`
			ResolvedPath string `json:"resolvedPath"`
			OriginURL    string `json:"originUrl,omitempty"`
			IsNewRepo    bool   `json:"isNewRepo"`
			IsCloned     bool   `json:"isCloned"`
		} `json:"destination"`
		RepoSequence              []SourceRepoSequenceEntry       `json:"repoSequence"`
		SkippedSources            []SkippedSourceItem             `json:"skippedSources,omitempty"`
		FileCollisions            []FileCollisionRecord           `json:"fileCollisions"`
		CommitMetadata            []SourceCommitMetadata          `json:"commitMetadata,omitempty"`
		FileProvenance            map[string]FileProvenanceRecord `json:"fileProvenance,omitempty"`
		UniqueFilesDirectlyPlaced int                             `json:"uniqueFilesDirectlyPlaced"`
	} `json:"data"`
}
```

#### 2.2.2 Manifest Serialization (`cli/cmdmergeai/merge_ai_manifest.go`)
Implement serialization of `merge-ai-manifest.json`:
- Indentation: 2 spaces.
- Target path: `<destDir>/merge-ai-manifest.json`.
- Permissions: mode `0644`.

#### 2.2.3 AI Instruction Generation (`cli/cmdmergeai/merge_ai_manifest.go`)
Implement `WriteMergeAIInstruction`:
- Header with clear warning and execution summary.
- Section 1: Non-Negotiable Rules:
  1. Preserve All Business Logic.
  2. Flatten Sequence Prefixes (`01_*`, `02_*` $\rightarrow$ canonical name).
  3. Reference `merge-ai-manifest.json`.
  4. Zero Code Deletion.
  5. Affirmative Verification.
- Section 2: Skipped Sources (if any).
- Section 3: Detected File Collisions table with canonical path, variant filename, source repo, and commit SHA.
- Section 4: Step-by-Step AI Action Checklist.

---

## 3. Step-by-Step Implementation Sequence

```
┌──────────────────────────────────────────────────────────────┐
│  Phase 1: Task-01 Destination Resolution & Cleanup Engine    │
├──────────────────────────────────────────────────────────────┤
│  1.1 Add slug tests in cli/cmdresolver/repo_feature_test.go   │
│      - "Icon Studio 2026 V1" -> "icon-studio-2026-v1"        │
│      - Punctuation, symbols, and empty string fallbacks       │
│  1.2 Verify SanitizeSlug implementation in repo_feature.go    │
│  1.3 Update resolveBareSlugFeature with positive booleans    │
│  1.4 Harden cleanDestStagingDirectory in merge_ai_staging.go │
│      - Preserve .git directory                               │
│      - Return wrapped error on failure                       │
└──────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────┐
│  Phase 2: Task-02 Manifest & Instruction Subsystem           │
├──────────────────────────────────────────────────────────────┤
│  2.1 Extend structs in cli/cmdmergeai/merge_ai_types.go      │
│      - SourceCommitMetadata, FileProvenanceRecord            │
│      - Positive booleans: isNewRepo, isCloned, hasCollision  │
│  2.2 Update WriteMergeAIManifest in merge_ai_manifest.go     │
│      - Serialize commitMetadata and fileProvenance           │
│  2.3 Update WriteMergeAIInstruction in merge_ai_manifest.go  │
│      - Format collision variants, commit context, dates      │
│      - Dynamic skipped source list                           │
│      - 8-step AI consolidation checklist                     │
│  2.4 Add unit tests in cli/cmdmergeai/merge_ai_test.go       │
│      - Manifest serialization and schema round-trip          │
│      - Instruction markdown formatting                       │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. Verification Matrix & Acceptance Criteria

| Checkpoint | Target Component | Input / Scenario | Expected Outcome | Verification Metric |
| :--- | :--- | :--- | :--- | :--- |
| **AC-01** | `cmdresolver.SanitizeSlug` | `"Icon Studio 2026 V1"` | Returns `"icon-studio-2026-v1"` | Unit test pass in `repo_feature_test.go` |
| **AC-02** | `cmdresolver.SanitizeSlug` | `"   Project__Core (2026)   "` | Returns `"project-core-2026"` | Unit test pass in `repo_feature_test.go` |
| **AC-03** | `cmdresolver.ResolveRepoFeature` | Non-existent bare slug `"icon-studio-2026-v1"` | Creates folder, initializes `git init -b main`, `isNewRepo: true`, `isCloned: false` | Directory exists, `.git` exists, correct flags |
| **AC-04** | `cleanDestStagingDirectory` | Directory containing `.git` and stale files | Stale files deleted, `.git` untouched | Filesystem inspection |
| **AC-05** | `WriteMergeAIManifest` | Structured `MergeAIManifest` object | Writes valid `<destDir>/merge-ai-manifest.json` with commit metadata and provenance | JSON unmarshal verification |
| **AC-06** | `WriteMergeAIInstruction` | Collision list with 2 variants | Writes `<destDir>/instruction.md` containing variant filenames, commit hashes, and checklist | Markdown substring verification |
| **AC-07** | Boolean Hygiene | All struct fields and JSON output | Strictly positive booleans (`isNewRepo`, `isCloned`, `hasCollision`, `isDirectPlacement`, `isConsolidated`) | Zero negative boolean keys |
