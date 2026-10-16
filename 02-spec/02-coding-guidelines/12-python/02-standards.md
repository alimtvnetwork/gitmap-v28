# Python Coding Standards (AI Execution Prompt)

> **/goal** Define mandatory static typing, data validation models, PEP-8 formatting standards, and robust exception handling for Python development.
> **/learn** Master type hints with generics/Union/Optional, Pydantic/dataclass schema validation, Black formatting limits, and typed domain exceptions instead of bare excepts.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Annotate every function argument and return type with explicit hints; eliminate `Any` across signatures.
- [ ] `/learn` Avoid raw dictionaries for complex inputs; use Pydantic models or dataclasses for structured validation.
- [ ] `/goal` Format code with Black (max 100 characters per line) and maintain strict PEP-8 compliance.
- [ ] `/learn` Catch specific exception classes only; never generate bare `except:` or `except Exception:` blocks.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Type Hinting

- Every public function MUST include Python type hints for all arguments and return values.
- Avoid `Any` type. Use generics, `Union`, or `Optional` when needed.

```python
# ❌ FORBIDDEN: Missing type hints and untyped parameters
def calculate_metrics(items, flag):
    return {"total": sum(items), "active": flag}

# ✅ REQUIRED: Explicit type hints and affirmative boolean naming
def calculate_metrics(items: list[float], isActive: bool) -> dict[str, float | bool]:
    totalAmount = sum(items)

    return {"total": totalAmount, "isActive": isActive}
```

## 2. Data Validation

- Use `pydantic` models for structured data validation at system boundaries (APIs, Database inputs, File reads).
- Define explicitly typed attributes within classes using `dataclasses` or `pydantic`.

## 3. PEP-8 Compliance

- Adhere strictly to PEP-8.
- Use `black` for auto-formatting.
- Maximum line length is 100 characters.

## 4. Error Handling

- Never use bare `except:` or `except Exception:`. Always catch specific exception classes.
- Wrap low-level exceptions with the application's domain-specific errors.

```python
# ❌ FORBIDDEN: Bare except catching all exceptions indiscriminately
try:
    process_record(record)
except:
    pass

# ✅ REQUIRED: Catching specific domain exceptions and wrapping
try:
    process_record(record)
except (ValueError, KeyError) as parseErr:
    logger.warning("Failed to parse record: %s", parseErr)
    raise RecordProcessingError("Invalid record payload") from parseErr
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PY-002: Python Coding Standards Conformance

**Given** Polyglot development guidelines and language standards.
**When** Audited against this language specification.
**Then** Zero non-compliant conventions or syntax patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

### AC-CG-PY-003: Python Dynamic Enums & Constants Conformance

**Given** Python scripts and modules defining configurable options or state machines.
**When** Inspected for hardcoded string constants and Cartesian string permutations.
**Then** Identifiers use `Enum` or `StrEnum` types, dynamic array builders eliminate repetitive string combinations, and magic literals are eliminated with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

### AC-CG-PY-004: Python DRY Caching & Architecture Conformance

**Given** Python CI/CD, automation, or linting scripts in `03-ai-scripts/`.
**When** Scanned for duplicate utility implementations or unshared helpers.
**Then** Scripts import shared logic from `03-ai-scripts/02-shared-engine.py`, enforce idempotent caching where applicable, and maintain unified exit codes with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.
