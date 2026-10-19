# Execution Plan: Screenshot Alignment and Installer Script Adjustment

- **Task Slug**: `screenshot-alignment-and-installer-script-adjustment`
- **Specification**: `02-spec/21-app/screenshot-alignment-and-installer-script-adjustment/spec.md`
- **Status**: IN_PROGRESS
- **Budget**: 400 Steps
- **Capacity**: A = 3 agents, H = 2 hands, Concurrency = 6

---

## 1. Verbatim User Request

```text
# Screenshot Alignment and Installer Script Adjustment: high priority instruction, non-negotiable task

I have given you three screenshots. First of all, the blue color that you see or created needs to be a little bit right-hand side so that it's aligned with the white one at the end. The second is that your numbers, one by five, four by five, need to be a little bit right-hand side to match with the probing section. When you start with the success, have a little bit of gap, like a new line before you start with the success step. Fix the TXT broken issue. Inside the circle, it looks nice but a little bit broken. Look into the circle that is created. A little bit broken, so let's calculate that. How long do you grade this? At the end, when you say, "Added source line," that will give a bunch of same bugs. Say, "Added source line," and then Enter, and then as a subtree, add these paths. It would make more sense and be much clearer. At the end, after you complete the tree, there should be a new line gap before you hand it over to the user or user types. Make sure the same happens for the shell installer script too.
```

---

## 2. Discrete Subtask Allocations

| Task ID | Title | Owned Files | Owner | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Task-01** | Spinner Indentation, Padding & `txt` Ghost Removal | `cli/cmdupdate/update_ui.go`, `cli/cmdupdate/update_ui_types.go` | Worker 01 | PENDING |
| **Task-02** | Card Circle Alignment & "Added source line" Subtree | `cli/cmdupdate/update_ui_notes.go`, `cli/completion/install.go`, `cli/constants/constants_completion.go` | Worker 02 | PENDING |
| **Task-03** | Installer Script Parity (`install.ps1`, `install.sh`) | `install.ps1`, `install.sh`, `cli/scripts/install.ps1`, `cli/scripts/install.sh`, `cli/fspath/install.ps1`, `cli/fspath/install.sh` | Worker 03 | PENDING |

---

## 3. Verification Contracts

- **Task-01**: Unit test in `cli/cmdupdate` verifying spinner indentation and clean stage exit without stray text.
- **Task-02**: Unit test in `cli/completion` verifying tree formatter for profile paths and `\n` spacing; unit test in `cli/cmdupdate` verifying card border widths across all lines match exactly.
- **Task-03**: Static inspection and PowerShell / Bash syntax validation verifying profile tree formatting and clean trailing newline.
