# Subtask Plan: CLI Routing, Help Menus, Terminal UI & Documentation

> **/goal** Implement the CLI routing, flag parsing, terminal UI table, interactive confirmation prompts with non-interactive safety guards, command discovery registration, and embedded markdown documentation for `gitmap fix release tags` (and aliases `fix-release-tags`, `frt`).
> **/learn** Parse multi-word argument sequences in `cmdfix.RunFix`, enforce positive boolean logic throughout flag structures, prevent non-interactive stdin hangs via `isInteractiveStdin()` guards, format 6-column high-contrast tables using `termout.PrintTable`, and integrate the command into the Muse 5 semantic help clusters.

- **Parent Plan:** `.ai-memory/plans/command-development-for-release-tag-management.md`
- **Spec Reference:** `02-spec/21-app/command-development-for-release-tag-management/02-cli-routing-ui-and-docs-spec.md`
- **Subtask ID:** `02-cli-routing-help-ui-and-documentation`
- **Status:** `QUEUED`
- **Subsystem Areas:** `cli/cmdfix/`, `cli/cmdfixreleasetags/`, `cli/cmd/`, `cli/constants/`, `cli/helpdoc/`

---

## 1. Scope & Responsibility Breakdown

This subtask implements the user-facing CLI interface and presentation layer for orphan and broken release tag remediation:
1. **Multi-Word Routing:** Intercept `gitmap fix release tags` (and `fix release-tags`, `fix releasetags`) in `cli/cmdfix/fix_cmd.go` and route top-level `gitmap fix-release-tags` and `gitmap frt` in `cli/cmd/rootrelease.go`.
2. **Flag Parser & Positive Booleans:** Support `--dry-run` (`-n`), `-y` / `--yes` / `--confirm`, `--json`, `--local-only`, `--remote-only`, and `--repo`.
3. **Diagnostic Table Preview:** Render a 6-column preview table (`TAG`, `COMMIT`, `RELEASE STATUS`, `ASSETS`, `CI/CD STATUS`, `ACTION`) using `termout.PrintTable`.
4. **Interactive Confirmation Safety:** Prompt `Delete %d orphan/broken release(s) and tag(s)? [y/N]: ` with `isInteractiveStdin()` guard protecting headless/CI environments from blocking.
5. **Two-Column Help Menu:** Implement `BuildFixReleaseTagsHelpMenu` rendered via `termout.HelpMenu`.
6. **Command Discovery & Docs:** Register constants in `cli/constants/constants_cli.go`, add help entries to `cli/cmd/rootusage_groups.go` and `cli/cmd/roothelp_clusters.go`, and author embedded markdown in `cli/helpdoc/fix-release-tags.md`.

---

## 2. Target File Modification Map

| Target File | Action | Symbol / Function Additions |
| :--- | :--- | :--- |
| `cli/constants/constants_cli.go` | Modify | `CmdFixReleaseTags`, `CmdFixReleaseTagsAlias`, `CmdFixReleaseTagsSolid`, `HelpFixReleaseTags` |
| `cli/cmdfix/fix_cmd.go` | Modify | `isFixReleaseTagsRequest(args []string) (bool, []string)`, integrate routing in `RunFix` |
| `cli/cmd/rootrelease.go` | Modify | Add `CmdFixReleaseTags` and `CmdFixReleaseTagsAlias` to `releaseDispatchEntries()` |
| `cli/cmd/rootusage_groups.go` | Modify | Add `renderLine(constants.HelpFixReleaseTags)` to `printGroupReleaseInfo()` |
| `cli/cmd/roothelp_clusters.go` | Modify | Add `CommandEntry` for `fix-release-tags, frt` to `resolveClusterEntriesRelease()` |
| `cli/cmdfixreleasetags/flags.go` | Create | `FixReleaseTagsFlags` struct, `parseFlags(args []string) (FixReleaseTagsFlags, error)` |
| `cli/cmdfixreleasetags/ui_table.go` | Create | `renderDiagnosticTable(rows []ReleaseTagAuditRow)`, `buildDiagnosticTableConfig()` |
| `cli/cmdfixreleasetags/ui_help.go` | Create | `BuildFixReleaseTagsHelpMenu() termout.HelpMenu`, `RenderFixReleaseTagsHelp()` |
| `cli/cmdfixreleasetags/prompt.go` | Create | `confirmDeletion(candidateCount int, isConfirmed bool) (bool, error)`, `isInteractiveStdin()` |
| `cli/cmdfixreleasetags/cli.go` | Create | `RunCLI(args []string) error` top-level orchestrator |
| `cli/helpdoc/fix-release-tags.md` | Create | Embedded documentation markdown for rich terminal help |
| `cli/cmdfix/fix_cmd_test.go` | Modify | Add tests for `isFixReleaseTagsRequest` and multi-word token parsing |
| `cli/cmdfixreleasetags/flags_test.go` | Create | Unit tests verifying flag permutations, defaults, and mutual exclusion |
| `cli/cmdfixreleasetags/prompt_test.go` | Create | Unit tests verifying non-interactive fail-safe and confirmation bypass |
| `cli/cmdfixreleasetags/ui_table_test.go` | Create | Unit tests verifying table column widths and formatters |

