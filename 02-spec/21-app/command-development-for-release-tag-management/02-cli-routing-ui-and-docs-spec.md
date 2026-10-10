# Component and CLI Specification: Release Tag & Broken Release Management (`gitmap fix release tags`)

> **/goal** Specify the CLI routing, flag specifications, Terminal UI components, interactive confirmation guards, root menu clusters, and embedded markdown documentation for `gitmap fix release tags` (aliases: `fix-release-tags`, `frt`).
> **/learn** Handle multi-word subcommand token streams in `cmdfix.RunFix`, enforce positive boolean logic across flag models, prevent non-interactive stdin hangs via `isInteractiveStdin()` guards, render high-contrast 6-column diagnostic tables using `termout.PrintTable`, and integrate cleanly into Muse 5-cluster help architecture.

- **Specification Slug:** `command-development-for-release-tag-management`
- **Spec Document:** `02-spec/21-app/command-development-for-release-tag-management/02-cli-routing-ui-and-docs-spec.md`
- **Spec Status:** `APPROVED`
- **Architecture Reference:** `02-spec/21-app/command-development-for-release-tag-management/01-release-tags-audit-and-deletion-engine-spec.md`
- **Parent Plan Reference:** `.ai-memory/plans/command-development-for-release-tag-management.md`
- **Subtask Plan Reference:** `.ai-memory/plans/subtasks/command-development-for-release-tag-management/02-cli-routing-help-ui-and-documentation.md`
- **Target Subsystem:** `cli/cmdfix/`, `cli/cmdfixreleasetags/`, `cli/cmd/`, `cli/constants/`, `cli/helpdoc/`

---

## 1. User Request (Verbatim)

```text
A, add a new command and also the help text for the command. The command should actually remove releases and release tags if that release tag does not have the proper binary or the release published, or at that time, if that has CI/CD, then the release tag will be removed. Can you please work on this as a command and also add that command to the new command list, and also add the example and help text for this, and also add this command information in the UI as well? Please bump the version, minor version, bump it and release it. And also at the same time, check and test it on top of the `gitmap`. The command format should be `gitmap` space fix space release tags. And that would also give a dry run if you want it to, so that it shows what will be the end result if you run it, and that would prompt it. But user can also do a hyphen Y or confirm, confirm to make sure that it works automatically.
```

---

## 2. Executive Summary & Architectural Scope

When maintaining continuous delivery repositories with automated GitHub Actions releases, workflow interruptions or failed build matrices often leave orphaned Git tags or broken GitHub releases. Specifically:
1. **Empty Asset Releases:** A GitHub release is generated, but the binary build fails, leaving zero downloadable binaries.
2. **Failed CI/CD Releases:** A release tag is created and pushed, but the pipeline errors out, leaving a broken tag at remote origin and GitHub.
3. **Draft / Stale Releases:** Partial releases remain stuck in draft or prerelease status without artifacts.

To remediate these issues cleanly, GitMap introduces a dedicated command:
`gitmap fix release tags` (along with top-level shortcuts `gitmap fix-release-tags` and `gitmap frt`).

This specification defines the complete presentation, routing, and user interaction layer:
- **Multi-Word Routing:** Seamlessly parse `gitmap fix release tags` in `cli/cmdfix/fix_cmd.go` while simultaneously supporting top-level invocations.
- **Flags Contract:** Complete flag support for `--dry-run` (`-n`), non-interactive bypass `-y` / `--yes` / `--confirm`, structured output `--json`, and scope scoping `--local-only` / `--remote-only`.
- **Terminal UI (TUI):** A structured 6-column preview table rendered through `termout.PrintTable`, styled two-column help menus via `termout.HelpMenu`, and a protected interactive confirmation prompt guarded by `isInteractiveStdin()`.
- **Command Discovery & Docs:** System-wide registration in `cli/constants/constants_cli.go`, root usage groups, Muse 5-cluster help menus, and comprehensive embedded documentation in `cli/helpdoc/fix-release-tags.md`.

---

