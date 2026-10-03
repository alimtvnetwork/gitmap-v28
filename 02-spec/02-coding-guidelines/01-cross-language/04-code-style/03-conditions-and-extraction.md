# Condition Extraction (AI Execution Prompt)

> **/goal** Eliminate complex inline multi-part conditional expressions in `if` statements across PHP, TypeScript, and Go by extracting them into named boolean variables, dedicated methods, or constants.
> **/learn** Understand single-intent control flow, 2+ operator extraction triggers (`&&`, `||`, `!`), affirmative boolean naming (`is*`, `has*`), and reusable type-guard function patterns.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Ban inline multi-part conditions containing 2 or more operators (`&&`, `||`, `!`) in `if` statements.
- [ ] `/learn` Extract local, one-off multi-part checks into positively named boolean variables (`is*`, `has*`).
- [ ] `/goal` Extract multi-part conditions reused across multiple places or carrying domain significance into dedicated helper functions or type guards.
- [ ] `/learn` Separate error checks from domain logic guards so error paths exit immediately before evaluating domain booleans.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Code Style](./readme.md)
> **Version:** 4.0.0
> **Updated:** 2026-03-31
> **Applies to:** PHP, TypeScript, Go
> **Rules covered:** 3

---

## Rule 3: Extract Complex Conditions — No Inline Multi-Part Checks

When an `if` condition contains **two or more operators** (`&&`, `||`, `!`), it **must** be extracted into one of:

1. **A named boolean variable** (`$is_*` / `$has_*` / `isX` / `hasX`) — for local, one-off checks
2. **A dedicated method/function** — for reusable or domain-meaningful checks
3. **A named constant** — for static flag combinations

The goal: every `if` reads as a **single intent**, not as implementation logic.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Inline multi-part condition
if ($error && in_array($error['type'], [E_ERROR, E_PARSE], true)) {
    $this->logger->fatal($error);
}

// ✅ REQUIRED: Extracted into a dedicated method
if (ErrorChecker::isFatalError($error)) {
    $this->logger->fatal($error);
}

// ❌ FORBIDDEN: Combinable conditions left inline
if ($request !== null && $request->has_param('file') && $request->get_param('file') !== '') {
    $this->process($request);
}

// ✅ REQUIRED: Named boolean for clarity
$hasFileParam = $request !== null
    && $request->hasParam('file')
    && $request->getParam('file') !== '';

if ($hasFileParam) {
    $this->process($request);
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Inline multi-part condition
if (response && response.status >= 400 && response.data?.code?.startsWith('E8')) {
    showDelegatedError(response);
}

// ✅ REQUIRED: Named boolean
const isDelegatedError = response != null
    && response.status >= 400
    && response.data?.code?.startsWith('E8');

if (isDelegatedError) {
    showDelegatedError(response);
}

// ✅ ALSO OK: Dedicated type-guard function for reusable checks
function isDelegatedError(res: ApiResponse | null): res is DelegatedErrorResponse {
    return res != null && res.status >= 400 && res.data?.code?.startsWith('E8');
}

if (isDelegatedError(response)) {
    showDelegatedError(response);
}
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Inline multi-part condition
if err != nil && resp != nil && resp.StatusCode >= 400 {
    handleUpstreamError(resp)
}

// ✅ REQUIRED: Separate error guard, then named boolean for domain check
if err != nil {
    return apperror.Wrap[UpstreamResult](err, errors.ErrUpstream, "upstream call failed")
}

isUpstreamError := resp != nil && resp.StatusCode >= 400

if isUpstreamError {
    handleUpstreamError(resp)
}
```

### When to Use Which Extraction

| Complexity | Extraction | Example |
|------------|-----------|---------|
| 2 conditions, used once | Named `$is_*` / `isX` variable | `$hasFile = $req !== null && $req->hasParam('file');` |
| 2+ conditions, used in multiple places | Dedicated method/function | `ErrorChecker::isFatalError($error)` |
| Static flag combination | Named constant | `const EDITABLE = 'PUT, PATCH';` |

---

*Part of [Code Style](./readme.md) — Rule 3*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-STYLE-003: Multi-Part Condition Extraction and Guard Inversion

**Given** Conditional logic expressions across PHP, TypeScript, and Go codebases.
**When** Linters and static analysis check `if` statements for multiple logical operators (`&&`, `||`, `!`).
**Then** All conditions with 2+ operators are extracted into named boolean variables or dedicated functions, reading as a single intent.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style --check-only
```
**Expected:** exit 0. Zero violations.
