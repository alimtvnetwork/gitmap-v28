# Citation Requirement for AI Agents

## Mandatory Citation Rule (CODE RED)

Whenever an AI agent generates code, explains a design decision, or enforces a standard, it **MUST** cite the specific `spec/` markdown file and line/section that justifies the action.

### Why This is Required

- It prevents agents from blending external training data with this repository's strict conventions.
- It provides human reviewers with an immediate paper trail to verify that the agent followed the house style.

### Examples of Valid Citations

- *"Implementing this as an early return to avoid nesting, per `spec/02-coding-guidelines/01-cross-language/01-zero-nesting.md`."*
- *"Returning a structured error with context, per `spec/03-error-manage/02-error-architecture/01-error-handling-reference.md`."*

### Violations

If an agent enforces a rule (e.g., "Variables must be named X") but cannot cite a spec file to back it up, it has failed the anti-hallucination contract. Human reviewers should reject such suggestions.
