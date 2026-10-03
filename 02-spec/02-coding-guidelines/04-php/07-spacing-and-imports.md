# PHP Spacing and Import Rules (AI Execution Prompt)

> **/goal** Eliminate formatting friction, missing vertical blank lines, leading backslash type references, and raw string log keys across all PHP files in the `RiseupAsia` namespace.
> **/learn** Master PSR-12 and architectural standards: mandatory blank line before `if` and `throw`, top-level `use` imports for global classes/exceptions, and `ResponseKeyType` enum usage for multi-file log context keys.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce mandatory blank line before `if` statements when preceded by executable statements.
- [ ] `/goal` Enforce mandatory blank line before `throw` statements when preceded by executable statements.
- [ ] `/learn` Eliminate all leading backslash global type references (`\RuntimeException`, `\Throwable`) via file-level `use` imports.
- [ ] `/learn` Replace all multi-file raw string log context keys with `ResponseKeyType` enum instances.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None
**Applies to:** All PHP files in the `RiseupAsia` namespace
**Source:** Consolidated from `04-coding-guidelines-wpon/07-php-standards/php-spacing-and-imports.md`

---

## Rule 1: Blank Line Before `if` When Preceded by Statements

When an `if` block is preceded by one or more statements, insert one blank line before the `if`.

**Exception:** No blank line when `if` is the first statement in a function body, or immediately follows another `}`.

```php
// ❌ FORBIDDEN: No blank line between statement and if
$existingRunning = $this->findRunningProcess();
if ($existingRunning !== null) {
    Logger::warning('Scan already running', array('existingId' => $existingRunning->id));

    throw new RuntimeException('A scan is already in progress', 14100);
}

// ✅ REQUIRED: Blank line separates setup from decision
$existingRunning = $this->findRunningProcess();

if ($existingRunning !== null) {
    Logger::warning('Scan already running', array('existingId' => $existingRunning->id));

    throw new RuntimeException('A scan is already in progress', 14100);
}
```

---

## Rule 2: Blank Line Before `throw` When Preceded by Statements

Same as `return`: if a `throw` is preceded by one or more statements in the same block, insert one blank line before it.

```php
// ❌ FORBIDDEN: Missing blank line before throw
if ($existingRunning !== null) {
    Logger::warning('Scan already running', array('existingId' => $existingRunning->id));
    throw new RuntimeException('A scan is already in progress', 14100);
}

// ✅ REQUIRED: Blank line before throw
if ($existingRunning !== null) {
    Logger::warning('Scan already running', array('existingId' => $existingRunning->id));

    throw new RuntimeException('A scan is already in progress', 14100);
}
```

---

## Rule 3: No Leading Backslash — Use `use` Import

In namespaced PHP files, **never** reference global types with a leading backslash. Add a `use` import at the top instead.

```php
// ❌ FORBIDDEN: Leading backslash global type reference
throw new \RuntimeException('...');
catch (\Throwable $e) { ... }

// ✅ REQUIRED: use import at file top
use RuntimeException;
use Throwable;

throw new RuntimeException('...');
catch (Throwable $e) { ... }
```

**Exemptions:**

- `Autoloader.php` — must be self-contained
- Main plugin bootstrap file — may use backslash before autoloader is registered

---

## Rule 4: Reusable Log Context Keys Must Use Enums

Log context keys follow camelCase. But reusable keys appearing in 3+ log calls across different files must use `ResponseKeyType` enum.

```php
// ❌ FORBIDDEN: Raw string used across multiple files
Logger::warning('Scan running', array('existingId' => $id));

// ✅ REQUIRED: Backed enum for reusable key
Logger::warning('Scan running', array(ResponseKeyType::ExistingId->value => $id));
```

---

## Combined Example — All Rules

```php
// ❌ FORBIDDEN: Four spacing and import violations
$existingRunning = $this->findRunningProcess();
if ($existingRunning !== null) {
    Logger::warning('Scan already running', ['existing_id' => $existingRunning->id]);
    throw new \RuntimeException('A scan is already in progress', 14100);
}

// ✅ REQUIRED: All spacing and import rules applied
$existingRunning = $this->findRunningProcess();

if ($existingRunning !== null) {
    Logger::warning('Scan already running', array('existingId' => $existingRunning->id));

    throw new RuntimeException('A scan is already in progress', 14100);
}
```

---

## Cross-References

- [Code Style](../01-cross-language/04-code-style/readme.md) — Rules R4, R5, R10
- [PHP Naming Conventions](../../01-spec-authoring-guide/03-naming-conventions.md) — Array key casing
- [PHP Forbidden Patterns](./03-forbidden-patterns.md) — Banned patterns

---

*PHP spacing and import rules — consolidated from WPOnboard coding guidelines.*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PHP-007: PHP Spacing, Blank Lines and Import Hygiene

**Given** PHP source code in the `RiseupAsia` namespace.
**When** Files are audited against vertical spacing, import hygiene, and log context key standards.
**Then** All `if` and `throw` blocks have proper preceding blank lines, global classes use top-level `use` imports without leading backslashes, and reusable log context keys utilize backed enums with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.
