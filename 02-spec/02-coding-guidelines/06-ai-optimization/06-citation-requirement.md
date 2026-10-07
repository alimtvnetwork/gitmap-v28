# Citation & Relative Path Requirement for AI Agents (AI Execution Prompt)

> **/goal** Enforce strict relative git repository paths and mandatory spec citation across all AI agent interactions, plans, code generation, and memory records.
> **/learn** Eliminate absolute filesystem paths, Windows drive letters, and `file:///` URIs; anchor all decisions and rules directly to canonical `02-spec/` relative paths.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Cite specific relative `02-spec/` or `.ai-memory/` paths and sections whenever generating code, plans, or design rationales.
- [ ] `/learn` Never output absolute filesystem paths (`/home/...`, `C:\...`) or `file:///` URIs in any repository documentation or plans.
- [ ] `/goal` Ensure every enforced rule or convention has a verifiable citation to an existing markdown specification.
- [ ] `/learn` Verify zero absolute path citations and 100% relative git root formatting across all generated artifacts.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## Mandatory Citation & Relative Path Rule (CODE RED)

Whenever an AI agent generates code, creates plans (`.ai-memory/plans/pending/`), breaks tasks into subtasks (`.ai-memory/plans/subtasks/`), writes memory logs (`.ai-memory/memory/issues/`), explains design decisions, or enforces standards, it **MUST** cite the specific `02-spec/` or `.ai-memory/` markdown file and line/section that justifies the action using **STRICTLY RELATIVE PATHS FROM THE GIT REPOSITORY ROOT**.

### 1. Total Ban on Absolute Paths & `file:///` URIs in Repository Files

- **TOTAL BAN:** NEVER write absolute filesystem paths (e.g. `/absolute/path/to/...`, `C:\Users\...`, `/home/...`) or absolute URI schemes (`file:///absolute/path/to/work/...`, `file:///absolute/path/to/`) inside markdown plans, subtask files, code comments, citations, or committed repository files.
- **STRICT RELATIVE GIT PATHS ONLY:** Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
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
- [SSH Commands](02-spec/13-generic-cli/readme.md) — Why: Defines required behavior.
- [App Error Docs](02-spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [cmd/main.go](cmd/main.go) — Why: Target file.
```

### Why This is Required

- It prevents agents from blending external training data with this repository's strict conventions.
- It provides human reviewers with an immediate paper trail to verify that the agent followed the house style.

### Examples of Valid Citations

- *"Implementing this as an early return to avoid nesting, per `02-spec/02-coding-guidelines/01-cross-language/01-zero-nesting.md`."*
- *"Returning a structured error with context, per `02-spec/03-error-manage/02-error-architecture/readme.md`."*

### Violations

If an agent enforces a rule (e.g., "Variables must be named X") but cannot cite a spec file to back it up, it has failed the anti-hallucination contract. Human reviewers should reject such suggestions.
If an agent outputs absolute file paths (`file:///` or drive letters) into repository markdown files or plans, or enforces a rule without citing a valid relative spec path, it has failed the anti-hallucination contract.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-ai-optimization/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-AI-006: Mandatory Citations & Relative Path Governance

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.
