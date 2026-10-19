# Specification: Screenshot Alignment and Installer Script Adjustment

## Status
- **Task ID**: `screenshot-alignment-and-installer-script-adjustment`
- **Specification Version**: 1.0.0
- **Status**: APPROVED
- **Date**: 2026-10-19

---

## 1. User Request (Verbatim)

```text
# Screenshot Alignment and Installer Script Adjustment: high priority instruction, non-negotiable task

I have given you three screenshots. First of all, the blue color that you see or created needs to be a little bit right-hand side so that it's aligned with the white one at the end. The second is that your numbers, one by five, four by five, need to be a little bit right-hand side to match with the probing section. When you start with the success, have a little bit of gap, like a new line before you start with the success step. Fix the TXT broken issue. Inside the circle, it looks nice but a little bit broken. Look into the circle that is created. A little bit broken, so let's calculate that. How long do you grade this? At the end, when you say, "Added source line," that will give a bunch of same bugs. Say, "Added source line," and then Enter, and then as a subtree, add these paths. It would make more sense and be much clearer. At the end, after you complete the tree, there should be a new line gap before you hand it over to the user or user types. Make sure the same happens for the shell installer script too.
```

---

## 2. Screenshot Forensics & Root Cause Analysis

### 2.1 Spinner Symbol & Stage Step Number Misalignment
- **Screenshot Evidence**: `assets/screenshots/screenshot-alignment-and-installer-script-adjustment-02.png`
- **Existing State**:
  - The probe message in `cli/constants/constants_update.go` is defined as:
    `"\n  ■ Probing sibling repos (%s-v%d..%s-v%d)...\n"`
    which indents the white square `■` by 2 spaces (columns 0-1 are spaces, column 2 is `■`, column 3 is space, column 4 starts `Probing`).
  - In contrast, the `pterm.DefaultSpinner` in `cli/cmdupdate/update_ui.go` rendered the cyan spinner square `■` at column 0 without leading indentation.
  - The stage step tags `[1/5]`, `[4/5]` began at column 2, directly underneath the white `■` instead of aligning with `Probing` at column 4.
- **Root Cause**: `pterm.DefaultSpinner` outputs sequences flush-left by default.
- **Target Resolution**: Indent the spinner sequence by 2 spaces (`  ■ `), placing the cyan square at column 2 (matching the white square) and the stage numbers `[1/5]` at column 4 (matching `Probing`).

### 2.2 `txt` Broken Ghost Text & Stage Ghost Duplication
- **Screenshot Evidence**: `assets/screenshots/screenshot-alignment-and-installer-script-adjustment-01.png`
- **Existing State**:
  - A stray literal string `txt` appears flush-left at column 0, sandwiched between two blocks of `SUCCESS` badges.
  - The first two `SUCCESS` lines are duplicated both above and below `txt`.
- **Root Cause**:
  - `updateUIRendererImpl` used `pterm.DefaultMultiPrinter` with dynamic writer allocation on Windows PowerShell.
  - During Stage 3 (`StageVerifying`), the spinner was active displaying `[3/5] Verifying SHA256 integrity against checksums.txt...` (length 58 chars).
  - When subsequent shorter messages or stage transitions occurred without clearing the full line width with ANSI `\033[2K\r` / space padding, trailing characters (`txt`) remained in the console buffer.
  - Concurrently, calling `r.multi.NewWriter()` while the background goroutine rendered caused partial frame flushes to remain on screen, followed by a second flush when `r.multi.Stop()` terminated.
- **Target Resolution**:
  - Eliminate `pterm.MultiPrinter` background contention. Drive stages sequentially using a dedicated clean single-spinner lifecycle: start on stage begin, stop and clear cleanly upon stage completion.
  - Explicitly clear terminal lines (`\r\033[2K` or line overwrite) to ensure no trailing string fragments (`txt`) survive.

### 2.3 Success Badge Top Padding
- **Screenshot Evidence**: `assets/screenshots/screenshot-alignment-and-installer-script-adjustment-01.png`, `-02.png`
- **Existing State**: The first `SUCCESS` badge appeared immediately below the active spinner line without vertical separation.
- **Root Cause**: In `cli/cmdupdate/update_ui.go`, `padFirstSuccess()` wrote an empty line into `r.multi.NewWriter()`, which was squashed by `MultiPrinter` frame clearing.
- **Target Resolution**: Ensure `\n` is printed directly before the first `SUCCESS` badge so a clean 1-line vertical gap precedes the success sequence.