---

## 3. Step-by-Step Implementation Sequence

### Step 1: Constants & Help Registration
- Edit `cli/constants/constants_cli.go`:
  - Define `CmdFixReleaseTags = "fix-release-tags"`.
  - Define `CmdFixReleaseTagsAlias = "frt"`.
  - Define `CmdFixReleaseTagsSolid = "fixreleasetags"`.
  - Define `HelpFixReleaseTags = "  fix-release-tags (frt) Delete orphan/broken GitHub releases & tags lacking assets or with failed CI"`.
- Edit `cli/cmd/rootusage_groups.go`:
  - Append `renderLine(constants.HelpFixReleaseTags)` into `printGroupReleaseInfo()`.
- Edit `cli/cmd/roothelp_clusters.go`:
  - Add `{Command: "fix-release-tags, frt", Description: "Audit and clean orphan or broken GitHub releases and tags lacking assets"}` into `resolveClusterEntriesRelease()`.

### Step 2: Multi-Word CLI Routing in `cli/cmdfix/fix_cmd.go`
- Implement `isFixReleaseTagsRequest(args []string) (bool, []string)`:
  - Check for two-word sequence: `strings.ToLower(args[0]) == "release" && strings.ToLower(args[1]) == "tags"`.
  - Check for hyphenated form: `strings.ToLower(args[0]) == "release-tags"`.
  - Check for concatenated form: `strings.ToLower(args[0]) == "releasetags"`.
  - Return remaining arguments stripped of subcommand tokens.
- Update `RunFix(args []string, aliasOverride string)`:
  - When `isFixReleaseTagsRequest` returns true, invoke `cmdfixreleasetags.RunCLI(subArgs)`.

### Step 3: Top-Level Command Dispatching in `cli/cmd/rootrelease.go`
- Register `constants.CmdFixReleaseTags`, `constants.CmdFixReleaseTagsAlias`, and `constants.CmdFixReleaseTagsSolid` in `releaseDispatchEntries()`.
- Direct invocation to `cmdfixreleasetags.RunCLI(argsTail())`.

### Step 4: Flags Model & Parsing in `cli/cmdfixreleasetags/flags.go`
- Define struct `FixReleaseTagsFlags` with positive booleans:
  - `IsDryRun bool` (`--dry-run`, `-n`)
  - `IsConfirmed bool` (`-y`, `--yes`, `--confirm`)
  - `IsJSON bool` (`--json`)
  - `IsLocalOnly bool` (`--local-only`)
  - `IsRemoteOnly bool` (`--remote-only`)
  - `IsVerbose bool` (`-v`, `--verbose`)
  - `IsHelpRequested bool` (`-h`, `--help`)
  - `TargetDirectory string` (`--repo`, `-r`, default: `.`)
- Implement `parseFlags(args []string) (FixReleaseTagsFlags, error)`:
  - Enforce validation: error if both `IsLocalOnly` and `IsRemoteOnly` are true.
  - Verify `TargetDirectory` exists and contains `.git` (or resolve through Split-DB).

### Step 5: Terminal UI Diagnostic Table in `cli/cmdfixreleasetags/ui_table.go`
- Build table configuration conforming to `termout.TableConfig`:
  - Columns: `TAG`, `COMMIT`, `RELEASE STATUS`, `ASSETS`, `CI/CD STATUS`, `ACTION`.
  - Header styling: `constants.ColorCyan`.
  - Border styling: `constants.ColorDim`.
- Color formatters:
  - `formatReleaseStatusCell`: Green for "Published", Yellow for "Draft", Red for "Missing Release".
  - `formatAssetsCell`: Red for "0 Assets", Green for `N Assets (XX MB)`.
  - `formatCICDStatusCell`: Red for "Failed", Yellow for "In-Progress", Green for "Passed", Dim for "No CI".
  - `formatActionCell`: Red for "Delete Release & Tag", Yellow for "Skip (Protected)", Dim for "Keep".
