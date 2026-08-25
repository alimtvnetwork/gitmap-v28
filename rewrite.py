import os

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
plan_file = ".lovable/plans/pending/01-ssh-login-and-join.md"

template = """---
plan: {plan_file}
domain: {domain}
phase: {phase}
target_files: ["{file}"]
depends_on: [{depends}]
citations:
  app_spec: "spec/19-ssh-executor/01-spec.md §Section"
  canonical_size: "spec/05-coding-guidelines/01-code-quality-improvement.md"
  language_guideline: "spec/05-coding-guidelines/02-go-code-style.md"
  boolean_styling: "spec/05-coding-guidelines/03-naming-conventions.md"
  folder_naming: "spec/05-coding-guidelines/05-file-project-structure.md"
  error_architecture: "spec/05-coding-guidelines/04-error-handling.md"
  error_codes: "spec/05-coding-guidelines/04-error-handling.md"
  logging_traces: "spec/05-coding-guidelines/07-logging-observability.md"
  response_envelope: "spec/05-coding-guidelines/10-api-design.md"
  golden_fixture: "spec/08-json-schemas/ssh-list.schema.json"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "spec/05-coding-guidelines/11-database-patterns.md"
  ui_surface: "n/a — cli tool"
  tests: "unit Test{sym}"
  ci_cd_guard: ".github/workflows/ci.yml"
  ambiguity: "n/a — spec defined"
  issue_rca: "n/a — feature work"
---
# Task {idx:03d} — {title}

## 1. Learn
- [Spec: SSH Commands](file:///d:/work/gitmap/.lovable/spec/commands/01-ssh-commands.md) — Why: To understand the required CLI syntax.
- [Spec: Troubleshooting](file:///d:/work/gitmap/.lovable/issues/06-ssh-troubleshooting-guide.md) — Why: To prevent syntax mistakes in SSH logic.
- [App Error Docs](file:///d:/work/gitmap/spec/05-coding-guidelines/04-error-handling.md) — Why: To properly return apperror.Result.
- [{file}](file:///d:/work/gitmap/{file}) — Why: Target file for implementation.

## 2. Goal
Implement the {phase} step for {sym} to support {title}. This enables the CLI to handle the required domain logic without affecting existing commands.

## 3. Inputs and Contracts
- **Input Types**: `string`, `context.Context`
- **Output Types**: `error`, `apperror.Result[bool]`
- **Error Codes**: `E_INVALID_ARGS`, `E_INTERNAL_ERROR`
- **Contract**:
  ```go
  func {sym}(...) error
  ```

## 4. Execute
1. Open `{file}`.
2. Define `{sym}` according to the {phase} requirements.
3. Ensure no external dependencies are unnecessarily imported.

## 5. Constraints
- **Canonical Size**: Must respect function size limits (rule: spec/05-coding-guidelines/01-code-quality-improvement.md).
- **Boolean Prefix**: Use `is` or `has` for booleans (rule: spec/05-coding-guidelines/03-naming-conventions.md).
- **No Globals**: Do not use global state (rule: .lovable/strictly-avoid.md).

## 6. Verify
Run the following command to verify the {phase} changes:
```bash
go test ./... -v -run {sym}
```
Expected output:
```text
PASS
ok      github.com/alimtvnetwork/gitmap-v28/gitmap/...
```

## 7. Done When
- [ ] 1. The function `{sym}` is successfully implemented.
- [ ] 2. The tests pass.
- [ ] 3. No canonical size violations exist.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.
"""

import json
units = []
with open("gen_ssh_plan.py", "r") as f:
    in_units = False
    block = ""
    for line in f:
        if line.strip() == "units = [":
            in_units = True
        if in_units:
            block += line
            if line.strip() == "]":
                break
    
    # Eval the block
    loc = {}
    exec(block, {}, loc)
    units = loc["units"][:50]

for i, u in enumerate(units, 1):
    depends = f"Task {i-1:03d}" if i > 1 else "None"
    content = template.format(idx=i, title=u["title"], domain=u["domain"], phase=u["phase"], sym=u["sym"], file=u["file"], plan_file=plan_file, depends=depends)
    with open(os.path.join(task_dir, f"{i:03d}-task.md"), "w", encoding="utf-8") as f:
        f.write(content)
print(f"Generated {len(units)} files.")
