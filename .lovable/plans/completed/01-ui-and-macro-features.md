# UI and Macro Features

Status: completed
Execution model: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.

## Context

- Release policy: One release at the end, do not bump version per task.
- Issues: .lovable/issues/01-clone-ui.md
- Spec gaps: .lovable/ambiguous-questions/01-new-ambiguity/01-spec-gaps.md

## CI/CD verification

None for now.

## Coding Guideline Checklist

| Topic                            | Single source file                                                          | Duplicates found |
| canonical size tier              | n/a             | none             |
| boolean naming prefixes          | n/a      | none |
| boolean guards + extraction      | n/a | none |
| boolean params + conditions      | n/a | none |
| boolean exemptions + api         | n/a   | none |
| boolean quick reference          | n/a      | none |
| boolean flag methods             | n/a      | none             |
| no negatives                     | n/a              | none             |
| braces + nesting                 | n/a | none      |
| conditions + extraction (style)  | n/a | none |
| blank lines + spacing            | n/a | none |
| function + type size             | n/a | none  |
| multi-line formatting            | n/a | none   |
| code-style checklist             | n/a   | none             |
| nesting resolution               | n/a | none           |
| cyclomatic complexity            | n/a     | none             |
| code mutation avoidance          | n/a   | none             |
| strict typing                    | n/a             | none             |
| null-pointer safety              | n/a       | none             |
| naming + casing (keys)           | n/a     | none             |
| file/folder naming               | n/a               | none             |
| testing                          | n/a | none             |
| error handling + codes           | n/a                   | none             |
| error code registry              | n/a                                | none             |
| logging + stack traces           | n/a             | none             |
| serialization/determinism        | n/a                                               | none             |
| ci/cd verification               | n/a                           | none             |
| ci guards                        | n/a        | none             |
| contract + e2e testing           | n/a       | none       |
| static analysis / sarif          | n/a          | none             |
