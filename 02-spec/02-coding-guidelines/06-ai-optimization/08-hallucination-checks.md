# Hallucination Checks & Prevention (AI Execution Prompt)

> **/goal** Eliminate AI code hallucinations by enforcing the "Read Before Write" protocol, strict compile-time typing, and immediate static analysis validation.
> **/learn** Prevent invented endpoints, non-existent database columns, and fictitious library APIs through rigorous grounded exploration and pre-commit checks.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Execute "Read Before Write": read routers, schemas, and manifests before generating calls or references.
- [ ] `/learn` Rely on strict typing (`strict: true`, `Result[T]`, `mypy`) to catch invented properties at compile time.
- [ ] `/goal` Verify generated code with local static analysis and linter checks prior to concluding tasks.
- [ ] `/learn` Never assume library APIs or external contracts exist without inspecting package dependencies or official specs.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Defining "Hallucination" in Code

An AI agent "hallucinates" when it:

- Invents API endpoints that do not exist in the backend.
- Proposes database columns or tables that were never defined in the schema.
- Uses third-party library functions that do not exist or are deprecated.

## 2. Prevention Strategies

### The "Read Before Write" Protocol

Agents MUST explicitly search and read the relevant definition files (e.g., router definitions, database schemas, `package.json`) before calling functions or endpoints.

### Strict Null and Type Enforcement

Ensure that the codebase enforces strict typing (TypeScript `strict: true`, Go `Result[T]`, Python `mypy`). A hallucinated property access will immediately fail CI/CD build steps, catching the hallucination before it reaches production.

### Verification Step

Agents are required to verify their own work by running static analysis or build commands (`npm run build`, `go build`, `cargo check`) immediately after generating a block of code.

```typescript
// ❌ FORBIDDEN: Hallucinating non-existent properties without inspecting schema
const payload = {
  userName: "alice",
  sendEmailNotification: true, // Hallucinated field
};

// ✅ REQUIRED: Grounded field calls strictly matching verified schema
const payload: CreateUserParams = {
  userName: "alice",
  isNotificationEnabled: true,
};
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-ai-optimization/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-AI-008: Hallucination Prevention & Static Verification

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.
