# Instruction (must follow): Coding Guideline Execution — Booleans, Naming & Enums

Trigger Keywords & Aliases: `cg-boolean`, `cg-execute boolean`, `fix boolean naming`, `execute naming guidelines`

```text
N = 100
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to scanning the codebase for non-compliant booleans, negative flags, generic garbage variable names, and unsuffixed enum types, writing the master execution plan to `.lovable/plans/pending/XX-boolean-naming-audit.md`, and decomposing it into subtasks in `.lovable/plans/subtasks/XX-boolean-naming/`.
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to executing those subtasks in an autonomous self-loop until all variables and enums follow strict naming conventions.
- [ ] /goal **Linter Mandate**: If an automated linter for boolean prefixes and generic variable names does not exist, you MUST create a custom naming linter script and connect it directly to the CI/CD local runner and workflows.
- [ ] /learn Ingest `.lovable/coding-guidelines/coding-guidelines.md`, `.lovable/strictly-avoid.md`, and `.lovable/memory/00-index.md` before touching any code.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## 1. Boolean & Naming Non-Negotiable Checklist

You MUST audit and strictly enforce every rule below across the entire codebase:

### A. Boolean Naming & Framing
1. **Mandatory Prefix**: Every boolean variable or struct field must begin with one of: `is`, `has`, `can`, `should`, `was`, `will`, `did`, `must`.
2. **Positive Framing Only**: Always use positive names (`isEnabled` yes, `isNotDisabled` no; `hasAccess` yes, `hasNoAccess` no). If the natural check is negative, invert it and flip the check site.
3. **No Explicit True/False Comparisons**: Never write `if isReady == true` or `if hasData == false`. Use `if isReady` or `if !hasData` (or invert to positive variable).
4. **No Boolean Flag Parameters**: Never pass raw boolean flags to functions (e.g. `render(true)` is banned). Split into two distinct functions (e.g. `renderExpanded()` and `renderCollapsed()`).
5. **No Tuples for Booleans**: Functions returning multiple values including a boolean must return a struct/wrapper object (`{ data, isSuccess }`), not raw `(T, bool)` tuples.

### B. Anti-Garbage Semantic Naming
1. **Zero Generic Garbage Variables**: Absolutely NO variable names like `temp`, `data`, `obj`, `val`, `item`, `comp_100`, `foo`, `bar`, `helper1`, `handler2`. Every identifier must clearly express its domain purpose.
2. **Semantic Unit Tests**: Test names must be behavior-driven and descriptive (e.g. `TestCreateUser_RejectsDuplicateEmail`), never generic (e.g. `TestHandler100` is an auto-reject).
3. **Acronyms in PascalCase**: Acronyms in identifiers must use PascalCase (`UserId`, `HttpServer`, `XmlParser`), not all-caps (`UserID`, `HTTPServer`).

### C. Enum Naming & Safety
1. **Mandatory `Type` Suffix**: Every enum and type definition representing a set of variants MUST end with the suffix `Type` (e.g. `UserRoleType`, `OrderStatusType`, `ErrorType`, `SeverityType`).
2. **No Magic Strings/Numbers**: Every status or state comparison must be against a registered enum constant.
3. **Dedicated Definition Files**: Enums and types must be defined in dedicated files (e.g. `types.go`, `UserRoleType.ts`), never inline alongside business logic.

---

## 2. Phase 1: Planning, Audit & Subtask Decomposition (Steps 1 .. N/2)

1. **Regex/AST Code Audit**: Scan the codebase for booleans lacking required prefixes, negative flags, generic variable names (`temp`, `data`), and enums missing the `Type` suffix.
2. **Master Plan**: Write a detailed execution plan to `.lovable/plans/pending/XX-boolean-naming-audit.md`.
3. **Subtask Files**: Decompose into subtask files in `.lovable/plans/subtasks/XX-boolean-naming/` (e.g. `01-task.md`, `02-task.md`, ...).
4. **Linter Connection**: Create or update the naming linter script and wire it into the CI/CD local runner (`.lovable/ai-fix-scripts/03-cicd-local-runner.py`).

---

## 3. Phase 2: Autonomous Execution Loop (Steps N/2+1 .. N)

1. **Refactor Identifiers**: Rename all non-compliant booleans, variables, and enums systematically across definitions and call sites.
2. **Global Blast Radius Verification**: Use ripgrep to ensure all references to renamed types and fields are updated across all packages.
3. **Run CI Gates**: Verify that `go vet ./...`, `go test ./...`, and CI/CD local runners pass cleanly.
4. **Update Status**: Mark completed tasks as `DONE`, move completed plans to `.lovable/plans/completed/`, and update `.lovable/plans/index.md`.

---

## 4. Pre-Commit Verification Checklist

- [ ] All booleans start with `is`, `has`, `can`, `should`, `was`, `will`, `did`, or `must`.
- [ ] No negative booleans or explicit `== true` comparisons exist.
- [ ] Zero generic garbage variable names (`temp`, `data`, `obj`, `comp_100`) remain.
- [ ] All enum names end with the `Type` suffix and live in dedicated files.
- [ ] Naming linter is integrated into CI/CD and passes with exit 0.
