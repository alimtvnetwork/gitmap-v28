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