### 2.4 Summary Card "Circle" (Top-Right Rounded Corner) Disconnection
- **Screenshot Evidence**: `assets/screenshots/screenshot-alignment-and-installer-script-adjustment-04.png`
- **Pixel Forensics**:
  - Top border rounded corner `╮`: x = 817 px (column 77).
  - Body lines 2–11 right border `│`: x = 818 px (column 77).
  - Bottom border rounded corner `╯`: x = 817 px (column 77).
  - Title line 1 right border `│`: x = 830 px (column 78 — displaced by +1 character cell / 12 pixels!).
- **Root Cause**:
  - In `cli/cmdupdate/update_ui_notes.go`, line 160:
    `title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaGreen)).Render("✔ Gitmap CLI Updated Successfully")`
  - In Windows Terminal, `✔` (`\u2714`, Heavy Check Mark) is rendered by the font fallback / DirectWrite as an emoji taking 2 console columns.
  - In `github.com/charmbracelet/lipgloss` (via `mattn/go-runewidth`), `runewidth.RuneWidth('✔') == 1`.
  - Lipgloss calculates the title line width as 35 columns, whereas the terminal renders it as 36 columns.
  - When Lipgloss pads the title line with spaces up to the box width, it adds 1 extra space, pushing `│` on line 1 out by 1 character cell (12 px to x=830), creating a visible disconnected "broken circle" gap with the top-right corner `╮`.
- **Target Resolution**: Replace `✔` (`\u2714`) with the standard text checkmark `✓` (`\u2713`, `\u2713\uFE0E`), which has deterministic 1-column width in both Lipgloss and terminal emulators, ensuring pixel-perfect alignment.

### 2.5 "Added source line" Subtree Representation
- **Screenshot Evidence**: `assets/screenshots/screenshot-alignment-and-installer-script-adjustment-02.png`
- **Existing State**:
  ```text
    ✓ Added source line to C:\Users\Administrator\Documents\WindowsPowerShell\profile.ps1
    ✓ Added source line to C:\Users\Administrator\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1
    ✓ Added source line to C:\Users\Administrator\Documents\PowerShell\profile.ps1
    ✓ Added source line to C:\Users\Administrator\Documents\PowerShell\Microsoft.PowerShell_profile.ps1
  ```
- **Root Cause**: `cli/completion/install.go:122` printed `MsgCompProfileWrite` per profile inside a loop.
- **Target Resolution**:
  Output "Added source line:" as a header, followed by each path in a clean tree structure:
  ```text
    ✓ Added source line:
      ├── C:\Users\Administrator\Documents\WindowsPowerShell\profile.ps1
      ├── C:\Users\Administrator\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1
      ├── C:\Users\Administrator\Documents\PowerShell\profile.ps1
      └── C:\Users\Administrator\Documents\PowerShell\Microsoft.PowerShell_profile.ps1
  ```

### 2.6 Trailing Blank Line Gap Before Shell Prompt
- **Screenshot Evidence**: `assets/screenshots/screenshot-alignment-and-installer-script-adjustment-03.png`
- **Existing State**: `PS C:\Users\Administrator>` prompt attached directly to the last profile path line without vertical breathing room.
- **Target Resolution**: Print a newline gap (`\n`) after the completion tree before handing control back to the shell prompt.

### 2.7 Shell Installer Script Parity (`install.ps1` & `install.sh`)
- Both PowerShell (`install.ps1`) and POSIX shell (`install.sh`) must reflect identical tree structure for profile updates and ensure trailing newline separation.

---

## 3. Disjoint Task Breakdown & File Ownership

- **Task-01**: `cli/cmdupdate/update_ui.go`, `cli/cmdupdate/update_ui_types.go` (Spinner alignment, margin, `txt` elimination, success top padding).
- **Task-02**: `cli/cmdupdate/update_ui_notes.go`, `cli/completion/install.go`, `cli/constants/constants_completion.go` (Title checkmark `✓` alignment, "Added source line" subtree formatter, trailing newline).
- **Task-03**: `install.ps1`, `install.sh`, `cli/scripts/install.ps1`, `cli/scripts/install.sh`, `cli/fspath/install.ps1`, `cli/fspath/install.sh` (Shell installer script parity and trailing spacing).
