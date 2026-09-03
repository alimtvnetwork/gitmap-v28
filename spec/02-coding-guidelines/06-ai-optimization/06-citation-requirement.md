# Citation & Relative Path Requirement for AI Agents

## Mandatory Citation & Relative Path Rule (CODE RED)

Whenever an AI agent generates code, creates plans (`.lovable/plans/pending/`), breaks tasks into subtasks (`.lovable/plans/subtasks/`), writes memory logs (`.lovable/memory/issues/`), explains design decisions, or enforces standards, it **MUST** cite the specific `spec/` or `.lovable/` markdown file and line/section that justifies the action using **STRICTLY RELATIVE PATHS FROM THE GIT REPOSITORY ROOT**.

### 1. Total Ban on Absolute Paths & `file:///` URIs in Repository Files

- **TOTAL BAN:** NEVER write absolute filesystem paths (e.g. `/absolute/path/to/...`, `C:\Users\...`, `/home/...`) or absolute URI schemes (`file:///absolute/path/to/work/...`, `file:///absolute/path/to/`) inside markdown plans, subtask files, code comments, citations, or committed repository files.
- **PORTABILITY REQUIREMENT:** All paths and markdown links within repository files MUST be relative to the git root so they work seamlessly across Windows, Linux, macOS, and CI/CD pipelines.

### 2. Concrete Examples

#### ❌ INVALID (Absolute Path / File URI):

```markdown
- [SSH Commands](file:///absolute/path/to/...) — Why: Defines required behavior.
- [App Error Docs](file:///absolute/path/to/02-spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [cmd/main.go](file:///absolute/path/to/cmd/main.go) — Why: Target file.
```

#### ✅ VALID (Strict Relative Git Path):

```markdown
- [SSH Commands](02-spec/13-generic-cli/01-index.md) — Why: Defines required behavior.
- [App Error Docs](02-spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [cmd/main.go](cmd/main.go) — Why: Target file.
```
### Why This is Required

- It prevents agents from blending external training data with this repository's strict conventions.
- It provides human reviewers with an immediate paper trail to verify that the agent followed the house style.

### Examples of Valid Citations

- *"Implementing this as an early return to avoid nesting, per `02-spec/02-coding-guidelines/01-cross-language/01-zero-nesting.md`."*
- *"Returning a structured error with context, per `02-spec/03-error-manage/02-error-architecture/01-index.md`."*

### Violations

If an agent enforces a rule (e.g., "Variables must be named X") but cannot cite a spec file to back it up, it has failed the anti-hallucination contract. Human reviewers should reject such suggestions.
If an agent outputs absolute file paths (`file:///` or drive letters) into repository markdown files or plans, or enforces a rule without citing a valid relative spec path, it has failed the anti-hallucination contract.
