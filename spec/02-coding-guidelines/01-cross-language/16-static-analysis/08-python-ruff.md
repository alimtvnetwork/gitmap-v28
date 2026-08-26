# Python — Ruff / Pylint / Flake8 Enforcement

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`python` · `ruff` · `pylint` · `flake8` · `mypy` · `pyright` · `static-analysis` · `linter`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Maps cross-language coding guidelines to **Ruff** rules (primary), with **Pylint** and **Flake8** equivalents. **mypy** or **pyright** handles type checking. Ruff is preferred as it replaces Flake8, isort, pyupgrade, and most Pylint checks in a single fast tool.

---

## Guideline → Ruff Rule Mapping

| # | Guideline | Ruff Rule | Code | Pylint Equiv. | Severity | Notes |
|---|-----------|-----------|------|---------------|----------|-------|
| 1 | Zero nested `if` | Collapsible if | SIM102 | R1705 | `error` | Flatten nested conditions |
| 2 | No else after return | No else after return/break/continue | RET505, RET506, RET507, RET508 | R1705 | `error` | Early exit pattern |
| 3 | Boolean naming (`is_/has_/can_/should_/was_/will_`) | — | — | C0103 | `warn` | Custom Ruff plugin or code review |
| 4 | No magic strings | Hardcoded string | — | W0012 | — | SonarQube S1192; extract to constants |
| 5 | No magic numbers | Magic value comparison | PLR2004 | R0903 | `error` | Use named constants |
| 6 | Max 15-line functions | Too many statements | PLR0915 | R0915 | `error` | Configure: `max-statements = 10` |
| 7 | Max 3 parameters | Too many arguments | PLR0913 | R0913 | `error` | Configure: `max-args = 3` |
| 8 | Cognitive complexity | McCabe complexity | C901 | — | `error` | Configure: `max-complexity = 10` |
| 9 | DRY — no duplicate code | — | — | R0801 | `warn` | Pylint `duplicate-code`; SonarQube S1871, S4144 |
| 10 | No `Any` type | — | — | — | — | Enforced by mypy `disallow_any_explicit` |
| 11 | Blank line before return | — | — | — | — | Enforced by formatter (Ruff format / Black) |
| 12 | `asyncio.gather` for independent calls | — | — | — | — | Code review; no direct lint |
| 13 | No wildcard imports | Wildcard import | F403 | W0401 | `error` | Explicit imports only |
| 14 | No unused imports | Unused import | F401 | W0611 | `error` | Auto-fixable |
| 15 | No unused variables | Unused variable | F841 | W0612 | `error` | — |

### Python-Specific Rules

| # | Guideline | Ruff Rule | Code | Severity | Notes |
|---|-----------|-----------|------|----------|-------|
| 16 | Use `pathlib` over `os.path` | Use `pathlib` | PTH100–PTH124 | `warn` | Modern path handling |
| 17 | Use f-strings over `.format()` | f-string | UP032 | `error` | Pyupgrade |
| 18 | Use `from __future__ import annotations` | Future annotations | FA100 | `warn` | PEP 563 |
| 19 | No bare `except:` | Bare except | E722 | `error` | Must specify exception type |
| 20 | No mutable default arguments | Mutable default | B006 | `error` | flake8-bugbear |
| 21 | Use `is` for `None` comparison | None comparison | E711 | `error` | `is None` not `== None` |
| 22 | Use type hints everywhere | — | — | — | mypy `strict = true` |
| 23 | No `assert` in production | Assert used | S101 | `error` | bandit / flake8-bandit |
| 24 | No `eval()` / `exec()` | Use of eval | S307 | `error` | Security: code injection |
| 25 | No hardcoded passwords | Hardcoded password | S105, S106 | `error` | bandit security check |
| 26 | Use `dataclass` or `pydantic` over raw dicts | — | — | — | Code review convention |
| 27 | Docstrings on all public functions | Missing docstring | D100–D107 | `warn` | pydocstyle |