## 3. CLI Routing Architecture & Multi-Word Parsing

### 3.1 Subcommand Dispatching Matrix

```
User Input Command Invocations
  ├── gitmap fix release tags [flags] ──┐
  ├── gitmap fix release-tags [flags] ──┼─► cli/cmdfix/fix_cmd.go:RunFix()
  ├── gitmap fix releasetags [flags]  ──┘       │
  │                                             ▼
  │                                     isFixReleaseTagsRequest()
  │                                             │ (extract flags, strip tokens)
  │                                             ▼
  ├── gitmap fix-release-tags [flags] ──┐   cli/cmdfixreleasetags.RunCLI(flags)
  └── gitmap frt [flags]              ──┴─► cli/cmd/rootrelease.go:dispatchRelease()
```

### 3.2 Multi-Word Parsing Logic in `cli/cmdfix/fix_cmd.go`

In standard Go CLI routing, spaces separate positional arguments. When a user runs `gitmap fix release tags`, `RunFix(args []string, aliasOverride string)` receives `args = ["release", "tags", ...]`.

To ensure smooth dispatching without conflicting with other `fix` targets (such as `fix repo-name stash`), `fix_cmd.go` introduces token sequence detection:

```go
// isFixReleaseTagsRequest inspects arguments for multi-word or hyphenated release tags fix requests.
// Returns a boolean match indicator and the remaining argument slice stripped of command tokens.
func isFixReleaseTagsRequest(args []string) (bool, []string) {
	if len(args) == 0 {
		return false, nil
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))

	// Case 1: Two-word positional sequence: "fix release tags"
	if first == "release" && len(args) >= 2 {
		second := strings.ToLower(strings.TrimSpace(args[1]))
		if second == "tags" || second == "tag" {
			remaining := make([]string, 0, len(args)-2)
			remaining = append(remaining, args[2:]...)
			return true, remaining
		}
	}

	// Case 2: Hyphenated / concatenated subcommands: "fix release-tags", "fix releasetags"
	if first == "release-tags" || first == "releasetags" || first == "release_tags" {
		remaining := make([]string, 0, len(args)-1)
		remaining = append(remaining, args[1:]...)
		return true, remaining
	}

	return false, nil
}
```

### 3.3 Integration into `RunFix` Dispatch Loop

In `cli/cmdfix/fix_cmd.go`, `RunFix` tests for the release tags request immediately after standard help and AGY checks:

```go
func RunFix(args []string, aliasOverride string) error {
	cmdName := "fix"
	if aliasOverride != "" {
		cmdName = aliasOverride
	}
	checkHelp(cmdName, args)

	if isFixAgyRequest(args) {
		return cmdagy.RunPipelineFixAgyCLI(args)
	}

	// Route "gitmap fix release tags" / "gitmap fix release-tags"
	if isMatch, subArgs := isFixReleaseTagsRequest(args); isMatch {
		return cmdfixreleasetags.RunCLI(subArgs)
	}

	if isFixIgnoreRequest(args) {
		return runFixIgnoreDispatch(args)
	}

	if isFixLsRequest(args) {
		return runFixLs(args)
	}
    // ... remaining flow
}
```

### 3.4 Top-Level Command Dispatching

In addition to the nested `fix release tags` syntax, top-level aliases are registered in `cli/cmd/rootrelease.go`:

```go
// In cli/cmd/rootrelease.go -> releaseDispatchEntries()
{
	[]string{
		constants.CmdFixReleaseTags,
		constants.CmdFixReleaseTagsAlias,
		constants.CmdFixReleaseTagsSolid,
	},
	func() error {
		return cmdfixreleasetags.RunCLI(argsTail())
	},
},
```

---

## 4. Flags Specification & Domain Models

### 4.1 Flag Definition Matrix

