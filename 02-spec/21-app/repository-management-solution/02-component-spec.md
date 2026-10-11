# Component Specification: Unified Fleet Control, Remediation Options & Table Rendering

- **Specification Slug:** `repository-management-solution`
- **Spec File:** `02-spec/21-app/repository-management-solution/02-component-spec.md`
- **Status:** `APPROVED`
- **Architecture Reference:** `02-spec/21-app/repository-management-solution/01-architecture-spec.md`
- **Subtask References:**
  - `.ai-memory/plans/subtasks/repository-management-solution/03-unified-fleet-control.md`
  - `.ai-memory/plans/subtasks/repository-management-solution/04-test-and-verify.md`
- **Target Module:** `cli/cmdpending`
- **Target Version:** `v6.529.0`

---

## 1. Executive Summary & Architectural Scope

This component specification defines the technical modifications to `cli/cmdpending` to upgrade repository remediation hints from dual-option to three-option unified fleet control, enforce strict terminal box width alignment, and guarantee 100% parity between CLI rendering and JSON telemetry payloads.

### 1.1 Core Objectives
1. **Three-Option Remediation**: Expand remediation options for dirty repositories by introducing Option 3 (`gitmap sends cp all "wip: save changes"`), providing immediate fleet-wide batch execution alongside single-repository commit/push (Option 1) and stash (Option 2).
2. **Subtree Connector Rendering**: Update `renderTreeRemediationHints` to render standard tree branches (`├──`) for intermediate items and a terminal branch (`└──`) for Option 3.
3. **Exact Rune-Width Formatting**: Guarantee that `formatTreeLine` strictly formats each line to 78 inner runes (80 characters outer box with border vertical bars `│`), avoiding terminal overflow and line wrapping.
4. **Table Footer & Summary Metrics**: Standardize the summary metrics banner and footer messaging to reflect fleet-wide remediation capabilities.

---

## 2. Component 1: Terminal Subtree Rendering (`cli/cmdpending/pending_commits_render.go`)

### 2.1 Subtree Rendering Algorithm (`renderTreeRemediationHints`)
When `rec.IsDirty` is true, the table renderer emits structured subtree lines beneath the repository row. Each option line begins with an indent and a tree connector:

```go
func renderTreeRemediationHints(rec RepoPendingCommitRecord, border, reset string) {
	totalOpts := len(rec.RemediationOptions)
	for idx, opt := range rec.RemediationOptions {
		prefix := "├──"
		if idx == totalOpts-1 {
			prefix = "└──"
		}
		rawLine := fmt.Sprintf("   %s Option %d: %s", prefix, opt.OptionNumber, opt.Command)
		formattedLine := formatTreeLine(rawLine)
		fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedLine, border, reset)
	}
}
```

#### Visual Layout Contract (80-Column Terminal Box)
```text
  ┌──────────────────────────────────────────────────────────────────────────────┐
  │ REPOSITORY              VER/BRANCH        UNCOMMITTED    UNPUSHED  STATUS    │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ my-service              v1.2.0/main                 3           0  ● PEND    │
  │    ├── Option 1: gitmap sends cp my-service "wip: save changes"              │
  │    ├── Option 2: gitmap stash my-service                                     │
  │    └── Option 3: gitmap sends cp all "wip: save changes"                     │
  └──────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Rune-Accurate Line Formatting (`formatTreeLine`)
Terminal box formatting requires exact rune length budgeting:
- **Inner Box Width:** Exactly 78 runes.
- **Short Lines:** Right-padded with spaces up to 78 runes.
- **Long Lines:** Truncated at 77 runes and terminated with a single-character ellipsis (`…`), totaling 78 runes.
- **Function Length:** Must not exceed 15 lines.

```go
func formatTreeLine(raw string) string {
	runes := []rune(raw)
	const boxWidth = 78
	if len(runes) > boxWidth {
		return string(runes[:boxWidth-1]) + "…"
	}
	return raw + strings.Repeat(" ", boxWidth-len(runes))
}
```

---

## 3. Component 2: Remediation Options Generator (`cli/cmdpending/pending_commits_cmd.go`)

### 3.1 Refactoring `createRemediationOptionsPair` to `createRemediationOptionsList`
The legacy constructor `createRemediationOptionsPair` only returned two options. It is refactored into `createRemediationOptionsList` to generate three options in strict sequential order.

```go
func buildRepoRemediationOptions(rec RepoPendingCommitRecord) []RemediationOption {
	if !rec.IsDirty {
		return []RemediationOption{}
	}
	targetPath := rec.RelativePath
	if targetPath == "" || targetPath == "." || looksLikeHashOrOpaqueID(targetPath) {
		targetPath = rec.RepoName
	}
	return createRemediationOptionsList(targetPath)
}
```

### 3.2 Modular Option Constructors (Function Length <= 15 Lines)
To comply with function size limits, option construction is decomposed into discrete builders:

```go
func createRemediationOptionsList(targetPath string) []RemediationOption {
	target := quoteRemediationTarget(targetPath)
	return []RemediationOption{
		buildOptionCommitPush(target),
		buildOptionStash(target),
		buildOptionFleetControl(),
	}
}

