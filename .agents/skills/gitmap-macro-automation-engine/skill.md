---
name: gitmap-macro-automation-engine
description: Autonomously extend and troubleshoot interactive macro recording, replay, cross-platform shell adapters, and automated task execution.
---

# GitMap Macro Automation Engine Skill (`gitmap-macro-automation-engine`)

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for developing, maintaining, and executing interactive multi-step macros, cross-platform shell command shims, and batch automation scripts in GitMap.

---

## 1. Key Architectural Components & Code Map

| Component | Primary Location | Key Responsibilities |
|---|---|---|
| **Macro Executor** | `cli/cmdmacro/`, `cli/macro/` | Multi-step macro command runner, live output streaming, step timeout protection, and exit code handling. |
| **Interactive Macro Builder** | `cli/macro/builder.go` | Interactive step recording, dynamic PWD header display, and in-builder helpers (`ls`, `find`, `search`, `replace`). |
| **Cross-Platform Shell Shims** | `cli/cmdmacro/macro_open.go` | OS-aware translation of desktop commands (e.g. `open` -> `cmd.exe /c start ""` on Windows, `xdg-open` on Linux, `/usr/bin/open` on macOS). |
| **Macro Storage & Export** | `cli/macro/export.go`, `cli/cmdmacro/` | YAML/JSON serialization of macro workflows, macro imports, and validation. |

---

## 2. Essential Commands

```bash
# Run a registered macro by name
gitmap macro run <macro-name>
gitmap <macro-name>

# Launch interactive macro builder to record steps
gitmap macro create <macro-name>

# List all registered macros
gitmap macro list

# Export macro to shareable definition file
gitmap macro export <macro-name> --output ./macro.yaml

# Import macro definition from file
gitmap macro import ./macro.yaml
```

---

## 3. Core Invariants & Engineering Guardrails

1. **Cross-Platform Desktop Command Shimming:** Never assume Unix utilities exist on Windows. Shell commands like `open` must route through OS-aware adapters (`start` on Windows, `xdg-open` on Linux, `/usr/bin/open` on macOS) as resolved in `.ai-memory/ambiguous-questions/02-ambiguity-resolved/02-macro-step-open-command-behavior.md`.
2. **Discrete Argument Passing:** Always pass discrete, structured arguments into `exec.Command` rather than passing concatenated shell strings to prevent Windows quote stripping and injection hazards.
3. **Detached Process Safety:** When launching detached UI applications (such as browsers or editors), ensure process handles are cleanly decoupled so the macro runner does not hang waiting for GUI windows to close.
4. **Step Timeout & Failure Isolation:** Every macro step must have a bounded timeout and clear error attribution using `*apperror.AppError` (`[E9000:EXECUTION]`).
5. **No Nested Ifs:** Keep macro execution pipelines shallow and maintainable using early returns and guard clauses (zero nesting rule).