| Flag Name | Short Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--dry-run` | `-n` | `bool` | `false` | Audit and preview orphan/broken tags and releases without deleting anything. |
| `--yes`, `--confirm` | `-y` | `bool` | `false` | Automatically confirm deletion prompts; bypasses interactive terminal stdin. |
| `--json` | - | `bool` | `false` | Output findings and remediation results as structured JSON (suppresses TUI tables). |
| `--local-only` | - | `bool` | `false` | Restrict scanning and deletion to local Git repository tags (`refs/tags/*`). |
| `--remote-only` | - | `bool` | `false` | Restrict scanning and deletion to GitHub releases and remote tracking tags (`origin`). |
| `--repo` | `-r` | `string` | `.` | Target repository directory or repository slug (defaults to current working directory). |
| `--verbose` | `-v` | `bool` | `false` | Emit detailed telemetry, API response payloads, and sub-process execution traces. |
| `--help` | `-h` | `bool` | `false` | Display the styled two-column command help menu. |

### 4.2 Positive Boolean Domain Model

In compliance with repository coding standards, all boolean fields are expressed with positive polarity:

```go
// Package cmdfixreleasetags models for command execution flags.
package cmdfixreleasetags

// FixReleaseTagsFlags defines parsed CLI options.
type FixReleaseTagsFlags struct {
	IsDryRun        bool   `json:"isDryRun"`
	IsConfirmed     bool   `json:"isConfirmed"`
	IsJSON          bool   `json:"isJSON"`
	IsLocalOnly     bool   `json:"isLocalOnly"`
	IsRemoteOnly    bool   `json:"isRemoteOnly"`
	IsVerbose       bool   `json:"isVerbose"`
	IsHelpRequested bool   `json:"isHelpRequested"`
	TargetDirectory string `json:"targetDirectory"`
}
```

### 4.3 Validation Rules & Error Conditions

1. **Mutual Exclusion:** If both `--local-only` and `--remote-only` are passed simultaneously:
   - Exit with user error (`E1024`): `cannot specify both --local-only and --remote-only flags simultaneously`.
2. **Directory Validation:** If `--repo` points to a non-existent or inaccessible folder:
   - Exit with error: `specified repository directory does not exist or is inaccessible: <path>`.
3. **Git Repository Verification:** If the target directory does not contain a `.git` root:
   - Exit with error: `directory is not a git repository: <path>`.

---

## 5. Terminal User Interface (TUI) Design

### 5.1 Diagnostic Table Preview via `termout.PrintTable`

Before any destructive deletion occurs, or whenever `--dry-run` is active, GitMap prints a structured preview table using `termout.PrintTable`.

#### 5.1.1 Column Specifications

| Column Index | Column Header | Min Width | Max Width | Alignment | Header / Cell Styling |
| :---: | :--- | :---: | :---: | :--- | :--- |
| **1** | `TAG` | 10 | 18 | `AlignLeft` | `constants.ColorCyan` (e.g. `v1.2.3`, `v6.529.0`) |
| **2** | `COMMIT` | 8 | 10 | `AlignLeft` | `constants.ColorDim` (e.g. `a1b2c3d`) |
| **3** | `RELEASE STATUS` | 16 | 22 | `AlignLeft` | Color-coded: Green `Published`, Yellow `Draft`, Red `Missing Release` |
| **4** | `ASSETS` | 10 | 16 | `AlignLeft` | Color-coded: Red `0 Assets`, Green `4 Assets (42MB)` |
| **5** | `CI/CD STATUS` | 14 | 18 | `AlignLeft` | Color-coded: Red `Failed`, Yellow `In-Progress`, Green `Passed`, Dim `No CI` |
| **6** | `ACTION` | 20 | 28 | `AlignLeft` | Red `Delete Release & Tag`, Yellow `Skip (Protected)`, Dim `Keep` |

#### 5.1.2 Go Table Configuration Constructor

```go
func buildDiagnosticTableConfig(rows []ReleaseTagAuditRow) termout.TableConfig {
	columns := []termout.Column{
		{Title: "TAG", MinWidth: 10, MaxWidth: 18, Align: termout.AlignLeft},
		{Title: "COMMIT", MinWidth: 8, MaxWidth: 10, Align: termout.AlignLeft},
		{Title: "RELEASE STATUS", MinWidth: 16, MaxWidth: 22, Align: termout.AlignLeft},
		{Title: "ASSETS", MinWidth: 10, MaxWidth: 16, Align: termout.AlignLeft},
		{Title: "CI/CD STATUS", MinWidth: 14, MaxWidth: 18, Align: termout.AlignLeft},
		{Title: "ACTION", MinWidth: 20, MaxWidth: 28, Align: termout.AlignLeft},
	}

	tableRows := make([]termout.Row, 0, len(rows))
	for _, r := range rows {
		tableRows = append(tableRows, termout.Row{
			Cells: []string{
				r.TagName,
				r.ShortCommit,
				formatReleaseStatusCell(r.ReleaseStatus),
				formatAssetsCell(r.AssetCount, r.TotalSizeBytes),
				formatCICDStatusCell(r.CICDStatus),
				formatActionCell(r.ActionDecision),
			},
		})
	}

	return termout.TableConfig{
		Columns:      columns,
		Rows:         tableRows,
		HeaderColor:  constants.ColorCyan,
		BorderColor:  constants.ColorDim,
		HasBorders:   true,
		EllipsisText: "...",
	}
}
```

#### 5.1.3 Visual Terminal Preview

```text
  [RELEASE-TAG-AUDIT] Inspecting tags for owner/gitmap
  Found 3 candidate tag(s) requiring remediation:

  ┌────────────┬──────────┬──────────────────┬─────────────┬──────────────┬──────────────────────────┐
  │ TAG        │ COMMIT   │ RELEASE STATUS   │ ASSETS      │ CI/CD STATUS │ ACTION                   │
  ├────────────┼──────────┼──────────────────┼─────────────┼──────────────┼──────────────────────────┤
  │ v6.528.1   │ 7e89ab1  │ Missing Release  │ 0 Assets    │ Failed       │ Delete Release & Tag     │
  │ v6.527.0   │ 3c4d5e6  │ Published (Stale)│ 0 Assets    │ Failed       │ Delete Release & Tag     │
  │ v6.526.0-d │ 9f0a1b2  │ Draft            │ 0 Assets    │ In-Progress  │ Skip (Grace Period)      │
  │ v6.525.0   │ 1122334  │ Published        │ 3 (45 MB)   │ Passed       │ Keep (Healthy Release)   │
  └────────────┴──────────┴──────────────────┴─────────────┴──────────────┴──────────────────────────┘
```

---

### 5.2 Two-Column Help Menu via `termout.HelpMenu`

When invoked with `-h` or `--help`, the command displays a two-column menu created via `termout.HelpMenu`:

```go
// BuildFixReleaseTagsHelpMenu constructs the command help menu.
func BuildFixReleaseTagsHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Release Tag & Broken Release Cleanup Engine (gitmap fix release tags)",
		UsageLines: []string{
			"gitmap fix release tags [flags]",
			"gitmap fix-release-tags [flags]",
			"gitmap frt [flags]",
		},
		Sections: []termout.HelpSection{
			buildAuditAndScopeSection(),
			buildExecutionControlSection(),
			buildOutputFormattingSection(),
		},
		FooterFlags: buildHelpFooterFlags(),
		Tips: []string{
			"Always run with '--dry-run' first to preview orphan tags without deleting.",
			"Use '-y' or '--confirm' in automated GitHub Actions / CI cleanup workflows.",
			"Tags with passing CI/CD or published binary assets are automatically protected.",
		},
	}
}

func buildAuditAndScopeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Audit Scope & Targets",
		Color: constants.ColorCyan,
		Entries: []termout.CommandEntry{
			{Command: "--repo, -r <path>", Description: "Target repository path or slug (default: current directory)"},
			{Command: "--local-only", Description: "Audit and delete local git tags only (skip remote/GitHub)"},
			{Command: "--remote-only", Description: "Audit and delete remote origin tags and GitHub releases only"},
		},
	}
}

func buildExecutionControlSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Execution & Safety Controls",
		Color: constants.ColorGreen,
		Entries: []termout.CommandEntry{
			{Command: "-n, --dry-run", Description: "Simulate audit and preview actions without modifying repository"},
			{Command: "-y, --yes, --confirm", Description: "Bypass interactive confirmation prompt and execute deletions"},
		},
	}
}

func buildOutputFormattingSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Output & Formatting",
		Color: constants.ColorYellow,
		Entries: []termout.CommandEntry{
			{Command: "--json", Description: "Emit audit report and deletion results as structured JSON"},
			{Command: "-v, --verbose", Description: "Display detailed API HTTP payloads and git command execution"},
		},
	}
}

func buildHelpFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-h, --help", Description: "Show this help menu and exit"},
	}
}
```

---

### 5.3 Interactive Confirmation Prompt & `isInteractiveStdin()` Guard

#### 5.3.1 Safety Guard Specification

To prevent scripts, pipes, and CI agents from hanging indefinitely on `os.Stdin.Read`, confirmation logic follows a multi-tier guard:

```go
func confirmDeletion(candidateCount int, isConfirmed bool) (bool, error) {
	// Tier 1: Explicit flag bypass (-y / --confirm)
	if isConfirmed {
		return true, nil
	}

	// Tier 2: Non-interactive environment check
	if !isInteractiveStdin() {
		return false, apperror.NewSimple(
			"Aborting: non-interactive environment detected without --confirm (-y). Pass --dry-run to preview or -y to proceed.",
			"E1025",
		)
	}

	// Tier 3: Interactive terminal prompt
	fmt.Printf("\n%sDelete %d orphan/broken release(s) and tag(s)? [y/N]: %s",
		constants.ColorYellow, candidateCount, constants.ColorReset)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, apperror.WrapSimple(err, "read confirmation response")
	}

	trimmed := strings.ToLower(strings.TrimSpace(input))
	if trimmed == "y" || trimmed == "yes" {
		return true, nil
	}

	return false, nil
}
```

#### 5.3.2 Non-Interactive Detection Standard

`isInteractiveStdin()` verifies:
1. Environment variables: `CI != ""` or `GITMAP_NON_INTERACTIVE == "1"`.
2. Terminal file descriptor: `(fi.Mode() & os.ModeCharDevice) != 0`.

---

## 6. Command Discovery, Help Clusters & Documentation

### 6.1 CLI Constants Registration (`cli/constants/constants_cli.go`)

The following constants are added to `cli/constants/constants_cli.go`:

```go
// Release tag cleanup commands
const (
	CmdFixReleaseTags      = "fix-release-tags"
	CmdFixReleaseTagsAlias = "frt"
	CmdFixReleaseTagsSolid = "fixreleasetags"
)

// Help string definition
const (
	HelpFixReleaseTags = "  fix-release-tags (frt) Delete orphan or broken GitHub releases and tags lacking assets or with failed CI/CD"
)
```

### 6.2 Root Usage Groups Integration (`cli/cmd/rootusage_groups.go`)

`HelpFixReleaseTags` is appended to the Release Information group:

```go
func printGroupReleaseInfo() {
	renderHeader(constants.HelpGroupReleaseInfo)
	renderLine(constants.HelpChangelog)
	renderLine(constants.HelpChangelogGen)
	renderLine(constants.HelpListVersions)
	renderLine(constants.HelpListReleases)
	renderLine(constants.HelpReleasePend)
	renderLine(constants.HelpFixReleaseTags) // Added
	renderLine(constants.HelpRevert)
	renderLine(constants.HelpClearReleaseJSON)
	renderLine(constants.HelpPrune)
}
```

### 6.3 Semantic Help Clusters Integration (`cli/cmd/roothelp_clusters.go`)

In `resolveClusterEntriesRelease()`, the command is listed as an official capability:

```go
func resolveClusterEntriesRelease() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "release, r", Description: "Automated semantic release ceremony, tagging, and publishing"},
		{Command: "fix-release-tags, frt", Description: "Audit and clean orphan or broken GitHub releases and tags lacking assets"},
		{Command: "changelog", Description: "Generate or view repository changelog from commits"},
		{Command: "list-versions", Description: "Display repository version history and release metadata"},
		{Command: "commit, co", Description: "Structured conventional commit runner with validation"},
		{Command: "cpf", Description: "Commit, push, and create feature branch in one pass"},
		{Command: "cpb", Description: "Commit and push bugfix branch with automated telemetry"},
		{Command: "cpr", Description: "Commit, push, and trigger release workflow orchestration"},
		{Command: "cpar", Description: "Commit and push across all connected repository workspaces"},
	}
}
```

---

## 7. Embedded Markdown Documentation Specification

### 7.1 Target File Location
- **Path:** `cli/helpdoc/fix-release-tags.md`

### 7.2 Document Structure & Content

```markdown
# gitmap fix release tags

Audit, identify, and remove orphan or broken GitHub releases and Git release tags that lack binary assets, have unpublished release notes, or encountered failed CI/CD pipelines.

## Aliases

- `gitmap fix-release-tags`
- `gitmap frt`
- `gitmap fix release-tags`
- `gitmap fix releasetags`

## Usage

```bash
gitmap fix release tags [flags]
gitmap fix-release-tags [flags]
gitmap frt [flags]
```

## Description

Automated continuous integration and release workflows occasionally create a Git tag before triggering the binary compilation and release publishing steps. If the build runner runs out of disk space, encounters compilation errors, or network limits prevent artifact uploads, the repository is left with:
- An orphan Git tag pointing to an unreleased commit.
- A GitHub release with zero downloadable binary assets.
- A failed CI/CD workflow run attached to the release tag.

`gitmap fix release tags` audits all tags against GitHub Releases API and GitHub Actions workflow runs. If a tag is identified as broken or orphaned, the command securely deletes:
1. The GitHub Release entity (via `gh release delete` or REST API).
2. The remote Git tag on `origin` (`git push origin :refs/tags/<tag>`).
3. The local Git tag (`git tag -d <tag>`).

## Options & Flags

| Flag | Shorthand | Description | Default |
| :--- | :---: | :--- | :--- |
| `--dry-run` | `-n` | Preview detected broken releases and tags without deleting. | `false` |
| `--yes`, `--confirm` | `-y` | Non-interactive execution; bypasses the `[y/N]` confirmation prompt. | `false` |
| `--json` | - | Emit machine-readable JSON output (disables ASCII tables). | `false` |
| `--local-only` | - | Audit and remove only local Git tags; skips remote/GitHub operations. | `false` |
| `--remote-only` | - | Audit and remove only remote origin tags and GitHub releases. | `false` |
| `--repo` | `-r` | Target repository directory or repository name. | `.` |
| `--verbose` | `-v` | Display raw GitHub API responses and Git command outputs. | `false` |
| `--help` | `-h` | Display the two-column interactive help card. | `false` |

## Safety Invariants

The audit engine enforces three strict safety gates:
1. **Healthy Release Immunity:** Tags that possess at least one valid binary asset and a passing/neutral CI/CD status are never deleted.
2. **Grace Window:** Tags created within the last 30 minutes with active "in-progress" or "queued" CI/CD runs are skipped.
3. **Active Binary Protection:** If the current running `gitmap` binary matches the tag version, deletion is blocked.

## Examples

### Example 1: Dry-Run Inspection

Preview all orphan release tags in the current repository without making modifications:

```bash
gitmap fix release tags --dry-run
```

### Example 2: Interactive Remediation

Review candidate tags in a table and confirm deletion when prompted:

```bash
gitmap fix release tags
```

### Example 3: Non-Interactive CI Automation

Clean up failed release tags in CI/CD runners or headless batch scripts:

```bash
gitmap fix-release-tags --confirm
```

### Example 4: JSON Output for External Tooling

Generate structured JSON telemetry for pipeline integrations:

```bash
gitmap frt --dry-run --json
```

## JSON Schema Output Example

```json
{
  "repository": "alimtvnetwork/gitmap-v28",
  "auditTimestamp": "2026-10-10T06:30:00Z",
  "isDryRun": true,
  "candidates": [
    {
      "tag": "v6.528.1",
      "commit": "7e89ab1",
      "releaseStatus": "missing",
      "assetCount": 0,
      "cicdStatus": "failure",
      "action": "delete",
      "reasons": [
        "Release entity missing on GitHub",
        "CI/CD workflow run 1928374 failed"
      ]
    }
  ],
  "summary": {
    "totalAudited": 142,
    "brokenFound": 1,
    "deletedCount": 0
  }
}
```

## See Also

- `gitmap release` — Perform semantic release ceremony and push tags.
- `gitmap list-versions` — Inspect version tags and published metadata.
- `gitmap fix-git` — Self-heal corrupted Git index, permissions, and lockfiles.
```

---

## 8. JSON Data Contract Specification

When `--json` is supplied, standard text and TUI tables are suppressed. The command emits a JSON document matching this contract:

```go
type FixReleaseTagsReport struct {
	Repository      string                      `json:"repository"`
	AuditTimestamp  string                      `json:"auditTimestamp"`
	IsDryRun        bool                        `json:"isDryRun"`
	IsLocalOnly     bool                        `json:"isLocalOnly"`
	IsRemoteOnly    bool                        `json:"isRemoteOnly"`
	Candidates      []ReleaseTagCandidateReport `json:"candidates"`
	Summary         FixReleaseTagsSummary       `json:"summary"`
	ExecutionErrors []string                    `json:"executionErrors,omitempty"`
}

type ReleaseTagCandidateReport struct {
	Tag           string   `json:"tag"`
	CommitSHA     string   `json:"commit"`
	ReleaseStatus string   `json:"releaseStatus"`
	AssetCount    int      `json:"assetCount"`
	TotalBytes    int64    `json:"totalBytes"`
	CICDStatus    string   `json:"cicdStatus"`
	Action        string   `json:"action"`
	Reasons       []string `json:"reasons"`
	IsDeleted     bool     `json:"isDeleted"`
}

type FixReleaseTagsSummary struct {
	TotalAudited int `json:"totalAudited"`
	BrokenFound  int `json:"brokenFound"`
	DeletedCount int `json:"deletedCount"`
	SkippedCount int `json:"skippedCount"`
}
```

---

## 9. Acceptance Criteria & Quality Gates

| ID | Criterion | Verification Method |
| :--- | :--- | :--- |
| **AC-01** | `gitmap fix release tags` routes to `cmdfixreleasetags.RunCLI` with flags preserved. | Unit test with argument combinations in `cli/cmdfix/fix_cmd_test.go`. |
| **AC-02** | `gitmap fix-release-tags` and `gitmap frt` route identically from top-level CLI. | Route execution test in `cli/cmd/rootrelease_test.go`. |
| **AC-03** | `--dry-run` and `-n` perform audit, print table, and execute zero deletion calls. | Verification of mock execution calls in `cmdfixreleasetags/cli_test.go`. |
| **AC-04** | Confirmation prompt displays formatted count and handles `y`/`yes` cleanly. | Interactive IO unit test in `cmdfixreleasetags/prompt_test.go`. |
| **AC-05** | Non-interactive stdin aborts safely with `E1025` unless `--confirm`/`-y` is passed. | Pipe execution test in headless environment. |
| **AC-06** | Preview table renders exactly the 6 required columns via `termout.PrintTable`. | Golden snapshot test in `cmdfixreleasetags/ui_table_test.go`. |
| **AC-07** | `--json` emits valid JSON conforming to `FixReleaseTagsReport`. | JSON unmarshal verification test in CI. |
| **AC-08** | Two-column help menu renders on `-h` / `--help`. | Output string matching in `cmdfixreleasetags/ui_help_test.go`. |
| **AC-09** | Strictly relative paths used across all specifications and implementation files. | Relative path audit script pass (`check-spec-cross-links.py`). |
| **AC-10** | All boolean fields and variables use positive polarity naming. | Code style and guideline verification. |
