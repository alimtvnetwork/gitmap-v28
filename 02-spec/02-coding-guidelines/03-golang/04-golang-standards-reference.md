# Golang Coding Standards (AI Execution Prompt)

> **/goal** Route all Go standards reference inquiries to the canonical modular documentation in `04-golang-standards-reference/`.
> **/learn** Master Go sizing limits (functions <= 15 lines), `appfault.AppError` returns, value semantics, and zero `interface{}` usage.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Route all Go standards inquiries to the partitioned reference files in `04-golang-standards-reference/`.
- [ ] `/learn` Enforce function length <= 15 lines, affirmative booleans, and parameter structs for >3 arguments.
- [ ] `/goal` Mandate `*appfault.AppError` and `Result[T]` container returns across all Go packages.
- [ ] `/learn` Preserve clean stub navigation and verify zero absolute path links across all Go specifications.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> ⚠️ **This file has been split into a subfolder.** See [04-golang-standards-reference/readme.md](./04-golang-standards-reference/readme.md)

All content now lives in [`04-golang-standards-reference/`](./04-golang-standards-reference/).

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-STUB-004: Go Standards Reference Pointer & Subfolder Delegation

**Given** The Go standards reference pointer stub in `02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference.md`.
**When** Audited against this reference specification and coding guidelines.
**Then** Pointer successfully delegates to `04-golang-standards-reference/readme.md` with valid relative links and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations.