- Invoke `termout.PrintTable(cfg)` to output formatted table.

### Step 6: Confirmation Guard in `cli/cmdfixreleasetags/prompt.go`
- Implement `confirmDeletion(candidateCount int, isConfirmed bool) (bool, error)`:
  - Return `true, nil` immediately if `isConfirmed` is true.
  - Guard: If `!isInteractiveStdin()`, return error `E1025` ("Aborting: non-interactive environment detected without --confirm (-y)").
  - Display prompt: `Delete %d orphan/broken release(s) and tag(s)? [y/N]: `.
  - Read input line via `bufio.NewReader(os.Stdin)`.
  - Accept "y" / "yes" (case-insensitive). Return `false, nil` for all other inputs.

### Step 7: Help Menu in `cli/cmdfixreleasetags/ui_help.go`
- Implement `BuildFixReleaseTagsHelpMenu() termout.HelpMenu`:
  - Set Title, UsageLines, and 3 help sections:
    1. "Audit Scope & Targets" (`--repo`, `--local-only`, `--remote-only`)
    2. "Execution & Safety Controls" (`-n, --dry-run`, `-y, --yes, --confirm`)
    3. "Output & Formatting" (`--json`, `-v, --verbose`)
  - Add footer flags (`-h, --help`) and usage tips.
- Implement `RenderFixReleaseTagsHelp()` delegating to `termout.RenderMenu()`.

### Step 8: CLI Orchestrator in `cli/cmdfixreleasetags/cli.go`
- Connect components:
  1. Parse flags via `parseFlags(args)`.
  2. If `flags.IsHelpRequested`: display help via `RenderFixReleaseTagsHelp()` and return nil.
  3. Invoke audit engine from Subtask 01 (`engine.AuditReleaseTags(flags)`).
  4. If `flags.IsJSON`: render JSON output via `json.NewEncoder(os.Stdout)` and return.
  5. Print header banner and diagnostic table via `renderDiagnosticTable`.
  6. If candidate count == 0: print "✓ All release tags are healthy and accompanied by valid release assets." and exit cleanly.
  7. If `flags.IsDryRun`: print dry-run summary note and return nil without deletions.
  8. Prompt confirmation via `confirmDeletion(candidateCount, flags.IsConfirmed)`.
  9. If confirmed: invoke deletion engine from Subtask 01 (`engine.DeleteCandidates(candidates, flags)`).

### Step 9: Embedded Help Documentation
- Author `cli/helpdoc/fix-release-tags.md` matching spec §7 with aliases, examples, option descriptions, and JSON contract.

### Step 10: Unit Testing Suite
- `cli/cmdfix/fix_cmd_test.go`:
  - Test `isFixReleaseTagsRequest` with:
    - `["release", "tags"]` -> `(true, [])`
    - `["release", "tags", "--dry-run"]` -> `(true, ["--dry-run"])`
    - `["release-tags", "-y"]` -> `(true, ["-y"])`
    - `["releasetags"]` -> `(true, [])`
    - `["my-repo", "stash"]` -> `(false, nil)`
- `cli/cmdfixreleasetags/flags_test.go`:
  - Test default flag values (all booleans false).
  - Test `--dry-run` and `-n` setting `IsDryRun = true`.
  - Test `-y`, `--yes`, `--confirm` setting `IsConfirmed = true`.
  - Test mutual exclusion of `--local-only` and `--remote-only`.
- `cli/cmdfixreleasetags/prompt_test.go`:
  - Test `isConfirmed` bypass returning true without stdin read.
  - Test interactive stdin mock with "y\n" returning true.
  - Test interactive stdin mock with "n\n" returning false.
- `cli/cmdfixreleasetags/ui_table_test.go`:
  - Test table configuration with 0, 1, and 5 rows.
  - Verify column titles and widths.

---

## 4. Verification & Quality Gates

1. **Relative Path Audit:** Confirm zero absolute paths or `file:///` URIs exist in any source or documentation files.
2. **Positive Booleans Audit:** Verify that all flags, struct fields, and functions use positive polarity naming (`isConfirmed`, `isDryRun`, `isLocalOnly`, `isRemoteOnly`).
3. **CLI Invocation Tests:**
   - `gitmap fix release tags --help`
   - `gitmap fix-release-tags --dry-run`
   - `gitmap frt -n --json`
4. **Non-Interactive CI Safety Test:** Execute `echo "" | gitmap fix release tags` to ensure the process exits cleanly with `E1025` without blocking.