func buildOptionCommitPush(target string) RemediationOption {
	return RemediationOption{
		OptionNumber: 1,
		Label:        "Commit & Push WIP",
		Command:      fmt.Sprintf(`gitmap sends cp %s "wip: save changes"`, target),
	}
}

func buildOptionStash(target string) RemediationOption {
	return RemediationOption{
		OptionNumber: 2,
		Label:        "Stash WIP",
		Command:      fmt.Sprintf(`gitmap stash %s`, target),
	}
}

func buildOptionFleetControl() RemediationOption {
	return RemediationOption{
		OptionNumber: 3,
		Label:        "Fleet Commit & Push All",
		Command:      `gitmap sends cp all "wip: save changes"`,
	}
}
```

### 3.3 Target Quoting Invariant
If a repository path contains whitespace, it must be quoted safely:

```go
func quoteRemediationTarget(s string) string {
	if strings.Contains(s, " ") && !strings.HasPrefix(s, "\"") {
		return fmt.Sprintf(`"%s"`, s)
	}
	return s
}
```

---

## 4. Component 3: Table Footer & Summary Metrics

### 4.1 Summary Banner (`renderTableSummary`)
The table summary displays four aggregated metrics across all inspected repositories:
- `TotalReposScanned`: Count of all repositories in workspace.
- `TotalDirtyRepos`: Count of repositories with uncommitted files.
- `TotalUncommittedFiles`: Total untracked, modified, and staged files.
- `TotalUnpushedCommits`: Total unpushed commits ahead of upstream.

```go
func renderTableSummary(payload PendingCommitsPayload, border, reset string) {
	summaryText := fmt.Sprintf(" Scanned: %d repos │ Dirty: %d repos │ Uncommitted: %d files │ Unpushed: %d",
		payload.TotalReposScanned, payload.TotalDirtyRepos, payload.TotalUncommittedFiles, payload.TotalUnpushedCommits)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(summaryText, 78), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}
```

### 4.2 Table Footer (`renderTableFooter`)
The table footer presents the canonical fleet-wide batch command:

```go
func renderTableFooter(border, reset string) {
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	footerMsg := ` Fleet Remediation: gitmap sends cp all "wip: save changes"`
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(footerMsg, 78), border, reset)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}
```

---

## 5. Telemetry & JSON Schema Parity

### 5.1 JSON Payload Contract (`pending_commits_types.go`)
CLI options and JSON telemetry share identical data structures. `RepoPendingCommitRecord` serializes the options under `remediationOptions`:

```go
type RemediationOption struct {
	OptionNumber int    `json:"optionNumber"`
	Label        string `json:"label"`
	Command      string `json:"command"`
}
```

### 5.2 JSON Telemetry Example
```json
{
  "repoName": "my-service",
  "relativePath": "services/my-service",
  "currentBranch": "main",
  "version": "v1.2.0",
  "shortVersionBranch": "v1.2.0/main",
  "isDirty": true,
  "isClean": false,
  "hasUncommitted": true,
  "hasUnpushed": false,
  "hasUpstream": true,
  "totalUncommitted": 3,
  "remediationOptions": [
    {
      "optionNumber": 1,
      "label": "Commit & Push WIP",
      "command": "gitmap sends cp services/my-service \"wip: save changes\""
    },
    {
      "optionNumber": 2,
      "label": "Stash WIP",
      "command": "gitmap stash services/my-service"
    },
    {
      "optionNumber": 3,
      "label": "Fleet Commit & Push All",
      "command": "gitmap sends cp all \"wip: save changes\""
    }
  ]
}
```

---

## 6. Coding Standards & Invariants Matrix

| Requirement | Enforced Rule | Verification Tool |
| :--- | :--- | :--- |
| **Max File Size** | All Go files must be <= 300 lines | `python 03-ai-scripts/52-file-size-check.py --limit 300` |
| **Max Function Length** | All functions must be <= 15 lines | `python linter-scripts/check-function-lengths.py` |
| **Positive Booleans** | Only `is` and `has` prefixes allowed (`isDirty`, `isClean`, `hasUpstream`) | `python linter-scripts/check-boolean-guidelines.py` |
| **No Nested Ifs** | Nested if depth > 1 is strictly forbidden | `python linter-scripts/check-nested-ifs.py` |
| **Strict Relative Paths** | Zero absolute paths or `file:///` URIs | `python linter-scripts/check-relative-paths.py` |
| **Git Command Ban** | Zero direct git command execution | CI pre-commit gate |
