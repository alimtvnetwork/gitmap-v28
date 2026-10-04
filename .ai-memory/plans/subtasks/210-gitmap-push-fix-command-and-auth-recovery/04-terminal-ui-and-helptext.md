# Subtask Plan 04: Terminal UI, Status Badges, Progress Cards & Helptext System

> **Task Reference:** `210-gitmap-push-fix-command-and-auth-recovery`  
> **Subtask ID:** `04-terminal-ui-and-helptext`  
> **Status:** Pending Execution  
> **Assigned Agent Role:** Terminal UI & Documentation Specialist  

---

## 1. Overview & Objectives

Implement the terminal user interface, status badge formatters, diagnostic breakdown cards, step-by-step progress cards, and embedded Markdown help documentation for `gitmap push-fix`.

### Key Deliverables:
1. Create `cli/helptext/push_fix.md` containing rich command documentation, flags, prerequisites, examples, and JSON help metadata.
2. Register the topic in `cli/helptext/catalog.go` (`topicSummaries`) and register aliases in `cli/helptext/print.go` (`helpAliases`).
3. Implement `RenderPushFixHelp()` leveraging `cli/termhelp` (`HelpMenu`, `PrintBanner`, `PrintUsage`, `PrintSection`, `CalculateMaxCommandWidth`).
4. Implement Status Badges supporting both `glyphs.ModeRich` (UTF-8 color badges) and `glyphs.ModeSafe` (ASCII bracket tags) in `cli/cmdpull/push_fix_ui.go`.
5. Implement Error Breakdown Table & Diagnostic Cards displaying issue category, root cause, and remediation plan.
6. Implement 4-Stage Progress Cards (`RenderStepCard`) tracking pipeline progress from inspection to remote push verification.
7. Author unit tests ensuring helptext embedding and UI rendering operate cleanly across all terminal modes.

---

## 2. Step-by-Step Implementation Steps

### Step 4.1: Author Embedded Helptext Documentation
1. **Create `cli/helptext/push_fix.md`**:
   - Provide complete markdown documentation with standard headers:
     - `# gitmap push-fix`
     - Description of autonomous push failure diagnosis and auth recovery.
     - `## Aliases`: `pf`, `pushfix`, `fix-push`, `push fix`.
     - `## Usage`:
       ```
       gitmap push-fix [flags]
       gitmap pf [flags]
       gitmap push fix [flags]
       ```
     - `## Flags`: Table documenting `--dry-run`, `--remote`, `--branch`, `--force`, `--yes`, `--ssh`, `--https`.
     - `## Prerequisites`: Working directory must be a Git repository.
     - `## Examples`:
       - Autonomous push fix.
       - Dry run preview.
       - SSH conversion and push.
       - Target custom remote and branch with `--force-with-lease`.
     - `## Scripting (JSON)`: Instructions for `gitmap help --json --filter push-fix`.

### Step 4.2: Register Topic in Catalog and Print Aliases
1. **Edit `cli/helptext/catalog.go`**:
   - In `topicSummaries`, register:
     ```go
     "push-fix": "Diagnose failed git pushes, repair authentication and transport protocols, and safely complete remote updates.",
     "pf":       "Diagnose failed git pushes, repair authentication and transport protocols, and safely complete remote updates.",
     "pushfix":  "Diagnose failed git pushes, repair authentication and transport protocols, and safely complete remote updates.",
     "fix-push": "Diagnose failed git pushes, repair authentication and transport protocols, and safely complete remote updates.",
     ```
2. **Edit `cli/helptext/print.go`**:
   - In `helpAliases`, register:
     ```go
     "pf":       "push_fix",
     "pushfix":  "push_fix",
     "fix-push": "push_fix",
     "push-fix": "push_fix",
     ```

### Step 4.3: Implement Rich Box Menu via `cli/termhelp`
1. **Create `cli/cmdpull/push_fix_help.go`**:
   - Define `BuildPushFixHelpMenu() termhelp.HelpMenu`:
     ```go
     package cmdpull

     import (
         "github.com/alimtvnetwork/gitmap-v28/cli/constants"
         "github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
     )

     func BuildPushFixHelpMenu() termhelp.HelpMenu {
         return termhelp.HelpMenu{
             Title: "GITMAP PUSH-FIX & AUTH RECOVERY",
             UsageLines: []string{
                 "gitmap push-fix [flags]",
                 "gitmap pf [flags]",
                 "gitmap push fix [flags]",
             },
             Sections: []termhelp.HelpSection{
                 {
                     Title: "Autonomous Remediation",
                     Color: constants.ColorGreen,
                     Entries: []termhelp.CommandEntry{
                         {Command: "gitmap pf", Description: "Diagnose and auto-repair push failures"},
                         {Command: "gitmap pf --dry-run", Description: "Simulate recovery actions safely"},
                         {Command: "gitmap pf --ssh", Description: "Convert remote to SSH and push"},
                         {Command: "gitmap pf --https", Description: "Convert remote to HTTPS and push"},
                         {Command: "gitmap pf --force", Description: "Safe push with --force-with-lease"},
                     },
                 },
                 {
                     Title: "Options",
                     Color: constants.ColorCyan,
                     Entries: []termhelp.CommandEntry{
                         {Command: "-n, --dry-run", Description: "Preview fixes without making changes"},
                         {Command: "-r, --remote <name>", Description: "Target remote (default: origin)"},
                         {Command: "-b, --branch <name>", Description: "Target branch (default: current)"},
                         {Command: "-f, --force", Description: "Push with --force-with-lease"},
                         {Command: "-y, --yes", Description: "Non-interactive bypass for prompts"},
                     },
                 },
             },
             Tips: []string{
                 "Run 'gitmap pf -n' to safely preview the repair plan before modifying remotes.",
                 "Use 'gitmap push fix' as an in-place muscle-memory alternative.",
             },
         }
     }

     func RenderPushFixHelp() {
         termhelp.RenderMenu(BuildPushFixHelpMenu())
     }
     ```

