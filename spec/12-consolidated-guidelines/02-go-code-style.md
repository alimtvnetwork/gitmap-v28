# 02 — Go Code Style

Go-specific coding conventions prioritizing readability and consistency.

## Positive Conditionals Only

Always write positive conditions. No `!`, no `!=`.

```go
// Correct
if len(args) > 0 { process(args) }

// Wrong
if len(args) != 0 { process(args) }
```

## Function Length — 8 to 15 Lines

Extract a helper when exceeding 15 lines. Each function does one thing.

## File Length — Max 200 Lines

| Signal | Action |
|--------|--------|
| 2+ unrelated function groups | Split into separate files |
| Large switch statement | Each case becomes a file |
| Mixed types and logic | Separate `model.go` from logic |

## Package Granularity — One Responsibility

| Package | Responsibility |
|---------|----------------|
| `cmd` | CLI routing and flag parsing |
| `config` | Config loading and merging |
| `constants` | All shared string literals |
| `model` | Shared data structures |
| `store` | Database access |

Rules: No circular imports. `cmd` orchestrates; others never import `cmd`. `model` and `constants` are leaf packages.

## No Magic Strings

Every literal used for comparison, defaults, or messages goes in `constants`.

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package | Lowercase, single word | `scanner` |
| Exported func | PascalCase, verb-led | `BuildRecords` |
| Unexported func | camelCase, verb-led | `parseFlags` |
| Constants | PascalCase | `DefaultBranch` |
| Files | Lowercase, single word | `terminal.go` |

## Constant Ownership — Organize by Package/Domain

Do **not** force artificial naming prefixes just to compensate for poor ownership.
Put constants in the package that owns the behavior:

| Owns the behavior | Put constants in | Examples |
|-------------------|------------------|----------|
| CLI routing / top-level commands | `gitmap-v27/constants/` when shared across dispatcher, completion, tests | `CmdScan`, `CmdReleaseAlias` |
| A single feature package | That feature package | template file extensions in `templates`, merge prompts in `movemerge` |
| A single cmd flow | `cmd` package/file-local constants | subcommand verbs, usage text, table headers |

Rules:

1. Use `gitmap-v27/constants/` only for literals shared across packages.
2. Keep package-private literals in the owning package instead of exporting them from `constants`.
3. Use clear PascalCase for exported constants and clear lowerCamelCase for unexported file/package-local constants.
4. Prefixes like `Cmd`, `Err`, `Flag`, and `Default` are allowed when they describe real semantics, but they are **not** mandatory for every new constant.

`.github/scripts/check-constants-naming.sh` enforces this at a lightweight level by rejecting newly added `gitmap-v27/constants/` names that are not ordinary PascalCase identifiers.

## Error Handling

Check errors immediately. Return errors up the stack. In `cmd` handlers: print and `os.Exit(1)`. Never `panic` for expected conditions.

---

Source: `spec/05-coding-guidelines/02-go-code-style.md`, `spec/03-general/06-code-style-rules.md`
