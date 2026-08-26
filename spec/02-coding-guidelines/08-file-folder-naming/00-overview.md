# File & Folder Naming Conventions

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**AI Confidence:** High  
**Ambiguity:** None

---

## Keywords

`file-naming` · `folder-naming` · `directory-structure` · `conventions` · `cross-language` · `wordpress` · `php` · `golang` · `typescript` · `rust` · `csharp`

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

Defines file and folder naming conventions for **every language** in the project. Each language has its own file with specific rules, examples, and forbidden patterns. This is the single source of truth for naming files and directories.

---

## Golden Rule

> **Consistency within a language ecosystem matters more than personal preference.** Follow the convention of the language/framework, not your own habits.

---

## Categories

| # | File | Language | Convention Summary |
|---|------|----------|-------------------|
| 01 | [01-cross-language.md](./01-cross-language.md) | All | Universal rules that apply everywhere |
| 02 | [02-php-wordpress.md](./02-php-wordpress.md) | PHP / WordPress | WordPress plugin/theme file & folder naming |
| 03 | [03-golang.md](./03-golang.md) | Go | Package-based naming, flat structure |
| 04 | [04-typescript-javascript.md](./04-typescript-javascript.md) | TypeScript / JS | Component files, hooks, utilities |
| 05 | [05-rust-csharp.md](./05-rust-csharp.md) | Rust / C# | snake_case (Rust) and PascalCase (C#) |

---

## Quick Reference

| Language | Files | Folders | Examples |
|----------|-------|---------|----------|
| **Universal** | lowercase, no spaces | lowercase, no spaces | `config.yaml`, `linter-scripts/` |
| **PHP / WordPress** | `kebab-case.php` (classes: `class-*.php`) | `kebab-case/` | `class-admin-settings.php`, `includes/` |
| **Go** | `snake_case.go` | `lowercase` (no hyphens) | `http_handler.go`, `internal/` |
| **TypeScript** | `kebab-case.ts` (components: `PascalCase.tsx`) | `kebab-case/` | `UserCard.tsx`, `use-auth.ts` |
| **Rust** | `snake_case.rs` | `snake_case/` | `http_client.rs`, `error_handling/` |
| **C#** | `PascalCase.cs` | `PascalCase/` | `UserService.cs`, `Models/` |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Parent Overview | [../00-overview.md](../00-overview.md) |
| Cross-Language Guidelines | [../01-cross-language/00-overview.md](../01-cross-language/00-overview.md) |
| PHP Standards | [../04-php/00-overview.md](../04-php/00-overview.md) |
| Golang Standards | [../03-golang/00-overview.md](../03-golang/00-overview.md) |
| TypeScript Standards | [../02-typescript/00-overview.md](../02-typescript/00-overview.md) |

---

*Single source of truth for file & folder naming across all languages.*
