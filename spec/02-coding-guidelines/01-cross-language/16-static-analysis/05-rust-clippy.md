# Rust — Clippy Lint Enforcement

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`rust` · `clippy` · `static-analysis` · `linter` · `rustfmt` · `cargo-clippy`

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

Maps cross-language coding guidelines to **Clippy lints** and **Rust compiler warnings** for enforcement. Clippy is the primary linter; `rustfmt` handles formatting.

---

## Guideline → Clippy Lint Mapping

| # | Guideline | Clippy Lint / Rust Warning | Severity | Notes |
|---|-----------|---------------------------|----------|-------|
| 1 | Zero nested `if` | `clippy::collapsible_if`, `clippy::collapsible_else_if` | `deny` | Flatten conditions; use early returns and guard clauses |
| 2 | No else after return/break/continue | `clippy::needless_return`, `clippy::redundant_else` | `deny` | Forces early-exit pattern |
| 3 | Boolean naming (`is_/has_/can_/should_/was_/will_`) | `clippy::wrong_self_convention` | `warn` | Convention check; custom lint recommended for full enforcement |
| 4 | No magic strings/numbers | `clippy::unreadable_literal`, `clippy::approx_constant` | `deny` | Use `const` declarations; SonarQube S1192 for string literals |
| 5 | Max 15-line functions | `clippy::too_many_lines` | `deny` | Configure: `too-many-lines-threshold = 15` |
| 6 | Max 3 parameters | `clippy::too_many_arguments` | `deny` | Configure: `too-many-arguments-threshold = 3` |
| 7 | No `unwrap()` in production | `clippy::unwrap_used`, `clippy::expect_used` | `deny` | Use `?` operator or explicit error handling |
| 8 | Handle all Results/Options | `clippy::let_underscore_must_use`, `#[must_use]` | `deny` | Compiler `unused_must_use` is `deny` by default |
| 9 | No `clone()` without justification | `clippy::clone_on_copy`, `clippy::redundant_clone` | `deny` | Prefer borrowing |
| 10 | DRY — no duplicate code | `clippy::same_functions_in_if_conditions` | `warn` | SonarQube S1871, S4144 for broader detection |
| 11 | Cognitive complexity | `clippy::cognitive_complexity` | `deny` | Configure: `cognitive-complexity-threshold = 10` |
| 12 | Use iterators over manual loops | `clippy::manual_filter_map`, `clippy::manual_find_map`, `clippy::needless_range_loop` | `deny` | Idiomatic Rust |
| 13 | No wildcard imports | `clippy::wildcard_imports` | `deny` | Explicit imports only |
| 14 | Single responsibility types | `clippy::struct_excessive_bools` | `warn` | Configure: `max-struct-bools = 3` |
| 15 | Blank line before return | — | — | Enforced by `rustfmt` custom config or code review |

---

## SonarQube Rule Mapping (sonar-rust — Community)

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S1066 | Collapsible if | Zero nesting |
| S1192 | Duplicated string literals | No magic strings |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |
| S107 | Too many parameters | Max 3 parameters |
| S4144 | Identical functions | DRY |
| S1126 | Return boolean directly | Boolean principles |

> **Note:** sonar-rust is a community plugin with limited coverage. Clippy is the primary enforcement tool.

---

## Lint Groups

Enable these Clippy lint groups project-wide:

| Group | Purpose |
|-------|---------|
| `clippy::pedantic` | Stricter lints beyond default — enable selectively |
| `clippy::nursery` | Newer lints — enable for cutting-edge checks |
| `clippy::cargo` | Cargo.toml hygiene |
| `clippy::unwrap_used` | Part of `clippy::restriction` — deny individually |

---

## Reference Configuration

### `clippy.toml`

```toml

# clippy.toml — project root

too-many-arguments-threshold = 3
too-many-lines-threshold = 15
cognitive-complexity-threshold = 10
max-struct-bools = 3
type-complexity-threshold = 200
```

### `lib.rs` / `main.rs` — Crate-Level Attributes

```rust
// Deny critical lints project-wide
#![deny(clippy::unwrap_used)]
#![deny(clippy::expect_used)]
#![deny(clippy::redundant_else)]
#![deny(clippy::collapsible_if)]
#![deny(clippy::collapsible_else_if)]
#![deny(clippy::wildcard_imports)]
#![deny(clippy::redundant_clone)]
#![deny(clippy::cognitive_complexity)]
#![deny(clippy::too_many_arguments)]
#![deny(clippy::too_many_lines)]
#![deny(clippy::needless_range_loop)]
#![deny(clippy::manual_filter_map)]
#![deny(clippy::manual_find_map)]
#![deny(clippy::let_underscore_must_use)]

// Warn on style lints
#![warn(clippy::pedantic)]
#![warn(clippy::wrong_self_convention)]
#![warn(clippy::struct_excessive_bools)]
#![warn(clippy::same_functions_in_if_conditions)]

// Compiler warnings as errors
#![deny(unused_must_use)]
#![deny(unused_variables)]
#![deny(dead_code)]
```

### `rustfmt.toml`

```toml
edition = "2021"
max_width = 100
tab_spaces = 4
newline_style = "Unix"
use_field_init_shorthand = true
```

### CI Command

```bash
cargo clippy --all-targets --all-features -- -D warnings
cargo fmt --all -- --check
```

---

## Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | Add `clippy.toml` to project root | 🔲 |
| 2 | Add crate-level `#![deny(...)]` attributes | 🔲 |
| 3 | Add `rustfmt.toml` to project root | 🔲 |
| 4 | CI runs `cargo clippy -- -D warnings` | 🔲 |
| 5 | CI runs `cargo fmt --all -- --check` | 🔲 |
| 6 | SonarQube sonar-rust plugin configured (optional) | 🔲 |
| 7 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Parent document
- [Rust Coding Standards](../../05-rust/00-overview.md) — Rust-specific guidelines
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Source rules
- [Master Coding Guidelines](../15-master-coding-guidelines/00-overview.md) — Full checklist

---

*Rust Clippy enforcement v1.0.0 — cross-language guideline mapping — 2026-04-01*
