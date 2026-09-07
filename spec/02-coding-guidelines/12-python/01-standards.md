# Python Coding Standards

## 1. Type Hinting

- Every public function MUST include Python type hints for all arguments and return values.
- Avoid `Any` type. Use generics, `Union`, or `Optional` when needed.

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

## 5. Cross-Platform Shell & Subprocess Execution

- Automation and helper scripts invoking shell binaries MUST NOT hardcode binary paths.
- Scripts MUST use centralized OS-detection helpers (such as `get_bash_path()` in `03-ai-scripts/02-shared-engine.py`).
- Functions implementing OS or executable detection MUST adhere to core coding guidelines:
  - Nesting depth <= 1 (no nested `if` statements).
  - Affirmative boolean naming (`is_windows`, `is_valid_executable`).
  - Centralized constants for OS names, fallback commands, and candidate paths (no magic strings).
  - Function length <= 15 lines.

