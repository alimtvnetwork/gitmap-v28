# File & Folder Naming — Rust / C# (AI Execution Prompt)

> **/goal** Master and enforce Rust (`snake_case.rs`, `snake_case/` modules) and C# (`PascalCase.cs`, `PascalCase/` namespaces) naming conventions.
> **/learn** Understand module file hierarchies (`mod.rs`), interface prefix standards (`IUserService.cs`), test project folders, and cross-platform case invariants.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce `snake_case.rs` and `snake_case/` module folders across all Rust crates.
- [ ] `/learn` Never use camelCase or kebab-case in Rust module filenames or directories.
- [ ] `/goal` Enforce `PascalCase.cs` matching class/interface names (`I*.cs`) in C# projects.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Rust

### File Naming — `snake_case.rs`

```
✅ http_client.rs
✅ error_handling.rs
✅ mod.rs
❌ httpClient.rs
❌ http-client.rs
❌ HttpClient.rs
```

### Folder Naming — `snake_case/`

```
✅ src/error_handling/
✅ src/http_client/
❌ src/error-handling/
❌ src/ErrorHandling/
```

### Standard Layout

```
my-crate/                        ← kebab-case crate name
├── src/
│   ├── main.rs                  ← binary entry
│   ├── lib.rs                   ← library entry
│   ├── error_handling/
│   │   ├── mod.rs               ← module root
│   │   └── custom_errors.rs
│   └── http_client/
│       ├── mod.rs
│       └── request_builder.rs
├── tests/
│   └── integration_test.rs
├── benches/
│   └── benchmark.rs
├── Cargo.toml
└── Cargo.lock
```

### Rust-Specific Rules

| Rule | Convention |
|------|-----------|
| Module files | `mod.rs` or `{module_name}.rs` |
| Test files | `tests/` folder or inline `#[cfg(test)]` |
| Crate name | `kebab-case` in Cargo.toml, `snake_case` in code |
| Macros | `snake_case!` |

---

## C#

### File Naming — `PascalCase.cs`

```
✅ UserService.cs
✅ HttpClientFactory.cs
✅ IUserRepository.cs            ← interfaces prefixed with I
❌ userService.cs
❌ user-service.cs
❌ user_service.cs
```

### Folder Naming — `PascalCase/`

C# is the **only language** that uses PascalCase folders:

```
✅ Models/
✅ Services/
✅ Controllers/
❌ models/
❌ services/
```

### Standard Layout

```
MyProject/                       ← PascalCase
├── MyProject.sln
├── src/
│   └── MyProject.Api/
│       ├── Controllers/
│       │   └── UserController.cs
│       ├── Models/
│       │   └── UserModel.cs
│       ├── Services/
│       │   ├── IUserService.cs
│       │   └── UserService.cs
│       └── Program.cs
└── tests/
    └── MyProject.Tests/
        └── UserServiceTests.cs
```

### C#-Specific Rules

| Rule | Convention |
|------|-----------|
| Interfaces | `I` prefix: `IUserService.cs` |
| Abstract classes | No special prefix |
| Test projects | `{Project}.Tests/` |
| One class per file | File name matches class name |

---

## Forbidden Patterns

| Language | Pattern | Why |
|----------|---------|-----|
| Rust | `camelCase.rs` | Violates Rust conventions |
| Rust | `kebab-case/` folders | Rust modules use `snake_case` |
| C# | `snake_case.cs` | Violates .NET conventions |
| C# | `lowercase/` folders | .NET uses PascalCase directories |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Rust Standards | [../05-rust/readme.md](../05-rust/readme.md) |
| C# Standards | [../07-csharp/readme.md](../07-csharp/readme.md) |
| Cross-Language Rules | [./02-cross-language.md](./02-cross-language.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/08-file-folder-naming/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-FILE-006: Rust and C# File and Directory Conventions

**Given** Repository file and directory structures across polyglot stacks.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase or invalid naming patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.