### Step 4.4: Implement Status Badges & Diagnostic Cards
1. **Create `cli/cmdpull/push_fix_ui.go`**:
   - Implement badge rendering supporting `glyphs.ModeRich` and `glyphs.ModeSafe`:
     ```go
     package cmdpull

     import (
         "fmt"
         "github.com/alimtvnetwork/gitmap-v28/cli/constants"
         "github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
     )

     type PushFixBadge int

     const (
         BadgeDiagnosing PushFixBadge = iota
         BadgeAuthError
         BadgeDiverged
         BadgeProtocol
         BadgeForceLease
         BadgeDryRun
         BadgeRecovered
     )

     func FormatBadge(badge PushFixBadge) string {
         isRich := glyphs.Resolve() == glyphs.ModeRich
         switch badge {
         case BadgeDiagnosing:
             if isRich {
                 return fmt.Sprintf("%s[ 🔍 DIAGNOSING ]%s", constants.ColorCyan, constants.ColorReset)
             }
             return "[DIAGNOSE]"
         case BadgeAuthError:
             if isRich {
                 return fmt.Sprintf("%s[ ✗ AUTH ERROR ]%s", constants.ColorRed, constants.ColorReset)
             }
             return "[AUTH-FAIL]"
         case BadgeDiverged:
             if isRich {
                 return fmt.Sprintf("%s[ ⚠ DIVERGED ]%s", constants.ColorYellow, constants.ColorReset)
             }
             return "[NON-FF]"
         case BadgeProtocol:
             if isRich {
                 return fmt.Sprintf("%s[ ⚡ PROTOCOL ]%s", constants.ColorMagenta, constants.ColorReset)
             }
             return "[TRANSPORT]"
         case BadgeForceLease:
             if isRich {
                 return fmt.Sprintf("%s[ 🛡️ FORCE-LEASE ]%s", constants.ColorYellow, constants.ColorReset)
             }
             return "[FORCE-LEASE]"
         case BadgeDryRun:
             if isRich {
                 return fmt.Sprintf("%s[ ℹ DRY-RUN ]%s", constants.ColorBlue, constants.ColorReset)
             }
             return "[DRY-RUN]"
         case BadgeRecovered:
             if isRich {
                 return fmt.Sprintf("%s[ ✓ RECOVERED ]%s", constants.ColorGreen, constants.ColorReset)
             }
             return "[SUCCESS]"
         default:
             return "[STATUS]"
         }
     }
     ```

### Step 4.5: Implement Error Breakdown & Diagnostic Card Renderer
1. **In `cli/cmdpull/push_fix_ui.go`**, add `RenderDiagnosticCard`:
   ```go
   type DiagnosticReport struct {
       Category      string
       RemoteTarget  string
       ActiveBranch  string
       RootCause     string
       PlannedAction string
   }

   func RenderDiagnosticCard(report DiagnosticReport) {
       fmt.Println()
       fmt.Println("  ┌── Push Failure Diagnosis ──────────────────────────────────────────┐")
       fmt.Printf("  │ Issue Category : %-48s │\n", truncateText(report.Category, 48))
       fmt.Printf("  │ Remote Target  : %-48s │\n", truncateText(report.RemoteTarget, 48))
       fmt.Printf("  │ Active Branch  : %-48s │\n", truncateText(report.ActiveBranch, 48))
       fmt.Printf("  │ Root Cause     : %-48s │\n", truncateText(report.RootCause, 48))
       fmt.Printf("  │ Planned Action : %-48s │\n", truncateText(report.PlannedAction, 48))
       fmt.Println("  └────────────────────────────────────────────────────────────────────┘")
       fmt.Println()
   }
   ```

### Step 4.6: Implement 4-Stage Step-by-Step Progress Cards
1. **In `cli/cmdpull/push_fix_ui.go`**, add `RenderStepCard`:
   ```go
   func RenderStepCard(currentStep, totalSteps int, description, status string) {
       statusColor := constants.ColorCyan
       switch status {
       case "DONE", "SUCCESS", "RECOVERED":
           statusColor = constants.ColorGreen
       case "FAILED", "ERROR":
           statusColor = constants.ColorRed
       case "SKIPPED":
           statusColor = constants.ColorYellow
       }
       fmt.Printf("  [%d/%d] %-52s %s%s%s\n", currentStep, totalSteps, description, statusColor, status, constants.ColorReset)
   }
   ```

### Step 4.7: Author Unit Tests for Helptext & UI
1. **Create `cli/cmdpull/push_fix_ui_test.go`**:
   - Test `FormatBadge` output in `ModeRich` and `ModeSafe`.
   - Test `RenderDiagnosticCard` formatting and line length constraints.
   - Test `BuildPushFixHelpMenu` structure and non-empty entries.
2. **Verify embedded markdown in `cli/helptext/coverage_test.go`**:
   - Confirm `cli/helptext/push_fix.md` loads without missing file errors via `helptext.ReadRaw("push_fix")`.

---

## 3. Verification & Quality Commands

```bash
# Verify helptext coverage and embedding
go test -v ./cli/helptext -run "TestCoverage"

# Verify UI rendering functions
go test -v ./cli/cmdpull -run "TestPushFixUI"

# Verify coding style and guidelines
python 03-ai-scripts/05-guideline-autofixer.py cli/helptext cli/cmdpull --check-only
python linter-scripts/check-relative-paths.py
```