---

## SonarQube Rule Mapping (sonar-python)

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S4144 | Identical functions | DRY |
| S1066 | Collapsible if | Zero nesting |
| S1192 | Duplicated string literals | No magic strings |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |
| S107 | Too many parameters | Max 3 parameters |
| S1481 | Unused local variables | No unused variables |
| S5717 | Mutable default argument | No mutable defaults |
| S1126 | Return boolean directly | Boolean principles |

---

## Type Checking — mypy Configuration

```toml

# pyproject.toml — [tool.mypy]

[tool.mypy]
python_version = "3.12"
strict = true
warn_return_any = true
warn_unused_configs = true
disallow_any_explicit = true
disallow_any_generics = true
disallow_untyped_defs = true
check_untyped_defs = true
no_implicit_optional = true
warn_redundant_casts = true
warn_unused_ignores = true
```

---

## Reference Configuration

### `pyproject.toml` — Ruff

```toml
[tool.ruff]
target-version = "py312"
line-length = 100

[tool.ruff.lint]
select = [
  "E",      # pycodestyle errors
  "W",      # pycodestyle warnings
  "F",      # pyflakes
  "I",      # isort
  "N",      # pep8-naming
  "UP",     # pyupgrade
  "S",      # flake8-bandit (security)
  "B",      # flake8-bugbear
  "C90",    # mccabe complexity
  "SIM",    # flake8-simplify
  "RET",    # flake8-return
  "PTH",    # flake8-use-pathlib
  "PLR",    # Pylint refactor
  "PLC",    # Pylint convention
  "PLE",    # Pylint error
  "FA",     # flake8-future-annotations
  "D",      # pydocstyle
]
ignore = [
  "D203",   # conflicts with D211
  "D213",   # conflicts with D212
]

[tool.ruff.lint.mccabe]
max-complexity = 10

[tool.ruff.lint.pylint]
max-args = 3
max-statements = 10

[tool.ruff.lint.per-file-ignores]
"tests/**/*.py" = ["S101", "PLR0913", "D100", "D103"]

[tool.ruff.format]
quote-style = "double"
indent-style = "space"
```

### CI Commands

```bash

# Lint

ruff check . --output-format=github

# Format check

ruff format --check .

# Type check

mypy . --strict
```

### Required Packages

```bash
pip install ruff mypy
```

---

## Tool Comparison

| Feature | Ruff | Pylint | Flake8 |
|---------|------|--------|--------|
| Speed | ⚡ 10–100× faster | Slow | Moderate |
| Auto-fix | ✅ Most rules | ❌ | ❌ |
| Replaces isort | ✅ | ❌ | Plugin |
| Replaces pyupgrade | ✅ | ❌ | Plugin |
| Replaces bandit | ✅ (S rules) | ❌ | Plugin |
| Duplicate detection | ❌ | ✅ (R0801) | ❌ |
| Custom plugins | 🔲 Limited | ✅ | ✅ |

**Recommendation:** Use Ruff as primary + Pylint `duplicate-code` (R0801) only for DRY enforcement.

---

## Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | Add `[tool.ruff]` config to `pyproject.toml` | 🔲 |
| 2 | Add `[tool.mypy]` config to `pyproject.toml` | 🔲 |
| 3 | Install `ruff` and `mypy` in dev dependencies | 🔲 |
| 4 | CI runs `ruff check . --output-format=github` | 🔲 |
| 5 | CI runs `ruff format --check .` | 🔲 |
| 6 | CI runs `mypy . --strict` | 🔲 |
| 7 | SonarQube sonar-python configured | 🔲 |
| 8 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Parent document
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Source rules
- [Master Coding Guidelines](../15-master-coding-guidelines/00-overview.md) — Full checklist
- [Node.js ESLint Enforcement](./07-nodejs-eslint.md) — Sibling server-side spec

---

*Python Ruff enforcement v1.0.0 — cross-language guideline mapping — 2026-04-01*
